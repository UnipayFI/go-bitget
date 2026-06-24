package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
	"github.com/shopspring/decimal"
)

// CreateSubAccountService -- POST /api/v3/user/create-sub (UTA mgt. read & write)
//
// Creates a virtual sub-account under the main account. The username generates a
// virtual email address (lowercase letters only, max 20 chars). Main account
// only.
type CreateSubAccountService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateSubAccountService(username string) *CreateSubAccountService {
	return &CreateSubAccountService{c: c, body: map[string]any{
		"username": username,
	}}
}

// SetAccountMode sets the sub-account type ("classic" or "unified").
func (s *CreateSubAccountService) SetAccountMode(accountMode string) *CreateSubAccountService {
	s.body["accountMode"] = accountMode
	return s
}

// SetNote sets the sub-account note (max 50 chars).
func (s *CreateSubAccountService) SetNote(note string) *CreateSubAccountService {
	s.body["note"] = note
	return s
}

func (s *CreateSubAccountService) Do(ctx context.Context) (*SubAccount, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/create-sub", s.body).WithSign()
	return request.Do[SubAccount](req)
}

type SubAccount struct {
	Username    string    `json:"username"`
	SubUID      string    `json:"subUid"`
	Status      string    `json:"status"` // normal
	Note        string    `json:"note"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

// FreezeSubAccountService -- POST /api/v3/user/freeze-sub (UTA mgt. read & write)
//
// Freezes or unfreezes a sub-account (operation is "freeze" or "unfreeze"). The
// reply data is an empty object. Main account only.
type FreezeSubAccountService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewFreezeSubAccountService(subUid, operation string) *FreezeSubAccountService {
	return &FreezeSubAccountService{c: c, body: map[string]any{
		"subUid":    subUid,
		"operation": operation,
	}}
}

func (s *FreezeSubAccountService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/freeze-sub", s.body).WithSign()
	return request.Do[any](req)
}

// GetSubAccountUnifiedAssetsService -- GET /api/v3/account/sub-unified-assets (UTA mgt. read)
//
// Returns the unified-account asset holdings of one sub-account, or of all
// sub-accounts when subUid is omitted, paginated by cursor.
type GetSubAccountUnifiedAssetsService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubAccountUnifiedAssetsService() *GetSubAccountUnifiedAssetsService {
	return &GetSubAccountUnifiedAssetsService{c: c, params: map[string]string{}}
}

// SetSubUid filters to a single sub-account (all sub-accounts when omitted).
func (s *GetSubAccountUnifiedAssetsService) SetSubUID(subUid string) *GetSubAccountUnifiedAssetsService {
	s.params["subUid"] = subUid
	return s
}

func (s *GetSubAccountUnifiedAssetsService) SetCursor(cursor string) *GetSubAccountUnifiedAssetsService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit sets the sub-accounts per page (default 10, max 50).
func (s *GetSubAccountUnifiedAssetsService) SetLimit(limit string) *GetSubAccountUnifiedAssetsService {
	s.params["limit"] = limit
	return s
}

func (s *GetSubAccountUnifiedAssetsService) Do(ctx context.Context) ([]SubAccountUnifiedAssets, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/sub-unified-assets", s.params).WithSign()
	resp, err := request.Do[[]SubAccountUnifiedAssets](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

type SubAccountUnifiedAssets struct {
	SubUID string            `json:"subUid"`
	Cursor string            `json:"cursor"`
	Assets []SubAccountAsset `json:"assets"`
}

type SubAccountAsset struct {
	Coin      string          `json:"coin"`
	Equity    decimal.Decimal `json:"equity"`
	USDValue  decimal.Decimal `json:"usdValue"`
	Balance   decimal.Decimal `json:"balance"`
	Debt      decimal.Decimal `json:"debt"`
	Available decimal.Decimal `json:"available"`
	Locked    decimal.Decimal `json:"locked"`
}

// GetSubAccountListService -- GET /api/v3/user/sub-list (UTA mgt. read)
//
// Returns the main account's sub-account list, paginated by cursor. Main account
// only.
type GetSubAccountListService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubAccountListService() *GetSubAccountListService {
	return &GetSubAccountListService{c: c, params: map[string]string{}}
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetSubAccountListService) SetLimit(limit string) *GetSubAccountListService {
	s.params["limit"] = limit
	return s
}

func (s *GetSubAccountListService) SetCursor(cursor string) *GetSubAccountListService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSubAccountListService) Do(ctx context.Context) (*SubAccountList, error) {
	req := request.Get(ctx, s.c, "/api/v3/user/sub-list", s.params).WithSign()
	return request.Do[SubAccountList](req)
}

type SubAccountList struct {
	List    []SubAccountInfo `json:"list"`
	HasNext bool             `json:"hasNext"`
	Cursor  string           `json:"cursor"`
}

type SubAccountInfo struct {
	SubUID      string    `json:"subUid"`
	Username    string    `json:"username"`
	Status      string    `json:"status"`      // normal, freeze
	AccountMode string    `json:"accountMode"` // CLASSIC, UNIFIED
	Type        string    `json:"type"`        // normal, virtual, custodian, agent, other
	Note        string    `json:"note"`
	CreatedTime time.Time `json:"createdTime"`
	UpdatedTime time.Time `json:"updatedTime"`
}

// CreateSubAccountAPIKeyService -- POST /api/v3/user/create-sub-api (UTA mgt. read & write)
//
// Creates an API key for a sub-account. type is "read_write" or "read_only";
// permissions are "uta_mgt"/"uta_trade"; ips is the IPv4 whitelist (max 30).
// Main account only.
type CreateSubAccountAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateSubAccountAPIKeyService(subUid, note, keyType, passphrase string, permissions, ips []string) *CreateSubAccountAPIKeyService {
	return &CreateSubAccountAPIKeyService{c: c, body: map[string]any{
		"subUid":      subUid,
		"note":        note,
		"type":        keyType,
		"passphrase":  passphrase,
		"permissions": permissions,
		"ips":         ips,
	}}
}

func (s *CreateSubAccountAPIKeyService) Do(ctx context.Context) (*SubAccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/create-sub-api", s.body).WithSign()
	return request.Do[SubAccountAPIKey](req)
}

type SubAccountAPIKey struct {
	Note        string   `json:"note"`
	APIKey      string   `json:"apiKey"`
	Secret      string   `json:"secret"`
	Type        string   `json:"type"` // read_write, read_only
	Permissions []string `json:"permissions"`
	Ips         []string `json:"ips"`
}

// UpdateSubAccountAPIKeyService -- POST /api/v3/user/update-sub-api (UTA mgt. read & write)
//
// Modifies a sub-account API key. type is required when permissions is set (and
// vice versa); passing an empty ips array deletes the whitelist. Main account
// only.
type UpdateSubAccountAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewUpdateSubAccountAPIKeyService(apikey, passphrase string) *UpdateSubAccountAPIKeyService {
	return &UpdateSubAccountAPIKeyService{c: c, body: map[string]any{
		"apikey":     apikey,
		"passphrase": passphrase,
	}}
}

// SetType sets the permission type ("read_write" or "read_only"; required when
// permissions is set).
func (s *UpdateSubAccountAPIKeyService) SetType(keyType string) *UpdateSubAccountAPIKeyService {
	s.body["type"] = keyType
	return s
}

// SetPermissions sets the permission values ("uta_mgt"/"uta_trade"; required
// when type is set).
func (s *UpdateSubAccountAPIKeyService) SetPermissions(permissions []string) *UpdateSubAccountAPIKeyService {
	s.body["permissions"] = permissions
	return s
}

// SetIps sets the IPv4 whitelist (max 30; an empty array deletes the whitelist).
func (s *UpdateSubAccountAPIKeyService) SetIps(ips []string) *UpdateSubAccountAPIKeyService {
	s.body["ips"] = ips
	return s
}

func (s *UpdateSubAccountAPIKeyService) Do(ctx context.Context) (*SubAccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/update-sub-api", s.body).WithSign()
	return request.Do[SubAccountAPIKey](req)
}

// DeleteSubAccountAPIKeyService -- POST /api/v3/user/delete-sub-api (UTA mgt. read & write)
//
// Deletes a sub-account API key, taking effect immediately. The reply data is an
// empty object. Main account only.
type DeleteSubAccountAPIKeyService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewDeleteSubAccountAPIKeyService(apiKey string) *DeleteSubAccountAPIKeyService {
	return &DeleteSubAccountAPIKeyService{c: c, body: map[string]any{
		"apiKey": apiKey,
	}}
}

func (s *DeleteSubAccountAPIKeyService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/delete-sub-api", s.body).WithSign()
	return request.Do[any](req)
}

// GetSubAccountAPIKeysService -- GET /api/v3/user/sub-api-list (UTA mgt. read)
//
// Returns the API keys of a sub-account, paginated by cursor. Main account only.
type GetSubAccountAPIKeysService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetSubAccountAPIKeysService(subUid string) *GetSubAccountAPIKeysService {
	return &GetSubAccountAPIKeysService{c: c, params: map[string]string{"subUid": subUid}}
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetSubAccountAPIKeysService) SetLimit(limit string) *GetSubAccountAPIKeysService {
	s.params["limit"] = limit
	return s
}

func (s *GetSubAccountAPIKeysService) SetCursor(cursor string) *GetSubAccountAPIKeysService {
	s.params["cursor"] = cursor
	return s
}

func (s *GetSubAccountAPIKeysService) Do(ctx context.Context) (*SubAccountAPIKeys, error) {
	req := request.Get(ctx, s.c, "/api/v3/user/sub-api-list", s.params).WithSign()
	return request.Do[SubAccountAPIKeys](req)
}

type SubAccountAPIKeys struct {
	Cursor  string                `json:"cursor"`
	HasNext bool                  `json:"hasNext"`
	List    []SubAccountAPIKeyRow `json:"list"`
}

type SubAccountAPIKeyRow struct {
	APIKey      string    `json:"apiKey"`
	Type        string    `json:"type"` // read_write, read_only
	Note        string    `json:"note"`
	Permissions []string  `json:"permissions"`
	Ips         []string  `json:"ips"`
	Ts          time.Time `json:"ts"`
}

// CreateAgentSubAccountService -- POST /api/v3/user/sub-account/agent-create (UTA mgt. read & write)
//
// Creates an agent sub-account together with its API key. The username generates
// a virtual email address (lowercase letters only, max 20 chars); passphrase is
// 8-32 alphanumeric chars. Main account only.
type CreateAgentSubAccountService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateAgentSubAccountService(username, passphrase string) *CreateAgentSubAccountService {
	return &CreateAgentSubAccountService{c: c, body: map[string]any{
		"username":   username,
		"passphrase": passphrase,
	}}
}

// SetNote sets the sub-account note.
func (s *CreateAgentSubAccountService) SetNote(note string) *CreateAgentSubAccountService {
	s.body["note"] = note
	return s
}

func (s *CreateAgentSubAccountService) Do(ctx context.Context) (*AgentSubAccount, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/sub-account/agent-create", s.body).WithSign()
	return request.Do[AgentSubAccount](req)
}

type AgentSubAccount struct {
	Username    string    `json:"username"`
	SubUID      string    `json:"subUid"`
	APIKey      string    `json:"apiKey"`
	Secret      string    `json:"secret"`
	Note        string    `json:"note"`
	CreatedTime time.Time `json:"createdTime"`
}
