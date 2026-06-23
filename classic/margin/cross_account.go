package margin

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCrossAccountAssetsService -- GET /api/v2/margin/crossed/account/assets (margin read)
//
// Returns the cross-margin account's per-coin balances, borrowings and interest.
type GetCrossAccountAssetsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossAccountAssetsService() *GetCrossAccountAssetsService {
	return &GetCrossAccountAssetsService{c: c, params: map[string]string{}}
}

// SetCoin filters the assets to a single coin (e.g. USDT).
func (s *GetCrossAccountAssetsService) SetCoin(coin string) *GetCrossAccountAssetsService {
	s.params["coin"] = coin
	return s
}

func (s *GetCrossAccountAssetsService) Do(ctx context.Context) ([]CrossAccountAsset, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/account/assets", s.params).WithSign()
	resp, err := request.Do[[]CrossAccountAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossAccountAsset is one coin's cross-margin balance snapshot.
type CrossAccountAsset struct {
	Coin        string          `json:"coin"`        // token name
	TotalAmount decimal.Decimal `json:"totalAmount"` // total amount
	Available   decimal.Decimal `json:"available"`   // available amount
	Frozen      decimal.Decimal `json:"frozen"`      // assets frozen
	Borrow      decimal.Decimal `json:"borrow"`      // borrowed amount
	Interest    decimal.Decimal `json:"interest"`    // accrued interest
	Net         decimal.Decimal `json:"net"`         // net assets = available + frozen - borrow - interest
	Coupon      decimal.Decimal `json:"coupon"`      // trading bonus
	CTime       time.Time       `json:"cTime"`       // creation time
	UTime       time.Time       `json:"uTime"`       // update time
}

// CrossBorrowService -- POST /api/v2/margin/crossed/account/borrow (margin write)
//
// Borrows the given amount of a coin into the cross-margin account.
type CrossBorrowService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCrossBorrowService(coin string, borrowAmount decimal.Decimal) *CrossBorrowService {
	return &CrossBorrowService{c: c, body: map[string]any{
		"coin":         coin,
		"borrowAmount": borrowAmount.String(),
	}}
}

// SetClientID sets a customer-defined order ID for the borrow request.
func (s *CrossBorrowService) SetClientID(clientID string) *CrossBorrowService {
	s.body["clientid"] = clientID
	return s
}

func (s *CrossBorrowService) Do(ctx context.Context) (*CrossBorrowResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/account/borrow", s.body).WithSign()
	return request.Do[CrossBorrowResult](req)
}

// CrossBorrowResult is the result of a cross-margin borrow.
type CrossBorrowResult struct {
	LoanID       string          `json:"loanId"`       // loan order ID
	Coin         string          `json:"coin"`         // borrowing coin
	BorrowAmount decimal.Decimal `json:"borrowAmount"` // borrowing amount
}

// CrossRepayService -- POST /api/v2/margin/crossed/account/repay (margin write)
//
// Repays the given amount of a borrowed coin in the cross-margin account.
type CrossRepayService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCrossRepayService(coin string, repayAmount decimal.Decimal) *CrossRepayService {
	return &CrossRepayService{c: c, body: map[string]any{
		"coin":        coin,
		"repayAmount": repayAmount.String(),
	}}
}

func (s *CrossRepayService) Do(ctx context.Context) (*CrossRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/account/repay", s.body).WithSign()
	return request.Do[CrossRepayResult](req)
}

// CrossRepayResult is the result of a cross-margin repayment.
type CrossRepayResult struct {
	Coin             string          `json:"coin"`             // coin
	RepayID          string          `json:"repayId"`          // repay ID
	RemainDebtAmount decimal.Decimal `json:"remainDebtAmount"` // remaining borrowings
	RepayAmount      decimal.Decimal `json:"repayAmount"`      // repayment amount
}

// GetCrossRiskRateService -- GET /api/v2/margin/crossed/account/risk-rate (margin read)
//
// Returns the cross-margin account risk ratio (total assets / total liabilities).
type GetCrossRiskRateService struct {
	c *MarginClient
}

func (c *MarginClient) NewGetCrossRiskRateService() *GetCrossRiskRateService {
	return &GetCrossRiskRateService{c: c}
}

func (s *GetCrossRiskRateService) Do(ctx context.Context) (*CrossRiskRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/account/risk-rate").WithSign()
	return request.Do[CrossRiskRate](req)
}

// CrossRiskRate is the cross-margin account risk ratio.
type CrossRiskRate struct {
	RiskRateRatio decimal.Decimal `json:"riskRateRatio"` // total assets / total liabilities under cross mode
}

// GetCrossMaxBorrowableService -- GET /api/v2/margin/crossed/account/max-borrowable-amount (margin read)
//
// Returns the maximum amount of a coin that can be borrowed in the cross-margin account.
type GetCrossMaxBorrowableService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossMaxBorrowableService(coin string) *GetCrossMaxBorrowableService {
	return &GetCrossMaxBorrowableService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetCrossMaxBorrowableService) Do(ctx context.Context) (*CrossMaxBorrowable, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/account/max-borrowable-amount", s.params).WithSign()
	return request.Do[CrossMaxBorrowable](req)
}

// CrossMaxBorrowable is the maximum borrowable amount for a coin.
type CrossMaxBorrowable struct {
	Coin                string          `json:"coin"`                // coin identifier
	MaxBorrowableAmount decimal.Decimal `json:"maxBorrowableAmount"` // maximum borrow amount (changes in real time)
}

// GetCrossMaxTransferOutService -- GET /api/v2/margin/crossed/account/max-transfer-out-amount (margin read)
//
// Returns the maximum amount of a coin that can be transferred out of the cross-margin account.
type GetCrossMaxTransferOutService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossMaxTransferOutService(coin string) *GetCrossMaxTransferOutService {
	return &GetCrossMaxTransferOutService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetCrossMaxTransferOutService) Do(ctx context.Context) (*CrossMaxTransferOut, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/account/max-transfer-out-amount", s.params).WithSign()
	return request.Do[CrossMaxTransferOut](req)
}

// CrossMaxTransferOut is the maximum transferable-out amount for a coin.
type CrossMaxTransferOut struct {
	Coin                 string          `json:"coin"`                 // coin identifier
	MaxTransferOutAmount decimal.Decimal `json:"maxTransferOutAmount"` // maximum transferable amount
}

// GetCrossInterestRateAndLimitService -- GET /api/v2/margin/crossed/interest-rate-and-limit (margin read)
//
// Returns the per-coin cross-margin interest rate, leverage and borrowing limit,
// including the per-VIP-tier rate table.
type GetCrossInterestRateAndLimitService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossInterestRateAndLimitService(coin string) *GetCrossInterestRateAndLimitService {
	return &GetCrossInterestRateAndLimitService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetCrossInterestRateAndLimitService) Do(ctx context.Context) ([]CrossInterestRateAndLimit, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/interest-rate-and-limit", s.params).WithSign()
	resp, err := request.Do[[]CrossInterestRateAndLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossInterestRateAndLimit is one coin's cross-margin interest rate and limit.
type CrossInterestRateAndLimit struct {
	Coin                string             `json:"coin"`                // asset symbol
	Leverage            string             `json:"leverage"`            // leverage multiplier (default 3, tiers 3/5/10)
	Transferable        bool               `json:"transferable"`        // whether the coin is transferable
	Borrowable          bool               `json:"borrowable"`          // whether the coin is borrowable
	DailyInterestRate   decimal.Decimal    `json:"dailyInterestRate"`   // non-VIP daily rate
	AnnualInterestRate  decimal.Decimal    `json:"annualInterestRate"`  // non-VIP annual rate
	MaxBorrowableAmount decimal.Decimal    `json:"maxBorrowableAmount"` // maximum borrowable quantity
	VipList             []CrossVipRateItem `json:"vipList"`             // per-VIP-tier rate table
}

// CrossVipRateItem is one VIP tier's cross-margin interest rate and limit.
type CrossVipRateItem struct {
	Level              string          `json:"level"`              // VIP tier level
	Limit              decimal.Decimal `json:"limit"`              // VIP borrowing limit
	DailyInterestRate  decimal.Decimal `json:"dailyInterestRate"`  // VIP daily interest rate
	AnnualInterestRate decimal.Decimal `json:"annualInterestRate"` // VIP annual rate
	DiscountRate       decimal.Decimal `json:"discountRate"`       // rate multiplier (1 = no discount)
}

// GetCrossTierDataService -- GET /api/v2/margin/crossed/tier-data (margin read)
//
// Returns the cross-margin tier configuration (leverage, max borrow and
// maintenance margin rate per tier) for a coin.
type GetCrossTierDataService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetCrossTierDataService(coin string) *GetCrossTierDataService {
	return &GetCrossTierDataService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetCrossTierDataService) Do(ctx context.Context) ([]CrossTierData, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/crossed/tier-data", s.params).WithSign()
	resp, err := request.Do[[]CrossTierData](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossTierData is one cross-margin risk tier configuration entry.
type CrossTierData struct {
	Tier                string          `json:"tier"`                // tier
	Leverage            string          `json:"leverage"`            // effective leverage (global default 3x)
	Coin                string          `json:"coin"`                // coin
	MaxBorrowableAmount decimal.Decimal `json:"maxBorrowableAmount"` // maximum borrow
	MaintainMarginRate  decimal.Decimal `json:"maintainMarginRate"`  // maintenance margin rate
}

// CrossFlashRepayService -- POST /api/v2/margin/crossed/account/flash-repay (margin write)
//
// Triggers a flash repayment of the cross-margin account. If coin is omitted the
// account is fully repaid.
type CrossFlashRepayService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewCrossFlashRepayService() *CrossFlashRepayService {
	return &CrossFlashRepayService{c: c, body: map[string]any{}}
}

// SetCoin limits the flash repayment to a single coin (default: full repayment).
func (s *CrossFlashRepayService) SetCoin(coin string) *CrossFlashRepayService {
	s.body["coin"] = coin
	return s
}

func (s *CrossFlashRepayService) Do(ctx context.Context) (*CrossFlashRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/account/flash-repay", s.body).WithSign()
	return request.Do[CrossFlashRepayResult](req)
}

// CrossFlashRepayResult is the result of a cross-margin flash repayment.
type CrossFlashRepayResult struct {
	RepayID string `json:"repayId"` // repayment identifier
	Coin    string `json:"coin"`    // repayment coin (empty on full repayment with no residual)
}

// QueryCrossFlashRepayStatusService -- POST /api/v2/margin/crossed/account/query-flash-repay-status (margin read)
//
// Returns the result status of one or more cross-margin flash repayments by ID
// (max 100). Although read-only, the endpoint is a POST whose body is the ID list.
type QueryCrossFlashRepayStatusService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewQueryCrossFlashRepayStatusService(idList []string) *QueryCrossFlashRepayStatusService {
	return &QueryCrossFlashRepayStatusService{c: c, body: map[string]any{"idList": idList}}
}

func (s *QueryCrossFlashRepayStatusService) Do(ctx context.Context) ([]CrossFlashRepayStatus, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/crossed/account/query-flash-repay-status", s.body).WithSign()
	resp, err := request.Do[[]CrossFlashRepayStatus](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CrossFlashRepayStatus is one flash-repayment result status.
type CrossFlashRepayStatus struct {
	RepayID string `json:"repayId"` // repayment identifier
	Status  string `json:"status"`  // repayment result status
}
