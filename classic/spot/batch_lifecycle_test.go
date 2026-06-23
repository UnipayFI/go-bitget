package spot

import (
	"os"
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

// TestSpotBatchLifecycle exercises the remaining state-changing spot trade
// endpoints — batch place, batch cancel, cancel-replace and cancel-symbol — and
// verifies their response structs. Gated behind BITGET_TEST_WRITE=1. Every
// order is a resting limit BUY at 30000 (well below the ~62k market), so nothing
// fills and nothing is spent; the orders are all cancelled before returning.
func TestSpotBatchLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run state-changing spot batch order tests")
	}
	c := NewSpotClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync server time: %v", err)
	}
	const symbol = "BTCUSDT"
	cx := apitest.Ctx(t)

	// 1. batch place two resting limit buys.
	items := []BatchOrderItem{
		{Symbol: symbol, Side: SideBuy, OrderType: OrderTypeLimit, Force: ForceGTC, Price: "30000", Size: "0.0002"},
		{Symbol: symbol, Side: SideBuy, OrderType: OrderTypeLimit, Force: ForceGTC, Price: "29000", Size: "0.0002"},
	}
	placed, err := c.NewBatchPlaceOrdersService(items).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("batch-orders: %v", err)
	}
	if len(placed.SuccessList) == 0 {
		t.Fatalf("batch-orders: no orders placed (failures: %+v)", placed.FailureList)
	}
	t.Logf("batch placed %d orders", len(placed.SuccessList))

	// 2. batch cancel the orders just placed.
	cancelItems := make([]BatchCancelOrderItem, 0, len(placed.SuccessList))
	for _, o := range placed.SuccessList {
		cancelItems = append(cancelItems, BatchCancelOrderItem{Symbol: symbol, OrderID: o.OrderID})
	}
	cancelled, err := c.NewBatchCancelOrderService(cancelItems).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("batch-cancel-order: %v", err)
	}
	t.Logf("batch cancelled %d orders", len(cancelled.SuccessList))

	// 3. cancel-replace: place one resting order, then cancel+replace it.
	repl, err := c.NewPlaceOrderService(symbol, SideBuy, OrderTypeLimit, ForceGTC, decimal.RequireFromString("0.0002")).
		SetPrice(decimal.RequireFromString("30000")).Do(cx)
	if err != nil {
		t.Fatalf("place for cancel-replace: %v", err)
	}
	cr, err := c.NewCancelReplaceOrderService(symbol, decimal.RequireFromString("29500"), decimal.RequireFromString("0.0002")).
		SetOrderID(repl.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("cancel-replace-order: %v", err)
	}
	t.Logf("cancel-replace success=%s newOrderId=%s", cr.Success, cr.OrderID)

	// 4. cancel-symbol-order: sweep any remaining open orders for the symbol.
	sweep, err := c.NewCancelSymbolOrderService(symbol).Do(cx)
	if err != nil {
		t.Fatalf("cancel-symbol-order: %v", err)
	}
	t.Logf("cancel-symbol-order swept symbol=%s", sweep.Symbol)
}
