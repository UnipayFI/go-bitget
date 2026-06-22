package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Data field is already decoded into the channel's typed slice.
type WsHandler[T any] func(*request.WsPush[[]T], error)

// SubscribeTickerService -- public "ticker" channel.
type SubscribeTickerService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribeTickerService(instType WsInstType, symbol string) *SubscribeTickerService {
	return &SubscribeTickerService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeTickerService) Do(ctx context.Context, cb WsHandler[WsTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTicker](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "ticker", Symbol: s.symbol}, cb)
}

type WsTicker struct {
	Symbol        string          `json:"symbol"`
	LastPrice     decimal.Decimal `json:"lastPrice"`
	OpenPrice24h  decimal.Decimal `json:"openPrice24h"`
	HighPrice24h  decimal.Decimal `json:"highPrice24h"`
	LowPrice24h   decimal.Decimal `json:"lowPrice24h"`
	Ask1Price     decimal.Decimal `json:"ask1Price"`
	Bid1Price     decimal.Decimal `json:"bid1Price"`
	Bid1Size      decimal.Decimal `json:"bid1Size"`
	Ask1Size      decimal.Decimal `json:"ask1Size"`
	Price24hPcnt  decimal.Decimal `json:"price24hPcnt"`
	Volume24h     decimal.Decimal `json:"volume24h"`
	Turnover24h   decimal.Decimal `json:"turnover24h"`
	IndexPrice    decimal.Decimal `json:"indexPrice"`
	MarkPrice     decimal.Decimal `json:"markPrice"`
	FundingRate   decimal.Decimal `json:"fundingRate"`
	OpenInterest  decimal.Decimal `json:"openInterest"`
	DeliveryStart time.Time       `json:"deliveryStartTime"`
	DeliveryTime  time.Time       `json:"deliveryTime"`
	Ts            time.Time       `json:"ts"`
}

// SubscribeKlineService -- public "kline" candlestick channel.
type SubscribeKlineService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
	interval KlineGranularity
}

func (c *UTAWebSocketClient) NewSubscribeKlineService(instType WsInstType, symbol string, interval KlineGranularity) *SubscribeKlineService {
	return &SubscribeKlineService{c: c, instType: instType, symbol: symbol, interval: interval}
}

func (s *SubscribeKlineService) Do(ctx context.Context, cb WsHandler[WsKline]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsKline](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "kline", Symbol: s.symbol, Interval: string(s.interval)}, cb)
}

type WsKline struct {
	Start    time.Time       `json:"start"`
	Open     decimal.Decimal `json:"open"`
	Close    decimal.Decimal `json:"close"`
	High     decimal.Decimal `json:"high"`
	Low      decimal.Decimal `json:"low"`
	Volume   decimal.Decimal `json:"volume"`
	Turnover decimal.Decimal `json:"turnover"`
}

// SubscribeOrderBookService -- public depth channel. depth selects the topic:
// "books" (full), "books1", "books5", or "books50".
type SubscribeOrderBookService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
	depth    string
}

func (c *UTAWebSocketClient) NewSubscribeOrderBookService(instType WsInstType, symbol string) *SubscribeOrderBookService {
	return &SubscribeOrderBookService{c: c, instType: instType, symbol: symbol, depth: "books"}
}

// SetDepth selects the depth topic: "books", "books1", "books5", or "books50".
func (s *SubscribeOrderBookService) SetDepth(depth string) *SubscribeOrderBookService {
	s.depth = depth
	return s
}

func (s *SubscribeOrderBookService) Do(ctx context.Context, cb WsHandler[WsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOrderBook](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: s.depth, Symbol: s.symbol}, cb)
}

type WsOrderBook struct {
	Asks     [][]decimal.Decimal `json:"a"`
	Bids     [][]decimal.Decimal `json:"b"`
	Seq      int64               `json:"seq"`
	Pseq     int64               `json:"pseq"`
	MaxDepth string              `json:"maxDepth"`
	Ts       time.Time           `json:"ts"`
}

// SubscribeTradeService -- public "publicTrade" channel (tick-by-tick fills).
type SubscribeTradeService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribeTradeService(instType WsInstType, symbol string) *SubscribeTradeService {
	return &SubscribeTradeService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeTradeService) Do(ctx context.Context, cb WsHandler[WsTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsTrade](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "publicTrade", Symbol: s.symbol}, cb)
}

type WsTrade struct {
	ExecID     string          `json:"i"`
	ExecLinkID string          `json:"L"`
	Price      decimal.Decimal `json:"p"`
	Size       decimal.Decimal `json:"v"`
	Side       Side            `json:"S"`
	Ts         time.Time       `json:"T"`
	IsRPI      string          `json:"isRPI"`
}

// SubscribeRPIOrderBookService -- public RPI depth channel. depth selects the
// topic: "rpi-books" (full), "rpi-books1", "rpi-books5", or "rpi-books50".
type SubscribeRPIOrderBookService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
	depth    string
}

func (c *UTAWebSocketClient) NewSubscribeRPIOrderBookService(instType WsInstType, symbol string) *SubscribeRPIOrderBookService {
	return &SubscribeRPIOrderBookService{c: c, instType: instType, symbol: symbol, depth: "rpi-books5"}
}

// SetDepth selects the depth topic: "rpi-books", "rpi-books1", "rpi-books5", or "rpi-books50".
func (s *SubscribeRPIOrderBookService) SetDepth(depth string) *SubscribeRPIOrderBookService {
	s.depth = depth
	return s
}

func (s *SubscribeRPIOrderBookService) Do(ctx context.Context, cb WsHandler[WsRPIOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsRPIOrderBook](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: s.depth, Symbol: s.symbol}, cb)
}

// WsRPIOrderBook rows are [price, non-RPI quantity, RPI quantity] triples.
type WsRPIOrderBook struct {
	Asks [][]decimal.Decimal `json:"a"`
	Bids [][]decimal.Decimal `json:"b"`
	Seq  int64               `json:"seq"`
	Pseq int64               `json:"pseq"`
	Ts   time.Time           `json:"ts"`
}

// SubscribeLiquidationService -- public "liquidation" channel (futures only).
type SubscribeLiquidationService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribeLiquidationService(instType WsInstType, symbol string) *SubscribeLiquidationService {
	return &SubscribeLiquidationService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeLiquidationService) Do(ctx context.Context, cb WsHandler[WsLiquidation]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsLiquidation](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "liquidation", Symbol: s.symbol}, cb)
}

type WsLiquidation struct {
	Symbol string          `json:"symbol"`
	Side   Side            `json:"side"`
	Price  decimal.Decimal `json:"price"`
	Amount decimal.Decimal `json:"amount"`
	Ts     time.Time       `json:"ts"`
}
