package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetRiskReserveService -- GET /api/v3/market/risk-reserve
//
// Returns the daily risk-reserve (insurance) fund history for a futures symbol.
type GetRiskReserveService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRiskReserveService(category Category, symbol string) *GetRiskReserveService {
	return &GetRiskReserveService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetRiskReserveService) SetMarginCoin(marginCoin string) *GetRiskReserveService {
	s.params["marginCoin"] = marginCoin
	return s
}

func (s *GetRiskReserveService) Do(ctx context.Context) (*RiskReserve, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/risk-reserve", s.params)
	return request.Do[RiskReserve](req)
}

// RiskReserve is the daily reserve-fund history for a single settlement coin.
type RiskReserve struct {
	TotalBalance       decimal.Decimal     `json:"totalBalance"` // deprecated
	Coin               string              `json:"coin"`
	RiskReserveRecords []RiskReserveRecord `json:"riskReserveRecords"`
}

type RiskReserveRecord struct {
	Type    string          `json:"type"` // deprecated
	Amount  decimal.Decimal `json:"amount"`
	Balance decimal.Decimal `json:"balance"`
	Ts      time.Time       `json:"ts"`
}

// GetRiskReserveHourService -- GET /api/v3/market/risk-reserve-hour
//
// Returns the hourly risk-reserve (insurance) fund history for a futures symbol.
type GetRiskReserveHourService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRiskReserveHourService(category Category, symbol string) *GetRiskReserveHourService {
	return &GetRiskReserveHourService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetRiskReserveHourService) SetMarginCoin(marginCoin string) *GetRiskReserveHourService {
	s.params["marginCoin"] = marginCoin
	return s
}

func (s *GetRiskReserveHourService) Do(ctx context.Context) (*RiskReserveHour, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/risk-reserve-hour", s.params)
	return request.Do[RiskReserveHour](req)
}

// RiskReserveHour is the hourly reserve-fund history for a single settlement
// coin. It omits the deprecated totalBalance/type fields of the daily endpoint.
type RiskReserveHour struct {
	Coin               string                  `json:"coin"`
	RiskReserveRecords []RiskReserveHourRecord `json:"riskReserveRecords"`
}

type RiskReserveHourRecord struct {
	Amount  decimal.Decimal `json:"amount"`
	Balance decimal.Decimal `json:"balance"`
	Ts      time.Time       `json:"ts"`
}

// GetRiskReserveAllService -- GET /api/v3/market/risk-reserve-all
//
// Returns the current risk-reserve (insurance) fund balances for every fund in
// a product category, with the trading pairs each fund backs.
type GetRiskReserveAllService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRiskReserveAllService(category Category) *GetRiskReserveAllService {
	return &GetRiskReserveAllService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetRiskReserveAllService) Do(ctx context.Context) (*RiskReserveAll, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/risk-reserve-all", s.params)
	return request.Do[RiskReserveAll](req)
}

type RiskReserveAll struct {
	List []RiskReserveFund `json:"list"`
}

type RiskReserveFund struct {
	Symbols []string        `json:"symbols"`
	Coin    string          `json:"coin"`
	Balance decimal.Decimal `json:"balance"`
}

// GetDiscountRateService -- GET /api/v3/market/discount-rate
//
// Returns the tiered collateral discount rates applied to each coin's value
// when used as unified-account margin.
type GetDiscountRateService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetDiscountRateService() *GetDiscountRateService {
	return &GetDiscountRateService{c: c}
}

