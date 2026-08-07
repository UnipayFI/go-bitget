package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetCopyTradingPairsService -- GET /api/v3/copy/futures/trading-pairs (Elite trading read)
//
// Returns the futures trading pairs available to an elite (lead) trader, with the
// per-margin-coin long/short position capacity for each symbol.
type GetCopyTradingPairsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCopyTradingPairsService() *GetCopyTradingPairsService {
	return &GetCopyTradingPairsService{c: c}
}

func (s *GetCopyTradingPairsService) Do(ctx context.Context) ([]CopyTradingPair, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/trading-pairs").WithSign()
	resp, err := request.Do[[]CopyTradingPair](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CopyTradingPair struct {
	Symbol        string             `json:"symbol"`
	Leverage      string             `json:"leverage"`
	MarginDetails []CopyMarginDetail `json:"marginDetails"`
}

type CopyMarginDetail struct {
	MarginCoin          string `json:"marginCoin"`
	MaxLongCount        string `json:"maxLongCount"`
	RemainingLongCount  string `json:"remainingLongCount"`
	MaxShortCount       string `json:"maxShortCount"`
	RemainingShortCount string `json:"remainingShortCount"`
}

// GetCopyPositionSummaryService -- GET /api/v3/copy/futures/position-summary (Elite trading read)
//
// Returns the elite (lead) trader's open futures positions with their margin,
// PnL and risk metrics.
type GetCopyPositionSummaryService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCopyPositionSummaryService() *GetCopyPositionSummaryService {
	return &GetCopyPositionSummaryService{c: c}
}

func (s *GetCopyPositionSummaryService) Do(ctx context.Context) ([]CopyPositionSummary, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/position-summary").WithSign()
	resp, err := request.Do[[]CopyPositionSummary](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type CopyPositionSummary struct {
	Symbol        string          `json:"symbol"`
	HoldSide      PosSide         `json:"holdSide"`
	HoldSize      decimal.Decimal `json:"holdSize"`
	AvgPrice      decimal.Decimal `json:"avgPrice"`
	MarkPrice     decimal.Decimal `json:"markPrice"`
	LiqPrice      decimal.Decimal `json:"liqPrice"`
	Leverage      string          `json:"leverage"`
	MarginMode    string          `json:"marginMode"` // isolated, cross
	Margin        decimal.Decimal `json:"margin"`
	PositionValue decimal.Decimal `json:"positionValue"`
	UnrealizedPnL decimal.Decimal `json:"unrealizedPnl"`
	RealizedPnL   decimal.Decimal `json:"realizedPnl"`
	ROI           decimal.Decimal `json:"roi"`
}

// GetCopyMaxTransferableService -- GET /api/v3/copy/futures/max-transferable (Elite trading read)
//
// Returns the maximum amount of a coin that can be transferred out of the elite
// (lead) account, alongside the available balance.
type GetCopyMaxTransferableService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyMaxTransferableService(coin string) *GetCopyMaxTransferableService {
	return &GetCopyMaxTransferableService{c: c, params: map[string]string{"coin": coin}}
}

func (s *GetCopyMaxTransferableService) Do(ctx context.Context) (*CopyMaxTransferable, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/max-transferable", s.params).WithSign()
	return request.Do[CopyMaxTransferable](req)
}

type CopyMaxTransferable struct {
	MaxTransferable decimal.Decimal `json:"maxTransferable"`
	Available       decimal.Decimal `json:"available"`
}

// CopyTransferService -- POST /api/v3/copy/futures/transfer (Elite trading read & write)
//
// Moves a coin between the elite (lead) account and the spot/funding account.
// type is "in" (spot/funding -> lead) or "out" (lead -> spot/funding).
type CopyTransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCopyTransferService(transferType, coin string, amount decimal.Decimal) *CopyTransferService {
	return &CopyTransferService{c: c, body: map[string]any{
		"type":   transferType,
		"coin":   coin,
		"amount": amount.String(),
	}}
}

// SetInAccountType restricts the source accounts a transfer-in may draw from:
// a comma-separated list of "funding", "uta" and "otc". All three are used when
// omitted, deducted in the order funding -> otc -> uta.
func (s *CopyTransferService) SetInAccountType(inAccountType string) *CopyTransferService {
	s.body["inAccountType"] = inAccountType
	return s
}

func (s *CopyTransferService) Do(ctx context.Context) (*CopyTransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/transfer", s.body).WithSign()
	return request.Do[CopyTransferResult](req)
}

type CopyTransferResult struct {
	TransferID string `json:"transferId"`
}

// GetCopyTransferRecordService -- GET /api/v3/copy/futures/transfer-record (Elite trading read)
//
// Returns the elite (lead) account's transfer records, paginated by cursor (pass
// the previous response's transferId as the cursor).
type GetCopyTransferRecordService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyTransferRecordService() *GetCopyTransferRecordService {
	return &GetCopyTransferRecordService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t.
func (s *GetCopyTransferRecordService) SetStartTime(t time.Time) *GetCopyTransferRecordService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t.
func (s *GetCopyTransferRecordService) SetEndTime(t time.Time) *GetCopyTransferRecordService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetCopyTransferRecordService) SetLimit(limit string) *GetCopyTransferRecordService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyTransferRecordService) SetCursor(cursor string) *GetCopyTransferRecordService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyTransferRecordService) Do(ctx context.Context) (*CopyTransferRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/transfer-record", s.params).WithSign()
	return request.Do[CopyTransferRecords](req)
}

type CopyTransferRecords struct {
	List []CopyTransferRecord `json:"list"`
}

type CopyTransferRecord struct {
	TransferID string `json:"transferId"`
	// FromType and ToType are spot, uta, lead or otc; a single record may carry
	// several account types, comma-separated.
	FromType    string          `json:"fromType"`
	ToType      string          `json:"toType"`
	Coin        string          `json:"coin"`
	Amount      decimal.Decimal `json:"amount"`
	Status      string          `json:"status"` // Successful, Failed, Processing
	CreatedTime time.Time       `json:"createdTime"`
}

// GetCopyCurrentFollowersService -- GET /api/v3/copy/futures/current-follower (Elite trading read)
//
// Returns the elite (lead) account's active followers, paginated by cursor
// (pass the previous response's endId as the cursor to page forward).
type GetCopyCurrentFollowersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyCurrentFollowersService() *GetCopyCurrentFollowersService {
	return &GetCopyCurrentFollowersService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetCopyCurrentFollowersService) SetLimit(limit string) *GetCopyCurrentFollowersService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyCurrentFollowersService) SetCursor(cursor string) *GetCopyCurrentFollowersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyCurrentFollowersService) Do(ctx context.Context) (*CopyCurrentFollowers, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/current-follower", s.params).WithSign()
	return request.Do[CopyCurrentFollowers](req)
}

