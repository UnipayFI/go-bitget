package earn

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestEarnSharkFin(t *testing.T) {
	c := NewEarnClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	tolCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "70102", "40029", "40913", "43012", "40054", "70231", "70102", "47002", "40808"}

	// SharkFin Products -- GET /api/v2/earn/sharkfin/product
	var productID string
	{
		params := map[string]string{"coin": "USDT"}
		resp, err := c.NewGetSharkFinProductService("USDT").Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sharkfin/product", err, tolCodes...) {
				t.Fatalf("sharkfin/product: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/sharkfin/product", params, true)
			apitest.AssertCovers(t, "sharkfin/product", raw, resp)
			if len(resp.ResultList) > 0 {
				productID = resp.ResultList[0].ProductID
			}
		}
	}

	// SharkFin Account -- GET /api/v2/earn/sharkfin/account
	{
		resp, err := c.NewGetSharkFinAccountService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sharkfin/account", err, tolCodes...) {
				t.Fatalf("sharkfin/account: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/sharkfin/account", nil, true)
			apitest.AssertCovers(t, "sharkfin/account", raw, resp)
		}
	}

	// SharkFin Assets -- GET /api/v2/earn/sharkfin/assets
	{
		params := map[string]string{"status": string(SharkFinAssetStatusSubscribed)}
		resp, err := c.NewGetSharkFinAssetsService(SharkFinAssetStatusSubscribed).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sharkfin/assets", err, tolCodes...) {
				t.Fatalf("sharkfin/assets: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/sharkfin/assets", params, true)
			apitest.AssertCovers(t, "sharkfin/assets", raw, resp)
		}
	}

	// SharkFin Records -- GET /api/v2/earn/sharkfin/records
	{
		params := map[string]string{"type": string(SharkFinRecordTypeSubscription)}
		resp, err := c.NewGetSharkFinRecordsService(SharkFinRecordTypeSubscription).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sharkfin/records", err, tolCodes...) {
				t.Fatalf("sharkfin/records: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/sharkfin/records", params, true)
			apitest.AssertCovers(t, "sharkfin/records", raw, resp)
		}
	}

	// SharkFin Subscription Detail -- GET /api/v2/earn/sharkfin/subscribe-info
	if productID != "" {
		params := map[string]string{"productId": productID}
		resp, err := c.NewGetSharkFinSubscribeInfoService(productID).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sharkfin/subscribe-info", err, tolCodes...) {
				t.Fatalf("sharkfin/subscribe-info: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/earn/sharkfin/subscribe-info", params, true)
			apitest.AssertCovers(t, "sharkfin/subscribe-info", raw, resp)
		}
	} else {
		t.Log("sharkfin/subscribe-info: no product available to query; skipping")
	}

	// subscribe (POST, state-changing) and subscribe-result (needs a real orderId)
	// are intentionally not tested here.
}
