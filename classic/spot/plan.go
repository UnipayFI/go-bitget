package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlanOrderStatus is the lifecycle state of a spot plan (trigger) order. Current
// plan orders report not_trigger/executing; history plan orders report
// executed/fail_execute/cancelled.
type PlanOrderStatus string

const (
	PlanOrderStatusNotTrigger  PlanOrderStatus = "not_trigger"
	PlanOrderStatusExecuting   PlanOrderStatus = "executing"
	PlanOrderStatusExecuted    PlanOrderStatus = "executed"
	PlanOrderStatusFailExecute PlanOrderStatus = "fail_execute"
	PlanOrderStatusCancelled   PlanOrderStatus = "cancelled"
)

// PlanSubOrderStatus is the trigger status of the spot order spawned by a plan
// order.
type PlanSubOrderStatus string

const (
	PlanSubOrderStatusSuccess            PlanSubOrderStatus = "success"
	PlanSubOrderStatusFail               PlanSubOrderStatus = "fail"
	PlanSubOrderStatusCancelled          PlanSubOrderStatus = "cancelled"
	PlanSubOrderStatusInProgress         PlanSubOrderStatus = "in_progress"
	PlanSubOrderStatusInProgressTracking PlanSubOrderStatus = "in_progress_tracking"
)

// OrderSource identifies where an order was entered from.
type OrderSource string

const (
	OrderSourceWeb     OrderSource = "WEB"
	OrderSourceAPI     OrderSource = "API"
	OrderSourceSys     OrderSource = "SYS"
	OrderSourceAndroid OrderSource = "ANDROID"
	OrderSourceIOS     OrderSource = "IOS"
)

// PlacePlanOrderService -- POST /api/v2/spot/trade/place-plan-order (spot trade)
//
// Places a spot plan (trigger) order that fires once triggerPrice is reached.
type PlacePlanOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewPlacePlanOrderService(symbol string, side Side, triggerPrice decimal.Decimal, orderType OrderType, size decimal.Decimal, triggerType TriggerType) *PlacePlanOrderService {
	return &PlacePlanOrderService{c: c, body: map[string]any{
		"symbol":       symbol,
		"side":         string(side),
		"triggerPrice": triggerPrice.String(),
		"orderType":    string(orderType),
		"size":         size.String(),
		"triggerType":  string(triggerType),
	}}
}

// SetExecutePrice sets the limit price (required when orderType=limit).
func (s *PlacePlanOrderService) SetExecutePrice(executePrice decimal.Decimal) *PlacePlanOrderService {
	s.body["executePrice"] = executePrice.String()
	return s
}

// SetPlanType selects how size is denominated (amount=base coin, total=quote
// coin; default amount).
func (s *PlacePlanOrderService) SetPlanType(planType PlanType) *PlacePlanOrderService {
	s.body["planType"] = string(planType)
	return s
}

// SetClientOid sets a user-defined order identifier.
func (s *PlacePlanOrderService) SetClientOid(clientOid string) *PlacePlanOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetStpMode sets the self-trade prevention mode (default none).
func (s *PlacePlanOrderService) SetStpMode(mode SelfTradePreventionMode) *PlacePlanOrderService {
	s.body["stpMode"] = string(mode)
	return s
}

func (s *PlacePlanOrderService) Do(ctx context.Context) (*PlacePlanOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/place-plan-order", s.body).WithSign()
	return request.Do[PlacePlanOrderResponse](req)
}

// PlacePlanOrderResponse is the result of placing a plan order.
type PlacePlanOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// ModifyPlanOrderService -- POST /api/v2/spot/trade/modify-plan-order (spot trade)
//
// Modifies an existing spot plan order. Identify the order by orderId or
// clientOid (set at least one).
type ModifyPlanOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewModifyPlanOrderService(triggerPrice decimal.Decimal, orderType OrderType, size decimal.Decimal) *ModifyPlanOrderService {
	return &ModifyPlanOrderService{c: c, body: map[string]any{
		"triggerPrice": triggerPrice.String(),
		"orderType":    string(orderType),
		"size":         size.String(),
	}}
}