func (s *GetDiscountRateService) Do(ctx context.Context) ([]DiscountRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/discount-rate")
	resp, err := request.Do[[]DiscountRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type DiscountRate struct {
	Coin string             `json:"coin"`
	List []DiscountRateTier `json:"list"`
}

type DiscountRateTier struct {
	TierStartValue decimal.Decimal `json:"tierStartValue"`
	DiscountRate   decimal.Decimal `json:"discountRate"`
}

// GetMarginLoansService -- GET /api/v3/market/margin-loans
//
// Returns the current daily/annual borrow interest rates and loan limit for a
// coin in cross-margin / unified-account borrowing.
type GetMarginLoansService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMarginLoansService(coin string) *GetMarginLoansService {
	return &GetMarginLoansService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetMarginLoansService) Do(ctx context.Context) (*MarginLoan, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/margin-loans", s.params)
	return request.Do[MarginLoan](req)
}

type MarginLoan struct {
	DailyInterest        decimal.Decimal `json:"dailyInterest"`
	AnnualInterest       decimal.Decimal `json:"annualInterest"`
	Limit                decimal.Decimal `json:"limit"`
	MasterSubLimit       decimal.Decimal `json:"masterSubLimit"`       // master/sub account borrow limit
	PlatformRemaingQuota decimal.Decimal `json:"platformRemaingQuota"` // platform remaining quota (Bitget's spelling)
}

// GetPositionTierService -- GET /api/v3/market/position-tier
//
// Returns the leverage/maintenance-margin tier ladder for a futures symbol or a
// margin coin.
type GetPositionTierService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetPositionTierService(category Category) *GetPositionTierService {
	return &GetPositionTierService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetPositionTierService) SetSymbol(symbol string) *GetPositionTierService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetPositionTierService) SetCoin(coin string) *GetPositionTierService {
	s.params["coin"] = coin
	return s
}

func (s *GetPositionTierService) Do(ctx context.Context) ([]PositionTier, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/position-tier", s.params)
	resp, err := request.Do[[]PositionTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PositionTier struct {
	Tier         string          `json:"tier"`
	MinTierValue decimal.Decimal `json:"minTierValue"`
	MaxTierValue decimal.Decimal `json:"maxTierValue"`
	Leverage     string          `json:"leverage"`
	Mmr          decimal.Decimal `json:"mmr"`
}

// GetIndexComponentsService -- GET /api/v3/market/index-components
//
// Returns the constituent exchanges and weights that compose a symbol's index
// price.
type GetIndexComponentsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetIndexComponentsService(symbol string) *GetIndexComponentsService {
	return &GetIndexComponentsService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIndexComponentsService) SetCategory(category Category) *GetIndexComponentsService {
	s.params["category"] = string(category)
	return s
}

func (s *GetIndexComponentsService) Do(ctx context.Context) (*IndexComponents, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/index-components", s.params)
	return request.Do[IndexComponents](req)
}

type IndexComponents struct {
	Symbol        string           `json:"symbol"`
	ComponentList []IndexComponent `json:"componentList"`
}

type IndexComponent struct {
	Exchange        string          `json:"exchange"`
	SpotPair        string          `json:"spotPair"`
	EquivalentPrice decimal.Decimal `json:"equivalentPrice"`
	Weight          decimal.Decimal `json:"weight"`
}

// GetProofOfReservesService -- GET /api/v3/market/proof-of-reserves
//
// Returns Bitget's published proof-of-reserves: the Merkle root and the
// per-asset user vs. platform balances with reserve ratios.
type GetProofOfReservesService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetProofOfReservesService() *GetProofOfReservesService {
	return &GetProofOfReservesService{c: c}
}

func (s *GetProofOfReservesService) Do(ctx context.Context) (*ProofOfReserves, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/proof-of-reserves")
	return request.Do[ProofOfReserves](req)
}

// ProofOfReserves carries percentage ratios as raw strings (e.g. "127%"), which
// are not numeric, and the coin balances as decimals.
type ProofOfReserves struct {
	MerkleRootHash    string         `json:"merkleRootHash"`
	TotalReserveRatio string         `json:"totalReserveRatio"`
	List              []ReserveAsset `json:"list"`
}

type ReserveAsset struct {
	Coin           string          `json:"coin"`
	UserAssets     decimal.Decimal `json:"userAssets"`
	PlatformAssets decimal.Decimal `json:"platformAssets"`
	ReserveRatio   string          `json:"reserveRatio"`
}
