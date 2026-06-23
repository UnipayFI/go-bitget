package mix

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixCandle(t *testing.T) {
	c := NewMixClient(apitest.PublicOptions()...)
	cx := apitest.Ctx(t)

	const (
		symbol      = "BTCUSDT"
		productType = ProductTypeUSDTFutures
		granularity = KlineGranularity1m
	)
	endTime := time.Now()
	endTimeStr := strconv.FormatInt(endTime.UnixMilli(), 10)

	// GET /api/v2/mix/market/candles
	{
		resp, err := c.NewGetCandlesService(symbol, productType, granularity).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("candles: empty response")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/candles", map[string]string{
			"symbol":      symbol,
			"productType": string(productType),
			"granularity": string(granularity),
			"limit":       "2",
		}, false)
		apitest.AssertCovers(t, "candles", raw, resp)
	}

	// GET /api/v2/mix/market/history-candles
	{
		resp, err := c.NewGetHistoryCandlesService(symbol, productType, granularity).SetEndTime(endTime).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("history-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("history-candles: empty response")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/history-candles", map[string]string{
			"symbol":      symbol,
			"productType": string(productType),
			"granularity": string(granularity),
			"endTime":     endTimeStr,
			"limit":       "2",
		}, false)
		apitest.AssertCovers(t, "history-candles", raw, resp)
	}

	// GET /api/v2/mix/market/history-index-candles
	{
		resp, err := c.NewGetHistoryIndexCandlesService(symbol, productType, granularity).SetEndTime(endTime).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("history-index-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("history-index-candles: empty response")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/history-index-candles", map[string]string{
			"symbol":      symbol,
			"productType": string(productType),
			"granularity": string(granularity),
			"endTime":     endTimeStr,
			"limit":       "2",
		}, false)
		apitest.AssertCovers(t, "history-index-candles", raw, resp)
	}

	// GET /api/v2/mix/market/history-mark-candles
	{
		resp, err := c.NewGetHistoryMarkCandlesService(symbol, productType, granularity).SetEndTime(endTime).SetLimit(2).Do(cx)
		if err != nil {
			t.Fatalf("history-mark-candles: %v", err)
		}
		if len(resp) == 0 {
			t.Fatalf("history-mark-candles: empty response")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/history-mark-candles", map[string]string{
			"symbol":      symbol,
			"productType": string(productType),
			"granularity": string(granularity),
			"endTime":     endTimeStr,
			"limit":       "2",
		}, false)
		apitest.AssertCovers(t, "history-mark-candles", raw, resp)
	}
}
