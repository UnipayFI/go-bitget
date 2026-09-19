package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CreateCopyService -- POST /api/v3/copy/futures/follower-settings (UTA trade read & write)
//
// Subscribes to a lead trading project as a follower. copyType is "fixed_ratio"
// (orders sized by the lead trader's margin-to-equity ratio) or "fixed_margin"
// (a fixed margin per order, set via SetMarginPerOrder); amount is the USDT
// copy amount (minimum 50). The reply data is the literal string "success".
type CreateCopyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateCopyService(projectID, copyType string, amount decimal.Decimal) *CreateCopyService {
	return &CreateCopyService{c: c, body: map[string]any{
		"projectId": projectID,
		"type":      copyType,
		"amount":    amount.String(),
	}}
}

// SetAccountType restricts the source accounts auto fund transfer may draw
// from: a comma-separated list of "funding", "uta" and "otc". All three are used
// when omitted, deducted in the order funding -> otc -> uta.
func (s *CreateCopyService) SetAccountType(accountType string) *CreateCopyService {
	s.body["accountType"] = accountType
	return s
}

// SetTradingPairList limits copying to the given symbols (all pairs when
// omitted); lead trades outside the list are not copied.
func (s *CreateCopyService) SetTradingPairList(tradingPairList []string) *CreateCopyService {
	s.body["tradingPairList"] = tradingPairList
	return s
}

// SetMarginPerOrder sets the USDT margin per order (required for fixed_margin).
func (s *CreateCopyService) SetMarginPerOrder(marginPerOrder decimal.Decimal) *CreateCopyService {
	s.body["marginPerOrder"] = marginPerOrder.String()
	return s
}

// SetAutoCopy sets whether new trading pairs are followed automatically ("on",
// the default, or "off").
func (s *CreateCopyService) SetAutoCopy(autoCopy string) *CreateCopyService {
	s.body["autoCopy"] = autoCopy
	return s
}

// SetLeverage sets the leverage in [1, 10] (follows the lead trader's leverage
// when omitted).
func (s *CreateCopyService) SetLeverage(leverage string) *CreateCopyService {
	s.body["leverage"] = leverage
	return s
}

// SetMaxEntrySlippage sets the maximum entry slippage in percent, in [0.1, 3]
// (0.1 means 0.1%).
func (s *CreateCopyService) SetMaxEntrySlippage(maxEntrySlippage decimal.Decimal) *CreateCopyService {
	s.body["maxEntrySlippage"] = maxEntrySlippage.String()
	return s
}

// SetMaxMarginRatio caps a single order's margin as a percentage of total
// assets, in [5, 95] (5 means 5%).
func (s *CreateCopyService) SetMaxMarginRatio(maxMarginRatio decimal.Decimal) *CreateCopyService {
	s.body["maxMarginRatio"] = maxMarginRatio.String()
	return s
}

// SetMaxPositionValue caps the total USDT position value (default and maximum
// 2,000,000); no new copy orders are placed once it is reached. The wire field
// is spelled maxPostionValue.
func (s *CreateCopyService) SetMaxPositionValue(maxPositionValue decimal.Decimal) *CreateCopyService {
	s.body["maxPostionValue"] = maxPositionValue.String()
	return s
}

func (s *CreateCopyService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/follower-settings", s.body).WithSign()
	return request.Do[string](req)
}

// ModifyFollowerSettingsService -- POST /api/v3/copy/futures/modify-follower-settings (UTA trade read & write)
//
// Amends an existing copy trading subscription. The reply data is the literal
// string "success".
type ModifyFollowerSettingsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyFollowerSettingsService(projectID string) *ModifyFollowerSettingsService {
	return &ModifyFollowerSettingsService{c: c, body: map[string]any{"projectId": projectID}}
}

// SetTradingPairList sets the symbols to copy; pass []string{"all"} to follow
// every pair.
func (s *ModifyFollowerSettingsService) SetTradingPairList(tradingPairList []string) *ModifyFollowerSettingsService {
	s.body["tradingPairList"] = tradingPairList
	return s
}

// SetMarginPerOrder sets the USDT margin per order.
func (s *ModifyFollowerSettingsService) SetMarginPerOrder(marginPerOrder decimal.Decimal) *ModifyFollowerSettingsService {
	s.body["marginPerOrder"] = marginPerOrder.String()
	return s
}

// SetAutoCopy sets whether new trading pairs are followed automatically ("on"
// or "off").
func (s *ModifyFollowerSettingsService) SetAutoCopy(autoCopy string) *ModifyFollowerSettingsService {
	s.body["autoCopy"] = autoCopy
	return s
}

