package ws

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/UnipayFI/go-bitget/classic/spot"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/shopspring/decimal"
)

// authWSOptions builds an authenticated WebSocket client from the env creds,
// honoring BITGET_PROXY. Skips the test if creds are absent.
func authWSOptions(t *testing.T) []client.WebSocketOptions {
	t.Helper()
	k, s, p := apitest.Creds(t)
	opts := []client.WebSocketOptions{client.WithWebSocketAuth(k, s, p)}
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return opts
}

// TestWsPrivateSpotOrders validates the private spot "orders" channel end-to-end:
// it logs in, subscribes, then places a resting limit order via REST and waits
// for the channel to push it (proving login + private subscribe + typed decode),
// then cancels the order. Gated behind BITGET_TEST_WRITE=1.
func TestWsPrivateSpotOrders(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run the state-changing private WS test")
	}
	wsC := NewWebSocketClient(authWSOptions(t)...)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pushes := make(chan *WsPush[[]SpotWsOrder], 8)
	done, _, err := wsC.NewSubscribeSpotOrdersService("default").Do(ctx, func(p *WsPush[[]SpotWsOrder], e error) {
		// Do not touch t here: this callback can fire from the read goroutine
		// after the test returns (connection-closed on teardown).
		if e != nil || len(p.Data) == 0 {
			return
		}
		select {
		case pushes <- p:
		default:
		}
	})
	if err != nil {
		t.Fatalf("subscribe orders: %v", err)
	}
	defer close(done)
	time.Sleep(2 * time.Second) // let login + subscribe settle

	// Place a resting limit BUY via REST to trigger an orders push.
	spotC := spot.NewSpotClient(apitest.AuthOptions(t)...)
	if err := spotC.SyncServerTime(ctx); err != nil {
		t.Fatalf("sync: %v", err)
	}
	placed, err := spotC.NewPlaceOrderService("BTCUSDT", spot.SideBuy, spot.OrderTypeLimit, spot.ForceGTC, decimal.RequireFromString("0.0002")).
		SetPrice(decimal.RequireFromString("30000")).Do(ctx)
	if err != nil {
		t.Fatalf("place: %v", err)
	}
	t.Logf("placed resting order %s; awaiting orders-channel push", placed.OrderID)
	defer func() {
		_, _ = spotC.NewCancelOrderService("BTCUSDT").SetOrderID(placed.OrderID).Do(context.Background())
	}()

	select {
	case p := <-pushes:
		t.Logf("orders channel pushed %d order(s); first status=%s instId=%s — login + private subscribe + decode OK",
			len(p.Data), p.Data[0].Status, p.Data[0].InstId)
	case <-ctx.Done():
		t.Fatalf("timed out waiting for an orders-channel push")
	}
}

// TestWsTradeSpot validates WebSocket order entry (op:"trade"): it dials a
// logged-in trade connection, places a resting limit order over WS, then cancels
// it over WS. Gated behind BITGET_TEST_WRITE=1.
func TestWsTradeSpot(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run the state-changing WS-trade test")
	}
	wsC := NewWebSocketClient(authWSOptions(t)...)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	tc, err := wsC.DialTrade(ctx)
	if err != nil {
		t.Fatalf("dial trade: %v", err)
	}
	defer tc.Close()

	placed, err := tc.PlaceOrder(ctx, InstTypeSpot, "BTCUSDT", map[string]any{
		"orderType": "limit",
		"side":      "buy",
		"size":      "0.0002",
		"price":     "30000",
		"force":     "gtc",
	})
	if err != nil {
		t.Fatalf("ws place-order: %v", err)
	}
	if len(placed.Arg) == 0 || placed.Arg[0].Params.OrderID == "" {
		t.Fatalf("ws place-order: no orderId in reply: %+v", placed)
	}
	orderID := placed.Arg[0].Params.OrderID
	t.Logf("ws placed order %s", orderID)

	cancelled, err := tc.CancelOrder(ctx, InstTypeSpot, "BTCUSDT", map[string]any{"orderId": orderID})
	if err != nil {
		t.Fatalf("ws cancel-order: %v", err)
	}
	t.Logf("ws cancelled order; reply orderId=%s", func() string {
		if len(cancelled.Arg) > 0 {
			return cancelled.Arg[0].Params.OrderID
		}
		return ""
	}())
}
