package earn

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

func TestEarnLoan(t *testing.T) {
	// Tolerable capability/empty codes for the permission-limited shared key.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

	// --- Public endpoints (no signing). ---
	pub := NewEarnClient(apitest.PublicOptions()...)
	pcx := apitest.Ctx(t)

	// GET /api/v2/earn/loan/public/coinInfos — Currency List (public).
	{
		const label = "earn/loan/public/coinInfos"
		params := map[string]string{"coin": "USDT"}
		resp, err := pub.NewGetLoanCoinInfosService().SetCoin("USDT").Do(pcx)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, pub, pcx, "/api/v2/earn/loan/public/coinInfos", params, false)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/public/hour-interest — Est. Interest and Borrowable (public).
	{
		const label = "earn/loan/public/hour-interest"
		params := map[string]string{"loanCoin": "USDT", "pledgeCoin": "BNB", "daily": "SEVEN", "pledgeAmount": "10"}
		resp, err := pub.NewGetLoanHourInterestService("USDT", "BNB", LoanDailySeven, decimal.RequireFromString("10")).Do(pcx)
		if err != nil {
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, pub, pcx, "/api/v2/earn/loan/public/hour-interest", params, false)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// --- Private read endpoints (signed). ---
	c := NewEarnClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	end := time.Now()
	start := end.AddDate(0, 0, -30)
	startMs := strconv.FormatInt(start.UnixMilli(), 10)
	endMs := strconv.FormatInt(end.UnixMilli(), 10)

	// GET /api/v2/earn/loan/ongoing-orders — Loan Orders (no required params).
	{
		const label = "earn/loan/ongoing-orders"
		params := map[string]string{}
		resp, err := c.NewGetLoanOngoingOrdersService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/ongoing-orders", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/debts — Debts (no required params).
	{
		const label = "earn/loan/debts"
		params := map[string]string{}
		resp, err := c.NewGetLoanDebtsService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/debts", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/repay-history — Repay History (startTime/endTime required).
	{
		const label = "earn/loan/repay-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs}
		resp, err := c.NewGetLoanRepayHistoryService(start, end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/repay-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/revise-history — Pledge Rate History (startTime/endTime required).
	{
		const label = "earn/loan/revise-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs}
		resp, err := c.NewGetLoanReviseHistoryService(start, end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/revise-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/borrow-history — Loan History (startTime/endTime required).
	{
		const label = "earn/loan/borrow-history"
		params := map[string]string{"startTime": startMs, "endTime": endMs}
		resp, err := c.NewGetLoanBorrowHistoryService(start, end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/borrow-history", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/earn/loan/reduces — Liquidation Records (startTime/endTime required).
	{
		const label = "earn/loan/reduces"
		params := map[string]string{"startTime": startMs, "endTime": endMs}
		resp, err := c.NewGetLoanReducesService(start, end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, okCodes...) {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/loan/reduces", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}
}