// SetLeverage sets the leverage in [1, 10]; only allowed with no open positions.
func (s *ModifyFollowerSettingsService) SetLeverage(leverage string) *ModifyFollowerSettingsService {
	s.body["leverage"] = leverage
	return s
}

// SetMaxEntrySlippage sets the maximum entry slippage in percent, in [0.1, 3]
// (0.1 means 0.1%).
func (s *ModifyFollowerSettingsService) SetMaxEntrySlippage(maxEntrySlippage decimal.Decimal) *ModifyFollowerSettingsService {
	s.body["maxEntrySlippage"] = maxEntrySlippage.String()
	return s
}

// SetMaxMarginRatio caps a single order's margin as a percentage of total
// assets, in [5, 95].
func (s *ModifyFollowerSettingsService) SetMaxMarginRatio(maxMarginRatio decimal.Decimal) *ModifyFollowerSettingsService {
	s.body["maxMarginRatio"] = maxMarginRatio.String()
	return s
}

// SetMaxPositionValue caps the total USDT position value (maximum 2,000,000).
// The wire field is spelled maxPostionValue.
func (s *ModifyFollowerSettingsService) SetMaxPositionValue(maxPositionValue decimal.Decimal) *ModifyFollowerSettingsService {
	s.body["maxPostionValue"] = maxPositionValue.String()
	return s
}

func (s *ModifyFollowerSettingsService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/modify-follower-settings", s.body).WithSign()
	return request.Do[string](req)
}

// CopyUnfollowService -- POST /api/v3/copy/futures/unfollow (UTA trade read & write)
//
// Cancels a copy trading subscription. The reply data is the literal string
// "success".
type CopyUnfollowService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCopyUnfollowService(projectID string) *CopyUnfollowService {
	return &CopyUnfollowService{c: c, body: map[string]any{"projectId": projectID}}
}

// SetCloseType sets how open positions are handled: "follow_close" (close
// alongside the lead trader) or "instant_close" (close immediately).
func (s *CopyUnfollowService) SetCloseType(closeType string) *CopyUnfollowService {
	s.body["closeType"] = closeType
	return s
}

func (s *CopyUnfollowService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/unfollow", s.body).WithSign()
	return request.Do[string](req)
}

// GetCopySettingsService -- GET /api/v3/copy/futures/copy-settings (UTA trade read)
//
// Returns the follower's copy trading settings for a project.
type GetCopySettingsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopySettingsService(projectID string) *GetCopySettingsService {
	return &GetCopySettingsService{c: c, params: map[string]string{"projectId": projectID}}
}

func (s *GetCopySettingsService) Do(ctx context.Context) (*CopySettings, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/copy-settings", s.params).WithSign()
	return request.Do[CopySettings](req)
}

type CopySettings struct {
	Type            string          `json:"type"`   // fixed_ratio, fixed_margin
	Amount          decimal.Decimal `json:"amount"` // copy amount, USDT
	TradingPairList []string        `json:"tradingPairList"`
	MarginPerOrder  decimal.Decimal `json:"marginPerOrder"` // USDT
	AutoCopy        string          `json:"autoCopy"`       // on, off
	// Leverage is empty when following the lead trader's leverage.
	Leverage         string          `json:"leverage"`
	MaxEntrySlippage decimal.Decimal `json:"maxEntrySlippage"` // percent, 1 means 1%
	MaxMarginRatio   decimal.Decimal `json:"maxMarginRatio"`   // percent, 90 means 90%
	MaxPositionValue decimal.Decimal `json:"maxPostionValue"`  // USDT; wire field is misspelled
}

// CopyFollowerTransferService -- POST /api/v3/copy/futures/copy-transfer (UTA mgt. read & write)
//
// Moves a coin into or out of the follower's copy trading account for a
// project. type is "in" or "out" (out only to the spot/funding account). The
// reply data is the literal string "success".
type CopyFollowerTransferService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCopyFollowerTransferService(projectID, transferType, coin string, amount decimal.Decimal) *CopyFollowerTransferService {
	return &CopyFollowerTransferService{c: c, body: map[string]any{
		"projectId": projectID,
		"type":      transferType,
		"coin":      coin,
		"amount":    amount.String(),
	}}
}

// SetInAccountType restricts the source accounts a transfer-in may draw from:
// a comma-separated list of "funding", "uta" and "otc". All three are used when
// omitted, deducted in the order funding -> otc -> uta.
func (s *CopyFollowerTransferService) SetInAccountType(inAccountType string) *CopyFollowerTransferService {
	s.body["inAccountType"] = inAccountType
	return s
}

