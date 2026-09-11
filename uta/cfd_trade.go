package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// PlaceCFDOrderService -- POST /api/v3/cfd/trade/place-order (UTA trade read & write)
//
// Submits a single CFD order. price is required for limit orders. Prices must
// be an integer multiple of the pair's tick size. symbol must match the trading
// mode the account is currently in -- spelling it for another mode is rejected.
// The reply carries only the request id; poll the unfilled/history order
// endpoints for the order id.
type PlaceCFDOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewPlaceCFDOrderService(symbol string, orderType OrderType, side Side, qty decimal.Decimal) *PlaceCFDOrderService {
	return &PlaceCFDOrderService{c: c, body: map[string]any{
		"symbol":    symbol,
		"orderType": string(orderType),
		"side":      string(side),
		"qty":       qty.String(),
	}}
}

// SetPrice sets the order price, required for limit orders.
func (s *PlaceCFDOrderService) SetPrice(price decimal.Decimal) *PlaceCFDOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTakeProfit sets the preset take-profit price.
func (s *PlaceCFDOrderService) SetTakeProfit(takeProfit decimal.Decimal) *PlaceCFDOrderService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the preset stop-loss price.
func (s *PlaceCFDOrderService) SetStopLoss(stopLoss decimal.Decimal) *PlaceCFDOrderService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

func (s *PlaceCFDOrderService) Do(ctx context.Context) (*CFDOrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/place-order", s.body).WithSign()
	return request.Do[CFDOrderRef](req)
}

// CFDOrderRef identifies an accepted CFD order request. Placement returns only
// TrxID (the request id); modification returns only OrderID.
type CFDOrderRef struct {
	TrxID   string `json:"trxId"`
	OrderID string `json:"orderId"`
}

// ModifyCFDOrderService -- POST /api/v3/cfd/trade/modify-order (UTA trade read & write)
//
// Amends an open CFD order's price and/or preset take-profit and stop-loss.
type ModifyCFDOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyCFDOrderService(orderID string) *ModifyCFDOrderService {
	return &ModifyCFDOrderService{c: c, body: map[string]any{"orderId": orderID}}
}

// SetPrice sets the new order price, required for limit orders.
func (s *ModifyCFDOrderService) SetPrice(price decimal.Decimal) *ModifyCFDOrderService {
	s.body["price"] = price.String()
	return s
}

// SetTakeProfit sets the new take-profit price.
func (s *ModifyCFDOrderService) SetTakeProfit(takeProfit decimal.Decimal) *ModifyCFDOrderService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the new stop-loss price.
func (s *ModifyCFDOrderService) SetStopLoss(stopLoss decimal.Decimal) *ModifyCFDOrderService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

func (s *ModifyCFDOrderService) Do(ctx context.Context) (*CFDOrderRef, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/modify-order", s.body).WithSign()
	return request.Do[CFDOrderRef](req)
}

// CancelCFDOrderService -- POST /api/v3/cfd/trade/cancel-order (UTA trade read & write)
//
// Cancels a single open CFD order. The reply data is null.
type CancelCFDOrderService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelCFDOrderService(orderID string) *CancelCFDOrderService {
	return &CancelCFDOrderService{c: c, body: map[string]any{"orderId": orderID}}
}

func (s *CancelCFDOrderService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/cancel-order", s.body).WithSign()
	return request.Do[any](req)
}

// CancelAllCFDOrdersService -- POST /api/v3/cfd/trade/cancel-all (UTA trade read & write)
//
// Cancels every open CFD order, or only those of one symbol when SetSymbol is
// used. The reply data is null.
type CancelAllCFDOrdersService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelAllCFDOrdersService() *CancelAllCFDOrdersService {
	return &CancelAllCFDOrdersService{c: c, body: map[string]any{}}
}

// SetSymbol restricts the cancellation to a single CFD trading pair.
func (s *CancelAllCFDOrdersService) SetSymbol(symbol string) *CancelAllCFDOrdersService {
	s.body["symbol"] = symbol
	return s
}

func (s *CancelAllCFDOrdersService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/cancel-all", s.body).WithSign()
	return request.Do[any](req)
}

// CloseCFDPositionsService -- POST /api/v3/cfd/trade/close-positions (UTA trade read & write)
//
// Closes qty of the position identified by positionId (from
// GetCFDPositionsService). Only market-price closing is currently supported.
// The reply data is null.
type CloseCFDPositionsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCloseCFDPositionsService(positionID string, qty decimal.Decimal) *CloseCFDPositionsService {
	return &CloseCFDPositionsService{c: c, body: map[string]any{
		"positionId": positionID,
		"qty":        qty.String(),
	}}
}

func (s *CloseCFDPositionsService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/close-positions", s.body).WithSign()
	return request.Do[any](req)
}

// CloseAllCFDPositionsService -- POST /api/v3/cfd/trade/close-all-positions (UTA trade read & write)
//
// Closes every CFD position, or only those of one symbol when SetSymbol is
// used. The reply data is null.
type CloseAllCFDPositionsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCloseAllCFDPositionsService() *CloseAllCFDPositionsService {
	return &CloseAllCFDPositionsService{c: c, body: map[string]any{}}
}

// SetSymbol restricts the close to a single CFD trading pair.
func (s *CloseAllCFDPositionsService) SetSymbol(symbol string) *CloseAllCFDPositionsService {
	s.body["symbol"] = symbol
	return s
}

