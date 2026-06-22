package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSpotWhaleNetFlowService -- GET /api/v3/market/spot-whale-flow
//
// Returns the per-period whale (large holder) buy/sell net flow volume for a
// spot trading pair.
type GetSpotWhaleNetFlowService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSpotWhaleNetFlowService(symbol string) *GetSpotWhaleNetFlowService {
	return &GetSpotWhaleNetFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetSpotWhaleNetFlowService) Do(ctx context.Context) ([]SpotWhaleNetFlow, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/spot-whale-flow", s.params)
	resp, err := request.Do[[]SpotWhaleNetFlow](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SpotWhaleNetFlow struct {
	Volume decimal.Decimal `json:"volume"`
	Date   time.Time       `json:"date"`
}

// GetSpotFundFlowService -- GET /api/v3/market/spot-fund-flow
//
// Returns the buy/sell volumes and ratios split by holder size (whale, dolphin,
// fish) for a spot trading pair over the requested period.
type GetSpotFundFlowService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSpotFundFlowService(symbol string) *GetSpotFundFlowService {
	return &GetSpotFundFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the query interval: 15m (default), 30m, 1h, 2h, 4h, 1d.
func (s *GetSpotFundFlowService) SetPeriod(period string) *GetSpotFundFlowService {
	s.params["period"] = period
	return s
}

func (s *GetSpotFundFlowService) Do(ctx context.Context) (*SpotFundFlow, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/spot-fund-flow", s.params)
	return request.Do[SpotFundFlow](req)
}

type SpotFundFlow struct {
	WhaleBuyVolume    decimal.Decimal `json:"whaleBuyVolume"`
	DolphinBuyVolume  decimal.Decimal `json:"dolphinBuyVolume"`
	FishBuyVolume     decimal.Decimal `json:"fishBuyVolume"`
	WhaleSellVolume   decimal.Decimal `json:"whaleSellVolume"`
	DolphinSellVolume decimal.Decimal `json:"dolphinSellVolume"`
	FishSellVolume    decimal.Decimal `json:"fishSellVolume"`
	WhaleBuyRatio     decimal.Decimal `json:"whaleBuyRatio"`
	DolphinBuyRatio   decimal.Decimal `json:"dolphinBuyRatio"`
	FishBuyRatio      decimal.Decimal `json:"fishBuyRatio"`
	WhaleSellRatio    decimal.Decimal `json:"whaleSellRatio"`
	DolphinSellRatio  decimal.Decimal `json:"dolphinSellRatio"`
	FishSellRatio     decimal.Decimal `json:"fishSellRatio"`
}

// GetSpotNetFlowService -- GET /api/v3/market/spot-net-flow
//
// Returns the 24H whale net capital inflow series for a spot trading pair.
type GetSpotNetFlowService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSpotNetFlowService(symbol string) *GetSpotNetFlowService {
	return &GetSpotNetFlowService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetSpotNetFlowService) Do(ctx context.Context) ([]SpotNetFlow, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/spot-net-flow", s.params)
	resp, err := request.Do[[]SpotNetFlow](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SpotNetFlow struct {
	NetFlow decimal.Decimal `json:"netFlow"`
	Ts      time.Time       `json:"ts"`
}

// GetMarginLongShortService -- GET /api/v3/market/margin-long-short
//
// Returns the margin long/short ratio series for a trading pair.
type GetMarginLongShortService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMarginLongShortService(symbol string) *GetMarginLongShortService {
	return &GetMarginLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time period: 24h (default) or 30d.
func (s *GetMarginLongShortService) SetPeriod(period string) *GetMarginLongShortService {
	s.params["period"] = period
	return s
}

// SetCoin sets the base or quote coin; defaults to the base coin.
func (s *GetMarginLongShortService) SetCoin(coin string) *GetMarginLongShortService {
	s.params["coin"] = coin
	return s
}

func (s *GetMarginLongShortService) Do(ctx context.Context) ([]MarginLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/margin-long-short", s.params)
	resp, err := request.Do[[]MarginLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type MarginLongShort struct {
	LongShortRatio decimal.Decimal `json:"longShortRatio"`
	Ts             time.Time       `json:"ts"`
}

// GetMarginLoanGrowthService -- GET /api/v3/market/margin-loan-growth
//
// Returns the margin loan growth-rate series for a trading pair.
type GetMarginLoanGrowthService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMarginLoanGrowthService(symbol string) *GetMarginLoanGrowthService {
	return &GetMarginLoanGrowthService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time interval: 24h (default) or 30d.
func (s *GetMarginLoanGrowthService) SetPeriod(period string) *GetMarginLoanGrowthService {
	s.params["period"] = period
	return s
}

// SetCoin sets the base or quote coin; defaults to the base coin.
func (s *GetMarginLoanGrowthService) SetCoin(coin string) *GetMarginLoanGrowthService {
	s.params["coin"] = coin
	return s
}

func (s *GetMarginLoanGrowthService) Do(ctx context.Context) ([]MarginLoanGrowth, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/margin-loan-growth", s.params)
	resp, err := request.Do[[]MarginLoanGrowth](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type MarginLoanGrowth struct {
	GrowthRate decimal.Decimal `json:"growthRate"`
	Ts         time.Time       `json:"ts"`
}

// GetMarginIsolatedBorrowService -- GET /api/v3/market/margin-isolated-borrow
//
// Returns the isolated-margin borrowing-ratio series for a trading pair.
type GetMarginIsolatedBorrowService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMarginIsolatedBorrowService(symbol string) *GetMarginIsolatedBorrowService {
	return &GetMarginIsolatedBorrowService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time period: 24h (default) or 30d.
func (s *GetMarginIsolatedBorrowService) SetPeriod(period string) *GetMarginIsolatedBorrowService {
	s.params["period"] = period
	return s
}

func (s *GetMarginIsolatedBorrowService) Do(ctx context.Context) ([]MarginIsolatedBorrow, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/margin-isolated-borrow", s.params)
	resp, err := request.Do[[]MarginIsolatedBorrow](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type MarginIsolatedBorrow struct {
	BorrowRate decimal.Decimal `json:"borrowRate"`
	Ts         time.Time       `json:"ts"`
}

// GetFuturesActiveBuySellService -- GET /api/v3/market/futures-active-buy-sell
//
// Returns the per-period taker (active) buy/sell volume series for a futures
// trading pair.
type GetFuturesActiveBuySellService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFuturesActiveBuySellService(symbol string) *GetFuturesActiveBuySellService {
	return &GetFuturesActiveBuySellService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time interval: 5m (default), 15m, 30m, 1h, 2h, 4h, 6h,
// 12h, 1d.
func (s *GetFuturesActiveBuySellService) SetPeriod(period string) *GetFuturesActiveBuySellService {
	s.params["period"] = period
	return s
}

func (s *GetFuturesActiveBuySellService) Do(ctx context.Context) ([]FuturesActiveBuySell, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/futures-active-buy-sell", s.params)
	resp, err := request.Do[[]FuturesActiveBuySell](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FuturesActiveBuySell struct {
	BuyVolume  decimal.Decimal `json:"buyVolume"`
	SellVolume decimal.Decimal `json:"sellVolume"`
	Ts         time.Time       `json:"ts"`
}

// GetFuturesLongShortService -- GET /api/v3/market/futures-long-short
//
// Returns the futures long/short ratio series for a trading pair.
type GetFuturesLongShortService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFuturesLongShortService(symbol string) *GetFuturesLongShortService {
	return &GetFuturesLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time interval: 5m (default), 15m, 30m, 1h, 2h, 4h, 6h,
// 12h, 1Dutc.
func (s *GetFuturesLongShortService) SetPeriod(period string) *GetFuturesLongShortService {
	s.params["period"] = period
	return s
}

func (s *GetFuturesLongShortService) Do(ctx context.Context) ([]FuturesLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/futures-long-short", s.params)
	resp, err := request.Do[[]FuturesLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FuturesLongShort struct {
	LongRatio      decimal.Decimal `json:"longRatio"`
	ShortRatio     decimal.Decimal `json:"shortRatio"`
	LongShortRatio decimal.Decimal `json:"longShortRatio"`
	Ts             time.Time       `json:"ts"`
}

// GetFuturesPositionLongShortService -- GET /api/v3/market/futures-position-long-short
//
// Returns the futures active long/short position ratio series for a trading
// pair.
type GetFuturesPositionLongShortService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFuturesPositionLongShortService(symbol string) *GetFuturesPositionLongShortService {
	return &GetFuturesPositionLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time interval: 5m (default), 15m, 30m, 1h, 2h, 4h, 6h,
// 12h, 1d.
func (s *GetFuturesPositionLongShortService) SetPeriod(period string) *GetFuturesPositionLongShortService {
	s.params["period"] = period
	return s
}

func (s *GetFuturesPositionLongShortService) Do(ctx context.Context) ([]FuturesPositionLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/futures-position-long-short", s.params)
	resp, err := request.Do[[]FuturesPositionLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FuturesPositionLongShort struct {
	LongPositionRatio      decimal.Decimal `json:"longPositionRatio"`
	ShortPositionRatio     decimal.Decimal `json:"shortPositionRatio"`
	LongShortPositionRatio decimal.Decimal `json:"longShortPositionRatio"`
	Ts                     time.Time       `json:"ts"`
}

// GetFuturesAccountLongShortService -- GET /api/v3/market/futures-account-long-short
//
// Returns the futures active long/short account ratio series for a trading
// pair.
type GetFuturesAccountLongShortService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFuturesAccountLongShortService(symbol string) *GetFuturesAccountLongShortService {
	return &GetFuturesAccountLongShortService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPeriod sets the time interval: 5m (default), 15m, 30m, 1h, 2h, 4h, 6h,
// 12h, 1d.
func (s *GetFuturesAccountLongShortService) SetPeriod(period string) *GetFuturesAccountLongShortService {
	s.params["period"] = period
	return s
}

func (s *GetFuturesAccountLongShortService) Do(ctx context.Context) ([]FuturesAccountLongShort, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/futures-account-long-short", s.params)
	resp, err := request.Do[[]FuturesAccountLongShort](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type FuturesAccountLongShort struct {
	LongAccountRatio      decimal.Decimal `json:"longAccountRatio"`
	ShortAccountRatio     decimal.Decimal `json:"shortAccountRatio"`
	LongShortAccountRatio decimal.Decimal `json:"longShortAccountRatio"`
	Ts                    time.Time       `json:"ts"`
}
