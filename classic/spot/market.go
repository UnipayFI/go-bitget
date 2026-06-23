package spot

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCoinsService -- GET /api/v2/spot/public/coins (public)
//
// Returns coin metadata and the supported deposit/withdraw chains, optionally
// filtered to a single coin.
type GetCoinsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetCoinsService() *GetCoinsService {
	return &GetCoinsService{c: c, params: map[string]string{}}
}

// SetCoin filters the result to a single coin (e.g. "BTC").
func (s *GetCoinsService) SetCoin(coin string) *GetCoinsService {
	s.params["coin"] = coin
	return s
}

func (s *GetCoinsService) Do(ctx context.Context) ([]Coin, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/public/coins", s.params)
	resp, err := request.Do[[]Coin](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Coin is a single coin's metadata together with its per-chain config.
type Coin struct {
	CoinID   string      `json:"coinId"`
	Coin     string      `json:"coin"`
	Transfer string      `json:"transfer"` // "true"/"false": transferable between accounts
	AreaCoin string      `json:"areaCoin"` // "yes"/"no": region-restricted coin
	Chains   []CoinChain `json:"chains"`
}

// CoinChain is the deposit/withdraw configuration for one chain of a coin.
type CoinChain struct {
	Chain             string          `json:"chain"`
	NeedTag           string          `json:"needTag"`      // "true"/"false": deposit requires a memo/tag
	Withdrawable      string          `json:"withdrawable"` // "true"/"false"
	Rechargeable      string          `json:"rechargeable"` // "true"/"false": depositable
	WithdrawFee       decimal.Decimal `json:"withdrawFee"`
	ExtraWithdrawFee  decimal.Decimal `json:"extraWithdrawFee"` // extra on-chain destruction fee
	DepositConfirm    string          `json:"depositConfirm"`   // confirmations to credit a deposit
	WithdrawConfirm   string          `json:"withdrawConfirm"`  // confirmations before withdraw is final
	MinDepositAmount  decimal.Decimal `json:"minDepositAmount"`
	MinWithdrawAmount decimal.Decimal `json:"minWithdrawAmount"`
	BrowserURL        string          `json:"browserUrl"`
	ContractAddress   string          `json:"contractAddress"`
	WithdrawStep      string          `json:"withdrawStep"`     // withdraw amount step (0 = no step limit)
	WithdrawMinScale  string          `json:"withdrawMinScale"` // decimal places allowed on withdraw amount
	Congestion        string          `json:"congestion"`       // chain status: "normal"/"congested"
}

// GetSymbolsService -- GET /api/v2/spot/public/symbols (public)
//
// Returns trading-pair metadata (precision, fees, limits, status), optionally
// filtered to a single symbol.
type GetSymbolsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSymbolsService() *GetSymbolsService {
	return &GetSymbolsService{c: c, params: map[string]string{}}
}

// SetSymbol filters the result to a single trading pair (e.g. "BTCUSDT").
func (s *GetSymbolsService) SetSymbol(symbol string) *GetSymbolsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetSymbolsService) Do(ctx context.Context) ([]Symbol, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/public/symbols", s.params)
	resp, err := request.Do[[]Symbol](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SymbolStatus is the trading state of a symbol.
type SymbolStatus string

const (
	SymbolStatusOffline SymbolStatus = "offline"
	SymbolStatusGray    SymbolStatus = "gray"
	SymbolStatusOnline  SymbolStatus = "online"
	SymbolStatusHalt    SymbolStatus = "halt"
)

// Symbol is a single spot trading pair's configuration.
type Symbol struct {
	Symbol              string          `json:"symbol"`
	BaseCoin            string          `json:"baseCoin"`
	QuoteCoin           string          `json:"quoteCoin"`
	MinTradeAmount      decimal.Decimal `json:"minTradeAmount"`
	MaxTradeAmount      decimal.Decimal `json:"maxTradeAmount"`
	TakerFeeRate        decimal.Decimal `json:"takerFeeRate"`
	MakerFeeRate        decimal.Decimal `json:"makerFeeRate"`
	PricePrecision      string          `json:"pricePrecision"`    // decimal places for price
	QuantityPrecision   string          `json:"quantityPrecision"` // decimal places for base-coin quantity
	QuotePrecision      string          `json:"quotePrecision"`    // decimal places for quote-coin amount
	Status              SymbolStatus    `json:"status"`
	MinTradeUSDT        decimal.Decimal `json:"minTradeUSDT"`
	BuyLimitPriceRatio  decimal.Decimal `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio decimal.Decimal `json:"sellLimitPriceRatio"`
	AreaSymbol          string          `json:"areaSymbol"`    // "yes"/"no": region-restricted pair
	OrderQuantity       string          `json:"orderQuantity"` // max number of open orders per symbol
	OpenTime            time.Time       `json:"openTime"`
	OffTime             time.Time       `json:"offTime"`
	MaxLimitOrderValue  decimal.Decimal `json:"maxLimitOrderValue"`
	MaxMarketOrderValue decimal.Decimal `json:"maxMarketOrderValue"`
}

// GetVIPFeeRateService -- GET /api/v2/spot/market/vip-fee-rate (public)
//
// Returns the spot VIP fee-rate tier schedule.
type GetVIPFeeRateService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetVIPFeeRateService() *GetVIPFeeRateService {
	return &GetVIPFeeRateService{c: c}
}

func (s *GetVIPFeeRateService) Do(ctx context.Context) ([]VIPFeeRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/vip-fee-rate")
	resp, err := request.Do[[]VIPFeeRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// VIPFeeRate is one VIP tier and its associated trade fees and withdraw caps.
type VIPFeeRate struct {
	Level              string          `json:"level"`
	DealAmount         decimal.Decimal `json:"dealAmount"`  // 30-day trade volume threshold (USDT)
	AssetAmount        decimal.Decimal `json:"assetAmount"` // asset balance threshold (USDT)
	TakerFeeRate       decimal.Decimal `json:"takerFeeRate"`
	MakerFeeRate       decimal.Decimal `json:"makerFeeRate"`
	BTCWithdrawAmount  decimal.Decimal `json:"btcWithdrawAmount"`  // 24h BTC-denominated withdraw cap
	USDTWithdrawAmount decimal.Decimal `json:"usdtWithdrawAmount"` // 24h USDT-denominated withdraw cap
}

// GetTickersService -- GET /api/v2/spot/market/tickers (public)
//
// Returns 24h ticker statistics for all symbols, or a single symbol when set.
type GetTickersService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTickersService() *GetTickersService {
	return &GetTickersService{c: c, params: map[string]string{}}
}

// SetSymbol filters the result to a single trading pair (e.g. "BTCUSDT").
func (s *GetTickersService) SetSymbol(symbol string) *GetTickersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetTickersService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/tickers", s.params)
	resp, err := request.Do[[]Ticker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Ticker is a single symbol's 24h ticker snapshot.
type Ticker struct {
	Symbol       string          `json:"symbol"`
	Open         decimal.Decimal `json:"open"` // price 24h ago
	High24h      decimal.Decimal `json:"high24h"`
	Low24h       decimal.Decimal `json:"low24h"`
	LastPr       decimal.Decimal `json:"lastPr"` // last traded price
	QuoteVolume  decimal.Decimal `json:"quoteVolume"`
	BaseVolume   decimal.Decimal `json:"baseVolume"`
	USDTVolume   decimal.Decimal `json:"usdtVolume"`
	Ts           time.Time       `json:"ts"`
	BidPr        decimal.Decimal `json:"bidPr"`   // best bid price
	AskPr        decimal.Decimal `json:"askPr"`   // best ask price
	BidSz        decimal.Decimal `json:"bidSz"`   // best bid size
	AskSz        decimal.Decimal `json:"askSz"`   // best ask size
	OpenUtc      decimal.Decimal `json:"openUtc"` // price at 00:00 UTC
	ChangeUtc24h decimal.Decimal `json:"changeUtc24h"`
	Change24h    decimal.Decimal `json:"change24h"`
}

// MergeDepthPrecision is the price-aggregation step for the merge-depth book.
type MergeDepthPrecision string

const (
	MergeDepthPrecisionScale0 MergeDepthPrecision = "scale0"
	MergeDepthPrecisionScale1 MergeDepthPrecision = "scale1"
	MergeDepthPrecisionScale2 MergeDepthPrecision = "scale2"
	MergeDepthPrecisionScale3 MergeDepthPrecision = "scale3"
)

// GetMergeDepthService -- GET /api/v2/spot/market/merge-depth (public)
//
// Returns an aggregated (merged) order book for a symbol at the chosen price
// precision.
type GetMergeDepthService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetMergeDepthService(symbol string) *GetMergeDepthService {
	return &GetMergeDepthService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetPrecision sets the price-aggregation level (scale0..scale3; default
// scale0).
func (s *GetMergeDepthService) SetPrecision(precision MergeDepthPrecision) *GetMergeDepthService {
	s.params["precision"] = string(precision)
	return s
}

// SetLimit caps the number of levels per side (1, 5, 15, 50, max, default 100).
func (s *GetMergeDepthService) SetLimit(limit string) *GetMergeDepthService {
	s.params["limit"] = limit
	return s
}

func (s *GetMergeDepthService) Do(ctx context.Context) (*MergeDepth, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/merge-depth", s.params)
	return request.Do[MergeDepth](req)
}

// MergeDepth is an aggregated order book snapshot. Asks and bids arrive as
// arrays of [price, size] pairs.
type MergeDepth struct {
	Asks           [][]decimal.Decimal `json:"asks"`
	Bids           [][]decimal.Decimal `json:"bids"`
	Ts             time.Time           `json:"ts"`
	Scale          decimal.Decimal     `json:"scale"`          // numeric price step for the chosen precision
	Precision      string              `json:"precision"`      // echoes the requested precision (scaleN)
	IsMaxPrecision string              `json:"isMaxPrecision"` // "YES"/"NO"
}

// OrderBookType is the price-aggregation step for the raw order book.
type OrderBookType string

const (
	OrderBookTypeStep0 OrderBookType = "step0"
	OrderBookTypeStep1 OrderBookType = "step1"
	OrderBookTypeStep2 OrderBookType = "step2"
	OrderBookTypeStep3 OrderBookType = "step3"
	OrderBookTypeStep4 OrderBookType = "step4"
	OrderBookTypeStep5 OrderBookType = "step5"
)

// GetOrderBookService -- GET /api/v2/spot/market/orderbook (public)
//
// Returns the order book depth (asks and bids) for a symbol.
type GetOrderBookService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetOrderBookService(symbol string) *GetOrderBookService {
	return &GetOrderBookService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetType sets the price-aggregation level (step0..step5; default step0).
func (s *GetOrderBookService) SetType(typ OrderBookType) *GetOrderBookService {
	s.params["type"] = string(typ)
	return s
}

// SetLimit caps the number of levels per side (1..150; default 100).
func (s *GetOrderBookService) SetLimit(limit int) *GetOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/orderbook", s.params)
	return request.Do[OrderBook](req)
}

// OrderBook is the order book depth snapshot. Asks and bids arrive as arrays of
// [price, size] string pairs.
type OrderBook struct {
	Asks [][]decimal.Decimal `json:"asks"`
	Bids [][]decimal.Decimal `json:"bids"`
	Ts   time.Time           `json:"ts"`
}

// Candle is one candlestick row. Bitget returns each candle as a fixed-position
// JSON array ([ts, open, high, low, close, baseVolume, quoteVolume,
// usdtVolume]); Candle parses that array into named fields and re-emits the same
// array shape on marshal.
type Candle struct {
	Ts          time.Time       `json:"ts"`          // array[0] -- candle start time (ms)
	Open        decimal.Decimal `json:"open"`        // array[1]
	High        decimal.Decimal `json:"high"`        // array[2]
	Low         decimal.Decimal `json:"low"`         // array[3]
	Close       decimal.Decimal `json:"close"`       // array[4]
	BaseVolume  decimal.Decimal `json:"baseVolume"`  // array[5] -- base coin volume
	QuoteVolume decimal.Decimal `json:"quoteVolume"` // array[6] -- quote coin turnover
	USDTVolume  decimal.Decimal `json:"usdtVolume"`  // array[7] -- USDT-denominated volume
}

// UnmarshalJSON decodes the 8-element positional array into named fields.
func (k *Candle) UnmarshalJSON(data []byte) error {
	var row []string
	if err := common.JSONUnmarshal(data, &row); err != nil {
		return err
	}
	if len(row) < 8 {
		return fmt.Errorf("spot: candle has %d columns, want 8", len(row))
	}
	ms, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		return fmt.Errorf("spot: candle timestamp %q: %w", row[0], err)
	}
	k.Ts = time.UnixMilli(ms)
	for i, dst := range []*decimal.Decimal{&k.Open, &k.High, &k.Low, &k.Close, &k.BaseVolume, &k.QuoteVolume, &k.USDTVolume} {
		d, err := decimal.NewFromString(row[i+1])
		if err != nil {
			return fmt.Errorf("spot: candle column %d %q: %w", i+1, row[i+1], err)
		}
		*dst = d
	}
	return nil
}

// MarshalJSON re-emits the candle as the positional array Bitget sends, so the
// round-trip preserves the wire shape.
func (k Candle) MarshalJSON() ([]byte, error) {
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

// GetCandlesService -- GET /api/v2/spot/market/candles (public)
//
// Returns the most recent candlesticks for a symbol at the given granularity.
type GetCandlesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetCandlesService(symbol string, granularity string) *GetCandlesService {
	return &GetCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"granularity": granularity,
	}}
}

// SetStartTime filters candles at or after t.
func (s *GetCandlesService) SetStartTime(t time.Time) *GetCandlesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters candles at or before t.
func (s *GetCandlesService) SetEndTime(t time.Time) *GetCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of candles returned (max 1000, default 100).
func (s *GetCandlesService) SetLimit(limit int) *GetCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetHistoryCandlesService -- GET /api/v2/spot/market/history-candles (public)
//
// Returns historical candlesticks for a symbol at the given granularity, ending
// at endTime.
type GetHistoryCandlesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetHistoryCandlesService(symbol string, granularity string) *GetHistoryCandlesService {
	return &GetHistoryCandlesService{c: c, params: map[string]string{
		"symbol":      symbol,
		"granularity": granularity,
	}}
}

// SetEndTime returns candles before this time (required for a populated result).
func (s *GetHistoryCandlesService) SetEndTime(t time.Time) *GetHistoryCandlesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of candles returned (max 200, default 100).
func (s *GetHistoryCandlesService) SetLimit(limit int) *GetHistoryCandlesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryCandlesService) Do(ctx context.Context) ([]Candle, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/history-candles", s.params)
	resp, err := request.Do[[]Candle](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetRecentFillsService -- GET /api/v2/spot/market/fills (public)
//
// Returns the most recent public trades (fills) for a symbol.
type GetRecentFillsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetRecentFillsService(symbol string) *GetRecentFillsService {
	return &GetRecentFillsService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetLimit caps the number of trades returned (max 500, default 100).
func (s *GetRecentFillsService) SetLimit(limit int) *GetRecentFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecentFillsService) Do(ctx context.Context) ([]MarketFill, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/fills", s.params)
	resp, err := request.Do[[]MarketFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarketFill is a single public trade. The fills endpoint reports side in
// lowercase ("buy"/"sell"); the fills-history endpoint reports it capitalized
// ("Buy"/"Sell"), so this is left as a plain string.
type MarketFill struct {
	Symbol  string          `json:"symbol"`
	TradeID string          `json:"tradeId"`
	Side    string          `json:"side"`
	Price   decimal.Decimal `json:"price"`
	Size    decimal.Decimal `json:"size"`
	Ts      time.Time       `json:"ts"`
}

// GetMarketTradesService -- GET /api/v2/spot/market/fills-history (public)
//
// Returns historical public trades (fills) for a symbol, paged via tradeId.
type GetMarketTradesService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetMarketTradesService(symbol string) *GetMarketTradesService {
	return &GetMarketTradesService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetLimit caps the number of trades returned (max 1000, default 500).
func (s *GetMarketTradesService) SetLimit(limit int) *GetMarketTradesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan returns trades with a tradeId less than this value (paging).
func (s *GetMarketTradesService) SetIDLessThan(tradeID string) *GetMarketTradesService {
	s.params["idLessThan"] = tradeID
	return s
}

// SetStartTime filters trades at or after t.
func (s *GetMarketTradesService) SetStartTime(t time.Time) *GetMarketTradesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters trades at or before t.
func (s *GetMarketTradesService) SetEndTime(t time.Time) *GetMarketTradesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMarketTradesService) Do(ctx context.Context) ([]MarketFill, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/fills-history", s.params)
	resp, err := request.Do[[]MarketFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AuctionStage is the call-auction phase for a newly listed symbol.
type AuctionStage string

const (
	AuctionStageCallAuction     AuctionStage = "callAuction"
	AuctionStageNoCancelAuction AuctionStage = "noCancelAuction"
	AuctionStageMatch           AuctionStage = "match"
	AuctionStageContinuousTrade AuctionStage = "continuousTrade"
)

// GetAuctionService -- GET /api/v2/spot/market/auction (public)
//
// Returns the call-auction information for a symbol. All fields are zero/empty
// when the symbol is not in an auction phase.
type GetAuctionService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAuctionService(symbol string) *GetAuctionService {
	return &GetAuctionService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetAuctionService) Do(ctx context.Context) (*Auction, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/market/auction", s.params)
	return request.Do[Auction](req)
}

// Auction is the call-auction snapshot for a symbol.
type Auction struct {
	Stage          AuctionStage    `json:"stage"`
	StageEndTime   time.Time       `json:"stageEndTime"`
	EstOpeningPr   decimal.Decimal `json:"estOpeningPrice"` // estimated opening price
	MatchedVolume  decimal.Decimal `json:"matchedVolume"`
	AuctionEndTime time.Time       `json:"auctionEndTime"`
}
