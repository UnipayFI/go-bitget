package copy

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// MixFollowerMarginType is how the margin per copied symbol is decided by a
// follower's per-symbol copy settings.
type MixFollowerMarginType string

const (
	MixFollowerMarginTypeTrader  MixFollowerMarginType = "trader"  // mirror the trader's margin
	MixFollowerMarginTypeSpecify MixFollowerMarginType = "specify" // use a follower-specified margin coin
)

// MixFollowerLeverType is how the leverage per copied symbol is decided by a
// follower's per-symbol copy settings.
type MixFollowerLeverType string

const (
	MixFollowerLeverTypePosition MixFollowerLeverType = "position" // keep the current position leverage
	MixFollowerLeverTypeSpecify  MixFollowerLeverType = "specify"  // use a follower-specified leverage
	MixFollowerLeverTypeTrader   MixFollowerLeverType = "trader"   // mirror the trader's leverage
)

// MixFollowerTraceType is how the copied order size is derived from the
// trader's order in a follower's per-symbol copy settings.
type MixFollowerTraceType string

const (
	MixFollowerTraceTypePercent MixFollowerTraceType = "percent" // size as a percentage of the trader's order
	MixFollowerTraceTypeAmount  MixFollowerTraceType = "amount"  // a fixed margin amount per order
	MixFollowerTraceTypeCount   MixFollowerTraceType = "count"   // a fixed contract count per order
)

// MixFollowerSettingsMode is the granularity of a follower's copy settings.
type MixFollowerSettingsMode string

const (
	MixFollowerSettingsModeBasic    MixFollowerSettingsMode = "basic"
	MixFollowerSettingsModeAdvanced MixFollowerSettingsMode = "advanced"
)

// MixFollowerSwitch is the on/off toggle value used by follower settings
// (autoCopy, equityGuardian, leverage mode).
type MixFollowerSwitch string

const (
	MixFollowerSwitchOn  MixFollowerSwitch = "on"
	MixFollowerSwitchOff MixFollowerSwitch = "off"
)

// GetMixFollowerCurrentOrdersService -- GET /api/v2/copy/mix-follower/query-current-orders (copy-trading follower read)
//
// Returns the follower's in-progress (currently tracking) copy orders for a
// futures product type.
type GetMixFollowerCurrentOrdersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixFollowerCurrentOrdersService(productType ProductType) *GetMixFollowerCurrentOrdersService {
	return &GetMixFollowerCurrentOrdersService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixFollowerCurrentOrdersService) SetSymbol(symbol string) *GetMixFollowerCurrentOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMixFollowerCurrentOrdersService) SetTraderID(traderID string) *GetMixFollowerCurrentOrdersService {
	s.params["traderId"] = traderID
	return s
}

func (s *GetMixFollowerCurrentOrdersService) SetStartTime(t time.Time) *GetMixFollowerCurrentOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerCurrentOrdersService) SetEndTime(t time.Time) *GetMixFollowerCurrentOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerCurrentOrdersService) SetLimit(limit int) *GetMixFollowerCurrentOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages to records older than the given tracking-order ID.
func (s *GetMixFollowerCurrentOrdersService) SetIDLessThan(id string) *GetMixFollowerCurrentOrdersService {
	s.params["idLessThan"] = id
	return s
}

// SetIDGreaterThan pages to records newer than the given tracking-order ID.
func (s *GetMixFollowerCurrentOrdersService) SetIDGreaterThan(id string) *GetMixFollowerCurrentOrdersService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixFollowerCurrentOrdersService) Do(ctx context.Context) (*MixFollowerCurrentOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-follower/query-current-orders", s.params).WithSign()
	return request.Do[MixFollowerCurrentOrders](req)
}

// MixFollowerCurrentOrders is the paged list of the follower's in-progress copy
// orders (data is an object {trackingList, endId}, not a bare array).
type MixFollowerCurrentOrders struct {
	TrackingList []MixFollowerCurrentOrder `json:"trackingList"`
	EndID        string                    `json:"endId"`
}

// MixFollowerCurrentOrder is one in-progress copy order for the follower.
type MixFollowerCurrentOrder struct {
	TrackingNo    string          `json:"trackingNo"`
	TraderID      string          `json:"traderId"`
	TraderName    string          `json:"traderName"`
	OpenOrderID   string          `json:"openOrderId"`
	CloseOrderID  string          `json:"closeOrderId"`
	Symbol        string          `json:"symbol"`
	PosSide       PosSide         `json:"posSide"`
	OpenLeverage  string          `json:"openLeverage"`
	OpenAvgPrice  decimal.Decimal `json:"openAvgPrice"`
	OpenTime      time.Time       `json:"openTime"`
	OpenSize      decimal.Decimal `json:"openSize"`
	OpenMarginSz  decimal.Decimal `json:"openMarginSz"`
	CloseAvgPrice decimal.Decimal `json:"closeAvgPrice"`
	CloseSize     decimal.Decimal `json:"closeSize"`
	CloseTime     time.Time       `json:"closeTime"`
}

