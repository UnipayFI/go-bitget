package ws

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
)

// This file wraps the classic v2 margin WebSocket channels (cross + isolated).
// All margin channels route with instType "MARGIN" (InstTypeMargin). Cross types
// and services are prefixed MarginCross, isolated ones MarginIsolated, to avoid
// collisions with the sibling spot/mix channel files in this package.

// -----------------------------------------------------------------------------
// index-price (PUBLIC, cross) -- channel "index-price", instId "default" (all symbols).
// https://www.bitget.com/api-doc/margin/cross/websocket/public/Margin-Index-Price
// -----------------------------------------------------------------------------

// SubscribeMarginCrossIndexPriceService -- public "index-price" channel (cross margin).
type SubscribeMarginCrossIndexPriceService struct {
	c      *WebSocketClient
	instId string
}

// NewSubscribeMarginCrossIndexPriceService subscribes to the cross-margin index
// price channel. instId is the symbol; pass "default" for all symbols.
func (c *WebSocketClient) NewSubscribeMarginCrossIndexPriceService(instId string) *SubscribeMarginCrossIndexPriceService {
	return &SubscribeMarginCrossIndexPriceService{c: c, instId: instId}
}

func (s *SubscribeMarginCrossIndexPriceService) Do(ctx context.Context, cb WsHandler[MarginCrossIndexPrice]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MarginCrossIndexPrice](ctx, s.c, false,
		WsArg{InstType: string(InstTypeMargin), Channel: "index-price", InstId: s.instId}, cb)
}

// MarginCrossIndexPrice is one element of the "index-price" push data array.
type MarginCrossIndexPrice struct {
	Symbol     string          `json:"symbol"`
	BaseCoin   string          `json:"baseCoin"`
	QuoteCoin  string          `json:"quoteCoin"`
	IndexPrice decimal.Decimal `json:"indexPrice"`
	Ts         time.Time       `json:"ts"`
}

// -----------------------------------------------------------------------------
// orders-crossed (PRIVATE) -- channel "orders-crossed", instId required (symbol).
// https://www.bitget.com/api-doc/margin/cross/websocket/private/Cross-Orders
// -----------------------------------------------------------------------------

// SubscribeMarginCrossOrdersService -- private "orders-crossed" channel (cross margin).
type SubscribeMarginCrossOrdersService struct {
	c      *WebSocketClient
	instId string
}

// NewSubscribeMarginCrossOrdersService subscribes to the cross-margin order
// channel. instId is the symbol (e.g. "BTCUSDT").
func (c *WebSocketClient) NewSubscribeMarginCrossOrdersService(instId string) *SubscribeMarginCrossOrdersService {
	return &SubscribeMarginCrossOrdersService{c: c, instId: instId}
}

func (s *SubscribeMarginCrossOrdersService) Do(ctx context.Context, cb WsHandler[MarginCrossOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MarginCrossOrder](ctx, s.c, true,
		WsArg{InstType: string(InstTypeMargin), Channel: "orders-crossed", InstId: s.instId}, cb)
}

// MarginCrossOrder is one element of the "orders-crossed" push data array.
type MarginCrossOrder struct {
	BaseSize         decimal.Decimal  `json:"baseSize"`         // quantity of base coins
	CTime            time.Time        `json:"cTime"`            // creation time
	UTime            time.Time        `json:"uTime"`            // update time
	ClientOid        string           `json:"clientOid"`        // client order id
	FillPrice        decimal.Decimal  `json:"fillPrice"`        // execution price
	BaseVolume       decimal.Decimal  `json:"baseVolume"`       // filled quantity
	FillTotalAmount  decimal.Decimal  `json:"fillTotalAmount"`  // total value of filled amount
	LoanType         string           `json:"loanType"`         // normal, autoLoan, autoRepay, autoLoanAndRepay
	OrderId          string           `json:"orderId"`          // order id
	OrderType        string           `json:"orderType"`        // limit or market
	Price            decimal.Decimal  `json:"price"`            // order price
	QuoteSize        decimal.Decimal  `json:"quoteSize"`        // quantity of denominated coins
	Side             string           `json:"side"`             // buy/sell
	EnterPointSource string           `json:"enterPointSource"` // WEB, API, SYS, ANDROID, IOS
	Status           string           `json:"status"`           // live, partially_filled, filled, cancelled
	Force            string           `json:"force"`            // order strategy
	StpMode          string           `json:"stpMode"`          // none, cancel_taker, cancel_maker, cancel_both
	FeeDetail        []MarginOrderFee `json:"feeDetail"`        // transaction fees
}

// MarginOrderFee is one element of a margin order's "feeDetail" array (shared by
// the cross and isolated order channels).
type MarginOrderFee struct {
	FeeCoin           string          `json:"feeCoin"`           // fee currency
	Deduction         string          `json:"deduction"`         // discount indicator
	TotalDeductionFee decimal.Decimal `json:"totalDeductionFee"` // discount fee amount
	TotalFee          decimal.Decimal `json:"totalFee"`          // platform transaction fee
}

// -----------------------------------------------------------------------------
// account-crossed (PRIVATE) -- channel "account-crossed", coin "default" (all coins).
// https://www.bitget.com/api-doc/margin/cross/websocket/private/Margin-Cross-Account-Assets
// -----------------------------------------------------------------------------

// SubscribeMarginCrossAccountService -- private "account-crossed" channel (cross margin assets).
type SubscribeMarginCrossAccountService struct {
	c    *WebSocketClient
	coin string
}

// NewSubscribeMarginCrossAccountService subscribes to the cross-margin account
// assets channel. coin selects the coin; pass "default" for all coins (the only
// value currently supported).
func (c *WebSocketClient) NewSubscribeMarginCrossAccountService(coin string) *SubscribeMarginCrossAccountService {
	return &SubscribeMarginCrossAccountService{c: c, coin: coin}
}

