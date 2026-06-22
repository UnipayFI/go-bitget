package uta

import "testing"

func TestAccountRecords(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Financial records.
	fin, err := c.NewGetFinancialRecordsService(CategoryUSDTFutures).SetCoin("USDT").Do(cx)
	if err != nil {
		t.Fatalf("financial records: %v", err)
	}
	t.Logf("financial records: %d cursor=%s", len(fin.List), fin.Cursor)
	raw := fetchRawGet(t, c, cx, "/api/v3/account/financial-records",
		map[string]string{"category": string(CategoryUSDTFutures), "coin": "USDT"}, true)
	assertCovers(t, "account/financial-records", raw, fin)

	// Convert records.
	conv, err := c.NewGetConvertRecordsService("USDT", "USDC").Do(cx)
	if err != nil {
		t.Fatalf("convert records: %v", err)
	}
	t.Logf("convert records: %d cursor=%s", len(conv.List), conv.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/convert-records",
		map[string]string{"fromCoin": "USDT", "toCoin": "USDC"}, true)
	assertCovers(t, "account/convert-records", raw, conv)

	// Repayable coins.
	repayable, err := c.NewGetRepayableCoinsService().Do(cx)
	if err != nil {
		t.Fatalf("repayable coins: %v", err)
	}
	t.Logf("repayable coins: %d maxSelection=%s", len(repayable.RepayableCoinList), repayable.MaxSelection)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/repayable-coins", nil, true)
	assertCovers(t, "account/repayable-coins", raw, repayable)

	// Payment coins.
	payment, err := c.NewGetPaymentCoinsService().Do(cx)
	if err != nil {
		t.Fatalf("payment coins: %v", err)
	}
	t.Logf("payment coins: %d maxSelection=%s", len(payment.PaymentCoinList), payment.MaxSelection)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/payment-coins", nil, true)
	assertCovers(t, "account/payment-coins", raw, payment)
}
