package ws

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/shopspring/decimal"
)

// SubscribeSpotTickerService -- public "ticker" channel (spot). instType SPOT,
// instId required.
type SubscribeSpotTickerService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotTickerService(symbol string) *SubscribeSpotTickerService {
	return &SubscribeSpotTickerService{c: c, symbol: symbol}
}

func (s *SubscribeSpotTickerService) Do(ctx context.Context, cb WsHandler[SpotWsTicker]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsTicker](ctx, s.c, false,
		WsArg{InstType: string(InstTypeSpot), Channel: "ticker", InstID: s.symbol}, cb)
}

// SpotWsTicker is one element of the spot "ticker" channel push data array.
type SpotWsTicker struct {
	InstID       string          `json:"instId"`       // product id
	LastPr       decimal.Decimal `json:"lastPr"`       // current market price
	AskPr        decimal.Decimal `json:"askPr"`        // best ask price
	BidPr        decimal.Decimal `json:"bidPr"`        // best bid price
	Open24h      decimal.Decimal `json:"open24h"`      // entry price of the last 24 hours
	High24h      decimal.Decimal `json:"high24h"`      // 24h high
	Low24h       decimal.Decimal `json:"low24h"`       // 24h low
	BaseVolume   decimal.Decimal `json:"baseVolume"`   // 24h volume in base (left) coin
	QuoteVolume  decimal.Decimal `json:"quoteVolume"`  // 24h volume in quote (right) coin
	OpenUtc      decimal.Decimal `json:"openUtc"`      // UTC+0 entry price
	ChangeUtc24h decimal.Decimal `json:"changeUtc24h"` // UTC+0 change (0.01 = 1%)
	BidSz        decimal.Decimal `json:"bidSz"`        // best bid size
	AskSz        decimal.Decimal `json:"askSz"`        // best ask size
	Change24h    decimal.Decimal `json:"change24h"`    // 24h change (0.01 = 1%)
	Ts           time.Time       `json:"ts"`           // event timestamp (ms)
}

// SubscribeSpotCandleService -- public candlestick channel (spot). The channel
// name encodes the interval, e.g. "candle1m". instType SPOT, instId required.
//
// Valid intervals: 1m, 5m, 15m, 30m, 1H, 4H, 6H, 12H, 1D, 3D, 1W, 1M (and their
// "...utc" variants, e.g. "1Dutc"). The channel sent is "candle" + interval.
type SubscribeSpotCandleService struct {
	c        *WebSocketClient
	symbol   string
	interval string
}

func (c *WebSocketClient) NewSubscribeSpotCandleService(symbol, interval string) *SubscribeSpotCandleService {
	return &SubscribeSpotCandleService{c: c, symbol: symbol, interval: interval}
}

func (s *SubscribeSpotCandleService) Do(ctx context.Context, cb WsHandler[SpotWsCandle]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsCandle](ctx, s.c, false,
		WsArg{InstType: string(InstTypeSpot), Channel: "candle" + s.interval, InstID: s.symbol}, cb)
}

// SpotWsCandle is one candlestick row. Bitget pushes each candle as a
// fixed-position JSON array of strings
// ([ts, open, high, low, close, baseVolume, quoteVolume, usdtVolume]);
// SpotWsCandle parses that array into named fields and re-emits the same shape.
type SpotWsCandle struct {
	Ts          time.Time       `json:"ts"`          // array[0] -- candle start time (ms)
	Open        decimal.Decimal `json:"open"`        // array[1]
	High        decimal.Decimal `json:"high"`        // array[2]
	Low         decimal.Decimal `json:"low"`         // array[3]
	Close       decimal.Decimal `json:"close"`       // array[4]
	BaseVolume  decimal.Decimal `json:"baseVolume"`  // array[5] -- volume in base coin
	QuoteVolume decimal.Decimal `json:"quoteVolume"` // array[6] -- volume in quote coin
	USDTVolume  decimal.Decimal `json:"usdtVolume"`  // array[7] -- volume in USDT
}

