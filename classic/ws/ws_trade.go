package ws

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
	"github.com/UnipayFI/go-bitget/request"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/gorilla/websocket"
)

// TradeConn is a persistent, logged-in private v2 connection for placing and
// cancelling orders over WebSocket. Unlike the subscription channels it is
// request/response: each call sends an op:"trade" frame whose arg carries a
// unique id, and waits for the reply echoing that id. One connection serves any
// number of concurrent calls. The classic v2 frame shape differs from the UTA
// v3 trade op (here the id lives INSIDE the arg, and channel is "place-order" /
// "cancel-order"), so it has its own connection type.
type TradeConn struct {
	conn      *websocket.Conn
	logger    log.Logger
	mu        sync.Mutex
	pending   map[string]chan []byte
	nextID    atomic.Int64
	done      chan struct{}
	closeOnce sync.Once
}

// DialTrade opens the private v2 gateway, logs in, and returns a ready trade
// connection. Close it when done.
func (c *WebSocketClient) DialTrade(ctx context.Context) (*TradeConn, error) {
	conn, err := request.DialPrivateLoggedIn(ctx, c)
	if err != nil {
		return nil, err
	}
	t := &TradeConn{
		conn:    conn,
		logger:  c.GetLogger(),
		pending: make(map[string]chan []byte),
		done:    make(chan struct{}),
	}
	go t.readLoop()
	go t.keepAlive()
	return t, nil
}

type tradeArg struct {
	ID       string         `json:"id"`
	InstType string         `json:"instType"`
	InstId   string         `json:"instId"`
	Channel  string         `json:"channel"`
	Params   map[string]any `json:"params"`
}

type tradeFrame struct {
	Op   string     `json:"op"`
	Args []tradeArg `json:"args"`
}

// TradeResponse is the reply to an op:"trade" frame. Arg carries the per-call
// echo (id + the order acknowledgement params).
type TradeResponse struct {
	Event string             `json:"event"`
	Arg   []TradeResponseArg `json:"arg"`
	Code  jsontext.Value     `json:"code"` // number (0) or string
	Msg   string             `json:"msg"`
	Ts    time.Time          `json:"ts"`
}

// TradeResponseArg is one echoed order acknowledgement.
type TradeResponseArg struct {
	ID       string           `json:"id"`
	InstType string           `json:"instType"`
	Channel  string           `json:"channel"`
	InstId   string           `json:"instId"`
	Params   TradeResultParam `json:"params"`
}

// TradeResultParam is the order ack payload.
type TradeResultParam struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

func (t *TradeConn) readLoop() {
	for {
		_, msg, err := t.conn.ReadMessage()
		if err != nil {
			t.failAll()
			return
		}
		if common.BytesToString(msg) == "pong" {
			continue
		}
		t.logger.Debugf("ws trade recv: %s", common.BytesToString(msg))
		// The reply id lives inside arg[0].id.
		var hdr struct {
			Arg []struct {
				ID string `json:"id"`
			} `json:"arg"`
		}
		if err := common.JSONUnmarshal(msg, &hdr); err != nil || len(hdr.Arg) == 0 || hdr.Arg[0].ID == "" {
			continue // login ack / other control frame
		}
		id := hdr.Arg[0].ID
		t.mu.Lock()
		ch := t.pending[id]
		delete(t.pending, id)
		t.mu.Unlock()
		if ch != nil {
			ch <- msg
		}
	}
}

func (t *TradeConn) keepAlive() {
	ticker := time.NewTicker(common.DEFAULT_KEEP_ALIVE_INTERVAL)
	defer ticker.Stop()
	for {
		select {
		case <-t.done:
			return
		case <-ticker.C:
			if err := t.conn.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				return
			}
		}
	}
}

func (t *TradeConn) failAll() {
	t.mu.Lock()
	for id, ch := range t.pending {
		close(ch)
		delete(t.pending, id)
	}
	t.mu.Unlock()
}

// Trade sends an op:"trade" frame for the given channel ("place-order" or
// "cancel-order") with the supplied params and blocks for the reply. It is the
// generic primitive behind PlaceOrder/CancelOrder.
func (t *TradeConn) Trade(ctx context.Context, instType InstType, instId, channel string, params map[string]any) (*TradeResponse, error) {
	id := strconv.FormatInt(t.nextID.Add(1), 10)
	ch := make(chan []byte, 1)
	t.mu.Lock()
	t.pending[id] = ch
	t.mu.Unlock()

	frame := tradeFrame{Op: "trade", Args: []tradeArg{{
		ID:       id,
		InstType: string(instType),
		InstId:   instId,
		Channel:  channel,
		Params:   params,
	}}}
	data, err := common.JSONMarshal(frame)
	if err != nil {
		t.clearPending(id)
		return nil, err
	}
	if err := t.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.clearPending(id)
		return nil, err
	}

	select {
	case msg, ok := <-ch:
		if !ok {
			return nil, errors.New("ws trade: connection closed")
		}
		var resp TradeResponse
		if err := common.JSONUnmarshal(msg, &resp); err != nil {
			return nil, err
		}
		if code := normalizeCode(resp.Code); code != "" && code != "0" && code != "00000" {
			return nil, &request.WsError{Code: code, Message: resp.Msg}
		}
		return &resp, nil
	case <-ctx.Done():
		t.clearPending(id)
		return nil, ctx.Err()
	case <-t.done:
		return nil, errors.New("ws trade: connection closed")
	}
}

// PlaceOrder places an order over the WebSocket. params holds the order fields
// (orderType, side, size, price, force, clientOid, ...) per the product's
// place-order channel doc.
func (t *TradeConn) PlaceOrder(ctx context.Context, instType InstType, instId string, params map[string]any) (*TradeResponse, error) {
	return t.Trade(ctx, instType, instId, "place-order", params)
}

// CancelOrder cancels an order over the WebSocket. params holds orderId or
// clientOid (and any product-specific fields).
func (t *TradeConn) CancelOrder(ctx context.Context, instType InstType, instId string, params map[string]any) (*TradeResponse, error) {
	return t.Trade(ctx, instType, instId, "cancel-order", params)
}

func (t *TradeConn) clearPending(id string) {
	t.mu.Lock()
	delete(t.pending, id)
	t.mu.Unlock()
}

// Close terminates the connection and fails any in-flight calls.
func (t *TradeConn) Close() error {
	var err error
	t.closeOnce.Do(func() {
		close(t.done)
		err = t.conn.Close()
	})
	return err
}

// normalizeCode renders the response code (a JSON number or quoted string) as a
// plain string.
func normalizeCode(c jsontext.Value) string {
	return strings.Trim(string(c), `"`)
}
