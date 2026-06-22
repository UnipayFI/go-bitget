package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetP2PAdListService -- GET /api/v3/p2p/ad-list (UTA P2P read)
//
// Returns the public marketplace advertisements for a token/fiat pair and trade
// side, optionally filtered by a target fiat amount, paginated by page number.
type GetP2PAdListService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PAdListService(token, fiat string, side Side, pageNum, limit string) *GetP2PAdListService {
	return &GetP2PAdListService{c: c, params: map[string]string{
		"token":   token,
		"fiat":    fiat,
		"side":    string(side),
		"pageNum": pageNum,
		"limit":   limit,
	}}
}

// SetAmount filters ads to those eligible for the given fiat trade amount.
func (s *GetP2PAdListService) SetAmount(amount string) *GetP2PAdListService {
	s.params["amount"] = amount
	return s
}

func (s *GetP2PAdListService) Do(ctx context.Context) ([]P2PAd, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/ad-list", s.params).WithSign()
	resp, err := request.Do[[]P2PAd](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// P2PAd is a single public marketplace advertisement.
type P2PAd struct {
	AdvID             string           `json:"advId"`
	MerchantID        string           `json:"merchantId"`
	MerchantName      string           `json:"merchantName"`
	Token             string           `json:"token"`
	Fiat              string           `json:"fiat"`
	Side              Side             `json:"side"`
	Price             decimal.Decimal  `json:"price"`
	MinAmount         decimal.Decimal  `json:"minAmount"`
	MaxAmount         decimal.Decimal  `json:"maxAmount"`
	Quantity          decimal.Decimal  `json:"quantity"`
	PayMethods        []P2PAdPayMethod `json:"payMethods"`
	CompletedOrderNum string           `json:"completedOrderNum"`
	CompletedRate     decimal.Decimal  `json:"completedRate"`
	AvgReleaseTime    string           `json:"avgReleaseTime"`
	CreatedTime       time.Time        `json:"createdTime"`
}

// P2PAdPayMethod is an accepted payment method on a public advertisement.
type P2PAdPayMethod struct {
	PayMethodID   string `json:"payMethodId"`
	PayMethodName string `json:"payMethodName"`
}

// GetP2PExchangeRateService -- GET /api/v3/p2p/exchange-rate (UTA P2P read)
//
// Returns the reference exchange rate for a token/fiat currency pair.
type GetP2PExchangeRateService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PExchangeRateService(token, fiat string) *GetP2PExchangeRateService {
	return &GetP2PExchangeRateService{c: c, params: map[string]string{
		"token": token,
		"fiat":  fiat,
	}}
}

func (s *GetP2PExchangeRateService) Do(ctx context.Context) (*P2PExchangeRate, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/exchange-rate", s.params).WithSign()
	return request.Do[P2PExchangeRate](req)
}

type P2PExchangeRate struct {
	Rate decimal.Decimal `json:"rate"`
}

// GetP2PAdLimitService -- GET /api/v3/p2p/ad-limit (UTA P2P read)
//
// Returns the minimum and maximum token quantities (in USDT) permitted for buy
// and sell advertisements on a token/fiat pair and side.
type GetP2PAdLimitService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PAdLimitService(token, fiat string, side Side) *GetP2PAdLimitService {
	return &GetP2PAdLimitService{c: c, params: map[string]string{
		"token": token,
		"fiat":  fiat,
		"side":  string(side),
	}}
}

func (s *GetP2PAdLimitService) Do(ctx context.Context) (*P2PAdLimit, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/ad-limit", s.params).WithSign()
	return request.Do[P2PAdLimit](req)
}

type P2PAdLimit struct {
	BuyMinTokenAmount  decimal.Decimal `json:"buyMinTokenAmount"`
	BuyMaxTokenAmount  decimal.Decimal `json:"buyMaxTokenAmount"`
	SellMinTokenAmount decimal.Decimal `json:"sellMinTokenAmount"`
	SellMaxTokenAmount decimal.Decimal `json:"sellMaxTokenAmount"`
}

// FeeSimulateP2PService -- POST /api/v3/p2p/fee-simulate (UTA P2P read & write)
//
// Estimates the fee and tax that would be charged for an advertisement of the
// given size on a token/fiat pair, side and market type ("public"/"private").
type FeeSimulateP2PService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewFeeSimulateP2PService(token, fiat string, side Side, marketType, amount string) *FeeSimulateP2PService {
	return &FeeSimulateP2PService{c: c, body: map[string]any{
		"token":      token,
		"fiat":       fiat,
		"side":       string(side),
		"marketType": marketType,
		"amount":     amount,
	}}
}

func (s *FeeSimulateP2PService) Do(ctx context.Context) (*P2PFeeSimulation, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/fee-simulate", s.body).WithSign()
	return request.Do[P2PFeeSimulation](req)
}

type P2PFeeSimulation struct {
	FeeToken  string          `json:"feeToken"`
	FeeAmount decimal.Decimal `json:"feeAmount"`
	TaxToken  string          `json:"taxToken"`
	TaxAmount decimal.Decimal `json:"taxAmount"`
}

// GetP2PMyAdsService -- GET /api/v3/p2p/my-ads (UTA P2P read)
//
// Returns the caller's own advertisements, paginated by cursor and bounded to a
// 90-day lookback window, with optional ad/token/fiat/status filters.
type GetP2PMyAdsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PMyAdsService() *GetP2PMyAdsService {
	return &GetP2PMyAdsService{c: c, params: map[string]string{}}
}

// SetStartTime filters ads at or after t (90-day lookback, 30-day query range).
func (s *GetP2PMyAdsService) SetStartTime(t time.Time) *GetP2PMyAdsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters ads at or before t.
func (s *GetP2PMyAdsService) SetEndTime(t time.Time) *GetP2PMyAdsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetP2PMyAdsService) SetLimit(limit string) *GetP2PMyAdsService {
	s.params["limit"] = limit
	return s
}

func (s *GetP2PMyAdsService) SetCursor(cursor string) *GetP2PMyAdsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetP2PMyAdsService) SetAdvID(advID string) *GetP2PMyAdsService {
	s.params["advId"] = advID
	return s
}

