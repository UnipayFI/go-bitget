package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// The Reality basic-data endpoints describe the company behind an RWA (Reality)
// trading pair rather than the pair itself, so they are keyed by the stock's
// ticker ("AAPL") — the code field of GetStockInfoService — not by symbol.
// Dates arrive as millisecond timestamps, with the exception noted on
// CompanyOverview.ListingDate.

// GetCompanyOverviewService -- GET /api/v3/reality/market/company-overview
//
// Returns the profile and headline valuation of the company behind a Reality
// trading pair.
type GetCompanyOverviewService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetCompanyOverviewService(code string) *GetCompanyOverviewService {
	return &GetCompanyOverviewService{c: c, params: map[string]string{"code": code}}
}

func (s *GetCompanyOverviewService) Do(ctx context.Context) (*CompanyOverview, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/company-overview", s.params)
	return request.Do[CompanyOverview](req)
}

// CompanyOverview is the profile of a listed company.
type CompanyOverview struct {
	Code       string          `json:"code"`
	Name       string          `json:"name"`
	PERatio    decimal.Decimal `json:"peRatio"`
	PBRatio    decimal.Decimal `json:"pbRatio"`
	MarketCap  decimal.Decimal `json:"marketCap"`
	High52Week decimal.Decimal `json:"high52Week"`
	Low52Week  decimal.Decimal `json:"low52Week"`
	// TotalShares is empty for some listings; use GetShareCapitalChangeService
	// for the authoritative share count.
	TotalShares decimal.Decimal `json:"totalShares"`
	// ListingDate is a plain yyyy-MM-dd date (e.g. "1980-12-12"), not the
	// millisecond timestamp the other Reality date fields carry. The docs spell
	// the field "ListingDate"; the live API sends "listingDate".
	ListingDate    string          `json:"listingDate"`
	Employees      decimal.Decimal `json:"employees"`
	CompanyAddress string          `json:"companyAddress"`
}

// GetValuationIndicatorsService -- GET /api/v3/reality/market/valuation-indicators
//
// Returns the full valuation-ratio set of a listed company. TTM ratios use the
// trailing twelve months, LYR ratios the last full-year report; Ed and Pd are
// the period-end and period-average variants.
type GetValuationIndicatorsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetValuationIndicatorsService(code string) *GetValuationIndicatorsService {
	return &GetValuationIndicatorsService{c: c, params: map[string]string{"code": code}}
}

func (s *GetValuationIndicatorsService) Do(ctx context.Context) (*ValuationIndicators, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/valuation-indicators", s.params)
	return request.Do[ValuationIndicators](req)
}

// ValuationIndicators is the valuation-ratio snapshot of a listed company.
type ValuationIndicators struct {
	Date time.Time `json:"date"`
	// PB, PE, PS and PCF are the price-to-book, price-to-earnings,
	// price-to-sales and price-to-cash-flow ratios.
	PB          decimal.Decimal `json:"pb"`
	PBEd        decimal.Decimal `json:"pbEd"`
	PBPd        decimal.Decimal `json:"pbPd"`
	PBMrq       decimal.Decimal `json:"pbMrq"` // most recent quarter
	PE          decimal.Decimal `json:"pe"`
	PELyr       decimal.Decimal `json:"peLyr"`
	PETtmEd     decimal.Decimal `json:"peTtmEd"`
	PETtmPd     decimal.Decimal `json:"peTtmPd"`
	PS          decimal.Decimal `json:"ps"`
	PSLyr       decimal.Decimal `json:"psLyr"`
	PSTtmEd     decimal.Decimal `json:"psTtmEd"`
	PSTtmPd     decimal.Decimal `json:"psTtmPd"`
	PCF         decimal.Decimal `json:"pcf"`
	PCFNet      decimal.Decimal `json:"pcfNet"`
	PCFNetTtmEd decimal.Decimal `json:"pcfNetTtmEd"`
	PCFNetTtmPd decimal.Decimal `json:"pcfNetTtmPd"`
	PCFTtmEd    decimal.Decimal `json:"pcfTtmEd"`
	PCFTtmPd    decimal.Decimal `json:"pcfTtmPd"`
	// EVEbitda is the enterprise-value to EBITDA multiple.
	EVEbitda decimal.Decimal `json:"evEbitda"`
	// EquityValue is null for some listings and then decodes as zero.
	EquityValue      decimal.Decimal `json:"equityValue"`
	GrossEV          decimal.Decimal `json:"grossEv"`
	NetEV            decimal.Decimal `json:"netEv"`
	TMVUSD           decimal.Decimal `json:"tmvUsd"`    // total market value
	CirTMVUSD        decimal.Decimal `json:"cirTmvUsd"` // circulating market value
	DividendYieldTtm decimal.Decimal `json:"dividendYieldTtm"`
}

