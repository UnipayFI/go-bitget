package margin

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginCrossAccount(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// codes tolerated across all cross-margin reads: capability/permission/empty.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	// GET /api/v2/margin/crossed/account/assets
	{
		resp, err := c.NewGetCrossAccountAssetsService().SetCoin("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-assets", err, okCodes...) {
				goto riskRate
			}
			t.Fatalf("cross-assets: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/account/assets", map[string]string{"coin": "USDT"}, true)
		apitest.AssertCovers(t, "cross-assets", raw, resp)
	}

riskRate:
	// GET /api/v2/margin/crossed/account/risk-rate
	{
		resp, err := c.NewGetCrossRiskRateService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-risk-rate", err, okCodes...) {
				goto maxBorrowable
			}
			t.Fatalf("cross-risk-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/account/risk-rate", nil, true)
		apitest.AssertCovers(t, "cross-risk-rate", raw, resp)
	}

maxBorrowable:
	// GET /api/v2/margin/crossed/account/max-borrowable-amount
	{
		resp, err := c.NewGetCrossMaxBorrowableService("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-max-borrowable", err, okCodes...) {
				goto maxTransferOut
			}
			t.Fatalf("cross-max-borrowable: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/account/max-borrowable-amount", map[string]string{"coin": "USDT"}, true)
		apitest.AssertCovers(t, "cross-max-borrowable", raw, resp)
	}

maxTransferOut:
	// GET /api/v2/margin/crossed/account/max-transfer-out-amount
	{
		resp, err := c.NewGetCrossMaxTransferOutService("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-max-transfer-out", err, okCodes...) {
				goto interestRate
			}
			t.Fatalf("cross-max-transfer-out: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/account/max-transfer-out-amount", map[string]string{"coin": "USDT"}, true)
		apitest.AssertCovers(t, "cross-max-transfer-out", raw, resp)
	}

interestRate:
	// GET /api/v2/margin/crossed/interest-rate-and-limit
	{
		resp, err := c.NewGetCrossInterestRateAndLimitService("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-interest-rate-and-limit", err, okCodes...) {
				goto tierData
			}
			t.Fatalf("cross-interest-rate-and-limit: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/interest-rate-and-limit", map[string]string{"coin": "USDT"}, true)
		apitest.AssertCovers(t, "cross-interest-rate-and-limit", raw, resp)
	}

tierData:
	// GET /api/v2/margin/crossed/tier-data
	{
		resp, err := c.NewGetCrossTierDataService("USDT").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "cross-tier-data", err, okCodes...) {
				return
			}
			t.Fatalf("cross-tier-data: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/crossed/tier-data", map[string]string{"coin": "USDT"}, true)
		apitest.AssertCovers(t, "cross-tier-data", raw, resp)
	}

	// Not tested (state-changing or require pre-existing resource state):
	//   POST /api/v2/margin/crossed/account/borrow
	//   POST /api/v2/margin/crossed/account/repay
	//   POST /api/v2/margin/crossed/account/flash-repay
	//   POST /api/v2/margin/crossed/account/query-flash-repay-status (needs a real repayId)
}
