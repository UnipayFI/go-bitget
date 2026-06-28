package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceRealityOrderService -- POST /api/v3/trade/place-reality-order (UTA trade read & write)
//
// Submits a limit or market order for Reality stock trading pairs (e.g.
// rAAPLUSDT). Restricted to whitelisted UIDs. symbol, side, orderType and qty
// are required; price is required for limit orders. For market buys qty is in
// the quote coin, otherwise in the base coin. The reply carries the new order's
// identifiers.
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
// clientOid (orderId wins if both are set).
type CancelRealityOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelRealityOrderService(symbol string) *CancelRealityOrderService {
	return &CancelRealityOrderService{c: c, body: map[string]any{
		"symbol": symbol,
	}}
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
