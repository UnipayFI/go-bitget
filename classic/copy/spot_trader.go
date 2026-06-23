package copy

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSpotTraderProfitSummarysService -- GET /api/v2/copy/spot-trader/profit-summarys (private)
//
// Returns the spot elite trader's profit-sharing overview together with the
// per-coin profit-distribution history.
type GetSpotTraderProfitSummarysService struct {
	c *CopyClient
}

func (c *CopyClient) NewGetSpotTraderProfitSummarysService() *GetSpotTraderProfitSummarysService {
	return &GetSpotTraderProfitSummarysService{c: c}
}

func (s *GetSpotTraderProfitSummarysService) Do(ctx context.Context) (*SpotTraderProfitSummary, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/profit-summarys").WithSign()
	return request.Do[SpotTraderProfitSummary](req)
}

// SpotTraderProfitSummary is the profit-sharing overview returned by
// profit-summarys.
type SpotTraderProfitSummary struct {
	ProfitSummarys    SpotTraderProfitSummaryOverview `json:"profitSummarys"`
	ProfitHistoryList []SpotTraderProfitHistory       `json:"profitHistoryList"`
}

// SpotTraderProfitSummaryOverview is the aggregate profit-sharing snapshot.
type SpotTraderProfitSummaryOverview struct {
	YesterdayProfit decimal.Decimal `json:"yesterdayProfit"`
	YesterdayTime   time.Time       `json:"yesterdayTime"`
	SumProfit       decimal.Decimal `json:"sumProfit"`
	WaitProfit      decimal.Decimal `json:"waitProfit"`
}

// SpotTraderProfitHistory is a per-coin profit-distribution history bucket.
type SpotTraderProfitHistory struct {
	Coin               string                          `json:"coin"`
	ProfitCount        decimal.Decimal                 `json:"profitCount"`
	LastProfitTime     time.Time                       `json:"lastProfitTime"`
	HistorysByDateList []SpotTraderProfitHistoryByDate `json:"historysByDateList"`
}

// SpotTraderProfitHistoryByDate is a single dated profit-distribution entry.
type SpotTraderProfitHistoryByDate struct {
	Profit     decimal.Decimal `json:"profit"`
	ProfitTime time.Time       `json:"profitTime"`
}

// GetSpotTraderProfitHistoryDetailsService -- GET /api/v2/copy/spot-trader/profit-history-details (private)
//
// Returns the spot elite trader's settled (historical) profit-sharing records.
type GetSpotTraderProfitHistoryDetailsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotTraderProfitHistoryDetailsService() *GetSpotTraderProfitHistoryDetailsService {
	return &GetSpotTraderProfitHistoryDetailsService{c: c, params: map[string]string{}}
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetIdLessThan(id string) *GetSpotTraderProfitHistoryDetailsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetIdGreaterThan(id string) *GetSpotTraderProfitHistoryDetailsService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetStartTime(t time.Time) *GetSpotTraderProfitHistoryDetailsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetEndTime(t time.Time) *GetSpotTraderProfitHistoryDetailsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetLimit(limit int) *GetSpotTraderProfitHistoryDetailsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) SetCoin(coin string) *GetSpotTraderProfitHistoryDetailsService {
	s.params["coin"] = coin
	return s
}

func (s *GetSpotTraderProfitHistoryDetailsService) Do(ctx context.Context) (*SpotTraderProfitHistoryDetails, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/profit-history-details", s.params).WithSign()
	return request.Do[SpotTraderProfitHistoryDetails](req)
}

// SpotTraderProfitHistoryDetails is the paginated settled profit-sharing list.
type SpotTraderProfitHistoryDetails struct {
	EndID      string                          `json:"endId"`
	ProfitList []SpotTraderProfitHistoryDetail `json:"profitList"`
}

// SpotTraderProfitHistoryDetail is a single settled profit-sharing record.
type SpotTraderProfitHistoryDetail struct {
	ProfitID        string          `json:"profitId"`
	Coin            string          `json:"coin"`
	DistributeRatio decimal.Decimal `json:"distributeRatio"`
	Profit          decimal.Decimal `json:"profit"`
	FollowerName    string          `json:"followerName"`
	ProfitTime      time.Time       `json:"profitTime"`
}

