package copy

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetMixTraderCurrentTrackService -- GET /api/v2/copy/mix-trader/order-current-track (signed)
//
// Returns the elite trader's currently-open (in-progress) tracking orders for a
// futures product line.
type GetMixTraderCurrentTrackService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderCurrentTrackService(productType ProductType) *GetMixTraderCurrentTrackService {
	return &GetMixTraderCurrentTrackService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixTraderCurrentTrackService) SetSymbol(symbol string) *GetMixTraderCurrentTrackService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMixTraderCurrentTrackService) SetStartTime(t time.Time) *GetMixTraderCurrentTrackService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderCurrentTrackService) SetEndTime(t time.Time) *GetMixTraderCurrentTrackService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderCurrentTrackService) SetLimit(limit int) *GetMixTraderCurrentTrackService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetMixTraderCurrentTrackService) SetIDGreaterThan(id string) *GetMixTraderCurrentTrackService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixTraderCurrentTrackService) SetIDLessThan(id string) *GetMixTraderCurrentTrackService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetMixTraderCurrentTrackService) Do(ctx context.Context) (*MixTraderCurrentTrackResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/order-current-track", s.params).WithSign()
	return request.Do[MixTraderCurrentTrackResult](req)
}

// MixTraderCurrentTrackResult is the paginated container for the trader's
// current tracking orders.
type MixTraderCurrentTrackResult struct {
	TrackingList []MixTraderCurrentTrack `json:"trackingList"`
	EndID        string                  `json:"endId"`
}

// MixTraderCurrentTrack is one in-progress tracking order opened by the elite
// trader.
type MixTraderCurrentTrack struct {
	TrackingNo             string          `json:"trackingNo"`
	OpenOrderID            string          `json:"openOrderId"`
	Symbol                 string          `json:"symbol"`
	PosSide                PosSide         `json:"posSide"`
	OpenLeverage           string          `json:"openLeverage"`
	OpenPriceAvg           decimal.Decimal `json:"openPriceAvg"`
	OpenTime               time.Time       `json:"openTime"`
	OpenSize               decimal.Decimal `json:"openSize"`
	PresetStopSurplusPrice decimal.Decimal `json:"presetStopSurplusPrice"`
	PresetStopLossPrice    decimal.Decimal `json:"presetStopLossPrice"`
	OpenFee                decimal.Decimal `json:"openFee"`
	FollowCount            string          `json:"followCount"`
}

// GetMixTraderHistoryTrackService -- GET /api/v2/copy/mix-trader/order-history-track (signed)
//
// Returns the elite trader's closed (terminal) tracking orders for a futures
// product line.
type GetMixTraderHistoryTrackService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderHistoryTrackService(productType ProductType) *GetMixTraderHistoryTrackService {
	return &GetMixTraderHistoryTrackService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixTraderHistoryTrackService) SetSymbol(symbol string) *GetMixTraderHistoryTrackService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetMixTraderHistoryTrackService) SetStartTime(t time.Time) *GetMixTraderHistoryTrackService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderHistoryTrackService) SetEndTime(t time.Time) *GetMixTraderHistoryTrackService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderHistoryTrackService) SetLimit(limit int) *GetMixTraderHistoryTrackService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetMixTraderHistoryTrackService) SetOrder(order string) *GetMixTraderHistoryTrackService {
	s.params["order"] = order
	return s
}

func (s *GetMixTraderHistoryTrackService) SetIDGreaterThan(id string) *GetMixTraderHistoryTrackService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixTraderHistoryTrackService) SetIDLessThan(id string) *GetMixTraderHistoryTrackService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetMixTraderHistoryTrackService) Do(ctx context.Context) (*MixTraderHistoryTrackResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/order-history-track", s.params).WithSign()
	return request.Do[MixTraderHistoryTrackResult](req)
}

// MixTraderHistoryTrackResult is the paginated container for the trader's
// closed tracking orders.
type MixTraderHistoryTrackResult struct {
	TrackingList []MixTraderHistoryTrack `json:"trackingList"`
	EndID        string                  `json:"endId"`
}

