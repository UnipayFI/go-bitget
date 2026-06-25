package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// FeeDetail is a single fee line in an order's or fill's fee breakdown.
type FeeDetail struct {
	FeeCoin string          `json:"feeCoin"`
	Fee     decimal.Decimal `json:"fee"`
}

// Order is the union of the order-info, unfilled-orders and history-orders
// shapes; fields not present for a given endpoint or category arrive empty.
type Order struct {
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderType     OrderType       `json:"orderType"`
	Side          Side            `json:"side"`
	Price         decimal.Decimal `json:"price"`
	Qty           decimal.Decimal `json:"qty"`
	Amount        decimal.Decimal `json:"amount"`
	CumExecQty    decimal.Decimal `json:"cumExecQty"`
	CumExecValue  decimal.Decimal `json:"cumExecValue"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	OrderStatus   OrderStatus     `json:"orderStatus"`
	PosSide       PosSide         `json:"posSide"`
	HoldMode      HoldMode        `json:"holdMode"`
	TradeSide     string          `json:"tradeSide"`
	DelegateType  string          `json:"delegateType"`
	ReduceOnly    ReduceOnly      `json:"reduceOnly"`
	MarginMode    MarginMode      `json:"marginMode"`
	StpMode       string          `json:"stpMode"`
	TakeProfit    decimal.Decimal `json:"takeProfit"`
	StopLoss      decimal.Decimal `json:"stopLoss"`
	TpTriggerBy   string          `json:"tpTriggerBy"`
	SlTriggerBy   string          `json:"slTriggerBy"`
	TpOrderType   OrderType       `json:"tpOrderType"`
	SlOrderType   OrderType       `json:"slOrderType"`
	TpLimitPrice  decimal.Decimal `json:"tpLimitPrice"`
	SlLimitPrice  decimal.Decimal `json:"slLimitPrice"`
	FeeDetail     []FeeDetail     `json:"feeDetail"`
	CancelReason  string          `json:"cancelReason"`
	ExecType      ExecType        `json:"execType"`
	CreatedTime   time.Time       `json:"createdTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
}

// GetOrderInfoService -- GET /api/v3/trade/order-info (UTA trade read)
//
// Returns the details of a single order, looked up by orderId or clientOid (one
// is required; orderId takes priority when both are present).
type GetOrderInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOrderInfoService() *GetOrderInfoService {
	return &GetOrderInfoService{c: c, params: map[string]string{}}
}

func (s *GetOrderInfoService) SetOrderID(orderID string) *GetOrderInfoService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetOrderInfoService) SetClientOrderID(clientOid string) *GetOrderInfoService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetOrderInfoService) Do(ctx context.Context) (*Order, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/order-info", s.params).WithSign()
	return request.Do[Order](req)
}

// GetOpenOrdersService -- GET /api/v3/trade/unfilled-orders (UTA trade read)
//
// Returns the account's currently open (unfilled / partially filled) orders,
// paginated by cursor.
type GetOpenOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOpenOrdersService() *GetOpenOrdersService {
	return &GetOpenOrdersService{c: c, params: map[string]string{}}
}

func (s *GetOpenOrdersService) SetCategory(category Category) *GetOpenOrdersService {
	s.params["category"] = string(category)
	return s
}

func (s *GetOpenOrdersService) SetSymbol(symbol string) *GetOpenOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters orders at or after t.
func (s *GetOpenOrdersService) SetStartTime(t time.Time) *GetOpenOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetOpenOrdersService) SetEndTime(t time.Time) *GetOpenOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOpenOrdersService) SetLimit(limit int) *GetOpenOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOpenOrdersService) SetCursor(cursor string) *GetOpenOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetOpenOrdersService) Do(ctx context.Context) (*OrderList, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/unfilled-orders", s.params).WithSign()
	return request.Do[OrderList](req)
}

// GetOrderHistoryService -- GET /api/v3/trade/history-orders (UTA trade read)
//
// Returns the account's historical (filled / cancelled) orders for a product
// category, paginated by cursor and bounded to a 90-day access window.
type GetOrderHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOrderHistoryService(category Category) *GetOrderHistoryService {
	return &GetOrderHistoryService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetOrderHistoryService) SetSymbol(symbol string) *GetOrderHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters orders at or after t (90-day access window).
func (s *GetOrderHistoryService) SetStartTime(t time.Time) *GetOrderHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t (max 30-day range from startTime).
func (s *GetOrderHistoryService) SetEndTime(t time.Time) *GetOrderHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOrderHistoryService) SetLimit(limit int) *GetOrderHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrderHistoryService) SetCursor(cursor string) *GetOrderHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetOrderHistoryService) Do(ctx context.Context) (*OrderList, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/history-orders", s.params).WithSign()
	return request.Do[OrderList](req)
}

// OrderList is the cursor-paginated order collection returned by the
// unfilled-orders and history-orders endpoints.
type OrderList struct {
	List   []Order `json:"list"`
	Cursor string  `json:"cursor"`
}

// GetFillHistoryService -- GET /api/v3/trade/fills (UTA trade read)
//
// Returns the account's trade fills for a product category, paginated by cursor
// and bounded to a 90-day access window.
type GetFillHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFillHistoryService() *GetFillHistoryService {
	return &GetFillHistoryService{c: c, params: map[string]string{}}
}

func (s *GetFillHistoryService) SetCategory(category Category) *GetFillHistoryService {
	s.params["category"] = string(category)
	return s
}

func (s *GetFillHistoryService) SetOrderID(orderID string) *GetFillHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetStartTime filters fills at or after t (90-day access window).
func (s *GetFillHistoryService) SetStartTime(t time.Time) *GetFillHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters fills at or before t (max 30-day range from startTime).
func (s *GetFillHistoryService) SetEndTime(t time.Time) *GetFillHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFillHistoryService) SetLimit(limit int) *GetFillHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFillHistoryService) SetCursor(cursor string) *GetFillHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetFillHistoryService) Do(ctx context.Context) (*FillList, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/fills", s.params).WithSign()
	return request.Do[FillList](req)
}

// FillList is the cursor-paginated fill collection returned by the fills
// endpoint.
type FillList struct {
	List   []Fill `json:"list"`
	Cursor string `json:"cursor"`
}

// Fill is a single trade execution.
type Fill struct {
	ExecID        string          `json:"execId"`
	ExecLinkID    string          `json:"execLinkId"`
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderType     OrderType       `json:"orderType"`
	Side          Side            `json:"side"`
	PosSide       PosSide         `json:"posSide"`
	ExecPrice     decimal.Decimal `json:"execPrice"`
	ExecQty       decimal.Decimal `json:"execQty"`
	ExecValue     decimal.Decimal `json:"execValue"`
	TradeScope    TradeScope      `json:"tradeScope"`
	TradeSide     string          `json:"tradeSide"`
	FeeDetail     []FeeDetail     `json:"feeDetail"`
	CreatedTime   time.Time       `json:"createdTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
	ExecPnL       decimal.Decimal `json:"execPnl"`
	IsRPI         string          `json:"isRPI"`
}