func (s *CopyFollowerTransferService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/copy-transfer", s.body).WithSign()
	return request.Do[string](req)
}

// GetCopyFollowerTransferRecordService -- GET /api/v3/copy/futures/copy-transfer-record (UTA mgt. read)
//
// Returns the follower's copy trading account transfer records for a project,
// paginated by cursor (pass the previous response's nextCursor).
type GetCopyFollowerTransferRecordService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyFollowerTransferRecordService(projectID string) *GetCopyFollowerTransferRecordService {
	return &GetCopyFollowerTransferRecordService{c: c, params: map[string]string{"projectId": projectID}}
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetCopyFollowerTransferRecordService) SetLimit(limit string) *GetCopyFollowerTransferRecordService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyFollowerTransferRecordService) SetCursor(cursor string) *GetCopyFollowerTransferRecordService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyFollowerTransferRecordService) Do(ctx context.Context) (*CopyFollowerTransferRecords, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/copy-transfer-record", s.params).WithSign()
	return request.Do[CopyFollowerTransferRecords](req)
}

type CopyFollowerTransferRecords struct {
	List       []CopyFollowerTransferRecord `json:"list"`
	NextCursor string                       `json:"nextCursor"`
}

type CopyFollowerTransferRecord struct {
	FromType    string          `json:"fromType"` // funding, uta, otc, copy
	ToType      string          `json:"toType"`   // funding, uta, otc, copy
	Amount      decimal.Decimal `json:"amount"`
	Coin        string          `json:"coin"`
	Status      string          `json:"status"` // successful, failed, processing
	CreatedTime time.Time       `json:"createdTime"`
}

// GetCurrentCopyService -- GET /api/v3/copy/futures/current-copy (UTA trade read)
//
// Returns the follower's current copy trading summary for a project.
type GetCurrentCopyService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCurrentCopyService(projectID string) *GetCurrentCopyService {
	return &GetCurrentCopyService{c: c, params: map[string]string{"projectId": projectID}}
}

func (s *GetCurrentCopyService) Do(ctx context.Context) (*CurrentCopy, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/current-copy", s.params).WithSign()
	return request.Do[CurrentCopy](req)
}

// CurrentCopy amounts are all in USDT.
type CurrentCopy struct {
	EliteTrader       string          `json:"eliteTrader"` // lead trader name
	EstNetProfit      decimal.Decimal `json:"estNetProfit"`
	ProfitShare       decimal.Decimal `json:"profitShare"`
	EstValue          decimal.Decimal `json:"estValue"` // incl. unrealized PnL
	Available         decimal.Decimal `json:"available"`
	CurrentInvestment decimal.Decimal `json:"currentInvestment"` // cumulative
}

// GetCopyFollowerProfitDetailsService -- GET /api/v3/copy/futures/copy-profit-details (UTA trade read)
//
// Returns the follower's profit-sharing settlements for a project, paginated
// by cursor (pass the previous response's nextCursor).
type GetCopyFollowerProfitDetailsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyFollowerProfitDetailsService(projectID string) *GetCopyFollowerProfitDetailsService {
	return &GetCopyFollowerProfitDetailsService{c: c, params: map[string]string{"projectId": projectID}}
}

// SetStartTime filters settlements at or after t.
func (s *GetCopyFollowerProfitDetailsService) SetStartTime(t time.Time) *GetCopyFollowerProfitDetailsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters settlements at or before t.
func (s *GetCopyFollowerProfitDetailsService) SetEndTime(t time.Time) *GetCopyFollowerProfitDetailsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetCopyFollowerProfitDetailsService) SetLimit(limit string) *GetCopyFollowerProfitDetailsService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyFollowerProfitDetailsService) SetCursor(cursor string) *GetCopyFollowerProfitDetailsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyFollowerProfitDetailsService) Do(ctx context.Context) (*CopyFollowerProfitDetails, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/copy-profit-details", s.params).WithSign()
	return request.Do[CopyFollowerProfitDetails](req)
}

type CopyFollowerProfitDetails struct {
	List       []CopyFollowerProfitDetail `json:"list"`
	NextCursor string                     `json:"nextCursor"`
}

type CopyFollowerProfitDetail struct {
	SettleTime   time.Time       `json:"settleTime"`
	Profit       decimal.Decimal `json:"profit"`       // follower PnL
	AllocatedPnL decimal.Decimal `json:"allocatedPnl"` // settled PnL
	PendingPnL   decimal.Decimal `json:"pendingPnl"`   // unsettled PnL
	ShareRatio   decimal.Decimal `json:"shareRatio"`   // profit-sharing ratio, decimal form
	ShareProfit  decimal.Decimal `json:"shareProfit"`
}