func (s *GetP2PMyAdsService) SetToken(token string) *GetP2PMyAdsService {
	s.params["token"] = token
	return s
}

func (s *GetP2PMyAdsService) SetFiat(fiat string) *GetP2PMyAdsService {
	s.params["fiat"] = fiat
	return s
}

// SetStatus filters by ad status ("publish", "delist", "remove").
func (s *GetP2PMyAdsService) SetStatus(status string) *GetP2PMyAdsService {
	s.params["status"] = status
	return s
}

func (s *GetP2PMyAdsService) Do(ctx context.Context) (*P2PMyAds, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/my-ads", s.params).WithSign()
	return request.Do[P2PMyAds](req)
}

type P2PMyAds struct {
	Items  []P2PMyAd `json:"items"`
	NextID string    `json:"nextId"`
}

// P2PMyAd is a single advertisement owned by the caller.
type P2PMyAd struct {
	AdvID       string             `json:"advId"`
	Side        Side               `json:"side"`
	PriceType   string             `json:"priceType"`
	Token       string             `json:"token"`
	Fiat        string             `json:"fiat"`
	Price       decimal.Decimal    `json:"price"`
	SoldAmount  decimal.Decimal    `json:"soldAmount"`
	LastAmount  decimal.Decimal    `json:"lastAmount"`
	PayMethods  []P2PMyAdPayMethod `json:"payMethods"`
	Status      string             `json:"status"`
	CreatedTime time.Time          `json:"createdTime"`
	UpdatedTime time.Time          `json:"updatedTime"`
}

// P2PMyAdPayMethod is a payment method reference on the caller's advertisement.
type P2PMyAdPayMethod struct {
	PayMethodID string `json:"payMethodId"`
}

// GetP2PAdInfoService -- GET /api/v3/p2p/ad-info (UTA P2P read)
//
// Returns the full detail of a single advertisement owned by the caller.
type GetP2PAdInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PAdInfoService(advID string) *GetP2PAdInfoService {
	return &GetP2PAdInfoService{c: c, params: map[string]string{"advId": advID}}
}

