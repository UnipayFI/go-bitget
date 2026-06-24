package broker

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
)

// ApikeyPermType is the authorization scope level of a sub-account API key.
type APIKeyPermType string

const (
	APIKeyPermTypeReadAndWrite APIKeyPermType = "read_and_write"
	APIKeyPermTypeReadonly     APIKeyPermType = "readonly"
)

// ApikeyPerm is a single permission scope granted to a sub-account API key.
type APIKeyPerm string

const (
	APIKeyPermContractOrder    APIKeyPerm = "contract_order"
	APIKeyPermContractPosition APIKeyPerm = "contract_position"
	APIKeyPermSpotTrade        APIKeyPerm = "spot_trade"
	APIKeyPermMarginTrade      APIKeyPerm = "margin_trade"
	APIKeyPermCopytradingTrade APIKeyPerm = "copytrading_trade"
	APIKeyPermWalletTransfer   APIKeyPerm = "wallet_transfer"
)

// CreateSubaccountApikeyService -- POST /api/v2/broker/manage/create-subaccount-apikey (private, state-changing)
//
// Creates a new API key for a broker sub-account.
type CreateSubaccountAPIKeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewCreateSubaccountAPIKeyService(subUid, passphrase, label string, permType APIKeyPermType, permList []APIKeyPerm) *CreateSubaccountAPIKeyService {
	return &CreateSubaccountAPIKeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"passphrase": passphrase,
		"label":      label,
		"permType":   string(permType),
		"permList":   permList,
	}}
}

// SetIpList sets the IP whitelist (max 30 entries). Required unless permType is readonly.
func (s *CreateSubaccountAPIKeyService) SetIPList(ipList []string) *CreateSubaccountAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *CreateSubaccountAPIKeyService) Do(ctx context.Context) (*SubaccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/manage/create-subaccount-apikey").SetBody(s.body).WithSign()
	return request.Do[SubaccountAPIKey](req)
}

// SubaccountApikey is a sub-account API key record. secretKey is populated only
// on creation and on the apikey-list query; the modify response omits it.
type SubaccountAPIKey struct {
	SubUID    string   `json:"subUid"`
	APIKey    string   `json:"apiKey"`
	SecretKey string   `json:"secretKey"`
	Label     string   `json:"label"`
	IPList    []string `json:"ipList"`
	PermType  string   `json:"permType"`
	PermList  []string `json:"permList"`
}

// ModifySubaccountApikeyService -- POST /api/v2/broker/manage/modify-subaccount-apikey (private, state-changing)
//
// Modifies an existing API key (permissions, IP whitelist, label) for a broker
// sub-account. The override semantics mean permType and permList fully replace
// the previous values.
type ModifySubaccountAPIKeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewModifySubaccountAPIKeyService(subUid, apiKey, passphrase string, permType APIKeyPermType, permList []APIKeyPerm) *ModifySubaccountAPIKeyService {
	return &ModifySubaccountAPIKeyService{c: c, body: map[string]any{
		"subUid":     subUid,
		"apiKey":     apiKey,
		"passphrase": passphrase,
		"permType":   string(permType),
		"permList":   permList,
	}}
}

// SetLabel sets the API key remark (max 20 characters).
func (s *ModifySubaccountAPIKeyService) SetLabel(label string) *ModifySubaccountAPIKeyService {
	s.body["label"] = label
	return s
}

// SetIpList overrides the IP whitelist (max 30 entries). An empty list means no change.
func (s *ModifySubaccountAPIKeyService) SetIPList(ipList []string) *ModifySubaccountAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *ModifySubaccountAPIKeyService) Do(ctx context.Context) (*SubaccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/broker/manage/modify-subaccount-apikey").SetBody(s.body).WithSign()
	return request.Do[SubaccountAPIKey](req)
}

// DeleteSubaccountApikeyService -- POST /api/v2/broker/manage/delete-subaccount-apikey (private, state-changing)
//
// Deletes an API key from a broker sub-account. The data field is the literal
// string "success" on completion.
type DeleteSubaccountAPIKeyService struct {
	c    *BrokerClient
	body map[string]any
}

func (c *BrokerClient) NewDeleteSubaccountAPIKeyService(subUid, apiKey string) *DeleteSubaccountAPIKeyService {
	return &DeleteSubaccountAPIKeyService{c: c, body: map[string]any{
		"subUid": subUid,
		"apiKey": apiKey,
	}}
}

func (s *DeleteSubaccountAPIKeyService) Do(ctx context.Context) (string, error) {
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
type GetSubaccountAPIKeyListService struct {
	c      *BrokerClient
	params map[string]string
}

func (c *BrokerClient) NewGetSubaccountAPIKeyListService(subUid string) *GetSubaccountAPIKeyListService {
	return &GetSubaccountAPIKeyListService{c: c, params: map[string]string{"subUid": subUid}}
}

func (s *GetSubaccountAPIKeyListService) Do(ctx context.Context) ([]SubaccountAPIKey, error) {
	req := request.Get(ctx, s.c, "/api/v2/broker/manage/subaccount-apikey-list", s.params).WithSign()
	resp, err := request.Do[[]SubaccountAPIKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
