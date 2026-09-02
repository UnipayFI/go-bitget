package request

import (
	"context"
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
