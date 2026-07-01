package uta

import "testing"

func TestTradeQuery(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Open (unfilled) orders.
	open, err := c.NewGetOpenOrdersService().
		SetCategory(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("unfilled orders: %v", err)
	}
	t.Logf("unfilled orders: %d cursor=%s", len(open.List), open.Cursor)
	raw := fetchRawGet(t, c, cx, "/api/v3/trade/unfilled-orders",
		map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}, true)
	assertCovers(t, "trade/unfilled-orders", raw, open)

	// Order history.
	hist, err := c.NewGetOrderHistoryService(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("history orders: %v", err)
	}
	t.Logf("history orders: %d cursor=%s", len(hist.List), hist.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/trade/history-orders",
		map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}, true)
	assertCovers(t, "trade/history-orders", raw, hist)

	// Fill history.
	fills, err := c.NewGetFillHistoryService().SetCategory(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("fills: %v", err)
	}
	t.Logf("fills: %d cursor=%s", len(fills.List), fills.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/trade/fills",
		map[string]string{"category": string(CategoryUSDTFutures)}, true)
	assertCovers(t, "trade/fills", raw, fills)

	// Loan data.
	loan, err := c.NewGetLoanDataService().Do(cx)
	if err != nil {
		t.Fatalf("loan data: %v", err)
	}
	t.Logf("loan data: currentLoans=%s debtCoins=%d", loan.CurrentLoans, len(loan.DebtCoinList))
	raw = fetchRawGet(t, c, cx, "/api/v3/trade/loan-data", nil, true)
	assertCovers(t, "trade/loan-data", raw, loan)
}
