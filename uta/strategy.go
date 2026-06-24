package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// StrategyType is the kind of strategy (plan) order.
type StrategyType string

const (
	StrategyTypeTPSL    StrategyType = "tpsl"
	StrategyTypeTrigger StrategyType = "trigger"
)

// TriggerBy is the price series a strategy order's trigger watches.
type TriggerBy string

const (
	TriggerByMarket TriggerBy = "market"
	TriggerByMark   TriggerBy = "mark"
)

// PlaceStrategyOrderService -- POST /api/v3/trade/place-strategy-order (UTA trade read & write)
//
// Places a take-profit/stop-loss ("tpsl") or trigger ("trigger") strategy order
// for a futures symbol. The reply data carries the assigned order identifiers.
type PlaceStrategyOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewPlaceStrategyOrderService(category Category, symbol string) *PlaceStrategyOrderService {
	return &PlaceStrategyOrderService{c: c, body: map[string]any{
		"category": string(category),
		"symbol":   symbol,
	}}
}

// SetClientOid sets the client order ID (6-hour idempotent validity).
func (s *PlaceStrategyOrderService) SetClientOrderID(clientOid string) *PlaceStrategyOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetType sets the strategy type (defaults to tpsl).
func (s *PlaceStrategyOrderService) SetType(strategyType StrategyType) *PlaceStrategyOrderService {
	s.body["type"] = string(strategyType)
	return s
}

// SetTpslMode sets the TPSL scope (e.g. "full" for all positions or "partial";
// defaults to full).
func (s *PlaceStrategyOrderService) SetTPSLMode(tpslMode string) *PlaceStrategyOrderService {
	s.body["tpslMode"] = tpslMode
	return s
}

// SetQty sets the order quantity in base coin (required for partial TPSL and
// trigger orders).
func (s *PlaceStrategyOrderService) SetQty(qty decimal.Decimal) *PlaceStrategyOrderService {
	s.body["qty"] = qty.String()
	return s
}

// SetSide sets the trade side.
func (s *PlaceStrategyOrderService) SetSide(side Side) *PlaceStrategyOrderService {
	s.body["side"] = string(side)
	return s
}

// SetPosSide sets the position side.
func (s *PlaceStrategyOrderService) SetPosSide(posSide PosSide) *PlaceStrategyOrderService {
	s.body["posSide"] = string(posSide)
	return s
}

// SetReduceOnly sets the reduce-only indicator ("yes" or "no").
func (s *PlaceStrategyOrderService) SetReduceOnly(reduceOnly string) *PlaceStrategyOrderService {
	s.body["reduceOnly"] = reduceOnly
	return s
}

// SetTpTriggerBy sets the take-profit trigger price type (defaults to market).
func (s *PlaceStrategyOrderService) SetTpTriggerBy(tpTriggerBy TriggerBy) *PlaceStrategyOrderService {
	s.body["tpTriggerBy"] = string(tpTriggerBy)
	return s
}

// SetSlTriggerBy sets the stop-loss trigger price type (defaults to market).
func (s *PlaceStrategyOrderService) SetSlTriggerBy(slTriggerBy TriggerBy) *PlaceStrategyOrderService {
	s.body["slTriggerBy"] = string(slTriggerBy)
	return s
}

// SetTakeProfit sets the take-profit trigger price.
func (s *PlaceStrategyOrderService) SetTakeProfit(takeProfit decimal.Decimal) *PlaceStrategyOrderService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the stop-loss trigger price.
func (s *PlaceStrategyOrderService) SetStopLoss(stopLoss decimal.Decimal) *PlaceStrategyOrderService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTpOrderType sets the take-profit order type (defaults to market).
func (s *PlaceStrategyOrderService) SetTpOrderType(tpOrderType OrderType) *PlaceStrategyOrderService {
	s.body["tpOrderType"] = string(tpOrderType)
	return s
}

// SetSlOrderType sets the stop-loss order type (defaults to market).
func (s *PlaceStrategyOrderService) SetSlOrderType(slOrderType OrderType) *PlaceStrategyOrderService {
	s.body["slOrderType"] = string(slOrderType)
	return s
}

// SetTpLimitPrice sets the take-profit execution price (limit orders only).
func (s *PlaceStrategyOrderService) SetTpLimitPrice(tpLimitPrice decimal.Decimal) *PlaceStrategyOrderService {
	s.body["tpLimitPrice"] = tpLimitPrice.String()
	return s
}

// SetSlLimitPrice sets the stop-loss execution price (limit orders only).
func (s *PlaceStrategyOrderService) SetSlLimitPrice(slLimitPrice decimal.Decimal) *PlaceStrategyOrderService {
	s.body["slLimitPrice"] = slLimitPrice.String()
	return s
}

// SetTriggerBy sets the trigger order price type (defaults to market).
func (s *PlaceStrategyOrderService) SetTriggerBy(triggerBy TriggerBy) *PlaceStrategyOrderService {
	s.body["triggerBy"] = string(triggerBy)
	return s
}

