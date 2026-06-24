package ws

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// This file wraps the classic-account v2 private FUTURES (mix) channels. Every
// exported identifier is prefixed with Mix to avoid collisions with the spot /
// margin sibling services in the same package. These channels all require login
// (private=true, handled by the framework) and are scoped by instType
// (USDT-FUTURES / COIN-FUTURES / USDC-FUTURES); the service constructors take an
// InstType so callers pick the product line, defaulting to USDT-FUTURES at the
// call site.

// SubscribeMixAccountService -- private "account" channel (futures account
// assets per margin coin). Subscribed with coin="default" (all coins).
type SubscribeMixAccountService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixAccountService(instType InstType) *SubscribeMixAccountService {
	return &SubscribeMixAccountService{c: c, instType: instType}
}

func (s *SubscribeMixAccountService) Do(ctx context.Context, cb WsHandler[MixWsAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsAccount](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "account", Coin: "default"}, cb)
}

type MixWsAccount struct {
	MarginCoin          string          `json:"marginCoin"`          // margin coin identifier
	Frozen              decimal.Decimal `json:"frozen"`              // locked quantity (margin coin)
	Available           decimal.Decimal `json:"available"`           // currently available assets
	MaxOpenPosAvailable decimal.Decimal `json:"maxOpenPosAvailable"` // max available balance to open positions
	MaxTransferOut      decimal.Decimal `json:"maxTransferOut"`      // max transferable amount
	Equity              decimal.Decimal `json:"equity"`              // account assets
	USDTEquity          decimal.Decimal `json:"usdtEquity"`          // account equity in USD
	CrossedRiskRate     decimal.Decimal `json:"crossedRiskRate"`     // risk ratio in cross margin mode
	UnrealizedPL        decimal.Decimal `json:"unrealizedPL"`        // unrealized PnL
	UnionTotalMargin    decimal.Decimal `json:"unionTotalMargin"`    // margin amount under union margin mode
	UnionAvailable      decimal.Decimal `json:"unionAvailable"`      // available balance under union margin mode
	UnionMm             decimal.Decimal `json:"unionMm"`             // maintenance margin under union margin mode
	AssetsMode          string          `json:"assetsMode"`          // union (unified) or single (single-currency)
}

// SubscribeMixEquityService -- private "equity" channel (account-level equity in
// BTC/USDT plus multi-asset margin figures).
type SubscribeMixEquityService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixEquityService(instType InstType) *SubscribeMixEquityService {
	return &SubscribeMixEquityService{c: c, instType: instType}
}

func (s *SubscribeMixEquityService) Do(ctx context.Context, cb WsHandler[MixWsEquity]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsEquity](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "equity"}, cb)
}

type MixWsEquity struct {
	BtcEquity        decimal.Decimal `json:"btcEquity"`        // account equity (BTC)
	USDTEquity       decimal.Decimal `json:"usdtEquity"`       // account equity (USDT)
	USDTUnrealized   decimal.Decimal `json:"usdtUnrealized"`   // unrealized PnL (USDT)
	UnionTotalMargin decimal.Decimal `json:"unionTotalMargin"` // total multi-asset margin
	UnionAvailable   decimal.Decimal `json:"unionAvailable"`   // available balance under multi-asset margin mode
	UnionMm          decimal.Decimal `json:"unionMm"`          // maintenance margin under multi-asset margin mode
}

// SubscribeMixPositionsService -- private "positions" channel (open positions).
// Subscribed with instId="default" (all symbols).
type SubscribeMixPositionsService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixPositionsService(instType InstType) *SubscribeMixPositionsService {
	return &SubscribeMixPositionsService{c: c, instType: instType}
}

func (s *SubscribeMixPositionsService) Do(ctx context.Context, cb WsHandler[MixWsPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsPosition](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "positions", InstID: "default"}, cb)
}