func (s *GetP2PAdInfoService) Do(ctx context.Context) (*P2PAdInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/ad-info", s.params).WithSign()
	return request.Do[P2PAdInfo](req)
}

type P2PAdInfo struct {
	AdvID          string               `json:"advId"`
	Side           Side                 `json:"side"`
	PriceType      string               `json:"priceType"`
	Token          string               `json:"token"`
	Fiat           string               `json:"fiat"`
	Price          decimal.Decimal      `json:"price"`
	Amount         decimal.Decimal      `json:"amount"`
	MinAmount      decimal.Decimal      `json:"minAmount"`
	MaxAmount      decimal.Decimal      `json:"maxAmount"`
	PayDuration    string               `json:"payDuration"`
	SoldAmount     decimal.Decimal      `json:"soldAmount"`
	LastAmount     decimal.Decimal      `json:"lastAmount"`
	Premium        decimal.Decimal      `json:"premium"`
	Turnover       decimal.Decimal      `json:"turnover"`
	Remark         string               `json:"remark"`
	TradeTerm      string               `json:"tradeTerm"`
	PayMethodIds   []P2PAdInfoPayMethod `json:"payMethodIds"`
	TokenPrecision string               `json:"tokenPrecision"`
	FiatPrecision  string               `json:"fiatPrecision"`
	Market         string               `json:"market"`
	AvgTime        string               `json:"avgTime"`
	CreatedTime    time.Time            `json:"createdTime"`
	UpdatedTime    time.Time            `json:"updatedTime"`
}

// P2PAdInfoPayMethod is a payment method reference on an advertisement detail.
type P2PAdInfoPayMethod struct {
	PayMethodID     string `json:"payMethodId"`
	UserPayMethodID string `json:"userPayMethodId"`
}

// CreateP2PAdService -- POST /api/v3/p2p/ad-create (UTA P2P read & write)
//
// Creates a new advertisement. price is required when priceType is "fixed";
// premium is required when priceType is "floating". payMethodIds entries carry a
// payMethodId and, for sell ads, a userPayMethodId.
type CreateP2PAdService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateP2PAdService(token, fiat string, side Side, priceType, minAmount, maxAmount, quantity string, payMethodIds []P2PPayMethodRef, payTimeLimit string) *CreateP2PAdService {
	return &CreateP2PAdService{c: c, body: map[string]any{
		"token":        token,
		"fiat":         fiat,
		"side":         string(side),
		"priceType":    priceType,
		"minAmount":    minAmount,
		"maxAmount":    maxAmount,
		"quantity":     quantity,
		"payMethodIds": payMethodIds,
		"payTimeLimit": payTimeLimit,
	}}
}

// SetPrice sets the fixed unit price (required when priceType is "fixed").
func (s *CreateP2PAdService) SetPrice(price string) *CreateP2PAdService {
	s.body["price"] = price
	return s
}

// SetPremium sets the premium rate (required when priceType is "floating").
func (s *CreateP2PAdService) SetPremium(premium string) *CreateP2PAdService {
	s.body["premium"] = premium
	return s
}

func (s *CreateP2PAdService) SetRemark(remark string) *CreateP2PAdService {
	s.body["remark"] = remark
	return s
}

func (s *CreateP2PAdService) SetTradeTerms(tradeTerms string) *CreateP2PAdService {
	s.body["tradeTerms"] = tradeTerms
	return s
}

func (s *CreateP2PAdService) Do(ctx context.Context) (*P2PCreatedAd, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/ad-create", s.body).WithSign()
	return request.Do[P2PCreatedAd](req)
}

type P2PCreatedAd struct {
	AdvID string `json:"advId"`
}

// P2PPayMethodRef is a payment-method reference supplied when creating or
// updating an advertisement. UserPayMethodID is required for sell ads.
type P2PPayMethodRef struct {
	PayMethodID     string `json:"payMethodId"`
	UserPayMethodID string `json:"userPayMethodId,omitempty"`
}

