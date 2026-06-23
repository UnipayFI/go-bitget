package broker

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// SubaccountStatus is the lifecycle state of a broker sub-account.
type SubaccountStatus string

const (
	SubaccountStatusNormal SubaccountStatus = "normal"
	SubaccountStatusFreeze SubaccountStatus = "freeze"
	SubaccountStatusDel    SubaccountStatus = "del"
)

// SubaccountSpotAssetType selects which spot balances the sub-account assets
// endpoint returns.
type SubaccountSpotAssetType string

const (
	SubaccountSpotAssetTypeHoldOnly SubaccountSpotAssetType = "hold_only"
	SubaccountSpotAssetTypeAll      SubaccountSpotAssetType = "all"
)

// SubaccountFutureProductType is the futures product line a sub-account assets
// query targets.
type SubaccountFutureProductType string

const (
	SubaccountFutureProductUSDTFutures SubaccountFutureProductType = "USDT-FUTURES"
	SubaccountFutureProductCoinFutures SubaccountFutureProductType = "COIN-FUTURES"
	SubaccountFutureProductUSDCFutures SubaccountFutureProductType = "USDC-FUTURES"
)

// SubaccountAutoTransferAccountType is the destination wallet a sub-account's
// deposits are auto-transferred to.
type SubaccountAutoTransferAccountType string

const (
	SubaccountAutoTransferSpot        SubaccountAutoTransferAccountType = "spot"
	SubaccountAutoTransferUSDTFutures SubaccountAutoTransferAccountType = "usdt-futures"
	SubaccountAutoTransferCoinFutures SubaccountAutoTransferAccountType = "coin-futures"
	SubaccountAutoTransferUSDCFutures SubaccountAutoTransferAccountType = "usdc-futures"
)

// SubaccountWithdrawalDest is the routing of a sub-account withdrawal.
type SubaccountWithdrawalDest string

const (
	SubaccountWithdrawalDestOnChain          SubaccountWithdrawalDest = "on_chain"
	SubaccountWithdrawalDestInternalTransfer SubaccountWithdrawalDest = "internal_transfer"
)

// SubDepWithdrawType filters the all-sub deposit/withdrawal records by direction.
type SubDepWithdrawType string

const (
	SubDepWithdrawTypeAll        SubDepWithdrawType = "all"
	SubDepWithdrawTypeDeposit    SubDepWithdrawType = "deposit"
	SubDepWithdrawTypeWithdrawal SubDepWithdrawType = "withdrawal"
)

// GetBrokerAccountInfoService -- GET /api/v2/broker/account/info (broker, signed)
//
// Returns the broker's sub-account quota usage.
type GetBrokerAccountInfoService struct {
	c *BrokerClient
}

func (c *BrokerClient) NewGetBrokerAccountInfoService() *GetBrokerAccountInfoService {
	return &GetBrokerAccountInfoService{c: c}
}

func (s *GetBrokerAccountInfoService) Do(ctx context.Context) (*BrokerAccountInfo, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/account/info").WithSign()
	return request.Do[BrokerAccountInfo](req)
}

// BrokerAccountInfo is the broker's sub-account quota usage.
type BrokerAccountInfo struct {
	SubAccountSize    string    `json:"subAccountSize"`
	MaxSubAccountSize string    `json:"maxSubAccountSize"`
	UTime             time.Time `json:"uTime"`
}

// CreateSubaccountService -- POST /api/v2/broker/account/create-subaccount (broker, signed)
//
// Creates a broker sub-account under the given (email-style) name.
type CreateSubaccountService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewCreateSubaccountService(subaccountName string) *CreateSubaccountService {
	return &CreateSubaccountService{c: c, body: map[string]any{"subaccountName": subaccountName}}
}

// SetLabel sets an optional remark (length < 20).
func (s *CreateSubaccountService) SetLabel(label string) *CreateSubaccountService {
	s.body["label"] = label
	return s
}

func (s *CreateSubaccountService) Do(ctx context.Context) (*Subaccount, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/create-subaccount", s.body).WithSign()
	return request.Do[Subaccount](req)
}