func (s *SubscribeMarginCrossAccountService) Do(ctx context.Context, cb WsHandler[MarginCrossAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MarginCrossAccount](ctx, s.c, true,
		WsArg{InstType: string(InstTypeMargin), Channel: "account-crossed", Coin: s.coin}, cb)
}

// MarginCrossAccount is one element of the "account-crossed" push data array.
type MarginCrossAccount struct {
	Available decimal.Decimal `json:"available"` // available amount
	Borrow    decimal.Decimal `json:"borrow"`    // borrow amount
	Coin      string          `json:"coin"`      // coin name
	Frozen    decimal.Decimal `json:"frozen"`    // amount frozen
	Coupon    decimal.Decimal `json:"coupon"`    // coupon
	Id        string          `json:"id"`        // id
	Interest  decimal.Decimal `json:"interest"`  // interest
	UTime     time.Time       `json:"uTime"`     // updated time
}

// -----------------------------------------------------------------------------
// orders-isolated (PRIVATE) -- channel "orders-isolated", instId required (symbol).
// https://www.bitget.com/api-doc/margin/isolated/websocket/private/Isolate-Orders
// -----------------------------------------------------------------------------

// SubscribeMarginIsolatedOrdersService -- private "orders-isolated" channel (isolated margin).
type SubscribeMarginIsolatedOrdersService struct {
	c      *WebSocketClient
	instId string
}

// NewSubscribeMarginIsolatedOrdersService subscribes to the isolated-margin order
// channel. instId is the symbol (e.g. "BTCUSDT").
func (c *WebSocketClient) NewSubscribeMarginIsolatedOrdersService(instId string) *SubscribeMarginIsolatedOrdersService {
	return &SubscribeMarginIsolatedOrdersService{c: c, instId: instId}
}

func (s *SubscribeMarginIsolatedOrdersService) Do(ctx context.Context, cb WsHandler[MarginIsolatedOrder]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MarginIsolatedOrder](ctx, s.c, true,
		WsArg{InstType: string(InstTypeMargin), Channel: "orders-isolated", InstId: s.instId}, cb)
}

// MarginIsolatedOrder is one element of the "orders-isolated" push data array.
type MarginIsolatedOrder struct {
	BaseSize         decimal.Decimal  `json:"baseSize"`         // number of base coins
	CTime            time.Time        `json:"cTime"`            // creation time
	UTime            time.Time        `json:"uTime"`            // update time
	ClientOid        string           `json:"clientOid"`        // client order id
	FillPrice        decimal.Decimal  `json:"fillPrice"`        // sale price
	BaseVolume       decimal.Decimal  `json:"baseVolume"`       // base coin quantity
	FillTotalAmount  decimal.Decimal  `json:"fillTotalAmount"`  // sum of money sold
	LoanType         string           `json:"loanType"`         // normal, autoLoan, autoRepay, autoLoanAndRepay
	OrderId          string           `json:"orderId"`          // order id
	OrderType        string           `json:"orderType"`        // limit or market
	Price            decimal.Decimal  `json:"price"`            // order price
	QuoteSize        decimal.Decimal  `json:"quoteSize"`        // number of denominated coins
	Side             string           `json:"side"`             // buy/sell
	EnterPointSource string           `json:"enterPointSource"` // WEB, API, SYS, ANDROID, IOS
	Status           string           `json:"status"`           // order status
	Force            string           `json:"force"`            // order strategy
	StpMode          string           `json:"stpMode"`          // none, cancel_taker, cancel_maker, cancel_both
	FeeDetail        []MarginOrderFee `json:"feeDetail"`        // transaction fees
}

// -----------------------------------------------------------------------------
// account-isolated (PRIVATE) -- channel "account-isolated", coin "default" (all coins).
// https://www.bitget.com/api-doc/margin/isolated/websocket/private/Margin-isolated-account-assets
// -----------------------------------------------------------------------------

// SubscribeMarginIsolatedAccountService -- private "account-isolated" channel (isolated margin assets).
type SubscribeMarginIsolatedAccountService struct {
	c    *WebSocketClient
	coin string
}

// NewSubscribeMarginIsolatedAccountService subscribes to the isolated-margin
// account assets channel. coin selects the coin; pass "default" for all coins
// (the only value currently supported).
func (c *WebSocketClient) NewSubscribeMarginIsolatedAccountService(coin string) *SubscribeMarginIsolatedAccountService {
	return &SubscribeMarginIsolatedAccountService{c: c, coin: coin}
}

func (s *SubscribeMarginIsolatedAccountService) Do(ctx context.Context, cb WsHandler[MarginIsolatedAccount]) (chan<- struct{}, <-chan struct{}, error) {
	return Subscribe[[]MarginIsolatedAccount](ctx, s.c, true,
		WsArg{InstType: string(InstTypeMargin), Channel: "account-isolated", Coin: s.coin}, cb)
}

// MarginIsolatedAccount is one element of the "account-isolated" push data array.
type MarginIsolatedAccount struct {
	UTime     time.Time       `json:"uTime"`     // updated time
	Id        string          `json:"id"`        // id
	Coin      string          `json:"coin"`      // coin name
	Symbol    string          `json:"symbol"`    // trading pair
	Available decimal.Decimal `json:"available"` // available amount
	Borrow    decimal.Decimal `json:"borrow"`    // borrow amount
	Frozen    decimal.Decimal `json:"frozen"`    // amount frozen
	Interest  decimal.Decimal `json:"interest"`  // interest accrued
	Coupon    decimal.Decimal `json:"coupon"`    // coupon value
}
