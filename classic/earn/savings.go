package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SavingsPeriodType selects the savings deposit term: flexible (demand) or fixed.
type SavingsPeriodType string

const (
	SavingsPeriodTypeFlexible SavingsPeriodType = "flexible"
	SavingsPeriodTypeFixed    SavingsPeriodType = "fixed"
)

// SavingsProductFilter filters the savings product list by the account's
// holding state relative to each product.
type SavingsProductFilter string

const (
	SavingsProductFilterAvailable        SavingsProductFilter = "available"
	SavingsProductFilterHeld             SavingsProductFilter = "held"
	SavingsProductFilterAvailableAndHeld SavingsProductFilter = "available_and_held"
	SavingsProductFilterAll              SavingsProductFilter = "all"
)

// SavingsRecordType classifies an entry in the savings records list.
type SavingsRecordType string

const (
	SavingsRecordTypeSubscribe   SavingsRecordType = "subscribe"
	SavingsRecordTypeRedeem      SavingsRecordType = "redeem"
	SavingsRecordTypePayInterest SavingsRecordType = "pay_interest"
	SavingsRecordTypeDeduction   SavingsRecordType = "deduction"
)

// GetSavingsProductService -- GET /api/v2/earn/savings/product (signed)
//
// Lists the available savings products, optionally filtered by coin and by the
// account's holding state.
type GetSavingsProductService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsProductService() *GetSavingsProductService {
	return &GetSavingsProductService{c: c, params: map[string]string{}}
}

// SetCoin filters products to a single subscription coin (e.g. BTC).
func (s *GetSavingsProductService) SetCoin(coin string) *GetSavingsProductService {
	s.params["coin"] = coin
	return s
}

// SetFilter filters products by holding state (available, held,
// available_and_held, all).
func (s *GetSavingsProductService) SetFilter(filter SavingsProductFilter) *GetSavingsProductService {
	s.params["filter"] = string(filter)
	return s
}

