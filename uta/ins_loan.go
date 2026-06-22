package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetInsLoanTransferedService -- GET /api/v3/ins-loan/transfered (UTA mgt. read)
//
// Returns the quantity of a coin transferred into the institutional loan
// account, for the master account or a sub-account.
type GetInsLoanTransferedService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanTransferedService(coin string) *GetInsLoanTransferedService {
	return &GetInsLoanTransferedService{c: c, params: map[string]string{"coin": coin}}
}

// SetUserId sets the user ID (master account or sub-account).
func (s *GetInsLoanTransferedService) SetUserId(userId string) *GetInsLoanTransferedService {
	s.params["userId"] = userId
	return s
}

func (s *GetInsLoanTransferedService) Do(ctx context.Context) (*InsLoanTransfered, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/transfered", s.params).WithSign()
	return request.Do[InsLoanTransfered](req)
}

type InsLoanTransfered struct {
	Coin       string          `json:"coin"`
	Transfered decimal.Decimal `json:"transfered"`
	UserId     string          `json:"userId"`
}

// GetInsLoanSymbolsService -- GET /api/v3/ins-loan/symbols (UTA mgt. read)
//
// Returns the tradable spot and futures symbols (with their leverage limits)
// available to an institutional loan product.
type GetInsLoanSymbolsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanSymbolsService(productId string) *GetInsLoanSymbolsService {
	return &GetInsLoanSymbolsService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetInsLoanSymbolsService) Do(ctx context.Context) (*InsLoanSymbols, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/symbols", s.params).WithSign()
	return request.Do[InsLoanSymbols](req)
}

type InsLoanSymbols struct {
	ProductId            string                `json:"productId"`
	SpotSymbols          []string              `json:"spotSymbols"`
	MarginLeverage       string                `json:"marginLeverage"`
	UsdtContractLeverage string                `json:"usdtContractLeverage"`
	CoinContractLeverage string                `json:"coinContractLeverage"`
	UsdcContractLeverage string                `json:"usdcContractLeverage"`
	UsdtContractSymbols  []InsLoanContractPair `json:"usdtContractSymbols"`
	CoinContractSymbols  []InsLoanContractPair `json:"coinContractSymbols"`
	UsdcContractSymbols  []InsLoanContractPair `json:"usdcContractSymbols"`
}

type InsLoanContractPair struct {
	Symbol   string `json:"symbol"`
	Leverage string `json:"leverage"`
}

// GetInsLoanRiskUnitService -- GET /api/v3/ins-loan/risk-unit (UTA mgt. read)
//
// Returns the risk unit IDs associated with the institutional loan account.
// Only the parent account API key can call this endpoint.
type GetInsLoanRiskUnitService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetInsLoanRiskUnitService() *GetInsLoanRiskUnitService {
	return &GetInsLoanRiskUnitService{c: c}
}

func (s *GetInsLoanRiskUnitService) Do(ctx context.Context) (*InsLoanRiskUnit, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/risk-unit").WithSign()
	return request.Do[InsLoanRiskUnit](req)
}

type InsLoanRiskUnit struct {
	RiskUnitId []string `json:"riskUnitId"`
}

// GetInsLoanRepaidHistoryService -- GET /api/v3/ins-loan/repaid-history (UTA mgt. read)
//
// Returns the institutional loan repayment orders, bounded to a two-year
// lookback window and limited to successfully completed repayments.
type GetInsLoanRepaidHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanRepaidHistoryService() *GetInsLoanRepaidHistoryService {
	return &GetInsLoanRepaidHistoryService{c: c, params: map[string]string{}}
}

func (s *GetInsLoanRepaidHistoryService) SetStartTime(t time.Time) *GetInsLoanRepaidHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInsLoanRepaidHistoryService) SetEndTime(t time.Time) *GetInsLoanRepaidHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetInsLoanRepaidHistoryService) SetLimit(limit string) *GetInsLoanRepaidHistoryService {
	s.params["limit"] = limit
	return s
}

