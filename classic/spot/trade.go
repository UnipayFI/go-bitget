package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceOrderService -- POST /api/v2/spot/trade/place-order (private, state-changing)
//
// Places a single spot order (limit or market, optionally with a SPOT TP/SL or
// preset take-profit/stop-loss).
type PlaceOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewPlaceOrderService(symbol string, side Side, orderType OrderType, force Force, size decimal.Decimal) *PlaceOrderService {
	return &PlaceOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"side":      string(side),
		"orderType": string(orderType),
		"force":     string(force),
		"size":      size.String(),
	}}
}

// SetPrice sets the limit price (required for limit orders).
func (s *PlaceOrderService) SetPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["price"] = price.String()
	return s
}

// SetClientOid sets a custom order id.
func (s *PlaceOrderService) SetClientOid(clientOid string) *PlaceOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetTriggerPrice sets the SPOT TP/SL trigger price (only for tpsl orders).
func (s *PlaceOrderService) SetTriggerPrice(triggerPrice decimal.Decimal) *PlaceOrderService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetTpslType selects a normal order or a SPOT TP/SL order (default normal).
func (s *PlaceOrderService) SetTpslType(tpslType TpslType) *PlaceOrderService {
	s.body["tpslType"] = string(tpslType)
	return s
}

// SetRequestTime sets the client request time.
func (s *PlaceOrderService) SetRequestTime(t time.Time) *PlaceOrderService {
	s.body["requestTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetReceiveWindow sets the valid time window in milliseconds.
func (s *PlaceOrderService) SetReceiveWindow(ms int64) *PlaceOrderService {
	s.body["receiveWindow"] = strconv.FormatInt(ms, 10)
	return s
}

// SetStpMode sets the self-trade prevention mode (default none).
func (s *PlaceOrderService) SetStpMode(mode SelfTradePreventionMode) *PlaceOrderService {
	s.body["stpMode"] = string(mode)
	return s
}

// SetPresetTakeProfitPrice sets the take-profit trigger price.
func (s *PlaceOrderService) SetPresetTakeProfitPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetTakeProfitPrice"] = price.String()
	return s
}

// SetExecuteTakeProfitPrice sets the take-profit execute price.
func (s *PlaceOrderService) SetExecuteTakeProfitPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["executeTakeProfitPrice"] = price.String()
	return s
}

// SetPresetStopLossPrice sets the stop-loss trigger price.
func (s *PlaceOrderService) SetPresetStopLossPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["presetStopLossPrice"] = price.String()
	return s
}

// SetExecuteStopLossPrice sets the stop-loss execute price.
func (s *PlaceOrderService) SetExecuteStopLossPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["executeStopLossPrice"] = price.String()
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*PlaceOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/place-order", s.body).WithSign()
	return request.Do[PlaceOrderResponse](req)
}

// PlaceOrderResponse is the result of placing a single order.
type PlaceOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CancelReplaceOrderService -- POST /api/v2/spot/trade/cancel-replace-order (private, state-changing)
//
// Cancels an existing order and submits a replacement order in one request.
type CancelReplaceOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCancelReplaceOrderService(symbol string, price, size decimal.Decimal) *CancelReplaceOrderService {
	return &CancelReplaceOrderService{c: c, body: map[string]any{
		"symbol": symbol,
		"price":  price.String(),
		"size":   size.String(),
	}}
}

// SetOrderID identifies the order to cancel (either orderId or clientOid).
func (s *CancelReplaceOrderService) SetOrderID(orderID string) *CancelReplaceOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the order to cancel (either orderId or clientOid).
func (s *CancelReplaceOrderService) SetClientOid(clientOid string) *CancelReplaceOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetNewClientOid sets a custom id for the replacement order.
func (s *CancelReplaceOrderService) SetNewClientOid(newClientOid string) *CancelReplaceOrderService {
	s.body["newClientOid"] = newClientOid
	return s
}

// SetPresetTakeProfitPrice sets the take-profit trigger price.
func (s *CancelReplaceOrderService) SetPresetTakeProfitPrice(price decimal.Decimal) *CancelReplaceOrderService {
	s.body["presetTakeProfitPrice"] = price.String()
	return s
}

