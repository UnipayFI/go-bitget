package common

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetFundingAssetsService -- GET /api/v2/account/funding-assets (private)
//
// Returns the balances held in the funding (P2P/fiat) account, optionally
// filtered to a single coin.
type GetFundingAssetsService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetFundingAssetsService() *GetFundingAssetsService {
	return &GetFundingAssetsService{c: c, params: map[string]string{}}
}

func (s *GetFundingAssetsService) SetCoin(coin string) *GetFundingAssetsService {
	s.params["coin"] = coin
	return s
}

func (s *GetFundingAssetsService) Do(ctx context.Context) ([]FundingAsset, error) {
	req := request.Get(ctx, s.c, "/api/v2/account/funding-assets", s.params).WithSign()
	resp, err := request.Do[[]FundingAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FundingAsset is a single coin balance in the funding account.
type FundingAsset struct {
	Coin      string          `json:"coin"`
	Available decimal.Decimal `json:"available"`
	Frozen    decimal.Decimal `json:"frozen"`
	USDTValue decimal.Decimal `json:"usdtValue"`
}

// BotAccountType is the bot product line selected by GetBotAssetsService.
type BotAccountType string

const (
	BotAccountTypeFutures BotAccountType = "futures"
	BotAccountTypeSpot    BotAccountType = "spot"
)

// GetBotAssetsService -- GET /api/v2/account/bot-assets (private)
//
// Returns the per-coin balances held in the bot (strategy/trading-bot) account,
// optionally filtered to a single bot product line.
type GetBotAssetsService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetBotAssetsService() *GetBotAssetsService {
	return &GetBotAssetsService{c: c, params: map[string]string{}}
}

func (s *GetBotAssetsService) SetAccountType(accountType BotAccountType) *GetBotAssetsService {
	s.params["accountType"] = string(accountType)
	return s
}

func (s *GetBotAssetsService) Do(ctx context.Context) ([]BotAsset, error) {
	req := request.Get(ctx, s.c, "/api/v2/account/bot-assets", s.params).WithSign()
	resp, err := request.Do[[]BotAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BotAsset is a single coin balance in the bot account.
type BotAsset struct {
	Coin      string          `json:"coin"`
	Available decimal.Decimal `json:"available"`
	Equity    decimal.Decimal `json:"equity"`
	Bonus     decimal.Decimal `json:"bonus"`
	Frozen    decimal.Decimal `json:"frozen"`
	USDTValue decimal.Decimal `json:"usdtValue"`
}

// GetAllAccountBalanceService -- GET /api/v2/account/all-account-balance (private)
//
// Returns the USDT-equivalent balance of each account type (spot, futures,
// funding, earn, bots, margin, ...). Rate-limited to 1 request/second per user.
type GetAllAccountBalanceService struct {
	c *CommonClient
}

func (c *CommonClient) NewGetAllAccountBalanceService() *GetAllAccountBalanceService {
	return &GetAllAccountBalanceService{c: c}
}

func (s *GetAllAccountBalanceService) Do(ctx context.Context) ([]AccountBalance, error) {
	req := request.Get(ctx, s.c, "/api/v2/account/all-account-balance").WithSign()
	resp, err := request.Do[[]AccountBalance](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountBalance is the USDT-equivalent balance of a single account type.
type AccountBalance struct {
	AccountType string          `json:"accountType"`
	USDTBalance decimal.Decimal `json:"usdtBalance"`
}
