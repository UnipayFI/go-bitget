package copy

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCopyMixBroker(t *testing.T) {
	c := NewCopyClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Broker reads are gated behind a broker partnership; tolerate the
	// not-a-broker / needs-child-perms / empty-data responses so the test still
	// proves the endpoint path + signing.
	brokerCodes := []string{"40029", "40014", "40099", "40034", "40068", "40054", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808", "400171", "400172", "40805", "31001", "31002", "40913", "40020", "43046", "43011", "43025", "40732", "40733", "40734"}

	// 1) query-traders -- no required params; also yields a real traderId for
	// the trace endpoints below.
	traderID := "1"
	traders, err := c.NewGetMixBrokerTradersService().SetPageSize(20).Do(cx)
	if err != nil {
		if !apitest.Tolerable(t, "mix-broker query-traders", err, brokerCodes...) {
			t.Fatalf("mix-broker query-traders: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-broker/query-traders", map[string]string{"pageSize": "20"}, true)
		apitest.AssertCovers(t, "mix-broker query-traders", raw, traders)
		if len(traders) > 0 && traders[0].TraderID != "" {
			traderID = traders[0].TraderID
		}
	}

	// 2) query-history-traces -- requires traderId + productType.
	histParams := map[string]string{"traderId": traderID, "productType": string(ProductTypeUSDTFutures), "limit": "100"}
	hist, err := c.NewGetMixBrokerHistoryTracesService(traderID, ProductTypeUSDTFutures).SetLimit(100).Do(cx)
	if err != nil {
		if !apitest.Tolerable(t, "mix-broker query-history-traces", err, brokerCodes...) {
			t.Fatalf("mix-broker query-history-traces: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-broker/query-history-traces", histParams, true)
		apitest.AssertCovers(t, "mix-broker query-history-traces", raw, hist)
	}

	// 3) query-current-traces -- requires traderId + productType.
	curParams := map[string]string{"traderId": traderID, "productType": string(ProductTypeUSDTFutures), "limit": "100"}
	cur, err := c.NewGetMixBrokerCurrentTracesService(traderID, ProductTypeUSDTFutures).SetLimit(100).Do(cx)
	if err != nil {
		if !apitest.Tolerable(t, "mix-broker query-current-traces", err, brokerCodes...) {
			t.Fatalf("mix-broker query-current-traces: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-broker/query-current-traces", curParams, true)
		apitest.AssertCovers(t, "mix-broker query-current-traces", raw, cur)
	}
}
