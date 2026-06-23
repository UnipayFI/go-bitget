package mix

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixPlan(t *testing.T) {
	c := NewMixClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	const product = ProductTypeUSDTFutures
	tolerable := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001"}

	// orders-plan-pending: open (untriggered) trigger orders.
	var planOrderID string
	t.Run("orders-plan-pending", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		params := map[string]string{
			"productType": string(product),
			"planType":    string(PlanTypeNormalPlan),
		}

		resp, err := c.NewGetOrdersPlanPendingService(product, PlanTypeNormalPlan).Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "orders-plan-pending", err, tolerable...) {
				return
			}
			t.Fatalf("orders-plan-pending Do: %v", err)
		}
		if len(resp.EntrustedList) > 0 {
			planOrderID = resp.EntrustedList[0].OrderID
		}

		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/order/orders-plan-pending", params, true)
		apitest.AssertCovers(t, "orders-plan-pending", raw, resp)
	})

	// orders-plan-history: triggered/cancelled trigger orders.
	t.Run("orders-plan-history", func(t *testing.T) {
		ctx := apitest.Ctx(t)
		params := map[string]string{
			"productType": string(product),
			"planType":    string(PlanTypeNormalPlan),
		}

		resp, err := c.NewGetOrdersPlanHistoryService(product, PlanTypeNormalPlan).Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "orders-plan-history", err, tolerable...) {
				return
			}
			t.Fatalf("orders-plan-history Do: %v", err)
		}
		if planOrderID == "" && len(resp.EntrustedList) > 0 {
			planOrderID = resp.EntrustedList[0].OrderID
		}

		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/order/orders-plan-history", params, true)
		apitest.AssertCovers(t, "orders-plan-history", raw, resp)
	})

	// plan-sub-order: needs a real trigger order id; skip when the account has none.
	t.Run("plan-sub-order", func(t *testing.T) {
		if planOrderID == "" {
			t.Skip("no trigger order id available to query sub orders")
		}
		ctx := apitest.Ctx(t)
		params := map[string]string{
			"productType": string(product),
			"planType":    string(PlanTypeNormalPlan),
			"planOrderId": planOrderID,
		}

		resp, err := c.NewGetPlanSubOrderService(product, PlanTypeNormalPlan, planOrderID).Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "plan-sub-order", err, tolerable...) {
				return
			}
			t.Fatalf("plan-sub-order Do: %v", err)
		}

		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/mix/order/plan-sub-order", params, true)
		apitest.AssertCovers(t, "plan-sub-order", raw, resp)
	})
}
