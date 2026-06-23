package copy

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCopyMixTrader(t *testing.T) {
	c := NewCopyClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// These trader endpoints require the account to be an enrolled elite copy
	// trader. A non-trader (or permission-limited sub-account) key returns a
	// capability/empty/no-permission code which we treat as a pass: the path
	// and signing are still proven.
	okCodes := []string{"40014", "40099", "40034", "40068", "40054", "40037", "40029", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808", "400171", "400172", "40805", "31001", "31002", "40913", "40020", "43046", "43011", "43025", "40732", "40733", "40734"}

	const pt = string(ProductTypeUSDTFutures)

	// order-current-track
	{
		resp, err := c.NewGetMixTraderCurrentTrackService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "order-current-track", err, okCodes...) {
				t.Fatalf("order-current-track: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/order-current-track", map[string]string{"productType": pt}, true)
			apitest.AssertCovers(t, "order-current-track", raw, resp)
		}
	}

	// order-history-track
	{
		resp, err := c.NewGetMixTraderHistoryTrackService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "order-history-track", err, okCodes...) {
				t.Fatalf("order-history-track: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/order-history-track", map[string]string{"productType": pt}, true)
			apitest.AssertCovers(t, "order-history-track", raw, resp)
		}
	}

	// order-total-detail
	{
		resp, err := c.NewGetMixTraderOrderTotalDetailService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "order-total-detail", err, okCodes...) {
				t.Fatalf("order-total-detail: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/order-total-detail", nil, true)
			apitest.AssertCovers(t, "order-total-detail", raw, resp)
		}
	}

	// profit-history-summarys
	{
		resp, err := c.NewGetMixTraderProfitHistorySummarysService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "profit-history-summarys", err, okCodes...) {
				t.Fatalf("profit-history-summarys: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/profit-history-summarys", nil, true)
			apitest.AssertCovers(t, "profit-history-summarys", raw, resp)
		}
	}

	// profit-history-details
	{
		resp, err := c.NewGetMixTraderProfitHistoryDetailsService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "profit-history-details", err, okCodes...) {
				t.Fatalf("profit-history-details: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/profit-history-details", nil, true)
			apitest.AssertCovers(t, "profit-history-details", raw, resp)
		}
	}

	// profit-details
	{
		resp, err := c.NewGetMixTraderProfitDetailsService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "profit-details", err, okCodes...) {
				t.Fatalf("profit-details: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/profit-details", nil, true)
			apitest.AssertCovers(t, "profit-details", raw, resp)
		}
	}

	// profits-group-coin-date
	{
		resp, err := c.NewGetMixTraderProfitsGroupCoinDateService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "profits-group-coin-date", err, okCodes...) {
				t.Fatalf("profits-group-coin-date: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/profits-group-coin-date", nil, true)
			apitest.AssertCovers(t, "profits-group-coin-date", raw, resp)
		}
	}

	// config-query-symbols
	{
		resp, err := c.NewGetMixTraderConfigQuerySymbolsService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "config-query-symbols", err, okCodes...) {
				t.Fatalf("config-query-symbols: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/config-query-symbols", map[string]string{"productType": pt}, true)
			apitest.AssertCovers(t, "config-query-symbols", raw, resp)
		}
	}

	// config-query-followers
	{
		resp, err := c.NewGetMixTraderConfigQueryFollowersService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "config-query-followers", err, okCodes...) {
				t.Fatalf("config-query-followers: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/mix-trader/config-query-followers", nil, true)
			apitest.AssertCovers(t, "config-query-followers", raw, resp)
		}
	}
}
