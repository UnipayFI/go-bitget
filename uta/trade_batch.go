package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// OrderResult is one entry in a batch order/cancel/modify response. The per-item
// code/msg are populated only when that individual order failed (batch endpoints
// allow partial success).
type OrderResult struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
	Code          string `json:"code"`
	Msg           string `json:"msg"`
}

// BatchOrderItem is a single order in a batch place request. category, symbol,
// qty, side and orderType are required; price is required for limit orders and
// timeInForce is required for limit orders.
type BatchOrderItem struct {
	Category      Category        `json:"category"`
	Symbol        string          `json:"symbol"`
	Qty           decimal.Decimal `json:"qty"`
	Price         decimal.Decimal `json:"price,omitempty"`
	Side          Side            `json:"side"`
	OrderType     OrderType       `json:"orderType"`
	TimeInForce   TimeInForce     `json:"timeInForce,omitempty"`
	PosSide       PosSide         `json:"posSide,omitempty"`
	ClientOrderID string          `json:"clientOid,omitempty"`
	StpMode       string          `json:"stpMode,omitempty"`
}

// PlaceBatchService -- POST /api/v3/trade/place-batch
//
// Places up to 20 orders in a single request; all orders must share the same
// category. The request body is a raw JSON array of order objects and the
// response is an array of per-order results (partial success allowed).
type PlaceBatchService struct {
	c     *UTAClient
	items []BatchOrderItem
}

func (c *UTAClient) NewPlaceBatchService(items []BatchOrderItem) *PlaceBatchService {
	return &PlaceBatchService{c: c, items: items}
}

func (s *PlaceBatchService) SetItems(items []BatchOrderItem) *PlaceBatchService {
	s.items = items
	return s
}

func (s *PlaceBatchService) Do(ctx context.Context) ([]OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/place-batch").SetBody(s.items).WithSign()
	resp, err := request.Do[[]OrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchModifyItem is a single order modification in a batch modify request.
// Either OrderID or ClientOid must be set (orderId takes priority if both are
// supplied).
type BatchModifyItem struct {
	OrderID       string          `json:"orderId,omitempty"`
	ClientOrderID string          `json:"clientOid,omitempty"`
	Qty           decimal.Decimal `json:"qty,omitempty"`
	Price         decimal.Decimal `json:"price,omitempty"`
	AutoCancel    string          `json:"autoCancel,omitempty"` // yes, no (default)
	Symbol        string          `json:"symbol,omitempty"`
	Category      Category        `json:"category,omitempty"`
}

// BatchModifyOrderService -- POST /api/v3/trade/batch-modify-order
//
// Modifies up to 20 orders in a single request. The request body is a raw JSON
// array of modification objects and the response is an array of per-order
// results.
type BatchModifyOrderService struct {
	c     *UTAClient
	items []BatchModifyItem
}

func (c *UTAClient) NewBatchModifyOrderService(items []BatchModifyItem) *BatchModifyOrderService {
	return &BatchModifyOrderService{c: c, items: items}
}

func (s *BatchModifyOrderService) SetItems(items []BatchModifyItem) *BatchModifyOrderService {
	s.items = items
	return s
}

func (s *BatchModifyOrderService) Do(ctx context.Context) ([]OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/batch-modify-order").SetBody(s.items).WithSign()
	resp, err := request.Do[[]OrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BatchCancelItem is a single order in a batch cancel request. category and
// symbol are required; either OrderID or ClientOid must be set.
type BatchCancelItem struct {
	OrderID       string   `json:"orderId,omitempty"`
	ClientOrderID string   `json:"clientOid,omitempty"`
	Category      Category `json:"category"`
	Symbol        string   `json:"symbol"`
}

// CancelBatchService -- POST /api/v3/trade/cancel-batch
//
// Cancels up to 20 orders in a single request. The request body is a raw JSON
// array of order objects and the response is an array of per-order results
// (partial success allowed).
type CancelBatchService struct {
	c     *UTAClient
	items []BatchCancelItem
}

func (c *UTAClient) NewCancelBatchService(items []BatchCancelItem) *CancelBatchService {
	return &CancelBatchService{c: c, items: items}
}

func (s *CancelBatchService) SetItems(items []BatchCancelItem) *CancelBatchService {
	s.items = items
	return s
}

func (s *CancelBatchService) Do(ctx context.Context) ([]OrderResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/cancel-batch").SetBody(s.items).WithSign()
	resp, err := request.Do[[]OrderResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ClosePositionsService -- POST /api/v3/trade/close-positions
//
// Market-closes positions for a futures category. Without symbol it closes all
// positions in the category; without posSide it closes both sides.
type ClosePositionsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewClosePositionsService(category Category) *ClosePositionsService {
	return &ClosePositionsService{c: c, body: map[string]any{"category": string(category)}}
}

func (s *ClosePositionsService) SetSymbol(symbol string) *ClosePositionsService {
	s.body["symbol"] = symbol
	return s
}

func (s *ClosePositionsService) SetPosSide(posSide PosSide) *ClosePositionsService {
	s.body["posSide"] = string(posSide)
	return s
}

func (s *ClosePositionsService) Do(ctx context.Context) (*ClosePositionsResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/close-positions", s.body).WithSign()
	return request.Do[ClosePositionsResult](req)
}

// ClosePositionsResult wraps the list of close orders that were submitted.
type ClosePositionsResult struct {
	List []OrderResult `json:"list"`
}

// MovePositionItem is a single position to transfer in a move-positions request.
type MovePositionItem struct {
	Symbol string          `json:"symbol"`
	Side   Side            `json:"side"`
	Qty    decimal.Decimal `json:"qty"`
}

// MovePositionsService -- POST /api/v3/account/move-positions
//
// Transfers up to 10 futures positions from one account to another. fromUid,
// toUid and category are required; positions are supplied via SetPositions.
// Business limit: 100 times/day/UID, minimum 30 seconds between requests.
type MovePositionsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewMovePositionsService(fromUid, toUid string, category Category) *MovePositionsService {
	return &MovePositionsService{c: c, body: map[string]any{
		"fromUid":  fromUid,
		"toUid":    toUid,
		"category": string(category),
	}}
}

func (s *MovePositionsService) SetPositions(positions []MovePositionItem) *MovePositionsService {
	s.body["positionList"] = positions
	return s
}

func (s *MovePositionsService) Do(ctx context.Context) (*MovePositionsResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/move-positions", s.body).WithSign()
	return request.Do[MovePositionsResult](req)
}

// MovePositionsResult holds the per-account order results: closePosition is the
// source account's result list and openPosition the target account's.
type MovePositionsResult struct {
	ClosePosition []OrderResult `json:"closePosition"`
	OpenPosition  []OrderResult `json:"openPosition"`
}