// UpdateP2PAdService -- POST /api/v3/p2p/ad-update (UTA P2P read & write)
//
// Updates an existing advertisement. Only advId and payTimeLimit are required;
// the remaining fields are applied when set. The reply data is the literal
// string "success".
type UpdateP2PAdService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewUpdateP2PAdService(advID, payTimeLimit string) *UpdateP2PAdService {
	return &UpdateP2PAdService{c: c, body: map[string]any{
		"advId":        advID,
		"payTimeLimit": payTimeLimit,
	}}
}

func (s *UpdateP2PAdService) SetPriceType(priceType string) *UpdateP2PAdService {
	s.body["priceType"] = priceType
	return s
}

// SetPrice sets the fixed unit price (required when priceType is "fixed").
func (s *UpdateP2PAdService) SetPrice(price string) *UpdateP2PAdService {
	s.body["price"] = price
	return s
}

// SetPremium sets the premium rate (required when priceType is "floating").
func (s *UpdateP2PAdService) SetPremium(premium string) *UpdateP2PAdService {
	s.body["premium"] = premium
	return s
}

func (s *UpdateP2PAdService) SetMinAmount(minAmount string) *UpdateP2PAdService {
	s.body["minAmount"] = minAmount
	return s
}

func (s *UpdateP2PAdService) SetMaxAmount(maxAmount string) *UpdateP2PAdService {
	s.body["maxAmount"] = maxAmount
	return s
}

func (s *UpdateP2PAdService) SetQuantity(quantity string) *UpdateP2PAdService {
	s.body["quantity"] = quantity
	return s
}

func (s *UpdateP2PAdService) SetPayMethodIds(payMethodIds []P2PPayMethodRef) *UpdateP2PAdService {
	s.body["payMethodIds"] = payMethodIds
	return s
}

func (s *UpdateP2PAdService) SetRemark(remark string) *UpdateP2PAdService {
	s.body["remark"] = remark
	return s
}

func (s *UpdateP2PAdService) SetTradeTerms(tradeTerms string) *UpdateP2PAdService {
	s.body["tradeTerms"] = tradeTerms
	return s
}

func (s *UpdateP2PAdService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/ad-update", s.body).WithSign()
	return request.Do[string](req)
}

// OperateP2PAdService -- POST /api/v3/p2p/ad-operate (UTA P2P read & write)
//
// Publishes ("publish"), delists ("delist") or removes ("remove") an
// advertisement. The reply data is the literal string "success".
type OperateP2PAdService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewOperateP2PAdService(advID, operation string) *OperateP2PAdService {
	return &OperateP2PAdService{c: c, body: map[string]any{
		"advId":     advID,
		"operation": operation,
	}}
}

func (s *OperateP2PAdService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/ad-operate", s.body).WithSign()
	return request.Do[string](req)
}

// GetP2PPendingOrdersService -- GET /api/v3/p2p/pending-orders (UTA P2P read)
//
// Returns the caller's in-progress P2P orders, paginated by cursor, with
// optional order/side/time filters.
type GetP2PPendingOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PPendingOrdersService() *GetP2PPendingOrdersService {
	return &GetP2PPendingOrdersService{c: c, params: map[string]string{}}
}

func (s *GetP2PPendingOrdersService) SetOrderID(orderID string) *GetP2PPendingOrdersService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetP2PPendingOrdersService) SetSide(side Side) *GetP2PPendingOrdersService {
	s.params["side"] = string(side)
	return s
}

func (s *GetP2PPendingOrdersService) SetStartTime(t time.Time) *GetP2PPendingOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetP2PPendingOrdersService) SetEndTime(t time.Time) *GetP2PPendingOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetP2PPendingOrdersService) SetCursor(cursor string) *GetP2PPendingOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetP2PPendingOrdersService) SetLimit(limit string) *GetP2PPendingOrdersService {
	s.params["limit"] = limit
	return s
}

func (s *GetP2PPendingOrdersService) Do(ctx context.Context) (*P2POrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/pending-orders", s.params).WithSign()
	return request.Do[P2POrders](req)
}

