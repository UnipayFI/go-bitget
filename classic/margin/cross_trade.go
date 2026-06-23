package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceCrossOrderService -- POST /api/v2/margin/crossed/place-order (cross-margin trade)
//
// Places a single cross-margin order.
type PlaceCrossOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewPlaceCrossOrderService(symbol string, orderType OrderType, loanType LoanType, side Side) *PlaceCrossOrderService {
	return &PlaceCrossOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"orderType": string(orderType),
		"loanType":  string(loanType),
		"side":      string(side),
	}}
}

// SetForce sets the time-in-force (required for limit orders, N/A for market).
func (s *PlaceCrossOrderService) SetForce(force Force) *PlaceCrossOrderService {
	s.body["force"] = string(force)
	return s
}

// SetPrice sets the order price (limit orders).
func (s *PlaceCrossOrderService) SetPrice(price decimal.Decimal) *PlaceCrossOrderService {
	s.body["price"] = price.String()
	return s
}

// SetBaseSize sets the size in base currency (required for limit/market sells).
func (s *PlaceCrossOrderService) SetBaseSize(size decimal.Decimal) *PlaceCrossOrderService {
	s.body["baseSize"] = size.String()
	return s
}

// SetQuoteSize sets the size in quote currency (required for market buys).
func (s *PlaceCrossOrderService) SetQuoteSize(size decimal.Decimal) *PlaceCrossOrderService {
	s.body["quoteSize"] = size.String()
	return s
}

// SetClientOid sets a custom order ID (6-hour idempotency window).
func (s *PlaceCrossOrderService) SetClientOid(clientOid string) *PlaceCrossOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetStpMode sets the self-trade prevention mode.
func (s *PlaceCrossOrderService) SetStpMode(mode SelfTradePreventionMode) *PlaceCrossOrderService {
	s.body["stpMode"] = string(mode)
	return s
}

func (s *PlaceCrossOrderService) Do(ctx context.Context) (*CrossPlaceOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/place-order", s.body).WithSign()
	return request.Do[CrossPlaceOrderResult](req)
}

// CrossPlaceOrderResult is the acknowledgement returned by a cross-margin order
// placement.
type CrossPlaceOrderResult struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CrossBatchOrderItem is one order in a cross-margin batch placement request.
type CrossBatchOrderItem struct {
	Side      Side                    `json:"side"`
	OrderType OrderType               `json:"orderType"`
	Force     Force                   `json:"force,omitempty"`
	Price     decimal.Decimal         `json:"price,omitzero"`
	BaseSize  decimal.Decimal         `json:"baseSize,omitzero"`
	QuoteSize decimal.Decimal         `json:"quoteSize,omitzero"`
	LoanType  LoanType                `json:"loanType"`
	ClientOid string                  `json:"clientOid,omitempty"`
	StpMode   SelfTradePreventionMode `json:"stpMode,omitempty"`
}

// PlaceCrossBatchOrderService -- POST /api/v2/margin/crossed/batch-place-order (cross-margin trade)
//
// Places up to 50 cross-margin orders for one symbol in a single request.
type PlaceCrossBatchOrderService struct {
	c         *MarginClient
	symbol    string
	orderList []CrossBatchOrderItem
}

func (c *MarginClient) NewPlaceCrossBatchOrderService(symbol string, orderList []CrossBatchOrderItem) *PlaceCrossBatchOrderService {
	return &PlaceCrossBatchOrderService{c: c, symbol: symbol, orderList: orderList}
}

func (s *PlaceCrossBatchOrderService) Do(ctx context.Context) (*CrossBatchOrderResult, error) {
	body := map[string]any{
		"symbol":    s.symbol,
		"orderList": s.orderList,
	}
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/batch-place-order").SetBody(body).WithSign()
	return request.Do[CrossBatchOrderResult](req)
}

// CrossBatchOrderResult splits a batch placement into the orders that succeeded
// and those that failed.
type CrossBatchOrderResult struct {
	SuccessList []CrossPlaceOrderResult  `json:"successList"`
	FailureList []CrossBatchOrderFailure `json:"failureList"`
}

// CrossBatchOrderFailure describes one order that could not be placed.
type CrossBatchOrderFailure struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode string `json:"errorCode"`
}

