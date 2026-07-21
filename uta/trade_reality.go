package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceRealityOrderService -- POST /api/v3/trade/place-reality-order (UTA trade read & write)
//
// Submits a limit or market order for Reality stock trading pairs (e.g.
// rAAPLUSDT). Restricted to whitelisted UIDs. symbol, side, orderType and qty
// are required; price is required for limit orders. For market buys qty is in
// the quote coin, otherwise in the base coin. category selects spot vs. margin
// trading (defaults to SPOT). The reply carries the new order's identifiers.
type PlaceRealityOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewPlaceRealityOrderService(symbol string, side Side, orderType OrderType, qty decimal.Decimal) *PlaceRealityOrderService {
	return &PlaceRealityOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"side":      string(side),
		"orderType": string(orderType),
		"qty":       qty.String(),
	}}
}

// SetCategory selects the trading mode: CategorySpot or CategoryMargin
// (defaults to CategorySpot).
func (s *PlaceRealityOrderService) SetCategory(category Category) *PlaceRealityOrderService {
	s.body["category"] = string(category)
	return s
}

// SetPrice sets the order price (required for limit orders).
func (s *PlaceRealityOrderService) SetPrice(price decimal.Decimal) *PlaceRealityOrderService {
	s.body["price"] = price.String()
	return s
}

// SetClientOrderID sets the client-generated order identifier.
func (s *PlaceRealityOrderService) SetClientOrderID(clientOid string) *PlaceRealityOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *PlaceRealityOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/place-reality-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// CancelRealityOrderService -- POST /api/v3/trade/cancel-reality-order (UTA trade read & write)
//
// Cancels an unfilled or partially filled Reality stock order. Restricted to
// whitelisted UIDs. symbol is required; identify the order by orderId or
// clientOid (orderId wins if both are set). category selects spot vs. margin
// trading (defaults to SPOT).
type CancelRealityOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelRealityOrderService(symbol string) *CancelRealityOrderService {
	return &CancelRealityOrderService{c: c, body: map[string]any{
		"symbol": symbol,
	}}
}

// SetCategory selects the trading mode: CategorySpot or CategoryMargin
// (defaults to CategorySpot).
func (s *CancelRealityOrderService) SetCategory(category Category) *CancelRealityOrderService {
	s.body["category"] = string(category)
	return s
}

// SetOrderID sets the order identifier (orderId or clientOid is required).
func (s *CancelRealityOrderService) SetOrderID(orderId string) *CancelRealityOrderService {
	s.body["orderId"] = orderId
	return s
}

// SetClientOrderID sets the client order identifier (orderId or clientOid is required).
func (s *CancelRealityOrderService) SetClientOrderID(clientOid string) *CancelRealityOrderService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *CancelRealityOrderService) Do(ctx context.Context) (*OrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/cancel-reality-order", s.body).WithSign()
	return request.Do[OrderRef](req)
}

// GetRealityOrderBookService -- GET /api/v3/account/reality-orderbook (UTA trade read)
//
// Returns the order book depth snapshot for a Reality stock trading pair (e.g.
// rAAPLUSDT). Requires API Key authentication and Reality access. symbol is
// required.
type GetRealityOrderBookService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRealityOrderBookService(symbol string) *GetRealityOrderBookService {
	return &GetRealityOrderBookService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetRealityOrderBookService) Do(ctx context.Context) (*RealityOrderBook, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/reality-orderbook", s.params).WithSign()
	return request.Do[RealityOrderBook](req)
}

// RealityOrderBook is a Reality trading pair's order book depth snapshot. Asks
// ("a") and bids ("b") arrive as arrays of [price, size] string pairs.
type RealityOrderBook struct {
	Symbol string              `json:"symbol"`
	Asks   [][]decimal.Decimal `json:"a"`
	Bids   [][]decimal.Decimal `json:"b"`
	Ts     time.Time           `json:"ts"`
}

// GetRealityFillsService -- GET /api/v3/account/reality-fills (UTA trade read)
//
// Returns the most recent platform fills for a Reality stock trading pair (e.g.
// rAAPLUSDT), covering the last 3 months. Requires API Key authentication and
// Reality access. symbol is required; limit defaults to 100 (max 100).
type GetRealityFillsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRealityFillsService(symbol string) *GetRealityFillsService {
	return &GetRealityFillsService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetRealityFillsService) SetLimit(limit int) *GetRealityFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRealityFillsService) Do(ctx context.Context) ([]RealityFill, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/reality-fills", s.params).WithSign()
	resp, err := request.Do[[]RealityFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// RealityFill is a single recent platform fill for a Reality trading pair.
type RealityFill struct {
	ExecID string          `json:"execId"`
	Price  decimal.Decimal `json:"price"`
	Size   decimal.Decimal `json:"size"`
	Side   Side            `json:"side"`
	Ts     time.Time       `json:"ts"`
}
