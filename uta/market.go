package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetServerTimeService -- GET /api/v2/public/time
//
// Returns the Bitget server time. Bitget exposes a single shared time endpoint
// (under the v2 common namespace); SyncServerTime uses it to correct clock
// drift before signing private requests.
type GetServerTimeService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v2/public/time")
	return request.Do[ServerTimeResponse](req)
}

type ServerTimeResponse struct {
	ServerTime time.Time `json:"serverTime"`
}

// GetInstrumentsService -- GET /api/v3/market/instruments
//
// Returns trading specifications for a product category, optionally filtered to
// a single symbol.
type GetInstrumentsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInstrumentsService(category Category) *GetInstrumentsService {
	return &GetInstrumentsService{c: c, params: map[string]string{"category": string(category)}}
}

func (s *GetInstrumentsService) SetSymbol(symbol string) *GetInstrumentsService {
	s.params["symbol"] = symbol
	return s
}

func (s *GetInstrumentsService) Do(ctx context.Context) ([]Instrument, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/instruments", s.params)
	resp, err := request.Do[[]Instrument](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Instrument is the union of the spot/margin and futures instrument shapes; the
// futures-only fields are empty for SPOT/MARGIN categories and vice versa.
type Instrument struct {
	Symbol     string   `json:"symbol"`
	Category   Category `json:"category"`
	BaseCoin   string   `json:"baseCoin"`
	QuoteCoin  string   `json:"quoteCoin"`
	SymbolType string   `json:"symbolType"` // crypto, metal, stock, commodity
	IsReality  string   `json:"isReality"`  // spot RWA flag (no/yes)
	IsRwa      string   `json:"isRwa"`      // futures RWA flag (NO/YES)

	BuyLimitPriceRatio  decimal.Decimal  `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio decimal.Decimal  `json:"sellLimitPriceRatio"`
	MinOrderQty         decimal.Decimal  `json:"minOrderQty"`
	MaxOrderQty         decimal.Decimal  `json:"maxOrderQty"`
	MinOrderAmount      decimal.Decimal  `json:"minOrderAmount"`
	PricePrecision      decimal.Decimal  `json:"pricePrecision"`
	QuantityPrecision   decimal.Decimal  `json:"quantityPrecision"`
	QuotePrecision      decimal.Decimal  `json:"quotePrecision"`
	MaxSymbolOrderNum   string           `json:"maxSymbolOrderNum"`
	MaxProductOrderNum  string           `json:"maxProductOrderNum"`
	MaxPositionNum      string           `json:"maxPositionNum"`
	Status              InstrumentStatus `json:"status"`
	MaintainTime        string           `json:"maintainTime"`
	LaunchTime          time.Time        `json:"launchTime"`

	// Futures-only fields.
	Type               SymbolType      `json:"type"` // perpetual, delivery
	FeeRateUpRatio     decimal.Decimal `json:"feeRateUpRatio"`
	MakerFeeRate       decimal.Decimal `json:"makerFeeRate"`
	TakerFeeRate       decimal.Decimal `json:"takerFeeRate"`
	OpenCostUpRatio    decimal.Decimal `json:"openCostUpRatio"`
	PriceMultiplier    decimal.Decimal `json:"priceMultiplier"`
	QuantityMultiplier decimal.Decimal `json:"quantityMultiplier"`
	MaxMarketOrderQty  decimal.Decimal `json:"maxMarketOrderQty"`
	FundInterval       string          `json:"fundInterval"`
	MinLeverage        decimal.Decimal `json:"minLeverage"`
	MaxLeverage        decimal.Decimal `json:"maxLeverage"`
	OffTime            time.Time       `json:"offTime"`
	LimitOpenTime      time.Time       `json:"limitOpenTime"`
	DeliveryTime       time.Time       `json:"deliveryTime"`
	DeliveryStartTime  time.Time       `json:"deliveryStartTime"`
	DeliveryPeriod     string          `json:"deliveryPeriod"`
}
