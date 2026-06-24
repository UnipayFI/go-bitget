package uta

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// wsCategory converts the REST uppercase Category (e.g. "SPOT") to the
// lowercase form the WebSocket trade ops require (e.g. "spot").
func wsCategory(category Category) string {
	return strings.ToLower(string(category))
}

// UTATradeConn is a persistent, logged-in WebSocket connection for placing,
// modifying and cancelling orders over the stream (a low-latency alternative to
// the REST trade endpoints). Obtain one with UTAWebSocketClient.DialTrade and
// Close it when done.
type UTATradeConn struct {
	*request.WsTradeConn
}

// DialTrade opens and logs in a WebSocket trade connection. Requires
// WithWebSocketAuth credentials.
func (c *UTAWebSocketClient) DialTrade(ctx context.Context) (*UTATradeConn, error) {
	conn, err := request.DialWsTrade(ctx, c)
	if err != nil {
		return nil, err
	}
	return &UTATradeConn{conn}, nil
}

// WsOrderAck acknowledges a WebSocket order op. Code/Msg are only set per-item
// on batch failures.
type WsOrderAck struct {
	Symbol        string    `json:"symbol"`
	OrderID       string    `json:"orderId"`
	ClientOrderID string    `json:"clientOid"`
	CTime         time.Time `json:"cTime"`
	Code          string    `json:"code"`
	Msg           string    `json:"msg"`
}

// WsNewOrder describes an order to place over the stream. Price is a pointer so
// it is omitted for market orders.
type WsNewOrder struct {
	Symbol        string           `json:"symbol"`
	Side          Side             `json:"side"`
	OrderType     OrderType        `json:"orderType"`
	Qty           decimal.Decimal  `json:"qty"`
	Price         *decimal.Decimal `json:"price,omitempty"`
	TimeInForce   TimeInForce      `json:"timeInForce,omitempty"`
	ClientOrderID string           `json:"clientOid,omitempty"`
	PosSide       PosSide          `json:"posSide,omitempty"`
	ReduceOnly    string           `json:"reduceOnly,omitempty"`
	MarginMode    MarginMode       `json:"marginMode,omitempty"`
	StpMode       string           `json:"stpMode,omitempty"`
}

// WsModifyOrder amends an open order. Identify it by OrderID or ClientOid.
type WsModifyOrder struct {
	Symbol        string           `json:"symbol,omitempty"`
	OrderID       string           `json:"orderId,omitempty"`
	ClientOrderID string           `json:"clientOid,omitempty"`
	Qty           *decimal.Decimal `json:"qty,omitempty"`
	Price         *decimal.Decimal `json:"price,omitempty"`
}

// WsCancelOrder identifies an order to cancel by OrderID or ClientOid.
type WsCancelOrder struct {
	Symbol        string `json:"symbol,omitempty"`
	OrderID       string `json:"orderId,omitempty"`
	ClientOrderID string `json:"clientOid,omitempty"`
}

// PlaceOrder places a single order over the stream.
func (t *UTATradeConn) PlaceOrder(ctx context.Context, category Category, order WsNewOrder) (*WsOrderAck, error) {
	return t.single(ctx, category, "place-order", order)
}

// BatchPlaceOrders places up to 20 same-category orders in one frame.
func (t *UTATradeConn) BatchPlaceOrders(ctx context.Context, category Category, orders []WsNewOrder) ([]WsOrderAck, error) {
	args := make([]any, len(orders))
	for i, o := range orders {
		args[i] = o
	}
	return t.batch(ctx, category, "batch-place-order", args)
}

// ModifyOrder amends a single open order.
func (t *UTATradeConn) ModifyOrder(ctx context.Context, category Category, modify WsModifyOrder) (*WsOrderAck, error) {
	return t.single(ctx, category, "modify-order", modify)
}

// BatchModifyOrders amends up to 20 same-category orders in one frame.
func (t *UTATradeConn) BatchModifyOrders(ctx context.Context, category Category, modifies []WsModifyOrder) ([]WsOrderAck, error) {
	args := make([]any, len(modifies))
	for i, m := range modifies {
		args[i] = m
	}
	return t.batch(ctx, category, "batch-modify-order", args)
}

// CancelOrder cancels a single order.
func (t *UTATradeConn) CancelOrder(ctx context.Context, category Category, cancel WsCancelOrder) (*WsOrderAck, error) {
	return t.single(ctx, category, "cancel-order", cancel)
}

// BatchCancelOrders cancels up to 20 same-category orders in one frame.
func (t *UTATradeConn) BatchCancelOrders(ctx context.Context, category Category, cancels []WsCancelOrder) ([]WsOrderAck, error) {
	args := make([]any, len(cancels))
	for i, c := range cancels {
		args[i] = c
	}
	return t.batch(ctx, category, "batch-cancel-order", args)
}

func (t *UTATradeConn) single(ctx context.Context, category Category, topic string, arg any) (*WsOrderAck, error) {
	resp, err := request.WsTradeCall[[]WsOrderAck](ctx, t.WsTradeConn, "trade", wsCategory(category), topic, []any{arg})
	if err != nil {
		return nil, err
	}
	if len(resp.Args) == 0 {
		return nil, fmt.Errorf("ws %s: empty acknowledgement", topic)
	}
	return &resp.Args[0], nil
}

func (t *UTATradeConn) batch(ctx context.Context, category Category, topic string, args []any) ([]WsOrderAck, error) {
	resp, err := request.WsTradeCall[[]WsOrderAck](ctx, t.WsTradeConn, "trade", wsCategory(category), topic, args)
	if err != nil {
		return nil, err
	}
	return resp.Args, nil
}
