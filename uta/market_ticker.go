package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetTickersService -- GET /api/v3/market/tickers
//
// Returns 24h ticker statistics for a product category, optionally filtered to a
// single symbol.
type GetTickersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetTickersService(category Category) *GetTickersService {
	return &GetTickersService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetTickersService) SetSymbol(symbol string) *GetTickersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetTickersService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/tickers", s.params)
	resp, err := request.Do[[]Ticker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Ticker is the union of the spot/margin and futures ticker shapes; the
// futures-only fields are empty for SPOT/MARGIN categories.
type Ticker struct {
	Category     Category        `json:"category"`
	Symbol       string          `json:"symbol"`
	Ts           time.Time       `json:"ts"`
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

	// Futures-only fields.
	IndexPrice        decimal.Decimal `json:"indexPrice"`
	MarkPrice         decimal.Decimal `json:"markPrice"`
	FundingRate       decimal.Decimal `json:"fundingRate"`
	OpenInterest      decimal.Decimal `json:"openInterest"`
	DeliveryStartTime time.Time       `json:"deliveryStartTime"`
	DeliveryTime      time.Time       `json:"deliveryTime"`
	DeliveryStatus    string          `json:"deliveryStatus"`
}

// GetOrderBookService -- GET /api/v3/market/orderbook
//
// Returns the order book depth (asks and bids) for a symbol.
type GetOrderBookService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetOrderBookService(category Category, symbol string) *GetOrderBookService {
	return &GetOrderBookService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

// SetLimit sets the depth limit (max 1000).
func (s *GetOrderBookService) SetLimit(limit int) *GetOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetOrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/orderbook", s.params)
	return request.Do[OrderBook](req)
}

// OrderBook is the order book depth snapshot. Asks ("a") and bids ("b") arrive
// as arrays of [price, size] string pairs.
type OrderBook struct {
	Asks [][]decimal.Decimal `json:"a"`
	Bids [][]decimal.Decimal `json:"b"`
	Ts   time.Time           `json:"ts"`
}

// GetRPIOrderBookService -- GET /api/v3/market/rpi-orderbook
//
// Returns the order book depth including RPI (Retail Price Improvement) liquidity
// for a symbol.
type GetRPIOrderBookService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRPIOrderBookService(category Category, symbol string) *GetRPIOrderBookService {
	return &GetRPIOrderBookService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

// SetLimit sets the depth limit (max 1000).
func (s *GetRPIOrderBookService) SetLimit(limit int) *GetRPIOrderBookService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRPIOrderBookService) Do(ctx context.Context) (*RPIOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/rpi-orderbook", s.params)
	return request.Do[RPIOrderBook](req)
}

// RPIOrderBook is the RPI order book depth snapshot. Asks ("a") and bids ("b")
// arrive as arrays of [price, size, rpiSize] string triples.
type RPIOrderBook struct {
	Asks [][]decimal.Decimal `json:"a"`
	Bids [][]decimal.Decimal `json:"b"`
	Ts   time.Time           `json:"ts"`
}

// GetRecentFillsService -- GET /api/v3/market/fills
//
// Returns the most recent public trades (fills) for a symbol.
type GetRecentFillsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRecentFillsService(category Category, symbol string) *GetRecentFillsService {
	return &GetRecentFillsService{c: c, params: map[string]string{
		"category": string(category),
		"symbol":   symbol,
	}}
}

func (s *GetRecentFillsService) SetLimit(limit int) *GetRecentFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecentFillsService) Do(ctx context.Context) ([]PublicFill, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/fills", s.params)
	resp, err := request.Do[[]PublicFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PublicFill is a single recent public trade.
type PublicFill struct {
	ExecID     string          `json:"execId"`
	Price      decimal.Decimal `json:"price"`
	Size       decimal.Decimal `json:"size"`
	Side       Side            `json:"side"`
	Ts         time.Time       `json:"ts"`
	ExecLinkID string          `json:"execLinkId"`
	IsRPI      string          `json:"isRPI"` // NO, YES
}

// GetRPISymbolsService -- GET /api/v3/market/rpi-symbols
//
// Returns the list of symbols that support RPI (Retail Price Improvement) orders,
// optionally filtered to a single product category.
type GetRPISymbolsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRPISymbolsService() *GetRPISymbolsService {
	return &GetRPISymbolsService{c: c, params: map[string]string{}}
}

func (s *GetRPISymbolsService) SetCategory(category Category) *GetRPISymbolsService {
	s.params["category"] = string(category)
	return s
}

func (s *GetRPISymbolsService) Do(ctx context.Context) ([]RPISymbol, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/rpi-symbols", s.params)
	resp, err := request.Do[[]RPISymbol](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// RPISymbol is a single RPI-enabled trading pair. Bitget returns the category in
// lowercase here (e.g. "spot", "usdt-futures").
type RPISymbol struct {
	Category string `json:"category"`
	Symbol   string `json:"symbol"`
}
