package earn

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// savingsTolerable lists the Bitget response codes that mean the shared
// permission-limited sub-account key cannot access this product line or simply
// has no savings data — the endpoint path and signing are still proven.
var savingsTolerable = []string{"40099", "40034", "40054", "40068", "40014", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

func TestEarnSavings(t *testing.T) {
	c := NewEarnClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Savings Product List (flexible) — needs a holding filter; "all" lists everything.
	var productID string
	{
		resp, err := c.NewGetSavingsProductService().SetFilter(SavingsProductFilterAll).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "savings/product", err, savingsTolerable...) {
				t.Fatalf("savings/product: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/savings/product", map[string]string{"filter": "all"}, true)
			apitest.AssertCovers(t, "savings/product", raw, resp)
			for _, p := range resp {
				if p.PeriodType == SavingsPeriodTypeFlexible && productID == "" {
					productID = p.ProductID
				}
			}
		}
	}

	// Savings Account (no params).
	{
		resp, err := c.NewGetSavingsAccountService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "savings/account", err, savingsTolerable...) {
				t.Fatalf("savings/account: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/savings/account", nil, true)
			apitest.AssertCovers(t, "savings/account", raw, resp)
		}
	}

	// Savings Assets (flexible).
	{
		params := map[string]string{"periodType": string(SavingsPeriodTypeFlexible)}
		resp, err := c.NewGetSavingsAssetsService(SavingsPeriodTypeFlexible).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "savings/assets", err, savingsTolerable...) {
				t.Fatalf("savings/assets: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/savings/assets", params, true)
			apitest.AssertCovers(t, "savings/assets", raw, resp)
		}
	}

	// Savings Records (flexible).
	{
		params := map[string]string{"periodType": string(SavingsPeriodTypeFlexible)}
		resp, err := c.NewGetSavingsRecordsService(SavingsPeriodTypeFlexible).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "savings/records", err, savingsTolerable...) {
				t.Fatalf("savings/records: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/savings/records", params, true)
			apitest.AssertCovers(t, "savings/records", raw, resp)
		}
	}

	// Savings Subscription Detail — needs a productId from the product list above.
	if productID != "" {
		params := map[string]string{"productId": productID, "periodType": string(SavingsPeriodTypeFlexible)}
		resp, err := c.NewGetSavingsSubscribeInfoService(productID, SavingsPeriodTypeFlexible).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "savings/subscribe-info", err, savingsTolerable...) {
				t.Fatalf("savings/subscribe-info: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/savings/subscribe-info", params, true)
			apitest.AssertCovers(t, "savings/subscribe-info", raw, resp)
		}
	} else {
		t.Log("savings/subscribe-info: no flexible productId available; skipping")
	}

	// Savings Subscription Result — keyed by the orderId returned from Subscribe
	// (Bitget's param table mislabels it productId; the live API + the doc's own
	// request example use orderId). No real subscription order exists to query,
	// so this is skipped; a bogus orderId returns 43111 which confirms the param
	// name is correct.
	t.Log("savings/subscribe-result: needs a real subscription orderId; skipped (state would have to be created)")

	// Savings Redemption Results — requires a real subscription orderId we cannot
	// safely create; exercising it without one would 4xx on a missing resource.
	t.Log("savings/redeem-result: needs a real subscription orderId; skipped (state would have to be created)")
}
