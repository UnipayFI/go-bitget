package copy

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SpotFollowerTraceType is the copy-sizing mode a follower selects per symbol.
// It is supplied on the settings request and echoed back in the
// query-settings response.
type SpotFollowerTraceType string

const (
	SpotFollowerTraceTypePercent SpotFollowerTraceType = "percent" // copy a percentage of the trader's order
	SpotFollowerTraceTypeAmount  SpotFollowerTraceType = "amount"  // copy a fixed investment amount
	SpotFollowerTraceTypeCount   SpotFollowerTraceType = "count"   // copy a fixed quantity
)

// SpotFollowerAutoCopy toggles whether the follower automatically copies new
// trading pairs the trader starts trading.
type SpotFollowerAutoCopy string

const (
	SpotFollowerAutoCopyOn  SpotFollowerAutoCopy = "on"
	SpotFollowerAutoCopyOff SpotFollowerAutoCopy = "off"
)

// SpotFollowerMode is the configuration granularity for follower settings.
type SpotFollowerMode string

const (
	SpotFollowerModeBasic    SpotFollowerMode = "basic"
	SpotFollowerModeAdvanced SpotFollowerMode = "advanced"
)

// StopSpotFollowerOrderService -- POST /api/v2/copy/spot-follower/stop-order (private, state-changing)
//
// Stops (sells out) one or more spot copy-trading orders by their tracking
// numbers. Execution is atomic: either all the listed orders are stopped or
// none are (up to 50 per call).
type StopSpotFollowerOrderService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewStopSpotFollowerOrderService(trackingNoList []string) *StopSpotFollowerOrderService {
	return &StopSpotFollowerOrderService{c: c, body: map[string]any{"trackingNoList": trackingNoList}}
}

func (s *StopSpotFollowerOrderService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-follower/stop-order").SetBody(s.body).WithSign()
	return request.Do[string](req)
}

// SpotFollowerSettingItem is one per-symbol following configuration carried in
// the AddSpotFollowerSettingsService body's settings array.
type SpotFollowerSettingItem struct {
	Symbol           string                `json:"symbol"`
	TraceType        SpotFollowerTraceType `json:"traceType"`
	MaxHoldSize      decimal.Decimal       `json:"maxHoldSize"`
	TraceValue       decimal.Decimal       `json:"traceValue"`
	StopLossRatio    decimal.Decimal       `json:"stopLossRatio,omitzero"`
	StopSurplusRatio decimal.Decimal       `json:"stopSurplusRatio,omitzero"`
}

// AddSpotFollowerSettingsService -- POST /api/v2/copy/spot-follower/settings (private, state-changing)
//
// Adds or modifies the follower's per-symbol following configurations for a
// given trader.
type AddSpotFollowerSettingsService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewAddSpotFollowerSettingsService(traderID string, settings []SpotFollowerSettingItem) *AddSpotFollowerSettingsService {
	return &AddSpotFollowerSettingsService{c: c, body: map[string]any{
		"traderId": traderID,
		"settings": settings,
	}}
}

// SetAutoCopy controls whether new symbols the trader trades are auto-copied.
func (s *AddSpotFollowerSettingsService) SetAutoCopy(v SpotFollowerAutoCopy) *AddSpotFollowerSettingsService {
	s.body["autoCopy"] = string(v)
	return s
}

// SetMode selects the configuration granularity (basic or advanced, default
// advanced).
func (s *AddSpotFollowerSettingsService) SetMode(v SpotFollowerMode) *AddSpotFollowerSettingsService {
	s.body["mode"] = string(v)
	return s
}

func (s *AddSpotFollowerSettingsService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-follower/settings").SetBody(s.body).WithSign()
	return request.Do[string](req)
}

// SetSpotFollowerTpslService -- POST /api/v2/copy/spot-follower/setting-tpsl (private, state-changing)
//
// Sets, updates or cancels the take-profit and/or stop-loss price for a single
// spot copy-trading order identified by its tracking number. An empty price
// leaves the side unchanged; "0" cancels it.
type SetSpotFollowerTpslService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewSetSpotFollowerTpslService(trackingNo string) *SetSpotFollowerTpslService {
	return &SetSpotFollowerTpslService{c: c, body: map[string]any{"trackingNo": trackingNo}}
}

// SetStopSurplusPrice sets the take-profit price.
func (s *SetSpotFollowerTpslService) SetStopSurplusPrice(v decimal.Decimal) *SetSpotFollowerTpslService {
	s.body["stopSurplusPrice"] = v.String()
	return s
}

// SetStopLossPrice sets the stop-loss price.
func (s *SetSpotFollowerTpslService) SetStopLossPrice(v decimal.Decimal) *SetSpotFollowerTpslService {
	s.body["stopLossPrice"] = v.String()
	return s
}

