package insloan

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetProductInfosService -- GET /api/v2/spot/ins-loan/product-infos (private)
//
// Returns the configuration of an institutional loan product: its leverage,
// supported contract types, and the risk-control lines for transfers, spot
// buys, position opening, and liquidation.
type GetProductInfosService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetProductInfosService(productId string) *GetProductInfosService {
	return &GetProductInfosService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetProductInfosService) Do(ctx context.Context) (*ProductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/product-infos", s.params).WithSign()
	return request.Do[ProductInfo](req)
}

// ProductInfo is the configuration of one institutional loan product. The
// support* / *ContractOpenLine fields are shared with the management-line
// product shape and are emitted when the product supports contract trading.
type ProductInfo struct {
	ProductID           string          `json:"productId"`
	Leverage            string          `json:"leverage"` // e.g. 2x/4x
	TransferLine        decimal.Decimal `json:"transferLine"`
	SpotBuyLine         decimal.Decimal `json:"spotBuyLine"`
	LiquidationLine     decimal.Decimal `json:"liquidationLine"`
	StopLiquidationLine decimal.Decimal `json:"stopLiquidationLine"`

	SupportUSDTContract  string          `json:"supportUsdtContract"` // YES, NO
	SupportCoinContract  string          `json:"supportCoinContract"` // YES, NO
	SupportUSDCContract  string          `json:"supportUsdcContract"` // YES, NO
	USDTContractOpenLine decimal.Decimal `json:"usdtContractOpenLine"`
	CoinContractOpenLine decimal.Decimal `json:"coinContractOpenLine"`
	USDCContractOpenLine decimal.Decimal `json:"usdcContractOpenLine"`
}

// GetEnsureCoinsConvertService -- GET /api/v2/spot/ins-loan/ensure-coins-convert (private)
//
// Returns the margin (collateral) coins of an institutional loan product, with
// their conversion ratios and the maximum convertible value per coin.
type GetEnsureCoinsConvertService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetEnsureCoinsConvertService(productId string) *GetEnsureCoinsConvertService {
	return &GetEnsureCoinsConvertService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetEnsureCoinsConvertService) Do(ctx context.Context) (*EnsureCoinsConvert, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/ensure-coins-convert", s.params).WithSign()
	return request.Do[EnsureCoinsConvert](req)
}

type EnsureCoinsConvert struct {
	ProductID string            `json:"productId"`
	CoinInfo  []EnsureCoinsInfo `json:"coinInfo"`
}

type EnsureCoinsInfo struct {
	Coin             string               `json:"coin"`
	ConvertRatio     decimal.Decimal      `json:"convertRatio"`
	MaxConvertValue  decimal.Decimal      `json:"maxConvertValue"` // max convert value (USDT)
	ConvertRatioList []ConvertRatioLadder `json:"convertRatioList"`
}

type ConvertRatioLadder struct {
	Ladder       string          `json:"ladder"`
	ConvertRatio decimal.Decimal `json:"convertRatio"`
}

// GetSymbolsService -- GET /api/v2/spot/ins-loan/symbols (private)
//
// Returns the tradable spot (and, when supported, contract) symbols with their
// leverage limits available to an institutional loan product.
type GetSymbolsService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetSymbolsService(productId string) *GetSymbolsService {
	return &GetSymbolsService{c: c, params: map[string]string{"productId": productId}}
}

func (s *GetSymbolsService) Do(ctx context.Context) (*Symbols, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/symbols", s.params).WithSign()
	return request.Do[Symbols](req)
}

type Symbols struct {
	ProductID            string         `json:"productId"`
	SpotSymbols          []string       `json:"spotSymbols"`
	MarginLeverage       string         `json:"marginLeverage"`
	USDTContractLeverage string         `json:"usdtContractLeverage"`
	CoinContractLeverage string         `json:"coinContractLeverage"`
	USDCContractLeverage string         `json:"usdcContractLeverage"`
	USDTContractSymbols  []ContractPair `json:"usdtContractSymbols"`
	CoinContractSymbols  []ContractPair `json:"coinContractSymbols"`
	USDCContractSymbols  []ContractPair `json:"usdcContractSymbols"`
}

