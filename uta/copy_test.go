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
}
