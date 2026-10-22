package uta

import "strings"

// Category is the product line a symbol or request targets.
type Category string

const (
	CategorySpot        Category = "SPOT"
	CategoryMargin      Category = "MARGIN"
	CategoryUSDTFutures Category = "USDT-FUTURES"
	CategoryCoinFutures Category = "COIN-FUTURES"
	CategoryUSDCFutures Category = "USDC-FUTURES"
)

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

// PosSide is the position direction for futures in hedge mode.
type PosSide string

const (
	PosSideLong  PosSide = "long"
	PosSideShort PosSide = "short"
)

// ReduceOnly is the reduce-only flag used by both the order-placement request
// (SetReduceOnly / WsNewOrder) and the order responses. The UTA wire values are
// lower-case "yes"/"no" (the classic mix endpoints use upper-case "YES"/"NO");
// responses may also arrive as "true"/"false" or "1"/"0", so Bool decodes
// tolerantly.
type ReduceOnly string

const (
	ReduceOnlyYes ReduceOnly = "yes"
	ReduceOnlyNo  ReduceOnly = "no"
)

// Bool reports whether the reduce-only flag is set, tolerating the several
// shapes Bitget returns ("yes"/"true"/"1", any case).
func (r ReduceOnly) Bool() bool {
	switch strings.ToLower(string(r)) {
	case "yes", "true", "1":
		return true
	default:
		return false
	}
}

// AutoBorrow is the spot auto-borrow switch used by the order-placement
// (SetAutoBorrow / WsNewOrder) and max-open-available requests. It defaults to
// disabled. When enabled, and the received coin does not support borrowing
// while the spent coin does, the system borrows the spent coin to cover the
// shortfall if its available balance is insufficient.
type AutoBorrow string

const (
	AutoBorrowYes AutoBorrow = "yes"
	AutoBorrowNo  AutoBorrow = "no"
)

// TimeInForce determines how long an order stays active.
type TimeInForce string

const (
	TimeInForceGTC      TimeInForce = "gtc"
	TimeInForcePostOnly TimeInForce = "post_only"
	TimeInForceFOK      TimeInForce = "fok"
	TimeInForceIOC      TimeInForce = "ioc"
)

// Venue is the order routing strategy for Reality (rtoken) pairs; it has no
// effect on other pairs. Orders default to VenueLitMarket.
type Venue string

const (
	// VenueInternalization matches only inside Bitget's internal (onshore)
	// pool and is never forwarded to the external broker.
	VenueInternalization Venue = "internalization"
	// VenueLitMarket matches internally first, routes any remainder to the
	// external pool, then rests what is still unmatched as a limit order.
	VenueLitMarket Venue = "lit_market"
)

// RTokenTimeInForce is the validity of a Reality (rtoken) order; it has no
// effect on other pairs. Orders default to RTokenTimeInForceGTD.
type RTokenTimeInForce string

const (
	// RTokenTimeInForceGTD is valid for the current session: an order placed
	// during US market hours (weekdays) lives until the end of after-hours
	// trading, one placed on a weekend until just before the next US open.
	RTokenTimeInForceGTD RTokenTimeInForce = "gtd"
	// RTokenTimeInForceGTC is valid for 90 days.
	RTokenTimeInForceGTC RTokenTimeInForce = "gtc"
)

// BBOLevel is the best-bid/offer price level of a limit order: level 1 or 5
// on the order's own side (queue) or on the opposite side (counterparty).
type BBOLevel string

const (
	BBOLevelQueue1        BBOLevel = "queue1"
	BBOLevelQueue5        BBOLevel = "queue5"
	BBOLevelCounterparty1 BBOLevel = "counterparty1"
	BBOLevelCounterparty5 BBOLevel = "counterparty5"
)

// HoldMode is the futures position mode (one-way vs hedge).
type HoldMode string

const (
	HoldModeOneWay HoldMode = "one_way_mode"
	HoldModeHedge  HoldMode = "hedge_mode"
)

// MarginMode is the margin management approach.
type MarginMode string

const (
	MarginModeCrossed  MarginMode = "crossed"
	MarginModeIsolated MarginMode = "isolated"
)

// OrderStatus is the order lifecycle state.
type OrderStatus string

const (
	OrderStatusLive            OrderStatus = "live"
	OrderStatusNew             OrderStatus = "new"
	OrderStatusPartiallyFilled OrderStatus = "partially_filled"
	OrderStatusFilled          OrderStatus = "filled"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

// TradeScope is the market participant role of a fill.
type TradeScope string

const (
	TradeScopeTaker TradeScope = "taker"
	TradeScopeMaker TradeScope = "maker"
)

// ExecType classifies how an order was executed.
type ExecType string

const (
	ExecTypeNormal      ExecType = "normal"
	ExecTypeOffset      ExecType = "offset"
	ExecTypeReduce      ExecType = "reduce"
	ExecTypeLiquidation ExecType = "liquidation"
	ExecTypeDelivery    ExecType = "delivery"
)

// KlineGranularity is a candlestick interval accepted by the market endpoints.
type KlineGranularity string

const (
	Granularity1m  KlineGranularity = "1m"
	Granularity3m  KlineGranularity = "3m"
	Granularity5m  KlineGranularity = "5m"
	Granularity15m KlineGranularity = "15m"
	Granularity30m KlineGranularity = "30m"
	Granularity1H  KlineGranularity = "1H"
	Granularity4H  KlineGranularity = "4H"
	Granularity6H  KlineGranularity = "6H"
	Granularity12H KlineGranularity = "12H"
	Granularity1D  KlineGranularity = "1D"
	Granularity3D  KlineGranularity = "3D"
	Granularity1W  KlineGranularity = "1W"
	Granularity1M  KlineGranularity = "1M"

	Granularity6Hutc  KlineGranularity = "6Hutc"
	Granularity12Hutc KlineGranularity = "12Hutc"
	Granularity1Dutc  KlineGranularity = "1Dutc"
	Granularity3Dutc  KlineGranularity = "3Dutc"
	Granularity1Wutc  KlineGranularity = "1Wutc"
	Granularity1Mutc  KlineGranularity = "1Mutc"
)

// KlineType selects which price series a candlestick query returns. The wire
// values are lower-case (verified live); "premium" is the premium-index series.
type KlineType string

const (
	KlineTypeMarket  KlineType = "market"
	KlineTypeMark    KlineType = "mark"
	KlineTypeIndex   KlineType = "index"
	KlineTypePremium KlineType = "premium"
)

// SymbolType distinguishes perpetual and delivery contracts.
type SymbolType string

const (
	SymbolTypePerpetual SymbolType = "perpetual"
	SymbolTypeDelivery  SymbolType = "delivery"
)

// InstrumentStatus is the listing condition of a trading pair.
type InstrumentStatus string

const (
	InstrumentStatusListed        InstrumentStatus = "listed"
	InstrumentStatusOnline        InstrumentStatus = "online"
	InstrumentStatusLimitOpen     InstrumentStatus = "limit_open"
	InstrumentStatusLimitClose    InstrumentStatus = "limit_close"
	InstrumentStatusOffline       InstrumentStatus = "offline"
	InstrumentStatusRestrictedAPI InstrumentStatus = "restrictedAPI"
)

// AccountMode is the unified-account margin mode.
type AccountMode string

const (
	AccountModeUnion    AccountMode = "UNION"    // multi-currency / unified margin
	AccountModeIsolated AccountMode = "ISOLATED" // isolated margin
	AccountModeClassic  AccountMode = "CLASSIC"  // classic (reserved for non-UTA)
)
