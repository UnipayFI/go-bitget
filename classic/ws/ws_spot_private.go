package ws

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// SubscribeSpotAccountService -- private "account" channel (spot balances).
//
// The channel is coin-scoped via the "coin" arg; Bitget currently only accepts
// "default" (all coins), so the service hardcodes it.
type SubscribeSpotAccountService struct {
	c *WebSocketClient
}

func (c *WebSocketClient) NewSubscribeSpotAccountService() *SubscribeSpotAccountService {
	return &SubscribeSpotAccountService{c: c}
}

func (s *SubscribeSpotAccountService) Do(ctx context.Context, cb WsHandler[SpotWsAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsAccount](ctx, s.c, true,
		WsArg{InstType: string(InstTypeSpot), Channel: "account", Coin: "default"}, cb)
}

// SpotWsAccount is one balance entry pushed on the spot "account" channel.
type SpotWsAccount struct {
	Coin           string          `json:"coin"`           // token/asset name
	Available      decimal.Decimal `json:"available"`      // spendable balance
	Frozen         decimal.Decimal `json:"frozen"`         // frozen (e.g. by open orders)
	Locked         decimal.Decimal `json:"locked"`         // locked assets
	LimitAvailable decimal.Decimal `json:"limitAvailable"` // restricted balance for spot copy trading
	UTime          time.Time       `json:"uTime"`          // last modification time
}

// SubscribeSpotOrdersService -- private "orders" channel (order lifecycle).
//
// instId narrows the subscription to a single symbol; pass "default" to receive
// updates for all symbols.
type SubscribeSpotOrdersService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotOrdersService(symbol string) *SubscribeSpotOrdersService {
	return &SubscribeSpotOrdersService{c: c, symbol: symbol}
}

func (s *SubscribeSpotOrdersService) Do(ctx context.Context, cb WsHandler[SpotWsOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsOrder](ctx, s.c, true,
		WsArg{InstType: string(InstTypeSpot), Channel: "orders", InstID: s.symbol}, cb)
}

// SpotWsOrder is one order update pushed on the spot "orders" channel.
type SpotWsOrder struct {
	InstID           string           `json:"instId"`           // product id, e.g. BTCUSDT
	OrderID          string           `json:"orderId"`          // order id
	ClientOrderID    string           `json:"clientOid"`        // user-specified order id
	Price            decimal.Decimal  `json:"price"`            // order price
	Size             decimal.Decimal  `json:"size"`             // order quantity (quote for buy, base for sell)
	NewSize          decimal.Decimal  `json:"newSize"`          // normalized quantity per order-type rules
	Notional         decimal.Decimal  `json:"notional"`         // purchase amount at market price
	OrderType        string           `json:"orderType"`        // market / limit
	Force            string           `json:"force"`            // GTC, post_only, FOK, IOC
	Side             string           `json:"side"`             // buy / sell
	FillPrice        decimal.Decimal  `json:"fillPrice"`        // most recent execution price
	TradeID          string           `json:"tradeId"`          // most recent trade id
	BaseVolume       decimal.Decimal  `json:"baseVolume"`       // quantity from latest execution
	FillTime         time.Time        `json:"fillTime"`         // latest transaction time (ms)
	FillFee          decimal.Decimal  `json:"fillFee"`          // latest transaction fee (negative)
	FillFeeCoin      string           `json:"fillFeeCoin"`      // fee currency
	TradeScope       string           `json:"tradeScope"`       // T (taker) / M (maker)
	AccBaseVolume    decimal.Decimal  `json:"accBaseVolume"`    // cumulative filled quantity
	PriceAvg         decimal.Decimal  `json:"priceAvg"`         // weighted average execution price
	Status           string           `json:"status"`           // live, partially_filled, filled, cancelled
	EnterPointSource string           `json:"enterPointSource"` // order origin source
	CTime            time.Time        `json:"cTime"`            // creation time (ms)
	UTime            time.Time        `json:"uTime"`            // update time (ms)
	StpMode          string           `json:"stpMode"`          // none, cancel_taker, cancel_maker, cancel_both
	FeeDetail        []SpotWsOrderFee `json:"feeDetail"`        // fee breakdown
}

