package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetAccountFundingAssetsService -- GET /api/v3/account/funding-assets (UTA mgt. read)
//
// Returns the funding (P2P) account balances, optionally filtered to one coin.
// Only coins with assets are returned.
type GetAccountFundingAssetsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAccountFundingAssetsService() *GetAccountFundingAssetsService {
	return &GetAccountFundingAssetsService{c: c, params: map[string]string{}}
}

func (s *GetAccountFundingAssetsService) SetCoin(coin string) *GetAccountFundingAssetsService {
	s.params["coin"] = coin
	return s
}

func (s *GetAccountFundingAssetsService) Do(ctx context.Context) ([]FundingAsset, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/funding-assets", s.params).WithSign()
	resp, err := request.Do[[]FundingAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FundingAsset struct {
	Coin      string          `json:"coin"`
	Available decimal.Decimal `json:"available"`
	Frozen    decimal.Decimal `json:"frozen"`
	Balance   decimal.Decimal `json:"balance"`
}

// GetMaxTransferableService -- GET /api/v3/account/max-transferable (UTA mgt. read)
//
// Returns the maximum amount of a coin transferable out of the unified account,
// including the portion that may be transferred against borrowed funds.
type GetMaxTransferableService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMaxTransferableService(coin string) *GetMaxTransferableService {
	return &GetMaxTransferableService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetMaxTransferableService) Do(ctx context.Context) (*MaxTransferable, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/max-transferable", s.params).WithSign()
	return request.Do[MaxTransferable](req)
}

type MaxTransferable struct {
	Coin              string          `json:"coin"`
	MaxTransfer       decimal.Decimal `json:"maxTransfer"`
	BorrowMaxTransfer decimal.Decimal `json:"borrowMaxTransfer"`
}

// GetMaxWithdrawalService -- GET /api/v3/account/max-withdrawal (UTA mgt. read)
//
// Returns the maximum withdrawable amount of a coin, broken down by sub-account
// type. The figures are computed in real time and are for reference only; a
// secondary check is applied at withdrawal time.
type GetMaxWithdrawalService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMaxWithdrawalService(coin string) *GetMaxWithdrawalService {
	return &GetMaxWithdrawalService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetMaxWithdrawalService) Do(ctx context.Context) (*MaxWithdrawal, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/max-withdrawal", s.params).WithSign()
	return request.Do[MaxWithdrawal](req)
}

type MaxWithdrawal struct {
	Coin               string          `json:"coin"`
	OtcMaxWithdrawal   decimal.Decimal `json:"otcMaxWithdrawal"`
	SpotMaxWithdrawal  decimal.Decimal `json:"spotMaxWithdrawal"`
	UtaMaxWithdrawal   decimal.Decimal `json:"utaMaxWithdrawal"`
	TotalMaxWithdrawal decimal.Decimal `json:"totalMaxWithdrawal"`
}
