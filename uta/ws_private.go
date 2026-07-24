package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SubscribeAccountService -- private "account" channel (equity + balances).
type SubscribeAccountService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeAccountService() *SubscribeAccountService {
	return &SubscribeAccountService{c: c}
}

func (s *SubscribeAccountService) Do(ctx context.Context, cb WsHandler[WsAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsAccount](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "account"}, cb)
}

type WsAccount struct {
	TotalEquity      decimal.Decimal `json:"totalEquity"`
	EffEquity        decimal.Decimal `json:"effEquity"`
	UnrealizedPnL    decimal.Decimal `json:"unrealisedPnL"`
	Imr              decimal.Decimal `json:"imr"`
	Mmr              decimal.Decimal `json:"mmr"`
	MgnRatio         decimal.Decimal `json:"mgnRatio"`
	PositionMgnRatio decimal.Decimal `json:"positionMgnRatio"`
	Bonus            decimal.Decimal `json:"bonus"` // USDT bonus amount
	Coin             []WsAccountCoin `json:"coin"`
}

type WsAccountCoin struct {
	Coin      string          `json:"coin"`
	Equity    decimal.Decimal `json:"equity"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
	Borrow    decimal.Decimal `json:"borrow"`
	Debts     decimal.Decimal `json:"debts"`
	USDValue  decimal.Decimal `json:"usdValue"`
}

// SubscribePositionService -- private "position" channel.
type SubscribePositionService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribePositionService() *SubscribePositionService {
	return &SubscribePositionService{c: c}
}

func (s *SubscribePositionService) Do(ctx context.Context, cb WsHandler[WsPosition]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsPosition](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "position"}, cb)
}

type WsPosition struct {
	Symbol          string          `json:"symbol"`
	MarginCoin      string          `json:"marginCoin"`
	HoldMode        HoldMode        `json:"holdMode"`
	PosSide         PosSide         `json:"posSide"`
	MarginMode      MarginMode      `json:"marginMode"`
	Size            decimal.Decimal `json:"size"`
	Available       decimal.Decimal `json:"available"`
	Frozen          decimal.Decimal `json:"frozen"`
	Leverage        decimal.Decimal `json:"leverage"`
	MarginSize      decimal.Decimal `json:"marginSize"`
	AvgPrice        decimal.Decimal `json:"avgPrice"`
	MarkPrice       decimal.Decimal `json:"markPrice"`
	BreakEvenPrice  decimal.Decimal `json:"breakEvenPrice"`
	LiqPrice        decimal.Decimal `json:"liqPrice"`
	Mmr             decimal.Decimal `json:"mmr"`
	UnrealizedPnL   decimal.Decimal `json:"unrealisedPnl"`
	CurRealisedPnL  decimal.Decimal `json:"curRealisedPnl"`
	ProfitRate      decimal.Decimal `json:"profitRate"`
	TotalFundingFee decimal.Decimal `json:"totalFundingFee"`
	OpenFeeTotal    decimal.Decimal `json:"openFeeTotal"`
	CloseFeeTotal   decimal.Decimal `json:"closeFeeTotal"`
	PositionStatus  string          `json:"positionStatus"`
	CreatedTime     time.Time       `json:"createdTime"`
	UpdatedTime     time.Time       `json:"updatedTime"`
}

// SubscribeOrderService -- private "order" channel (order lifecycle updates).
type SubscribeOrderService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeOrderService() *SubscribeOrderService {
	return &SubscribeOrderService{c: c}
}

func (s *SubscribeOrderService) Do(ctx context.Context, cb WsHandler[WsOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrder](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "order"}, cb)
}

type WsOrder struct {
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	Price         decimal.Decimal `json:"price"`
	Qty           decimal.Decimal `json:"qty"`
	Amount        decimal.Decimal `json:"amount"`
	HoldMode      HoldMode        `json:"holdMode"`
	HoldSide      PosSide         `json:"holdSide"`
	TradeSide     string          `json:"tradeSide"`
	DelegateType  string          `json:"delegateType"`
	OrderType     OrderType       `json:"orderType"`
	TimeInForce   TimeInForce     `json:"timeInForce"`
	Side          Side            `json:"side"`
	MarginMode    MarginMode      `json:"marginMode"`
	MarginCoin    string          `json:"marginCoin"`
	ReduceOnly    ReduceOnly      `json:"reduceOnly"`
	CumExecQty    decimal.Decimal `json:"cumExecQty"`
	CumExecValue  decimal.Decimal `json:"cumExecValue"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	TotalProfit   decimal.Decimal `json:"totalProfit"`
	OrderStatus   OrderStatus     `json:"orderStatus"`
	CancelReason  string          `json:"cancelReason"`
	Leverage      decimal.Decimal `json:"leverage"`
	StpMode       string          `json:"stpMode"`
	FeeDetail     []FeeDetail     `json:"feeDetail"`
	CreatedTime   time.Time       `json:"createdTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
}

// SubscribeFillService -- private "fill" channel (real-time executions).
type SubscribeFillService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeFillService() *SubscribeFillService {
	return &SubscribeFillService{c: c}
}

func (s *SubscribeFillService) Do(ctx context.Context, cb WsHandler[WsFill]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFill](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "fill"}, cb)
}

type WsFill struct {
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	ExecID        string          `json:"execId"`
	ExecLinkID    string          `json:"execLinkId"`
	OrderType     OrderType       `json:"orderType"`
	Side          Side            `json:"side"`
	HoldSide      PosSide         `json:"holdSide"`
	TradeSide     string          `json:"tradeSide"`
	TradeScope    TradeScope      `json:"tradeScope"`
	ExecPrice     decimal.Decimal `json:"execPrice"`
	ExecQty       decimal.Decimal `json:"execQty"`
	ExecValue     decimal.Decimal `json:"execValue"`
	ExecPnL       decimal.Decimal `json:"execPnl"`
	FeeDetail     []FeeDetail     `json:"feeDetail"`
	IsRPI         string          `json:"isRPI"`
	ExecTime      time.Time       `json:"execTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
}