// SpotWsOrderFee is one fee line in a spot order's fee breakdown.
type SpotWsOrderFee struct {
	FeeCoin string          `json:"feeCoin"` // fee currency
	Fee     decimal.Decimal `json:"fee"`     // charged fee amount
}

// SubscribeSpotFillService -- private "fill" channel (real-time executions).
//
// instId narrows the subscription to a single symbol; pass "default" to receive
// fills for all symbols.
type SubscribeSpotFillService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotFillService(symbol string) *SubscribeSpotFillService {
	return &SubscribeSpotFillService{c: c, symbol: symbol}
}

func (s *SubscribeSpotFillService) Do(ctx context.Context, cb WsHandler[SpotWsFill]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsFill](ctx, s.c, true,
		WsArg{InstType: string(InstTypeSpot), Channel: "fill", InstID: s.symbol}, cb)
}

// SpotWsFill is one execution pushed on the spot "fill" channel.
type SpotWsFill struct {
	OrderID    string          `json:"orderId"`    // order id
	TradeID    string          `json:"tradeId"`    // trade id
	Symbol     string          `json:"symbol"`     // trading pair symbol
	OrderType  string          `json:"orderType"`  // limit / market
	Side       string          `json:"side"`       // buy / sell
	PriceAvg   decimal.Decimal `json:"priceAvg"`   // total average filled price
	Size       decimal.Decimal `json:"size"`       // filled quantity
	Amount     decimal.Decimal `json:"amount"`     // accumulated filled size
	TradeScope string          `json:"tradeScope"` // taker / maker
	FeeDetail  []SpotWsFillFee `json:"feeDetail"`  // fee breakdown
	CTime      time.Time       `json:"cTime"`      // creation time (ms)
	UTime      time.Time       `json:"uTime"`      // update time (ms)
}

// SpotWsFillFee is one fee line in a spot fill's fee breakdown.
type SpotWsFillFee struct {
	FeeCoin           string          `json:"feeCoin"`           // fee currency
	Deduction         string          `json:"deduction"`         // yes / no
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"` // deducted fee amount
	TotalFee          decimal.Decimal `json:"totalFee"`          // total fee charged
}

// SubscribeSpotOrdersAlgoService -- private "orders-algo" channel (plan/trigger
// orders).
//
// instId narrows the subscription to a single symbol; pass "default" to receive
// updates for all symbols.
type SubscribeSpotOrdersAlgoService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotOrdersAlgoService(symbol string) *SubscribeSpotOrdersAlgoService {
	return &SubscribeSpotOrdersAlgoService{c: c, symbol: symbol}
}

func (s *SubscribeSpotOrdersAlgoService) Do(ctx context.Context, cb WsHandler[SpotWsOrderAlgo]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsOrderAlgo](ctx, s.c, true,
		WsArg{InstType: string(InstTypeSpot), Channel: "orders-algo", InstID: s.symbol}, cb)
}

// SpotWsOrderAlgo is one plan/trigger order update pushed on the spot
// "orders-algo" channel.
type SpotWsOrderAlgo struct {
	InstID           string          `json:"instId"`           // product id
	OrderID          string          `json:"orderId"`          // plan order id
	ClientOrderID    string          `json:"clientOid"`        // customized plan order id
	TriggerPrice     decimal.Decimal `json:"triggerPrice"`     // trigger price
	TriggerType      string          `json:"triggerType"`      // fill_price / mark_price
	PlanType         string          `json:"planType"`         // amount / total
	Price            decimal.Decimal `json:"price"`            // order price
	Size             decimal.Decimal `json:"size"`             // original order amount in coin
	ActualSize       decimal.Decimal `json:"actualSize"`       // actual number of orders in coin
	OrderType        string          `json:"orderType"`        // limit / market
	Side             string          `json:"side"`             // order direction
	Status           string          `json:"status"`           // order status
	ExecutePrice     decimal.Decimal `json:"executePrice"`     // execute price
	EnterPointSource string          `json:"enterPointSource"` // WEB, API, SYS, ANDROID, IOS
	CTime            time.Time       `json:"cTime"`            // create time (ms)
	UTime            time.Time       `json:"uTime"`            // update time (ms)
	StpMode          string          `json:"stpMode"`          // none, cancel_taker, cancel_maker, cancel_both
}