// GetEarningsForecastService -- GET /api/v3/reality/market/earnings-forecast
//
// Returns the analyst earnings forecast (or the reported actuals, see IsActual)
// of a listed company for one fiscal year.
type GetEarningsForecastService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetEarningsForecastService(code string) *GetEarningsForecastService {
	return &GetEarningsForecastService{c: c, params: map[string]string{"code": code}}
}

func (s *GetEarningsForecastService) Do(ctx context.Context) (*EarningsForecast, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/earnings-forecast", s.params)
	return request.Do[EarningsForecast](req)
}

// EarningsForecast is one fiscal year of forecast (or actual) company earnings.
type EarningsForecast struct {
	FiscalYear          string    `json:"fiscalYear"`
	PublicationDeadline time.Time `json:"publicationDeadline"`
	// IsActual distinguishes a reported figure from a forecast; it is a real
	// JSON boolean, not one of Bitget's usual quoted flags.
	IsActual bool `json:"isActual"`
	// Revenue, EBIT and NetIncomeParent are in millions of the reporting
	// Currency; NetIncomeParent is attributable to the parent's shareholders.
	Revenue         decimal.Decimal `json:"revenue"`
	EBIT            decimal.Decimal `json:"ebit"`
	NetIncomeParent decimal.Decimal `json:"netIncomeParent"`
	// EPS, BPS and CFPS are the per-share earnings, book value and cash flow.
	EPS      decimal.Decimal `json:"eps"`
	BPS      decimal.Decimal `json:"bps"`
	CFPS     decimal.Decimal `json:"cfps"`
	ROE      decimal.Decimal `json:"roe"`
	ROA      decimal.Decimal `json:"roa"`
	Currency string          `json:"currency"`
}

// GetSuspensionResumptionInfoService -- GET /api/v3/reality/market/suspension-resumption-info
//
// Returns the latest trading halt and resumption of a listed company.
type GetSuspensionResumptionInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSuspensionResumptionInfoService(code string) *GetSuspensionResumptionInfoService {
	return &GetSuspensionResumptionInfoService{c: c, params: map[string]string{"code": code}}
}

func (s *GetSuspensionResumptionInfoService) Do(ctx context.Context) (*SuspensionResumptionInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/suspension-resumption-info", s.params)
	return request.Do[SuspensionResumptionInfo](req)
}

// SuspensionResumptionInfo is a trading halt and its resumption. The dates are
// millisecond timestamps while the times are local HH:mm:ss clock strings.
type SuspensionResumptionInfo struct {
	Code             string          `json:"code"`
	Name             string          `json:"name"` // localized by the locale header
	SuspensionDate   time.Time       `json:"suspensionDate"`
	SuspensionTime   string          `json:"suspensionTime"`
	SuspensionReason string          `json:"suspensionReason"`
	SuspensionPrice  decimal.Decimal `json:"suspensionPrice"`
	ResumptionDate   time.Time       `json:"resumptionDate"`
	// ResumptionQuoteTime is when quoting reopens, ResumptionTradingTime when
	// the security is released for trading.
	ResumptionQuoteTime   string `json:"resumptionQuoteTime"`
	ResumptionTradingTime string `json:"resumptionTradingTime"`
}