// MixTraderHistoryTrack is one closed tracking order opened by the elite
// trader.
type MixTraderHistoryTrack struct {
	TrackingNo    string          `json:"trackingNo"`
	Symbol        string          `json:"symbol"`
	ProductType   ProductType     `json:"productType"`
	OpenOrderID   string          `json:"openOrderId"`
	CloseOrderID  string          `json:"closeOrderId"`
	PosSide       PosSide         `json:"posSide"`
	OpenLeverage  string          `json:"openLeverage"`
	OpenPriceAvg  decimal.Decimal `json:"openPriceAvg"`
	OpenTime      time.Time       `json:"openTime"`
	OpenSize      decimal.Decimal `json:"openSize"`
	CloseSize     decimal.Decimal `json:"closeSize"`
	CloseTime     time.Time       `json:"closeTime"`
	ClosePriceAvg decimal.Decimal `json:"closePriceAvg"`
	StopType      TraceStatus     `json:"stopType"`
	AchievedPL    decimal.Decimal `json:"achievedPL"`
	OpenFee       decimal.Decimal `json:"openFee"`
	CloseFee      decimal.Decimal `json:"closeFee"`
	CTime         time.Time       `json:"cTime"`
}

// CreateMixTraderCopyAPIService -- POST /api/v2/copy/mix-trader/create-copy-api (signed, state-changing)
//
// Creates the elite trader's dedicated copy-trading API key. This can only be
// done once per trader.
type CreateMixTraderCopyAPIService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCreateMixTraderCopyAPIService(passphrase string) *CreateMixTraderCopyAPIService {
	return &CreateMixTraderCopyAPIService{c: c, body: map[string]any{"passphrase": passphrase}}
}

