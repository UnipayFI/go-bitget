package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// AffiliateKYCResult is a referred customer's KYC/KYB verification outcome.
type AffiliateKYCResult string

const (
	AffiliateKYCResultPassed    AffiliateKYCResult = "passed"
	AffiliateKYCResultNotPassed AffiliateKYCResult = "not_passed"
)

// SubAffiliateRebateType is the product line (and regular/API/trading-expert
// variant) a sub-affiliate rebate rate applies to.
type SubAffiliateRebateType string

const (
	SubAffiliateRebateTypeFutures          SubAffiliateRebateType = "futures"
	SubAffiliateRebateTypeFuturesAPI       SubAffiliateRebateType = "futures_api"
	SubAffiliateRebateTypeFuturesTrader    SubAffiliateRebateType = "futures_trader"
	SubAffiliateRebateTypeFuturesTraderAPI SubAffiliateRebateType = "futures_trader_api"
	SubAffiliateRebateTypeSpot             SubAffiliateRebateType = "spot"
	SubAffiliateRebateTypeSpotAPI          SubAffiliateRebateType = "spot_api"
	SubAffiliateRebateTypeSpotTrader       SubAffiliateRebateType = "spot_trader"
	SubAffiliateRebateTypeSpotTraderAPI    SubAffiliateRebateType = "spot_trader_api"
	SubAffiliateRebateTypeOnchain          SubAffiliateRebateType = "onchain"
	SubAffiliateRebateTypeOnchainAPI       SubAffiliateRebateType = "onchain_api"
	SubAffiliateRebateTypeOnchainTrader    SubAffiliateRebateType = "onchain_trader"
	SubAffiliateRebateTypeOnchainTraderAPI SubAffiliateRebateType = "onchain_trader_api"
	SubAffiliateRebateTypeCFD              SubAffiliateRebateType = "cfd"
	SubAffiliateRebateTypeCFDAPI           SubAffiliateRebateType = "cfd_api"
	SubAffiliateRebateTypeCFDTrader        SubAffiliateRebateType = "cfd_trader"
	SubAffiliateRebateTypeCFDTraderAPI     SubAffiliateRebateType = "cfd_trader_api"
	SubAffiliateRebateTypeStock            SubAffiliateRebateType = "stock"
	SubAffiliateRebateTypeStockAPI         SubAffiliateRebateType = "stock_api"
	SubAffiliateRebateTypeStockTrader      SubAffiliateRebateType = "stock_trader"
	SubAffiliateRebateTypeStockTraderAPI   SubAffiliateRebateType = "stock_trader_api"
)

// SubAffiliateVerificationStatus is a sub-affiliate's identity verification
// state; passing either KYC or KYB counts as verified.
type SubAffiliateVerificationStatus string

const (
	SubAffiliateVerificationStatusVerified    SubAffiliateVerificationStatus = "verified"
	SubAffiliateVerificationStatusNotVerified SubAffiliateVerificationStatus = "not_verified"
)

// SubAffiliateCommissionStatus is the settlement state of a sub-affiliate
// commission record.
type SubAffiliateCommissionStatus string

const (
	SubAffiliateCommissionStatusSettled   SubAffiliateCommissionStatus = "settled"
	SubAffiliateCommissionStatusUnsettled SubAffiliateCommissionStatus = "unsettled"
	SubAffiliateCommissionStatusNotIssued SubAffiliateCommissionStatus = "notIssued"
)

// GetAffiliateCustomerListService -- GET /api/v3/affiliate/customer-list (affiliate)
//
// Returns the affiliate's direct customers with their registration times. With
// no time range every registered customer is returned (no 90-day limit); pass a
// uid to check whether a specific user is your referral. Real-time data.
type GetAffiliateCustomerListService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerListService() *GetAffiliateCustomerListService {
	return &GetAffiliateCustomerListService{c: c, params: map[string]string{}}
}

// SetStartTime filters customers registered at or after t. startTime and
// endTime must be set together and span at most 30 days.
func (s *GetAffiliateCustomerListService) SetStartTime(t time.Time) *GetAffiliateCustomerListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters customers registered at or before t (see SetStartTime).
func (s *GetAffiliateCustomerListService) SetEndTime(t time.Time) *GetAffiliateCustomerListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAffiliateCustomerListService) SetPageNo(pageNo string) *GetAffiliateCustomerListService {
	s.params["pageNo"] = pageNo
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetAffiliateCustomerListService) SetPageSize(pageSize string) *GetAffiliateCustomerListService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetAffiliateCustomerListService) SetUID(uid string) *GetAffiliateCustomerListService {
	s.params["uid"] = uid
	return s
}

