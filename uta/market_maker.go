package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetScoreWeightsService -- GET /api/v3/market/score-weights
//
// Returns the market-maker score weights (required spread, minimum maker
// volume and weighting) per symbol, optionally filtered to a business line.
type GetScoreWeightsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetScoreWeightsService() *GetScoreWeightsService {
	return &GetScoreWeightsService{c: c, params: map[string]string{}}
}

// SetCategory filters to a business line; the endpoint accepts "SPOT" or
// "FUTURES" and returns every line when omitted.
func (s *GetScoreWeightsService) SetCategory(category string) *GetScoreWeightsService {
	s.params["category"] = category
	return s
}

func (s *GetScoreWeightsService) Do(ctx context.Context) ([]ScoreWeight, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/score-weights", s.params)
	resp, err := request.Do[[]ScoreWeight](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ScoreWeight is one (symbol, spread tier) score-weight row. Category is the
// raw business line as returned (e.g. "SPOT", "USDT_FUTURES").
type ScoreWeight struct {
	Category       string          `json:"category"`
	Label          string          `json:"label"`
	Symbol         string          `json:"symbol"`
	RequiredSpread decimal.Decimal `json:"requiredSpread"`
	MinMakerVolume decimal.Decimal `json:"minMakerVolume"`
	Weight         decimal.Decimal `json:"weight"`
}

// GetFeeGroupService -- GET /api/v3/market/fee-group
//
// Returns the market-maker fee groups for a business line: the symbol labels
// (with weights) that make up each group and the per-level maker fee tiers.
type GetFeeGroupService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFeeGroupService(category string) *GetFeeGroupService {
	return &GetFeeGroupService{c: c, params: map[string]string{"category": category}}
}

// SetGroup filters to a single group ("GROUP_A", "GROUP_B", ...); all groups
// are returned when omitted.
func (s *GetFeeGroupService) SetGroup(group string) *GetFeeGroupService {
	s.params["group"] = group
	return s
}

func (s *GetFeeGroupService) Do(ctx context.Context) ([]FeeGroup, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/fee-group", s.params)
	resp, err := request.Do[[]FeeGroup](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FeeGroup is one market-maker fee group for a business line.
type FeeGroup struct {
	Category  string          `json:"category"`
	Group     string          `json:"group"`
	LabelList []FeeGroupLabel `json:"labelList"`
	TierList  []FeeGroupTier  `json:"tierList"`
}

// FeeGroupLabel is a named bucket of symbols carrying a score weight.
type FeeGroupLabel struct {
	Weight  decimal.Decimal `json:"weight"`
	Label   string          `json:"label"`
	Symbols []string        `json:"symbols"`
}

// FeeGroupTier is one maker-fee tier within a group. Level is the tier label
// ("MM1" through "MM5"); a negative MakerFeeRate is a maker rebate.
type FeeGroupTier struct {
	Level        string          `json:"level"`
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
}

// GetCashDividendRecordsService -- GET /api/v3/market/cash-dividend-records
//
// Returns the cash dividend schedule for a tokenized-stock symbol, either
// upcoming ("pending") or settled ("paid").
type GetCashDividendRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCashDividendRecordsService(symbol, dividendType string) *GetCashDividendRecordsService {
	return &GetCashDividendRecordsService{c: c, params: map[string]string{
		"symbol": symbol,
		"type":   dividendType,
	}}
}

// SetCursor sets the pagination cursor (default 1).
func (s *GetCashDividendRecordsService) SetCursor(cursor string) *GetCashDividendRecordsService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit caps the page size (max 100, default 20).
func (s *GetCashDividendRecordsService) SetLimit(limit int) *GetCashDividendRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCashDividendRecordsService) Do(ctx context.Context) (*CashDividendRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/cash-dividend-records", s.params)
	return request.Do[CashDividendRecords](req)
}

// CashDividendRecords wraps the paginated dividend list payload.
type CashDividendRecords struct {
	List []CashDividendRecord `json:"list"`
}

// CashDividendRecord is one scheduled or settled cash dividend. ExDividendDate
// is a plain calendar date ("2026-08-20"), while CashDividendTimestamp is the
// settlement time in milliseconds.
type CashDividendRecord struct {
	ExDividendDate        string          `json:"exDividendDate"`
	CashDividendPerShare  decimal.Decimal `json:"cashDividendPerShare"`
	CashDividendTimestamp time.Time       `json:"cashDividendTimestamp"`
}
