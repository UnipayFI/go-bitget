package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetLoanCoinsService -- GET /api/v3/loan/coins (UTA mgt. read)
//
// Returns the coins available for crypto-loan borrowing and collateral, with
// their rate tiers and borrow/pledge limits, optionally filtered to one coin.
type GetLoanCoinsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetLoanCoinsService() *GetLoanCoinsService {
	return &GetLoanCoinsService{c: c, params: map[string]string{}}
}

func (s *GetLoanCoinsService) SetCoin(coin string) *GetLoanCoinsService {
	s.params["coin"] = coin
	return s
}

func (s *GetLoanCoinsService) Do(ctx context.Context) (*LoanCoins, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/coins", s.params).WithSign()
	return request.Do[LoanCoins](req)
}

type LoanCoins struct {
	LoanInfos   []LoanCoinInfo   `json:"loanInfos"`
	PledgeInfos []PledgeCoinInfo `json:"pledgeInfos"`
}

type LoanCoinInfo struct {
	Coin             string          `json:"coin"`
	HourRateFlexible decimal.Decimal `json:"hourRateFlexible"`
	RateFlexible     decimal.Decimal `json:"rateFlexible"`
	HourRate7D       decimal.Decimal `json:"hourRate7D"`
	Rate7D           decimal.Decimal `json:"rate7D"`
	HourRate30D      decimal.Decimal `json:"hourRate30D"`
	Rate30D          decimal.Decimal `json:"rate30D"`
	MinBorrowAmount  decimal.Decimal `json:"minBorrowAmount"`
	MaxBorrowAmount  decimal.Decimal `json:"maxBorrowAmount"`
	MinBorrowLimit   decimal.Decimal `json:"minBorrowLimit"`
	MaxBorrowLimit   decimal.Decimal `json:"maxBorrowLimit"`
}

type PledgeCoinInfo struct {
	Coin            string          `json:"coin"`
	InitRate        decimal.Decimal `json:"initRate"`
	SupRate         decimal.Decimal `json:"supRate"`
	ForceRate       decimal.Decimal `json:"forceRate"`
	MinPledgeAmount decimal.Decimal `json:"minPledgeAmount"`
	MaxPledgeAmount decimal.Decimal `json:"maxPledgeAmount"`
}

// GetLoanInterestService -- GET /api/v3/loan/interest (UTA mgt. read)
//
// Estimates the hourly interest and borrowable amount for a loan/collateral
// coin pair, pledge term, and collateral amount.
type GetLoanInterestService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetLoanInterestService(loanCoin, pledgeCoin, daily string, pledgeAmount decimal.Decimal) *GetLoanInterestService {
	return &GetLoanInterestService{c: c, params: map[string]string{
		"loanCoin":     loanCoin,
		"pledgeCoin":   pledgeCoin,
		"daily":        daily,
		"pledgeAmount": pledgeAmount.String(),
	}}
}

func (s *GetLoanInterestService) Do(ctx context.Context) (*LoanInterest, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/interest", s.params).WithSign()
	return request.Do[LoanInterest](req)
}

type LoanInterest struct {
	HourInterest decimal.Decimal `json:"hourInterest"`
	LoanAmount   decimal.Decimal `json:"loanAmount"`
}

// BorrowCoinsService -- POST /api/v3/loan/borrow (UTA mgt. read & write)
//
// Opens a crypto loan against collateral. Supply exactly one of pledgeAmount or
// loanAmount; the exchange derives the other. The reply carries the order id.
type BorrowCoinsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewBorrowCoinsService(loanCoin, pledgeCoin, daily string) *BorrowCoinsService {
	return &BorrowCoinsService{c: c, body: map[string]any{
		"loanCoin":   loanCoin,
		"pledgeCoin": pledgeCoin,
		"daily":      daily,
	}}
}

// SetPledgeAmount sets the collateral amount (choose one of pledgeAmount or
// loanAmount).
func (s *BorrowCoinsService) SetPledgeAmount(pledgeAmount decimal.Decimal) *BorrowCoinsService {
	s.body["pledgeAmount"] = pledgeAmount.String()
	return s
}

// SetLoanAmount sets the borrow amount (choose one of loanAmount or
// pledgeAmount).
func (s *BorrowCoinsService) SetLoanAmount(loanAmount decimal.Decimal) *BorrowCoinsService {
	s.body["loanAmount"] = loanAmount.String()
	return s
}

func (s *BorrowCoinsService) Do(ctx context.Context) (*BorrowResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/loan/borrow", s.body).WithSign()
	return request.Do[BorrowResult](req)
}

type BorrowResult struct {
	OrderID string `json:"orderId"`
}

