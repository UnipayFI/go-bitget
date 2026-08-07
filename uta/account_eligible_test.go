package uta

import "testing"

// The eligible-* endpoints are whitelist-gated: a non-whitelisted account gets
// an empty list (or 40054), which still exercises the request path + signing.
func TestAccountEligible(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Eligible symbols.
	if symbols, err := c.NewGetEligibleSymbolsService().Do(cx); err != nil {
		if !tolerable(t, "account/eligible-symbols", err) {
			t.Fatalf("eligible symbols: %v", err)
		}
	} else {
		t.Logf("eligible symbols: %d", len(symbols))
		raw := fetchRawGet(t, c, cx, "/api/v3/account/eligible-symbols", nil, true)
		assertCovers(t, "account/eligible-symbols", raw, symbols)
	}

	// Eligible margin tier.
	if tiers, err := c.NewGetEligibleMarginTierService().Do(cx); err != nil {
		if !tolerable(t, "account/eligible-margin-tier", err) {
			t.Fatalf("eligible margin tier: %v", err)
		}
	} else {
		t.Logf("eligible margin tiers: %d", len(tiers))
		raw := fetchRawGet(t, c, cx, "/api/v3/account/eligible-margin-tier", nil, true)
		assertCovers(t, "account/eligible-margin-tier", raw, tiers)
	}

	// Eligible loan info.
	if loans, err := c.NewGetEligibleLoanInfoService().Do(cx); err != nil {
		if !tolerable(t, "account/eligible-loan-info", err) {
			t.Fatalf("eligible loan info: %v", err)
		}
	} else {
		t.Logf("eligible loan info: %d", len(loans))
		raw := fetchRawGet(t, c, cx, "/api/v3/account/eligible-loan-info", nil, true)
		assertCovers(t, "account/eligible-loan-info", raw, loans)
	}
}
