package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// EliteRedeemType is the redemption mode for an Elite (on-chain) product.
type EliteRedeemType string

const (
	EliteRedeemTypeFast     EliteRedeemType = "fast"
	EliteRedeemTypeStandard EliteRedeemType = "standard"
)

// EliteAccount is a subscription/redemption funding or crediting account.
type EliteAccount string

const (
	EliteAccountSpot    EliteAccount = "spot"    // Funding account
	EliteAccountUnified EliteAccount = "unified" // Unified account
)

// EliteStatus is the lifecycle state of an Elite subscribe/redeem order.
type EliteStatus string

const (
	EliteStatusSettled  EliteStatus = "settled"
	EliteStatusPending  EliteStatus = "pending"
	EliteStatusRejected EliteStatus = "rejected"
)

// EliteRecordType selects which Elite record kind to query.
type EliteRecordType string

const (
	EliteRecordTypeSubscribe EliteRecordType = "subscribe"
	EliteRecordTypeRedeem    EliteRecordType = "redeem"
	EliteRecordTypeInterest  EliteRecordType = "interest"
)

// GetEliteProductService -- GET /api/v3/earn/elite-product (UTA mgt. read)
//
// Returns the catalogue of Elite (on-chain) earn products, each with its APR
// range, eligible subscription coins and stock status.
type GetEliteProductService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetEliteProductService() *GetEliteProductService {
	return &GetEliteProductService{c: c}
}

func (s *GetEliteProductService) Do(ctx context.Context) ([]EliteProduct, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-product").WithSign()
	resp, err := request.Do[[]EliteProduct](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type EliteProduct struct {
	ProductID            string                      `json:"productId"`
	Coin                 string                      `json:"coin"`
	MinApr               decimal.Decimal             `json:"minApr"`
	MaxApr               decimal.Decimal             `json:"maxApr"`
	SubscriptionCoinList []EliteProductSubscribeCoin `json:"subscriptionCoinList"`
	SellOut              string                      `json:"sellOut"` // YES, NO
}

type EliteProductSubscribeCoin struct {
	SubscriptionCoin string          `json:"subscriptionCoin"`
	Precision        string          `json:"precision"`
	FeeRate          decimal.Decimal `json:"feeRate"`
	ExchangeRate     decimal.Decimal `json:"exchangeRate"`
}

// GetEliteAssetsService -- GET /api/v3/earn/elite-assets (UTA mgt. read)
//
// Returns the unified account's current Elite (on-chain) holdings, with the
// holding amount, USDT-equivalent value and yield for each product.
type GetEliteAssetsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetEliteAssetsService() *GetEliteAssetsService {
	return &GetEliteAssetsService{c: c}
}

func (s *GetEliteAssetsService) Do(ctx context.Context) (*EliteAssets, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-assets").WithSign()
	return request.Do[EliteAssets](req)
}

type EliteAssets struct {
	ResultList []EliteAsset `json:"resultList"`
}

type EliteAsset struct {
	ProductID         string              `json:"productId"`
	ProductCoin       string              `json:"productCoin"`
	HoldingAmount     decimal.Decimal     `json:"holdingAmount"`
	UsdtHoldingAmount decimal.Decimal     `json:"usdtHoldingAmount"`
	ExchangeRate      decimal.Decimal     `json:"exchangeRate"`
	Apr               decimal.Decimal     `json:"apr"`    // only applicable to BGSOL / BGBTC
	MinApy            decimal.Decimal     `json:"minApy"` // only applicable to BGUSD
	MaxApy            decimal.Decimal     `json:"maxApy"` // only applicable to BGUSD
	SubscriptionCoin  string              `json:"subscriptionCoin"`
	ExchangeAmount    decimal.Decimal     `json:"exchangeAmount"`
	ProjectList       []EliteAssetProject `json:"projectList"`       // only applicable to BGBTC
	UnsettledBGPoints string              `json:"unsettledBGPoints"` // only applicable to BGBTC
	InterestCoin      string              `json:"interestCoin"`      // only applicable to BGUSD
	TotalProfit       decimal.Decimal     `json:"totalProfit"`       // only applicable to BGUSD
}

type EliteAssetProject struct {
	ProjectName string `json:"projectName"`
}

// GetEliteSubscribeInfoService -- GET /api/v3/earn/elite-subscribe-info (UTA mgt. read)
//
// Returns the subscription terms for an Elite product: minimum amount, remaining
// quota, fee rate and the per-coin options available for the product.
type GetEliteSubscribeInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEliteSubscribeInfoService(productId string) *GetEliteSubscribeInfoService {
	return &GetEliteSubscribeInfoService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetEliteSubscribeInfoService) Do(ctx context.Context) (*EliteSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-subscribe-info", s.params).WithSign()
	return request.Do[EliteSubscribeInfo](req)
}

type EliteSubscribeInfo struct {
	ProductSubID         string                   `json:"productSubId"`
	MinAmount            decimal.Decimal          `json:"minAmount"`
	RemainQuota          decimal.Decimal          `json:"remainQuota"`
	ExchangeRate         decimal.Decimal          `json:"exchangeRate"` // applicable to BGBTC / BGSOL only
	ProductCoin          string                   `json:"productCoin"`
	InterestTime         time.Time                `json:"interestTime"` // only applicable to BGUSD / BGBTC
	SettleTime           time.Time                `json:"settleTime"`   // only applicable to BGUSD / BGBTC
	Precision            string                   `json:"precision"`
	FeeRate              decimal.Decimal          `json:"feeRate"`
	SubscriptionCoinList []EliteSubscribeInfoCoin `json:"subscriptionCoinList"` // BGUSD only
}

type EliteSubscribeInfoCoin struct {
	SubscriptionCoin string          `json:"subscriptionCoin"`
	RemainQuota      decimal.Decimal `json:"remainQuota"`
	Precision        string          `json:"precision"`
	FeeRate          decimal.Decimal `json:"feeRate"`
	ExchangeRate     decimal.Decimal `json:"exchangeRate"`
	MinAmount        decimal.Decimal `json:"minAmount"`
}

// GetEliteRecordsService -- GET /api/v3/earn/elite-records (UTA mgt. read)
//
// Returns the unified account's Elite subscribe/redeem/interest records for a
// record type, paginated by cursor (endId) and bounded to a 3-month window.
type GetEliteRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEliteRecordsService(recordType EliteRecordType) *GetEliteRecordsService {
	return &GetEliteRecordsService{c: c, params: map[string]string{"type": string(recordType)}}
}

// SetStartTime filters records at or after t (max 3-month range; defaults to 3
// months prior when omitted).
func (s *GetEliteRecordsService) SetStartTime(t time.Time) *GetEliteRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetEliteRecordsService) SetEndTime(t time.Time) *GetEliteRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetEliteRecordsService) SetLimit(limit string) *GetEliteRecordsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (use endId from the previous response).
func (s *GetEliteRecordsService) SetCursor(cursor string) *GetEliteRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetEliteRecordsService) Do(ctx context.Context) (*EliteRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-records", s.params).WithSign()
	return request.Do[EliteRecords](req)
}

