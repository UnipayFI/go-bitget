package uta

import "testing"

func TestAccountFunding(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	coin := "USDT"

	// Funding account assets.
	assets, err := c.NewGetAccountFundingAssetsService().SetCoin(coin).Do(cx)
	if err != nil {
		t.Fatalf("funding assets: %v", err)
	}
	t.Logf("funding assets: %d", len(assets))
	for _, a := range assets {
		t.Logf("  %s balance=%s available=%s frozen=%s", a.Coin, a.Balance, a.Available, a.Frozen)
	}
	raw := fetchRawGet(t, c, cx, "/api/v3/account/funding-assets", map[string]string{"coin": coin}, true)
	assertCovers(t, "account/funding-assets", raw, assets)

	// Max transferable.
	transfer, err := c.NewGetMaxTransferableService(coin).Do(cx)
	if err != nil {
		t.Fatalf("max transferable: %v", err)
	}
	t.Logf("max transferable: %+v", transfer)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/max-transferable", map[string]string{"coin": coin}, true)
	assertCovers(t, "account/max-transferable", raw, transfer)

	// Max withdrawal.
	withdrawal, err := c.NewGetMaxWithdrawalService(coin).Do(cx)
	if err != nil {
		t.Fatalf("max withdrawal: %v", err)
	}
	t.Logf("max withdrawal: %+v", withdrawal)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/max-withdrawal", map[string]string{"coin": coin}, true)
	assertCovers(t, "account/max-withdrawal", raw, withdrawal)
}
