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
	TransferID  string          `json:"transferId"`
	FromType    string          `json:"fromType"` // spot, uta, lead
	ToType      string          `json:"toType"`   // spot, uta, lead
	Coin        string          `json:"coin"`
	Amount      decimal.Decimal `json:"amount"`
	Status      string          `json:"status"` // Successful, Failed, Processing
	CreatedTime time.Time       `json:"createdTime"`
}
