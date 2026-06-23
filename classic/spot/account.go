package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// GetAccountInfoService -- GET /api/v2/spot/account/info (private)
//
// Returns the account's identity and permission metadata (user/inviter IDs, IP
// whitelist, authorities, affiliate info, register time).
type GetAccountInfoService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetAccountInfoService() *GetAccountInfoService {
	return &GetAccountInfoService{c: c}
}

func (s *GetAccountInfoService) Do(ctx context.Context) (*AccountInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/info").WithSign()
	return request.Do[AccountInfo](req)
}

// AccountInfo is the account identity and permission metadata.
type AccountInfo struct {
	UserId      string    `json:"userId"`
	InviterId   string    `json:"inviterId"`
	Ips         []string  `json:"ips"`
	Authorities []string  `json:"authorities"`
	ParentId    int64     `json:"parentId"`
	TraderType  string    `json:"traderType"` // trader, not_trader
	ChannelCode string    `json:"channelCode"`
	Channel     string    `json:"channel"`
	RegisTime   time.Time `json:"regisTime"`
}

// GetAccountAssetsService -- GET /api/v2/spot/account/assets (private)
//
// Returns the spot account's per-coin balances. With no filter it returns the
// coins currently held; SetAssetType("all") returns every supported coin.
type GetAccountAssetsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAccountAssetsService() *GetAccountAssetsService {
	return &GetAccountAssetsService{c: c, params: map[string]string{}}
}

// SetCoin filters the result to a single coin (e.g. USDT).
func (s *GetAccountAssetsService) SetCoin(coin string) *GetAccountAssetsService {
	s.params["coin"] = coin
	return s
}

// SetAssetType selects which coins to return: hold_only (default) or all.
func (s *GetAccountAssetsService) SetAssetType(assetType string) *GetAccountAssetsService {
	s.params["assetType"] = assetType
	return s
}

