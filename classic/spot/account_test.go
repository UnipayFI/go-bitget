package spot

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestSpotAccount(t *testing.T) {
	c := NewSpotClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync server time: %v", err)
	}

	// GET /api/v2/spot/account/info
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetAccountInfoService().Do(ctx)
		if err != nil {
			t.Fatalf("account info: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/info", nil, true)
		apitest.AssertCovers(t, "AccountInfo", raw, resp)
	}

	// GET /api/v2/spot/account/assets
	{
		ctx := apitest.Ctx(t)
		params := map[string]string{"assetType": "all"}
		resp, err := c.NewGetAccountAssetsService().SetAssetType("all").Do(ctx)
		if err != nil {
			t.Fatalf("account assets: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/assets", params, true)
		apitest.AssertCovers(t, "AccountAsset", raw, resp)
	}

	// GET /api/v2/spot/account/subaccount-assets
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetSubaccountAssetsService().Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "SubaccountAssets", err, "40029", "40037", "22001", "40068", "40014") {
				goto bills
			}
			t.Fatalf("subaccount assets: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/subaccount-assets", nil, true)
		apitest.AssertCovers(t, "SubaccountAssets", raw, resp)
	}

bills:
	// GET /api/v2/spot/account/bills
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetAccountBillsService().Do(ctx)
		if err != nil {
			t.Fatalf("account bills: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/bills", nil, true)
		apitest.AssertCovers(t, "AccountBill", raw, resp)
	}

	// GET /api/v2/spot/wallet/transfer-coin-info (data is an array of strings)
	{
		ctx := apitest.Ctx(t)
		if _, err := c.NewGetTransferCoinInfoService("spot", "usdt_futures").Do(ctx); err != nil {
			t.Fatalf("transfer coin info: %v", err)
		}
	}

	// GET /api/v2/spot/account/deduct-info
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetDeductInfoService().Do(ctx)
		if err != nil {
			t.Fatalf("deduct info: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/deduct-info", nil, true)
		apitest.AssertCovers(t, "DeductInfo", raw, resp)
	}

	// GET /api/v2/spot/account/transferRecords (coin required)
	{
		ctx := apitest.Ctx(t)
		params := map[string]string{"coin": "USDT"}
		resp, err := c.NewGetTransferRecordsService("USDT").Do(ctx)
		if err != nil {
			t.Fatalf("transfer records: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/transferRecords", params, true)
		apitest.AssertCovers(t, "TransferRecord", raw, resp)
	}

	// GET /api/v2/spot/account/sub-main-trans-record
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetSubMainTransRecordService().Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "SubMainTransferRecord", err, "40029", "40037", "22001", "40068", "40014", "40054") {
				goto upgradeStatus
			}
			t.Fatalf("sub-main transfer record: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/sub-main-trans-record", nil, true)
		apitest.AssertCovers(t, "SubMainTransferRecord", raw, resp)
	}

upgradeStatus:
	// GET /api/v2/spot/account/upgrade-status
	{
		ctx := apitest.Ctx(t)
		resp, err := c.NewGetUpgradeStatusService().Do(ctx)
		if err != nil {
			if apitest.Tolerable(t, "UpgradeStatus", err, "40029", "40037", "22001", "40068", "40014", "40054") {
				return
			}
			t.Fatalf("upgrade status: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/account/upgrade-status", nil, true)
		apitest.AssertCovers(t, "UpgradeStatus", raw, resp)
	}
}
