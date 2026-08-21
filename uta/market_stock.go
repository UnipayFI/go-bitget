package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSplitRecordsService -- GET /api/v3/market/split-records
//
// Returns the stock split and reverse split records of the RWA (Reality)
// trading pairs. The endpoint takes no parameters and covers every symbol with
// a recorded split.
type GetSplitRecordsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSplitRecordsService() *GetSplitRecordsService {
	return &GetSplitRecordsService{c: c}
}

func (s *GetSplitRecordsService) Do(ctx context.Context) ([]SplitRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/split-records")
	resp, err := request.Do[[]SplitRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SplitRecord is a single stock split or reverse split of an RWA trading pair.
type SplitRecord struct {
	Symbol string `json:"symbol"`
	Type   string `json:"type"`   // split, reverse_split
	Status string `json:"status"` // pending, ongoing, completed
	// AdjustmentRatio is the share adjustment factor: greater than 1 for a
	// split, smaller than 1 for a reverse split.
	AdjustmentRatio decimal.Decimal `json:"adjustmentRatio"`
	// ExDividendDate is the effective date. Bitget returns it as yyyyMMdd
	// (e.g. "20260715"), not the yyyy-MM-dd shown in the docs.
	ExDividendDate         string    `json:"exDividendDate"`
	ExDividendDateTimezone string    `json:"exDividendDateTimezone"` // ET
	TradingHaltStartTime   time.Time `json:"tradingHaltStartTime"`
	TradingHaltEndTime     time.Time `json:"tradingHaltEndTime"`
}

// GetStockInfoService -- GET /api/v3/reality/market/stock-info
//
// Returns the underlying stock information of the Reality trading pairs.
type GetStockInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetStockInfoService() *GetStockInfoService {
	return &GetStockInfoService{c: c, params: map[string]string{}}
}

// SetSymbol limits the reply to a single Reality trading pair (e.g. RAAPLUSDT);
// omitting it returns every pair.
func (s *GetStockInfoService) SetSymbol(symbol string) *GetStockInfoService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetStockInfoService) Do(ctx context.Context) ([]StockInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/stock-info", s.params)
	resp, err := request.Do[[]StockInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StockInfo describes the stock underlying a Reality trading pair.
type StockInfo struct {
	Symbol string `json:"symbol"`
	// Code is the ticker of the underlying stock (e.g. "AAPL").
	Code string `json:"code"`
	// Name is the company name; Bitget returns null for some listings.
	Name string `json:"name"`
	// TradingPeriod lists the tradable sessions: overnight, pre_market,
	// regular, after_hours.
	TradingPeriod   []string `json:"tradingPeriod"`
	WeekendTradable string   `json:"weekendTradable"` // no, yes
}

// GetMarketStatesService -- GET /api/v3/reality/market/states
//
// Returns the trading sessions of the stock market behind the Reality trading
// pairs, in that market's local time.
type GetMarketStatesService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetMarketStatesService() *GetMarketStatesService {
	return &GetMarketStatesService{c: c}
}

func (s *GetMarketStatesService) Do(ctx context.Context) (*MarketState, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/states")
	return request.Do[MarketState](req)
}

// MarketState is the session schedule of a stock market. The docs type data as
// a list, but the live API returns a single object.
type MarketState struct {
	Market       string             `json:"market"`       // US
	DaylightType string             `json:"daylightType"` // standard (winter), dst (summer)
	StateList    []MarketStateEntry `json:"stateList"`
}

// MarketStateEntry is one trading session window.
type MarketStateEntry struct {
	State     string `json:"state"`     // pre_market, regular, after_hours, overnight
	TimeZone  string `json:"timeZone"`  // EST, EDT
	StartTime string `json:"startTime"` // HH:mm
	EndTime   string `json:"endTime"`   // HH:mm
}

// GetMarketCalendarService -- GET /api/v3/reality/market/calendar
//
// Returns the closure calendar of the stock market behind the Reality trading
// pairs: recurring weekend closures plus one-off holiday closures.
type GetMarketCalendarService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetMarketCalendarService() *GetMarketCalendarService {
	return &GetMarketCalendarService{c: c}
}

func (s *GetMarketCalendarService) Do(ctx context.Context) (*MarketCalendar, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/calendar")
	return request.Do[MarketCalendar](req)
}

// MarketCalendar is the closure calendar of a stock market.
type MarketCalendar struct {
	TimeZone string `json:"timeZone"` // EST
	// SpecificConfig lists the non-regular closures, e.g. holidays.
	SpecificConfig []MarketClosure `json:"specificConfig"`
	// RegularConfig lists the weekdays that are always closed, e.g. SATURDAY.
	RegularConfig []string `json:"regularConfig"`
}

// MarketClosure is a single non-regular market closure window.
type MarketClosure struct {
	Remark    string `json:"remark"`
	StartTime string `json:"startTime"` // yyyy-MM-dd HH:mm
	EndTime   string `json:"endTime"`   // yyyy-MM-dd HH:mm
}
