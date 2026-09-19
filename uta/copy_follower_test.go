package uta

import (
	"os"
	"testing"
)

// TestCopyFollower reads back a followed copy trading project. Bitget has no
// "list followed projects" endpoint, so the project ID has to come from the
// environment: set BITGET_COPY_PROJECT_ID.
func TestCopyFollower(t *testing.T) {
	projectID := os.Getenv("BITGET_COPY_PROJECT_ID")
	if projectID == "" {
		t.Skip("BITGET_COPY_PROJECT_ID not set; skipping copy follower test")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)
	params := map[string]string{"projectId": projectID}

	// Copy settings.
	settings, err := c.NewGetCopySettingsService(projectID).Do(cx)
	if err != nil {
		t.Fatalf("copy settings: %v", err)
	}
	t.Logf("copy settings: %+v", settings)
	raw := fetchRawGet(t, c, cx, "/api/v3/copy/futures/copy-settings", params, true)
	assertCovers(t, "copy/futures/copy-settings", raw, settings)

	// Current copy.
	current, err := c.NewGetCurrentCopyService(projectID).Do(cx)
	if err != nil {
		t.Fatalf("current copy: %v", err)
	}
	t.Logf("current copy: %+v", current)
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/current-copy", params, true)
	assertCovers(t, "copy/futures/current-copy", raw, current)

	// Transfer records.
	records, err := c.NewGetCopyFollowerTransferRecordService(projectID).SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("copy transfer record: %v", err)
	}
	t.Logf("copy transfer records: %d", len(records.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/copy-transfer-record",
		map[string]string{"projectId": projectID, "limit": "20"}, true)
	assertCovers(t, "copy/futures/copy-transfer-record", raw, records)

	// Profit details.
	details, err := c.NewGetCopyFollowerProfitDetailsService(projectID).SetLimit("20").Do(cx)
	if err != nil {
		t.Fatalf("copy profit details: %v", err)
	}
	t.Logf("copy profit details: %d", len(details.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/copy-profit-details",
		map[string]string{"projectId": projectID, "limit": "20"}, true)
	assertCovers(t, "copy/futures/copy-profit-details", raw, details)

	// Current positions.
	positions, err := c.NewGetCopyCurrentPositionsService(projectID).Do(cx)
	if err != nil {
		t.Fatalf("copy current positions: %v", err)
	}
	t.Logf("copy current positions: %d", len(positions.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/current-positions", params, true)
	assertCovers(t, "copy/futures/current-positions", raw, positions)

	// Current TP/SL orders.
	tpsl, err := c.NewGetCopyCurrentTPSLOrdersService(projectID).Do(cx)
	if err != nil {
		t.Fatalf("copy current tpsl orders: %v", err)
	}
	t.Logf("copy current tpsl orders: %d", len(tpsl.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/current-tpsl-orders", params, true)
	assertCovers(t, "copy/futures/current-tpsl-orders", raw, tpsl)

	// TP/SL order history.
	history, err := c.NewGetCopyTPSLOrderHistoryService(projectID).Do(cx)
	if err != nil {
		t.Fatalf("copy tpsl order history: %v", err)
	}
	t.Logf("copy tpsl order history: %d", len(history.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/copy/futures/tpsl-order-history", params, true)
	assertCovers(t, "copy/futures/tpsl-order-history", raw, history)
}
