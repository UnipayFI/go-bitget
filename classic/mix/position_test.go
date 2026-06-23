package mix

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestMixPosition(t *testing.T) {
	c := NewMixClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	cx := apitest.Ctx(t)

	const (
		productType = string(ProductTypeUSDTFutures)
		symbol      = "BTCUSDT"
		marginCoin  = "USDT"
	)
	okCodes := []string{"40099", "40034", "40054", "40068", "40014", "40037", "22001", "47001"}

	// single-position: requires symbol + productType + marginCoin.
	t.Run("single-position", func(t *testing.T) {
		resp, err := c.NewGetSinglePositionService(ProductTypeUSDTFutures, symbol, marginCoin).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "single-position", err, okCodes...) {
				return
			}
			t.Fatalf("single-position: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/position/single-position", map[string]string{
			"productType": productType,
			"symbol":      symbol,
			"marginCoin":  marginCoin,
		}, true)
		apitest.AssertCovers(t, "single-position", raw, resp)
	})

	// all-position: requires productType.
	t.Run("all-position", func(t *testing.T) {
		resp, err := c.NewGetAllPositionService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "all-position", err, okCodes...) {
				return
			}
			t.Fatalf("all-position: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/position/all-position", map[string]string{
			"productType": productType,
		}, true)
		apitest.AssertCovers(t, "all-position", raw, resp)
	})

	// adlRank: requires productType.
	t.Run("adlRank", func(t *testing.T) {
		resp, err := c.NewGetPositionADLRankService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "adlRank", err, okCodes...) {
				return
			}
			t.Fatalf("adlRank: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/position/adlRank", map[string]string{
			"productType": productType,
		}, true)
		apitest.AssertCovers(t, "adlRank", raw, resp)
	})

	// history-position: requires productType.
	t.Run("history-position", func(t *testing.T) {
		resp, err := c.NewGetHistoryPositionService(ProductTypeUSDTFutures).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "history-position", err, okCodes...) {
				return
			}
			t.Fatalf("history-position: %v", err)
		}
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/mix/position/history-position", map[string]string{
			"productType": productType,
		}, true)
		apitest.AssertCovers(t, "history-position", raw, resp)
	})
}
