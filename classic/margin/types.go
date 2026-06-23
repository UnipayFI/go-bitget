package margin

// Side is the order direction. The basic values buy/sell are accepted on order
// placement; the order-history response may additionally report system-driven
// directions (forced liquidation and system repayment). The "systemRepay-selll"
// spelling is reproduced verbatim from the Bitget documentation.
type Side string

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"

	SideLiquidationBuy   Side = "liquidation-buy"
	SideLiquidationSell  Side = "liquidation-sell"
	SideSystemRepayBuy   Side = "systemRepay-buy"
	SideSystemRepaySelll Side = "systemRepay-selll"
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

// OrderStatus is the order lifecycle state. The "partially_fill" and "reject"
// wire values are reproduced verbatim from the Bitget documentation (they differ
// from the UTA spellings "partially_filled"/"cancelled").
type OrderStatus string

const (
	OrderStatusLive          OrderStatus = "live"
	OrderStatusPartiallyFill OrderStatus = "partially_fill"
	OrderStatusFilled        OrderStatus = "filled"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusReject        OrderStatus = "reject"
)

// LoanType is the margin order model: whether borrowing and/or repayment is
// handled automatically when placing the order.
type LoanType string

const (
	LoanTypeNormal           LoanType = "normal"
	LoanTypeAutoLoan         LoanType = "autoLoan"
	LoanTypeAutoRepay        LoanType = "autoRepay"
	LoanTypeAutoLoanAndRepay LoanType = "autoLoanAndRepay"
)

// SelfTradePreventionMode controls how matching against the trader's own resting
// orders is prevented.
type SelfTradePreventionMode string

const (
	SelfTradePreventionModeNone        SelfTradePreventionMode = "none"
	SelfTradePreventionModeCancelTaker SelfTradePreventionMode = "cancel_taker"
	SelfTradePreventionModeCancelMaker SelfTradePreventionMode = "cancel_maker"
	SelfTradePreventionModeCancelBoth  SelfTradePreventionMode = "cancel_both"
)

// EnterPointSource identifies the channel through which an order was created.
type EnterPointSource string

const (
	EnterPointSourceWeb     EnterPointSource = "WEB"
	EnterPointSourceAPI     EnterPointSource = "API"
	EnterPointSourceSys     EnterPointSource = "SYS"
	EnterPointSourceAndroid EnterPointSource = "ANDROID"
	EnterPointSourceIOS     EnterPointSource = "IOS"
)
