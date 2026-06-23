package spot

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// AccountType is a wallet/transfer account bucket on the classic account.
type AccountType string

const (
	AccountTypeSpot           AccountType = "spot"
	AccountTypeP2P            AccountType = "p2p"
	AccountTypeCoinFutures    AccountType = "coin_futures"
	AccountTypeUsdtFutures    AccountType = "usdt_futures"
	AccountTypeUsdcFutures    AccountType = "usdc_futures"
	AccountTypeCrossedMargin  AccountType = "crossed_margin"
	AccountTypeIsolatedMargin AccountType = "isolated_margin"
)

// DepositAccountType is the destination account for a coin's deposits, as used
// by the Modify Deposit Account endpoint (note the hyphenated futures spellings,
// which differ from AccountType).
type DepositAccountType string

const (
	DepositAccountTypeSpot        DepositAccountType = "spot"
	DepositAccountTypeFunding     DepositAccountType = "funding"
	DepositAccountTypeCoinFutures DepositAccountType = "coin-futures"
	DepositAccountTypeUsdtFutures DepositAccountType = "usdt-futures"
	DepositAccountTypeUsdcFutures DepositAccountType = "usdc-futures"
)

// WithdrawalTransferType selects on-chain vs internal (within-Bitget) withdrawal.
type WithdrawalTransferType string

const (
	WithdrawalTransferTypeOnChain          WithdrawalTransferType = "on_chain"
	WithdrawalTransferTypeInternalTransfer WithdrawalTransferType = "internal_transfer"
)

// WithdrawalInnerToType identifies the recipient form for an internal transfer.
type WithdrawalInnerToType string

const (
	WithdrawalInnerToTypeEmail  WithdrawalInnerToType = "email"
	WithdrawalInnerToTypeMobile WithdrawalInnerToType = "mobile"
	WithdrawalInnerToTypeUid    WithdrawalInnerToType = "uid"
)

// RecordDest is the channel a deposit/withdrawal travelled over.
type RecordDest string

const (
	RecordDestOnChain          RecordDest = "on_chain"
	RecordDestInternalTransfer RecordDest = "internal_transfer"
)

// RecordStatus is the lifecycle state of a deposit or withdrawal record.
type RecordStatus string

const (
	RecordStatusPending RecordStatus = "pending"
	RecordStatusFail    RecordStatus = "fail"
	RecordStatusSuccess RecordStatus = "success"
)

// TransferService -- POST /api/v2/spot/wallet/transfer (signed)
//
// Moves funds between the caller's own account buckets (spot, futures, margin,
// p2p, ...).
type TransferService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewTransferService(fromType, toType AccountType, amount decimal.Decimal, coin string) *TransferService {
	return &TransferService{c: c, body: map[string]any{
		"fromType": string(fromType),
		"toType":   string(toType),
		"amount":   amount.String(),
		"coin":     coin,
	}}
}

// SetSymbol sets the trading pair, required when transferring to/from an
// isolated-margin account.
func (s *TransferService) SetSymbol(symbol string) *TransferService {
	s.body["symbol"] = symbol
	return s
}

// SetClientOid sets a user-defined identifier; replaying the same value returns
// the original transfer result.
func (s *TransferService) SetClientOid(clientOid string) *TransferService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *TransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/wallet/transfer", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// TransferResult is returned by both the own-account and sub-account transfer
// endpoints.
type TransferResult struct {
	TransferId string `json:"transferId"`
	ClientOid  string `json:"clientOid"`
}

// SubAccountTransferService -- POST /api/v2/spot/wallet/subaccount-transfer (signed)
//
// Moves funds between a main account and its sub-accounts (or between two
// sub-accounts), identified by UID.
type SubAccountTransferService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewSubAccountTransferService(fromUserId, toUserId string, fromType, toType AccountType, amount decimal.Decimal, coin string) *SubAccountTransferService {
	return &SubAccountTransferService{c: c, body: map[string]any{
		"fromUserId": fromUserId,
		"toUserId":   toUserId,
		"fromType":   string(fromType),
		"toType":     string(toType),
		"amount":     amount.String(),
		"coin":       coin,
	}}
}

// SetSymbol sets the trading pair, required when transferring to/from an
// isolated-margin (spot) account.
func (s *SubAccountTransferService) SetSymbol(symbol string) *SubAccountTransferService {
	s.body["symbol"] = symbol
	return s
}

