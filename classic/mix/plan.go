package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlanOrderStatus is the lifecycle state of a futures trigger (plan) order.
// Pending trigger orders report live; history trigger orders report
// executed/fail_execute/cancelled.
type PlanOrderStatus string

const (
	PlanOrderStatusLive        PlanOrderStatus = "live"
	PlanOrderStatusExecuted    PlanOrderStatus = "executed"
	PlanOrderStatusFailExecute PlanOrderStatus = "fail_execute"
	PlanOrderStatusCancelled   PlanOrderStatus = "cancelled"
)

// PlanSubOrderStatus is the trigger status of the futures order spawned by a
// trigger (plan) order.
type PlanSubOrderStatus string

const (
	PlanSubOrderStatusSuccess            PlanSubOrderStatus = "success"
	PlanSubOrderStatusFail               PlanSubOrderStatus = "fail"
	PlanSubOrderStatusCancelled          PlanSubOrderStatus = "cancelled"
	PlanSubOrderStatusInProgress         PlanSubOrderStatus = "in_progress"
	PlanSubOrderStatusInProgressTracking PlanSubOrderStatus = "in_progress_tracking"
)

// PosSide is the position direction reported on trigger orders. Beyond the
// long/short of a hedge-mode position it can be net for one-way mode.
type PosSide string

const (
	PosSideLong  PosSide = "long"
	PosSideShort PosSide = "short"
	PosSideNet   PosSide = "net"
)

// OrderSource identifies where a trigger order was entered from.
type OrderSource string

const (
	OrderSourceWeb     OrderSource = "WEB"
	OrderSourceAPI     OrderSource = "API"
	OrderSourceSys     OrderSource = "SYS"
	OrderSourceAndroid OrderSource = "ANDROID"
	OrderSourceIOS     OrderSource = "IOS"
)

// PlaceTpslOrderService -- POST /api/v2/mix/order/place-tpsl-order (mix trade)
//
// Places a stop-profit/stop-loss plan order against a futures position.
type PlaceTpslOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewPlaceTpslOrderService(productType ProductType, symbol, marginCoin string, planType PlanType, triggerPrice decimal.Decimal, holdSide HoldSide) *PlaceTpslOrderService {
	return &PlaceTpslOrderService{c: c, body: map[string]any{
		"productType":  string(productType),
		"symbol":       symbol,
		"marginCoin":   marginCoin,
		"planType":     string(planType),
		"triggerPrice": triggerPrice.String(),
		"holdSide":     string(holdSide),
	}}
}

// SetTriggerType selects the price series the trigger is measured against
// (default fill_price).
func (s *PlaceTpslOrderService) SetTriggerType(triggerType TriggerType) *PlaceTpslOrderService {
	s.body["triggerType"] = string(triggerType)
	return s
}

// SetExecutePrice sets the execution price (0 or empty = market, >0 = limit;
// omit for moving_plan).
func (s *PlaceTpslOrderService) SetExecutePrice(executePrice decimal.Decimal) *PlaceTpslOrderService {
	s.body["executePrice"] = executePrice.String()
	return s
}

// SetSize sets the order quantity in base coin (required for
// profit/loss/moving plans, omit for whole-position plans).
func (s *PlaceTpslOrderService) SetSize(size decimal.Decimal) *PlaceTpslOrderService {
	s.body["size"] = size.String()
	return s
}

// SetRangeRate sets the callback range (required only for moving_plan).
func (s *PlaceTpslOrderService) SetRangeRate(rangeRate decimal.Decimal) *PlaceTpslOrderService {
	s.body["rangeRate"] = rangeRate.String()
	return s
}

// SetClientOid sets a user-defined order identifier.
func (s *PlaceTpslOrderService) SetClientOid(clientOid string) *PlaceTpslOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetStpMode sets the self-trade prevention mode (default none).
func (s *PlaceTpslOrderService) SetStpMode(mode SelfTradePreventionMode) *PlaceTpslOrderService {
	s.body["stpMode"] = string(mode)
	return s
}

func (s *PlaceTpslOrderService) Do(ctx context.Context) (*PlaceTpslOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/place-tpsl-order", s.body).WithSign()
	return request.Do[PlaceTpslOrderResponse](req)
}

