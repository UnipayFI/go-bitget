package spot

import (
	"os"
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

// TestSpotPlanLifecycle exercises the state-changing spot plan (trigger) order
// endpoints and reconciles the populated PlanOrder struct that the read-only
// test cannot verify on an account with no resting plan orders. Gated behind
// BITGET_TEST_WRITE=1. It places a BUY trigger order with a triggerPrice far
// ABOVE market (200000) so it stays untriggered (no fill, no spend), then
// modifies and cancels it.
func TestSpotPlanLifecycle(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") != "1" {
		t.Skip("set BITGET_TEST_WRITE=1 to run state-changing spot plan order tests")
	}
	c := NewSpotClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync server time: %v", err)
	}
	const symbol = "BTCUSDT"
	cx := apitest.Ctx(t)

	// 1. Place a BUY trigger order well above market -> stays not_trigger.
	placed, err := c.NewPlacePlanOrderService(symbol, SideBuy,
		decimal.RequireFromString("200000"), OrderTypeLimit,
		decimal.RequireFromString("0.0002"), TriggerTypeFillPrice).
		SetExecutePrice(decimal.RequireFromString("190000")).
		Do(cx)
	if err != nil {
		t.Fatalf("place plan order: %v", err)
	}
	if placed.OrderID == "" {
		t.Fatalf("place plan order: empty orderId")
	}
	orderID := placed.OrderID
	t.Logf("placed resting plan order %s", orderID)

	// 2. current-plan-order now contains the resting plan order.
	{
		resp, err := c.NewGetCurrentPlanOrderService().SetSymbol(symbol).Do(cx)
		if err != nil {
			t.Fatalf("current-plan-order: %v", err)
		}
		if len(resp.OrderList) == 0 {
			t.Fatalf("current-plan-order: expected the resting plan order, got none")
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/trade/current-plan-order", map[string]string{"symbol": symbol}, true)
		apitest.AssertCovers(t, "PlanOrderList", raw, resp)
	}

	// 3. modify the plan order.
	{
		mod, err := c.NewModifyPlanOrderService(decimal.RequireFromString("210000"), OrderTypeLimit, decimal.RequireFromString("0.0003")).
			SetOrderID(orderID).
			SetExecutePrice(decimal.RequireFromString("195000")).
			Do(cx)
		if err != nil {
			t.Fatalf("modify-plan-order: %v", err)
		}
		t.Logf("modified plan order -> %s", mod.OrderID)
	}

	// 4. cancel the plan order.
	{
		cancelled, err := c.NewCancelPlanOrderService().SetOrderID(orderID).Do(cx)
		if err != nil {
			t.Fatalf("cancel-plan-order: %v", err)
		}
		t.Logf("cancel-plan-order result=%s", cancelled.Result)
	}
}
