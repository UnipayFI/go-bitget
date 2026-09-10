package uta

import (
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/common"
)

func TestCFDCandleUnmarshal(t *testing.T) {
	const payload = `[["1597026383085","2610.34","2650.12","2605.10","2620.55"]]`
	var candles []CFDCandle
	if err := common.JSONUnmarshal([]byte(payload), &candles); err != nil {
		t.Fatal(err)
	}
	if len(candles) != 1 {
		t.Fatalf("candles = %d, want 1", len(candles))
	}
	k := candles[0]
	if got := k.Ts.UnixMilli(); got != 1597026383085 {
		t.Fatalf("ts = %d, want 1597026383085", got)
	}
	if k.Open.String() != "2610.34" || k.High.String() != "2650.12" ||
		k.Low.String() != "2605.1" || k.Close.String() != "2620.55" {
		t.Fatalf("ohlc = %s/%s/%s/%s", k.Open, k.High, k.Low, k.Close)
	}

	// The wire shape must survive the round-trip as the same positional array
	// (decimals come back normalised, so 2605.10 re-emits as 2605.1).
	const want = `[["1597026383085","2610.34","2650.12","2605.1","2620.55"]]`
	out, err := common.JSONMarshal(candles)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != want {
		t.Fatalf("re-marshal = %s, want %s", out, want)
	}
}

// TestCFD exercises the read side of the CFD endpoints. They all require a CFD
// account, which is opened separately from the Bitget website or app, so the
// whole test is gated on fund-detail answering.
func TestCFD(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	fund, err := c.NewGetCFDFundDetailService().Do(cx)
	if err != nil {
		t.Skipf("no CFD account on this UID (%v); skipping CFD endpoints", err)
	}
	t.Logf("fund detail: %+v", fund)
	raw := fetchRawGet(t, c, cx, "/api/v3/cfd/account/fund-detail", nil, true)
	assertCovers(t, "cfd/account/fund-detail", raw, fund)

	// Instruments; the first one seeds the per-symbol queries below.
	instruments, err := c.NewGetCFDInstrumentsService().Do(cx)
	if err != nil {
		t.Fatalf("cfd instruments: %v", err)
	}
	t.Logf("cfd instruments: %d", len(instruments))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/account/instruments", nil, true)
	assertCovers(t, "cfd/account/instruments", raw, instruments)
	if len(instruments) == 0 {
		t.Skip("no CFD instruments listed; skipping the symbol-scoped endpoints")
	}
	symbol := instruments[0].Symbol

	// Tickers.
	tickers, err := c.NewGetCFDTickersService().SetSymbol(symbol).Do(cx)
	if err != nil {
		t.Fatalf("cfd tickers: %v", err)
	}
	t.Logf("cfd tickers: %d", len(tickers))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/market/tickers", map[string]string{"symbol": symbol}, true)
	assertCovers(t, "cfd/market/tickers", raw, tickers)

	// History candles (positional arrays, so assertCovers has no keys to diff).
	candles, err := c.NewGetCFDHistoryCandlesService(symbol, CFDGranularity1m, SideBuy).SetLimit(10).Do(cx)
	if err != nil {
		t.Fatalf("cfd history candles: %v", err)
	}
	t.Logf("cfd candles: %d", len(candles))

	// Unfilled orders.
	unfilled, err := c.NewGetCFDUnfilledOrdersService(symbol).Do(cx)
	if err != nil {
		t.Fatalf("cfd unfilled orders: %v", err)
	}
	t.Logf("cfd unfilled orders: %d", len(unfilled))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/trade/unfilled-order", map[string]string{"symbol": symbol}, true)
	assertCovers(t, "cfd/trade/unfilled-order", raw, unfilled)

	// Order history (30-day window inside the 90-day lookback).
	end := time.Now()
	start := end.AddDate(0, 0, -7)
	history, err := c.NewGetCFDOrderHistoryService(symbol).SetStartTime(start).SetEndTime(end).Do(cx)
	if err != nil {
		t.Fatalf("cfd order history: %v", err)
	}
	t.Logf("cfd history orders: %d", len(history.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/trade/history-order", map[string]string{"symbol": symbol}, true)
	assertCovers(t, "cfd/trade/history-order", raw, history)

	// Current positions.
	positions, err := c.NewGetCFDPositionsService().Do(cx)
	if err != nil {
		t.Fatalf("cfd positions: %v", err)
	}
	t.Logf("cfd positions: %d", len(positions))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/trade/current-positions", nil, true)
	assertCovers(t, "cfd/trade/current-positions", raw, positions)

	// Transfer records.
	transfers, err := c.NewGetCFDTransferRecordsService().Do(cx)
	if err != nil {
		t.Fatalf("cfd transfer records: %v", err)
	}
	t.Logf("cfd transfer records: %d", len(transfers.List))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/account/transfer-records", nil, true)
	assertCovers(t, "cfd/account/transfer-records", raw, transfers)

	// Financial records.
	records, err := c.NewGetCFDFinancialRecordsService().Do(cx)
	if err != nil {
		t.Fatalf("cfd financial records: %v", err)
	}
	t.Logf("cfd financial records: %d", len(records))
	raw = fetchRawGet(t, c, cx, "/api/v3/cfd/account/financial-records", nil, true)
	assertCovers(t, "cfd/account/financial-records", raw, records)
}