// CancelCrossOrderService -- POST /api/v2/margin/crossed/cancel-order (cross-margin trade)
//
// Cancels a single cross-margin order by orderId or clientOid.
type CancelCrossOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCancelCrossOrderService(symbol string) *CancelCrossOrderService {
	return &CancelCrossOrderService{c: c, body: map[string]any{"symbol": symbol}}
}

// SetOrderID identifies the order to cancel by exchange order ID.
func (s *CancelCrossOrderService) SetOrderID(orderID string) *CancelCrossOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the order to cancel by client order ID.
func (s *CancelCrossOrderService) SetClientOid(clientOid string) *CancelCrossOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelCrossOrderService) Do(ctx context.Context) (*CrossCancelOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/cancel-order", s.body).WithSign()
	return request.Do[CrossCancelOrderResult](req)
}

// CrossCancelOrderResult is the acknowledgement returned by a cross-margin
// cancellation.
type CrossCancelOrderResult struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CrossBatchCancelItem identifies one order to cancel in a batch request, by
// orderId and/or clientOid.
type CrossBatchCancelItem struct {
	OrderID   string `json:"orderId,omitempty"`
	ClientOid string `json:"clientOid,omitempty"`
}

// CancelCrossBatchOrderService -- POST /api/v2/margin/crossed/batch-cancel-order (cross-margin trade)
//
// Cancels multiple cross-margin orders for one symbol in a single request.
type CancelCrossBatchOrderService struct {
	c           *MarginClient
	symbol      string
	orderIDList []CrossBatchCancelItem
}

func (c *MarginClient) NewCancelCrossBatchOrderService(symbol string, orderIDList []CrossBatchCancelItem) *CancelCrossBatchOrderService {
	return &CancelCrossBatchOrderService{c: c, symbol: symbol, orderIDList: orderIDList}
}

func (s *CancelCrossBatchOrderService) Do(ctx context.Context) (*CrossBatchCancelResult, error) {
	body := map[string]any{
		"symbol":      s.symbol,
		"orderIdList": s.orderIDList,
	}
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/batch-cancel-order").SetBody(body).WithSign()
	return request.Do[CrossBatchCancelResult](req)
}

// CrossBatchCancelResult splits a batch cancellation into the orders that were
// cancelled and those that failed.
type CrossBatchCancelResult struct {
	SuccessList []CrossCancelOrderResult  `json:"successList"`
	FailureList []CrossBatchCancelFailure `json:"failureList"`
}

// CrossBatchCancelFailure describes one order that could not be cancelled.
type CrossBatchCancelFailure struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode string `json:"errorCode"`
}

// GetCrossOpenOrdersService -- GET /api/v2/margin/crossed/open-orders (cross-margin trade)
//
// Returns the current (unfilled) cross-margin orders for a symbol.
type GetCrossOpenOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossOpenOrdersService(symbol string) *GetCrossOpenOrdersService {
	return &GetCrossOpenOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetOrderID filters to a single order by exchange order ID.
func (s *GetCrossOpenOrdersService) SetOrderID(orderID string) *GetCrossOpenOrdersService {
	s.params["orderId"] = orderID
	return s
}

// SetClientOid filters to a single order by client order ID.
func (s *GetCrossOpenOrdersService) SetClientOid(clientOid string) *GetCrossOpenOrdersService {
	s.params["clientOid"] = clientOid
	return s
}

// SetStartTime sets the start of the query window.
func (s *GetCrossOpenOrdersService) SetStartTime(t time.Time) *GetCrossOpenOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the end of the query window.
func (s *GetCrossOpenOrdersService) SetEndTime(t time.Time) *GetCrossOpenOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders returned (default 100).
func (s *GetCrossOpenOrdersService) SetLimit(limit int) *GetCrossOpenOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards using the last orderId from a prior query.
func (s *GetCrossOpenOrdersService) SetIDLessThan(id string) *GetCrossOpenOrdersService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetCrossOpenOrdersService) Do(ctx context.Context) (*CrossOpenOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/open-orders", s.params).WithSign()
	return request.Do[CrossOpenOrders](req)
}

// CrossOpenOrders is the paged list of current cross-margin orders.
type CrossOpenOrders struct {
	OrderList []CrossOrder `json:"orderList"`
	MaxID     string       `json:"maxId"`
	MinID     string       `json:"minId"`
}

