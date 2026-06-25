package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// OrderRef identifies a single order by its exchange and client identifiers. It
// is the reply shape for place/click-backhand/modify/cancel single-order
// operations.
type OrderRef struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
}

// PlaceOrderService -- POST /api/v2/mix/order/place-order (private, state-changing)
//
// Submits a single futures order. symbol, productType, marginMode, marginCoin,
// size, side and orderType are required; price/force are required for limit
// orders and tradeSide is required in hedge mode.
type PlaceOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewPlaceOrderService(symbol string, productType ProductType, marginMode MarginMode, marginCoin string, size decimal.Decimal, side Side, orderType OrderType) *PlaceOrderService {
	return &PlaceOrderService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
		"marginMode":  string(marginMode),
		"marginCoin":  marginCoin,
		"size":        size.String(),
		"side":        string(side),
		"orderType":   string(orderType),
	}}
}

// SetPrice sets the order price (required for limit orders).
func (s *PlaceOrderService) SetPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["price"] = price.String()
	return s
}

// SetForce sets the time-in-force policy (required for limit orders).
func (s *PlaceOrderService) SetForce(force Force) *PlaceOrderService {
	s.body["force"] = string(force)
	return s
}

// SetTradeSide sets whether the order opens or closes a position (required in
// hedge mode).
func (s *PlaceOrderService) SetTradeSide(tradeSide TradeSide) *PlaceOrderService {
	s.body["tradeSide"] = string(tradeSide)
	return s
}

// SetClientOid sets the client-generated order identifier.
func (s *PlaceOrderService) SetClientOrderID(clientOid string) *PlaceOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetReduceOnly sets the reduce-only flag (one-way mode only).
func (s *PlaceOrderService) SetReduceOnly(reduceOnly ReduceOnly) *PlaceOrderService {
	s.body["reduceOnly"] = string(reduceOnly)
	return s
}

// SetPresetStopSurplusPrice sets the preset take-profit trigger price.
func (s *PlaceOrderService) SetPresetStopSurplusPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetStopSurplusPrice"] = price.String()
	return s
}

// SetPresetStopLossPrice sets the preset stop-loss trigger price.
func (s *PlaceOrderService) SetPresetStopLossPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetStopLossPrice"] = price.String()
	return s
}

// SetPresetStopSurplusExecutePrice sets the preset take-profit execution price.
func (s *PlaceOrderService) SetPresetStopSurplusExecutePrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetStopSurplusExecutePrice"] = price.String()
	return s
}

// SetPresetStopLossExecutePrice sets the preset stop-loss execution price.
func (s *PlaceOrderService) SetPresetStopLossExecutePrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetStopLossExecutePrice"] = price.String()
	return s
}

// SetStpMode sets the self-trade prevention mode.
func (s *PlaceOrderService) SetStpMode(stpMode SelfTradePreventionMode) *PlaceOrderService {
	s.body["stpMode"] = string(stpMode)
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/place-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// ClickBackhandService -- POST /api/v2/mix/order/click-backhand (private, state-changing)
//
// Reversal: market-closes the current position and immediately re-opens an equal
// position in the opposite direction. symbol, marginCoin, productType and side
// are required; tradeSide is required in hedge mode.
type ClickBackhandService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewClickBackhandService(symbol string, marginCoin string, productType ProductType, side Side) *ClickBackhandService {
	return &ClickBackhandService{c: c, body: map[string]any{
		"symbol":      symbol,
		"marginCoin":  marginCoin,
		"productType": string(productType),
		"side":        string(side),
	}}
}

// SetSize sets the reversal amount (defaults to the full current position).
func (s *ClickBackhandService) SetSize(size decimal.Decimal) *ClickBackhandService {
	s.body["size"] = size.String()
	return s
}

// SetTradeSide sets whether the order opens or closes a position (required in
// hedge mode).
func (s *ClickBackhandService) SetTradeSide(tradeSide TradeSide) *ClickBackhandService {
	s.body["tradeSide"] = string(tradeSide)
	return s
}

// SetClientOid sets the client-generated order identifier.
func (s *ClickBackhandService) SetClientOrderID(clientOid string) *ClickBackhandService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *ClickBackhandService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/click-backhand", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// BatchOrderItem is a single order in a batch-place request. size, side and
// orderType are required; price/force are required for limit orders.
type BatchOrderItem struct {
	Size                   decimal.Decimal         `json:"size"`
	Price                  decimal.Decimal         `json:"price,omitzero"`
	Side                   Side                    `json:"side"`
	TradeSide              TradeSide               `json:"tradeSide,omitempty"`
	OrderType              OrderType               `json:"orderType"`
	Force                  Force                   `json:"force,omitempty"`
	ClientOrderID          string                  `json:"clientOid,omitempty"`
	ReduceOnly             ReduceOnly              `json:"reduceOnly,omitempty"`
	PresetStopSurplusPrice decimal.Decimal         `json:"presetStopSurplusPrice,omitzero"`
	PresetStopLossPrice    decimal.Decimal         `json:"presetStopLossPrice,omitzero"`
	StpMode                SelfTradePreventionMode `json:"stpMode,omitempty"`
}

// BatchPlaceOrderService -- POST /api/v2/mix/order/batch-place-order (private, state-changing)
//
// Places up to 50 orders in a single request; all orders share the same symbol,
// productType, marginMode and marginCoin. The orderList is sent as a nested JSON
// array of order objects and the reply splits into successList and failureList.
type BatchPlaceOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewBatchPlaceOrderService(symbol string, productType ProductType, marginMode MarginMode, marginCoin string, orderList []BatchOrderItem) *BatchPlaceOrderService {
	return &BatchPlaceOrderService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
		"marginMode":  string(marginMode),
		"marginCoin":  marginCoin,
		"orderList":   orderList,
	}}
}

