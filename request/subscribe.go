package request

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/gorilla/websocket"
)

// WsClient is what the subscribe framework needs from a *client.WebSocketClient.
type WsClient interface {
	GetPublicURL() string
	GetPrivateURL() string
	GetAPIKey() string
	GetAPISecret() string
	GetPassphrase() string
	GetSignFn() SignFn
	GetLogger() log.Logger
	GetDialer() *websocket.Dialer
}

// WsArg identifies a channel subscription. instType routes by product
// ("spot", "usdt-futures", "coin-futures", "usdc-futures", or "UTA" for the
// account channel); topic is the channel name; symbol/coin narrow it.
type WsArg struct {
	InstType string `json:"instType"`
	Topic    string `json:"topic"`
	Symbol   string `json:"symbol,omitempty"`
	Coin     string `json:"coin,omitempty"`
	Interval string `json:"interval,omitempty"` // candlestick (kline) channel only
}

// WsAction classifies a data push as a full snapshot or an incremental update.
type WsAction string

const (
	WsActionSnapshot WsAction = "snapshot" // full state
	WsActionUpdate   WsAction = "update"   // incremental change
)

// WsPush is the envelope Bitget pushes for a data event.
type WsPush[T any] struct {
	Action WsAction  `json:"action"`
	Arg    WsArg     `json:"arg"`
	Data   T         `json:"data"`
	Ts     time.Time `json:"ts"`
}

type wsOp struct {
	Op   string `json:"op"`
	Args []any  `json:"args"`
}

type wsLoginOp struct {
	Op   string       `json:"op"`
	Args []wsLoginArg `json:"args"`
}

type wsLoginArg struct {
	APIKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

// wsHeader is a lightweight view used to classify an inbound frame before
// committing to a typed decode. Code is read raw because the stream encodes it
// as a JSON number (0, 30001) whereas REST uses a quoted string.
type wsHeader struct {
	Event  string         `json:"event"`
	Action string         `json:"action"`
	Code   jsontext.Value `json:"code"`
	Msg    string         `json:"msg"`
}

// codeString normalizes the raw code token to a string ("" when absent).
func (h wsHeader) codeString() string {
	return strings.Trim(string(h.Code), `"`)
}

// ok reports whether the header's code is success/absent.
func (h wsHeader) ok() bool {
	c := h.codeString()
	return c == "" || c == "0" || c == "00000"
}

// Subscribe opens a dedicated connection to the public or private gateway, logs
// in when private, subscribes to arg, and invokes cb for every data push.
// Returns a done channel (close to stop) and a stop channel (closed when the
// reader exits). The typed Data field of the push is decoded into *T.
func Subscribe[T any](ctx context.Context, client WsClient, private bool, arg WsArg, cb func(*WsPush[T], error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, private, arg, func(message []byte, e error) {
		if e != nil {
			cb(nil, e)
			return
		}
		var push WsPush[T]
		if err := common.JSONUnmarshal(message, &push); err != nil {
			cb(nil, err)
			return
		}
		cb(&push, nil)
	})
}

// SubscribeRaw is like Subscribe but delivers each data frame's raw bytes,
// for channels whose payload shape the caller wants to decode itself.
func SubscribeRaw(ctx context.Context, client WsClient, private bool, arg WsArg, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, private, arg, cb)
}

// SubscribeRawArg is like SubscribeRaw but accepts an arbitrary subscription arg
// value (any JSON-serializable shape), not just the v3 WsArg. The classic v2
// streams use a different arg shape ({instType, channel, instId}), so they pass
// their own arg type here while reusing the shared connection/login/keepalive
// machinery.
func SubscribeRawArg(ctx context.Context, client WsClient, private bool, arg any, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, private, arg, cb)
}

