package uta

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

// TestStrategyOrder places a futures "trigger" strategy order whose trigger
// price sits far below market (so it never fires) and then cancels it — a
// reversible exercise of the strategy write path. Gated behind
// BITGET_TEST_WRITE=1.
func TestStrategyOrder(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run live order tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	ref, err := c.NewPlaceStrategyOrderService(CategoryUSDTFutures, "BTCUSDT").
		SetType(StrategyTypeTrigger).
		SetTriggerBy(TriggerByMarket).
		SetTriggerPrice(decimal.RequireFromString("32000")). // far below market; never fires
		SetSide(SideBuy).
		SetQty(decimal.RequireFromString("0.001")).
		SetTriggerOrderType(OrderTypeMarket).
		Do(cx)
	if err != nil {
		t.Fatalf("place-strategy-order: %v", err)
	}
	t.Logf("placed strategy orderId=%s clientOid=%s", ref.OrderID, ref.ClientOid)

	// Confirm it appears in the unfilled strategy list.
	open, err := c.NewGetUnfilledStrategyOrdersService(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("unfilled-strategy-orders: %v", err)
	}
	t.Logf("unfilled strategy orders: %d", len(open))

	// Cancel it.
	if _, err := c.NewCancelStrategyOrderService().SetOrderId(ref.OrderID).Do(cx); err != nil {
		t.Fatalf("cancel-strategy-order: %v", err)
	}
	t.Logf("cancelled strategy order %s", ref.OrderID)
}
