package affiliate

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CommissionBizType is the trade type a commission record belongs to
// (agent-commission detail).
type CommissionBizType string

const (
	CommissionBizTypeSpot    CommissionBizType = "spot"
	CommissionBizTypeFutures CommissionBizType = "futures"
)

// CommissionSubBizType is the product-line breakdown of a commission record
// (agent-commission detail).
type CommissionSubBizType string

const (
	CommissionSubBizTypeSpot        CommissionSubBizType = "spot"
	CommissionSubBizTypeMargin      CommissionSubBizType = "margin"
	CommissionSubBizTypeUSDTFutures CommissionSubBizType = "usdt_futures"
	CommissionSubBizTypeCoinFutures CommissionSubBizType = "coin_futures"
	CommissionSubBizTypeUSDCFutures CommissionSubBizType = "usdc_futures"
)

// CommissionTraderType distinguishes regular users from copy-trade traders.
type CommissionTraderType string

const (
	CommissionTraderTypeUser   CommissionTraderType = "user"
	CommissionTraderTypeTrader CommissionTraderType = "trader"
)

// CommissionApiType reports whether the customer trades via API.
type CommissionApiType string

const (
	CommissionApiTypeAPI    CommissionApiType = "api"
	CommissionApiTypeNonAPI CommissionApiType = "non_api"
)

// CommissionStatus is the settlement state of a commission record.
type CommissionStatus string

const (
	CommissionStatusSettled   CommissionStatus = "settled"
	CommissionStatusUnsettled CommissionStatus = "unsettled"
	CommissionStatusNotIssued CommissionStatus = "notIssued"
)

// KycResult is a customer's KYC verification outcome.
type KycResult string

const (
	KycResultPassed    KycResult = "passed"
	KycResultNotPassed KycResult = "not_passed"
)

// GetCustomerCommissionsService -- GET /api/v2/broker/customer-commissions (affiliate/agent)
//
// Returns the agent's direct commission records for referred customers, paged by
// idLessThan over a window of up to 30 days.
type GetCustomerCommissionsService struct {
	c      *AffiliateClient
	params map[string]string
}

func (c *AffiliateClient) NewGetCustomerCommissionsService() *GetCustomerCommissionsService {
	return &GetCustomerCommissionsService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t (max 30-day range).
func (s *GetCustomerCommissionsService) SetStartTime(t time.Time) *GetCustomerCommissionsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range).
func (s *GetCustomerCommissionsService) SetEndTime(t time.Time) *GetCustomerCommissionsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backward, returning data with an id older than this value.
func (s *GetCustomerCommissionsService) SetIDLessThan(id string) *GetCustomerCommissionsService {
	s.params["idLessThan"] = id
	return s
}

// SetLimit caps the number of records returned (default 100, max 1000).
func (s *GetCustomerCommissionsService) SetLimit(limit int) *GetCustomerCommissionsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerCommissionsService) SetUID(uid string) *GetCustomerCommissionsService {
	s.params["uid"] = uid
	return s
}

// SetCoin filters to a single coin (e.g. BTC).
func (s *GetCustomerCommissionsService) SetCoin(coin string) *GetCustomerCommissionsService {
	s.params["coin"] = coin
	return s
}

// SetSymbol filters to a single trading symbol (e.g. BGBUSDT_SPBL).
func (s *GetCustomerCommissionsService) SetSymbol(symbol string) *GetCustomerCommissionsService {
	s.params["symbol"] = symbol
	return s
}

// SetShowSub toggles inclusion of subordinate user info ("yes" or "no").
func (s *GetCustomerCommissionsService) SetShowSub(showSub string) *GetCustomerCommissionsService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetCustomerCommissionsService) Do(ctx context.Context) (*CustomerCommissions, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/customer-commissions", s.params).WithSign()
	return request.Do[CustomerCommissions](req)
}

// CustomerCommissions is the paged direct-commission response.
type CustomerCommissions struct {
	EndID          string                     `json:"endId"`
	CommissionList []CustomerCommissionRecord `json:"commissionList"`
}

