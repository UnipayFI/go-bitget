package uta

import (
	"context"
	"strconv"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CreateBrokerSubService -- POST /api/v3/broker/create-sub (ND Broker master)
//
// Creates a new broker sub-account under the master account. Only a master
// account with the ND Broker user type may call this endpoint.
type CreateBrokerSubService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateBrokerSubService(subaccountName, label string) *CreateBrokerSubService {
	return &CreateBrokerSubService{c: c, body: map[string]any{
		"subaccountName": subaccountName,
		"label":          label,
	}}
}

func (s *CreateBrokerSubService) Do(ctx context.Context) (*BrokerSubAccount, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/create-sub", s.body).WithSign()
	return request.Do[BrokerSubAccount](req)
}

// BrokerSubAccount describes a broker sub-account. The status is "normal" or
// "freeze"; permList entries are "withdraw", "transfer", "spot_trade",
// "contract_trade", "margin_trade", "deposit".
type BrokerSubAccount struct {
	SubUID          string    `json:"subUid"`
	SubaccountName  string    `json:"subaccountName"`
	SubaccountEmail string    `json:"subaccountEmail"`
	Status          string    `json:"status"`
	PermList        []string  `json:"permList"`
	Label           string    `json:"label"`
	Language        string    `json:"language"`
	CTime           time.Time `json:"cTime,format:unixmilli"`
	UTime           time.Time `json:"uTime,format:unixmilli"`
}

// GetBrokerSubListService -- GET /api/v3/broker/sub-list (ND Broker master)
//
// Returns the broker sub-account list, paginated by cursor (pass the last
// subUid of the previous page).
type GetBrokerSubListService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerSubListService() *GetBrokerSubListService {
	return &GetBrokerSubListService{c: c, params: map[string]string{}}
}

func (s *GetBrokerSubListService) SetLimit(limit string) *GetBrokerSubListService {
	s.params["limit"] = limit
	return s
}

// SetCursor sets the pagination cursor (the last subUid of the previous page).
func (s *GetBrokerSubListService) SetCursor(cursor string) *GetBrokerSubListService {
	s.params["cursor"] = cursor
	return s
}

// SetStatus filters by account status ("normal" or "freeze").
func (s *GetBrokerSubListService) SetStatus(status string) *GetBrokerSubListService {
	s.params["status"] = status
	return s
}

func (s *GetBrokerSubListService) Do(ctx context.Context) (*BrokerSubList, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/sub-list", s.params).WithSign()
	return request.Do[BrokerSubList](req)
}

type BrokerSubList struct {
	SubList []BrokerSubAccount `json:"subList"`
}

// ModifyBrokerSubService -- POST /api/v3/broker/modify-sub (ND Broker master)
//
// Updates a broker sub-account's status and/or permission list.
type ModifyBrokerSubService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyBrokerSubService(subUid string) *ModifyBrokerSubService {
	return &ModifyBrokerSubService{c: c, body: map[string]any{
		"subUid": subUid,
	}}
}

// SetStatus sets the account status ("normal" or "freeze").
func (s *ModifyBrokerSubService) SetStatus(status string) *ModifyBrokerSubService {
	s.body["status"] = status
	return s
}

// SetPermList sets the permission list ("withdraw", "transfer", "spot_trade",
// "contract_trade", "margin_trade", "deposit").
func (s *ModifyBrokerSubService) SetPermList(permList []string) *ModifyBrokerSubService {
	s.body["permList"] = permList
	return s
}

func (s *ModifyBrokerSubService) Do(ctx context.Context) (*BrokerSubAccount, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/modify-sub", s.body).WithSign()
	return request.Do[BrokerSubAccount](req)
}

// BrokerSubWithdrawalService -- POST /api/v3/broker/sub-withdrawal (ND Broker master)
//
// Withdraws funds from a broker sub-account, either on-chain or via internal
// transfer (dest is "on_chain" or "internal_transfer"; for internal transfers
// the address is the recipient UID).
type BrokerSubWithdrawalService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewBrokerSubWithdrawalService(subUid, coin, dest, address string, amount decimal.Decimal) *BrokerSubWithdrawalService {
	return &BrokerSubWithdrawalService{c: c, body: map[string]any{
		"subUid":  subUid,
		"coin":    coin,
		"dest":    dest,
		"address": address,
		"amount":  amount.String(),
	}}
}

// SetChain sets the chain name (defaults to the coin's main chain when omitted).
func (s *BrokerSubWithdrawalService) SetChain(chain string) *BrokerSubWithdrawalService {
	s.body["chain"] = chain
	return s
}

// SetTag sets the tag for chains that require one (e.g. EOS memo, TON comment).
func (s *BrokerSubWithdrawalService) SetTag(tag string) *BrokerSubWithdrawalService {
	s.body["tag"] = tag
	return s
}

func (s *BrokerSubWithdrawalService) SetClientOrderID(clientOid string) *BrokerSubWithdrawalService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *BrokerSubWithdrawalService) Do(ctx context.Context) (*BrokerSubWithdrawal, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/sub-withdrawal", s.body).WithSign()
	return request.Do[BrokerSubWithdrawal](req)
}

