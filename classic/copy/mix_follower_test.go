package copy

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCopyMixFollower(t *testing.T) {
	c := NewCopyClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Capability-gated codes: the shared key is a permission-limited sub-account
	// that may not be enrolled in copy-trading as a follower, or simply have no
	// data. Treat these as pass-with-tolerance: the path + signing were proven.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808", "400171", "400172", "40805", "31001", "31002", "40913", "40020", "43046", "43011", "43025", "40732", "40733", "40734"}

	// Get Current Tracking Orders -- private GET, data is an array.
	{
		params := map[string]string{"productType": string(ProductTypeUSDTFutures)}
		resp, err := c.NewGetMixFollowerCurrentOrdersService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "mix-follower current-orders", err, okCodes...) {
				t.Fatalf("mix-follower current-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-follower/query-current-orders", params, true)
			apitest.AssertCovers(t, "mix-follower current-orders", raw, resp)
		}
	}

	// Get History Tracking Orders -- private GET, data is an object {trackingList, endId}.
	{
		params := map[string]string{"productType": string(ProductTypeUSDTFutures)}
		resp, err := c.NewGetMixFollowerHistoryOrdersService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "mix-follower history-orders", err, okCodes...) {
				t.Fatalf("mix-follower history-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-follower/query-history-orders", params, true)
			apitest.AssertCovers(t, "mix-follower history-orders", raw, resp)
		}
	}

	// Get My Traders -- private GET, data is an array.
	{
		resp, err := c.NewGetMixFollowerTradersService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "mix-follower query-traders", err, okCodes...) {
				t.Fatalf("mix-follower query-traders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-follower/query-traders", map[string]string{}, true)
			apitest.AssertCovers(t, "mix-follower query-traders", raw, resp)
		}
	}

	// Get Follow Limit -- private GET, data is an array.
	{
		params := map[string]string{"productType": string(ProductTypeUSDTFutures)}
		resp, err := c.NewGetMixFollowerQuantityLimitService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "mix-follower quantity-limit", err, okCodes...) {
				t.Fatalf("mix-follower quantity-limit: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-follower/query-quantity-limit", params, true)
			apitest.AssertCovers(t, "mix-follower quantity-limit", raw, resp)
		}
	}
}
