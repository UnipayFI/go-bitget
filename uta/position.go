package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetPositionService -- GET /api/v3/position/current-position (UTA trade read)
//
// Returns the account's open futures positions for a product category,
// optionally filtered to a single symbol and/or position side.
type GetPositionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetPositionService(category Category) *GetPositionService {
	return &GetPositionService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetPositionService) SetSymbol(symbol string) *GetPositionService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetPositionService) SetPosSide(posSide PosSide) *GetPositionService {
	s.params["posSide"] = string(posSide)
	return s
}

func (s *GetPositionService) Do(ctx context.Context) ([]Position, error) {
	req := request.Get(ctx, s.c, "/api/v3/position/current-position", s.params).WithSign()
	resp, err := request.Do[positionList](req)
	if err != nil {
		return nil, err
	}
	return resp.List, nil
}

// positionList wraps the current-position payload, which is an object
// {"list":[...]} rather than a bare array.
type positionList struct {
	List []Position `json:"list"`
}

// Position is a single open futures position.
type Position struct {
	Category         Category        `json:"category"`
	Symbol           string          `json:"symbol"`
	MarginCoin       string          `json:"marginCoin"`
	PosSide          PosSide         `json:"posSide"`
	PositionBalance  decimal.Decimal `json:"positionBalance"`
	Available        decimal.Decimal `json:"available"`
	Frozen           decimal.Decimal `json:"frozen"`
	Total            decimal.Decimal `json:"total"`
	Leverage         decimal.Decimal `json:"leverage"`
	CurRealisedPnL   decimal.Decimal `json:"curRealisedPnl"`
	AvgPrice         decimal.Decimal `json:"avgPrice"`
	MarginMode       MarginMode      `json:"marginMode"`
	PositionStatus   string          `json:"positionStatus"` // normal
	HoldMode         HoldMode        `json:"holdMode"`
	UnrealizedPnL    decimal.Decimal `json:"unrealisedPnl"`
	LiquidationPrice decimal.Decimal `json:"liquidationPrice"`
	Mmr              decimal.Decimal `json:"mmr"`
	ProfitRate       decimal.Decimal `json:"profitRate"`
	MarkPrice        decimal.Decimal `json:"markPrice"`
	BreakEvenPrice   decimal.Decimal `json:"breakEvenPrice"`
	TotalFunding     decimal.Decimal `json:"totalFunding"`
	OpenFeeTotal     decimal.Decimal `json:"openFeeTotal"`
	CloseFeeTotal    decimal.Decimal `json:"closeFeeTotal"`
	CreatedTime      time.Time       `json:"createdTime"`
	UpdatedTime      time.Time       `json:"updatedTime"`
}

// GetPositionHistoryService -- GET /api/v3/position/history-position (UTA trade read)
//
// Returns closed/historical futures positions for a product category,
// paginated by cursor and bounded to a 90-day lookback window.
type GetPositionHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetPositionHistoryService(category Category) *GetPositionHistoryService {
	return &GetPositionHistoryService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetPositionHistoryService) SetSymbol(symbol string) *GetPositionHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters positions at or after t (90-day lookback window).
func (s *GetPositionHistoryService) SetStartTime(t time.Time) *GetPositionHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters positions at or before t (max 30-day range from startTime).
func (s *GetPositionHistoryService) SetEndTime(t time.Time) *GetPositionHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetPositionHistoryService) SetLimit(limit string) *GetPositionHistoryService {
	s.params["limit"] = limit
	return s
}

func (s *GetPositionHistoryService) SetCursor(cursor string) *GetPositionHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetPositionHistoryService) Do(ctx context.Context) (*PositionHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/position/history-position", s.params).WithSign()
	return request.Do[PositionHistory](req)
}

type PositionHistory struct {
	List   []HistoryPosition `json:"list"`
	Cursor string            `json:"cursor"`
}

// HistoryPosition is a single closed/historical futures position.
type HistoryPosition struct {
	PositionID     string          `json:"positionId"`
	Category       Category        `json:"category"`
	Symbol         string          `json:"symbol"`
	MarginCoin     string          `json:"marginCoin"`
	HoldMode       HoldMode        `json:"holdMode"`
	PosSide        PosSide         `json:"posSide"`
	MarginMode     MarginMode      `json:"marginMode"`
	OpenPriceAvg   decimal.Decimal `json:"openPriceAvg"`
	ClosePriceAvg  decimal.Decimal `json:"closePriceAvg"`
	OpenTotalPos   decimal.Decimal `json:"openTotalPos"`
	CloseTotalPos  decimal.Decimal `json:"closeTotalPos"`
	CumRealisedPnL decimal.Decimal `json:"cumRealisedPnl"`
	NetProfit      decimal.Decimal `json:"netProfit"`
	TotalFunding   decimal.Decimal `json:"totalFunding"`
	OpenFeeTotal   decimal.Decimal `json:"openFeeTotal"`
	CloseFeeTotal  decimal.Decimal `json:"closeFeeTotal"`
	CreatedTime    time.Time       `json:"createdTime"`
	UpdatedTime    time.Time       `json:"updatedTime"`
}