// GetMixFollowerHistoryOrdersService -- GET /api/v2/copy/mix-follower/query-history-orders (copy-trading follower read)
//
// Returns the follower's closed (historical) copy orders for a futures product
// type.
type GetMixFollowerHistoryOrdersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixFollowerHistoryOrdersService(productType ProductType) *GetMixFollowerHistoryOrdersService {
	return &GetMixFollowerHistoryOrdersService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixFollowerHistoryOrdersService) SetSymbol(symbol string) *GetMixFollowerHistoryOrdersService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMixFollowerHistoryOrdersService) SetTraderID(traderID string) *GetMixFollowerHistoryOrdersService {
	s.params["traderId"] = traderID
	return s
}

func (s *GetMixFollowerHistoryOrdersService) SetStartTime(t time.Time) *GetMixFollowerHistoryOrdersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerHistoryOrdersService) SetEndTime(t time.Time) *GetMixFollowerHistoryOrdersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerHistoryOrdersService) SetLimit(limit int) *GetMixFollowerHistoryOrdersService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages to records older than the given tracking-order ID.
func (s *GetMixFollowerHistoryOrdersService) SetIDLessThan(id string) *GetMixFollowerHistoryOrdersService {
	s.params["idLessThan"] = id
	return s
}

// SetIDGreaterThan pages to records newer than the given tracking-order ID.
func (s *GetMixFollowerHistoryOrdersService) SetIDGreaterThan(id string) *GetMixFollowerHistoryOrdersService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixFollowerHistoryOrdersService) Do(ctx context.Context) (*MixFollowerHistoryOrders, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-follower/query-history-orders", s.params).WithSign()
	return request.Do[MixFollowerHistoryOrders](req)
}

// MixFollowerHistoryOrders is the paged list of the follower's closed copy
// orders.
type MixFollowerHistoryOrders struct {
	TrackingList []MixFollowerHistoryOrder `json:"trackingList"`
	EndID        string                    `json:"endId"`
}

// MixFollowerHistoryOrder is one closed copy order for the follower.
type MixFollowerHistoryOrder struct {
	TrackingNo    string          `json:"trackingNo"`
	TraderID      string          `json:"traderId"`
	OpenOrderID   string          `json:"openOrderId"`
	CloseOrderID  string          `json:"closeOrderId"`
	ProductType   ProductType     `json:"productType"`
	Symbol        string          `json:"symbol"`
	PosSide       PosSide         `json:"posSide"`
	OpenLeverage  string          `json:"openLeverage"`
	OpenPriceAvg  decimal.Decimal `json:"openPriceAvg"`
	OpenTime      time.Time       `json:"openTime"`
	OpenSize      decimal.Decimal `json:"openSize"`
	ClosePriceAvg decimal.Decimal `json:"closePriceAvg"`
	CloseFee      decimal.Decimal `json:"closeFee"`
	OpenFee       decimal.Decimal `json:"openFee"`
	CloseSize     decimal.Decimal `json:"closeSize"`
	CloseTime     time.Time       `json:"closeTime"`
	ProfitRate    decimal.Decimal `json:"profitRate"`
	NetProfit     decimal.Decimal `json:"netProfit"`
	AchievedPL    decimal.Decimal `json:"achievedPL"`
}

// SetMixFollowerTPSLService -- POST /api/v2/copy/mix-follower/setting-tpsl (copy-trading follower, state-changing)
//
// Sets, updates or cancels the take-profit / stop-loss prices on one of the
// follower's tracking orders. Per Bitget: an empty price is ignored, "0"
// cancels an existing TPSL, and a positive price sets/updates it.
type SetMixFollowerTPSLService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewSetMixFollowerTPSLService(trackingNo string, productType ProductType) *SetMixFollowerTPSLService {
	return &SetMixFollowerTPSLService{c: c, body: map[string]any{
		"trackingNo":  trackingNo,
		"productType": string(productType),
	}}
}

func (s *SetMixFollowerTPSLService) SetSymbol(symbol string) *SetMixFollowerTPSLService {
	s.body["symbol"] = symbol
	return s
}

// SetStopSurplusPrice sets the take-profit trigger price.
func (s *SetMixFollowerTPSLService) SetStopSurplusPrice(price decimal.Decimal) *SetMixFollowerTPSLService {
	s.body["stopSurplusPrice"] = price.String()
	return s
}

