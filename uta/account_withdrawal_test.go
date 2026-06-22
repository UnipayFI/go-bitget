package uta

import (
	"strconv"
	"testing"
	"time"
)

func TestAccountWithdrawal(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	endTime := time.Now()
	startTime := endTime.Add(-30 * 24 * time.Hour)
	startMs := strconv.FormatInt(startTime.UnixMilli(), 10)
	endMs := strconv.FormatInt(endTime.UnixMilli(), 10)

	// Withdrawal records.
	records, err := c.NewGetWithdrawalRecordsService(startTime, endTime).Do(cx)
	if err != nil {
		t.Fatalf("withdrawal records: %v", err)
	}
	t.Logf("withdrawal records: %d", len(records))
	raw := fetchRawGet(t, c, cx, "/api/v3/account/withdrawal-records",
		map[string]string{"startTime": startMs, "endTime": endMs}, true)
	assertCovers(t, "account/withdrawal-records", raw, records)

	// Withdraw address book.
	book, err := c.NewGetWithdrawAddressService().Do(cx)
	if err != nil {
		t.Fatalf("withdraw address: %v", err)
	}
	t.Logf("withdraw address: %d cursor=%s", len(book.AddressList), book.Cursor)
	raw = fetchRawGet(t, c, cx, "/api/v3/account/withdraw-address", nil, true)
	assertCovers(t, "account/withdraw-address", raw, book)
}
