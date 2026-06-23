package copy

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// tolerable Bitget codes for a permission-limited sub-account: capability or
// data gaps that still prove the endpoint path + signing are correct.
var spotFollowerTolerable = []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808", "400171", "400172", "40805", "31001", "31002", "40913", "40020", "43046", "43011", "43025", "40732", "40733", "40734"}

func TestCopySpotFollower(t *testing.T) {
	c := NewCopyClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// A traderId we may discover from the my-trader list, used to exercise the
	// trader-scoped reads below.
	var traderID string

	// query-traders -- no required params.
	{
		resp, err := c.NewGetSpotFollowerTradersService().SetPageSize("10").Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "query-traders", err, spotFollowerTolerable...) {
				t.Fatalf("query-traders: %v", err)
			}
		} else {
			if len(resp.ResultList) > 0 {
				traderID = resp.ResultList[0].TraderID
			}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-follower/query-traders", map[string]string{"pageSize": "10"}, true)
			apitest.AssertCovers(t, "query-traders", raw, resp)
		}
	}

	// query-current-orders -- no required params.
	{
		resp, err := c.NewGetSpotFollowerCurrentOrdersService().SetLimit("10").Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "query-current-orders", err, spotFollowerTolerable...) {
				t.Fatalf("query-current-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-follower/query-current-orders", map[string]string{"limit": "10"}, true)
			apitest.AssertCovers(t, "query-current-orders", raw, resp)
		}
	}

	// query-history-orders -- no required params.
	{
		resp, err := c.NewGetSpotFollowerHistoryOrdersService().SetLimit("10").Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "query-history-orders", err, spotFollowerTolerable...) {
				t.Fatalf("query-history-orders: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-follower/query-history-orders", map[string]string{"limit": "10"}, true)
			apitest.AssertCovers(t, "query-history-orders", raw, resp)
		}
	}

	if traderID == "" {
		t.Skip("query-trader-symbols/query-settings: no followed trader available to derive a traderId; trader-scoped reads skipped")
	}

	// query-trader-symbols -- requires traderId.
	{
		resp, err := c.NewGetSpotFollowerTraderSymbolsService(traderID).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "query-trader-symbols", err, spotFollowerTolerable...) {
				t.Fatalf("query-trader-symbols: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-follower/query-trader-symbols", map[string]string{"traderId": traderID}, true)
			apitest.AssertCovers(t, "query-trader-symbols", raw, resp)
		}
	}

	// query-settings -- requires traderId.
	{
		resp, err := c.NewGetSpotFollowerSettingsService(traderID).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "query-settings", err, spotFollowerTolerable...) {
				t.Fatalf("query-settings: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-follower/query-settings", map[string]string{"traderId": traderID}, true)
			apitest.AssertCovers(t, "query-settings", raw, resp)
		}
	}
}
