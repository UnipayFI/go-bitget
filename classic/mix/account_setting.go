package mix

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SetAutoMarginService -- POST /api/v2/mix/account/set-auto-margin (signed)
//
// Enables or disables automatic margin top-up for an isolated-margin position.
// The reply data is the literal string "success".
type SetAutoMarginService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetAutoMarginService(symbol string, autoMargin AutoMargin, marginCoin string, holdSide HoldSide) *SetAutoMarginService {
	return &SetAutoMarginService{c: c, body: map[string]any{
		"symbol":     symbol,
		"autoMargin": string(autoMargin),
		"marginCoin": marginCoin,
		"holdSide":   string(holdSide),
	}}
}

func (s *SetAutoMarginService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-auto-margin", s.body).WithSign()
	return request.Do[string](req)
}

// SetLeverageService -- POST /api/v2/mix/account/set-leverage (signed)
//
// Adjusts the leverage for a single symbol. In cross margin (or one-way
// isolated) pass leverage; in hedge-mode isolated pass longLeverage/shortLeverage
// together with holdSide.
type SetLeverageService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetLeverageService(symbol string, productType ProductType, marginCoin string) *SetLeverageService {
	return &SetLeverageService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
	}}
}

// SetLeverage sets the leverage ratio for cross margin and one-way isolated
// positions.
func (s *SetLeverageService) SetLeverage(leverage string) *SetLeverageService {
	s.body["leverage"] = leverage
	return s
}

// SetLongLeverage sets the long-position leverage (hedge-mode isolated margin).
func (s *SetLeverageService) SetLongLeverage(longLeverage string) *SetLeverageService {
	s.body["longLeverage"] = longLeverage
	return s
}

// SetShortLeverage sets the short-position leverage (hedge-mode isolated margin).
func (s *SetLeverageService) SetShortLeverage(shortLeverage string) *SetLeverageService {
	s.body["shortLeverage"] = shortLeverage
	return s
}

// SetHoldSide sets the position direction (required for hedge-mode isolated
// margin; omit for cross margin).
func (s *SetLeverageService) SetHoldSide(holdSide HoldSide) *SetLeverageService {
	s.body["holdSide"] = string(holdSide)
	return s
}

func (s *SetLeverageService) Do(ctx context.Context) (*Leverage, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-leverage", s.body).WithSign()
	return request.Do[Leverage](req)
}

// Leverage is the per-symbol leverage configuration returned by Set-Leverage.
type Leverage struct {
	Symbol              string     `json:"symbol"`
	MarginCoin          string     `json:"marginCoin"`
	LongLeverage        string     `json:"longLeverage"`
	ShortLeverage       string     `json:"shortLeveage"` // NB: Bitget misspells the JSON key as "shortLeveage"
	CrossMarginLeverage string     `json:"crossMarginLeverage"`
	MarginMode          MarginMode `json:"marginMode"`
}

// SetAllLeverageService -- POST /api/v2/mix/account/set-all-leverage (signed)
//
// Sets the same leverage across an entire product line (applies only to symbols
// that currently have open positions). The reply data is the literal string
// "success".
type SetAllLeverageService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetAllLeverageService(productType ProductType, leverage string) *SetAllLeverageService {
	return &SetAllLeverageService{c: c, body: map[string]any{
		"productType": string(productType),
		"leverage":    leverage,
	}}
}

// SetSymbol restricts the change to a single trading pair.
func (s *SetAllLeverageService) SetSymbol(symbol string) *SetAllLeverageService {
	s.body["symbol"] = symbol
	return s
}

// SetMarginCoin sets the margin coin to apply the leverage to.
func (s *SetAllLeverageService) SetMarginCoin(marginCoin string) *SetAllLeverageService {
	s.body["marginCoin"] = marginCoin
	return s
}

