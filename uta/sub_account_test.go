package uta

import "testing"

func TestSubAccount(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Sub-account list (no required params; works without pre-existing state).
	list, err := c.NewGetSubAccountListService().Do(cx)
	if err != nil {
		t.Fatalf("sub-list: %v", err)
	}
	t.Logf("sub-list: %d hasNext=%v cursor=%s", len(list.List), list.HasNext, list.Cursor)
	raw := fetchRawGet(t, c, cx, "/api/v3/user/sub-list", nil, true)
	assertCovers(t, "user/sub-list", raw, list)

	// Sub-account unified assets (subUid omitted -> all sub-accounts).
	assets, err := c.NewGetSubAccountUnifiedAssetsService().Do(cx)
	if err != nil {
		t.Fatalf("sub-unified-assets: %v", err)
	}
	t.Logf("sub-unified-assets: %d sub-account(s)", len(assets))
	raw = fetchRawGet(t, c, cx, "/api/v3/account/sub-unified-assets", nil, true)
	assertCovers(t, "account/sub-unified-assets", raw, assets)
}