type P2POrders struct {
	Items  []P2POrder `json:"items"`
	NextID string     `json:"nextId"`
}

// P2POrder is a single P2P order summary.
type P2POrder struct {
	OrderID      string          `json:"orderId"`
	Side         Side            `json:"side"`
	Token        string          `json:"token"`
	Fiat         string          `json:"fiat"`
	Price        decimal.Decimal `json:"price"`
	Amount       decimal.Decimal `json:"amount"`
	Quantity     decimal.Decimal `json:"quantity"`
	Fee          decimal.Decimal `json:"fee"`
	Counterparty string          `json:"counterparty"`
	Status       string          `json:"status"`
	CreatedTime  time.Time       `json:"createdTime"`
	UpdatedTime  time.Time       `json:"updatedTime"`
}

// GetP2PAllOrdersService -- GET /api/v3/p2p/all-orders (UTA P2P read)
//
// Returns the caller's P2P orders across all states, paginated by cursor and
// bounded to a 90-day lookback window, with optional order/side/status filters.
type GetP2PAllOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PAllOrdersService() *GetP2PAllOrdersService {
	return &GetP2PAllOrdersService{c: c, params: map[string]string{}}
}

// SetStartTime filters orders at or after t (90-day lookback window).
func (s *GetP2PAllOrdersService) SetStartTime(t time.Time) *GetP2PAllOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetP2PAllOrdersService) SetEndTime(t time.Time) *GetP2PAllOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetP2PAllOrdersService) SetLimit(limit string) *GetP2PAllOrdersService {
	s.params["limit"] = limit
	return s
}

func (s *GetP2PAllOrdersService) SetCursor(cursor string) *GetP2PAllOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetP2PAllOrdersService) SetOrderID(orderID string) *GetP2PAllOrdersService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetP2PAllOrdersService) SetSide(side Side) *GetP2PAllOrdersService {
	s.params["side"] = string(side)
	return s
}

// SetStatus filters by order status ("pending_payment", "pending_release",
// "completed", "cancelled", "in_appeal").
func (s *GetP2PAllOrdersService) SetStatus(status string) *GetP2PAllOrdersService {
	s.params["status"] = status
	return s
}

func (s *GetP2PAllOrdersService) Do(ctx context.Context) (*P2POrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/all-orders", s.params).WithSign()
	return request.Do[P2POrders](req)
}

// GetP2POrderInfoService -- GET /api/v3/p2p/order-info (UTA P2P read)
//
// Returns the full detail of a single P2P order, including payment method and
// counterparty information.
type GetP2POrderInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2POrderInfoService(orderID string) *GetP2POrderInfoService {
	return &GetP2POrderInfoService{c: c, params: map[string]string{"orderId": orderID}}
}

func (s *GetP2POrderInfoService) Do(ctx context.Context) (*P2POrderInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/order-info", s.params).WithSign()
	return request.Do[P2POrderInfo](req)
}

type P2POrderInfo struct {
	OrderID              string            `json:"orderId"`
	Side                 Side              `json:"side"`
	Token                string            `json:"token"`
	Fiat                 string            `json:"fiat"`
	Price                decimal.Decimal   `json:"price"`
	Quantity             decimal.Decimal   `json:"quantity"`
	Amount               decimal.Decimal   `json:"amount"`
	Status               string            `json:"status"`
	AdvID                string            `json:"advId"`
	Remark               string            `json:"remark"`
	OrderCancelCountdown string            `json:"orderCancelCountdown"`
	Fee                  decimal.Decimal   `json:"fee"`
	VerifyCapital        string            `json:"verifyCapital"`
	PayMethodDetail      P2POrderPayMethod `json:"payMethodDetail"`
	SellUserInfo         P2POrderUserInfo  `json:"sellUserInfo"`
	BuyUserInfo          P2POrderUserInfo  `json:"buyUserInfo"`
	CreatedTime          time.Time         `json:"createdTime"`
	UpdatedTime          time.Time         `json:"updatedTime"`
}

