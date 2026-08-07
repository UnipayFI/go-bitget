package uta

import "testing"

func TestCopy(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Trading pairs.
	pairs, err := c.NewGetCopyTradingPairsService().Do(cx)
	if err != nil {
		if tolerable(t, "copy/futures/trading-pairs", err) {
			t.Skip("API key lacks copy-trading permission")
		}
		t.Fatalf("trading pairs: %v", err)
	}
	t.Logf("copy trading pairs: %d", len(pairs))
	raw := fetchRawGet(t, c, cx, "/api/v3/copy/futures/trading-pairs", nil, true)
	assertCovers(t, "copy/futures/trading-pairs", raw, pairs)

	// Position summary.
	positions, err := c.NewGetCopyPositionSummaryService().Do(cx)
	if err != nil {
		t.Fatalf("position summary: %v", err)
	}
	t.Logf("copy positions: %d", len(positions))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/position-summary", nil, true)
	assertCovers(t, "copy/futures/position-summary", raw, positions)

	// Max transferable.
	maxT, err := c.NewGetCopyMaxTransferableService("USDT").Do(cx)
	if err != nil {
		t.Fatalf("max transferable: %v", err)
	}
	t.Logf("max transferable: %s available=%s", maxT.MaxTransferable, maxT.Available)
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/max-transferable",
		map[string]string{"coin": "USDT"}, true)
	assertCovers(t, "copy/futures/max-transferable", raw, maxT)

	// Transfer records.
	records, err := c.NewGetCopyTransferRecordService().SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("transfer record: %v", err)
	}
	t.Logf("copy transfer records: %d", len(records.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/transfer-record",
		map[string]string{"limit": "20"}, true)
	assertCovers(t, "copy/futures/transfer-record", raw, records)

	// Current followers.
	current, err := c.NewGetCopyCurrentFollowersService().SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("current followers: %v", err)
	}
	t.Logf("current followers: %d", len(current.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/current-follower",
		map[string]string{"limit": "20"}, true)
	assertCovers(t, "copy/futures/current-follower", raw, current)

	// History followers.
	history, err := c.NewGetCopyHistoryFollowersService().SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("history followers: %v", err)
	}
	t.Logf("history followers: %d", len(history.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/history-follower",
		map[string]string{"limit": "20"}, true)
	assertCovers(t, "copy/futures/history-follower", raw, history)

	// Profit summary.
	summary, err := c.NewGetCopyProfitSummaryService().Do(cx)
	if err != nil {
		t.Fatalf("profit summary: %v", err)
	}
	t.Logf("profit summary: total=%s allocated=%s pending=%s",
		summary.TotalProfit, summary.TotalAllocatedProfit, summary.TotalPendingProfit)
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/profit-summary", nil, true)
	assertCovers(t, "copy/futures/profit-summary", raw, summary)

	// Profit details.
	details, err := c.NewGetCopyProfitDetailsService().SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("profit details: %v", err)
	}
	t.Logf("profit details: %d", len(details.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/profit-details",
		map[string]string{"limit": "20"}, true)
	assertCovers(t, "copy/futures/profit-details", raw, details)
}
