package uta

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
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

// WsTicker is a symbol's snapshot pushed on the "ticker" channel. Its items
// carry neither the symbol nor a timestamp: read them from the envelope's
// WsPush.Arg.Symbol and WsPush.Ts.
type WsTicker struct {
	LastPrice    decimal.Decimal `json:"lastPrice"`
	OpenPrice24h decimal.Decimal `json:"openPrice24h"`
	HighPrice24h decimal.Decimal `json:"highPrice24h"`
	LowPrice24h  decimal.Decimal `json:"lowPrice24h"`
	Ask1Price    decimal.Decimal `json:"ask1Price"`
	Bid1Price    decimal.Decimal `json:"bid1Price"`
	Bid1Size     decimal.Decimal `json:"bid1Size"`
	Ask1Size     decimal.Decimal `json:"ask1Size"`
	Price24hPcnt decimal.Decimal `json:"price24hPcnt"`
	Volume24h    decimal.Decimal `json:"volume24h"`
	Turnover24h  decimal.Decimal `json:"turnover24h"`
	// PlatformTurnover24h is the 24h platform volume, only available for rtoken.
	PlatformTurnover24h decimal.Decimal `json:"platformTurnover24h"`
	IndexPrice          decimal.Decimal `json:"indexPrice"`
	MarkPrice           decimal.Decimal `json:"markPrice"`
	FundingRate         decimal.Decimal `json:"fundingRate"`
	NextFundingTime     time.Time       `json:"nextFundingTime,format:unixmilli"` // perp: next funding settlement time
	OpenInterest        decimal.Decimal `json:"openInterest"`
	DeliveryStart       time.Time       `json:"deliveryStartTime,format:unixmilli"`
	DeliveryTime        time.Time       `json:"deliveryTime,format:unixmilli"`
	DeliveryStatus      string          `json:"deliveryStatus"`
}

// SubscribeKlineService -- public "kline" candlestick channel. Reality stock
// (rtoken) symbols such as rAAPLUSDT are also supported.
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
	Start    time.Time       `json:"start,format:unixmilli"`
	Open     decimal.Decimal `json:"open"`
	Close    decimal.Decimal `json:"close"`
	High     decimal.Decimal `json:"high"`
	Low      decimal.Decimal `json:"low"`
	Volume   decimal.Decimal `json:"volume"`
	Turnover decimal.Decimal `json:"turnover"`
}

// SubscribeOrderBookService -- public depth channel. depth selects the topic:
// "books" (full), "books1", "books5", or "books50".
//
// For Reality (rtoken) pairs, books1 pushes at all times, while books5,
// books50 and books push only outside US market hours (internal depth only),
// so they stay silent during weekday market hours. Once a pair migrates to the
// new rToken matching engine (2026-10-14 to 2026-10-20) its book may briefly
// show crossed prices. Use SubscribeRealityOrderBookService for full rtoken
// depth.
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
	Ts       time.Time           `json:"ts,format:unixmilli"`
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
	Ts         time.Time       `json:"T,format:unixmilli"`
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
	Ts   time.Time           `json:"ts,format:unixmilli"`
}

// SubscribeRealityOrderBookService -- "reality-orderbook" channel: raw
// (unaggregated) depth for a Reality (US stock) trading pair, e.g. RAAPLUSDT,
// merging onshore (internal) and offshore (external) quantities; see
// SubscribeRealityOrderBookSegmentedService for the split. Every push is a
// full snapshot of up to 85 levels, with complete depth in both trading and
// non-trading sessions, sent every 500ms while the depth changes or every 10s
// otherwise. Once the pair migrates to the new rToken matching engine
// (2026-10-14 to 2026-10-20) the book may briefly show crossed prices.
//
// This channel requires API Key authentication and whitelist access (apply via
// your BD contact), so the subscription is sent over the authenticated
// connection.
type SubscribeRealityOrderBookService struct {
	c      *UTAWebSocketClient
	symbol string
}

func (c *UTAWebSocketClient) NewSubscribeRealityOrderBookService(symbol string) *SubscribeRealityOrderBookService {
	return &SubscribeRealityOrderBookService{c: c, symbol: symbol}
}

func (s *SubscribeRealityOrderBookService) Do(ctx context.Context, cb WsHandler[WsRealityOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return subscribeObjectOrArray(ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "reality-orderbook", Symbol: s.symbol}, cb)
}

