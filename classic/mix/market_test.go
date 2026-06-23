package mix

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixMarket(t *testing.T) {
	c := NewMixClient(apitest.PublicOptions()...)
	cx := apitest.Ctx(t)

	const (
		symbol = "BTCUSDT"
		pt     = ProductTypeUSDTFutures
	)
	ptParams := map[string]string{"productType": string(pt)}
	symParams := map[string]string{"symbol": symbol, "productType": string(pt)}

	// merge-depth
	{
		resp, err := c.NewGetMergeDepthService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("merge-depth: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/merge-depth", symParams, false)
		apitest.AssertCovers(t, "merge-depth", raw, resp)
	}

	// ticker
	{
		resp, err := c.NewGetTickerService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("ticker: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/ticker", symParams, false)
		apitest.AssertCovers(t, "ticker", raw, resp)
	}

	// tickers (all)
	{
		resp, err := c.NewGetAllTickersService(pt).Do(cx)
		if err != nil {
			t.Fatalf("tickers: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/tickers", ptParams, false)
		apitest.AssertCovers(t, "tickers", raw, resp)
	}

	// fills (recent)
	{
		resp, err := c.NewGetRecentFillsService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("fills: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/fills", symParams, false)
		apitest.AssertCovers(t, "fills", raw, resp)
	}

	// fills-history
	{
		resp, err := c.NewGetFillsHistoryService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("fills-history: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/fills-history", symParams, false)
		apitest.AssertCovers(t, "fills-history", raw, resp)
	}

	// open-interest
	{
		resp, err := c.NewGetOpenInterestService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("open-interest: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/open-interest", symParams, false)
		apitest.AssertCovers(t, "open-interest", raw, resp)
	}

	// funding-time
	{
		resp, err := c.NewGetFundingTimeService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("funding-time: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/funding-time", symParams, false)
		apitest.AssertCovers(t, "funding-time", raw, resp)
	}

	// symbol-price
	{
		resp, err := c.NewGetSymbolPriceService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("symbol-price: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/symbol-price", symParams, false)
		apitest.AssertCovers(t, "symbol-price", raw, resp)
	}

	// oi-limit
	{
		resp, err := c.NewGetOiLimitService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("oi-limit: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/oi-limit", symParams, false)
		apitest.AssertCovers(t, "oi-limit", raw, resp)
	}

	// contracts
	{
		resp, err := c.NewGetContractsService(pt).SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("contracts: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/contracts", symParams, false)
		apitest.AssertCovers(t, "contracts", raw, resp)
	}

	// query-position-lever
	{
		resp, err := c.NewGetQueryPositionLeverService(symbol, pt).Do(cx)
		if err != nil {
			t.Fatalf("query-position-lever: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/query-position-lever", symParams, false)
		apitest.AssertCovers(t, "query-position-lever", raw, resp)
	}
}
