package mix

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

func TestMixAccount(t *testing.T) {
	c := NewMixClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const (
		productType = ProductTypeUSDTFutures
		symbol      = "BTCUSDT"
		marginCoin  = "USDT"
	)
	tol := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001", "40847"}

	// GET /api/v2/mix/account/account
	{
		params := map[string]string{"symbol": symbol, "productType": string(productType), "marginCoin": marginCoin}
		resp, err := c.NewGetSingleAccountService(symbol, productType, marginCoin).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "single-account", err, tol...) {
				t.Fatalf("single-account: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/account", params, true)
			apitest.AssertCovers(t, "single-account", raw, resp)
		}
	}

	// GET /api/v2/mix/account/accounts
	{
		params := map[string]string{"productType": string(productType)}
		resp, err := c.NewGetAccountListService(productType).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "accounts", err, tol...) {
				t.Fatalf("accounts: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/accounts", params, true)
			apitest.AssertCovers(t, "accounts", raw, resp)
		}
	}

	// GET /api/v2/mix/account/sub-account-assets
	{
		params := map[string]string{"productType": string(productType)}
		resp, err := c.NewGetSubAccountAssetsService(productType).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "sub-account-assets", err, tol...) {
				t.Fatalf("sub-account-assets: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/sub-account-assets", params, true)
			apitest.AssertCovers(t, "sub-account-assets", raw, resp)
		}
	}

	// GET /api/v2/mix/account/interest-history (requires a startTime+endTime window)
	{
		end := time.Now()
		start := end.Add(-7 * 24 * time.Hour)
		params := map[string]string{
			"productType": string(productType),
			"startTime":   strconv.FormatInt(start.UnixMilli(), 10),
			"endTime":     strconv.FormatInt(end.UnixMilli(), 10),
		}
		resp, err := c.NewGetInterestHistoryService(productType).SetStartTime(start).SetEndTime(end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "interest-history", err, tol...) {
				t.Fatalf("interest-history: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/interest-history", params, true)
			apitest.AssertCovers(t, "interest-history", raw, resp)
		}
	}

	// GET /api/v2/mix/account/max-open
	{
		params := map[string]string{"symbol": symbol, "productType": string(productType), "marginCoin": marginCoin, "posSide": string(HoldSideLong), "orderType": string(OrderTypeMarket)}
		resp, err := c.NewGetMaxOpenService(symbol, productType, marginCoin, HoldSideLong, OrderTypeMarket).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "max-open", err, tol...) {
				t.Fatalf("max-open: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/max-open", params, true)
			apitest.AssertCovers(t, "max-open", raw, resp)
		}
	}

	// GET /api/v2/mix/account/liq-price
	{
		openAmount := decimal.RequireFromString("10")
		params := map[string]string{"symbol": symbol, "productType": string(productType), "marginCoin": marginCoin, "posSide": string(HoldSideLong), "orderType": string(OrderTypeMarket), "openAmount": openAmount.String()}
		resp, err := c.NewGetLiquidationPriceService(symbol, productType, marginCoin, HoldSideLong, OrderTypeMarket, openAmount).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "liq-price", err, tol...) {
				t.Fatalf("liq-price: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/liq-price", params, true)
			apitest.AssertCovers(t, "liq-price", raw, resp)
		}
	}

	// GET /api/v2/mix/account/open-count
	{
		openAmount := decimal.RequireFromString("10")
		openPrice := decimal.RequireFromString("50000")
		params := map[string]string{"symbol": symbol, "productType": string(productType), "marginCoin": marginCoin, "openAmount": openAmount.String(), "openPrice": openPrice.String()}
		resp, err := c.NewGetOpenCountService(symbol, productType, marginCoin, openAmount, openPrice).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "open-count", err, tol...) {
				t.Fatalf("open-count: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/open-count", params, true)
			apitest.AssertCovers(t, "open-count", raw, resp)
		}
	}

	// GET /api/v2/mix/account/bill
	{
		params := map[string]string{"productType": string(productType)}
		resp, err := c.NewGetAccountBillService(productType).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "bill", err, tol...) {
				t.Fatalf("bill: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/bill", params, true)
			apitest.AssertCovers(t, "bill", raw, resp)
		}
	}

	// GET /api/v2/mix/account/transfer-limits
	{
		params := map[string]string{"coin": "USDT"}
		resp, err := c.NewGetTransferLimitsService("USDT").Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "transfer-limits", err, tol...) {
				t.Fatalf("transfer-limits: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/transfer-limits", params, true)
			apitest.AssertCovers(t, "transfer-limits", raw, resp)
		}
	}

	// GET /api/v2/mix/account/union-config
	{
		resp, err := c.NewGetUnionConfigService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "union-config", err, tol...) {
				t.Fatalf("union-config: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/union-config", nil, true)
			apitest.AssertCovers(t, "union-config", raw, resp)
		}
	}

	// GET /api/v2/mix/account/switch-union-usdt
	{
		resp, err := c.NewGetSwitchUnionUSDTService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "switch-union-usdt", err, tol...) {
				t.Fatalf("switch-union-usdt: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/switch-union-usdt", nil, true)
			apitest.AssertCovers(t, "switch-union-usdt", raw, resp)
		}
	}

	// GET /api/v2/mix/account/isolated-symbols
	{
		params := map[string]string{"productType": string(productType)}
		resp, err := c.NewGetIsolatedSymbolsService(productType).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "isolated-symbols", err, tol...) {
				t.Fatalf("isolated-symbols: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/account/isolated-symbols", params, true)
			apitest.AssertCovers(t, "isolated-symbols", raw, resp)
		}
	}
}