type ContractPair struct {
	Symbol   string `json:"symbol"`
	Leverage string `json:"leverage"`
}

// GetLTVConvertService -- GET /api/v2/spot/ins-loan/ltv-convert (private)
//
// Returns the loan-to-value risk rate of a risk unit, with the breakdown of
// outstanding debt and asset holdings converted to USDT. riskUnitId is
// mandatory for parent accounts and optional for risk-sub-unit accounts.
type GetLTVConvertService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetLTVConvertService() *GetLTVConvertService {
	return &GetLTVConvertService{c: c, params: map[string]string{}}
}

// SetRiskUnitId sets the risk unit ID (required for parent account calls).
func (s *GetLTVConvertService) SetRiskUnitID(riskUnitId string) *GetLTVConvertService {
	s.params["riskUnitId"] = riskUnitId
	return s
}

func (s *GetLTVConvertService) Do(ctx context.Context) (*LTVConvert, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/ltv-convert", s.params).WithSign()
	return request.Do[LTVConvert](req)
}

type LTVConvert struct {
	Ltv              decimal.Decimal `json:"ltv"`
	SubAccountUids   []string        `json:"subAccountUids"`
	USDTBalance      decimal.Decimal `json:"usdtBalance"`
	UnpaidUSDTAmount decimal.Decimal `json:"unpaidUsdtAmount"`
	UnpaidInfo       []UnpaidInfo    `json:"unpaidInfo"`
	BalanceInfo      []BalanceInfo   `json:"balanceInfo"`
}

type UnpaidInfo struct {
	Coin           string          `json:"coin"`
	UnpaidQty      decimal.Decimal `json:"unpaidQty"`
	UnpaidInterest decimal.Decimal `json:"unpaidInterest"`
}

type BalanceInfo struct {
	Coin                string          `json:"coin"`
	Price               decimal.Decimal `json:"price"`
	Amount              decimal.Decimal `json:"amount"`
	ConvertedUSDTAmount decimal.Decimal `json:"convertedUsdtAmount"`
}

// GetTransferedService -- GET /api/v2/spot/ins-loan/transfered (private)
//
// Returns the quantity of a coin transferred into the institutional loan
// account, for the master account or a sub-account.
type GetTransferedService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetTransferedService(coin string) *GetTransferedService {
	return &GetTransferedService{c: c, params: map[string]string{"coin": coin}}
}

// SetUserId sets the user ID (master account or sub-account).
func (s *GetTransferedService) SetUserID(userId string) *GetTransferedService {
	s.params["userId"] = userId
	return s
}

func (s *GetTransferedService) Do(ctx context.Context) (*Transfered, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/transfered", s.params).WithSign()
	return request.Do[Transfered](req)
}

type Transfered struct {
	Coin       string          `json:"coin"`
	Transfered decimal.Decimal `json:"transfered"`
	UserID     string          `json:"userId"`
}

// GetRiskUnitService -- GET /api/v2/spot/ins-loan/risk-unit (private)
//
// Returns the risk unit IDs associated with the institutional loan account.
// Only the parent account API key can call this endpoint.
type GetRiskUnitService struct {
	c *InsLoanClient
}

func (c *InsLoanClient) NewGetRiskUnitService() *GetRiskUnitService {
	return &GetRiskUnitService{c: c}
}

func (s *GetRiskUnitService) Do(ctx context.Context) (*RiskUnit, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/risk-unit").WithSign()
	return request.Do[RiskUnit](req)
}

type RiskUnit struct {
	RiskUnitID []string `json:"riskUnitId"`
}

// BindUidOperate is the bind/unbind action for a risk unit UID binding.
type BindUIDOperate string

