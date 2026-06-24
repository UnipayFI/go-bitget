package margin

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetIsolatedAssetsService -- GET /api/v2/margin/isolated/account/assets (signed)
//
// Returns the per-symbol isolated-margin balances (one row per trading pair),
// including borrowed amount, accrued interest and net assets.
type GetIsolatedAssetsService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedAssetsService() *GetIsolatedAssetsService {
	return &GetIsolatedAssetsService{c: c, params: map[string]string{}}
}

func (s *GetIsolatedAssetsService) SetSymbol(symbol string) *GetIsolatedAssetsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetIsolatedAssetsService) Do(ctx context.Context) ([]IsolatedAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/account/assets", s.params).WithSign()
	resp, err := request.Do[[]IsolatedAssets](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedAssets is one trading pair's isolated-margin balance.
type IsolatedAssets struct {
	Coin        string          `json:"coin"`
	Symbol      string          `json:"symbol"`
	TotalAmount decimal.Decimal `json:"totalAmount"`
	Available   decimal.Decimal `json:"available"`
	Frozen      decimal.Decimal `json:"frozen"`
	Borrow      decimal.Decimal `json:"borrow"`
	Interest    decimal.Decimal `json:"interest"`
	Net         decimal.Decimal `json:"net"`
	Coupon      decimal.Decimal `json:"coupon"`
	CTime       time.Time       `json:"cTime"`
	UTime       time.Time       `json:"uTime"`
}

// IsolatedBorrowService -- POST /api/v2/margin/isolated/account/borrow (signed, state-changing)
//
// Borrows an asset against an isolated-margin trading pair.
type IsolatedBorrowService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewIsolatedBorrowService(symbol, coin string, borrowAmount decimal.Decimal) *IsolatedBorrowService {
	return &IsolatedBorrowService{c: c, body: map[string]any{
		"symbol":       symbol,
		"coin":         coin,
		"borrowAmount": borrowAmount.String(),
	}}
}

func (s *IsolatedBorrowService) SetClientOrderID(clientOid string) *IsolatedBorrowService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *IsolatedBorrowService) Do(ctx context.Context) (*IsolatedBorrowResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/account/borrow", s.body).WithSign()
	return request.Do[IsolatedBorrowResult](req)
}

// IsolatedBorrowResult is the outcome of an isolated borrow request.
type IsolatedBorrowResult struct {
	LoanID       string          `json:"loanId"`
	Symbol       string          `json:"symbol"`
	Coin         string          `json:"coin"`
	BorrowAmount decimal.Decimal `json:"borrowAmount"`
}

// IsolatedRepayService -- POST /api/v2/margin/isolated/account/repay (signed, state-changing)
//
// Repays borrowed funds for an isolated-margin trading pair.
type IsolatedRepayService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewIsolatedRepayService(symbol, coin string, repayAmount decimal.Decimal) *IsolatedRepayService {
	return &IsolatedRepayService{c: c, body: map[string]any{
		"symbol":      symbol,
		"coin":        coin,
		"repayAmount": repayAmount.String(),
	}}
}

func (s *IsolatedRepayService) SetClientOrderID(clientOid string) *IsolatedRepayService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *IsolatedRepayService) Do(ctx context.Context) (*IsolatedRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/account/repay", s.body).WithSign()
	return request.Do[IsolatedRepayResult](req)
}

// IsolatedRepayResult is the outcome of an isolated repay request.
type IsolatedRepayResult struct {
	RemainDebtAmount decimal.Decimal `json:"remainDebtAmount"`
	RepayID          string          `json:"repayId"`
	Symbol           string          `json:"symbol"`
	Coin             string          `json:"coin"`
	RepayAmount      decimal.Decimal `json:"repayAmount"`
}

// GetIsolatedRiskRateService -- GET /api/v2/margin/isolated/account/risk-rate (signed)
//
// Returns the per-symbol isolated-margin risk ratio (total assets / total
// liabilities under isolated mode).
type GetIsolatedRiskRateService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedRiskRateService() *GetIsolatedRiskRateService {
	return &GetIsolatedRiskRateService{c: c, params: map[string]string{}}
}

func (s *GetIsolatedRiskRateService) SetSymbol(symbol string) *GetIsolatedRiskRateService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetIsolatedRiskRateService) SetPageNum(pageNum int) *GetIsolatedRiskRateService {
	s.params["pageNum"] = strconv.Itoa(pageNum)
	return s
}