// UnmarshalJSON decodes the 8-element positional array into named fields.
func (k *SpotWsCandle) UnmarshalJSON(data []byte) error {
	var row []string
	if err := common.JSONUnmarshal(data, &row); err != nil {
		return err
	}
	if len(row) < 8 {
		return fmt.Errorf("ws: spot candle has %d columns, want 8", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("ws: spot candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	for i, dst := range []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close, &k.BaseVolume, &k.QuoteVolume, &k.USDTVolume} {
		d, err := decimal.NewFromString(row[i+1])
		if err != nil {
			return fmt.Errorf("ws: spot candle column %d %q: %w", i+1, row[i+1], err)
		}
		*dst = d
	}
	return nil
}

// MarshalJSON re-emits the candle as the positional array Bitget sends.
func (k SpotWsCandle) MarshalJSON() ([]byte, error) {
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

// SubscribeSpotTradeService -- public "trade" channel (spot). instType SPOT,
// instId required.
type SubscribeSpotTradeService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotTradeService(symbol string) *SubscribeSpotTradeService {
	return &SubscribeSpotTradeService{c: c, symbol: symbol}
}

func (s *SubscribeSpotTradeService) Do(ctx context.Context, cb WsHandler[SpotWsTrade]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsTrade](ctx, s.c, false,
		WsArg{InstType: string(InstTypeSpot), Channel: "trade", InstID: s.symbol}, cb)
}

// SpotWsTrade is one element of the spot "trade" channel push data array.
type SpotWsTrade struct {
	Ts      time.Time       `json:"ts"`      // transaction time (ms)
	TradeID string          `json:"tradeId"` // transaction id
	Price   decimal.Decimal `json:"price"`   // transaction price
	Size    decimal.Decimal `json:"size"`    // transaction quantity
	Side    string          `json:"side"`    // "buy" or "sell"
}

// SubscribeSpotOrderBookService -- public depth channel (spot). depth selects
// the channel: "books" (full, snapshot then updates), "books1", "books5" or
// "books15" (full snapshot each push). instType SPOT, instId required.
type SubscribeSpotOrderBookService struct {
	c      *WebSocketClient
	symbol string
	depth  string
}

func (c *WebSocketClient) NewSubscribeSpotOrderBookService(symbol string) *SubscribeSpotOrderBookService {
	return &SubscribeSpotOrderBookService{c: c, symbol: symbol, depth: "books"}
}

// SetDepth selects the depth channel: "books", "books1", "books5" or "books15".
func (s *SubscribeSpotOrderBookService) SetDepth(depth string) *SubscribeSpotOrderBookService {
	s.depth = depth
	return s
}

func (s *SubscribeSpotOrderBookService) Do(ctx context.Context, cb WsHandler[SpotWsOrderBook]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsOrderBook](ctx, s.c, false,
		WsArg{InstType: string(InstTypeSpot), Channel: s.depth, InstID: s.symbol}, cb)
}

// SpotWsOrderBook is one element of the spot depth channel push data array.
// Asks/Bids rows are [price, size] pairs.
type SpotWsOrderBook struct {
	Asks     [][]decimal.Decimal `json:"asks"`     // seller depth [price, size]
	Bids     [][]decimal.Decimal `json:"bids"`     // buyer depth [price, size]
	Checksum int64               `json:"checksum"` // CRC32 validation value
	Seq      int64               `json:"seq"`      // serial number, increments with updates
	Ts       time.Time           `json:"ts"`       // matching engine timestamp (ms)
}

// SubscribeSpotAuctionService -- public "auction" call-auction channel (spot).
// instType SPOT, instId required.
type SubscribeSpotAuctionService struct {
	c      *WebSocketClient
	symbol string
}

func (c *WebSocketClient) NewSubscribeSpotAuctionService(symbol string) *SubscribeSpotAuctionService {
	return &SubscribeSpotAuctionService{c: c, symbol: symbol}
}

func (s *SubscribeSpotAuctionService) Do(ctx context.Context, cb WsHandler[SpotWsAuction]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]SpotWsAuction](ctx, s.c, false,
		WsArg{InstType: string(InstTypeSpot), Channel: "auction", InstID: s.symbol}, cb)
}

// SpotWsAuction is one element of the spot "auction" channel push data array.
type SpotWsAuction struct {
	Stage           string          `json:"stage"`           // pre_market, stage_1, stage_2, stage_3, success, failure
	StageEndTime    time.Time       `json:"stageEndTime"`    // current phase end time (ms)
	EstOpeningPrice decimal.Decimal `json:"estOpeningPrice"` // estimated opening price
	MatchedVolume   decimal.Decimal `json:"matchedVolume"`   // matched volume, base coin
	AuctionEndTime  time.Time       `json:"auctionEndTime"`  // call auction end time (ms)
}
