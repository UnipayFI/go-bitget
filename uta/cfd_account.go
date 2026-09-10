package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCFDInstrumentsService -- GET /api/v3/cfd/account/instruments (UTA trade read)
//
// Returns the tradable CFD instruments, optionally filtered to a single symbol.
type GetCFDInstrumentsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDInstrumentsService() *GetCFDInstrumentsService {
	return &GetCFDInstrumentsService{c: c, params: map[string]string{}}
}

// SetSymbol filters to a single CFD trading pair (e.g. XAUUSD).
func (s *GetCFDInstrumentsService) SetSymbol(symbol string) *GetCFDInstrumentsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetCFDInstrumentsService) Do(ctx context.Context) ([]CFDInstrument, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/account/instruments", s.params).WithSign()
	resp, err := request.Do[[]CFDInstrument](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CFDInstrument is the trading specification of one CFD pair. CFD contracts
// carry their own margin/commission model rather than the spot/futures one, so
// this shares no fields with Instrument.
type CFDInstrument struct {
	Symbol     string    `json:"symbol"`
	OnlineTime time.Time `json:"onlineTime"`
	Currency   string    `json:"currency"` // base currency
	// Digits is the number of price decimal places.
	Digits       decimal.Decimal `json:"digits"`
	Spread       decimal.Decimal `json:"spread"`
	ContractSize decimal.Decimal `json:"contractSize"`
	// MarginCurrency and ProfitCurrency are the currencies margin is posted in
	// and profit is settled in.
	MarginCurrency string          `json:"marginCurrency"`
	ProfitCurrency string          `json:"profitCurrency"`
	StopLevel      decimal.Decimal `json:"stopLevel"`
	// Enable is the trading state: 0 disabled, 1 close-only, 2 enabled.
	Enable        string          `json:"enable"`
	MarginPercent decimal.Decimal `json:"marginPercent"`
	InitialMargin decimal.Decimal `json:"initialMargin"`
	MarginMode    string          `json:"marginMode"` // Forex, CFD
	MaintMargin   decimal.Decimal `json:"maintMargin"`
	// MaintMarginRateBuy/Sell are the reciprocal of leverage (0.002 at 500x).
	MaintMarginRateBuy    decimal.Decimal `json:"maintMarginRateBuy"`
	MaintMarginRateSell   decimal.Decimal `json:"maintMarginRateSell"`
	InitialMarginRateBuy  decimal.Decimal `json:"initialMarginRateBuy"`
	InitialMarginRateSell decimal.Decimal `json:"initialMarginRateSell"`
	MaxVolume             decimal.Decimal `json:"maxVolume"`
	MinVolume             decimal.Decimal `json:"minVolume"`
	StepVolume            decimal.Decimal `json:"stepVolume"`
	PriceCurrency         string          `json:"priceCurrency"` // quote currency
	TickSize              decimal.Decimal `json:"tickSize"`
	TickValue             decimal.Decimal `json:"tickValue"`
	Leverage              decimal.Decimal `json:"leverage"`
	TradeTime             string          `json:"tradeTime"` // tradable hours, e.g. 24/5
	// CommissionType is how commission is charged: 0 per lot, 1 as a percentage
	// of trade value. CommissionRate is read according to it.
	CommissionType string          `json:"commissionType"`
	CommissionRate decimal.Decimal `json:"commissionRate"`
	// SymbolCategory is the symbol name without its contract suffix, so
	// XAUUSD.s has SymbolCategory "XAUUSD".
	SymbolCategory string `json:"symbolCategory"`
	// ExchangeRate converts PriceCurrency to USD, MarginUSDRate converts
	// MarginCurrency to USD.
	ExchangeRate  decimal.Decimal `json:"exchangeRate"`
	MarginUSDRate decimal.Decimal `json:"marginUsdRate"`
}

// GetCFDFundDetailService -- GET /api/v3/cfd/account/fund-detail (UTA account read)
//
// Returns the CFD account's balance, margin usage and running profit.
type GetCFDFundDetailService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCFDFundDetailService() *GetCFDFundDetailService {
	return &GetCFDFundDetailService{c: c}
}

func (s *GetCFDFundDetailService) Do(ctx context.Context) (*CFDFundDetail, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/account/fund-detail").WithSign()
	return request.Do[CFDFundDetail](req)
}

type CFDFundDetail struct {
	Leverage decimal.Decimal `json:"leverage"`
	// Balance excludes the credit line; Equity includes it.
	Balance decimal.Decimal `json:"balance"`
	Credit  decimal.Decimal `json:"credit"`
	// MarginStopOut, MarginCall and MarginLevel are percentages.
	MarginStopOut decimal.Decimal `json:"marginStopOut"`
	MarginUsed    decimal.Decimal `json:"marginUsed"`
	MarginCall    decimal.Decimal `json:"marginCall"`
	MarginFree    decimal.Decimal `json:"marginFree"`
	MarginLevel   decimal.Decimal `json:"marginLevel"`
	Equity        decimal.Decimal `json:"equity"`
	PnL           decimal.Decimal `json:"pnl"`
	Status        string          `json:"status"`
	// Swap is the accrued overnight interest, also called storage or rollover
	// fee.
	Swap     decimal.Decimal `json:"swap"`
	Frozen   decimal.Decimal `json:"frozen"`
	Currency string          `json:"currency"`
}

// CFDTransferDirection is the direction of a CFD account transfer.
type CFDTransferDirection string

const (
	CFDTransferIn  CFDTransferDirection = "in"  // deposit into the CFD account
	CFDTransferOut CFDTransferDirection = "out" // withdraw from the CFD account
)

// CFDTransferService -- POST /api/v3/cfd/account/transfer (UTA account read & write)
//
// Moves a coin into or out of the CFD account. accountType is the counterpart
// account: "funding" (spot/funding) or "uta" (unified trading account).
type CFDTransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCFDTransferService(coin string, amount decimal.Decimal, accountType string, direction CFDTransferDirection) *CFDTransferService {
	return &CFDTransferService{c: c, body: map[string]any{
		"coin":        coin,
		"amount":      amount.String(),
		"accountType": accountType,
		"direction":   string(direction),
	}}
}

func (s *CFDTransferService) Do(ctx context.Context) (*CFDTransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/cfd/account/transfer", s.body).WithSign()
	return request.Do[CFDTransferResult](req)
}

type CFDTransferResult struct {
	TransferID string `json:"transferId"`
}

// GetCFDTransferRecordsService -- GET /api/v3/cfd/account/transfer-records (UTA account read)
//
// Returns the CFD account's transfer records, paginated by cursor. A single
// query may span at most 30 days within a 90-day lookback window.
type GetCFDTransferRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDTransferRecordsService() *GetCFDTransferRecordsService {
	return &GetCFDTransferRecordsService{c: c, params: map[string]string{}}
}

// SetTransferID filters to a single transfer.
func (s *GetCFDTransferRecordsService) SetTransferID(transferID string) *GetCFDTransferRecordsService {
	s.params["transferId"] = transferID
	return s
}

// SetSubUID queries a sub-account's transfer records instead of the caller's.
func (s *GetCFDTransferRecordsService) SetSubUID(subUid string) *GetCFDTransferRecordsService {
	s.params["subUid"] = subUid
	return s
}

// SetStartTime filters records at or after t (90-day lookback window).
func (s *GetCFDTransferRecordsService) SetStartTime(t time.Time) *GetCFDTransferRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetCFDTransferRecordsService) SetEndTime(t time.Time) *GetCFDTransferRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetDirection filters to deposits ("in") or withdrawals ("out"); both are
// returned when unset.
func (s *GetCFDTransferRecordsService) SetDirection(direction CFDTransferDirection) *GetCFDTransferRecordsService {
	s.params["direction"] = string(direction)
	return s
}

// SetLimit caps the number of records per page (default 10, range [1, 50]).
func (s *GetCFDTransferRecordsService) SetLimit(limit int) *GetCFDTransferRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetCursor pages forward using the cursor of the previous page.
func (s *GetCFDTransferRecordsService) SetCursor(cursor string) *GetCFDTransferRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCFDTransferRecordsService) Do(ctx context.Context) (*CFDTransferRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/account/transfer-records", s.params).WithSign()
	return request.Do[CFDTransferRecords](req)
}

