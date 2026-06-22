package uta

import "testing"

func TestMarketCandle(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// GET /api/v3/market/candles
	{
		params := map[string]string{
			"category": string(CategoryUSDTFutures),
			"symbol":   "BTCUSDT",
			"interval": string(Granularity1m),
			"limit":    "2",
		}
		candles, err := c.NewGetCandlesService(CategoryUSDTFutures, "BTCUSDT", Granularity1m).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("candles: %v", err)
		}
		if len(candles) == 0 {
			t.Fatal("no candles returned")
		}
		first := candles[0]
		t.Logf("candle: ts=%s open=%s high=%s low=%s close=%s vol=%s turnover=%s",
			first.Ts, first.Open, first.High, first.Low, first.Close, first.Volume, first.Turnover)
		if first.Ts.IsZero() || first.Close.IsZero() {
			t.Fatal("candle ts/close not parsed")
		}
		raw := fetchRawGet(t, c, cx, "/api/v3/market/candles", params, false)
		assertCovers(t, "market/candles", raw, candles)
	}

	// GET /api/v3/market/history-candles
	{
		params := map[string]string{
			"category": string(CategorySpot),
			"symbol":   "BTCUSDT",
			"interval": string(Granularity1m),
			"limit":    "2",
		}
		candles, err := c.NewGetHistoryCandlesService(CategorySpot, "BTCUSDT", Granularity1m).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("history-candles: %v", err)
		}
		if len(candles) == 0 {
			t.Fatal("no history candles returned")
		}
		first := candles[0]
		t.Logf("history candle: ts=%s close=%s turnover=%s", first.Ts, first.Close, first.Turnover)
		if first.Ts.IsZero() || first.Close.IsZero() {
			t.Fatal("history candle ts/close not parsed")
		}
		raw := fetchRawGet(t, c, cx, "/api/v3/market/history-candles", params, false)
		assertCovers(t, "market/history-candles", raw, candles)
	}
}
