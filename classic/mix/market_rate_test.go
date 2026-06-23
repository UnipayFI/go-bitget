package mix

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixMarketRate(t *testing.T) {
	c := NewMixClient(apitest.PublicOptions()...)
	cx := apitest.Ctx(t)

	// vip-fee-rate -- public, no required params.
	{
		resp, err := c.NewGetVIPFeeRateService().Do(cx)
		if err != nil {
			t.Fatalf("vip-fee-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/vip-fee-rate", map[string]string{}, false)
		apitest.AssertCovers(t, "vip-fee-rate", raw, resp)
	}

	// union-interest-rate-history -- public, requires coin.
	{
		params := map[string]string{"coin": "USDT"}
		resp, err := c.NewGetInterestRateHistoryService("USDT").Do(cx)
		if err != nil {
			t.Fatalf("interest-rate-history: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/union-interest-rate-history", params, false)
		apitest.AssertCovers(t, "interest-rate-history", raw, resp)
	}

	// exchange-rate -- public, no params.
	{
		resp, err := c.NewGetExchangeRateService().Do(cx)
		if err != nil {
			t.Fatalf("exchange-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/exchange-rate", nil, false)
		apitest.AssertCovers(t, "exchange-rate", raw, resp)
	}

	// discount-rate -- public, no required params.
	{
		resp, err := c.NewGetDiscountRateService().Do(cx)
		if err != nil {
			t.Fatalf("discount-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/discount-rate", map[string]string{}, false)
		apitest.AssertCovers(t, "discount-rate", raw, resp)
	}

	// history-fund-rate -- public, requires symbol + productType.
	{
		params := map[string]string{"symbol": "BTCUSDT", "productType": string(ProductTypeUSDTFutures)}
		resp, err := c.NewGetHistoryFundRateService("BTCUSDT", ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			t.Fatalf("history-fund-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/history-fund-rate", params, false)
		apitest.AssertCovers(t, "history-fund-rate", raw, resp)
	}

	// current-fund-rate -- public, requires symbol + productType.
	{
		params := map[string]string{"symbol": "BTCUSDT", "productType": string(ProductTypeUSDTFutures)}
		resp, err := c.NewGetCurrentFundRateService("BTCUSDT", ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			t.Fatalf("current-fund-rate: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/market/current-fund-rate", params, false)
		apitest.AssertCovers(t, "current-fund-rate", raw, resp)
	}
}