// CopyClosePositionsService -- POST /api/v3/copy/futures/close-positions (UTA trade read & write)
//
// Manually closes qty of a position in the follower's copy trading account.
// The reply data is the literal string "success".
type CopyClosePositionsService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCopyClosePositionsService(projectID, symbol string, qty decimal.Decimal) *CopyClosePositionsService {
	return &CopyClosePositionsService{c: c, body: map[string]any{
		"projectId": projectID,
		"symbol":    symbol,
		"qty":       qty.String(),
	}}
}

// SetHoldSide sets the position side to close (required in hedge mode).
func (s *CopyClosePositionsService) SetHoldSide(holdSide PosSide) *CopyClosePositionsService {
	s.body["holdSide"] = string(holdSide)
	return s
}

func (s *CopyClosePositionsService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/close-positions", s.body).WithSign()
	return request.Do[string](req)
}

// CopyCloseAllService -- POST /api/v3/copy/futures/close-all (UTA trade read & write)
//
// Closes every position in the follower's copy trading account for a project.
// The reply data is the literal string "success".
type CopyCloseAllService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCopyCloseAllService(projectID string) *CopyCloseAllService {
	return &CopyCloseAllService{c: c, body: map[string]any{"projectId": projectID}}
}

func (s *CopyCloseAllService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/close-all", s.body).WithSign()
	return request.Do[string](req)
}

// GetCopyCurrentPositionsService -- GET /api/v3/copy/futures/current-positions (UTA trade read)
//
// Returns the open positions in the follower's copy trading account for a
// project, optionally filtered by symbol and side.
type GetCopyCurrentPositionsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyCurrentPositionsService(projectID string) *GetCopyCurrentPositionsService {
	return &GetCopyCurrentPositionsService{c: c, params: map[string]string{"projectId": projectID}}
}

func (s *GetCopyCurrentPositionsService) SetSymbol(symbol string) *GetCopyCurrentPositionsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetCopyCurrentPositionsService) SetPosSide(posSide PosSide) *GetCopyCurrentPositionsService {
	s.params["posSide"] = string(posSide)
	return s
}

func (s *GetCopyCurrentPositionsService) Do(ctx context.Context) (*CopyCurrentPositions, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/current-positions", s.params).WithSign()
	return request.Do[CopyCurrentPositions](req)
}

type CopyCurrentPositions struct {
	List []CopyCurrentPosition `json:"list"`
}

type CopyCurrentPosition struct {
	Symbol      string          `json:"symbol"`
	MarginCoin  string          `json:"marginCoin"`
	PosSide     PosSide         `json:"posSide"`
	Total       decimal.Decimal `json:"total"` // position quantity
	Leverage    string          `json:"leverage"`
	AvgPrice    decimal.Decimal `json:"avgPrice"`
	MarginMode  MarginMode      `json:"marginMode"`
	HoldMode    HoldMode        `json:"holdMode"`
	CreatedTime time.Time       `json:"createdTime"`
	PositionID  string          `json:"positionId"`
}

// PlaceCopyTPSLService -- POST /api/v3/copy/futures/place-tpsl (UTA trade read & write)
//
// Sets take-profit and stop-loss on a position in the follower's copy trading
// account. The trigger types are market or mark.
type PlaceCopyTPSLService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewPlaceCopyTPSLService(projectID, positionID string, tpTriggerBy, slTriggerBy TriggerBy, takeProfit, stopLoss decimal.Decimal) *PlaceCopyTPSLService {
	return &PlaceCopyTPSLService{c: c, body: map[string]any{
		"projectId":   projectID,
		"positionId":  positionID,
		"tpTriggerBy": string(tpTriggerBy),
		"slTriggerBy": string(slTriggerBy),
		"takeProfit":  takeProfit.String(),
		"stopLoss":    stopLoss.String(),
	}}
}

func (s *PlaceCopyTPSLService) Do(ctx context.Context) (*CopyTPSLResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/place-tpsl", s.body).WithSign()
	return request.Do[CopyTPSLResult](req)
}

type CopyTPSLResult struct {
	StrategyID string `json:"strategyId"`
}

