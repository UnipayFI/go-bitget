package p2p

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// P2PSide is the transaction direction of a P2P advertisement or order.
type P2PSide string

const (
	P2PSideBuy  P2PSide = "buy"
	P2PSideSell P2PSide = "sell"
)

// P2POnline is the merchant online-status filter accepted by merchantList.
type P2POnline string

const (
	P2POnlineYes P2POnline = "yes"
	P2POnlineNo  P2POnline = "no"
)

// P2PAdvStatus is the lifecycle state of a P2P advertisement.
type P2PAdvStatus string

const (
	P2PAdvStatusOnline    P2PAdvStatus = "online"
	P2PAdvStatusOffline   P2PAdvStatus = "offline"
	P2PAdvStatusEditing   P2PAdvStatus = "editing"
	P2PAdvStatusCompleted P2PAdvStatus = "completed"
)

// P2POrderStatus is the lifecycle state of a P2P order.
type P2POrderStatus string

const (
	P2POrderStatusPendingPay P2POrderStatus = "pending_pay"
	P2POrderStatusPaid       P2POrderStatus = "Paid"
	P2POrderStatusAppeal     P2POrderStatus = "Appeal"
	P2POrderStatusCompleted  P2POrderStatus = "Completed"
	P2POrderStatusCancelled  P2POrderStatus = "cancelled"
)

// P2PAdvSourceType selects whose advertisements advList returns.
type P2PAdvSourceType string

const (
	P2PAdvSourceTypeOwner              P2PAdvSourceType = "owner"
	P2PAdvSourceTypeCompetitor         P2PAdvSourceType = "competitior"
	P2PAdvSourceTypeOwnerAndCompetitor P2PAdvSourceType = "ownerAndCompetitior"
)

// GetP2PMerchantListService -- GET /api/v2/p2p/merchantList (private)
//
// Returns the paged list of P2P merchants, optionally filtered by online status.
type GetP2PMerchantListService struct {
	c      *P2PClient
	params map[string]string
}

func (c *P2PClient) NewGetP2PMerchantListService() *GetP2PMerchantListService {
	return &GetP2PMerchantListService{c: c, params: map[string]string{}}
}

// SetOnline filters merchants by online status ("yes"/"no").
func (s *GetP2PMerchantListService) SetOnline(online P2POnline) *GetP2PMerchantListService {
	s.params["online"] = string(online)
	return s
}

// SetIDLessThan pages backwards using the minMerchantId of the previous query.
func (s *GetP2PMerchantListService) SetIDLessThan(idLessThan string) *GetP2PMerchantListService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of merchants returned (default 100).
func (s *GetP2PMerchantListService) SetLimit(limit int) *GetP2PMerchantListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetP2PMerchantListService) Do(ctx context.Context) (*P2PMerchantList, error) {
	req := request.Get(ctx, s.c, "/api/v2/p2p/merchantList", s.params).WithSign()
	return request.Do[P2PMerchantList](req)
}

// P2PMerchantList is the paged merchant-list payload.
type P2PMerchantList struct {
	MerchantList  []P2PMerchant `json:"merchantList"`
	MinMerchantID string        `json:"minMerchantId"`
}

// P2PMerchant is a single merchant's public profile and trading statistics.
type P2PMerchant struct {
	RegisterTime        time.Time       `json:"registerTime"`
	NickName            string          `json:"nickName"`
	IsOnline            string          `json:"isOnline"`
	AvgPaymentTime      string          `json:"avgPaymentTime"` // minutes
	AvgReleaseTime      string          `json:"avgReleaseTime"` // minutes
	TotalTrades         string          `json:"totalTrades"`
	TotalBuy            string          `json:"totalBuy"`
	TotalSell           string          `json:"totalSell"`
	TotalCompletionRate decimal.Decimal `json:"totalCompletionRate"`
	Trades30d           string          `json:"trades30d"`
	Sell30d             string          `json:"sell30d"`
	Buy30d              string          `json:"buy30d"`
	CompletionRate30d   decimal.Decimal `json:"completionRate30d"`
}

