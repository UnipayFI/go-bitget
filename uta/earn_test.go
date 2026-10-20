package uta

import "testing"

func TestEarn(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Elite products (no params).
	products, err := c.NewGetEliteProductService().Do(cx)
	if err != nil {
		t.Fatalf("elite product: %v", err)
	}
	t.Logf("elite products: %d", len(products))
	raw := fetchRawGet(t, c, cx, "/api/v3/earn/elite-product", nil, true)
	assertCovers(t, "earn/elite-product", raw, products)

	// Elite assets (no params).
	assets, err := c.NewGetEliteAssetsService().Do(cx)
	if err != nil {
		t.Fatalf("elite assets: %v", err)
	}
	t.Logf("elite assets: %d", len(assets.ResultList))
	raw = fetchRawGet(t, c, cx, "/api/v3/earn/elite-assets", nil, true)
	assertCovers(t, "earn/elite-assets", raw, assets)

	// Elite records (type required; subscribe history may be empty).
	records, err := c.NewGetEliteRecordsService(EliteRecordTypeSubscribe).Do(cx)
	if err != nil {
		t.Fatalf("elite records: %v", err)
	}
	t.Logf("elite records: %d endId=%s", len(records.RecordList), records.EndID)
	raw = fetchRawGet(t, c, cx, "/api/v3/earn/elite-records",
		map[string]string{"type": string(EliteRecordTypeSubscribe)}, true)
	assertCovers(t, "earn/elite-records", raw, records)

	// Subscribe info and redeem info both require a productId; derive one from
	// the product catalogue when available, otherwise skip (read-only either way).
	if len(products) > 0 {
		productId := products[0].ProductID

		info, err := c.NewGetEliteSubscribeInfoService(productId).Do(cx)
		if err != nil {
			t.Fatalf("elite subscribe info: %v", err)
		}
		t.Logf("subscribe info: productSubId=%s minAmount=%s", info.ProductSubID, info.MinAmount)
		raw = fetchRawGet(t, c, cx, "/api/v3/earn/elite-subscribe-info",
			map[string]string{"productId": productId}, true)
		assertCovers(t, "earn/elite-subscribe-info", raw, info)

		redeemInfo, err := c.NewGetRedeemInfoService(productId).Do(cx)
		if err != nil {
			t.Fatalf("elite redeem info: %v", err)
		}
		t.Logf("redeem info: productSubId=%s modes=%d", redeemInfo.ProductSubID, len(redeemInfo.RedeemModeList))
		raw = fetchRawGet(t, c, cx, "/api/v3/earn/elite-redeem-info",
			map[string]string{"productId": productId}, true)
		assertCovers(t, "earn/elite-redeem-info", raw, redeemInfo)
	}
}

func TestCashPlus(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	subscribeInfo, err := c.NewGetCashPlusSubscribeInfoService().Do(cx)
	if err != nil {
		t.Fatalf("cash+ subscribe info: %v", err)
	}
	t.Logf("cash+ subscribe info: apr=%s coins=%d", subscribeInfo.CashPlusApr, len(subscribeInfo.SubscribeCoinList))
	raw := fetchRawGet(t, c, cx, "/api/v3/earn/cashplus/subscribe-info", nil, true)
	assertCovers(t, "earn/cashplus/subscribe-info", raw, subscribeInfo)

	redeemInfo, err := c.NewGetCashPlusRedeemInfoService().Do(cx)
	if err != nil {
		t.Fatalf("cash+ redeem info: %v", err)
	}
	t.Logf("cash+ redeem info: coins=%d quota=%s", len(redeemInfo.RedeemCoinList), redeemInfo.UserRemainingQuota)
	raw = fetchRawGet(t, c, cx, "/api/v3/earn/cashplus/redeem-info", nil, true)
	assertCovers(t, "earn/cashplus/redeem-info", raw, redeemInfo)

	positions, err := c.NewGetCashPlusPositionsService().Do(cx)
	if err != nil {
		t.Fatalf("cash+ positions: %v", err)
	}
	t.Logf("cash+ positions: balance=%s totalProfit=%s", positions.Balance, positions.TotalProfit)
	raw = fetchRawGet(t, c, cx, "/api/v3/earn/cashplus/positions", nil, true)
	assertCovers(t, "earn/cashplus/positions", raw, positions)

	// History may be empty for an account that never used Cash+.
	history, err := c.NewGetCashPlusHistoryOrdersService().Do(cx)
	if err != nil {
		t.Fatalf("cash+ history orders: %v", err)
	}
	t.Logf("cash+ history orders: %d nextCursor=%s", len(history.Items), history.NextCursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/earn/cashplus/history-orders", nil, true)
	assertCovers(t, "earn/cashplus/history-orders", raw, history)
}