// SetOrderList replaces the list of orders to place.
func (s *BatchPlaceOrderService) SetOrderList(orderList []BatchOrderItem) *BatchPlaceOrderService {
	s.body["orderList"] = orderList
	return s
}

func (s *BatchPlaceOrderService) Do(ctx context.Context) (*BatchOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/batch-place-order").SetBody(s.body).WithSign()
	return request.Do[BatchOrderResult](req)
}

// BatchOrderResult is the partial-success reply shared by batch place/cancel and
// cancel-all: successList holds the accepted orders, failureList the rejected
// ones with their error code/message.
type BatchOrderResult struct {
	SuccessList []OrderRef         `json:"successList"`
	FailureList []OrderFailureItem `json:"failureList"`
}

// OrderFailureItem is one rejected order in a batch operation.
type OrderFailureItem struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
	ErrorMsg      string `json:"errorMsg"`
	ErrorCode     string `json:"errorCode"`
}

// ModifyOrderService -- POST /api/v2/mix/order/modify-order (private, state-changing)
//
// Amends an open order. symbol, productType, marginCoin and newClientOid are
// required; identify the order by orderId or clientOid (orderId wins). At least
// one of newSize/newPrice should be supplied to change the order.
type ModifyOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewModifyOrderService(symbol string, productType ProductType, marginCoin string, newClientOid string) *ModifyOrderService {
	return &ModifyOrderService{c: c, body: map[string]any{
		"symbol":       symbol,
		"productType":  string(productType),
		"marginCoin":   marginCoin,
		"newClientOid": newClientOid,
	}}
}

// SetOrderId sets the order identifier (orderId or clientOid is required).
func (s *ModifyOrderService) SetOrderID(orderId string) *ModifyOrderService {
	s.body["orderId"] = orderId
	return s
}

// SetClientOid sets the client order identifier (orderId or clientOid is required).
func (s *ModifyOrderService) SetClientOrderID(clientOid string) *ModifyOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetNewSize sets the modified order amount.
func (s *ModifyOrderService) SetNewSize(newSize decimal.Decimal) *ModifyOrderService {
	s.body["newSize"] = newSize.String()
	return s
}

// SetNewPrice sets the modified order price.
func (s *ModifyOrderService) SetNewPrice(newPrice decimal.Decimal) *ModifyOrderService {
	s.body["newPrice"] = newPrice.String()
	return s
}

// SetNewPresetStopSurplusPrice sets the modified take-profit price (0 deletes it).
func (s *ModifyOrderService) SetNewPresetStopSurplusPrice(price decimal.Decimal) *ModifyOrderService {
	s.body["newPresetStopSurplusPrice"] = price.String()
	return s
}