// GetSpotTraderProfitDetailsService -- GET /api/v2/copy/spot-trader/profit-details (private)
//
// Returns the spot elite trader's unrealized (pending) profit-sharing details.
type GetSpotTraderProfitDetailsService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotTraderProfitDetailsService() *GetSpotTraderProfitDetailsService {
	return &GetSpotTraderProfitDetailsService{c: c, params: map[string]string{}}
}

func (s *GetSpotTraderProfitDetailsService) SetCoin(coin string) *GetSpotTraderProfitDetailsService {
	s.params["coin"] = coin
	return s
}

func (s *GetSpotTraderProfitDetailsService) SetPageNo(pageNo int) *GetSpotTraderProfitDetailsService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetSpotTraderProfitDetailsService) SetPageSize(pageSize int) *GetSpotTraderProfitDetailsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetSpotTraderProfitDetailsService) Do(ctx context.Context) ([]SpotTraderProfitDetail, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/profit-details", s.params).WithSign()
	resp, err := request.Do[[]SpotTraderProfitDetail](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotTraderProfitDetail is a single unrealized profit-sharing entry.
type SpotTraderProfitDetail struct {
	DistributeRatio decimal.Decimal `json:"distributeRatio"`
	Coin            string          `json:"coin"`
	Profit          decimal.Decimal `json:"profit"`
	FollowerName    string          `json:"followerName"`
}

// GetSpotTraderOrderTotalDetailService -- GET /api/v2/copy/spot-trader/order-total-detail (private)
//
// Returns the spot elite trader's aggregate data-indicator statistics
// (followers, win rate, equity, and rolling ROI / profit series).
type GetSpotTraderOrderTotalDetailService struct {
	c *CopyClient
}

func (c *CopyClient) NewGetSpotTraderOrderTotalDetailService() *GetSpotTraderOrderTotalDetailService {
	return &GetSpotTraderOrderTotalDetailService{c: c}
}

func (s *GetSpotTraderOrderTotalDetailService) Do(ctx context.Context) (*SpotTraderOrderTotalDetail, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/order-total-detail").WithSign()
	return request.Do[SpotTraderOrderTotalDetail](req)
}

// SpotTraderOrderTotalDetail is the trader's statistics dashboard payload.
type SpotTraderOrderTotalDetail struct {
	TotalFollowerNum    string                  `json:"totalFollowerNum"`
	CurrentFollowerNum  string                  `json:"currentFollowerNum"`
	MaxFollowerNum      string                  `json:"maxFollowerNum"`
	TradingOrderNum     string                  `json:"tradingOrderNum"`
	Totalpl             decimal.Decimal         `json:"totalpl"`
	GainNum             string                  `json:"gainNum"`
	LossNum             string                  `json:"lossNum"`
	TotalEquity         decimal.Decimal         `json:"totalEquity"`
	WinRate             decimal.Decimal         `json:"winRate"`
	LastWeekRoiList     []SpotTraderRoiPoint    `json:"lastWeekRoiList"`
	LastMonthRoiList    []SpotTraderRoiPoint    `json:"lastMonthRoiList"`
	LastWeekProfitList  []SpotTraderProfitPoint `json:"lastWeekProfitList"`
	LastMonthProfitList []SpotTraderProfitPoint `json:"lastMonthProfitList"`
}

// SpotTraderRoiPoint is a single timestamped ROI sample in a rolling series.
type SpotTraderRoiPoint struct {
	Rate  decimal.Decimal `json:"rate"`
	Ctime time.Time       `json:"ctime"`
}

// SpotTraderProfitPoint is a single timestamped profit sample in a rolling
// series.
type SpotTraderProfitPoint struct {
	Amount decimal.Decimal `json:"amount"`
	Ctime  time.Time       `json:"ctime"`
}

// ModifySpotTraderOrderTpslService -- POST /api/v2/copy/spot-trader/order-modify-tpsl (private, state-changing)
//
// Modifies the take-profit and/or stop-loss prices of a spot tracking order.
type ModifySpotTraderOrderTpslService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewModifySpotTraderOrderTpslService(trackingNo string) *ModifySpotTraderOrderTpslService {
	return &ModifySpotTraderOrderTpslService{c: c, body: map[string]any{"trackingNo": trackingNo}}
}

func (s *ModifySpotTraderOrderTpslService) SetStopSurplusPrice(price decimal.Decimal) *ModifySpotTraderOrderTpslService {
	s.body["stopSurplusPrice"] = price.String()
	return s
}

func (s *ModifySpotTraderOrderTpslService) SetStopLossPrice(price decimal.Decimal) *ModifySpotTraderOrderTpslService {
	s.body["stopLossPrice"] = price.String()
	return s
}

func (s *ModifySpotTraderOrderTpslService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-trader/order-modify-tpsl", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetSpotTraderOrderHistoryTrackService -- GET /api/v2/copy/spot-trader/order-history-track (private)
//
// Returns the spot elite trader's closed (historical) tracking orders.
type GetSpotTraderOrderHistoryTrackService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotTraderOrderHistoryTrackService() *GetSpotTraderOrderHistoryTrackService {
	return &GetSpotTraderOrderHistoryTrackService{c: c, params: map[string]string{}}
}

func (s *GetSpotTraderOrderHistoryTrackService) SetIdLessThan(id string) *GetSpotTraderOrderHistoryTrackService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) SetIdGreaterThan(id string) *GetSpotTraderOrderHistoryTrackService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) SetStartTime(t time.Time) *GetSpotTraderOrderHistoryTrackService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) SetEndTime(t time.Time) *GetSpotTraderOrderHistoryTrackService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) SetLimit(limit int) *GetSpotTraderOrderHistoryTrackService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) SetSymbol(symbol string) *GetSpotTraderOrderHistoryTrackService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetSpotTraderOrderHistoryTrackService) Do(ctx context.Context) (*SpotTraderOrderHistoryTrack, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/order-history-track", s.params).WithSign()
	return request.Do[SpotTraderOrderHistoryTrack](req)
}