// SubscribeFastFillService -- private "fast-fill" channel (low-latency fills).
type SubscribeFastFillService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeFastFillService() *SubscribeFastFillService {
	return &SubscribeFastFillService{c: c}
}

func (s *SubscribeFastFillService) Do(ctx context.Context, cb WsHandler[WsFastFill]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFastFill](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "fast-fill"}, cb)
}

type WsFastFill struct {
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	ExecID        string          `json:"execId"`
	Side          Side            `json:"side"`
	HoldSide      PosSide         `json:"holdSide"`
	TradeScope    TradeScope      `json:"tradeScope"`
	ExecPrice     decimal.Decimal `json:"execPrice"`
	ExecQty       decimal.Decimal `json:"execQty"`
	ExecTime      time.Time       `json:"execTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
}

// SubscribeStrategyOrderService -- private "strategy-order" channel.
type SubscribeStrategyOrderService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeStrategyOrderService() *SubscribeStrategyOrderService {
	return &SubscribeStrategyOrderService{c: c}
}

func (s *SubscribeStrategyOrderService) Do(ctx context.Context, cb WsHandler[WsStrategyOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsStrategyOrder](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "strategy-order"}, cb)
}

type WsStrategyOrder struct {
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	OrderID       string          `json:"orderId"`
	ClientOrderID string          `json:"clientOid"`
	Qty           decimal.Decimal `json:"qty"`
	Side          Side            `json:"side"`
	PosSide       PosSide         `json:"posSide"`
	Type          StrategyType    `json:"type"`
	Status        string          `json:"status"`
	TriggerType   string          `json:"triggerType"`
	TakeProfit    decimal.Decimal `json:"takeProfit"`
	StopLoss      decimal.Decimal `json:"stopLoss"`
	TpOrderType   OrderType       `json:"tpOrderType"`
	SlOrderType   OrderType       `json:"slOrderType"`
	TriggerPrice  decimal.Decimal `json:"triggerPrice"`
	CreatedTime   time.Time       `json:"createdTime"`
	UpdatedTime   time.Time       `json:"updatedTime"`
}

// SubscribeADLNotificationService -- private "adl-notification" channel.
type SubscribeADLNotificationService struct {
	c *UTAWebSocketClient
}

func (c *UTAWebSocketClient) NewSubscribeADLNotificationService() *SubscribeADLNotificationService {
	return &SubscribeADLNotificationService{c: c}
}

func (s *SubscribeADLNotificationService) Do(ctx context.Context, cb WsHandler[WsADLNotification]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsADLNotification](ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "adl-notification"}, cb)
}

type WsADLNotification struct {
	Symbol string          `json:"symbol"`
	Side   Side            `json:"side"`
	Status string          `json:"status"`
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
	Ts     time.Time       `json:"ts"`
}