func (s *GetAffiliateCustomerListService) SetReferralCode(referralCode string) *GetAffiliateCustomerListService {
	s.params["referralCode"] = referralCode
	return s
}

func (s *GetAffiliateCustomerListService) Do(ctx context.Context) ([]AffiliateCustomer, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-list", s.params).WithSign()
	resp, err := request.Do[[]AffiliateCustomer](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AffiliateCustomer is one referred customer. uid is a bare JSON number here
// (the doc example and the live v2 customer-list both send it unquoted), so it
// is typed int64.
type AffiliateCustomer struct {
	UID          int64     `json:"uid"`
	RegisterTime time.Time `json:"registerTime,format:unixmilli"`
}

// GetAffiliateCustomerAssetService -- GET /api/v3/affiliate/customer-asset (affiliate)
//
// Returns account-balance snapshots of the affiliate's direct customers.
// Refreshed every 10 minutes.
type GetAffiliateCustomerAssetService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerAssetService() *GetAffiliateCustomerAssetService {
	return &GetAffiliateCustomerAssetService{c: c, params: map[string]string{}}
}

func (s *GetAffiliateCustomerAssetService) SetPageNo(pageNo string) *GetAffiliateCustomerAssetService {
	s.params["pageNo"] = pageNo
	return s
}

// SetPageSize sets the page size (default 100, max 500).
func (s *GetAffiliateCustomerAssetService) SetPageSize(pageSize string) *GetAffiliateCustomerAssetService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetAffiliateCustomerAssetService) SetUID(uid string) *GetAffiliateCustomerAssetService {
	s.params["uid"] = uid
	return s
}

// SetShowSub sets whether direct customers' subordinate users are included
// ("yes" or "no").
func (s *GetAffiliateCustomerAssetService) SetShowSub(showSub string) *GetAffiliateCustomerAssetService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateCustomerAssetService) Do(ctx context.Context) ([]AffiliateCustomerAsset, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-asset", s.params).WithSign()
	resp, err := request.Do[[]AffiliateCustomerAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AffiliateCustomerAsset is one direct customer's balance snapshot. Remark
// explains a missing balance (e.g. "sub account exceed 5").
type AffiliateCustomerAsset struct {
	UID     string          `json:"uid"`
	Balance decimal.Decimal `json:"balance"`
	UTime   time.Time       `json:"uTime,format:unixmilli"`
	Remark  string          `json:"remark"`
}

// GetAffiliateCustomerTradeVolumeService -- GET /api/v3/affiliate/customer-trade-volume (affiliate)
//
// Returns the spot and futures trading volume of the affiliate's direct
// customers. Only the last 90 days are queryable; with no time range it returns
// yesterday (00:00-23:59 UTC+8). Refreshed hourly.
type GetAffiliateCustomerTradeVolumeService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerTradeVolumeService() *GetAffiliateCustomerTradeVolumeService {
	return &GetAffiliateCustomerTradeVolumeService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t. startTime and endTime must be set
// together and span at most 30 days.
func (s *GetAffiliateCustomerTradeVolumeService) SetStartTime(t time.Time) *GetAffiliateCustomerTradeVolumeService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (see SetStartTime).
func (s *GetAffiliateCustomerTradeVolumeService) SetEndTime(t time.Time) *GetAffiliateCustomerTradeVolumeService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAffiliateCustomerTradeVolumeService) SetPageNo(pageNo string) *GetAffiliateCustomerTradeVolumeService {
	s.params["pageNo"] = pageNo
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetAffiliateCustomerTradeVolumeService) SetPageSize(pageSize string) *GetAffiliateCustomerTradeVolumeService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetAffiliateCustomerTradeVolumeService) SetUID(uid string) *GetAffiliateCustomerTradeVolumeService {
	s.params["uid"] = uid
	return s
}

// SetShowSub sets whether sub-customer information is included ("yes" or
// "no").
func (s *GetAffiliateCustomerTradeVolumeService) SetShowSub(showSub string) *GetAffiliateCustomerTradeVolumeService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateCustomerTradeVolumeService) Do(ctx context.Context) ([]AffiliateCustomerTradeVolume, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-trade-volume", s.params).WithSign()
	resp, err := request.Do[[]AffiliateCustomerTradeVolume](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AffiliateCustomerTradeVolume is one direct customer's trading volume.
// Volumn (spelled that way on the wire) is the spot + futures total.
type AffiliateCustomerTradeVolume struct {
	UID           string          `json:"uid"`
	Volumn        decimal.Decimal `json:"volumn"`
	SpotVolume    decimal.Decimal `json:"spotVolume"`
	FuturesVolume decimal.Decimal `json:"futuresVolume"`
	Time          time.Time       `json:"time,format:unixmilli"`
}

// GetAffiliateCustomerDepositService -- GET /api/v3/affiliate/customer-deposit (affiliate)
//
// Returns the on-chain deposits and internal transfers of the affiliate's
// direct customers. Only the last 90 days are queryable; with no time range it
// returns yesterday (00:00-23:59 UTC+8). Real-time data.
type GetAffiliateCustomerDepositService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerDepositService() *GetAffiliateCustomerDepositService {
	return &GetAffiliateCustomerDepositService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t. startTime and endTime must be set
// together and span at most 30 days.
func (s *GetAffiliateCustomerDepositService) SetStartTime(t time.Time) *GetAffiliateCustomerDepositService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (see SetStartTime).
func (s *GetAffiliateCustomerDepositService) SetEndTime(t time.Time) *GetAffiliateCustomerDepositService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAffiliateCustomerDepositService) SetPageNo(pageNo string) *GetAffiliateCustomerDepositService {
	s.params["pageNo"] = pageNo
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetAffiliateCustomerDepositService) SetPageSize(pageSize string) *GetAffiliateCustomerDepositService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetAffiliateCustomerDepositService) SetUID(uid string) *GetAffiliateCustomerDepositService {
	s.params["uid"] = uid
	return s
}

// SetShowSub sets whether direct customers' subordinate users are included
// ("yes" or "no").
func (s *GetAffiliateCustomerDepositService) SetShowSub(showSub string) *GetAffiliateCustomerDepositService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateCustomerDepositService) Do(ctx context.Context) ([]AffiliateCustomerDeposit, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-deposit", s.params).WithSign()
	resp, err := request.Do[[]AffiliateCustomerDeposit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AffiliateCustomerDeposit is one deposit of a direct customer.
type AffiliateCustomerDeposit struct {
	UID           string          `json:"uid"`
	OrderID       string          `json:"orderId"`
	DepositTime   time.Time       `json:"depositTime,format:unixmilli"`
	DepositCoin   string          `json:"depositCoin"`
	DepositAmount decimal.Decimal `json:"depositAmount"`
}

// GetAffiliateCustomerKYCResultService -- GET /api/v3/affiliate/customer-kyc-result (affiliate)
//
// Returns the KYC/KYB result of the affiliate's direct customers, paginated by
// cursor (endId). With no time range it returns yesterday's results
// (00:00-23:59 UTC+8). Refreshed every 5 minutes.
type GetAffiliateCustomerKYCResultService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerKYCResultService() *GetAffiliateCustomerKYCResultService {
	return &GetAffiliateCustomerKYCResultService{c: c, params: map[string]string{}}
}

// SetStartTime filters results at or after t. startTime and endTime must be set
// together and span at most 90 days.
func (s *GetAffiliateCustomerKYCResultService) SetStartTime(t time.Time) *GetAffiliateCustomerKYCResultService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters results at or before t (see SetStartTime).
func (s *GetAffiliateCustomerKYCResultService) SetEndTime(t time.Time) *GetAffiliateCustomerKYCResultService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetAffiliateCustomerKYCResultService) SetCursor(cursor string) *GetAffiliateCustomerKYCResultService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit sets the page size (default 100, max 1000).
func (s *GetAffiliateCustomerKYCResultService) SetLimit(limit string) *GetAffiliateCustomerKYCResultService {
	s.params["limit"] = limit
	return s
}

func (s *GetAffiliateCustomerKYCResultService) SetUID(uid string) *GetAffiliateCustomerKYCResultService {
	s.params["uid"] = uid
	return s
}

// SetShowSub sets whether direct customers' subordinate users are included
// ("yes" or "no").
func (s *GetAffiliateCustomerKYCResultService) SetShowSub(showSub string) *GetAffiliateCustomerKYCResultService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateCustomerKYCResultService) Do(ctx context.Context) (*AffiliateCustomerKYCResults, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-kyc-result", s.params).WithSign()
	return request.Do[AffiliateCustomerKYCResults](req)
}