type CopyCurrentFollowers struct {
	List []CopyCurrentFollower `json:"list"`
}

type CopyCurrentFollower struct {
	FollowerName     string          `json:"followerName"`
	EstimateAssets   decimal.Decimal `json:"estimateAssets"` // estimated assets, USDT
	TotalAssets      decimal.Decimal `json:"totalAssets"`    // account equity, USDT
	TotalProfit      decimal.Decimal `json:"totalProfit"`
	TotalShareProfit decimal.Decimal `json:"totalShareProfit"`
	TotalInvestment  decimal.Decimal `json:"totalInvestment"`
	CanRemove        string          `json:"canRemove"` // yes, no
	FollowDays       string          `json:"followDays"`
	StartTime        time.Time       `json:"startTime"`
}

// GetCopyHistoryFollowersService -- GET /api/v3/copy/futures/history-follower (Elite trading read)
//
// Returns the elite (lead) account's past followers, paginated by cursor (pass
// the previous response's endId as the cursor to page forward).
type GetCopyHistoryFollowersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyHistoryFollowersService() *GetCopyHistoryFollowersService {
	return &GetCopyHistoryFollowersService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetCopyHistoryFollowersService) SetLimit(limit string) *GetCopyHistoryFollowersService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyHistoryFollowersService) SetCursor(cursor string) *GetCopyHistoryFollowersService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyHistoryFollowersService) Do(ctx context.Context) (*CopyHistoryFollowers, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/history-follower", s.params).WithSign()
	return request.Do[CopyHistoryFollowers](req)
}