const (
	BindUIDOperateBind   BindUIDOperate = "bind"
	BindUIDOperateUnbind BindUIDOperate = "unbind"
)

// BindUidService -- POST /api/v2/spot/ins-loan/bind-uid (private, state-changing)
//
// Binds or unbinds a sub-account UID to a risk unit. riskUnitId is required for
// parent account calls only; max 50 UIDs per risk unit. The response data is the
// literal string "success" on completion.
type BindUIDService struct {
	c    *InsLoanClient
	body map[string]any
}

func (c *InsLoanClient) NewBindUIDService(uid string, operate BindUIDOperate) *BindUIDService {
	return &BindUIDService{c: c, body: map[string]any{
		"uid":     uid,
		"operate": string(operate),
	}}
}

// SetRiskUnitId sets the risk unit ID (required for parent account calls only).
func (s *BindUIDService) SetRiskUnitID(riskUnitId string) *BindUIDService {
	s.body["riskUnitId"] = riskUnitId
	return s
}

func (s *BindUIDService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/ins-loan/bind-uid", s.body).WithSign()
	return request.Do[string](req)
}

// GetLoanOrderService -- GET /api/v2/spot/ins-loan/loan-order (private)
//
// Returns the institutional loan orders. Omitting orderId returns all orders.
// startTime/endTime bound the lookback window (max 30-day span).
type GetLoanOrderService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetLoanOrderService() *GetLoanOrderService {
	return &GetLoanOrderService{c: c, params: map[string]string{}}
}

func (s *GetLoanOrderService) SetOrderID(orderId string) *GetLoanOrderService {
	s.params["orderId"] = orderId
	return s
}

func (s *GetLoanOrderService) SetStartTime(t time.Time) *GetLoanOrderService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLoanOrderService) SetEndTime(t time.Time) *GetLoanOrderService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetLoanOrderService) Do(ctx context.Context) ([]LoanOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/loan-order", s.params).WithSign()
	resp, err := request.Do[[]LoanOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type LoanOrder struct {
	OrderID        string          `json:"orderId"`
	OrderProductID string          `json:"orderProductId"`
	UID            string          `json:"uid"`
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

// GetRepaidHistoryService -- GET /api/v2/spot/ins-loan/repaid-history (private)
//
// Returns the institutional loan repayment orders within the given time window.
type GetRepaidHistoryService struct {
	c      *InsLoanClient
	params map[string]string
}

func (c *InsLoanClient) NewGetRepaidHistoryService() *GetRepaidHistoryService {
	return &GetRepaidHistoryService{c: c, params: map[string]string{}}
}

func (s *GetRepaidHistoryService) SetStartTime(t time.Time) *GetRepaidHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetRepaidHistoryService) SetEndTime(t time.Time) *GetRepaidHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetRepaidHistoryService) SetLimit(limit string) *GetRepaidHistoryService {
	s.params["limit"] = limit
	return s
}

func (s *GetRepaidHistoryService) Do(ctx context.Context) ([]RepaidOrder, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/ins-loan/repaid-history", s.params).WithSign()
	resp, err := request.Do[[]RepaidOrder](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// RepaidOrder is one institutional loan repayment order. Bitget has emitted the
// principal/interest under two key spellings across versions of this endpoint
// (repaidAmount/repaidInterest and repayAmount/repayInterest), so both variants
// are captured to cover whichever shape the live API returns.
type RepaidOrder struct {
	RepayOrderID   string          `json:"repayOrderId"`
	BusinessType   string          `json:"businessType"` // normal, liquidation
	RepayType      string          `json:"repayType"`    // all, part
	RepaidTime     time.Time       `json:"repaidTime"`
	Coin           string          `json:"coin"`
	RepaidAmount   decimal.Decimal `json:"repaidAmount"`
	RepaidInterest decimal.Decimal `json:"repaidInterest"`
	RepayAmount    decimal.Decimal `json:"repayAmount"`
	RepayInterest  decimal.Decimal `json:"repayInterest"`
}
