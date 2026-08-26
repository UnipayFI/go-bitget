package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GridType is the direction of a spot/futures grid strategy. It is required
// for futures grids.
type GridType string

const (
	GridTypeLong  GridType = "long"
	GridTypeShort GridType = "short"
	// GridTypeNeutral is spelled "netural" on the wire; the typo is Bitget's.
	GridTypeNeutral GridType = "netural"
)

// GridOrderMode is how the grid levels are spaced between the lower and upper
// limit price.
type GridOrderMode string

const (
	GridOrderModeArithmetic GridOrderMode = "arithmetic"
	GridOrderModeGeometric  GridOrderMode = "geometric"
)

// GridFundsSource is an account a grid strategy draws its investment from.
type GridFundsSource string

const (
	GridFundsSourceFunding GridFundsSource = "funding"
	GridFundsSourceUTA     GridFundsSource = "uta"
	GridFundsSourceOTC     GridFundsSource = "otc"
)

// GridTriggerCondition decides when a created grid bot starts placing orders.
type GridTriggerCondition string

const (
	GridTriggerConditionInstant GridTriggerCondition = "instant"
	GridTriggerConditionPrice   GridTriggerCondition = "price"
	GridTriggerConditionRSI     GridTriggerCondition = "rsi"
	GridTriggerConditionBoll    GridTriggerCondition = "boll"
)

// GridTerminationCondition decides when a running grid bot stops by itself.
type GridTerminationCondition string

const (
	GridTerminationConditionRSI  GridTerminationCondition = "rsi"
	GridTerminationConditionBoll GridTerminationCondition = "boll"
)

// GridBotStatus is the lifecycle state of a grid bot.
type GridBotStatus string

const (
	GridBotStatusInit        GridBotStatus = "init"
	GridBotStatusWaiting     GridBotStatus = "waiting"
	GridBotStatusRunning     GridBotStatus = "running"
	GridBotStatusTerminating GridBotStatus = "terminating"
	GridBotStatusTerminated  GridBotStatus = "terminated"
)

// GridAdjustType is the direction of a futures grid investment adjustment.
type GridAdjustType string

const (
	GridAdjustTypeIncrease GridAdjustType = "increase"
	GridAdjustTypeDecrease GridAdjustType = "decrease"
)

// GridInvestment is one coin/amount pair of a grid strategy's investment.
// Futures grids only accept the quote coin (USDT or USDC).
type GridInvestment struct {
	Coin   string          `json:"coin"`
	Amount decimal.Decimal `json:"amount"`
}

// GridIndicatorParams configures the RSI or BOLL indicator behind a grid
// trigger or termination condition.
type GridIndicatorParams struct {
	IndicatorLength string          `json:"indicatorLength"`
	Threshold       decimal.Decimal `json:"threshold"`
	Multiplier      decimal.Decimal `json:"multiplier"`
	Interval        string          `json:"interval"`
}

// yesNo renders the yes/no toggles Bitget uses for the grid switches.
func yesNo(on bool) string {
	if on {
		return "yes"
	}
	return "no"
}

// ValidateGridParamsService -- POST /api/v3/trade/grid/validate (UTA trade read & write)
//
// Checks whether a spot or futures grid configuration (price range, grid
// number, investment, …) is legal without creating the strategy.
type ValidateGridParamsService struct {
	c    *UTAClient
	body map[string]any
}

// NewValidateGridParamsService starts a validation request. category, symbol,
// the price range, gridNum, gridOrderMode and investment are required;
// everything else is optional and mirrors NewCreateGridBotService.
func (c *UTAClient) NewValidateGridParamsService(category Category, symbol string, minPrice, maxPrice decimal.Decimal, gridNum string, gridOrderMode GridOrderMode, investment []GridInvestment) *ValidateGridParamsService {
	return &ValidateGridParamsService{c: c, body: map[string]any{
		"category":         string(category),
		"symbol":           symbol,
		"minPrice":         minPrice.String(),
		"maxPrice":         maxPrice.String(),
		"gridNum":          gridNum,
		"gridOrderMode":    string(gridOrderMode),
		"investmentAmount": investment,
	}}
}

// SetGridType sets the grid direction (required for futures).
func (s *ValidateGridParamsService) SetGridType(gridType GridType) *ValidateGridParamsService {
	s.body["gridType"] = string(gridType)
	return s
}

