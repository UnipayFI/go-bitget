package request

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/gorilla/websocket"
)

type testWSClient struct {
	url     string
	private bool
}

func (c *testWSClient) GetPublicURL() string         { return c.url }
func (c *testWSClient) GetPrivateURL() string        { return c.url }
func (c *testWSClient) GetAPIKey() string            { return "key" }
func (c *testWSClient) GetAPISecret() string         { return "secret" }
func (c *testWSClient) GetPassphrase() string        { return "passphrase" }
func (c *testWSClient) GetSignFn() SignFn            { return nil }
func (c *testWSClient) GetLogger() log.Logger        { return log.GetDefaultLogger() }
func (c *testWSClient) GetDialer() *websocket.Dialer { return websocket.DefaultDialer }

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
