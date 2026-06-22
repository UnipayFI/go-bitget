package uta

import (
	"strconv"
	"testing"
	"time"
)

func TestAccountDeposit(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	end := time.Now()
	start := end.Add(-30 * 24 * time.Hour)

	// Deposit records (safe time-window listing, no pre-existing state needed).
	records, err := c.NewGetDepositRecordsService(start, end).Do(cx)
	if err != nil {
		t.Fatalf("deposit records: %v", err)
	}
	t.Logf("deposit records: %d", len(records))
	raw := fetchRawGet(t, c, cx, "/api/v3/account/deposit-records", map[string]string{
		"startTime": strconv.FormatInt(start.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
	}, true)
	assertCovers(t, "account/deposit-records", raw, records)
}
