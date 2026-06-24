package spot

// Side is the order direction.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"
)

// OrderType is the order execution method.
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

// Force is the execution strategy (time in force) for an order.
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
	OrderStatusCancelled       OrderStatus = "cancelled"
)

// SelfTradePreventionMode controls how self-trades are prevented.
type SelfTradePreventionMode string

const (
	SelfTradePreventionModeNone        SelfTradePreventionMode = "none"
	SelfTradePreventionModeCancelTaker SelfTradePreventionMode = "cancel_taker"
	SelfTradePreventionModeCancelMaker SelfTradePreventionMode = "cancel_maker"
	SelfTradePreventionModeCancelBoth  SelfTradePreventionMode = "cancel_both"
)

// TriggerType is the price reference used to trigger a plan order.
type TriggerType string

const (
	TriggerTypeFillPrice TriggerType = "fill_price"
	TriggerTypeMarkPrice TriggerType = "mark_price"
)

// PlanType is how a plan order's size is denominated: by base coin quantity
// (amount) or by quote coin value (total).
type PlanType string

const (
	PlanTypeAmount PlanType = "amount"
	PlanTypeTotal  PlanType = "total"
)

// TpslType distinguishes a normal order from a take-profit/stop-loss order.
type TPSLType string

const (
	TPSLTypeNormal TPSLType = "normal"
	TPSLTypeTPSL   TPSLType = "tpsl"
)
