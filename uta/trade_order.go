package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceOrderService -- POST /api/v3/trade/place-order (UTA trade read & write)
//
// Submits a single order. category, symbol, qty, side and orderType are
// required; price is required for limit orders. The reply data carries the new
// order's identifiers.
type PlaceOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewPlaceOrderService(category Category, symbol string, qty decimal.Decimal, side Side, orderType OrderType) *PlaceOrderService {
	return &PlaceOrderService{c: c, body: map[string]any{
		"category":  string(category),
		"symbol":    symbol,
		"qty":       qty.String(),
		"side":      string(side),
		"orderType": string(orderType),
	}}
}

// SetPrice sets the order price (required for limit orders).
func (s *PlaceOrderService) SetPrice(price decimal.Decimal) *PlaceOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTimeInForce sets the time-in-force policy (defaults to gtc for limit orders).
func (s *PlaceOrderService) SetTimeInForce(timeInForce TimeInForce) *PlaceOrderService {
	s.body["timeInForce"] = string(timeInForce)
	return s
}

// SetPosSide sets the position side (required for hedge-mode futures).
func (s *PlaceOrderService) SetPosSide(posSide PosSide) *PlaceOrderService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetClientOid sets the client-generated order identifier (1-32 chars).
func (s *PlaceOrderService) SetClientOrderID(clientOid string) *PlaceOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetReduceOnly sets the reduce-only flag (defaults to ReduceOnlyNo).
func (s *PlaceOrderService) SetReduceOnly(reduceOnly ReduceOnly) *PlaceOrderService {
	s.body["reduceOnly"] = string(reduceOnly)
	return s
}

// SetStpMode sets the self-trade prevention mode (none, cancel_taker,
// cancel_maker, cancel_both).
func (s *PlaceOrderService) SetStpMode(stpMode string) *PlaceOrderService {
	s.body["stpMode"] = stpMode
	return s
}

// SetMarginMode sets the margin mode (defaults to crossed).
func (s *PlaceOrderService) SetMarginMode(marginMode MarginMode) *PlaceOrderService {
	s.body["marginMode"] = string(marginMode)
	return s
}

// SetTpTriggerBy sets the take-profit trigger price type (market or mark).
func (s *PlaceOrderService) SetTpTriggerBy(tpTriggerBy string) *PlaceOrderService {
	s.body["tpTriggerBy"] = tpTriggerBy
	return s
}

// SetSlTriggerBy sets the stop-loss trigger price type (market or mark).
func (s *PlaceOrderService) SetSlTriggerBy(slTriggerBy string) *PlaceOrderService {
	s.body["slTriggerBy"] = slTriggerBy
	return s
}

// SetTakeProfit sets the preset take-profit trigger price.
func (s *PlaceOrderService) SetTakeProfit(takeProfit decimal.Decimal) *PlaceOrderService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the preset stop-loss trigger price.
func (s *PlaceOrderService) SetStopLoss(stopLoss decimal.Decimal) *PlaceOrderService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTpOrderType sets the take-profit order type (limit or market).
func (s *PlaceOrderService) SetTpOrderType(tpOrderType OrderType) *PlaceOrderService {
	s.body["tpOrderType"] = string(tpOrderType)
	return s
}

// SetSlOrderType sets the stop-loss order type (limit or market).
func (s *PlaceOrderService) SetSlOrderType(slOrderType OrderType) *PlaceOrderService {
	s.body["slOrderType"] = string(slOrderType)
	return s
}

// SetTpLimitPrice sets the take-profit execution price for limit orders.
func (s *PlaceOrderService) SetTpLimitPrice(tpLimitPrice decimal.Decimal) *PlaceOrderService {
	s.body["tpLimitPrice"] = tpLimitPrice.String()
	return s
}

// SetSlLimitPrice sets the stop-loss execution price for limit orders.
func (s *PlaceOrderService) SetSlLimitPrice(slLimitPrice decimal.Decimal) *PlaceOrderService {
	s.body["slLimitPrice"] = slLimitPrice.String()
	return s
}

