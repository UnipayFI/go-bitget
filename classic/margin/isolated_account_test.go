package margin

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginIsolatedAccount(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const symbol = "BTCUSDT"
	tolerable := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "50021", "50067", "50001"}

	// GET /api/v2/margin/isolated/account/assets
	{
		resp, err := c.NewGetIsolatedAssetsService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-assets", err, tolerable...) {
				t.Fatalf("isolated-assets: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/account/assets", map[string]string{}, true)
			apitest.AssertCovers(t, "isolated-assets", raw, resp)
		}
	}

	// GET /api/v2/margin/isolated/account/risk-rate
	{
		resp, err := c.NewGetIsolatedRiskRateService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-risk-rate", err, tolerable...) {
				t.Fatalf("isolated-risk-rate: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/account/risk-rate", map[string]string{}, true)
			apitest.AssertCovers(t, "isolated-risk-rate", raw, resp)
		}
	}

	// GET /api/v2/margin/isolated/interest-rate-and-limit
	{
		params := map[string]string{"symbol": symbol}
		resp, err := c.NewGetIsolatedInterestRateAndLimitService(symbol).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-interest-rate-and-limit", err, tolerable...) {
				t.Fatalf("isolated-interest-rate-and-limit: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/interest-rate-and-limit", params, true)
			apitest.AssertCovers(t, "isolated-interest-rate-and-limit", raw, resp)
		}
	}

	// GET /api/v2/margin/isolated/tier-data
	{
		params := map[string]string{"symbol": symbol}
		resp, err := c.NewGetIsolatedTierDataService(symbol).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-tier-data", err, tolerable...) {
				t.Fatalf("isolated-tier-data: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/tier-data", params, true)
			apitest.AssertCovers(t, "isolated-tier-data", raw, resp)
		}
	}

	// GET /api/v2/margin/isolated/account/max-borrowable-amount
	{
		params := map[string]string{"symbol": symbol}
		resp, err := c.NewGetIsolatedMaxBorrowableService(symbol).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-max-borrowable", err, tolerable...) {
				t.Fatalf("isolated-max-borrowable: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/account/max-borrowable-amount", params, true)
			apitest.AssertCovers(t, "isolated-max-borrowable", raw, resp)
		}
	}

	// GET /api/v2/margin/isolated/account/max-transfer-out-amount
	{
		params := map[string]string{"symbol": symbol}
		resp, err := c.NewGetIsolatedMaxTransferOutService(symbol).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-max-transfer-out", err, tolerable...) {
				t.Fatalf("isolated-max-transfer-out: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/isolated/account/max-transfer-out-amount", params, true)
			apitest.AssertCovers(t, "isolated-max-transfer-out", raw, resp)
		}
	}
}