type EliteRecords struct {
	RecordList []EliteRecord `json:"recordList"`
	EndID      string        `json:"endId"`
}

type EliteRecord struct {
	RecordID               string          `json:"recordId"`
	ProductID              string          `json:"productId"`
	Coin                   string          `json:"coin"`
	Status                 EliteStatus     `json:"status"`
	ExchangeRate           decimal.Decimal `json:"exchangeRate"`
	ReceivedCoin           string          `json:"receivedCoin"`
	ReceivedAmount         decimal.Decimal `json:"receivedAmount"`
	InvestAmount           decimal.Decimal `json:"investAmount"`
	FeeRate                decimal.Decimal `json:"feeRate"`
	RedeemType             []string        `json:"redeemType"`             // fast, standard (redeem only)
	ReceivingAccount       string          `json:"receivingAccount"`       // spot, unified (redeem only)
	ActualReceivingAccount string          `json:"actualReceivingAccount"` // redeem only
	PaymentAccount         []string        `json:"paymentAccount"`         // spot, unified (subscribe only)
	SettlePoints           string          `json:"settlePoints"`           // BGBTC (interest only)
	Fee                    decimal.Decimal `json:"fee"`                    // subscribe/redeem only
}

// EliteSubscribeService -- POST /api/v3/earn/elite-subscribe (UTA mgt. read & write)
//
// Subscribes to an Elite (on-chain) product sub-ID with the given amount,
// returning the created subscription order ID.
type EliteSubscribeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewEliteSubscribeService(productSubId string, amount decimal.Decimal) *EliteSubscribeService {
	return &EliteSubscribeService{c: c, body: map[string]any{
		"productSubId": productSubId,
		"amount":       amount.String(),
	}}
}

// SetCoin sets the subscription coin (required for BGUSD subscriptions; USDT or
// USDC for classic accounts, USDC only for unified accounts).
func (s *EliteSubscribeService) SetCoin(coin string) *EliteSubscribeService {
	s.body["coin"] = coin
	return s
}

// SetPaymentAccount sets the debit account source (defaults to the funding
// account when omitted; "unified" is unified mode only).
func (s *EliteSubscribeService) SetPaymentAccount(paymentAccount EliteAccount) *EliteSubscribeService {
	s.body["paymentAccount"] = string(paymentAccount)
	return s
}

