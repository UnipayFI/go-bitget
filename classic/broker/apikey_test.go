package broker

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestBrokerAPIKey(t *testing.T) {
	c := NewBrokerClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Only the apikey-list GET is a safe read. It needs an ND-broker main
	// account; the shared sub-account key typically lacks the capability, which
	// we tolerate (the request path + signing are still proven).
	const subUid = "1"
	params := map[string]string{"subUid": subUid}

	resp, err := c.NewGetSubaccountAPIKeyListService(subUid).Do(cx)
	if err != nil {
		if apitest.Tolerable(t, "subaccount-apikey-list", err,
			"40029", "40099", "40034", "40068", "40014", "40037", "40054", "47001", "22001") {
			return
		}
		t.Fatalf("subaccount-apikey-list: %v", err)
	}
	raw := apitest.FetchRawGet(t, c, cx, "/api/v2/broker/manage/subaccount-apikey-list", params, true)
	apitest.AssertCovers(t, "subaccount-apikey-list", raw, resp)
}