func subscribeBytes(ctx context.Context, client WsClient, private bool, arg any, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	endpoint := client.GetPublicURL()
	if private {
		endpoint = client.GetPrivateURL()
	}
	conn, _, err := client.GetDialer().DialContext(ctx, endpoint, nil)
	if err != nil {
		return nil, nil, err
	}
	conn.SetReadLimit(10 << 20)

	if private {
		if err := wsLogin(client, conn); err != nil {
			conn.Close()
			return nil, nil, err
		}
	}

	sub := wsOp{Op: "subscribe", Args: []any{arg}}
	data, _ := common.JSONMarshal(sub)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		conn.Close()
		return nil, nil, err
	}

	doneC := make(chan struct{})
	stopC := make(chan struct{})
	// silent suppresses callback delivery on the deliberate-shutdown path: once the caller
	// closes doneC to stop, the ReadMessage error from the watcher closing the conn must not
	// reach the callback. atomic because the watcher and reader run on different goroutines.
	var silent atomic.Bool

	go keepAlive(conn, common.DEFAULT_KEEP_ALIVE_INTERVAL)
	go func() {
		select {
		case <-stopC:
			silent.Store(true)
		case <-doneC:
			silent.Store(true)
		}
		// Best-effort unsubscribe before closing.
		unsub := wsOp{Op: "unsubscribe", Args: []any{arg}}
		if b, e := common.JSONMarshal(unsub); e == nil {
			_ = conn.WriteMessage(websocket.TextMessage, b)
		}
		conn.Close()
	}()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !silent.Load() {
					cb(nil, err)
				}
				close(stopC)
				return
			}
			if common.BytesToString(message) == "pong" {
				continue
			}
			client.GetLogger().Debugf("ws recv: %s", common.BytesToString(message))

			var hdr wsHeader
			if err := common.JSONUnmarshal(message, &hdr); err != nil {
				cb(nil, err)
				continue
			}
			switch {
			case hdr.Event == "error":
				// A server error control frame (subscription rejected 30001/30016, relogin
				// failure, etc.) is protocol-level fatal: under the one-connection-per-subscription
				// model this connection is permanently dead, and if the reader kept looping the
				// stream would be "fake-alive" (connected, never delivering data, never reconnecting).
				// Deliver the error, then close(stopC)+return to terminate the reader, matching the
				// transport-level ReadMessage error path so the caller drives reconnect/resubscribe
				// off the stop-close (the documented done/stop contract).
				if !silent.Load() {
					cb(nil, &WsError{Code: hdr.codeString(), Message: hdr.Msg})
				}
				close(stopC)
				return
			case hdr.Action != "":
				cb(message, nil)
			default:
				// subscribe/unsubscribe/login acks and other control frames.
			}
		}
	}()
	return doneC, stopC, nil
}

// DialPrivateLoggedIn dials the private WebSocket gateway and completes the
// login handshake, returning a ready connection. WebSocket order-entry
// (op:"trade") connections build on this. The caller owns and must Close the
// returned connection.
func DialPrivateLoggedIn(ctx context.Context, client WsClient) (*websocket.Conn, error) {
	conn, _, err := client.GetDialer().DialContext(ctx, client.GetPrivateURL(), nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(10 << 20)
	if err := wsLogin(client, conn); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// wsLogin performs the private-stream login handshake and blocks until the
// server acknowledges (or rejects) it.
func wsLogin(client WsClient, conn *websocket.Conn) error {
	apiKey := client.GetAPIKey()
	secret := client.GetAPISecret()
	passphrase := client.GetPassphrase()
	if apiKey == "" || secret == "" || passphrase == "" {
		return errors.New("ws login: missing credentials (WithWebSocketAuth)")
	}
	ts := strconv.FormatInt(time.Now().Unix(), 10)
	prehash := ts + "GET" + "/user/verify"
	var (
		sign string
		err  error
	)
	if fn := client.GetSignFn(); fn != nil {
		sign, err = fn(secret, prehash)
	} else {
		sign, err = HMACSign(secret, prehash)
	}
	if err != nil {
		return err
	}

	login := wsLoginOp{Op: "login", Args: []wsLoginArg{{
		APIKey:     apiKey,
		Passphrase: passphrase,
		Timestamp:  ts,
		Sign:       sign,
	}}}
	data, _ := common.JSONMarshal(login)
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return err
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	defer conn.SetReadDeadline(time.Time{})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		if common.BytesToString(message) == "pong" {
			continue
		}
		var hdr wsHeader
		if err := common.JSONUnmarshal(message, &hdr); err != nil {
			return err
		}
		switch hdr.Event {
		case "login":
			if !hdr.ok() {
				return &WsError{Code: hdr.codeString(), Message: hdr.Msg}
			}
			return nil
		case "error":
			return &WsError{Code: hdr.codeString(), Message: hdr.Msg}
		}
	}
}

// keepAlive sends Bitget's literal "ping" text frame on an interval; the server
// replies "pong" (handled in the read loop).
func keepAlive(conn *websocket.Conn, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
			return
		}
	}
}

// WsError is a Bitget WebSocket control-frame error.
type WsError struct {
	Code    string
	Message string
}

func (e *WsError) Error() string {
	return "<WsError> code=" + e.Code + ", msg=" + e.Message
}
