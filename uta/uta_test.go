package uta

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/client"
)

// testClient builds an authenticated UTA client from environment variables.
// Tests that need private endpoints skip themselves when credentials are
// absent so the suite stays runnable without secrets.
func testClient(t *testing.T) *UTAClient {
	t.Helper()
	apiKey := os.Getenv("BITGET_API_KEY")
	apiSecret := os.Getenv("BITGET_API_SECRET")
	passphrase := os.Getenv("BITGET_PASSPHRASE")
	if apiKey == "" || apiSecret == "" || passphrase == "" {
		t.Skip("BITGET_API_KEY/SECRET/PASSPHRASE not set; skipping private test")
	}
	opts := []client.Options{
		client.WithAuth(apiKey, apiSecret, passphrase),
	}
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewUTAClient(opts...)
}

func testPublicClient() *UTAClient {
	opts := []client.Options{}
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	return NewUTAClient(opts...)
}

func ctx(t *testing.T) context.Context {
	t.Helper()
	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	t.Cleanup(cancel)
	return c
}

func TestServerTime(t *testing.T) {
	c := testPublicClient()
	resp, err := c.NewGetServerTimeService().Do(ctx(t))
	if err != nil {
		t.Fatalf("server time: %v", err)
	}
	t.Logf("serverTime=%s (%d)", resp.ServerTime, resp.ServerTime.UnixMilli())
	if resp.ServerTime.IsZero() {
		t.Fatal("server time is zero")
	}
}

func TestAccountAssets(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)
	resp, err := c.NewGetAccountAssetsService().Do(cx)
	if err != nil {
		t.Fatalf("account assets: %v", err)
	}
	t.Logf("accountEquity=%s usdtEquity=%s assets=%d",
		resp.AccountEquity, resp.USDTEquity, len(resp.Assets))
	for _, a := range resp.Assets {
		t.Logf("  %s balance=%s available=%s usdValue=%s", a.Coin, a.Balance, a.Available, a.USDValue)
	}
	raw := fetchRawGet(t, c, cx, "/api/v3/account/assets", nil, true)
	assertCovers(t, "account/assets", raw, resp)
}

func TestInstruments(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)
	params := map[string]string{"category": string(CategoryUSDTFutures)}
	list, err := c.NewGetInstrumentsService(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("instruments: %v", err)
	}
	if len(list) == 0 {
		t.Fatal("no instruments returned")
	}
	t.Logf("first: %+v", list[0])
	raw := fetchRawGet(t, c, cx, "/api/v3/market/instruments", params, false)
	assertCovers(t, "market/instruments", raw, list)
}
