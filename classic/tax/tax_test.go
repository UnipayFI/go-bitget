package tax

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestTax(t *testing.T) {
	c := NewTaxClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	// Tax records require a startTime+endTime window of at most 30 days.
	end := time.Now()
	start := end.Add(-7 * 24 * time.Hour)
	window := map[string]string{
		"startTime": strconv.FormatInt(start.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
	}

	// Codes the permission-limited shared sub-account key may return for these
	// product lines; treat as pass (endpoint + signing proven).
	okCodes := []string{"40068", "40014", "40054", "40099", "40034", "40029", "40037", "47001", "22001"}

	// spot-record
	{
		spot, err := c.NewGetSpotRecordService(start, end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "spot-record", err, okCodes...) {
				t.Fatalf("spot-record: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/tax/spot-record", window, true)
			apitest.AssertCovers(t, "spot-record", raw, spot)
		}
	}

	// future-record
	{
		fut, err := c.NewGetFutureRecordService(start, end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "future-record", err, okCodes...) {
				t.Fatalf("future-record: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/tax/future-record", window, true)
			apitest.AssertCovers(t, "future-record", raw, fut)
		}
	}

	// margin-record
	{
		mgn, err := c.NewGetMarginRecordService(start, end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "margin-record", err, okCodes...) {
				t.Fatalf("margin-record: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/tax/margin-record", window, true)
			apitest.AssertCovers(t, "margin-record", raw, mgn)
		}
	}

	// p2p-record
	{
		p2p, err := c.NewGetP2PRecordService(start, end).Do(cx)
		if err != nil {
			if !apitest.Tolerable(t, "p2p-record", err, okCodes...) {
				t.Fatalf("p2p-record: %v", err)
			}
		} else {
			raw := apitest.FetchRawGet(t, c, cx, "/api/v2/tax/p2p-record", window, true)
			apitest.AssertCovers(t, "p2p-record", raw, p2p)
		}
	}
}
