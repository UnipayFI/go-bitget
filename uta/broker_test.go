package uta

import "testing"

func TestBroker(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Sub-account list.
	subs, err := c.NewGetBrokerSubListService().Do(cx)
	if err != nil {
		if tolerable(t, "broker/sub-list", err) {
			t.Skip("account is not a broker")
		}
		t.Fatalf("broker sub-list: %v", err)
	}
	t.Logf("broker sub-list: %d", len(subs.SubList))
	raw := fetchRawGet(t, c, cx, "/api/v3/broker/sub-list", nil, true)
	assertCovers(t, "broker/sub-list", raw, subs)

	// All sub-account deposit/withdrawal records.
	dw, err := c.NewGetAllBrokerSubDepositWithdrawalService().Do(cx)
	if err != nil {
		t.Fatalf("broker all-sub-deposit-withdrawal: %v", err)
	}
	t.Logf("broker all-sub-deposit-withdrawal: %d endId=%s", len(dw.List), dw.EndId)
	raw = fetchRawGet(t, c, cx, "/api/v3/broker/all-sub-deposit-withdrawal", nil, true)
	assertCovers(t, "broker/all-sub-deposit-withdrawal", raw, dw)

	// Commission records.
	comm, err := c.NewGetBrokerCommissionService().Do(cx)
	if err != nil {
		t.Fatalf("broker commission: %v", err)
	}
	t.Logf("broker commission: %d", len(comm))
	raw = fetchRawGet(t, c, cx, "/api/v3/broker/commission", nil, true)
	assertCovers(t, "broker/commission", raw, comm)
}