func (s *GetSavingsProductService) Do(ctx context.Context) ([]SavingsProduct, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/product", s.params).WithSign()
	resp, err := request.Do[[]SavingsProduct](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SavingsProduct is one entry in the savings product list.
type SavingsProduct struct {
	ProductID     string              `json:"productId"`
	Coin          string              `json:"coin"`
	PeriodType    SavingsPeriodType   `json:"periodType"`
	Period        string              `json:"period"`  // empty for flexible
	APYType       string              `json:"apyType"` // single, ladder
	AdvanceRedeem string              `json:"advanceRedeem"`
	SettleMethod  string              `json:"settleMethod"` // daily, maturity
	APYList       []SavingsProductAPY `json:"apyList"`
	Status        string              `json:"status"`
	ProductLevel  string              `json:"productLevel"`
}

// SavingsProductAPY is one tier of a savings product's (possibly laddered)
// interest rate schedule.
type SavingsProductAPY struct {
	RateLevel  string          `json:"rateLevel"`
	MinStepVal decimal.Decimal `json:"minStepVal"`
	MaxStepVal decimal.Decimal `json:"maxStepVal"`
	CurrentApy decimal.Decimal `json:"currentApy"`
}

// GetSavingsAccountService -- GET /api/v2/earn/savings/account (signed)
//
// Returns the aggregate savings-account valuation and earnings.
type GetSavingsAccountService struct {
	c *EarnClient
}

func (c *EarnClient) NewGetSavingsAccountService() *GetSavingsAccountService {
	return &GetSavingsAccountService{c: c}
}

func (s *GetSavingsAccountService) Do(ctx context.Context) (*SavingsAccount, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/account").WithSign()
	return request.Do[SavingsAccount](req)
}

// SavingsAccount is the aggregate savings-account valuation.
type SavingsAccount struct {
	BtcAmount        decimal.Decimal `json:"btcAmount"`
	UsdtAmount       decimal.Decimal `json:"usdtAmount"`
	Btc24hEarning    decimal.Decimal `json:"btc24hEarning"`
	Usdt24hEarning   decimal.Decimal `json:"usdt24hEarning"`
	BtcTotalEarning  decimal.Decimal `json:"btcTotalEarning"`
	UsdtTotalEarning decimal.Decimal `json:"usdtTotalEarning"`
}

// GetSavingsAssetsService -- GET /api/v2/earn/savings/assets (signed)
//
// Returns the account's held savings positions for a given period type.
type GetSavingsAssetsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsAssetsService(periodType SavingsPeriodType) *GetSavingsAssetsService {
	return &GetSavingsAssetsService{c: c, params: map[string]string{"periodType": string(periodType)}}
}

// SetStartTime filters positions at or after t (max 3-month span).
func (s *GetSavingsAssetsService) SetStartTime(t time.Time) *GetSavingsAssetsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters positions at or before t (defaults to now).
func (s *GetSavingsAssetsService) SetEndTime(t time.Time) *GetSavingsAssetsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of rows returned (default 20, max 100).
func (s *GetSavingsAssetsService) SetLimit(limit int) *GetSavingsAssetsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan returns rows older than the given pagination cursor.
func (s *GetSavingsAssetsService) SetIDLessThan(id string) *GetSavingsAssetsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSavingsAssetsService) Do(ctx context.Context) (*SavingsAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/assets", s.params).WithSign()
	return request.Do[SavingsAssets](req)
}

// SavingsAssets is the paginated list of held savings positions.
type SavingsAssets struct {
	ResultList []SavingsAsset `json:"resultList"`
	EndID      string         `json:"endId"`
}

// SavingsAsset is a single held savings position.
type SavingsAsset struct {
	ProductID       string            `json:"productId"`
	OrderID         string            `json:"orderId"` // only for fixed
	ProductCoin     string            `json:"productCoin"`
	InterestCoin    string            `json:"interestCoin"`
	PeriodType      SavingsPeriodType `json:"periodType"`
	Period          string            `json:"period"`
	HoldAmount      decimal.Decimal   `json:"holdAmount"`
	LastProfit      decimal.Decimal   `json:"lastProfit"`
	TotalProfit     decimal.Decimal   `json:"totalProfit"`
	HoldDays        string            `json:"holdDays"`
	Status          string            `json:"status"`
	AllowRedemption string            `json:"allowRedemption"`
	ProductLevel    string            `json:"productLevel"`
	Apy             []SavingsAssetAPY `json:"apy"`
}

// SavingsAssetAPY is one tier of a held position's interest rate schedule.
type SavingsAssetAPY struct {
	RateLevel  string          `json:"rateLevel"`
	MinApy     decimal.Decimal `json:"minApy"`
	MaxApy     decimal.Decimal `json:"maxApy"`
	CurrentApy decimal.Decimal `json:"currentApy"`
}

// GetSavingsRecordsService -- GET /api/v2/earn/savings/records (signed)
//
// Returns the account's savings activity (subscribe / redeem / interest /
// deduction) for a given period type.
type GetSavingsRecordsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsRecordsService(periodType SavingsPeriodType) *GetSavingsRecordsService {
	return &GetSavingsRecordsService{c: c, params: map[string]string{"periodType": string(periodType)}}
}

// SetCoin filters records to a single subscription coin.
func (s *GetSavingsRecordsService) SetCoin(coin string) *GetSavingsRecordsService {
	s.params["coin"] = coin
	return s
}

// SetOrderType filters by record type (subscribe, redeem, pay_interest,
// deduction).
func (s *GetSavingsRecordsService) SetOrderType(orderType SavingsRecordType) *GetSavingsRecordsService {
	s.params["orderType"] = string(orderType)
	return s
}

// SetStartTime filters records at or after t (max 3-month span).
func (s *GetSavingsRecordsService) SetStartTime(t time.Time) *GetSavingsRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 3-month span).
func (s *GetSavingsRecordsService) SetEndTime(t time.Time) *GetSavingsRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of rows returned (default 20, max 100).
func (s *GetSavingsRecordsService) SetLimit(limit int) *GetSavingsRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan returns rows older than the given pagination cursor.
func (s *GetSavingsRecordsService) SetIDLessThan(id string) *GetSavingsRecordsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSavingsRecordsService) Do(ctx context.Context) (*SavingsRecords, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/records", s.params).WithSign()
	return request.Do[SavingsRecords](req)
}

// SavingsRecords is the paginated list of savings activity entries.
type SavingsRecords struct {
	ResultList []SavingsRecord `json:"resultList"`
	EndID      string          `json:"endId"`
}

// SavingsRecord is a single savings activity entry.
type SavingsRecord struct {
	OrderID        string            `json:"orderId"`
	CoinName       string            `json:"coinName"`
	SettleCoinName string            `json:"settleCoinName"`
	ProductType    SavingsPeriodType `json:"productType"`
	Period         string            `json:"period"` // not returned for flexible
	ProductLevel   string            `json:"productLevel"`
	Amount         decimal.Decimal   `json:"amount"`
	Ts             time.Time         `json:"ts"`
	OrderType      SavingsRecordType `json:"orderType"`
}

