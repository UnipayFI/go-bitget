package uta

import (
	"testing"

	"github.com/UnipayFI/go-bitget/common"
)

func TestStockPlusMaxPurchaseUnmarshal(t *testing.T) {
	const payload = `{"cashMaxQty":"100","marginMaxQty":""}`
	var got StockPlusMaxPurchase
	if err := common.JSONUnmarshal([]byte(payload), &got); err != nil {
		t.Fatal(err)
	}
	if got.CashMaxQty.String() != "100" || !got.MarginMaxQty.IsZero() {
		t.Fatalf("got cash=%s margin=%s, want 100/0", got.CashMaxQty, got.MarginMaxQty)
	}
}

// TestStockPlusMaxPurchase needs a whitelisted UID with a Stock+ securities
// account, so it skips when the endpoint refuses the caller.
func TestStockPlusMaxPurchase(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	const symbol = "AAPL.US"
	maxPurchase, err := c.NewGetStockPlusMaxPurchaseService(symbol, StockPlusOrderTypeMarket, StockPlusSideBuy).Do(cx)
	if err != nil {
		t.Skipf("Stock+ not available on this UID (%v); skipping", err)
	}
	t.Logf("stockplus max-purchase: %+v", maxPurchase)
	raw := fetchRawGet(t, c, cx, "/api/v3/stockplus/trade/max-purchase", map[string]string{
		"symbol":    symbol,
		"orderType": string(StockPlusOrderTypeMarket),
		"side":      string(StockPlusSideBuy),
	}, true)
	assertCovers(t, "stockplus/trade/max-purchase", raw, maxPurchase)
}