// AffiliateCustomerKYCResults is a page of KYC results. The doc's field table
// flattens uid/kycResult into data, but its example (and the live v2
// customer-kyc-result) nest them in userList, which is what is modelled here.
type AffiliateCustomerKYCResults struct {
	UserList []AffiliateCustomerKYC `json:"userList"`
	EndID    string                 `json:"endId"`
}

type AffiliateCustomerKYC struct {
	UID       string             `json:"uid"`
	KYCResult AffiliateKYCResult `json:"kycResult"` // passed, not_passed
}

// GetAffiliateCustomerCommissionService -- GET /api/v3/affiliate/customer-commission (affiliate)
//
// Returns the affiliate's direct commission records per referred customer,
// paginated by cursor (endId). Only the last 90 days are queryable; with no
// time range it returns yesterday (00:00-23:59 UTC+8). Updated T+1 (UTC+8).
type GetAffiliateCustomerCommissionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateCustomerCommissionService() *GetAffiliateCustomerCommissionService {
	return &GetAffiliateCustomerCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t. startTime and endTime must be set
// together and span at most 30 days.
func (s *GetAffiliateCustomerCommissionService) SetStartTime(t time.Time) *GetAffiliateCustomerCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (see SetStartTime).
func (s *GetAffiliateCustomerCommissionService) SetEndTime(t time.Time) *GetAffiliateCustomerCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetAffiliateCustomerCommissionService) SetCursor(cursor string) *GetAffiliateCustomerCommissionService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit sets the page size (default 100, max 1000).
func (s *GetAffiliateCustomerCommissionService) SetLimit(limit string) *GetAffiliateCustomerCommissionService {
	s.params["limit"] = limit
	return s
}

