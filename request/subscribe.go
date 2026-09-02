package request

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"sync"
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

// RawSubscription owns one long-lived websocket connection. Subscribe and
// Unsubscribe update that connection in place; Stop is closed when its reader
// exits. Close is idempotent.
type RawSubscription struct {
	inner *rawSubscription
}

type rawSubscription struct {
	conn      *websocket.Conn
	cb        func([]byte, error)
	writeMu   sync.Mutex
	lastWrite time.Time
	opMu      sync.Mutex
	active    map[string]any
	closeC    chan struct{}
	stopC     chan struct{}
	closeOnce sync.Once
	stopOnce  sync.Once
	silent    atomic.Bool
}

// OpenRawSubscription opens one public or private connection and installs the
// initial channel set. Later updates reuse the same connection.
func OpenRawSubscription(ctx context.Context, client WsClient, private bool, args []WsArg, cb func(message []byte, err error)) (*RawSubscription, error) {
	if len(args) == 0 {
		return nil, errors.New("ws subscribe: no args")
	}
	values := make([]any, len(args))
	for i := range args {
		values[i] = args[i]
	}
	sub, err := openRawSubscription(ctx, client, private, values, cb)
	if err != nil {
		return nil, err
	}
	return &RawSubscription{inner: sub}, nil
}

// Subscribe adds channels to the existing connection. Already-active args are
// ignored, so callback-only updates do not consume a subscription request.
func (s *RawSubscription) Subscribe(args []WsArg) error {
	if s == nil || s.inner == nil {
		return errors.New("ws subscription is nil")
	}
	values := make([]any, len(args))
	for i := range args {
		values[i] = args[i]
	}
	return s.inner.update("subscribe", values)
}

// Unsubscribe removes channels from the existing connection. Unknown args are
// ignored.
func (s *RawSubscription) Unsubscribe(args []WsArg) error {
	if s == nil || s.inner == nil {
		return errors.New("ws subscription is nil")
	}
	values := make([]any, len(args))
	for i := range args {
		values[i] = args[i]
	}
	return s.inner.update("unsubscribe", values)
}

func (s *RawSubscription) Close() error {
	if s == nil || s.inner == nil {
		return nil
	}
	s.inner.close()
	return nil
}

func (s *RawSubscription) Stop() <-chan struct{} {
	if s == nil || s.inner == nil {
		closed := make(chan struct{})
		close(closed)
		return closed
	}
	return s.inner.stopC
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
	return subscribeBytes(ctx, client, private, []any{arg}, func(message []byte, e error) {
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
	return subscribeBytes(ctx, client, private, []any{arg}, cb)
}

// SubscribeRawArgs subscribes multiple v3 channels over one connection and
// delivers each data frame as raw bytes. Bitget accepts all args in a single
// subscribe operation; callers should still shard large channel sets.
func SubscribeRawArgs(ctx context.Context, client WsClient, private bool, args []WsArg, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	if len(args) == 0 {
		return nil, nil, errors.New("ws subscribe: no args")
	}
	values := make([]any, len(args))
	for i := range args {
		values[i] = args[i]
	}
	return subscribeBytes(ctx, client, private, values, cb)
}

// SubscribeRawArg is like SubscribeRaw but accepts an arbitrary subscription arg
// value (any JSON-serializable shape), not just the v3 WsArg. The classic v2
// streams use a different arg shape ({instType, channel, instId}), so they pass
// their own arg type here while reusing the shared connection/login/keepalive
// machinery.
func SubscribeRawArg(ctx context.Context, client WsClient, private bool, arg any, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return subscribeBytes(ctx, client, private, []any{arg}, cb)
}

func subscribeBytes(ctx context.Context, client WsClient, private bool, args []any, cb func(message []byte, err error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	sub, err := openRawSubscription(ctx, client, private, args, cb)
	if err != nil {
		return nil, nil, err
	}
	doneC := make(chan struct{})
	go func() {
		select {
		case <-doneC:
			sub.close()
		case <-sub.stopC:
		}
	}()
	return doneC, sub.stopC, nil
}

func openRawSubscription(ctx context.Context, client WsClient, private bool, args []any, cb func(message []byte, err error)) (*rawSubscription, error) {
	if len(args) == 0 {
		return nil, errors.New("ws subscribe: no args")
	}
	endpoint := client.GetPublicURL()
	if private {
		endpoint = client.GetPrivateURL()
	}
	conn, _, err := client.GetDialer().DialContext(ctx, endpoint, nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(10 << 20)

	if private {
		if err := wsLogin(client, conn); err != nil {
			conn.Close()
			return nil, err
		}
	}

	sub := wsOp{Op: "subscribe", Args: args}
	data, _ := common.JSONMarshal(sub)
	_ = conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		conn.Close()
		return nil, err
	}
	_ = conn.SetWriteDeadline(time.Time{})

	active := make(map[string]any, len(args))
	for _, arg := range args {
		key, err := wsArgKey(arg)
		if err != nil {
			conn.Close()
			return nil, err
		}
		active[key] = arg
	}
	s := &rawSubscription{conn: conn, cb: cb, active: active, closeC: make(chan struct{}), stopC: make(chan struct{}), lastWrite: time.Now()}

	go keepAlive(s, common.DEFAULT_KEEP_ALIVE_INTERVAL)
	go func() {
		select {
		case <-s.stopC:
		case <-s.closeC:
		case <-ctx.Done():
		}
		s.silent.Store(true)
		// Best-effort unsubscribe before closing.
		if s.opMu.TryLock() {
			remaining := make([]any, 0, len(s.active))
			for _, arg := range s.active {
				remaining = append(remaining, arg)
			}
			unsub := wsOp{Op: "unsubscribe", Args: remaining}
			if b, e := common.JSONMarshal(unsub); e == nil && len(remaining) > 0 && s.writeMu.TryLock() {
				_ = conn.SetWriteDeadline(time.Now().Add(time.Second))
				_ = conn.WriteMessage(websocket.TextMessage, b)
				s.writeMu.Unlock()
			}
			s.opMu.Unlock()
		}
		conn.Close()
	}()
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				if !s.silent.Load() {
					cb(nil, err)
				}
				s.stopOnce.Do(func() { close(s.stopC) })
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
				if !s.silent.Load() {
					cb(nil, &WsError{Code: hdr.codeString(), Message: hdr.Msg})
				}
				s.stopOnce.Do(func() { close(s.stopC) })
				return
			case hdr.Action != "":
				cb(message, nil)
			default:
				// subscribe/unsubscribe/login acks and other control frames.
			}
		}
	}()
	return s, nil
}

