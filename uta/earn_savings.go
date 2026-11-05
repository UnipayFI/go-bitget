package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SavingsPeriodType is a savings product's term: flexible (demand) or fixed.
type SavingsPeriodType string

const (
	SavingsPeriodTypeFlexible SavingsPeriodType = "flexible"
	SavingsPeriodTypeFixed    SavingsPeriodType = "fixed"
)

// SavingsProductFilter filters the savings product list by the account's
// relation to each product.
type SavingsProductFilter string

const (
	SavingsProductFilterAvailable        SavingsProductFilter = "available"
	SavingsProductFilterHeld             SavingsProductFilter = "held"
	SavingsProductFilterAvailableAndHeld SavingsProductFilter = "available_and_held"
	SavingsProductFilterAll              SavingsProductFilter = "all" // includes delisted products
)

// SavingsProductStatus is the lifecycle state of a savings product.
type SavingsProductStatus string

const (
	SavingsProductStatusNotStarted SavingsProductStatus = "not_started" // warming up
	SavingsProductStatusInProgress SavingsProductStatus = "in_progress" // subscribing
	SavingsProductStatusPaused     SavingsProductStatus = "paused"
	SavingsProductStatusCompleted  SavingsProductStatus = "completed"
	SavingsProductStatusSoldOut    SavingsProductStatus = "sold_out"
	SavingsProductStatusOffLine    SavingsProductStatus = "off_line"
)

// SavingsRecordType classifies a savings record.
type SavingsRecordType string

const (
	SavingsRecordTypeSubscribe   SavingsRecordType = "subscribe"
	SavingsRecordTypeRedeem      SavingsRecordType = "redeem"
	SavingsRecordTypePayInterest SavingsRecordType = "pay_interest"
	SavingsRecordTypeDeduction   SavingsRecordType = "deduction" // penalty interest, fixed term only
)

// GetSavingsProductService -- GET /api/v3/earn/savings/product (signed)
//
// Returns the savings product list, optionally filtered by coin and by the
// account's holding state.
type GetSavingsProductService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsProductService() *GetSavingsProductService {
	return &GetSavingsProductService{c: c, params: map[string]string{}}
}

func (s *GetSavingsProductService) SetCoin(coin string) *GetSavingsProductService {
	s.params["coin"] = coin
	return s
}

func (s *GetSavingsProductService) SetFilter(filter SavingsProductFilter) *GetSavingsProductService {
	s.params["filter"] = string(filter)
	return s
}

func (s *GetSavingsProductService) Do(ctx context.Context) ([]SavingsProduct, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/product", s.params).WithSign()
	resp, err := request.Do[[]SavingsProduct](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SavingsProduct is one savings product. Period is empty for flexible
// products; AdvanceRedeem is "Yes" or "No".
type SavingsProduct struct {
	ProductID     string               `json:"productId"`
	Coin          string               `json:"coin"`
	PeriodType    SavingsPeriodType    `json:"periodType"`
	Period        string               `json:"period"`
	ApyType       string               `json:"apyType"` // single, ladder
	AdvanceRedeem string               `json:"advanceRedeem"`
	SettleMethod  string               `json:"settleMethod"` // daily, maturity
	ApyList       []SavingsApyTier     `json:"apyList"`
	Status        SavingsProductStatus `json:"status"`
	ProductLevel  string               `json:"productLevel"` // beginner, normal, VIP
}

// SavingsApyTier is one tier of a (possibly laddered) savings rate schedule.
type SavingsApyTier struct {
	RateLevel  string          `json:"rateLevel"`
	MinStepVal decimal.Decimal `json:"minStepVal"`
	MaxStepVal decimal.Decimal `json:"maxStepVal"`
	CurrentApy decimal.Decimal `json:"currentApy"`
}

// GetSavingsAccountService -- GET /api/v3/earn/savings/account (signed)
//
// Returns the savings account's estimated value and earnings in BTC and USDT.
type GetSavingsAccountService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSavingsAccountService() *GetSavingsAccountService {
	return &GetSavingsAccountService{c: c}
}

func (s *GetSavingsAccountService) Do(ctx context.Context) (*SavingsAccount, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/account").WithSign()
	return request.Do[SavingsAccount](req)
}

type SavingsAccount struct {
	BtcAmount        decimal.Decimal `json:"btcAmount"`
	USDTAmount       decimal.Decimal `json:"usdtAmount"`
	Btc24hEarning    decimal.Decimal `json:"btc24hEarning"`
	USDT24hEarning   decimal.Decimal `json:"usdt24hEarning"`
	BtcTotalEarning  decimal.Decimal `json:"btcTotalEarning"`
	USDTTotalEarning decimal.Decimal `json:"usdtTotalEarning"`
}

// GetSavingsAssetsService -- GET /api/v3/earn/savings/assets (signed)
//
// Returns the account's held savings positions of a period type, paginated by
// cursor (endId). The time range filters by creation time, spans at most three
// months and defaults to the last three months.
type GetSavingsAssetsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsAssetsService(periodType SavingsPeriodType) *GetSavingsAssetsService {
	return &GetSavingsAssetsService{c: c, params: map[string]string{"periodType": string(periodType)}}
}

func (s *GetSavingsAssetsService) SetStartTime(t time.Time) *GetSavingsAssetsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSavingsAssetsService) SetEndTime(t time.Time) *GetSavingsAssetsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSavingsAssetsService) SetLimit(limit string) *GetSavingsAssetsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSavingsAssetsService) SetCursor(cursor string) *GetSavingsAssetsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSavingsAssetsService) Do(ctx context.Context) (*SavingsAssets, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/assets", s.params).WithSign()
	return request.Do[SavingsAssets](req)
}