// GetMerchantInfoService -- GET /api/v2/p2p/merchantInfo (private)
//
// Returns the authenticated merchant's own profile, statistics and KYC/contact
// verification status.
type GetMerchantInfoService struct {
	c *P2PClient
}

func (c *P2PClient) NewGetMerchantInfoService() *GetMerchantInfoService {
	return &GetMerchantInfoService{c: c}
}

func (s *GetMerchantInfoService) Do(ctx context.Context) (*MerchantInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/p2p/merchantInfo").WithSign()
	return request.Do[MerchantInfo](req)
}

// MerchantInfo is the authenticated merchant's own profile.
type MerchantInfo struct {
	RegisterTime        time.Time       `json:"registerTime"`
	NickName            string          `json:"nickName"`
	MerchantID          string          `json:"merchantId"`
	AvgPaymentTime      string          `json:"avgPaymentTime"` // minutes
	AvgReleaseTime      string          `json:"avgReleaseTime"` // minutes
	TotalTrades         string          `json:"totalTrades"`
	TotalBuy            string          `json:"totalBuy"`
	TotalSell           string          `json:"totalSell"`
	TotalCompletionRate decimal.Decimal `json:"totalCompletionRate"`
	Trades30d           string          `json:"trades30d"`
	Sell30d             string          `json:"sell30d"`
	Buy30d              string          `json:"buy30d"`
	CompletionRate30d   decimal.Decimal `json:"completionRate30d"`
	KycStatus           bool            `json:"kycStatus"`
	EmailBindStatus     bool            `json:"emailBindStatus"`
	MobileBindStatus    bool            `json:"mobileBindStatus"`
	Email               string          `json:"email"`
	Mobile              string          `json:"mobile"`
}

// GetP2PAdvListService -- GET /api/v2/p2p/advList (private)
//
// Returns the paged list of P2P advertisements within a time window (max 90-day
// range; startTime is required).
type GetP2PAdvListService struct {
	c      *P2PClient
	params map[string]string
}

