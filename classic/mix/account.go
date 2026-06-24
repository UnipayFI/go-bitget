package mix

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetSingleAccountService -- GET /api/v2/mix/account/account (private)
//
// Returns a single futures account's equity, margin and leverage detail for one
// symbol + margin coin under a product type.
type GetSingleAccountService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetSingleAccountService(symbol string, productType ProductType, marginCoin string) *GetSingleAccountService {
	return &GetSingleAccountService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
	}}
}

func (s *GetSingleAccountService) Do(ctx context.Context) (*Account, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/account", s.params).WithSign()
	return request.Do[Account](req)
}

// Account is a single futures account's balance, margin and leverage snapshot.
type Account struct {
	MarginCoin            string          `json:"marginCoin"`
	Locked                decimal.Decimal `json:"locked"`
	Available             decimal.Decimal `json:"available"`
	CrossedMaxAvailable   decimal.Decimal `json:"crossedMaxAvailable"`
	IsolatedMaxAvailable  decimal.Decimal `json:"isolatedMaxAvailable"`
	MaxTransferOut        decimal.Decimal `json:"maxTransferOut"`
	AccountEquity         decimal.Decimal `json:"accountEquity"`
	USDTEquity            decimal.Decimal `json:"usdtEquity"`
	BtcEquity             decimal.Decimal `json:"btcEquity"`
	CrossedRiskRate       decimal.Decimal `json:"crossedRiskRate"`
	CrossedMarginLeverage decimal.Decimal `json:"crossedMarginLeverage"`
	IsolatedLongLever     decimal.Decimal `json:"isolatedLongLever"`
	IsolatedShortLever    decimal.Decimal `json:"isolatedShortLever"`
	MarginMode            MarginMode      `json:"marginMode"`
	PosMode               PositionMode    `json:"posMode"`
	UnrealizedPL          decimal.Decimal `json:"unrealizedPL"`
	Coupon                decimal.Decimal `json:"coupon"`
	CrossedUnrealizedPL   decimal.Decimal `json:"crossedUnrealizedPL"`
	IsolatedUnrealizedPL  decimal.Decimal `json:"isolatedUnrealizedPL"`
	AssetMode             AssetMode       `json:"assetMode"`
	Grant                 decimal.Decimal `json:"grant"`
	IsolatedMargin        decimal.Decimal `json:"isolatedMargin"`
	CrossedMargin         decimal.Decimal `json:"crossedMargin"`
}

