package spot

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestSpotPlan(t *testing.T) {
	c := NewSpotClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	// current-plan-order: open (untriggered) plan orders.
	var planOrderID string
	t.Run("current-plan-order", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/current-plan-order", nil, true)
		apitest.AssertCovers(t, "current-plan-order", raw, &PlanOrderList{})

		resp, err := c.NewGetCurrentPlanOrderService().Do(ctx)
		if err != nil {
			t.Fatalf("current-plan-order Do: %v", err)
		}
		if len(resp.OrderList) > 0 {
			planOrderID = resp.OrderList[0].OrderID
		}
	})

	// history-plan-order: triggered/cancelled plan orders.
	t.Run("history-plan-order", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/history-plan-order", nil, true)
		apitest.AssertCovers(t, "history-plan-order", raw, &PlanOrderList{})

		if planOrderID == "" {
			resp, err := c.NewGetHistoryPlanOrderService().Do(ctx)
			if err != nil {
				t.Fatalf("history-plan-order Do: %v", err)
			}
			if len(resp.OrderList) > 0 {
				planOrderID = resp.OrderList[0].OrderID
			}
		}
	})

	// plan-sub-order: needs a real plan order id; skip when the account has none.
	t.Run("plan-sub-order", func(t *testing.T) {
		if planOrderID == "" {
			t.Skip("no plan order id available to query sub orders")
		}
		ctx := apitest.Ctx(t)
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/trade/plan-sub-order", map[string]string{"planOrderId": planOrderID}, true)
		apitest.AssertCovers(t, "plan-sub-order", raw, []PlanSubOrder(nil))
	})
}
