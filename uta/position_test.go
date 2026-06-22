package uta

import "testing"

func TestPosition(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Current positions.
	positions, err := c.NewGetPositionService(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("current position: %v", err)
	}
	t.Logf("current positions: %d", len(positions))
	// current-position returns {"list":[...]}; Position field coverage is
	// validated in the order-lifecycle test once a real position exists.

	// Positions history.
	history, err := c.NewGetPositionHistoryService(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("history position: %v", err)
	}
	t.Logf("history positions: %d cursor=%s", len(history.List), history.Cursor)
	raw := fetchRawGet(t, c, cx, "/api/v3/position/history-position",
		map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}, true)
	assertCovers(t, "position/history-position", raw, history)

	// ADL rank.
	adl, err := c.NewGetPositionADLRankService().Do(cx)
	if err != nil {
		t.Fatalf("adl rank: %v", err)
	}
	t.Logf("adl rank: %d", len(adl))
	raw = fetchRawGet(t, c, cx, "/api/v3/position/adlRank", nil, true)
	assertCovers(t, "position/adlRank", raw, adl)
}