// SetNewPresetStopLossPrice sets the modified stop-loss price (0 deletes it).
func (s *ModifyOrderService) SetNewPresetStopLossPrice(price decimal.Decimal) *ModifyOrderService {
	s.body["newPresetStopLossPrice"] = price.String()
	return s
}

func (s *ModifyOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/modify-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// CancelOrderService -- POST /api/v2/mix/order/cancel-order (private, state-changing)
//
// Cancels a single open order. symbol and productType are required; identify the
// order by orderId or clientOid (orderId wins if both are set).
type CancelOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewCancelOrderService(symbol string, productType ProductType) *CancelOrderService {
	return &CancelOrderService{c: c, body: map[string]any{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetMarginCoin sets the margin coin (capitalized).
func (s *CancelOrderService) SetMarginCoin(marginCoin string) *CancelOrderService {
	s.body["marginCoin"] = marginCoin
	return s
}

// SetOrderId sets the order identifier (orderId or clientOid is required).
func (s *CancelOrderService) SetOrderID(orderId string) *CancelOrderService {
	s.body["orderId"] = orderId
	return s
}

// SetClientOid sets the client order identifier (orderId or clientOid is required).
func (s *CancelOrderService) SetClientOrderID(clientOid string) *CancelOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/cancel-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// BatchCancelItem identifies a single order to cancel in a batch-cancel request.
// Either OrderId or ClientOid must be set (orderId takes precedence).
type BatchCancelItem struct {
	OrderID       string `json:"orderId,omitempty"`
	ClientOrderID string `json:"clientOid,omitempty"`
}

// BatchCancelOrdersService -- POST /api/v2/mix/order/batch-cancel-orders (private, state-changing)
//
// Cancels up to 50 orders in a single request. productType is required; symbol
// is required when orderIdList is provided. The reply splits into successList
// and failureList.
type BatchCancelOrdersService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewBatchCancelOrdersService(productType ProductType) *BatchCancelOrdersService {
	return &BatchCancelOrdersService{c: c, body: map[string]any{
		"productType": string(productType),
	}}
}

// SetSymbol sets the trading pair (required when orderIdList is provided).
func (s *BatchCancelOrdersService) SetSymbol(symbol string) *BatchCancelOrdersService {
	s.body["symbol"] = symbol
	return s
}

// SetMarginCoin sets the margin coin (capitalized).
func (s *BatchCancelOrdersService) SetMarginCoin(marginCoin string) *BatchCancelOrdersService {
	s.body["marginCoin"] = marginCoin
	return s
}

// SetOrderIdList sets the list of orders to cancel (max 50). When omitted, all
// open orders matching the other filters are cancelled.
func (s *BatchCancelOrdersService) SetOrderIDList(orderIdList []BatchCancelItem) *BatchCancelOrdersService {
	s.body["orderIdList"] = orderIdList
	return s
}

func (s *BatchCancelOrdersService) Do(ctx context.Context) (*BatchOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/batch-cancel-orders").SetBody(s.body).WithSign()
	return request.Do[BatchOrderResult](req)
}

// ClosePositionsService -- POST /api/v2/mix/order/close-positions (private, state-changing)
//
// Flash-closes positions at market price. productType is required; without
// symbol it closes all positions, without holdSide it closes both sides. The
// reply lists the close orders that were and weren't submitted.
type ClosePositionsService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewClosePositionsService(productType ProductType) *ClosePositionsService {
	return &ClosePositionsService{c: c, body: map[string]any{
		"productType": string(productType),
	}}
}

// SetSymbol limits the flash-close to a single trading pair.
func (s *ClosePositionsService) SetSymbol(symbol string) *ClosePositionsService {
	s.body["symbol"] = symbol
	return s
}

// SetHoldSide limits the flash-close to one position direction (hedge mode).
func (s *ClosePositionsService) SetHoldSide(holdSide HoldSide) *ClosePositionsService {
	s.body["holdSide"] = string(holdSide)
	return s
}

func (s *ClosePositionsService) Do(ctx context.Context) (*ClosePositionsResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/close-positions", s.body).WithSign()
	return request.Do[ClosePositionsResult](req)
}

// ClosePositionsResult is the flash-close reply: successList holds the submitted
// close orders, failureList the positions that could not be closed.
type ClosePositionsResult struct {
	SuccessList []CloseOrderRef     `json:"successList"`
	FailureList []CloseOrderFailure `json:"failureList"`
}

// CloseOrderRef is one submitted flash-close order.
type CloseOrderRef struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
	Symbol        string `json:"symbol"`
}

// CloseOrderFailure is one position that could not be flash-closed.
type CloseOrderFailure struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
	Symbol        string `json:"symbol"`
	ErrorMsg      string `json:"errorMsg"`
	ErrorCode     string `json:"errorCode"`
}

// CancelAllOrdersService -- POST /api/v2/mix/order/cancel-all-orders (private, state-changing)
//
// Cancels every open order in a product line. productType is required. The reply
// splits into successList and failureList.
type CancelAllOrdersService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewCancelAllOrdersService(productType ProductType) *CancelAllOrdersService {
	return &CancelAllOrdersService{c: c, body: map[string]any{
		"productType": string(productType),
	}}
}

// SetMarginCoin sets the margin coin (capitalized).
func (s *CancelAllOrdersService) SetMarginCoin(marginCoin string) *CancelAllOrdersService {
	s.body["marginCoin"] = marginCoin
	return s
}

// SetRequestTime sets the request timestamp (ms).
func (s *CancelAllOrdersService) SetRequestTime(t time.Time) *CancelAllOrdersService {
	s.body["requestTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetReceiveWindow sets the valid window period (ms).
func (s *CancelAllOrdersService) SetReceiveWindow(receiveWindow string) *CancelAllOrdersService {
	s.body["receiveWindow"] = receiveWindow
	return s
}

func (s *CancelAllOrdersService) Do(ctx context.Context) (*BatchOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/cancel-all-orders", s.body).WithSign()
	return request.Do[BatchOrderResult](req)
}

// MixOrder is the union of the order-detail, orders-pending and orders-history
// shapes; fields not present for a given endpoint arrive empty. feeDetail is a
// JSON-encoded string on these order endpoints (it is a nested object on the
// fill endpoints).
type MixOrder struct {
	Symbol                        string          `json:"symbol"`
	Size                          decimal.Decimal `json:"size"`
	OrderID                       string          `json:"orderId"`
	ClientOrderID                 string          `json:"clientOid"`
	BaseVolume                    decimal.Decimal `json:"baseVolume"`
	Fee                           decimal.Decimal `json:"fee"`
	Price                         decimal.Decimal `json:"price"`
	PriceAvg                      decimal.Decimal `json:"priceAvg"`
	Status                        OrderStatus     `json:"status"` // orders-pending key
	State                         OrderStatus     `json:"state"`  // order-detail key (same value, different name)
	Side                          Side            `json:"side"`
	NewTradeSide                  string          `json:"newTradeSide"` // order-detail only
	Force                         Force           `json:"force"`
	TotalProfits                  decimal.Decimal `json:"totalProfits"`
	PosSide                       string          `json:"posSide"`
	MarginCoin                    string          `json:"marginCoin"`
	PresetStopSurplusPrice        decimal.Decimal `json:"presetStopSurplusPrice"`
	PresetStopSurplusType         TriggerType     `json:"presetStopSurplusType"`        // order-detail key
	PresetStopSurplusTriggerType  TriggerType     `json:"presetStopSurplusTriggerType"` // orders-pending key
	PresetStopSurplusExecutePrice decimal.Decimal `json:"presetStopSurplusExecutePrice"`
	PresetStopLossPrice           decimal.Decimal `json:"presetStopLossPrice"`
	PresetStopLossType            TriggerType     `json:"presetStopLossType"`        // order-detail key
	PresetStopLossTriggerType     TriggerType     `json:"presetStopLossTriggerType"` // orders-pending key
	PresetStopLossExecutePrice    decimal.Decimal `json:"presetStopLossExecutePrice"`
	QuoteVolume                   decimal.Decimal `json:"quoteVolume"`
	OrderType                     OrderType       `json:"orderType"`
	Leverage                      decimal.Decimal `json:"leverage"`
	MarginMode                    MarginMode      `json:"marginMode"`
	ReduceOnly                    ReduceOnly      `json:"reduceOnly"`
	EnterPointSource              string          `json:"enterPointSource"`
	TradeSide                     string          `json:"tradeSide"`
	PosMode                       PositionMode    `json:"posMode"`
	PosAvg                        decimal.Decimal `json:"posAvg"`
	OrderSource                   string          `json:"orderSource"`
	LiqPrice                      decimal.Decimal `json:"liqPrice"`
	CancelReason                  string          `json:"cancelReason"`
	FeeDetail                     string          `json:"feeDetail"`
	CTime                         time.Time       `json:"cTime"`
	UTime                         time.Time       `json:"uTime"`
}

// EffectiveStatus returns the order status regardless of which endpoint shape
// produced the value: orders-pending / orders-history fill Status, order-detail
// fills State (the same semantic value under a different JSON key).
func (o *MixOrder) EffectiveStatus() OrderStatus {
	if o.Status != "" {
		return o.Status
	}
	return o.State
}

// GetOrderDetailService -- GET /api/v2/mix/order/detail (private)
//
// Returns the details of a single order, looked up by orderId or clientOid (one
// is required). symbol and productType are required.
type GetOrderDetailService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrderDetailService(symbol string, productType ProductType) *GetOrderDetailService {
	return &GetOrderDetailService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetOrderId sets the order identifier (orderId or clientOid is required).
func (s *GetOrderDetailService) SetOrderID(orderId string) *GetOrderDetailService {
	s.params["orderId"] = orderId
	return s
}

// SetClientOid sets the client order identifier (orderId or clientOid is required).
func (s *GetOrderDetailService) SetClientOrderID(clientOid string) *GetOrderDetailService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetOrderDetailService) Do(ctx context.Context) (*MixOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/detail", s.params).WithSign()
	return request.Do[MixOrder](req)
}

// GetOrdersPendingService -- GET /api/v2/mix/order/orders-pending (private)
//
// Returns the account's currently open (live / partially-filled) orders for a
// product line, paginated by idLessThan. productType is required.
type GetOrdersPendingService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrdersPendingService(productType ProductType) *GetOrdersPendingService {
	return &GetOrdersPendingService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetOrdersPendingService) SetOrderID(orderId string) *GetOrdersPendingService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetOrdersPendingService) SetClientOrderID(clientOid string) *GetOrdersPendingService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetOrdersPendingService) SetSymbol(symbol string) *GetOrdersPendingService {
	s.params["symbol"] = symbol
	return s
}

// SetStatus filters by order status (live or partially_filled).
func (s *GetOrdersPendingService) SetStatus(status OrderStatus) *GetOrdersPendingService {
	s.params["status"] = string(status)
	return s
}

// SetIdLessThan pages to orders older than the given order id.
func (s *GetOrdersPendingService) SetIDLessThan(idLessThan string) *GetOrdersPendingService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters orders at or after t.
func (s *GetOrdersPendingService) SetStartTime(t time.Time) *GetOrdersPendingService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetOrdersPendingService) SetEndTime(t time.Time) *GetOrdersPendingService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOrdersPendingService) SetLimit(limit int) *GetOrdersPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersPendingService) Do(ctx context.Context) (*MixOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/orders-pending", s.params).WithSign()
	return request.Do[MixOrderList](req)
}

