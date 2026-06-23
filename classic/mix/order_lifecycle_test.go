package mix

import (
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

// TestMixOrderLifecycle exercises the state-changing futures endpoints against
// the live account with the minimum trade size (0.0001 BTC ~ 6 USDT notional,
// crossed margin) and reconciles the populated order/position/fill structs that
// the read-only tests cannot verify with no open orders/positions. Gated behind
// BITGET_TEST_WRITE=1.
//
// Flow:
//   - set-leverage (account-setting write)
//   - resting limit (open) far below market -> place / orders-pending / detail /
//     cancel — no fill, no position.
//   - market open long -> single-position + order fills populated -> flash close.
func TestMixOrderLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run state-changing futures order tests")
	}
	c := NewMixClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync server time: %v", err)
	}
	const (
		productType = ProductTypeUSDTFutures
		symbol      = "BTCUSDT"
		marginCoin  = "USDT"
		size        = "0.0001" // == minTradeNum for BTCUSDT
	)
	cx := apitest.Ctx(t)

	// 1. set-leverage (idempotent: 10x).
	if _, err := c.NewSetLeverageService(symbol, productType, marginCoin).SetLeverage("10").Do(cx); err != nil {
		t.Fatalf("set-leverage: %v", err)
	}
	t.Log("set leverage 10x")

	// 2. Resting limit open BUY far below market (30000) -> stays unfilled.
	placed, err := c.NewPlaceOrderService(symbol, productType, MarginModeCrossed, marginCoin,
		decimal.RequireFromString(size), SideBuy, OrderTypeLimit).
		// 55000 is below the ~63k market (rests, never fills) yet keeps the
		// notional above the 5 USDT minimum (0.0001 * 55000 = 5.5).
		SetPrice(decimal.RequireFromString("55000")).
		SetForce(ForceGTC).
		SetTradeSide(TradeSideOpen).
		Do(cx)
	if err != nil {
		t.Fatalf("place limit: %v", err)
	}
	orderID := placed.OrderId
	t.Logf("placed resting limit %s", orderID)

	// 3. orders-pending now contains it.
	{
		resp, err := c.NewGetOrdersPendingService(productType).SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("orders-pending: %v", err)
		}
		if len(resp.EntrustedList) == 0 {
			t.Fatalf("orders-pending: expected the resting order")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/orders-pending", map[string]string{"productType": string(productType), "symbol": symbol}, true)
		apitest.AssertCovers(t, "MixOrderList(pending)", raw, resp)
	}

	// 4. order detail by orderId.
	{
		resp, err := c.NewGetOrderDetailService(symbol, productType).SetOrderId(orderID).Do(cx)
		if err != nil {
			t.Fatalf("order detail: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/detail", map[string]string{"symbol": symbol, "productType": string(productType), "orderId": orderID}, true)
		apitest.AssertCovers(t, "MixOrder(detail)", raw, resp)
	}

	// 5. cancel it.
	if _, err := c.NewCancelOrderService(symbol, productType).SetOrderId(orderID).Do(cx); err != nil {
		t.Fatalf("cancel-order: %v", err)
	}
	t.Logf("cancelled resting limit %s", orderID)

	// 6. Market open long -> opens a position.
	if _, err := c.NewPlaceOrderService(symbol, productType, MarginModeCrossed, marginCoin,
		decimal.RequireFromString(size), SideBuy, OrderTypeMarket).
		SetTradeSide(TradeSideOpen).Do(cx); err != nil {
		t.Fatalf("market open: %v", err)
	}
	t.Log("opened market long")
	time.Sleep(2 * time.Second)

	// 7. single-position now populated.
	{
		cx := apitest.Ctx(t)
		resp, err := c.NewGetSinglePositionService(productType, symbol, marginCoin).Do(cx)
		if err != nil {
			t.Fatalf("single-position: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/position/single-position", map[string]string{"productType": string(productType), "symbol": symbol, "marginCoin": marginCoin}, true)
		apitest.AssertCovers(t, "Position", raw, resp)
	}

	// 8. order fills populated.
	{
		cx := apitest.Ctx(t)
		resp, err := c.NewGetOrderFillsService(productType).SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("order fills: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/fills", map[string]string{"productType": string(productType), "symbol": symbol}, true)
		apitest.AssertCovers(t, "MixFillList", raw, resp)
	}

	// 9. flash-close the long position.
	{
		cx := apitest.Ctx(t)
		res, err := c.NewClosePositionsService(productType).SetSymbol(symbol).SetHoldSide(HoldSideLong).Do(cx)
		if err != nil {
			t.Fatalf("close-positions: %v", err)
		}
		t.Logf("close-positions: %d closed, %d failed", len(res.SuccessList), len(res.FailureList))
	}
}
