package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CrossBorrowType is the way a cross-margin loan was created: automatically when
// placing an autoLoan order, or manually via the borrow endpoint.
type CrossBorrowType string

const (
	CrossBorrowTypeAutoLoan   CrossBorrowType = "auto_loan"
	CrossBorrowTypeManualLoan CrossBorrowType = "manual_loan"
)

// CrossRepayType is how a cross-margin repayment was triggered.
type CrossRepayType string

const (
	CrossRepayTypeAutoRepay   CrossRepayType = "auto_repay"
	CrossRepayTypeManualRepay CrossRepayType = "manual_repay"
	CrossRepayTypeLiqRepay    CrossRepayType = "liq_repay"
	CrossRepayTypeForceRepay  CrossRepayType = "force_repay"
)

// CrossInterestType distinguishes the first interest charged at borrowing time
// from the recurring scheduled interest.
type CrossInterestType string

const (
	CrossInterestTypeFirst     CrossInterestType = "first"
	CrossInterestTypeScheduled CrossInterestType = "scheduled"
)

// CrossMarginFlowType is the capital-flow category of a cross-margin financial
// record (transfer/borrow/repay/liquidation/exchange flows).
type CrossMarginFlowType string

const (
	CrossMarginFlowTypeTransferIn     CrossMarginFlowType = "transfer_in"
	CrossMarginFlowTypeTransferOut    CrossMarginFlowType = "transfer_out"
	CrossMarginFlowTypeBorrow         CrossMarginFlowType = "borrow"
	CrossMarginFlowTypeRepay          CrossMarginFlowType = "repay"
	CrossMarginFlowTypeLiquidationFee CrossMarginFlowType = "liquidation_fee"
	CrossMarginFlowTypeCompensate     CrossMarginFlowType = "compensate"
	CrossMarginFlowTypeDealIn         CrossMarginFlowType = "deal_in"
	CrossMarginFlowTypeDealOut        CrossMarginFlowType = "deal_out"
	CrossMarginFlowTypeConfiscated    CrossMarginFlowType = "confiscated"
	CrossMarginFlowTypeExchangeIn     CrossMarginFlowType = "exchange_in"
	CrossMarginFlowTypeExchangeOut    CrossMarginFlowType = "exchange_out"
	CrossMarginFlowTypeSysExchangeIn  CrossMarginFlowType = "sys_exchange_in"
	CrossMarginFlowTypeSysExchangeOut CrossMarginFlowType = "sys_exchange_out"
)

// GetCrossBorrowHistoryService -- GET /api/v2/margin/crossed/borrow-history (cross margin read)
//
// Returns the cross-margin borrow (loan) records within a time window.
type GetCrossBorrowHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossBorrowHistoryService(startTime time.Time) *GetCrossBorrowHistoryService {
	return &GetCrossBorrowHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetCrossBorrowHistoryService) SetLoanId(loanId string) *GetCrossBorrowHistoryService {
	s.params["loanId"] = loanId
	return s
}

func (s *GetCrossBorrowHistoryService) SetCoin(coin string) *GetCrossBorrowHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetCrossBorrowHistoryService) SetEndTime(endTime time.Time) *GetCrossBorrowHistoryService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetCrossBorrowHistoryService) SetLimit(limit int) *GetCrossBorrowHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCrossBorrowHistoryService) SetIdLessThan(idLessThan string) *GetCrossBorrowHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetCrossBorrowHistoryService) Do(ctx context.Context) (*CrossBorrowHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/borrow-history", s.params).WithSign()
	return request.Do[CrossBorrowHistoryResponse](req)
}

// CrossBorrowHistoryResponse is the paginated cross-margin borrow history page.
type CrossBorrowHistoryResponse struct {
	ResultList []CrossBorrowRecord `json:"resultList"`
	MaxId      string              `json:"maxId"`
	MinId      string              `json:"minId"`
}

