package margin

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginCrossTrade(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const symbol = "BTCUSDT"
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	startMs := strconv.FormatInt(start.UnixMilli(), 10)
	endMs := strconv.FormatInt(end.UnixMilli(), 10)

	// tolerable: capability-gated sub-account / empty-data / no-permission codes.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	// Cross open orders.
	{
		resp, err := c.NewGetCrossOpenOrdersService(symbol).SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "cross-open-orders", err, okCodes...) {
				t.Fatalf("cross-open-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/open-orders", map[string]string{
				"symbol":    symbol,
				"startTime": startMs,
				"endTime":   endMs,
			}, true)
			apitest.AssertCovers(t, "cross-open-orders", raw, resp)
		}
	}

	// Cross history orders.
	{
		resp, err := c.NewGetCrossHistoryOrdersService(symbol).SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "cross-history-orders", err, okCodes...) {
				t.Fatalf("cross-history-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/history-orders", map[string]string{
				"symbol":    symbol,
				"startTime": startMs,
				"endTime":   endMs,
			}, true)
			apitest.AssertCovers(t, "cross-history-orders", raw, resp)
		}
	}

	// Cross fills.
	{
		resp, err := c.NewGetCrossFillsService(symbol).SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "cross-fills", err, okCodes...) {
				t.Fatalf("cross-fills: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/fills", map[string]string{
				"symbol":    symbol,
				"startTime": startMs,
				"endTime":   endMs,
			}, true)
			apitest.AssertCovers(t, "cross-fills", raw, resp)
		}
	}

	// Cross liquidation orders.
	{
		resp, err := c.NewGetCrossLiquidationOrdersService().SetSymbol(symbol).SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "cross-liquidation-orders", err, okCodes...) {
				t.Fatalf("cross-liquidation-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/liquidation-order", map[string]string{
				"symbol":    symbol,
				"startTime": startMs,
				"endTime":   endMs,
			}, true)
			apitest.AssertCovers(t, "cross-liquidation-orders", raw, resp)
		}
	}
}
