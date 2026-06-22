package uta

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

// TimeInForce determines how long an order stays active.
type TimeInForce string

const (
	TimeInForceGTC      TimeInForce = "gtc"
	TimeInForcePostOnly TimeInForce = "post_only"
	TimeInForceFOK      TimeInForce = "fok"
	TimeInForceIOC      TimeInForce = "ioc"
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