// CrossBorrowRecord is a single cross-margin loan record.
type CrossBorrowRecord struct {
	LoanId       string          `json:"loanId"`
	Coin         string          `json:"coin"`
	BorrowAmount decimal.Decimal `json:"borrowAmount"`
	BorrowType   CrossBorrowType `json:"borrowType"`
	CTime        time.Time       `json:"cTime"`
	UTime        time.Time       `json:"uTime"`
}

// GetCrossRepayHistoryService -- GET /api/v2/margin/crossed/repay-history (cross margin read)
//
// Returns the cross-margin repayment records within a time window.
type GetCrossRepayHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossRepayHistoryService(startTime time.Time) *GetCrossRepayHistoryService {
	return &GetCrossRepayHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetCrossRepayHistoryService) SetRepayId(repayId string) *GetCrossRepayHistoryService {
	s.params["repayId"] = repayId
	return s
}

func (s *GetCrossRepayHistoryService) SetCoin(coin string) *GetCrossRepayHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetCrossRepayHistoryService) SetEndTime(endTime time.Time) *GetCrossRepayHistoryService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetCrossRepayHistoryService) SetLimit(limit int) *GetCrossRepayHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCrossRepayHistoryService) SetIdLessThan(idLessThan string) *GetCrossRepayHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetCrossRepayHistoryService) Do(ctx context.Context) (*CrossRepayHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/repay-history", s.params).WithSign()
	return request.Do[CrossRepayHistoryResponse](req)
}

// CrossRepayHistoryResponse is the paginated cross-margin repay history page.
type CrossRepayHistoryResponse struct {
	ResultList []CrossRepayRecord `json:"resultList"`
	MaxId      string             `json:"maxId"`
	MinId      string             `json:"minId"`
}

// CrossRepayRecord is a single cross-margin repayment record.
type CrossRepayRecord struct {
	RepayId        string          `json:"repayId"`
	Coin           string          `json:"coin"`
	RepayPrincipal decimal.Decimal `json:"repayPrincipal"`
	RepayAmount    decimal.Decimal `json:"repayAmount"`
	RepayInterest  decimal.Decimal `json:"repayInterest"`
	RepayType      CrossRepayType  `json:"repayType"`
	CTime          time.Time       `json:"cTime"`
	UTime          time.Time       `json:"uTime"`
}

// GetCrossInterestHistoryService -- GET /api/v2/margin/crossed/interest-history (cross margin read)
//
// Returns the cross-margin interest accrual records within a time window.
type GetCrossInterestHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossInterestHistoryService(startTime time.Time) *GetCrossInterestHistoryService {
	return &GetCrossInterestHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetCrossInterestHistoryService) SetCoin(coin string) *GetCrossInterestHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetCrossInterestHistoryService) SetEndTime(endTime time.Time) *GetCrossInterestHistoryService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetCrossInterestHistoryService) SetLimit(limit int) *GetCrossInterestHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCrossInterestHistoryService) SetIdLessThan(idLessThan string) *GetCrossInterestHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetCrossInterestHistoryService) Do(ctx context.Context) (*CrossInterestHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/interest-history", s.params).WithSign()
	return request.Do[CrossInterestHistoryResponse](req)
}

// CrossInterestHistoryResponse is the paginated cross-margin interest history page.
type CrossInterestHistoryResponse struct {
	ResultList []CrossInterestRecord `json:"resultList"`
	MaxId      string                `json:"maxId"`
	MinId      string                `json:"minId"`
}

// CrossInterestRecord is a single cross-margin interest accrual record.
type CrossInterestRecord struct {
	InterestId        string            `json:"interestId"`
	LoanCoin          string            `json:"loanCoin"`
	InterestCoin      string            `json:"interestCoin"`
	DailyInterestRate decimal.Decimal   `json:"dailyInterestRate"`
	InterestAmount    decimal.Decimal   `json:"interestAmount"`
	InterstType       CrossInterestType `json:"interstType"` // key spelled "interstType" verbatim per Bitget docs
	CTime             time.Time         `json:"cTime"`
	UTime             time.Time         `json:"uTime"`
}

