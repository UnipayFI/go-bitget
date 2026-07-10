package uta

import (
	"testing"

	"github.com/UnipayFI/go-bitget/client"
)

// TestRealityMarketData exercises the authenticated reality-orderbook and
// reality-fills endpoints. Both require Reality (US stock) access; a whitelist
// or permission error skips the test rather than failing it.
func TestRealityMarketData(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)
	symbol := "RAAPLUSDT"

	// Reality order book.
	ob, err := c.NewGetRealityOrderBookService(symbol).Do(cx)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			t.Skipf("reality orderbook not available (whitelist required): %v", apiErr)
		}
		t.Fatalf("reality orderbook: %v", err)
	}
	t.Logf("reality orderbook: symbol=%s asks=%d bids=%d ts=%s", ob.Symbol, len(ob.Asks), len(ob.Bids), ob.Ts)
	raw := fetchRawGet(t, c, cx, "/api/v3/account/reality-orderbook", map[string]string{"symbol": symbol}, true)
	assertCovers(t, "account/reality-orderbook", raw, ob)

	// Reality fills.
	fills, err := c.NewGetRealityFillsService(symbol).Do(cx)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			t.Skipf("reality fills not available (whitelist required): %v", apiErr)
		}
		t.Fatalf("reality fills: %v", err)
	}
	t.Logf("reality fills: %d", len(fills))
	raw = fetchRawGet(t, c, cx, "/api/v3/account/reality-fills", map[string]string{"symbol": symbol}, true)
	assertCovers(t, "account/reality-fills", raw, fills)
}
