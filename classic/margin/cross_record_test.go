package margin

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginCrossRecord(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	end := time.Now()
	start := end.AddDate(0, 0, -7)
	startMs := strconv.FormatInt(start.UnixMilli(), 10)
	endMs := strconv.FormatInt(end.UnixMilli(), 10)
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	// GET /api/v2/margin/crossed/borrow-history — Cross Loan Records (startTime required).
	{
		const label = "margin/crossed/borrow-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs, "coin": "USDT"}
		resp, err := c.NewGetCrossBorrowHistoryService(start).SetEndTime(end).SetCoin("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/borrow-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/margin/crossed/repay-history — Cross Repay History (startTime required).
	{
		const label = "margin/crossed/repay-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs, "coin": "USDT"}
		resp, err := c.NewGetCrossRepayHistoryService(start).SetEndTime(end).SetCoin("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/repay-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/margin/crossed/interest-history — Cross Interest Records (startTime required).
	{
		const label = "margin/crossed/interest-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs, "coin": "USDT"}
		resp, err := c.NewGetCrossInterestHistoryService(start).SetEndTime(end).SetCoin("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/interest-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/margin/crossed/liquidation-history — Cross Liquidation Records (startTime required).
	{
		const label = "margin/crossed/liquidation-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs}
		resp, err := c.NewGetCrossLiquidationHistoryService(start).SetEndTime(end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/liquidation-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/margin/crossed/financial-records — Cross Finance Flow History (startTime required).
	{
		const label = "margin/crossed/financial-records"
		params := map[string]string{"startTime": startMs, "endTime": endMs, "coin": "USDT"}
		resp, err := c.NewGetCrossFinancialRecordsService(start).SetEndTime(end).SetCoin("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/financial-records", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}
}
