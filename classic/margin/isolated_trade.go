package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceIsolatedOrderService -- POST /api/v2/margin/isolated/place-order (margin trade)
//
// Places a single isolated-margin order.
type PlaceIsolatedOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewPlaceIsolatedOrderService(symbol string, orderType OrderType, loanType LoanType, force Force, side Side) *PlaceIsolatedOrderService {
	return &PlaceIsolatedOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"orderType": string(orderType),
		"loanType":  string(loanType),
		"force":     string(force),
		"side":      string(side),
	}}
}

func (s *PlaceIsolatedOrderService) SetPrice(price decimal.Decimal) *PlaceIsolatedOrderService {
	s.body["price"] = price.String()
	return s
}

func (s *PlaceIsolatedOrderService) SetBaseSize(baseSize decimal.Decimal) *PlaceIsolatedOrderService {
	s.body["baseSize"] = baseSize.String()
	return s
}

func (s *PlaceIsolatedOrderService) SetQuoteSize(quoteSize decimal.Decimal) *PlaceIsolatedOrderService {
	s.body["quoteSize"] = quoteSize.String()
	return s
}

func (s *PlaceIsolatedOrderService) SetClientOid(clientOid string) *PlaceIsolatedOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *PlaceIsolatedOrderService) SetStpMode(stpMode SelfTradePreventionMode) *PlaceIsolatedOrderService {
	s.body["stpMode"] = string(stpMode)
	return s
}

func (s *PlaceIsolatedOrderService) Do(ctx context.Context) (*IsolatedPlaceOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/place-order", s.body).WithSign()
	return request.Do[IsolatedPlaceOrderResult](req)
}

// IsolatedPlaceOrderResult is the result of placing a single isolated order.
type IsolatedPlaceOrderResult struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// BatchPlaceIsolatedOrderService -- POST /api/v2/margin/isolated/batch-place-order (margin trade)
//
// Places up to 50 isolated-margin orders for a single symbol in one request.
type BatchPlaceIsolatedOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewBatchPlaceIsolatedOrderService(symbol string, orderList []IsolatedBatchOrderItem) *BatchPlaceIsolatedOrderService {
	return &BatchPlaceIsolatedOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"orderList": orderList,
	}}
}

func (s *BatchPlaceIsolatedOrderService) Do(ctx context.Context) (*IsolatedBatchOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/batch-place-order", s.body).WithSign()
	return request.Do[IsolatedBatchOrderResult](req)
}

// IsolatedBatchOrderItem is a single order within a batch-place request. Optional
// decimal fields use ,omitzero so a zero value is dropped (a serialized "0"
// would break, e.g., market orders that omit price).
type IsolatedBatchOrderItem struct {
	Side      Side                    `json:"side"`
	OrderType OrderType               `json:"orderType"`
	Force     Force                   `json:"force"`
	LoanType  LoanType                `json:"loanType"`
	Price     decimal.Decimal         `json:"price,omitzero"`
	BaseSize  decimal.Decimal         `json:"baseSize,omitzero"`
	QuoteSize decimal.Decimal         `json:"quoteSize,omitzero"`
	ClientOid string                  `json:"clientOid,omitempty"`
	StpMode   SelfTradePreventionMode `json:"stpMode,omitempty"`
}

// IsolatedBatchOrderResult is the per-order outcome of a batch-place request.
type IsolatedBatchOrderResult struct {
	SuccessList []IsolatedBatchOrderSuccess `json:"successList"`
	FailureList []IsolatedBatchOrderFailure `json:"failureList"`
}

// IsolatedBatchOrderSuccess is one successfully accepted order from a batch.
type IsolatedBatchOrderSuccess struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// IsolatedBatchOrderFailure is one rejected order from a batch.
type IsolatedBatchOrderFailure struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
}

// CancelIsolatedOrderService -- POST /api/v2/margin/isolated/cancel-order (margin trade)
//
// Cancels a single isolated-margin order, identified by orderId or clientOid.
type CancelIsolatedOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCancelIsolatedOrderService(symbol string) *CancelIsolatedOrderService {
	return &CancelIsolatedOrderService{c: c, body: map[string]any{
		"symbol": symbol,
	}}
}

func (s *CancelIsolatedOrderService) SetOrderId(orderId string) *CancelIsolatedOrderService {
	s.body["orderId"] = orderId
	return s
}

func (s *CancelIsolatedOrderService) SetClientOid(clientOid string) *CancelIsolatedOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelIsolatedOrderService) Do(ctx context.Context) (*IsolatedCancelOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/cancel-order", s.body).WithSign()
	return request.Do[IsolatedCancelOrderResult](req)
}

// IsolatedCancelOrderResult is the result of cancelling a single isolated order.
type IsolatedCancelOrderResult struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// BatchCancelIsolatedOrderService -- POST /api/v2/margin/isolated/batch-cancel-order (margin trade)
//
// Cancels multiple isolated-margin orders for a single symbol in one request.
type BatchCancelIsolatedOrderService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewBatchCancelIsolatedOrderService(symbol string) *BatchCancelIsolatedOrderService {
	return &BatchCancelIsolatedOrderService{c: c, body: map[string]any{
		"symbol": symbol,
	}}
}