// SpotTraderOrderHistoryTrack is the paginated closed tracking-order list.
type SpotTraderOrderHistoryTrack struct {
	EndID        string                           `json:"endId"`
	TrackingList []SpotTraderHistoryTrackingOrder `json:"trackingList"`
}

// SpotTraderHistoryTrackingOrder is a single closed spot tracking order.
type SpotTraderHistoryTrackingOrder struct {
	TrackingNo  string          `json:"trackingNo"`
	FillSize    decimal.Decimal `json:"fillSize"`
	BuyPrice    decimal.Decimal `json:"buyPrice"`
	SellPrice   decimal.Decimal `json:"sellPrice"`
	AchievedPL  decimal.Decimal `json:"achievedPL"`
	BuyTime     time.Time       `json:"buyTime"`
	SellTime    time.Time       `json:"sellTime"`
	BuyFee      decimal.Decimal `json:"buyFee"`
	SellFee     decimal.Decimal `json:"sellFee"`
	AchievedPLR decimal.Decimal `json:"achievedPLR"`
	Symbol      string          `json:"symbol"`
	NetProfit   decimal.Decimal `json:"netProfit"`
	FollowCount string          `json:"followCount"`
}

// GetSpotTraderOrderCurrentTrackService -- GET /api/v2/copy/spot-trader/order-current-track (private)
//
// Returns the spot elite trader's open (current) tracking orders.
type GetSpotTraderOrderCurrentTrackService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotTraderOrderCurrentTrackService() *GetSpotTraderOrderCurrentTrackService {
	return &GetSpotTraderOrderCurrentTrackService{c: c, params: map[string]string{}}
}

