package uta

import (
	"context"
	"strconv"
	"strings"
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
	USDTHoldingAmount decimal.Decimal     `json:"usdtHoldingAmount"`
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
	InterestTime         time.Time                `json:"interestTime,format:unixmilli"` // only applicable to BGUSD / BGBTC
	SettleTime           time.Time                `json:"settleTime,format:unixmilli"`   // only applicable to BGUSD / BGBTC
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
	BizTime                time.Time       `json:"bizTime,format:unixmilli"`
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
	RedeemTime      time.Time       `json:"redeemTime,format:unixmilli"`
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

// CashPlusAccount is a Cash+ debit (subscribe) or credit (redeem) account.
type CashPlusAccount string

const (
	CashPlusAccountFunding CashPlusAccount = "funding" // Funding account
	CashPlusAccountUTA     CashPlusAccount = "uta"     // Unified account
	CashPlusAccountOTC     CashPlusAccount = "otc"     // OTC account (subscribe debit only)
)

// CashPlusStatus is the lifecycle state of a Cash+ order or record.
type CashPlusStatus string

const (
	CashPlusStatusProcessing CashPlusStatus = "processing"
	CashPlusStatusSuccess    CashPlusStatus = "success"
	CashPlusStatusFailed     CashPlusStatus = "failed"
)

// CashPlusRecordType selects which Cash+ history record kind to query.
type CashPlusRecordType string

const (
	CashPlusRecordTypeSubscribe         CashPlusRecordType = "subscribe"
	CashPlusRecordTypeRedeem            CashPlusRecordType = "redeem"
	CashPlusRecordTypeCompoundInterest  CashPlusRecordType = "compound_interest"
	CashPlusRecordTypeLiquidationRedeem CashPlusRecordType = "liquidation_redeem"
)

// GetCashPlusSubscribeInfoService -- GET /api/v3/earn/cashplus/subscribe-info (UTA mgt. read)
//
// Returns the Cash+ subscription prerequisites: the current annualized rate and
// the subscribable coins with their product sub-IDs. Callable by both master and
// sub-accounts.
type GetCashPlusSubscribeInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCashPlusSubscribeInfoService() *GetCashPlusSubscribeInfoService {
	return &GetCashPlusSubscribeInfoService{c: c}
}

func (s *GetCashPlusSubscribeInfoService) Do(ctx context.Context) (*CashPlusSubscribeInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/subscribe-info").WithSign()
	return request.Do[CashPlusSubscribeInfo](req)
}

type CashPlusSubscribeInfo struct {
	ProfitSourceCoin  string                  `json:"profitSourceCoin"` // currently USDGO
	CashPlusApr       decimal.Decimal         `json:"cashPlusApr"`      // 0.04 means 4%
	SubscribeCoinList []CashPlusSubscribeCoin `json:"subscribeCoinList"`
	InterestTime      time.Time               `json:"interestTime,format:unixmilli"`
	PayInterestTime   time.Time               `json:"payInterestTime,format:unixmilli"`
	SubscribeTime     time.Time               `json:"subscribeTime,format:unixmilli"`
}

type CashPlusSubscribeCoin struct {
	ProductSubID   string          `json:"productSubId"`
	Coin           string          `json:"coin"`
	MinAmount      decimal.Decimal `json:"minAmount"`
	ConversionRate decimal.Decimal `json:"conversionRate"` // Cash+ per unit of coin
	SubscribeScale string          `json:"subscribeScale"`
	ProfitScale    string          `json:"profitScale"`
	SoldOut        string          `json:"soldOut"` // yes, no
}

// CashPlusSubscribeService -- POST /api/v3/earn/cashplus/subscribe (UTA mgt. read & write)
//
// Subscribes to Cash+ with the coin of the given product sub-ID (from
// subscribe-info), returning the created subscription order ID.
type CashPlusSubscribeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCashPlusSubscribeService(productSubID string, amount decimal.Decimal) *CashPlusSubscribeService {
	return &CashPlusSubscribeService{c: c, body: map[string]any{
		"productSubId": productSubID,
		"amount":       amount.String(),
	}}
}

// SetAccountType sets the debit account(s); several may be combined (e.g.
// funding and uta). Defaults to the funding account when omitted.
func (s *CashPlusSubscribeService) SetAccountType(accountTypes ...CashPlusAccount) *CashPlusSubscribeService {
	types := make([]string, len(accountTypes))
	for i, t := range accountTypes {
		types[i] = string(t)
	}
	s.body["accountType"] = strings.Join(types, ",")
	return s
}

// SetIdempotentNo sets an idempotency ID to guard against duplicate submissions.
func (s *CashPlusSubscribeService) SetIdempotentNo(idempotentNo string) *CashPlusSubscribeService {
	s.body["idempotentNo"] = idempotentNo
	return s
}

func (s *CashPlusSubscribeService) Do(ctx context.Context) (*CashPlusSubscribeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/cashplus/subscribe", s.body).WithSign()
	return request.Do[CashPlusSubscribeResult](req)
}

type CashPlusSubscribeResult struct {
	OrderID string `json:"orderId"`
}

// GetCashPlusSubscribeResultService -- GET /api/v3/earn/cashplus/subscribe-result (UTA mgt. read)
//
// Returns the state of a Cash+ subscription order and how it was paid.
type GetCashPlusSubscribeResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCashPlusSubscribeResultService(orderID string) *GetCashPlusSubscribeResultService {
	return &GetCashPlusSubscribeResultService{c: c, params: map[string]string{"orderId": orderID}}
}

func (s *GetCashPlusSubscribeResultService) Do(ctx context.Context) (*CashPlusSubscribeStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/subscribe-result", s.params).WithSign()
	return request.Do[CashPlusSubscribeStatus](req)
}

type CashPlusSubscribeStatus struct {
	OrderID            string            `json:"orderId"`
	Status             CashPlusStatus    `json:"status"` // processing, success, failed
	Coin               string            `json:"coin"`
	Amount             decimal.Decimal   `json:"amount"`
	SubscribeTime      time.Time         `json:"subscribeTime,format:unixmilli"`
	InterestTime       time.Time         `json:"interestTime,format:unixmilli"`
	DividendTime       time.Time         `json:"dividendTime,format:unixmilli"` // first interest distribution
	PaymentAccountList []CashPlusPayment `json:"paymentAccountList"`
}

type CashPlusPayment struct {
	AccountType CashPlusAccount `json:"accountType"`
	PayAmount   decimal.Decimal `json:"payAmount"`
}

// GetCashPlusRedeemInfoService -- GET /api/v3/earn/cashplus/redeem-info (UTA mgt. read)
//
// Returns the Cash+ redemption prerequisites: the coins Cash+ can be redeemed
// into with the redeemable position, and today's remaining redemption quota.
// Callable by both master and sub-accounts.
type GetCashPlusRedeemInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCashPlusRedeemInfoService() *GetCashPlusRedeemInfoService {
	return &GetCashPlusRedeemInfoService{c: c}
}

func (s *GetCashPlusRedeemInfoService) Do(ctx context.Context) (*CashPlusRedeemInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/redeem-info").WithSign()
	return request.Do[CashPlusRedeemInfo](req)
}

type CashPlusRedeemInfo struct {
	RedeemCoinList     []CashPlusRedeemCoin `json:"redeemCoinList"`
	UserRemainingQuota decimal.Decimal      `json:"userRemainingQuota"` // today's remaining quota
}

type CashPlusRedeemCoin struct {
	ProductSubID    string          `json:"productSubId"`
	Coin            string          `json:"coin"`
	ConversionRate  decimal.Decimal `json:"conversionRate"` // coin per unit of Cash+
	RedeemAvailable decimal.Decimal `json:"redeemAvailable"`
	RedeemMinAmount decimal.Decimal `json:"redeemMinAmount"`
	RedeemScale     string          `json:"redeemScale"`
	RedeemTime      time.Time       `json:"redeemTime,format:unixmilli"`
	ReceivedTime    time.Time       `json:"receivedTime,format:unixmilli"` // estimated arrival
}

// CashPlusRedeemService -- POST /api/v3/earn/cashplus/redeem (UTA mgt. read & write)
//
// Redeems the given amount of Cash+ into the coin of the product sub-ID (from
// redeem-info), returning the redemption order ID and estimated arrival time.
type CashPlusRedeemService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCashPlusRedeemService(productSubID string, amount decimal.Decimal) *CashPlusRedeemService {
	return &CashPlusRedeemService{c: c, body: map[string]any{
		"productSubId": productSubID,
		"amount":       amount.String(),
	}}
}

// SetAccountType sets the credit account (funding or uta; defaults to the
// funding account when omitted).
func (s *CashPlusRedeemService) SetAccountType(accountType CashPlusAccount) *CashPlusRedeemService {
	s.body["accountType"] = string(accountType)
	return s
}

func (s *CashPlusRedeemService) Do(ctx context.Context) (*CashPlusRedeemResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/earn/cashplus/redeem", s.body).WithSign()
	return request.Do[CashPlusRedeemResult](req)
}

type CashPlusRedeemResult struct {
	OrderID     string    `json:"orderId"`
	RedeemTime  time.Time `json:"redeemTime,format:unixmilli"`
	ReceiveTime time.Time `json:"receiveTime,format:unixmilli"` // estimated arrival
}

// GetCashPlusRedeemResultService -- GET /api/v3/earn/cashplus/redeem-result (UTA mgt. read)
//
// Returns the state of a Cash+ redemption order.
type GetCashPlusRedeemResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCashPlusRedeemResultService(orderID string) *GetCashPlusRedeemResultService {
	return &GetCashPlusRedeemResultService{c: c, params: map[string]string{"orderId": orderID}}
}

func (s *GetCashPlusRedeemResultService) Do(ctx context.Context) (*CashPlusRedeemStatus, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/redeem-result", s.params).WithSign()
	return request.Do[CashPlusRedeemStatus](req)
}

type CashPlusRedeemStatus struct {
	OrderID     string         `json:"orderId"`
	Status      CashPlusStatus `json:"status"` // processing, success, failed
	RedeemTime  time.Time      `json:"redeemTime,format:unixmilli"`
	ReceiveTime time.Time      `json:"receiveTime,format:unixmilli"` // estimated arrival
}

// GetCashPlusPositionsService -- GET /api/v3/earn/cashplus/positions (UTA mgt. read)
//
// Returns the account's Cash+ holding, cumulative and yesterday's profit (all
// in USD) and the reference annualized rate.
type GetCashPlusPositionsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCashPlusPositionsService() *GetCashPlusPositionsService {
	return &GetCashPlusPositionsService{c: c}
}

func (s *GetCashPlusPositionsService) Do(ctx context.Context) (*CashPlusPositions, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/positions").WithSign()
	return request.Do[CashPlusPositions](req)
}

type CashPlusPositions struct {
	Balance         decimal.Decimal `json:"balance"`         // USD
	TotalProfit     decimal.Decimal `json:"totalProfit"`     // USD
	YesterdayProfit decimal.Decimal `json:"yesterdayProfit"` // USD
	Apr             decimal.Decimal `json:"apr"`             // compound annual rate, for reference
}

// GetCashPlusHistoryOrdersService -- GET /api/v3/earn/cashplus/history-orders (UTA mgt. read)
//
// Returns the account's Cash+ subscribe/redeem/interest records, paginated by
// cursor (nextCursor).
type GetCashPlusHistoryOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCashPlusHistoryOrdersService() *GetCashPlusHistoryOrdersService {
	return &GetCashPlusHistoryOrdersService{c: c, params: map[string]string{}}
}

func (s *GetCashPlusHistoryOrdersService) SetStartTime(t time.Time) *GetCashPlusHistoryOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetCashPlusHistoryOrdersService) SetEndTime(t time.Time) *GetCashPlusHistoryOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetCashPlusHistoryOrdersService) SetRecordType(recordType CashPlusRecordType) *GetCashPlusHistoryOrdersService {
	s.params["recordType"] = string(recordType)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetCashPlusHistoryOrdersService) SetLimit(limit string) *GetCashPlusHistoryOrdersService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (use nextCursor from the previous
// response; omit for the first page).
func (s *GetCashPlusHistoryOrdersService) SetCursor(cursor string) *GetCashPlusHistoryOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCashPlusHistoryOrdersService) Do(ctx context.Context) (*CashPlusHistoryOrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/earn/cashplus/history-orders", s.params).WithSign()
	return request.Do[CashPlusHistoryOrders](req)
}

type CashPlusHistoryOrders struct {
	Items      []CashPlusHistoryOrder `json:"items"`
	NextCursor string                 `json:"nextCursor"`
}

type CashPlusHistoryOrder struct {
	RecordID    string             `json:"recordId"`
	RecordType  CashPlusRecordType `json:"recordType"`
	Status      CashPlusStatus     `json:"status"`
	Amount      decimal.Decimal    `json:"amount"`
	Coin        string             `json:"coin"`
	CreatedTime time.Time          `json:"createdTime,format:unixmilli"`
	AccountType CashPlusAccount    `json:"accountType"` // debit account for subscribe, credit account otherwise
}
