package uta

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestOrderBatch exercises the batch + cancel-all write paths with tiny,
// far-from-market resting SPOT limit orders (fully reversible). Gated behind
// BITGET_TEST_WRITE=1. Covers place-batch, cancel-batch, cancel-symbol-order
// and countdown-cancel-all.
func TestOrderBatch(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run live order tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	const symbol = "BTCUSDT"
	base := time.Now().UnixNano()
	item := func(price string, i int) BatchOrderItem {
		return BatchOrderItem{
			Category:    CategorySpot,
			Symbol:      symbol,
			Qty:         decimal.RequireFromString("0.00004"),
			Price:       decimal.RequireFromString(price),
			Side:        SideBuy,
			OrderType:   OrderTypeLimit,
			TimeInForce: TimeInForcePostOnly,
			ClientOid:   "gobitget-b" + strconv.FormatInt(base+int64(i), 10),
		}
	}

	// 1) place-batch: two resting limit buys.
	placed, err := c.NewPlaceBatchService([]BatchOrderItem{item("32000", 1), item("31000", 2)}).Do(cx)
	if err != nil {
		t.Fatalf("place-batch: %v", err)
	}
	var ids []string
	for _, r := range placed {
		t.Logf("batch placed orderId=%s code=%s msg=%s", r.OrderID, r.Code, r.Msg)
		if r.OrderID != "" {
			ids = append(ids, r.OrderID)
		}
	}
	if len(ids) != 2 {
		t.Fatalf("expected 2 placed orders, got %d", len(ids))
	}

	// 2) cancel-batch: cancel both.
	cancelItems := make([]BatchCancelItem, 0, len(ids))
	for _, id := range ids {
		cancelItems = append(cancelItems, BatchCancelItem{OrderID: id, Category: CategorySpot, Symbol: symbol})
	}
	cancelled, err := c.NewCancelBatchService(cancelItems).Do(cx)
	if err != nil {
		t.Fatalf("cancel-batch: %v", err)
	}
	for _, r := range cancelled {
		t.Logf("batch cancelled orderId=%s code=%s msg=%s", r.OrderID, r.Code, r.Msg)
	}

	// 3) cancel-symbol-order: place one resting order then cancel all for the symbol.
	ref, err := c.NewPlaceOrderService(CategorySpot, symbol, decimal.RequireFromString("0.00004"), SideBuy, OrderTypeLimit).
		SetPrice(decimal.RequireFromString("30000")).
		SetTimeInForce(TimeInForcePostOnly).
		Do(cx)
	if err != nil {
		t.Fatalf("place for cancel-symbol: %v", err)
	}
	t.Logf("placed %s for cancel-symbol-order", ref.OrderId)
	res, err := c.NewCancelSymbolOrderService(CategorySpot).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("cancel-symbol-order: %v", err)
	}
	t.Logf("cancel-symbol-order cancelled %d order(s)", len(res.List))
	for _, r := range res.List {
		t.Logf("  cancelled orderId=%s code=%s", r.OrderId, r.Code)
	}

	// 4) countdown-cancel-all("0"): disable the dead-man switch (no-op, safe).
	cd, err := c.NewCountDownCancelAllService("0").Do(cx)
	if err != nil {
		t.Fatalf("countdown-cancel-all: %v", err)
	}
	t.Logf("countdown-cancel-all -> %q", *cd)
}