type SavingsAssets struct {
	ResultList []SavingsAsset `json:"resultList"`
	EndID      string         `json:"endId"`
}

// SavingsAsset is one held savings position. OrderID is set for fixed
// positions only (pass it to SavingsRedeemService.SetOrderID); Period is empty
// for flexible ones. AllowRedemption is "Yes" or "No".
type SavingsAsset struct {
	ProductID       string            `json:"productId"`
	OrderID         string            `json:"orderId"`
	ProductCoin     string            `json:"productCoin"`
	InterestCoin    string            `json:"interestCoin"`
	PeriodType      SavingsPeriodType `json:"periodType"`
	Period          string            `json:"period"`
	HoldAmount      decimal.Decimal   `json:"holdAmount"`
	LastProfit      decimal.Decimal   `json:"lastProfit"`
	TotalProfit     decimal.Decimal   `json:"totalProfit"`
	HoldDays        string            `json:"holdDays"`
	Status          string            `json:"status"` // in_holding
	AllowRedemption string            `json:"allowRedemption"`
	ProductLevel    string            `json:"productLevel"` // beginner, normal, VIP
	Apy             []SavingsAssetApy `json:"apy"`
}

// SavingsAssetApy is one tier of a held position's rate schedule.
type SavingsAssetApy struct {
	RateLevel  string          `json:"rateLevel"`
	MinApy     decimal.Decimal `json:"minApy"`
	MaxApy     decimal.Decimal `json:"maxApy"`
	CurrentApy decimal.Decimal `json:"currentApy"`
}

// GetSavingsRecordsService -- GET /api/v3/earn/savings/records (signed)
//
// Returns the account's savings subscribe/redeem/interest/penalty records of a
// period type, paginated by cursor (endId). The time range spans at most three
// months and defaults to the last three months.
type GetSavingsRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsRecordsService(periodType SavingsPeriodType) *GetSavingsRecordsService {
	return &GetSavingsRecordsService{c: c, params: map[string]string{"periodType": string(periodType)}}
}

func (s *GetSavingsRecordsService) SetCoin(coin string) *GetSavingsRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetSavingsRecordsService) SetOrderType(orderType SavingsRecordType) *GetSavingsRecordsService {
	s.params["orderType"] = string(orderType)
	return s
}

func (s *GetSavingsRecordsService) SetStartTime(t time.Time) *GetSavingsRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSavingsRecordsService) SetEndTime(t time.Time) *GetSavingsRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSavingsRecordsService) SetLimit(limit string) *GetSavingsRecordsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSavingsRecordsService) SetCursor(cursor string) *GetSavingsRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSavingsRecordsService) Do(ctx context.Context) (*SavingsRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/records", s.params).WithSign()
	return request.Do[SavingsRecords](req)
}

type SavingsRecords struct {
	ResultList []SavingsRecord `json:"resultList"`
	EndID      string          `json:"endId"`
}

// SavingsRecord is one savings record; ProductType carries the period type and
// Period is empty for flexible products.
type SavingsRecord struct {
	OrderID        string            `json:"orderId"`
	CoinName       string            `json:"coinName"`
	SettleCoinName string            `json:"settleCoinName"`
	ProductType    SavingsPeriodType `json:"productType"`
	Period         string            `json:"period"`
	ProductLevel   string            `json:"productLevel"` // beginner, normal, VIP
	Amount         decimal.Decimal   `json:"amount"`
	Ts             time.Time         `json:"ts,format:unixmilli"`
	OrderType      SavingsRecordType `json:"orderType"`
}