// SetOrderIdList sets the list of orders to cancel; each item is identified by
// orderId or clientOid.
func (s *BatchCancelIsolatedOrderService) SetOrderIdList(orderIdList []IsolatedCancelOrderItem) *BatchCancelIsolatedOrderService {
	s.body["orderIdList"] = orderIdList
	return s
}

func (s *BatchCancelIsolatedOrderService) Do(ctx context.Context) (*IsolatedBatchCancelResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/batch-cancel-order", s.body).WithSign()
	return request.Do[IsolatedBatchCancelResult](req)
}

// IsolatedCancelOrderItem identifies one order to cancel within a batch request.
type IsolatedCancelOrderItem struct {
	OrderId   string `json:"orderId,omitempty"`
	ClientOid string `json:"clientOid,omitempty"`
}

// IsolatedBatchCancelResult is the per-order outcome of a batch-cancel request.
type IsolatedBatchCancelResult struct {
	SuccessList []IsolatedBatchCancelSuccess `json:"successList"`
	FailureList []IsolatedBatchCancelFailure `json:"failureList"`
}

// IsolatedBatchCancelSuccess is one successfully cancelled order from a batch.
type IsolatedBatchCancelSuccess struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// IsolatedBatchCancelFailure is one order that failed to cancel in a batch.
type IsolatedBatchCancelFailure struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
}

// GetIsolatedOpenOrdersService -- GET /api/v2/margin/isolated/open-orders (margin read)
//
// Returns the trader's currently open (unfilled or partially filled)
// isolated-margin orders for a symbol.
type GetIsolatedOpenOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedOpenOrdersService(symbol string, startTime time.Time) *GetIsolatedOpenOrdersService {
	return &GetIsolatedOpenOrdersService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedOpenOrdersService) SetOrderId(orderId string) *GetIsolatedOpenOrdersService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetIsolatedOpenOrdersService) SetClientOid(clientOid string) *GetIsolatedOpenOrdersService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetIsolatedOpenOrdersService) SetEndTime(endTime time.Time) *GetIsolatedOpenOrdersService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedOpenOrdersService) SetLimit(limit int) *GetIsolatedOpenOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIsolatedOpenOrdersService) Do(ctx context.Context) (*IsolatedOpenOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/open-orders", s.params).WithSign()
	return request.Do[IsolatedOpenOrders](req)
}

// IsolatedOpenOrders is the paged list of open isolated-margin orders.
type IsolatedOpenOrders struct {
	OrderList []IsolatedOrder `json:"orderList"`
	MaxId     string          `json:"maxId"`
	MinId     string          `json:"minId"`
}

// IsolatedOrder is a single isolated-margin order as returned by the open-orders
// and history-orders endpoints.
type IsolatedOrder struct {
	Symbol           string           `json:"symbol"`
	OrderId          string           `json:"orderId"`
	ClientOid        string           `json:"clientOid"`
	OrderType        OrderType        `json:"orderType"`
	Side             Side             `json:"side"`
	Price            decimal.Decimal  `json:"price"`
	BaseSize         decimal.Decimal  `json:"baseSize"`
	QuoteSize        decimal.Decimal  `json:"quoteSize"`
	Size             decimal.Decimal  `json:"size"`
	Amount           decimal.Decimal  `json:"amount"`
	PriceAvg         decimal.Decimal  `json:"priceAvg"`
	Force            Force            `json:"force"`
	Status           OrderStatus      `json:"status"`
	LoanType         LoanType         `json:"loanType"`
	EnterPointSource EnterPointSource `json:"enterPointSource"`
	CTime            time.Time        `json:"cTime"`
	UTime            time.Time        `json:"uTime"`
}

// GetIsolatedHistoryOrdersService -- GET /api/v2/margin/isolated/history-orders (margin read)
//
// Returns the trader's historical (completed or cancelled) isolated-margin orders
// for a symbol within a time window.
type GetIsolatedHistoryOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedHistoryOrdersService(symbol string, startTime time.Time) *GetIsolatedHistoryOrdersService {
	return &GetIsolatedHistoryOrdersService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedHistoryOrdersService) SetOrderId(orderId string) *GetIsolatedHistoryOrdersService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetIsolatedHistoryOrdersService) SetClientOid(clientOid string) *GetIsolatedHistoryOrdersService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetIsolatedHistoryOrdersService) SetEnterPointSource(enterPointSource EnterPointSource) *GetIsolatedHistoryOrdersService {
	s.params["enterPointSource"] = string(enterPointSource)
	return s
}

func (s *GetIsolatedHistoryOrdersService) SetEndTime(endTime time.Time) *GetIsolatedHistoryOrdersService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedHistoryOrdersService) SetLimit(limit int) *GetIsolatedHistoryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIsolatedHistoryOrdersService) SetIdLessThan(idLessThan string) *GetIsolatedHistoryOrdersService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedHistoryOrdersService) Do(ctx context.Context) (*IsolatedHistoryOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/history-orders", s.params).WithSign()
	return request.Do[IsolatedHistoryOrders](req)
}