// PlaceTpslOrderResponse is the result of placing a TP/SL plan order.
type PlaceTpslOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// PlacePosTpslService -- POST /api/v2/mix/order/place-pos-tpsl (mix trade)
//
// Places take-profit and stop-loss plan orders for a position in one request.
// Set at least one of the take-profit / stop-loss trigger prices.
type PlacePosTpslService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewPlacePosTpslService(productType ProductType, symbol, marginCoin string, holdSide HoldSide) *PlacePosTpslService {
	return &PlacePosTpslService{c: c, body: map[string]any{
		"productType": string(productType),
		"symbol":      symbol,
		"marginCoin":  marginCoin,
		"holdSide":    string(holdSide),
	}}
}

// SetStopSurplusTriggerPrice sets the take-profit trigger price.
func (s *PlacePosTpslService) SetStopSurplusTriggerPrice(price decimal.Decimal) *PlacePosTpslService {
	s.body["stopSurplusTriggerPrice"] = price.String()
	return s
}

// SetStopSurplusSize sets the take-profit quantity in base coin (filled =>
// profit_plan, empty => pos_profit).
func (s *PlacePosTpslService) SetStopSurplusSize(size decimal.Decimal) *PlacePosTpslService {
	s.body["stopSurplusSize"] = size.String()
	return s
}

// SetStopSurplusTriggerType selects the take-profit trigger price series
// (default fill_price).
func (s *PlacePosTpslService) SetStopSurplusTriggerType(triggerType TriggerType) *PlacePosTpslService {
	s.body["stopSurplusTriggerType"] = string(triggerType)
	return s
}

// SetStopSurplusExecutePrice sets the take-profit execution price (0 or empty =
// market, >0 = limit).
func (s *PlacePosTpslService) SetStopSurplusExecutePrice(price decimal.Decimal) *PlacePosTpslService {
	s.body["stopSurplusExecutePrice"] = price.String()
	return s
}

// SetStopLossTriggerPrice sets the stop-loss trigger price.
func (s *PlacePosTpslService) SetStopLossTriggerPrice(price decimal.Decimal) *PlacePosTpslService {
	s.body["stopLossTriggerPrice"] = price.String()
	return s
}

// SetStopLossSize sets the stop-loss quantity in base coin (filled => loss_plan,
// empty => pos_loss).
func (s *PlacePosTpslService) SetStopLossSize(size decimal.Decimal) *PlacePosTpslService {
	s.body["stopLossSize"] = size.String()
	return s
}

// SetStopLossTriggerType selects the stop-loss trigger price series (default
// fill_price).
func (s *PlacePosTpslService) SetStopLossTriggerType(triggerType TriggerType) *PlacePosTpslService {
	s.body["stopLossTriggerType"] = string(triggerType)
	return s
}

// SetStopLossExecutePrice sets the stop-loss execution price (0 or empty =
// market, >0 = limit).
func (s *PlacePosTpslService) SetStopLossExecutePrice(price decimal.Decimal) *PlacePosTpslService {
	s.body["stopLossExecutePrice"] = price.String()
	return s
}

// SetStpMode sets the self-trade prevention mode (default none).
func (s *PlacePosTpslService) SetStpMode(mode SelfTradePreventionMode) *PlacePosTpslService {
	s.body["stpMode"] = string(mode)
	return s
}

// SetStopSurplusClientOid sets a user-defined identifier for the take-profit order.
func (s *PlacePosTpslService) SetStopSurplusClientOid(clientOid string) *PlacePosTpslService {
	s.body["stopSurplusClientOid"] = clientOid
	return s
}

// SetStopLossClientOid sets a user-defined identifier for the stop-loss order.
func (s *PlacePosTpslService) SetStopLossClientOid(clientOid string) *PlacePosTpslService {
	s.body["stopLossClientOid"] = clientOid
	return s
}

