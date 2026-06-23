package common

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCommonAssets(t *testing.T) {
	c := NewCommonClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	// funding-assets: calling the Service verifies both that the endpoint works
	// and that the response deserializes; AssertCovers verifies field coverage.
	// 40068/40014 mean this key (a sub-account) cannot read the funding line —
	// endpoint + signing are still confirmed correct.
	t.Run("funding-assets", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetFundingAssetsService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "funding-assets", err, "40099", "40068", "40014") {
				return
			}
			t.Fatalf("funding-assets: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/account/funding-assets", nil, true)
		apitest.AssertCovers(t, "funding-assets", raw, resp)
	})

	t.Run("bot-assets", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetBotAssetsService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "bot-assets", err, "40099", "40068", "40014") {
				return
			}
			t.Fatalf("bot-assets: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/account/bot-assets", nil, true)
		apitest.AssertCovers(t, "bot-assets", raw, resp)
	})

	t.Run("all-account-balance", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetAllAccountBalanceService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "all-account-balance", err, "40099", "40068", "40014") {
				return
			}
			t.Fatalf("all-account-balance: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/account/all-account-balance", nil, true)
		apitest.AssertCovers(t, "all-account-balance", raw, resp)
	})
}