// GetCrossLiquidationHistoryService -- GET /api/v2/margin/crossed/liquidation-history (cross margin read)
//
// Returns the cross-margin liquidation records within a time window.
type GetCrossLiquidationHistoryService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossLiquidationHistoryService(startTime time.Time) *GetCrossLiquidationHistoryService {
	return &GetCrossLiquidationHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetCrossLiquidationHistoryService) SetEndTime(endTime time.Time) *GetCrossLiquidationHistoryService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetCrossLiquidationHistoryService) SetLimit(limit int) *GetCrossLiquidationHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCrossLiquidationHistoryService) SetIdLessThan(idLessThan string) *GetCrossLiquidationHistoryService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetCrossLiquidationHistoryService) Do(ctx context.Context) (*CrossLiquidationHistoryResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/liquidation-history", s.params).WithSign()
	return request.Do[CrossLiquidationHistoryResponse](req)
}

// CrossLiquidationHistoryResponse is the paginated cross-margin liquidation history page.
type CrossLiquidationHistoryResponse struct {
	ResultList []CrossLiquidationRecord `json:"resultList"`
	MaxId      string                   `json:"maxId"`
	MinId      string                   `json:"minId"`
}

// CrossLiquidationRecord is a single cross-margin liquidation record. Total
// assets/debt are denominated in USDT.
type CrossLiquidationRecord struct {
	LiqId        string          `json:"liqId"`
	LiqStartTime time.Time       `json:"liqStartTime"`
	LiqEndTime   time.Time       `json:"liqEndTime"`
	LiqRiskRatio decimal.Decimal `json:"liqRiskRatio"`
	TotalAssets  decimal.Decimal `json:"totalAssets"`
	TotalDebt    decimal.Decimal `json:"totalDebt"`
	LiqFee       decimal.Decimal `json:"liqFee"`
	CTime        time.Time       `json:"cTime"`
	UTime        time.Time       `json:"uTime"`
}

// GetCrossFinancialRecordsService -- GET /api/v2/margin/crossed/financial-records (cross margin read)
//
// Returns the cross-margin capital-flow (finance) records within a time window.
type GetCrossFinancialRecordsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossFinancialRecordsService(startTime time.Time) *GetCrossFinancialRecordsService {
	return &GetCrossFinancialRecordsService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
	}}
}

func (s *GetCrossFinancialRecordsService) SetMarginType(marginType CrossMarginFlowType) *GetCrossFinancialRecordsService {
	s.params["marginType"] = string(marginType)
	return s
}

func (s *GetCrossFinancialRecordsService) SetCoin(coin string) *GetCrossFinancialRecordsService {
	s.params["coin"] = coin
	return s
}

func (s *GetCrossFinancialRecordsService) SetEndTime(endTime time.Time) *GetCrossFinancialRecordsService {
	s.params["endTime"] = strconv.FormatInt(endTime.UnixMilli(), 10)
	return s
}

func (s *GetCrossFinancialRecordsService) SetLimit(limit int) *GetCrossFinancialRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetCrossFinancialRecordsService) SetIdLessThan(idLessThan string) *GetCrossFinancialRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetCrossFinancialRecordsService) Do(ctx context.Context) (*CrossFinancialRecordsResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/financial-records", s.params).WithSign()
	return request.Do[CrossFinancialRecordsResponse](req)
}

// CrossFinancialRecordsResponse is the paginated cross-margin financial-flow page.
type CrossFinancialRecordsResponse struct {
	ResultList []CrossFinancialRecord `json:"resultList"`
	MaxId      string                 `json:"maxId"`
	MinId      string                 `json:"minId"`
}

// CrossFinancialRecord is a single cross-margin capital-flow record.
type CrossFinancialRecord struct {
	MarginId   string              `json:"marginId"`
	Amount     decimal.Decimal     `json:"amount"`
	Coin       string              `json:"coin"`
	Balance    decimal.Decimal     `json:"balance"`
	Fee        decimal.Decimal     `json:"fee"`
	MarginType CrossMarginFlowType `json:"marginType"`
	UTime      time.Time           `json:"uTime"`
	CTime      time.Time           `json:"cTime"`
}
