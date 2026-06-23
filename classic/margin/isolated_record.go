package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// IsolatedRepayType is the kind of repayment recorded on an isolated-margin
// repay-history entry.
type IsolatedRepayType string

const (
	IsolatedRepayTypeAuto   IsolatedRepayType = "auto-repay"
	IsolatedRepayTypeManual IsolatedRepayType = "manual-repay"
	IsolatedRepayTypeLiq    IsolatedRepayType = "liq-repay"
	IsolatedRepayTypeForce  IsolatedRepayType = "force-repay"
)

// IsolatedBorrowType is whether an isolated-margin loan was taken out
// automatically (alongside an auto-loan order) or manually.
type IsolatedBorrowType string

const (
	IsolatedBorrowTypeAuto   IsolatedBorrowType = "auto_loan"
	IsolatedBorrowTypeManual IsolatedBorrowType = "manual_loan"
)

// IsolatedInterestType is whether an isolated-margin interest charge is the
// initial charge taken at borrow time ("first") or a later periodic accrual
// ("scheduled").
type IsolatedInterestType string

const (
	IsolatedInterestTypeFirst     IsolatedInterestType = "first"
	IsolatedInterestTypeScheduled IsolatedInterestType = "scheduled"
)

// IsolatedMarginType is the capital-flow category reported on an
// isolated-margin financial record.
type IsolatedMarginType string

const (
	IsolatedMarginTypeTransferIn     IsolatedMarginType = "transfer_in"
	IsolatedMarginTypeTransferOut    IsolatedMarginType = "transfer_out"
	IsolatedMarginTypeBorrow         IsolatedMarginType = "borrow"
	IsolatedMarginTypeRepay          IsolatedMarginType = "repay"
	IsolatedMarginTypeLiquidationFee IsolatedMarginType = "liquidation_fee"
	IsolatedMarginTypeCompensate     IsolatedMarginType = "compensate"
	IsolatedMarginTypeDealIn         IsolatedMarginType = "deal_in"
	IsolatedMarginTypeDealOut        IsolatedMarginType = "deal_out"
	IsolatedMarginTypeConfiscated    IsolatedMarginType = "confiscated"
	IsolatedMarginTypeExchangeIn     IsolatedMarginType = "exchange_in"
	IsolatedMarginTypeExchangeOut    IsolatedMarginType = "exchange_out"
	IsolatedMarginTypeSysExchangeIn  IsolatedMarginType = "sys_exchange_in"
	IsolatedMarginTypeSysExchangeOut IsolatedMarginType = "sys_exchange_out"
)

// GetIsolatedRepayHistoryService -- GET /api/v2/margin/isolated/repay-history (private)
//
// Returns the isolated-margin repayment history for a trading pair over a time
// window, paginated by idLessThan.
type GetIsolatedRepayHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedRepayHistoryService(symbol string, startTime time.Time) *GetIsolatedRepayHistoryService {
	return &GetIsolatedRepayHistoryService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedRepayHistoryService) SetRepayID(repayID string) *GetIsolatedRepayHistoryService {
	s.params["repayId"] = repayID
	return s
}

func (s *GetIsolatedRepayHistoryService) SetCoin(coin string) *GetIsolatedRepayHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetIsolatedRepayHistoryService) SetEndTime(t time.Time) *GetIsolatedRepayHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedRepayHistoryService) SetLimit(limit int) *GetIsolatedRepayHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest repayId from the previous page.
func (s *GetIsolatedRepayHistoryService) SetIDLessThan(idLessThan string) *GetIsolatedRepayHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedRepayHistoryService) Do(ctx context.Context) (*IsolatedRepayHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/repay-history", s.params).WithSign()
	return request.Do[IsolatedRepayHistory](req)
}

// IsolatedRepayHistory is one page of isolated-margin repayment records.
type IsolatedRepayHistory struct {
	ResultList []IsolatedRepayRecord `json:"resultList"`
	MaxID      string                `json:"maxId"`
	MinID      string                `json:"minId"`
}

