package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCollateralTypeService -- GET /api/v3/account/collateral-type (UTA mgt. read)
//
// Returns the account's current collateral configuration mode and, when the mode
// is custom, the comma-separated list of coins selected as collateral.
type GetCollateralTypeService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCollateralTypeService() *GetCollateralTypeService {
	return &GetCollateralTypeService{c: c}
}

func (s *GetCollateralTypeService) Do(ctx context.Context) (*CollateralType, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/collateral-type").WithSign()
	return request.Do[CollateralType](req)
}

type CollateralType struct {
	// CollateralType is the collateral mode: mainstream (USDT, USDC, BTC, DOGE,
	// ETH, SOL), all (every available coin), or custom (user-specified coins).
	CollateralType string `json:"collateralType"`
	// CollateralCoins is the comma-separated coin list, returned only when
	// CollateralType is custom.
	CollateralCoins string `json:"collateralCoins"`
}

// SetCollateralTypeService -- POST /api/v3/account/set-collateral-type (UTA mgt. read & write)
//
// Sets the account's collateral mode. The reply data is the literal string
// "success".
type SetCollateralTypeService struct {
	c    *UTAClient
	body map[string]any
}

// NewSetCollateralTypeService starts a request. collateralType is one of
// mainstream, all, or custom.
func (c *UTAClient) NewSetCollateralTypeService(collateralType string) *SetCollateralTypeService {
	return &SetCollateralTypeService{c: c, body: map[string]any{
		"collateralType": collateralType,
	}}
}

// SetCollateralCoins sets the comma-separated coin list (required when
// collateralType is custom).
func (s *SetCollateralTypeService) SetCollateralCoins(collateralCoins string) *SetCollateralTypeService {
	s.body["collateralCoins"] = collateralCoins
	return s
}

func (s *SetCollateralTypeService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/set-collateral-type", s.body).WithSign()
	return request.Do[string](req)
}

// GetCustomCollateralCoinsService -- GET /api/v3/account/custom-collateral-coins (UTA mgt. read)
//
// Returns the coins the platform accepts as custom collateral. Filter to a
// single coin with SetCollateralCoin, or omit it for the full list.
type GetCustomCollateralCoinsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCustomCollateralCoinsService() *GetCustomCollateralCoinsService {
	return &GetCustomCollateralCoinsService{c: c, params: map[string]string{}}
}

// SetCollateralCoin filters to a single collateral coin.
func (s *GetCustomCollateralCoinsService) SetCollateralCoin(collateralCoin string) *GetCustomCollateralCoinsService {
	s.params["collateralCoin"] = collateralCoin
	return s
}

func (s *GetCustomCollateralCoinsService) Do(ctx context.Context) ([]CustomCollateralCoin, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/custom-collateral-coins", s.params).WithSign()
	resp, err := request.Do[[]CustomCollateralCoin](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CustomCollateralCoin struct {
	CollateralCoin string `json:"collateralCoin"`
}

// PreSetLeverageService -- GET /api/v3/account/pre-set-leverage (UTA mgt. read)
//
// Previews the margin and tradable/borrowable impact of a leverage adjustment
// without applying it.
type PreSetLeverageService struct {
	c      *UTAClient
	params map[string]string
}

// NewPreSetLeverageService starts a preview. marginMode is isolated or crossed.
func (c *UTAClient) NewPreSetLeverageService(category Category, marginMode MarginMode) *PreSetLeverageService {
	return &PreSetLeverageService{c: c, params: map[string]string{
		"category":   string(category),
		"marginMode": string(marginMode),
	}}
}

// SetSymbol sets the symbol (required for futures leverage changes).
func (s *PreSetLeverageService) SetSymbol(symbol string) *PreSetLeverageService {
	s.params["symbol"] = symbol
	return s
}

// SetCoin sets the coin (required for margin trading leverage changes).
func (s *PreSetLeverageService) SetCoin(coin string) *PreSetLeverageService {
	s.params["coin"] = coin
	return s
}

// SetLeverage sets the leverage multiple for cross margin, one-way isolated, or
// uniform two-way isolated positions.
func (s *PreSetLeverageService) SetLeverage(leverage string) *PreSetLeverageService {
	s.params["leverage"] = leverage
	return s
}

// SetLongLeverage sets the long position leverage (isolated margin, two-way mode
// with differing leverage).
func (s *PreSetLeverageService) SetLongLeverage(longLeverage string) *PreSetLeverageService {
	s.params["longLeverage"] = longLeverage
	return s
}

// SetShortLeverage sets the short position leverage (isolated margin, two-way
// mode with differing leverage).
func (s *PreSetLeverageService) SetShortLeverage(shortLeverage string) *PreSetLeverageService {
	s.params["shortLeverage"] = shortLeverage
	return s
}

func (s *PreSetLeverageService) Do(ctx context.Context) (*PreSetLeverage, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/pre-set-leverage", s.params).WithSign()
	return request.Do[PreSetLeverage](req)
}

type PreSetLeverage struct {
	// EstMaxOpen is the projected maximum tradable quantity after the adjustment
	// (futures only).
	EstMaxOpen decimal.Decimal `json:"estMaxOpen"`
	// EstMaxBorrowable is the projected maximum borrowable quantity after the
	// adjustment, in the specified coin (margin only).
	EstMaxBorrowable decimal.Decimal `json:"estMaxBorrowable"`
	// RequiredMargin is the margin requirement denominated in USD.
	RequiredMargin decimal.Decimal `json:"requiredMargin"`
	// MarginChange is the change in margin usage (positive = increase).
	MarginChange decimal.Decimal `json:"marginChange"`
}
