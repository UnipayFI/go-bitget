package earn

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// LoanDailyTerm is the crypto-loan mortgage term (borrow duration).
type LoanDailyTerm string

const (
	LoanDailySeven    LoanDailyTerm = "SEVEN"    // 7-day fixed term
	LoanDailyThirty   LoanDailyTerm = "THIRTY"   // 30-day fixed term
	LoanDailyFlexible LoanDailyTerm = "FLEXIBLE" // flexible term
)

// LoanBorrowStatus is the terminal status of a finished loan order.
type LoanBorrowStatus string

const (
	LoanBorrowRollback LoanBorrowStatus = "ROLLBACK" // borrow failed / rolled back
	LoanBorrowForce    LoanBorrowStatus = "FORCE"    // force liquidated
	LoanBorrowRepay    LoanBorrowStatus = "REPAY"    // repaid
)

// LoanReviseType is the collateral adjustment direction for revise-pledge.
type LoanReviseType string

const (
	LoanReviseIn  LoanReviseType = "IN"  // supplement collateral
	LoanReviseOut LoanReviseType = "OUT" // withdraw collateral
)

// LoanReviseSide is the collateral adjustment direction reported in the pledge-rate history.
type LoanReviseSide string

const (
	LoanReviseSideDown LoanReviseSide = "down" // add collateral
	LoanReviseSideUp   LoanReviseSide = "up"   // withdraw collateral
)

// LoanReduceStatus is the status of a liquidation (reduce) record.
type LoanReduceStatus string

const (
	LoanReduceComplete LoanReduceStatus = "COMPLETE" // liquidation finished
	LoanReduceWait     LoanReduceStatus = "WAIT"     // liquidating
)

// GetLoanCoinInfosService -- GET /api/v2/earn/loan/public/coinInfos (public)
//
// Returns the crypto-loan currency list: per-loan-coin interest rates and limits
// and per-pledge-coin collateral rates and limits.
type GetLoanCoinInfosService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanCoinInfosService() *GetLoanCoinInfosService {
	return &GetLoanCoinInfosService{c: c, params: map[string]string{}}
}

// SetCoin filters the currency list to a single coin.
func (s *GetLoanCoinInfosService) SetCoin(coin string) *GetLoanCoinInfosService {
	s.params["coin"] = coin
	return s
}

func (s *GetLoanCoinInfosService) Do(ctx context.Context) (*LoanCoinInfos, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/public/coinInfos", s.params)
	return request.Do[LoanCoinInfos](req)
}

// LoanCoinInfos is the crypto-loan currency configuration.
type LoanCoinInfos struct {
	LoanInfos   []LoanCoinInfo   `json:"loanInfos"`   // borrowable-coin rate/limit table
	PledgeInfos []PledgeCoinInfo `json:"pledgeInfos"` // collateral-coin rate/limit table
}

// LoanCoinInfo is one borrowable coin's interest rates and borrowing limits.
type LoanCoinInfo struct {
	Coin        string          `json:"coin"`        // loan currency name
	HourRate7D  decimal.Decimal `json:"hourRate7D"`  // hourly rate for 7-day fixed term
	Rate7D      decimal.Decimal `json:"rate7D"`      // annualized rate for 7-day fixed term
	HourRate30D decimal.Decimal `json:"hourRate30D"` // hourly rate for 30-day fixed term
	Rate30D     decimal.Decimal `json:"rate30D"`     // annualized rate for 30-day fixed term
	MinUSDT     decimal.Decimal `json:"minUsdt"`     // minimum borrowable amount in USDT
	MaxUSDT     decimal.Decimal `json:"maxUsdt"`     // maximum borrowable amount in USDT
	Min         decimal.Decimal `json:"min"`         // minimum borrowing limit
	Max         decimal.Decimal `json:"max"`         // maximum borrowing limit
}