// SetStopLossPrice sets the stop-loss trigger price.
func (s *SetMixFollowerTPSLService) SetStopLossPrice(price decimal.Decimal) *SetMixFollowerTPSLService {
	s.body["stopLossPrice"] = price.String()
	return s
}

func (s *SetMixFollowerTPSLService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-follower/setting-tpsl", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// SetMixFollowerCopySettingsService -- POST /api/v2/copy/mix-follower/settings (copy-trading follower, state-changing)
//
// Adds or modifies the follower's per-symbol copy-trading configuration for a
// trader being followed (up to 10 settings entries).
type SetMixFollowerCopySettingsService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewSetMixFollowerCopySettingsService(traderID string, settings []MixFollowerSettingItem) *SetMixFollowerCopySettingsService {
	return &SetMixFollowerCopySettingsService{c: c, body: map[string]any{
		"traderId": traderID,
		"settings": settings,
	}}
}

// SetAutoCopy toggles auto-following of newly added symbols (on / off).
func (s *SetMixFollowerCopySettingsService) SetAutoCopy(v MixFollowerSwitch) *SetMixFollowerCopySettingsService {
	s.body["autoCopy"] = string(v)
	return s
}

// SetMode selects the configuration granularity (basic / advanced).
func (s *SetMixFollowerCopySettingsService) SetMode(mode MixFollowerSettingsMode) *SetMixFollowerCopySettingsService {
	s.body["mode"] = string(mode)
	return s
}

func (s *SetMixFollowerCopySettingsService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-follower/settings", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// MixFollowerSettingItem is one per-symbol copy-trading configuration entry in
// a Set Copy Trade Settings request.
type MixFollowerSettingItem struct {
	Symbol           string                `json:"symbol"`
	ProductType      ProductType           `json:"productType"`
	MarginType       MixFollowerMarginType `json:"marginType"`
	MarginCoin       string                `json:"marginCoin,omitempty"`
	LeverType        MixFollowerLeverType  `json:"leverType"`
	LongLeverage     string                `json:"longLeverage,omitempty"`
	ShortLeverage    string                `json:"shortLeverage,omitempty"`
	TraceType        MixFollowerTraceType  `json:"traceType"`
	TraceValue       string                `json:"traceValue"`
	MaxHoldSize      string                `json:"maxHoldSize,omitempty"`
	StopSurplusRatio string                `json:"stopSurplusRatio,omitempty"`
	StopLossRatio    string                `json:"stopLossRatio,omitempty"`
}

// GetMixFollowerSettingsService -- GET /api/v2/copy/mix-follower/query-settings (copy-trading follower read)
//
// Returns the follower's per-symbol copy-trading configuration for a specific
// trader being followed.
type GetMixFollowerSettingsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixFollowerSettingsService(traderID string) *GetMixFollowerSettingsService {
	return &GetMixFollowerSettingsService{c: c, params: map[string]string{"traderId": traderID}}
}

func (s *GetMixFollowerSettingsService) Do(ctx context.Context) (*MixFollowerSettings, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-follower/query-settings", s.params).WithSign()
	return request.Do[MixFollowerSettings](req)
}

// MixFollowerSettings is the follower's copy-trading configuration for a
// trader.
type MixFollowerSettings struct {
	FollowerEnable string                      `json:"followerEnable"` // YES, NO
	DetailList     []MixFollowerSettingsDetail `json:"detailList"`
}

// MixFollowerSettingsDetail is one per-symbol copy configuration entry.
type MixFollowerSettingsDetail struct {
	Symbol           string                `json:"symbol"`
	ProductType      ProductType           `json:"productType"`
	MarginType       MixFollowerMarginType `json:"marginType"`
	MarginCoin       string                `json:"marginCoin"`
	LeverType        MixFollowerLeverType  `json:"leverType"`
	LongLeverage     string                `json:"longLeverage"`
	ShortLeverage    string                `json:"shortLeverage"`
	TraceType        MixFollowerTraceType  `json:"traceType"`
	TraceValue       decimal.Decimal       `json:"traceValue"`
	MaxHoldSize      string                `json:"maxHoldSize"`
	StopSurplusRatio decimal.Decimal       `json:"stopSurplusRatio"`
	StopLossRatio    decimal.Decimal       `json:"stopLossRatio"`
}

// GetMixFollowerTradersService -- GET /api/v2/copy/mix-follower/query-traders (copy-trading follower read)
//
// Returns the list of traders the follower is currently following, with
// aggregate follow stats.
type GetMixFollowerTradersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixFollowerTradersService() *GetMixFollowerTradersService {
	return &GetMixFollowerTradersService{c: c, params: map[string]string{}}
}

func (s *GetMixFollowerTradersService) SetStartTime(t time.Time) *GetMixFollowerTradersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerTradersService) SetEndTime(t time.Time) *GetMixFollowerTradersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixFollowerTradersService) SetPageNo(pageNo int) *GetMixFollowerTradersService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetMixFollowerTradersService) SetPageSize(pageSize int) *GetMixFollowerTradersService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetMixFollowerTradersService) Do(ctx context.Context) ([]MixFollowerTrader, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-follower/query-traders", s.params).WithSign()
	resp, err := request.Do[[]MixFollowerTrader](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixFollowerTrader is one trader the follower is following.
type MixFollowerTrader struct {
	CertificationType      string          `json:"certificationType"` // Uncertified, Certified
	TraderID               string          `json:"traderId"`
	TraderName             string          `json:"traderName"`
	MaxFollowLimit         string          `json:"maxFollowLimit"`
	BgbMaxFollowLimit      string          `json:"bgbMaxFollowLimit"`
	FollowCount            string          `json:"followCount"`
	BgbFollowCount         string          `json:"bgbFollowCount"`
	TraceTotalMarginAmount decimal.Decimal `json:"traceTotalMarginAmount"`
	TraceTotalNetProfit    decimal.Decimal `json:"traceTotalNetProfit"`
	TraceTotalProfit       decimal.Decimal `json:"traceTotalProfit"`
	CurrentTradingPairs    []string        `json:"currentTradingPairs"`
	FollowerTime           time.Time       `json:"followerTime"`
}

// CloseMixFollowerPositionsService -- POST /api/v2/copy/mix-follower/close-positions (copy-trading follower, state-changing)
//
// Closes one or more of the follower's copy positions (flash close at market).
type CloseMixFollowerPositionsService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCloseMixFollowerPositionsService(productType ProductType) *CloseMixFollowerPositionsService {
	return &CloseMixFollowerPositionsService{c: c, body: map[string]any{"productType": string(productType)}}
}

func (s *CloseMixFollowerPositionsService) SetTrackingNo(trackingNo string) *CloseMixFollowerPositionsService {
	s.body["trackingNo"] = trackingNo
	return s
}

func (s *CloseMixFollowerPositionsService) SetSymbol(symbol string) *CloseMixFollowerPositionsService {
	s.body["symbol"] = symbol
	return s
}

func (s *CloseMixFollowerPositionsService) SetMarginCoin(marginCoin string) *CloseMixFollowerPositionsService {
	s.body["marginCoin"] = marginCoin
	return s
}

func (s *CloseMixFollowerPositionsService) SetMarginMode(mode MarginMode) *CloseMixFollowerPositionsService {
	s.body["marginMode"] = string(mode)
	return s
}

// SetHoldSide selects the position direction to close (required in hedge mode).
func (s *CloseMixFollowerPositionsService) SetHoldSide(side HoldSide) *CloseMixFollowerPositionsService {
	s.body["holdSide"] = string(side)
	return s
}

func (s *CloseMixFollowerPositionsService) Do(ctx context.Context) (*MixFollowerClosePositionsResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-follower/close-positions", s.body).WithSign()
	return request.Do[MixFollowerClosePositionsResult](req)
}

// MixFollowerClosePositionsResult lists the close orders the request created.
type MixFollowerClosePositionsResult struct {
	OrderIDList []string `json:"orderIdList"`
}

// GetMixFollowerQuantityLimitService -- GET /api/v2/copy/mix-follower/query-quantity-limit (copy-trading follower read)
//
// Returns the per-symbol min/max copy order size the follower may place.
type GetMixFollowerQuantityLimitService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixFollowerQuantityLimitService(productType ProductType) *GetMixFollowerQuantityLimitService {
	return &GetMixFollowerQuantityLimitService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixFollowerQuantityLimitService) SetSymbol(symbol string) *GetMixFollowerQuantityLimitService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMixFollowerQuantityLimitService) Do(ctx context.Context) ([]MixFollowerQuantityLimit, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-follower/query-quantity-limit", s.params).WithSign()
	resp, err := request.Do[[]MixFollowerQuantityLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixFollowerQuantityLimit is the follow size range for one symbol.
type MixFollowerQuantityLimit struct {
	Symbol        string          `json:"symbol"`
	MaxFollowSize decimal.Decimal `json:"maxFollowSize"`
	MinFollowSize decimal.Decimal `json:"minFollowSize"`
}

// CancelMixFollowerTraderService -- POST /api/v2/copy/mix-follower/cancel-trader (copy-trading follower, state-changing)
//
// Unfollows (cancels following of) a trader.
type CancelMixFollowerTraderService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCancelMixFollowerTraderService(traderID string) *CancelMixFollowerTraderService {
	return &CancelMixFollowerTraderService{c: c, body: map[string]any{"traderId": traderID}}
}

func (s *CancelMixFollowerTraderService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-follower/cancel-trader", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}