// WsRealityOrderBook rows are [price, size] pairs: asks ascending, bids
// descending.
type WsRealityOrderBook struct {
	Symbol string              `json:"symbol"`
	Asks   [][]decimal.Decimal `json:"a"`
	Bids   [][]decimal.Decimal `json:"b"`
}

// SubscribeRealityOrderBookSegmentedService -- "reality-orderbook-segmented"
// channel: the reality-orderbook depth with each level's onshore (internal)
// and offshore (external) quantities reported separately. Every push is a full
// snapshot of up to 85 levels, sent every 500ms while the depth changes or
// every 10s otherwise; the book may briefly show crossed prices.
// Only pairs migrated to the new rToken matching engine push data (live from
// 2026-10-15).
//
// This channel requires API Key authentication and whitelist access (apply via
// your BD contact), so the subscription is sent over the authenticated
// connection.
type SubscribeRealityOrderBookSegmentedService struct {
	c      *UTAWebSocketClient
	symbol string
}

func (c *UTAWebSocketClient) NewSubscribeRealityOrderBookSegmentedService(symbol string) *SubscribeRealityOrderBookSegmentedService {
	return &SubscribeRealityOrderBookSegmentedService{c: c, symbol: symbol}
}

func (s *SubscribeRealityOrderBookSegmentedService) Do(ctx context.Context, cb WsHandler[WsRealityOrderBookSegmented]) (chan<- struct{}, <-chan struct{}, error) {
	return subscribeObjectOrArray(ctx, s.c, true,
		request.WsArg{InstType: wsInstTypeUTA, Topic: "reality-orderbook-segmented", Symbol: s.symbol}, cb)
}

// WsRealityOrderBookSegmented rows are [price, onshore size, offshore size]
// triples: asks ascending, bids descending.
type WsRealityOrderBookSegmented struct {
	Asks [][]decimal.Decimal `json:"a"`
	Bids [][]decimal.Decimal `json:"b"`
}

// subscribeObjectOrArray is request.Subscribe for channels whose data arrives
// (per the docs or live) as a single object instead of the usual array: it
// decodes either shape and always hands cb a slice.
func subscribeObjectOrArray[T any](ctx context.Context, c *UTAWebSocketClient, private bool, arg request.WsArg, cb WsHandler[T]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe(ctx, c, private, arg, func(p *request.WsPush[objectOrArray[T]], err error) {
		if err != nil {
			cb(nil, err)
			return
		}
		cb(&request.WsPush[[]T]{Action: p.Action, Arg: p.Arg, Data: p.Data, Ts: p.Ts}, nil)
	})
}

// objectOrArray decodes a JSON array of T, or a single T object as a
// one-element slice.
type objectOrArray[T any] []T

func (s *objectOrArray[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != '{' {
		return json.UnmarshalDecode(dec, (*[]T)(s))
	}
	var v T
	if err := json.UnmarshalDecode(dec, &v); err != nil {
		return err
	}
	*s = objectOrArray[T]{v}
	return nil
}

// SubscribeLiquidationService -- public "liquidation" channel (futures only).
//
// Pushes once per second, carrying only the highest-quantity liquidation record
// for the long and short side of each trading pair in that second.
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
	Ts     time.Time       `json:"ts,format:unixmilli"`
}

// SubscribeInstrumentsService -- public "instruments" channel: pushes a
// product line's trading parameters (price ratios, TP/SL price protection,
// precisions, order limits) and trading-pair status whenever they change.
// Nothing is pushed on subscribe, so fetch the full set with
// NewGetInstrumentsService first and apply these pushes as changes.
//
// Modelled from the docs only: no push was seen live, so the shape is
// unverified and data is decoded whether sent as an array or a single object.
type SubscribeInstrumentsService struct {
	c        *UTAWebSocketClient
	instType WsInstType
}

func (c *UTAWebSocketClient) NewSubscribeInstrumentsService(instType WsInstType) *SubscribeInstrumentsService {
	return &SubscribeInstrumentsService{c: c, instType: instType}
}

func (s *SubscribeInstrumentsService) Do(ctx context.Context, cb WsHandler[WsInstrument]) (chan<- struct{}, <-chan struct{}, error) {
	return subscribeObjectOrArray(ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "instruments"}, cb)
}

