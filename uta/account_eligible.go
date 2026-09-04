package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetEligibleSymbolsService -- GET /api/v3/account/eligible-symbols (UTA mgt. read)
//
// Returns the margin trading pairs the account is eligible (whitelisted) to
// trade. Non-whitelisted accounts get an empty list.
type GetEligibleSymbolsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEligibleSymbolsService() *GetEligibleSymbolsService {
	return &GetEligibleSymbolsService{c: c, params: map[string]string{}}
}

func (s *GetEligibleSymbolsService) SetSymbol(symbol string) *GetEligibleSymbolsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetEligibleSymbolsService) Do(ctx context.Context) ([]EligibleSymbol, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/eligible-symbols", s.params).WithSign()
	resp, err := request.Do[[]EligibleSymbol](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// EligibleSymbol is the trading specification of a whitelist-only margin pair.
type EligibleSymbol struct {
	Symbol     string   `json:"symbol"`
	Category   Category `json:"category"` // MARGIN
	BaseCoin   string   `json:"baseCoin"`
	QuoteCoin  string   `json:"quoteCoin"`
	SymbolType string   `json:"symbolType"` // crypto, metal, stock, commodity
	IsReality  string   `json:"isReality"`  // Reality stock token flag (no/yes)
	AreaSymbol string   `json:"areaSymbol"` // regionally restricted pair (no/yes)

	BuyLimitPriceRatio  decimal.Decimal  `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio decimal.Decimal  `json:"sellLimitPriceRatio"`
	MinOrderQty         decimal.Decimal  `json:"minOrderQty"`
	MaxOrderQty         decimal.Decimal  `json:"maxOrderQty"`
	MinOrderAmount      decimal.Decimal  `json:"minOrderAmount"`
	PricePrecision      decimal.Decimal  `json:"pricePrecision"`
	QuantityPrecision   decimal.Decimal  `json:"quantityPrecision"`
	QuotePrecision      decimal.Decimal  `json:"quotePrecision"`
	MaxSymbolOrderNum   string           `json:"maxSymbolOrderNum"`
	MaxProductOrderNum  string           `json:"maxProductOrderNum"`
	MaxPositionNum      string           `json:"maxPositionNum"`
	Status              InstrumentStatus `json:"status"`
	MaintainTime        string           `json:"maintainTime"`
	LaunchTime          time.Time        `json:"launchTime"`
}

// GetEligibleMarginTierService -- GET /api/v3/account/eligible-margin-tier (UTA mgt. read)
//
// Returns the leverage / maintenance-margin tier ladder that applies to the
// account under its whitelist eligibility, per margin coin.
type GetEligibleMarginTierService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEligibleMarginTierService() *GetEligibleMarginTierService {
	return &GetEligibleMarginTierService{c: c, params: map[string]string{}}
}

func (s *GetEligibleMarginTierService) SetCoin(coin string) *GetEligibleMarginTierService {
	s.params["coin"] = coin
	return s
}

func (s *GetEligibleMarginTierService) Do(ctx context.Context) ([]EligibleMarginTier, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/eligible-margin-tier", s.params).WithSign()
	resp, err := request.Do[[]EligibleMarginTier](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type EligibleMarginTier struct {
	Coin  string         `json:"coin"`
	Tiers []PositionTier `json:"tiers"`
}

// GetEligibleLoanInfoService -- GET /api/v3/account/eligible-loan-info (UTA mgt. read)
//
// Returns the margin borrow limits that apply to the account under its
// whitelist eligibility, per coin.
type GetEligibleLoanInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEligibleLoanInfoService() *GetEligibleLoanInfoService {
	return &GetEligibleLoanInfoService{c: c, params: map[string]string{}}
}

func (s *GetEligibleLoanInfoService) SetCoin(coin string) *GetEligibleLoanInfoService {
	s.params["coin"] = coin
	return s
}

func (s *GetEligibleLoanInfoService) Do(ctx context.Context) ([]EligibleLoanInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/eligible-loan-info", s.params).WithSign()
	resp, err := request.Do[[]EligibleLoanInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type EligibleLoanInfo struct {
	Limit                decimal.Decimal `json:"limit"`                // single UID borrow limit
	MasterSubLimit       decimal.Decimal `json:"masterSubLimit"`       // master/sub account borrow limit
	PlatformRemaingQuota decimal.Decimal `json:"platformRemaingQuota"` // platform remaining quota (Bitget's spelling)
}

// GetEligibleDiscountRateService -- GET /api/v3/account/eligible-discount-rate (UTA mgt. read)
//
// Returns the collateral discount-rate tiers the account is eligible for under
// its whitelist, per coin. Unlike GET /api/v3/market/discount-rate, which
// quotes the platform-wide ladder, this one is account-specific.
type GetEligibleDiscountRateService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEligibleDiscountRateService() *GetEligibleDiscountRateService {
	return &GetEligibleDiscountRateService{c: c, params: map[string]string{}}
}

func (s *GetEligibleDiscountRateService) SetCoin(coin string) *GetEligibleDiscountRateService {
	s.params["coin"] = coin
	return s
}

func (s *GetEligibleDiscountRateService) Do(ctx context.Context) ([]DiscountRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/eligible-discount-rate", s.params).WithSign()
	resp, err := request.Do[[]DiscountRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
