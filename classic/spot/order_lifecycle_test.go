package spot

import (
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

// TestSpotOrderLifecycle exercises the state-changing spot trade endpoints
// against the live account with tiny BTC amounts, and reconciles the populated
// order/fill structs (OrderInfo, UnfilledOrder, Fill, HistoryOrder) that the
// read-only tests cannot verify on a fresh account. It is gated behind
// BITGET_TEST_WRITE=1 because it places (and fills) real orders.
//
// Flow:
//  1. Limit BUY far below market -> rests unfilled (no spend): verifies
//     place-order, orderInfo, unfilled-orders, then cancel-order.
//  2. Market BUY of ~2 USDT -> fills immediately: verifies fills + history-orders
//     against real executed data. Leaves a few cents of BTC dust.
func TestSpotOrderLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run state-changing spot order tests")
	}
	c := NewSpotClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync server time: %v", err)
	}
	const symbol = "BTCUSDT"

	// 1. Resting limit BUY at 30000 (well below the ~62k market) — never fills.
	cx := apitest.Ctx(t)
	placed, err := c.NewPlaceOrderService(symbol, SideBuy, OrderTypeLimit, ForceGTC, decimal.RequireFromString("0.0002")).
		SetPrice(decimal.RequireFromString("30000")).
		Do(cx)
	if err != nil {
		t.Fatalf("place limit order: %v", err)
	}
	if placed.OrderID == "" {
		t.Fatalf("place limit order: empty orderId")
	}
	orderID := placed.OrderID
	t.Logf("placed resting limit order %s", orderID)

	// 2. orderInfo for the resting order.
	{
		info, err := c.NewGetOrderInfoService().SetOrderID(orderID).Do(cx)
		if err != nil {
			t.Fatalf("orderInfo: %v", err)
		}
		if len(info) == 0 {
			t.Fatalf("orderInfo: empty")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/trade/orderInfo", map[string]string{"orderId": orderID}, true)
		apitest.AssertCovers(t, "OrderInfo", raw, info)
	}

	// 3. unfilled-orders now contains the resting order.
	{
		unfilled, err := c.NewGetUnfilledOrdersService().SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("unfilled-orders: %v", err)
		}
		if len(unfilled) == 0 {
			t.Fatalf("unfilled-orders: expected the resting order, got none")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/trade/unfilled-orders", map[string]string{"symbol": symbol}, true)
		apitest.AssertCovers(t, "UnfilledOrder", raw, unfilled)
	}

	// 4. cancel the resting order.
	{
		cancelled, err := c.NewCancelOrderService(symbol).SetOrderID(orderID).Do(cx)
		if err != nil {
			t.Fatalf("cancel-order: %v", err)
		}
		if cancelled.OrderID != orderID {
			t.Fatalf("cancel-order: echoed orderId %q != %q", cancelled.OrderID, orderID)
		}
		t.Logf("cancelled resting order %s", orderID)
	}

	// 5. Market BUY of ~2 USDT (size is the quote amount for spot market buys).
	{
		mkt, err := c.NewPlaceOrderService(symbol, SideBuy, OrderTypeMarket, ForceGTC, decimal.RequireFromString("2")).Do(cx)
		if err != nil {
			t.Fatalf("market buy: %v", err)
		}
		t.Logf("placed market buy %s (~2 USDT of BTC)", mkt.OrderID)
		time.Sleep(2 * time.Second) // let the fill settle
	}

	// 6. fills now contains the executed trade.
	{
		cx := apitest.Ctx(t)
		fills, err := c.NewGetFillsService().SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("fills: %v", err)
		}
		if len(fills) == 0 {
			t.Fatalf("fills: expected the market-buy fill, got none")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/trade/fills", map[string]string{"symbol": symbol}, true)
		apitest.AssertCovers(t, "Fill", raw, fills)
	}

	// 7. history-orders now contains the filled market order.
	{
		cx := apitest.Ctx(t)
		hist, err := c.NewGetHistoryOrdersService().SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("history-orders: %v", err)
		}
		if len(hist) == 0 {
			t.Fatalf("history-orders: expected the filled order, got none")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/trade/history-orders", map[string]string{"symbol": symbol}, true)
		apitest.AssertCovers(t, "HistoryOrder", raw, hist)
	}
}