// PledgeCoinInfo is one collateral coin's pledge rates and limits.
type PledgeCoinInfo struct {
	Coin      string          `json:"coin"`      // collateral currency name
	InitRate  decimal.Decimal `json:"initRate"`  // initial collateral rate
	SupRate   decimal.Decimal `json:"supRate"`   // supplementary collateral rate
	ForceRate decimal.Decimal `json:"forceRate"` // forced-liquidation collateral rate
	MinUSDT   decimal.Decimal `json:"minUsdt"`   // minimum collateral limit in USDT
	MaxUSDT   decimal.Decimal `json:"maxUsdt"`   // maximum collateral limit in USDT
}

// GetLoanHourInterestService -- GET /api/v2/earn/loan/public/hour-interest (public)
//
// Returns the estimated hourly interest and the borrowable amount for a given
// loan/pledge pair, term and pledge amount.
type GetLoanHourInterestService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanHourInterestService(loanCoin, pledgeCoin string, daily LoanDailyTerm, pledgeAmount decimal.Decimal) *GetLoanHourInterestService {
	return &GetLoanHourInterestService{c: c, params: map[string]string{
		"loanCoin":     loanCoin,
		"pledgeCoin":   pledgeCoin,
		"daily":        string(daily),
		"pledgeAmount": pledgeAmount.String(),
	}}
}

func (s *GetLoanHourInterestService) Do(ctx context.Context) (*LoanHourInterest, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/public/hour-interest", s.params)
	return request.Do[LoanHourInterest](req)
}

// LoanHourInterest is the estimated hourly interest and borrowable amount.
type LoanHourInterest struct {
	HourInterest decimal.Decimal `json:"hourInterest"` // estimated interest amount per hour
	LoanAmount   decimal.Decimal `json:"loanAmount"`   // borrowable amount
}

// LoanBorrowService -- POST /api/v2/earn/loan/borrow (earn write)
//
// Borrows a coin against pledged collateral. Exactly one of pledgeAmount or
// loanAmount must be supplied.
type LoanBorrowService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewLoanBorrowService(loanCoin, pledgeCoin string, daily LoanDailyTerm) *LoanBorrowService {
	return &LoanBorrowService{c: c, body: map[string]any{
		"loanCoin":   loanCoin,
		"pledgeCoin": pledgeCoin,
		"daily":      string(daily),
	}}
}

// SetPledgeAmount sets the collateral amount (supply this OR loanAmount).
func (s *LoanBorrowService) SetPledgeAmount(amount decimal.Decimal) *LoanBorrowService {
	s.body["pledgeAmount"] = amount.String()
	return s
}

// SetLoanAmount sets the borrow amount (supply this OR pledgeAmount).
func (s *LoanBorrowService) SetLoanAmount(amount decimal.Decimal) *LoanBorrowService {
	s.body["loanAmount"] = amount.String()
	return s
}

func (s *LoanBorrowService) Do(ctx context.Context) (*LoanBorrowResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/loan/borrow", s.body).WithSign()
	return request.Do[LoanBorrowResult](req)
}

// LoanBorrowResult is the result of a borrow request.
type LoanBorrowResult struct {
	OrderID string `json:"orderId"` // loan order ID
}

// GetLoanOngoingOrdersService -- GET /api/v2/earn/loan/ongoing-orders (earn read)
//
// Returns the account's ongoing (active) loan orders.
type GetLoanOngoingOrdersService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanOngoingOrdersService() *GetLoanOngoingOrdersService {
	return &GetLoanOngoingOrdersService{c: c, params: map[string]string{}}
}

// SetOrderID filters to a single loan order.
func (s *GetLoanOngoingOrdersService) SetOrderID(orderID string) *GetLoanOngoingOrdersService {
	s.params["orderId"] = orderID
	return s
}

// SetLoanCoin filters by borrowed coin.
func (s *GetLoanOngoingOrdersService) SetLoanCoin(loanCoin string) *GetLoanOngoingOrdersService {
	s.params["loanCoin"] = loanCoin
	return s
}