type CFDTransferRecords struct {
	List   []CFDTransferRecord `json:"list"`
	Cursor string              `json:"cursor"`
}

type CFDTransferRecord struct {
	TransferID  string               `json:"transferId"`
	Coin        string               `json:"coin"`
	Amount      decimal.Decimal      `json:"amount"`
	Direction   CFDTransferDirection `json:"direction"`
	Status      string               `json:"status"`
	AccountType string               `json:"accountType"` // funding, uta
	CreatedTime time.Time            `json:"createdTime"`
	UpdatedTime time.Time            `json:"updatedTime"`
}

// GetCFDFinancialRecordsService -- GET /api/v3/cfd/account/financial-records (UTA account read)
//
// Returns the CFD account's ledger records, paginated by cursor. A single query
// may span at most 30 days within a 90-day lookback window.
type GetCFDFinancialRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCFDFinancialRecordsService() *GetCFDFinancialRecordsService {
	return &GetCFDFinancialRecordsService{c: c, params: map[string]string{}}
}

// SetType filters to a single financial-record type.
func (s *GetCFDFinancialRecordsService) SetType(recordType string) *GetCFDFinancialRecordsService {
	s.params["type"] = recordType
	return s
}

// SetStartTime filters records at or after t (90-day lookback window).
func (s *GetCFDFinancialRecordsService) SetStartTime(t time.Time) *GetCFDFinancialRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetCFDFinancialRecordsService) SetEndTime(t time.Time) *GetCFDFinancialRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records per page (default 10, range [1, 50]).
func (s *GetCFDFinancialRecordsService) SetLimit(limit int) *GetCFDFinancialRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetCursor pages forward using the cursor of the previous page.
func (s *GetCFDFinancialRecordsService) SetCursor(cursor string) *GetCFDFinancialRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCFDFinancialRecordsService) Do(ctx context.Context) ([]CFDFinancialRecord, error) {
	req := request.Get(ctx, s.c, "/api/v3/cfd/account/financial-records", s.params).WithSign()
	resp, err := request.Do[[]CFDFinancialRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CFDFinancialRecord is one CFD ledger entry. Unlike the rest of the REST API,
// this endpoint returns bare JSON numbers rather than quoted strings; the
// decimal fields accept both.
type CFDFinancialRecord struct {
	ID     int64  `json:"id"`
	Symbol string `json:"symbol"`
	// Side is the trade side, documented and delivered as a number rather than
	// the buy/sell strings the other endpoints use.
	Side int             `json:"side"`
	Qty  decimal.Decimal `json:"qty"`
	// Swap is the overnight interest and Fee the commission; both are excluded
	// from CashFlow and included in BalanceChange.
	Swap decimal.Decimal `json:"swap"`
	Fee  decimal.Decimal `json:"fee"`
	// CashFlow is settled PnL, transfers, compensation and the like.
	CashFlow      decimal.Decimal `json:"cashFlow"`
	BalanceBefore decimal.Decimal `json:"balanceBefore"`
	// BalanceChange is CashFlow + Swap + Fee.
	BalanceChange decimal.Decimal `json:"balanceChange"`
	BalanceAfter  decimal.Decimal `json:"balanceAfter"`
	CreditBefore  decimal.Decimal `json:"creditBefore"`
	CreditChange  decimal.Decimal `json:"creditChange"`
	CreditAfter   decimal.Decimal `json:"creditAfter"`
	OpenPrice     decimal.Decimal `json:"openPrice"`
	ClosePrice    decimal.Decimal `json:"closePrice"`
	OrderID       int64           `json:"orderId"`
	Ts            time.Time       `json:"ts"`
}
