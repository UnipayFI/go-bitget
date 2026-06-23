package margin

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMarginCommon(t *testing.T) {
	c := NewMarginClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// GET /api/v2/margin/currencies — Margin Support Currencies (no params).
	{
		const label = "margin/currencies"
		resp, err := c.NewGetMarginCurrenciesService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, "40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/currencies", nil, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	// GET /api/v2/margin/interest-rate-record — Leverage Interest Rate (coin required).
	{
		const label = "margin/interest-rate-record"
		params := map[string]string{"coin": "BTC"}
		resp, err := c.NewGetInterestRateRecordService("BTC").Do(cx)
		if err != nil {
			if apitest.Tolerable(t, label, err, "40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/margin/interest-rate-record", params, true)
		apitest.AssertCovers(t, label, raw, resp)
	}
}