func (s *CloseAllCFDPositionsService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/trade/close-all-positions", s.body).WithSign()
	return request.Do[any](req)
}

// GetCFDUnfilledOrdersService -- GET /api/v3/cfd/trade/unfilled-order (UTA trade read)
//
// Returns the open (unfilled) CFD orders for a symbol.
type GetCFDUnfilledOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDUnfilledOrdersService(symbol string) *GetCFDUnfilledOrdersService {
	return &GetCFDUnfilledOrdersService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetSubUID queries a sub-account's unfilled orders instead of the caller's.
func (s *GetCFDUnfilledOrdersService) SetSubUID(subUid string) *GetCFDUnfilledOrdersService {
	s.params["subUid"] = subUid
	return s
}

func (s *GetCFDUnfilledOrdersService) Do(ctx context.Context) ([]CFDOrder, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/trade/unfilled-order", s.params).WithSign()
	resp, err := request.Do[[]CFDOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CFDOrder struct {
	Symbol     string          `json:"symbol"`
	OrderType  OrderType       `json:"orderType"`
	Side       Side            `json:"side"`
	Price      decimal.Decimal `json:"price"`
	Qty        decimal.Decimal `json:"qty"`
	TakeProfit decimal.Decimal `json:"takeProfit"`
	StopLoss   decimal.Decimal `json:"stopLoss"`
	TrxID      string          `json:"trxId"`
	OrderID    string          `json:"orderId"`
}

// GetCFDOrderHistoryService -- GET /api/v3/cfd/trade/history-order (UTA trade read)
//
// Returns historical CFD orders for a symbol, paginated by page number. A
// single query may span at most 30 days within a 90-day lookback window.
type GetCFDOrderHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDOrderHistoryService(symbol string) *GetCFDOrderHistoryService {
	return &GetCFDOrderHistoryService{c: c, params: map[string]string{"symbol": symbol}}
}

// SetSubUID queries a sub-account's historical orders instead of the caller's.
func (s *GetCFDOrderHistoryService) SetSubUID(subUid string) *GetCFDOrderHistoryService {
	s.params["subUid"] = subUid
	return s
}

// SetStartTime filters orders at or after t (90-day lookback window).
func (s *GetCFDOrderHistoryService) SetStartTime(t time.Time) *GetCFDOrderHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t (max 30-day range from startTime).
func (s *GetCFDOrderHistoryService) SetEndTime(t time.Time) *GetCFDOrderHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of orders per page (default 10, range [1, 50]).
func (s *GetCFDOrderHistoryService) SetLimit(limit int) *GetCFDOrderHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetPageNo selects the page to return.
func (s *GetCFDOrderHistoryService) SetPageNo(pageNo int) *GetCFDOrderHistoryService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetCFDOrderHistoryService) Do(ctx context.Context) (*CFDOrderHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/trade/history-order", s.params).WithSign()
	return request.Do[CFDOrderHistory](req)
}

type CFDOrderHistory struct {
	List []CFDHistoryOrder `json:"list"`
}

// CFDHistoryOrder is a settled CFD order: a CFDOrder plus its fill progress and
// final status.
type CFDHistoryOrder struct {
	Symbol     string          `json:"symbol"`
	OrderType  OrderType       `json:"orderType"`
	Side       Side            `json:"side"`
	Price      decimal.Decimal `json:"price"`
	Qty        decimal.Decimal `json:"qty"`
	CumExecQty decimal.Decimal `json:"cumExecQty"`
	TakeProfit decimal.Decimal `json:"takeProfit"`
	StopLoss   decimal.Decimal `json:"stopLoss"`
	Status     string          `json:"status"`
	TrxID      string          `json:"trxId"`
	OrderID    string          `json:"orderId"`
}

// GetCFDPositionsService -- GET /api/v3/cfd/trade/current-positions (UTA trade read)
//
// Returns the open CFD positions, for every symbol unless SetSymbol is used.
type GetCFDPositionsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDPositionsService() *GetCFDPositionsService {
	return &GetCFDPositionsService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single CFD trading pair.
func (s *GetCFDPositionsService) SetSymbol(symbol string) *GetCFDPositionsService {
	s.params["symbol"] = symbol
	return s
}

// SetSubUID queries a sub-account's positions instead of the caller's.
func (s *GetCFDPositionsService) SetSubUID(subUid string) *GetCFDPositionsService {
	s.params["subUid"] = subUid
	return s
}

func (s *GetCFDPositionsService) Do(ctx context.Context) ([]CFDPosition, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/trade/current-positions", s.params).WithSign()
	resp, err := request.Do[[]CFDPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CFDPosition struct {
	PositionID string          `json:"positionId"`
	Symbol     string          `json:"symbol"`
	OrderType  OrderType       `json:"orderType"`
	Side       Side            `json:"side"`
	Qty        decimal.Decimal `json:"qty"`
	OpenPrice  decimal.Decimal `json:"openPrice"`
	TakeProfit decimal.Decimal `json:"takeProfit"`
	StopLoss   decimal.Decimal `json:"stopLoss"`
	// Interest is the accrued overnight interest (swap) on the position.
	Interest      decimal.Decimal `json:"interest"`
	UnrealizedPnL decimal.Decimal `json:"unrealizedPnl"`
	// TotalProfit is Interest + UnrealizedPnL.
	TotalProfit decimal.Decimal `json:"totalProfit"`
}
