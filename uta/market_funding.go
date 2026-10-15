package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCurrentFundingRateService -- GET /api/v3/market/current-fund-rate
//
// Returns the current funding rate for a futures category, optionally filtered
// to a single symbol.
type GetCurrentFundingRateService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCurrentFundingRateService(category Category) *GetCurrentFundingRateService {
	return &GetCurrentFundingRateService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetCurrentFundingRateService) SetSymbol(symbol string) *GetCurrentFundingRateService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetCurrentFundingRateService) Do(ctx context.Context) ([]CurrentFundingRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/current-fund-rate", s.params)
	resp, err := request.Do[[]CurrentFundingRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CurrentFundingRate struct {
	Symbol                 string          `json:"symbol"`
	FundingRate            decimal.Decimal `json:"fundingRate"`
	FundingRateInterval    string          `json:"fundingRateInterval"`
	NextUpdate             time.Time       `json:"nextUpdate,format:unixmilli"`
	MinFundingRate         decimal.Decimal `json:"minFundingRate"`
	MaxFundingRate         decimal.Decimal `json:"maxFundingRate"`
	CashDividend           decimal.Decimal `json:"cashDividend"`
	CashDividendNextUpdate time.Time       `json:"cashDividendNextUpdate,format:unixmilli"`
}

// GetFundingRateHistoryService -- GET /api/v3/market/history-fund-rate
//
// Returns the historical funding rate records for a futures symbol, most recent
// first, with cursor-based pagination.
type GetFundingRateHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFundingRateHistoryService(category Category, symbol string) *GetFundingRateHistoryService {
	return &GetFundingRateHistoryService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetFundingRateHistoryService) SetCursor(cursor int) *GetFundingRateHistoryService {
	s.params["cursor"] = strconv.Itoa(cursor)
	return s
}

func (s *GetFundingRateHistoryService) SetLimit(limit int) *GetFundingRateHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFundingRateHistoryService) Do(ctx context.Context) (*FundingRateHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/history-fund-rate", s.params)
	return request.Do[FundingRateHistoryResponse](req)
}

type FundingRateHistoryResponse struct {
	ResultList []FundingRateHistory `json:"resultList"`
}

type FundingRateHistory struct {
	Symbol               string          `json:"symbol"`
	FundingRate          decimal.Decimal `json:"fundingRate"`
	FundingRateTimestamp time.Time       `json:"fundingRateTimestamp,format:unixmilli"`
}

// GetOpenInterestService -- GET /api/v3/market/open-interest
//
// Returns the open interest for a futures category, optionally filtered to a
// single symbol.
type GetOpenInterestService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOpenInterestService(category Category) *GetOpenInterestService {
	return &GetOpenInterestService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetOpenInterestService) SetSymbol(symbol string) *GetOpenInterestService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetOpenInterestService) Do(ctx context.Context) (*OpenInterestResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/open-interest", s.params)
	return request.Do[OpenInterestResponse](req)
}

type OpenInterestResponse struct {
	List []OpenInterest `json:"list"`
	Ts   time.Time      `json:"ts,format:unixmilli"`
}

type OpenInterest struct {
	Symbol       string          `json:"symbol"`
	OpenInterest decimal.Decimal `json:"openInterest"`
}

// GetOpenInterestHistoryService -- GET /api/v3/market/open-interest-history
//
// Returns the recent open-interest series of a futures symbol, newest first.
// period is required: 5m, 15m, 30m, 1h, 2h, 4h, 6h, 12h or 1d.
type GetOpenInterestHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOpenInterestHistoryService(symbol, period string) *GetOpenInterestHistoryService {
	return &GetOpenInterestHistoryService{c: c, params: map[string]string{"symbol": symbol, "period": period}}
}

func (s *GetOpenInterestHistoryService) Do(ctx context.Context) (*OpenInterestHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/open-interest-history", s.params)
	return request.Do[OpenInterestHistory](req)
}

type OpenInterestHistory struct {
	List []OpenInterestRecord `json:"list"`
}

// OpenInterestRecord is the open interest at the end of one period. Quantities
// are in base coin; values are in quote coin (openInterest * markPrice). The
// one-side figures count a single side of the book.
type OpenInterestRecord struct {
	Symbol                   string          `json:"symbol"`
	OpenInterest             decimal.Decimal `json:"openInterest"`
	OneSideOpenInterest      decimal.Decimal `json:"oneSideOpenInterest"`
	OpenInterestValue        decimal.Decimal `json:"openInterestValue"`
	OneSideOpenInterestValue decimal.Decimal `json:"oneSideOpenInterestValue"`
	Timestamp                time.Time       `json:"timestamp,format:unixmilli"` // period end
}

// GetOpenInterestLimitService -- GET /api/v3/market/oi-limit
//
// Returns the per-symbol and total notional open-interest limits for a futures
// category, optionally filtered to a single symbol.
type GetOpenInterestLimitService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOpenInterestLimitService(category Category) *GetOpenInterestLimitService {
	return &GetOpenInterestLimitService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetOpenInterestLimitService) SetSymbol(symbol string) *GetOpenInterestLimitService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetOpenInterestLimitService) Do(ctx context.Context) ([]OpenInterestLimit, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/oi-limit", s.params)
	resp, err := request.Do[[]OpenInterestLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type OpenInterestLimit struct {
	Symbol             string          `json:"symbol"`
	NotionalValue      decimal.Decimal `json:"notionalValue"`
	TotalNotionalValue decimal.Decimal `json:"totalNotionalValue"`
}