// SetLeverage sets the leverage (futures only, defaults to 1).
func (s *ValidateGridParamsService) SetLeverage(leverage string) *ValidateGridParamsService {
	s.body["leverage"] = leverage
	return s
}

// SetReservedMargin sets the margin held back from the grid (futures only).
func (s *ValidateGridParamsService) SetReservedMargin(reservedMargin decimal.Decimal) *ValidateGridParamsService {
	s.body["reservedMargin"] = reservedMargin.String()
	return s
}

// SetTriggerCondition sets the start condition (defaults to instant).
func (s *ValidateGridParamsService) SetTriggerCondition(triggerCondition GridTriggerCondition) *ValidateGridParamsService {
	s.body["triggerCondition"] = string(triggerCondition)
	return s
}

// SetTriggerParams sets the indicator parameters of an rsi/boll start condition.
func (s *ValidateGridParamsService) SetTriggerParams(triggerParams []GridIndicatorParams) *ValidateGridParamsService {
	s.body["triggerParams"] = triggerParams
	return s
}

// SetTriggerPrice sets the start price of a price start condition.
func (s *ValidateGridParamsService) SetTriggerPrice(triggerPrice decimal.Decimal) *ValidateGridParamsService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetTerminationCondition sets the self-termination condition (rsi or boll).
func (s *ValidateGridParamsService) SetTerminationCondition(terminationCondition GridTerminationCondition) *ValidateGridParamsService {
	s.body["terminationCondition"] = string(terminationCondition)
	return s
}

// SetTerminationParams sets the indicator parameters of the termination condition.
func (s *ValidateGridParamsService) SetTerminationParams(terminationParams []GridIndicatorParams) *ValidateGridParamsService {
	s.body["terminationParams"] = terminationParams
	return s
}

// SetStopLoss sets the stop-loss price of the strategy.
func (s *ValidateGridParamsService) SetStopLoss(stopLoss decimal.Decimal) *ValidateGridParamsService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTakeProfit sets the take-profit price of the strategy.
func (s *ValidateGridParamsService) SetTakeProfit(takeProfit decimal.Decimal) *ValidateGridParamsService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetTrailingGrid enables the trailing grid (defaults to off).
func (s *ValidateGridParamsService) SetTrailingGrid(trailingGrid bool) *ValidateGridParamsService {
	s.body["trailingGrid"] = yesNo(trailingGrid)
	return s
}

// SetMovingAverageGains sets the moving-average gain of a trailing grid.
func (s *ValidateGridParamsService) SetMovingAverageGains(movingAverageGains decimal.Decimal) *ValidateGridParamsService {
	s.body["movingAverageGains"] = movingAverageGains.String()
	return s
}

// SetStopUpwardPrice sets the price at which the grid stops trailing upward.
func (s *ValidateGridParamsService) SetStopUpwardPrice(stopUpwardPrice decimal.Decimal) *ValidateGridParamsService {
	s.body["stopUpwardPrice"] = stopUpwardPrice.String()
	return s
}

// SetHodlMode enables HODL mode, which keeps the base coin on termination
// (spot only, defaults to off).
func (s *ValidateGridParamsService) SetHodlMode(hodlMode bool) *ValidateGridParamsService {
	s.body["hodlMode"] = yesNo(hodlMode)
	return s
}

// SetMarketOpen controls whether the first order opens at market (futures only,
// enabled by default).
func (s *ValidateGridParamsService) SetMarketOpen(marketOpen bool) *ValidateGridParamsService {
	s.body["marketOpen"] = yesNo(marketOpen)
	return s
}

// SetLossReserve controls the loss reserve (futures only, enabled by default).
func (s *ValidateGridParamsService) SetLossReserve(lossReserve bool) *ValidateGridParamsService {
	s.body["lossReserve"] = yesNo(lossReserve)
	return s
}

// SetAutoTransferProfits controls whether realized profits are transferred out
// (defaults to off).
func (s *ValidateGridParamsService) SetAutoTransferProfits(autoTransferProfits bool) *ValidateGridParamsService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *ValidateGridParamsService) Do(ctx context.Context) (*GridValidation, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/validate", s.body).WithSign()
	return request.Do[GridValidation](req)
}

// GridValidation is the verdict returned by the two grid validate endpoints.
type GridValidation struct {
	Valid  string `json:"valid"` // yes, no
	Reason string `json:"reason"`
	// MinInvestment is the smallest investment the configuration accepts.
	MinInvestment decimal.Decimal `json:"minInvestment"`
}

