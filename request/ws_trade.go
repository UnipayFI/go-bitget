package request

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/gorilla/websocket"
)

// WsTradeConn is a persistent, logged-in private connection for placing orders
// over WebSocket. Unlike the subscription channels it is request/response:
// each call sends an "op":"trade" frame tagged with a unique id and waits for
// the matching reply. One connection serves any number of concurrent calls.
type WsTradeConn struct {
	conn    *websocket.Conn
	logger  log.Logger
	mu      sync.Mutex
	pending map[string]chan []byte
	nextID  atomic.Int64
	done    chan struct{}
	closeMu sync.Once
}

type wsTradeOp struct {
	Op       string `json:"op"`
	ID       string `json:"id"`
	Category string `json:"category"`
	Topic    string `json:"topic"`
	Args     []any  `json:"args"`
}

// WsTradeResponse is the reply to a trade op. Args carries the topic-specific
// payload (e.g. the order acknowledgements).
type WsTradeResponse[T any] struct {
	Event  string    `json:"event"`
	ID     string    `json:"id"`
	Topic  string    `json:"topic"`
	Code   string    `json:"code"`
	Msg    string    `json:"msg"`
	Args   T         `json:"args"`
	ConnID string    `json:"connId"`
	Ts     time.Time `json:"ts"`
}

// DialWsTrade opens the private gateway, logs in, and returns a ready trade
// connection. Close it when done.
func DialWsTrade(ctx context.Context, client WsClient) (*WsTradeConn, error) {
	conn, _, err := client.GetDialer().DialContext(ctx, client.GetPrivateURL(), nil)
	if err != nil {
		return nil, err
	}
	conn.SetReadLimit(10 << 20)
	if err := wsLogin(client, conn); err != nil {
		conn.Close()
		return nil, err
	}
	t := &WsTradeConn{
		conn:    conn,
		logger:  client.GetLogger(),
		pending: make(map[string]chan []byte),
		done:    make(chan struct{}),
	}
	go t.readLoop()
	go t.keepAlive()
	return t, nil
}

func (t *WsTradeConn) readLoop() {
	for {
		_, msg, err := t.conn.ReadMessage()
		if err != nil {
			t.failAll(err)
			return
		}
		if common.BytesToString(msg) == "pong" {
			continue
		}
		t.logger.Debugf("ws trade recv: %s", common.BytesToString(msg))
		var hdr struct {
			ID string `json:"id"`
		}
		if err := common.JSONUnmarshal(msg, &hdr); err != nil || hdr.ID == "" {
			continue // non-response control frame
		}
		t.mu.Lock()
		ch := t.pending[hdr.ID]
		delete(t.pending, hdr.ID)
		t.mu.Unlock()
		if ch != nil {
			ch <- msg
		}
	}
}

func (t *WsTradeConn) keepAlive() {
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

func (t *WsTradeConn) failAll(err error) {
	t.mu.Lock()
	for id, ch := range t.pending {
		close(ch)
		delete(t.pending, id)
	}
	t.mu.Unlock()
}

// call sends a trade op and blocks for its reply (or ctx cancellation).
func (t *WsTradeConn) call(ctx context.Context, op, category, topic string, args []any) ([]byte, error) {
	id := strconv.FormatInt(t.nextID.Add(1), 10)
	ch := make(chan []byte, 1)
	t.mu.Lock()
	t.pending[id] = ch
	t.mu.Unlock()

	req := wsTradeOp{Op: op, ID: id, Category: category, Topic: topic, Args: args}
	data, err := common.JSONMarshal(req)
	if err != nil {
		return nil, err
	}
	if err := t.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.mu.Lock()
		delete(t.pending, id)
		t.mu.Unlock()
		return nil, err
	}

	select {
	case msg, ok := <-ch:
		if !ok {
			return nil, errors.New("ws trade: connection closed")
		}
		return msg, nil
	case <-ctx.Done():
		t.mu.Lock()
		delete(t.pending, id)
		t.mu.Unlock()
		return nil, ctx.Err()
	case <-t.done:
		return nil, errors.New("ws trade: connection closed")
	}
}

// WsTradeCall sends a trade op and decodes the reply, returning a *WsError on a
// non-success code.
func WsTradeCall[T any](ctx context.Context, t *WsTradeConn, op, category, topic string, args []any) (*WsTradeResponse[T], error) {
	msg, err := t.call(ctx, op, category, topic, args)
	if err != nil {
		return nil, err
	}
	var resp WsTradeResponse[T]
	if err := common.JSONUnmarshal(msg, &resp); err != nil {
		return nil, err
	}
	if resp.Code != "" && resp.Code != "0" && resp.Code != "00000" {
		return nil, &WsError{Code: resp.Code, Message: resp.Msg}
	}
	return &resp, nil
}

// Close terminates the connection and fails any in-flight calls.
func (t *WsTradeConn) Close() error {
	var err error
	t.closeMu.Do(func() {
		close(t.done)
		err = t.conn.Close()
	})
	return err
}