// SetTriggerPrice sets the trigger order trigger price.
func (s *PlaceStrategyOrderService) SetTriggerPrice(triggerPrice decimal.Decimal) *PlaceStrategyOrderService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetTriggerOrderType sets the trigger order type.
func (s *PlaceStrategyOrderService) SetTriggerOrderType(triggerOrderType OrderType) *PlaceStrategyOrderService {
	s.body["triggerOrderType"] = string(triggerOrderType)
	return s
}

// SetTriggerOrderPrice sets the trigger order execution price (limit orders only).
func (s *PlaceStrategyOrderService) SetTriggerOrderPrice(triggerOrderPrice decimal.Decimal) *PlaceStrategyOrderService {
	s.body["triggerOrderPrice"] = triggerOrderPrice.String()
	return s
}

func (s *PlaceStrategyOrderService) Do(ctx context.Context) (*StrategyOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/place-strategy-order", s.body).WithSign()
	return request.Do[StrategyOrderResult](req)
}

// StrategyOrderResult is the identifier pair returned by the place/modify
// strategy-order endpoints.
type StrategyOrderResult struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
}

// ModifyStrategyOrderService -- POST /api/v3/trade/modify-strategy-order (UTA trade read & write)
//
// Modifies an existing strategy order identified by orderId or clientOid. The
// reply data carries the order identifiers.
type ModifyStrategyOrderService struct {
	c    *UTAClient
	body map[string]any
}

// NewModifyStrategyOrderService starts a modify request. Identify the order with
// SetOrderId or SetClientOid (one is required; orderId takes priority); qty is
// required and set here.
func (c *UTAClient) NewModifyStrategyOrderService(qty decimal.Decimal) *ModifyStrategyOrderService {
	return &ModifyStrategyOrderService{c: c, body: map[string]any{
		"qty": qty.String(),
	}}
}

// SetOrderId sets the order ID (either orderId or clientOid is required; orderId
// takes priority).
func (s *ModifyStrategyOrderService) SetOrderID(orderId string) *ModifyStrategyOrderService {
	s.body["orderId"] = orderId
	return s
}

// SetClientOid sets the client order ID (either orderId or clientOid is required).
func (s *ModifyStrategyOrderService) SetClientOrderID(clientOid string) *ModifyStrategyOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetTpTriggerBy sets the take-profit trigger price type.
func (s *ModifyStrategyOrderService) SetTpTriggerBy(tpTriggerBy TriggerBy) *ModifyStrategyOrderService {
	s.body["tpTriggerBy"] = string(tpTriggerBy)
	return s
}

// SetSlTriggerBy sets the stop-loss trigger price type.
func (s *ModifyStrategyOrderService) SetSlTriggerBy(slTriggerBy TriggerBy) *ModifyStrategyOrderService {
	s.body["slTriggerBy"] = string(slTriggerBy)
	return s
}

// SetTakeProfit sets the take-profit trigger price.
func (s *ModifyStrategyOrderService) SetTakeProfit(takeProfit decimal.Decimal) *ModifyStrategyOrderService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the stop-loss trigger price.
func (s *ModifyStrategyOrderService) SetStopLoss(stopLoss decimal.Decimal) *ModifyStrategyOrderService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTpOrderType sets the take-profit order type.
func (s *ModifyStrategyOrderService) SetTpOrderType(tpOrderType OrderType) *ModifyStrategyOrderService {
	s.body["tpOrderType"] = string(tpOrderType)
	return s
}

// SetSlOrderType sets the stop-loss order type.
func (s *ModifyStrategyOrderService) SetSlOrderType(slOrderType OrderType) *ModifyStrategyOrderService {
	s.body["slOrderType"] = string(slOrderType)
	return s
}

// SetTpLimitPrice sets the take-profit execution price (limit orders only).
func (s *ModifyStrategyOrderService) SetTpLimitPrice(tpLimitPrice decimal.Decimal) *ModifyStrategyOrderService {
	s.body["tpLimitPrice"] = tpLimitPrice.String()
	return s
}

// SetSlLimitPrice sets the stop-loss execution price (limit orders only).
func (s *ModifyStrategyOrderService) SetSlLimitPrice(slLimitPrice decimal.Decimal) *ModifyStrategyOrderService {
	s.body["slLimitPrice"] = slLimitPrice.String()
	return s
}

// SetTriggerBy sets the trigger price type.
func (s *ModifyStrategyOrderService) SetTriggerBy(triggerBy TriggerBy) *ModifyStrategyOrderService {
	s.body["triggerBy"] = string(triggerBy)
	return s
}

// SetTriggerPrice sets the trigger price.
func (s *ModifyStrategyOrderService) SetTriggerPrice(triggerPrice decimal.Decimal) *ModifyStrategyOrderService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetTriggerOrderType sets the trigger order type.
func (s *ModifyStrategyOrderService) SetTriggerOrderType(triggerOrderType OrderType) *ModifyStrategyOrderService {
	s.body["triggerOrderType"] = string(triggerOrderType)
	return s
}

// SetTriggerOrderPrice sets the trigger execution price (limit orders only).
func (s *ModifyStrategyOrderService) SetTriggerOrderPrice(triggerOrderPrice decimal.Decimal) *ModifyStrategyOrderService {
	s.body["triggerOrderPrice"] = triggerOrderPrice.String()
	return s
}

