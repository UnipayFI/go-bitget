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
	t.Logf("broker all-sub-deposit-withdrawal: %d endId=%s", len(dw.List), dw.EndID)
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

	// Daily totals and per-fill commission (both default to yesterday UTC+8).
	total, err := c.NewGetBrokerTotalCommissionService().Do(cx)
	if err != nil {
		t.Fatalf("broker total-commission: %v", err)
	}
	t.Logf("broker total-commission: %d", len(total))
	raw = fetchRawGet(t, c, cx, "/api/v3/broker/total-commission", nil, true)
	assertCovers(t, "broker/total-commission", raw, total)

	orderComm, err := c.NewGetBrokerOrderCommissionService().Do(cx)
	if err != nil {
		t.Fatalf("broker order-commission: %v", err)
	}
	t.Logf("broker order-commission: %d endId=%s", len(orderComm.CommissionList), orderComm.EndID)
	raw = fetchRawGet(t, c, cx, "/api/v3/broker/order-commission", nil, true)
	assertCovers(t, "broker/order-commission", raw, orderComm)

	// Trade volume needs separate access granted by BD/RM.
	volume, err := c.NewGetBrokerTradeVolumeService().Do(cx)
	switch {
	case err == nil:
		t.Logf("broker trade-volume: %d", len(volume))
		raw = fetchRawGet(t, c, cx, "/api/v3/broker/trade-volume", nil, true)
		assertCovers(t, "broker/trade-volume", raw, volume)
	case !tolerable(t, "broker/trade-volume", err):
		t.Fatalf("broker trade-volume: %v", err)
	}

	// Email and rebate info need a sub-account UID.
	if len(subs.SubList) > 0 {
		subUID := subs.SubList[0].SubUID

		email, err := c.NewGetBrokerSubEmailService(subUID).Do(cx)
		if err != nil {
			t.Fatalf("broker subaccount-email: %v", err)
		}
		t.Logf("broker subaccount-email: subUid=%s", email.SubUID)
		raw = fetchRawGet(t, c, cx, "/api/v3/broker/subaccount-email", map[string]string{"subUid": subUID}, true)
		assertCovers(t, "broker/subaccount-email", raw, email)

		rebate, err := c.NewGetBrokerRebateInfoService(subUID).Do(cx)
		if err != nil {
			t.Fatalf("broker rebate-info: %v", err)
		}
		t.Logf("broker rebate-info: %s %s", rebate.AffiliationType, rebate.UserLevel)
		raw = fetchRawGet(t, c, cx, "/api/v3/broker/rebate-info", map[string]string{"uid": subUID}, true)
		assertCovers(t, "broker/rebate-info", raw, rebate)
	}
}