func (s *SetSpotFollowerTpslService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-follower/setting-tpsl", s.body).WithSign()
	return request.Do[string](req)
}

// GetSpotFollowerTradersService -- GET /api/v2/copy/spot-follower/query-traders (private)
//
// Returns the list of traders the account is currently following on spot
// copy-trading.
type GetSpotFollowerTradersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotFollowerTradersService() *GetSpotFollowerTradersService {
	return &GetSpotFollowerTradersService{c: c, params: map[string]string{}}
}

// SetPageNo sets the 1-based page number (default 1).
func (s *GetSpotFollowerTradersService) SetPageNo(v string) *GetSpotFollowerTradersService {
	s.params["pageNo"] = v
	return s
}

// SetPageSize sets the page size (default 20, max 50).
func (s *GetSpotFollowerTradersService) SetPageSize(v string) *GetSpotFollowerTradersService {
	s.params["pageSize"] = v
	return s
}

// SetStartTime filters traders followed at or after t (millisecond timestamp).
func (s *GetSpotFollowerTradersService) SetStartTime(v string) *GetSpotFollowerTradersService {
	s.params["startTime"] = v
	return s
}

// SetEndTime filters traders followed at or before t (millisecond timestamp).
func (s *GetSpotFollowerTradersService) SetEndTime(v string) *GetSpotFollowerTradersService {
	s.params["endTime"] = v
	return s
}

func (s *GetSpotFollowerTradersService) Do(ctx context.Context) (*SpotFollowerTraderList, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-follower/query-traders", s.params).WithSign()
	return request.Do[SpotFollowerTraderList](req)
}

// SpotFollowerTraderList is the paginated my-trader-list payload.
type SpotFollowerTraderList struct {
	ResultList []SpotFollowerTrader `json:"resultList"`
}

// SpotFollowerTrader is a single trader the account follows.
type SpotFollowerTrader struct {
	TraderID            string          `json:"traderId"`
	TraderName          string          `json:"traderName"`
	CertificationType   string          `json:"certificationType"` // Uncertified, Certified
	TraceTotalAmount    decimal.Decimal `json:"traceTotalAmount"`
	TraceTotalNetProfit decimal.Decimal `json:"traceTotalNetProfit"`
	TraceTotalProfit    decimal.Decimal `json:"traceTotalProfit"`
	MaxFollowLimit      string          `json:"maxFollowLimit"`
	BgbMaxFollowLimit   string          `json:"bgbMaxFollowLimit"`
	FollowCount         string          `json:"followCount"`
	BgbFollowCount      string          `json:"bgbFollowCount"`
	FollowerTime        time.Time       `json:"followerTime"`
}

// GetSpotFollowerTraderSymbolsService -- GET /api/v2/copy/spot-follower/query-trader-symbols (private)
//
// Returns the trading pairs a followed trader is currently copy-trading.
type GetSpotFollowerTraderSymbolsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotFollowerTraderSymbolsService(traderID string) *GetSpotFollowerTraderSymbolsService {
	return &GetSpotFollowerTraderSymbolsService{c: c, params: map[string]string{"traderId": traderID}}
}

func (s *GetSpotFollowerTraderSymbolsService) Do(ctx context.Context) (*SpotFollowerTraderSymbols, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-follower/query-trader-symbols", s.params).WithSign()
	return request.Do[SpotFollowerTraderSymbols](req)
}

// SpotFollowerTraderSymbols is the trader's current trading-pair list.
type SpotFollowerTraderSymbols struct {
	CurrentTradingList []string `json:"currentTradingList"`
}

// GetSpotFollowerSettingsService -- GET /api/v2/copy/spot-follower/query-settings (private)
//
// Returns the follower's configuration for a given trader, including the active
// per-symbol settings and the per-symbol system limits.
type GetSpotFollowerSettingsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotFollowerSettingsService(traderID string) *GetSpotFollowerSettingsService {
	return &GetSpotFollowerSettingsService{c: c, params: map[string]string{"traderId": traderID}}
}

func (s *GetSpotFollowerSettingsService) Do(ctx context.Context) (*SpotFollowerSettings, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-follower/query-settings", s.params).WithSign()
	return request.Do[SpotFollowerSettings](req)
}

// SpotFollowerSettings is the follow-configuration payload for a trader.
type SpotFollowerSettings struct {
	Enable                 string                           `json:"enable"`
	ProfitRate             decimal.Decimal                  `json:"profitRate"`
	SettledInDays          string                           `json:"settledInDays"`
	TraderHeadPic          string                           `json:"traderHeadPic"`
	TraderName             string                           `json:"traderName"`
	TradeSettingList       []SpotFollowerTradeSetting       `json:"tradeSettingList"`
	TradeSymbolSettingList []SpotFollowerTradeSymbolSetting `json:"tradeSymbolSettingList"`
	EndID                  string                           `json:"endId"`
}

