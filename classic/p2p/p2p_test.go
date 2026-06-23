package p2p

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestP2P(t *testing.T) {
	c := NewP2PClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	// merchantList is the only P2P read that does not require pre-existing
	// merchant/order state. Calling the Service verifies the endpoint works and
	// that the response deserializes; AssertCovers verifies field coverage. The
	// endpoints are P2P-merchant-gated, so 60039/40014 (and the usual
	// capability/empty codes) are tolerated — signing + path are still proven.
	t.Run("merchantList", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetP2PMerchantListService().SetLimit(10).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "merchantList", err, "60039", "40014", "40099", "40034", "40068", "40054") {
				return
			}
			t.Fatalf("merchantList: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/p2p/merchantList", map[string]string{"limit": "10"}, true)
		apitest.AssertCovers(t, "merchantList", raw, resp)
	})
}
