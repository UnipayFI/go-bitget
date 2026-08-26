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
	StrategyTypeTPSL         StrategyType = "tpsl"
	StrategyTypeTrigger      StrategyType = "trigger"
	StrategyTypeOCO          StrategyType = "oco"
	StrategyTypeTrailingStop StrategyType = "trailing_stop"
	StrategyTypeIceberg      StrategyType = "iceberg"
	StrategyTypeTWAP         StrategyType = "twap"
)

// TriggerBy is the price series a strategy order's trigger watches.
type TriggerBy string

const (
	TriggerByMarket TriggerBy = "market"
	TriggerByMark   TriggerBy = "mark"
	// TriggerByIndex is accepted only as a trailing stop's activationType, and
	// only on futures.
	TriggerByIndex TriggerBy = "index"
)

// TrailType is how a trailing stop measures its distance from the best price.
type TrailType string

const (
	TrailTypeRatio  TrailType = "ratio"
	TrailTypeSpread TrailType = "spread"
)

// IcebergSplitMode is how an iceberg order divides its total quantity.
type IcebergSplitMode string

const (
	// IcebergSplitModeQuantity splits by quantity per sub-order (qtyPerOrder),
	// IcebergSplitModeOrder by sub-order count (splitOrderNumbers).
	IcebergSplitModeQuantity IcebergSplitMode = "quantity"
	IcebergSplitModeOrder    IcebergSplitMode = "order"
)

// IcebergOrderPreference is how an iceberg order prices its sub-orders.
type IcebergOrderPreference string

const (
	IcebergOrderPreferenceFasterExecution IcebergOrderPreference = "faster_execution"
	IcebergOrderPreferenceFixedDistance   IcebergOrderPreference = "fixed_distance"
	IcebergOrderPreferenceFixedPrice      IcebergOrderPreference = "fixed_price"
)

// IcebergExecutionStrategy picks the side of the book an iceberg sub-order
// joins. It applies only to IcebergOrderPreferenceFasterExecution.
type IcebergExecutionStrategy string

const (
	// IcebergExecutionStrategyQueue1 queues at the best price on our own side,
	// IcebergExecutionStrategyCounterparty1 takes the best opposing price.
	IcebergExecutionStrategyQueue1        IcebergExecutionStrategy = "queue1"
	IcebergExecutionStrategyCounterparty1 IcebergExecutionStrategy = "counterparty1"
)

// OffsetType is how a price offset is expressed — as a percentage of the
// reference price or as an absolute spread. It is shared by the iceberg
// fixedDistanceType and the TWAP limitOffsetType.
type OffsetType string

const (
	OffsetTypePercentage OffsetType = "percentage"
	OffsetTypeSpread     OffsetType = "spread"
)

// The four param sets below configure the oco, trailing_stop, iceberg and twap
// strategy types. Bitget types each of them as a list in both the request and
// the response, so they are modelled as slices even though a strategy order
// carries exactly one param set.

// StrategyOCOParams configures an OCO ("one cancels the other") strategy order:
// a resting limit order paired with a conditional order, whichever triggers
// first cancelling the other.
type StrategyOCOParams struct {
	// OCOLimitPrice is the limit leg — bottom-fishing for a buy, take-profit
	// for a sell. OCOTriggerPrice is the conditional leg — chase-up for a buy,
	// stop-loss for a sell.
	OCOLimitPrice   decimal.Decimal `json:"ocoLimitPrice"`
	OCOTriggerPrice decimal.Decimal `json:"ocoTriggerPrice"`
	OCOOrderType    OrderType       `json:"ocoOrderType"`
	// OCOOrderPrice is the conditional leg's execution price, required when
	// OCOOrderType is limit.
	OCOOrderPrice decimal.Decimal `json:"ocoOrderPrice"`
}

// StrategyTrailingStopParams configures a trailing stop: once the market
// reaches ActivationPrice the order follows the best price at TrailVariance,
// firing when the price retraces by that much.
type StrategyTrailingStopParams struct {
	ActivationPrice decimal.Decimal `json:"activationPrice"`
	// ActivationType is the price series watched for activation; index is
	// futures-only.
	ActivationType TriggerBy `json:"activationType"`
	TrailType      TrailType `json:"trailType"`
	// TrailVariance is a percentage for TrailTypeRatio (spot [0.1-20], futures
	// [0.1-10]) and an absolute price distance for TrailTypeSpread.
	TrailVariance decimal.Decimal `json:"trailVariance"`
	// PreOrderType is the type of the order placed when the trail fires;
	// PreOrderPrice is required when it is limit.
	PreOrderType  OrderType       `json:"preOrderType"`
	PreOrderPrice decimal.Decimal `json:"preOrderPrice"`
}

