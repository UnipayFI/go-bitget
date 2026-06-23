package common

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// InsightPeriod is the candle/bucket interval shared by the futures (mix)
// big-data insight feeds. The default when omitted is 5m.
type InsightPeriod string

const (
	InsightPeriod5m  InsightPeriod = "5m"
	InsightPeriod15m InsightPeriod = "15m"
	InsightPeriod30m InsightPeriod = "30m"
	InsightPeriod1h  InsightPeriod = "1h"
	InsightPeriod2h  InsightPeriod = "2h"
	InsightPeriod4h  InsightPeriod = "4h"
	InsightPeriod6h  InsightPeriod = "6h"
	InsightPeriod12h InsightPeriod = "12h"
	InsightPeriod1d  InsightPeriod = "1d"
)

// LongShortPeriod is the bucket interval for the futures long/short ratio feed,
// which differs from the other mix feeds by offering a UTC-aligned daily bucket
// (1Dutc) instead of 1d. The default when omitted is 5m.
type LongShortPeriod string

const (
	LongShortPeriod5m    LongShortPeriod = "5m"
	LongShortPeriod15m   LongShortPeriod = "15m"
	LongShortPeriod30m   LongShortPeriod = "30m"
	LongShortPeriod1h    LongShortPeriod = "1h"
	LongShortPeriod2h    LongShortPeriod = "2h"
	LongShortPeriod4h    LongShortPeriod = "4h"
	LongShortPeriod6h    LongShortPeriod = "6h"
	LongShortPeriod12h   LongShortPeriod = "12h"
	LongShortPeriod1Dutc LongShortPeriod = "1Dutc"
)

// MarginInsightPeriod is the period for the margin big-data feeds (long/short
// ratio, loan growth, isolated borrow rate). The default when omitted is 24h.
type MarginInsightPeriod string

const (
	MarginInsightPeriod24h MarginInsightPeriod = "24h"
	MarginInsightPeriod30d MarginInsightPeriod = "30d"
)

// FundFlowPeriod is the period for the spot fund-flow feed. The default when
// omitted is 15m.
type FundFlowPeriod string

const (
	FundFlowPeriod15m FundFlowPeriod = "15m"
	FundFlowPeriod30m FundFlowPeriod = "30m"
	FundFlowPeriod1h  FundFlowPeriod = "1h"
	FundFlowPeriod2h  FundFlowPeriod = "2h"
	FundFlowPeriod4h  FundFlowPeriod = "4h"
	FundFlowPeriod1d  FundFlowPeriod = "1d"
)