func (c *P2PClient) NewGetP2PAdvListService(startTime time.Time) *GetP2PAdvListService {
	return &GetP2PAdvListService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

// SetEndTime caps the time window (max 90 days from startTime).
func (s *GetP2PAdvListService) SetEndTime(t time.Time) *GetP2PAdvListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backwards using the minAdvId of the previous query.
func (s *GetP2PAdvListService) SetIDLessThan(idLessThan string) *GetP2PAdvListService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of advertisements returned (default 100).
func (s *GetP2PAdvListService) SetLimit(limit int) *GetP2PAdvListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetStatus filters by advertisement status.
func (s *GetP2PAdvListService) SetStatus(status P2PAdvStatus) *GetP2PAdvListService {
	s.params["status"] = string(status)
	return s
}

// SetAdvNo filters to a single advertisement order number.
func (s *GetP2PAdvListService) SetAdvNo(advNo string) *GetP2PAdvListService {
	s.params["advNo"] = advNo
	return s
}

// SetSide filters by transaction type (buy/sell).
func (s *GetP2PAdvListService) SetSide(side P2PSide) *GetP2PAdvListService {
	s.params["side"] = string(side)
	return s
}

// SetCoin filters by digital currency (e.g. USDT).
func (s *GetP2PAdvListService) SetCoin(coin string) *GetP2PAdvListService {
	s.params["coin"] = coin
	return s
}

// SetLanguage filters by language ("zh-CN"/"en-US").
func (s *GetP2PAdvListService) SetLanguage(language string) *GetP2PAdvListService {
	s.params["language"] = language
	return s
}

// SetFiat filters by fiat currency (e.g. USD).
func (s *GetP2PAdvListService) SetFiat(fiat string) *GetP2PAdvListService {
	s.params["fiat"] = fiat
	return s
}

// SetOrderBy sorts results by "createTime" or "price" (descending by default).
func (s *GetP2PAdvListService) SetOrderBy(orderBy string) *GetP2PAdvListService {
	s.params["orderBy"] = orderBy
	return s
}

// SetPayMethodID filters by payment method ID.
func (s *GetP2PAdvListService) SetPayMethodID(payMethodID string) *GetP2PAdvListService {
	s.params["payMethodId"] = payMethodID
	return s
}

// SetSourceType selects whose advertisements to return (owner/competitior/ownerAndCompetitior).
func (s *GetP2PAdvListService) SetSourceType(sourceType P2PAdvSourceType) *GetP2PAdvListService {
	s.params["sourceType"] = string(sourceType)
	return s
}

func (s *GetP2PAdvListService) Do(ctx context.Context) (*P2PAdvList, error) {
	req := request.Get(ctx, s.c, "/api/v2/p2p/advList", s.params).WithSign()
	return request.Do[P2PAdvList](req)
}

// P2PAdvList is the paged advertisement-list payload. The advertisements arrive
// under the "merchantList" key.
type P2PAdvList struct {
	MerchantList []P2PAdv `json:"merchantList"`
	MinAdvID     string   `json:"minAdvId"`
}

// P2PAdv is a single P2P advertisement.
type P2PAdv struct {
	AdvID                 string                 `json:"advId"`
	AdvNo                 string                 `json:"advNo"`
	Side                  P2PSide                `json:"side"`
	AdvSize               decimal.Decimal        `json:"advSize"`
	Size                  decimal.Decimal        `json:"size"`
	Coin                  string                 `json:"coin"`
	Price                 decimal.Decimal        `json:"price"`
	CoinPrecision         string                 `json:"coinPrecision"`
	Fiat                  string                 `json:"fiat"`
	FiatPrecision         string                 `json:"fiatPrecision"`
	FiatSymbol            string                 `json:"fiatSymbol"`
	Status                P2PAdvStatus           `json:"status"`
	Hide                  string                 `json:"hide"` // yes, no
	MaxTradeAmount        decimal.Decimal        `json:"maxTradeAmount"`
	MinTradeAmount        decimal.Decimal        `json:"minTradeAmount"`
	PayDuration           string                 `json:"payDuration"` // minutes
	TurnoverNum           decimal.Decimal        `json:"turnoverNum"`
	TurnoverRate          decimal.Decimal        `json:"turnoverRate"`
	Label                 string                 `json:"label"`
	Ctime                 time.Time              `json:"ctime"`
	Utime                 time.Time              `json:"utime"`
	RegisterTime          time.Time              `json:"registerTime"`
	UserLimitList         []P2PAdvUserLimit      `json:"userLimitList"`
	PaymentMethod         []P2PAdvPaymentMethod  `json:"paymentMethod"`
	MerchantCertifiedList []P2PMerchantCertified `json:"merchantCertifiedList"`
}

// P2PAdvUserLimit is a single counterparty-restriction rule on an advertisement.
type P2PAdvUserLimit struct {
	MinCompleteNum     string          `json:"minCompleteNum"`
	MaxCompleteNum     string          `json:"maxCompleteNum"`
	PlaceOrderNum      string          `json:"placeOrderNum"`
	AllowMerchantPlace string          `json:"allowMerchantPlace"` // yes, no
	ThirtyCompleteRate decimal.Decimal `json:"thirtyCompleteRate"`
	Country            string          `json:"country"`
}

// P2PAdvPaymentMethod is a payment method offered on an advertisement.
type P2PAdvPaymentMethod struct {
	PaymentMethod string                  `json:"paymentMethod"`
	PaymentID     string                  `json:"paymentId"`
	PaymentInfo   []P2PAdvPaymentInfoItem `json:"paymentInfo"`
}

// P2PAdvPaymentInfoItem is a single field describing a payment method.
type P2PAdvPaymentInfoItem struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
	Type     string `json:"type"` // number, file
}

// P2PMerchantCertified is a single merchant certification badge.
type P2PMerchantCertified struct {
	ImageUrl string `json:"imageUrl"`
	Desc     string `json:"desc"`
}

// GetP2POrderListService -- GET /api/v2/p2p/orderList (private)
//
// Returns the paged list of the merchant's P2P orders, optionally filtered by a
// time window (max 90-day range), status, side and other criteria.
type GetP2POrderListService struct {
	c      *P2PClient
	params map[string]string
}

