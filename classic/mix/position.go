package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSinglePositionService -- GET /api/v2/mix/position/single-position (private)
//
// Returns the open position(s) for a single symbol; in hedge mode both the long
// and short legs are returned.
type GetSinglePositionService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetSinglePositionService(productType ProductType, symbol, marginCoin string) *GetSinglePositionService {
	return &GetSinglePositionService{c: c, params: map[string]string{
		"productType": string(productType),
		"symbol":      symbol,
		"marginCoin":  marginCoin,
	}}
}

func (s *GetSinglePositionService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/position/single-position", s.params).WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetAllPositionService -- GET /api/v2/mix/position/all-position (private)
//
// Returns every open position for a product line, optionally filtered to a
// single margin coin.
type GetAllPositionService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetAllPositionService(productType ProductType) *GetAllPositionService {
	return &GetAllPositionService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetAllPositionService) SetMarginCoin(marginCoin string) *GetAllPositionService {
	s.params["marginCoin"] = marginCoin
	return s
}

func (s *GetAllPositionService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/position/all-position", s.params).WithSign()
	resp, err := request.Do[[]Position](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Position is a single open futures position. It is the union of the
// single-position and all-position response shapes (single-position adds
// liquidationPrice/marginRatio context fields that all-position also returns).
type Position struct {
	Symbol           string          `json:"symbol"`
	MarginCoin       string          `json:"marginCoin"`
	HoldSide         HoldSide        `json:"holdSide"`
	OpenDelegateSize decimal.Decimal `json:"openDelegateSize"`
	MarginSize       decimal.Decimal `json:"marginSize"`
	Available        decimal.Decimal `json:"available"`
	Locked           decimal.Decimal `json:"locked"`
	Total            decimal.Decimal `json:"total"`
	Leverage         string          `json:"leverage"`
	AchievedProfits  decimal.Decimal `json:"achievedProfits"`
	OpenPriceAvg     decimal.Decimal `json:"openPriceAvg"`
	MarginMode       MarginMode      `json:"marginMode"`
	PosMode          PositionMode    `json:"posMode"`
	UnrealizedPL     decimal.Decimal `json:"unrealizedPL"`
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"`
	KeepMarginRate   decimal.Decimal `json:"keepMarginRate"`
	MarkPrice        decimal.Decimal `json:"markPrice"`
	MarginRatio      decimal.Decimal `json:"marginRatio"`
	BreakEvenPrice   decimal.Decimal `json:"breakEvenPrice"`
	TotalFee         decimal.Decimal `json:"totalFee"`
	DeductedFee      decimal.Decimal `json:"deductedFee"`
	TakeProfit       decimal.Decimal `json:"takeProfit"`
	StopLoss         decimal.Decimal `json:"stopLoss"`
	TakeProfitID     string          `json:"takeProfitId"`
	StopLossID       string          `json:"stopLossId"`
	AssetMode        AssetMode       `json:"assetMode"`
	AutoMargin       AutoMargin      `json:"autoMargin"`
	Grant            decimal.Decimal `json:"grant"`
	CTime            time.Time       `json:"cTime"`
	UTime            time.Time       `json:"uTime"`
}

// GetPositionADLRankService -- GET /api/v2/mix/position/adlRank (private)
//
// Returns the auto-deleveraging (ADL) queue rank of the account's open
// positions for a product line.
type GetPositionADLRankService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetPositionADLRankService(productType ProductType) *GetPositionADLRankService {
	return &GetPositionADLRankService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetPositionADLRankService) Do(ctx context.Context) ([]PositionADLRank, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/position/adlRank", s.params).WithSign()
	resp, err := request.Do[[]PositionADLRank](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PositionADLRank is one position's ADL queue rank. A rank closer to 1 means the
// position is more likely to be auto-deleveraged.
type PositionADLRank struct {
	Symbol     string   `json:"symbol"`
	MarginCoin string   `json:"marginCoin"`
	ADLRank    string   `json:"adlRank"` // deprecated; use Rank
	Rank       string   `json:"rank"`
	HoldSide   HoldSide `json:"holdSide"`
}

// GetHistoryPositionService -- GET /api/v2/mix/position/history-position (private)
//
// Returns closed (historical) positions, paged most-recent first. Without a time
// window the most recent 3 months are returned.
type GetHistoryPositionService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetHistoryPositionService(productType ProductType) *GetHistoryPositionService {
	return &GetHistoryPositionService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetHistoryPositionService) SetSymbol(symbol string) *GetHistoryPositionService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetHistoryPositionService) SetIDLessThan(idLessThan string) *GetHistoryPositionService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetHistoryPositionService) SetStartTime(t time.Time) *GetHistoryPositionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetHistoryPositionService) SetEndTime(t time.Time) *GetHistoryPositionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetHistoryPositionService) SetLimit(limit int) *GetHistoryPositionService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetHistoryPositionService) Do(ctx context.Context) (*HistoryPositionResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/position/history-position", s.params).WithSign()
	return request.Do[HistoryPositionResponse](req)
}

// HistoryPositionResponse is the paged historical-position payload.
type HistoryPositionResponse struct {
	List  []HistoryPosition `json:"list"`
	EndID string            `json:"endId"`
}

// HistoryPosition is a single closed position record.
type HistoryPosition struct {
	PositionID    string          `json:"positionId"`
	Symbol        string          `json:"symbol"`
	MarginCoin    string          `json:"marginCoin"`
	HoldSide      HoldSide        `json:"holdSide"`
	PosMode       PositionMode    `json:"posMode"`
	OpenAvgPrice  decimal.Decimal `json:"openAvgPrice"`
	CloseAvgPrice decimal.Decimal `json:"closeAvgPrice"`
	MarginMode    MarginMode      `json:"marginMode"`
	OpenTotalPos  decimal.Decimal `json:"openTotalPos"`
	CloseTotalPos decimal.Decimal `json:"closeTotalPos"`
	Pnl           decimal.Decimal `json:"pnl"`
	NetProfit     decimal.Decimal `json:"netProfit"`
	TotalFunding  decimal.Decimal `json:"totalFunding"`
	OpenFee       decimal.Decimal `json:"openFee"`
	CloseFee      decimal.Decimal `json:"closeFee"`
	CTime         time.Time       `json:"ctime"`
	UTime         time.Time       `json:"utime"`
}
