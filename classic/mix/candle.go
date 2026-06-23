package mix

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// KlineGranularity is the candlestick interval. The plain values are
// exchange-local; the "utc" suffixed variants align the daily/weekly/monthly
// boundary to UTC+0 instead of UTC+8.
type KlineGranularity string

const (
	KlineGranularity1m     KlineGranularity = "1m"
	KlineGranularity3m     KlineGranularity = "3m"
	KlineGranularity5m     KlineGranularity = "5m"
	KlineGranularity15m    KlineGranularity = "15m"
	KlineGranularity30m    KlineGranularity = "30m"
	KlineGranularity1H     KlineGranularity = "1H"
	KlineGranularity4H     KlineGranularity = "4H"
	KlineGranularity6H     KlineGranularity = "6H"
	KlineGranularity12H    KlineGranularity = "12H"
	KlineGranularity1D     KlineGranularity = "1D"
	KlineGranularity3D     KlineGranularity = "3D"
	KlineGranularity1W     KlineGranularity = "1W"
	KlineGranularity1M     KlineGranularity = "1M"
	KlineGranularity6Hutc  KlineGranularity = "6Hutc"
	KlineGranularity12Hutc KlineGranularity = "12Hutc"
	KlineGranularity1Dutc  KlineGranularity = "1Dutc"
	KlineGranularity3Dutc  KlineGranularity = "3Dutc"
	KlineGranularity1Wutc  KlineGranularity = "1Wutc"
	KlineGranularity1Mutc  KlineGranularity = "1Mutc"
)

// KlineType selects which price series the candles endpoint returns (the
// traded market price, the mark price, or the index price).
type KlineType string

const (
	KlineTypeMarket KlineType = "MARKET"
	KlineTypeMark   KlineType = "MARK"
	KlineTypeIndex  KlineType = "INDEX"
)

// Candle is one candlestick row. Bitget returns each candle as a fixed-position
// JSON array ([ts, open, high, low, close, volume, turnover]); Candle parses
// that array into named fields and re-emits the same array shape on marshal.
// Index- and mark-price candles share this 7-column shape but report zero for
// the volume and turnover columns.
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
		return fmt.Errorf("mix: candle has %d columns, want 7", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("mix: candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	for i, dst := range []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close, &k.Volume, &k.Turnover} {
		d, err := decimal.NewFromString(row[i+1])
		if err != nil {
			return fmt.Errorf("mix: candle column %d %q: %w", i+1, row[i+1], err)
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

// GetCandlesService -- GET /api/v2/mix/market/candles (public)
//
// Returns the most recent candlesticks for a contract at the given interval.
type GetCandlesService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetCandlesService(symbol string, productType ProductType, granularity KlineGranularity) *GetCandlesService {
	return &GetCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"granularity": string(granularity),
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

// SetKlineType selects the price series (MARKET default, MARK, INDEX).
func (s *GetCandlesService) SetKlineType(klineType KlineType) *GetCandlesService {
	s.params["kLineType"] = string(klineType)
	return s
}

// SetLimit caps the number of candles returned (default 100, max 1000).
func (s *GetCandlesService) SetLimit(limit int) *GetCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetHistoryCandlesService -- GET /api/v2/mix/market/history-candles (public)
//
// Returns historical candlesticks for a contract at the given interval (max
// 90-day range per query).
type GetHistoryCandlesService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetHistoryCandlesService(symbol string, productType ProductType, granularity KlineGranularity) *GetHistoryCandlesService {
	return &GetHistoryCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"granularity": string(granularity),
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

// SetLimit caps the number of candles returned (default 100, max 200).
func (s *GetHistoryCandlesService) SetLimit(limit int) *GetHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/history-candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetHistoryIndexCandlesService -- GET /api/v2/mix/market/history-index-candles (public)
//
// Returns historical index-price candlesticks for a contract at the given
// interval (max 90-day range per query). The volume and turnover columns are
// always zero for index-price candles.
type GetHistoryIndexCandlesService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetHistoryIndexCandlesService(symbol string, productType ProductType, granularity KlineGranularity) *GetHistoryIndexCandlesService {
	return &GetHistoryIndexCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"granularity": string(granularity),
	}}
}

// SetStartTime filters candles at or after t.
func (s *GetHistoryIndexCandlesService) SetStartTime(t time.Time) *GetHistoryIndexCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t.
func (s *GetHistoryIndexCandlesService) SetEndTime(t time.Time) *GetHistoryIndexCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of candles returned (default 100, max 200).
func (s *GetHistoryIndexCandlesService) SetLimit(limit int) *GetHistoryIndexCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryIndexCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/history-index-candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetHistoryMarkCandlesService -- GET /api/v2/mix/market/history-mark-candles (public)
//
// Returns historical mark-price candlesticks for a contract at the given
// interval (max 90-day range per query). The volume and turnover columns are
// always zero for mark-price candles.
type GetHistoryMarkCandlesService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetHistoryMarkCandlesService(symbol string, productType ProductType, granularity KlineGranularity) *GetHistoryMarkCandlesService {
	return &GetHistoryMarkCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"granularity": string(granularity),
	}}
}

// SetStartTime filters candles at or after t.
func (s *GetHistoryMarkCandlesService) SetStartTime(t time.Time) *GetHistoryMarkCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t.
func (s *GetHistoryMarkCandlesService) SetEndTime(t time.Time) *GetHistoryMarkCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of candles returned (default 100, max 200).
func (s *GetHistoryMarkCandlesService) SetLimit(limit int) *GetHistoryMarkCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryMarkCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/history-mark-candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
