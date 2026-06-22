package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SetLeverageService -- POST /api/v3/account/set-leverage (UTA mgt. read & write)
//
// Adjusts the leverage multiple for a margin coin or futures symbol. The reply
// data is the literal string "success".
type SetLeverageService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSetLeverageService(category Category, leverage string) *SetLeverageService {
	return &SetLeverageService{c: c, body: map[string]any{
		"category": string(category),
		"leverage": leverage,
	}}
}

// SetSymbol sets the symbol (required for futures leverage adjustment).
func (s *SetLeverageService) SetSymbol(symbol string) *SetLeverageService {
	s.body["symbol"] = symbol
	return s
}

// SetCoin sets the coin (required for margin trading).
func (s *SetLeverageService) SetCoin(coin string) *SetLeverageService {
	s.body["coin"] = coin
	return s
}

// SetPosSide sets the position side (required for isolated margin).
func (s *SetLeverageService) SetPosSide(posSide PosSide) *SetLeverageService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetMarginMode sets the margin mode (defaults to cross margin).
func (s *SetLeverageService) SetMarginMode(marginMode MarginMode) *SetLeverageService {
	s.body["marginMode"] = string(marginMode)
	return s
}

// SetLongLeverage sets the long position leverage (isolated margin, two-way
// mode only; takes precedence over the leverage parameter).
func (s *SetLeverageService) SetLongLeverage(longLeverage string) *SetLeverageService {
	s.body["longLeverage"] = longLeverage
	return s
}

// SetShortLeverage sets the short position leverage (isolated margin, two-way
// mode only; takes precedence over the leverage parameter).
func (s *SetLeverageService) SetShortLeverage(shortLeverage string) *SetLeverageService {
	s.body["shortLeverage"] = shortLeverage
	return s
}

func (s *SetLeverageService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/set-leverage", s.body).WithSign()
	return request.Do[string](req)
}

// SetHoldModeService -- POST /api/v3/account/set-hold-mode (UTA mgt. read & write)
//
// Switches the futures position holding mode between one-way and hedge. The
// reply data is the literal string "success".
type SetHoldModeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSetHoldModeService(holdMode HoldMode) *SetHoldModeService {
	return &SetHoldModeService{c: c, body: map[string]any{
		"holdMode": string(holdMode),
	}}
}

func (s *SetHoldModeService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/set-hold-mode", s.body).WithSign()
	return request.Do[string](req)
}

// SetMarginService -- POST /api/v3/account/set-margin (UTA mgt. read & write)
//
// Adds or removes isolated-position margin for a futures symbol. The reply data
// is the literal string "success".
type SetMarginService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSetMarginService(category Category, symbol string, posSide PosSide, operation string, amount decimal.Decimal) *SetMarginService {
	return &SetMarginService{c: c, body: map[string]any{
		"category":  string(category),
		"symbol":    symbol,
		"posSide":   string(posSide),
		"operation": operation,
		"amount":    amount.String(),
	}}
}

func (s *SetMarginService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/set-margin", s.body).WithSign()
	return request.Do[string](req)
}

// AdjustAccountModeService -- POST /api/v3/account/adjust-account-mode (UTA mgt. read & write)
//
// Switches the unified account's trading mode (basic, advanced, delta, or
// isolated). The reply data is null on success.
type AdjustAccountModeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewAdjustAccountModeService(mode string) *AdjustAccountModeService {
	return &AdjustAccountModeService{c: c, body: map[string]any{
		"mode": mode,
	}}
}

// SetTargetUid sets the target account UID for sub-account operations (defaults
// to the current account when omitted).
func (s *AdjustAccountModeService) SetTargetUid(targetUid string) *AdjustAccountModeService {
	s.body["targetUid"] = targetUid
	return s
}

func (s *AdjustAccountModeService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/adjust-account-mode", s.body).WithSign()
	return request.Do[any](req)
}

// SwitchDeductService -- POST /api/v3/account/switch-deduct (UTA mgt. read & write)
//
// Enables or disables BGB fee deduction (deduct is "on" or "off"). The reply
// data is a boolean that is true when the change is applied successfully.
type SwitchDeductService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSwitchDeductService(deduct string) *SwitchDeductService {
	return &SwitchDeductService{c: c, body: map[string]any{
		"deduct": deduct,
	}}
}

func (s *SwitchDeductService) Do(ctx context.Context) (*bool, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/switch-deduct", s.body).WithSign()
	return request.Do[bool](req)
}

// SetDepositAccountService -- POST /api/v3/account/deposit-account (UTA mgt. read & write)
//
// Sets the target account a coin's deposits are credited to. The reply data is
// the literal string "success".
type SetDepositAccountService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSetDepositAccountService(coin, accountType string) *SetDepositAccountService {
	return &SetDepositAccountService{c: c, body: map[string]any{
		"coin":        coin,
		"accountType": accountType,
	}}
}

func (s *SetDepositAccountService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/deposit-account", s.body).WithSign()
	return request.Do[string](req)
}

// SwitchAccountService -- POST /api/v3/account/switch (UTA mgt. read & write)
//
// Switches a parent account to classic account mode. The reply data is null;
// switching completes asynchronously (poll the switch-status endpoint). Only
// available for parent accounts.
type SwitchAccountService struct {
	c *UTAClient
}

func (c *UTAClient) NewSwitchAccountService() *SwitchAccountService {
	return &SwitchAccountService{c: c}
}

func (s *SwitchAccountService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/switch").WithSign()
	return request.Do[any](req)
}
