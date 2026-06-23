package earn

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetEarnAccountAssetsService -- GET /api/v2/earn/account/assets (earn account read)
//
// Returns the per-coin Earn (savings) assets overview for the account, optionally
// filtered to a single coin.
type GetEarnAccountAssetsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetEarnAccountAssetsService() *GetEarnAccountAssetsService {
	return &GetEarnAccountAssetsService{c: c, params: map[string]string{}}
}

// SetCoin filters the overview to a single asset coin.
func (s *GetEarnAccountAssetsService) SetCoin(coin string) *GetEarnAccountAssetsService {
	s.params["coin"] = coin
	return s
}

func (s *GetEarnAccountAssetsService) Do(ctx context.Context) ([]EarnAccountAsset, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/account/assets", s.params).WithSign()
	resp, err := request.Do[[]EarnAccountAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// EarnAccountAsset is one coin's Earn (savings) balance in the assets overview.
type EarnAccountAsset struct {
	Coin   string          `json:"coin"`
	Amount decimal.Decimal `json:"amount"`
}