// CreateGridBotService -- POST /api/v3/trade/grid/create-bot (UTA trade read & write)
//
// Creates a spot or futures grid trading bot. The reply data carries the new
// bot's ID.
type CreateGridBotService struct {
	c    *UTAClient
	body map[string]any
}

// NewCreateGridBotService starts a create request. category, symbol, the price
// range, gridNum, gridOrderMode, investment and fundsSource are required.
func (c *UTAClient) NewCreateGridBotService(category Category, symbol string, minPrice, maxPrice decimal.Decimal, gridNum string, gridOrderMode GridOrderMode, investment []GridInvestment, fundsSource []GridFundsSource) *CreateGridBotService {
	return &CreateGridBotService{c: c, body: map[string]any{
		"category":         string(category),
		"symbol":           symbol,
		"minPrice":         minPrice.String(),
		"maxPrice":         maxPrice.String(),
		"gridNum":          gridNum,
		"gridOrderMode":    string(gridOrderMode),
		"investmentAmount": investment,
		"fundsSource":      fundsSource,
	}}
}

// SetGridType sets the grid direction (required for futures).
func (s *CreateGridBotService) SetGridType(gridType GridType) *CreateGridBotService {
	s.body["gridType"] = string(gridType)
	return s
}

// SetLeverage sets the leverage (futures only, defaults to 1).
func (s *CreateGridBotService) SetLeverage(leverage string) *CreateGridBotService {
	s.body["leverage"] = leverage
	return s
}

// SetAutoReserveMargin enables automatic margin reservation.
func (s *CreateGridBotService) SetAutoReserveMargin(autoReserveMargin bool) *CreateGridBotService {
	s.body["autoReserveMargin"] = yesNo(autoReserveMargin)
	return s
}

// SetReservedMargin sets the margin held back from the grid (futures only).
func (s *CreateGridBotService) SetReservedMargin(reservedMargin decimal.Decimal) *CreateGridBotService {
	s.body["reservedMargin"] = reservedMargin.String()
	return s
}

// SetTriggerCondition sets the start condition (defaults to instant).
func (s *CreateGridBotService) SetTriggerCondition(triggerCondition GridTriggerCondition) *CreateGridBotService {
	s.body["triggerCondition"] = string(triggerCondition)
	return s
}

// SetTriggerParams sets the indicator parameters of an rsi/boll start condition.
func (s *CreateGridBotService) SetTriggerParams(triggerParams []GridIndicatorParams) *CreateGridBotService {
	s.body["triggerParams"] = triggerParams
	return s
}