// SetOrderID identifies the plan order to modify by order id.
func (s *ModifyPlanOrderService) SetOrderID(orderID string) *ModifyPlanOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the plan order to modify by client order id.
func (s *ModifyPlanOrderService) SetClientOid(clientOid string) *ModifyPlanOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetExecutePrice sets the limit price (required when orderType=limit).
func (s *ModifyPlanOrderService) SetExecutePrice(executePrice decimal.Decimal) *ModifyPlanOrderService {
	s.body["executePrice"] = executePrice.String()
	return s
}

func (s *ModifyPlanOrderService) Do(ctx context.Context) (*ModifyPlanOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/modify-plan-order", s.body).WithSign()
	return request.Do[ModifyPlanOrderResponse](req)
}

// ModifyPlanOrderResponse is the result of modifying a plan order.
type ModifyPlanOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CancelPlanOrderService -- POST /api/v2/spot/trade/cancel-plan-order (spot trade)
//
// Cancels a single spot plan order. Identify the order by orderId or clientOid
// (set at least one).
type CancelPlanOrderService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCancelPlanOrderService() *CancelPlanOrderService {
	return &CancelPlanOrderService{c: c, body: map[string]any{}}
}

// SetOrderID identifies the plan order to cancel by order id.
func (s *CancelPlanOrderService) SetOrderID(orderID string) *CancelPlanOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the plan order to cancel by client order id.
func (s *CancelPlanOrderService) SetClientOid(clientOid string) *CancelPlanOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelPlanOrderService) Do(ctx context.Context) (*CancelPlanOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/cancel-plan-order", s.body).WithSign()
	return request.Do[CancelPlanOrderResponse](req)
}

// CancelPlanOrderResponse is the result of cancelling a plan order.
type CancelPlanOrderResponse struct {
	Result string `json:"result"` // success or failure
}

// BatchCancelPlanOrderService -- POST /api/v2/spot/trade/batch-cancel-plan-order (spot trade)
//
// Cancels spot plan orders in batch, optionally filtered to a set of symbols.
// When no symbols are set, all spot plan orders are cancelled.
type BatchCancelPlanOrderService struct {
	c          *SpotClient
	symbolList []string
}

func (c *SpotClient) NewBatchCancelPlanOrderService() *BatchCancelPlanOrderService {
	return &BatchCancelPlanOrderService{c: c}
}

// SetSymbolList restricts the cancellation to the given trading pairs.
func (s *BatchCancelPlanOrderService) SetSymbolList(symbolList []string) *BatchCancelPlanOrderService {
	s.symbolList = symbolList
	return s
}

func (s *BatchCancelPlanOrderService) Do(ctx context.Context) (*BatchCancelPlanOrderResponse, error) {
	body := map[string]any{}
	if len(s.symbolList) > 0 {
		body["symbolList"] = s.symbolList
	}
	req := request.Post(ctx, s.c, "/api/v2/spot/trade/batch-cancel-plan-order").SetBody(body).WithSign()
	return request.Do[BatchCancelPlanOrderResponse](req)
}

// BatchCancelPlanOrderResponse holds the per-order cancellation outcomes.
type BatchCancelPlanOrderResponse struct {
	SuccessList []BatchCancelPlanOrderSuccess `json:"successList"`
	FailureList []BatchCancelPlanOrderFailure `json:"failureList"`
}

// BatchCancelPlanOrderSuccess is a successfully cancelled plan order.
type BatchCancelPlanOrderSuccess struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// BatchCancelPlanOrderFailure is a plan order that failed to cancel.
type BatchCancelPlanOrderFailure struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
}

