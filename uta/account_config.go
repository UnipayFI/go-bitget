package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetAccountSettingsService -- GET /api/v3/account/settings (UTA mgt. read)
//
// Returns the unified account's mode settings and per-symbol / per-coin
// leverage and margin configurations.
type GetAccountSettingsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetAccountSettingsService() *GetAccountSettingsService {
	return &GetAccountSettingsService{c: c}
}

func (s *GetAccountSettingsService) Do(ctx context.Context) (*AccountSettings, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/settings").WithSign()
	return request.Do[AccountSettings](req)
}

type AccountSettings struct {
	Uid              string                 `json:"uid"`
	AccountMode      string                 `json:"accountMode"`  // unified, hybrid, upgrading, switching
	AssetMode        string                 `json:"assetMode"`    // multi_assets
	AccountLevel     string                 `json:"accountLevel"` // basic, advanced, isolated, delta
	HoldMode         HoldMode               `json:"holdMode"`
	StpMode          string                 `json:"stpMode"` // none, cancel_taker, cancel_maker, cancel_both
	SymbolConfigList []SymbolLeverageConfig `json:"symbolConfigList"`
	CoinConfigList   []CoinLeverageConfig   `json:"coinConfigList"`
}

type SymbolLeverageConfig struct {
	Category   Category   `json:"category"`
	Symbol     string     `json:"symbol"`
	MarginMode MarginMode `json:"marginMode"`
	Leverage   string     `json:"leverage"`
}

type CoinLeverageConfig struct {
	Coin     string `json:"coin"`
	Leverage string `json:"leverage"`
}

// GetAccountInfoService -- GET /api/v3/account/info (UTA mgt. read)
//
// Returns identity and permission metadata for the calling account.
type GetAccountInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetAccountInfoService() *GetAccountInfoService {
	return &GetAccountInfoService{c: c}
}

func (s *GetAccountInfoService) Do(ctx context.Context) (*AccountInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/info").WithSign()
	return request.Do[AccountInfo](req)
}

type AccountInfo struct {
	UserId      string    `json:"userId"`
	InviterId   string    `json:"inviterId"`
	ParentId    string    `json:"parentId"`
	ChannelCode string    `json:"channelCode"`
	Channel     string    `json:"channel"`
	Ips         string    `json:"ips"`
	PermType    string    `json:"permType"`    // read-only, read-and-write
	Permissions []string  `json:"permissions"` // uta_mgt, uta_trade, withdraw, copy_futures_position, copy_futures_order
	RegisTime   time.Time `json:"regisTime"`
}

// GetDeltaInfoService -- GET /api/v3/account/delta-info (UTA mgt. read)
//
// Returns the account's delta-neutral status and per-coin net position ratios
// used for ADL ranking.
type GetDeltaInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetDeltaInfoService() *GetDeltaInfoService {
	return &GetDeltaInfoService{c: c}
}

func (s *GetDeltaInfoService) Do(ctx context.Context) (*DeltaInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/delta-info").WithSign()
	return request.Do[DeltaInfo](req)
}

type DeltaInfo struct {
	DeltaEquityRatio  decimal.Decimal `json:"deltaEquityRatio"`
	DeltaThreshold    decimal.Decimal `json:"deltaThreshold"`
	PositionThreshold decimal.Decimal `json:"positionThreshold"`
	List              []DeltaCoinInfo `json:"list"`
}

type DeltaCoinInfo struct {
	Coin          string          `json:"coin"`
	PositionRatio decimal.Decimal `json:"positionRatio"`
}

// GetAccountFeeRateService -- GET /api/v3/account/fee-rate (UTA mgt. read)
//
// Returns the account's maker/taker trading fee rates for a symbol.
type GetAccountFeeRateService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAccountFeeRateService(category Category, symbol string) *GetAccountFeeRateService {
	return &GetAccountFeeRateService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetAccountFeeRateService) Do(ctx context.Context) (*AccountFeeRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/fee-rate", s.params).WithSign()
	return request.Do[AccountFeeRate](req)
}

type AccountFeeRate struct {
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate decimal.Decimal `json:"takerFeeRate"`
}

// GetOILimitService -- GET /api/v3/account/open-interest-limit (UTA mgt. read)
//
// Returns the open-interest (position) quantity limits for a futures symbol.
type GetOILimitService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOILimitService(category Category, symbol string) *GetOILimitService {
	return &GetOILimitService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetOILimitService) Do(ctx context.Context) (*OILimit, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/open-interest-limit", s.params).WithSign()
	return request.Do[OILimit](req)
}

type OILimit struct {
	Symbol           string          `json:"symbol"`
	SingleUserLimit  decimal.Decimal `json:"singleUserLimit"`
	MasterSubLimit   decimal.Decimal `json:"masterSubLimit"`
	MarketMakerLimit decimal.Decimal `json:"marketMakerLimit"`
}

// GetDeductInfoService -- GET /api/v3/account/deduct-info (UTA mgt. read)
//
// Returns whether BGB fee deduction is enabled for the account.
type GetDeductInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetDeductInfoService() *GetDeductInfoService {
	return &GetDeductInfoService{c: c}
}

func (s *GetDeductInfoService) Do(ctx context.Context) (*DeductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/deduct-info").WithSign()
	return request.Do[DeductInfo](req)
}

type DeductInfo struct {
	Deduct string `json:"deduct"` // on, off
}

// GetSwitchStatusService -- GET /api/v3/account/switch-status (UTA mgt. read)
//
// Returns the status of an in-progress account-mode switch. Only supported for
// parent accounts.
type GetSwitchStatusService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSwitchStatusService() *GetSwitchStatusService {
	return &GetSwitchStatusService{c: c}
}

func (s *GetSwitchStatusService) Do(ctx context.Context) (*SwitchStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/switch-status").WithSign()
	return request.Do[SwitchStatus](req)
}

type SwitchStatus struct {
	Status string `json:"status"` // process, success, fail
	Reason string `json:"reason"` // only present when status = fail
}