// SetClientOid sets a user-defined identifier.
func (s *SubAccountTransferService) SetClientOid(clientOid string) *SubAccountTransferService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *SubAccountTransferService) Do(ctx context.Context) (*TransferResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/wallet/subaccount-transfer", s.body).WithSign()
	return request.Do[TransferResult](req)
}

// ModifyDepositAccountService -- POST /api/v2/spot/wallet/modify-deposit-account (signed)
//
// Sets which account type a given coin's deposits are credited to. The response
// data is the bare string "success" or "fail".
type ModifyDepositAccountService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewModifyDepositAccountService(accountType DepositAccountType, coin string) *ModifyDepositAccountService {
	return &ModifyDepositAccountService{c: c, body: map[string]any{
		"accountType": string(accountType),
		"coin":        coin,
	}}
}

func (s *ModifyDepositAccountService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/wallet/modify-deposit-account", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// WithdrawalService -- POST /api/v2/spot/wallet/withdrawal (signed)
//
// Initiates an on-chain or internal withdrawal of a coin.
type WithdrawalService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewWithdrawalService(coin string, transferType WithdrawalTransferType, address string, size decimal.Decimal) *WithdrawalService {
	return &WithdrawalService{c: c, body: map[string]any{
		"coin":         coin,
		"transferType": string(transferType),
		"address":      address,
		"size":         size.String(),
	}}
}

// SetChain sets the network type (e.g. erc20, trc20); required for on-chain
// withdrawals.
func (s *WithdrawalService) SetChain(chain string) *WithdrawalService {
	s.body["chain"] = chain
	return s
}

// SetInnerToType sets the recipient form for an internal transfer (default uid).
func (s *WithdrawalService) SetInnerToType(innerToType WithdrawalInnerToType) *WithdrawalService {
	s.body["innerToType"] = string(innerToType)
	return s
}

// SetAreaCode sets the phone area code, required when innerToType is mobile.
func (s *WithdrawalService) SetAreaCode(areaCode string) *WithdrawalService {
	s.body["areaCode"] = areaCode
	return s
}

// SetTag sets the address tag/memo required by some coins (e.g. EOS).
func (s *WithdrawalService) SetTag(tag string) *WithdrawalService {
	s.body["tag"] = tag
	return s
}

// SetRemark sets a free-form note.
func (s *WithdrawalService) SetRemark(remark string) *WithdrawalService {
	s.body["remark"] = remark
	return s
}

// SetClientOid sets a client-generated unique identifier.
func (s *WithdrawalService) SetClientOid(clientOid string) *WithdrawalService {
	s.body["clientOid"] = clientOid
	return s
}

// SetMemberCode sets the Korea-KYC exchange code (bithumb, korbit, coinone).
func (s *WithdrawalService) SetMemberCode(memberCode string) *WithdrawalService {
	s.body["memberCode"] = memberCode
	return s
}

// SetIdentityType sets the Korea-KYC identity type (user or company).
func (s *WithdrawalService) SetIdentityType(identityType string) *WithdrawalService {
	s.body["identityType"] = identityType
	return s
}

// SetCompanyName sets the Korea-KYC company name.
func (s *WithdrawalService) SetCompanyName(companyName string) *WithdrawalService {
	s.body["companyName"] = companyName
	return s
}

// SetFirstName sets the Korea-KYC first name.
func (s *WithdrawalService) SetFirstName(firstName string) *WithdrawalService {
	s.body["firstName"] = firstName
	return s
}

// SetLastName sets the Korea-KYC last name.
func (s *WithdrawalService) SetLastName(lastName string) *WithdrawalService {
	s.body["lastName"] = lastName
	return s
}

func (s *WithdrawalService) Do(ctx context.Context) (*WithdrawalResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/wallet/withdrawal", s.body).WithSign()
	return request.Do[WithdrawalResult](req)
}

// WithdrawalResult identifies a newly-created withdrawal.
type WithdrawalResult struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
}

// CancelWithdrawalService -- POST /api/v2/spot/wallet/cancel-withdrawal (signed)
//
// Cancels a withdrawal during its brief grace period. The response data is the
// bare string "success" or "fail".
type CancelWithdrawalService struct {
	c    *SpotClient
	body map[string]any
}

func (c *SpotClient) NewCancelWithdrawalService(orderId string) *CancelWithdrawalService {
	return &CancelWithdrawalService{c: c, body: map[string]any{
		"orderId": orderId,
	}}
}