func (s *GetAccountAssetsService) Do(ctx context.Context) ([]AccountAsset, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/assets", s.params).WithSign()
	resp, err := request.Do[[]AccountAsset](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountAsset is one coin's spot balance breakdown.
type AccountAsset struct {
	Coin           string          `json:"coin"`
	Available      decimal.Decimal `json:"available"`
	Frozen         decimal.Decimal `json:"frozen"`
	Locked         decimal.Decimal `json:"locked"`
	LimitAvailable decimal.Decimal `json:"limitAvailable"`
	UTime          time.Time       `json:"uTime"`
}

// GetSubaccountAssetsService -- GET /api/v2/spot/account/subaccount-assets (private)
//
// Returns the spot balances of every sub-account holding non-zero assets,
// paginated by sub-account.
type GetSubaccountAssetsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSubaccountAssetsService() *GetSubaccountAssetsService {
	return &GetSubaccountAssetsService{c: c, params: map[string]string{}}
}

// SetIdLessThan sets the pagination cursor (the last id returned previously).
func (s *GetSubaccountAssetsService) SetIdLessThan(id string) *GetSubaccountAssetsService {
	s.params["idLessThan"] = id
	return s
}

// SetLimit caps the number of sub-accounts per page (default 10, max 50).
func (s *GetSubaccountAssetsService) SetLimit(limit int) *GetSubaccountAssetsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSubaccountAssetsService) Do(ctx context.Context) ([]SubaccountAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/subaccount-assets", s.params).WithSign()
	resp, err := request.Do[[]SubaccountAssets](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubaccountAssets is one sub-account's spot balances.
type SubaccountAssets struct {
	Id         string            `json:"id"`
	UserId     string            `json:"userId"`
	AssetsList []SubaccountAsset `json:"assetsList"`
}

// SubaccountAsset is one coin's balance within a sub-account.
type SubaccountAsset struct {
	Coin           string          `json:"coin"`
	Available      decimal.Decimal `json:"available"`
	LimitAvailable decimal.Decimal `json:"limitAvailable"`
	Frozen         decimal.Decimal `json:"frozen"`
	Locked         decimal.Decimal `json:"locked"`
	UTime          time.Time       `json:"uTime"`
}

// GetAccountBillsService -- GET /api/v2/spot/account/bills (private)
//
// Returns the spot account's financial flow records (deposits, withdrawals,
// trades, transfers, fees, ...), filterable by coin/type and a <=90-day window.
type GetAccountBillsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetAccountBillsService() *GetAccountBillsService {
	return &GetAccountBillsService{c: c, params: map[string]string{}}
}

// SetCoin filters bills to a single coin.
func (s *GetAccountBillsService) SetCoin(coin string) *GetAccountBillsService {
	s.params["coin"] = coin
	return s
}

// SetGroupType filters by billing group (deposit, withdraw, transaction,
// transfer, loan, financial, fait, convert, c2c, pre_c2c, on_chain, strategy,
// other).
func (s *GetAccountBillsService) SetGroupType(groupType string) *GetAccountBillsService {
	s.params["groupType"] = groupType
	return s
}

// SetBusinessType filters by business type (DEPOSIT, WITHDRAW, BUY, SELL,
// DEDUCTION_HANDLING_FEE, TRANSFER_IN, TRANSFER_OUT, ...).
func (s *GetAccountBillsService) SetBusinessType(businessType string) *GetAccountBillsService {
	s.params["businessType"] = businessType
	return s
}

// SetStartTime filters bills at or after t (interval with endTime <= 90 days).
func (s *GetAccountBillsService) SetStartTime(t time.Time) *GetAccountBillsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters bills at or before t (interval with startTime <= 90 days).
func (s *GetAccountBillsService) SetEndTime(t time.Time) *GetAccountBillsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records returned (default 100, max 500).
func (s *GetAccountBillsService) SetLimit(limit int) *GetAccountBillsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIdLessThan sets the pagination cursor (billId for older data).
func (s *GetAccountBillsService) SetIdLessThan(id string) *GetAccountBillsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetAccountBillsService) Do(ctx context.Context) ([]AccountBill, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/bills", s.params).WithSign()
	resp, err := request.Do[[]AccountBill](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// AccountBill is one financial flow record.
type AccountBill struct {
	CTime        time.Time       `json:"cTime"`
	Coin         string          `json:"coin"`
	GroupType    string          `json:"groupType"`
	BusinessType string          `json:"businessType"`
	Size         decimal.Decimal `json:"size"`
	Balance      decimal.Decimal `json:"balance"`
	Fees         decimal.Decimal `json:"fees"`
	BillId       string          `json:"billId"`
	BizOrderId   string          `json:"bizOrderId"`
}

// GetTransferCoinInfoService -- GET /api/v2/spot/wallet/transfer-coin-info (private)
//
// Returns the list of coins that can be transferred between the two given
// account types (the intersection of coins each side supports).
type GetTransferCoinInfoService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTransferCoinInfoService(fromType, toType string) *GetTransferCoinInfoService {
	return &GetTransferCoinInfoService{c: c, params: map[string]string{
		"fromType": fromType,
		"toType":   toType,
	}}
}

func (s *GetTransferCoinInfoService) Do(ctx context.Context) ([]string, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/transfer-coin-info", s.params).WithSign()
	resp, err := request.Do[[]string](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetDeductInfoService -- GET /api/v2/spot/account/deduct-info (private)
//
// Returns whether trading-fee deduction with BGB is currently enabled.
type GetDeductInfoService struct {
	c *SpotClient
}

func (c *SpotClient) NewGetDeductInfoService() *GetDeductInfoService {
	return &GetDeductInfoService{c: c}
}

func (s *GetDeductInfoService) Do(ctx context.Context) (*DeductInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/deduct-info").WithSign()
	return request.Do[DeductInfo](req)
}

// DeductInfo reports the BGB fee-deduction switch state.
type DeductInfo struct {
	Deduct string `json:"deduct"` // on, off
}

// SwitchDeductService -- POST /api/v2/spot/account/switch-deduct (private, state-changing)
//
// Turns BGB trading-fee deduction on or off.
type SwitchDeductService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewSwitchDeductService(deduct string) *SwitchDeductService {
	return &SwitchDeductService{c: c, body: map[string]any{"deduct": deduct}}
}

func (s *SwitchDeductService) Do(ctx context.Context) (*bool, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/account/switch-deduct", s.body).WithSign()
	return request.Do[bool](req)
}

// GetTransferRecordsService -- GET /api/v2/spot/account/transferRecords (private)
//
// Returns the account's internal transfer records for a coin, optionally
// filtered by source account type and a <=90-day window.
type GetTransferRecordsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetTransferRecordsService(coin string) *GetTransferRecordsService {
	return &GetTransferRecordsService{c: c, params: map[string]string{"coin": coin}}
}

// SetFromType filters by source account type (spot, p2p, coin_futures,
// usdt_futures, usdc_futures, crossed_margin, isolated_margin).
func (s *GetTransferRecordsService) SetFromType(fromType string) *GetTransferRecordsService {
	s.params["fromType"] = fromType
	return s
}

// SetStartTime filters records at or after t (interval with endTime <= 90 days).
func (s *GetTransferRecordsService) SetStartTime(t time.Time) *GetTransferRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t.
func (s *GetTransferRecordsService) SetEndTime(t time.Time) *GetTransferRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetClientOid filters by user-customized order ID.
func (s *GetTransferRecordsService) SetClientOid(clientOid string) *GetTransferRecordsService {
	s.params["clientOid"] = clientOid
	return s
}

// SetPageNum sets the (deprecated) page number (default 1, max 1000).
func (s *GetTransferRecordsService) SetPageNum(pageNum int) *GetTransferRecordsService {
	s.params["pageNum"] = strconv.Itoa(pageNum)
	return s
}

// SetLimit caps the number of records returned (default 100, max 500).
func (s *GetTransferRecordsService) SetLimit(limit int) *GetTransferRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIdLessThan sets the pagination cursor (transferId for older data).
func (s *GetTransferRecordsService) SetIdLessThan(id string) *GetTransferRecordsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetTransferRecordsService) Do(ctx context.Context) ([]TransferRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/transferRecords", s.params).WithSign()
	resp, err := request.Do[[]TransferRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// TransferRecord is one internal transfer. fromSymbol/toSymbol are only set when
// the corresponding side is an isolated_margin account.
type TransferRecord struct {
	Coin       string          `json:"coin"`
	Status     string          `json:"status"` // Successful, Failed, Processing
	ToType     string          `json:"toType"`
	ToSymbol   string          `json:"toSymbol"`
	FromType   string          `json:"fromType"`
	FromSymbol string          `json:"fromSymbol"`
	Size       decimal.Decimal `json:"size"`
	Ts         time.Time       `json:"ts"`
	ClientOid  string          `json:"clientOid"`
	TransferId string          `json:"transferId"`
}

// GetSubMainTransRecordService -- GET /api/v2/spot/account/sub-main-trans-record (private)
//
// Returns transfer records between the main account and its sub-accounts.
type GetSubMainTransRecordService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSubMainTransRecordService() *GetSubMainTransRecordService {
	return &GetSubMainTransRecordService{c: c, params: map[string]string{}}
}

// SetCoin filters by token name.
func (s *GetSubMainTransRecordService) SetCoin(coin string) *GetSubMainTransRecordService {
	s.params["coin"] = coin
	return s
}

// SetRole filters by the caller's role in the transfer: initiator (default) or
// receiver.
func (s *GetSubMainTransRecordService) SetRole(role string) *GetSubMainTransRecordService {
	s.params["role"] = role
	return s
}

// SetSubUid filters to a single sub-account UID (main accounts only).
func (s *GetSubMainTransRecordService) SetSubUid(subUid string) *GetSubMainTransRecordService {
	s.params["subUid"] = subUid
	return s
}

// SetStartTime filters records at or after t (interval with endTime <= 90 days).
func (s *GetSubMainTransRecordService) SetStartTime(t time.Time) *GetSubMainTransRecordService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t.
func (s *GetSubMainTransRecordService) SetEndTime(t time.Time) *GetSubMainTransRecordService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetClientOid filters by user-defined order ID.
func (s *GetSubMainTransRecordService) SetClientOid(clientOid string) *GetSubMainTransRecordService {
	s.params["clientOid"] = clientOid
	return s
}

// SetLimit caps the number of records returned (default 100, max 100).
func (s *GetSubMainTransRecordService) SetLimit(limit int) *GetSubMainTransRecordService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIdLessThan sets the pagination cursor (transferId for older data).
func (s *GetSubMainTransRecordService) SetIdLessThan(id string) *GetSubMainTransRecordService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSubMainTransRecordService) Do(ctx context.Context) ([]SubMainTransferRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/sub-main-trans-record", s.params).WithSign()
	resp, err := request.Do[[]SubMainTransferRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubMainTransferRecord is one main<->sub transfer record.
type SubMainTransferRecord struct {
	Coin          string          `json:"coin"`
	Status        string          `json:"status"` // Successful, Failed, Processing
	ToType        string          `json:"toType"`
	FromType      string          `json:"fromType"`
	Size          decimal.Decimal `json:"size"`
	Ts            time.Time       `json:"ts"`
	ClientOid     string          `json:"clientOid"`
	TransferId    string          `json:"transferId"`
	NewTransferId string          `json:"newTransferId"`
	FromUserId    string          `json:"fromUserId"`
	ToUserId      string          `json:"toUserId"`
}

// UpgradeAccountService -- POST /api/v2/spot/account/upgrade (private, state-changing)
//
// Triggers the one-way upgrade of the account. The request takes no body; poll
// GetUpgradeStatusService for the outcome.
type UpgradeAccountService struct {
	c *SpotClient
}

func (c *SpotClient) NewUpgradeAccountService() *UpgradeAccountService {
	return &UpgradeAccountService{c: c}
}

func (s *UpgradeAccountService) Do(ctx context.Context) (*UpgradeAccountResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/account/upgrade").WithSign()
	return request.Do[UpgradeAccountResult](req)
}

// UpgradeAccountResult reports the upgrade request outcome.
type UpgradeAccountResult struct {
	Status string `json:"status"` // process, success, fail
}

// GetUpgradeStatusService -- GET /api/v2/spot/account/upgrade-status (private)
//
// Returns the current account-upgrade status.
type GetUpgradeStatusService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetUpgradeStatusService() *GetUpgradeStatusService {
	return &GetUpgradeStatusService{c: c, params: map[string]string{}}
}

// SetSubUid queries the upgrade status of a specific sub-account.
func (s *GetUpgradeStatusService) SetSubUid(subUid string) *GetUpgradeStatusService {
	s.params["subUid"] = subUid
	return s
}

func (s *GetUpgradeStatusService) Do(ctx context.Context) (*UpgradeStatus, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/account/upgrade-status", s.params).WithSign()
	return request.Do[UpgradeStatus](req)
}

// UpgradeStatus is the account-upgrade state. Reason is set only on failure.
type UpgradeStatus struct {
	Status string `json:"status"` // process, success, fail
	Reason string `json:"reason"`
}
