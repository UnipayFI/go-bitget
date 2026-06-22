package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetAccountAssetsService -- GET /api/v3/account/assets (UTA mgt. read)
//
// Returns the unified account's aggregate equity and per-coin balances.
type GetAccountAssetsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetAccountAssetsService() *GetAccountAssetsService {
	return &GetAccountAssetsService{c: c}
}

func (s *GetAccountAssetsService) Do(ctx context.Context) (*AccountAssets, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/assets").WithSign()
	return request.Do[AccountAssets](req)
}

type AccountAssets struct {
	AccountEquity     decimal.Decimal `json:"accountEquity"`
	UsdtEquity        decimal.Decimal `json:"usdtEquity"`
	BtcEquity         decimal.Decimal `json:"btcEquity"`
	UnrealisedPnl     decimal.Decimal `json:"unrealisedPnl"`
	UsdtUnrealisedPnl decimal.Decimal `json:"usdtUnrealisedPnl"`
	BtcUnrealizedPnl  decimal.Decimal `json:"btcUnrealizedPnl"`
	EffEquity         decimal.Decimal `json:"effEquity"`
	Mmr               decimal.Decimal `json:"mmr"`
	Imr               decimal.Decimal `json:"imr"`
	MgnRatio          decimal.Decimal `json:"mgnRatio"`
	PositionMgnRatio  decimal.Decimal `json:"positionMgnRatio"`
	Assets            []CoinAsset     `json:"assets"`
}

type CoinAsset struct {
	Coin      string          `json:"coin"`
	Equity    decimal.Decimal `json:"equity"`
	UsdValue  decimal.Decimal `json:"usdValue"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
	Debt      decimal.Decimal `json:"debt"`
	Locked    decimal.Decimal `json:"locked"`
}
