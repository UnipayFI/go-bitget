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