// GetAccountListService -- GET /api/v2/mix/account/accounts (private)
//
// Returns every futures account (one per margin coin) under a product type.
type GetAccountListService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetAccountListService(productType ProductType) *GetAccountListService {
	return &GetAccountListService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetAccountListService) Do(ctx context.Context) ([]AccountListItem, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/accounts", s.params).WithSign()
	resp, err := request.Do[[]AccountListItem](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountListItem is one margin-coin account in the account list. It carries the
// union (multi-asset) aggregates in addition to the per-coin balances.
type AccountListItem struct {
	MarginCoin           string             `json:"marginCoin"`
	Locked               decimal.Decimal    `json:"locked"`
	Available            decimal.Decimal    `json:"available"`
	CrossedMaxAvailable  decimal.Decimal    `json:"crossedMaxAvailable"`
	IsolatedMaxAvailable decimal.Decimal    `json:"isolatedMaxAvailable"`
	MaxTransferOut       decimal.Decimal    `json:"maxTransferOut"`
	AccountEquity        decimal.Decimal    `json:"accountEquity"`
	USDTEquity           decimal.Decimal    `json:"usdtEquity"`
	BtcEquity            decimal.Decimal    `json:"btcEquity"`
	CrossedRiskRate      decimal.Decimal    `json:"crossedRiskRate"`
	UnrealizedPL         decimal.Decimal    `json:"unrealizedPL"`
	Coupon               decimal.Decimal    `json:"coupon"`
	UnionTotalMargin     decimal.Decimal    `json:"unionTotalMargin"`
	UnionAvailable       decimal.Decimal    `json:"unionAvailable"`
	UnionMm              decimal.Decimal    `json:"unionMm"`
	AssetList            []AccountAssetItem `json:"assetList"`
	IsolatedMargin       decimal.Decimal    `json:"isolatedMargin"`
	CrossedMargin        decimal.Decimal    `json:"crossedMargin"`
	CrossedUnrealizedPL  decimal.Decimal    `json:"crossedUnrealizedPL"`
	IsolatedUnrealizedPL decimal.Decimal    `json:"isolatedUnrealizedPL"`
	AssetMode            AssetMode          `json:"assetMode"`
	Grant                decimal.Decimal    `json:"grant"`
}

// AccountAssetItem is one coin holding under the union (multi-asset) mode.
type AccountAssetItem struct {
	Coin      string          `json:"coin"`
	Balance   decimal.Decimal `json:"balance"`
	Available decimal.Decimal `json:"available"`
}

// GetSubAccountAssetsService -- GET /api/v2/mix/account/sub-account-assets (private)
//
// Returns the futures assets of every sub-account under a product type.
type GetSubAccountAssetsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetSubAccountAssetsService(productType ProductType) *GetSubAccountAssetsService {
	return &GetSubAccountAssetsService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetSubAccountAssetsService) Do(ctx context.Context) ([]SubAccountAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/sub-account-assets", s.params).WithSign()
	resp, err := request.Do[[]SubAccountAssets](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubAccountAssets is one sub-account's futures asset list.
type SubAccountAssets struct {
	UserID    string                `json:"userId"`
	AssetList []SubAccountAssetItem `json:"assetList"`
}

// SubAccountAssetItem is one margin-coin balance within a sub-account.
type SubAccountAssetItem struct {
	MarginCoin           string          `json:"marginCoin"`
	Locked               decimal.Decimal `json:"locked"`
	Available            decimal.Decimal `json:"available"`
	CrossedMaxAvailable  decimal.Decimal `json:"crossedMaxAvailable"`
	IsolatedMaxAvailable decimal.Decimal `json:"isolatedMaxAvailable"`
	MaxTransferOut       decimal.Decimal `json:"maxTransferOut"`
	AccountEquity        decimal.Decimal `json:"accountEquity"`
	USDTEquity           decimal.Decimal `json:"usdtEquity"`
	BtcEquity            decimal.Decimal `json:"btcEquity"`
	UnrealizedPL         decimal.Decimal `json:"unrealizedPL"`
	Coupon               decimal.Decimal `json:"coupon"`
	CrossedUnrealizedPL  decimal.Decimal `json:"crossedUnrealizedPL"`
	IsolatedUnrealizedPL decimal.Decimal `json:"isolatedUnrealizedPL"`
	Grant                decimal.Decimal `json:"grant"`
	AssetMode            AssetMode       `json:"assetMode"`
	IsolatedMargin       decimal.Decimal `json:"isolatedMargin"`
	CrossedMargin        decimal.Decimal `json:"crossedMargin"`
}

// GetInterestHistoryService -- GET /api/v2/mix/account/interest-history (private)
//
// Returns the USDT-M futures interest (borrowing) history for the account, with
// the current debt and loan limit alongside the per-coin interest records.
type GetInterestHistoryService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetInterestHistoryService(productType ProductType) *GetInterestHistoryService {
	return &GetInterestHistoryService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetInterestHistoryService) SetCoin(coin string) *GetInterestHistoryService {
	s.params["coin"] = coin
	return s
}

func (s *GetInterestHistoryService) SetIDLessThan(id string) *GetInterestHistoryService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetInterestHistoryService) SetStartTime(t time.Time) *GetInterestHistoryService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInterestHistoryService) SetEndTime(t time.Time) *GetInterestHistoryService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetInterestHistoryService) SetLimit(limit int) *GetInterestHistoryService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetInterestHistoryService) Do(ctx context.Context) (*InterestHistory, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/interest-history", s.params).WithSign()
	return request.Do[InterestHistory](req)
}

// InterestHistory is the futures borrowing summary plus its interest records.
type InterestHistory struct {
	NextSettleTime time.Time       `json:"nextSettleTime"`
	BorrowAmount   decimal.Decimal `json:"borrowAmount"`
	BorrowLimit    decimal.Decimal `json:"borrowLimit"`
	InterestList   []InterestItem  `json:"interestList"`
	EndID          string          `json:"endId"`
}

// InterestItem is a single per-coin interest record.
type InterestItem struct {
	Coin              string          `json:"coin"`
	Liability         decimal.Decimal `json:"liability"`
	InterestFreeLimit decimal.Decimal `json:"interestFreeLimit"`
	InterestLimit     decimal.Decimal `json:"interestLimit"`
	HourInterestRate  decimal.Decimal `json:"hourInterestRate"`
	Interest          decimal.Decimal `json:"interest"`
	CTime             time.Time       `json:"cTime"`
}