// CustomerCommissionRecord is one direct-commission row for a referred customer.
type CustomerCommissionRecord struct {
	UID                    string          `json:"uid"`
	Date                   time.Time       `json:"date"` // commission date, UTC+8 (ms)
	Coin                   string          `json:"coin"`
	Symbol                 string          `json:"symbol"`
	ProductType            string          `json:"productType"` // SPOT, MARGIN, USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
	DealAmount             decimal.Decimal `json:"dealAmount"`
	Fee                    decimal.Decimal `json:"fee"`
	FeeDeduction           decimal.Decimal `json:"feeDeduction"`
	ActivityBonusDeduct    decimal.Decimal `json:"activityBonusDeduct"`
	SpotCouponDeduct       decimal.Decimal `json:"spotCouponDeduct"`
	FuturesCouponDeduct    decimal.Decimal `json:"futuresCouponDeduct"`
	SpotFeeDiscountDeduct  decimal.Decimal `json:"spotFeeDiscountDeduct"`
	NegativeMakerFeeDeduct decimal.Decimal `json:"negativeMakerFeeDeduct"`
	FeePaid                decimal.Decimal `json:"feePaid"`
	RebateAmount           decimal.Decimal `json:"rebateAmount"`
	UserTotalRebateAmount  decimal.Decimal `json:"userTotalRebateAmount"`
	DayTotalRebateAmount   decimal.Decimal `json:"dayTotalRebateAmount"`
	TotalRebateAmount      decimal.Decimal `json:"totalRebateAmount"`
}

// GetCustomerTradeVolumeService -- POST /api/v2/broker/customer-trade-volume (affiliate/agent)
//
// Returns the daily trade volume of the agent's referred customers.
type GetCustomerTradeVolumeService struct {
	c    *AffiliateClient
	body map[string]any
}

func (c *AffiliateClient) NewGetCustomerTradeVolumeService() *GetCustomerTradeVolumeService {
	return &GetCustomerTradeVolumeService{c: c, body: map[string]any{}}
}

