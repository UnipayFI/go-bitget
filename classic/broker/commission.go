package broker

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CommissionBizType is the top-level business line a broker commission row
// belongs to (order-commission endpoint).
type CommissionBizType string

const (
	CommissionBizTypeSpot    CommissionBizType = "spot"
	CommissionBizTypeFutures CommissionBizType = "futures"
)

// CommissionSubBizType is the finer-grained product line of a broker commission
// row (order-commission endpoint).
type CommissionSubBizType string

const (
	CommissionSubBizTypeSpotTrade   CommissionSubBizType = "spot_trade"
	CommissionSubBizTypeSpotMargin  CommissionSubBizType = "spot_margin"
	CommissionSubBizTypeUSDTFutures CommissionSubBizType = "usdt_futures"
	CommissionSubBizTypeCoinFutures CommissionSubBizType = "coin_futures"
	CommissionSubBizTypeUSDCFutures CommissionSubBizType = "usdc_futures"
)

// RebateAffiliationType is the broker/client rebate relationship type
// (rebate-info endpoint).
type RebateAffiliationType string

const (
	RebateAffiliationTypeAffiliate RebateAffiliationType = "affiliate"
	RebateAffiliationTypeOfficial  RebateAffiliationType = "official"
)

// GetTotalCommissionService -- GET /api/v2/broker/total-commission (private; broker-gated)
//
// Returns the broker's daily aggregate commission, trading volume and active
// trader counts, broken down by spot and futures business lines.
type GetTotalCommissionService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetTotalCommissionService() *GetTotalCommissionService {
	return &GetTotalCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime sets the window start. startTime and endTime must be set
// together; the range cannot exceed 180 days.
func (s *GetTotalCommissionService) SetStartTime(t time.Time) *GetTotalCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the window end. startTime and endTime must be set together;
// the range cannot exceed 180 days.
func (s *GetTotalCommissionService) SetEndTime(t time.Time) *GetTotalCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetTotalCommissionService) Do(ctx context.Context) ([]TotalCommission, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/total-commission", s.params).WithSign()
	resp, err := request.Do[[]TotalCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TotalCommission is one day's aggregate broker commission record.
type TotalCommission struct {
	Date               string                 `json:"date"`
	TotalTradingVolume decimal.Decimal        `json:"totalTradingVolume"`
	TotalActiveTraders string                 `json:"totalActiveTraders"`
	TotalCommission    decimal.Decimal        `json:"totalCommission"`
	Spot               TotalCommissionSpot    `json:"spot"`
	Futures            TotalCommissionFutures `json:"futures"`
}

// TotalCommissionSpot is the spot business-line breakdown of a daily total.
type TotalCommissionSpot struct {
	SpotTradingVolume  decimal.Decimal `json:"spotTradingVolume"`
	SpotTradingFee     decimal.Decimal `json:"spotTradingFee"`
	SpotPureTradingFee decimal.Decimal `json:"spotPureTradingFee"`
	SpotCommission     decimal.Decimal `json:"spotCommission"`
}

// TotalCommissionFutures is the futures business-line breakdown of a daily total.
type TotalCommissionFutures struct {
	FuturesTradingVolume  decimal.Decimal `json:"futuresTradingVolume"`
	FuturesTradingFee     decimal.Decimal `json:"futuresTradingFee"`
	FuturesPureTradingFee decimal.Decimal `json:"futuresPureTradingFee"`
	FuturesCommission     decimal.Decimal `json:"futuresCommission"`
}

// GetOrderCommissionService -- GET /api/v2/broker/order-commission (private; broker-gated)
//
// Returns the per-fill broker commission list within a time window, with cursor
// pagination.
type GetOrderCommissionService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetOrderCommissionService() *GetOrderCommissionService {
	return &GetOrderCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime sets the window start (max range 180 days).
func (s *GetOrderCommissionService) SetStartTime(t time.Time) *GetOrderCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the window end (max range 180 days).
func (s *GetOrderCommissionService) SetEndTime(t time.Time) *GetOrderCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (max 500, default 500).
func (s *GetOrderCommissionService) SetLimit(limit int) *GetOrderCommissionService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetUID filters to a single user ID.
func (s *GetOrderCommissionService) SetUID(uid string) *GetOrderCommissionService {
	s.params["uid"] = uid
	return s
}

// SetOrderID filters to a single order ID.
func (s *GetOrderCommissionService) SetOrderID(orderID string) *GetOrderCommissionService {
	s.params["orderid"] = orderID
	return s
}

// SetIDLessThan sets the pagination cursor (the previous response's endId).
func (s *GetOrderCommissionService) SetIDLessThan(idLessThan string) *GetOrderCommissionService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetOrderCommissionService) Do(ctx context.Context) (*OrderCommission, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/order-commission", s.params).WithSign()
	return request.Do[OrderCommission](req)
}

// OrderCommission is the paginated per-fill broker commission response.
type OrderCommission struct {
	CommissionList []OrderCommissionItem `json:"commissionlist"`
	EndID          string                `json:"endId"`
}

// OrderCommissionItem is a single fill's broker commission record.
type OrderCommissionItem struct {
	FillID        string               `json:"fillId"`
	OrderID       string               `json:"orderId"`
	Ts            time.Time            `json:"ts"`
	ClientOrderID string               `json:"clientOid"`
	BizType       CommissionBizType    `json:"bizType"`
	SubBizType    CommissionSubBizType `json:"subBizType"`
	Symbol        string               `json:"symbol"`
	Volume        decimal.Decimal      `json:"volume"`
	Fee           decimal.Decimal      `json:"fee"`
	PureFee       decimal.Decimal      `json:"pureFee"`
	RebateAmount  decimal.Decimal      `json:"rebateAmount"`
}

// GetRebateInfoService -- GET /api/v2/broker/rebate-info (private; broker-gated)
//
// Returns the broker's rebate relationship type, user tier and the spot/futures
// rebate ratios available for a given user.
type GetRebateInfoService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetRebateInfoService(uid string) *GetRebateInfoService {
	return &GetRebateInfoService{c: c, params: map[string]string{"uid": uid}}
}

func (s *GetRebateInfoService) Do(ctx context.Context) (*RebateInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/rebate-info", s.params).WithSign()
	return request.Do[RebateInfo](req)
}

// RebateInfo is the broker rebate configuration for a user.
type RebateInfo struct {
	AffiliationType          RebateAffiliationType `json:"affiliationType"`
	UserLevel                string                `json:"userLevel"` // VIP0..VIP7 / PRO1..PRO6
	ClientSpotRebateRatio    decimal.Decimal       `json:"clientSpotRebateRatio"`
	ClientFuturesRebateRatio decimal.Decimal       `json:"clientFuturesRebateRatio"`
}

// SubAffiliateRebateType is the product line (and regular/API/trading-expert
// variant) a sub-affiliate rebate rate applies to (sub-affiliate-info endpoint).
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

// GetSubAffiliateInfoService -- GET /api/v2/broker/sub-affiliate-info (private; broker-gated)
//
// Returns the broker's sub-affiliates with their rebate rates, referral counts
// and USDT-denominated volume/deposit/fee/rebate totals (the fields shown on the
// affiliate dashboard), with cursor pagination.
type GetSubAffiliateInfoService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubAffiliateInfoService() *GetSubAffiliateInfoService {
	return &GetSubAffiliateInfoService{c: c, params: map[string]string{}}
}

// SetUID filters to a single sub-affiliate.
func (s *GetSubAffiliateInfoService) SetUID(uid string) *GetSubAffiliateInfoService {
	s.params["uid"] = uid
	return s
}

// SetStartTime sets the window start (defaults to the last 7 days when unset).
func (s *GetSubAffiliateInfoService) SetStartTime(t time.Time) *GetSubAffiliateInfoService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the window end (max range 30 days).
func (s *GetSubAffiliateInfoService) SetEndTime(t time.Time) *GetSubAffiliateInfoService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIDLessThan sets the pagination cursor (the previous response's endId).
func (s *GetSubAffiliateInfoService) SetIDLessThan(idLessThan string) *GetSubAffiliateInfoService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit sets the page size (max 100, default 100).
func (s *GetSubAffiliateInfoService) SetLimit(limit int) *GetSubAffiliateInfoService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIncludeDownLine sets whether downline users are included (default yes).
func (s *GetSubAffiliateInfoService) SetIncludeDownLine(include bool) *GetSubAffiliateInfoService {
	if include {
		s.params["includeDownLine"] = "yes"
	} else {
		s.params["includeDownLine"] = "no"
	}
	return s
}

func (s *GetSubAffiliateInfoService) Do(ctx context.Context) (*SubAffiliateInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/sub-affiliate-info", s.params).WithSign()
	return request.Do[SubAffiliateInfo](req)
}

// SubAffiliateInfo is the paginated sub-affiliate list.
type SubAffiliateInfo struct {
	List  []SubAffiliate `json:"list"`
	EndID string         `json:"endId"`
}

// SubAffiliate is one sub-affiliate's dashboard summary. Volume, Deposit,
// TransactionFee and Rebate are in USDT; Rebate is what the sub-affiliate paid
// back to the channel.
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

// SubAffiliateRebateRate is a sub-affiliate's rebate rate for one product line;
// RebateRate is a percentage (5 means 5%).
type SubAffiliateRebateRate struct {
	RebateType SubAffiliateRebateType `json:"rebateType"`
	RebateRate decimal.Decimal        `json:"rebateRate"`
}
