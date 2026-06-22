package uta

import "testing"

func TestMarketTicker(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// Tickers -- use USDT-FUTURES so the futures-only fields are covered too.
	{
		params := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}
		list, err := c.NewGetTickersService(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
		if err != nil {
			t.Fatalf("tickers: %v", err)
		}
		if len(list) == 0 {
			t.Fatal("no tickers returned")
		}
		t.Logf("ticker: %+v", list[0])
		raw := fetchRawGet(t, c, cx, "/api/v3/market/tickers", params, false)
		assertCovers(t, "market/tickers", raw, list)
	}

	// OrderBook.
	{
		params := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT", "limit": "5"}
		ob, err := c.NewGetOrderBookService(CategoryUSDTFutures, "BTCUSDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("orderbook: %v", err)
		}
		if len(ob.Asks) == 0 || len(ob.Bids) == 0 {
			t.Fatal("empty orderbook")
		}
		t.Logf("orderbook: asks=%d bids=%d ts=%s", len(ob.Asks), len(ob.Bids), ob.Ts)
		raw := fetchRawGet(t, c, cx, "/api/v3/market/orderbook", params, false)
		assertCovers(t, "market/orderbook", raw, ob)
	}

	// RPI OrderBook.
	{
		params := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT", "limit": "5"}
		ob, err := c.NewGetRPIOrderBookService(CategoryUSDTFutures, "BTCUSDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("rpi-orderbook: %v", err)
		}
		t.Logf("rpi-orderbook: asks=%d bids=%d ts=%s", len(ob.Asks), len(ob.Bids), ob.Ts)
		raw := fetchRawGet(t, c, cx, "/api/v3/market/rpi-orderbook", params, false)
		assertCovers(t, "market/rpi-orderbook", raw, ob)
	}

	// Recent public fills.
	{
		params := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT", "limit": "5"}
		fills, err := c.NewGetRecentFillsService(CategoryUSDTFutures, "BTCUSDT").SetLimit(5).Do(cx)
		if err != nil {
			t.Fatalf("fills: %v", err)
		}
		if len(fills) == 0 {
			t.Fatal("no fills returned")
		}
		t.Logf("fill: %+v", fills[0])
		raw := fetchRawGet(t, c, cx, "/api/v3/market/fills", params, false)
		assertCovers(t, "market/fills", raw, fills)
	}

	// RPI symbols.
	{
		syms, err := c.NewGetRPISymbolsService().Do(cx)
		if err != nil {
			t.Fatalf("rpi-symbols: %v", err)
		}
		if len(syms) == 0 {
			t.Fatal("no rpi-symbols returned")
		}
		t.Logf("rpi-symbols: %d (first %+v)", len(syms), syms[0])
		raw := fetchRawGet(t, c, cx, "/api/v3/market/rpi-symbols", nil, false)
		assertCovers(t, "market/rpi-symbols", raw, syms)
	}
}
