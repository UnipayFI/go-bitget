package mix

// ProductType is the futures product line a symbol or request targets.
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
)

// MarginMode is the margin management approach for a position.
type MarginMode string

const (
	MarginModeIsolated MarginMode = "isolated"
	MarginModeCrossed  MarginMode = "crossed"
)

// Side is the order direction.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// TradeSide is whether an order opens or closes a position. It is only
// meaningful in hedge (two-way) position mode.
type TradeSide string

const (
	TradeSideOpen  TradeSide = "open"
	TradeSideClose TradeSide = "close"
)

// HoldSide is the position direction.
type HoldSide string

const (
	HoldSideLong  HoldSide = "long"
	HoldSideShort HoldSide = "short"
)

// PositionMode is the account-level futures position mode (one-way vs hedge).
type PositionMode string

const (
	PositionModeOneWay PositionMode = "one_way_mode"
	PositionModeHedge  PositionMode = "hedge_mode"
)

// OrderType is the order execution method.
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// Force determines how long an order stays active (time in force).
type Force string

const (
	ForceGTC      Force = "gtc"
	ForcePostOnly Force = "post_only"
	ForceFOK      Force = "fok"
	ForceIOC      Force = "ioc"
)

// OrderStatus is the order lifecycle state.
type OrderStatus string

const (
	OrderStatusLive            OrderStatus = "live"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCanceled        OrderStatus = "canceled"
)

// PlanType is the trigger/TP-SL plan order category.
type PlanType string

const (
	PlanTypeNormalPlan PlanType = "normal_plan" // standard trigger order
	PlanTypeTrackPlan  PlanType = "track_plan"  // trailing stop (trigger order)
	PlanTypeProfitPlan PlanType = "profit_plan" // partial take-profit
	PlanTypeLossPlan   PlanType = "loss_plan"   // partial stop-loss
	PlanTypePosProfit  PlanType = "pos_profit"  // whole-position take-profit
	PlanTypePosLoss    PlanType = "pos_loss"    // whole-position stop-loss
	PlanTypeMovingPlan PlanType = "moving_plan" // trailing stop (TP-SL)
)

// TriggerType is the price series a plan order's trigger is measured against.
type TriggerType string

const (
	TriggerTypeFillPrice TriggerType = "fill_price"
	TriggerTypeMarkPrice TriggerType = "mark_price"
)

// SelfTradePreventionMode controls how self-matching orders are handled.
type SelfTradePreventionMode string

const (
	SelfTradePreventionModeNone        SelfTradePreventionMode = "none"
	SelfTradePreventionModeCancelTaker SelfTradePreventionMode = "cancel_taker"
	SelfTradePreventionModeCancelMaker SelfTradePreventionMode = "cancel_maker"
	SelfTradePreventionModeCancelBoth  SelfTradePreventionMode = "cancel_both"
)

// ReduceOnly flags whether an order may only reduce an existing position. The
// wire values are upper-case "YES"/"NO" on the trade endpoints.
type ReduceOnly string

const (
	ReduceOnlyYes ReduceOnly = "YES"
	ReduceOnlyNo  ReduceOnly = "NO"
)

// AssetMode is the futures account's asset/collateral mode: union (multi-asset,
// shared collateral pool) or single (per-coin margin).
type AssetMode string

const (
	AssetModeUnion  AssetMode = "union"
	AssetModeSingle AssetMode = "single"
)

// AutoMargin is the isolated-position auto-margin top-up toggle.
type AutoMargin string

const (
	AutoMarginOn  AutoMargin = "on"
	AutoMarginOff AutoMargin = "off"
)