// WsInstrument items carry no symbol per the docs: identify the pair by
// BaseCoin and QuoteCoin.
type WsInstrument struct {
	BaseCoin            string          `json:"baseCoin"`
	QuoteCoin           string          `json:"quoteCoin"`
	BuyLimitPriceRatio  decimal.Decimal `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio decimal.Decimal `json:"sellLimitPriceRatio"`
	MarketPriceRatio    decimal.Decimal `json:"marketPriceRatio"`
	// PriceProtectThreshold is the max TP/SL trigger price deviation from the
	// last/mark price, as a fraction (0.1 means 10%).
	PriceProtectThreshold decimal.Decimal  `json:"priceProtectThreshold"`
	FeeRateUpRatio        decimal.Decimal  `json:"feeRateUpRatio"`
	OpenCostUpRatio       decimal.Decimal  `json:"openCostUpRatio"`
	MinOrderQty           decimal.Decimal  `json:"minOrderQty"`
	MaxOrderQty           decimal.Decimal  `json:"maxOrderQty"` // 0 means no limit
	MinOrderAmount        decimal.Decimal  `json:"minOrderAmount"`
	MaxMarketOrderQty     decimal.Decimal  `json:"maxMarketOrderQty"`
	PricePrecision        decimal.Decimal  `json:"pricePrecision"`
	QuantityPrecision     decimal.Decimal  `json:"quantityPrecision"`
	QuotePrecision        decimal.Decimal  `json:"quotePrecision"`
	PriceMultiplier       decimal.Decimal  `json:"priceMultiplier"`
	QuantityMultiplier    decimal.Decimal  `json:"quantityMultiplier"`
	Type                  SymbolType       `json:"type"` // futures: perpetual, delivery
	Status                InstrumentStatus `json:"status"`
	OpenTime              time.Time        `json:"openTime,format:unixmilli"`
	OffTime               time.Time        `json:"offTime,format:unixmilli"` // futures only; zero when not configured
	// LimitOpenTime, when set, is when the pair enters maintenance and
	// trading is disabled.
	LimitOpenTime time.Time `json:"limitOpenTime,format:unixmilli"`
}

// SubscribeOpenInterestLimitService -- public "open-interest-limit" channel
// (futures only): pushes a full snapshot of a symbol's single-user open
// interest limit on subscribe, then again whenever the limit changes. Stock
// contracts allow up to 20% per UID and 40% across a master account and its
// sub-accounts.
type SubscribeOpenInterestLimitService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribeOpenInterestLimitService(instType WsInstType, symbol string) *SubscribeOpenInterestLimitService {
	return &SubscribeOpenInterestLimitService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeOpenInterestLimitService) Do(ctx context.Context, cb WsHandler[WsOpenInterestLimit]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsOpenInterestLimit](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "open-interest-limit", Symbol: s.symbol}, cb)
}

type WsOpenInterestLimit struct {
	Symbol string `json:"symbol"`
	// The live wire key is singleNotionalValue; the docs list it as
	// notionalValue.
	SingleNotionalValue decimal.Decimal `json:"singleNotionalValue"`
	TotalNotionalValue  decimal.Decimal `json:"totalNotionalValue"`
	// Position-ratio caps as fractions (0.1 means 10%): per UID, and across a
	// master account and its sub-accounts.
	SinglePositionRatio decimal.Decimal `json:"singlePositionRatio"`
	TotalPositionRatio  decimal.Decimal `json:"totalPositionRatio"`
}

// SubscribePositionTierService -- public "position-tier" channel (futures
// only): pushes a full snapshot of a symbol's position tiers on subscribe,
// then again whenever a tier's value range changes. Data arrives live as a
// single object (the docs show an array).
type SubscribePositionTierService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribePositionTierService(instType WsInstType, symbol string) *SubscribePositionTierService {
	return &SubscribePositionTierService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribePositionTierService) Do(ctx context.Context, cb WsHandler[WsPositionTier]) (chan<- struct{}, <-chan struct{}, error) {
	return subscribeObjectOrArray(ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "position-tier", Symbol: s.symbol}, cb)
}

type WsPositionTier struct {
	Symbol string         `json:"symbol"`
	List   []PositionTier `json:"list"`
}

// SubscribeFundingRateService -- public "funding-rate" channel (futures only):
// pushes a full snapshot on subscribe, again when the settlement period
// temporarily switches to 1h or reverts, on a cash dividend event, and as a
// heartbeat every 60s.
type SubscribeFundingRateService struct {
	c        *UTAWebSocketClient
	instType WsInstType
	symbol   string
}

func (c *UTAWebSocketClient) NewSubscribeFundingRateService(instType WsInstType, symbol string) *SubscribeFundingRateService {
	return &SubscribeFundingRateService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeFundingRateService) Do(ctx context.Context, cb WsHandler[WsFundingRate]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsFundingRate](ctx, s.c, false,
		request.WsArg{InstType: string(s.instType), Topic: "funding-rate", Symbol: s.symbol}, cb)
}

type WsFundingRate struct {
	Symbol              string          `json:"symbol"`
	FundingRate         decimal.Decimal `json:"fundingRate"`
	FundingRateInterval string          `json:"fundingRateInterval"` // settlement period in hours
	NextUpdate          time.Time       `json:"nextUpdate,format:unixmilli"`
	MinFundingRate      decimal.Decimal `json:"minFundingRate"`
	MaxFundingRate      decimal.Decimal `json:"maxFundingRate"`
	LastFundingRate     decimal.Decimal `json:"lastFundingRate"`
	NextFundingRate     decimal.Decimal `json:"nextFundingRate"`
	// Unlike the REST current-fund-rate response, the push spells this key
	// correctly.
	BaseInterestRate     decimal.Decimal `json:"baseInterestRate"`
	ImpactMarginNotional decimal.Decimal `json:"impactMarginNotional"` // quote coin
	// CashDividend is the per-share USDT dividend paid by RWA stock contracts;
	// it and CashDividendNextUpdate are zero outside a cash dividend event.
	CashDividend           decimal.Decimal `json:"cashDividend"`
	CashDividendNextUpdate time.Time       `json:"cashDividendNextUpdate,format:unixmilli"`
}

// SubscribeMarginLoanService -- public "margin-loan" channel: pushes a full
// snapshot of a coin's margin loan parameters on subscribe, again when its
// loan limits, interest-free quota or interest discount change or margin
// borrowing of the coin is enabled/disabled, and as a heartbeat every 60s.
type SubscribeMarginLoanService struct {
	c    *UTAWebSocketClient
	coin string
}

func (c *UTAWebSocketClient) NewSubscribeMarginLoanService(coin string) *SubscribeMarginLoanService {
	return &SubscribeMarginLoanService{c: c, coin: coin}
}

func (s *SubscribeMarginLoanService) Do(ctx context.Context, cb WsHandler[WsMarginLoan]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsMarginLoan](ctx, s.c, false,
		request.WsArg{InstType: wsInstTypeUTAPublic, Topic: "margin-loan", Coin: s.coin}, cb)
}

// WsMarginLoan rates are fractions (0.01 means 1%); limits and quotas are in
// the loan coin.
type WsMarginLoan struct {
	Coin              string          `json:"coin"`
	DailyInterest     decimal.Decimal `json:"dailyInterest"`
	AnnualInterest    decimal.Decimal `json:"annualInterest"`
	Limit             decimal.Decimal `json:"limit"`             // per-account borrow limit
	MasterSubLimit    decimal.Decimal `json:"masterSubLimit"`    // master/sub account borrow limit
	InterestFreeLimit decimal.Decimal `json:"interestFreeLimit"` // zero when the coin has none
	// The live wire key is platformRemainingQuota; the docs misspell it
	// platformRemaingQuota (as the REST margin-loans response does).
	PlatformRemainingQuota decimal.Decimal `json:"platformRemainingQuota"`
	Status                 string          `json:"status"` // enabled, disabled
}

// SubscribeMarginCoinsService -- public "margin-coins" channel: pushes a full
// snapshot of a coin's collateral (margin coin) settings on subscribe, then
// again when it is enabled/disabled or its discount rate tiers change. Live
// data items are flat margin-coin records; the docs wrap them in a
// {coin, list} object.
type SubscribeMarginCoinsService struct {
	c    *UTAWebSocketClient
	coin string
}

func (c *UTAWebSocketClient) NewSubscribeMarginCoinsService(coin string) *SubscribeMarginCoinsService {
	return &SubscribeMarginCoinsService{c: c, coin: coin}
}

func (s *SubscribeMarginCoinsService) Do(ctx context.Context, cb WsHandler[WsMarginCoin]) (chan<- struct{}, <-chan struct{}, error) {
	return request.Subscribe[[]WsMarginCoin](ctx, s.c, false,
		request.WsArg{InstType: wsInstTypeUTAPublic, Topic: "margin-coins", Coin: s.coin}, cb)
}

type WsMarginCoin struct {
	MarginCoin       string             `json:"marginCoin"`
	DiscountRateList []DiscountRateTier `json:"discountRateList"`
	Status           string             `json:"status"` // enabled, disabled
}