func (s *CancelWithdrawalService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/spot/wallet/cancel-withdrawal", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetDepositAddressService -- GET /api/v2/spot/wallet/deposit-address (signed)
//
// Returns the deposit address for a coin on a given chain.
type GetDepositAddressService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetDepositAddressService(coin string) *GetDepositAddressService {
	return &GetDepositAddressService{c: c, params: map[string]string{"coin": coin}}
}

// SetChain selects the chain (e.g. trc20); defaults to the coin's primary chain.
func (s *GetDepositAddressService) SetChain(chain string) *GetDepositAddressService {
	s.params["chain"] = chain
	return s
}

// SetSize sets the Bitcoin Lightning Network withdrawal amount (0.000001 - 0.01).
func (s *GetDepositAddressService) SetSize(size decimal.Decimal) *GetDepositAddressService {
	s.params["size"] = size.String()
	return s
}

func (s *GetDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/deposit-address", s.params).WithSign()
	return request.Do[DepositAddress](req)
}

// DepositAddress is a coin's deposit address on a specific chain.
type DepositAddress struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag"`
	Url     string `json:"url"`
}

// GetSubAccountDepositAddressService -- GET /api/v2/spot/wallet/subaccount-deposit-address (signed)
//
// Returns the deposit address for a sub-account's coin on a given chain.
type GetSubAccountDepositAddressService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSubAccountDepositAddressService(subUid, coin string) *GetSubAccountDepositAddressService {
	return &GetSubAccountDepositAddressService{c: c, params: map[string]string{
		"subUid": subUid,
		"coin":   coin,
	}}
}

// SetChain selects the chain (e.g. trc20); defaults to the coin's primary chain.
func (s *GetSubAccountDepositAddressService) SetChain(chain string) *GetSubAccountDepositAddressService {
	s.params["chain"] = chain
	return s
}

// SetSize sets the Bitcoin Lightning Network withdrawal amount (0.000001 - 0.01).
func (s *GetSubAccountDepositAddressService) SetSize(size decimal.Decimal) *GetSubAccountDepositAddressService {
	s.params["size"] = size.String()
	return s
}

func (s *GetSubAccountDepositAddressService) Do(ctx context.Context) (*DepositAddress, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/subaccount-deposit-address", s.params).WithSign()
	return request.Do[DepositAddress](req)
}

// GetDepositRecordsService -- GET /api/v2/spot/wallet/deposit-records (signed)
//
// Returns the caller's deposit history within a time window.
type GetDepositRecordsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetDepositRecordsService(startTime, endTime int64) *GetDepositRecordsService {
	return &GetDepositRecordsService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime, 10),
		"endTime":   strconv.FormatInt(endTime, 10),
	}}
}

// SetCoin filters by coin.
func (s *GetDepositRecordsService) SetCoin(coin string) *GetDepositRecordsService {
	s.params["coin"] = coin
	return s
}

// SetOrderId filters by a specific order ID.
func (s *GetDepositRecordsService) SetOrderId(orderId string) *GetDepositRecordsService {
	s.params["orderId"] = orderId
	return s
}

