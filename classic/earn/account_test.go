package earn

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestEarnAccount(t *testing.T) {
	c := NewEarnClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	tol := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

	// GET /api/v2/earn/account/assets
	{
		resp, err := c.NewGetEarnAccountAssetsService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "earn-account-assets", err, tol...) {
				t.Fatalf("earn-account-assets: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/account/assets", nil, true)
			apitest.AssertCovers(t, "earn-account-assets", raw, resp)
		}
	}
}
