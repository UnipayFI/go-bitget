package broker

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// brokerOKCodes are the capability/empty-data responses we tolerate: the shared
// test key is a permission-limited sub-account, so broker endpoints commonly
// answer "not enrolled / no data / no permission" even though signing worked.
var brokerOKCodes = []string{
	"40068", // sub-account access disabled
	"40014", // needs child-account perms
	"40054", // data empty
	"40099", // no permission for this product line
	"40034", // no permission for this product line
	"40029", // not a broker
	"40037", // sub-account not found
	"47001",
	"22001",
	"70231", // broker partner channel not enabled for this account
	"70102", // broker identity authentication failed
}

func TestBrokerSubaccount(t *testing.T) {
	c := NewBrokerClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// GET /api/v2/broker/account/info
	{
		resp, err := c.NewGetBrokerAccountInfoService().Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "account/info", err, brokerOKCodes...) {
				t.Fatalf("account/info: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/account/info", nil, true)
			apitest.AssertCovers(t, "account/info", raw, resp)
		}
	}

	// GET /api/v2/broker/account/subaccount-list
	{
		params := map[string]string{"limit": "10"}
		resp, err := c.NewGetSubaccountListService().SetLimit(10).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "subaccount-list", err, brokerOKCodes...) {
				t.Fatalf("subaccount-list: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/account/subaccount-list", params, true)
			apitest.AssertCovers(t, "subaccount-list", raw, resp)
		}
	}

	// GET /api/v2/broker/subaccount-deposit (records, bare array)
	{
		resp, err := c.NewGetSubaccountDepositRecordsService().SetLimit(20).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "subaccount-deposit", err, brokerOKCodes...) {
				t.Fatalf("subaccount-deposit: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/subaccount-deposit", map[string]string{"limit": "20"}, true)
			apitest.AssertCovers(t, "subaccount-deposit", raw, resp)
		}
	}

	// GET /api/v2/broker/subaccount-withdrawal (records, object{resultList})
	{
		resp, err := c.NewGetSubaccountWithdrawalRecordsService().SetLimit(20).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "subaccount-withdrawal", err, brokerOKCodes...) {
				t.Fatalf("subaccount-withdrawal: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/subaccount-withdrawal", map[string]string{"limit": "20"}, true)
			apitest.AssertCovers(t, "subaccount-withdrawal", raw, resp)
		}
	}

	// GET /api/v2/broker/all-sub-deposit-withdrawal
	{
		resp, err := c.NewGetAllSubDepositWithdrawalService().SetLimit(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "all-sub-deposit-withdrawal", err, brokerOKCodes...) {
				t.Fatalf("all-sub-deposit-withdrawal: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/all-sub-deposit-withdrawal", map[string]string{"limit": "100"}, true)
			apitest.AssertCovers(t, "all-sub-deposit-withdrawal", raw, resp)
		}
	}

	// GET /api/v2/broker/subaccounts
	{
		resp, err := c.NewGetBrokerSubaccountsService().SetPageSize(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "subaccounts", err, brokerOKCodes...) {
				t.Fatalf("subaccounts: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/subaccounts", map[string]string{"pageSize": "100"}, true)
			apitest.AssertCovers(t, "subaccounts", raw, resp)
		}
	}

	// GET /api/v2/broker/commissions
	{
		resp, err := c.NewGetBrokerCommissionsService().SetPageSize(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "commissions", err, brokerOKCodes...) {
				t.Fatalf("commissions: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/commissions", map[string]string{"pageSize": "100"}, true)
			apitest.AssertCovers(t, "commissions", raw, resp)
		}
	}

	// GET /api/v2/broker/trade-volume
	{
		resp, err := c.NewGetBrokerTradeVolumeService().SetPageSize(100).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "trade-volume", err, brokerOKCodes...) {
				t.Fatalf("trade-volume: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/trade-volume", map[string]string{"pageSize": "100"}, true)
			apitest.AssertCovers(t, "trade-volume", raw, resp)
		}
	}
}