// SetExecuteTakeProfitPrice sets the take-profit execute price.
func (s *CancelReplaceOrderService) SetExecuteTakeProfitPrice(price decimal.Decimal) *CancelReplaceOrderService {
	s.body["executeTakeProfitPrice"] = price.String()
	return s
}

// SetPresetStopLossPrice sets the stop-loss trigger price.
func (s *CancelReplaceOrderService) SetPresetStopLossPrice(price decimal.Decimal) *CancelReplaceOrderService {
	s.body["presetStopLossPrice"] = price.String()
	return s
}

// SetExecuteStopLossPrice sets the stop-loss execute price.
func (s *CancelReplaceOrderService) SetExecuteStopLossPrice(price decimal.Decimal) *CancelReplaceOrderService {
	s.body["executeStopLossPrice"] = price.String()
	return s
}

func (s *CancelReplaceOrderService) Do(ctx context.Context) (*CancelReplaceOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/cancel-replace-order", s.body).WithSign()
	return request.Do[CancelReplaceOrderResponse](req)
}

// CancelReplaceOrderResponse is the result of a cancel-replace operation.
type CancelReplaceOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	Success   string `json:"success"` // success, failure
	Msg       string `json:"msg"`
}

// CancelReplaceOrderItem is one cancel-replace instruction in a batch request.
type CancelReplaceOrderItem struct {
	Symbol                 string `json:"symbol"`
	Price                  string `json:"price"`
	Size                   string `json:"size"`
	OrderID                string `json:"orderId,omitempty"`
	ClientOid              string `json:"clientOid,omitempty"`
	NewClientOid           string `json:"newClientOid,omitempty"`
	PresetTakeProfitPrice  string `json:"presetTakeProfitPrice,omitempty"`
	ExecuteTakeProfitPrice string `json:"executeTakeProfitPrice,omitempty"`
	PresetStopLossPrice    string `json:"presetStopLossPrice,omitempty"`
	ExecuteStopLossPrice   string `json:"executeStopLossPrice,omitempty"`
}

// BatchCancelReplaceOrderService -- POST /api/v2/spot/trade/batch-cancel-replace-order (private, state-changing)
//
// Cancels and replaces up to 50 orders in one request. The body is a JSON array
// of items under "orderList".
type BatchCancelReplaceOrderService struct {
	c         *SpotClient
	orderList []CancelReplaceOrderItem
}

func (c *SpotClient) NewBatchCancelReplaceOrderService(orderList []CancelReplaceOrderItem) *BatchCancelReplaceOrderService {
	return &BatchCancelReplaceOrderService{c: c, orderList: orderList}
}