// SetHoldSide sets the position direction.
func (s *SetAllLeverageService) SetHoldSide(holdSide HoldSide) *SetAllLeverageService {
	s.body["holdSide"] = string(holdSide)
	return s
}

func (s *SetAllLeverageService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-all-leverage", s.body).WithSign()
	return request.Do[string](req)
}

// SetMarginService -- POST /api/v2/mix/account/set-margin (signed)
//
// Adds or removes margin for an isolated-margin position. A positive amount adds
// margin, a negative amount removes it. The reply data is an empty string on
// success.
type SetMarginService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetMarginService(symbol string, productType ProductType, marginCoin string, holdSide HoldSide, amount decimal.Decimal) *SetMarginService {
	return &SetMarginService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
		"holdSide":    string(holdSide),
		"amount":      amount.String(),
	}}
}

func (s *SetMarginService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-margin", s.body).WithSign()
	return request.Do[string](req)
}

// SetAssetModeService -- POST /api/v2/mix/account/set-asset-mode (signed)
//
// Switches the USDT-M futures asset mode between single-asset and multi-asset
// (union) margining. The reply data is the literal string "success".
type SetAssetModeService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetAssetModeService(productType ProductType, assetMode AssetMode) *SetAssetModeService {
	return &SetAssetModeService{c: c, body: map[string]any{
		"productType": string(productType),
		"assetMode":   string(assetMode),
	}}
}

func (s *SetAssetModeService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-asset-mode", s.body).WithSign()
	return request.Do[string](req)
}

// SetMarginModeService -- POST /api/v2/mix/account/set-margin-mode (signed)
//
// Switches a symbol between isolated and cross margin. Cannot be used while any
// position or order is open for the symbol.
type SetMarginModeService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetMarginModeService(symbol string, productType ProductType, marginCoin string, marginMode MarginMode) *SetMarginModeService {
	return &SetMarginModeService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
		"marginMode":  string(marginMode),
	}}
}

func (s *SetMarginModeService) Do(ctx context.Context) (*MarginModeResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-margin-mode", s.body).WithSign()
	return request.Do[MarginModeResult](req)
}

// MarginModeResult is the margin configuration returned by Set-Margin-Mode.
type MarginModeResult struct {
	Symbol        string     `json:"symbol"`
	MarginCoin    string     `json:"marginCoin"`
	LongLeverage  string     `json:"longLeverage"`
	ShortLeverage string     `json:"shortLeverage"`
	MarginMode    MarginMode `json:"marginMode"`
}

// UnionConvertService -- POST /api/v2/mix/account/union-convert (signed)
//
// Converts a coin balance into USDT within the unified (multi-asset) futures
// account.
type UnionConvertService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewUnionConvertService(coin string, amount decimal.Decimal) *UnionConvertService {
	return &UnionConvertService{c: c, body: map[string]any{
		"coin":   coin,
		"amount": amount.String(),
	}}
}

func (s *UnionConvertService) Do(ctx context.Context) (*UnionConvertResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/union-convert", s.body).WithSign()
	return request.Do[UnionConvertResult](req)
}

// UnionConvertResult reports the USDT received from a union-convert.
type UnionConvertResult struct {
	USDTAmount decimal.Decimal `json:"usdtAmount"`
}

// SetPositionModeService -- POST /api/v2/mix/account/set-position-mode (signed)
//
// Switches the product line between one-way and hedge (two-way) position mode.
type SetPositionModeService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewSetPositionModeService(productType ProductType, posMode PositionMode) *SetPositionModeService {
	return &SetPositionModeService{c: c, body: map[string]any{
		"productType": string(productType),
		"posMode":     string(posMode),
	}}
}

func (s *SetPositionModeService) Do(ctx context.Context) (*PositionModeResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/account/set-position-mode", s.body).WithSign()
	return request.Do[PositionModeResult](req)
}

// PositionModeResult confirms the applied position mode.
type PositionModeResult struct {
	PosMode PositionMode `json:"posMode"`
}