// StrategyIcebergParams configures an iceberg order: one large order worked as
// a series of smaller sub-orders so the full size never shows on the book.
type StrategyIcebergParams struct {
	SplitMode IcebergSplitMode `json:"splitMode"`
	// QtyPerOrder is required for IcebergSplitModeQuantity, SplitOrderNumbers
	// (range [1-100]) for IcebergSplitModeOrder.
	QtyPerOrder       decimal.Decimal        `json:"qtyPerOrder"`
	SplitOrderNumbers string                 `json:"splitOrderNumbers"`
	OrderPreference   IcebergOrderPreference `json:"orderPreference"`
	// ExecutionStrategy applies to IcebergOrderPreferenceFasterExecution.
	ExecutionStrategy IcebergExecutionStrategy `json:"executionStrategy"`
	// FixedDistanceType and Distance apply to
	// IcebergOrderPreferenceFixedDistance, FixedPrice to
	// IcebergOrderPreferenceFixedPrice.
	FixedDistanceType OffsetType      `json:"fixedDistanceType"`
	Distance          decimal.Decimal `json:"distance"`
	FixedPrice        decimal.Decimal `json:"fixedPrice"`
	// PriceLimit caps how far the sub-orders may chase the market; it applies
	// to the faster_execution and fixed_distance preferences.
	PriceLimit decimal.Decimal `json:"priceLimit"`
}

// StrategyTWAPParams configures a TWAP order: the total quantity is sliced into
// equal sub-orders released at a fixed interval over Duration.
type StrategyTWAPParams struct {
	// Duration is the total run time in minutes, range [1-1440]. Interval is
	// the sub-order frequency in seconds: 5, 10, 20, 30 or 60.
	Duration string `json:"duration"`
	Interval string `json:"interval"`
	// OrderType is the sub-order type (defaults to market).
	OrderType OrderType `json:"orderType"`
	// LimitOffsetType picks which of LimitOffsetPercentage (range
	// [0.001-0.1]) and LimitOffsetSpread prices a limit sub-order away from the
	// market. The spread may not exceed 20%.
	LimitOffsetType       OffsetType      `json:"limitOffsetType"`
	LimitOffsetPercentage decimal.Decimal `json:"limitOffsetPercentage"`
	LimitOffsetSpread     decimal.Decimal `json:"limitOffsetSpread"`
	// TWAPTriggerPrice delays the strategy until the market reaches it;
	// TWAPTerminationPrice ends the strategy early when the market reaches it.
	TWAPTriggerPrice     decimal.Decimal `json:"twapTriggerPrice"`
	TWAPTerminationPrice decimal.Decimal `json:"twapTerminationPrice"`
}

// PlaceStrategyOrderService -- POST /api/v3/trade/place-strategy-order (UTA trade read & write)
//
// Places a strategy order of any StrategyType: take-profit/stop-loss ("tpsl"),
// trigger ("trigger"), OCO ("oco"), trailing stop ("trailing_stop"), iceberg
// ("iceberg") or TWAP ("twap"). The last four are configured through their own
// param set — SetOCOParams, SetTrailingStopParams, SetIcebergParams,
// SetTWAPParams — which is required once SetType selects them. Supported
// business lines: spot, margin, and futures. The reply data carries the
// assigned order identifiers.
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

// SetClientOid sets the client order ID (6-hour idempotent validity). Bitget
// honours it only when the strategy type is tpsl.
func (s *PlaceStrategyOrderService) SetClientOrderID(clientOid string) *PlaceStrategyOrderService {
	s.body["clientOid"] = clientOid
	return s
}

// SetType sets the strategy type (defaults to tpsl).
func (s *PlaceStrategyOrderService) SetType(strategyType StrategyType) *PlaceStrategyOrderService {
	s.body["type"] = string(strategyType)
	return s
}

// SetOCOParams sets the OCO configuration, required when the type is oco.
func (s *PlaceStrategyOrderService) SetOCOParams(ocoParams []StrategyOCOParams) *PlaceStrategyOrderService {
	s.body["ocoParams"] = ocoParams
	return s
}

