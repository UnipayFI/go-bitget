package common

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestInsights(t *testing.T) {
	c := NewCommonClient(apitest.PublicOptions()...)
	ctx := apitest.Ctx(t)

	const symbol = "BTCUSDT"

	// All endpoints in this group are public GETs, so each is live-tested.
	raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/whale-net-flow", map[string]string{"symbol": symbol}, false)
	apitest.AssertCovers(t, "WhaleNetFlow", raw, []WhaleNetFlow{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/market/taker-buy-sell", map[string]string{"symbol": symbol, "period": "5m"}, false)
	apitest.AssertCovers(t, "TakerBuySell", raw, []TakerBuySell{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/market/position-long-short", map[string]string{"symbol": symbol, "period": "5m"}, false)
	apitest.AssertCovers(t, "PositionLongShort", raw, []PositionLongShort{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/margin/market/long-short-ratio", map[string]string{"symbol": symbol, "period": "24h"}, false)
	apitest.AssertCovers(t, "MarginLongShortRatio", raw, []MarginLongShortRatio{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/margin/market/loan-growth", map[string]string{"symbol": symbol, "period": "24h"}, false)
	apitest.AssertCovers(t, "MarginLoanGrowth", raw, []MarginLoanGrowth{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/margin/market/isolated-borrow-rate", map[string]string{"symbol": symbol, "period": "24h"}, false)
	apitest.AssertCovers(t, "IsolatedBorrowRate", raw, []IsolatedBorrowRate{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/market/long-short", map[string]string{"symbol": symbol, "period": "5m"}, false)
	apitest.AssertCovers(t, "LongShort", raw, []LongShort{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/fund-flow", map[string]string{"symbol": symbol, "period": "15m"}, false)
	apitest.AssertCovers(t, "SpotFundFlow", raw, SpotFundFlow{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/support-symbols", nil, false)
	apitest.AssertCovers(t, "SupportSymbols", raw, SupportSymbols{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/market/fund-net-flow", map[string]string{"symbol": symbol}, false)
	apitest.AssertCovers(t, "FundNetFlow", raw, []FundNetFlow{})

	raw = apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/market/account-long-short", map[string]string{"symbol": symbol, "period": "5m"}, false)
	apitest.AssertCovers(t, "AccountLongShort", raw, []AccountLongShort{})
}