func (s *GetIsolatedRiskRateService) SetPageSize(pageSize int) *GetIsolatedRiskRateService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetIsolatedRiskRateService) Do(ctx context.Context) ([]IsolatedRiskRate, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/account/risk-rate", s.params).WithSign()
	resp, err := request.Do[[]IsolatedRiskRate](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedRiskRate is one trading pair's isolated-margin risk ratio.
type IsolatedRiskRate struct {
	Symbol        string          `json:"symbol"`
	RiskRateRatio decimal.Decimal `json:"riskRateRatio"`
}

// GetIsolatedInterestRateAndLimitService -- GET /api/v2/margin/isolated/interest-rate-and-limit (signed)
//
// Returns the per-symbol isolated-margin interest rates and maximum borrowable
// amounts for both the base and quote coins, including the VIP tier schedules.
type GetIsolatedInterestRateAndLimitService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedInterestRateAndLimitService(symbol string) *GetIsolatedInterestRateAndLimitService {
	return &GetIsolatedInterestRateAndLimitService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIsolatedInterestRateAndLimitService) Do(ctx context.Context) ([]IsolatedInterestRateAndLimit, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/interest-rate-and-limit", s.params).WithSign()
	resp, err := request.Do[[]IsolatedInterestRateAndLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedInterestRateAndLimit holds the interest-rate and borrow-limit
// configuration for one isolated-margin trading pair.
type IsolatedInterestRateAndLimit struct {
	Symbol   string `json:"symbol"`
	Leverage string `json:"leverage"`

	BaseCoin                  string                  `json:"baseCoin"`
	BaseTransferable          bool                    `json:"baseTransferable"`
	BaseBorrowable            bool                    `json:"baseBorrowable"`
	BaseDailyInterestRate     decimal.Decimal         `json:"baseDailyInterestRate"`
	BaseAnnuallyInterestRate  decimal.Decimal         `json:"baseAnnuallyInterestRate"`
	BaseMaxBorrowableAmount   decimal.Decimal         `json:"baseMaxBorrowableAmount"`
	BaseVIPList               []IsolatedMarginVIPItem `json:"baseVipList"`
	QuoteCoin                 string                  `json:"quoteCoin"`
	QuoteTransferable         bool                    `json:"quoteTransferable"`
	QuoteBorrowable           bool                    `json:"quoteBorrowable"`
	QuoteDailyInterestRate    decimal.Decimal         `json:"quoteDailyInterestRate"`
	QuoteAnnuallyInterestRate decimal.Decimal         `json:"quoteAnnuallyInterestRate"`
	QuoteMaxBorrowableAmount  decimal.Decimal         `json:"quoteMaxBorrowableAmount"`
	QuoteList                 []IsolatedMarginVIPItem `json:"quoteList"`
}

// IsolatedMarginVipItem is one VIP tier of the isolated-margin interest-rate
// schedule.
type IsolatedMarginVIPItem struct {
	Level                string          `json:"level"`
	DailyInterestRate    decimal.Decimal `json:"dailyInterestRate"`
	Limit                decimal.Decimal `json:"limit"`
	AnnuallyInterestRate decimal.Decimal `json:"annuallyInterestRate"`
	DiscountRate         decimal.Decimal `json:"discountRate"`
}

// GetIsolatedTierDataService -- GET /api/v2/margin/isolated/tier-data (signed)
//
// Returns the isolated-margin tier (leverage ladder) configuration for a
// trading pair.
type GetIsolatedTierDataService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedTierDataService(symbol string) *GetIsolatedTierDataService {
	return &GetIsolatedTierDataService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIsolatedTierDataService) Do(ctx context.Context) ([]IsolatedTierData, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/tier-data", s.params).WithSign()
	resp, err := request.Do[[]IsolatedTierData](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedTierData is one tier (leverage step) of an isolated-margin pair.
type IsolatedTierData struct {
	Tier                     string          `json:"tier"`
	Symbol                   string          `json:"symbol"`
	Leverage                 string          `json:"leverage"`
	BaseCoin                 string          `json:"baseCoin"`
	QuoteCoin                string          `json:"quoteCoin"`
	BaseMaxBorrowableAmount  decimal.Decimal `json:"baseMaxBorrowableAmount"`
	QuoteMaxBorrowableAmount decimal.Decimal `json:"quoteMaxBorrowableAmount"`
	MaintainMarginRate       decimal.Decimal `json:"maintainMarginRate"`
	InitRate                 decimal.Decimal `json:"initRate"`
}

// GetIsolatedMaxBorrowableService -- GET /api/v2/margin/isolated/account/max-borrowable-amount (signed)
//
// Returns the real-time maximum borrowable amounts for the base and quote coins
// of an isolated-margin trading pair.
type GetIsolatedMaxBorrowableService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedMaxBorrowableService(symbol string) *GetIsolatedMaxBorrowableService {
	return &GetIsolatedMaxBorrowableService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIsolatedMaxBorrowableService) Do(ctx context.Context) (*IsolatedMaxBorrowable, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/account/max-borrowable-amount", s.params).WithSign()
	return request.Do[IsolatedMaxBorrowable](req)
}

// IsolatedMaxBorrowable is the maximum borrowable amounts for one pair.
type IsolatedMaxBorrowable struct {
	Symbol                   string          `json:"symbol"`
	BaseCoin                 string          `json:"baseCoin"`
	BaseCoinMaxBorrowAmount  decimal.Decimal `json:"baseCoinMaxBorrowAmount"`
	QuoteCoin                string          `json:"quoteCoin"`
	QuoteCoinMaxBorrowAmount decimal.Decimal `json:"quoteCoinMaxBorrowAmount"`
}

// GetIsolatedMaxTransferOutService -- GET /api/v2/margin/isolated/account/max-transfer-out-amount (signed)
//
// Returns the maximum amounts that can be transferred out of an isolated-margin
// trading pair for its base and quote coins.
type GetIsolatedMaxTransferOutService struct {
	c      *MarginClient
	params map[string]string
}

func (c *MarginClient) NewGetIsolatedMaxTransferOutService(symbol string) *GetIsolatedMaxTransferOutService {
	return &GetIsolatedMaxTransferOutService{c: c, params: map[string]string{"symbol": symbol}}
}

func (s *GetIsolatedMaxTransferOutService) Do(ctx context.Context) (*IsolatedMaxTransferOut, error) {
	req := request.Get(ctx, s.c, "/api/v2/margin/isolated/account/max-transfer-out-amount", s.params).WithSign()
	return request.Do[IsolatedMaxTransferOut](req)
}

// IsolatedMaxTransferOut is the maximum transferable-out amounts for one pair.
type IsolatedMaxTransferOut struct {
	Symbol                        string          `json:"symbol"`
	BaseCoin                      string          `json:"baseCoin"`
	QuoteCoin                     string          `json:"quoteCoin"`
	BaseCoinMaxTransferOutAmount  decimal.Decimal `json:"baseCoinMaxTransferOutAmount"`
	QuoteCoinMaxTransferOutAmount decimal.Decimal `json:"quoteCoinMaxTransferOutAmount"`
}

// IsolatedFlashRepayService -- POST /api/v2/margin/isolated/account/flash-repay (signed, state-changing)
//
// Triggers a one-click ("flash") repayment for the given isolated-margin trading
// pairs; when symbolList is omitted every borrowed pair is repaid.
type IsolatedFlashRepayService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewIsolatedFlashRepayService() *IsolatedFlashRepayService {
	return &IsolatedFlashRepayService{c: c, body: map[string]any{}}
}

func (s *IsolatedFlashRepayService) SetSymbolList(symbolList []string) *IsolatedFlashRepayService {
	s.body["symbolList"] = symbolList
	return s
}

func (s *IsolatedFlashRepayService) Do(ctx context.Context) ([]IsolatedFlashRepayResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/account/flash-repay", s.body).WithSign()
	resp, err := request.Do[[]IsolatedFlashRepayResult](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedFlashRepayResult is the per-pair outcome of a flash repay request.
type IsolatedFlashRepayResult struct {
	RepayID string `json:"repayId"`
	Symbol  string `json:"symbol"`
	Result  string `json:"result"` // success / failure
}

// QueryIsolatedFlashRepayStatusService -- POST /api/v2/margin/isolated/account/query-flash-repay-status (signed)
//
// Queries the result status of previously-submitted isolated flash repayments by
// their repayment ids (up to 100 per request).
type QueryIsolatedFlashRepayStatusService struct {
	c    *MarginClient
	body map[string]any
}

func (c *MarginClient) NewQueryIsolatedFlashRepayStatusService(idList []string) *QueryIsolatedFlashRepayStatusService {
	return &QueryIsolatedFlashRepayStatusService{c: c, body: map[string]any{"idList": idList}}
}

func (s *QueryIsolatedFlashRepayStatusService) Do(ctx context.Context) ([]IsolatedFlashRepayStatus, error) {
	req := request.Post(ctx, s.c, "/api/v2/margin/isolated/account/query-flash-repay-status", s.body).WithSign()
	resp, err := request.Do[[]IsolatedFlashRepayStatus](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedFlashRepayStatus is one repayment's result status.
type IsolatedFlashRepayStatus struct {
	RepayID string `json:"repayId"`
	Status  string `json:"status"`
}
