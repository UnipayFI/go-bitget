package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSmallAssetsService -- GET /api/v3/convert/small-assets (UTA mgt. read)
//
// Returns the account's dust balances that are eligible for one-click
// conversion, each quoted with the target coin (BGB), the estimated proceeds
// and the conversion fee.
type GetSmallAssetsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSmallAssetsService() *GetSmallAssetsService {
	return &GetSmallAssetsService{c: c}
}

func (s *GetSmallAssetsService) Do(ctx context.Context) ([]SmallAsset, error) {
	req := request.Get(ctx, s.c, "/api/v3/convert/small-assets").WithSign()
	resp, err := request.Do[[]SmallAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SmallAsset struct {
	Coin            string          `json:"coin"`
	Available       decimal.Decimal `json:"available"`
	EstimatedCoin   string          `json:"estimatedCoin"` // conversion target, BGB
	EstimatedAmount decimal.Decimal `json:"estimatedAmount"`
	FeeDetail       SmallAssetFee   `json:"feeDetail"`
	CTime           time.Time       `json:"cTime"`
}

// SmallAssetFee is the fee quoted for converting one small asset.
type SmallAssetFee struct {
	FeeRate decimal.Decimal `json:"feeRate"`
	Fee     decimal.Decimal `json:"fee"`
}

// SmallAssetsTradeService -- POST /api/v3/convert/small-assets-trade (UTA trade)
//
// Converts the given dust coins into the target coin (BGB) in one order.
type SmallAssetsTradeService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSmallAssetsTradeService(fromCoinList []string) *SmallAssetsTradeService {
	return &SmallAssetsTradeService{c: c, body: map[string]any{"fromCoinList": fromCoinList}}
}

// SetFromCoinList sets the coins to convert.
func (s *SmallAssetsTradeService) SetFromCoinList(fromCoinList []string) *SmallAssetsTradeService {
	s.body["fromCoinList"] = fromCoinList
	return s
}

func (s *SmallAssetsTradeService) Do(ctx context.Context) (*SmallAssetsTradeResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/convert/small-assets-trade", s.body).WithSign()
	return request.Do[SmallAssetsTradeResult](req)
}

type SmallAssetsTradeResult struct {
	FromCoin []string `json:"fromCoin"`
	ToCoin   string   `json:"toCoin"`
	OrderID  string   `json:"orderId"`
}

// GetSmallAssetsHistoryService -- GET /api/v3/convert/small-assets-history (UTA mgt. read)
//
// Returns past small-asset conversions, paginated by cursor and bounded to a
// 90-day access window (max 30 days per query).
type GetSmallAssetsHistoryService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSmallAssetsHistoryService() *GetSmallAssetsHistoryService {
	return &GetSmallAssetsHistoryService{c: c, params: map[string]string{}}
}

// SetOrderID filters history by the conversion order ID.
func (s *GetSmallAssetsHistoryService) SetOrderID(orderId string) *GetSmallAssetsHistoryService {
	s.params["orderId"] = orderId
	return s
}

// SetStartTime filters records at or after t (90-day access window).
func (s *GetSmallAssetsHistoryService) SetStartTime(t time.Time) *GetSmallAssetsHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetSmallAssetsHistoryService) SetEndTime(t time.Time) *GetSmallAssetsHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetSmallAssetsHistoryService) SetLimit(limit int) *GetSmallAssetsHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetCursor pages forward using the smallest orderId of the previous page.
func (s *GetSmallAssetsHistoryService) SetCursor(cursor string) *GetSmallAssetsHistoryService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSmallAssetsHistoryService) Do(ctx context.Context) (*SmallAssetsHistory, error) {
	req := request.Get(ctx, s.c, "/api/v3/convert/small-assets-history", s.params).WithSign()
	return request.Do[SmallAssetsHistory](req)
}

type SmallAssetsHistory struct {
	List   []SmallAssetConversion `json:"list"`
	Cursor string                 `json:"cursor"`
}

type SmallAssetConversion struct {
	OrderID       string          `json:"orderId"`
	FromCoin      string          `json:"fromCoin"`
	FromAmount    decimal.Decimal `json:"fromAmount"`
	FromCoinPrice decimal.Decimal `json:"fromCoinPrice"`
	ToCoin        string          `json:"toCoin"`
	ToAmount      decimal.Decimal `json:"toAmount"`
	ToCoinPrice   decimal.Decimal `json:"toCoinPrice"`
	FeeDetail     FeeDetail       `json:"feeDetail"`
	Status        string          `json:"status"` // success
	CreatedTime   time.Time       `json:"createdTime"`
}
