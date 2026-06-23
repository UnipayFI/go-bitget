package common

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
	"github.com/shopspring/decimal"
)

func TestConvert(t *testing.T) {
	c := NewCommonClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))

	// Convert currencies (no params).
	t.Run("currencies", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetConvertCurrenciesService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "convert/currencies", err, "40099", "40034") {
				return
			}
			t.Fatalf("convert currencies: %v", err)
		}
		t.Logf("convert currencies: %d coins", len(resp))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/convert/currencies", nil, true)
		apitest.AssertCovers(t, "convert/currencies", raw, resp)
	})

	// Quoted price (RFQ; read-only, does not execute the conversion).
	t.Run("quoted-price", func(t *testing.T) {
		cx := apitest.Ctx(t)
		// Exactly one of fromCoinSize/toCoinSize must be supplied (code 20002).
		params := map[string]string{"fromCoin": "USDT", "toCoin": "BTC", "fromCoinSize": "10"}
		resp, err := c.NewGetConvertQuotedPriceService("USDT", "BTC").
			SetFromCoinSize(decimal.RequireFromString("10")).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "convert/quoted-price", err, "40099", "40034", "47001") {
				return
			}
			t.Fatalf("convert quoted-price: %v", err)
		}
		t.Logf("convert quoted-price: traceId=%s price=%s", resp.TraceID, resp.CnvtPrice)
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/convert/quoted-price", params, true)
		apitest.AssertCovers(t, "convert/quoted-price", raw, resp)
	})

	// Convert history (startTime/endTime required, max 90-day span).
	t.Run("convert-record", func(t *testing.T) {
		cx := apitest.Ctx(t)
		end := time.Now()
		start := end.Add(-90 * 24 * time.Hour)
		resp, err := c.NewGetConvertRecordService(start, end).Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "convert/convert-record", err, "40099", "40034") {
				return
			}
			t.Fatalf("convert record: %v", err)
		}
		t.Logf("convert record: %d records endId=%s", len(resp.DataList), resp.EndID)
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/convert/convert-record", map[string]string{
			"startTime": strconv.FormatInt(start.UnixMilli(), 10),
			"endTime":   strconv.FormatInt(end.UnixMilli(), 10),
		}, true)
		apitest.AssertCovers(t, "convert/convert-record", raw, resp)
	})

	// BGB convert coin list (no params).
	t.Run("bgb-convert-coin-list", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetBGBConvertCoinListService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "convert/bgb-convert-coin-list", err, "40099", "40034") {
				return
			}
			t.Fatalf("bgb convert coin list: %v", err)
		}
		t.Logf("bgb convert coin list: %d coins", len(resp.CoinList))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/convert/bgb-convert-coin-list", nil, true)
		apitest.AssertCovers(t, "convert/bgb-convert-coin-list", raw, resp)
	})

	// BGB convert history (no required params; live path is plural "records").
	t.Run("bgb-convert-records", func(t *testing.T) {
		cx := apitest.Ctx(t)
		resp, err := c.NewGetBGBConvertRecordService().Do(cx)
		if err != nil {
			if apitest.Tolerable(t, "convert/bgb-convert-records", err, "40099", "40034") {
				return
			}
			t.Fatalf("bgb convert records: %v", err)
		}
		t.Logf("bgb convert records: %d records", len(resp))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/convert/bgb-convert-records", nil, true)
		apitest.AssertCovers(t, "convert/bgb-convert-records", raw, resp)
	})
}
