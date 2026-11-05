package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SharkFinProductStatus is the lifecycle state of a SharkFin product.
type SharkFinProductStatus string

const (
	SharkFinProductStatusNotStarted      SharkFinProductStatus = "not_started" // warming up
	SharkFinProductStatusInProgress      SharkFinProductStatus = "in_progress" // subscribing
	SharkFinProductStatusPaused          SharkFinProductStatus = "paused"
	SharkFinProductStatusInterestPending SharkFinProductStatus = "interest_pending"
	SharkFinProductStatusSettlePending   SharkFinProductStatus = "settle_pending"
	SharkFinProductStatusSettled         SharkFinProductStatus = "settled"
	SharkFinProductStatusRedeemed        SharkFinProductStatus = "redeemed"
	SharkFinProductStatusSoldOut         SharkFinProductStatus = "sold_out"
)

// SharkFinAssetStatus selects held (subscribed) or settled SharkFin assets.
type SharkFinAssetStatus string

const (
	SharkFinAssetStatusSubscribed SharkFinAssetStatus = "subscribed"
	SharkFinAssetStatusSettled    SharkFinAssetStatus = "settled"
)

// SharkFinAssetProductStatus is the state of the product behind a held
// SharkFin asset. These values differ from SharkFinProductStatus.
type SharkFinAssetProductStatus string

const (
	SharkFinAssetProductStatusRunning      SharkFinAssetProductStatus = "running"
	SharkFinAssetProductStatusPause        SharkFinAssetProductStatus = "pause"
	SharkFinAssetProductStatusWaitInterest SharkFinAssetProductStatus = "wait_interest"
	SharkFinAssetProductStatusWaitSettle   SharkFinAssetProductStatus = "wait_settle"
	SharkFinAssetProductStatusSettle       SharkFinAssetProductStatus = "settle"
	SharkFinAssetProductStatusRedeem       SharkFinAssetProductStatus = "redeem"
	SharkFinAssetProductStatusSellOut      SharkFinAssetProductStatus = "sell_out"
)

// SharkFinRecordType is the transaction type of a SharkFin record.
type SharkFinRecordType string

const (
	SharkFinRecordTypeSubscription SharkFinRecordType = "subscription"
	SharkFinRecordTypeRedemption   SharkFinRecordType = "redemption"
	SharkFinRecordTypeInterest     SharkFinRecordType = "interest"
)

// GetSharkFinProductService -- GET /api/v3/earn/sharkfin/product (signed)
//
// Returns the SharkFin products for a coin, paginated by cursor (endId).
type GetSharkFinProductService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSharkFinProductService(coin string) *GetSharkFinProductService {
	return &GetSharkFinProductService{c: c, params: map[string]string{"coin": coin}}
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSharkFinProductService) SetLimit(limit string) *GetSharkFinProductService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSharkFinProductService) SetCursor(cursor string) *GetSharkFinProductService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSharkFinProductService) Do(ctx context.Context) (*SharkFinProducts, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/product", s.params).WithSign()
	return request.Do[SharkFinProducts](req)
}

type SharkFinProducts struct {
	ResultList []SharkFinProduct `json:"resultList"`
	EndID      string            `json:"endId"`
}

// SharkFinProduct is one SharkFin product. LimitAmount is a number or
// "unlimited", so it stays a string.
type SharkFinProduct struct {
	ProductID         string                `json:"productId"`
	ProductName       string                `json:"productName"`
	ProductCoin       string                `json:"productCoin"`
	SubscribeCoin     string                `json:"subscribeCoin"`
	FarmingStartTime  time.Time             `json:"farmingStartTime,format:unixmilli"`
	FarmingEndTime    time.Time             `json:"farmingEndTime,format:unixmilli"`
	LowerRate         decimal.Decimal       `json:"lowerRate"`
	DefaultRate       decimal.Decimal       `json:"defaultRate"`
	UpperRate         decimal.Decimal       `json:"upperRate"`
	Period            string                `json:"period"`
	InterestStartTime time.Time             `json:"interestStartTime,format:unixmilli"`
	Status            SharkFinProductStatus `json:"status"`
	MinAmount         decimal.Decimal       `json:"minAmount"`
	LimitAmount       string                `json:"limitAmount"`
	SoldAmount        decimal.Decimal       `json:"soldAmount"`
	StartTime         time.Time             `json:"startTime,format:unixmilli"`
	EndTime           time.Time             `json:"endTime,format:unixmilli"`
}

