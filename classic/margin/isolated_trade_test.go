package margin

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginIsolatedTrade(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const symbol = "BTCUSDT"
	endTime := time.Now()
	startTime := endTime.Add(-72 * time.Hour)
	startMs := strconv.FormatInt(startTime.UnixMilli(), 10)

	// Tolerated capability/empty codes for margin sub-account reads.
	tolerable := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	// open-orders
	{
		resp, err := c.NewGetIsolatedOpenOrdersService(symbol, startTime).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "isolated open-orders", err, tolerable...) {
				goto historyOrders
			}
			t.Fatalf("isolated open-orders: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/open-orders",
			map[string]string{"symbol": symbol, "startTime": startMs}, true)
		apitest.AssertCovers(t, "isolated open-orders", raw, resp)
	}

historyOrders:
	// history-orders
	{
		resp, err := c.NewGetIsolatedHistoryOrdersService(symbol, startTime).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "isolated history-orders", err, tolerable...) {
				goto fills
			}
			t.Fatalf("isolated history-orders: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/history-orders",
			map[string]string{"symbol": symbol, "startTime": startMs}, true)
		apitest.AssertCovers(t, "isolated history-orders", raw, resp)
	}

fills:
	// fills
	{
		resp, err := c.NewGetIsolatedFillsService(symbol, startTime).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "isolated fills", err, tolerable...) {
				goto liquidationOrders
			}
			t.Fatalf("isolated fills: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/fills",
			map[string]string{"symbol": symbol, "startTime": startMs}, true)
		apitest.AssertCovers(t, "isolated fills", raw, resp)
	}

liquidationOrders:
	// liquidation-order
	{
		resp, err := c.NewGetIsolatedLiquidationOrdersService().SetSymbol(symbol).
			SetStartTime(startTime).SetEndTime(endTime).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "isolated liquidation-order", err, tolerable...) {
				return
			}
			t.Fatalf("isolated liquidation-order: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/liquidation-order",
			map[string]string{
				"symbol":    symbol,
				"startTime": startMs,
				"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
			}, true)
		apitest.AssertCovers(t, "isolated liquidation-order", raw, resp)
	}
}
