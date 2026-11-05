package uta

import "testing"

func TestAffiliate(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Direct customer list (gates the rest: a non-affiliate key is a skip).
	customers, err := c.NewGetAffiliateCustomerListService().Do(cx)
	if err != nil {
		if tolerable(t, "affiliate/customer-list", err) {
			t.Skip("account is not an affiliate")
		}
		t.Fatalf("affiliate customer-list: %v", err)
	}
	t.Logf("affiliate customer-list: %d", len(customers))
	raw := fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-list", nil, true)
	assertCovers(t, "affiliate/customer-list", raw, customers)

	assets, err := c.NewGetAffiliateCustomerAssetService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate customer-asset: %v", err)
	}
	t.Logf("affiliate customer-asset: %d", len(assets))
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-asset", nil, true)
	assertCovers(t, "affiliate/customer-asset", raw, assets)

	// Defaults to yesterday (UTC+8) without a time range.
	volumes, err := c.NewGetAffiliateCustomerTradeVolumeService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate customer-trade-volume: %v", err)
	}
	t.Logf("affiliate customer-trade-volume: %d", len(volumes))
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-trade-volume", nil, true)
	assertCovers(t, "affiliate/customer-trade-volume", raw, volumes)

	deposits, err := c.NewGetAffiliateCustomerDepositService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate customer-deposit: %v", err)
	}
	t.Logf("affiliate customer-deposit: %d", len(deposits))
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-deposit", nil, true)
	assertCovers(t, "affiliate/customer-deposit", raw, deposits)

	kyc, err := c.NewGetAffiliateCustomerKYCResultService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate customer-kyc-result: %v", err)
	}
	t.Logf("affiliate customer-kyc-result: %d endId=%s", len(kyc.UserList), kyc.EndID)
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-kyc-result", nil, true)
	assertCovers(t, "affiliate/customer-kyc-result", raw, kyc)

	comm, err := c.NewGetAffiliateCustomerCommissionService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate customer-commission: %v", err)
	}
	t.Logf("affiliate customer-commission: %d endId=%s", len(comm.CommissionList), comm.EndID)
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/customer-commission", nil, true)
	assertCovers(t, "affiliate/customer-commission", raw, comm)

	subCustomers, err := c.NewGetAffiliateSubCustomerListService().Do(cx)
	if err != nil {
		t.Fatalf("affiliate sub-customer-list: %v", err)
	}
	t.Logf("affiliate sub-customer-list: %d minId=%s", len(subCustomers.List), subCustomers.MinID)
	raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/sub-customer-list", nil, true)
	assertCovers(t, "affiliate/sub-customer-list", raw, subCustomers)

	// Sub-affiliate views need a partner with downline affiliates; tolerate a
	// capability error there rather than failing the whole test.
	subInfo, err := c.NewGetSubAffiliateInfoService().Do(cx)
	switch {
	case err == nil:
		t.Logf("affiliate sub-affiliate-info: %d endId=%s", len(subInfo.List), subInfo.EndID)
		raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/sub-affiliate-info", nil, true)
		assertCovers(t, "affiliate/sub-affiliate-info", raw, subInfo)
	case !tolerable(t, "affiliate/sub-affiliate-info", err):
		t.Fatalf("affiliate sub-affiliate-info: %v", err)
	}

	subComm, err := c.NewGetSubAffiliateCommissionService().Do(cx)
	switch {
	case err == nil:
		t.Logf("affiliate sub-affiliate-commission: %d endId=%s", len(subComm.CommissionList), subComm.EndID)
		raw = fetchRawGet(t, c, cx, "/api/v3/affiliate/sub-affiliate-commission", nil, true)
		assertCovers(t, "affiliate/sub-affiliate-commission", raw, subComm)
	case !tolerable(t, "affiliate/sub-affiliate-commission", err):
		t.Fatalf("affiliate sub-affiliate-commission: %v", err)
	}
}