type BrokerSubWithdrawal struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOid"`
}

// BrokerSubDepositAddressService -- POST /api/v3/broker/sub-deposit-address (ND Broker master)
//
// Returns the deposit address of a coin for a broker sub-account.
type BrokerSubDepositAddressService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewBrokerSubDepositAddressService(subUid, coin string) *BrokerSubDepositAddressService {
	return &BrokerSubDepositAddressService{c: c, body: map[string]any{
		"subUid": subUid,
		"coin":   coin,
	}}
}

// SetChain sets the chain name (defaults to the coin's primary chain when omitted).
func (s *BrokerSubDepositAddressService) SetChain(chain string) *BrokerSubDepositAddressService {
	s.body["chain"] = chain
	return s
}

func (s *BrokerSubDepositAddressService) Do(ctx context.Context) (*BrokerSubDepositAddress, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/sub-deposit-address", s.body).WithSign()
	return request.Do[BrokerSubDepositAddress](req)
}

type BrokerSubDepositAddress struct {
	SubUID  string    `json:"subUid"`
	Coin    string    `json:"coin"`
	Address string    `json:"address"`
	Chain   string    `json:"chain"`
	Tag     string    `json:"tag"`
	URL     string    `json:"url"`
	CTime   time.Time `json:"cTime,format:unixmilli"`
}

// GetAllBrokerSubDepositWithdrawalService -- GET /api/v3/broker/all-sub-deposit-withdrawal (ND Broker master)
//
// Returns the deposit and withdrawal records across all broker sub-accounts,
// paginated by cursor (endId) and bounded to a 7-day range.
type GetAllBrokerSubDepositWithdrawalService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetAllBrokerSubDepositWithdrawalService() *GetAllBrokerSubDepositWithdrawalService {
	return &GetAllBrokerSubDepositWithdrawalService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t (defaults to yesterday UTC when
// omitted; max 7-day range).
func (s *GetAllBrokerSubDepositWithdrawalService) SetStartTime(t time.Time) *GetAllBrokerSubDepositWithdrawalService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 7-day range from startTime).
func (s *GetAllBrokerSubDepositWithdrawalService) SetEndTime(t time.Time) *GetAllBrokerSubDepositWithdrawalService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetAllBrokerSubDepositWithdrawalService) SetLimit(limit string) *GetAllBrokerSubDepositWithdrawalService {
	s.params["limit"] = limit
	return s
}

func (s *GetAllBrokerSubDepositWithdrawalService) SetCursor(cursor string) *GetAllBrokerSubDepositWithdrawalService {
	s.params["cursor"] = cursor
	return s
}

// SetStatus filters by record status ("pending", "fail", or "success").
func (s *GetAllBrokerSubDepositWithdrawalService) SetStatus(status string) *GetAllBrokerSubDepositWithdrawalService {
	s.params["status"] = status
	return s
}

func (s *GetAllBrokerSubDepositWithdrawalService) Do(ctx context.Context) (*BrokerSubDepositWithdrawals, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/all-sub-deposit-withdrawal", s.params).WithSign()
	return request.Do[BrokerSubDepositWithdrawals](req)
}

type BrokerSubDepositWithdrawals struct {
	List  []BrokerSubDepositWithdrawal `json:"list"`
	EndID string                       `json:"endId"`
}