func (s *GetAffiliateCustomerCommissionService) SetUID(uid string) *GetAffiliateCustomerCommissionService {
	s.params["uid"] = uid
	return s
}

func (s *GetAffiliateCustomerCommissionService) SetCoin(coin string) *GetAffiliateCustomerCommissionService {
	s.params["coin"] = coin
	return s
}

// SetSymbol filters by symbol in the legacy suffixed form: BGBUSDT_SPBL (spot),
// BTCUSDT_SPBL-MABL (spot margin), BTCUSDT_UMCBL (USDT-Futures), BTCUSD_DMCBL
// (COIN-Futures), BTCPERP_CMCBL (USDC-Futures).
func (s *GetAffiliateCustomerCommissionService) SetSymbol(symbol string) *GetAffiliateCustomerCommissionService {
	s.params["symbol"] = symbol
	return s
}

// SetShowSub sets whether direct customers' subordinate users are included
// ("yes" or "no").
func (s *GetAffiliateCustomerCommissionService) SetShowSub(showSub string) *GetAffiliateCustomerCommissionService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateCustomerCommissionService) Do(ctx context.Context) (*AffiliateCustomerCommissions, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/customer-commission", s.params).WithSign()
	return request.Do[AffiliateCustomerCommissions](req)
}

type AffiliateCustomerCommissions struct {
	CommissionList []AffiliateCustomerCommission `json:"commissionList"`
	EndID          string                        `json:"endId"`
}

