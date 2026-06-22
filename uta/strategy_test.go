package uta

import "testing"

func TestStrategy(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Unfilled strategy orders.
	unfilled, err := c.NewGetUnfilledStrategyOrdersService(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("unfilled strategy orders: %v", err)
	}
	t.Logf("unfilled strategy orders: %d", len(unfilled))
	raw := fetchRawGet(t, c, cx, "/api/v3/trade/unfilled-strategy-orders",
		map[string]string{"category": string(CategoryUSDTFutures)}, true)
	assertCovers(t, "trade/unfilled-strategy-orders", raw, unfilled)

	// History strategy orders.
	history, err := c.NewGetHistoryStrategyOrdersService(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("history strategy orders: %v", err)
	}
	t.Logf("history strategy orders: %d cursor=%s", len(history.List), history.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/trade/history-strategy-orders",
		map[string]string{"category": string(CategoryUSDTFutures)}, true)
	assertCovers(t, "trade/history-strategy-orders", raw, history)
}