// BrokerSubDepositWithdrawal is a single deposit or withdrawal record. The type
// is "deposit" or "withdrawal"; subType is "onchain", "internal", or "fast";
// status is "pending", "fail", or "success".
type BrokerSubDepositWithdrawal struct {
	UID     string          `json:"uid"`
	TxID    string          `json:"txId"`
	Type    string          `json:"type"`
	SubType string          `json:"subType"`
	Coin    string          `json:"coin"`
	Amount  decimal.Decimal `json:"amount"`
	Status  string          `json:"status"`
	Ts      time.Time       `json:"ts,format:unixmilli"`
}

// GetBrokerCommissionService -- GET /api/v3/broker/commission (ND Broker master)
//
// Returns the broker's commission records, paginated by page number and bounded
// to a 30-day range.
type GetBrokerCommissionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerCommissionService() *GetBrokerCommissionService {
	return &GetBrokerCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t (defaults to yesterday UTC when
// omitted; max 30-day range).
func (s *GetBrokerCommissionService) SetStartTime(t time.Time) *GetBrokerCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (max 30-day range from startTime).
func (s *GetBrokerCommissionService) SetEndTime(t time.Time) *GetBrokerCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBrokerCommissionService) SetPageSize(pageSize string) *GetBrokerCommissionService {
	s.params["pageSize"] = pageSize
	return s
}

func (s *GetBrokerCommissionService) SetPageNo(pageNo string) *GetBrokerCommissionService {
	s.params["pageNo"] = pageNo
	return s
}

// SetBizType filters by business type ("spot" or "futures").
func (s *GetBrokerCommissionService) SetBizType(bizType string) *GetBrokerCommissionService {
	s.params["bizType"] = bizType
	return s
}

// SetSubBizType filters by sub business type ("spot_trade", "spot_margin",
// "usdt_futures", "usdc_futures", "coin_futures").
func (s *GetBrokerCommissionService) SetSubBizType(subBizType string) *GetBrokerCommissionService {
	s.params["subBizType"] = subBizType
	return s
}

func (s *GetBrokerCommissionService) Do(ctx context.Context) ([]BrokerCommission, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/commission", s.params).WithSign()
	resp, err := request.Do[[]BrokerCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type BrokerCommission struct {
	UID             string          `json:"uid"`
	Coin            string          `json:"coin"`
	Symbol          string          `json:"symbol"`
	DealtAmount     decimal.Decimal `json:"dealtAmount"`
	TotalFee        decimal.Decimal `json:"totalFee"`
	DeductedFee     decimal.Decimal `json:"deductedFee"`
	PaidFee         decimal.Decimal `json:"paidFee"`
	MarkUpFee       decimal.Decimal `json:"markUpFee"`
	TotalCommission decimal.Decimal `json:"totalCommission"`
}

// CreateBrokerSubApikeyService -- POST /api/v3/broker/create-sub-apikey (ND Broker master)
//
// Creates an API key for a broker sub-account. permType is "read_write" or
// "read_only"; permList entries are "uta_trade", "uta_mgt", "withdraw"
// ("withdraw" requires read_write).
type CreateBrokerSubAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateBrokerSubAPIKeyService(subUid, passphrase, label string, ipList []string, permType string, permList []string) *CreateBrokerSubAPIKeyService {
	return &CreateBrokerSubAPIKeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"label":      label,
		"ipList":     ipList,
		"permType":   permType,
		"permList":   permList,
	}}
}

func (s *CreateBrokerSubAPIKeyService) Do(ctx context.Context) (*BrokerSubAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/create-sub-apikey", s.body).WithSign()
	return request.Do[BrokerSubAPIKey](req)
}

// BrokerSubApikey describes a broker sub-account API key. secretKey is only
// returned on creation. permType is "read_write" or "read_only"; permList
// entries are "uta_trade", "uta_mgt", "withdraw".
type BrokerSubAPIKey struct {
	SubUID    string   `json:"subUid"`
	Label     string   `json:"label"`
	APIKey    string   `json:"apiKey"`
	SecretKey string   `json:"secretKey"`
	PermType  string   `json:"permType"`
	PermList  []string `json:"permList"`
	IPList    []string `json:"ipList"`
}

// ModifyBrokerSubApikeyService -- POST /api/v3/broker/modify-sub-apikey (ND Broker master)
//
// Updates a broker sub-account's API key (label, IP whitelist, permission type
// and permission list).
type ModifyBrokerSubAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyBrokerSubAPIKeyService(subUid, passphrase, apiKey string) *ModifyBrokerSubAPIKeyService {
	return &ModifyBrokerSubAPIKeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"apiKey":     apiKey,
	}}
}