// GetOrdersHistoryService -- GET /api/v2/mix/order/orders-history (private)
//
// Returns the account's historical (filled / cancelled) orders for a product
// line, paginated by idLessThan. productType is required.
type GetOrdersHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrdersHistoryService(productType ProductType) *GetOrdersHistoryService {
	return &GetOrdersHistoryService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetOrdersHistoryService) SetOrderID(orderId string) *GetOrdersHistoryService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetOrdersHistoryService) SetClientOrderID(clientOid string) *GetOrdersHistoryService {
	s.params["clientOid"] = clientOid
	return s
}

func (s *GetOrdersHistoryService) SetSymbol(symbol string) *GetOrdersHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetIdLessThan pages to orders older than the given order id.
func (s *GetOrdersHistoryService) SetIDLessThan(idLessThan string) *GetOrdersHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetOrderSource filters by order origin (normal, market, profit_market, ...).
func (s *GetOrdersHistoryService) SetOrderSource(orderSource string) *GetOrdersHistoryService {
	s.params["orderSource"] = orderSource
	return s
}

// SetStartTime filters orders at or after t.
func (s *GetOrdersHistoryService) SetStartTime(t time.Time) *GetOrdersHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetOrdersHistoryService) SetEndTime(t time.Time) *GetOrdersHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOrdersHistoryService) SetLimit(limit int) *GetOrdersHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersHistoryService) Do(ctx context.Context) (*MixOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/orders-history", s.params).WithSign()
	return request.Do[MixOrderList](req)
}