// AffiliateCustomerCommission is one day's commission from a referred
// customer. Amounts are in Coin and cover maker+taker. RebateAmount goes to the
// affiliate, UserTotalRebateAmount back to the customer (only when user rebates
// are configured), DayTotalRebateAmount is their sum, and TotalRebateAmount is
// the rolling 90-day total.
type AffiliateCustomerCommission struct {
	UID                    string          `json:"uid"`
	Date                   time.Time       `json:"date,format:unixmilli"` // commission date (UTC+8)
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

// GetAffiliateSubCustomerListService -- GET /api/v3/affiliate/sub-customer-list (affiliate)
//
// Returns the affiliate's sub-customers with their registration times,
// paginated by cursor (minId). With no time range every registered customer is
// returned (no 90-day limit).
type GetAffiliateSubCustomerListService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAffiliateSubCustomerListService() *GetAffiliateSubCustomerListService {
	return &GetAffiliateSubCustomerListService{c: c, params: map[string]string{}}
}

// SetStartTime filters customers registered at or after t. startTime and
// endTime must be set together and span at most 30 days.
func (s *GetAffiliateSubCustomerListService) SetStartTime(t time.Time) *GetAffiliateSubCustomerListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters customers registered at or before t (see SetStartTime).
func (s *GetAffiliateSubCustomerListService) SetEndTime(t time.Time) *GetAffiliateSubCustomerListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetCursor sets the pagination cursor (minId from the previous response; omit
// on the first call).
func (s *GetAffiliateSubCustomerListService) SetCursor(cursor string) *GetAffiliateSubCustomerListService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetAffiliateSubCustomerListService) SetLimit(limit string) *GetAffiliateSubCustomerListService {
	s.params["limit"] = limit
	return s
}

func (s *GetAffiliateSubCustomerListService) SetUID(uid string) *GetAffiliateSubCustomerListService {
	s.params["uid"] = uid
	return s
}

// SetShowSub sets whether direct customers' subordinate users are included
// ("yes" or "no").
func (s *GetAffiliateSubCustomerListService) SetShowSub(showSub string) *GetAffiliateSubCustomerListService {
	s.params["showSub"] = showSub
	return s
}

func (s *GetAffiliateSubCustomerListService) Do(ctx context.Context) (*AffiliateSubCustomers, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/sub-customer-list", s.params).WithSign()
	return request.Do[AffiliateSubCustomers](req)
}

// AffiliateSubCustomers is a page of sub-customers. The doc's field table
// flattens uid/registerTime into data, but its example nests them in list
// beside the minId cursor, which is what is modelled here.
type AffiliateSubCustomers struct {
	List  []AffiliateCustomer `json:"list"`
	MinID string              `json:"minId"`
}

// GetSubAffiliateInfoService -- GET /api/v3/affiliate/sub-affiliate-info (affiliate)
//
// Returns the affiliate's sub-affiliates with their rebate rates, referral
// counts and USDT-denominated volume/deposit/fee/rebate totals, paginated by
// idLessThan (endId).
type GetSubAffiliateInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubAffiliateInfoService() *GetSubAffiliateInfoService {
	return &GetSubAffiliateInfoService{c: c, params: map[string]string{}}
}

// SetUID filters to a single sub-affiliate.
func (s *GetSubAffiliateInfoService) SetUID(uid string) *GetSubAffiliateInfoService {
	s.params["uid"] = uid
	return s
}

// SetStartTime sets the statistics window start (window at most 30 days).
func (s *GetSubAffiliateInfoService) SetStartTime(t time.Time) *GetSubAffiliateInfoService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the statistics window end (window at most 30 days).
func (s *GetSubAffiliateInfoService) SetEndTime(t time.Time) *GetSubAffiliateInfoService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan sets the pagination cursor (endId from the previous response).
func (s *GetSubAffiliateInfoService) SetIDLessThan(idLessThan string) *GetSubAffiliateInfoService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetSubAffiliateInfoService) SetLimit(limit string) *GetSubAffiliateInfoService {
	s.params["limit"] = limit
	return s
}

// SetIncludeDownLine sets whether downline users are included ("yes" or "no",
// default "yes").
func (s *GetSubAffiliateInfoService) SetIncludeDownLine(includeDownLine string) *GetSubAffiliateInfoService {
	s.params["includeDownLine"] = includeDownLine
	return s
}

func (s *GetSubAffiliateInfoService) Do(ctx context.Context) (*SubAffiliateInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/sub-affiliate-info", s.params).WithSign()
	return request.Do[SubAffiliateInfo](req)
}

type SubAffiliateInfo struct {
	List  []SubAffiliate `json:"list"`
	EndID string         `json:"endId"`
}

// SubAffiliate is one sub-affiliate's summary. Volume, Deposit, TransactionFee
// and Rebate are in USDT; Rebate is what the sub-affiliate paid back to the
// channel.
type SubAffiliate struct {
	UID                         string                         `json:"uid"`
	UplineName                  string                         `json:"uplineName"`
	RebateRateList              []SubAffiliateRebateRate       `json:"rebateRateList"`
	VerificationStatus          SubAffiliateVerificationStatus `json:"verificationStatus"`
	DirectReferrals             string                         `json:"directReferrals"`
	SubAffiliates               string                         `json:"subAffiliates"`
	SubAffiliateDirectReferrals string                         `json:"subAffiliateDirectReferrals"`
	Volume                      decimal.Decimal                `json:"volume"`
	Deposit                     decimal.Decimal                `json:"deposit"`
	TransactionFee              decimal.Decimal                `json:"transactionFee"`
	Rebate                      decimal.Decimal                `json:"rebate"`
}