func (s *GetSpotTraderOrderCurrentTrackService) SetSymbol(symbol string) *GetSpotTraderOrderCurrentTrackService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) SetIdLessThan(id string) *GetSpotTraderOrderCurrentTrackService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) SetIdGreaterThan(id string) *GetSpotTraderOrderCurrentTrackService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) SetStartTime(t time.Time) *GetSpotTraderOrderCurrentTrackService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) SetEndTime(t time.Time) *GetSpotTraderOrderCurrentTrackService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) SetLimit(limit int) *GetSpotTraderOrderCurrentTrackService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotTraderOrderCurrentTrackService) Do(ctx context.Context) (*SpotTraderOrderCurrentTrack, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/order-current-track", s.params).WithSign()
	return request.Do[SpotTraderOrderCurrentTrack](req)
}

// SpotTraderOrderCurrentTrack is the paginated open tracking-order list.
type SpotTraderOrderCurrentTrack struct {
	EndID        string                           `json:"endId"`
	TrackingList []SpotTraderCurrentTrackingOrder `json:"trackingList"`
}

// SpotTraderCurrentTrackingOrder is a single open spot tracking order.
type SpotTraderCurrentTrackingOrder struct {
	TrackingNo       string          `json:"trackingNo"`
	OrderID          string          `json:"orderId"`
	BuyFillSize      decimal.Decimal `json:"buyFillSize"`
	BuyDelegateSize  decimal.Decimal `json:"buyDelegateSize"`
	BuyPrice         decimal.Decimal `json:"buyPrice"`
	UnrealizedPL     decimal.Decimal `json:"unrealizedPL"`
	BuyTime          time.Time       `json:"buyTime"`
	BuyFee           decimal.Decimal `json:"buyFee"`
	UnrealizedPLR    decimal.Decimal `json:"unrealizedPLR"`
	Symbol           string          `json:"symbol"`
	StopLossPrice    decimal.Decimal `json:"stopLossPrice"`
	StopSurplusPrice decimal.Decimal `json:"stopSurplusPrice"`
	FollowCount      string          `json:"followCount"`
}

// CloseSpotTraderOrderTrackingService -- POST /api/v2/copy/spot-trader/order-close-tracking (private, state-changing)
//
// Closes (sells) one or more spot tracking orders for a single trading pair, in
// batch. All succeed or all fail.
type CloseSpotTraderOrderTrackingService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewCloseSpotTraderOrderTrackingService(symbol string, trackingNoList []string) *CloseSpotTraderOrderTrackingService {
	return &CloseSpotTraderOrderTrackingService{c: c, body: map[string]any{
		"symbol":         symbol,
		"trackingNoList": trackingNoList,
	}}
}

