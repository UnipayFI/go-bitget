package uta

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

// TestOrderFill validates market orders and the Fill struct against a real
// execution: a tiny ~2 USDT SPOT market buy of BTC, then a sell-back of the
// resulting balance to restore the account (minus negligible fees). Gated
// behind BITGET_TEST_WRITE=1. Per the project rule, fills are tested on a big
// coin (BTC) with a minimal amount.
func TestOrderFill(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run live order tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	const symbol = "BTCUSDT"

	// 1) Market BUY — for a spot market buy, qty is the quote amount (USDT).
	ref, err := c.NewPlaceOrderService(CategorySpot, symbol, decimal.RequireFromString("2"), SideBuy, OrderTypeMarket).Do(cx)
	if err != nil {
		t.Fatalf("market buy: %v", err)
	}
	t.Logf("market buy orderId=%s", ref.OrderID)

	// 2) order-info: confirm it filled.
	order, err := c.NewGetOrderInfoService().SetOrderID(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("order-info: %v", err)
	}
	t.Logf("buy: status=%s cumExecQty=%s avgPrice=%s cumExecValue=%s", order.OrderStatus, order.CumExecQty, order.AvgPrice, order.CumExecValue)

	// 3) fills by orderId — validate the Fill struct against the real execution.
	fills, err := c.NewGetFillHistoryService().SetCategory(CategorySpot).SetOrderID(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("fills: %v", err)
	}
	t.Logf("fills: %d", len(fills.List))
	for _, f := range fills.List {
		t.Logf("  fill execId=%s price=%s qty=%s scope=%s fee=%v", f.ExecID, f.ExecPrice, f.ExecQty, f.TradeScope, f.FeeDetail)
	}
	raw := fetchRawGet(t, c, cx, "/api/v3/trade/fills",
		map[string]string{"category": string(CategorySpot), "orderId": ref.OrderID}, true)
	assertCovers(t, "trade/fills", raw, fills)

	// 4) Sell back the BTC we now hold to restore the account. Use the live
	// available balance (net of fee) truncated to the symbol's precision.
	assets, err := c.NewGetAccountAssetsService().Do(cx)
	if err != nil {
		t.Fatalf("assets: %v", err)
	}
	var btc decimal.Decimal
	for _, a := range assets.Assets {
		if a.Coin == "BTC" {
			btc = a.Available
		}
	}
	sellQty := btc.Truncate(6)
	t.Logf("BTC available=%s, selling back %s", btc, sellQty)
	if sellQty.IsZero() {
		t.Skip("no BTC available to sell back (nothing to clean up)")
	}
	sref, err := c.NewPlaceOrderService(CategorySpot, symbol, sellQty, SideSell, OrderTypeMarket).Do(cx)
	if err != nil {
		t.Logf("sell-back failed (residual BTC left, harmless): %v", err)
		return
	}
	t.Logf("sold back orderId=%s", sref.OrderID)
}