// GetMaxOpenService -- GET /api/v2/mix/account/max-open (private)
//
// Returns the maximum openable quantity for a symbol given the side, order type
// and (for limit) open price, accounting for current positions and open orders.
type GetMaxOpenService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetMaxOpenService(symbol string, productType ProductType, marginCoin string, posSide HoldSide, orderType OrderType) *GetMaxOpenService {
	return &GetMaxOpenService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
		"posSide":     string(posSide),
		"orderType":   string(orderType),
	}}
}

// SetOpenPrice sets the entry price; required when orderType is limit.
func (s *GetMaxOpenService) SetOpenPrice(price decimal.Decimal) *GetMaxOpenService {
	s.params["openPrice"] = price.String()
	return s
}

func (s *GetMaxOpenService) Do(ctx context.Context) (*MaxOpen, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/max-open", s.params).WithSign()
	return request.Do[MaxOpen](req)
}

// MaxOpen is the estimated maximum openable quantity.
type MaxOpen struct {
	MaxOpen decimal.Decimal `json:"maxOpen"`
}

// GetLiquidationPriceService -- GET /api/v2/mix/account/liq-price (private)
//
// Returns the estimated liquidation price for a hypothetical position of the
// given size, side and (for limit) open price.
type GetLiquidationPriceService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetLiquidationPriceService(symbol string, productType ProductType, marginCoin string, posSide HoldSide, orderType OrderType, openAmount decimal.Decimal) *GetLiquidationPriceService {
	return &GetLiquidationPriceService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
		"posSide":     string(posSide),
		"orderType":   string(orderType),
		"openAmount":  openAmount.String(),
	}}
}

// SetOpenPrice sets the entry price; required when orderType is limit.
func (s *GetLiquidationPriceService) SetOpenPrice(price decimal.Decimal) *GetLiquidationPriceService {
	s.params["openPrice"] = price.String()
	return s
}

func (s *GetLiquidationPriceService) Do(ctx context.Context) (*LiquidationPrice, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/liq-price", s.params).WithSign()
	return request.Do[LiquidationPrice](req)
}

// LiquidationPrice is the estimated liquidation price.
type LiquidationPrice struct {
	LiqPrice decimal.Decimal `json:"liqPrice"`
}

// GetOpenCountService -- GET /api/v2/mix/account/open-count (private)
//
// Returns the estimated open size for a hypothetical position given the margin
// amount, order price and leverage.
type GetOpenCountService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetOpenCountService(symbol string, productType ProductType, marginCoin string, openAmount, openPrice decimal.Decimal) *GetOpenCountService {
	return &GetOpenCountService{c: c, params: map[string]string{
		"symbol":      symbol,
		"productType": string(productType),
		"marginCoin":  marginCoin,
		"openAmount":  openAmount.String(),
		"openPrice":   openPrice.String(),
	}}
}

// SetLeverage sets the leverage used for the estimate (default 20).
func (s *GetOpenCountService) SetLeverage(leverage string) *GetOpenCountService {
	s.params["leverage"] = leverage
	return s
}

func (s *GetOpenCountService) Do(ctx context.Context) (*OpenCount, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/open-count", s.params).WithSign()
	return request.Do[OpenCount](req)
}

// OpenCount is the estimated open size.
type OpenCount struct {
	Size decimal.Decimal `json:"size"`
}

// GetAccountBillService -- GET /api/v2/mix/account/bill (private)
//
// Returns the account's funding-flow bill records (fees, settlements, transfers,
// opens/closes) under a product type, paginated by id.
type GetAccountBillService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetAccountBillService(productType ProductType) *GetAccountBillService {
	return &GetAccountBillService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetAccountBillService) SetCoin(coin string) *GetAccountBillService {
	s.params["coin"] = coin
	return s
}

func (s *GetAccountBillService) SetBusinessType(businessType string) *GetAccountBillService {
	s.params["businessType"] = businessType
	return s
}

// SetOnlyFunding restricts the result to funding bills only ("yes"/"no").
func (s *GetAccountBillService) SetOnlyFunding(onlyFunding string) *GetAccountBillService {
	s.params["onlyFunding"] = onlyFunding
	return s
}