func (s *ModifyBrokerSubAPIKeyService) SetLabel(label string) *ModifyBrokerSubAPIKeyService {
	s.body["label"] = label
	return s
}

// SetIpList sets the IP whitelist (max 30 entries).
func (s *ModifyBrokerSubAPIKeyService) SetIPList(ipList []string) *ModifyBrokerSubAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

// SetPermType sets the permission type ("read_write" or "read_only").
func (s *ModifyBrokerSubAPIKeyService) SetPermType(permType string) *ModifyBrokerSubAPIKeyService {
	s.body["permType"] = permType
	return s
}

// SetPermList sets the permission list ("uta_trade", "uta_mgt", "withdraw";
// "withdraw" requires read_write permType).
func (s *ModifyBrokerSubAPIKeyService) SetPermList(permList []string) *ModifyBrokerSubAPIKeyService {
	s.body["permList"] = permList
	return s
}

func (s *ModifyBrokerSubAPIKeyService) Do(ctx context.Context) (*BrokerSubAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/modify-sub-apikey", s.body).WithSign()
	return request.Do[BrokerSubAPIKey](req)
}

// DeleteBrokerSubApikeyService -- POST /api/v3/broker/delete-sub-apikey (ND Broker master)
//
// Deletes an API key of a broker sub-account. The reply data is null on success.
type DeleteBrokerSubAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewDeleteBrokerSubAPIKeyService(subUid, apiKey string) *DeleteBrokerSubAPIKeyService {
	return &DeleteBrokerSubAPIKeyService{c: c, body: map[string]any{
		"subUid": subUid,
		"apiKey": apiKey,
	}}
}

func (s *DeleteBrokerSubAPIKeyService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/delete-sub-apikey", s.body).WithSign()
	return request.Do[any](req)
}

// GetBrokerSubApikeyService -- GET /api/v3/broker/query-sub-apikey (ND Broker master)
//
// Returns the API key of a broker sub-account.
type GetBrokerSubAPIKeyService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerSubAPIKeyService(subUid string) *GetBrokerSubAPIKeyService {
	return &GetBrokerSubAPIKeyService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetBrokerSubAPIKeyService) Do(ctx context.Context) (*BrokerSubAPIKey, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/query-sub-apikey", s.params).WithSign()
	return request.Do[BrokerSubAPIKey](req)
}

// ModifyBrokerSubEmailService -- POST /api/v3/broker/modify-subaccount-email (broker)
//
// Binds or changes the email of a broker sub-account. The response data is the
// literal "success".
type ModifyBrokerSubEmailService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyBrokerSubEmailService(subUid, subaccountEmail string) *ModifyBrokerSubEmailService {
	return &ModifyBrokerSubEmailService{c: c, body: map[string]any{
		"subUid":          subUid,
		"subaccountEmail": subaccountEmail,
	}}
}

