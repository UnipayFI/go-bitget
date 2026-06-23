package ws

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/shopspring/decimal"
)

// This file wraps the classic v2 public futures (mix) market channels. Every
// New<Service> constructor takes an InstType so callers pick the product line
// (USDT-FUTURES, COIN-FUTURES, or USDC-FUTURES); pass InstTypeUSDTFutures for
// the common USDT-margined perpetuals.

// SubscribeMixTickerService -- public "ticker" channel (futures).
//
// docs: https://www.bitget.com/api-doc/contract/websocket/public/Tickers-Channel
type SubscribeMixTickerService struct {
	c        *WebSocketClient
	instType InstType
	symbol   string
}

func (c *WebSocketClient) NewSubscribeMixTickerService(instType InstType, symbol string) *SubscribeMixTickerService {
	return &SubscribeMixTickerService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeMixTickerService) Do(ctx context.Context, cb WsHandler[MixWsTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsTicker](ctx, s.c, false,
		WsArg{InstType: string(s.instType), Channel: "ticker", InstId: s.symbol}, cb)
}

// MixWsTicker is one element of the futures "ticker" channel push data array.
type MixWsTicker struct {
	InstId          string          `json:"instId"`          // product id, e.g. BTCUSDT
	LastPr          decimal.Decimal `json:"lastPr"`          // most recent transaction price
	AskPr           decimal.Decimal `json:"askPr"`           // best ask price
	BidPr           decimal.Decimal `json:"bidPr"`           // best bid price
	BidSz           decimal.Decimal `json:"bidSz"`           // best bid size
	AskSz           decimal.Decimal `json:"askSz"`           // best ask size
	High24h         decimal.Decimal `json:"high24h"`         // 24h high
	Low24h          decimal.Decimal `json:"low24h"`          // 24h low
	Change24h       decimal.Decimal `json:"change24h"`       // 24h change ratio
	FundingRate     decimal.Decimal `json:"fundingRate"`     // current funding rate
	NextFundingTime time.Time       `json:"nextFundingTime"` // next funding settlement time (ms)
	MarkPrice       decimal.Decimal `json:"markPrice"`       // mark price
	IndexPrice      decimal.Decimal `json:"indexPrice"`      // index price
	HoldingAmount   decimal.Decimal `json:"holdingAmount"`   // open interest amount
	BaseVolume      decimal.Decimal `json:"baseVolume"`      // base coin trading volume
	QuoteVolume     decimal.Decimal `json:"quoteVolume"`     // quote currency volume
	OpenUtc         decimal.Decimal `json:"openUtc"`         // 00:00 UTC opening price
	SymbolType      string          `json:"symbolType"`      // "1" perpetual, "2" delivery
	Symbol          string          `json:"symbol"`          // trading pair name
	DeliveryPrice   decimal.Decimal `json:"deliveryPrice"`   // delivery price (0 for perpetual)
	Open24h         decimal.Decimal `json:"open24h"`         // price 24h ago
	Ts              time.Time       `json:"ts"`              // data timestamp (ms)
}

// SubscribeMixCandleService -- public candlestick channel (futures). The
// interval is encoded in the channel name, e.g. "1m","5m","15m","30m","1H",
// "4H","6H","12H","1D","3D","1W","1M" (and their "...utc" variants); the channel
// sent is "candle"+interval.
//
// docs: https://www.bitget.com/api-doc/contract/websocket/public/Candlesticks-Channel
type SubscribeMixCandleService struct {
	c        *WebSocketClient
	instType InstType
	symbol   string
	interval string
}

func (c *WebSocketClient) NewSubscribeMixCandleService(instType InstType, symbol, interval string) *SubscribeMixCandleService {
	return &SubscribeMixCandleService{c: c, instType: instType, symbol: symbol, interval: interval}
}

func (s *SubscribeMixCandleService) Do(ctx context.Context, cb WsHandler[MixWsCandle]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsCandle](ctx, s.c, false,
		WsArg{InstType: string(s.instType), Channel: "candle" + s.interval, InstId: s.symbol}, cb)
}

// MixWsCandle is one candlestick row. Bitget pushes each candle as a
// fixed-position JSON array of strings
// ([ts, open, high, low, close, baseVolume, quoteVolume, usdtVolume]);
// MixWsCandle parses that array into named fields and re-emits the same shape.
type MixWsCandle struct {
	Ts          time.Time       `json:"ts"`          // array[0] -- candle start time (ms)
	Open        decimal.Decimal `json:"open"`        // array[1]
	High        decimal.Decimal `json:"high"`        // array[2]
	Low         decimal.Decimal `json:"low"`         // array[3]
	Close       decimal.Decimal `json:"close"`       // array[4]
	BaseVolume  decimal.Decimal `json:"baseVolume"`  // array[5] -- base coin volume
	QuoteVolume decimal.Decimal `json:"quoteVolume"` // array[6] -- quote currency volume
	USDTVolume  decimal.Decimal `json:"usdtVolume"`  // array[7] -- USDT-denominated volume
}

