package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetServerTimeService -- GET /api/v3/market/time
//
// Returns the Bitget server time; SyncServerTime uses it to correct clock drift
// before signing private requests.
type GetServerTimeService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetServerTimeService() *GetServerTimeService {
	return &GetServerTimeService{c: c}
}

func (s *GetServerTimeService) Do(ctx context.Context) (*ServerTimeResponse, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/time")
	return request.Do[ServerTimeResponse](req)
}

// ServerTimeResponse carries the same instant at millisecond and microsecond
// precision.
type ServerTimeResponse struct {
	ServerTime       time.Time `json:"serverTime,format:unixmilli"`
	ServerTimeMicros time.Time `json:"serverTimeMicros,format:unixmicro"`
}

// GetCoinsService -- GET /api/v3/market/coins
//
// Returns coin metadata and the supported deposit/withdraw chains, for every
// coin unless SetCoin is used.
type GetCoinsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCoinsService() *GetCoinsService {
	return &GetCoinsService{c: c, params: map[string]string{}}
}

// SetCoin filters the result to a single coin (e.g. "BTC").
func (s *GetCoinsService) SetCoin(coin string) *GetCoinsService {
	s.params["coin"] = coin
	return s
}

func (s *GetCoinsService) Do(ctx context.Context) ([]CoinInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/market/coins", s.params)
	resp, err := request.Do[[]CoinInfo](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// CoinInfo is a single coin's metadata together with its per-chain config.
type CoinInfo struct {
	CoinID   string      `json:"coinId"`
	Coin     string      `json:"coin"`
	Transfer string      `json:"transfer"` // "true"/"false": transferable between accounts
	AreaCoin string      `json:"areaCoin"` // "yes"/"no": region-restricted coin
	Chains   []CoinChain `json:"chains"`
}

// CoinChain is the deposit/withdraw configuration for one chain of a coin.
type CoinChain struct {
	Chain             string          `json:"chain"`
	NeedTag           string          `json:"needTag"`      // "true"/"false": deposit requires a memo/tag
	Withdrawable      string          `json:"withdrawable"` // "true"/"false"
	Rechargeable      string          `json:"rechargeable"` // "true"/"false": depositable
	WithdrawFee       decimal.Decimal `json:"withdrawFee"`
	ExtraWithdrawFee  decimal.Decimal `json:"extraWithdrawFee"` // extra on-chain destruction fee, 0.1 = 10%
	DepositConfirm    string          `json:"depositConfirm"`   // confirmations to credit a deposit
	WithdrawConfirm   string          `json:"withdrawConfirm"`  // confirmations before withdraw is final
	MinDepositAmount  decimal.Decimal `json:"minDepositAmount"`
	MinWithdrawAmount decimal.Decimal `json:"minWithdrawAmount"`
	BrowserURL        string          `json:"browserUrl"`
	ContractAddress   string          `json:"contractAddress"`
	WithdrawStep      string          `json:"withdrawStep"` // withdraw amount step (0 = no step limit)
	// WithdrawMinScale is the decimal places allowed on a withdraw amount (the
	// docs call it withdrawPrecision; the API sends withdrawMinScale).
	WithdrawMinScale string `json:"withdrawMinScale"`
	Congestion       string `json:"congestion"` // chain status: "normal"/"congested"
}

// AnnouncementType enumerates the announcement categories accepted by the
// annType filter (and echoed back in each announcement's annType field).
type AnnouncementType string

const (
	AnnouncementTypeLatestNews               AnnouncementType = "latest_news"
	AnnouncementTypeCoinListings             AnnouncementType = "coin_listings"
	AnnouncementTypeProductUpdates           AnnouncementType = "product_updates"
	AnnouncementTypeSecurity                 AnnouncementType = "security"
	AnnouncementTypeAPITrading               AnnouncementType = "api_trading"
	AnnouncementTypeMaintenanceSystemUpdates AnnouncementType = "maintenance_system_updates"
	AnnouncementTypeSymbolDelisting          AnnouncementType = "symbol_delisting"
)

// AnnouncementLanguage enumerates the localization codes accepted by the
// required language parameter; unsupported languages fall back to English.
type AnnouncementLanguage string

const (
	AnnouncementLanguageZhCN AnnouncementLanguage = "zh_CN"
	AnnouncementLanguageEnUS AnnouncementLanguage = "en_US"
)

// GetAnnouncementsService -- GET /api/v3/public/annoucements
//
// Returns Bitget platform announcements for a language, newest first,
// optionally filtered by category and display-time window and paginated by
// cursor (max 10 per page). Note the path uses Bitget's misspelling
// "annoucements".
type GetAnnouncementsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAnnouncementsService(language AnnouncementLanguage) *GetAnnouncementsService {
	return &GetAnnouncementsService{c: c, params: map[string]string{"language": string(language)}}
}

func (s *GetAnnouncementsService) SetAnnType(annType AnnouncementType) *GetAnnouncementsService {
	s.params["annType"] = string(annType)
	return s
}

// SetStartTime filters announcements displayed at or after t.
func (s *GetAnnouncementsService) SetStartTime(t time.Time) *GetAnnouncementsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters announcements displayed at or before t.
func (s *GetAnnouncementsService) SetEndTime(t time.Time) *GetAnnouncementsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetCursor pages on: pass the last annId of the previous page.
func (s *GetAnnouncementsService) SetCursor(cursor string) *GetAnnouncementsService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit sets the page size (default 10, max 10).
func (s *GetAnnouncementsService) SetLimit(limit int) *GetAnnouncementsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAnnouncementsService) Do(ctx context.Context) ([]Announcement, error) {
	req := request.Get(ctx, s.c, "/api/v3/public/annoucements", s.params)
	resp, err := request.Do[[]Announcement](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// Announcement is a single platform notice. The docs list its release time as
// createdTime; the API sends cTime.
type Announcement struct {
	AnnID      string           `json:"annId"`
	AnnTitle   string           `json:"annTitle"`
	AnnDesc    string           `json:"annDesc"`
	AnnType    AnnouncementType `json:"annType"`
	AnnSubType string           `json:"annSubType"`
	AnnURL     string           `json:"annUrl"`
	Language   string           `json:"language"`
	CTime      time.Time        `json:"cTime,format:unixmilli"`
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
	AreaSymbol string   `json:"areaSymbol"` // spot/margin region-restricted flag (yes/no); spot sends only yes

	// Traditional-finance identity, returned for both spot and futures.
	// FullName is the pair/company/ETF name (may carry an ADR/ST/QUANTO mark);
	// MarketRegion (US, CN, HK, KR, JP) and ISIN are empty for crypto pairs.
	FullName       string `json:"fullName"`
	InstrumentType string `json:"instrumentType"` // crypto, metal, energy, stock, etf, fx, index
	MarketRegion   string `json:"marketRegion"`
	ISIN           string `json:"isin"`

	BuyLimitPriceRatio   decimal.Decimal  `json:"buyLimitPriceRatio"`
	SellLimitPriceRatio  decimal.Decimal  `json:"sellLimitPriceRatio"`
	MinOrderQty          decimal.Decimal  `json:"minOrderQty"`
	MaxOrderQty          decimal.Decimal  `json:"maxOrderQty"`
	MinOrderAmount       decimal.Decimal  `json:"minOrderAmount"`       // in USDT
	MaxMarketOrderAmount decimal.Decimal  `json:"maxMarketOrderAmount"` // spot only -- max single market order amount in quote coin
	PricePrecision       decimal.Decimal  `json:"pricePrecision"`
	QuantityPrecision    decimal.Decimal  `json:"quantityPrecision"`
	QuotePrecision       decimal.Decimal  `json:"quotePrecision"`
	MaxSymbolOrderNum    string           `json:"maxSymbolOrderNum"`
	MaxProductOrderNum   string           `json:"maxProductOrderNum"`
	MaxPositionNum       string           `json:"maxPositionNum"`
	Status               InstrumentStatus `json:"status"`
	MaintainTime         string           `json:"maintainTime"`
	LaunchTime           time.Time        `json:"launchTime,format:unixmilli"`

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
	OffTime            time.Time       `json:"offTime,format:unixmilli"`
	LimitOpenTime      time.Time       `json:"limitOpenTime,format:unixmilli"`
	DeliveryTime       time.Time       `json:"deliveryTime,format:unixmilli"`
	DeliveryStartTime  time.Time       `json:"deliveryStartTime,format:unixmilli"`
	DeliveryPeriod     string          `json:"deliveryPeriod"`

	// Margin-only fields.
	IsIsolatedBaseBorrowable   string          `json:"isIsolatedBaseBorrowable"`   // YES/NO
	IsIsolatedQuotedBorrowable string          `json:"isIsolatedQuotedBorrowable"` // YES/NO
	WarningRiskRatio           decimal.Decimal `json:"warningRiskRatio"`
	LiquidationRiskRatio       decimal.Decimal `json:"liquidationRiskRatio"`
	MaxCrossedLeverage         decimal.Decimal `json:"maxCrossedLeverage"`
	MaxIsolatedLeverage        decimal.Decimal `json:"maxIsolatedLeverage"`
	UserMinBorrow              decimal.Decimal `json:"userMinBorrow"`
}
