package uta

import "testing"

func TestMarketFunding(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	cat := string(CategoryUSDTFutures)
	symbol := "BTCUSDT"

	// Current funding rate.
	curr, err := c.NewGetCurrentFundingRateService(CategoryUSDTFutures).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("current funding rate: %v", err)
	}
	if len(curr) == 0 {
		t.Fatal("no current funding rate returned")
	}
	t.Logf("current: %+v", curr[0])
	raw := fetchRawGet(t, c, cx, "/api/v3/market/current-fund-rate", map[string]string{"category": cat, "symbol": symbol}, false)
	assertCovers(t, "market/current-fund-rate", raw, curr)

	// Funding rate history.
	hist, err := c.NewGetFundingRateHistoryService(CategoryUSDTFutures, symbol).Do(cx)
	if err != nil {
		t.Fatalf("funding rate history: %v", err)
	}
	if len(hist.ResultList) == 0 {
		t.Fatal("no funding rate history returned")
	}
	t.Logf("history first: %+v", hist.ResultList[0])
	raw = fetchRawGet(t, c, cx, "/api/v3/market/history-fund-rate", map[string]string{"category": cat, "symbol": symbol}, false)
	assertCovers(t, "market/history-fund-rate", raw, hist)

	// Open interest.
	oi, err := c.NewGetOpenInterestService(CategoryUSDTFutures).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("open interest: %v", err)
	}
	if len(oi.List) == 0 {
		t.Fatal("no open interest returned")
	}
	t.Logf("open interest: %+v ts=%s", oi.List[0], oi.Ts)
	raw = fetchRawGet(t, c, cx, "/api/v3/market/open-interest", map[string]string{"category": cat, "symbol": symbol}, false)
	assertCovers(t, "market/open-interest", raw, oi)

	// Open interest limit.
	oiLimit, err := c.NewGetOpenInterestLimitService(CategoryUSDTFutures).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("open interest limit: %v", err)
	}
	if len(oiLimit) == 0 {
		t.Fatal("no open interest limit returned")
	}
	t.Logf("oi limit: %+v", oiLimit[0])
	raw = fetchRawGet(t, c, cx, "/api/v3/market/oi-limit", map[string]string{"category": cat, "symbol": symbol}, false)
	assertCovers(t, "market/oi-limit", raw, oiLimit)
}
