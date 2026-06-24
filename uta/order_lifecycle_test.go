package uta

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

// TestOrderLifecycle exercises the full single-order write path against the live
// account with a tiny, far-from-market SPOT limit order that rests unfilled and
// is then cancelled — so it is fully reversible and costs nothing but the
// briefly-reserved ~$1.3. Gated behind BITGET_TEST_WRITE=1.
//
// Flow: place (post_only, ~50% below market) -> order-info (validates the Order
// struct against a real order) -> unfilled-orders -> modify price -> order-info
// -> cancel -> history-orders.
func TestOrderLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run live order tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	const symbol = "BTCUSDT"
	clientOid := "gobitget-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	price := decimal.RequireFromString("32000") // far below market; rests as maker
	qty := decimal.RequireFromString("0.00004") // ~1.28 USDT notional, above $1 min

	// 1) Place a resting post_only limit buy.
	ref, err := c.NewPlaceOrderService(CategorySpot, symbol, qty, SideBuy, OrderTypeLimit).
		SetPrice(price).
		SetTimeInForce(TimeInForcePostOnly).
		SetClientOrderID(clientOid).
		Do(cx)
	if err != nil {
		t.Fatalf("place-order: %v", err)
	}
	t.Logf("placed orderId=%s clientOid=%s", ref.OrderID, ref.ClientOrderID)
	if ref.OrderID == "" {
		t.Fatal("place-order returned empty orderId")
	}

	// Ensure we always clean up even if a later step fails.
	defer func() {
		_, _ = c.NewCancelOrderService().SetOrderID(ref.OrderID).SetCategory(CategorySpot).Do(ctx(t))
	}()

	// 2) order-info by orderId — validate the Order struct against the real order.
	order, err := c.NewGetOrderInfoService().SetOrderID(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("order-info: %v", err)
	}
	t.Logf("order: status=%s price=%s qty=%s side=%s type=%s", order.OrderStatus, order.Price, order.Qty, order.Side, order.OrderType)
	raw := fetchRawGet(t, c, cx, "/api/v3/trade/order-info", map[string]string{"orderId": ref.OrderID}, true)
	assertCovers(t, "trade/order-info", raw, order)

	// 3) unfilled-orders should include it.
	open, err := c.NewGetOpenOrdersService().SetCategory(CategorySpot).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("unfilled-orders: %v", err)
	}
	found := false
	for _, o := range open.List {
		if o.OrderID == ref.OrderID {
			found = true
		}
	}
	t.Logf("unfilled-orders: %d (found ours=%v)", len(open.List), found)
	if !found {
		t.Errorf("placed order %s not present in unfilled-orders", ref.OrderID)
	}

	// 4) modify the price (still far below market).
	newPrice := decimal.RequireFromString("31000")
	if _, err := c.NewModifyOrderService().
		SetOrderID(ref.OrderID).
		SetSymbol(symbol).
		SetCategory(CategorySpot).
		SetPrice(newPrice).
		Do(cx); err != nil {
		t.Fatalf("modify-order: %v", err)
	}
	t.Logf("modified price -> %s", newPrice)

	// 5) order-info reflects the new price.
	order2, err := c.NewGetOrderInfoService().SetOrderID(ref.OrderID).Do(cx)
	if err != nil {
		t.Fatalf("order-info(2): %v", err)
	}
	t.Logf("after modify: price=%s status=%s", order2.Price, order2.OrderStatus)
	if !order2.Price.Equal(newPrice) {
		t.Errorf("modify did not take effect: price=%s want=%s", order2.Price, newPrice)
	}

	// 6) cancel.
	if _, err := c.NewCancelOrderService().SetOrderID(ref.OrderID).SetCategory(CategorySpot).Do(cx); err != nil {
		t.Fatalf("cancel-order: %v", err)
	}
	t.Logf("cancelled %s", ref.OrderID)

	// 7) history-orders should now include the cancelled order.
	hist, err := c.NewGetOrderHistoryService(CategorySpot).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("history-orders: %v", err)
	}
	inHist := false
	for _, o := range hist.List {
		if o.OrderID == ref.OrderID {
			inHist = true
		}
	}
	t.Logf("history-orders: %d (found ours=%v)", len(hist.List), inHist)
}
