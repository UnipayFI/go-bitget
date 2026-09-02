package request

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/gorilla/websocket"
)

type testWSClient struct {
	url     string
	private bool
	dialer  *websocket.Dialer
}

func (c *testWSClient) GetPublicURL() string  { return c.url }
func (c *testWSClient) GetPrivateURL() string { return c.url }
func (c *testWSClient) GetAPIKey() string     { return "key" }
func (c *testWSClient) GetAPISecret() string  { return "secret" }
func (c *testWSClient) GetPassphrase() string { return "passphrase" }
func (c *testWSClient) GetSignFn() SignFn     { return nil }
func (c *testWSClient) GetLogger() log.Logger { return log.GetDefaultLogger() }
func (c *testWSClient) GetDialer() *websocket.Dialer {
	if c.dialer != nil {
		return c.dialer
	}
	return websocket.DefaultDialer
}

type blockedWriteConn struct {
	net.Conn
	blocked         atomic.Bool
	release         chan struct{}
	releaseOnce     sync.Once
	deadlineMu      sync.Mutex
	writeDeadline   time.Time
	deadlineChanged chan struct{}
}

type failWriteConn struct {
	net.Conn
	fail atomic.Bool
}

func (c *failWriteConn) Write(p []byte) (int, error) {
	if c.fail.Load() {
		return 0, net.ErrClosed
	}
	return c.Conn.Write(p)
}

func newBlockedWriteConn(conn net.Conn) *blockedWriteConn {
	return &blockedWriteConn{Conn: conn, release: make(chan struct{}), deadlineChanged: make(chan struct{}, 1)}
}

func (c *blockedWriteConn) Write(p []byte) (int, error) {
	if !c.blocked.Load() {
		return c.Conn.Write(p)
	}
	for {
		c.deadlineMu.Lock()
		deadline := c.writeDeadline
		c.deadlineMu.Unlock()
		if deadline.IsZero() {
			select {
			case <-c.release:
				return 0, net.ErrClosed
			case <-c.deadlineChanged:
			}
			continue
		}
		timer := time.NewTimer(time.Until(deadline))
		select {
		case <-timer.C:
			return 0, os.ErrDeadlineExceeded
		case <-c.release:
			timer.Stop()
			return 0, net.ErrClosed
		case <-c.deadlineChanged:
			timer.Stop()
		}
	}
}

func (c *blockedWriteConn) SetWriteDeadline(deadline time.Time) error {
	c.deadlineMu.Lock()
	c.writeDeadline = deadline
	c.deadlineMu.Unlock()
	select {
	case c.deadlineChanged <- struct{}{}:
	default:
	}
	return c.Conn.SetWriteDeadline(deadline)
}

func (c *blockedWriteConn) releaseWrite() {
	c.releaseOnce.Do(func() { close(c.release) })
}

