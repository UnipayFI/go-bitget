package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetMarginCurrenciesService -- GET /api/v2/margin/currencies (margin read, signed)
//
// Returns the margin-tradable trading pairs together with their leverage,
// risk-ratio, fee and borrowing configuration.
type GetMarginCurrenciesService struct {
	c *MarginClient
}

func (c *MarginClient) NewGetMarginCurrenciesService() *GetMarginCurrenciesService {
	return &GetMarginCurrenciesService{c: c}
}

func (s *GetMarginCurrenciesService) Do(ctx context.Context) ([]MarginCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/currencies").WithSign()
	resp, err := request.Do[[]MarginCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginCurrency is the margin configuration for a single trading pair.
type MarginCurrency struct {
	Symbol               string          `json:"symbol"`
	BaseCoin             string          `json:"baseCoin"`
	QuoteCoin            string          `json:"quoteCoin"`
	MaxCrossedLeverage   string          `json:"maxCrossedLeverage"`
	MaxIsolatedLeverage  string          `json:"maxIsolatedLeverage"`
	WarningRiskRatio     decimal.Decimal `json:"warningRiskRatio"`
	LiquidationRiskRatio decimal.Decimal `json:"liquidationRiskRatio"`
	MinTradeAmount       decimal.Decimal `json:"minTradeAmount"`
	MaxTradeAmount       decimal.Decimal `json:"maxTradeAmount"`
	TakerFeeRate         decimal.Decimal `json:"takerFeeRate"`
	MakerFeeRate         decimal.Decimal `json:"makerFeeRate"`
	PricePrecision       string          `json:"pricePrecision"`
	QuantityPrecision    string          `json:"quantityPrecision"`
	MinTradeUSDT         decimal.Decimal `json:"minTradeUSDT"`
	IsBorrowable         bool            `json:"isBorrowable"`
	UserMinBorrow        decimal.Decimal `json:"userMinBorrow"`
	Status               string          `json:"status"` // 1: tradable, 2: under temporary maintenance
	// The example response returns these as bare booleans (the doc field table
	// mislabels them String).
	IsIsolatedBaseBorrowable  bool `json:"isIsolatedBaseBorrowable"`
	IsIsolatedQuoteBorrowable bool `json:"isIsolatedQuoteBorrowable"`
	IsCrossBorrowable         bool `json:"isCrossBorrowable"`
}

// GetInterestRateRecordService -- GET /api/v2/margin/interest-rate-record (margin read, signed)
//
// Returns the current leverage (borrow) interest rate for a single coin.
type GetInterestRateRecordService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetInterestRateRecordService(coin string) *GetInterestRateRecordService {
	return &GetInterestRateRecordService{c: c, params: map[string]string{"coin": coin}}
}

// SetStartTime narrows the query to records on or after t. The endpoint
// primarily returns the current rate for coin; the optional time window is
// accepted for history-style queries.
func (s *GetInterestRateRecordService) SetStartTime(t time.Time) *GetInterestRateRecordService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime narrows the query to records on or before t.
func (s *GetInterestRateRecordService) SetEndTime(t time.Time) *GetInterestRateRecordService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInterestRateRecordService) Do(ctx context.Context) (*InterestRateRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/interest-rate-record", s.params).WithSign()
	return request.Do[InterestRateRecord](req)
}

// InterestRateRecord is the current leverage interest rate for a coin.
type InterestRateRecord struct {
	Coin               string          `json:"coin"`
	DailyInterestRate  decimal.Decimal `json:"dailyInterestRate"`
	AnnualInterestRate decimal.Decimal `json:"annualInterestRate"`
	UpdatedTime        time.Time       `json:"updatedTime"`
}