// Subaccount describes a broker sub-account.
type Subaccount struct {
	SubUid         string           `json:"subUid"`
	SubaccountName string           `json:"subaccountName"`
	Status         SubaccountStatus `json:"status"`
	PermList       []string         `json:"permList"`
	Label          string           `json:"label"`
	Language       string           `json:"language"`
	CTime          time.Time        `json:"cTime"`
	UTime          time.Time        `json:"uTime"`
}

// GetSubaccountListService -- GET /api/v2/broker/account/subaccount-list (broker, signed)
//
// Returns a page of the broker's sub-accounts.
type GetSubaccountListService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountListService() *GetSubaccountListService {
	return &GetSubaccountListService{c: c, params: map[string]string{}}
}

// SetLimit caps the number of results (default 10, max 100).
func (s *GetSubaccountListService) SetLimit(limit int) *GetSubaccountListService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards from this subUid (older data).
func (s *GetSubaccountListService) SetIDLessThan(id string) *GetSubaccountListService {
	s.params["idLessThan"] = id
	return s
}

// SetStatus filters by sub-account status.
func (s *GetSubaccountListService) SetStatus(status SubaccountStatus) *GetSubaccountListService {
	s.params["status"] = string(status)
	return s
}

// SetStartTime filters to sub-accounts created at or after t.
func (s *GetSubaccountListService) SetStartTime(t time.Time) *GetSubaccountListService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to sub-accounts created at or before t.
func (s *GetSubaccountListService) SetEndTime(t time.Time) *GetSubaccountListService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetSubaccountListService) Do(ctx context.Context) (*SubaccountList, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/account/subaccount-list", s.params).WithSign()
	return request.Do[SubaccountList](req)
}

// SubaccountList is a paged list of broker sub-accounts.
type SubaccountList struct {
	HasNextPage bool         `json:"hasNextPage"`
	IDLessThan  string       `json:"idLessThan"`
	SubList     []Subaccount `json:"subList"`
}

// ModifySubaccountService -- POST /api/v2/broker/account/modify-subaccount (broker, signed)
//
// Updates a sub-account's permissions and status.
type ModifySubaccountService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewModifySubaccountService(subUid string, permList []string, status SubaccountStatus) *ModifySubaccountService {
	return &ModifySubaccountService{c: c, body: map[string]any{
		"subUid":   subUid,
		"permList": permList,
		"status":   string(status),
	}}
}

// SetLanguage sets the sub-account's language code.
func (s *ModifySubaccountService) SetLanguage(language string) *ModifySubaccountService {
	s.body["language"] = language
	return s
}

func (s *ModifySubaccountService) Do(ctx context.Context) (*Subaccount, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/modify-subaccount", s.body).WithSign()
	return request.Do[Subaccount](req)
}

// ModifySubaccountEmailService -- POST /api/v2/broker/account/modify-subaccount-email (broker, signed)
//
// Updates a sub-account's email. The response data is the literal "success".
type ModifySubaccountEmailService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewModifySubaccountEmailService(subUid, subaccountEmail string) *ModifySubaccountEmailService {
	return &ModifySubaccountEmailService{c: c, body: map[string]any{
		"subUid":          subUid,
		"subaccountEmail": subaccountEmail,
	}}
}