// SetPledgeCoin filters by collateral coin.
func (s *GetLoanOngoingOrdersService) SetPledgeCoin(pledgeCoin string) *GetLoanOngoingOrdersService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

func (s *GetLoanOngoingOrdersService) Do(ctx context.Context) ([]LoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/ongoing-orders", s.params).WithSign()
	resp, err := request.Do[[]LoanOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LoanOrder is one ongoing loan order.
type LoanOrder struct {
	OrderID          string          `json:"orderId"`          // order ID
	LoanCoin         string          `json:"loanCoin"`         // borrowed coin
	LoanAmount       decimal.Decimal `json:"loanAmount"`       // loan amount
	InterestAmount   decimal.Decimal `json:"interestAmount"`   // accrued interest amount
	HourInterestRate decimal.Decimal `json:"hourInterestRate"` // hourly interest rate
	PledgeCoin       string          `json:"pledgeCoin"`       // collateral coin
	PledgeAmount     decimal.Decimal `json:"pledgeAmount"`     // collateral amount
	PledgeRate       decimal.Decimal `json:"pledgeRate"`       // collateral rate
	SupRate          decimal.Decimal `json:"supRate"`          // supplementary collateral rate
	ForceRate        decimal.Decimal `json:"forceRate"`        // forced-liquidation collateral rate
	BorrowTime       time.Time       `json:"borrowTime"`       // borrow time
	ExpireTime       time.Time       `json:"expireTime"`       // expiry time
}

// LoanRepayService -- POST /api/v2/earn/loan/repay (earn write)
//
// Repays a loan order, optionally in full and optionally redeeming collateral.
type LoanRepayService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewLoanRepayService(orderID string, repayAll bool) *LoanRepayService {
	all := "no"
	if repayAll {
		all = "yes"
	}
	return &LoanRepayService{c: c, body: map[string]any{
		"orderId":  orderID,
		"repayAll": all,
	}}
}

// SetAmount sets the repayment amount (used when repayAll is no).
func (s *LoanRepayService) SetAmount(amount decimal.Decimal) *LoanRepayService {
	s.body["amount"] = amount.String()
	return s
}

// SetRepayUnlock controls whether collateral is redeemed after repay (default yes).
func (s *LoanRepayService) SetRepayUnlock(unlock bool) *LoanRepayService {
	if unlock {
		s.body["repayUnlock"] = "yes"
	} else {
		s.body["repayUnlock"] = "no"
	}
	return s
}

func (s *LoanRepayService) Do(ctx context.Context) (*LoanRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/loan/repay", s.body).WithSign()
	return request.Do[LoanRepayResult](req)
}

// LoanRepayResult is the result of a repayment.
type LoanRepayResult struct {
	LoanCoin          string          `json:"loanCoin"`          // borrowed coin
	PledgeCoin        string          `json:"pledgeCoin"`        // collateral coin
	RepayAmount       decimal.Decimal `json:"repayAmount"`       // total repayment amount
	PayInterest       decimal.Decimal `json:"payInterest"`       // interest paid
	RepayLoanAmount   decimal.Decimal `json:"repayLoanAmount"`   // principal repaid
	RepayUnlockAmount decimal.Decimal `json:"repayUnlockAmount"` // collateral redeemed
}

// GetLoanRepayHistoryService -- GET /api/v2/earn/loan/repay-history (earn read)
//
// Returns the account's repayment history (startTime/endTime required, past 3
// months only).
type GetLoanRepayHistoryService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanRepayHistoryService(startTime, endTime time.Time) *GetLoanRepayHistoryService {
	return &GetLoanRepayHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetOrderID filters to a single loan order.
func (s *GetLoanRepayHistoryService) SetOrderID(orderID string) *GetLoanRepayHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetLoanCoin filters by borrowed coin.
func (s *GetLoanRepayHistoryService) SetLoanCoin(loanCoin string) *GetLoanRepayHistoryService {
	s.params["loanCoin"] = loanCoin
	return s
}

// SetPledgeCoin filters by collateral coin.
func (s *GetLoanRepayHistoryService) SetPledgeCoin(pledgeCoin string) *GetLoanRepayHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetPageNo sets the 1-based page number (default 1).
func (s *GetLoanRepayHistoryService) SetPageNo(pageNo int) *GetLoanRepayHistoryService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 10, max 100).
func (s *GetLoanRepayHistoryService) SetPageSize(pageSize int) *GetLoanRepayHistoryService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetLoanRepayHistoryService) Do(ctx context.Context) ([]LoanRepayRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/repay-history", s.params).WithSign()
	resp, err := request.Do[[]LoanRepayRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LoanRepayRecord is one repayment-history entry.
type LoanRepayRecord struct {
	OrderID           string          `json:"orderId"`           // order ID
	LoanCoin          string          `json:"loanCoin"`          // borrowed coin
	PledgeCoin        string          `json:"pledgeCoin"`        // collateral coin
	RepayAmount       decimal.Decimal `json:"repayAmount"`       // repayment amount
	PayInterest       decimal.Decimal `json:"payInterest"`       // interest paid
	RepayLoanAmount   decimal.Decimal `json:"repayLoanAmount"`   // principal repaid
	RepayUnlockAmount decimal.Decimal `json:"repayUnlockAmount"` // collateral released
	RepayTime         time.Time       `json:"repayTime"`         // repayment time
}

// LoanRevisePledgeService -- POST /api/v2/earn/loan/revise-pledge (earn write)
//
// Adjusts the collateral (pledge) on a loan order, supplementing (IN) or
// withdrawing (OUT) the given amount.
type LoanRevisePledgeService struct {
	c    *EarnClient
	body map[string]any
}

func (c *EarnClient) NewLoanRevisePledgeService(orderID, pledgeCoin string, reviseType LoanReviseType, amount decimal.Decimal) *LoanRevisePledgeService {
	return &LoanRevisePledgeService{c: c, body: map[string]any{
		"orderId":    orderID,
		"pledgeCoin": pledgeCoin,
		"reviseType": string(reviseType),
		"amount":     amount.String(),
	}}
}

func (s *LoanRevisePledgeService) Do(ctx context.Context) (*LoanRevisePledgeResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/earn/loan/revise-pledge", s.body).WithSign()
	return request.Do[LoanRevisePledgeResult](req)
}

// LoanRevisePledgeResult is the result of a pledge-rate adjustment.
type LoanRevisePledgeResult struct {
	LoanCoin        string          `json:"loanCoin"`        // borrowed coin
	PledgeCoin      string          `json:"pledgeCoin"`      // collateral coin
	AfterPledgeRate decimal.Decimal `json:"afterPledgeRate"` // collateral rate after adjustment
}

// GetLoanReviseHistoryService -- GET /api/v2/earn/loan/revise-history (earn read)
//
// Returns the pledge-rate adjustment history (startTime/endTime required, past 3
// months only).
type GetLoanReviseHistoryService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanReviseHistoryService(startTime, endTime time.Time) *GetLoanReviseHistoryService {
	return &GetLoanReviseHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetOrderID filters to a single loan order.
func (s *GetLoanReviseHistoryService) SetOrderID(orderID string) *GetLoanReviseHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetReviseSide filters by adjustment direction (down = add, up = withdraw).
func (s *GetLoanReviseHistoryService) SetReviseSide(side LoanReviseSide) *GetLoanReviseHistoryService {
	s.params["reviseSide"] = string(side)
	return s
}

// SetPledgeCoin filters by collateral coin.
func (s *GetLoanReviseHistoryService) SetPledgeCoin(pledgeCoin string) *GetLoanReviseHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetPageNo sets the 1-based page number (default 1).
func (s *GetLoanReviseHistoryService) SetPageNo(pageNo int) *GetLoanReviseHistoryService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 10, max 100).
func (s *GetLoanReviseHistoryService) SetPageSize(pageSize int) *GetLoanReviseHistoryService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetLoanReviseHistoryService) Do(ctx context.Context) ([]LoanReviseRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/revise-history", s.params).WithSign()
	resp, err := request.Do[[]LoanReviseRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LoanReviseRecord is one pledge-rate adjustment record.
type LoanReviseRecord struct {
	LoanCoin         string          `json:"loanCoin"`         // borrowed coin
	PledgeCoin       string          `json:"pledgeCoin"`       // collateral coin
	OrderID          string          `json:"orderId"`          // loan order ID
	ReviseTime       time.Time       `json:"reviseTime"`       // adjustment time
	ReviseSide       LoanReviseSide  `json:"reviseSide"`       // adjustment direction
	ReviseAmount     decimal.Decimal `json:"reviseAmount"`     // adjustment amount
	AfterPledgeRate  decimal.Decimal `json:"afterPledgeRate"`  // collateral rate after adjustment
	BeforePledgeRate decimal.Decimal `json:"beforePledgeRate"` // collateral rate before adjustment
}

// GetLoanBorrowHistoryService -- GET /api/v2/earn/loan/borrow-history (earn read)
//
// Returns the account's loan (borrow) history (startTime/endTime required, past
// 90 days only).
type GetLoanBorrowHistoryService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanBorrowHistoryService(startTime, endTime time.Time) *GetLoanBorrowHistoryService {
	return &GetLoanBorrowHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetOrderID filters to a single loan order.
func (s *GetLoanBorrowHistoryService) SetOrderID(orderID string) *GetLoanBorrowHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetLoanCoin filters by borrowed coin.
func (s *GetLoanBorrowHistoryService) SetLoanCoin(loanCoin string) *GetLoanBorrowHistoryService {
	s.params["loanCoin"] = loanCoin
	return s
}

// SetPledgeCoin filters by collateral coin.
func (s *GetLoanBorrowHistoryService) SetPledgeCoin(pledgeCoin string) *GetLoanBorrowHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetStatus filters by terminal loan status (ROLLBACK, FORCE, REPAY).
func (s *GetLoanBorrowHistoryService) SetStatus(status LoanBorrowStatus) *GetLoanBorrowHistoryService {
	s.params["status"] = string(status)
	return s
}

// SetPageNo sets the 1-based page number (default 1).
func (s *GetLoanBorrowHistoryService) SetPageNo(pageNo int) *GetLoanBorrowHistoryService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 10, max 100).
func (s *GetLoanBorrowHistoryService) SetPageSize(pageSize int) *GetLoanBorrowHistoryService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetLoanBorrowHistoryService) Do(ctx context.Context) ([]LoanBorrowRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/borrow-history", s.params).WithSign()
	resp, err := request.Do[[]LoanBorrowRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LoanBorrowRecord is one loan-history entry.
type LoanBorrowRecord struct {
	OrderID          string           `json:"orderId"`          // order ID
	LoanCoin         string           `json:"loanCoin"`         // borrowed coin
	PledgeCoin       string           `json:"pledgeCoin"`       // collateral coin
	InitPledgeAmount decimal.Decimal  `json:"initPledgeAmount"` // initial collateral amount
	InitLoanAmount   decimal.Decimal  `json:"initLoanAmount"`   // initial loan amount
	HourRate         decimal.Decimal  `json:"hourRate"`         // hourly interest rate
	Daily            LoanDailyTerm    `json:"daily"`            // pledge duration term
	BorrowTime       time.Time        `json:"borrowTime"`       // borrow time
	Status           LoanBorrowStatus `json:"status"`           // terminal status
}

// GetLoanDebtsService -- GET /api/v2/earn/loan/debts (earn read)
//
// Returns the account's outstanding loan debts and pledged collateral.
type GetLoanDebtsService struct {
	c *EarnClient
}

func (c *EarnClient) NewGetLoanDebtsService() *GetLoanDebtsService {
	return &GetLoanDebtsService{c: c}
}

func (s *GetLoanDebtsService) Do(ctx context.Context) (*LoanDebts, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/debts").WithSign()
	return request.Do[LoanDebts](req)
}

// LoanDebts is the account's outstanding debt and collateral snapshot.
type LoanDebts struct {
	PledgeInfos []LoanDebtItem `json:"pledgeInfos"` // pledged collateral per coin
	LoanInfos   []LoanDebtItem `json:"loanInfos"`   // outstanding loan per coin
}

// LoanDebtItem is one coin's debt or collateral balance.
type LoanDebtItem struct {
	Coin       string          `json:"coin"`       // coin symbol
	Amount     decimal.Decimal `json:"amount"`     // amount
	AmountUSDT decimal.Decimal `json:"amountUsdt"` // USDT-equivalent amount
}

// GetLoanReducesService -- GET /api/v2/earn/loan/reduces (earn read)
//
// Returns the account's liquidation (reduce) records (startTime/endTime
// required, past 3 months only).
type GetLoanReducesService struct {
	c      *EarnClient
	params map[string]string
}

func (c *EarnClient) NewGetLoanReducesService(startTime, endTime time.Time) *GetLoanReducesService {
	return &GetLoanReducesService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetOrderID filters to a single loan order.
func (s *GetLoanReducesService) SetOrderID(orderID string) *GetLoanReducesService {
	s.params["orderId"] = orderID
	return s
}

// SetLoanCoin filters by borrowed coin.
func (s *GetLoanReducesService) SetLoanCoin(loanCoin string) *GetLoanReducesService {
	s.params["loanCoin"] = loanCoin
	return s
}

// SetPledgeCoin filters by collateral coin.
func (s *GetLoanReducesService) SetPledgeCoin(pledgeCoin string) *GetLoanReducesService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetStatus filters by liquidation status (COMPLETE, WAIT).
func (s *GetLoanReducesService) SetStatus(status LoanReduceStatus) *GetLoanReducesService {
	s.params["status"] = string(status)
	return s
}

// SetPageNo sets the 1-based page number (default 1).
func (s *GetLoanReducesService) SetPageNo(pageNo int) *GetLoanReducesService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 10, max 100).
func (s *GetLoanReducesService) SetPageSize(pageSize int) *GetLoanReducesService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetLoanReducesService) Do(ctx context.Context) ([]LoanReduceRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/earn/loan/reduces", s.params).WithSign()
	resp, err := request.Do[[]LoanReduceRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// LoanReduceRecord is one liquidation (reduce) record.
type LoanReduceRecord struct {
	OrderID         string           `json:"orderId"`         // order ID
	LoanCoin        string           `json:"loanCoin"`        // borrowed coin
	PledgeCoin      string           `json:"pledgeCoin"`      // collateral coin
	ReduceTime      time.Time        `json:"reduceTime"`      // liquidation time
	PledgeRate      decimal.Decimal  `json:"pledgeRate"`      // collateral rate at liquidation
	PledgePrice     decimal.Decimal  `json:"pledgePrice"`     // collateral price at liquidation
	Status          LoanReduceStatus `json:"status"`          // liquidation status
	PledgeAmount    decimal.Decimal  `json:"pledgeAmount"`    // liquidated collateral amount
	ReduceFee       decimal.Decimal  `json:"reduceFee"`       // liquidation fee
	ResidueAmount   decimal.Decimal  `json:"residueAmount"`   // remaining collateral balance
	RunlockAmount   decimal.Decimal  `json:"runlockAmount"`   // released collateral amount
	RepayLoanAmount decimal.Decimal  `json:"repayLoanAmount"` // loan repaid
}