// SubAffiliateRebateRate is a rebate rate for one product line; RebateRate is a
// percentage (5 means 5%).
type SubAffiliateRebateRate struct {
	RebateType SubAffiliateRebateType `json:"rebateType"`
	RebateRate decimal.Decimal        `json:"rebateRate"`
}

// GetSubAffiliateCommissionService -- GET /api/v3/affiliate/sub-affiliate-commission (affiliate)
//
// Returns the commissions of the affiliate's sub-affiliates, paginated by
// cursor (endId). Only the last 90 days are queryable; with no time range it
// returns yesterday (00:00-23:59 UTC+8).
type GetSubAffiliateCommissionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubAffiliateCommissionService() *GetSubAffiliateCommissionService {
	return &GetSubAffiliateCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t. startTime and endTime must be set
// together and span at most 7 days.
func (s *GetSubAffiliateCommissionService) SetStartTime(t time.Time) *GetSubAffiliateCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (see SetStartTime).
func (s *GetSubAffiliateCommissionService) SetEndTime(t time.Time) *GetSubAffiliateCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetSubAffiliateCommissionService) SetLimit(limit string) *GetSubAffiliateCommissionService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (endId from the previous response; omit
// on the first call).
func (s *GetSubAffiliateCommissionService) SetCursor(cursor string) *GetSubAffiliateCommissionService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSubAffiliateCommissionService) Do(ctx context.Context) (*SubAffiliateCommissions, error) {
	req := request.Get(ctx, s.c, "/api/v3/affiliate/sub-affiliate-commission", s.params).WithSign()
	return request.Do[SubAffiliateCommissions](req)
}

type SubAffiliateCommissions struct {
	CommissionList []SubAffiliateCommission `json:"commissionList"`
	EndID          string                   `json:"endId"`
}

// SubAffiliateCommission is one sub-affiliate commission row; amounts are
// totals over the requested window. DirectCommission comes from the partner's
// first-level customers, SubCommission from deeper levels, PartnerCommission is
// their sum, and PartnerActualCommission is what remains after risk-control and
// special-case deductions.
type SubAffiliateCommission struct {
	UID                     string                       `json:"uid"`
	BizType                 string                       `json:"bizType"`    // spot, futures
	SubBizType              string                       `json:"subBizType"` // spot, margin, usdt_futures, coin_futures, usdc_futures
	Symbol                  string                       `json:"symbol"`
	Coin                    string                       `json:"coin"`
	Fee                     decimal.Decimal              `json:"fee"`
	Volume                  decimal.Decimal              `json:"volume"`
	ActivityBonusDeduct     decimal.Decimal              `json:"activityBonusDeduct"`
	SpotCouponDeduct        decimal.Decimal              `json:"spotCouponDeduct"`
	FuturesCouponDeduct     decimal.Decimal              `json:"futuresCouponDeduct"`
	SpotFeeDiscountDeduct   decimal.Decimal              `json:"spotFeeDiscountDeduct"`
	NegativeMakerFeeDeduct  decimal.Decimal              `json:"negativeMakerFeeDeduct"`
	FeePaid                 decimal.Decimal              `json:"feePaid"`
	DirectCommission        decimal.Decimal              `json:"directCommission"`
	SubCommission           decimal.Decimal              `json:"subCommission"`
	PartnerCommission       decimal.Decimal              `json:"partnerCommission"`
	PartnerActualCommission decimal.Decimal              `json:"partnerActualCommission"`
	TraderType              string                       `json:"traderType"` // user, trader
	APIType                 string                       `json:"apiType"`    // api, non_api
	Status                  SubAffiliateCommissionStatus `json:"status"`     // settled, unsettled, notIssued
	StartCalculationTime    time.Time                    `json:"startCalculationTime,format:unixmilli"`
	EndCalculationTime      time.Time                    `json:"endCalculationTime,format:unixmilli"`
}