// GetDividendsService -- GET /api/v3/reality/market/dividends
//
// Returns the dividend and split history of a listed company, newest first and
// paginated by cursor.
type GetDividendsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetDividendsService(code string) *GetDividendsService {
	return &GetDividendsService{c: c, params: map[string]string{"code": code}}
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetDividendsService) SetLimit(limit string) *GetDividendsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetDividendsService) SetCursor(cursor string) *GetDividendsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetDividendsService) Do(ctx context.Context) (*Dividends, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/dividends", s.params)
	return request.Do[Dividends](req)
}

type Dividends struct {
	List   []DividendRecord `json:"list"`
	Cursor string           `json:"cursor"`
}

// DividendRecord is one dividend, stock dividend or split. Only the fields of
// the matching Type are populated; the rest arrive as null and decode as zero.
type DividendRecord struct {
	Type             string    `json:"type"` // cash_dividend, stock_split, stock_dividend
	AnnouncementDate time.Time `json:"announcementDate"`
	RecordDate       time.Time `json:"recordDate"`
	ExrightDate      time.Time `json:"exrightDate"`
	DividendDate     time.Time `json:"dividendDate"`
	// DividendPerShare is the cash paid per share, StockDividendPerShare the
	// shares granted per share.
	DividendPerShare      decimal.Decimal `json:"dividendPerShare"`
	StockDividendPerShare decimal.Decimal `json:"stockDividendPerShare"`
	SplitValidDate        time.Time       `json:"splitValidDate"`
	// SplitNumerator over SplitDenominator is the split ratio: 10 over 1 is a
	// ten-for-one split.
	SplitNumerator   decimal.Decimal `json:"splitNumerator"`
	SplitDenominator decimal.Decimal `json:"splitDenominator"`
}

// GetShareCapitalChangeService -- GET /api/v3/reality/market/share-capital-change
//
// Returns the latest share capital structure of a listed company and why it
// last changed.
type GetShareCapitalChangeService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetShareCapitalChangeService(code string) *GetShareCapitalChangeService {
	return &GetShareCapitalChangeService{c: c, params: map[string]string{"code": code}}
}

func (s *GetShareCapitalChangeService) Do(ctx context.Context) (*ShareCapitalChange, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/share-capital-change", s.params)
	return request.Do[ShareCapitalChange](req)
}

// ShareCapitalChange is a company's share structure after its latest change.
type ShareCapitalChange struct {
	AnnouncementDate time.Time       `json:"announcementDate"`
	ChangeDate       time.Time       `json:"changeDate"`
	TotalShares      decimal.Decimal `json:"totalShares"`
	CommonShares     decimal.Decimal `json:"commonShares"`
	PreferredShares  decimal.Decimal `json:"preferredShares"`
	OtherShares      decimal.Decimal `json:"otherShares"`
	SpecialExplain   string          `json:"specialExplain"`
	ChangeReason     string          `json:"changeReason"`
}

// GetInnerTradesService -- GET /api/v3/reality/market/inner-trades
//
// Returns the insider trades reported for a listed company, newest first and
// paginated by cursor.
type GetInnerTradesService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetInnerTradesService(code string) *GetInnerTradesService {
	return &GetInnerTradesService{c: c, params: map[string]string{"code": code}}
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetInnerTradesService) SetLimit(limit string) *GetInnerTradesService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetInnerTradesService) SetCursor(cursor string) *GetInnerTradesService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetInnerTradesService) Do(ctx context.Context) (*InnerTrades, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/inner-trades", s.params)
	return request.Do[InnerTrades](req)
}

type InnerTrades struct {
	List   []InnerTrade `json:"list"`
	Cursor string       `json:"cursor"`
}

