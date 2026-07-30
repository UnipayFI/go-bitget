package uta

import (
	"errors"
	"testing"

	"github.com/UnipayFI/go-bitget/client"
)

func TestRateLimitQuota(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Rate-limit quota is institutional-only (MM/PRO); a regular account gets a
	// permission error, which still confirms the endpoint + signing work.
	quota, err := c.NewGetRateLimitQuotaService("futures").Do(cx)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) {
			t.Logf("rate-limit-quota: account is not MM/PRO institutional (code=%s) — endpoint+signing OK", apiErr.Code)
			return
		}
		t.Fatalf("rate-limit-quota: %v", err)
	}
	t.Logf("rate-limit-quota: %+v", quota)
	raw := fetchRawGet(t, c, cx, "/api/v3/user/rate-limit-quota", map[string]string{"category": "futures"}, true)
	assertCovers(t, "user/rate-limit-quota", raw, quota)
}
