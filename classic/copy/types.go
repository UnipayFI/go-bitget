package copy

// ProductType is the futures product line a copy-trading request or record
// targets. It is sent on essentially every futures copy-trading endpoint
// (trader and follower) and echoed back in their responses.
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
)

// Side is the order direction used by the order-placement and order-query
// copy-trading endpoints.
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

// PosSide is the position direction as reported by the tracking-order
// responses (Trader-Order-Current-Track, follower current/history tracking
// orders). It is meaningful in hedge mode.
type PosSide string

const (
	PosSideLong  PosSide = "long"
	PosSideShort PosSide = "short"
)

// HoldSide is the position direction supplied on request bodies such as
// close-positions. It carries the same wire values as PosSide; both names
// exist because the API uses distinct field names for the request and
// response sides of the same concept.
type HoldSide string

const (
	HoldSideLong  HoldSide = "long"
	HoldSideShort HoldSide = "short"
)

// MarginMode is the margin management approach for a copy-traded position.
// The wire value for cross margin is "crossed" (matching the classic mix
// futures product line).
type MarginMode string

const (
	MarginModeIsolated MarginMode = "isolated"
	MarginModeCrossed  MarginMode = "crossed"
)

// TraceStatus is the lifecycle state of a copy-trading (trace) order, shared
// across the trader and follower order-query endpoints. Open/current-order
// endpoints report the in-progress states; history endpoints report the
// terminal states.
type TraceStatus string

const (
	TraceStatusTracking TraceStatus = "tracking" // order is being actively copied
	TraceStatusProfit   TraceStatus = "profit"   // closed in profit
	TraceStatusLoss     TraceStatus = "loss"     // closed at a loss
)
