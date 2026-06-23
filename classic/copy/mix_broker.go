package copy

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetMixBrokerTradersService -- GET /api/v2/copy/mix-broker/query-traders (broker)
//
// Returns the platform's futures copy-trading traders the broker can surface to
// its followers, with each trader's aggregate performance metrics.
type GetMixBrokerTradersService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixBrokerTradersService() *GetMixBrokerTradersService {
	return &GetMixBrokerTradersService{c: c, params: map[string]string{}}
}

// SetTraderID filters to a single trader UID.
func (s *GetMixBrokerTradersService) SetTraderID(traderID string) *GetMixBrokerTradersService {
	s.params["traderId"] = traderID
	return s
}

// SetTraderName fuzzy-matches a trader nickname.
func (s *GetMixBrokerTradersService) SetTraderName(traderName string) *GetMixBrokerTradersService {
	s.params["traderName"] = traderName
	return s
}

// SetFullStatus filters by whether the trader's follower slots are full (Full
// for at-capacity traders, All for everyone; default All).
func (s *GetMixBrokerTradersService) SetFullStatus(fullStatus string) *GetMixBrokerTradersService {
	s.params["fullStatus"] = fullStatus
	return s
}

// SetSortRule orders the results (Composite, ROI, totalPL, AUM; default
// Composite).
func (s *GetMixBrokerTradersService) SetSortRule(sortRule string) *GetMixBrokerTradersService {
	s.params["sortRule"] = sortRule
	return s
}

// SetSortFlag selects ascending or descending order (Asc, Desc; default Desc).
func (s *GetMixBrokerTradersService) SetSortFlag(sortFlag string) *GetMixBrokerTradersService {
	s.params["sortFlag"] = sortFlag
	return s
}

// SetLanguage selects the locale of the columnList descriptions (zh-CN, en-US).
func (s *GetMixBrokerTradersService) SetLanguage(language string) *GetMixBrokerTradersService {
	s.params["language"] = language
	return s
}