// InnerTrade is one reported insider transaction.
type InnerTrade struct {
	ReporterName     string    `json:"reporterName"`
	IssueOrgName     string    `json:"issueOrgName"`
	EndDate          time.Time `json:"endDate"`
	AnnouncementDate time.Time `json:"announcementDate"`
	// InnerType is the insider's relationship to the issuer (Director,
	// Officer, …) and Position the specific role.
	InnerType              string    `json:"innerType"`
	Position               string    `json:"position"`
	SecurityTitle          string    `json:"securityTitle"`
	TradeDate              time.Time `json:"tradeDate"`
	DesignatedExerciseDate time.Time `json:"designatedExerciseDate"`
	// InnerTradeType is the SEC transaction code: P purchase, S sale, A award,
	// F tax withholding, M exempt, C conversion, G gift, J other.
	InnerTradeType string          `json:"innerTradeType"`
	ActionType     string          `json:"actionType"` // add, decrease
	TradeVolume    decimal.Decimal `json:"tradeVolume"`
	TradePrice     decimal.Decimal `json:"tradePrice"`
	HoldingType    string          `json:"holdingType"` // direct, indirect
	// IndirectHoldingStatement explains an indirect holding.
	IndirectHoldingStatement string          `json:"indirectHoldingStatement"`
	PostTradeQuantity        decimal.Decimal `json:"postTradeQuantity"`
}

// GetExecutiveShareholdingsService -- GET /api/v3/reality/market/executive-shareholdings
//
// Returns the shares held by a listed company's directors and officers,
// paginated by cursor.
type GetExecutiveShareholdingsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetExecutiveShareholdingsService(code string) *GetExecutiveShareholdingsService {
	return &GetExecutiveShareholdingsService{c: c, params: map[string]string{"code": code}}
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetExecutiveShareholdingsService) SetLimit(limit string) *GetExecutiveShareholdingsService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetExecutiveShareholdingsService) SetCursor(cursor string) *GetExecutiveShareholdingsService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetExecutiveShareholdingsService) Do(ctx context.Context) (*ExecutiveShareholdings, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/executive-shareholdings", s.params)
	return request.Do[ExecutiveShareholdings](req)
}

type ExecutiveShareholdings struct {
	List   []ExecutiveShareholding `json:"list"`
	Cursor string                  `json:"cursor"`
}

// ExecutiveShareholding is one executive's stake as of a reporting period.
type ExecutiveShareholding struct {
	Name       string          `json:"name"`
	HoldingNum decimal.Decimal `json:"holdingNum"`
	// HoldingRatio and VoteRatio are percentages: 9.63 means 9.63%.
	HoldingRatio     decimal.Decimal `json:"holdingRatio"`
	VoteRatio        decimal.Decimal `json:"voteRatio"`
	Period           string          `json:"period"` // localized by the locale header
	AnnouncementDate time.Time       `json:"announcementDate"`
}

// GetShareholdDetailService -- GET /api/v3/reality/market/sharehold-detail
//
// Returns the major shareholders of a listed company, paginated by cursor.
type GetShareholdDetailService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetShareholdDetailService(code string) *GetShareholdDetailService {
	return &GetShareholdDetailService{c: c, params: map[string]string{"code": code}}
}

// SetLimit sets the page size (default 20, max 100).
func (s *GetShareholdDetailService) SetLimit(limit string) *GetShareholdDetailService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor from a previous response.
func (s *GetShareholdDetailService) SetCursor(cursor string) *GetShareholdDetailService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetShareholdDetailService) Do(ctx context.Context) (*ShareholdDetail, error) {
	req := request.Get(ctx, s.c, "/api/v3/reality/market/sharehold-detail", s.params)
	return request.Do[ShareholdDetail](req)
}

type ShareholdDetail struct {
	List   []ShareholderHolding `json:"list"`
	Cursor string               `json:"cursor"`
}

// ShareholderHolding is one major shareholder's stake. The same holder can
// appear several times, once per reporting period.
type ShareholderHolding struct {
	Name       string          `json:"name"`
	HoldingNum decimal.Decimal `json:"holdingNum"`
	// HoldingRatio is a percentage: 9.63 means 9.63%.
	HoldingRatio decimal.Decimal `json:"holdingRatio"`
}
