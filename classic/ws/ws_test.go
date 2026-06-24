package ws

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/client"
)

// publicWSOptions honors BITGET_PROXY for egress; public WS has no IP allowlist.
func publicWSOptions() []client.WebSocketOptions {
	var opts []client.WebSocketOptions
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return opts
}

// TestWsSpotTickerSmoke validates the classic v2 WebSocket plumbing end-to-end:
// it connects to the public v2 gateway, subscribes to the spot ticker channel
// for BTCUSDT, and waits for one data push. This proves the v2 URL, subscribe-op
// format and {instType, channel, instId} arg shape are correct.
func TestWsSpotTickerSmoke(t *testing.T) {
	c := NewWebSocketClient(publicWSOptions()...)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	got := make(chan []byte, 1)
	done, _, err := c.SubscribeRaw(ctx, false, WsArg{InstType: string(InstTypeSpot), Channel: "ticker", InstID: "BTCUSDT"}, func(msg []byte, e error) {
		if e != nil {
			return
		}
		select {
		case got <- msg:
		default:
		}
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}
	defer close(done)

	select {
	case msg := <-got:
		s := string(msg)
		if len(s) > 200 {
			s = s[:200] + "..."
		}
		t.Logf("received spot ticker push: %s", s)
	case <-ctx.Done():
		t.Fatalf("timed out waiting for a spot ticker push")
	}
}
