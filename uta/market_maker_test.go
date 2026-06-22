package uta

import "testing"

func TestMarketMaker(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// Score weights (public, no required params).
	weights, err := c.NewGetScoreWeightsService().Do(cx)
	if err != nil {
		t.Fatalf("score-weights: %v", err)
	}
	if len(weights) == 0 {
		t.Fatal("no score weights returned")
	}
	t.Logf("score-weights: %d rows, first: %+v", len(weights), weights[0])
	raw := fetchRawGet(t, c, cx, "/api/v3/market/score-weights", nil, false)
	assertCovers(t, "market/score-weights", raw, weights)

	// Fee group (category required).
	feeParams := map[string]string{"category": "SPOT"}
	groups, err := c.NewGetFeeGroupService("SPOT").Do(cx)
	if err != nil {
		t.Fatalf("fee-group: %v", err)
	}
	if len(groups) == 0 {
		t.Fatal("no fee groups returned")
	}
	t.Logf("fee-group: %d groups, first: %+v", len(groups), groups[0])
	raw = fetchRawGet(t, c, cx, "/api/v3/market/fee-group", feeParams, false)
	assertCovers(t, "market/fee-group", raw, groups)

	// Cash dividend records (symbol + type required). MSFTUSDT/pending has rows.
	divParams := map[string]string{"symbol": "MSFTUSDT", "type": "pending"}
	div, err := c.NewGetCashDividendRecordsService("MSFTUSDT", "pending").Do(cx)
	if err != nil {
		t.Fatalf("cash-dividend-records: %v", err)
	}
	t.Logf("cash-dividend-records: %d rows", len(div.List))
	for _, r := range div.List {
		t.Logf("  exDividend=%s perShare=%s ts=%s", r.ExDividendDate, r.CashDividendPerShare, r.CashDividendTimestamp)
	}
	raw = fetchRawGet(t, c, cx, "/api/v3/market/cash-dividend-records", divParams, false)
	assertCovers(t, "market/cash-dividend-records", raw, div)
}