// GetMovePositionHistoryService -- GET /api/v3/account/move-position-history (UTA trade read)
//
// Returns the history of position-move (transfer) records for a product
// category, paginated by cursor and bounded to a 90-day lookback window.
type GetMovePositionHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetMovePositionHistoryService(category Category) *GetMovePositionHistoryService {
	return &GetMovePositionHistoryService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetMovePositionHistoryService) SetSymbol(symbol string) *GetMovePositionHistoryService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters records at or after t.
func (s *GetMovePositionHistoryService) SetStartTime(t time.Time) *GetMovePositionHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (defaults to last 30 days; max 90-day span).
func (s *GetMovePositionHistoryService) SetEndTime(t time.Time) *GetMovePositionHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMovePositionHistoryService) SetLimit(limit string) *GetMovePositionHistoryService {
	s.params["limit"] = limit
	return s
}

func (s *GetMovePositionHistoryService) SetCursor(cursor string) *GetMovePositionHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetMovePositionHistoryService) Do(ctx context.Context) (*MovePositionHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/move-position-history", s.params).WithSign()
	return request.Do[MovePositionHistory](req)
}

type MovePositionHistory struct {
	List   []MovePosition `json:"list"`
	Cursor string         `json:"cursor"`
}

// MovePosition is a single position-move (transfer) record.
type MovePosition struct {
	Category    Category        `json:"category"`
	FromUID     string          `json:"fromUid"`
	ToUID       string          `json:"toUid"`
	OrderID     string          `json:"orderId"`
	OpenExecID  string          `json:"openExecId"`
	CloseExecID string          `json:"closeExecId"`
	Symbol      string          `json:"symbol"`
	PosSide     PosSide         `json:"posSide"`
	Qty         decimal.Decimal `json:"qty"`
	Price       decimal.Decimal `json:"price"`
	Status      string          `json:"status"` // processing, completed, failed
	CreatedTime time.Time       `json:"createdTime"`
	UpdatedTime time.Time       `json:"updatedTime"`
}

// GetPositionADLRankService -- GET /api/v3/position/adlRank (UTA trade read)
//
// Returns the auto-deleveraging (ADL) rank for each open position. A rank closer
// to 1 means a higher probability of being reduced during liquidation events.
type GetPositionADLRankService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetPositionADLRankService() *GetPositionADLRankService {
	return &GetPositionADLRankService{c: c}
}

func (s *GetPositionADLRankService) Do(ctx context.Context) ([]PositionADLRank, error) {
	req := request.Get(ctx, s.c, "/api/v3/position/adlRank").WithSign()
	resp, err := request.Do[[]PositionADLRank](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PositionADLRank is a single position's ADL queue ranking.
type PositionADLRank struct {
	Symbol     string          `json:"symbol"`
	MarginCoin string          `json:"marginCoin"`
	ADLRank    decimal.Decimal `json:"adlRank"`
	HoldSide   PosSide         `json:"holdSide"`
}

// GetMaxOpenAvailableService -- POST /api/v3/account/max-open-available (UTA trade read)
//
// Returns the maximum openable/available quantities and the quote-coin cost to
// open a position for the given symbol, order type and side.
type GetMaxOpenAvailableService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewGetMaxOpenAvailableService(category Category, symbol string, orderType OrderType, side Side) *GetMaxOpenAvailableService {
	return &GetMaxOpenAvailableService{c: c, body: map[string]any{
		"category":  string(category),
		"symbol":    symbol,
		"orderType": string(orderType),
		"side":      string(side),
	}}
}

// SetPrice sets the order price (required when orderType is limit).
func (s *GetMaxOpenAvailableService) SetPrice(price decimal.Decimal) *GetMaxOpenAvailableService {
	s.body["price"] = price.String()
	return s
}

// SetSize sets the order quantity in base coin.
func (s *GetMaxOpenAvailableService) SetSize(size decimal.Decimal) *GetMaxOpenAvailableService {
	s.body["size"] = size.String()
	return s
}

func (s *GetMaxOpenAvailableService) Do(ctx context.Context) (*MaxOpenAvailable, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/max-open-available", s.body).WithSign()
	return request.Do[MaxOpenAvailable](req)
}

type MaxOpenAvailable struct {
	Available        decimal.Decimal `json:"available"`
	MaxOpen          decimal.Decimal `json:"maxOpen"`
	BuyOpenCost      decimal.Decimal `json:"buyOpenCost"`
	SellOpenCost     decimal.Decimal `json:"sellOpenCost"`
	MaxBuyOpen       decimal.Decimal `json:"maxBuyOpen"`
	MaxSellOpen      decimal.Decimal `json:"maxSellOpen"`
	MaxBuyAvailable  decimal.Decimal `json:"maxBuyAvailable"`
	MaxSellAvailable decimal.Decimal `json:"maxSellAvailable"`
}