func (s *EliteSubscribeService) Do(ctx context.Context) (*EliteSubscribeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/elite-subscribe", s.body).WithSign()
	return request.Do[EliteSubscribeResult](req)
}

type EliteSubscribeResult struct {
	OrderID string `json:"orderId"`
}

// GetEliteSubscribeResultService -- GET /api/v3/earn/elite-subscribe-result (UTA mgt. read)
//
// Returns the settlement state of an Elite subscription order.
type GetEliteSubscribeResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEliteSubscribeResultService(orderId string) *GetEliteSubscribeResultService {
	return &GetEliteSubscribeResultService{c: c, params: map[string]string{"orderId": orderId}}
}

func (s *GetEliteSubscribeResultService) Do(ctx context.Context) (*EliteSubscribeStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-subscribe-result", s.params).WithSign()
	return request.Do[EliteSubscribeStatus](req)
}

type EliteSubscribeStatus struct {
	Result EliteStatus `json:"result"` // settled, pending, rejected
}

// GetRedeemInfoService -- GET /api/v3/earn/elite-redeem-info (UTA mgt. read)
//
// Returns the redemption terms for an Elite product: profit coin, exchange rate,
// outstanding interest and the available redemption modes (fast/standard).
type GetRedeemInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRedeemInfoService(productId string) *GetRedeemInfoService {
	return &GetRedeemInfoService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetRedeemInfoService) Do(ctx context.Context) (*EliteRedeemInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/elite-redeem-info", s.params).WithSign()
	return request.Do[EliteRedeemInfo](req)
}

type EliteRedeemInfo struct {
	ProductID                string                  `json:"productId"`
	ProductSubID             string                  `json:"productSubId"`
	ProductCoin              string                  `json:"productCoin"`
	SubscriptionCoin         string                  `json:"subscriptionCoin"`
	ProfitCoin               string                  `json:"profitCoin"`
	ExchangeRate             decimal.Decimal         `json:"exchangeRate"` // used for BGBTC / BGSOL
	TotalUnPayInterestAmount decimal.Decimal         `json:"totalUnPayInterestAmount"`
	PreSettleApr             decimal.Decimal         `json:"preSettleApr"`
	ReceivedCoin             string                  `json:"receivedCoin"`    // used for BGBTC / BGSOL
	UnsettledPoints          string                  `json:"unsettledPoints"` // available for BGBTC
	BgusdReceiveCoinList     []EliteBgusdReceiveCoin `json:"bgusdReceiveCoinList"`
	RedeemModeList           []EliteRedeemMode       `json:"redeemModeList"`
}

type EliteBgusdReceiveCoin struct {
	BgusdReceiveCoin  string          `json:"bgusdReceiveCoin"`
	BgusdExchangeRate decimal.Decimal `json:"bgusdExchangeRate"`
}

type EliteRedeemMode struct {
	RedeemFeeRate   decimal.Decimal `json:"redeemFeeRate"`
	RemainQuota     decimal.Decimal `json:"remainQuota"`
	RedeemType      EliteRedeemType `json:"redeemType"` // fast, standard
	RedeemScale     string          `json:"redeemScale"`
	RedeemDelayDate string          `json:"redeemDelayDate"`
	MinRedeemAmount decimal.Decimal `json:"minRedeemAmount"`
	RedeemTime      time.Time       `json:"redeemTime"`
}

// EliteRedeemService -- POST /api/v3/earn/elite-redeem (UTA mgt. read & write)
//
// Redeems an Elite (on-chain) product holding via the chosen redemption mode,
// returning the created redemption order ID.
type EliteRedeemService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewEliteRedeemService(productId, productSubId string, redeemType EliteRedeemType, amount decimal.Decimal, receiveAccount EliteAccount) *EliteRedeemService {
	return &EliteRedeemService{c: c, body: map[string]any{
		"productId":      productId,
		"productSubId":   productSubId,
		"redeemType":     string(redeemType),
		"amount":         amount.String(),
		"receiveAccount": string(receiveAccount),
	}}
}

// SetCoin sets the subscription coin (required for BGUSD redemptions).
func (s *EliteRedeemService) SetCoin(coin string) *EliteRedeemService {
	s.body["coin"] = coin
	return s
}

// SetAdvancedSettle enables early redemption ("yes" or "no", default "no"; only
// for BGBTC).
func (s *EliteRedeemService) SetAdvancedSettle(advancedSettle string) *EliteRedeemService {
	s.body["advancedSettle"] = advancedSettle
	return s
}

func (s *EliteRedeemService) Do(ctx context.Context) (*EliteRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/elite-redeem", s.body).WithSign()
	return request.Do[EliteRedeemResult](req)
}

type EliteRedeemResult struct {
	OrderID string `json:"orderId"`
}
