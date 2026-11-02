package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// StockPlusOrderType is the order type of a Stock+ stock or option order.
type StockPlusOrderType string

const (
	StockPlusOrderTypeLimit                StockPlusOrderType = "LO"      // Limit Order
	StockPlusOrderTypeMarket               StockPlusOrderType = "MO"      // Market Order
	StockPlusOrderTypeLimitIfTouched       StockPlusOrderType = "LIT"     // Limit If Touched
	StockPlusOrderTypeMarketIfTouched      StockPlusOrderType = "MIT"     // Market If Touched
	StockPlusOrderTypeTrailingStopLimitAmt StockPlusOrderType = "TSLPAMT" // Trailing Stop Limit Order (trailing amount)
	StockPlusOrderTypeTrailingStopLimitPct StockPlusOrderType = "TSLPPCT" // Trailing Stop Limit Order (trailing percentage)
)

// StockPlusSide is the order side of a Stock+ order. Unlike Side, the wire
// values are capitalised.
type StockPlusSide string

const (
	StockPlusSideBuy  StockPlusSide = "Buy"
	StockPlusSideSell StockPlusSide = "Sell"
)

// GetStockPlusMaxPurchaseService -- GET /api/v3/stockplus/trade/max-purchase
//
// Returns the estimated maximum purchasable quantity of a Stock+ stock or
// option. Only whitelisted users who have completed securities account opening
// can call it; contact your BD to apply. Rate limit: 30 requests/30s per UID,
// minimum interval 0.02s.
type GetStockPlusMaxPurchaseService struct {
	c      *UTAClient
	params map[string]string
}

// NewGetStockPlusMaxPurchaseService takes the symbol in `ticker.region` format
// (e.g. "AAPL.US"). StockPlusSideSell is only supported for US stock short
// selling.
func (c *UTAClient) NewGetStockPlusMaxPurchaseService(symbol string, orderType StockPlusOrderType, side StockPlusSide) *GetStockPlusMaxPurchaseService {
	return &GetStockPlusMaxPurchaseService{c: c, params: map[string]string{
		"symbol":    symbol,
		"orderType": string(orderType),
		"side":      string(side),
	}}
}

// SetPrice sets the estimated order price.
func (s *GetStockPlusMaxPurchaseService) SetPrice(price decimal.Decimal) *GetStockPlusMaxPurchaseService {
	s.params["price"] = price.String()
	return s
}

// SetCurrency sets the settlement currency.
func (s *GetStockPlusMaxPurchaseService) SetCurrency(currency string) *GetStockPlusMaxPurchaseService {
	s.params["currency"] = currency
	return s
}

// SetOrderID sets the order being modified; required when estimating the max
// purchasable quantity for an order modification.
func (s *GetStockPlusMaxPurchaseService) SetOrderID(orderID string) *GetStockPlusMaxPurchaseService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetStockPlusMaxPurchaseService) Do(ctx context.Context) (*StockPlusMaxPurchase, error) {
	req := request.Get(ctx, s.c, "/api/v3/stockplus/trade/max-purchase", s.params).WithSign()
	return request.Do[StockPlusMaxPurchase](req)
}

// StockPlusMaxPurchase is the estimated maximum purchasable quantity. Either
// field is zero when the API sends its "" default.
type StockPlusMaxPurchase struct {
	CashMaxQty   decimal.Decimal `json:"cashMaxQty"`   // max quantity payable in cash
	MarginMaxQty decimal.Decimal `json:"marginMaxQty"` // max quantity payable on margin
}