func (s *ModifySubaccountEmailService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/modify-subaccount-email", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetSubaccountEmailService -- GET /api/v2/broker/account/subaccount-email (broker, signed)
//
// Returns the email bound to a sub-account.
type GetSubaccountEmailService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountEmailService(subUid string) *GetSubaccountEmailService {
	return &GetSubaccountEmailService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetSubaccountEmailService) Do(ctx context.Context) (*SubaccountEmail, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/account/subaccount-email", s.params).WithSign()
	return request.Do[SubaccountEmail](req)
}

// SubaccountEmail is a sub-account's email binding.
type SubaccountEmail struct {
	SubUid          string    `json:"subUid"`
	SubaccountName  string    `json:"subaccountName"`
	SubaccountEmail string    `json:"subaccountEmail"`
	CTime           time.Time `json:"cTime"`
	UTime           time.Time `json:"uTime"`
}

// GetSubaccountSpotAssetsService -- GET /api/v2/broker/account/subaccount-spot-assets (broker, signed)
//
// Returns a sub-account's spot balances.
type GetSubaccountSpotAssetsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountSpotAssetsService(subUid string) *GetSubaccountSpotAssetsService {
	return &GetSubaccountSpotAssetsService{c: c, params: map[string]string{"subUid": subUid}}
}

// SetCoin filters to a single coin.
func (s *GetSubaccountSpotAssetsService) SetCoin(coin string) *GetSubaccountSpotAssetsService {
	s.params["coin"] = coin
	return s
}

// SetAssetType selects which balances to return (default hold_only).
func (s *GetSubaccountSpotAssetsService) SetAssetType(assetType SubaccountSpotAssetType) *GetSubaccountSpotAssetsService {
	s.params["assetType"] = string(assetType)
	return s
}

func (s *GetSubaccountSpotAssetsService) Do(ctx context.Context) (*SubaccountSpotAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/account/subaccount-spot-assets", s.params).WithSign()
	return request.Do[SubaccountSpotAssets](req)
}

// SubaccountSpotAssets wraps the sub-account's spot balance list.
type SubaccountSpotAssets struct {
	AssetsList []SubaccountSpotAsset `json:"assetsList"`
}

// SubaccountSpotAsset is one spot balance of a sub-account.
type SubaccountSpotAsset struct {
	Coin      string          `json:"coin"`
	Available decimal.Decimal `json:"available"`
	Frozen    decimal.Decimal `json:"frozen"`
	Locked    decimal.Decimal `json:"locked"`
	UTime     time.Time       `json:"uTime"`
}

// GetSubaccountFutureAssetsService -- GET /api/v2/broker/account/subaccount-future-assets (broker, signed)
//
// Returns a sub-account's futures balances for a product line.
type GetSubaccountFutureAssetsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountFutureAssetsService(subUid string, productType SubaccountFutureProductType) *GetSubaccountFutureAssetsService {
	return &GetSubaccountFutureAssetsService{c: c, params: map[string]string{
		"subUid":      subUid,
		"productType": string(productType),
	}}
}

func (s *GetSubaccountFutureAssetsService) Do(ctx context.Context) (*SubaccountFutureAssets, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/account/subaccount-future-assets", s.params).WithSign()
	return request.Do[SubaccountFutureAssets](req)
}

// SubaccountFutureAssets wraps the sub-account's futures balance list.
type SubaccountFutureAssets struct {
	AssetsList []SubaccountFutureAsset `json:"assetsList"`
}

// SubaccountFutureAsset is one futures balance of a sub-account.
type SubaccountFutureAsset struct {
	MarginCoin           string          `json:"marginCoin"`
	Available            decimal.Decimal `json:"available"`
	Frozen               decimal.Decimal `json:"frozen"`
	Locked               decimal.Decimal `json:"locked"`
	CrossedMaxAvailable  decimal.Decimal `json:"crossedMaxAvailable"`
	IsolatedMaxAvailable decimal.Decimal `json:"isolatedMaxAvailable"`
	MaxTransferOut       decimal.Decimal `json:"maxTransferOut"`
	AccountEquity        decimal.Decimal `json:"accountEquity"`
	UsdtEquity           decimal.Decimal `json:"usdtEquity"`
	BtcEquity            decimal.Decimal `json:"btcEquity"`
	UnrealizedPL         decimal.Decimal `json:"unrealizedPL"`
	UTime                time.Time       `json:"uTime"`
}

// CreateSubaccountDepositAddressService -- POST /api/v2/broker/account/subaccount-address (broker, signed)
//
// Creates (or fetches) a sub-account's deposit address for a coin.
type CreateSubaccountDepositAddressService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewCreateSubaccountDepositAddressService(subUid, coin string) *CreateSubaccountDepositAddressService {
	return &CreateSubaccountDepositAddressService{c: c, body: map[string]any{
		"subUid": subUid,
		"coin":   coin,
	}}
}

// SetChain selects a deposit network (defaults to the coin's main chain).
func (s *CreateSubaccountDepositAddressService) SetChain(chain string) *CreateSubaccountDepositAddressService {
	s.body["chain"] = chain
	return s
}

func (s *CreateSubaccountDepositAddressService) Do(ctx context.Context) (*SubaccountDepositAddress, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/subaccount-address", s.body).WithSign()
	return request.Do[SubaccountDepositAddress](req)
}

// SubaccountDepositAddress is a sub-account's deposit address for a coin.
type SubaccountDepositAddress struct {
	SubUid  string    `json:"subUid"`
	Address string    `json:"address"`
	Chain   string    `json:"chain"`
	Coin    string    `json:"coin"`
	Tag     string    `json:"tag"`
	URL     string    `json:"url"`
	CTime   time.Time `json:"cTime"`
}

// SubaccountWithdrawalService -- POST /api/v2/broker/account/subaccount-withdrawal (broker, signed)
//
// Withdraws from a sub-account, on-chain or via internal transfer.
type SubaccountWithdrawalService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewSubaccountWithdrawalService(subUid, coin string, dest SubaccountWithdrawalDest, address string, amount decimal.Decimal) *SubaccountWithdrawalService {
	return &SubaccountWithdrawalService{c: c, body: map[string]any{
		"subUid":  subUid,
		"coin":    coin,
		"dest":    string(dest),
		"address": address,
		"amount":  amount.String(),
	}}
}

// SetChain selects the withdrawal network (defaults to the coin's main chain).
func (s *SubaccountWithdrawalService) SetChain(chain string) *SubaccountWithdrawalService {
	s.body["chain"] = chain
	return s
}

// SetTag sets the destination tag/memo for chains that require it.
func (s *SubaccountWithdrawalService) SetTag(tag string) *SubaccountWithdrawalService {
	s.body["tag"] = tag
	return s
}

// SetClientOid sets a client-supplied order id.
func (s *SubaccountWithdrawalService) SetClientOid(clientOid string) *SubaccountWithdrawalService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *SubaccountWithdrawalService) Do(ctx context.Context) (*SubaccountWithdrawalResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/subaccount-withdrawal", s.body).WithSign()
	return request.Do[SubaccountWithdrawalResult](req)
}