func (s *BatchCancelReplaceOrderService) Do(ctx context.Context) ([]CancelReplaceOrderResponse, error) {
	body := map[string]any{"orderList": s.orderList}
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/batch-cancel-replace-order").SetBody(body).WithSign()
	resp, err := request.Do[[]CancelReplaceOrderResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchOrderItem is one order in a batch place request.
type BatchOrderItem struct {
	Symbol                 string    `json:"symbol,omitempty"`
	Side                   Side      `json:"side"`
	OrderType              OrderType `json:"orderType"`
	Force                  Force     `json:"force"`
	Price                  string    `json:"price,omitempty"`
	Size                   string    `json:"size"`
	ClientOid              string    `json:"clientOid,omitempty"`
	StpMode                string    `json:"stpMode,omitempty"`
	PresetTakeProfitPrice  string    `json:"presetTakeProfitPrice,omitempty"`
	ExecuteTakeProfitPrice string    `json:"executeTakeProfitPrice,omitempty"`
	PresetStopLossPrice    string    `json:"presetStopLossPrice,omitempty"`
	ExecuteStopLossPrice   string    `json:"executeStopLossPrice,omitempty"`
}

// BatchPlaceOrdersService -- POST /api/v2/spot/trade/batch-orders (private, state-changing)
//
// Places up to 50 orders in one request. The orders are carried as a JSON array
// under "orderList"; symbol/batchMode are optional top-level fields.
type BatchPlaceOrdersService struct {
	c         *SpotClient
	body      map[string]any
	orderList []BatchOrderItem
}

func (c *SpotClient) NewBatchPlaceOrdersService(orderList []BatchOrderItem) *BatchPlaceOrdersService {
	return &BatchPlaceOrdersService{c: c, body: map[string]any{}, orderList: orderList}
}

// SetSymbol sets the shared trading pair for all orders in the batch.
func (s *BatchPlaceOrdersService) SetSymbol(symbol string) *BatchPlaceOrdersService {
	s.body["symbol"] = symbol
	return s
}

// SetBatchMode selects the batch order mode: single (default) or multiple.
func (s *BatchPlaceOrdersService) SetBatchMode(batchMode string) *BatchPlaceOrdersService {
	s.body["batchMode"] = batchMode
	return s
}

func (s *BatchPlaceOrdersService) Do(ctx context.Context) (*BatchOrdersResponse, error) {
	s.body["orderList"] = s.orderList
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/batch-orders").SetBody(s.body).WithSign()
	return request.Do[BatchOrdersResponse](req)
}

// BatchOrdersResponse is the result of a batch place/cancel request, split into
// the orders that succeeded and the ones that failed.
type BatchOrdersResponse struct {
	SuccessList []BatchOrderSuccess `json:"successList"`
	FailureList []BatchOrderFailure `json:"failureList"`
}

// BatchOrderSuccess is one successfully processed order in a batch response.
type BatchOrderSuccess struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// BatchOrderFailure is one failed order in a batch response.
type BatchOrderFailure struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode string `json:"errorCode"`
}

// CancelOrderService -- POST /api/v2/spot/trade/cancel-order (private, state-changing)
//
// Cancels a single order by orderId or clientOid.
type CancelOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCancelOrderService(symbol string) *CancelOrderService {
	return &CancelOrderService{c: c, body: map[string]any{"symbol": symbol}}
}

// SetOrderID identifies the order to cancel (either orderId or clientOid).
func (s *CancelOrderService) SetOrderID(orderID string) *CancelOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the order to cancel (either orderId or clientOid).
func (s *CancelOrderService) SetClientOid(clientOid string) *CancelOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetTpslType selects a normal order or a SPOT TP/SL order (default normal).
func (s *CancelOrderService) SetTpslType(tpslType TpslType) *CancelOrderService {
	s.body["tpslType"] = string(tpslType)
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*CancelOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/cancel-order", s.body).WithSign()
	return request.Do[CancelOrderResponse](req)
}

// CancelOrderResponse is the result of cancelling a single order.
type CancelOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// BatchCancelOrderItem is one cancel instruction in a batch cancel request.
type BatchCancelOrderItem struct {
	OrderID   string `json:"orderId,omitempty"`
	ClientOid string `json:"clientOid,omitempty"`
	Symbol    string `json:"symbol,omitempty"`
}

// BatchCancelOrderService -- POST /api/v2/spot/trade/batch-cancel-order (private, state-changing)
//
// Cancels up to 50 orders in one request. The orders are carried as a JSON array
// under "orderList"; symbol/batchMode are optional top-level fields.
type BatchCancelOrderService struct {
	c         *SpotClient
	body      map[string]any
	orderList []BatchCancelOrderItem
}

func (c *SpotClient) NewBatchCancelOrderService(orderList []BatchCancelOrderItem) *BatchCancelOrderService {
	return &BatchCancelOrderService{c: c, body: map[string]any{}, orderList: orderList}
}

// SetSymbol sets the shared trading pair for all cancels in the batch.
func (s *BatchCancelOrderService) SetSymbol(symbol string) *BatchCancelOrderService {
	s.body["symbol"] = symbol
	return s
}

// SetBatchMode selects the batch cancel mode: single (default) or multiple.
func (s *BatchCancelOrderService) SetBatchMode(batchMode string) *BatchCancelOrderService {
	s.body["batchMode"] = batchMode
	return s
}

func (s *BatchCancelOrderService) Do(ctx context.Context) (*BatchOrdersResponse, error) {
	s.body["orderList"] = s.orderList
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/batch-cancel-order").SetBody(s.body).WithSign()
	return request.Do[BatchOrdersResponse](req)
}

// CancelSymbolOrderService -- POST /api/v2/spot/trade/cancel-symbol-order (private, state-changing)
//
// Cancels all open orders for a symbol. The cancellation is processed
// asynchronously; confirm via Get History Orders.
type CancelSymbolOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCancelSymbolOrderService(symbol string) *CancelSymbolOrderService {
	return &CancelSymbolOrderService{c: c, body: map[string]any{"symbol": symbol}}
}

func (s *CancelSymbolOrderService) Do(ctx context.Context) (*CancelSymbolOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/cancel-symbol-order", s.body).WithSign()
	return request.Do[CancelSymbolOrderResponse](req)
}

// CancelSymbolOrderResponse echoes the symbol whose orders are being cancelled.
type CancelSymbolOrderResponse struct {
	Symbol string `json:"symbol"`
}

// GetOrderInfoService -- GET /api/v2/spot/trade/orderInfo (private)
//
// Returns full details for a single order, identified by orderId or clientOid.
type GetOrderInfoService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOrderInfoService() *GetOrderInfoService {
	return &GetOrderInfoService{c: c, params: map[string]string{}}
}

// SetOrderID identifies the order to fetch (either orderId or clientOid).
func (s *GetOrderInfoService) SetOrderID(orderID string) *GetOrderInfoService {
	s.params["orderId"] = orderID
	return s
}

// SetClientOid identifies the order to fetch (either orderId or clientOid).
func (s *GetOrderInfoService) SetClientOid(clientOid string) *GetOrderInfoService {
	s.params["clientOid"] = clientOid
	return s
}

// SetRequestTime sets the client request time.
func (s *GetOrderInfoService) SetRequestTime(t time.Time) *GetOrderInfoService {
	s.params["requestTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetReceiveWindow sets the valid time window in milliseconds.
func (s *GetOrderInfoService) SetReceiveWindow(ms int64) *GetOrderInfoService {
	s.params["receiveWindow"] = strconv.FormatInt(ms, 10)
	return s
}

func (s *GetOrderInfoService) Do(ctx context.Context) ([]OrderInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/orderInfo", s.params).WithSign()
	resp, err := request.Do[[]OrderInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OrderInfo is the full detail of a single order. feeDetail is a Bitget-encoded
// JSON string (not a nested object), so it is kept as a raw string.
type OrderInfo struct {
	UserID           string          `json:"userId"`
	Symbol           string          `json:"symbol"`
	OrderID          string          `json:"orderId"`
	ClientOid        string          `json:"clientOid"`
	Price            decimal.Decimal `json:"price"`
	Size             decimal.Decimal `json:"size"`
	OrderType        OrderType       `json:"orderType"`
	Side             Side            `json:"side"`
	Status           OrderStatus     `json:"status"`
	PriceAvg         decimal.Decimal `json:"priceAvg"`
	BaseVolume       decimal.Decimal `json:"baseVolume"`
	QuoteVolume      decimal.Decimal `json:"quoteVolume"`
	BaseCoin         string          `json:"baseCoin"`
	QuoteCoin        string          `json:"quoteCoin"`
	EnterPointSource string          `json:"enterPointSource"`
	FeeDetail        string          `json:"feeDetail"` // JSON-encoded fee breakdown string
	OrderSource      string          `json:"orderSource"`
	TpslType         TpslType        `json:"tpslType"`
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`
	CancelReason     string          `json:"cancelReason"` // normal_cancel, stp_cancel
	CTime            time.Time       `json:"cTime"`
	UTime            time.Time       `json:"uTime"`
}

// GetUnfilledOrdersService -- GET /api/v2/spot/trade/unfilled-orders (private)
//
// Returns the account's current (open / unfilled) orders.
type GetUnfilledOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetUnfilledOrdersService() *GetUnfilledOrdersService {
	return &GetUnfilledOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single trading pair.
func (s *GetUnfilledOrdersService) SetSymbol(symbol string) *GetUnfilledOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters orders created at or after t.
func (s *GetUnfilledOrdersService) SetStartTime(t time.Time) *GetUnfilledOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t.
func (s *GetUnfilledOrdersService) SetEndTime(t time.Time) *GetUnfilledOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backwards: return orders with orderId older than this.
func (s *GetUnfilledOrdersService) SetIDLessThan(idLessThan string) *GetUnfilledOrdersService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of orders returned (default 100, max 100).
func (s *GetUnfilledOrdersService) SetLimit(limit int) *GetUnfilledOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOrderID filters to a single order id.
func (s *GetUnfilledOrdersService) SetOrderID(orderID string) *GetUnfilledOrdersService {
	s.params["orderId"] = orderID
	return s
}

// SetTpslType selects normal orders or SPOT TP/SL orders (default normal).
func (s *GetUnfilledOrdersService) SetTpslType(tpslType TpslType) *GetUnfilledOrdersService {
	s.params["tpslType"] = string(tpslType)
	return s
}

// SetRequestTime sets the client request time.
func (s *GetUnfilledOrdersService) SetRequestTime(t time.Time) *GetUnfilledOrdersService {
	s.params["requestTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetReceiveWindow sets the valid time window in milliseconds.
func (s *GetUnfilledOrdersService) SetReceiveWindow(ms int64) *GetUnfilledOrdersService {
	s.params["receiveWindow"] = strconv.FormatInt(ms, 10)
	return s
}

func (s *GetUnfilledOrdersService) Do(ctx context.Context) ([]UnfilledOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/unfilled-orders", s.params).WithSign()
	resp, err := request.Do[[]UnfilledOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// UnfilledOrder is one current (open) order.
type UnfilledOrder struct {
	UserID                 string          `json:"userId"`
	Symbol                 string          `json:"symbol"`
	OrderID                string          `json:"orderId"`
	ClientOid              string          `json:"clientOid"`
	PriceAvg               decimal.Decimal `json:"priceAvg"`
	Size                   decimal.Decimal `json:"size"`
	OrderType              OrderType       `json:"orderType"`
	Side                   Side            `json:"side"`
	Status                 OrderStatus     `json:"status"`
	Force                  Force           `json:"force"`
	BasePrice              decimal.Decimal `json:"basePrice"`
	BaseVolume             decimal.Decimal `json:"baseVolume"`
	QuoteVolume            decimal.Decimal `json:"quoteVolume"`
	EnterPointSource       string          `json:"enterPointSource"`
	OrderSource            string          `json:"orderSource"`
	PresetTakeProfitPrice  decimal.Decimal `json:"presetTakeProfitPrice"`
	ExecuteTakeProfitPrice decimal.Decimal `json:"executeTakeProfitPrice"`
	PresetStopLossPrice    decimal.Decimal `json:"presetStopLossPrice"`
	ExecuteStopLossPrice   decimal.Decimal `json:"executeStopLossPrice"`
	TriggerPrice           decimal.Decimal `json:"triggerPrice"`
	TpslType               TpslType        `json:"tpslType"`
	CTime                  time.Time       `json:"cTime"`
	UTime                  time.Time       `json:"uTime"`
}

// GetHistoryOrdersService -- GET /api/v2/spot/trade/history-orders (private)
//
// Returns the account's historical (completed / cancelled) orders.
type GetHistoryOrdersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetHistoryOrdersService() *GetHistoryOrdersService {
	return &GetHistoryOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single trading pair.
func (s *GetHistoryOrdersService) SetSymbol(symbol string) *GetHistoryOrdersService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters orders created at or after t.
func (s *GetHistoryOrdersService) SetStartTime(t time.Time) *GetHistoryOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t.
func (s *GetHistoryOrdersService) SetEndTime(t time.Time) *GetHistoryOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backwards: return orders with orderId older than this.
func (s *GetHistoryOrdersService) SetIDLessThan(idLessThan string) *GetHistoryOrdersService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of orders returned (default 100, max 100).
func (s *GetHistoryOrdersService) SetLimit(limit int) *GetHistoryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetOrderID filters to a single order id.
func (s *GetHistoryOrdersService) SetOrderID(orderID string) *GetHistoryOrdersService {
	s.params["orderId"] = orderID
	return s
}

// SetTpslType selects normal orders or SPOT TP/SL orders (default normal).
func (s *GetHistoryOrdersService) SetTpslType(tpslType TpslType) *GetHistoryOrdersService {
	s.params["tpslType"] = string(tpslType)
	return s
}

// SetRequestTime sets the client request time.
func (s *GetHistoryOrdersService) SetRequestTime(t time.Time) *GetHistoryOrdersService {
	s.params["requestTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetReceiveWindow sets the valid time window in milliseconds.
func (s *GetHistoryOrdersService) SetReceiveWindow(ms int64) *GetHistoryOrdersService {
	s.params["receiveWindow"] = strconv.FormatInt(ms, 10)
	return s
}

func (s *GetHistoryOrdersService) Do(ctx context.Context) ([]HistoryOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/history-orders", s.params).WithSign()
	resp, err := request.Do[[]HistoryOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// HistoryOrder is one historical order. feeDetail is a Bitget-encoded JSON
// string (not a nested object), so it is kept as a raw string.
type HistoryOrder struct {
	UserID           string          `json:"userId"`
	Symbol           string          `json:"symbol"`
	OrderID          string          `json:"orderId"`
	ClientOid        string          `json:"clientOid"`
	Price            decimal.Decimal `json:"price"`
	Size             decimal.Decimal `json:"size"`
	OrderType        OrderType       `json:"orderType"`
	Side             Side            `json:"side"`
	Status           OrderStatus     `json:"status"`
	PriceAvg         decimal.Decimal `json:"priceAvg"`
	BaseVolume       decimal.Decimal `json:"baseVolume"`
	QuoteVolume      decimal.Decimal `json:"quoteVolume"`
	BaseCoin         string          `json:"baseCoin"`
	QuoteCoin        string          `json:"quoteCoin"`
	EnterPointSource string          `json:"enterPointSource"`
	OrderSource      string          `json:"orderSource"`
	FeeDetail        string          `json:"feeDetail"` // JSON-encoded fee breakdown string
	TpslType         TpslType        `json:"tpslType"`
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`
	CancelReason     string          `json:"cancelReason"` // normal_cancel, stp_cancel, ""
	CTime            time.Time       `json:"cTime"`
	UTime            time.Time       `json:"uTime"`
}

// GetFillsService -- GET /api/v2/spot/trade/fills (private)
//
// Returns the account's trade fills (executed transaction details).
type GetFillsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetFillsService() *GetFillsService {
	return &GetFillsService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single trading pair.
func (s *GetFillsService) SetSymbol(symbol string) *GetFillsService {
	s.params["symbol"] = symbol
	return s
}

// SetOrderID filters to fills of a single order.
func (s *GetFillsService) SetOrderID(orderID string) *GetFillsService {
	s.params["orderId"] = orderID
	return s
}

// SetStartTime filters fills created at or after t.
func (s *GetFillsService) SetStartTime(t time.Time) *GetFillsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters fills created at or before t.
func (s *GetFillsService) SetEndTime(t time.Time) *GetFillsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of fills returned (default 100, max 100).
func (s *GetFillsService) SetLimit(limit int) *GetFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards: return fills with tradeId older than this.
func (s *GetFillsService) SetIDLessThan(idLessThan string) *GetFillsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetFillsService) Do(ctx context.Context) ([]Fill, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/fills", s.params).WithSign()
	resp, err := request.Do[[]Fill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Fill is one executed trade.
type Fill struct {
	UserID     string          `json:"userId"`
	Symbol     string          `json:"symbol"`
	OrderID    string          `json:"orderId"`
	TradeID    string          `json:"tradeId"`
	OrderType  OrderType       `json:"orderType"`
	Side       Side            `json:"side"`
	PriceAvg   decimal.Decimal `json:"priceAvg"`
	Size       decimal.Decimal `json:"size"`
	Amount     decimal.Decimal `json:"amount"`
	FeeDetail  FillFeeDetail   `json:"feeDetail"`
	TradeScope string          `json:"tradeScope"` // taker, maker
	CTime      time.Time       `json:"cTime"`
	UTime      time.Time       `json:"uTime"`
}

// FillFeeDetail is the per-fill fee breakdown returned (as an object) by the
// fills endpoint.
type FillFeeDetail struct {
	Deduction         string          `json:"deduction"`
	FeeCoin           string          `json:"feeCoin"`
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"`
	TotalFee          decimal.Decimal `json:"totalFee"`
}
