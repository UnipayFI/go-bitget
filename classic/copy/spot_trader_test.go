package copy

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCopySpotTrader(t *testing.T) {
	c := NewCopyClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Capability/empty tolerance codes for the spot copy-trader product line:
	// sub-account / no-permission / empty-data / not-enrolled.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808", "400171", "400172", "40805", "31001", "31002", "40913", "40020", "43046", "43011", "43025", "40732", "40733", "40734"}

	// profit-summarys (object, no params)
	if resp, err := c.NewGetSpotTraderProfitSummarysService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/profit-summarys", err, okCodes...) {
			t.Fatalf("spot-trader/profit-summarys: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/profit-summarys", nil, true)
		apitest.AssertCovers(t, "spot-trader/profit-summarys", raw, resp)
	}

	// profit-history-details (object, optional params)
	if resp, err := c.NewGetSpotTraderProfitHistoryDetailsService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/profit-history-details", err, okCodes...) {
			t.Fatalf("spot-trader/profit-history-details: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/profit-history-details", nil, true)
		apitest.AssertCovers(t, "spot-trader/profit-history-details", raw, resp)
	}

	// profit-details (array, optional params)
	if resp, err := c.NewGetSpotTraderProfitDetailsService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/profit-details", err, okCodes...) {
			t.Fatalf("spot-trader/profit-details: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/profit-details", nil, true)
		apitest.AssertCovers(t, "spot-trader/profit-details", raw, resp)
	}

	// order-total-detail (object, no params)
	if resp, err := c.NewGetSpotTraderOrderTotalDetailService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/order-total-detail", err, okCodes...) {
			t.Fatalf("spot-trader/order-total-detail: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/order-total-detail", nil, true)
		apitest.AssertCovers(t, "spot-trader/order-total-detail", raw, resp)
	}

	// order-history-track (object, optional params)
	if resp, err := c.NewGetSpotTraderOrderHistoryTrackService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/order-history-track", err, okCodes...) {
			t.Fatalf("spot-trader/order-history-track: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/order-history-track", nil, true)
		apitest.AssertCovers(t, "spot-trader/order-history-track", raw, resp)
	}

	// order-current-track (object, optional params)
	if resp, err := c.NewGetSpotTraderOrderCurrentTrackService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/order-current-track", err, okCodes...) {
			t.Fatalf("spot-trader/order-current-track: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/order-current-track", nil, true)
		apitest.AssertCovers(t, "spot-trader/order-current-track", raw, resp)
	}

	// config-query-settings (object, no params)
	if resp, err := c.NewGetSpotTraderConfigSettingsService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/config-query-settings", err, okCodes...) {
			t.Fatalf("spot-trader/config-query-settings: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/config-query-settings", nil, true)
		apitest.AssertCovers(t, "spot-trader/config-query-settings", raw, resp)
	}

	// config-query-followers (array, optional params)
	if resp, err := c.NewGetSpotTraderConfigFollowersService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "spot-trader/config-query-followers", err, okCodes...) {
			t.Fatalf("spot-trader/config-query-followers: %v", err)
		}
	} else {
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/copy/spot-trader/config-query-followers", nil, true)
		apitest.AssertCovers(t, "spot-trader/config-query-followers", raw, resp)
	}
}