// MixOrderList is the page returned by orders-pending and orders-history:
// entrustedList holds the orders and endId is the pagination cursor.
type MixOrderList struct {
	EntrustedList []MixOrder `json:"entrustedList"`
	EndID         string     `json:"endId"`
}

// MixFillFeeDetail is one fee line in a fill's fee breakdown. It is returned as a
// nested object (not a JSON-encoded string) on the fill endpoints.
type MixFillFeeDetail struct {
	Deduction         string          `json:"deduction"`
	FeeCoin           string          `json:"feeCoin"`
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"`
	TotalFee          decimal.Decimal `json:"totalFee"`
}

// MixFill is a single trade execution returned by the fills and fill-history
// endpoints.
type MixFill struct {
	TradeID          string             `json:"tradeId"`
	Symbol           string             `json:"symbol"`
	MarginCoin       string             `json:"marginCoin"`
	OrderID          string             `json:"orderId"`
	Price            decimal.Decimal    `json:"price"`
	BaseVolume       decimal.Decimal    `json:"baseVolume"`
	QuoteVolume      decimal.Decimal    `json:"quoteVolume"`
	Side             Side               `json:"side"`
	Profit           decimal.Decimal    `json:"profit"`
	EnterPointSource string             `json:"enterPointSource"`
	TradeSide        string             `json:"tradeSide"`
	PosMode          PositionMode       `json:"posMode"`
	TradeScope       string             `json:"tradeScope"`
	FeeDetail        []MixFillFeeDetail `json:"feeDetail"`
	CTime            time.Time          `json:"cTime"`
}