func (s *PlacePosTpslService) Do(ctx context.Context) ([]PlacePosTpslResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/place-pos-tpsl", s.body).WithSign()
	resp, err := request.Do[[]PlacePosTpslResponse](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PlacePosTpslResponse is one of the trigger orders created by a position
// TP/SL request.
type PlacePosTpslResponse struct {
	OrderID              string `json:"orderId"`
	StopSurplusClientOid string `json:"stopSurplusClientOid"`
	StopLossClientOid    string `json:"stopLossClientOid"`
}

// PlacePlanOrderService -- POST /api/v2/mix/order/place-plan-order (mix trade)
//
// Places a futures trigger (plan) order that fires once triggerPrice is reached.
type PlacePlanOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewPlacePlanOrderService(productType ProductType, symbol, marginCoin string, marginMode MarginMode, planType PlanType, size decimal.Decimal, triggerPrice decimal.Decimal, triggerType TriggerType, side Side, orderType OrderType) *PlacePlanOrderService {
	return &PlacePlanOrderService{c: c, body: map[string]any{
		"productType":  string(productType),
		"symbol":       symbol,
		"marginCoin":   marginCoin,
		"marginMode":   string(marginMode),
		"planType":     string(planType),
		"size":         size.String(),
		"triggerPrice": triggerPrice.String(),
		"triggerType":  string(triggerType),
		"side":         string(side),
		"orderType":    string(orderType),
	}}
}

// SetPrice sets the order price (required for limit orders in normal_plan;
// empty for market orders / track_plan).
func (s *PlacePlanOrderService) SetPrice(price decimal.Decimal) *PlacePlanOrderService {
	s.body["price"] = price.String()
	return s
}

// SetCallbackRatio sets the retracement percentage for trailing stops (track_plan).
func (s *PlacePlanOrderService) SetCallbackRatio(callbackRatio decimal.Decimal) *PlacePlanOrderService {
	s.body["callbackRatio"] = callbackRatio.String()
	return s
}

// SetTradeSide sets whether the order opens or closes a position (required in
// hedge mode).
func (s *PlacePlanOrderService) SetTradeSide(tradeSide TradeSide) *PlacePlanOrderService {
	s.body["tradeSide"] = string(tradeSide)
	return s
}

// SetClientOid sets a user-defined order identifier.
func (s *PlacePlanOrderService) SetClientOid(clientOid string) *PlacePlanOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetReduceOnly restricts the order to only reducing an existing position
// (one-way mode only).
func (s *PlacePlanOrderService) SetReduceOnly(reduceOnly ReduceOnly) *PlacePlanOrderService {
	s.body["reduceOnly"] = string(reduceOnly)
	return s
}

// SetStopSurplusTriggerPrice sets the attached take-profit trigger price.
func (s *PlacePlanOrderService) SetStopSurplusTriggerPrice(price decimal.Decimal) *PlacePlanOrderService {
	s.body["stopSurplusTriggerPrice"] = price.String()
	return s
}

// SetStopSurplusExecutePrice sets the attached take-profit execution price
// (empty triggers market execution for normal_plan).
func (s *PlacePlanOrderService) SetStopSurplusExecutePrice(price decimal.Decimal) *PlacePlanOrderService {
	s.body["stopSurplusExecutePrice"] = price.String()
	return s
}

// SetStopSurplusTriggerType selects the attached take-profit trigger price
// series (track_plan only accepts fill_price).
func (s *PlacePlanOrderService) SetStopSurplusTriggerType(triggerType TriggerType) *PlacePlanOrderService {
	s.body["stopSurplusTriggerType"] = string(triggerType)
	return s
}

// SetStopLossTriggerPrice sets the attached stop-loss trigger price.
func (s *PlacePlanOrderService) SetStopLossTriggerPrice(price decimal.Decimal) *PlacePlanOrderService {
	s.body["stopLossTriggerPrice"] = price.String()
	return s
}

// SetStopLossExecutePrice sets the attached stop-loss execution price (empty
// triggers market execution for normal_plan).
func (s *PlacePlanOrderService) SetStopLossExecutePrice(price decimal.Decimal) *PlacePlanOrderService {
	s.body["stopLossExecutePrice"] = price.String()
	return s
}

// SetStopLossTriggerType selects the attached stop-loss trigger price series
// (track_plan only accepts fill_price).
func (s *PlacePlanOrderService) SetStopLossTriggerType(triggerType TriggerType) *PlacePlanOrderService {
	s.body["stopLossTriggerType"] = string(triggerType)
	return s
}

// SetStpMode sets the self-trade prevention mode (default none).
func (s *PlacePlanOrderService) SetStpMode(mode SelfTradePreventionMode) *PlacePlanOrderService {
	s.body["stpMode"] = string(mode)
	return s
}

func (s *PlacePlanOrderService) Do(ctx context.Context) (*PlacePlanOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/place-plan-order", s.body).WithSign()
	return request.Do[PlacePlanOrderResponse](req)
}

// PlacePlanOrderResponse is the result of placing a trigger (plan) order.
type PlacePlanOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// ModifyTpslOrderService -- POST /api/v2/mix/order/modify-tpsl-order (mix trade)
//
// Modifies an existing TP/SL plan order. Identify the order by orderId or
// clientOid (set at least one).
type ModifyTpslOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewModifyTpslOrderService(productType ProductType, symbol, marginCoin string, triggerPrice decimal.Decimal, size decimal.Decimal) *ModifyTpslOrderService {
	return &ModifyTpslOrderService{c: c, body: map[string]any{
		"productType":  string(productType),
		"symbol":       symbol,
		"marginCoin":   marginCoin,
		"triggerPrice": triggerPrice.String(),
		"size":         size.String(),
	}}
}

// SetOrderID identifies the TP/SL order to modify by order id.
func (s *ModifyTpslOrderService) SetOrderID(orderID string) *ModifyTpslOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the TP/SL order to modify by client order id.
func (s *ModifyTpslOrderService) SetClientOid(clientOid string) *ModifyTpslOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetTriggerType selects the price series the trigger is measured against.
func (s *ModifyTpslOrderService) SetTriggerType(triggerType TriggerType) *ModifyTpslOrderService {
	s.body["triggerType"] = string(triggerType)
	return s
}

// SetExecutePrice sets the execution price (0 or empty = market, >0 = limit).
func (s *ModifyTpslOrderService) SetExecutePrice(executePrice decimal.Decimal) *ModifyTpslOrderService {
	s.body["executePrice"] = executePrice.String()
	return s
}

// SetRangeRate sets the callback range percentage.
func (s *ModifyTpslOrderService) SetRangeRate(rangeRate decimal.Decimal) *ModifyTpslOrderService {
	s.body["rangeRate"] = rangeRate.String()
	return s
}

func (s *ModifyTpslOrderService) Do(ctx context.Context) (*ModifyTpslOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/modify-tpsl-order", s.body).WithSign()
	return request.Do[ModifyTpslOrderResponse](req)
}

// ModifyTpslOrderResponse is the result of modifying a TP/SL plan order.
type ModifyTpslOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// ModifyPlanOrderService -- POST /api/v2/mix/order/modify-plan-order (mix trade)
//
// Modifies an existing trigger (plan) order. Identify the order by orderId or
// clientOid (set at least one); unset fields keep their current value.
type ModifyPlanOrderService struct {
	c    *MixClient
	body map[string]any
}

func (c *MixClient) NewModifyPlanOrderService(productType ProductType) *ModifyPlanOrderService {
	return &ModifyPlanOrderService{c: c, body: map[string]any{
		"productType": string(productType),
	}}
}

// SetOrderID identifies the trigger order to modify by order id.
func (s *ModifyPlanOrderService) SetOrderID(orderID string) *ModifyPlanOrderService {
	s.body["orderId"] = orderID
	return s
}

// SetClientOid identifies the trigger order to modify by client order id.
func (s *ModifyPlanOrderService) SetClientOid(clientOid string) *ModifyPlanOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetNewSize updates the order quantity.
func (s *ModifyPlanOrderService) SetNewSize(newSize decimal.Decimal) *ModifyPlanOrderService {
	s.body["newSize"] = newSize.String()
	return s
}

// SetNewPrice updates the execution price (empty for market / trailing orders).
func (s *ModifyPlanOrderService) SetNewPrice(newPrice decimal.Decimal) *ModifyPlanOrderService {
	s.body["newPrice"] = newPrice.String()
	return s
}

// SetNewCallbackRatio updates the callback rate for trailing stops (max 10).
func (s *ModifyPlanOrderService) SetNewCallbackRatio(newCallbackRatio decimal.Decimal) *ModifyPlanOrderService {
	s.body["newCallbackRatio"] = newCallbackRatio.String()
	return s
}

// SetNewTriggerPrice updates the trigger price.
func (s *ModifyPlanOrderService) SetNewTriggerPrice(newTriggerPrice decimal.Decimal) *ModifyPlanOrderService {
	s.body["newTriggerPrice"] = newTriggerPrice.String()
	return s
}

// SetNewTriggerType updates the trigger price series (requires newTriggerPrice).
func (s *ModifyPlanOrderService) SetNewTriggerType(newTriggerType TriggerType) *ModifyPlanOrderService {
	s.body["newTriggerType"] = string(newTriggerType)
	return s
}

// SetNewStopSurplusTriggerPrice updates the take-profit trigger price (0 removes it).
func (s *ModifyPlanOrderService) SetNewStopSurplusTriggerPrice(price decimal.Decimal) *ModifyPlanOrderService {
	s.body["newStopSurplusTriggerPrice"] = price.String()
	return s
}

// SetNewStopSurplusExecutePrice updates the take-profit execution price (0 removes it).
func (s *ModifyPlanOrderService) SetNewStopSurplusExecutePrice(price decimal.Decimal) *ModifyPlanOrderService {
	s.body["newStopSurplusExecutePrice"] = price.String()
	return s
}

// SetNewStopSurplusTriggerType updates the take-profit trigger series (requires
// newStopSurplusTriggerPrice).
func (s *ModifyPlanOrderService) SetNewStopSurplusTriggerType(triggerType TriggerType) *ModifyPlanOrderService {
	s.body["newStopSurplusTriggerType"] = string(triggerType)
	return s
}

// SetNewStopLossTriggerPrice updates the stop-loss trigger price (0 removes it).
func (s *ModifyPlanOrderService) SetNewStopLossTriggerPrice(price decimal.Decimal) *ModifyPlanOrderService {
	s.body["newStopLossTriggerPrice"] = price.String()
	return s
}

// SetNewStopLossExecutePrice updates the stop-loss execution price (0 removes it).
func (s *ModifyPlanOrderService) SetNewStopLossExecutePrice(price decimal.Decimal) *ModifyPlanOrderService {
	s.body["newStopLossExecutePrice"] = price.String()
	return s
}

// SetNewStopLossTriggerType updates the stop-loss trigger series (requires
// newStopLossTriggerPrice).
func (s *ModifyPlanOrderService) SetNewStopLossTriggerType(triggerType TriggerType) *ModifyPlanOrderService {
	s.body["newStopLossTriggerType"] = string(triggerType)
	return s
}

func (s *ModifyPlanOrderService) Do(ctx context.Context) (*ModifyPlanOrderResponse, error) {
	req := request.Post(ctx, s.c, "/api/v2/mix/order/modify-plan-order", s.body).WithSign()
	return request.Do[ModifyPlanOrderResponse](req)
}

// ModifyPlanOrderResponse is the result of modifying a trigger (plan) order.
type ModifyPlanOrderResponse struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CancelPlanOrderService -- POST /api/v2/mix/order/cancel-plan-order (mix trade)
//
// Cancels trigger (plan) orders. Set an explicit orderIdList to cancel specific
// orders (symbol then required), or leave it empty and use SetSymbol/SetPlanType
// to cancel all matching trigger orders.
type CancelPlanOrderService struct {
	c           *MixClient
	productType ProductType
	orderIDList []CancelPlanOrderID
	symbol      string
	marginCoin  string
	planType    PlanType
}

// CancelPlanOrderID identifies a single trigger order to cancel by order id or
// client order id (set at least one).
type CancelPlanOrderID struct {
	OrderID   string `json:"orderId,omitempty"`
	ClientOid string `json:"clientOid,omitempty"`
}

func (c *MixClient) NewCancelPlanOrderService(productType ProductType) *CancelPlanOrderService {
	return &CancelPlanOrderService{c: c, productType: productType}
}

// SetOrderIDList sets the explicit list of trigger orders to cancel.
func (s *CancelPlanOrderService) SetOrderIDList(orderIDList []CancelPlanOrderID) *CancelPlanOrderService {
	s.orderIDList = orderIDList
	return s
}

// SetSymbol restricts the cancellation to a single trading pair (required when
// orderIdList is provided).
func (s *CancelPlanOrderService) SetSymbol(symbol string) *CancelPlanOrderService {
	s.symbol = symbol
	return s
}

// SetMarginCoin sets the margin coin (capitalized).
func (s *CancelPlanOrderService) SetMarginCoin(marginCoin string) *CancelPlanOrderService {
	s.marginCoin = marginCoin
	return s
}

// SetPlanType restricts the cancellation to a single trigger order category.
func (s *CancelPlanOrderService) SetPlanType(planType PlanType) *CancelPlanOrderService {
	s.planType = planType
	return s
}

func (s *CancelPlanOrderService) Do(ctx context.Context) (*CancelPlanOrderResponse, error) {
	body := map[string]any{"productType": string(s.productType)}
	if len(s.orderIDList) > 0 {
		body["orderIdList"] = s.orderIDList
	}
	if s.symbol != "" {
		body["symbol"] = s.symbol
	}
	if s.marginCoin != "" {
		body["marginCoin"] = s.marginCoin
	}
	if s.planType != "" {
		body["planType"] = string(s.planType)
	}
	req := request.Post(ctx, s.c, "/api/v2/mix/order/cancel-plan-order").SetBody(body).WithSign()
	return request.Do[CancelPlanOrderResponse](req)
}

// CancelPlanOrderResponse holds the per-order cancellation outcomes.
type CancelPlanOrderResponse struct {
	SuccessList []CancelPlanOrderSuccess `json:"successList"`
	FailureList []CancelPlanOrderFailure `json:"failureList"`
}

// CancelPlanOrderSuccess is a successfully cancelled trigger order.
type CancelPlanOrderSuccess struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CancelPlanOrderFailure is a trigger order that failed to cancel.
type CancelPlanOrderFailure struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	ErrorMsg  string `json:"errorMsg"`
}

// GetOrdersPlanPendingService -- GET /api/v2/mix/order/orders-plan-pending (mix trade)
//
// Returns the account's open (untriggered) futures trigger orders for a product
// type and plan category.
type GetOrdersPlanPendingService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrdersPlanPendingService(productType ProductType, planType PlanType) *GetOrdersPlanPendingService {
	return &GetOrdersPlanPendingService{c: c, params: map[string]string{
		"productType": string(productType),
		"planType":    string(planType),
	}}
}

// SetOrderID filters to a single trigger order by order id.
func (s *GetOrdersPlanPendingService) SetOrderID(orderID string) *GetOrdersPlanPendingService {
	s.params["orderId"] = orderID
	return s
}

// SetClientOid filters to a single trigger order by client order id.
func (s *GetOrdersPlanPendingService) SetClientOid(clientOid string) *GetOrdersPlanPendingService {
	s.params["clientOid"] = clientOid
	return s
}

// SetSymbol filters to a single trading pair.
func (s *GetOrdersPlanPendingService) SetSymbol(symbol string) *GetOrdersPlanPendingService {
	s.params["symbol"] = symbol
	return s
}

// SetIDLessThan returns orders with an orderId older than the given cursor.
func (s *GetOrdersPlanPendingService) SetIDLessThan(idLessThan string) *GetOrdersPlanPendingService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters orders created at or after t.
func (s *GetOrdersPlanPendingService) SetStartTime(t time.Time) *GetOrdersPlanPendingService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t.
func (s *GetOrdersPlanPendingService) SetEndTime(t time.Time) *GetOrdersPlanPendingService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders returned (default 100, max 100).
func (s *GetOrdersPlanPendingService) SetLimit(limit int) *GetOrdersPlanPendingService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersPlanPendingService) Do(ctx context.Context) (*PlanOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/orders-plan-pending", s.params).WithSign()
	return request.Do[PlanOrderList](req)
}

