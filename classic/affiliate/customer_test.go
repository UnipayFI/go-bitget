package affiliate

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestAffiliateCustomer(t *testing.T) {
	c := NewAffiliateClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// tolerable: this product line is agent/affiliate-gated, so a non-agent key
	// (or one lacking child-account perms) and empty data are expected passes.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	startMs := strconv.FormatInt(start.UnixMilli(), 10)
	endMs := strconv.FormatInt(end.UnixMilli(), 10)

	// customer-commissions (GET)
	{
		resp, err := c.NewGetCustomerCommissionsService().
			SetStartTime(start).SetEndTime(end).SetLimit(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "customer-commissions", err, okCodes...) {
				t.Fatalf("customer-commissions: %v", err)
			}
		} else {
			params := map[string]string{"startTime": startMs, "endTime": endMs, "limit": "100"}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/customer-commissions", params, true)
			apitest.AssertCovers(t, "customer-commissions", raw, resp)
		}
	}

	// customer-kyc-result (GET)
	{
		resp, err := c.NewGetCustomerKycResultService().
			SetStartTime(start).SetEndTime(end).SetLimit(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "customer-kyc-result", err, okCodes...) {
				t.Fatalf("customer-kyc-result: %v", err)
			}
		} else {
			params := map[string]string{"startTime": startMs, "endTime": endMs, "limit": "100"}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/customer-kyc-result", params, true)
			apitest.AssertCovers(t, "customer-kyc-result", raw, resp)
		}
	}

	// agent-commission (GET)
	{
		resp, err := c.NewGetAgentCommissionService().
			SetStartTime(start).SetEndTime(end).SetLimit(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "agent-commission", err, okCodes...) {
				t.Fatalf("agent-commission: %v", err)
			}
		} else {
			params := map[string]string{"startTime": startMs, "endTime": endMs, "limit": "100"}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/agent-commission", params, true)
			apitest.AssertCovers(t, "agent-commission", raw, resp)
		}
	}
}