// GetCurrentPlanOrderService -- GET /api/v2/spot/trade/current-plan-order (spot trade)
//
// Returns the account's open (untriggered) spot plan orders.
type GetCurrentPlanOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetCurrentPlanOrderService() *GetCurrentPlanOrderService {
	return &GetCurrentPlanOrderService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single trading pair.
func (s *GetCurrentPlanOrderService) SetSymbol(symbol string) *GetCurrentPlanOrderService {
	s.params["symbol"] = symbol
	return s
}

// SetLimit caps the number of orders returned (default 20, max 100).
func (s *GetCurrentPlanOrderService) SetLimit(limit int) *GetCurrentPlanOrderService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan returns orders with an orderId older than the given cursor.
func (s *GetCurrentPlanOrderService) SetIDLessThan(idLessThan string) *GetCurrentPlanOrderService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters orders created at or after t.
func (s *GetCurrentPlanOrderService) SetStartTime(t time.Time) *GetCurrentPlanOrderService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t.
func (s *GetCurrentPlanOrderService) SetEndTime(t time.Time) *GetCurrentPlanOrderService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetCurrentPlanOrderService) Do(ctx context.Context) (*PlanOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/current-plan-order", s.params).WithSign()
	return request.Do[PlanOrderList](req)
}

// PlanOrderList is a page of plan orders with the pagination cursor.
type PlanOrderList struct {
	NextFlag   bool        `json:"nextFlag"`
	IDLessThan string      `json:"idLessThan"`
	OrderList  []PlanOrder `json:"orderList"`
}

// PlanOrder is a single spot plan (trigger) order.
type PlanOrder struct {
	OrderID          string          `json:"orderId"`
	ClientOid        string          `json:"clientOid"`
	Symbol           string          `json:"symbol"`
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`
	OrderType        OrderType       `json:"orderType"`
	ExecutePrice     decimal.Decimal `json:"executePrice"`
	PlanType         PlanType        `json:"planType"`
	Size             decimal.Decimal `json:"size"`
	Status           PlanOrderStatus `json:"status"`
	Side             Side            `json:"side"`
	TriggerType      TriggerType     `json:"triggerType"`
	EnterPointSource OrderSource     `json:"enterPointSource"`
	CTime            time.Time       `json:"cTime"`
	UTime            time.Time       `json:"uTime"`
}

// GetPlanSubOrderService -- GET /api/v2/spot/trade/plan-sub-order (spot trade)
//
// Returns the spot order(s) spawned by the given plan order once triggered.
type GetPlanSubOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetPlanSubOrderService(planOrderID string) *GetPlanSubOrderService {
	return &GetPlanSubOrderService{c: c, params: map[string]string{"planOrderId": planOrderID}}
}

func (s *GetPlanSubOrderService) Do(ctx context.Context) ([]PlanSubOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/plan-sub-order", s.params).WithSign()
	resp, err := request.Do[[]PlanSubOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PlanSubOrder is a spot order triggered by a plan order.
type PlanSubOrder struct {
	OrderID string             `json:"orderId"`
	Price   decimal.Decimal    `json:"price"`
	Type    OrderType          `json:"type"`
	Status  PlanSubOrderStatus `json:"status"`
}

// GetHistoryPlanOrderService -- GET /api/v2/spot/trade/history-plan-order (spot trade)
//
// Returns the account's triggered/cancelled spot plan orders (history).
type GetHistoryPlanOrderService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetHistoryPlanOrderService() *GetHistoryPlanOrderService {
	return &GetHistoryPlanOrderService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single trading pair.
func (s *GetHistoryPlanOrderService) SetSymbol(symbol string) *GetHistoryPlanOrderService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters orders created at or after t (within 90 days of endTime).
func (s *GetHistoryPlanOrderService) SetStartTime(t time.Time) *GetHistoryPlanOrderService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t (within 90 days of startTime).
func (s *GetHistoryPlanOrderService) SetEndTime(t time.Time) *GetHistoryPlanOrderService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan returns orders with an orderId older than the given cursor.
func (s *GetHistoryPlanOrderService) SetIDLessThan(idLessThan string) *GetHistoryPlanOrderService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the number of orders returned (default 100, max 100).
func (s *GetHistoryPlanOrderService) SetLimit(limit int) *GetHistoryPlanOrderService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryPlanOrderService) Do(ctx context.Context) (*PlanOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/trade/history-plan-order", s.params).WithSign()
	return request.Do[PlanOrderList](req)
}