// GetOrdersPlanHistoryService -- GET /api/v2/mix/order/orders-plan-history (mix trade)
//
// Returns the account's triggered/cancelled futures trigger orders (history) for
// a product type and plan category (max 3-month window).
type GetOrdersPlanHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOrdersPlanHistoryService(productType ProductType, planType PlanType) *GetOrdersPlanHistoryService {
	return &GetOrdersPlanHistoryService{c: c, params: map[string]string{
		"productType": string(productType),
		"planType":    string(planType),
	}}
}

// SetOrderID filters to a single trigger order by order id.
func (s *GetOrdersPlanHistoryService) SetOrderID(orderID string) *GetOrdersPlanHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetClientOid filters to a single trigger order by client order id.
func (s *GetOrdersPlanHistoryService) SetClientOid(clientOid string) *GetOrdersPlanHistoryService {
	s.params["clientOid"] = clientOid
	return s
}

// SetPlanStatus filters to a single trigger order status (executed,
// fail_execute, cancelled).
func (s *GetOrdersPlanHistoryService) SetPlanStatus(planStatus PlanOrderStatus) *GetOrdersPlanHistoryService {
	s.params["planStatus"] = string(planStatus)
	return s
}

// SetSymbol filters to a single trading pair.
func (s *GetOrdersPlanHistoryService) SetSymbol(symbol string) *GetOrdersPlanHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetIDLessThan returns orders with an orderId older than the given cursor.
func (s *GetOrdersPlanHistoryService) SetIDLessThan(idLessThan string) *GetOrdersPlanHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetStartTime filters orders created at or after t (within 3 months of endTime).
func (s *GetOrdersPlanHistoryService) SetStartTime(t time.Time) *GetOrdersPlanHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders created at or before t (within 3 months of startTime).
func (s *GetOrdersPlanHistoryService) SetEndTime(t time.Time) *GetOrdersPlanHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders returned (default 100, max 100).
func (s *GetOrdersPlanHistoryService) SetLimit(limit int) *GetOrdersPlanHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrdersPlanHistoryService) Do(ctx context.Context) (*PlanOrderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/orders-plan-history", s.params).WithSign()
	return request.Do[PlanOrderList](req)
}