// CrossOrder is a single cross-margin order, as returned by the open-orders and
// history-orders endpoints.
type CrossOrder struct {
	OrderID          string           `json:"orderId"`
	Symbol           string           `json:"symbol"`
	OrderType        OrderType        `json:"orderType"`
	EnterPointSource EnterPointSource `json:"enterPointSource"`
	ClientOid        string           `json:"clientOid"`
	LoanType         LoanType         `json:"loanType"`
	Price            decimal.Decimal  `json:"price"`
	Side             Side             `json:"side"`
	Status           OrderStatus      `json:"status"`
	BaseSize         decimal.Decimal  `json:"baseSize"`
	QuoteSize        decimal.Decimal  `json:"quoteSize"`
	PriceAvg         decimal.Decimal  `json:"priceAvg"`
	Size             decimal.Decimal  `json:"size"`
	Amount           decimal.Decimal  `json:"amount"`
	Force            Force            `json:"force"`
	CTime            time.Time        `json:"cTime"`
	UTime            time.Time        `json:"uTime"`
}

// GetCrossHistoryOrdersService -- GET /api/v2/margin/crossed/history-orders (cross-margin trade)
//
// Returns historical cross-margin orders for a symbol (max 90-day window).
type GetCrossHistoryOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossHistoryOrdersService(symbol string) *GetCrossHistoryOrdersService {
	return &GetCrossHistoryOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetOrderID filters to a single order by exchange order ID.
func (s *GetCrossHistoryOrdersService) SetOrderID(orderID string) *GetCrossHistoryOrdersService {
	s.params["orderId"] = orderID
	return s
}

// SetEnterPointSource filters by the channel the order was created through.
func (s *GetCrossHistoryOrdersService) SetEnterPointSource(src EnterPointSource) *GetCrossHistoryOrdersService {
	s.params["enterPointSource"] = string(src)
	return s
}

// SetClientOid filters to a single order by client order ID.
func (s *GetCrossHistoryOrdersService) SetClientOid(clientOid string) *GetCrossHistoryOrdersService {
	s.params["clientOid"] = clientOid
	return s
}

// SetStartTime sets the start of the query window.
func (s *GetCrossHistoryOrdersService) SetStartTime(t time.Time) *GetCrossHistoryOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the end of the query window.
func (s *GetCrossHistoryOrdersService) SetEndTime(t time.Time) *GetCrossHistoryOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders returned (default 100, max 500).
func (s *GetCrossHistoryOrdersService) SetLimit(limit int) *GetCrossHistoryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards using the last orderId from a prior query.
func (s *GetCrossHistoryOrdersService) SetIDLessThan(id string) *GetCrossHistoryOrdersService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetCrossHistoryOrdersService) Do(ctx context.Context) (*CrossHistoryOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/history-orders", s.params).WithSign()
	return request.Do[CrossHistoryOrders](req)
}

// CrossHistoryOrders is the paged list of historical cross-margin orders.
type CrossHistoryOrders struct {
	OrderList []CrossOrder `json:"orderList"`
	MaxID     string       `json:"maxId"`
	MinID     string       `json:"minId"`
}

// GetCrossFillsService -- GET /api/v2/margin/crossed/fills (cross-margin trade)
//
// Returns cross-margin transaction (fill) details for a symbol (max 90-day
// window).
type GetCrossFillsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossFillsService(symbol string) *GetCrossFillsService {
	return &GetCrossFillsService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetOrderID filters fills to a single order.
func (s *GetCrossFillsService) SetOrderID(orderID string) *GetCrossFillsService {
	s.params["orderId"] = orderID
	return s
}

// SetStartTime sets the start of the query window.
func (s *GetCrossFillsService) SetStartTime(t time.Time) *GetCrossFillsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the end of the query window.
func (s *GetCrossFillsService) SetEndTime(t time.Time) *GetCrossFillsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of fills returned (default 100, max 500).
func (s *GetCrossFillsService) SetLimit(limit int) *GetCrossFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards using the last tradeId from a prior query.
func (s *GetCrossFillsService) SetIDLessThan(id string) *GetCrossFillsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetCrossFillsService) Do(ctx context.Context) (*CrossFills, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/fills", s.params).WithSign()
	return request.Do[CrossFills](req)
}

// CrossFills is the paged list of cross-margin fills.
type CrossFills struct {
	Fills []CrossFill `json:"fills"`
	MaxID string      `json:"maxId"`
	MinID string      `json:"minId"`
}

