package common

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetConvertCurrenciesService -- GET /api/v2/convert/currencies (convert read)
//
// Returns the coins available for flash conversion (Convert) along with the
// account's available balance and the per-coin min/max convertible amounts.
type GetConvertCurrenciesService struct {
	c *CommonClient
}

func (c *CommonClient) NewGetConvertCurrenciesService() *GetConvertCurrenciesService {
	return &GetConvertCurrenciesService{c: c}
}

func (s *GetConvertCurrenciesService) Do(ctx context.Context) ([]ConvertCurrency, error) {
	req := request.Get(ctx, s.c, "/api/v2/convert/currencies").WithSign()
	resp, err := request.Do[[]ConvertCurrency](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// ConvertCurrency is one coin that can be used in a Convert (flash) trade.
type ConvertCurrency struct {
	Coin      string          `json:"coin"`
	Available decimal.Decimal `json:"available"`
	MaxAmount decimal.Decimal `json:"maxAmount"`
	MinAmount decimal.Decimal `json:"minAmount"`
}

// GetConvertQuotedPriceService -- GET /api/v2/convert/quoted-price (convert read)
//
// Requests a quote (RFQ) for converting fromCoin into toCoin. Supply exactly one
// of fromCoinSize / toCoinSize to fix the side being sized; the quote returns the
// traceId that the subsequent Convert (trade) call must echo back within 8s.
type GetConvertQuotedPriceService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetConvertQuotedPriceService(fromCoin, toCoin string) *GetConvertQuotedPriceService {
	return &GetConvertQuotedPriceService{c: c, params: map[string]string{
		"fromCoin": fromCoin,
		"toCoin":   toCoin,
	}}
}

// SetFromCoinSize sets the amount of the quote (source) coin to convert; only
// valid together with toCoinSize.
func (s *GetConvertQuotedPriceService) SetFromCoinSize(fromCoinSize decimal.Decimal) *GetConvertQuotedPriceService {
	s.params["fromCoinSize"] = fromCoinSize.String()
	return s
}

// SetToCoinSize sets the amount of the target coin to receive; only valid
// together with fromCoinSize.
func (s *GetConvertQuotedPriceService) SetToCoinSize(toCoinSize decimal.Decimal) *GetConvertQuotedPriceService {
	s.params["toCoinSize"] = toCoinSize.String()
	return s
}

func (s *GetConvertQuotedPriceService) Do(ctx context.Context) (*ConvertQuotedPrice, error) {
	req := request.Get(ctx, s.c, "/api/v2/convert/quoted-price", s.params).WithSign()
	return request.Do[ConvertQuotedPrice](req)
}

// ConvertQuotedPrice is the RFQ result for a Convert trade.
type ConvertQuotedPrice struct {
	Fee          decimal.Decimal `json:"fee"`
	FromCoinSize decimal.Decimal `json:"fromCoinSize"`
	FromCoin     string          `json:"fromCoin"`
	CnvtPrice    decimal.Decimal `json:"cnvtPrice"` // quote-coin price / target-coin price
	ToCoinSize   decimal.Decimal `json:"toCoinSize"`
	ToCoin       string          `json:"toCoin"`
	TraceID      string          `json:"traceId"` // RFQ id, valid for 8 seconds
}

// ConvertTradeService -- POST /api/v2/convert/trade (convert trade, STATE-CHANGING)
//
// Executes a Convert (flash) trade against a quote previously obtained from
// GetConvertQuotedPriceService. All fields are required and must match the quote;
// the traceId is only valid for 8 seconds after the RFQ.
type ConvertTradeService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewConvertTradeService(fromCoin string, fromCoinSize, cnvtPrice decimal.Decimal, toCoin string, toCoinSize decimal.Decimal, traceID string) *ConvertTradeService {
	return &ConvertTradeService{c: c, body: map[string]any{
		"fromCoin":     fromCoin,
		"fromCoinSize": fromCoinSize.String(),
		"cnvtPrice":    cnvtPrice.String(),
		"toCoin":       toCoin,
		"toCoinSize":   toCoinSize.String(),
		"traceId":      traceID,
	}}
}

func (s *ConvertTradeService) Do(ctx context.Context) (*ConvertTradeResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/convert/trade", s.body).WithSign()
	return request.Do[ConvertTradeResult](req)
}

// ConvertTradeResult is the confirmation of an executed Convert trade.
type ConvertTradeResult struct {
	Ts         time.Time       `json:"ts"`
	CnvtPrice  decimal.Decimal `json:"cnvtPrice"`
	ToCoinSize decimal.Decimal `json:"toCoinSize"`
	ToCoin     string          `json:"toCoin"`
}

// GetConvertRecordService -- GET /api/v2/convert/convert-record (convert read)
//
// Returns the account's Convert (flash) trade history. The startTime/endTime
// window is required and may span at most 90 days; results paginate via
// idLessThan using the previous page's smallest orderId.
type GetConvertRecordService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetConvertRecordService(startTime, endTime time.Time) *GetConvertRecordService {
	return &GetConvertRecordService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetConvertRecordService) SetLimit(limit int) *GetConvertRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest orderId from the previous page.
func (s *GetConvertRecordService) SetIDLessThan(idLessThan string) *GetConvertRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetConvertRecordService) Do(ctx context.Context) (*ConvertRecords, error) {
	req := request.Get(ctx, s.c, "/api/v2/convert/convert-record", s.params).WithSign()
	return request.Do[ConvertRecords](req)
}

// ConvertRecords is a page of Convert trade history.
type ConvertRecords struct {
	DataList []ConvertRecord `json:"dataList"`
	EndID    string          `json:"endId"` // pagination cursor
}

// ConvertRecord is one historical Convert (flash) trade.
type ConvertRecord struct {
	ID           string          `json:"id"`
	Ts           time.Time       `json:"ts"`
	CnvtPrice    decimal.Decimal `json:"cnvtPrice"`
	Fee          decimal.Decimal `json:"fee"`
	FromCoinSize decimal.Decimal `json:"fromCoinSize"`
	FromCoin     string          `json:"fromCoin"`
	ToCoinSize   decimal.Decimal `json:"toCoinSize"`
	ToCoin       string          `json:"toCoin"`
}

// GetBGBConvertCoinListService -- GET /api/v2/convert/bgb-convert-coin-list (convert read)
//
// Returns the small-balance coins eligible to be converted into BGB, with the
// estimated BGB redeemable, the BGB scale (precision) and the per-coin fee tiers.
type GetBGBConvertCoinListService struct {
	c *CommonClient
}

func (c *CommonClient) NewGetBGBConvertCoinListService() *GetBGBConvertCoinListService {
	return &GetBGBConvertCoinListService{c: c}
}

func (s *GetBGBConvertCoinListService) Do(ctx context.Context) (*BGBConvertCoinList, error) {
	req := request.Get(ctx, s.c, "/api/v2/convert/bgb-convert-coin-list").WithSign()
	return request.Do[BGBConvertCoinList](req)
}

// BGBConvertCoinList wraps the coins convertible into BGB.
type BGBConvertCoinList struct {
	CoinList []BGBConvertCoin `json:"coinList"`
}

// BGBConvertCoin is one coin that can be converted into BGB.
type BGBConvertCoin struct {
	Coin         string              `json:"coin"`
	Available    decimal.Decimal     `json:"available"`
	BgbEstAmount decimal.Decimal     `json:"bgbEstAmount"`
	Precision    string              `json:"precision"` // BGB scale
	FeeDetail    []BGBConvertCoinFee `json:"feeDetail"`
	CTime        time.Time           `json:"cTime"`
}

// BGBConvertCoinFee is one fee tier for converting a coin into BGB.
type BGBConvertCoinFee struct {
	FeeRate decimal.Decimal `json:"feeRate"`
	Fee     decimal.Decimal `json:"fee"`
}

// ConvertBGBService -- POST /api/v2/convert/bgb-convert (convert trade, STATE-CHANGING)
//
// Converts the given small-balance coins into BGB in one shot. The response lists
// one swap order per coin.
type ConvertBGBService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewConvertBGBService(coinList []string) *ConvertBGBService {
	return &ConvertBGBService{c: c, body: map[string]any{"coinList": coinList}}
}

// SetCoinList replaces the set of coins to convert into BGB.
func (s *ConvertBGBService) SetCoinList(coinList []string) *ConvertBGBService {
	s.body["coinList"] = coinList
	return s
}

func (s *ConvertBGBService) Do(ctx context.Context) (*ConvertBGBResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/convert/bgb-convert", s.body).WithSign()
	return request.Do[ConvertBGBResult](req)
}

// ConvertBGBResult wraps the per-coin swap orders created by a BGB conversion.
type ConvertBGBResult struct {
	OrderList []BGBConvertOrder `json:"orderList"`
}

// BGBConvertOrder is one coin's swap order from a BGB conversion.
type BGBConvertOrder struct {
	Coin    string `json:"coin"`
	OrderID string `json:"orderId"`
}

// GetBGBConvertRecordService -- GET /api/v2/convert/bgb-convert-records (convert read)
//
// Returns the BGB conversion history. All filters are optional; the
// startTime/endTime window may span at most 90 days and results paginate via
// idLessThan.
//
// NOTE: the live endpoint path is the plural "bgb-convert-records"; the singular
// "bgb-convert-record" returns 40404 Request URL NOT FOUND.
type GetBGBConvertRecordService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetBGBConvertRecordService() *GetBGBConvertRecordService {
	return &GetBGBConvertRecordService{c: c, params: map[string]string{}}
}

// SetOrderID filters to a single BGB conversion order.
func (s *GetBGBConvertRecordService) SetOrderID(orderID string) *GetBGBConvertRecordService {
	s.params["orderId"] = orderID
	return s
}

// SetStartTime filters records at or after t (max 90-day span).
func (s *GetBGBConvertRecordService) SetStartTime(t time.Time) *GetBGBConvertRecordService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 90-day span).
func (s *GetBGBConvertRecordService) SetEndTime(t time.Time) *GetBGBConvertRecordService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 20, max 100).
func (s *GetBGBConvertRecordService) SetLimit(limit int) *GetBGBConvertRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan paginates: set it to the smallest orderId from the previous page.
func (s *GetBGBConvertRecordService) SetIDLessThan(idLessThan string) *GetBGBConvertRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetBGBConvertRecordService) Do(ctx context.Context) ([]BGBConvertRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/convert/bgb-convert-records", s.params).WithSign()
	resp, err := request.Do[[]BGBConvertRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BGBConvertRecord is one historical BGB conversion.
type BGBConvertRecord struct {
	OrderID       string                `json:"orderId"`
	FromCoin      string                `json:"fromCoin"`
	FromAmount    decimal.Decimal       `json:"fromAmount"`
	FromCoinPrice decimal.Decimal       `json:"fromCoinPrice"`
	ToCoin        string                `json:"toCoin"`
	ToAmount      decimal.Decimal       `json:"toAmount"`
	ToCoinPrice   decimal.Decimal       `json:"toCoinPrice"`
	FeeDetail     []BGBConvertRecordFee `json:"feeDetail"`
	Status        string                `json:"status"`
	Ctime         time.Time             `json:"ctime"`
}

// BGBConvertRecordFee is one fee line of a BGB conversion record.
type BGBConvertRecordFee struct {
	FeeCoin string          `json:"feeCoin"`
	Fee     decimal.Decimal `json:"fee"`
}