func (s *CloseSpotTraderOrderTrackingService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-trader/order-close-tracking", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// SpotTraderSettingType is the operation applied by config-setting-symbols: add
// or delete the listed trading pairs from the trader's copy-trade symbol set.
type SpotTraderSettingType string

const (
	SpotTraderSettingTypeAdd    SpotTraderSettingType = "add"
	SpotTraderSettingTypeDelete SpotTraderSettingType = "delete"
)

// SetSpotTraderConfigSymbolsService -- POST /api/v2/copy/spot-trader/config-setting-symbols (private, state-changing)
//
// Adds or removes trading pairs from the spot elite trader's copy-trade symbol
// configuration.
type SetSpotTraderConfigSymbolsService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewSetSpotTraderConfigSymbolsService(symbolList []string, settingType SpotTraderSettingType) *SetSpotTraderConfigSymbolsService {
	return &SetSpotTraderConfigSymbolsService{c: c, body: map[string]any{
		"symbolList":  symbolList,
		"settingType": string(settingType),
	}}
}

func (s *SetSpotTraderConfigSymbolsService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-trader/config-setting-symbols", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// RemoveSpotTraderFollowerService -- POST /api/v2/copy/spot-trader/config-remove-follower (private, state-changing)
//
// Removes a follower from the spot elite trader's follower list.
type RemoveSpotTraderFollowerService struct {
	c    *CopyClient
	body map[string]any
}

func (c *CopyClient) NewRemoveSpotTraderFollowerService(followerUid string) *RemoveSpotTraderFollowerService {
	return &RemoveSpotTraderFollowerService{c: c, body: map[string]any{"followerUid": followerUid}}
}

func (s *RemoveSpotTraderFollowerService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/copy/spot-trader/config-remove-follower", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetSpotTraderConfigSettingsService -- GET /api/v2/copy/spot-trader/config-query-settings (private)
//
// Returns the spot elite trader's copy-trade configuration: per-symbol trading
// limits, labels, visibility toggles, and the enabled trace symbol set.
type GetSpotTraderConfigSettingsService struct {
	c *CopyClient
}

func (c *CopyClient) NewGetSpotTraderConfigSettingsService() *GetSpotTraderConfigSettingsService {
	return &GetSpotTraderConfigSettingsService{c: c}
}

func (s *GetSpotTraderConfigSettingsService) Do(ctx context.Context) (*SpotTraderConfigSettings, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/config-query-settings").WithSign()
	return request.Do[SpotTraderConfigSettings](req)
}

// SpotTraderConfigSettings is the trader's copy-trade configuration payload.
type SpotTraderConfigSettings struct {
	RemoveLimitUsdt decimal.Decimal         `json:"removeLimitUsdt"`
	SpotInfoList    []SpotTraderSymbolInfo  `json:"spotInfoList"`
	LabelList       []SpotTraderLabel       `json:"labelList"`
	Enable          string                  `json:"enable"`
	ShowAssetsMap   string                  `json:"showAssetsMap"`
	ShowEquity      string                  `json:"showEquity"`
	TraceSymbolList []SpotTraderTraceSymbol `json:"traceSymbolList"`
}

// SpotTraderSymbolInfo is a per-symbol trading-amount entry.
type SpotTraderSymbolInfo struct {
	MaxQuoteSize     decimal.Decimal `json:"maxQuoteSize"`
	SurplusQuoteSize decimal.Decimal `json:"surplusQuoteSize"`
	Symbol           string          `json:"symbol"`
}

// SpotTraderLabel is a single copy-trade label/tag.
type SpotTraderLabel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SpotTraderTraceSymbol is a single trace-enabled trading pair entry.
type SpotTraderTraceSymbol struct {
	Enable       string          `json:"enable"`
	Symbol       string          `json:"symbol"`
	MinOpenCount decimal.Decimal `json:"minOpenCount"`
}

// GetSpotTraderConfigFollowersService -- GET /api/v2/copy/spot-trader/config-query-followers (private)
//
// Returns the spot elite trader's follower list.
type GetSpotTraderConfigFollowersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetSpotTraderConfigFollowersService() *GetSpotTraderConfigFollowersService {
	return &GetSpotTraderConfigFollowersService{c: c, params: map[string]string{}}
}

func (s *GetSpotTraderConfigFollowersService) SetPageNo(pageNo int) *GetSpotTraderConfigFollowersService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetSpotTraderConfigFollowersService) SetPageSize(pageSize int) *GetSpotTraderConfigFollowersService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

func (s *GetSpotTraderConfigFollowersService) SetStartTime(t time.Time) *GetSpotTraderConfigFollowersService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderConfigFollowersService) SetEndTime(t time.Time) *GetSpotTraderConfigFollowersService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSpotTraderConfigFollowersService) Do(ctx context.Context) ([]SpotTraderFollower, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/spot-trader/config-query-followers", s.params).WithSign()
	resp, err := request.Do[[]SpotTraderFollower](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotTraderFollower is a single follower entry in the trader's follower list.
type SpotTraderFollower struct {
	AccountEquity   decimal.Decimal `json:"accountEquity"`
	IsRemove        string          `json:"isRemove"`
	FollowerHeadPic string          `json:"followerHeadPic"`
	FollowerName    string          `json:"followerName"`
	FollowerUid     string          `json:"followerUid"`
	FollowerTime    time.Time       `json:"followerTime"`
}