// IsolatedRepayRecord is a single isolated-margin repayment.
type IsolatedRepayRecord struct {
	RepayID        string            `json:"repayId"`
	Symbol         string            `json:"symbol"`
	Coin           string            `json:"coin"`
	RepayPrincipal decimal.Decimal   `json:"repayPrincipal"`
	RepayAmount    decimal.Decimal   `json:"repayAmount"`
	RepayInterest  decimal.Decimal   `json:"repayInterest"`
	RepayType      IsolatedRepayType `json:"repayType"`
	CTime          time.Time         `json:"cTime"`
	UTime          time.Time         `json:"uTime"`
}

// GetIsolatedBorrowHistoryService -- GET /api/v2/margin/isolated/borrow-history (private)
//
// Returns the isolated-margin loan (borrow) history for a trading pair over a
// time window, paginated by idLessThan.
type GetIsolatedBorrowHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedBorrowHistoryService(symbol string, startTime time.Time) *GetIsolatedBorrowHistoryService {
	return &GetIsolatedBorrowHistoryService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedBorrowHistoryService) SetLoanID(loanID string) *GetIsolatedBorrowHistoryService {
	s.params["loanId"] = loanID
	return s
}

func (s *GetIsolatedBorrowHistoryService) SetCoin(coin string) *GetIsolatedBorrowHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetIsolatedBorrowHistoryService) SetEndTime(t time.Time) *GetIsolatedBorrowHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedBorrowHistoryService) SetLimit(limit int) *GetIsolatedBorrowHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest loanId from the previous page.
func (s *GetIsolatedBorrowHistoryService) SetIDLessThan(idLessThan string) *GetIsolatedBorrowHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedBorrowHistoryService) Do(ctx context.Context) (*IsolatedBorrowHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/borrow-history", s.params).WithSign()
	return request.Do[IsolatedBorrowHistory](req)
}

// IsolatedBorrowHistory is one page of isolated-margin loan records.
type IsolatedBorrowHistory struct {
	ResultList []IsolatedBorrowRecord `json:"resultList"`
	MaxID      string                 `json:"maxId"`
	MinID      string                 `json:"minId"`
}

// IsolatedBorrowRecord is a single isolated-margin loan.
type IsolatedBorrowRecord struct {
	LoanID       string             `json:"loanId"`
	Coin         string             `json:"coin"`
	BorrowAmount decimal.Decimal    `json:"borrowAmount"`
	BorrowType   IsolatedBorrowType `json:"borrowType"`
	Symbol       string             `json:"symbol"`
	CTime        time.Time          `json:"cTime"`
	UTime        time.Time          `json:"uTime"`
}

// GetIsolatedInterestHistoryService -- GET /api/v2/margin/isolated/interest-history (private)
//
// Returns the isolated-margin interest accrual history for a trading pair over a
// time window, paginated by idLessThan.
type GetIsolatedInterestHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedInterestHistoryService(symbol string, startTime time.Time) *GetIsolatedInterestHistoryService {
	return &GetIsolatedInterestHistoryService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedInterestHistoryService) SetCoin(coin string) *GetIsolatedInterestHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetIsolatedInterestHistoryService) SetEndTime(t time.Time) *GetIsolatedInterestHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedInterestHistoryService) SetLimit(limit int) *GetIsolatedInterestHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest interestId from the previous page.
func (s *GetIsolatedInterestHistoryService) SetIDLessThan(idLessThan string) *GetIsolatedInterestHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedInterestHistoryService) Do(ctx context.Context) (*IsolatedInterestHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/interest-history", s.params).WithSign()
	return request.Do[IsolatedInterestHistory](req)
}

// IsolatedInterestHistory is one page of isolated-margin interest records.
type IsolatedInterestHistory struct {
	ResultList []IsolatedInterestRecord `json:"resultList"`
	MaxID      string                   `json:"maxId"`
	MinID      string                   `json:"minId"`
}

// IsolatedInterestRecord is a single isolated-margin interest charge. The
// "interstType" key is reproduced verbatim from the Bitget documentation.
type IsolatedInterestRecord struct {
	Symbol            string               `json:"symbol"`
	InterestID        string               `json:"interestId"`
	InterestAmount    decimal.Decimal      `json:"interestAmount"`
	DailyInterestRate decimal.Decimal      `json:"dailyInterestRate"`
	InterstType       IsolatedInterestType `json:"interstType"`
	InterestCoin      string               `json:"interestCoin"`
	LoanCoin          string               `json:"loanCoin"`
	CTime             time.Time            `json:"cTime"`
	UTime             time.Time            `json:"uTime"`
}