// P2POrderPayMethod is the payment method attached to an order (omitted for
// asset-verification orders).
type P2POrderPayMethod struct {
	PayMethodID       string              `json:"payMethodId"`
	PayMethodName     string              `json:"payMethodName"`
	PayMethodUserName string              `json:"payMethodUserName"`
	Postscript        string              `json:"postscript"`
	PayMethodInfo     []P2PPayMethodField `json:"payMethodInfo"`
}

// P2POrderUserInfo identifies a party to an order.
type P2POrderUserInfo struct {
	NickName string `json:"nickName"`
	RealName string `json:"realName"`
}

// ReleaseP2PAssetService -- POST /api/v3/p2p/order-release (UTA P2P read & write)
//
// Releases the crypto for a P2P sell order after the buyer's payment is
// verified. The reply data is the literal string "success".
type ReleaseP2PAssetService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewReleaseP2PAssetService(orderID string) *ReleaseP2PAssetService {
	return &ReleaseP2PAssetService{c: c, body: map[string]any{"orderId": orderID}}
}

func (s *ReleaseP2PAssetService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/order-release", s.body).WithSign()
	return request.Do[string](req)
}

// ConfirmP2PPaymentService -- POST /api/v3/p2p/order-pay (UTA P2P read & write)
//
// Marks a P2P buy order as paid. The reply data is the literal string
// "success".
type ConfirmP2PPaymentService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewConfirmP2PPaymentService(orderID string) *ConfirmP2PPaymentService {
	return &ConfirmP2PPaymentService{c: c, body: map[string]any{"orderId": orderID}}
}

func (s *ConfirmP2PPaymentService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/p2p/order-pay", s.body).WithSign()
	return request.Do[string](req)
}

// GetP2PUserInfoService -- GET /api/v3/p2p/user-info (UTA P2P read)
//
// Returns the caller's P2P merchant profile, including order statistics and
// advertisement/order limits.
type GetP2PUserInfoService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetP2PUserInfoService() *GetP2PUserInfoService {
	return &GetP2PUserInfoService{c: c}
}

func (s *GetP2PUserInfoService) Do(ctx context.Context) (*P2PUserInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/user-info").WithSign()
	return request.Do[P2PUserInfo](req)
}

type P2PUserInfo struct {
	UID                  string              `json:"uid"`
	NickName             string              `json:"nickName"`
	AccountLevel         string              `json:"accountLevel"`
	CompletedOrderNum    string              `json:"completedOrderNum"`
	PositiveRate         decimal.Decimal     `json:"positiveRate"`
	CompletedOrderNum30D string              `json:"completedOrderNum30D"`
	CompletedRate30D     decimal.Decimal     `json:"completedRate30D"`
	AvgPayTime30D        string              `json:"avgPayTime30D"`
	AvgReleaseTime30D    string              `json:"avgReleaseTime30D"`
	EquityDetail         P2PUserEquityDetail `json:"equityDetail"`
	RegisterTime         time.Time           `json:"registerTime"`
}

// P2PUserEquityDetail is the caller's advertisement and pending-order capacity.
type P2PUserEquityDetail struct {
	MaxAdvSellNum        string          `json:"maxAdvSellNum"`
	MaxAdvSellLimit      decimal.Decimal `json:"maxAdvSellLimit"`
	MaxAdvBuyNum         string          `json:"maxAdvBuyNum"`
	MaxAdvBuyLimit       decimal.Decimal `json:"maxAdvBuyLimit"`
	TotalMaxPendingOrder string          `json:"totalMaxPendingOrder"`
	AdvMaxPendingOrder   string          `json:"advMaxPendingOrder"`
}

// GetP2PCurrenciesService -- GET /api/v3/p2p/currencies (UTA P2P read)
//
// Returns the tokens and fiat currencies supported by the P2P marketplace,
// including per-fiat rate ranges, ad limits and payment methods.
type GetP2PCurrenciesService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetP2PCurrenciesService() *GetP2PCurrenciesService {
	return &GetP2PCurrenciesService{c: c}
}