type MixWsPosition struct {
	PosID            string          `json:"posId"`            // position ID
	InstID           string          `json:"instId"`           // product ID
	MarginCoin       string          `json:"marginCoin"`       // currency of occupied margin
	MarginSize       decimal.Decimal `json:"marginSize"`       // occupied margin (amount)
	MarginMode       string          `json:"marginMode"`       // margin mode (crossed / isolated)
	HoldSide         string          `json:"holdSide"`         // position direction (long / short)
	PosMode          string          `json:"posMode"`          // position mode (one_way_mode / hedge_mode)
	Total            decimal.Decimal `json:"total"`            // open position size
	Available        decimal.Decimal `json:"available"`        // size of positions that can be closed
	Frozen           decimal.Decimal `json:"frozen"`           // amount of frozen margin
	OpenPriceAvg     decimal.Decimal `json:"openPriceAvg"`     // average entry price
	Leverage         decimal.Decimal `json:"leverage"`         // leverage
	AchievedProfits  decimal.Decimal `json:"achievedProfits"`  // realized PnL
	UnrealizedPL     decimal.Decimal `json:"unrealizedPL"`     // unrealized PnL
	UnrealizedPLR    decimal.Decimal `json:"unrealizedPLR"`    // unrealized ROI
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"` // estimated liquidation price
	KeepMarginRate   decimal.Decimal `json:"keepMarginRate"`   // maintenance margin rate
	MarginRate       decimal.Decimal `json:"marginRate"`       // occupancy rate of margin
	BreakEvenPrice   decimal.Decimal `json:"breakEvenPrice"`   // position breakeven price
	TotalFee         decimal.Decimal `json:"totalFee"`         // accumulated funding fee during position
	DeductedFee      decimal.Decimal `json:"deductedFee"`      // deducted transaction fees
	MarkPrice        decimal.Decimal `json:"markPrice"`        // mark price
	AssetMode        string          `json:"assetMode"`        // account mode (union / single)
	CTime            time.Time       `json:"cTime"`            // position creation time
	UTime            time.Time       `json:"uTime"`            // latest position update time
}

// SubscribeMixPositionsHistoryService -- private "positions-history" channel
// (pushed when a position is fully closed). Subscribed with instId="default".
type SubscribeMixPositionsHistoryService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixPositionsHistoryService(instType InstType) *SubscribeMixPositionsHistoryService {
	return &SubscribeMixPositionsHistoryService{c: c, instType: instType}
}

func (s *SubscribeMixPositionsHistoryService) Do(ctx context.Context, cb WsHandler[MixWsPositionHistory]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsPositionHistory](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "positions-history", InstID: "default"}, cb)
}

type MixWsPositionHistory struct {
	PosID           string          `json:"posId"`           // position identifier
	InstID          string          `json:"instId"`          // product ID
	MarginCoin      string          `json:"marginCoin"`      // margin currency
	MarginMode      string          `json:"marginMode"`      // fixed (isolated) or crossed
	HoldSide        string          `json:"holdSide"`        // direction of position
	PosMode         string          `json:"posMode"`         // position mode
	OpenPriceAvg    decimal.Decimal `json:"openPriceAvg"`    // average entry price
	ClosePriceAvg   decimal.Decimal `json:"closePriceAvg"`   // average close price
	OpenSize        decimal.Decimal `json:"openSize"`        // quantity opened
	CloseSize       decimal.Decimal `json:"closeSize"`       // quantity closed
	AchievedProfits decimal.Decimal `json:"achievedProfits"` // realized PnL
	SettleFee       decimal.Decimal `json:"settleFee"`       // settlement fees
	OpenFee         decimal.Decimal `json:"openFee"`         // total opening fees
	CloseFee        decimal.Decimal `json:"closeFee"`        // total closing fees
	CTime           time.Time       `json:"cTime"`           // position creation time
	UTime           time.Time       `json:"uTime"`           // latest position update time
}

// SubscribeMixFillService -- private "fill" channel (real-time executions).
// instId is optional; subscribe with "default" for all symbols.
type SubscribeMixFillService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixFillService(instType InstType) *SubscribeMixFillService {
	return &SubscribeMixFillService{c: c, instType: instType}
}

func (s *SubscribeMixFillService) Do(ctx context.Context, cb WsHandler[MixWsFill]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsFill](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "fill", InstID: "default"}, cb)
}

type MixWsFill struct {
	OrderID       string          `json:"orderId"`     // order identifier
	ClientOrderID string          `json:"clientOid"`   // user-defined order ID
	TradeID       string          `json:"tradeId"`     // trade identifier
	Symbol        string          `json:"symbol"`      // trading pair name
	Side          string          `json:"side"`        // buy / sell
	OrderType     string          `json:"orderType"`   // limit / market
	PosMode       string          `json:"posMode"`     // one_way_mode / hedge_mode
	Price         decimal.Decimal `json:"price"`       // execution price
	BaseVolume    decimal.Decimal `json:"baseVolume"`  // base asset quantity traded
	QuoteVolume   decimal.Decimal `json:"quoteVolume"` // quote asset quantity traded
	Profit        decimal.Decimal `json:"profit"`      // realized PnL
	TradeSide     string          `json:"tradeSide"`   // trade classification (open/close/...)
	TradeScope    string          `json:"tradeScope"`  // taker / maker
	FeeDetail     []MixWsFillFee  `json:"feeDetail"`   // transaction fee breakdown
	CTime         time.Time       `json:"cTime"`       // creation time
	UTime         time.Time       `json:"uTime"`       // update time
}

