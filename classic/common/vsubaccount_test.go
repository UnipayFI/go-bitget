package common

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// TestVirtualSubaccount live-tests only the safe READ endpoints of the virtual
// sub-account group: virtual-subaccount-list and (when a sub-account exists)
// virtual-subaccount-apikey-list. Both are capability-gated — an account
// without sub-account management enabled returns a permission error, which is
// tolerated. The POST create/modify endpoints are state-changing and never
// executed.
func TestVirtualSubaccount(t *testing.T) {
	c := NewCommonClient(apitest.AuthOptions(t)...)
	cx := apitest.Ctx(t)
	_ = c.SyncServerTime(cx)

	// virtual-subaccount-list (signed GET, no required params).
	listResp, err := c.NewGetVirtualSubaccountListService().Do(cx)
	if err != nil {
		if apitest.Tolerable(t, "virtual-subaccount-list", err, "40099", "40034", "40037", "40806", "00172", "40014", "40068") {
			return
		}
		t.Fatalf("virtual-subaccount-list: %v", err)
	}
	raw := apitest.FetchRawGet(t, c, cx, "/api/v2/user/virtual-subaccount-list", nil, true)
	apitest.AssertCovers(t, "virtual-subaccount-list", raw, listResp)

	// virtual-subaccount-apikey-list needs a sub-account uid; only test it when
	// the list returned at least one.
	if listResp == nil || len(listResp.SubAccountList) == 0 {
		t.Log("virtual-subaccount-apikey-list: no virtual sub-account to query; skipping")
		return
	}
	uid := listResp.SubAccountList[0].SubAccountUid
	keyResp, err := c.NewGetVirtualSubaccountApikeyListService(uid).Do(cx)
	if err != nil {
		if apitest.Tolerable(t, "virtual-subaccount-apikey-list", err, "40099", "40034", "40037", "40806", "00172", "40014", "40068") {
			return
		}
		t.Fatalf("virtual-subaccount-apikey-list: %v", err)
	}
	rawKey := apitest.FetchRawGet(t, c, cx, "/api/v2/user/virtual-subaccount-apikey-list", map[string]string{"subAccountUid": uid}, true)
	apitest.AssertCovers(t, "virtual-subaccount-apikey-list", rawKey, keyResp)
}