// GetSavingsSubscribeInfoService -- GET /api/v3/earn/savings/subscribe-info (signed)
//
// Returns a savings product's subscription limits, precision, schedule and
// rate tiers, for use before subscribing.
type GetSavingsSubscribeInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsSubscribeInfoService(productID string, periodType SavingsPeriodType) *GetSavingsSubscribeInfoService {
	return &GetSavingsSubscribeInfoService{c: c, params: map[string]string{
		"productId":  productID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsSubscribeInfoService) Do(ctx context.Context) (*SavingsSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/subscribe-info", s.params).WithSign()
	return request.Do[SavingsSubscribeInfo](req)
}

// SavingsSubscribeInfo describes a product's subscription terms. RedeemDelay is
// the redemption arrival delay ("D+1" means one day).
type SavingsSubscribeInfo struct {
	SingleMinAmount    decimal.Decimal  `json:"singleMinAmount"`
	SingleMaxAmount    decimal.Decimal  `json:"singleMaxAmount"`
	RemainingAmount    decimal.Decimal  `json:"remainingAmount"`
	SubscribePrecision string           `json:"subscribePrecision"`
	ProfitPrecision    string           `json:"profitPrecision"`
	SubscribeTime      time.Time        `json:"subscribeTime,format:unixmilli"`
	InterestTime       time.Time        `json:"interestTime,format:unixmilli"`
	SettleTime         time.Time        `json:"settleTime,format:unixmilli"`
	ExpireTime         time.Time        `json:"expireTime,format:unixmilli"`
	RedeemTime         time.Time        `json:"redeemTime,format:unixmilli"`
	SettleMethod       string           `json:"settleMethod"`
	ApyList            []SavingsApyTier `json:"apyList"`
	RedeemDelay        string           `json:"redeemDelay"`
}

// SavingsSubscribeService -- POST /api/v3/earn/savings/subscribe (signed)
//
// Subscribes to a savings product, returning the subscription order ID.
type SavingsSubscribeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSavingsSubscribeService(productID string, periodType SavingsPeriodType, amount decimal.Decimal) *SavingsSubscribeService {
	return &SavingsSubscribeService{c: c, body: map[string]any{
		"productId":  productID,
		"periodType": string(periodType),
		"amount":     amount.String(),
	}}
}

func (s *SavingsSubscribeService) Do(ctx context.Context) (*SavingsSubscribeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/savings/subscribe", s.body).WithSign()
	return request.Do[SavingsSubscribeResult](req)
}

type SavingsSubscribeResult struct {
	OrderID string `json:"orderId"`
}

// GetSavingsSubscribeResultService -- GET /api/v3/earn/savings/subscribe-result (signed)
//
// Returns the outcome of a savings subscription order.
type GetSavingsSubscribeResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsSubscribeResultService(orderID string, periodType SavingsPeriodType) *GetSavingsSubscribeResultService {
	return &GetSavingsSubscribeResultService{c: c, params: map[string]string{
		"orderId":    orderID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsSubscribeResultService) Do(ctx context.Context) (*SavingsOrderStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/subscribe-result", s.params).WithSign()
	return request.Do[SavingsOrderStatus](req)
}

// SavingsOrderStatus is the outcome of a savings subscribe or redeem order; Msg
// carries the error message when Result is "fail".
type SavingsOrderStatus struct {
	Result string `json:"result"` // success, fail
	Msg    string `json:"msg"`
}

// SavingsRedeemService -- POST /api/v3/earn/savings/redeem (signed)
//
// Redeems a savings position. Redemptions must be at least one minute apart.
type SavingsRedeemService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSavingsRedeemService(productID string, periodType SavingsPeriodType, amount decimal.Decimal) *SavingsRedeemService {
	return &SavingsRedeemService{c: c, body: map[string]any{
		"productId":  productID,
		"periodType": string(periodType),
		"amount":     amount.String(),
	}}
}

// SetOrderID sets the assets order ID to redeem (SavingsAsset.OrderID, fixed
// positions).
func (s *SavingsRedeemService) SetOrderID(orderID string) *SavingsRedeemService {
	s.body["orderId"] = orderID
	return s
}

func (s *SavingsRedeemService) Do(ctx context.Context) (*SavingsRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/savings/redeem", s.body).WithSign()
	return request.Do[SavingsRedeemResult](req)
}

type SavingsRedeemResult struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

// GetSavingsRedeemResultService -- GET /api/v3/earn/savings/redeem-result (signed)
//
// Returns the outcome of a savings redemption order.
type GetSavingsRedeemResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSavingsRedeemResultService(orderID string, periodType SavingsPeriodType) *GetSavingsRedeemResultService {
	return &GetSavingsRedeemResultService{c: c, params: map[string]string{
		"orderId":    orderID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsRedeemResultService) Do(ctx context.Context) (*SavingsOrderStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/savings/redeem-result", s.params).WithSign()
	return request.Do[SavingsOrderStatus](req)
}