func (s *CreateMixTraderCopyAPIService) Do(ctx context.Context) ([]MixTraderCopyAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/create-copy-api", s.body).WithSign()
	resp, err := request.Do[[]MixTraderCopyAPIKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderCopyAPIKey is the freshly-created copy-trading API credential.
type MixTraderCopyAPIKey struct {
	APIKey      string   `json:"apikey"`
	Secret      string   `json:"secret"`
	Permissions []string `json:"permissions"` // contract_trade, copytrading_trade
}

// ModifyMixTraderOrderTPSLService -- POST /api/v2/copy/mix-trader/order-modify-tpsl (signed, state-changing)
//
// Updates the take-profit and/or stop-loss prices on an open tracking order.
// At least one of stopSurplusPrice / stopLossPrice must be set.
type ModifyMixTraderOrderTPSLService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewModifyMixTraderOrderTPSLService(trackingNo string, productType ProductType) *ModifyMixTraderOrderTPSLService {
	return &ModifyMixTraderOrderTPSLService{c: c, body: map[string]any{
		"trackingNo":  trackingNo,
		"productType": string(productType),
	}}
}

func (s *ModifyMixTraderOrderTPSLService) SetStopSurplusPrice(price decimal.Decimal) *ModifyMixTraderOrderTPSLService {
	s.body["stopSurplusPrice"] = price.String()
	return s
}

func (s *ModifyMixTraderOrderTPSLService) SetStopLossPrice(price decimal.Decimal) *ModifyMixTraderOrderTPSLService {
	s.body["stopLossPrice"] = price.String()
	return s
}

func (s *ModifyMixTraderOrderTPSLService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/order-modify-tpsl", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetMixTraderOrderTotalDetailService -- GET /api/v2/copy/mix-trader/order-total-detail (signed)
//
// Returns the elite trader's aggregate tracking-order statistics (ROI, follower
// counts, win rate, and rolling weekly/monthly profit series).
type GetMixTraderOrderTotalDetailService struct {
	c *CopyClient
}

func (c *CopyClient) NewGetMixTraderOrderTotalDetailService() *GetMixTraderOrderTotalDetailService {
	return &GetMixTraderOrderTotalDetailService{c: c}
}

func (s *GetMixTraderOrderTotalDetailService) Do(ctx context.Context) (*MixTraderOrderTotalDetail, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/order-total-detail").WithSign()
	return request.Do[MixTraderOrderTotalDetail](req)
}

// MixTraderOrderTotalDetail is the trader's aggregate performance summary.
type MixTraderOrderTotalDetail struct {
	Roi                       decimal.Decimal        `json:"roi"`
	TradingOrderNum           string                 `json:"tradingOrderNum"`
	TotalFollowerNum          string                 `json:"totalFollowerNum"`
	CurrentFollowerNum        string                 `json:"currentFollowerNum"`
	TotalPL                   decimal.Decimal        `json:"totalpl"`
	GainNum                   string                 `json:"gainNum"`
	LossNum                   string                 `json:"lossNum"`
	WinRate                   decimal.Decimal        `json:"winRate"`
	TotalEquity               decimal.Decimal        `json:"totalEquity"`
	TradingPairsAvailableList []string               `json:"tradingPairsAvailableList"`
	LastWeekRoiList           []MixTraderRoiPoint    `json:"lastWeekRoiList"`
	LastWeekProfitList        []MixTraderProfitPoint `json:"lastWeekProfitList"`
	LastMonthRoiList          []MixTraderRoiPoint    `json:"lastMonthRoiList"`
	LastMonthProfitList       []MixTraderProfitPoint `json:"lastMonthProfitList"`
}

// MixTraderRoiPoint is one point on a rolling ROI series.
type MixTraderRoiPoint struct {
	Rate  decimal.Decimal `json:"rate"`
	Ctime time.Time       `json:"ctime"`
}

// MixTraderProfitPoint is one point on a rolling profit series.
type MixTraderProfitPoint struct {
	Amount decimal.Decimal `json:"amount"`
	Ctime  time.Time       `json:"ctime"`
}

// GetMixTraderProfitHistorySummarysService -- GET /api/v2/copy/mix-trader/profit-history-summarys (signed)
//
// Returns the elite trader's profit-share overview plus per-coin historical
// profit-share totals.
type GetMixTraderProfitHistorySummarysService struct {
	c *CopyClient
}

func (c *CopyClient) NewGetMixTraderProfitHistorySummarysService() *GetMixTraderProfitHistorySummarysService {
	return &GetMixTraderProfitHistorySummarysService{c: c}
}

func (s *GetMixTraderProfitHistorySummarysService) Do(ctx context.Context) (*MixTraderProfitHistorySummary, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/profit-history-summarys").WithSign()
	return request.Do[MixTraderProfitHistorySummary](req)
}

// MixTraderProfitHistorySummary is the trader's profit-share summary container.
type MixTraderProfitHistorySummary struct {
	ProfitSummary     MixTraderProfitOverview      `json:"profitSummary"`
	ProfitHistoryList []MixTraderProfitHistoryItem `json:"profitHistoryList"`
}

// MixTraderProfitOverview is the high-level profit-share rollup.
type MixTraderProfitOverview struct {
	YesterdayProfit decimal.Decimal `json:"yesterdayProfit"`
	YesterdayTime   time.Time       `json:"yesterdayTime"`
	SumProfit       decimal.Decimal `json:"sumProfit"`
	WaitProfit      decimal.Decimal `json:"waitProfit"`
}

// MixTraderProfitHistoryItem is one per-coin historical profit-share total.
type MixTraderProfitHistoryItem struct {
	Coin           string          `json:"coin"`
	ProfitCount    decimal.Decimal `json:"profitCount"`
	LastProfitTime time.Time       `json:"lastProfitTime"`
}

// GetMixTraderProfitHistoryDetailsService -- GET /api/v2/copy/mix-trader/profit-history-details (signed)
//
// Returns the elite trader's individual historical profit-share distribution
// records.
type GetMixTraderProfitHistoryDetailsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderProfitHistoryDetailsService() *GetMixTraderProfitHistoryDetailsService {
	return &GetMixTraderProfitHistoryDetailsService{c: c, params: map[string]string{}}
}

func (s *GetMixTraderProfitHistoryDetailsService) SetCoin(coin string) *GetMixTraderProfitHistoryDetailsService {
	s.params["coin"] = coin
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) SetStartTime(t time.Time) *GetMixTraderProfitHistoryDetailsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) SetEndTime(t time.Time) *GetMixTraderProfitHistoryDetailsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) SetLimit(limit int) *GetMixTraderProfitHistoryDetailsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) SetIDGreaterThan(id string) *GetMixTraderProfitHistoryDetailsService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) SetIDLessThan(id string) *GetMixTraderProfitHistoryDetailsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetMixTraderProfitHistoryDetailsService) Do(ctx context.Context) (*MixTraderProfitHistoryDetailsResult, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/profit-history-details", s.params).WithSign()
	return request.Do[MixTraderProfitHistoryDetailsResult](req)
}

// MixTraderProfitHistoryDetailsResult is the paginated container for historical
// profit-share details.
type MixTraderProfitHistoryDetailsResult struct {
	ProfitList []MixTraderProfitHistoryDetail `json:"profitList"`
	EndID      string                         `json:"endId"`
}

// MixTraderProfitHistoryDetail is one settled profit-share distribution record.
type MixTraderProfitHistoryDetail struct {
	ProfitID   string          `json:"profitId"`
	Coin       string          `json:"coin"`
	Profit     decimal.Decimal `json:"profit"`
	NickName   string          `json:"nickName"`
	ProfitTime time.Time       `json:"profitTime"`
}

// CloseMixTraderPositionsService -- POST /api/v2/copy/mix-trader/order-close-positions (signed, state-changing)
//
// Closes the trader's tracking position(s). When only productType is supplied,
// all positions on that product line are closed.
type CloseMixTraderPositionsService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCloseMixTraderPositionsService(productType ProductType) *CloseMixTraderPositionsService {
	return &CloseMixTraderPositionsService{c: c, body: map[string]any{"productType": string(productType)}}
}