// SetTrailingStopParams sets the trailing-stop configuration, required when the
// type is trailing_stop.
func (s *PlaceStrategyOrderService) SetTrailingStopParams(trailingStopParams []StrategyTrailingStopParams) *PlaceStrategyOrderService {
	s.body["trailingStopParams"] = trailingStopParams
	return s
}

// SetIcebergParams sets the iceberg configuration, required when the type is
// iceberg.
func (s *PlaceStrategyOrderService) SetIcebergParams(icebergParams []StrategyIcebergParams) *PlaceStrategyOrderService {
	s.body["icebergParams"] = icebergParams
	return s
}

// SetTWAPParams sets the TWAP configuration, required when the type is twap.
func (s *PlaceStrategyOrderService) SetTWAPParams(twapParams []StrategyTWAPParams) *PlaceStrategyOrderService {
	s.body["twapParams"] = twapParams
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

// SetReduceOnly sets the reduce-only indicator (defaults to ReduceOnlyNo).
func (s *PlaceStrategyOrderService) SetReduceOnly(reduceOnly ReduceOnly) *PlaceStrategyOrderService {
	s.body["reduceOnly"] = string(reduceOnly)
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
// Returns the account's open (unfilled) strategy orders for a spot, margin, or
// futures category, optionally filtered to a single strategy type.
type GetUnfilledStrategyOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetUnfilledStrategyOrdersService(category Category) *GetUnfilledStrategyOrdersService {
	return &GetUnfilledStrategyOrdersService{c: c, params: map[string]string{"category": string(category)}}
}

// SetType filters by strategy type (tpsl, trigger, oco, trailing_stop,
// iceberg or twap).
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
	// Exactly one of the four param sets below is populated, matching the
	// strategy type the order was placed with.
	OCOParams          []StrategyOCOParams          `json:"ocoParams"`
	TrailingStopParams []StrategyTrailingStopParams `json:"trailingStopParams"`
	IcebergParams      []StrategyIcebergParams      `json:"icebergParams"`
	TWAPParams         []StrategyTWAPParams         `json:"twapParams"`
	CreatedTime        time.Time                    `json:"createdTime"`
	UpdatedTime        time.Time                    `json:"updatedTime"`
}

// GetHistoryStrategyOrdersService -- GET /api/v3/trade/history-strategy-orders (UTA trade read)
//
// Returns the account's historical (completed/cancelled) strategy orders for a
// spot, margin, or futures category, paginated by cursor.
type GetHistoryStrategyOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetHistoryStrategyOrdersService(category Category) *GetHistoryStrategyOrdersService {
	return &GetHistoryStrategyOrdersService{c: c, params: map[string]string{"category": string(category)}}
}

// SetType filters by strategy type (tpsl, trigger, oco, trailing_stop,
// iceberg or twap).
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

// GetStrategySubOrdersService -- GET /api/v3/trade/strategy-sub-orders (UTA trade read)
//
// Returns the sub-orders a strategy order has generated, paginated by cursor.
// Only the strategy types that work an order in slices — iceberg and twap —
// produce more than one.
type GetStrategySubOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetStrategySubOrdersService(orderID string) *GetStrategySubOrdersService {
	return &GetStrategySubOrdersService{c: c, params: map[string]string{"orderId": orderID}}
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetStrategySubOrdersService) SetLimit(limit string) *GetStrategySubOrdersService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetStrategySubOrdersService) SetCursor(cursor string) *GetStrategySubOrdersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetStrategySubOrdersService) Do(ctx context.Context) (*StrategySubOrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/strategy-sub-orders", s.params).WithSign()
	return request.Do[StrategySubOrders](req)
}

type StrategySubOrders struct {
	List   []StrategySubOrder `json:"list"`
	Cursor string             `json:"cursor"`
}

// StrategySubOrder is one order a strategy order placed on the book.
type StrategySubOrder struct {
	SubOrderID       string          `json:"subOrderId"`
	SubClientOrderID string          `json:"subClientOid"`
	Category         Category        `json:"category"`
	Symbol           string          `json:"symbol"`
	Price            decimal.Decimal `json:"price"`
	Qty              decimal.Decimal `json:"qty"`
	CumExecQty       decimal.Decimal `json:"cumExecQty"`
	AvgPrice         decimal.Decimal `json:"avgPrice"`
	Side             Side            `json:"side"`
	PosSide          PosSide         `json:"posSide"` // futures only
	Status           string          `json:"status"`  // filled, cancelled, failed
	CreatedTime      time.Time       `json:"createdTime"`
	UpdatedTime      time.Time       `json:"updatedTime"`
}