func (c *P2PClient) NewGetP2POrderListService() *GetP2POrderListService {
	return &GetP2POrderListService{c: c, params: map[string]string{}}
}

// SetStartTime filters orders at or after t.
func (s *GetP2POrderListService) SetStartTime(t time.Time) *GetP2POrderListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t (max 90 days from startTime).
func (s *GetP2POrderListService) SetEndTime(t time.Time) *GetP2POrderListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backwards using the minOrderId of the previous query.
func (s *GetP2POrderListService) SetIDLessThan(idLessThan string) *GetP2POrderListService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of orders returned (default 100).
func (s *GetP2POrderListService) SetLimit(limit int) *GetP2POrderListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetStatus filters by order status.
func (s *GetP2POrderListService) SetStatus(status P2POrderStatus) *GetP2POrderListService {
	s.params["status"] = string(status)
	return s
}

// SetAdvNo filters to a single advertisement order number.
func (s *GetP2POrderListService) SetAdvNo(advNo string) *GetP2POrderListService {
	s.params["advNo"] = advNo
	return s
}

// SetSide filters by transaction type (buy/sell).
func (s *GetP2POrderListService) SetSide(side P2PSide) *GetP2POrderListService {
	s.params["side"] = string(side)
	return s
}

// SetCoin filters by digital currency (e.g. USDT).
func (s *GetP2POrderListService) SetCoin(coin string) *GetP2POrderListService {
	s.params["coin"] = coin
	return s
}

// SetLanguage filters by language ("zh-CN"/"en-US").
func (s *GetP2POrderListService) SetLanguage(language string) *GetP2POrderListService {
	s.params["language"] = language
	return s
}

// SetFiat filters by fiat currency (e.g. USD).
func (s *GetP2POrderListService) SetFiat(fiat string) *GetP2POrderListService {
	s.params["fiat"] = fiat
	return s
}

// SetOrderNo filters to a single order number.
func (s *GetP2POrderListService) SetOrderNo(orderNo string) *GetP2POrderListService {
	s.params["orderNo"] = orderNo
	return s
}

func (s *GetP2POrderListService) Do(ctx context.Context) (*P2POrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/p2p/orderList", s.params).WithSign()
	return request.Do[P2POrderList](req)
}

// P2POrderList is the paged order-list payload.
type P2POrderList struct {
	OrderList  []P2POrder `json:"orderList"`
	MinOrderID string     `json:"minOrderId"`
}

// P2POrder is a single P2P order.
type P2POrder struct {
	OrderID        string              `json:"orderId"`
	OrderNo        string              `json:"orderNo"`
	AdvNo          string              `json:"advNo"`
	Price          decimal.Decimal     `json:"price"`
	Count          decimal.Decimal     `json:"count"`
	Side           P2PSide             `json:"side"`
	Fiat           string              `json:"fiat"`
	Coin           string              `json:"coin"`
	WithdrawTime   time.Time           `json:"withdrawTime"`
	RepresentTime  time.Time           `json:"representTime"`
	PaymentTime    time.Time           `json:"paymentTime"`
	ReleaseTime    time.Time           `json:"releaseTime"`
	Amount         decimal.Decimal     `json:"amount"`
	BuyerRealName  string              `json:"buyerRealName"`
	SellerRealName string              `json:"sellerRealName"`
	Status         P2POrderStatus      `json:"status"`
	Ctime          time.Time           `json:"ctime"`
	Utime          time.Time           `json:"utime"`
	PaymentInfo    P2POrderPaymentInfo `json:"paymentInfo"`
}

// P2POrderPaymentInfo is the payment method attached to an order.
type P2POrderPaymentInfo struct {
	PaymethodName string                  `json:"paymethodName"`
	PaymethodID   string                  `json:"paymethodId"`
	PaymethodInfo []P2POrderPaymethodItem `json:"paymethodInfo"`
}

// P2POrderPaymethodItem is a single field of an order's payment method.
type P2POrderPaymethodItem struct {
	Name     string `json:"name"`
	Required string `json:"required"` // yes, no
	Type     string `json:"type"`     // number, file
	Value    string `json:"value"`
}
