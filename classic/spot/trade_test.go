package spot

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// TestSpotTrade live-tests the safe private GET reads in the spot trade group:
// current (unfilled) orders, history orders and fills. These work without any
// pre-existing order state (they simply return an empty array on a clean
// account). The POST endpoints (place/cancel/batch/cancel-replace/cancel-symbol)
// are state-changing and are NOT exercised here.
func TestSpotTrade(t *testing.T) {
	c := NewSpotClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	t.Run("unfilled-orders", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		params := map[string]string{"symbol": "BTCUSDT", "limit": "100"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/unfilled-orders", params, true)

		resp, err := c.NewGetUnfilledOrdersService().SetSymbol("BTCUSDT").SetLimit(100).Do(ctx)
		if err != nil {
			t.Fatalf("unfilled-orders: %v", err)
		}
		apitest.AssertCovers(t, "unfilled-orders", raw, resp)
	})

	t.Run("history-orders", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		params := map[string]string{"symbol": "BTCUSDT", "limit": "100"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/history-orders", params, true)

		resp, err := c.NewGetHistoryOrdersService().SetSymbol("BTCUSDT").SetLimit(100).Do(ctx)
		if err != nil {
			t.Fatalf("history-orders: %v", err)
		}
		apitest.AssertCovers(t, "history-orders", raw, resp)
	})

	t.Run("fills", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		params := map[string]string{"symbol": "BTCUSDT", "limit": "100"}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/fills", params, true)

		resp, err := c.NewGetFillsService().SetSymbol("BTCUSDT").SetLimit(100).Do(ctx)
		if err != nil {
			t.Fatalf("fills: %v", err)
		}
		apitest.AssertCovers(t, "fills", raw, resp)
	})
}
