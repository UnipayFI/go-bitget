package uta

import "testing"

func TestSmallAssets(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Convertible dust balances.
	assets, err := c.NewGetSmallAssetsService().Do(cx)
	if err != nil {
		t.Fatalf("small assets: %v", err)
	}
	t.Logf("small assets: %d", len(assets))
	for _, a := range assets {
		t.Logf("  %s available=%s -> %s %s fee=%s",
			a.Coin, a.Available, a.EstimatedAmount, a.EstimatedCoin, a.FeeDetail.Fee)
	}
	raw := fetchRawGet(t, c, cx, "/api/v3/convert/small-assets", nil, true)
	assertCovers(t, "convert/small-assets", raw, assets)

	// Conversion history.
	history, err := c.NewGetSmallAssetsHistoryService().SetLimit(20).Do(cx)
	if err != nil {
		t.Fatalf("small assets history: %v", err)
	}
	t.Logf("small assets history: %d cursor=%s", len(history.List), history.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/convert/small-assets-history",
		map[string]string{"limit": "20"}, true)
	assertCovers(t, "convert/small-assets-history", raw, history)
}