func (s *ModifyStrategyOrderService) Do(ctx context.Context) (*StrategyOrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/modify-strategy-order", s.body).WithSign()
	return request.Do[StrategyOrderResult](req)
}

// CancelStrategyOrderService -- POST /api/v3/trade/cancel-strategy-order (UTA trade read & write)
//
// Cancels a strategy order identified by orderId or clientOid. The reply data is
// null on success.
type CancelStrategyOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelStrategyOrderService() *CancelStrategyOrderService {
	return &CancelStrategyOrderService{c: c, body: map[string]any{}}
}

// SetOrderId sets the order ID (either orderId or clientOid is required; orderId
// takes priority).
func (s *CancelStrategyOrderService) SetOrderID(orderId string) *CancelStrategyOrderService {
	s.body["orderId"] = orderId
	return s
}

// SetClientOid sets the client order ID (either orderId or clientOid is required).
func (s *CancelStrategyOrderService) SetClientOrderID(clientOid string) *CancelStrategyOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelStrategyOrderService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/cancel-strategy-order", s.body).WithSign()
	return request.Do[any](req)
}

// GetUnfilledStrategyOrdersService -- GET /api/v3/trade/unfilled-strategy-orders (UTA trade read)
//
// Returns the account's open (unfilled) strategy orders for a futures category,
// optionally filtered to a single strategy type.
type GetUnfilledStrategyOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetUnfilledStrategyOrdersService(category Category) *GetUnfilledStrategyOrdersService {
	return &GetUnfilledStrategyOrdersService{c: c, params: map[string]string{"category": string(category)}}
}

// SetType filters by strategy type (tpsl or trigger).
func (s *GetUnfilledStrategyOrdersService) SetType(strategyType StrategyType) *GetUnfilledStrategyOrdersService {
	s.params["type"] = string(strategyType)
	return s
}

func (s *GetUnfilledStrategyOrdersService) Do(ctx context.Context) ([]StrategyOrder, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/unfilled-strategy-orders", s.params).WithSign()
	resp, err := request.Do[[]StrategyOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// StrategyOrder is a single strategy (plan) order, shared by the unfilled and
// history listings.
type StrategyOrder struct {
	OrderID           string          `json:"orderId"`
	ClientOrderID     string          `json:"clientOid"`
	Category          Category        `json:"category"`
	Symbol            string          `json:"symbol"`
	Qty               decimal.Decimal `json:"qty"`
	PosSide           PosSide         `json:"posSide"`
	Status            string          `json:"status"` // pending, success, failed, cancelled, submitting
	TpTriggerBy       TriggerBy       `json:"tpTriggerBy"`
	SlTriggerBy       TriggerBy       `json:"slTriggerBy"`
	TakeProfit        decimal.Decimal `json:"takeProfit"`
	StopLoss          decimal.Decimal `json:"stopLoss"`
	TpOrderType       OrderType       `json:"tpOrderType"`
	SlOrderType       OrderType       `json:"slOrderType"`
	TpLimitPrice      decimal.Decimal `json:"tpLimitPrice"`
	SlLimitPrice      decimal.Decimal `json:"slLimitPrice"`
	TriggerBy         TriggerBy       `json:"triggerBy"`
	TriggerPrice      decimal.Decimal `json:"triggerPrice"`
	TriggerOrderType  OrderType       `json:"triggerOrderType"`
	TriggerOrderPrice decimal.Decimal `json:"triggerOrderPrice"`
	CreatedTime       time.Time       `json:"createdTime"`
	UpdatedTime       time.Time       `json:"updatedTime"`
}

// GetHistoryStrategyOrdersService -- GET /api/v3/trade/history-strategy-orders (UTA trade read)
//
// Returns the account's historical (completed/cancelled) strategy orders for a
// futures category, paginated by cursor.
type GetHistoryStrategyOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetHistoryStrategyOrdersService(category Category) *GetHistoryStrategyOrdersService {
	return &GetHistoryStrategyOrdersService{c: c, params: map[string]string{"category": string(category)}}
}

// SetType filters by strategy type (tpsl or trigger).
func (s *GetHistoryStrategyOrdersService) SetType(strategyType StrategyType) *GetHistoryStrategyOrdersService {
	s.params["type"] = string(strategyType)
	return s
}

// SetStartTime filters orders at or after t.
func (s *GetHistoryStrategyOrdersService) SetStartTime(t time.Time) *GetHistoryStrategyOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetHistoryStrategyOrdersService) SetEndTime(t time.Time) *GetHistoryStrategyOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetHistoryStrategyOrdersService) SetLimit(limit string) *GetHistoryStrategyOrdersService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetHistoryStrategyOrdersService) SetCursor(cursor string) *GetHistoryStrategyOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetHistoryStrategyOrdersService) Do(ctx context.Context) (*HistoryStrategyOrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/history-strategy-orders", s.params).WithSign()
	return request.Do[HistoryStrategyOrders](req)
}

type HistoryStrategyOrders struct {
	List   []StrategyOrder `json:"list"`
	Cursor string          `json:"cursor"`
}