func (s *PlaceOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/place-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// OrderRef identifies a single order by its exchange and client identifiers. It
// is the reply shape for place/modify/cancel single-order operations.
type OrderRef struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
}

// ModifyOrderService -- POST /api/v3/trade/modify-order (UTA trade read & write)
//
// Amends an open order's quantity and/or price. Identify the order by orderId
// or clientOid (orderId wins if both are set); at least one of qty or price
// must be provided.
type ModifyOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyOrderService() *ModifyOrderService {
	return &ModifyOrderService{c: c, body: map[string]any{}}
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

// SetQty sets the new order quantity (qty or price is required).
func (s *ModifyOrderService) SetQty(qty decimal.Decimal) *ModifyOrderService {
	s.body["qty"] = qty.String()
	return s
}

// SetPrice sets the new order price (qty or price is required).
func (s *ModifyOrderService) SetPrice(price decimal.Decimal) *ModifyOrderService {
	s.body["price"] = price.String()
	return s
}

// SetAutoCancel sets whether the original order is canceled when the
// modification fails ("yes" to cancel, "no" not to cancel; defaults to "no").
func (s *ModifyOrderService) SetAutoCancel(autoCancel string) *ModifyOrderService {
	s.body["autoCancel"] = autoCancel
	return s
}

// SetSymbol sets the symbol name (e.g. BTCUSDT).
func (s *ModifyOrderService) SetSymbol(symbol string) *ModifyOrderService {
	s.body["symbol"] = symbol
	return s
}

// SetCategory sets the product category.
func (s *ModifyOrderService) SetCategory(category Category) *ModifyOrderService {
	s.body["category"] = string(category)
	return s
}

func (s *ModifyOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/modify-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// CancelOrderService -- POST /api/v3/trade/cancel-order (UTA trade read & write)
//
// Cancels a single open order. Identify the order by orderId or clientOid
// (orderId wins if both are set).
type CancelOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelOrderService() *CancelOrderService {
	return &CancelOrderService{c: c, body: map[string]any{}}
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

// SetCategory sets the product category.
func (s *CancelOrderService) SetCategory(category Category) *CancelOrderService {
	s.body["category"] = string(category)
	return s
}

func (s *CancelOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/cancel-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// CancelSymbolOrderService -- POST /api/v3/trade/cancel-symbol-order (UTA trade read & write)
//
// Cancels all open orders in a category, optionally limited to a single symbol.
// The reply lists each attempted cancellation with its per-order result code.
type CancelSymbolOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelSymbolOrderService(category Category) *CancelSymbolOrderService {
	return &CancelSymbolOrderService{c: c, body: map[string]any{
		"category": string(category),
	}}
}

// SetSymbol limits the cancellation to a single trading pair (e.g. BTCUSDT).
// When omitted, all pending orders in the category are cancelled.
func (s *CancelSymbolOrderService) SetSymbol(symbol string) *CancelSymbolOrderService {
	s.body["symbol"] = symbol
	return s
}

func (s *CancelSymbolOrderService) Do(ctx context.Context) (*CancelSymbolOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/cancel-symbol-order", s.body).WithSign()
	return request.Do[CancelSymbolOrderResult](req)
}

// CancelSymbolOrderResult wraps the per-order cancellation results.
type CancelSymbolOrderResult struct {
	List []CancelResult `json:"list"`
}

// CancelResult is the outcome of cancelling a single order in a batch.
type CancelResult struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
}

// CountDownCancelAllService -- POST /api/v3/trade/countdown-cancel-all (UTA trade read & write)
//
// Arms a dead-man's-switch: if no further countdown request arrives within the
// reconnect window, all open orders are cancelled. countdown is the window in
// seconds (positive integer in [5, 60]); set 0 to disable. The reply data is
// the literal string "success".
type CountDownCancelAllService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCountDownCancelAllService(countdown string) *CountDownCancelAllService {
	return &CountDownCancelAllService{c: c, body: map[string]any{
		"countdown": countdown,
	}}
}

func (s *CountDownCancelAllService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/countdown-cancel-all", s.body).WithSign()
	return request.Do[string](req)
}