// GetSharkFinAccountService -- GET /api/v3/earn/sharkfin/account (signed)
//
// Returns the SharkFin account's subscribed, historical and earned totals in
// BTC and USDT.
type GetSharkFinAccountService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSharkFinAccountService() *GetSharkFinAccountService {
	return &GetSharkFinAccountService{c: c}
}

func (s *GetSharkFinAccountService) Do(ctx context.Context) (*SharkFinAccount, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/account").WithSign()
	return request.Do[SharkFinAccount](req)
}

type SharkFinAccount struct {
	BtcSubscribeAmount   decimal.Decimal `json:"btcSubscribeAmount"`
	USDTSubscribeAmount  decimal.Decimal `json:"usdtSubscribeAmount"`
	BtcHistoricalAmount  decimal.Decimal `json:"btcHistoricalAmount"`
	USDTHistoricalAmount decimal.Decimal `json:"usdtHistoricalAmount"`
	BtcTotalEarning      decimal.Decimal `json:"btcTotalEarning"`
	USDTTotalEarning     decimal.Decimal `json:"usdtTotalEarning"`
}

// GetSharkFinAssetsService -- GET /api/v3/earn/sharkfin/assets (signed)
//
// Returns the account's subscribed or settled SharkFin assets, paginated by
// cursor (endId). The time range filters by creation time, spans at most three
// months and defaults to the last three months.
type GetSharkFinAssetsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSharkFinAssetsService(status SharkFinAssetStatus) *GetSharkFinAssetsService {
	return &GetSharkFinAssetsService{c: c, params: map[string]string{"status": string(status)}}
}

func (s *GetSharkFinAssetsService) SetStartTime(t time.Time) *GetSharkFinAssetsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSharkFinAssetsService) SetEndTime(t time.Time) *GetSharkFinAssetsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSharkFinAssetsService) SetLimit(limit string) *GetSharkFinAssetsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSharkFinAssetsService) SetCursor(cursor string) *GetSharkFinAssetsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSharkFinAssetsService) Do(ctx context.Context) (*SharkFinAssets, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/assets", s.params).WithSign()
	return request.Do[SharkFinAssets](req)
}

type SharkFinAssets struct {
	ResultList []SharkFinAsset `json:"resultList"`
	EndID      string          `json:"endId"`
}

// SharkFinAsset is one held SharkFin position.
type SharkFinAsset struct {
	ProductID         string                     `json:"productId"`
	InterestStartTime time.Time                  `json:"interestStartTime,format:unixmilli"`
	InterestEndTime   time.Time                  `json:"interestEndTime,format:unixmilli"`
	ProductCoin       string                     `json:"productCoin"`
	SubscribeCoin     string                     `json:"subscribeCoin"`
	Trend             string                     `json:"trend"` // up, down
	SettleTime        time.Time                  `json:"settleTime,format:unixmilli"`
	InterestAmount    decimal.Decimal            `json:"interestAmount"`
	ProductStatus     SharkFinAssetProductStatus `json:"productStatus"`
}

// GetSharkFinRecordsService -- GET /api/v3/earn/sharkfin/records (signed)
//
// Returns the account's SharkFin subscription/redemption/interest records of a
// type, paginated by cursor (endId). The time range spans at most three months
// and defaults to the last three months.
type GetSharkFinRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSharkFinRecordsService(recordType SharkFinRecordType) *GetSharkFinRecordsService {
	return &GetSharkFinRecordsService{c: c, params: map[string]string{"type": string(recordType)}}
}

func (s *GetSharkFinRecordsService) SetCoin(coin string) *GetSharkFinRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetSharkFinRecordsService) SetStartTime(t time.Time) *GetSharkFinRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSharkFinRecordsService) SetEndTime(t time.Time) *GetSharkFinRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSharkFinRecordsService) SetLimit(limit string) *GetSharkFinRecordsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSharkFinRecordsService) SetCursor(cursor string) *GetSharkFinRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSharkFinRecordsService) Do(ctx context.Context) (*SharkFinRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/records", s.params).WithSign()
	return request.Do[SharkFinRecords](req)
}