func (s *CloseMixTraderPositionsService) SetTrackingNo(trackingNo string) *CloseMixTraderPositionsService {
	s.body["trackingNo"] = trackingNo
	return s
}

func (s *CloseMixTraderPositionsService) SetSymbol(symbol string) *CloseMixTraderPositionsService {
	s.body["symbol"] = symbol
	return s
}

func (s *CloseMixTraderPositionsService) Do(ctx context.Context) ([]MixTraderClosedPosition, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/order-close-positions", s.body).WithSign()
	resp, err := request.Do[[]MixTraderClosedPosition](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderClosedPosition is one tracking position that was closed.
type MixTraderClosedPosition struct {
	TrackingNo  string      `json:"trackingNo"`
	Symbol      string      `json:"symbol"`
	ProductType ProductType `json:"productType"`
}

// GetMixTraderProfitDetailsService -- GET /api/v2/copy/mix-trader/profit-details (signed)
//
// Returns the elite trader's per-follower unrealized profit-share details.
type GetMixTraderProfitDetailsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderProfitDetailsService() *GetMixTraderProfitDetailsService {
	return &GetMixTraderProfitDetailsService{c: c, params: map[string]string{}}
}

func (s *GetMixTraderProfitDetailsService) SetCoin(coin string) *GetMixTraderProfitDetailsService {
	s.params["coin"] = coin
	return s
}

func (s *GetMixTraderProfitDetailsService) SetPageSize(pageSize int) *GetMixTraderProfitDetailsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetMixTraderProfitDetailsService) SetPageNo(pageNo int) *GetMixTraderProfitDetailsService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetMixTraderProfitDetailsService) Do(ctx context.Context) ([]MixTraderProfitDetail, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/profit-details", s.params).WithSign()
	resp, err := request.Do[[]MixTraderProfitDetail](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderProfitDetail is one unrealized profit-share record.
type MixTraderProfitDetail struct {
	Coin     string          `json:"coin"`
	Profit   decimal.Decimal `json:"profit"`
	NickName string          `json:"nickName"`
}

// GetMixTraderProfitsGroupCoinDateService -- GET /api/v2/copy/mix-trader/profits-group-coin-date (signed)
//
// Returns the trader's profit-share totals grouped by coin and date.
type GetMixTraderProfitsGroupCoinDateService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderProfitsGroupCoinDateService() *GetMixTraderProfitsGroupCoinDateService {
	return &GetMixTraderProfitsGroupCoinDateService{c: c, params: map[string]string{}}
}

func (s *GetMixTraderProfitsGroupCoinDateService) SetPageSize(pageSize int) *GetMixTraderProfitsGroupCoinDateService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetMixTraderProfitsGroupCoinDateService) SetPageNo(pageNo int) *GetMixTraderProfitsGroupCoinDateService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetMixTraderProfitsGroupCoinDateService) Do(ctx context.Context) ([]MixTraderProfitGroup, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/profits-group-coin-date", s.params).WithSign()
	resp, err := request.Do[[]MixTraderProfitGroup](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderProfitGroup is one coin-and-date profit-share total.
type MixTraderProfitGroup struct {
	Coin       string          `json:"coin"`
	Profit     decimal.Decimal `json:"profit"`
	ProfitTime time.Time       `json:"profitTime"`
}

// GetMixTraderConfigQuerySymbolsService -- GET /api/v2/copy/mix-trader/config-query-symbols (signed)
//
// Returns the trader's copy-trade symbol settings (which symbols are open for
// copying, leverage caps, and TP/SL ratios).
type GetMixTraderConfigQuerySymbolsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderConfigQuerySymbolsService(productType ProductType) *GetMixTraderConfigQuerySymbolsService {
	return &GetMixTraderConfigQuerySymbolsService{c: c, params: map[string]string{"productType": string(productType)}}
}

func (s *GetMixTraderConfigQuerySymbolsService) Do(ctx context.Context) ([]MixTraderSymbolConfig, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/config-query-symbols", s.params).WithSign()
	resp, err := request.Do[[]MixTraderSymbolConfig](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderSymbolConfig is one symbol's copy-trade configuration.
type MixTraderSymbolConfig struct {
	Symbol           string          `json:"symbol"`
	OpenTrader       string          `json:"openTrader"` // YES, NO
	MinOpenCount     string          `json:"minOpenCount"`
	MaxLeverage      string          `json:"maxLeverage"`
	StopSurplusRatio decimal.Decimal `json:"stopSurplusRatio"`
	StopLossRatio    decimal.Decimal `json:"stopLossRatio"`
}

// SetMixTraderConfigSymbolsService -- POST /api/v2/copy/mix-trader/config-setting-symbols (signed, state-changing)
//
// Adds, updates, or deletes the trader's copy-trade symbol settings (up to 50
// entries per call).
type SetMixTraderConfigSymbolsService struct {
	c       *CopyClient
	setting []MixTraderSymbolSetting
}

func (c *CopyClient) NewSetMixTraderConfigSymbolsService(settingList []MixTraderSymbolSetting) *SetMixTraderConfigSymbolsService {
	return &SetMixTraderConfigSymbolsService{c: c, setting: settingList}
}

func (s *SetMixTraderConfigSymbolsService) Do(ctx context.Context) (string, error) {
	body := map[string]any{"settingList": s.setting}
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/config-setting-symbols").SetBody(body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// MixTraderSymbolSetting is one symbol-setting change in a
// config-setting-symbols request.
type MixTraderSymbolSetting struct {
	Symbol           string          `json:"symbol"`
	ProductType      ProductType     `json:"productType"`
	SettingType      string          `json:"settingType"` // ADD, DELETE, UPDATE
	StopSurplusRatio decimal.Decimal `json:"stopSurplusRatio,omitzero"`
	StopLossRatio    decimal.Decimal `json:"stopLossRatio,omitzero"`
}

// SetMixTraderConfigBaseService -- POST /api/v2/copy/mix-trader/config-settings-base (signed, state-changing)
//
// Changes the trader's global copy-trade settings (enable elite trading, show
// total equity, show TP/SL). At least one field must be set.
type SetMixTraderConfigBaseService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewSetMixTraderConfigBaseService() *SetMixTraderConfigBaseService {
	return &SetMixTraderConfigBaseService{c: c, body: map[string]any{}}
}

func (s *SetMixTraderConfigBaseService) SetEnable(enable string) *SetMixTraderConfigBaseService {
	s.body["enable"] = enable
	return s
}

func (s *SetMixTraderConfigBaseService) SetShowTotalEquity(show string) *SetMixTraderConfigBaseService {
	s.body["showTotalEquity"] = show
	return s
}

func (s *SetMixTraderConfigBaseService) SetShowTpsl(show string) *SetMixTraderConfigBaseService {
	s.body["showTpsl"] = show
	return s
}

func (s *SetMixTraderConfigBaseService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/config-settings-base", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetMixTraderConfigQueryFollowersService -- GET /api/v2/copy/mix-trader/config-query-followers (signed)
//
// Returns the trader's current followers.
type GetMixTraderConfigQueryFollowersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixTraderConfigQueryFollowersService() *GetMixTraderConfigQueryFollowersService {
	return &GetMixTraderConfigQueryFollowersService{c: c, params: map[string]string{}}
}

func (s *GetMixTraderConfigQueryFollowersService) SetPageNo(pageNo int) *GetMixTraderConfigQueryFollowersService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetMixTraderConfigQueryFollowersService) SetPageSize(pageSize int) *GetMixTraderConfigQueryFollowersService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetMixTraderConfigQueryFollowersService) SetStartTime(t time.Time) *GetMixTraderConfigQueryFollowersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderConfigQueryFollowersService) SetEndTime(t time.Time) *GetMixTraderConfigQueryFollowersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetMixTraderConfigQueryFollowersService) Do(ctx context.Context) ([]MixTraderFollower, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-trader/config-query-followers", s.params).WithSign()
	resp, err := request.Do[[]MixTraderFollower](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixTraderFollower is one of the trader's followers.
type MixTraderFollower struct {
	AccountEquity   decimal.Decimal `json:"accountEquity"`
	IsRemove        string          `json:"isRemove"` // yes, no
	FollowerHeadPic string          `json:"followerHeadPic"`
	FollowerName    string          `json:"followerName"`
	FollowerUID     string          `json:"followerUid"`
	FollowerTime    time.Time       `json:"followerTime"`
}

// RemoveMixTraderFollowerService -- POST /api/v2/copy/mix-trader/config-remove-follower (signed, state-changing)
//
// Removes a follower from the trader's copy-trading relationship.
type RemoveMixTraderFollowerService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewRemoveMixTraderFollowerService(followerUID string) *RemoveMixTraderFollowerService {
	return &RemoveMixTraderFollowerService{c: c, body: map[string]any{"followerUid": followerUID}}
}

func (s *RemoveMixTraderFollowerService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/mix-trader/config-remove-follower", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}
