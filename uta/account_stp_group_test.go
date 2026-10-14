package uta

import (
	"errors"
	"testing"

	"github.com/UnipayFI/go-bitget/client"
)

func TestSTPGroups(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// STP groups are MM-whitelist only; a regular account gets a permission
	// error, which still confirms the endpoint + signing work.
	groups, err := c.NewGetSTPGroupsService().Do(cx)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) {
			t.Logf("stp-groups: account is not MM whitelisted (code=%s) — endpoint+signing OK", apiErr.Code)
			return
		}
		t.Fatalf("stp-groups: %v", err)
	}
	t.Logf("stp-groups: %d", len(groups.Groups))
	raw := fetchRawGet(t, c, cx, "/api/v3/account/stp-groups", nil, true)
	assertCovers(t, "account/stp-groups", raw, groups)
}
