package tax

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SpotTaxType is the spot transaction-record category (Bitget returns dozens of
// values, e.g. Deposit, Withdrawal, Buy, Sell, Interest, Airdrop Reward-A).
type SpotTaxType string

// FutureTaxType is the futures transaction-record category, e.g. TRANSFER_IN,
// OPEN_LONG, CLOSE_SHORT, FORCE_CLOSE_LONG, RISK_LIQ_USER_IN.
type FutureTaxType string

const (
	FutureTaxTransferIn                   FutureTaxType = "TRANSFER_IN"
	FutureTaxTransferOut                  FutureTaxType = "TRANSFER_OUT"
	FutureTaxOrderDealtFrozenOut          FutureTaxType = "ORDER_DEALT_FROZEN_OUT"
	FutureTaxOrderDealtIn                 FutureTaxType = "ORDER_DEALT_IN"
	FutureTaxOrderPlfFeeOut               FutureTaxType = "ORDER_PLF_FEE_OUT"
	FutureTaxExchangeSourceTokenUserOut   FutureTaxType = "EXCHANGE_SOURCE_TOKEN_USER_OUT"
	FutureTaxExchangeTargetTokenUserIn    FutureTaxType = "EXCHANGE_TARGET_TOKEN_USER_IN"
	FutureTaxOpenLong                     FutureTaxType = "OPEN_LONG"
	FutureTaxOpenShort                    FutureTaxType = "OPEN_SHORT"
	FutureTaxBuyDeal                      FutureTaxType = "BUY_DEAL"
	FutureTaxSellDeal                     FutureTaxType = "SELL_DEAL"
	FutureTaxCloseLong                    FutureTaxType = "CLOSE_LONG"
	FutureTaxCloseShort                   FutureTaxType = "CLOSE_SHORT"
	FutureTaxForceCloseLong               FutureTaxType = "FORCE_CLOSE_LONG"
	FutureTaxForceCloseShort              FutureTaxType = "FORCE_CLOSE_SHORT"
	FutureTaxBurstCloseLong               FutureTaxType = "BURST_CLOSE_LONG"
	FutureTaxBurstCloseShort              FutureTaxType = "BURST_CLOSE_SHORT"
	FutureTaxOffsetReduceCloseLong        FutureTaxType = "OFFSET_REDUCE_CLOSE_LONG"
	FutureTaxOffsetReduceCloseShort       FutureTaxType = "OFFSET_REDUCE_CLOSE_SHORT"
	FutureTaxForceBuySSM                  FutureTaxType = "FORCE_BUY_SSM"
	FutureTaxForceSellSSM                 FutureTaxType = "FORCE_SELL_SSM"
	FutureTaxBurstBuySSM                  FutureTaxType = "BURST_BUY_SSM"
	FutureTaxBurstSellSSM                 FutureTaxType = "BURST_SELL_SSM"
	FutureTaxRiskLiqUserIn                FutureTaxType = "RISK_LIQ_USER_IN"
	FutureTaxRiskLiqUserOut               FutureTaxType = "RISK_LIQ_USER_OUT"
	FutureTaxInterestSettlementOut        FutureTaxType = "INTEREST_SETTLEMENT_OUT"
	FutureTaxContractMainSettleFeeUserIn  FutureTaxType = "CONTRACT_MAIN_SETTLE_FEE_USER_IN"
	FutureTaxContractMainSettleFeeUserOut FutureTaxType = "CONTRACT_MAIN_SETTLE_FEE_USER_OUT"
)

// MarginTaxType is the margin transaction-record category.
type MarginTaxType string

const (
	MarginTaxTransferIn     MarginTaxType = "transfer_in"
	MarginTaxTransferOut    MarginTaxType = "transfer_out"
	MarginTaxBorrow         MarginTaxType = "borrow"
	MarginTaxRepay          MarginTaxType = "repay"
	MarginTaxLiquidationFee MarginTaxType = "liquidation_fee"
	MarginTaxCompensate     MarginTaxType = "compensate"
	MarginTaxDealIn         MarginTaxType = "deal_in"
	MarginTaxDealOut        MarginTaxType = "deal_out"
	MarginTaxInterestRepay  MarginTaxType = "interest_repay"
	MarginTaxConfiscated    MarginTaxType = "confiscated"
	MarginTaxExchangeIn     MarginTaxType = "exchange_in"
	MarginTaxExchangeOut    MarginTaxType = "exchange_out"
)

// P2PTaxType is the P2P transaction-record category.
type P2PTaxType string

const (
	P2PTaxTransferIn  P2PTaxType = "transfer_in"
	P2PTaxTransferOut P2PTaxType = "transfer_out"
	P2PTaxSell        P2PTaxType = "sell"
	P2PTaxBuy         P2PTaxType = "buy"
)

// TaxFutureProductType is the futures product line for the future-record query.
type TaxFutureProductType string

const (
	TaxProductUSDTFutures TaxFutureProductType = "USDT-FUTURES"
	TaxProductCoinFutures TaxFutureProductType = "COIN-FUTURES"
	TaxProductUSDCFutures TaxFutureProductType = "USDC-FUTURES"
)

// TaxMarginType is the leverage type for the margin-record query.
type TaxMarginType string

const (
	TaxMarginIsolated TaxMarginType = "isolated"
	TaxMarginCrossed  TaxMarginType = "crossed"
)

// GetSpotRecordService -- GET /api/v2/tax/spot-record (signed; tax read)
//
// Returns spot-account transaction records for tax reporting within a time
// window (startTime and endTime, max 30-day interval).
type GetSpotRecordService struct {
	c      *TaxClient
	params map[string]string
}

func (c *TaxClient) NewGetSpotRecordService(startTime, endTime time.Time) *GetSpotRecordService {
	return &GetSpotRecordService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetSpotRecordService) SetCoin(coin string) *GetSpotRecordService {
	s.params["coin"] = coin
	return s
}

func (s *GetSpotRecordService) SetLimit(limit int) *GetSpotRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSpotRecordService) SetIDLessThan(idLessThan string) *GetSpotRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetSpotRecordService) Do(ctx context.Context) ([]SpotRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/tax/spot-record", s.params).WithSign()
	resp, err := request.Do[[]SpotRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SpotRecord is a single spot-account transaction record.
type SpotRecord struct {
	ID          string          `json:"id"`
	Coin        string          `json:"coin"`
	SpotTaxType SpotTaxType     `json:"spotTaxType"`
	Amount      decimal.Decimal `json:"amount"`
	Fee         decimal.Decimal `json:"fee"`
	Balance     decimal.Decimal `json:"balance"`
	BizOrderID  string          `json:"bizOrderId"`
	Ts          time.Time       `json:"ts"`
}

// GetFutureRecordService -- GET /api/v2/tax/future-record (signed; tax read)
//
// Returns futures-account transaction records for tax reporting within a time
// window (startTime and endTime, max 30-day interval).
type GetFutureRecordService struct {
	c      *TaxClient
	params map[string]string
}

func (c *TaxClient) NewGetFutureRecordService(startTime, endTime time.Time) *GetFutureRecordService {
	return &GetFutureRecordService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetFutureRecordService) SetProductType(productType TaxFutureProductType) *GetFutureRecordService {
	s.params["productType"] = string(productType)
	return s
}

func (s *GetFutureRecordService) SetMarginCoin(marginCoin string) *GetFutureRecordService {
	s.params["marginCoin"] = marginCoin
	return s
}

func (s *GetFutureRecordService) SetLimit(limit int) *GetFutureRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetFutureRecordService) SetIDLessThan(idLessThan string) *GetFutureRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetFutureRecordService) Do(ctx context.Context) ([]FutureRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/tax/future-record", s.params).WithSign()
	resp, err := request.Do[[]FutureRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// FutureRecord is a single futures-account transaction record.
type FutureRecord struct {
	ID            string          `json:"id"`
	Symbol        string          `json:"symbol"`
	MarginCoin    string          `json:"marginCoin"`
	FutureTaxType FutureTaxType   `json:"futureTaxType"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
	Ts            time.Time       `json:"ts"`
}

// GetMarginRecordService -- GET /api/v2/tax/margin-record (signed; tax read)
//
// Returns margin-account transaction records for tax reporting within a time
// window (startTime and endTime, max 30-day interval).
type GetMarginRecordService struct {
	c      *TaxClient
	params map[string]string
}

func (c *TaxClient) NewGetMarginRecordService(startTime, endTime time.Time) *GetMarginRecordService {
	return &GetMarginRecordService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetMarginRecordService) SetMarginType(marginType TaxMarginType) *GetMarginRecordService {
	s.params["marginType"] = string(marginType)
	return s
}

func (s *GetMarginRecordService) SetCoin(coin string) *GetMarginRecordService {
	s.params["coin"] = coin
	return s
}

func (s *GetMarginRecordService) SetLimit(limit int) *GetMarginRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetMarginRecordService) SetIDLessThan(idLessThan string) *GetMarginRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetMarginRecordService) Do(ctx context.Context) ([]MarginRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/tax/margin-record", s.params).WithSign()
	resp, err := request.Do[[]MarginRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// MarginRecord is a single margin-account transaction record.
type MarginRecord struct {
	ID            string          `json:"id"`
	Coin          string          `json:"coin"`
	Symbol        string          `json:"symbol"`
	MarginTaxType MarginTaxType   `json:"marginTaxType"`
	Amount        decimal.Decimal `json:"amount"`
	Fee           decimal.Decimal `json:"fee"`
	Total         decimal.Decimal `json:"total"`
	Ts            time.Time       `json:"ts"`
}

// GetP2PRecordService -- GET /api/v2/tax/p2p-record (signed; tax read)
//
// Returns P2P-account transaction records for tax reporting within a time
// window (startTime and endTime, max 30-day interval).
type GetP2PRecordService struct {
	c      *TaxClient
	params map[string]string
}

func (c *TaxClient) NewGetP2PRecordService(startTime, endTime time.Time) *GetP2PRecordService {
	return &GetP2PRecordService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}}
}

func (s *GetP2PRecordService) SetCoin(coin string) *GetP2PRecordService {
	s.params["coin"] = coin
	return s
}

func (s *GetP2PRecordService) SetLimit(limit int) *GetP2PRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetP2PRecordService) SetIDLessThan(idLessThan string) *GetP2PRecordService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetP2PRecordService) Do(ctx context.Context) ([]P2PRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/tax/p2p-record", s.params).WithSign()
	resp, err := request.Do[[]P2PRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// P2PRecord is a single P2P-account transaction record.
type P2PRecord struct {
	ID         string          `json:"id"`
	Coin       string          `json:"coin"`
	P2PTaxType P2PTaxType      `json:"p2pTaxType"`
	Balance    decimal.Decimal `json:"balance"`
	Ts         time.Time       `json:"ts"`
}