// GetBorrowOngoingService -- GET /api/v3/loan/borrow-ongoing (UTA mgt. read)
//
// Returns the unified account's open (outstanding) crypto loans, optionally
// filtered by order id, loan coin, or collateral coin.
type GetBorrowOngoingService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBorrowOngoingService() *GetBorrowOngoingService {
	return &GetBorrowOngoingService{c: c, params: map[string]string{}}
}

func (s *GetBorrowOngoingService) SetOrderID(orderID string) *GetBorrowOngoingService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetBorrowOngoingService) SetLoanCoin(loanCoin string) *GetBorrowOngoingService {
	s.params["loanCoin"] = loanCoin
	return s
}

func (s *GetBorrowOngoingService) SetPledgeCoin(pledgeCoin string) *GetBorrowOngoingService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

func (s *GetBorrowOngoingService) Do(ctx context.Context) ([]BorrowOngoing, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/borrow-ongoing", s.params).WithSign()
	resp, err := request.Do[[]BorrowOngoing](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type BorrowOngoing struct {
	OrderID          string          `json:"orderId"`
	LoanCoin         string          `json:"loanCoin"`
	LoanAmount       decimal.Decimal `json:"loanAmount"`
	InterestAmount   decimal.Decimal `json:"interestAmount"`
	HourInterestRate decimal.Decimal `json:"hourInterestRate"`
	PledgeCoin       string          `json:"pledgeCoin"`
	PledgeAmount     decimal.Decimal `json:"pledgeAmount"`
	PledgeRate       decimal.Decimal `json:"pledgeRate"`
	SupRate          decimal.Decimal `json:"supRate"`
	ForceRate        decimal.Decimal `json:"forceRate"`
	BorrowTime       time.Time       `json:"borrowTime"`
	ExpireTime       time.Time       `json:"expireTime"`
}

// GetBorrowHistoryService -- GET /api/v3/loan/borrow-history (UTA mgt. read)
//
// Returns the unified account's settled (no longer ongoing) crypto loans within
// a time window, paged. The window is bounded to the last 3 months.
type GetBorrowHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBorrowHistoryService(startTime, endTime time.Time) *GetBorrowHistoryService {
	return &GetBorrowHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetBorrowHistoryService) SetOrderID(orderID string) *GetBorrowHistoryService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetBorrowHistoryService) SetLoanCoin(loanCoin string) *GetBorrowHistoryService {
	s.params["loanCoin"] = loanCoin
	return s
}

func (s *GetBorrowHistoryService) SetPledgeCoin(pledgeCoin string) *GetBorrowHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetStatus filters by loan status (ROLLBACK, FORCE, REPAY).
func (s *GetBorrowHistoryService) SetStatus(status string) *GetBorrowHistoryService {
	s.params["status"] = status
	return s
}

func (s *GetBorrowHistoryService) SetPageNum(pageNum string) *GetBorrowHistoryService {
	s.params["pageNum"] = pageNum
	return s
}

func (s *GetBorrowHistoryService) SetPageSize(pageSize string) *GetBorrowHistoryService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetBorrowHistoryService) Do(ctx context.Context) ([]BorrowHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/borrow-history", s.params).WithSign()
	resp, err := request.Do[[]BorrowHistory](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type BorrowHistory struct {
	OrderID          string          `json:"orderId"`
	LoanCoin         string          `json:"loanCoin"`
	PledgeCoin       string          `json:"pledgeCoin"`
	InitPledgeAmount decimal.Decimal `json:"initPledgeAmount"`
	InitLoanAmount   decimal.Decimal `json:"initLoanAmount"`
	HourRate         decimal.Decimal `json:"hourRate"`
	PledgeDays       string          `json:"pledgeDays"`
	BorrowTime       time.Time       `json:"borrowTime"`
	Status           string          `json:"status"` // ROLLBACK, FORCE, REPAY
	Daily            string          `json:"daily"`  // FLEXIBLE, SEVEN, THIRTY
}

// RepayCoinsService -- POST /api/v3/loan/repay (UTA mgt. read & write)
//
// Repays an open crypto loan. With repayAll "yes" the amount is ignored and the
// full debt is settled; with "no" the amount field is required. repayUnlock
// controls whether collateral is redeemed.
type RepayCoinsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewRepayCoinsService(orderID, repayAll string) *RepayCoinsService {
	return &RepayCoinsService{c: c, body: map[string]any{
		"orderId":  orderID,
		"repayAll": repayAll,
	}}
}

// SetAmount sets the repayment amount (required when repayAll is "no").
func (s *RepayCoinsService) SetAmount(amount decimal.Decimal) *RepayCoinsService {
	s.body["amount"] = amount.String()
	return s
}

// SetRepayUnlock sets the collateral redemption option ("yes" redeems, "no"
// retains; defaults to "no", ineffective when repayAll is "yes").
func (s *RepayCoinsService) SetRepayUnlock(repayUnlock string) *RepayCoinsService {
	s.body["repayUnlock"] = repayUnlock
	return s
}

func (s *RepayCoinsService) Do(ctx context.Context) (*RepayCoinsResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/loan/repay", s.body).WithSign()
	return request.Do[RepayCoinsResult](req)
}

type RepayCoinsResult struct {
	LoanCoin          string          `json:"loanCoin"`
	PledgeCoin        string          `json:"pledgeCoin"`
	RepayAmount       decimal.Decimal `json:"repayAmount"`
	PayInterest       decimal.Decimal `json:"payInterest"`
	RepayLoanAmount   decimal.Decimal `json:"repayLoanAmount"`
	RepayUnlockAmount decimal.Decimal `json:"repayUnlockAmount"`
}

// GetRepayHistoryService -- GET /api/v3/loan/repay-history (UTA mgt. read)
//
// Returns the unified account's crypto-loan repayment records within a time
// window, paged. The window is bounded to the last 3 months.
type GetRepayHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRepayHistoryService(startTime, endTime time.Time) *GetRepayHistoryService {
	return &GetRepayHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetRepayHistoryService) SetOrderID(orderID string) *GetRepayHistoryService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetRepayHistoryService) SetLoanCoin(loanCoin string) *GetRepayHistoryService {
	s.params["loanCoin"] = loanCoin
	return s
}

func (s *GetRepayHistoryService) SetPledgeCoin(pledgeCoin string) *GetRepayHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

func (s *GetRepayHistoryService) SetPageNum(pageNum string) *GetRepayHistoryService {
	s.params["pageNum"] = pageNum
	return s
}

func (s *GetRepayHistoryService) SetPageSize(pageSize string) *GetRepayHistoryService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetRepayHistoryService) Do(ctx context.Context) ([]RepayHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/repay-history", s.params).WithSign()
	resp, err := request.Do[[]RepayHistory](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type RepayHistory struct {
	OrderID           string          `json:"orderId"`
	LoanCoin          string          `json:"loanCoin"`
	PledgeCoin        string          `json:"pledgeCoin"`
	RepayAmount       decimal.Decimal `json:"repayAmount"`
	PayInterest       decimal.Decimal `json:"payInterest"`
	RepayLoanAmount   decimal.Decimal `json:"repayLoanAmount"`
	RepayUnlockAmount decimal.Decimal `json:"repayUnlockAmount"`
	RepayTime         time.Time       `json:"repayTime"`
}

// RevisePledgeService -- POST /api/v3/loan/revise-pledge (UTA mgt. read & write)
//
// Adjusts the collateral on an open crypto loan. reviseType "IN" adds
// collateral, "OUT" withdraws it. The reply carries the resulting pledge rate.
type RevisePledgeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewRevisePledgeService(orderID, pledgeCoin string, amount decimal.Decimal) *RevisePledgeService {
	return &RevisePledgeService{c: c, body: map[string]any{
		"orderId":    orderID,
		"pledgeCoin": pledgeCoin,
		"amount":     amount.String(),
	}}
}

// SetReviseType sets the adjustment direction ("IN" adds collateral, "OUT"
// withdraws collateral).
func (s *RevisePledgeService) SetReviseType(reviseType string) *RevisePledgeService {
	s.body["reviseType"] = reviseType
	return s
}

func (s *RevisePledgeService) Do(ctx context.Context) (*RevisePledgeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/loan/revise-pledge", s.body).WithSign()
	return request.Do[RevisePledgeResult](req)
}

type RevisePledgeResult struct {
	LoanCoin        string          `json:"loanCoin"`
	PledgeCoin      string          `json:"pledgeCoin"`
	AfterPledgeRate decimal.Decimal `json:"afterPledgeRate"`
}

// GetPledgeRateHistoryService -- GET /api/v3/loan/pledge-rate-history (UTA mgt. read)
//
// Returns the unified account's collateral-adjustment records within a time
// window, paged. The window is bounded to the last 3 months.
type GetPledgeRateHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetPledgeRateHistoryService(startTime, endTime time.Time) *GetPledgeRateHistoryService {
	return &GetPledgeRateHistoryService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetPledgeRateHistoryService) SetOrderID(orderID string) *GetPledgeRateHistoryService {
	s.params["orderId"] = orderID
	return s
}

// SetReviseSide filters by adjustment direction ("down" transfer in/decrease,
// "up" transfer out/increase).
func (s *GetPledgeRateHistoryService) SetReviseSide(reviseSide string) *GetPledgeRateHistoryService {
	s.params["reviseSide"] = reviseSide
	return s
}

func (s *GetPledgeRateHistoryService) SetPledgeCoin(pledgeCoin string) *GetPledgeRateHistoryService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

func (s *GetPledgeRateHistoryService) SetPageNum(pageNum string) *GetPledgeRateHistoryService {
	s.params["pageNum"] = pageNum
	return s
}

func (s *GetPledgeRateHistoryService) SetPageSize(pageSize string) *GetPledgeRateHistoryService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetPledgeRateHistoryService) Do(ctx context.Context) ([]PledgeRateHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/pledge-rate-history", s.params).WithSign()
	resp, err := request.Do[[]PledgeRateHistory](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type PledgeRateHistory struct {
	LoanCoin         string          `json:"loanCoin"`
	PledgeCoin       string          `json:"pledgeCoin"`
	OrderID          string          `json:"orderId"`
	ReviseTime       time.Time       `json:"reviseTime"`
	ReviseSide       string          `json:"reviseSide"` // down, up
	ReviseAmount     decimal.Decimal `json:"reviseAmount"`
	AfterPledgeRate  decimal.Decimal `json:"afterPledgeRate"`
	BeforePledgeRate decimal.Decimal `json:"beforePledgeRate"`
}

// GetLoanDebtsService -- GET /api/v3/loan/debts (UTA mgt. read)
//
// Returns the unified account's current crypto-loan collateral and liability
// positions, each with the coin amount and its USDT-equivalent value.
type GetLoanDebtsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetLoanDebtsService() *GetLoanDebtsService {
	return &GetLoanDebtsService{c: c}
}

func (s *GetLoanDebtsService) Do(ctx context.Context) (*LoanDebts, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/debts").WithSign()
	return request.Do[LoanDebts](req)
}

type LoanDebts struct {
	PledgeInfos []LoanDebtInfo `json:"pledgeInfos"`
	LoanInfos   []LoanDebtInfo `json:"loanInfos"`
}

type LoanDebtInfo struct {
	Coin       string          `json:"coin"`
	Amount     decimal.Decimal `json:"amount"`
	AmountUsdt decimal.Decimal `json:"amountUsdt"`
}

// GetLoanReducesService -- GET /api/v3/loan/reduces (UTA mgt. read)
//
// Returns the unified account's crypto-loan liquidation (collateral reduction)
// records within a time window, paged. The window is bounded to the last 3
// months.
type GetLoanReducesService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetLoanReducesService(startTime, endTime time.Time) *GetLoanReducesService {
	return &GetLoanReducesService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetLoanReducesService) SetOrderID(orderID string) *GetLoanReducesService {
	s.params["orderId"] = orderID
	return s
}

func (s *GetLoanReducesService) SetLoanCoin(loanCoin string) *GetLoanReducesService {
	s.params["loanCoin"] = loanCoin
	return s
}

func (s *GetLoanReducesService) SetPledgeCoin(pledgeCoin string) *GetLoanReducesService {
	s.params["pledgeCoin"] = pledgeCoin
	return s
}

// SetStatus filters by liquidation status (COMPLETE, WAIT).
func (s *GetLoanReducesService) SetStatus(status string) *GetLoanReducesService {
	s.params["status"] = status
	return s
}

func (s *GetLoanReducesService) SetPageNum(pageNum string) *GetLoanReducesService {
	s.params["pageNum"] = pageNum
	return s
}

func (s *GetLoanReducesService) SetPageSize(pageSize string) *GetLoanReducesService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetLoanReducesService) Do(ctx context.Context) ([]LoanReduce, error) {
	req := request.Get(ctx, s.c, "/api/v3/loan/reduces", s.params).WithSign()
	resp, err := request.Do[[]LoanReduce](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type LoanReduce struct {
	OrderID         string          `json:"orderId"`
	LoanCoin        string          `json:"loanCoin"`
	PledgeCoin      string          `json:"pledgeCoin"`
	ReduceTime      time.Time       `json:"reduceTime"`
	PledgeRate      decimal.Decimal `json:"pledgeRate"`
	PledgePrice     decimal.Decimal `json:"pledgePrice"`
	Status          string          `json:"status"` // COMPLETE, WAIT
	PledgeAmount    decimal.Decimal `json:"pledgeAmount"`
	ReduceFee       decimal.Decimal `json:"reduceFee"`
	ResidueAmount   decimal.Decimal `json:"residueAmount"`
	RunlockAmount   decimal.Decimal `json:"runlockAmount"`
	RepayLoanAmount decimal.Decimal `json:"repayLoanAmount"`
}
