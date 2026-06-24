package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SharkFinProductStatus is the lifecycle status of a SharkFin product.
type SharkFinProductStatus string

const (
	SharkFinProductStatusNotStarted      SharkFinProductStatus = "not_started"
	SharkFinProductStatusInProgress      SharkFinProductStatus = "in_progress"
	SharkFinProductStatusPaused          SharkFinProductStatus = "paused"
	SharkFinProductStatusInterestPending SharkFinProductStatus = "interest_pending"
	SharkFinProductStatusSettlePending   SharkFinProductStatus = "settle_pending"
	SharkFinProductStatusSettled         SharkFinProductStatus = "settled"
	SharkFinProductStatusRedeemed        SharkFinProductStatus = "redeemed"
	SharkFinProductStatusSoldOut         SharkFinProductStatus = "sold_out"
)

// SharkFinAssetStatus is the asset filter status for the SharkFin assets query.
type SharkFinAssetStatus string

const (
	SharkFinAssetStatusSubscribed SharkFinAssetStatus = "subscribed"
	SharkFinAssetStatusSettled    SharkFinAssetStatus = "settled"
)

// SharkFinAssetProductStatus is the running status of a held SharkFin asset.
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

// SharkFinRecordType is the transaction type for a SharkFin record query.
type SharkFinRecordType string

const (
	SharkFinRecordTypeSubscription SharkFinRecordType = "subscription"
	SharkFinRecordTypeRedemption   SharkFinRecordType = "redemption"
	SharkFinRecordTypeInterest     SharkFinRecordType = "interest"
)

// GetSharkFinProductService -- GET /api/v2/earn/sharkfin/product (earn read)
//
// Returns the list of available SharkFin products for a coin.
type GetSharkFinProductService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSharkFinProductService(coin string) *GetSharkFinProductService {
	return &GetSharkFinProductService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetSharkFinProductService) SetLimit(limit string) *GetSharkFinProductService {
	s.params["limit"] = limit
	return s
}

func (s *GetSharkFinProductService) SetIDLessThan(idLessThan string) *GetSharkFinProductService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetSharkFinProductService) Do(ctx context.Context) (*SharkFinProductResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/product", s.params).WithSign()
	return request.Do[SharkFinProductResult](req)
}

// SharkFinProductResult is the paginated SharkFin product list.
type SharkFinProductResult struct {
	ResultList []SharkFinProduct `json:"resultList"`
	EndID      string            `json:"endId"`
}

// SharkFinProduct is a single SharkFin product offering.
type SharkFinProduct struct {
	ProductID         string                `json:"productId"`
	ProductName       string                `json:"productName"`
	ProductCoin       string                `json:"productCoin"`
	SubscribeCoin     string                `json:"subscribeCoin"`
	FarmingStartTime  time.Time             `json:"farmingStartTime"`
	FarmingEndTime    time.Time             `json:"farmingEndTime"`
	LowerRate         decimal.Decimal       `json:"lowerRate"`
	DefaultRate       decimal.Decimal       `json:"defaultRate"`
	UpperRate         decimal.Decimal       `json:"upperRate"`
	Period            string                `json:"period"`
	InterestStartTime time.Time             `json:"interestStartTime"`
	Status            SharkFinProductStatus `json:"status"`
	MinAmount         decimal.Decimal       `json:"minAmount"`
	LimitAmount       string                `json:"limitAmount"` // numeric amount or "unlimited"
	SoldAmount        decimal.Decimal       `json:"soldAmount"`
	StartTime         time.Time             `json:"startTime"`
	EndTime           time.Time             `json:"endTime"`
}

// GetSharkFinAccountService -- GET /api/v2/earn/sharkfin/account (earn read)
//
// Returns the SharkFin account summary (subscription and earning totals).
type GetSharkFinAccountService struct {
	c *EarnClient
}

func (c *EarnClient) NewGetSharkFinAccountService() *GetSharkFinAccountService {
	return &GetSharkFinAccountService{c: c}
}

func (s *GetSharkFinAccountService) Do(ctx context.Context) (*SharkFinAccount, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/account").WithSign()
	return request.Do[SharkFinAccount](req)
}

// SharkFinAccount is the aggregate SharkFin account summary.
type SharkFinAccount struct {
	BtcSubscribeAmount   decimal.Decimal `json:"btcSubscribeAmount"`
	USDTSubscribeAmount  decimal.Decimal `json:"usdtSubscribeAmount"`
	BtcHistoricalAmount  decimal.Decimal `json:"btcHistoricalAmount"`
	USDTHistoricalAmount decimal.Decimal `json:"usdtHistoricalAmount"`
	BtcTotalEarning      decimal.Decimal `json:"btcTotalEarning"`
	USDTTotalEarning     decimal.Decimal `json:"usdtTotalEarning"`
}

// GetSharkFinAssetsService -- GET /api/v2/earn/sharkfin/assets (earn read)
//
// Returns the user's held SharkFin assets filtered by status.
type GetSharkFinAssetsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSharkFinAssetsService(status SharkFinAssetStatus) *GetSharkFinAssetsService {
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

func (s *GetSharkFinAssetsService) SetLimit(limit string) *GetSharkFinAssetsService {
	s.params["limit"] = limit
	return s
}

func (s *GetSharkFinAssetsService) SetIDLessThan(idLessThan string) *GetSharkFinAssetsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetSharkFinAssetsService) Do(ctx context.Context) (*SharkFinAssetsResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/assets", s.params).WithSign()
	return request.Do[SharkFinAssetsResult](req)
}