// SpotFollowerTradeSetting is one active per-symbol follow configuration.
type SpotFollowerTradeSetting struct {
	Symbol            string                `json:"symbol"`
	MaxTraceAmount    decimal.Decimal       `json:"maxTraceAmount"`
	StopLossRation    decimal.Decimal       `json:"stopLossRation"`
	StopSurplusRation decimal.Decimal       `json:"stopSurplusRation"`
	TraceType         SpotFollowerTraceType `json:"traceType"`
}

// SpotFollowerTradeSymbolSetting is the per-symbol system limit envelope shown
// when configuring a follow.
type SpotFollowerTradeSymbolSetting struct {
	Symbol                    string          `json:"symbol"`
	MinTraceAmount            decimal.Decimal `json:"minTraceAmount"`
	MaxTraceAmount            decimal.Decimal `json:"maxTraceAmount"`
	MaxTraceAmountSystem      decimal.Decimal `json:"maxTraceAmountSystem"`
	MinTraceSize              decimal.Decimal `json:"minTraceSize"`
	MaxTraceSize              decimal.Decimal `json:"maxTraceSize"`
	MinTraceRation            decimal.Decimal `json:"minTraceRation"`
	MaxTraceRation            decimal.Decimal `json:"maxTraceRation"`
	MinStopLossRation         decimal.Decimal `json:"minStopLossRation"`
	MaxStopLossRation         decimal.Decimal `json:"maxStopLossRation"`
	SliderMaxStopLossRatio    decimal.Decimal `json:"sliderMaxStopLossRatio"`
	MinStopSurplusRation      decimal.Decimal `json:"minStopSurplusRation"`
	MaxStopSurplusRation      decimal.Decimal `json:"maxStopSurplusRation"`
	SliderMaxStopSurplusRatio decimal.Decimal `json:"sliderMaxStopSurplusRatio"`
}