type MixWsFillFee struct {
	FeeCoin           string          `json:"feeCoin"`           // fee currency
	Deduction         string          `json:"deduction"`         // yes / no
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"` // deducted fee amount
	TotalFee          decimal.Decimal `json:"totalFee"`          // complete fee charged
}

// SubscribeMixOrdersService -- private "orders" channel (order lifecycle
// updates). instId is optional; subscribe with "default" for all symbols.
type SubscribeMixOrdersService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixOrdersService(instType InstType) *SubscribeMixOrdersService {
	return &SubscribeMixOrdersService{c: c, instType: instType}
}

func (s *SubscribeMixOrdersService) Do(ctx context.Context, cb WsHandler[MixWsOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsOrder](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "orders", InstID: "default"}, cb)
}

type MixWsOrder struct {
	InstID           string          `json:"instId"`           // product ID, e.g. ETHUSDT
	OrderID          string          `json:"orderId"`          // order ID
	ClientOrderID    string          `json:"clientOid"`        // customized order ID
	Price            decimal.Decimal `json:"price"`            // order price
	Size             decimal.Decimal `json:"size"`             // original order amount in coin
	PosMode          string          `json:"posMode"`          // one_way_mode / hedge_mode
	EnterPointSource string          `json:"enterPointSource"` // order source (WEB, API, SYS, ANDROID, IOS)
	TradeSide        string          `json:"tradeSide"`        // direction (open, close, reduce_close_long, ...)
	NotionalUSD      decimal.Decimal `json:"notionalUsd"`      // estimated USD value of orders
	OrderType        string          `json:"orderType"`        // limit / market
	Force            string          `json:"force"`            // order validity period
	Side             string          `json:"side"`             // order direction
	PosSide          string          `json:"posSide"`          // position direction (long, short, net)
	MarginMode       string          `json:"marginMode"`       // crossed / isolated
	MarginCoin       string          `json:"marginCoin"`       // margin coin
	FillPrice        decimal.Decimal `json:"fillPrice"`        // latest filled price
	TradeID          string          `json:"tradeId"`          // latest transaction ID
	BaseVolume       decimal.Decimal `json:"baseVolume"`       // number of latest filled orders
	FillTime         time.Time       `json:"fillTime"`         // latest transaction time
	FillFee          decimal.Decimal `json:"fillFee"`          // transaction fee of latest transaction
	FillFeeCoin      string          `json:"fillFeeCoin"`      // currency of transaction fee
	TradeScope       string          `json:"tradeScope"`       // liquidity direction (T: taker, M: maker)
	AccBaseVolume    decimal.Decimal `json:"accBaseVolume"`    // total filled quantity
	FillNotionalUSD  decimal.Decimal `json:"fillNotionalUsd"`  // USD value of filled orders
	PriceAvg         decimal.Decimal `json:"priceAvg"`         // average filled price
	Status           string          `json:"status"`           // live, partially_filled, filled, canceled
	CancelReason     string          `json:"cancelReason"`     // cancellation reason
	Leverage         decimal.Decimal `json:"leverage"`         // leverage
	FeeDetail        []MixWsOrderFee `json:"feeDetail"`        // transaction fee details

	PnL                           decimal.Decimal `json:"pnl"`                           // profit
	UTime                         time.Time       `json:"uTime"`                         // order update time
	CTime                         time.Time       `json:"cTime"`                         // order creation time
	ReduceOnly                    string          `json:"reduceOnly"`                    // reduce-only status (yes / no)
	PresetStopSurplusPrice        decimal.Decimal `json:"presetStopSurplusPrice"`        // take-profit price
	PresetStopLossPrice           decimal.Decimal `json:"presetStopLossPrice"`           // stop-loss price
	StpMode                       string          `json:"stpMode"`                       // none, cancel_taker, cancel_maker, cancel_both
	TotalProfits                  decimal.Decimal `json:"totalProfits"`                  // total profits
	PresetStopSurplusExecutePrice decimal.Decimal `json:"presetStopSurplusExecutePrice"` // TP execution price
	PresetStopLossExecutePrice    decimal.Decimal `json:"presetStopLossExecutePrice"`    // SL execution price
}

