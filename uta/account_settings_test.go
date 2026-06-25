package uta

import (
	"os"
	"testing"
)

// TestAccountSettingsWrite exercises the SAFE, reversible state-changing account
// endpoints against the live account. It is gated behind BITGET_TEST_WRITE=1 so
// it never runs by accident. The genuinely destructive endpoints in this file —
// adjust-account-mode (changes account-wide margin mode), set-margin (needs an
// open isolated position), set-deposit-account (changes deposit routing) and
// especially switch-account (downgrades UTA to classic) — are intentionally NOT
// executed here; they are implemented and compile-checked only.
func TestAccountSettingsWrite(t *testing.T) {
	if os.Getenv("BITGET_TEST_WRITE") == "" {
		t.Skip("set BITGET_TEST_WRITE=1 to run state-changing account tests")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}

	// set-leverage: the account holds no positions, so adjusting a futures
	// symbol's leverage is fully reversible and side-effect free.
	res, err := c.NewSetLeverageService(CategoryUSDTFutures, 5).
		SetSymbol("ETHUSDT").
		SetMarginMode(MarginModeCrossed).
		Do(ctx(t))
	if err != nil {
		t.Fatalf("set-leverage: %v", err)
	}
	t.Logf("set-leverage -> %q", *res)

	// switch-deduct: read the current BGB-deduction state, then write it back
	// unchanged so the account ends exactly as it started.
	cur, err := c.NewGetDeductInfoService().Do(ctx(t))
	if err != nil {
		t.Fatalf("deduct-info: %v", err)
	}
	state := cur.Deduct
	if state == "" {
		state = "off"
	}
	got, err := c.NewSwitchDeductService(state).Do(ctx(t))
	if err != nil {
		t.Fatalf("switch-deduct: %v", err)
	}
	t.Logf("switch-deduct(%s) -> %v", state, *got)
}