// GetSpotFollowerHistoryOrdersService -- GET /api/v2/copy/spot-follower/query-history-orders (private)
//
// Returns the follower's finished (closed) spot copy-trading orders.
type GetSpotFollowerHistoryOrdersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotFollowerHistoryOrdersService() *GetSpotFollowerHistoryOrdersService {
	return &GetSpotFollowerHistoryOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol filters by trading pair.
func (s *GetSpotFollowerHistoryOrdersService) SetSymbol(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["symbol"] = v
	return s
}

// SetTraderID filters by trader.
func (s *GetSpotFollowerHistoryOrdersService) SetTraderID(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["traderId"] = v
	return s
}

// SetIDLessThan pages to records before the given id.
func (s *GetSpotFollowerHistoryOrdersService) SetIDLessThan(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["idLessThan"] = v
	return s
}

// SetIDGreaterThan pages to records after the given id.
func (s *GetSpotFollowerHistoryOrdersService) SetIDGreaterThan(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["idGreaterThan"] = v
	return s
}

// SetStartTime filters records at or after the given millisecond timestamp.
func (s *GetSpotFollowerHistoryOrdersService) SetStartTime(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["startTime"] = v
	return s
}

// SetEndTime filters records at or before the given millisecond timestamp.
func (s *GetSpotFollowerHistoryOrdersService) SetEndTime(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["endTime"] = v
	return s
}

// SetLimit caps the number of records returned (default 20, max 50).
func (s *GetSpotFollowerHistoryOrdersService) SetLimit(v string) *GetSpotFollowerHistoryOrdersService {
	s.params["limit"] = v
	return s
}

func (s *GetSpotFollowerHistoryOrdersService) Do(ctx context.Context) (*SpotFollowerHistoryOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-follower/query-history-orders", s.params).WithSign()
	return request.Do[SpotFollowerHistoryOrders](req)
}

// SpotFollowerHistoryOrders is the paginated history-tracking-order payload.
type SpotFollowerHistoryOrders struct {
	EndID        string                     `json:"endId"`
	TrackingList []SpotFollowerHistoryOrder `json:"trackingList"`
}

// SpotFollowerHistoryOrder is a single closed spot copy-trading order.
type SpotFollowerHistoryOrder struct {
	TrackingNo  string          `json:"trackingNo"`
	TraderID    string          `json:"traderId"`
	Symbol      string          `json:"symbol"`
	FillSize    decimal.Decimal `json:"fillSize"`
	BuyPrice    decimal.Decimal `json:"buyPrice"`
	SellPrice   decimal.Decimal `json:"sellPrice"`
	BuyFee      decimal.Decimal `json:"buyFee"`
	SellFee     decimal.Decimal `json:"sellFee"`
	AchievedPL  decimal.Decimal `json:"achievedPL"`
	AchievedPLR decimal.Decimal `json:"achievedPLR"`
	BuyTime     time.Time       `json:"buyTime"`
	SellTime    time.Time       `json:"sellTime"`
}

// GetSpotFollowerCurrentOrdersService -- GET /api/v2/copy/spot-follower/query-current-orders (private)
//
// Returns the follower's currently open spot copy-trading orders.
type GetSpotFollowerCurrentOrdersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotFollowerCurrentOrdersService() *GetSpotFollowerCurrentOrdersService {
	return &GetSpotFollowerCurrentOrdersService{c: c, params: map[string]string{}}
}

// SetSymbol filters by trading pair.
func (s *GetSpotFollowerCurrentOrdersService) SetSymbol(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["symbol"] = v
	return s
}

// SetTraderID filters by trader.
func (s *GetSpotFollowerCurrentOrdersService) SetTraderID(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["traderId"] = v
	return s
}

// SetIDLessThan pages to records before the given id.
func (s *GetSpotFollowerCurrentOrdersService) SetIDLessThan(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["idLessThan"] = v
	return s
}

// SetIDGreaterThan pages to records after the given id.
func (s *GetSpotFollowerCurrentOrdersService) SetIDGreaterThan(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["idGreaterThan"] = v
	return s
}

// SetStartTime filters records at or after the given millisecond timestamp.
func (s *GetSpotFollowerCurrentOrdersService) SetStartTime(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["startTime"] = v
	return s
}

// SetEndTime filters records at or before the given millisecond timestamp.
func (s *GetSpotFollowerCurrentOrdersService) SetEndTime(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["endTime"] = v
	return s
}

// SetLimit caps the number of records returned (default 20, max 50).
func (s *GetSpotFollowerCurrentOrdersService) SetLimit(v string) *GetSpotFollowerCurrentOrdersService {
	s.params["limit"] = v
	return s
}

func (s *GetSpotFollowerCurrentOrdersService) Do(ctx context.Context) (*SpotFollowerCurrentOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-follower/query-current-orders", s.params).WithSign()
	return request.Do[SpotFollowerCurrentOrders](req)
}

// SpotFollowerCurrentOrders is the paginated current-copy-trade-order payload.
type SpotFollowerCurrentOrders struct {
	EndID        string                     `json:"endId"`
	TrackingList []SpotFollowerCurrentOrder `json:"trackingList"`
}

// SpotFollowerCurrentOrder is a single open spot copy-trading order.
type SpotFollowerCurrentOrder struct {
	TrackingNo       string          `json:"trackingNo"`
	TraderID         string          `json:"traderId"`
	Symbol           string          `json:"symbol"`
	BuyFillSize      decimal.Decimal `json:"buyFillSize"`
	BuyDelegateSize  decimal.Decimal `json:"buyDelegateSize"`
	BuyPrice         decimal.Decimal `json:"buyPrice"`
	BuyFee           decimal.Decimal `json:"buyFee"`
	UnrealizedPL     decimal.Decimal `json:"unrealizedPL"`
	UnrealizedPLR    decimal.Decimal `json:"unrealizedPLR"`
	StopSurplusPrice decimal.Decimal `json:"stopSurplusPrice"`
	StopLossPrice    decimal.Decimal `json:"stopLossPrice"`
	BuyTime          time.Time       `json:"buyTime"`
}

// CloseSpotFollowerOrderService -- POST /api/v2/copy/spot-follower/order-close-tracking (private, state-changing)
//
// Sells out (closes) one or more open spot copy-trading orders for a single
// trading pair. Execution is atomic across the listed tracking numbers (up to
// 50 per call).
type CloseSpotFollowerOrderService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCloseSpotFollowerOrderService(symbol string, trackingNoList []string) *CloseSpotFollowerOrderService {
	return &CloseSpotFollowerOrderService{c: c, body: map[string]any{
		"symbol":         symbol,
		"trackingNoList": trackingNoList,
	}}
}

func (s *CloseSpotFollowerOrderService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-follower/order-close-tracking").SetBody(s.body).WithSign()
	return request.Do[string](req)
}

// CancelSpotFollowerTraderService -- POST /api/v2/copy/spot-follower/cancel-trader (private, state-changing)
//
// Cancels following (unfollows) a trader on spot copy-trading.
type CancelSpotFollowerTraderService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCancelSpotFollowerTraderService(traderID string) *CancelSpotFollowerTraderService {
	return &CancelSpotFollowerTraderService{c: c, body: map[string]any{"traderId": traderID}}
}

func (s *CancelSpotFollowerTraderService) Do(ctx context.Context) (*string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-follower/cancel-trader", s.body).WithSign()
	return request.Do[string](req)
}