// GetWhaleNetFlowService -- GET /api/v2/spot/market/whale-net-flow (public)
//
// Returns the daily net buy/sell volume of whale (large) spot transactions for
// a symbol.
type GetWhaleNetFlowService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetWhaleNetFlowService(symbol string) *GetWhaleNetFlowService {
	return &GetWhaleNetFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetWhaleNetFlowService) Do(ctx context.Context) ([]WhaleNetFlow, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/whale-net-flow", s.params)
	resp, err := request.Do[[]WhaleNetFlow](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// WhaleNetFlow is one whale net-flow record. The timestamp arrives under the
// "date" key (millisecond string), not the usual "ts".
type WhaleNetFlow struct {
	Volume decimal.Decimal `json:"volume"`
	Date   time.Time       `json:"date"`
}

// GetTakerBuySellService -- GET /api/v2/mix/market/taker-buy-sell (public)
//
// Returns the futures taker (active) buy and sell volume per interval for a
// symbol.
type GetTakerBuySellService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetTakerBuySellService(symbol string) *GetTakerBuySellService {
	return &GetTakerBuySellService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 5m).
func (s *GetTakerBuySellService) SetPeriod(period InsightPeriod) *GetTakerBuySellService {
	s.params["period"] = string(period)
	return s
}

func (s *GetTakerBuySellService) Do(ctx context.Context) ([]TakerBuySell, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/taker-buy-sell", s.params)
	resp, err := request.Do[[]TakerBuySell](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TakerBuySell is one futures taker buy/sell volume bucket.
type TakerBuySell struct {
	BuyVolume  decimal.Decimal `json:"buyVolume"`
	SellVolume decimal.Decimal `json:"sellVolume"`
	Ts         time.Time       `json:"ts"`
}

// GetPositionLongShortService -- GET /api/v2/mix/market/position-long-short (public)
//
// Returns the proportion of futures long vs short positions per interval for a
// symbol.
type GetPositionLongShortService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetPositionLongShortService(symbol string) *GetPositionLongShortService {
	return &GetPositionLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 5m).
func (s *GetPositionLongShortService) SetPeriod(period InsightPeriod) *GetPositionLongShortService {
	s.params["period"] = string(period)
	return s
}

func (s *GetPositionLongShortService) Do(ctx context.Context) ([]PositionLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/position-long-short", s.params)
	resp, err := request.Do[[]PositionLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PositionLongShort is one futures long/short position-ratio bucket.
type PositionLongShort struct {
	LongPositionRatio      decimal.Decimal `json:"longPositionRatio"`
	ShortPositionRatio     decimal.Decimal `json:"shortPositionRatio"`
	LongShortPositionRatio decimal.Decimal `json:"longShortPositionRatio"`
	Ts                     time.Time       `json:"ts"`
}

// GetMarginLongShortRatioService -- GET /api/v2/margin/market/long-short-ratio (public)
//
// Returns the leveraged (margin) long-short ratio per interval for a symbol.
type GetMarginLongShortRatioService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetMarginLongShortRatioService(symbol string) *GetMarginLongShortRatioService {
	return &GetMarginLongShortRatioService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 24h).
func (s *GetMarginLongShortRatioService) SetPeriod(period MarginInsightPeriod) *GetMarginLongShortRatioService {
	s.params["period"] = string(period)
	return s
}

// SetCoin filters by base or quote coin (defaults to the base coin).
func (s *GetMarginLongShortRatioService) SetCoin(coin string) *GetMarginLongShortRatioService {
	s.params["coin"] = coin
	return s
}

func (s *GetMarginLongShortRatioService) Do(ctx context.Context) ([]MarginLongShortRatio, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/market/long-short-ratio", s.params)
	resp, err := request.Do[[]MarginLongShortRatio](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginLongShortRatio is one margin long-short ratio bucket.
type MarginLongShortRatio struct {
	LongShortRatio decimal.Decimal `json:"longShortRatio"`
	Ts             time.Time       `json:"ts"`
}

// GetMarginLoanGrowthService -- GET /api/v2/margin/market/loan-growth (public)
//
// Returns the margin loan-growth rate per interval for a symbol.
type GetMarginLoanGrowthService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetMarginLoanGrowthService(symbol string) *GetMarginLoanGrowthService {
	return &GetMarginLoanGrowthService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 24h).
func (s *GetMarginLoanGrowthService) SetPeriod(period MarginInsightPeriod) *GetMarginLoanGrowthService {
	s.params["period"] = string(period)
	return s
}

// SetCoin filters by base or quote coin (defaults to the base coin).
func (s *GetMarginLoanGrowthService) SetCoin(coin string) *GetMarginLoanGrowthService {
	s.params["coin"] = coin
	return s
}

func (s *GetMarginLoanGrowthService) Do(ctx context.Context) ([]MarginLoanGrowth, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/market/loan-growth", s.params)
	resp, err := request.Do[[]MarginLoanGrowth](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginLoanGrowth is one margin loan-growth-rate bucket.
type MarginLoanGrowth struct {
	GrowthRate decimal.Decimal `json:"growthRate"`
	Ts         time.Time       `json:"ts"`
}

// GetIsolatedBorrowRateService -- GET /api/v2/margin/market/isolated-borrow-rate (public)
//
// Returns the isolated-margin borrowing ratio per interval for a symbol.
type GetIsolatedBorrowRateService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetIsolatedBorrowRateService(symbol string) *GetIsolatedBorrowRateService {
	return &GetIsolatedBorrowRateService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 24h).
func (s *GetIsolatedBorrowRateService) SetPeriod(period MarginInsightPeriod) *GetIsolatedBorrowRateService {
	s.params["period"] = string(period)
	return s
}

func (s *GetIsolatedBorrowRateService) Do(ctx context.Context) ([]IsolatedBorrowRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/market/isolated-borrow-rate", s.params)
	resp, err := request.Do[[]IsolatedBorrowRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedBorrowRate is one isolated-margin borrow-ratio bucket.
type IsolatedBorrowRate struct {
	BorrowRate decimal.Decimal `json:"borrowRate"`
	Ts         time.Time       `json:"ts"`
}

// GetLongShortService -- GET /api/v2/mix/market/long-short (public)
//
// Returns the futures long and short ratio per interval for a symbol.
type GetLongShortService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetLongShortService(symbol string) *GetLongShortService {
	return &GetLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 5m).
func (s *GetLongShortService) SetPeriod(period LongShortPeriod) *GetLongShortService {
	s.params["period"] = string(period)
	return s
}

func (s *GetLongShortService) Do(ctx context.Context) ([]LongShort, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/long-short", s.params)
	resp, err := request.Do[[]LongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LongShort is one futures long/short ratio bucket.
type LongShort struct {
	LongRatio      decimal.Decimal `json:"longRatio"`
	ShortRatio     decimal.Decimal `json:"shortRatio"`
	LongShortRatio decimal.Decimal `json:"longShortRatio"`
	Ts             time.Time       `json:"ts"`
}

// GetSpotFundFlowService -- GET /api/v2/spot/market/fund-flow (public)
//
// Returns the spot fund flow for a symbol, broken down by trader tier (whale /
// dolphin / fish) and side (buy / sell) over the requested period.
type GetSpotFundFlowService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetSpotFundFlowService(symbol string) *GetSpotFundFlowService {
	return &GetSpotFundFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the query period (default 15m).
func (s *GetSpotFundFlowService) SetPeriod(period FundFlowPeriod) *GetSpotFundFlowService {
	s.params["period"] = string(period)
	return s
}

func (s *GetSpotFundFlowService) Do(ctx context.Context) (*SpotFundFlow, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/fund-flow", s.params)
	return request.Do[SpotFundFlow](req)
}

// SpotFundFlow is the spot fund-flow breakdown by trader tier and side.
type SpotFundFlow struct {
	WhaleBuyVolume    decimal.Decimal `json:"whaleBuyVolume"`
	DolphinBuyVolume  decimal.Decimal `json:"dolphinBuyVolume"`
	FishBuyVolume     decimal.Decimal `json:"fishBuyVolume"`
	WhaleSellVolume   decimal.Decimal `json:"whaleSellVolume"`
	DolphinSellVolume decimal.Decimal `json:"dolphinSellVolume"`
	FishSellVolume    decimal.Decimal `json:"fishSellVolume"`
	WhaleBuyRatio     decimal.Decimal `json:"whaleBuyRatio"`
	DolphinBuyRatio   decimal.Decimal `json:"dolphinBuyRatio"`
	FishBuyRatio      decimal.Decimal `json:"fishBuyRatio"`
	WhaleSellRatio    decimal.Decimal `json:"whaleSellRatio"`
	DolphinSellRatio  decimal.Decimal `json:"dolphinSellRatio"`
	FishSellRatio     decimal.Decimal `json:"fishSellRatio"`
}

// GetSupportSymbolsService -- GET /api/v2/spot/market/support-symbols (public)
//
// Returns the spot and futures symbols for which the big-data trading-insight
// feeds are available.
type GetSupportSymbolsService struct {
	c *CommonClient
}

func (c *CommonClient) NewGetSupportSymbolsService() *GetSupportSymbolsService {
	return &GetSupportSymbolsService{c: c}
}

func (s *GetSupportSymbolsService) Do(ctx context.Context) (*SupportSymbols, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/support-symbols")
	return request.Do[SupportSymbols](req)
}

// SupportSymbols lists the symbols supported by the big-data feeds.
type SupportSymbols struct {
	SpotList   []string `json:"spotList"`
	FutureList []string `json:"futureList"`
}

// GetFundNetFlowService -- GET /api/v2/spot/market/fund-net-flow (public)
//
// Returns the spot 24h net capital inflow series for a symbol.
type GetFundNetFlowService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetFundNetFlowService(symbol string) *GetFundNetFlowService {
	return &GetFundNetFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetFundNetFlowService) Do(ctx context.Context) ([]FundNetFlow, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/fund-net-flow", s.params)
	resp, err := request.Do[[]FundNetFlow](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FundNetFlow is one spot net-capital-inflow bucket.
type FundNetFlow struct {
	NetFlow decimal.Decimal `json:"netFlow"`
	Ts      time.Time       `json:"ts"`
}

// GetAccountLongShortService -- GET /api/v2/mix/market/account-long-short (public)
//
// Returns the proportion of futures accounts holding long vs short positions
// per interval for a symbol.
type GetAccountLongShortService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetAccountLongShortService(symbol string) *GetAccountLongShortService {
	return &GetAccountLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the bucket interval (default 5m).
func (s *GetAccountLongShortService) SetPeriod(period InsightPeriod) *GetAccountLongShortService {
	s.params["period"] = string(period)
	return s
}

func (s *GetAccountLongShortService) Do(ctx context.Context) ([]AccountLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/account-long-short", s.params)
	resp, err := request.Do[[]AccountLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountLongShort is one futures long/short account-ratio bucket.
type AccountLongShort struct {
	LongAccountRatio      decimal.Decimal `json:"longAccountRatio"`
	ShortAccountRatio     decimal.Decimal `json:"shortAccountRatio"`
	LongShortAccountRatio decimal.Decimal `json:"longShortAccountRatio"`
	Ts                    time.Time       `json:"ts"`
}