// SetStartTime filters records at or after t (ms).
func (s *GetCustomerTradeVolumeService) SetStartTime(t time.Time) *GetCustomerTradeVolumeService {
	s.body["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (ms).
func (s *GetCustomerTradeVolumeService) SetEndTime(t time.Time) *GetCustomerTradeVolumeService {
	s.body["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageNo selects the page number.
func (s *GetCustomerTradeVolumeService) SetPageNo(pageNo int) *GetCustomerTradeVolumeService {
	s.body["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetCustomerTradeVolumeService) SetPageSize(pageSize int) *GetCustomerTradeVolumeService {
	s.body["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerTradeVolumeService) SetUID(uid string) *GetCustomerTradeVolumeService {
	s.body["uid"] = uid
	return s
}

func (s *GetCustomerTradeVolumeService) Do(ctx context.Context) ([]CustomerTradeVolume, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/customer-trade-volume", s.body).WithSign()
	resp, err := request.Do[[]CustomerTradeVolume](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CustomerTradeVolume is one customer's trade volume on a given day.
type CustomerTradeVolume struct {
	UID    string          `json:"uid"`
	Volumn decimal.Decimal `json:"volumn"` // trade volume (spelled "volumn" on the wire)
	Time   time.Time       `json:"time"`   // timestamp (ms)
}

// GetCustomerListService -- POST /api/v2/broker/customer-list (affiliate/agent)
//
// Returns the agent's referred customer list with registration times.
type GetCustomerListService struct {
	c    *AffiliateClient
	body map[string]any
}

func (c *AffiliateClient) NewGetCustomerListService() *GetCustomerListService {
	return &GetCustomerListService{c: c, body: map[string]any{}}
}

// SetStartTime filters records at or after t (ms).
func (s *GetCustomerListService) SetStartTime(t time.Time) *GetCustomerListService {
	s.body["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (ms).
func (s *GetCustomerListService) SetEndTime(t time.Time) *GetCustomerListService {
	s.body["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageNo selects the page number.
func (s *GetCustomerListService) SetPageNo(pageNo int) *GetCustomerListService {
	s.body["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetCustomerListService) SetPageSize(pageSize int) *GetCustomerListService {
	s.body["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerListService) SetUID(uid string) *GetCustomerListService {
	s.body["uid"] = uid
	return s
}

// SetReferralCode filters to customers referred under a specific referral code.
func (s *GetCustomerListService) SetReferralCode(code string) *GetCustomerListService {
	s.body["referralCode"] = code
	return s
}

// SetShowSub toggles inclusion of subordinate user info ("yes" or "no").
func (s *GetCustomerListService) SetShowSub(showSub string) *GetCustomerListService {
	s.body["showSub"] = showSub
	return s
}

func (s *GetCustomerListService) Do(ctx context.Context) ([]CustomerListEntry, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/customer-list", s.body).WithSign()
	resp, err := request.Do[[]CustomerListEntry](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CustomerListEntry is one referred customer. Bitget returns uid here as a bare
// JSON number (not a quoted string), so it is typed int64.
type CustomerListEntry struct {
	UID          int64     `json:"uid"`
	RegisterTime time.Time `json:"registerTime"` // registration time (ms)
}

// GetCustomerKycResultService -- GET /api/v2/broker/customer-kyc-result (affiliate/agent)
//
// Returns the KYC verification result for the agent's referred customers, paged
// by idLessThan over a window of up to 90 days.
type GetCustomerKycResultService struct {
	c      *AffiliateClient
	params map[string]string
}

func (c *AffiliateClient) NewGetCustomerKycResultService() *GetCustomerKycResultService {
	return &GetCustomerKycResultService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t (max 90-day range).
func (s *GetCustomerKycResultService) SetStartTime(t time.Time) *GetCustomerKycResultService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 90-day range).
func (s *GetCustomerKycResultService) SetEndTime(t time.Time) *GetCustomerKycResultService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan pages backward, returning data with an id older than this value.
func (s *GetCustomerKycResultService) SetIDLessThan(id string) *GetCustomerKycResultService {
	s.params["idLessThan"] = id
	return s
}

// SetLimit caps the number of records returned (default 100, max 1000).
func (s *GetCustomerKycResultService) SetLimit(limit int) *GetCustomerKycResultService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerKycResultService) SetUID(uid string) *GetCustomerKycResultService {
	s.params["uid"] = uid
	return s
}

// SetShowSub toggles inclusion of subordinate user info ("yes" or "no").
func (s *GetCustomerKycResultService) SetShowSub(showSub string) *GetCustomerKycResultService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetCustomerKycResultService) Do(ctx context.Context) (*CustomerKycResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/customer-kyc-result", s.params).WithSign()
	return request.Do[CustomerKycResult](req)
}

// CustomerKycResult is the paged KYC-result response.
type CustomerKycResult struct {
	UserList []CustomerKycRecord `json:"userList"`
	EndID    string              `json:"endId"`
}

// CustomerKycRecord is one customer's KYC outcome.
type CustomerKycRecord struct {
	UID       string    `json:"uid"`
	KycResult KycResult `json:"kycResult"` // passed, not_passed
}

// GetCustomerDepositService -- POST /api/v2/broker/customer-deposit (affiliate/agent)
//
// Returns deposit records of the agent's referred customers.
type GetCustomerDepositService struct {
	c    *AffiliateClient
	body map[string]any
}

func (c *AffiliateClient) NewGetCustomerDepositService() *GetCustomerDepositService {
	return &GetCustomerDepositService{c: c, body: map[string]any{}}
}

// SetStartTime filters records at or after t (ms).
func (s *GetCustomerDepositService) SetStartTime(t time.Time) *GetCustomerDepositService {
	s.body["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (ms).
func (s *GetCustomerDepositService) SetEndTime(t time.Time) *GetCustomerDepositService {
	s.body["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageNo selects the page number.
func (s *GetCustomerDepositService) SetPageNo(pageNo int) *GetCustomerDepositService {
	s.body["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetCustomerDepositService) SetPageSize(pageSize int) *GetCustomerDepositService {
	s.body["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerDepositService) SetUID(uid string) *GetCustomerDepositService {
	s.body["uid"] = uid
	return s
}

func (s *GetCustomerDepositService) Do(ctx context.Context) ([]CustomerDeposit, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/customer-deposit", s.body).WithSign()
	resp, err := request.Do[[]CustomerDeposit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CustomerDeposit is one deposit record of a referred customer.
type CustomerDeposit struct {
	OrderID       string          `json:"orderId"`
	UID           string          `json:"uid"`
	DepositTime   time.Time       `json:"depositTime"` // deposit time (ms)
	DepositCoin   string          `json:"depositCoin"`
	DepositAmount decimal.Decimal `json:"depositAmount"`
}

// GetCustomerAssetService -- POST /api/v2/broker/customer-asset (affiliate/agent)
//
// Returns account-balance snapshots of the agent's referred customers (refreshed
// roughly every 10 minutes).
type GetCustomerAssetService struct {
	c    *AffiliateClient
	body map[string]any
}

func (c *AffiliateClient) NewGetCustomerAssetService() *GetCustomerAssetService {
	return &GetCustomerAssetService{c: c, body: map[string]any{}}
}

// SetPageNo selects the page number.
func (s *GetCustomerAssetService) SetPageNo(pageNo int) *GetCustomerAssetService {
	s.body["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetPageSize sets the page size (default 100, max 500).
func (s *GetCustomerAssetService) SetPageSize(pageSize int) *GetCustomerAssetService {
	s.body["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetUID filters to a single referred customer UID.
func (s *GetCustomerAssetService) SetUID(uid string) *GetCustomerAssetService {
	s.body["uid"] = uid
	return s
}

func (s *GetCustomerAssetService) Do(ctx context.Context) ([]CustomerAsset, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/customer-asset", s.body).WithSign()
	resp, err := request.Do[[]CustomerAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CustomerAsset is one referred customer's account-balance snapshot.
type CustomerAsset struct {
	UID     string          `json:"uid"`
	Balance decimal.Decimal `json:"balance"`
	UTime   time.Time       `json:"uTime"`  // last update time (ms)
	Remark  string          `json:"remark"` // e.g. "sub account exceed 5"
}

// GetAgentCommissionService -- GET /api/v2/broker/agent-commission (affiliate/agent)
//
// Returns the agent's own commission detail records, paged by idLessThan.
type GetAgentCommissionService struct {
	c      *AffiliateClient
	params map[string]string
}

func (c *AffiliateClient) NewGetAgentCommissionService() *GetAgentCommissionService {
	return &GetAgentCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t (ms).
func (s *GetAgentCommissionService) SetStartTime(t time.Time) *GetAgentCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (ms).
func (s *GetAgentCommissionService) SetEndTime(t time.Time) *GetAgentCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *GetAgentCommissionService) SetLimit(limit int) *GetAgentCommissionService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backward, returning data with an id older than this value.
func (s *GetAgentCommissionService) SetIDLessThan(id string) *GetAgentCommissionService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetAgentCommissionService) Do(ctx context.Context) (*AgentCommission, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/agent-commission", s.params).WithSign()
	return request.Do[AgentCommission](req)
}

// AgentCommission is the paged agent-commission-detail response.
type AgentCommission struct {
	EndID          string                  `json:"endId"`
	CommissionList []AgentCommissionRecord `json:"commissionList"`
}

// AgentCommissionRecord is one agent-commission-detail row.
type AgentCommissionRecord struct {
	UID                     string               `json:"uid"`
	BizType                 CommissionBizType    `json:"bizType"`    // spot, futures
	SubBizType              CommissionSubBizType `json:"subBizType"` // spot, margin, usdt_futures, coin_futures, usdc_futures
	Symbol                  string               `json:"symbol"`
	Coin                    string               `json:"coin"`
	Fee                     decimal.Decimal      `json:"fee"`
	Volume                  decimal.Decimal      `json:"volume"`
	ActivityBonusDeduct     decimal.Decimal      `json:"activityBonusDeduct"`
	SpotCouponDeduct        decimal.Decimal      `json:"spotCouponDeduct"`
	FuturesCouponDeduct     decimal.Decimal      `json:"futuresCouponDeduct"`
	SpotFeeDiscountDeduct   decimal.Decimal      `json:"spotFeeDiscountDeduct"`
	NegativeMakerFeeDeduct  decimal.Decimal      `json:"negativeMakerFeeDeduct"`
	FeePaid                 decimal.Decimal      `json:"feePaid"`
	DirectCommission        decimal.Decimal      `json:"directCommission"`
	SubCommission           decimal.Decimal      `json:"subCommission"`
	PartnerCommission       decimal.Decimal      `json:"partnerCommission"`
	PartnerActualCommission decimal.Decimal      `json:"partnerActualCommission"`
	TraderType              CommissionTraderType `json:"traderType"` // user, trader
	ApiType                 CommissionApiType    `json:"apiType"`    // api, non_api
	Status                  CommissionStatus     `json:"status"`     // settled, unsettled, notIssued
	StartCalculationTime    time.Time            `json:"startCalculationTime"`
	EndCalculationTime      time.Time            `json:"endCalculationTime"`
}
