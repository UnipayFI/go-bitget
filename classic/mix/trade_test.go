package mix

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixTrade(t *testing.T) {
	c := NewMixClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const product = ProductTypeUSDTFutures

	// tolerable: sub-account capability gates and empty data.
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001"}

	// orders-pending
	{
		resp, err := c.NewGetOrdersPendingService(product).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "orders-pending", err, okCodes...) {
				t.Fatalf("orders-pending: %v", err)
			}
		} else {
			params := map[string]string{"productType": string(product)}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/orders-pending", params, true)
			apitest.AssertCovers(t, "orders-pending", raw, resp)
		}
	}

	// orders-history
	{
		resp, err := c.NewGetOrdersHistoryService(product).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "orders-history", err, okCodes...) {
				t.Fatalf("orders-history: %v", err)
			}
		} else {
			params := map[string]string{"productType": string(product)}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/orders-history", params, true)
			apitest.AssertCovers(t, "orders-history", raw, resp)
		}
	}

	// fills
	{
		resp, err := c.NewGetOrderFillsService(product).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "fills", err, okCodes...) {
				t.Fatalf("fills: %v", err)
			}
		} else {
			params := map[string]string{"productType": string(product)}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/fills", params, true)
			apitest.AssertCovers(t, "fills", raw, resp)
		}
	}

	// fill-history
	{
		resp, err := c.NewGetFillHistoryService(product).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "fill-history", err, okCodes...) {
				t.Fatalf("fill-history: %v", err)
			}
		} else {
			params := map[string]string{"productType": string(product)}
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/order/fill-history", params, true)
			apitest.AssertCovers(t, "fill-history", raw, resp)
		}
	}
}
