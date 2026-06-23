package common

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// TradeRateBusinessType is the business line (product) selector shared by the
// trade-rate endpoints.
type TradeRateBusinessType string

const (
	TradeRateBusinessTypeSpot   TradeRateBusinessType = "spot"
	TradeRateBusinessTypeMix    TradeRateBusinessType = "mix"
	TradeRateBusinessTypeMargin TradeRateBusinessType = "margin"
)

// GetServerTimeService -- GET /api/v2/public/time (public)
//
// Returns the Bitget server time. Bitget exposes a single shared time endpoint
// (under the v2 common namespace); SyncServerTime uses it to correct clock
// drift before signing private requests.
type GetServerTimeService struct {
	c *CommonClient
}

func (c *CommonClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/public/time")
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime time.Time `json:"serverTime"`
}

// GetTradeRateService -- GET /api/v2/common/trade-rate (private)
//
// Returns the account's maker/taker fee rates for one symbol on a given
// business line (spot, mix, or margin).
type GetTradeRateService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetTradeRateService(symbol string, businessType TradeRateBusinessType) *GetTradeRateService {
	return &GetTradeRateService{c: c, params: map[string]string{
		"symbol":       symbol,
		"businessType": string(businessType),
	}}
}

func (s *GetTradeRateService) Do(ctx context.Context) (*TradeRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/common/trade-rate", s.params).WithSign()
	return request.Do[TradeRate](req)
}

type TradeRate struct {
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate decimal.Decimal `json:"takerFeeRate"`
}

// GetAllTradeRateService -- GET /api/v2/common/all-trade-rate (private)
//
// Returns the account's maker/taker fee rates across symbols for a given
// business line (spot, mix, or margin).
type GetAllTradeRateService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetAllTradeRateService(symbol string, businessType TradeRateBusinessType) *GetAllTradeRateService {
	return &GetAllTradeRateService{c: c, params: map[string]string{
		"symbol":       symbol,
		"businessType": string(businessType),
	}}
}

func (s *GetAllTradeRateService) Do(ctx context.Context) ([]SymbolTradeRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/common/all-trade-rate", s.params).WithSign()
	resp, err := request.Do[[]SymbolTradeRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SymbolTradeRate struct {
	Symbol       string          `json:"symbol"`
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate decimal.Decimal `json:"takerFeeRate"`
}