// SubaccountWithdrawalResult identifies a submitted sub-account withdrawal.
type SubaccountWithdrawalResult struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// SetSubaccountAutoTransferService -- POST /api/v2/broker/account/set-subaccount-autotransfer (broker, signed)
//
// Configures auto-transfer of a sub-account's deposits to a destination wallet.
// The response data is the literal "success" or "fail".
type SetSubaccountAutoTransferService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewSetSubaccountAutoTransferService(subUid, coin string, toAccountType SubaccountAutoTransferAccountType) *SetSubaccountAutoTransferService {
	return &SetSubaccountAutoTransferService{c: c, body: map[string]any{
		"subUid":        subUid,
		"coin":          coin,
		"toAccountType": string(toAccountType),
	}}
}

func (s *SetSubaccountAutoTransferService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/account/set-subaccount-autotransfer", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetSubaccountDepositRecordsService -- GET /api/v2/broker/subaccount-deposit (broker, signed)
//
// Returns sub-account deposit records (max 3-month window).
type GetSubaccountDepositRecordsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountDepositRecordsService() *GetSubaccountDepositRecordsService {
	return &GetSubaccountDepositRecordsService{c: c, params: map[string]string{}}
}

// SetOrderId filters to a single deposit record id.
func (s *GetSubaccountDepositRecordsService) SetOrderId(orderId string) *GetSubaccountDepositRecordsService {
	s.params["orderId"] = orderId
	return s
}

// SetUserId filters to a single sub-account UID.
func (s *GetSubaccountDepositRecordsService) SetUserId(userId string) *GetSubaccountDepositRecordsService {
	s.params["userId"] = userId
	return s
}

// SetStartTime filters to records at or after t.
func (s *GetSubaccountDepositRecordsService) SetStartTime(t time.Time) *GetSubaccountDepositRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t.
func (s *GetSubaccountDepositRecordsService) SetEndTime(t time.Time) *GetSubaccountDepositRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records (default 20, max 100).
func (s *GetSubaccountDepositRecordsService) SetLimit(limit int) *GetSubaccountDepositRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards from this id (older data).
func (s *GetSubaccountDepositRecordsService) SetIDLessThan(id string) *GetSubaccountDepositRecordsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSubaccountDepositRecordsService) Do(ctx context.Context) ([]SubaccountDepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/subaccount-deposit", s.params).WithSign()
	resp, err := request.Do[[]SubaccountDepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// SubaccountDepositRecord is one sub-account deposit.
type SubaccountDepositRecord struct {
	OrderId     string          `json:"orderId"`
	TxId        string          `json:"txId"`
	Coin        string          `json:"coin"`
	Type        string          `json:"type"`
	Dest        string          `json:"dest"`
	Amount      decimal.Decimal `json:"amount"`
	Status      string          `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Fee         decimal.Decimal `json:"fee"`
	Chain       string          `json:"chain"`
	Confirm     string          `json:"confirm"`
	Tag         string          `json:"tag"`
	UserId      string          `json:"userId"`
	CTime       time.Time       `json:"cTime"`
	UTime       time.Time       `json:"uTime"`
}

// GetSubaccountWithdrawalRecordsService -- GET /api/v2/broker/subaccount-withdrawal (broker, signed)
//
// Returns sub-account withdrawal records (max 3-month window).
type GetSubaccountWithdrawalRecordsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountWithdrawalRecordsService() *GetSubaccountWithdrawalRecordsService {
	return &GetSubaccountWithdrawalRecordsService{c: c, params: map[string]string{}}
}

// SetOrderId filters to a single withdrawal record id.
func (s *GetSubaccountWithdrawalRecordsService) SetOrderId(orderId string) *GetSubaccountWithdrawalRecordsService {
	s.params["orderId"] = orderId
	return s
}

// SetUserId filters to a single sub-account UID.
func (s *GetSubaccountWithdrawalRecordsService) SetUserId(userId string) *GetSubaccountWithdrawalRecordsService {
	s.params["userId"] = userId
	return s
}

// SetStartTime filters to records at or after t.
func (s *GetSubaccountWithdrawalRecordsService) SetStartTime(t time.Time) *GetSubaccountWithdrawalRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t.
func (s *GetSubaccountWithdrawalRecordsService) SetEndTime(t time.Time) *GetSubaccountWithdrawalRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records (default 20, max 100).
func (s *GetSubaccountWithdrawalRecordsService) SetLimit(limit int) *GetSubaccountWithdrawalRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards from this id (older data).
func (s *GetSubaccountWithdrawalRecordsService) SetIDLessThan(id string) *GetSubaccountWithdrawalRecordsService {
	s.params["idLessThan"] = id
	return s
}

func (s *GetSubaccountWithdrawalRecordsService) Do(ctx context.Context) (*SubaccountWithdrawalRecords, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/subaccount-withdrawal", s.params).WithSign()
	return request.Do[SubaccountWithdrawalRecords](req)
}

// SubaccountWithdrawalRecords is a paged list of sub-account withdrawals.
type SubaccountWithdrawalRecords struct {
	ResultList []SubaccountWithdrawalRecord `json:"resultList"`
	EndId      string                       `json:"endId"`
}

// SubaccountWithdrawalRecord is one sub-account withdrawal.
type SubaccountWithdrawalRecord struct {
	OrderId     string          `json:"orderId"`
	TxId        string          `json:"txId"`
	Coin        string          `json:"coin"`
	Type        string          `json:"type"`
	Dest        string          `json:"dest"`
	Amount      decimal.Decimal `json:"amount"`
	Status      string          `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Fee         decimal.Decimal `json:"fee"`
	Chain       string          `json:"chain"`
	Confirm     string          `json:"confirm"`
	Tag         string          `json:"tag"`
	UserId      string          `json:"userId"`
	CTime       time.Time       `json:"cTime"`
	UTime       time.Time       `json:"uTime"`
}

// GetAllSubDepositWithdrawalService -- GET /api/v2/broker/all-sub-deposit-withdrawal (broker, signed)
//
// Returns deposit and withdrawal records across all sub-accounts (max 7-day
// window; defaults to yesterday UTC when both times omitted).
type GetAllSubDepositWithdrawalService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetAllSubDepositWithdrawalService() *GetAllSubDepositWithdrawalService {
	return &GetAllSubDepositWithdrawalService{c: c, params: map[string]string{}}
}

// SetStartTime filters to records at or after t (must pair with SetEndTime).
func (s *GetAllSubDepositWithdrawalService) SetStartTime(t time.Time) *GetAllSubDepositWithdrawalService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t (must pair with SetStartTime).
func (s *GetAllSubDepositWithdrawalService) SetEndTime(t time.Time) *GetAllSubDepositWithdrawalService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit caps the number of records (default 100, max 100).
func (s *GetAllSubDepositWithdrawalService) SetLimit(limit int) *GetAllSubDepositWithdrawalService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

// SetIDLessThan pages backwards from this id (older data).
func (s *GetAllSubDepositWithdrawalService) SetIDLessThan(id string) *GetAllSubDepositWithdrawalService {
	s.params["idLessThan"] = id
	return s
}

// SetType filters by direction (all, deposit, withdrawal).
func (s *GetAllSubDepositWithdrawalService) SetType(t SubDepWithdrawType) *GetAllSubDepositWithdrawalService {
	s.params["type"] = string(t)
	return s
}

func (s *GetAllSubDepositWithdrawalService) Do(ctx context.Context) (*AllSubDepositWithdrawal, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/all-sub-deposit-withdrawal", s.params).WithSign()
	return request.Do[AllSubDepositWithdrawal](req)
}

// AllSubDepositWithdrawal is a paged list of cross-sub-account transfers.
type AllSubDepositWithdrawal struct {
	List  []SubDepositWithdrawalRecord `json:"list"`
	EndId string                       `json:"endId"`
}

// SubDepositWithdrawalRecord is one cross-sub-account deposit or withdrawal.
type SubDepositWithdrawalRecord struct {
	Uid     string          `json:"uid"`
	TxId    string          `json:"txId"`
	Type    string          `json:"type"`    // deposit, withdrawal
	SubType string          `json:"subType"` // onchain, internal, fast
	Coin    string          `json:"coin"`
	Amount  decimal.Decimal `json:"amount"`
	Status  string          `json:"status"` // pending, fail, success
	Ts      time.Time       `json:"ts"`
}

// GetBrokerSubaccountsService -- GET /api/v2/broker/subaccounts (broker, signed)
//
// Returns broker sub-account summary stats (asset, first deposit/trade, register
// time) over a window (max 30 days; defaults to yesterday UTC).
type GetBrokerSubaccountsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetBrokerSubaccountsService() *GetBrokerSubaccountsService {
	return &GetBrokerSubaccountsService{c: c, params: map[string]string{}}
}

// SetStartTime filters to records at or after t (must pair with SetEndTime).
func (s *GetBrokerSubaccountsService) SetStartTime(t time.Time) *GetBrokerSubaccountsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t (must pair with SetStartTime).
func (s *GetBrokerSubaccountsService) SetEndTime(t time.Time) *GetBrokerSubaccountsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageSize sets items per page (default 100, max 1000).
func (s *GetBrokerSubaccountsService) SetPageSize(pageSize int) *GetBrokerSubaccountsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo sets the page number (default 1).
func (s *GetBrokerSubaccountsService) SetPageNo(pageNo int) *GetBrokerSubaccountsService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetBrokerSubaccountsService) Do(ctx context.Context) ([]BrokerSubaccountStat, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/subaccounts", s.params).WithSign()
	resp, err := request.Do[[]BrokerSubaccountStat](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerSubaccountStat is the summary stats for one broker sub-account.
type BrokerSubaccountStat struct {
	Uid              string          `json:"uid"`
	Asset            decimal.Decimal `json:"asset"`
	FirstTimeDeposit time.Time       `json:"firstTimeDeposit"`
	FirstTimeTrade   time.Time       `json:"firstTimeTrade"`
	RegisterTime     time.Time       `json:"registerTime"`
}

// GetBrokerCommissionsService -- GET /api/v2/broker/commissions (broker, signed)
//
// Returns per-sub-account commission breakdowns over a window (max 30 days;
// defaults to yesterday UTC).
type GetBrokerCommissionsService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetBrokerCommissionsService() *GetBrokerCommissionsService {
	return &GetBrokerCommissionsService{c: c, params: map[string]string{}}
}

// SetStartTime filters to records at or after t (must pair with SetEndTime).
func (s *GetBrokerCommissionsService) SetStartTime(t time.Time) *GetBrokerCommissionsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t (must pair with SetStartTime).
func (s *GetBrokerCommissionsService) SetEndTime(t time.Time) *GetBrokerCommissionsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageSize sets items per page (default 100, max 1000).
func (s *GetBrokerCommissionsService) SetPageSize(pageSize int) *GetBrokerCommissionsService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo sets the page number (default 1).
func (s *GetBrokerCommissionsService) SetPageNo(pageNo int) *GetBrokerCommissionsService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

// SetBizType filters by business line (spot, futures).
func (s *GetBrokerCommissionsService) SetBizType(bizType CommissionBizType) *GetBrokerCommissionsService {
	s.params["bizType"] = string(bizType)
	return s
}

// SetSubBizType filters by sub business line (spot_trade, spot_margin,
// usdt_futures, usdc_futures, coin_futures).
func (s *GetBrokerCommissionsService) SetSubBizType(subBizType CommissionSubBizType) *GetBrokerCommissionsService {
	s.params["subBizType"] = string(subBizType)
	return s
}

func (s *GetBrokerCommissionsService) Do(ctx context.Context) ([]BrokerCommission, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/commissions", s.params).WithSign()
	resp, err := request.Do[[]BrokerCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerCommission is one sub-account commission breakdown row.
type BrokerCommission struct {
	Uid             string          `json:"uid"`
	Coin            string          `json:"coin"`
	Symbol          string          `json:"symbol"`
	DealtAmount     decimal.Decimal `json:"dealtAmount"`
	TotalFee        decimal.Decimal `json:"totalFee"`
	DeductedFee     decimal.Decimal `json:"deductedFee"`
	PaidFee         decimal.Decimal `json:"paidFee"`
	MarkUpFee       decimal.Decimal `json:"markUpFee"`
	TotalCommission decimal.Decimal `json:"totalCommission"`
}

// GetBrokerTradeVolumeService -- GET /api/v2/broker/trade-volume (broker, signed)
//
// Returns per-sub-account trade volume over a window (max 30 days; defaults to
// yesterday UTC).
type GetBrokerTradeVolumeService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetBrokerTradeVolumeService() *GetBrokerTradeVolumeService {
	return &GetBrokerTradeVolumeService{c: c, params: map[string]string{}}
}

// SetStartTime filters to records at or after t (must pair with SetEndTime).
func (s *GetBrokerTradeVolumeService) SetStartTime(t time.Time) *GetBrokerTradeVolumeService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters to records at or before t (must pair with SetStartTime).
func (s *GetBrokerTradeVolumeService) SetEndTime(t time.Time) *GetBrokerTradeVolumeService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageSize sets items per page (default 100, max 1000).
func (s *GetBrokerTradeVolumeService) SetPageSize(pageSize int) *GetBrokerTradeVolumeService {
	s.params["pageSize"] = strconv.Itoa(pageSize)
	return s
}

// SetPageNo sets the page number (default 1).
func (s *GetBrokerTradeVolumeService) SetPageNo(pageNo int) *GetBrokerTradeVolumeService {
	s.params["pageNo"] = strconv.Itoa(pageNo)
	return s
}

func (s *GetBrokerTradeVolumeService) Do(ctx context.Context) ([]BrokerTradeVolume, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/trade-volume", s.params).WithSign()
	resp, err := request.Do[[]BrokerTradeVolume](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerTradeVolume is one sub-account's trade volume breakdown.
type BrokerTradeVolume struct {
	Uid          string          `json:"uid"`
	Volume       decimal.Decimal `json:"volume"`
	SpotVolume   decimal.Decimal `json:"spotVolume"`
	FutureVolume decimal.Decimal `json:"futureVolume"`
}