func (s *GetAccountBillService) SetIDLessThan(id string) *GetAccountBillService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetAccountBillService) SetStartTime(t time.Time) *GetAccountBillService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAccountBillService) SetEndTime(t time.Time) *GetAccountBillService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAccountBillService) SetLimit(limit int) *GetAccountBillService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetAccountBillService) Do(ctx context.Context) (*AccountBill, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/bill", s.params).WithSign()
	return request.Do[AccountBill](req)
}

// AccountBill is a page of bill records plus the pagination marker.
type AccountBill struct {
	Bills []BillItem `json:"bills"`
	EndID string     `json:"endId"`
}

// BillItem is a single account bill (funding flow) record.
type BillItem struct {
	BillID       string          `json:"billId"`
	Symbol       string          `json:"symbol"`
	Amount       decimal.Decimal `json:"amount"`
	Fee          decimal.Decimal `json:"fee"`
	FeeByCoupon  decimal.Decimal `json:"feeByCoupon"`
	BusinessType string          `json:"businessType"`
	Coin         string          `json:"coin"`
	Balance      decimal.Decimal `json:"balance"`
	CTime        time.Time       `json:"cTime"`
}

// GetTransferLimitsService -- GET /api/v2/mix/account/transfer-limits (private)
//
// Returns the union (multi-asset) maximum transfer-in quantity for a coin.
type GetTransferLimitsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetTransferLimitsService(coin string) *GetTransferLimitsService {
	return &GetTransferLimitsService{c: c, params: map[string]string{
		"coin": coin,
	}}
}

func (s *GetTransferLimitsService) Do(ctx context.Context) (*TransferLimits, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/transfer-limits", s.params).WithSign()
	return request.Do[TransferLimits](req)
}

// TransferLimits is the union transfer-in limit for a coin.
type TransferLimits struct {
	Coin          string          `json:"coin"`
	MaxTransferIn decimal.Decimal `json:"maxTransferIn"`
}

// GetUnionConfigService -- GET /api/v2/mix/account/union-config (private)
//
// Returns the union (multi-asset) account liability margin-rate configuration.
type GetUnionConfigService struct {
	c *MixClient
}

func (c *MixClient) NewGetUnionConfigService() *GetUnionConfigService {
	return &GetUnionConfigService{c: c}
}

func (s *GetUnionConfigService) Do(ctx context.Context) (*UnionConfig, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/union-config").WithSign()
	return request.Do[UnionConfig](req)
}

// UnionConfig is the union-mode liability margin-rate configuration.
type UnionConfig struct {
	Imr                  decimal.Decimal `json:"imr"`
	Mmr                  decimal.Decimal `json:"mmr"`
	IndividualLimit      decimal.Decimal `json:"individualLimit"`
	IndividualLimitRatio decimal.Decimal `json:"individualLimitRatio"`
}

// GetSwitchUnionUSDTService -- GET /api/v2/mix/account/switch-union-usdt (private)
//
// Returns the USDT quota required to switch the account into union (multi-asset)
// mode.
type GetSwitchUnionUSDTService struct {
	c *MixClient
}

func (c *MixClient) NewGetSwitchUnionUSDTService() *GetSwitchUnionUSDTService {
	return &GetSwitchUnionUSDTService{c: c}
}

func (s *GetSwitchUnionUSDTService) Do(ctx context.Context) (*SwitchUnionUSDT, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/switch-union-usdt").WithSign()
	return request.Do[SwitchUnionUSDT](req)
}

// SwitchUnionUSDT is the USDT quota needed for the union-mode switch.
type SwitchUnionUSDT struct {
	USDTAmount decimal.Decimal `json:"usdtAmount"`
}

// GetIsolatedSymbolsService -- GET /api/v2/mix/account/isolated-symbols (private)
//
// Returns the symbols currently configured for isolated margin under a product
// type.
type GetIsolatedSymbolsService struct {
	c      *MixClient
	params map[string]string
}

func (c *MixClient) NewGetIsolatedSymbolsService(productType ProductType) *GetIsolatedSymbolsService {
	return &GetIsolatedSymbolsService{c: c, params: map[string]string{
		"productType": string(productType),
	}}
}

func (s *GetIsolatedSymbolsService) Do(ctx context.Context) ([]IsolatedSymbol, error) {
	req := request.Get(ctx, s.c, "/api/v2/mix/account/isolated-symbols", s.params).WithSign()
	resp, err := request.Do[[]IsolatedSymbol](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// IsolatedSymbol is one isolated-margin trading pair.
type IsolatedSymbol struct {
	Symbol     string     `json:"symbol"`
	MarginMode MarginMode `json:"marginMode"`
}
