package spot

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestSpotMarket(t *testing.T) {
	c := NewSpotClient(apitest.PublicOptions()...)
	ctx := apitest.Ctx(t)

	// Get Coin Info.
	{
		params := map[string]string{"coin": "BTC"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/public/coins", params, false)
		resp, err := c.NewGetCoinsService().SetCoin("BTC").Do(ctx)
		if err != nil {
			t.Fatalf("coins: %v", err)
		}
		apitest.AssertCovers(t, "spot coins", raw, resp)
	}

	// Get Symbol Info.
	{
		params := map[string]string{"symbol": "BTCUSDT"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/public/symbols", params, false)
		resp, err := c.NewGetSymbolsService().SetSymbol("BTCUSDT").Do(ctx)
		if err != nil {
			t.Fatalf("symbols: %v", err)
		}
		apitest.AssertCovers(t, "spot symbols", raw, resp)
	}

	// Get VIP Fee Rate.
	{
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/vip-fee-rate", nil, false)
		resp, err := c.NewGetVIPFeeRateService().Do(ctx)
		if err != nil {
			t.Fatalf("vip-fee-rate: %v", err)
		}
		apitest.AssertCovers(t, "spot vip-fee-rate", raw, resp)
	}

	// Get Tickers.
	{
		params := map[string]string{"symbol": "BTCUSDT"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/tickers", params, false)
		resp, err := c.NewGetTickersService().SetSymbol("BTCUSDT").Do(ctx)
		if err != nil {
			t.Fatalf("tickers: %v", err)
		}
		apitest.AssertCovers(t, "spot tickers", raw, resp)
	}

	// Get Merge Depth.
	{
		params := map[string]string{"symbol": "BTCUSDT", "precision": "scale0", "limit": "5"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/merge-depth", params, false)
		resp, err := c.NewGetMergeDepthService("BTCUSDT").
			SetPrecision(MergeDepthPrecisionScale0).SetLimit("5").Do(ctx)
		if err != nil {
			t.Fatalf("merge-depth: %v", err)
		}
		apitest.AssertCovers(t, "spot merge-depth", raw, resp)
	}

	// Get OrderBook Depth.
	{
		params := map[string]string{"symbol": "BTCUSDT", "type": "step0", "limit": "5"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/orderbook", params, false)
		resp, err := c.NewGetOrderBookService("BTCUSDT").
			SetType(OrderBookTypeStep0).SetLimit(5).Do(ctx)
		if err != nil {
			t.Fatalf("orderbook: %v", err)
		}
		apitest.AssertCovers(t, "spot orderbook", raw, resp)
	}

	// Get Candlestick Data.
	{
		resp, err := c.NewGetCandlesService("BTCUSDT", "1min").SetLimit(5).Do(ctx)
		if err != nil {
			t.Fatalf("candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("candles: empty result")
		}
	}

	// Get History Candlestick Data. endTime is required for this endpoint.
	{
		endTime := time.Now()
		resp, err := c.NewGetHistoryCandlesService("BTCUSDT", "1min").SetEndTime(endTime).SetLimit(5).Do(ctx)
		if err != nil {
			t.Fatalf("history-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("history-candles: empty result")
		}
	}

	// Get Recent Trades.
	{
		params := map[string]string{"symbol": "BTCUSDT", "limit": "5"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/fills", params, false)
		resp, err := c.NewGetRecentFillsService("BTCUSDT").SetLimit(5).Do(ctx)
		if err != nil {
			t.Fatalf("fills: %v", err)
		}
		apitest.AssertCovers(t, "spot fills", raw, resp)
	}

	// Get Market Trades.
	{
		params := map[string]string{"symbol": "BTCUSDT", "limit": "5"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/fills-history", params, false)
		resp, err := c.NewGetMarketTradesService("BTCUSDT").SetLimit(5).Do(ctx)
		if err != nil {
			t.Fatalf("fills-history: %v", err)
		}
		apitest.AssertCovers(t, "spot fills-history", raw, resp)
	}

	// Get Call Auction Information.
	{
		params := map[string]string{"symbol": "BTCUSDT"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/auction", params, false)
		resp, err := c.NewGetAuctionService("BTCUSDT").Do(ctx)
		if err != nil {
			t.Fatalf("auction: %v", err)
		}
		apitest.AssertCovers(t, "spot auction", raw, resp)
	}
}