func (s *GetInsLoanRepaidHistoryService) Do(ctx context.Context) ([]InsLoanRepaidOrder, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/repaid-history", s.params).WithSign()
	resp, err := request.Do[[]InsLoanRepaidOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type InsLoanRepaidOrder struct {
	RepayOrderId  string          `json:"repayOrderId"`
	BusinessType  string          `json:"businessType"` // normal, liquidation
	RepayType     string          `json:"repayType"`    // all, part
	RepaidTime    time.Time       `json:"repaidTime"`
	Coin          string          `json:"coin"`
	RepayAmount   decimal.Decimal `json:"repayAmount"`
	RepayInterest decimal.Decimal `json:"repayInterest"`
}

// GetInsLoanProductInfosService -- GET /api/v3/ins-loan/product-infos (UTA mgt. read)
//
// Returns the configuration of an institutional loan product: its leverage,
// supported contract types, and the risk-control lines for transfers, buys,
// position opening and liquidation.
type GetInsLoanProductInfosService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanProductInfosService(productId string) *GetInsLoanProductInfosService {
	return &GetInsLoanProductInfosService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetInsLoanProductInfosService) Do(ctx context.Context) (*InsLoanProductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/product-infos", s.params).WithSign()
	return request.Do[InsLoanProductInfo](req)
}

type InsLoanProductInfo struct {
	ProductId            string          `json:"productId"`
	Leverage             string          `json:"leverage"`            // e.g. 2x/4x
	SupportUsdtContract  string          `json:"supportUsdtContract"` // YES, NO
	SupportCoinContract  string          `json:"supportCoinContract"` // YES, NO
	SupportUsdcContract  string          `json:"supportUsdcContract"` // YES, NO
	TransferLine         decimal.Decimal `json:"transferLine"`
	SpotBuyLine          decimal.Decimal `json:"spotBuyLine"`
	UsdtContractOpenLine decimal.Decimal `json:"usdtContractOpenLine"`
	CoinContractOpenLine decimal.Decimal `json:"coinContractOpenLine"`
	UsdcContractOpenLine decimal.Decimal `json:"usdcContractOpenLine"`
	LiquidationLine      decimal.Decimal `json:"liquidationLine"`
	StopLiquidationLine  decimal.Decimal `json:"stopLiquidationLine"`
}

// GetInsLoanOrderService -- GET /api/v3/ins-loan/loan-order (UTA mgt. read)
//
// Returns the institutional loan orders. Omitting orderId returns all orders
// sorted by loanTime descending.
type GetInsLoanOrderService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanOrderService() *GetInsLoanOrderService {
	return &GetInsLoanOrderService{c: c, params: map[string]string{}}
}

func (s *GetInsLoanOrderService) SetOrderId(orderId string) *GetInsLoanOrderService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetInsLoanOrderService) SetStartTime(t time.Time) *GetInsLoanOrderService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInsLoanOrderService) SetEndTime(t time.Time) *GetInsLoanOrderService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInsLoanOrderService) Do(ctx context.Context) ([]InsLoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/loan-order", s.params).WithSign()
	resp, err := request.Do[[]InsLoanOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type InsLoanOrder struct {
	OrderId        string          `json:"orderId"`
	OrderProductId string          `json:"orderProductId"`
	Uid            string          `json:"uid"`
	LoanTime       time.Time       `json:"loanTime"`
	LoanCoin       string          `json:"loanCoin"`
	LoanAmount     decimal.Decimal `json:"loanAmount"`
	UnpaidAmount   decimal.Decimal `json:"unpaidAmount"`
	UnpaidInterest decimal.Decimal `json:"unpaidInterest"`
	RepaidAmount   decimal.Decimal `json:"repaidAmount"`
	RepaidInterest decimal.Decimal `json:"repaidInterest"`
	Reserve        decimal.Decimal `json:"reserve"`
	Status         string          `json:"status"` // not_paid_off, paid_off
}

// GetInsLoanLTVConvertService -- GET /api/v3/ins-loan/ltv-convert (UTA mgt. read)
//
// Returns the loan-to-value risk rate of a risk unit, with the breakdown of
// outstanding debt and asset holdings converted to USDT.
type GetInsLoanLTVConvertService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanLTVConvertService() *GetInsLoanLTVConvertService {
	return &GetInsLoanLTVConvertService{c: c, params: map[string]string{}}
}

// SetRiskUnitId sets the risk unit ID (required for parent account calls).
func (s *GetInsLoanLTVConvertService) SetRiskUnitId(riskUnitId string) *GetInsLoanLTVConvertService {
	s.params["riskUnitId"] = riskUnitId
	return s
}

func (s *GetInsLoanLTVConvertService) Do(ctx context.Context) (*InsLoanLTVConvert, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/ltv-convert", s.params).WithSign()
	return request.Do[InsLoanLTVConvert](req)
}

type InsLoanLTVConvert struct {
	Ltv              decimal.Decimal      `json:"ltv"`
	SubAccountUids   []string             `json:"subAccountUids"`
	UsdtBalance      decimal.Decimal      `json:"usdtBalance"`
	UnpaidUsdtAmount decimal.Decimal      `json:"unpaidUsdtAmount"`
	UnpaidInfo       []InsLoanUnpaidInfo  `json:"unpaidInfo"`
	BalanceInfo      []InsLoanBalanceInfo `json:"balanceInfo"`
}

type InsLoanUnpaidInfo struct {
	Coin           string          `json:"coin"`
	UnpaidQty      decimal.Decimal `json:"unpaidQty"`
	UnpaidInterest decimal.Decimal `json:"unpaidInterest"`
}

type InsLoanBalanceInfo struct {
	Coin                string          `json:"coin"`
	Price               decimal.Decimal `json:"price"`
	Amount              decimal.Decimal `json:"amount"`
	ConvertedUsdtAmount decimal.Decimal `json:"convertedUsdtAmount"`
}

// GetInsLoanEnsureCoinsConvertService -- GET /api/v3/ins-loan/ensure-coins-convert (UTA mgt. read)
//
// Returns the margin (collateral) coins of an institutional loan product, with
// their conversion ratios and the tiered ladder of ratios by value bracket.
type GetInsLoanEnsureCoinsConvertService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInsLoanEnsureCoinsConvertService(productId string) *GetInsLoanEnsureCoinsConvertService {
	return &GetInsLoanEnsureCoinsConvertService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetInsLoanEnsureCoinsConvertService) Do(ctx context.Context) (*InsLoanEnsureCoins, error) {
	req := request.Get(ctx, s.c, "/api/v3/ins-loan/ensure-coins-convert", s.params).WithSign()
	return request.Do[InsLoanEnsureCoins](req)
}

type InsLoanEnsureCoins struct {
	ProductId string                  `json:"productId"`
	CoinInfo  []InsLoanEnsureCoinInfo `json:"coinInfo"`
}

type InsLoanEnsureCoinInfo struct {
	Coin             string                      `json:"coin"`
	ConvertRatio     decimal.Decimal             `json:"convertRatio"`
	MaxConvertValue  decimal.Decimal             `json:"maxConvertValue"`
	ConvertRatioList []InsLoanConvertRatioLadder `json:"convertRatioList"`
}

type InsLoanConvertRatioLadder struct {
	Ladder       string          `json:"ladder"`
	ConvertRatio decimal.Decimal `json:"convertRatio"`
}

// BindInsLoanUidService -- POST /api/v3/ins-loan/bind-uid (UTA mgt. read & write)
//
// Binds or unbinds a sub-account UID to a risk unit (operate is "bind" or
// "unbind"). Advanced account mode is required; max 50 UIDs per risk unit.
type BindInsLoanUidService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewBindInsLoanUidService(uid, operate string) *BindInsLoanUidService {
	return &BindInsLoanUidService{c: c, body: map[string]any{
		"uid":     uid,
		"operate": operate,
	}}
}

// SetRiskUnitId sets the risk unit ID (required for parent account calls only).
func (s *BindInsLoanUidService) SetRiskUnitId(riskUnitId string) *BindInsLoanUidService {
	s.body["riskUnitId"] = riskUnitId
	return s
}

func (s *BindInsLoanUidService) Do(ctx context.Context) (*InsLoanBindResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/ins-loan/bind-uid", s.body).WithSign()
	return request.Do[InsLoanBindResult](req)
}

type InsLoanBindResult struct {
	RiskUnitId string `json:"riskUnitId"`
	Uid        string `json:"uid"`
	Operate    string `json:"operate"` // bind, unbind
}
