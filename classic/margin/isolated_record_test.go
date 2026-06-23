package margin

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginIsolatedRecord(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const symbol = "BTCUSDT"
	end := time.Now()
	start := end.Add(-30 * 24 * time.Hour)
	// Tolerable capability/empty codes for permission-limited sub-account keys.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	window := map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(start.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
	}

	// repay-history
	{
		resp, err := c.NewGetIsolatedRepayHistoryService(symbol, start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated/repay-history", err, okCodes...) {
				t.Fatalf("isolated/repay-history: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/repay-history", window, true)
			apitest.AssertCovers(t, "isolated/repay-history", raw, resp)
		}
	}

	// borrow-history
	{
		resp, err := c.NewGetIsolatedBorrowHistoryService(symbol, start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated/borrow-history", err, okCodes...) {
				t.Fatalf("isolated/borrow-history: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/borrow-history", window, true)
			apitest.AssertCovers(t, "isolated/borrow-history", raw, resp)
		}
	}

	// interest-history
	{
		resp, err := c.NewGetIsolatedInterestHistoryService(symbol, start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated/interest-history", err, okCodes...) {
				t.Fatalf("isolated/interest-history: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/interest-history", window, true)
			apitest.AssertCovers(t, "isolated/interest-history", raw, resp)
		}
	}

	// liquidation-history
	{
		resp, err := c.NewGetIsolatedLiquidationHistoryService(symbol, start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated/liquidation-history", err, okCodes...) {
				t.Fatalf("isolated/liquidation-history: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/liquidation-history", window, true)
			apitest.AssertCovers(t, "isolated/liquidation-history", raw, resp)
		}
	}

	// financial-records
	{
		resp, err := c.NewGetIsolatedFinancialRecordsService(symbol, start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated/financial-records", err, okCodes...) {
				t.Fatalf("isolated/financial-records: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/financial-records", window, true)
			apitest.AssertCovers(t, "isolated/financial-records", raw, resp)
		}
	}
}
