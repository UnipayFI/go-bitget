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
	USDTEquity        decimal.Decimal `json:"usdtEquity"`
	BtcEquity         decimal.Decimal `json:"btcEquity"`
	UnrealizedPnL     decimal.Decimal `json:"unrealisedPnl"`
	USDTUnrealizedPnL decimal.Decimal `json:"usdtUnrealisedPnl"`
	BtcUnrealizedPnL  decimal.Decimal `json:"btcUnrealizedPnl"`
	EffEquity         decimal.Decimal `json:"effEquity"`
	Mmr               decimal.Decimal `json:"mmr"`              // maintenance margin AMOUNT (quote ccy), not a ratio
	Imr               decimal.Decimal `json:"imr"`              // initial margin AMOUNT (quote ccy), not a ratio
	MgnRatio          decimal.Decimal `json:"mgnRatio"`         // margin ratio (≈ mmr/equity)
	PositionMgnRatio  decimal.Decimal `json:"positionMgnRatio"` // position-only margin ratio
	PositionValue     decimal.Decimal `json:"positionValue"`    // position value (USD)
	Leverage          decimal.Decimal `json:"leverage"`         // account leverage (non-negative)
	Assets            []CoinAsset     `json:"assets"`
}

type CoinAsset struct {
	Coin      string          `json:"coin"`
	Equity    decimal.Decimal `json:"equity"`
	USDValue  decimal.Decimal `json:"usdValue"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
	Debt      decimal.Decimal `json:"debt"`
	Locked    decimal.Decimal `json:"locked"`
}

// GetAllFeeRateService -- GET /api/v3/account/all-fee-rate (UTA mgt. read)
//
// Returns the account's maker/taker trading fee rates across all trading pairs
// in a product category, optionally filtered to a single symbol.
type GetAllFeeRateService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAllFeeRateService(category Category) *GetAllFeeRateService {
	return &GetAllFeeRateService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetAllFeeRateService) SetSymbol(symbol string) *GetAllFeeRateService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetAllFeeRateService) Do(ctx context.Context) ([]SymbolFeeRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/all-fee-rate", s.params).WithSign()
	resp, err := request.Do[[]SymbolFeeRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SymbolFeeRate struct {
	Symbol       string          `json:"symbol"`
	MakerFeeRate decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate decimal.Decimal `json:"takerFeeRate"`
}
