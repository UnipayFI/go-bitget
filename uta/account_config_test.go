package uta

import (
	"testing"

	"github.com/UnipayFI/go-bitget/client"
)

func TestAccountConfig(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	cat := string(CategoryUSDTFutures)
	symbol := "BTCUSDT"

	// Account settings.
	settings, err := c.NewGetAccountSettingsService().Do(cx)
	if err != nil {
		t.Fatalf("account settings: %v", err)
	}
	t.Logf("settings: %+v", settings)
	raw := fetchRawGet(t, c, cx, "/api/v3/account/settings", nil, true)
	assertCovers(t, "account/settings", raw, settings)

	// Account info.
	info, err := c.NewGetAccountInfoService().Do(cx)
	if err != nil {
		t.Fatalf("account info: %v", err)
	}
	t.Logf("info: %+v", info)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/info", nil, true)
	assertCovers(t, "account/info", raw, info)

	// Delta info.
	delta, err := c.NewGetDeltaInfoService().Do(cx)
	if err != nil {
		t.Fatalf("delta info: %v", err)
	}
	t.Logf("delta: %+v", delta)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/delta-info", nil, true)
	assertCovers(t, "account/delta-info", raw, delta)

	// Account fee rate.
	feeRate, err := c.NewGetAccountFeeRateService(CategoryUSDTFutures, symbol).Do(cx)
	if err != nil {
		t.Fatalf("account fee rate: %v", err)
	}
	t.Logf("fee rate: %+v", feeRate)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/fee-rate", map[string]string{"category": cat, "symbol": symbol}, true)
	assertCovers(t, "account/fee-rate", raw, feeRate)

	// All fee rates.
	allFeeRate, err := c.NewGetAllFeeRateService(CategoryUSDTFutures).SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("all fee rate: %v", err)
	}
	t.Logf("all fee rate: %d", len(allFeeRate))
	raw = fetchRawGet(t, c, cx, "/api/v3/account/all-fee-rate", map[string]string{"category": cat, "symbol": symbol}, true)
	assertCovers(t, "account/all-fee-rate", raw, allFeeRate)

	// Open-interest limit.
	oiLimit, err := c.NewGetOILimitService(CategoryUSDTFutures, symbol).Do(cx)
	if err != nil {
		t.Fatalf("oi limit: %v", err)
	}
	t.Logf("oi limit: %+v", oiLimit)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/open-interest-limit", map[string]string{"category": cat, "symbol": symbol}, true)
	assertCovers(t, "account/open-interest-limit", raw, oiLimit)

	// Deduct info.
	deduct, err := c.NewGetDeductInfoService().Do(cx)
	if err != nil {
		t.Fatalf("deduct info: %v", err)
	}
	t.Logf("deduct: %+v", deduct)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/deduct-info", nil, true)
	assertCovers(t, "account/deduct-info", raw, deduct)

	// Switch status. Only returns data when a parent account has a UTA<->classic
	// switch in progress; otherwise Bitget replies 40054 ("data is empty"),
	// which still confirms the endpoint + signing work.
	switchStatus, err := c.NewGetSwitchStatusService().Do(cx)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok && apiErr.Code == "40054" {
			t.Logf("switch status: no switch in progress (40054), endpoint OK")
		} else {
			t.Fatalf("switch status: %v", err)
		}
	} else {
		t.Logf("switch status: %+v", switchStatus)
		raw = fetchRawGet(t, c, cx, "/api/v3/account/switch-status", nil, true)
		assertCovers(t, "account/switch-status", raw, switchStatus)
	}
}