type CopyHistoryFollowers struct {
	List []CopyHistoryFollower `json:"list"`
}

type CopyHistoryFollower struct {
	FollowerName     string          `json:"followerName"`
	TotalProfit      decimal.Decimal `json:"totalProfit"`
	TotalShareProfit decimal.Decimal `json:"totalShareProfit"`
	TotalInvestment  decimal.Decimal `json:"totalInvestment"`
	StartTime        time.Time       `json:"startTime"`
	EndTime          time.Time       `json:"endTime"`
}

// GetCopyProfitSummaryService -- GET /api/v3/copy/futures/profit-summary (Elite trading read)
//
// Returns the elite (lead) account's cumulative profit-sharing totals.
type GetCopyProfitSummaryService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetCopyProfitSummaryService() *GetCopyProfitSummaryService {
	return &GetCopyProfitSummaryService{c: c}
}

func (s *GetCopyProfitSummaryService) Do(ctx context.Context) (*CopyProfitSummary, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/profit-summary").WithSign()
	return request.Do[CopyProfitSummary](req)
}

type CopyProfitSummary struct {
	TotalProfit          decimal.Decimal `json:"totalProfit"`          // cumulative profit, USDT
	TotalAllocatedProfit decimal.Decimal `json:"totalAllocatedProfit"` // profit share already settled, USDT
	TotalPendingProfit   decimal.Decimal `json:"totalPendingProfit"`   // profit share awaiting settlement, USDT
}

// GetCopyProfitDetailsService -- GET /api/v3/copy/futures/profit-details (Elite trading read)
//
// Returns the elite (lead) account's per-follower profit-sharing settlements,
// paginated by cursor (pass the previous response's nextCursor).
type GetCopyProfitDetailsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyProfitDetailsService() *GetCopyProfitDetailsService {
	return &GetCopyProfitDetailsService{c: c, params: map[string]string{}}
}

// SetStartTime filters settlements at or after t.
func (s *GetCopyProfitDetailsService) SetStartTime(t time.Time) *GetCopyProfitDetailsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters settlements at or before t.
func (s *GetCopyProfitDetailsService) SetEndTime(t time.Time) *GetCopyProfitDetailsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetCopyProfitDetailsService) SetLimit(limit string) *GetCopyProfitDetailsService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyProfitDetailsService) SetCursor(cursor string) *GetCopyProfitDetailsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyProfitDetailsService) Do(ctx context.Context) (*CopyProfitDetails, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/profit-details", s.params).WithSign()
	return request.Do[CopyProfitDetails](req)
}

type CopyProfitDetails struct {
	List       []CopyProfitDetail `json:"list"`
	NextCursor string             `json:"nextCursor"`
}

type CopyProfitDetail struct {
	FollowerName string          `json:"followerName"`
	Profit       decimal.Decimal `json:"profit"`       // follower PnL
	AllocatedPnL decimal.Decimal `json:"allocatedPnl"` // settled PnL
	PendingPnL   decimal.Decimal `json:"pendingPnl"`   // unsettled PnL
	ShareRatio   decimal.Decimal `json:"shareRatio"`   // profit-sharing ratio, decimal form
	ShareProfit  decimal.Decimal `json:"shareProfit"`
	Reason       string          `json:"reason"` // period, unfollow
	SettleTime   time.Time       `json:"settleTime"`
}