// SetIdLessThan pages backwards: returns records with an ID older than this one.
func (s *GetDepositRecordsService) SetIdLessThan(idLessThan string) *GetDepositRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the page size (default 20, max 100).
func (s *GetDepositRecordsService) SetLimit(limit int) *GetDepositRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetDepositRecordsService) Do(ctx context.Context) ([]DepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/deposit-records", s.params).WithSign()
	resp, err := request.Do[[]DepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// DepositRecord is a single deposit. It is the union of the main-account and
// sub-account deposit-record shapes (clientOid/tag/confirm appear on the
// sub-account variant).
type DepositRecord struct {
	OrderId     string          `json:"orderId"`
	TradeId     string          `json:"tradeId"`
	Coin        string          `json:"coin"`
	ClientOid   string          `json:"clientOid"`
	Type        string          `json:"type"` // fixed: deposit
	Size        decimal.Decimal `json:"size"`
	Status      RecordStatus    `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Chain       string          `json:"chain"`
	Confirm     string          `json:"confirm"`
	Dest        RecordDest      `json:"dest"`
	Tag         string          `json:"tag"`
	CTime       time.Time       `json:"cTime"`
	UTime       time.Time       `json:"uTime"`
}

// GetSubAccountDepositRecordsService -- GET /api/v2/spot/wallet/subaccount-deposit-records (signed)
//
// Returns a sub-account's deposit history.
type GetSubAccountDepositRecordsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetSubAccountDepositRecordsService(subUid string) *GetSubAccountDepositRecordsService {
	return &GetSubAccountDepositRecordsService{c: c, params: map[string]string{"subUid": subUid}}
}

// SetCoin filters by coin.
func (s *GetSubAccountDepositRecordsService) SetCoin(coin string) *GetSubAccountDepositRecordsService {
	s.params["coin"] = coin
	return s
}

// SetStartTime filters records at or after t.
func (s *GetSubAccountDepositRecordsService) SetStartTime(t time.Time) *GetSubAccountDepositRecordsService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t.
func (s *GetSubAccountDepositRecordsService) SetEndTime(t time.Time) *GetSubAccountDepositRecordsService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetIdLessThan pages backwards: returns records with an orderId older than this.
func (s *GetSubAccountDepositRecordsService) SetIdLessThan(idLessThan string) *GetSubAccountDepositRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the page size (default 20, max 100).
func (s *GetSubAccountDepositRecordsService) SetLimit(limit int) *GetSubAccountDepositRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetSubAccountDepositRecordsService) Do(ctx context.Context) ([]DepositRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/subaccount-deposit-records", s.params).WithSign()
	resp, err := request.Do[[]DepositRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// GetWithdrawalRecordsService -- GET /api/v2/spot/wallet/withdrawal-records (signed)
//
// Returns the caller's withdrawal history within a time window.
type GetWithdrawalRecordsService struct {
	c      *SpotClient
	params map[string]string
}

func (c *SpotClient) NewGetWithdrawalRecordsService(startTime, endTime int64) *GetWithdrawalRecordsService {
	return &GetWithdrawalRecordsService{c: c, params: map[string]string{
		"startTime": strconv.FormatInt(startTime, 10),
		"endTime":   strconv.FormatInt(endTime, 10),
	}}
}

// SetCoin filters by coin.
func (s *GetWithdrawalRecordsService) SetCoin(coin string) *GetWithdrawalRecordsService {
	s.params["coin"] = coin
	return s
}

// SetClientOid filters by a user-defined order identifier.
func (s *GetWithdrawalRecordsService) SetClientOid(clientOid string) *GetWithdrawalRecordsService {
	s.params["clientOid"] = clientOid
	return s
}

// SetOrderId filters by a specific order ID.
func (s *GetWithdrawalRecordsService) SetOrderId(orderId string) *GetWithdrawalRecordsService {
	s.params["orderId"] = orderId
	return s
}

// SetIdLessThan pages backwards: returns records with an ID older than this one.
func (s *GetWithdrawalRecordsService) SetIdLessThan(idLessThan string) *GetWithdrawalRecordsService {
	s.params["idLessThan"] = idLessThan
	return s
}

// SetLimit caps the page size (default 20, max 100).
func (s *GetWithdrawalRecordsService) SetLimit(limit int) *GetWithdrawalRecordsService {
	s.params["limit"] = strconv.Itoa(limit)
	return s
}

func (s *GetWithdrawalRecordsService) Do(ctx context.Context) ([]WithdrawalRecord, error) {
	req := request.Get(ctx, s.c, "/api/v2/spot/wallet/withdrawal-records", s.params).WithSign()
	resp, err := request.Do[[]WithdrawalRecord](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// WithdrawalRecord is a single withdrawal.
type WithdrawalRecord struct {
	OrderId     string          `json:"orderId"`
	TradeId     string          `json:"tradeId"`
	Coin        string          `json:"coin"`
	ClientOid   string          `json:"clientOid"`
	Type        string          `json:"type"` // fixed: withdraw
	Dest        RecordDest      `json:"dest"`
	Size        decimal.Decimal `json:"size"`
	Fee         decimal.Decimal `json:"fee"`
	Status      RecordStatus    `json:"status"`
	FromAddress string          `json:"fromAddress"`
	ToAddress   string          `json:"toAddress"`
	Chain       string          `json:"chain"`
	Confirm     string          `json:"confirm"`
	Tag         string          `json:"tag"`
	CTime       time.Time       `json:"cTime"`
	UTime       time.Time       `json:"uTime"`
}