// SharkFinAssetsResult is the paginated SharkFin assets list.
type SharkFinAssetsResult struct {
	ResultList []SharkFinAsset `json:"resultList"`
	EndID      string          `json:"endId"`
}

// SharkFinAsset is a single held SharkFin position.
type SharkFinAsset struct {
	ProductID         string                     `json:"productId"`
	InterestStartTime time.Time                  `json:"interestStartTime"`
	InterestEndTime   time.Time                  `json:"interestEndTime"`
	ProductCoin       string                     `json:"productCoin"`
	SubscribeCoin     string                     `json:"subscribeCoin"`
	Trend             string                     `json:"trend"` // up, down
	SettleTime        time.Time                  `json:"settleTime"`
	InterestAmount    decimal.Decimal            `json:"interestAmount"`
	ProductStatus     SharkFinAssetProductStatus `json:"productStatus"`
}

// GetSharkFinRecordsService -- GET /api/v2/earn/sharkfin/records (earn read)
//
// Returns SharkFin transaction records (subscription, redemption, interest).
type GetSharkFinRecordsService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSharkFinRecordsService(recordType SharkFinRecordType) *GetSharkFinRecordsService {
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

func (s *GetSharkFinRecordsService) SetLimit(limit string) *GetSharkFinRecordsService {
	s.params["limit"] = limit
	return s
}

func (s *GetSharkFinRecordsService) SetIDLessThan(idLessThan string) *GetSharkFinRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetSharkFinRecordsService) Do(ctx context.Context) (*SharkFinRecordsResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/records", s.params).WithSign()
	return request.Do[SharkFinRecordsResult](req)
}

// SharkFinRecordsResult is the paginated SharkFin records list.
type SharkFinRecordsResult struct {
	ResultList []SharkFinRecord `json:"resultList"`
	EndID      string           `json:"endId"`
}

// SharkFinRecord is a single SharkFin transaction record.
type SharkFinRecord struct {
	OrderID string             `json:"orderId"`
	Product string             `json:"product"`
	Period  string             `json:"period"`
	Amount  decimal.Decimal    `json:"amount"`
	Ts      time.Time          `json:"ts"`
	Type    SharkFinRecordType `json:"type"`
}

// GetSharkFinSubscribeInfoService -- GET /api/v2/earn/sharkfin/subscribe-info (earn read)
//
// Returns the subscription detail for a single SharkFin product.
type GetSharkFinSubscribeInfoService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSharkFinSubscribeInfoService(productID string) *GetSharkFinSubscribeInfoService {
	return &GetSharkFinSubscribeInfoService{c: c, params: map[string]string{"productId": productID}}
}

func (s *GetSharkFinSubscribeInfoService) Do(ctx context.Context) (*SharkFinSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/subscribe-info", s.params).WithSign()
	return request.Do[SharkFinSubscribeInfo](req)
}

// SharkFinSubscribeInfo is the subscription detail for a SharkFin product.
type SharkFinSubscribeInfo struct {
	ProductCoin        string          `json:"productCoin"`
	SubscribeCoin      string          `json:"subscribeCoin"`
	InterestTime       time.Time       `json:"interestTime"`
	ExpirationTime     time.Time       `json:"expirationTime"`
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

// SubscribeSharkFinService -- POST /api/v2/earn/sharkfin/subscribe (earn trade)
//
// Subscribes the requested amount to a SharkFin product. STATE-CHANGING.
type SubscribeSharkFinService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewSubscribeSharkFinService(productID string, amount decimal.Decimal) *SubscribeSharkFinService {
	return &SubscribeSharkFinService{c: c, body: map[string]any{
		"productId": productID,
		"amount":    amount.String(),
	}}
}

func (s *SubscribeSharkFinService) Do(ctx context.Context) (*SharkFinSubscribeResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/sharkfin/subscribe", s.body).WithSign()
	return request.Do[SharkFinSubscribeResponse](req)
}

// SharkFinSubscribeResponse is the result of a SharkFin subscription request.
type SharkFinSubscribeResponse struct {
	OrderID string `json:"orderId"`
	Status  string `json:"status"`
}

// GetSharkFinSubscribeResultService -- GET /api/v2/earn/sharkfin/subscribe-result (earn read)
//
// Returns the result of a SharkFin subscription order.
type GetSharkFinSubscribeResultService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetSharkFinSubscribeResultService(orderID string) *GetSharkFinSubscribeResultService {
	return &GetSharkFinSubscribeResultService{c: c, params: map[string]string{"orderId": orderID}}
}

func (s *GetSharkFinSubscribeResultService) Do(ctx context.Context) (*SharkFinSubscribeResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/sharkfin/subscribe-result", s.params).WithSign()
	return request.Do[SharkFinSubscribeResult](req)
}

// SharkFinSubscribeResult is the outcome of a SharkFin subscription order.
type SharkFinSubscribeResult struct {
	Result string `json:"result"` // success, fail
	Msg    string `json:"msg"`    // error message when result is fail
}
