package common

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestCommonRate(t *testing.T) {
	// Public: server time.
	{
		c := NewCommonClient(apitest.PublicOptions()...)
		ctx := apitest.Ctx(t)
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/public/time", nil, false)
		resp, err := c.NewGetServerTimeService().Do(ctx)
		if err != nil {
			t.Fatalf("GetServerTime: %v", err)
		}
		apitest.AssertCovers(t, "ServerTime", raw, resp)
	}

	// Private reads: fee rates (account config, no pre-existing order state).
	c := NewCommonClient(apitest.AuthOptions(t)...)
	ctx := apitest.Ctx(t)
	if err := c.SyncServerTime(ctx); err != nil {
		t.Fatalf("SyncServerTime: %v", err)
	}

	tradeRateParams := map[string]string{"symbol": "BTCUSDT", "businessType": "spot"}

	{
		label := "TradeRate"
		resp, err := c.NewGetTradeRateService("BTCUSDT", TradeRateBusinessTypeSpot).Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, label, err, "40034", "40309", "22002") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/common/trade-rate", tradeRateParams, true)
		apitest.AssertCovers(t, label, raw, resp)
	}

	{
		label := "AllTradeRate"
		resp, err := c.NewGetAllTradeRateService("BTCUSDT", TradeRateBusinessTypeSpot).Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, label, err, "40034", "40309", "22002") {
				return
			}
			t.Fatalf("%s: %v", label, err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/common/all-trade-rate", tradeRateParams, true)
		apitest.AssertCovers(t, label, raw, resp)
	}
}
