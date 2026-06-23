package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetVIPFeeRateService -- GET /api/v2/mix/market/vip-fee-rate (public)
//
// Returns the futures VIP fee-rate tiers (trading-volume / asset thresholds and
// the corresponding maker/taker rates and withdrawal limits).
type GetVIPFeeRateService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetVIPFeeRateService() *GetVIPFeeRateService {
	return &GetVIPFeeRateService{c: c, params: map[string]string{}}
}

// SetProductType filters the tiers to a single product line.
func (s *GetVIPFeeRateService) SetProductType(productType ProductType) *GetVIPFeeRateService {
	s.params["productType"] = string(productType)
	return s
}

func (s *GetVIPFeeRateService) Do(ctx context.Context) ([]VIPFeeRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/vip-fee-rate", s.params)
	resp, err := request.Do[[]VIPFeeRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// VIPFeeRate is one VIP fee-rate tier.
type VIPFeeRate struct {
	Level              string          `json:"level"`
	DealAmount         decimal.Decimal `json:"dealAmount"`
	AssetAmount        decimal.Decimal `json:"assetAmount"`
	TakerFeeRate       decimal.Decimal `json:"takerFeeRate"`
	MakerFeeRate       decimal.Decimal `json:"makerFeeRate"`
	BtcWithdrawAmount  decimal.Decimal `json:"btcWithdrawAmount"`
	UsdtWithdrawAmount decimal.Decimal `json:"usdtWithdrawAmount"`
}

// GetInterestRateHistoryService -- GET /api/v2/mix/market/union-interest-rate-history (public)
//
// Returns the historical unified-margin interest rates for a coin.
type GetInterestRateHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetInterestRateHistoryService(coin string) *GetInterestRateHistoryService {
	return &GetInterestRateHistoryService{c: c, params: map[string]string{"coin": coin}}
}

// SetPageSize caps the number of history rows returned.
func (s *GetInterestRateHistoryService) SetPageSize(pageSize int) *GetInterestRateHistoryService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo selects the page of history rows to return.
func (s *GetInterestRateHistoryService) SetPageNo(pageNo int) *GetInterestRateHistoryService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetInterestRateHistoryService) Do(ctx context.Context) (*InterestRateHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/union-interest-rate-history", s.params)
	return request.Do[InterestRateHistory](req)
}

// InterestRateHistory is the interest-rate history for one coin.
type InterestRateHistory struct {
	Coin                    string                `json:"coin"`
	HistoryInterestRateList []HistoryInterestRate `json:"historyInterestRateList"`
}

// HistoryInterestRate is one interest-rate history snapshot.
type HistoryInterestRate struct {
	Ts                 time.Time       `json:"ts"`
	AnnualInterestRate decimal.Decimal `json:"annualInterestRate"`
	DailyInterestRate  decimal.Decimal `json:"dailyInterestRate"`
}

// GetExchangeRateService -- GET /api/v2/mix/market/exchange-rate (public)
//
// Returns the per-coin interest exchange-rate tiers used when a borrowed asset
// is converted for margin valuation.
type GetExchangeRateService struct {
	c *MixClient
}

func (c *MixClient) NewGetExchangeRateService() *GetExchangeRateService {
	return &GetExchangeRateService{c: c}
}

func (s *GetExchangeRateService) Do(ctx context.Context) ([]ExchangeRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/exchange-rate")
	resp, err := request.Do[[]ExchangeRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ExchangeRate is the exchange-rate tier list for a single coin. The nested-list
// json key is misspelled "excahngeRateList" in the live API; the tag matches the
// wire exactly.
type ExchangeRate struct {
	Coin             string             `json:"coin"`
	ExchangeRateList []ExchangeRateTier `json:"excahngeRateList"`
}

// ExchangeRateTier is one exchange-rate tier.
type ExchangeRateTier struct {
	Tier         string          `json:"tier"`
	MinAmount    decimal.Decimal `json:"minAmount"`
	MaxAmount    decimal.Decimal `json:"maxAmount"`
	ExchangeRate decimal.Decimal `json:"exchangeRate"`
}

// GetDiscountRateService -- GET /api/v2/mix/market/discount-rate (public)
//
// Returns the per-coin collateral discount-rate tiers and the per-user / total
// collateral limits.
type GetDiscountRateService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetDiscountRateService() *GetDiscountRateService {
	return &GetDiscountRateService{c: c, params: map[string]string{}}
}

// SetProductType filters the tiers to a single product line.
func (s *GetDiscountRateService) SetProductType(productType ProductType) *GetDiscountRateService {
	s.params["productType"] = string(productType)
	return s
}

func (s *GetDiscountRateService) Do(ctx context.Context) ([]DiscountRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/discount-rate", s.params)
	resp, err := request.Do[[]DiscountRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DiscountRate is the collateral discount-rate tier list for a single coin.
type DiscountRate struct {
	Coin             string             `json:"coin"`
	UserLimit        decimal.Decimal    `json:"userLimit"`
	TotalLimit       decimal.Decimal    `json:"totalLimit"`
	DiscountRateList []DiscountRateTier `json:"discountRateList"`
}

// DiscountRateTier is one collateral discount-rate tier.
type DiscountRateTier struct {
	Tier         string          `json:"tier"`
	MinAmount    decimal.Decimal `json:"minAmount"`
	MaxAmount    decimal.Decimal `json:"maxAmount"`
	DiscountRate decimal.Decimal `json:"discountRate"`
}

// GetHistoryFundRateService -- GET /api/v2/mix/market/history-fund-rate (public)
//
// Returns the historical funding rates for a symbol.
type GetHistoryFundRateService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetHistoryFundRateService(symbol string, productType ProductType) *GetHistoryFundRateService {
	return &GetHistoryFundRateService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetPageSize caps the number of funding-rate rows returned.
func (s *GetHistoryFundRateService) SetPageSize(pageSize int) *GetHistoryFundRateService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo selects the page of funding-rate rows to return.
func (s *GetHistoryFundRateService) SetPageNo(pageNo int) *GetHistoryFundRateService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetHistoryFundRateService) Do(ctx context.Context) ([]HistoryFundRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/history-fund-rate", s.params)
	resp, err := request.Do[[]HistoryFundRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// HistoryFundRate is one historical funding-rate settlement.
type HistoryFundRate struct {
	Symbol      string          `json:"symbol"`
	FundingRate decimal.Decimal `json:"fundingRate"`
	FundingTime time.Time       `json:"fundingTime"`
}

// GetCurrentFundRateService -- GET /api/v2/mix/market/current-fund-rate (public)
//
// Returns the current funding rate for a symbol, including the funding interval
// and the funding-rate caps.
type GetCurrentFundRateService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetCurrentFundRateService(symbol string, productType ProductType) *GetCurrentFundRateService {
	return &GetCurrentFundRateService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetCurrentFundRateService) Do(ctx context.Context) ([]CurrentFundRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/current-fund-rate", s.params)
	resp, err := request.Do[[]CurrentFundRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CurrentFundRate is the current funding rate for a single symbol.
type CurrentFundRate struct {
	Symbol              string          `json:"symbol"`
	FundingRate         decimal.Decimal `json:"fundingRate"`
	FundingRateInterval string          `json:"fundingRateInterval"` // hours between settlements
	NextUpdate          time.Time       `json:"nextUpdate"`
	MinFundingRate      decimal.Decimal `json:"minFundingRate"`
	MaxFundingRate      decimal.Decimal `json:"maxFundingRate"`
}