// CrossFill is a single cross-margin transaction (fill).
type CrossFill struct {
	OrderID    string             `json:"orderId"`
	TradeID    string             `json:"tradeId"`
	Symbol     string             `json:"symbol"`
	OrderType  OrderType          `json:"orderType"`
	Side       Side               `json:"side"`
	PriceAvg   decimal.Decimal    `json:"priceAvg"`
	Size       decimal.Decimal    `json:"size"`
	Amount     decimal.Decimal    `json:"amount"`
	TradeScope string             `json:"tradeScope"` // taker, maker
	FeeDetail  CrossFillFeeDetail `json:"feeDetail"`
	CTime      time.Time          `json:"cTime"`
	UTime      time.Time          `json:"uTime"`
}

// CrossFillFeeDetail breaks down the fees charged on a cross-margin fill.
type CrossFillFeeDetail struct {
	Deduction         string          `json:"deduction"`
	FeeCoin           string          `json:"feeCoin"`
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"`
	TotalFee          decimal.Decimal `json:"totalFee"`
}

// GetCrossLiquidationOrdersService -- GET /api/v2/margin/crossed/liquidation-order (cross-margin trade)
//
// Returns cross-margin liquidation orders (forced liquidation and forced swaps).
type GetCrossLiquidationOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossLiquidationOrdersService() *GetCrossLiquidationOrdersService {
	return &GetCrossLiquidationOrdersService{c: c, params: map[string]string{}}
}

// SetType selects the liquidation record type ("place_order" default, or "swap").
func (s *GetCrossLiquidationOrdersService) SetType(t string) *GetCrossLiquidationOrdersService {
	s.params["type"] = t
	return s
}

// SetSymbol filters by trading pair (applies only when type=place_order).
func (s *GetCrossLiquidationOrdersService) SetSymbol(symbol string) *GetCrossLiquidationOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetFromCoin filters by source currency (applies only when type=swap).
func (s *GetCrossLiquidationOrdersService) SetFromCoin(coin string) *GetCrossLiquidationOrdersService {
	s.params["fromCoin"] = coin
	return s
}

// SetToCoin filters by target currency (applies only when type=swap).
func (s *GetCrossLiquidationOrdersService) SetToCoin(coin string) *GetCrossLiquidationOrdersService {
	s.params["toCoin"] = coin
	return s
}

// SetStartTime sets the start of the query window.
func (s *GetCrossLiquidationOrdersService) SetStartTime(t time.Time) *GetCrossLiquidationOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the end of the query window.
func (s *GetCrossLiquidationOrdersService) SetEndTime(t time.Time) *GetCrossLiquidationOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 100, max 500).
func (s *GetCrossLiquidationOrdersService) SetLimit(limit int) *GetCrossLiquidationOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards using the last idLessThan from a prior query.
func (s *GetCrossLiquidationOrdersService) SetIDLessThan(id string) *GetCrossLiquidationOrdersService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetCrossLiquidationOrdersService) Do(ctx context.Context) (*CrossLiquidationOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/liquidation-order", s.params).WithSign()
	return request.Do[CrossLiquidationOrders](req)
}

// CrossLiquidationOrders is the paged list of cross-margin liquidation orders.
type CrossLiquidationOrders struct {
	ResultList []CrossLiquidationOrder `json:"resultList"`
	IDLessThan string                  `json:"idLessThan"`
	MaxID      string                  `json:"maxId"`
	MinID      string                  `json:"minId"`
}

// CrossLiquidationOrder is a single cross-margin liquidation record (a forced
// order or a forced coin swap).
type CrossLiquidationOrder struct {
	OrderID   string          `json:"orderId"`
	Symbol    string          `json:"symbol"`
	OrderType OrderType       `json:"orderType"`
	Side      Side            `json:"side"`
	PriceAvg  decimal.Decimal `json:"priceAvg"`
	Price     decimal.Decimal `json:"price"`
	FillSize  decimal.Decimal `json:"fillSize"`
	Size      decimal.Decimal `json:"size"`
	Amount    decimal.Decimal `json:"amount"`
	FromCoin  string          `json:"fromCoin"`
	FromSize  decimal.Decimal `json:"fromSize"`
	ToCoin    string          `json:"toCoin"`
	ToSize    decimal.Decimal `json:"toSize"`
	CTime     time.Time       `json:"cTime"`
	UTime     time.Time       `json:"uTime"`
}
