package uta

import "testing"

func TestP2P(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Supported currencies.
	currencies, err := c.NewGetP2PCurrenciesService().Do(cx)
	if err != nil {
		if tolerable(t, "p2p/currencies", err) {
			t.Skip("account is not registered as a P2P user")
		}
		t.Fatalf("p2p currencies: %v", err)
	}
	t.Logf("p2p currencies: tokens=%d fiats=%d", len(currencies.TokenDetailList), len(currencies.FiatDetailList))
	raw := fetchRawGet(t, c, cx, "/api/v3/p2p/currencies", nil, true)
	assertCovers(t, "p2p/currencies", raw, currencies)

	// Caller's payment methods.
	payMethods, err := c.NewGetP2PPayMethodsService().Do(cx)
	if err != nil {
		t.Fatalf("p2p pay methods: %v", err)
	}
	t.Logf("p2p pay methods: %d", len(payMethods))
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/pay-method", nil, true)
	assertCovers(t, "p2p/pay-method", raw, payMethods)

	// Caller's merchant profile.
	userInfo, err := c.NewGetP2PUserInfoService().Do(cx)
	if err != nil {
		t.Fatalf("p2p user info: %v", err)
	}
	t.Logf("p2p user info: uid=%s level=%s", userInfo.UID, userInfo.AccountLevel)
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/user-info", nil, true)
	assertCovers(t, "p2p/user-info", raw, userInfo)

	// P2P balance for USDT.
	balance, err := c.NewGetP2PBalanceService("USDT").Do(cx)
	if err != nil {
		t.Fatalf("p2p balance: %v", err)
	}
	t.Logf("p2p balance: token=%s available=%s", balance.Token, balance.AvailableBalance)
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/balance", map[string]string{"token": "USDT"}, true)
	assertCovers(t, "p2p/balance", raw, balance)

	// Caller's own ads.
	myAds, err := c.NewGetP2PMyAdsService().Do(cx)
	if err != nil {
		t.Fatalf("p2p my ads: %v", err)
	}
	t.Logf("p2p my ads: %d nextId=%s", len(myAds.Items), myAds.NextID)
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/my-ads", nil, true)
	assertCovers(t, "p2p/my-ads", raw, myAds)

	// Pending orders.
	pending, err := c.NewGetP2PPendingOrdersService().Do(cx)
	if err != nil {
		t.Fatalf("p2p pending orders: %v", err)
	}
	t.Logf("p2p pending orders: %d nextId=%s", len(pending.Items), pending.NextID)
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/pending-orders", nil, true)
	assertCovers(t, "p2p/pending-orders", raw, pending)

	// All orders.
	allOrders, err := c.NewGetP2PAllOrdersService().Do(cx)
	if err != nil {
		t.Fatalf("p2p all orders: %v", err)
	}
	t.Logf("p2p all orders: %d nextId=%s", len(allOrders.Items), allOrders.NextID)
	raw = fetchRawGet(t, c, cx, "/api/v3/p2p/all-orders", nil, true)
	assertCovers(t, "p2p/all-orders", raw, allOrders)
}