// GetIsolatedLiquidationHistoryService -- GET /api/v2/margin/isolated/liquidation-history (private)
//
// Returns the isolated-margin liquidation history for a trading pair over a time
// window, paginated by idLessThan.
type GetIsolatedLiquidationHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedLiquidationHistoryService(symbol string, startTime time.Time) *GetIsolatedLiquidationHistoryService {
	return &GetIsolatedLiquidationHistoryService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedLiquidationHistoryService) SetEndTime(t time.Time) *GetIsolatedLiquidationHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedLiquidationHistoryService) SetLimit(limit int) *GetIsolatedLiquidationHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest liqId from the previous page.
func (s *GetIsolatedLiquidationHistoryService) SetIDLessThan(idLessThan string) *GetIsolatedLiquidationHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedLiquidationHistoryService) Do(ctx context.Context) (*IsolatedLiquidationHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/liquidation-history", s.params).WithSign()
	return request.Do[IsolatedLiquidationHistory](req)
}

// IsolatedLiquidationHistory is one page of isolated-margin liquidation records.
type IsolatedLiquidationHistory struct {
	ResultList []IsolatedLiquidationRecord `json:"resultList"`
	MaxID      string                      `json:"maxId"`
	MinID      string                      `json:"minId"`
}

// IsolatedLiquidationRecord is a single isolated-margin liquidation event.
type IsolatedLiquidationRecord struct {
	LiqID        string          `json:"liqId"`
	Symbol       string          `json:"symbol"`
	LiqStartTime time.Time       `json:"liqStartTime"`
	LiqEndTime   time.Time       `json:"liqEndTime"`
	LiqRiskRatio decimal.Decimal `json:"liqRiskRatio"`
	TotalAssets  decimal.Decimal `json:"totalAssets"`
	TotalDebt    decimal.Decimal `json:"totalDebt"`
	LiqFee       decimal.Decimal `json:"liqFee"`
	CTime        time.Time       `json:"cTime"`
	UTime        time.Time       `json:"uTime"`
}

// GetIsolatedFinancialRecordsService -- GET /api/v2/margin/isolated/financial-records (private)
//
// Returns the isolated-margin capital-flow (financial) records for a trading
// pair over a time window, paginated by idLessThan.
type GetIsolatedFinancialRecordsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedFinancialRecordsService(symbol string, startTime time.Time) *GetIsolatedFinancialRecordsService {
	return &GetIsolatedFinancialRecordsService{c: c, params: map[string]string{
		"symbol":    symbol,
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetIsolatedFinancialRecordsService) SetMarginType(marginType IsolatedMarginType) *GetIsolatedFinancialRecordsService {
	s.params["marginType"] = string(marginType)
	return s
}

func (s *GetIsolatedFinancialRecordsService) SetCoin(coin string) *GetIsolatedFinancialRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetIsolatedFinancialRecordsService) SetEndTime(t time.Time) *GetIsolatedFinancialRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetIsolatedFinancialRecordsService) SetLimit(limit int) *GetIsolatedFinancialRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the minId from the previous page.
func (s *GetIsolatedFinancialRecordsService) SetIDLessThan(idLessThan string) *GetIsolatedFinancialRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetIsolatedFinancialRecordsService) Do(ctx context.Context) (*IsolatedFinancialRecords, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/financial-records", s.params).WithSign()
	return request.Do[IsolatedFinancialRecords](req)
}

// IsolatedFinancialRecords is one page of isolated-margin capital-flow records.
type IsolatedFinancialRecords struct {
	ResultList []IsolatedFinancialRecord `json:"resultList"`
	MaxID      string                    `json:"maxId"`
	MinID      string                    `json:"minId"`
}

// IsolatedFinancialRecord is a single isolated-margin capital-flow entry.
type IsolatedFinancialRecord struct {
	Symbol     string             `json:"symbol"`
	Coin       string             `json:"coin"`
	MarginID   string             `json:"marginId"`
	MarginType IsolatedMarginType `json:"marginType"`
	Amount     decimal.Decimal    `json:"amount"`
	Balance    decimal.Decimal    `json:"balance"`
	Fee        decimal.Decimal    `json:"fee"`
	CTime      time.Time          `json:"cTime"`
	UTime      time.Time          `json:"uTime"`
}
