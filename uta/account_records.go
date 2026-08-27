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
	Category        Category            `json:"category"`
	ID              string              `json:"id"`
	Symbol          string              `json:"symbol"`
	Coin            string              `json:"coin"`
	Type            FinancialRecordType `json:"type"`
	PositionType    string              `json:"positionType"` // crossed, isolated
	PositionAmount  decimal.Decimal     `json:"positionAmount"`
	PositionBalance decimal.Decimal     `json:"positionBalance"`
	Amount          decimal.Decimal     `json:"amount"`
	Fee             decimal.Decimal     `json:"fee"`
	Balance         decimal.Decimal     `json:"balance"`
	Ts              time.Time           `json:"ts"`
}

// FinancialRecordType classifies a financial-records entry. The constants below
// cover the contract settlement-fee (funding) entries SDK callers filter on; the
// full vocabulary is large and overlaps the tax endpoint
// (classic/tax.FutureTaxType uses the same upper-case identifiers).
type FinancialRecordType string

const (
	FinancialRecordContractMainSettleFeeUserIn  FinancialRecordType = "CONTRACT_MAIN_SETTLE_FEE_USER_IN"
	FinancialRecordContractMainSettleFeeUserOut FinancialRecordType = "CONTRACT_MAIN_SETTLE_FEE_USER_OUT"
	// RWA cash-dividend cross-margin funding-fee entries.
	FinancialRecordRWAContractMainSettleFeeUserIn  FinancialRecordType = "RWA_CONTRACT_MAIN_SETTLE_FEE_USER_IN"
	FinancialRecordRWAContractMainSettleFeeUserOut FinancialRecordType = "RWA_CONTRACT_MAIN_SETTLE_FEE_USER_OUT"
	// RWA contract rebase entries (position open/close and in-SSM buy/sell).
	FinancialRecordRWAContractRebaseUserOpenLong   FinancialRecordType = "RWA_CONTRACT_REBASE_USER_OPEN_LONG"
	FinancialRecordRWAContractRebaseUserOpenShort  FinancialRecordType = "RWA_CONTRACT_REBASE_USER_OPEN_SHORT"
	FinancialRecordRWAContractRebaseUserCloseLong  FinancialRecordType = "RWA_CONTRACT_REBASE_USER_CLOSE_LONG"
	FinancialRecordRWAContractRebaseUserCloseShort FinancialRecordType = "RWA_CONTRACT_REBASE_USER_CLOSE_SHORT"
	FinancialRecordRWAContractRebaseUserBuyInSSM   FinancialRecordType = "RWA_CONTRACT_REBASE_USER_BUY_IN_SSM"
	FinancialRecordRWAContractRebaseUserSellInSSM  FinancialRecordType = "RWA_CONTRACT_REBASE_USER_SELL_IN_SSM"
	// Copy-trading (trace) transfer entries between the UTA and funding account.
	FinancialRecordTraceTransferUserOut  FinancialRecordType = "TRACE_TRANSFER_USER_OUT"
	FinancialRecordTraceTransferUserIn   FinancialRecordType = "TRACE_TRANSFER_USER_IN"
	FinancialRecordTraceTransferRefundIn FinancialRecordType = "TRACE_TRANSFER_REFUND_IN"
)

// GetFundingFinancialRecordsService -- GET /api/v3/account/funding-financial-records (UTA mgt. read)
//
// Returns the funding account's financial (ledger) records, paginated by cursor
// and bounded to a 90-day access window.
type GetFundingFinancialRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetFundingFinancialRecordsService() *GetFundingFinancialRecordsService {
	return &GetFundingFinancialRecordsService{c: c, params: map[string]string{}}
}

func (s *GetFundingFinancialRecordsService) SetCoin(coin string) *GetFundingFinancialRecordsService {
	s.params["coin"] = coin
	return s
}

// SetType filters by the record type (buy, sell, deposit, withdraw,
// transfer_in, interest, dividend, ...).
func (s *GetFundingFinancialRecordsService) SetType(recordType string) *GetFundingFinancialRecordsService {
	s.params["type"] = recordType
	return s
}

// SetStartTime filters records at or after t (90-day access window).
func (s *GetFundingFinancialRecordsService) SetStartTime(t time.Time) *GetFundingFinancialRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetFundingFinancialRecordsService) SetEndTime(t time.Time) *GetFundingFinancialRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFundingFinancialRecordsService) SetLimit(limit int) *GetFundingFinancialRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFundingFinancialRecordsService) SetCursor(cursor string) *GetFundingFinancialRecordsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetFundingFinancialRecordsService) Do(ctx context.Context) (*FundingFinancialRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/funding-financial-records", s.params).WithSign()
	return request.Do[FundingFinancialRecords](req)
}

type FundingFinancialRecords struct {
	List   []FundingFinancialRecord `json:"list"`
	Cursor string                   `json:"cursor"`
}

type FundingFinancialRecord struct {
	ID        string          `json:"id"`
	Coin      string          `json:"coin"`
	GroupType string          `json:"groupType"` // transaction, deposit, withdraw, transfer, financial, loan, convert, rwa, stock, ...
	Type      string          `json:"type"`      // buy, sell, deposit, withdraw, transfer_in, interest, dividend, ...
	Amount    decimal.Decimal `json:"amount"`
	Balance   decimal.Decimal `json:"balance"`
	Ts        time.Time       `json:"ts"`
}

// GetConvertRecordsService -- GET /api/v3/account/convert-records (UTA mgt. read)
//
// Returns the unified account's coin-conversion records, paginated by cursor and
// bounded to a 90-day access window. fromCoin and toCoin are optional filters.
type GetConvertRecordsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetConvertRecordsService() *GetConvertRecordsService {
	return &GetConvertRecordsService{c: c, params: map[string]string{}}
}

// SetFromCoin filters records by the source coin.
func (s *GetConvertRecordsService) SetFromCoin(fromCoin string) *GetConvertRecordsService {
	s.params["fromCoin"] = fromCoin
	return s
}

// SetToCoin filters records by the target coin.
func (s *GetConvertRecordsService) SetToCoin(toCoin string) *GetConvertRecordsService {
	s.params["toCoin"] = toCoin
	return s
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