func TestSubscribeRawArgsUsesOneConnection(t *testing.T) {
	for _, private := range []bool{false, true} {
		t.Run(map[bool]string{false: "public", true: "private"}[private], func(t *testing.T) {
			args := []WsArg{
				{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"},
				{InstType: "usdt-futures", Topic: "ticker", Symbol: "ETHUSDT"},
			}
			receivedSubscribe := make(chan []WsArg, 1)
			upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				conn, err := upgrader.Upgrade(w, r, nil)
				if err != nil {
					return
				}
				defer conn.Close()
				if private {
					_, msg, err := conn.ReadMessage()
					if err != nil {
						return
					}
					var login struct {
						Op string `json:"op"`
					}
					if common.JSONUnmarshal(msg, &login) != nil || login.Op != "login" {
						return
					}
					_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"event":"login","code":0}`))
				}
				_, msg, err := conn.ReadMessage()
				if err != nil {
					return
				}
				var sub struct {
					Op   string  `json:"op"`
					Args []WsArg `json:"args"`
				}
				if common.JSONUnmarshal(msg, &sub) != nil || sub.Op != "subscribe" {
					return
				}
				receivedSubscribe <- sub.Args
				for _, arg := range sub.Args {
					push := map[string]any{"action": "snapshot", "arg": arg, "data": []any{map[string]string{"symbol": arg.Symbol}}}
					body, _ := common.JSONMarshal(push)
					_ = conn.WriteMessage(websocket.TextMessage, body)
				}
				_, _, _ = conn.ReadMessage()
			}))
			defer server.Close()

			client := &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http"), private: private}
			pushes := make(chan []byte, len(args))
			done, stop, err := SubscribeRawArgs(context.Background(), client, private, args, func(msg []byte, err error) {
				if err == nil {
					pushes <- msg
				}
			})
			if err != nil {
				t.Fatal(err)
			}
			if got := <-receivedSubscribe; len(got) != len(args) {
				t.Fatalf("subscribe args = %d, want %d", len(got), len(args))
			}
			for range args {
				select {
				case <-pushes:
				case <-time.After(time.Second):
					t.Fatal("timed out waiting for push")
				}
			}
			close(done)
			select {
			case <-stop:
			case <-time.After(time.Second):
				t.Fatal("timed out waiting for stop")
			}
		})
	}
}

func TestRawSubscriptionUpdatesOneConnection(t *testing.T) {
	type operation struct {
		Op   string  `json:"op"`
		Args []WsArg `json:"args"`
	}
	operations := make(chan operation, 3)
	var connections atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		connections.Add(1)
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var op operation
			if common.JSONUnmarshal(message, &op) == nil && (op.Op == "subscribe" || op.Op == "unsubscribe") {
				operations <- op
			}
		}
	}))
	defer server.Close()

	client := &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	initial := WsArg{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}
	sub, err := OpenRawSubscription(context.Background(), client, false, []WsArg{initial}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()

	added := []WsArg{
		{InstType: "spot", Topic: "ticker", Symbol: "ETHUSDT"},
		{InstType: "usdt-futures", Topic: "ticker", Symbol: "ETHUSDT"},
	}
	if err := sub.Subscribe(added); err != nil {
		t.Fatal(err)
	}
	if err := sub.Unsubscribe([]WsArg{initial}); err != nil {
		t.Fatal(err)
	}

	want := []operation{{Op: "subscribe", Args: []WsArg{initial}}, {Op: "subscribe", Args: added}, {Op: "unsubscribe", Args: []WsArg{initial}}}
	for i := range want {
		select {
		case got := <-operations:
			if got.Op != want[i].Op || fmt.Sprint(got.Args) != fmt.Sprint(want[i].Args) {
				t.Fatalf("operation %d = %+v, want %+v", i, got, want[i])
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for operation %d", i)
		}
	}
	if got := connections.Load(); got != 1 {
		t.Fatalf("connections = %d, want 1", got)
	}
}

func TestRawSubscriptionConcurrentSubscribeFrames(t *testing.T) {
	const updates = 20
	received := make(chan struct{}, updates+1)
	var connections atomic.Int32
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		connections.Add(1)
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var op struct {
				Op   string  `json:"op"`
				Args []WsArg `json:"args"`
			}
			if common.JSONUnmarshal(message, &op) != nil || op.Op != "subscribe" || len(op.Args) == 0 {
				continue
			}
			received <- struct{}{}
		}
	}))
	defer server.Close()

	client := &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	sub, err := OpenRawSubscription(context.Background(), client, false, []WsArg{{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()

	var wg sync.WaitGroup
	for i := 0; i < updates; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if err := sub.Subscribe([]WsArg{{InstType: "spot", Topic: "ticker", Symbol: fmt.Sprintf("COIN%dUSDT", i)}}); err != nil {
				t.Errorf("subscribe %d: %v", i, err)
			}
		}(i)
	}
	wg.Wait()
	for i := 0; i < updates+1; i++ {
		select {
		case <-received:
		case <-time.After(time.Second):
			t.Fatalf("received %d/%d subscribe frames", i, updates+1)
		}
	}
	if got := connections.Load(); got != 1 {
		t.Fatalf("connections = %d, want 1", got)
	}
}

func TestSubscribeRawArgsCloseSurvivesBlockedWrite(t *testing.T) {
	subscribed := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		close(subscribed)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	dialer := *websocket.DefaultDialer
	var blocked *blockedWriteConn
	dialer.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err != nil {
			return nil, err
		}
		blocked = newBlockedWriteConn(conn)
		return blocked, nil
	}
	client := &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http"), dialer: &dialer}
	done, stop, err := SubscribeRawArgs(context.Background(), client, false, []WsArg{{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer blocked.releaseWrite()
	select {
	case <-subscribed:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscription")
	}
	blocked.blocked.Store(true)
	close(done)
	select {
	case <-stop:
	case <-time.After(2 * time.Second):
		t.Fatal("close blocked behind a stuck websocket write")
	}
}

func TestSubscribeRawArgsContextOnlyControlsSetup(t *testing.T) {
	subscribed := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		close(subscribed)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	done, stop, err := SubscribeRawArgs(ctx, &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http")}, false,
		[]WsArg{{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-subscribed:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscription")
	}
	cancel()
	select {
	case <-stop:
		t.Fatal("subscription stopped when setup context was cancelled")
	case <-time.After(100 * time.Millisecond):
	}
	close(done)
	select {
	case <-stop:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for stop")
	}
}

func TestRawSubscriptionDeduplicatesEachBatch(t *testing.T) {
	type operation struct {
		Op   string  `json:"op"`
		Args []WsArg `json:"args"`
	}
	operations := make(chan operation, 3)
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				return
			}
			var op operation
			if common.JSONUnmarshal(message, &op) == nil && (op.Op == "subscribe" || op.Op == "unsubscribe") {
				operations <- op
			}
		}
	}))
	defer server.Close()

	client := &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http")}
	initial := WsArg{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}
	sub, err := OpenRawSubscription(context.Background(), client, false, []WsArg{initial, initial}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	added := WsArg{InstType: "spot", Topic: "ticker", Symbol: "ETHUSDT"}
	if err := sub.Subscribe([]WsArg{added, added}); err != nil {
		t.Fatal(err)
	}
	if err := sub.Unsubscribe([]WsArg{added, added}); err != nil {
		t.Fatal(err)
	}
	for i, want := range []operation{{Op: "subscribe", Args: []WsArg{initial}}, {Op: "subscribe", Args: []WsArg{added}}, {Op: "unsubscribe", Args: []WsArg{added}}} {
		select {
		case got := <-operations:
			if got.Op != want.Op || fmt.Sprint(got.Args) != fmt.Sprint(want.Args) {
				t.Fatalf("operation %d = %+v, want %+v", i, got, want)
			}
		case <-time.After(time.Second):
			t.Fatalf("timed out waiting for operation %d", i)
		}
	}
}

func TestRawSubscriptionStopsOnServerDisconnect(t *testing.T) {
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		_, _, _ = conn.ReadMessage()
		_ = conn.Close()
	}))
	defer server.Close()

	done, stop, err := SubscribeRawArgs(context.Background(), &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http")}, false,
		[]WsArg{{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer close(done)
	select {
	case <-stop:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for stop after server disconnect")
	}
}

func TestRawSubscriptionStopsOnHeartbeatWriteFailure(t *testing.T) {
	subscribed := make(chan struct{})
	upgrader := websocket.Upgrader{CheckOrigin: func(*http.Request) bool { return true }}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
		close(subscribed)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}))
	defer server.Close()

	dialer := *websocket.DefaultDialer
	var failed *failWriteConn
	dialer.NetDialContext = func(ctx context.Context, network, address string) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, network, address)
		if err != nil {
			return nil, err
		}
		failed = &failWriteConn{Conn: conn}
		return failed, nil
	}
	sub, err := OpenRawSubscription(context.Background(), &testWSClient{url: "ws" + strings.TrimPrefix(server.URL, "http"), dialer: &dialer}, false,
		[]WsArg{{InstType: "spot", Topic: "ticker", Symbol: "BTCUSDT"}}, func([]byte, error) {})
	if err != nil {
		t.Fatal(err)
	}
	defer sub.Close()
	select {
	case <-subscribed:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for subscription")
	}
	failed.fail.Store(true)
	go keepAlive(sub.inner, time.Millisecond)
	select {
	case <-sub.Stop():
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for stop after heartbeat write failure")
	}
}
