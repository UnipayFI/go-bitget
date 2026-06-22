package uta

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestWsTradeLifecycle places a tiny, far-from-market resting SPOT limit order
// over the WebSocket trade connection and cancels it over the same connection.
// Fully reversible. Gated behind BITGET_TEST_WRITE=1.
func TestWsTradeLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run live ws trade tests")
	}
	ws := testWsClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	tc, err := ws.DialTrade(ctx)
	if err != nil {
		t.Fatalf("dial trade: %v", err)
	}
	defer tc.Close()

	price := decimal.RequireFromString("30000")
	clientOid := "gobitget-ws" + strconv.FormatInt(time.Now().UnixNano(), 10)
	ack, err := tc.PlaceOrder(ctx, CategorySpot, WsNewOrder{
		Symbol:      "BTCUSDT",
		Side:        SideBuy,
		OrderType:   OrderTypeLimit,
		Qty:         decimal.RequireFromString("0.00004"),
		Price:       &price,
		TimeInForce: TimeInForcePostOnly,
		ClientOid:   clientOid,
	})
	if err != nil {
		t.Fatalf("ws place-order: %v", err)
	}
	t.Logf("ws placed orderId=%s clientOid=%s cTime=%s", ack.OrderID, ack.ClientOid, ack.CTime)
	if ack.OrderID == "" {
		t.Fatal("ws place-order returned empty orderId")
	}

	cancelAck, err := tc.CancelOrder(ctx, CategorySpot, WsCancelOrder{
		Symbol:  "BTCUSDT",
		OrderID: ack.OrderID,
	})
	if err != nil {
		t.Fatalf("ws cancel-order: %v", err)
	}
	t.Logf("ws cancelled orderId=%s", cancelAck.OrderID)
}