// IsolatedHistoryOrders is the paged list of historical isolated-margin orders.
type IsolatedHistoryOrders struct {
	OrderList []IsolatedOrder `json:"orderList"`
	MaxId     string          `json:"maxId"`
	MinId     string          `json:"minId"`
}

// GetIsolatedFillsService -- GET /api/v2/margin/isolated/fills (margin read)
//
// Returns the trader's isolated-margin transaction (fill) details for a symbol
// within a time window.
type GetIsolatedFillsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedFillsService(symbol string, startTime time.Time) *GetIsolatedFillsService {
	return &GetIsolatedFillsService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedFillsService) SetOrderId(orderId string) *GetIsolatedFillsService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetIsolatedFillsService) SetEndTime(endTime time.Time) *GetIsolatedFillsService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedFillsService) SetLimit(limit int) *GetIsolatedFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIsolatedFillsService) SetIdLessThan(idLessThan string) *GetIsolatedFillsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedFillsService) Do(ctx context.Context) (*IsolatedFills, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/fills", s.params).WithSign()
	return request.Do[IsolatedFills](req)
}

// IsolatedFills is the paged list of isolated-margin fills.
type IsolatedFills struct {
	Fills []IsolatedFill `json:"fills"`
	MaxId string         `json:"maxId"`
	MinId string         `json:"minId"`
}

// IsolatedFill is a single isolated-margin fill (trade execution).
type IsolatedFill struct {
	OrderId    string                `json:"orderId"`
	TradeId    string                `json:"tradeId"`
	Symbol     string                `json:"symbol"`
	OrderType  OrderType             `json:"orderType"`
	Side       Side                  `json:"side"`
	PriceAvg   decimal.Decimal       `json:"priceAvg"`
	Size       decimal.Decimal       `json:"size"`
	Amount     decimal.Decimal       `json:"amount"`
	TradeScope string                `json:"tradeScope"` // taker, maker
	FeeDetail  IsolatedFillFeeDetail `json:"feeDetail"`
	CTime      time.Time             `json:"cTime"`
	UTime      time.Time             `json:"uTime"`
}

// IsolatedFillFeeDetail is the fee breakdown for an isolated-margin fill.
type IsolatedFillFeeDetail struct {
	Deduction         string          `json:"deduction"` // yes, no
	FeeCoin           string          `json:"feeCoin"`
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"`
	TotalFee          decimal.Decimal `json:"totalFee"`
}

// GetIsolatedLiquidationOrdersService -- GET /api/v2/margin/isolated/liquidation-order (margin read)
//
// Returns the trader's isolated-margin forced-liquidation orders.
type GetIsolatedLiquidationOrdersService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedLiquidationOrdersService() *GetIsolatedLiquidationOrdersService {
	return &GetIsolatedLiquidationOrdersService{c: c, params: map[string]string{}}
}

// SetType filters by liquidation type: swap or place_order (default place_order).
func (s *GetIsolatedLiquidationOrdersService) SetType(typ string) *GetIsolatedLiquidationOrdersService {
	s.params["type"] = typ
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetSymbol(symbol string) *GetIsolatedLiquidationOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetFromCoin(fromCoin string) *GetIsolatedLiquidationOrdersService {
	s.params["fromCoin"] = fromCoin
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetToCoin(toCoin string) *GetIsolatedLiquidationOrdersService {
	s.params["toCoin"] = toCoin
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetStartTime(startTime time.Time) *GetIsolatedLiquidationOrdersService {
	s.params["startTime"] = strconv.FormatInt(startTime.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetEndTime(endTime time.Time) *GetIsolatedLiquidationOrdersService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetLimit(limit int) *GetIsolatedLiquidationOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetIsolatedLiquidationOrdersService) SetIdLessThan(idLessThan string) *GetIsolatedLiquidationOrdersService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedLiquidationOrdersService) Do(ctx context.Context) (*IsolatedLiquidationOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/liquidation-order", s.params).WithSign()
	return request.Do[IsolatedLiquidationOrders](req)
}

// IsolatedLiquidationOrders is the paged list of isolated-margin liquidation orders.
type IsolatedLiquidationOrders struct {
	ResultList []IsolatedLiquidationOrder `json:"resultList"`
	IdLessThan string                     `json:"idLessThan"`
	MaxId      string                     `json:"maxId"`
	MinId      string                     `json:"minId"`
}

// IsolatedLiquidationOrder is a single isolated-margin forced-liquidation order.
// The side here reports the system-driven direction liquidation_sell /
// liquidation_buy (underscore form), so it is kept as a plain string rather than
// the buy/sell Side enum.
type IsolatedLiquidationOrder struct {
	OrderId   string          `json:"orderId"`
	Symbol    string          `json:"symbol"`
	OrderType OrderType       `json:"orderType"`
	Side      string          `json:"side"` // liquidation_sell, liquidation_buy
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
