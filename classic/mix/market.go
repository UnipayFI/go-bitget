package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetMergeDepthService -- GET /api/v2/mix/market/merge-depth (public)
//
// Returns the merged (price-aggregated) order book depth for a contract.
type GetMergeDepthService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetMergeDepthService(symbol string, productType ProductType) *GetMergeDepthService {
	return &GetMergeDepthService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetPrecision selects the price-aggregation step (scale0..scale3, default
// scale0).
func (s *GetMergeDepthService) SetPrecision(precision string) *GetMergeDepthService {
	s.params["precision"] = precision
	return s
}

// SetLimit caps the number of depth levels returned (1, 5, 15, 50, max, default
// 100).
func (s *GetMergeDepthService) SetLimit(limit string) *GetMergeDepthService {
	s.params["limit"] = limit
	return s
}

func (s *GetMergeDepthService) Do(ctx context.Context) (*MergeDepth, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/merge-depth", s.params)
	return request.Do[MergeDepth](req)
}

// MergeDepth is the merged order book snapshot. Asks and bids arrive as arrays
// of [price, size] pairs (bare JSON numbers).
type MergeDepth struct {
	Asks           [][]decimal.Decimal `json:"asks"`
	Bids           [][]decimal.Decimal `json:"bids"`
	Ts             time.Time           `json:"ts"`
	Scale          decimal.Decimal     `json:"scale"`
	Precision      string              `json:"precision"`
	IsMaxPrecision string              `json:"isMaxPrecision"` // YES, NO
}

// GetTickerService -- GET /api/v2/mix/market/ticker (public)
//
// Returns the 24h ticker statistics for a single contract.
type GetTickerService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetTickerService(symbol string, productType ProductType) *GetTickerService {
	return &GetTickerService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetTickerService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/ticker", s.params)
	resp, err := request.Do[[]Ticker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetAllTickersService -- GET /api/v2/mix/market/tickers (public)
//
// Returns 24h ticker statistics for every contract in a product type.
type GetAllTickersService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetAllTickersService(productType ProductType) *GetAllTickersService {
	return &GetAllTickersService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetAllTickersService) Do(ctx context.Context) ([]Ticker, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/tickers", s.params)
	resp, err := request.Do[[]Ticker](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Ticker is one contract's 24h ticker snapshot.
type Ticker struct {
	Symbol            string          `json:"symbol"`
	LastPr            decimal.Decimal `json:"lastPr"`
	AskPr             decimal.Decimal `json:"askPr"`
	BidPr             decimal.Decimal `json:"bidPr"`
	BidSz             decimal.Decimal `json:"bidSz"`
	AskSz             decimal.Decimal `json:"askSz"`
	High24h           decimal.Decimal `json:"high24h"`
	Low24h            decimal.Decimal `json:"low24h"`
	Ts                time.Time       `json:"ts"`
	Change24h         decimal.Decimal `json:"change24h"`
	BaseVolume        decimal.Decimal `json:"baseVolume"`
	QuoteVolume       decimal.Decimal `json:"quoteVolume"`
	USDTVolume        decimal.Decimal `json:"usdtVolume"`
	OpenUtc           decimal.Decimal `json:"openUtc"`
	ChangeUtc24h      decimal.Decimal `json:"changeUtc24h"`
	IndexPrice        decimal.Decimal `json:"indexPrice"`
	FundingRate       decimal.Decimal `json:"fundingRate"`
	HoldingAmount     decimal.Decimal `json:"holdingAmount"`
	DeliveryStartTime time.Time       `json:"deliveryStartTime"`
	DeliveryTime      time.Time       `json:"deliveryTime"`
	DeliveryStatus    string          `json:"deliveryStatus"`
	Open24h           decimal.Decimal `json:"open24h"`
	MarkPrice         decimal.Decimal `json:"markPrice"`
}

// GetRecentFillsService -- GET /api/v2/mix/market/fills (public)
//
// Returns the most recent public trades (fills) for a contract.
type GetRecentFillsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetRecentFillsService(symbol string, productType ProductType) *GetRecentFillsService {
	return &GetRecentFillsService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetLimit caps the number of fills returned (default 100, max 1000).
func (s *GetRecentFillsService) SetLimit(limit int) *GetRecentFillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetRecentFillsService) Do(ctx context.Context) ([]MarketFill, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/fills", s.params)
	resp, err := request.Do[[]MarketFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetFillsHistoryService -- GET /api/v2/mix/market/fills-history (public)
//
// Returns historical public trades (fills) for a contract, paged by trade id /
// time window.
type GetFillsHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetFillsHistoryService(symbol string, productType ProductType) *GetFillsHistoryService {
	return &GetFillsHistoryService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

// SetLimit caps the number of fills returned (default 500, max 1000).
func (s *GetFillsHistoryService) SetLimit(limit int) *GetFillsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan returns fills with a trade id smaller than tradeID (older
// trades).
func (s *GetFillsHistoryService) SetIDLessThan(tradeID string) *GetFillsHistoryService {
	s.params["idLessThan"] = tradeID
	return s
}

// SetStartTime filters fills at or after t.
func (s *GetFillsHistoryService) SetStartTime(t time.Time) *GetFillsHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters fills at or before t.
func (s *GetFillsHistoryService) SetEndTime(t time.Time) *GetFillsHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetFillsHistoryService) Do(ctx context.Context) ([]MarketFill, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/fills-history", s.params)
	resp, err := request.Do[[]MarketFill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarketFill is a single public trade. The side casing differs across endpoints
// ("buy"/"sell" on fills, "Buy"/"Sell" on fills-history), so it is kept as a
// plain string rather than the Side enum.
type MarketFill struct {
	TradeID string          `json:"tradeId"`
	Price   decimal.Decimal `json:"price"`
	Size    decimal.Decimal `json:"size"`
	Side    string          `json:"side"`
	Ts      time.Time       `json:"ts"`
	Symbol  string          `json:"symbol"`
}

// GetOpenInterestService -- GET /api/v2/mix/market/open-interest (public)
//
// Returns the total open interest (open position size) for a contract.
type GetOpenInterestService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOpenInterestService(symbol string, productType ProductType) *GetOpenInterestService {
	return &GetOpenInterestService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetOpenInterestService) Do(ctx context.Context) (*OpenInterest, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/open-interest", s.params)
	return request.Do[OpenInterest](req)
}

// OpenInterest is the open-interest snapshot for the queried product type.
type OpenInterest struct {
	OpenInterestList []OpenInterestItem `json:"openInterestList"`
	Ts               time.Time          `json:"ts"`
}

// OpenInterestItem is one contract's open interest (in base coin).
type OpenInterestItem struct {
	Symbol string          `json:"symbol"`
	Size   decimal.Decimal `json:"size"`
}

// GetFundingTimeService -- GET /api/v2/mix/market/funding-time (public)
//
// Returns the next funding settlement time and funding interval for a contract.
type GetFundingTimeService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetFundingTimeService(symbol string, productType ProductType) *GetFundingTimeService {
	return &GetFundingTimeService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetFundingTimeService) Do(ctx context.Context) ([]FundingTime, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/funding-time", s.params)
	resp, err := request.Do[[]FundingTime](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FundingTime is the next funding settlement for a contract.
type FundingTime struct {
	Symbol          string    `json:"symbol"`
	NextFundingTime time.Time `json:"nextFundingTime"`
	RatePeriod      string    `json:"ratePeriod"` // funding interval in hours
}

// GetSymbolPriceService -- GET /api/v2/mix/market/symbol-price (public)
//
// Returns the latest traded, index and mark prices for a contract.
type GetSymbolPriceService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetSymbolPriceService(symbol string, productType ProductType) *GetSymbolPriceService {
	return &GetSymbolPriceService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetSymbolPriceService) Do(ctx context.Context) ([]SymbolPrice, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/symbol-price", s.params)
	resp, err := request.Do[[]SymbolPrice](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SymbolPrice is the price snapshot (traded, index, mark) for a contract.
type SymbolPrice struct {
	Symbol     string          `json:"symbol"`
	Price      decimal.Decimal `json:"price"`
	IndexPrice decimal.Decimal `json:"indexPrice"`
	MarkPrice  decimal.Decimal `json:"markPrice"`
	Ts         time.Time       `json:"ts"`
}

// GetOiLimitService -- GET /api/v2/mix/market/oi-limit (public)
//
// Returns the per-contract and product-line open-interest (notional) limits.
type GetOiLimitService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOiLimitService(symbol string, productType ProductType) *GetOiLimitService {
	return &GetOiLimitService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetOiLimitService) Do(ctx context.Context) ([]OiLimit, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/oi-limit", s.params)
	resp, err := request.Do[[]OiLimit](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// OiLimit is one contract's open-interest notional limit.
type OiLimit struct {
	Symbol             string          `json:"symbol"`
	NotionalValue      decimal.Decimal `json:"notionalValue"`
	TotalNotionalValue decimal.Decimal `json:"totalNotionalValue"`
}

// GetContractsService -- GET /api/v2/mix/market/contracts (public)
//
// Returns the trading configuration (precision, limits, fee rates, leverage) of
// the contracts in a product type, optionally filtered to a single symbol.
type GetContractsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetContractsService(productType ProductType) *GetContractsService {
	return &GetContractsService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

// SetSymbol filters the result to a single contract.
func (s *GetContractsService) SetSymbol(symbol string) *GetContractsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetContractsService) Do(ctx context.Context) ([]Contract, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/contracts", s.params)
	resp, err := request.Do[[]Contract](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Contract is one contract's trading configuration.
type Contract struct {
	Symbol              string          `json:"symbol"`
	BaseCoin            string          `json:"baseCoin"`
	QuoteCoin           string          `json:"quoteCoin"`
	BuyLimitPriceRatio  decimal.Decimal `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio decimal.Decimal `json:"sellLimitPriceRatio"`
	FeeRateUpRatio      decimal.Decimal `json:"feeRateUpRatio"`
	MakerFeeRate        decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate        decimal.Decimal `json:"takerFeeRate"`
	OpenCostUpRatio     decimal.Decimal `json:"openCostUpRatio"`
	SupportMarginCoins  []string        `json:"supportMarginCoins"`
	MinTradeNum         decimal.Decimal `json:"minTradeNum"`
	PriceEndStep        string          `json:"priceEndStep"`
	VolumePlace         string          `json:"volumePlace"`
	PricePlace          string          `json:"pricePlace"`
	SizeMultiplier      decimal.Decimal `json:"sizeMultiplier"`
	SymbolType          string          `json:"symbolType"` // perpetual, delivery
	MinTradeUSDT        decimal.Decimal `json:"minTradeUSDT"`
	MaxSymbolOrderNum   string          `json:"maxSymbolOrderNum"`
	MaxProductOrderNum  string          `json:"maxProductOrderNum"`
	MaxPositionNum      string          `json:"maxPositionNum"`
	SymbolStatus        string          `json:"symbolStatus"` // listed, normal, maintain, limit_open, restrictedAPI, off
	OffTime             time.Time       `json:"offTime"`
	LimitOpenTime       time.Time       `json:"limitOpenTime"`
	DeliveryTime        time.Time       `json:"deliveryTime"`
	DeliveryStartTime   time.Time       `json:"deliveryStartTime"`
	DeliveryPeriod      string          `json:"deliveryPeriod"`
	LaunchTime          time.Time       `json:"launchTime"`
	FundInterval        string          `json:"fundInterval"`
	MinLever            string          `json:"minLever"`
	MaxLever            string          `json:"maxLever"`
	PosLimit            decimal.Decimal `json:"posLimit"`
	MaintainTime        time.Time       `json:"maintainTime"`
	OpenTime            time.Time       `json:"openTime"`
	MaxMarketOrderQty   decimal.Decimal `json:"maxMarketOrderQty"`
	MaxOrderQty         decimal.Decimal `json:"maxOrderQty"`
	IsRwa               string          `json:"isRwa"` // YES, NO
}

// GetQueryPositionLeverService -- GET /api/v2/mix/market/query-position-lever (public)
//
// Returns the leverage tier (position tier) table for a contract: per-tier
// notional bands, max leverage and maintenance margin rate.
type GetQueryPositionLeverService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetQueryPositionLeverService(symbol string, productType ProductType) *GetQueryPositionLeverService {
	return &GetQueryPositionLeverService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
	}}
}

func (s *GetQueryPositionLeverService) Do(ctx context.Context) ([]PositionLever, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/market/query-position-lever", s.params)
	resp, err := request.Do[[]PositionLever](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// PositionLever is one leverage tier for a contract.
type PositionLever struct {
	Symbol         string          `json:"symbol"`
	Level          string          `json:"level"`
	StartUnit      decimal.Decimal `json:"startUnit"`
	EndUnit        decimal.Decimal `json:"endUnit"`
	Leverage       string          `json:"leverage"`
	KeepMarginRate decimal.Decimal `json:"keepMarginRate"`
}