// PlanOrderList is a page of trigger (plan) orders with the pagination cursor.
type PlanOrderList struct {
	EntrustedList []PlanOrder `json:"entrustedList"`
	EndID         string      `json:"endId"`
}

// PlanOrder is a single futures trigger (plan) order. It is the union of the
// pending and history shapes; history-only fields (executeOrderId, baseVolume)
// are empty on pending orders and the orderSource field is empty on history
// orders.
type PlanOrder struct {
	PlanType                PlanType        `json:"planType"`
	Symbol                  string          `json:"symbol"`
	Size                    decimal.Decimal `json:"size"`
	OrderID                 string          `json:"orderId"`
	ExecuteOrderID          string          `json:"executeOrderId"`
	ClientOid               string          `json:"clientOid"`
	PlanStatus              PlanOrderStatus `json:"planStatus"`
	Price                   decimal.Decimal `json:"price"`
	ExecutePrice            decimal.Decimal `json:"executePrice"`
	BaseVolume              decimal.Decimal `json:"baseVolume"`
	CallbackRatio           decimal.Decimal `json:"callbackRatio"`
	TriggerPrice            decimal.Decimal `json:"triggerPrice"`
	TriggerType             TriggerType     `json:"triggerType"`
	Side                    Side            `json:"side"`
	PosSide                 PosSide         `json:"posSide"`
	MarginCoin              string          `json:"marginCoin"`
	MarginMode              MarginMode      `json:"marginMode"`
	EnterPointSource        OrderSource     `json:"enterPointSource"`
	TradeSide               TradeSide       `json:"tradeSide"`
	PosMode                 PositionMode    `json:"posMode"`
	OrderType               OrderType       `json:"orderType"`
	OrderSource             string          `json:"orderSource"`
	CTime                   time.Time       `json:"cTime"`
	UTime                   time.Time       `json:"uTime"`
	StopSurplusExecutePrice decimal.Decimal `json:"stopSurplusExecutePrice"`
	StopSurplusTriggerPrice decimal.Decimal `json:"stopSurplusTriggerPrice"`
	StopSurplusTriggerType  TriggerType     `json:"stopSurplusTriggerType"`
	StopLossExecutePrice    decimal.Decimal `json:"stopLossExecutePrice"`
	StopLossTriggerPrice    decimal.Decimal `json:"stopLossTriggerPrice"`
	StopLossTriggerType     TriggerType     `json:"stopLossTriggerType"`
}

// GetPlanSubOrderService -- GET /api/v2/mix/order/plan-sub-order (mix trade)
//
// Returns the futures order(s) spawned by the given trigger (plan) order once it
// triggered.
type GetPlanSubOrderService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetPlanSubOrderService(productType ProductType, planType PlanType, planOrderID string) *GetPlanSubOrderService {
	return &GetPlanSubOrderService{c: c, params: map[string]string{
		"productType": string(productType),
		"planType":    string(planType),
		"planOrderId": planOrderID,
	}}
}

func (s *GetPlanSubOrderService) Do(ctx context.Context) ([]PlanSubOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/order/plan-sub-order", s.params).WithSign()
	resp, err := request.Do[[]PlanSubOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PlanSubOrder is a futures order triggered by a trigger (plan) order.
type PlanSubOrder struct {
	OrderID string             `json:"orderId"`
	Price   decimal.Decimal    `json:"price"`
	Type    OrderType          `json:"type"`
	Status  PlanSubOrderStatus `json:"status"`
}
