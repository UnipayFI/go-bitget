package broker

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestBrokerCommission(t *testing.T) {
	c := NewBrokerClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// broker-gated tolerances: 40029 (not a broker) plus the shared
	// capability/empty-data codes the test key (a permission-limited
	// sub-account) may return.
	okCodes := []string{"40029", "40068", "40014", "40054", "40099", "40034", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	window := map[string]string{
		"startTime": strconv.FormatInt(start.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
	}

	// Get Total Commission -- GET /api/v2/broker/total-commission
	{
		resp, err := c.NewGetTotalCommissionService().SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "total-commission", err, okCodes...) {
				t.Fatalf("total-commission: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/total-commission", window, true)
			apitest.AssertCovers(t, "total-commission", raw, resp)
		}
	}

	// Get Order Commission -- GET /api/v2/broker/order-commission
	{
		resp, err := c.NewGetOrderCommissionService().SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "order-commission", err, okCodes...) {
				t.Fatalf("order-commission: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/order-commission", window, true)
			apitest.AssertCovers(t, "order-commission", raw, resp)
		}
	}

	// Get Rebate Info -- GET /api/v2/broker/rebate-info (requires uid)
	// We have no known-valid uid for the shared key, so this verifies the
	// endpoint path + signing; tolerate broker/capability and invalid-uid
	// (40037 sub-account not found) responses.
	{
		const uid = "0"
		resp, err := c.NewGetRebateInfoService(uid).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "rebate-info", err, okCodes...) {
				t.Fatalf("rebate-info: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/rebate-info", map[string]string{"uid": uid}, true)
			apitest.AssertCovers(t, "rebate-info", raw, resp)
		}
	}
}
