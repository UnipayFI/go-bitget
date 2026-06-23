package broker

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
)

// ApikeyPermType is the authorization scope level of a sub-account API key.
type ApikeyPermType string

const (
	ApikeyPermTypeReadAndWrite ApikeyPermType = "read_and_write"
	ApikeyPermTypeReadonly     ApikeyPermType = "readonly"
)

// ApikeyPerm is a single permission scope granted to a sub-account API key.
type ApikeyPerm string

const (
	ApikeyPermContractOrder    ApikeyPerm = "contract_order"
	ApikeyPermContractPosition ApikeyPerm = "contract_position"
	ApikeyPermSpotTrade        ApikeyPerm = "spot_trade"
	ApikeyPermMarginTrade      ApikeyPerm = "margin_trade"
	ApikeyPermCopytradingTrade ApikeyPerm = "copytrading_trade"
	ApikeyPermWalletTransfer   ApikeyPerm = "wallet_transfer"
)

// CreateSubaccountApikeyService -- POST /api/v2/broker/manage/create-subaccount-apikey (private, state-changing)
//
// Creates a new API key for a broker sub-account.
type CreateSubaccountApikeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewCreateSubaccountApikeyService(subUid, passphrase, label string, permType ApikeyPermType, permList []ApikeyPerm) *CreateSubaccountApikeyService {
	return &CreateSubaccountApikeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"label":      label,
		"permType":   string(permType),
		"permList":   permList,
	}}
}

// SetIpList sets the IP whitelist (max 30 entries). Required unless permType is readonly.
func (s *CreateSubaccountApikeyService) SetIpList(ipList []string) *CreateSubaccountApikeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *CreateSubaccountApikeyService) Do(ctx context.Context) (*SubaccountApikey, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/manage/create-subaccount-apikey").SetBody(s.body).WithSign()
	return request.Do[SubaccountApikey](req)
}

// SubaccountApikey is a sub-account API key record. secretKey is populated only
// on creation and on the apikey-list query; the modify response omits it.
type SubaccountApikey struct {
	SubUid    string   `json:"subUid"`
	ApiKey    string   `json:"apiKey"`
	SecretKey string   `json:"secretKey"`
	Label     string   `json:"label"`
	IpList    []string `json:"ipList"`
	PermType  string   `json:"permType"`
	PermList  []string `json:"permList"`
}

// ModifySubaccountApikeyService -- POST /api/v2/broker/manage/modify-subaccount-apikey (private, state-changing)
//
// Modifies an existing API key (permissions, IP whitelist, label) for a broker
// sub-account. The override semantics mean permType and permList fully replace
// the previous values.
type ModifySubaccountApikeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewModifySubaccountApikeyService(subUid, apiKey, passphrase string, permType ApikeyPermType, permList []ApikeyPerm) *ModifySubaccountApikeyService {
	return &ModifySubaccountApikeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"apiKey":     apiKey,
		"passphrase": passphrase,
		"permType":   string(permType),
		"permList":   permList,
	}}
}

// SetLabel sets the API key remark (max 20 characters).
func (s *ModifySubaccountApikeyService) SetLabel(label string) *ModifySubaccountApikeyService {
	s.body["label"] = label
	return s
}

// SetIpList overrides the IP whitelist (max 30 entries). An empty list means no change.
func (s *ModifySubaccountApikeyService) SetIpList(ipList []string) *ModifySubaccountApikeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *ModifySubaccountApikeyService) Do(ctx context.Context) (*SubaccountApikey, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/manage/modify-subaccount-apikey").SetBody(s.body).WithSign()
	return request.Do[SubaccountApikey](req)
}

// DeleteSubaccountApikeyService -- POST /api/v2/broker/manage/delete-subaccount-apikey (private, state-changing)
//
// Deletes an API key from a broker sub-account. The data field is the literal
// string "success" on completion.
type DeleteSubaccountApikeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewDeleteSubaccountApikeyService(subUid, apiKey string) *DeleteSubaccountApikeyService {
	return &DeleteSubaccountApikeyService{c: c, body: map[string]any{
		"subUid": subUid,
		"apiKey": apiKey,
	}}
}

func (s *DeleteSubaccountApikeyService) Do(ctx context.Context) (string, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/manage/delete-subaccount-apikey", s.body).WithSign()
	resp, err := request.Do[string](req)
	if err != nil {
		return "", err
	}
	return *resp, nil
}

// GetSubaccountApikeyListService -- GET /api/v2/broker/manage/subaccount-apikey-list (private)
//
// Returns the list of API keys for a broker sub-account.
type GetSubaccountApikeyListService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountApikeyListService(subUid string) *GetSubaccountApikeyListService {
	return &GetSubaccountApikeyListService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetSubaccountApikeyListService) Do(ctx context.Context) ([]SubaccountApikey, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/manage/subaccount-apikey-list", s.params).WithSign()
	resp, err := request.Do[[]SubaccountApikey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
