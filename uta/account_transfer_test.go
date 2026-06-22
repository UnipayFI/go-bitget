package uta

import "testing"

func TestAccountTransfer(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Transferable coins.
	coins, err := c.NewGetTransferableCoinsService("spot", "usdt_futures").Do(cx)
	if err != nil {
		t.Fatalf("transferable coins: %v", err)
	}
	t.Logf("transferable coins: %d", len(coins))
	raw := fetchRawGet(t, c, cx, "/api/v3/account/transferable-coins",
		map[string]string{"fromType": "spot", "toType": "usdt_futures"}, true)
	assertCovers(t, "account/transferable-coins", raw, coins)

	// Main/sub transfer records.
	records, err := c.NewGetSubTransferRecordsService().Do(cx)
	if err != nil {
		t.Fatalf("sub transfer records: %v", err)
	}
	t.Logf("sub transfer records: %d cursor=%s", len(records.List), records.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/sub-transfer-record", nil, true)
	assertCovers(t, "account/sub-transfer-record", raw, records)
}