// SetTriggerPrice sets the start price of a price start condition.
func (s *CreateGridBotService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateGridBotService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetTerminationCondition sets the self-termination condition (rsi or boll).
func (s *CreateGridBotService) SetTerminationCondition(terminationCondition GridTerminationCondition) *CreateGridBotService {
	s.body["terminationCondition"] = string(terminationCondition)
	return s
}

// SetTerminationParams sets the indicator parameters of the termination condition.
func (s *CreateGridBotService) SetTerminationParams(terminationParams []GridIndicatorParams) *CreateGridBotService {
	s.body["terminationParams"] = terminationParams
	return s
}

// SetTerminationSell controls whether the base coin is sold when the strategy
// terminates (spot grid only).
func (s *CreateGridBotService) SetTerminationSell(terminationSell bool) *CreateGridBotService {
	s.body["terminationSell"] = yesNo(terminationSell)
	return s
}

// SetStopLoss sets the stop-loss price of the strategy.
func (s *CreateGridBotService) SetStopLoss(stopLoss decimal.Decimal) *CreateGridBotService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTakeProfit sets the take-profit price of the strategy.
func (s *CreateGridBotService) SetTakeProfit(takeProfit decimal.Decimal) *CreateGridBotService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetSlippage sets the slippage tolerance of the grid's orders.
func (s *CreateGridBotService) SetSlippage(slippage decimal.Decimal) *CreateGridBotService {
	s.body["slippage"] = slippage.String()
	return s
}

// SetTrailingGrid enables the trailing grid (defaults to off).
func (s *CreateGridBotService) SetTrailingGrid(trailingGrid bool) *CreateGridBotService {
	s.body["trailingGrid"] = yesNo(trailingGrid)
	return s
}

// SetMovingAverageGains sets the moving-average gain of a trailing grid.
func (s *CreateGridBotService) SetMovingAverageGains(movingAverageGains decimal.Decimal) *CreateGridBotService {
	s.body["movingAverageGains"] = movingAverageGains.String()
	return s
}

// SetStopUpwardPrice sets the price at which the grid stops trailing upward.
func (s *CreateGridBotService) SetStopUpwardPrice(stopUpwardPrice decimal.Decimal) *CreateGridBotService {
	s.body["stopUpwardPrice"] = stopUpwardPrice.String()
	return s
}

// SetHodlMode enables HODL mode, which keeps the base coin on termination
// (spot only, defaults to off).
func (s *CreateGridBotService) SetHodlMode(hodlMode bool) *CreateGridBotService {
	s.body["hodlMode"] = yesNo(hodlMode)
	return s
}

// SetMarketOpen controls whether the first order opens at market (futures only,
// enabled by default).
func (s *CreateGridBotService) SetMarketOpen(marketOpen bool) *CreateGridBotService {
	s.body["marketOpen"] = yesNo(marketOpen)
	return s
}

// SetLossReserve controls the loss reserve (futures only, enabled by default).
func (s *CreateGridBotService) SetLossReserve(lossReserve bool) *CreateGridBotService {
	s.body["lossReserve"] = yesNo(lossReserve)
	return s
}

// SetAutoTransferProfits controls whether realized profits are transferred out
// (defaults to off).
func (s *CreateGridBotService) SetAutoTransferProfits(autoTransferProfits bool) *CreateGridBotService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *CreateGridBotService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/create-bot", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// GridBotResult is the bot identifier returned by the grid create/modify
// endpoints.
type GridBotResult struct {
	BotID string `json:"botId"`
}

// ModifyGridBotService -- POST /api/v3/trade/grid/modify-bot (UTA trade read & write)
//
// Modifies the basic parameters of a running spot/futures grid bot. Take-profit
// and stop-loss are cleared when omitted, so pass the current value to keep it.
type ModifyGridBotService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyGridBotService(botID string) *ModifyGridBotService {
	return &ModifyGridBotService{c: c, body: map[string]any{"botId": botID}}
}

// SetTakeProfit sets the take-profit price; omitting it clears the current one.
func (s *ModifyGridBotService) SetTakeProfit(takeProfit decimal.Decimal) *ModifyGridBotService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price; omitting it clears the current one.
func (s *ModifyGridBotService) SetStopLoss(stopLoss decimal.Decimal) *ModifyGridBotService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTerminationCondition sets the self-termination condition (rsi or boll).
func (s *ModifyGridBotService) SetTerminationCondition(terminationCondition GridTerminationCondition) *ModifyGridBotService {
	s.body["terminationCondition"] = string(terminationCondition)
	return s
}

// SetTerminationParams sets the indicator parameters of the termination condition.
func (s *ModifyGridBotService) SetTerminationParams(terminationParams []GridIndicatorParams) *ModifyGridBotService {
	s.body["terminationParams"] = terminationParams
	return s
}

// SetHodlMode toggles HODL mode (spot only).
func (s *ModifyGridBotService) SetHodlMode(hodlMode bool) *ModifyGridBotService {
	s.body["hodlMode"] = yesNo(hodlMode)
	return s
}

// SetAutoTransferProfits toggles transferring realized profits out.
func (s *ModifyGridBotService) SetAutoTransferProfits(autoTransferProfits bool) *ModifyGridBotService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *ModifyGridBotService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/modify-bot", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// ModifyGridIntervalService -- POST /api/v3/trade/grid/modify-grid-interval (UTA trade read & write)
//
// Moves the price range and grid number of a running spot/futures grid bot.
type ModifyGridIntervalService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyGridIntervalService(category Category, botID string, minPrice, maxPrice decimal.Decimal, gridNum string) *ModifyGridIntervalService {
	return &ModifyGridIntervalService{c: c, body: map[string]any{
		"category": string(category),
		"botId":    botID,
		"minPrice": minPrice.String(),
		"maxPrice": maxPrice.String(),
		"gridNum":  gridNum,
	}}
}

func (s *ModifyGridIntervalService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/modify-grid-interval", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// AddGridInvestmentService -- POST /api/v3/trade/grid/add-investment (UTA trade read & write)
//
// Adds (or, for futures, withdraws) investment on a running grid bot. The
// business line must match the bot ID.
type AddGridInvestmentService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewAddGridInvestmentService(category Category, botID, coin string, size decimal.Decimal, fundsSource []GridFundsSource) *AddGridInvestmentService {
	return &AddGridInvestmentService{c: c, body: map[string]any{
		"category":    string(category),
		"botId":       botID,
		"coin":        coin,
		"size":        size.String(),
		"fundsSource": fundsSource,
	}}
}

// SetAdjustType sets the adjustment direction (futures grid only, where it is
// required).
func (s *AddGridInvestmentService) SetAdjustType(adjustType GridAdjustType) *AddGridInvestmentService {
	s.body["adjustType"] = string(adjustType)
	return s
}

// SetReinvestProfit controls whether profits are reinvested (spot grid only,
// enabled by default).
func (s *AddGridInvestmentService) SetReinvestProfit(reinvestProfit bool) *AddGridInvestmentService {
	s.body["reinvestProfit"] = yesNo(reinvestProfit)
	return s
}

func (s *AddGridInvestmentService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/add-investment", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// CloseGridBotService -- POST /api/v3/trade/grid/close-bot (UTA trade read & write)
//
// Terminates a bot: order placement stops and the unfilled orders and positions
// are handled according to the strategy's configuration. The reply data is null
// on success.
type CloseGridBotService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCloseGridBotService(botID string) *CloseGridBotService {
	return &CloseGridBotService{c: c, body: map[string]any{"botId": botID}}
}

func (s *CloseGridBotService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/close-bot", s.body).WithSign()
	return request.Do[any](req)
}

// GetGridBotDetailService -- GET /api/v3/trade/grid/bot-detail (UTA trade read)
//
// Returns the full configuration and running profit of one spot/futures grid
// strategy.
type GetGridBotDetailService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetGridBotDetailService(botID string) *GetGridBotDetailService {
	return &GetGridBotDetailService{c: c, params: map[string]string{"botId": botID}}
}

func (s *GetGridBotDetailService) Do(ctx context.Context) (*GridBotDetail, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/grid/bot-detail", s.params).WithSign()
	return request.Do[GridBotDetail](req)
}

// GridBotDetail is the configuration and running state of a spot/futures grid
// strategy.
type GridBotDetail struct {
	BotID  string        `json:"botId"`
	Symbol string        `json:"symbol"`
	Status GridBotStatus `json:"status"`
	// TotalProfit is the total profit in USDT and ROI its decimal ratio.
	TotalProfit        decimal.Decimal `json:"totalProfit"`
	ROI                decimal.Decimal `json:"roi"`
	GridProfit         decimal.Decimal `json:"gridProfit"`
	GridProfitRate     decimal.Decimal `json:"gridProfitRate"`
	UnpairedProfit     decimal.Decimal `json:"unpairedProfit"`
	UnpairedProfitRate decimal.Decimal `json:"unpairedProfitRate"`
	Margin             decimal.Decimal `json:"margin"` // futures only
	// ArbitrageAPR and TotalAPR are annualized as decimal ratios; a bot running
	// for less than a day is annualized as if it had run one full day.
	ArbitrageAPR decimal.Decimal `json:"arbitrageAPR"`
	TotalAPR     decimal.Decimal `json:"totalAPR"`
	CreatedTime  time.Time       `json:"createdTime"`
	// RunningTime is the running duration in milliseconds, not a timestamp.
	RunningTime             decimal.Decimal          `json:"runningTime"`
	CurrentBaseBalance      decimal.Decimal          `json:"currentBaseBalance"`
	CurrentQuoteBalance     decimal.Decimal          `json:"currentQuoteBalance"`
	InitialBaseHoldings     decimal.Decimal          `json:"initialBaseHoldings"`
	InitialQuoteHoldings    decimal.Decimal          `json:"initialQuoteHoldings"`
	ReservedBaseTradingFee  decimal.Decimal          `json:"reservedBaseTradingFee"`
	ReservedQuoteTradingFee decimal.Decimal          `json:"reservedQuoteTradingFee"`
	GridStartPrice          decimal.Decimal          `json:"gridStartPrice"`
	GridType                GridType                 `json:"gridType"`
	MaxPrice                decimal.Decimal          `json:"maxPrice"`
	MinPrice                decimal.Decimal          `json:"minPrice"`
	GridNum                 string                   `json:"gridNum"`
	GridOrderMode           GridOrderMode            `json:"gridOrderMode"`
	BaseInvestmentCoin      string                   `json:"baseInvestmentCoin"`
	BaseInvestmentAmount    decimal.Decimal          `json:"baseInvestmentAmount"`
	QuoteInvestmentCoin     string                   `json:"quoteInvestmentCoin"`
	QuoteInvestmentAmount   decimal.Decimal          `json:"quoteInvestmentAmount"`
	Slippage                string                   `json:"slippage"` // 1%, 2%, none
	TriggerCondition        GridTriggerCondition     `json:"triggerCondition"`
	TriggerParams           []GridIndicatorParams    `json:"triggerParams"`
	TriggerPrice            decimal.Decimal          `json:"triggerPrice"`
	TerminationCondition    GridTerminationCondition `json:"terminationCondition"`
	TerminationParams       []GridIndicatorParams    `json:"terminationParams"`
	TerminationSell         string                   `json:"terminationSell"` // yes, no
	StopLoss                decimal.Decimal          `json:"stopLoss"`
	TakeProfit              decimal.Decimal          `json:"takeProfit"`
	HodlMode                string                   `json:"hodlMode"`            // yes, no
	AutoTransferProfits     string                   `json:"autoTransferProfits"` // yes, no
}

// GetGridBotOrderDetailsService -- GET /api/v3/trade/grid/list-details (UTA trade read)
//
// Returns the resting buy and sell orders of one spot/futures grid strategy.
type GetGridBotOrderDetailsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetGridBotOrderDetailsService(botID string) *GetGridBotOrderDetailsService {
	return &GetGridBotOrderDetailsService{c: c, params: map[string]string{"botId": botID}}
}

func (s *GetGridBotOrderDetailsService) Do(ctx context.Context) (*GridBotOrderDetails, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/grid/list-details", s.params).WithSign()
	return request.Do[GridBotOrderDetails](req)
}

// GridBotOrderDetails is the resting order book of a grid strategy, shared by
// the normal and neutral grid endpoints.
type GridBotOrderDetails struct {
	BotID         string      `json:"botId"`
	Symbol        string      `json:"symbol"`
	BuyOrderList  []GridOrder `json:"buyOrderList"`
	SellOrderList []GridOrder `json:"sellOrderList"`
}

// GridOrder is one resting order of a grid strategy.
type GridOrder struct {
	OrderID       string          `json:"orderId"`
	DelegateCount decimal.Decimal `json:"delegateCount"`
	DelegatePrice decimal.Decimal `json:"delegatePrice"`
	// ChangeRequired is the price move needed to fill the order, in percent with
	// the sign stripped: 1 means 1%.
	ChangeRequired decimal.Decimal `json:"changeRequired"`
}

// ValidateNeutralGridParamsService -- POST /api/v3/trade/grid/validate-neutral (UTA trade read & write)
//
// Checks whether a neutral futures grid configuration is legal without creating
// the strategy.
type ValidateNeutralGridParamsService struct {
	c    *UTAClient
	body map[string]any
}

// NewValidateNeutralGridParamsService starts a validation request. Neutral
// grids are futures-only (USDT-FUTURES or USDC-FUTURES).
func (c *UTAClient) NewValidateNeutralGridParamsService(category Category, symbol string, minPrice, maxPrice decimal.Decimal, gridNum string, gridOrderMode GridOrderMode) *ValidateNeutralGridParamsService {
	return &ValidateNeutralGridParamsService{c: c, body: map[string]any{
		"category":      string(category),
		"symbol":        symbol,
		"minPrice":      minPrice.String(),
		"maxPrice":      maxPrice.String(),
		"gridNum":       gridNum,
		"gridOrderMode": string(gridOrderMode),
	}}
}

// SetLeverage sets the leverage.
func (s *ValidateNeutralGridParamsService) SetLeverage(leverage string) *ValidateNeutralGridParamsService {
	s.body["leverage"] = leverage
	return s
}

// SetInvestment sets the investment; futures grids only accept the quote coin.
func (s *ValidateNeutralGridParamsService) SetInvestment(investment []GridInvestment) *ValidateNeutralGridParamsService {
	s.body["investmentAmount"] = investment
	return s
}

// SetTriggerPrice sets the price at which the strategy starts.
func (s *ValidateNeutralGridParamsService) SetTriggerPrice(triggerPrice decimal.Decimal) *ValidateNeutralGridParamsService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetStopLoss sets the stop-loss price of the strategy.
func (s *ValidateNeutralGridParamsService) SetStopLoss(stopLoss decimal.Decimal) *ValidateNeutralGridParamsService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTakeProfit sets the take-profit price of the strategy.
func (s *ValidateNeutralGridParamsService) SetTakeProfit(takeProfit decimal.Decimal) *ValidateNeutralGridParamsService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetAutoTransferProfits controls whether realized profits are transferred out.
func (s *ValidateNeutralGridParamsService) SetAutoTransferProfits(autoTransferProfits bool) *ValidateNeutralGridParamsService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *ValidateNeutralGridParamsService) Do(ctx context.Context) (*GridValidation, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/validate-neutral", s.body).WithSign()
	return request.Do[GridValidation](req)
}

// CreateNeutralGridBotService -- POST /api/v3/trade/grid/create-neutral-bot (UTA trade read & write)
//
// Creates a neutral futures grid trading bot, which holds both directions
// inside the price range. The reply data carries the new bot's ID.
type CreateNeutralGridBotService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateNeutralGridBotService(category Category, symbol string, minPrice, maxPrice decimal.Decimal, gridNum string, gridOrderMode GridOrderMode, fundsSource []GridFundsSource) *CreateNeutralGridBotService {
	return &CreateNeutralGridBotService{c: c, body: map[string]any{
		"category":      string(category),
		"symbol":        symbol,
		"minPrice":      minPrice.String(),
		"maxPrice":      maxPrice.String(),
		"gridNum":       gridNum,
		"gridOrderMode": string(gridOrderMode),
		"fundsSource":   fundsSource,
	}}
}

// SetLeverage sets the leverage.
func (s *CreateNeutralGridBotService) SetLeverage(leverage string) *CreateNeutralGridBotService {
	s.body["leverage"] = leverage
	return s
}

// SetInvestment sets the investment; futures grids only accept the quote coin.
func (s *CreateNeutralGridBotService) SetInvestment(investment []GridInvestment) *CreateNeutralGridBotService {
	s.body["investmentAmount"] = investment
	return s
}

// SetTriggerPrice sets the price at which the strategy starts.
func (s *CreateNeutralGridBotService) SetTriggerPrice(triggerPrice decimal.Decimal) *CreateNeutralGridBotService {
	s.body["triggerPrice"] = triggerPrice.String()
	return s
}

// SetStopLoss sets the stop-loss price of the strategy.
func (s *CreateNeutralGridBotService) SetStopLoss(stopLoss decimal.Decimal) *CreateNeutralGridBotService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetTakeProfit sets the take-profit price of the strategy.
func (s *CreateNeutralGridBotService) SetTakeProfit(takeProfit decimal.Decimal) *CreateNeutralGridBotService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetAutoTransferProfits controls whether realized profits are transferred out.
func (s *CreateNeutralGridBotService) SetAutoTransferProfits(autoTransferProfits bool) *CreateNeutralGridBotService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *CreateNeutralGridBotService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/create-neutral-bot", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// ModifyNeutralGridBotService -- POST /api/v3/trade/grid/modify-neutral-bot (UTA trade read & write)
//
// Modifies the basic parameters of a running neutral futures grid bot.
// Take-profit and stop-loss are cleared when omitted, so pass the current value
// to keep it.
type ModifyNeutralGridBotService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyNeutralGridBotService(category Category, botID string) *ModifyNeutralGridBotService {
	return &ModifyNeutralGridBotService{c: c, body: map[string]any{
		"category": string(category),
		"botId":    botID,
	}}
}

// SetTakeProfit sets the take-profit price; omitting it clears the current one.
func (s *ModifyNeutralGridBotService) SetTakeProfit(takeProfit decimal.Decimal) *ModifyNeutralGridBotService {
	s.body["takeProfit"] = takeProfit.String()
	return s
}

// SetStopLoss sets the stop-loss price; omitting it clears the current one.
func (s *ModifyNeutralGridBotService) SetStopLoss(stopLoss decimal.Decimal) *ModifyNeutralGridBotService {
	s.body["stopLoss"] = stopLoss.String()
	return s
}

// SetAutoTransferProfits toggles transferring realized profits out.
func (s *ModifyNeutralGridBotService) SetAutoTransferProfits(autoTransferProfits bool) *ModifyNeutralGridBotService {
	s.body["autoTransferProfits"] = yesNo(autoTransferProfits)
	return s
}

func (s *ModifyNeutralGridBotService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/modify-neutral-bot", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// ModifyNeutralGridIntervalService -- POST /api/v3/trade/grid/modify-neutral-grid-interval (UTA trade read & write)
//
// Moves the price range and grid number of a running neutral futures grid bot.
type ModifyNeutralGridIntervalService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyNeutralGridIntervalService(category Category, botID string, minPrice, maxPrice decimal.Decimal, gridNum string) *ModifyNeutralGridIntervalService {
	return &ModifyNeutralGridIntervalService{c: c, body: map[string]any{
		"category": string(category),
		"botId":    botID,
		"minPrice": minPrice.String(),
		"maxPrice": maxPrice.String(),
		"gridNum":  gridNum,
	}}
}

func (s *ModifyNeutralGridIntervalService) Do(ctx context.Context) (*GridBotResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/trade/grid/modify-neutral-grid-interval", s.body).WithSign()
	return request.Do[GridBotResult](req)
}

// GetNeutralGridBotDetailService -- GET /api/v3/trade/grid/neutral-bot-detail (UTA trade read)
//
// Returns the full configuration and running profit of one neutral futures grid
// strategy.
type GetNeutralGridBotDetailService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetNeutralGridBotDetailService(botID string) *GetNeutralGridBotDetailService {
	return &GetNeutralGridBotDetailService{c: c, params: map[string]string{"botId": botID}}
}

func (s *GetNeutralGridBotDetailService) Do(ctx context.Context) (*NeutralGridBotDetail, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/grid/neutral-bot-detail", s.params).WithSign()
	return request.Do[NeutralGridBotDetail](req)
}

// NeutralGridBotDetail is the configuration and running state of a neutral
// futures grid strategy.
type NeutralGridBotDetail struct {
	BotID       string          `json:"botId"`
	Symbol      string          `json:"symbol"`
	Status      GridBotStatus   `json:"status"`
	TotalProfit decimal.Decimal `json:"totalProfit"`
	ROI         decimal.Decimal `json:"roi"`
	// GridProfit is the realized arbitrage profit net of fees, in USDT.
	GridProfit         decimal.Decimal `json:"gridProfit"`
	GridProfitRate     decimal.Decimal `json:"gridProfitRate"`
	UnpairedProfit     decimal.Decimal `json:"unpairedProfit"`
	UnpairedProfitRate decimal.Decimal `json:"unpairedProfitRate"`
	Margin             decimal.Decimal `json:"margin"`
	// ArbitrageAPR and TotalAPR are annualized as decimal ratios; a bot running
	// for less than a day is annualized as if it had run one full day.
	ArbitrageAPR decimal.Decimal `json:"arbitrageAPR"`
	TotalAPR     decimal.Decimal `json:"totalAPR"`
	CreatedTime  time.Time       `json:"createdTime"`
	// RunningTime is the running duration in milliseconds, not a timestamp.
	RunningTime         decimal.Decimal `json:"runningTime"`
	GridStartPrice      decimal.Decimal `json:"gridStartPrice"`
	MaxPrice            decimal.Decimal `json:"maxPrice"`
	MinPrice            decimal.Decimal `json:"minPrice"`
	GridNum             string          `json:"gridNum"`
	GridOrderMode       GridOrderMode   `json:"gridOrderMode"`
	TriggerPrice        decimal.Decimal `json:"triggerPrice"`
	Leverage            string          `json:"leverage"`
	StopLoss            decimal.Decimal `json:"stopLoss"`
	TakeProfit          decimal.Decimal `json:"takeProfit"`
	HoldAveragePrice    decimal.Decimal `json:"holdAveragePrice"`
	HoldPosition        decimal.Decimal `json:"holdPosition"`
	OneGridMinProfit    decimal.Decimal `json:"oneGridMinProfit"`
	OneGridMaxProfit    decimal.Decimal `json:"oneGridMaxProfit"`
	AutoTransferProfits string          `json:"autoTransferProfits"` // yes, no
}

// GetNeutralGridBotOrderDetailsService -- GET /api/v3/trade/grid/neutral-list-details (UTA trade read)
//
// Returns the resting buy and sell orders of one neutral futures grid strategy.
type GetNeutralGridBotOrderDetailsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetNeutralGridBotOrderDetailsService(botID string) *GetNeutralGridBotOrderDetailsService {
	return &GetNeutralGridBotOrderDetailsService{c: c, params: map[string]string{"botId": botID}}
}

func (s *GetNeutralGridBotOrderDetailsService) Do(ctx context.Context) (*GridBotOrderDetails, error) {
	req := request.Get(ctx, s.c, "/api/v3/trade/grid/neutral-list-details", s.params).WithSign()
	return request.Do[GridBotOrderDetails](req)
}