// UnmarshalJSON decodes the positional candle array into named fields. The USDT
// volume column (array[7]) is optional and left zero when absent.
func (k *MixWsCandle) UnmarshalJSON(data []byte) error {
	var row []string
	if err := common.JSONUnmarshal(data, &row); err != nil {
		return err
	}
	if len(row) < 7 {
		return fmt.Errorf("mix: candle has %d columns, want at least 7", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("mix: candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	cols := []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close, &k.BaseVolume, &k.QuoteVolume, &k.USDTVolume}
	for i, dst := range cols {
		idx := i + 1
		if idx >= len(row) {
			break
		}
		d, err := decimal.NewFromString(row[idx])
		if err != nil {
			return fmt.Errorf("mix: candle column %d %q: %w", idx, row[idx], err)
		}
		*dst = d
	}
	return nil
}

// MarshalJSON re-emits the candle as the positional array Bitget sends, so the
// round-trip preserves the wire shape.
func (k MixWsCandle) MarshalJSON() ([]byte, error) {
	row := []string{
		strconv.FormatInt(k.Ts.UnixMilli(), 10),
		k.Open.String(),
		k.High.String(),
		k.Low.String(),
		k.Close.String(),
		k.BaseVolume.String(),
		k.QuoteVolume.String(),
		k.USDTVolume.String(),
	}
	return common.JSONMarshal(row)
}

// SubscribeMixOrderBookService -- public depth channel (futures). depth selects
// the channel: "books" (full, 150ms), "books1" (10ms), "books5" or "books15".
//
// docs: https://www.bitget.com/api-doc/contract/websocket/public/Order-Book-Channel
type SubscribeMixOrderBookService struct {
	c        *WebSocketClient
	instType InstType
	symbol   string
	depth    string
}

func (c *WebSocketClient) NewSubscribeMixOrderBookService(instType InstType, symbol string) *SubscribeMixOrderBookService {
	return &SubscribeMixOrderBookService{c: c, instType: instType, symbol: symbol, depth: "books"}
}

// SetDepth selects the depth channel: "books", "books1", "books5", or "books15".
func (s *SubscribeMixOrderBookService) SetDepth(depth string) *SubscribeMixOrderBookService {
	s.depth = depth
	return s
}

func (s *SubscribeMixOrderBookService) Do(ctx context.Context, cb WsHandler[MixWsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsOrderBook](ctx, s.c, false,
		WsArg{InstType: string(s.instType), Channel: s.depth, InstId: s.symbol}, cb)
}

// MixWsOrderBook is one element of the futures depth channel push data array.
// Asks and bids arrive as arrays of [price, size] string pairs.
type MixWsOrderBook struct {
	Asks     [][]decimal.Decimal `json:"asks"`     // seller depth, [price, size] pairs
	Bids     [][]decimal.Decimal `json:"bids"`     // buyer depth, [price, size] pairs
	Checksum int32               `json:"checksum"` // CRC32 checksum for validation
	Ts       time.Time           `json:"ts"`       // match engine timestamp (ms)
	Seq      int64               `json:"seq"`      // serial number, increments per update
}

// SubscribeMixTradeService -- public "trade" channel (tick-by-tick fills).
//
// docs: https://www.bitget.com/api-doc/contract/websocket/public/New-Trades-Channel
type SubscribeMixTradeService struct {
	c        *WebSocketClient
	instType InstType
	symbol   string
}

func (c *WebSocketClient) NewSubscribeMixTradeService(instType InstType, symbol string) *SubscribeMixTradeService {
	return &SubscribeMixTradeService{c: c, instType: instType, symbol: symbol}
}

func (s *SubscribeMixTradeService) Do(ctx context.Context, cb WsHandler[MixWsTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MixWsTrade](ctx, s.c, false,
		WsArg{InstType: string(s.instType), Channel: "trade", InstId: s.symbol}, cb)
}

// MixWsTrade is one element of the futures "trade" channel push data array.
type MixWsTrade struct {
	Ts      time.Time       `json:"ts"`      // fill time (ms)
	Price   decimal.Decimal `json:"price"`   // filled price
	Size    decimal.Decimal `json:"size"`    // filled quantity
	Side    string          `json:"side"`    // trade direction: "buy" or "sell"
	TradeId string          `json:"tradeId"` // trade identifier
}