type MixWsOrderFee struct {
	FeeCoin string          `json:"feeCoin"` // transaction fee currency
	Fee     decimal.Decimal `json:"fee"`     // platform transaction fee charged
}

// SubscribeMixOrdersAlgoService -- private "orders-algo" channel (trigger / plan
// orders). instId is optional; subscribe with "default" for all symbols.
type SubscribeMixOrdersAlgoService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixOrdersAlgoService(instType InstType) *SubscribeMixOrdersAlgoService {
	return &SubscribeMixOrdersAlgoService{c: c, instType: instType}
}

func (s *SubscribeMixOrdersAlgoService) Do(ctx context.Context, cb WsHandler[MixWsOrderAlgo]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsOrderAlgo](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "orders-algo", InstID: "default"}, cb)
}

type MixWsOrderAlgo struct {
	InstID           string          `json:"instId"`           // product ID
	OrderID          string          `json:"orderId"`          // bot order identifier
	ClientOrderID    string          `json:"clientOid"`        // custom bot order identifier
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`     // price level that activates the order
	TriggerType      string          `json:"triggerType"`      // fill_price / mark_price
	TriggerTime      time.Time       `json:"triggerTime"`      // activation timestamp
	PlanType         string          `json:"planType"`         // pl, tp, sl, ptp, psl, track, mtpsl
	Price            decimal.Decimal `json:"price"`            // order execution price
	ExecutePrice     decimal.Decimal `json:"executePrice"`     // actual execution price
	Size             decimal.Decimal `json:"size"`             // original order amount in coin
	ActualSize       decimal.Decimal `json:"actualSize"`       // actual filled amount in coin
	OrderType        string          `json:"orderType"`        // limit / market
	Side             string          `json:"side"`             // order direction
	TradeSide        string          `json:"tradeSide"`        // trading direction
	PosSide          string          `json:"posSide"`          // position direction
	MarginCoin       string          `json:"marginCoin"`       // collateral currency
	Status           string          `json:"status"`           // live, executed, fail_execute, cancelled, executing
	PosMode          string          `json:"posMode"`          // one_way_mode / hedge_mode
	EnterPointSource string          `json:"enterPointSource"` // origin: WEB, API, SYS, ANDROID, IOS

	StopSurplusTriggerPrice decimal.Decimal `json:"stopSurplusTriggerPrice"` // take-profit trigger price
	StopSurplusPrice        decimal.Decimal `json:"stopSurplusPrice"`        // take-profit execution price
	StopSurplusTriggerType  string          `json:"stopSurplusTriggerType"`  // take-profit trigger type
	StopLossTriggerPrice    decimal.Decimal `json:"stopLossTriggerPrice"`    // stop-loss trigger price
	StopLossPrice           decimal.Decimal `json:"stopLossPrice"`           // stop-loss execution price
	StopLossTriggerType     string          `json:"stopLossTriggerType"`     // stop-loss trigger type
	StpMode                 string          `json:"stpMode"`                 // self-trade prevention mode
	CTime                   time.Time       `json:"cTime"`                   // creation timestamp
	UTime                   time.Time       `json:"uTime"`                   // last update timestamp
}

// SubscribeMixADLNotificationService -- private "adl-noti" channel (auto-
// deleveraging notifications). instId is optional; subscribe with "default".
type SubscribeMixADLNotificationService struct {
	c        *WebSocketClient
	instType InstType
}

func (c *WebSocketClient) NewSubscribeMixADLNotificationService(instType InstType) *SubscribeMixADLNotificationService {
	return &SubscribeMixADLNotificationService{c: c, instType: instType}
}

func (s *SubscribeMixADLNotificationService) Do(ctx context.Context, cb WsHandler[MixWsADLNotification]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsADLNotification](ctx, s.c, true,
		WsArg{InstType: string(s.instType), Channel: "adl-noti", InstID: "default"}, cb)
}

type MixWsADLNotification struct {
	Symbol string          `json:"symbol"` // symbol name
	Side   string          `json:"side"`   // position side: buy / sell
	Status string          `json:"status"` // ADL status; currently "triggered"
	Price  decimal.Decimal `json:"price"`  // price at which ADL was executed
	Amount decimal.Decimal `json:"amount"` // execution quantity in quote coin units
	Ts     time.Time       `json:"ts"`     // start timestamp
}
