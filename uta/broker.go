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
	SubUid          string    `json:"subUid"`
	SubaccountName  string    `json:"subaccountName"`
	SubaccountEmail string    `json:"subaccountEmail"`
	Status          string    `json:"status"`
	PermList        []string  `json:"permList"`
	Label           string    `json:"label"`
	Language        string    `json:"language"`
	CTime           time.Time `json:"cTime"`
	UTime           time.Time `json:"uTime"`
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

func (s *BrokerSubWithdrawalService) SetClientOid(clientOid string) *BrokerSubWithdrawalService {
	s.body["clientOid"] = clientOid
	return s
}

func (s *BrokerSubWithdrawalService) Do(ctx context.Context) (*BrokerSubWithdrawal, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/sub-withdrawal", s.body).WithSign()
	return request.Do[BrokerSubWithdrawal](req)
}

type BrokerSubWithdrawal struct {
	OrderId   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
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
	SubUid  string    `json:"subUid"`
	Coin    string    `json:"coin"`
	Address string    `json:"address"`
	Chain   string    `json:"chain"`
	Tag     string    `json:"tag"`
	Url     string    `json:"url"`
	CTime   time.Time `json:"cTime"`
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
	EndId string                       `json:"endId"`
}

// BrokerSubDepositWithdrawal is a single deposit or withdrawal record. The type
// is "deposit" or "withdrawal"; subType is "onchain", "internal", or "fast";
// status is "pending", "fail", or "success".
type BrokerSubDepositWithdrawal struct {
	Uid     string          `json:"uid"`
	TxId    string          `json:"txId"`
	Type    string          `json:"type"`
	SubType string          `json:"subType"`
	Coin    string          `json:"coin"`
	Amount  decimal.Decimal `json:"amount"`
	Status  string          `json:"status"`
	Ts      time.Time       `json:"ts"`
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

// CreateBrokerSubApikeyService -- POST /api/v3/broker/create-sub-apikey (ND Broker master)
//
// Creates an API key for a broker sub-account. permType is "read_write" or
// "read_only"; permList entries are "uta_trade", "uta_mgt", "withdraw"
// ("withdraw" requires read_write).
type CreateBrokerSubApikeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateBrokerSubApikeyService(subUid, passphrase, label string, ipList []string, permType string, permList []string) *CreateBrokerSubApikeyService {
	return &CreateBrokerSubApikeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"label":      label,
		"ipList":     ipList,
		"permType":   permType,
		"permList":   permList,
	}}
}

func (s *CreateBrokerSubApikeyService) Do(ctx context.Context) (*BrokerSubApikey, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/create-sub-apikey", s.body).WithSign()
	return request.Do[BrokerSubApikey](req)
}

// BrokerSubApikey describes a broker sub-account API key. secretKey is only
// returned on creation. permType is "read_write" or "read_only"; permList
// entries are "uta_trade", "uta_mgt", "withdraw".
type BrokerSubApikey struct {
	SubUid    string   `json:"subUid"`
	Label     string   `json:"label"`
	ApiKey    string   `json:"apiKey"`
	SecretKey string   `json:"secretKey"`
	PermType  string   `json:"permType"`
	PermList  []string `json:"permList"`
	IpList    []string `json:"ipList"`
}

// ModifyBrokerSubApikeyService -- POST /api/v3/broker/modify-sub-apikey (ND Broker master)
//
// Updates a broker sub-account's API key (label, IP whitelist, permission type
// and permission list).
type ModifyBrokerSubApikeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewModifyBrokerSubApikeyService(subUid, passphrase, apiKey string) *ModifyBrokerSubApikeyService {
	return &ModifyBrokerSubApikeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"apiKey":     apiKey,
	}}
}

func (s *ModifyBrokerSubApikeyService) SetLabel(label string) *ModifyBrokerSubApikeyService {
	s.body["label"] = label
	return s
}

// SetIpList sets the IP whitelist (max 30 entries).
func (s *ModifyBrokerSubApikeyService) SetIpList(ipList []string) *ModifyBrokerSubApikeyService {
	s.body["ipList"] = ipList
	return s
}

// SetPermType sets the permission type ("read_write" or "read_only").
func (s *ModifyBrokerSubApikeyService) SetPermType(permType string) *ModifyBrokerSubApikeyService {
	s.body["permType"] = permType
	return s
}

// SetPermList sets the permission list ("uta_trade", "uta_mgt", "withdraw";
// "withdraw" requires read_write permType).
func (s *ModifyBrokerSubApikeyService) SetPermList(permList []string) *ModifyBrokerSubApikeyService {
	s.body["permList"] = permList
	return s
}

func (s *ModifyBrokerSubApikeyService) Do(ctx context.Context) (*BrokerSubApikey, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/modify-sub-apikey", s.body).WithSign()
	return request.Do[BrokerSubApikey](req)
}

// DeleteBrokerSubApikeyService -- POST /api/v3/broker/delete-sub-apikey (ND Broker master)
//
// Deletes an API key of a broker sub-account. The reply data is null on success.
type DeleteBrokerSubApikeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewDeleteBrokerSubApikeyService(subUid, apiKey string) *DeleteBrokerSubApikeyService {
	return &DeleteBrokerSubApikeyService{c: c, body: map[string]any{
		"subUid": subUid,
		"apiKey": apiKey,
	}}
}

func (s *DeleteBrokerSubApikeyService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/broker/delete-sub-apikey", s.body).WithSign()
	return request.Do[any](req)
}

// GetBrokerSubApikeyService -- GET /api/v3/broker/query-sub-apikey (ND Broker master)
//
// Returns the API key of a broker sub-account.
type GetBrokerSubApikeyService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetBrokerSubApikeyService(subUid string) *GetBrokerSubApikeyService {
	return &GetBrokerSubApikeyService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetBrokerSubApikeyService) Do(ctx context.Context) (*BrokerSubApikey, error) {
	req := request.Get(ctx, s.c, "/api/v3/broker/query-sub-apikey", s.params).WithSign()
	return request.Do[BrokerSubApikey](req)
}
