package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetFinancialRecordsService -- GET /api/v3/account/financial-records (UTA mgt. read)
//
// Returns the unified account's financial (ledger) records for a product
// category, paginated by cursor and bounded to a 90-day lookback window.
type GetFinancialRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFinancialRecordsService(category Category) *GetFinancialRecordsService {
	return &GetFinancialRecordsService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetFinancialRecordsService) SetCoin(coin string) *GetFinancialRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetFinancialRecordsService) SetType(recordType string) *GetFinancialRecordsService {
	s.params["type"] = recordType
	return s
}

// SetStartTime filters records at or after t (90-day lookback window).
func (s *GetFinancialRecordsService) SetStartTime(t time.Time) *GetFinancialRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetFinancialRecordsService) SetEndTime(t time.Time) *GetFinancialRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFinancialRecordsService) SetLimit(limit int) *GetFinancialRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFinancialRecordsService) SetCursor(cursor string) *GetFinancialRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetFinancialRecordsService) Do(ctx context.Context) (*FinancialRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/financial-records", s.params).WithSign()
	return request.Do[FinancialRecords](req)
}

type FinancialRecords struct {
	List   []FinancialRecord `json:"list"`
	Cursor string            `json:"cursor"`
}

type FinancialRecord struct {
	Category Category            `json:"category"`
	ID       string              `json:"id"`
	Symbol   string              `json:"symbol"`
	Coin     string              `json:"coin"`
	Type     FinancialRecordType `json:"type"`
	Amount   decimal.Decimal     `json:"amount"`
	Fee      decimal.Decimal     `json:"fee"`
	Balance  decimal.Decimal     `json:"balance"`
	Ts       time.Time           `json:"ts"`
}

// FinancialRecordType classifies a financial-records entry. The constants below
// cover the contract settlement-fee (funding) entries SDK callers filter on; the
// full vocabulary is large and overlaps the tax endpoint
// (classic/tax.FutureTaxType uses the same upper-case identifiers).
type FinancialRecordType string

const (
	FinancialRecordContractMainSettleFeeUserIn  FinancialRecordType = "CONTRACT_MAIN_SETTLE_FEE_USER_IN"
	FinancialRecordContractMainSettleFeeUserOut FinancialRecordType = "CONTRACT_MAIN_SETTLE_FEE_USER_OUT"
)

// GetConvertRecordsService -- GET /api/v3/account/convert-records (UTA mgt. read)
//
// Returns the unified account's coin-conversion records for a from/to coin
// pair, paginated by cursor and bounded to a 90-day access window.
type GetConvertRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetConvertRecordsService(fromCoin, toCoin string) *GetConvertRecordsService {
	return &GetConvertRecordsService{c: c, params: map[string]string{"fromCoin": fromCoin, "toCoin": toCoin}}
}

// SetStartTime filters records at or after t (90-day access window).
func (s *GetConvertRecordsService) SetStartTime(t time.Time) *GetConvertRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetConvertRecordsService) SetEndTime(t time.Time) *GetConvertRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetConvertRecordsService) SetLimit(limit int) *GetConvertRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetConvertRecordsService) SetCursor(cursor string) *GetConvertRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetConvertRecordsService) Do(ctx context.Context) (*ConvertRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/convert-records", s.params).WithSign()
	return request.Do[ConvertRecords](req)
}

type ConvertRecords struct {
	List   []ConvertRecord `json:"list"`
	Cursor string          `json:"cursor"`
}

type ConvertRecord struct {
	FromCoin     string          `json:"fromCoin"`
	FromCoinSize decimal.Decimal `json:"fromCoinSize"`
	ToCoin       string          `json:"toCoin"`
	ToCoinSize   decimal.Decimal `json:"toCoinSize"`
	Price        decimal.Decimal `json:"price"`
	Ts           time.Time       `json:"ts"`
}

// GetRepayableCoinsService -- GET /api/v3/account/repayable-coins (UTA mgt. read)
//
// Returns the coins the unified account currently owes and can repay, with the
// repayable size and its USD-equivalent amount.
type GetRepayableCoinsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetRepayableCoinsService() *GetRepayableCoinsService {
	return &GetRepayableCoinsService{c: c}
}

func (s *GetRepayableCoinsService) Do(ctx context.Context) (*RepayableCoins, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/repayable-coins").WithSign()
	return request.Do[RepayableCoins](req)
}

type RepayableCoins struct {
	RepayableCoinList []RepayableCoin `json:"repayableCoinList"`
	MaxSelection      string          `json:"maxSelection"`
}

type RepayableCoin struct {
	Coin   string          `json:"coin"`
	Size   decimal.Decimal `json:"size"`
	Amount decimal.Decimal `json:"amount"`
}

// GetPaymentCoinsService -- GET /api/v3/account/payment-coins (UTA mgt. read)
//
// Returns the coins the unified account can use to fund a repayment, with the
// available size and its USD-equivalent amount.
type GetPaymentCoinsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetPaymentCoinsService() *GetPaymentCoinsService {
	return &GetPaymentCoinsService{c: c}
}

func (s *GetPaymentCoinsService) Do(ctx context.Context) (*PaymentCoins, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/payment-coins").WithSign()
	return request.Do[PaymentCoins](req)
}

type PaymentCoins struct {
	PaymentCoinList []PaymentCoin `json:"paymentCoinList"`
	MaxSelection    string        `json:"maxSelection"`
}

type PaymentCoin struct {
	Coin   string          `json:"coin"`
	Size   decimal.Decimal `json:"size"`
	Amount decimal.Decimal `json:"amount"`
}

// RepayService -- POST /api/v3/account/repay (UTA mgt. read & write)
//
// Repays the given debt coins using the given payment coins. Both lists are
// coin names; the exchange settles the conversion.
type RepayService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewRepayService(repayableCoinList, paymentCoinList []string) *RepayService {
	return &RepayService{c: c, body: map[string]any{
		"repayableCoinList": repayableCoinList,
		"paymentCoinList":   paymentCoinList,
	}}
}

func (s *RepayService) Do(ctx context.Context) (*RepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/repay", s.body).WithSign()
	return request.Do[RepayResult](req)
}

type RepayResult struct {
	Result      string          `json:"result"`
	RepayAmount decimal.Decimal `json:"repayAmount"`
}