// ModifyCopyTPSLService -- POST /api/v3/copy/futures/modify-tpsl (UTA trade read & write)
//
// Amends a copy trading position's take-profit/stop-loss strategy.
type ModifyCopyTPSLService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyCopyTPSLService(projectID, strategyID string, tpTriggerBy, slTriggerBy TriggerBy, takeProfit, stopLoss decimal.Decimal) *ModifyCopyTPSLService {
	return &ModifyCopyTPSLService{c: c, body: map[string]any{
		"projectId":   projectID,
		"strategyId":  strategyID,
		"tpTriggerBy": string(tpTriggerBy),
		"slTriggerBy": string(slTriggerBy),
		"takeProfit":  takeProfit.String(),
		"stopLoss":    stopLoss.String(),
	}}
}

func (s *ModifyCopyTPSLService) Do(ctx context.Context) (*CopyTPSLResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/modify-tpsl", s.body).WithSign()
	return request.Do[CopyTPSLResult](req)
}

// CancelCopyTPSLService -- POST /api/v3/copy/futures/cancel-tpsl (UTA trade read & write)
//
// Cancels a copy trading position's take-profit/stop-loss strategy. The reply
// data is the literal string "success".
type CancelCopyTPSLService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCancelCopyTPSLService(projectID, strategyID string) *CancelCopyTPSLService {
	return &CancelCopyTPSLService{c: c, body: map[string]any{
		"projectId":  projectID,
		"strategyId": strategyID,
	}}
}

func (s *CancelCopyTPSLService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v3/copy/futures/cancel-tpsl", s.body).WithSign()
	return request.Do[string](req)
}

// GetCopyCurrentTPSLOrdersService -- GET /api/v3/copy/futures/current-tpsl-orders (UTA trade read)
//
// Returns the follower's active take-profit/stop-loss orders for a project.
type GetCopyCurrentTPSLOrdersService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyCurrentTPSLOrdersService(projectID string) *GetCopyCurrentTPSLOrdersService {
	return &GetCopyCurrentTPSLOrdersService{c: c, params: map[string]string{"projectId": projectID}}
}

func (s *GetCopyCurrentTPSLOrdersService) Do(ctx context.Context) (*CopyTPSLOrders, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/current-tpsl-orders", s.params).WithSign()
	return request.Do[CopyTPSLOrders](req)
}

type CopyTPSLOrders struct {
	List []CopyTPSLOrder `json:"list"`
}

// CopyTPSLOrder is a copy trading take-profit/stop-loss order, shared by the
// current and history listings.
type CopyTPSLOrder struct {
	StrategyID  string          `json:"strategyId"`
	Category    Category        `json:"category"`
	Symbol      string          `json:"symbol"`
	Qty         decimal.Decimal `json:"qty"`
	PosSide     PosSide         `json:"posSide"`
	Status      string          `json:"status"` // pending, success, failed, cancelled, submitting
	TpTriggerBy TriggerBy       `json:"tpTriggerBy"`
	SlTriggerBy TriggerBy       `json:"slTriggerBy"`
	TakeProfit  decimal.Decimal `json:"takeProfit"`
	StopLoss    decimal.Decimal `json:"stopLoss"`
	TpOrderType OrderType       `json:"tpOrderType"`
	SlOrderType OrderType       `json:"slOrderType"`
}

// GetCopyTPSLOrderHistoryService -- GET /api/v3/copy/futures/tpsl-order-history (UTA trade read)
//
// Returns the follower's finished take-profit/stop-loss orders for a project
// (last 30 days by default, 90-day max window), paginated by cursor (pass the
// previous response's cursor).
type GetCopyTPSLOrderHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCopyTPSLOrderHistoryService(projectID string) *GetCopyTPSLOrderHistoryService {
	return &GetCopyTPSLOrderHistoryService{c: c, params: map[string]string{"projectId": projectID}}
}

// SetStartTime filters orders at or after t.
func (s *GetCopyTPSLOrderHistoryService) SetStartTime(t time.Time) *GetCopyTPSLOrderHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters orders at or before t.
func (s *GetCopyTPSLOrderHistoryService) SetEndTime(t time.Time) *GetCopyTPSLOrderHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default and max 100).
func (s *GetCopyTPSLOrderHistoryService) SetLimit(limit string) *GetCopyTPSLOrderHistoryService {
	s.params["limit"] = limit
	return s
}

func (s *GetCopyTPSLOrderHistoryService) SetCursor(cursor string) *GetCopyTPSLOrderHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetCopyTPSLOrderHistoryService) Do(ctx context.Context) (*CopyTPSLOrderHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/copy/futures/tpsl-order-history", s.params).WithSign()
	return request.Do[CopyTPSLOrderHistory](req)
}

type CopyTPSLOrderHistory struct {
	List   []CopyTPSLOrder `json:"list"`
	Cursor string          `json:"cursor"`
}