func (s *GetP2PCurrenciesService) Do(ctx context.Context) (*P2PCurrencies, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/currencies").WithSign()
	return request.Do[P2PCurrencies](req)
}

type P2PCurrencies struct {
	TokenDetailList []P2PTokenDetail `json:"tokenDetailList"`
	FiatDetailList  []P2PFiatDetail  `json:"fiatDetailList"`
}

// P2PTokenDetail is a supported cryptocurrency.
type P2PTokenDetail struct {
	Token     string `json:"token"`
	Precision string `json:"precision"`
}

// P2PFiatDetail is a supported fiat currency with its P2P parameters.
type P2PFiatDetail struct {
	Fiat            string             `json:"fiat"`
	Precision       string             `json:"precision"`
	FloatRateMin    decimal.Decimal    `json:"floatRateMin"`
	FloatRateMax    decimal.Decimal    `json:"floatRateMax"`
	FixedRateMin    decimal.Decimal    `json:"fixedRateMin"`
	FixedRateMax    decimal.Decimal    `json:"fixedRateMax"`
	AdMinLimit      decimal.Decimal    `json:"adMinLimit"`
	AdMaxLimit      decimal.Decimal    `json:"adMaxLimit"`
	PayMethods      []P2PFiatPayMethod `json:"payMethods"`
	WithdrawalLimit P2PWithdrawalLimit `json:"withdrawalLimit"`
	PayTimeLimit    []int              `json:"payTimeLimit"`
}

// P2PFiatPayMethod is a payment method available for a fiat currency.
type P2PFiatPayMethod struct {
	PayMethodName string `json:"payMethodName"`
	IconURL       string `json:"iconUrl"`
	PayMethodID   string `json:"payMethodId"`
}

// P2PWithdrawalLimit is a fiat withdrawal restriction descriptor.
type P2PWithdrawalLimit struct {
	WithdrawLimitKey   string `json:"withdrawLimitKey"`
	WithdrawLimitValue string `json:"withdrawLimitValue"`
}

// GetP2PPayMethodsService -- GET /api/v3/p2p/pay-method (UTA P2P read)
//
// Returns the caller's configured payment methods, used when creating or
// updating advertisements.
type GetP2PPayMethodsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetP2PPayMethodsService() *GetP2PPayMethodsService {
	return &GetP2PPayMethodsService{c: c}
}

func (s *GetP2PPayMethodsService) Do(ctx context.Context) ([]P2PPayMethod, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/pay-method").WithSign()
	resp, err := request.Do[[]P2PPayMethod](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// P2PPayMethod is one of the caller's configured payment methods.
type P2PPayMethod struct {
	PayMethodID       string              `json:"payMethodId"`
	PayMethodName     string              `json:"payMethodName"`
	PayMethodUserName string              `json:"payMethodUserName"`
	UserPayMethodID   string              `json:"userPayMethodId"`
	PayMethodInfo     []P2PPayMethodField `json:"payMethodInfo"`
}

// P2PPayMethodField is a single field of a payment method's configuration.
type P2PPayMethodField struct {
	Name  string `json:"name"`
	Type  string `json:"type"` // file, number, txt
	Value string `json:"value"`
}

// GetP2PBalanceService -- GET /api/v3/p2p/balance (UTA P2P read)
//
// Returns the caller's P2P-tradeable balance for a token.
type GetP2PBalanceService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetP2PBalanceService(token string) *GetP2PBalanceService {
	return &GetP2PBalanceService{c: c, params: map[string]string{"token": token}}
}

func (s *GetP2PBalanceService) Do(ctx context.Context) (*P2PBalance, error) {
	req := request.Get(ctx, s.c, "/api/v3/p2p/balance", s.params).WithSign()
	return request.Do[P2PBalance](req)
}

type P2PBalance struct {
	Token            string          `json:"token"`
	TotalBalance     decimal.Decimal `json:"totalBalance"`
	FrozenBalance    decimal.Decimal `json:"frozenBalance"`
	AvailableBalance decimal.Decimal `json:"availableBalance"`
}