// GetSavingsSubscribeInfoService -- GET /api/v2/earn/savings/subscribe-info (signed)
//
// Returns the subscription limits, precision, schedule and rate tiers for a
// savings product.
type GetSavingsSubscribeInfoService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsSubscribeInfoService(productID string, periodType SavingsPeriodType) *GetSavingsSubscribeInfoService {
	return &GetSavingsSubscribeInfoService{c: c, params: map[string]string{
		"productId":  productID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsSubscribeInfoService) Do(ctx context.Context) (*SavingsSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/subscribe-info", s.params).WithSign()
	return request.Do[SavingsSubscribeInfo](req)
}

// SavingsSubscribeInfo describes the subscription parameters for a product.
type SavingsSubscribeInfo struct {
	SingleMinAmount    decimal.Decimal     `json:"singleMinAmount"`
	SingleMaxAmount    decimal.Decimal     `json:"singleMaxAmount"`
	RemainingAmount    decimal.Decimal     `json:"remainingAmount"`
	SubscribePrecision string              `json:"subscribePrecision"`
	ProfitPrecision    string              `json:"profitPrecision"`
	SubscribeTime      time.Time           `json:"subscribeTime"`
	InterestTime       time.Time           `json:"interestTime"`
	SettleTime         time.Time           `json:"settleTime"`
	ExpireTime         time.Time           `json:"expireTime"`
	RedeemTime         time.Time           `json:"redeemTime"`
	SettleMethod       string              `json:"settleMethod"`
	APYList            []SavingsProductAPY `json:"apyList"`
	RedeemDelay        string              `json:"redeemDelay"` // e.g. "D+1"
}

// SubscribeSavingsService -- POST /api/v2/earn/savings/subscribe (signed, state-changing)
//
// Subscribes the account to a savings product. State-changing: not exercised by
// tests.
type SubscribeSavingsService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewSubscribeSavingsService(productID string, periodType SavingsPeriodType, amount decimal.Decimal) *SubscribeSavingsService {
	return &SubscribeSavingsService{c: c, body: map[string]any{
		"productId":  productID,
		"periodType": string(periodType),
		"amount":     amount.String(),
	}}
}

func (s *SubscribeSavingsService) Do(ctx context.Context) (*SavingsSubscribeResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/savings/subscribe", s.body).WithSign()
	return request.Do[SavingsSubscribeResult](req)
}

// SavingsSubscribeResult is the response to a subscription request.
type SavingsSubscribeResult struct {
	OrderID string `json:"orderId"`
}

// GetSavingsSubscribeResultService -- GET /api/v2/earn/savings/subscribe-result (signed)
//
// Returns the outcome of a savings subscription, identified by the orderId
// returned from Subscribe. (Bitget's doc parameter table mislabels this as
// productId, but the live API and the doc's own request example use orderId.)
type GetSavingsSubscribeResultService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsSubscribeResultService(orderID string, periodType SavingsPeriodType) *GetSavingsSubscribeResultService {
	return &GetSavingsSubscribeResultService{c: c, params: map[string]string{
		"orderId":    orderID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsSubscribeResultService) Do(ctx context.Context) (*SavingsSubscribeStatus, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/subscribe-result", s.params).WithSign()
	return request.Do[SavingsSubscribeStatus](req)
}

// SavingsSubscribeStatus is the outcome of a subscription request.
type SavingsSubscribeStatus struct {
	Result string `json:"result"` // success, fail
	Msg    string `json:"msg"`    // error message when result is fail
}

// RedeemSavingsService -- POST /api/v2/earn/savings/redeem (signed, state-changing)
//
// Redeems a held savings position. State-changing: not exercised by tests.
type RedeemSavingsService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewRedeemSavingsService(productID string, periodType SavingsPeriodType, amount decimal.Decimal) *RedeemSavingsService {
	return &RedeemSavingsService{c: c, body: map[string]any{
		"productId":  productID,
		"periodType": string(periodType),
		"amount":     amount.String(),
	}}
}

// SetOrderID targets a specific assets order ID (required for fixed positions;
// obtained from the savings assets endpoint).
func (s *RedeemSavingsService) SetOrderID(orderID string) *RedeemSavingsService {
	s.body["orderId"] = orderID
	return s
}

func (s *RedeemSavingsService) Do(ctx context.Context) (*SavingsRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/savings/redeem", s.body).WithSign()
	return request.Do[SavingsRedeemResult](req)
}

// SavingsRedeemResult is the response to a redemption request.
type SavingsRedeemResult struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

// GetSavingsRedeemResultService -- GET /api/v2/earn/savings/redeem-result (signed)
//
// Returns the outcome of a savings redemption.
type GetSavingsRedeemResultService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSavingsRedeemResultService(orderID string, periodType SavingsPeriodType) *GetSavingsRedeemResultService {
	return &GetSavingsRedeemResultService{c: c, params: map[string]string{
		"orderId":    orderID,
		"periodType": string(periodType),
	}}
}

func (s *GetSavingsRedeemResultService) Do(ctx context.Context) (*SavingsRedeemStatus, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/savings/redeem-result", s.params).WithSign()
	return request.Do[SavingsRedeemStatus](req)
}

// SavingsRedeemStatus is the outcome of a redemption request.
type SavingsRedeemStatus struct {
	Result string `json:"result"` // success, fail
	Msg    string `json:"msg"`    // error message when result is fail
}
