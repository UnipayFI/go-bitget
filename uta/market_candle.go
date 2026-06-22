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

// Candle is one candlestick row. Bitget returns each candle as a fixed-position
// JSON array ([ts, open, high, low, close, volume, turnover]); Candle parses
// that array into named fields and re-emits the same array shape on marshal.
type Candle struct {
	Ts       time.Time       `json:"ts"`       // array[0] -- candle start time (ms)
	Open     decimal.Decimal `json:"open"`     // array[1]
	High     decimal.Decimal `json:"high"`     // array[2]
	Low      decimal.Decimal `json:"low"`      // array[3]
	Close    decimal.Decimal `json:"close"`    // array[4]
	Volume   decimal.Decimal `json:"volume"`   // array[5] -- base coin volume
	Turnover decimal.Decimal `json:"turnover"` // array[6] -- quote coin turnover
}

// UnmarshalJSON decodes the 7-element positional array into named fields.
func (k *Candle) UnmarshalJSON(data []byte) error {
	var row []string
	if err := common.JSONUnmarshal(data, &row); err != nil {
		return err
	}
	if len(row) < 7 {
		return fmt.Errorf("uta: candle has %d columns, want 7", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("uta: candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	for i, dst := range []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close, &k.Volume, &k.Turnover} {
		d, err := decimal.NewFromString(row[i+1])
		if err != nil {
			return fmt.Errorf("uta: candle column %d %q: %w", i+1, row[i+1], err)
		}
		*dst = d
	}
	return nil
}

// MarshalJSON re-emits the candle as the positional array Bitget sends, so the
// round-trip preserves the wire shape.
func (k Candle) MarshalJSON() ([]byte, error) {
	row := []string{
		strconv.FormatInt(k.Ts.UnixMilli(), 10),
		k.Open.String(),
		k.High.String(),
		k.Low.String(),
		k.Close.String(),
		k.Volume.String(),
		k.Turnover.String(),
	}
	return common.JSONMarshal(row)
}

// GetCandlesService -- GET /api/v3/market/candles
//
// Returns the most recent candlesticks for a symbol at the given interval.
type GetCandlesService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCandlesService(category Category, symbol string, interval KlineGranularity) *GetCandlesService {
	return &GetCandlesService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
		"interval": string(interval),
	}}
}

// SetStartTime filters candles at or after t.
func (s *GetCandlesService) SetStartTime(t time.Time) *GetCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t.
func (s *GetCandlesService) SetEndTime(t time.Time) *GetCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetType selects the price series (market, mark, index, premium; default
// market).
func (s *GetCandlesService) SetType(klineType KlineType) *GetCandlesService {
	s.params["type"] = string(klineType)
	return s
}

// SetLimit caps the number of candles returned (max 100).
func (s *GetCandlesService) SetLimit(limit int) *GetCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetHistoryCandlesService -- GET /api/v3/market/history-candles
//
// Returns historical candlesticks for a symbol at the given interval (max
// 90-day range per query).
type GetHistoryCandlesService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetHistoryCandlesService(category Category, symbol string, interval KlineGranularity) *GetHistoryCandlesService {
	return &GetHistoryCandlesService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
		"interval": string(interval),
	}}
}

// SetStartTime filters candles at or after t.
func (s *GetHistoryCandlesService) SetStartTime(t time.Time) *GetHistoryCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t.
func (s *GetHistoryCandlesService) SetEndTime(t time.Time) *GetHistoryCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetType selects the price series (market, mark, index, premium; default
// market).
func (s *GetHistoryCandlesService) SetType(klineType KlineType) *GetHistoryCandlesService {
	s.params["type"] = string(klineType)
	return s
}

// SetLimit caps the number of candles returned (max 100).
func (s *GetHistoryCandlesService) SetLimit(limit int) *GetHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/history-candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
