package uta

import "testing"

// TestAccountCollateral covers the collateral-type, custom-collateral-coins, and
// pre-set-leverage read endpoints. Skipped without credentials.
func TestAccountCollateral(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Collateral type.
	collateral, err := c.NewGetCollateralTypeService().Do(cx)
	if err != nil {
		t.Fatalf("collateral-type: %v", err)
	}
	t.Logf("collateral type: %+v", collateral)
	raw := fetchRawGet(t, c, cx, "/api/v3/account/collateral-type", nil, true)
	assertCovers(t, "account/collateral-type", raw, collateral)

	// Custom collateral coins.
	coins, err := c.NewGetCustomCollateralCoinsService().Do(cx)
	if err != nil {
		t.Fatalf("custom-collateral-coins: %v", err)
	}
	t.Logf("custom collateral coins: %d", len(coins))
	raw = fetchRawGet(t, c, cx, "/api/v3/account/custom-collateral-coins", nil, true)
	assertCovers(t, "account/custom-collateral-coins", raw, coins)

	// Pre-set leverage preview.
	params := map[string]string{
		"category":   string(CategoryUSDTFutures),
		"marginMode": string(MarginModeCrossed),
		"symbol":     "BTCUSDT",
		"leverage":   "20",
	}
	preview, err := c.NewPreSetLeverageService(CategoryUSDTFutures, MarginModeCrossed).
		SetSymbol("BTCUSDT").SetLeverage("20").Do(cx)
	if err != nil {
		t.Fatalf("pre-set-leverage: %v", err)
	}
	t.Logf("pre-set leverage: %+v", preview)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/pre-set-leverage", params, true)
	assertCovers(t, "account/pre-set-leverage", raw, preview)
}