func (s *ModifyBrokerSubEmailService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/modify-subaccount-email", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetBrokerSubEmailService -- GET /api/v3/broker/subaccount-email (broker)
//
// Returns the email bound to a broker sub-account.
type GetBrokerSubEmailService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerSubEmailService(subUid string) *GetBrokerSubEmailService {
	return &GetBrokerSubEmailService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetBrokerSubEmailService) Do(ctx context.Context) (*BrokerSubEmail, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/subaccount-email", s.params).WithSign()
	return request.Do[BrokerSubEmail](req)
}

// BrokerSubEmail is a sub-account's email binding. SubaccountName is the
// generated email-format username (subaccountName as the prefix).
type BrokerSubEmail struct {
	SubUID          string    `json:"subUid"`
	SubaccountName  string    `json:"subaccountName"`
	SubaccountEmail string    `json:"subaccountEmail"`
	CTime           time.Time `json:"cTime,format:unixmilli"`
	UTime           time.Time `json:"uTime,format:unixmilli"`
}

// GetBrokerTradeVolumeService -- GET /api/v3/broker/trade-volume (broker)
//
// Returns the trading volume of broker sub-accounts. Access must be requested
// from your BD or RM. With no time range it returns yesterday (00:00-23:59
// UTC+0).
type GetBrokerTradeVolumeService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerTradeVolumeService() *GetBrokerTradeVolumeService {
	return &GetBrokerTradeVolumeService{c: c, params: map[string]string{}}
}

// SetStartTime filters records at or after t. startTime and endTime must be set
// together and span at most 30 days.
func (s *GetBrokerTradeVolumeService) SetStartTime(t time.Time) *GetBrokerTradeVolumeService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime filters records at or before t (see SetStartTime).
func (s *GetBrokerTradeVolumeService) SetEndTime(t time.Time) *GetBrokerTradeVolumeService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetPageSize sets the page size (default 100, max 1000).
func (s *GetBrokerTradeVolumeService) SetPageSize(pageSize string) *GetBrokerTradeVolumeService {
	s.params["pageSize"] = pageSize
	return s
}

// SetPageNo sets the page number (default 1).
func (s *GetBrokerTradeVolumeService) SetPageNo(pageNo string) *GetBrokerTradeVolumeService {
	s.params["pageNo"] = pageNo
	return s
}

func (s *GetBrokerTradeVolumeService) Do(ctx context.Context) ([]BrokerTradeVolume, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/trade-volume", s.params).WithSign()
	resp, err := request.Do[[]BrokerTradeVolume](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerTradeVolume is one sub-account's trading volume (spot, futures and
// their total).
type BrokerTradeVolume struct {
	UID          string          `json:"uid"`
	Volume       decimal.Decimal `json:"volume"`
	SpotVolume   decimal.Decimal `json:"spotVolume"`
	FutureVolume decimal.Decimal `json:"futureVolume"`
}

// GetBrokerTotalCommissionService -- GET /api/v3/broker/total-commission (broker)
//
// Returns the broker's daily commission, trading volume and active-trader
// totals, split into spot and futures. Data is daily (UTC+8), updated T+1, kept
// for 365 days and available from 2025-06-01. With no time range it returns
// yesterday (00:00-23:59 UTC+8).
type GetBrokerTotalCommissionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerTotalCommissionService() *GetBrokerTotalCommissionService {
	return &GetBrokerTotalCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime sets the window start. startTime and endTime must be set
// together and span at most 180 days.
func (s *GetBrokerTotalCommissionService) SetStartTime(t time.Time) *GetBrokerTotalCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the window end (see SetStartTime).
func (s *GetBrokerTotalCommissionService) SetEndTime(t time.Time) *GetBrokerTotalCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

func (s *GetBrokerTotalCommissionService) Do(ctx context.Context) ([]BrokerTotalCommission, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/total-commission", s.params).WithSign()
	resp, err := request.Do[[]BrokerTotalCommission](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// BrokerTotalCommission is one day's broker totals; amounts are in USDT. Date
// stays a string because the doc gives no unit and its example
// ("17609760000000", 14 digits) is not a valid millisecond timestamp.
type BrokerTotalCommission struct {
	Date               string                       `json:"date"`
	TotalTradingVolume decimal.Decimal              `json:"totalTradingVolume"`
	TotalActiveTraders string                       `json:"totalActiveTraders"`
	TotalCommission    decimal.Decimal              `json:"totalCommission"`
	Spot               BrokerTotalCommissionSpot    `json:"spot"`
	Futures            BrokerTotalCommissionFutures `json:"futures"`
}

// BrokerTotalCommissionSpot is the spot part of a daily total. The pure
// trading fee is the fee after negative-maker and similar deductions, and is
// the commission base.
type BrokerTotalCommissionSpot struct {
	SpotTradingVolume  decimal.Decimal `json:"spotTradingVolume"`
	SpotTradingFee     decimal.Decimal `json:"spotTradingFee"`
	SpotPureTradingFee decimal.Decimal `json:"spotPureTradingFee"`
	SpotCommission     decimal.Decimal `json:"spotCommission"`
}

// BrokerTotalCommissionFutures is the futures part of a daily total (see
// BrokerTotalCommissionSpot).
type BrokerTotalCommissionFutures struct {
	FuturesTradingVolume  decimal.Decimal `json:"futuresTradingVolume"`
	FuturesTradingFee     decimal.Decimal `json:"futuresTradingFee"`
	FuturesPureTradingFee decimal.Decimal `json:"futuresPureTradingFee"`
	FuturesCommission     decimal.Decimal `json:"futuresCommission"`
}

// GetBrokerOrderCommissionService -- GET /api/v3/broker/order-commission (broker)
//
// Returns per-fill broker commission (rebate) records, paginated by idLessThan
// (endId). Only fills tagged with the broker channel ID appear. Data is daily
// (UTC+8), updated T+1, kept for 30 days and available from 2025-09-01. With no
// time range it returns yesterday (00:00-23:59 UTC+8).
type GetBrokerOrderCommissionService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerOrderCommissionService() *GetBrokerOrderCommissionService {
	return &GetBrokerOrderCommissionService{c: c, params: map[string]string{}}
}

// SetStartTime sets the window start. startTime and endTime must be set
// together; the endpoint notes say the span is at most 30 days, the parameter
// table says 180.
func (s *GetBrokerOrderCommissionService) SetStartTime(t time.Time) *GetBrokerOrderCommissionService {
	s.params["startTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetEndTime sets the window end (see SetStartTime).
func (s *GetBrokerOrderCommissionService) SetEndTime(t time.Time) *GetBrokerOrderCommissionService {
	s.params["endTime"] = strconv.FormatInt(t.UnixMilli(), 10)
	return s
}

// SetLimit sets the page size (default 500, max 500).
func (s *GetBrokerOrderCommissionService) SetLimit(limit string) *GetBrokerOrderCommissionService {
	s.params["limit"] = limit
	return s
}

func (s *GetBrokerOrderCommissionService) SetUID(uid string) *GetBrokerOrderCommissionService {
	s.params["uid"] = uid
	return s
}

// SetOrderID filters to one order (sent as the all-lowercase "orderid").
func (s *GetBrokerOrderCommissionService) SetOrderID(orderID string) *GetBrokerOrderCommissionService {
	s.params["orderid"] = orderID
	return s
}

// SetIDLessThan sets the pagination cursor (endId from the previous response;
// omit on the first call).
func (s *GetBrokerOrderCommissionService) SetIDLessThan(idLessThan string) *GetBrokerOrderCommissionService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetBrokerOrderCommissionService) Do(ctx context.Context) (*BrokerOrderCommissions, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/order-commission", s.params).WithSign()
	return request.Do[BrokerOrderCommissions](req)
}

// BrokerOrderCommissions is a page of per-fill commission records; the list key
// is the all-lowercase "commissionlist" on the wire.
type BrokerOrderCommissions struct {
	CommissionList []BrokerOrderCommission `json:"commissionlist"`
	EndID          string                  `json:"endId"`
}

// BrokerOrderCommission is one fill's broker commission; Volume, Fee, PureFee
// and RebateAmount are in USDT at full precision.
type BrokerOrderCommission struct {
	FillID        string          `json:"fillId"`
	OrderID       string          `json:"orderId"`
	Ts            time.Time       `json:"ts,format:unixmilli"`
	ClientOrderID string          `json:"clientOid"`
	BizType       string          `json:"bizType"`    // spot, futures
	SubBizType    string          `json:"subBizType"` // spot_trade, spot_margin, usdt_futures, coin_futures, usdc_futures
	Symbol        string          `json:"symbol"`
	Volume        decimal.Decimal `json:"volume"`
	Fee           decimal.Decimal `json:"fee"`
	PureFee       decimal.Decimal `json:"pureFee"`
	RebateAmount  decimal.Decimal `json:"rebateAmount"`
}

// GetBrokerRebateInfoService -- GET /api/v3/broker/rebate-info (broker)
//
// Returns a user's broker rebate relationship, VIP level and the spot/futures
// rebate ratios available to them. Real-time data.
type GetBrokerRebateInfoService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerRebateInfoService(uid string) *GetBrokerRebateInfoService {
	return &GetBrokerRebateInfoService{c: c, params: map[string]string{"uid": uid}}
}

func (s *GetBrokerRebateInfoService) Do(ctx context.Context) (*BrokerRebateInfo, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/rebate-info", s.params).WithSign()
	return request.Do[BrokerRebateInfo](req)
}

// BrokerRebateInfo is a user's broker rebate configuration. AffiliationType is
// "affiliate" (client channel) or "official" (website channel).
type BrokerRebateInfo struct {
	AffiliationType          string          `json:"affiliationType"`
	UserLevel                string          `json:"userLevel"` // VIP0..VIP7, PRO1..PRO6
	ClientSpotRebateRatio    decimal.Decimal `json:"clientSpotRebateRatio"`
	ClientFuturesRebateRatio decimal.Decimal `json:"clientFuturesRebateRatio"`
}
