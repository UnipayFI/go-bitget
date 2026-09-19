package uta

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CFD (contract for difference) trading covers non-crypto underlyings such as
// gold, crude oil and forex indices. It runs on a separate account that must be
// opened from the Bitget website or app before these endpoints work -- there is
// no API for opening one. Every CFD endpoint is signed.
//
// A pair is named after the underlying plus a suffix that encodes the trading
// mode: ECN carries no suffix (XAUUSD), zero-fee ends in ".s" (XAUUSD.s) and
// Pro ends in ".pro" (XAUUSD.pro). An account is in exactly one mode at a time,
// and requests must spell the symbol the way that mode does.

// GetCFDTickersService -- GET /api/v3/cfd/market/tickers (UTA trade read)
//
// Returns ticker information for CFD trading pairs, optionally filtered to a
// single symbol.
type GetCFDTickersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDTickersService() *GetCFDTickersService {
	return &GetCFDTickersService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single CFD trading pair (e.g. XAUUSD).
func (s *GetCFDTickersService) SetSymbol(symbol string) *GetCFDTickersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetCFDTickersService) Do(ctx context.Context) ([]CFDTicker, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/market/tickers", s.params).WithSign()
	resp, err := request.Do[[]CFDTicker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CFDTicker struct {
	Symbol string `json:"symbol"`
	// AskOpenPriceChange and BidOpenPriceChange are the changes since the ask
	// and bid opening prices as decimal fractions (0.1 means 10%).
	AskOpenPriceChange decimal.Decimal `json:"askOpenPriceChange"`
	BidOpenPriceChange decimal.Decimal `json:"bidOpenPriceChange"`
	HighPrice          decimal.Decimal `json:"highPrice"`
	LowPrice           decimal.Decimal `json:"lowPrice"`
	Ask1               decimal.Decimal `json:"ask1"` // best ask price
	Bid1               decimal.Decimal `json:"bid1"` // best bid price
	QuoteTime          time.Time       `json:"quoteTime"`
}

// CFDCandle is one CFD candlestick row. Bitget returns each candle as a
// fixed-position JSON array ([ts, open, high, low, close]) -- five columns, with
// no volume or turnover, unlike the spot/futures Candle.
type CFDCandle struct {
	Ts    time.Time       `json:"ts"`    // array[0] -- candle start time (ms)
	Open  decimal.Decimal `json:"open"`  // array[1]
	High  decimal.Decimal `json:"high"`  // array[2]
	Low   decimal.Decimal `json:"low"`   // array[3]
	Close decimal.Decimal `json:"close"` // array[4]
}

// UnmarshalJSON decodes the 5-element positional array into named fields.
func (k *CFDCandle) UnmarshalJSON(data []byte) error {
	var row []string
	if err := common.JSONUnmarshal(data, &row); err != nil {
		return err
	}
	if len(row) < 5 {
		return fmt.Errorf("uta: cfd candle has %d columns, want 5", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("uta: cfd candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	for i, dst := range []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close} {
		d, err := decimal.NewFromString(row[i+1])
		if err != nil {
			return fmt.Errorf("uta: cfd candle column %d %q: %w", i+1, row[i+1], err)
		}
		*dst = d
	}
	return nil
}

// MarshalJSON re-emits the candle as the positional array Bitget sends, so the
// round-trip preserves the wire shape.
func (k CFDCandle) MarshalJSON() ([]byte, error) {
	row := []string{
		strconv.FormatInt(k.Ts.UnixMilli(), 10),
		k.Open.String(),
		k.High.String(),
		k.Low.String(),
		k.Close.String(),
	}
	return common.JSONMarshal(row)
}

// CFDKlineGranularity is a CFD candlestick interval. CFD uses its own, smaller
// vocabulary in lower case, so the spot/futures KlineGranularity values (1H, 1D,
// ...) are not interchangeable with these.
type CFDKlineGranularity string

const (
	CFDGranularity1m  CFDKlineGranularity = "1m"
	CFDGranularity15m CFDKlineGranularity = "15m"
	CFDGranularity1h  CFDKlineGranularity = "1h"
	CFDGranularity4h  CFDKlineGranularity = "4h"
	CFDGranularity1d  CFDKlineGranularity = "1d"
)

// GetCFDHistoryCandlesService -- GET /api/v3/cfd/market/history-candlestick (UTA trade read)
//
// Returns historical candlesticks for a CFD pair. Unlike the spot/futures
// candles, the series is per trade side (CFD quotes bid and ask separately) and
// the granularity vocabulary is limited to 1m, 15m, 1h, 4h and 1d. Queries are
// bounded to a 90-day range.
type GetCFDHistoryCandlesService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDHistoryCandlesService(symbol string, interval CFDKlineGranularity, side Side) *GetCFDHistoryCandlesService {
	return &GetCFDHistoryCandlesService{c: c, params: map[string]string{
		"symbol":   symbol,
		"interval": string(interval),
		"side":     string(side),
	}}
}

// SetStartTime filters candles at or after t (90-day lookback window).
func (s *GetCFDHistoryCandlesService) SetStartTime(t time.Time) *GetCFDHistoryCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t (90-day lookback window).
func (s *GetCFDHistoryCandlesService) SetEndTime(t time.Time) *GetCFDHistoryCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of candles returned (default 100, max 100).
func (s *GetCFDHistoryCandlesService) SetLimit(limit int) *GetCFDHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCFDHistoryCandlesService) Do(ctx context.Context) ([]CFDCandle, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/market/history-candlestick", s.params).WithSign()
	resp, err := request.Do[[]CFDCandle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
