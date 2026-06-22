package uta

import (
	"strconv"
	"testing"
	"time"
)

func TestTax(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)

	records, err := c.NewGetTaxRecordsService(string(CategoryUSDTFutures), start, end).Do(cx)
	if err != nil {
		t.Fatalf("tax records: %v", err)
	}
	t.Logf("tax records: %d", len(records))
	raw := fetchRawGet(t, c, cx, "/api/v3/tax/records", map[string]string{
		"bizType":   string(CategoryUSDTFutures),
		"startTime": strconv.FormatInt(start.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
	}, true)
	assertCovers(t, "tax/records", raw, records)
}