// GetOrderFillsService -- GET /api/v2/mix/order/fills (private)
//
// Returns the account's recent trade fills for a product line, paginated by
// idLessThan. productType is required.
type GetOrderFillsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrderFillsService(productType ProductType) *GetOrderFillsService {
	return &GetOrderFillsService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetOrderFillsService) SetOrderID(orderId string) *GetOrderFillsService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetOrderFillsService) SetSymbol(symbol string) *GetOrderFillsService {
	s.params["symbol"] = symbol
	return s
}

// SetIdLessThan pages to fills older than the given tradeId.
func (s *GetOrderFillsService) SetIDLessThan(idLessThan string) *GetOrderFillsService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters fills at or after t.
func (s *GetOrderFillsService) SetStartTime(t time.Time) *GetOrderFillsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters fills at or before t.
func (s *GetOrderFillsService) SetEndTime(t time.Time) *GetOrderFillsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetOrderFillsService) SetLimit(limit int) *GetOrderFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrderFillsService) Do(ctx context.Context) (*MixFillList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/fills", s.params).WithSign()
	return request.Do[MixFillList](req)
}

// GetFillHistoryService -- GET /api/v2/mix/order/fill-history (private)
//
// Returns the account's historical trade fills for a product line, paginated by
// idLessThan (max one-week range per query). productType is required.
type GetFillHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetFillHistoryService(productType ProductType) *GetFillHistoryService {
	return &GetFillHistoryService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetFillHistoryService) SetOrderID(orderId string) *GetFillHistoryService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetFillHistoryService) SetSymbol(symbol string) *GetFillHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetIdLessThan pages to fills older than the given tradeId.
func (s *GetFillHistoryService) SetIDLessThan(idLessThan string) *GetFillHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters fills at or after t (max one-week range).
func (s *GetFillHistoryService) SetStartTime(t time.Time) *GetFillHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters fills at or before t (max one-week range).
func (s *GetFillHistoryService) SetEndTime(t time.Time) *GetFillHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFillHistoryService) SetLimit(limit int) *GetFillHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFillHistoryService) Do(ctx context.Context) (*MixFillList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/fill-history", s.params).WithSign()
	return request.Do[MixFillList](req)
}

// MixFillList is the page returned by fills and fill-history: fillList holds the
// executions and endId is the pagination cursor.
type MixFillList struct {
	FillList []MixFill `json:"fillList"`
	EndID    string    `json:"endId"`
}