type SharkFinRecords struct {
	ResultList []SharkFinRecord `json:"resultList"`
	EndID      string           `json:"endId"`
}

// SharkFinRecord is one SharkFin record; Product is the product name.
type SharkFinRecord struct {
	OrderID string             `json:"orderId"`
	Product string             `json:"product"`
	Period  string             `json:"period"`
	Amount  decimal.Decimal    `json:"amount"`
	Ts      time.Time          `json:"ts,format:unixmilli"`
	Type    SharkFinRecordType `json:"type"`
}

// GetSharkFinSubscribeInfoService -- GET /api/v3/earn/sharkfin/subscribe-info (signed)
//
// Returns a SharkFin product's price range, rates, quotas and precision, for
// use before subscribing.
type GetSharkFinSubscribeInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSharkFinSubscribeInfoService(productID string) *GetSharkFinSubscribeInfoService {
	return &GetSharkFinSubscribeInfoService{c: c, params: map[string]string{"productId": productID}}
}

func (s *GetSharkFinSubscribeInfoService) Do(ctx context.Context) (*SharkFinSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/subscribe-info", s.params).WithSign()
	return request.Do[SharkFinSubscribeInfo](req)
}

// SharkFinSubscribeInfo describes a SharkFin product's subscription terms;
// InterestTime starts interest and ExpirationTime is settlement.
type SharkFinSubscribeInfo struct {
	ProductCoin        string          `json:"productCoin"`
	SubscribeCoin      string          `json:"subscribeCoin"`
	InterestTime       time.Time       `json:"interestTime,format:unixmilli"`
	ExpirationTime     time.Time       `json:"expirationTime,format:unixmilli"`
	MinPrice           decimal.Decimal `json:"minPrice"`
	CurrentPrice       decimal.Decimal `json:"currentPrice"`
	MaxPrice           decimal.Decimal `json:"maxPrice"`
	MinRate            decimal.Decimal `json:"minRate"`
	DefaultRate        decimal.Decimal `json:"defaultRate"`
	MaxRate            decimal.Decimal `json:"maxRate"`
	Period             string          `json:"period"`
	ProductMinAmount   decimal.Decimal `json:"productMinAmount"`
	AvailableBalance   decimal.Decimal `json:"availableBalance"`
	UserAmount         decimal.Decimal `json:"userAmount"`
	RemainingAmount    decimal.Decimal `json:"remainingAmount"`
	ProfitPrecision    string          `json:"profitPrecision"`
	SubscribePrecision string          `json:"subscribePrecision"`
}

// SharkFinSubscribeService -- POST /api/v3/earn/sharkfin/subscribe (signed)
//
// Subscribes to a SharkFin product, returning the subscription order.
type SharkFinSubscribeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSharkFinSubscribeService(productID string, amount decimal.Decimal) *SharkFinSubscribeService {
	return &SharkFinSubscribeService{c: c, body: map[string]any{
		"productId": productID,
		"amount":    amount.String(),
	}}
}

func (s *SharkFinSubscribeService) Do(ctx context.Context) (*SharkFinSubscribeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/sharkfin/subscribe", s.body).WithSign()
	return request.Do[SharkFinSubscribeResult](req)
}

type SharkFinSubscribeResult struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

// GetSharkFinSubscribeResultService -- GET /api/v3/earn/sharkfin/subscribe-result (signed)
//
// Returns the outcome of a SharkFin subscription order.
type GetSharkFinSubscribeResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSharkFinSubscribeResultService(orderID string) *GetSharkFinSubscribeResultService {
	return &GetSharkFinSubscribeResultService{c: c, params: map[string]string{"orderId": orderID}}
}

func (s *GetSharkFinSubscribeResultService) Do(ctx context.Context) (*SharkFinSubscribeStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/sharkfin/subscribe-result", s.params).WithSign()
	return request.Do[SharkFinSubscribeStatus](req)
}

// SharkFinSubscribeStatus is the outcome of a SharkFin subscription; Msg
// carries the error message when Result is "fail".
type SharkFinSubscribeStatus struct {
	Result string `json:"result"` // success, fail
	Msg    string `json:"msg"`
}