// SetPageSize caps the number of traders returned (default 20, max 20).
func (s *GetMixBrokerTradersService) SetPageSize(pageSize int) *GetMixBrokerTradersService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo selects the page (default 1).
func (s *GetMixBrokerTradersService) SetPageNo(pageNo int) *GetMixBrokerTradersService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetMixBrokerTradersService) Do(ctx context.Context) ([]MixBrokerTrader, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-broker/query-traders", s.params).WithSign()
	resp, err := request.Do[[]MixBrokerTrader](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixBrokerTrader is one platform trader together with its performance metrics.
type MixBrokerTrader struct {
	CanTrace            string                      `json:"canTrace"`   // No, Yes
	TraderID            string                      `json:"traderId"`   //
	TraderName          string                      `json:"traderName"` //
	MaxLimit            string                      `json:"maxLimit"`   // max followers allowed
	BgbMaxFollowLimit   string                      `json:"bgbMaxFollowLimit"`
	FollowCount         string                      `json:"followCount"`    // current followers
	BgbFollowCount      string                      `json:"bgbFollowCount"` //
	TraderStatus        string                      `json:"traderStatus"`   // is this one of my traders: No, Yes
	TotalEquity         string                      `json:"totalEquity"`    // trader equity; Bitget masks it as "****" for privacy, so string not decimal
	CurrentTradingList  []string                    `json:"currentTradingList"`
	ColumnList          []MixBrokerTraderColumn     `json:"columnList"`
	TotalFollowers      string                      `json:"totalFollowers"`
	ProfitCount         string                      `json:"profitCount"`
	LossCount           string                      `json:"lossCount"`
	TradeCount          string                      `json:"tradeCount"`
	TraderPic           string                      `json:"traderPic"`
	MaxCallbackRate     decimal.Decimal             `json:"maxCallbackRate"` // max drawdown rate
	AverageWinRate      decimal.Decimal             `json:"averageWinRate"`
	DailyProfitRateList []MixBrokerTraderProfitRate `json:"dailyProfitRateList"` // 30-day win rate
	DailyProfitList     []MixBrokerTraderProfit     `json:"dailyProfitList"`     // 30-day profit
	ProfitRate24hList   []MixBrokerTraderProfitRate `json:"profitRate24hList"`   // 24h win rate
	Profit24hList       []MixBrokerTraderProfit     `json:"profit24hList"`       // 24h profit
	FollowerTotalProfit decimal.Decimal             `json:"followerTotalProfit"`
	LastTradeTime       time.Time                   `json:"lastTradeTime"`
	TradeDays           string                      `json:"tradeDays"`
}

// MixBrokerTraderColumn is one labelled performance metric. The value is a
// pre-formatted display string (e.g. "$12,376.16"), so it is kept as a string.
type MixBrokerTraderColumn struct {
	Describe string `json:"describe"`
	Value    string `json:"value"`
}

// MixBrokerTraderProfitRate is one (rate, time) point in a trader's win-rate
// series.
type MixBrokerTraderProfitRate struct {
	Rate  decimal.Decimal `json:"rate"`
	CTime time.Time       `json:"cTime"`
}

// MixBrokerTraderProfit is one (amount, time) point in a trader's profit series.
type MixBrokerTraderProfit struct {
	Amount decimal.Decimal `json:"amount"`
	CTime  time.Time       `json:"cTime"`
}

// GetMixBrokerHistoryTracesService -- GET /api/v2/copy/mix-broker/query-history-traces (broker)
//
// Returns a trader's closed (historical) copy-trading orders. Defaults to the
// last three months when the time window is omitted (max span three months).
type GetMixBrokerHistoryTracesService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixBrokerHistoryTracesService(traderID string, productType ProductType) *GetMixBrokerHistoryTracesService {
	return &GetMixBrokerHistoryTracesService{c: c, params: map[string]string{
		"traderId":    traderID,
		"productType": string(productType),
	}}
}

// SetSymbol filters to a single trading pair.
func (s *GetMixBrokerHistoryTracesService) SetSymbol(symbol string) *GetMixBrokerHistoryTracesService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters to traces created at or after t.
func (s *GetMixBrokerHistoryTracesService) SetStartTime(t time.Time) *GetMixBrokerHistoryTracesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to traces created at or before t.
func (s *GetMixBrokerHistoryTracesService) SetEndTime(t time.Time) *GetMixBrokerHistoryTracesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of traces returned (default 100, max 100).
func (s *GetMixBrokerHistoryTracesService) SetLimit(limit int) *GetMixBrokerHistoryTracesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards (older than the given endId).
func (s *GetMixBrokerHistoryTracesService) SetIDLessThan(id string) *GetMixBrokerHistoryTracesService {
	s.params["idLessThan"] = id
	return s
}

// SetIDGreaterThan pages forwards (newer than the given endId).
func (s *GetMixBrokerHistoryTracesService) SetIDGreaterThan(id string) *GetMixBrokerHistoryTracesService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixBrokerHistoryTracesService) Do(ctx context.Context) (*MixBrokerHistoryTraces, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-broker/query-history-traces", s.params).WithSign()
	return request.Do[MixBrokerHistoryTraces](req)
}

// MixBrokerHistoryTraces is the paged set of a trader's closed copy-trading
// orders.
type MixBrokerHistoryTraces struct {
	TrackingList []MixBrokerHistoryTrace `json:"trackingList"`
	EndID        string                  `json:"endId"` // cursor for idLessThan/idGreaterThan
}

// MixBrokerHistoryTrace is one closed copy-trading order. marginMode carries the
// wire values "isolated"/"cross" (note: "cross", not the order-side "crossed").
type MixBrokerHistoryTrace struct {
	TrackingNo    string          `json:"trackingNo"`   // trace (copy) order id
	OpenOrderID   string          `json:"openOrderId"`  //
	CloseOrderID  string          `json:"closeOrderId"` //
	MarginMode    string          `json:"marginMode"`   // isolated, cross
	PosSide       PosSide         `json:"posSide"`      // long, short
	Symbol        string          `json:"symbol"`
	OpenLeverage  string          `json:"openLeverage"`
	OpenPriceAvg  decimal.Decimal `json:"openPriceAvg"`
	OpenTime      time.Time       `json:"openTime"`
	OpenSize      decimal.Decimal `json:"openSize"`
	ClosePriceAvg decimal.Decimal `json:"closePriceAvg"`
	CloseTime     time.Time       `json:"closeTime"`
	CloseSize     decimal.Decimal `json:"closeSize"`
	OpenFee       decimal.Decimal `json:"openFee"`  // excludes discounts
	CloseFee      decimal.Decimal `json:"closeFee"` // excludes discounts
	MarginAmount  decimal.Decimal `json:"marginAmount"`
	FollowCount   string          `json:"followCount"` // followers on this order
	CTime         time.Time       `json:"cTime"`       // trace creation time
}

// GetMixBrokerCurrentTracesService -- GET /api/v2/copy/mix-broker/query-current-traces (broker)
//
// Returns a trader's open (current/pending) copy-trading orders.
type GetMixBrokerCurrentTracesService struct {
	c      *CopyClient
	params map[string]string
}

func (c *CopyClient) NewGetMixBrokerCurrentTracesService(traderID string, productType ProductType) *GetMixBrokerCurrentTracesService {
	return &GetMixBrokerCurrentTracesService{c: c, params: map[string]string{
		"traderId":    traderID,
		"productType": string(productType),
	}}
}

// SetSymbol filters to a single trading pair.
func (s *GetMixBrokerCurrentTracesService) SetSymbol(symbol string) *GetMixBrokerCurrentTracesService {
	s.params["symbol"] = symbol
	return s
}

// SetStartTime filters to traces created at or after t.
func (s *GetMixBrokerCurrentTracesService) SetStartTime(t time.Time) *GetMixBrokerCurrentTracesService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to traces created at or before t.
func (s *GetMixBrokerCurrentTracesService) SetEndTime(t time.Time) *GetMixBrokerCurrentTracesService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of traces returned (default 20, max 100).
func (s *GetMixBrokerCurrentTracesService) SetLimit(limit int) *GetMixBrokerCurrentTracesService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards (older than the given endId).
func (s *GetMixBrokerCurrentTracesService) SetIDLessThan(id string) *GetMixBrokerCurrentTracesService {
	s.params["idLessThan"] = id
	return s
}

// SetIDGreaterThan pages forwards (newer than the given endId).
func (s *GetMixBrokerCurrentTracesService) SetIDGreaterThan(id string) *GetMixBrokerCurrentTracesService {
	s.params["idGreaterThan"] = id
	return s
}

func (s *GetMixBrokerCurrentTracesService) Do(ctx context.Context) ([]MixBrokerCurrentTrace, error) {
	req := request.Get(ctx, s.c, "/api/v2/copy/mix-broker/query-current-traces", s.params).WithSign()
	resp, err := request.Do[[]MixBrokerCurrentTrace](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MixBrokerCurrentTrace is one open copy-trading order. marginMode carries the
// wire values "isolated"/"cross" (note: "cross", not the order-side "crossed").
type MixBrokerCurrentTrace struct {
	TrackingNo       string          `json:"trackingNo"`  // trace (copy) order id
	OpenOrderID      string          `json:"openOrderId"` //
	MarginMode       string          `json:"marginMode"`  // isolated, cross
	PosSide          PosSide         `json:"posSide"`     // long, short
	Symbol           string          `json:"symbol"`
	OpenLeverage     string          `json:"openLeverage"`
	OpenPriceAvg     decimal.Decimal `json:"openPriceAvg"`
	OpenTime         time.Time       `json:"openTime"`
	OpenSize         decimal.Decimal `json:"openSize"`
	OpenFee          decimal.Decimal `json:"openFee"` // USDT only, excludes discounts
	MarginAmount     decimal.Decimal `json:"marginAmount"`
	FollowCount      string          `json:"followCount"`      // followers on this order
	StopSurplusPrice decimal.Decimal `json:"stopSurplusPrice"` // take-profit price
	StopLossPrice    decimal.Decimal `json:"stopLossPrice"`    // stop-loss price
	CTime            time.Time       `json:"cTime"`            // trace creation time
}