func (s *rawSubscription) update(op string, args []any) error {
	if len(args) == 0 {
		return nil
	}
	s.opMu.Lock()
	defer s.opMu.Unlock()
	select {
	case <-s.closeC:
		return errors.New("ws subscription closed")
	case <-s.stopC:
		return errors.New("ws subscription stopped")
	default:
	}
	filtered := make([]any, 0, len(args))
	keys := make([]string, 0, len(args))
	for _, arg := range args {
		key, err := wsArgKey(arg)
		if err != nil {
			return err
		}
		_, exists := s.active[key]
		if (op == "subscribe" && exists) || (op == "unsubscribe" && !exists) {
			continue
		}
		filtered = append(filtered, arg)
		keys = append(keys, key)
	}
	if len(filtered) == 0 {
		return nil
	}
	data, err := common.JSONMarshal(wsOp{Op: op, Args: filtered})
	if err != nil {
		return err
	}
	if err := s.writeText(data); err != nil {
		s.conn.Close()
		return err
	}
	for i, key := range keys {
		if op == "subscribe" {
			s.active[key] = filtered[i]
		} else {
			delete(s.active, key)
		}
	}
	return nil
}

func (s *rawSubscription) close() {
	s.silent.Store(true)
	s.closeOnce.Do(func() { close(s.closeC) })
}

func wsArgKey(arg any) (string, error) {
	data, err := common.JSONMarshal(arg)
	return common.BytesToString(data), err
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
func keepAlive(sub *rawSubscription, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := sub.writeText([]byte("ping")); err != nil {
			sub.conn.Close()
			return
		}
	}
}

func (s *rawSubscription) writeText(data []byte) error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	// Bitget permits 10 client messages/second/connection, including ping and
	// subscription updates. Serialize at a slightly lower rate so bursts queue.
	if wait := 110*time.Millisecond - time.Since(s.lastWrite); wait > 0 {
		timer := time.NewTimer(wait)
		select {
		case <-timer.C:
		case <-s.closeC:
			timer.Stop()
			return errors.New("ws subscription closed")
		}
	}
	_ = s.conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	err := s.conn.WriteMessage(websocket.TextMessage, data)
	s.lastWrite = time.Now()
	_ = s.conn.SetWriteDeadline(time.Time{})
	return err
}

// WsError is a Bitget WebSocket control-frame error.
type WsError struct {
	Code    string
	Message string
}

func (e *WsError) Error() string {
	return "<WsError> code=" + e.Code + ", msg=" + e.Message
}
