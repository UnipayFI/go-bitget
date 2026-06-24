package common

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/request"
)

// VirtualSubaccountStatus is the lifecycle state of a virtual sub-account.
type VirtualSubaccountStatus string

const (
	VirtualSubaccountStatusNormal VirtualSubaccountStatus = "normal"
	VirtualSubaccountStatusFreeze VirtualSubaccountStatus = "freeze"
	VirtualSubaccountStatusDel    VirtualSubaccountStatus = "del"
)

// VirtualSubaccountPerm is a permission grantable to a virtual sub-account or
// its ApiKey.
type VirtualSubaccountPerm string

const (
	VirtualSubaccountPermSpotTrade     VirtualSubaccountPerm = "spot_trade"
	VirtualSubaccountPermMarginTrade   VirtualSubaccountPerm = "margin_trade"
	VirtualSubaccountPermContractTrade VirtualSubaccountPerm = "contract_trade"
	VirtualSubaccountPermTransfer      VirtualSubaccountPerm = "transfer"
	VirtualSubaccountPermRead          VirtualSubaccountPerm = "read"
)

// CreateVirtualSubaccountService -- POST /api/v2/user/create-virtual-subaccount (signed; state-changing)
//
// Creates one or more virtual sub-accounts from a list of 8-letter aliases.
type CreateVirtualSubaccountService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewCreateVirtualSubaccountService(subAccountList []string) *CreateVirtualSubaccountService {
	return &CreateVirtualSubaccountService{c: c, body: map[string]any{
		"subAccountList": subAccountList,
	}}
}

func (s *CreateVirtualSubaccountService) Do(ctx context.Context) (*CreateVirtualSubaccountResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/create-virtual-subaccount", s.body).WithSign()
	return request.Do[CreateVirtualSubaccountResult](req)
}

// CreateVirtualSubaccountResult is the create-virtual-subaccount payload,
// splitting created and rejected aliases.
type CreateVirtualSubaccountResult struct {
	FailureList []VirtualSubaccountCreateFailure `json:"failureList"`
	SuccessList []VirtualSubaccountCreateSuccess `json:"successList"`
}

// VirtualSubaccountCreateFailure is an alias that could not be created (already
// exists or the sub-account limit was reached).
type VirtualSubaccountCreateFailure struct {
	SubAccountName string `json:"subaAccountName"`
}

// VirtualSubaccountCreateSuccess is a newly created virtual sub-account.
type VirtualSubaccountCreateSuccess struct {
	SubAccountUID  string                  `json:"subaAccountUid"`
	SubAccountName string                  `json:"subaAccountName"`
	Status         VirtualSubaccountStatus `json:"status"`
	Label          string                  `json:"label"`
	PermList       []VirtualSubaccountPerm `json:"permList"`
	CTime          time.Time               `json:"cTime"`
	UTime          time.Time               `json:"uTime"`
}

// ModifyVirtualSubaccountService -- POST /api/v2/user/modify-virtual-subaccount (signed; state-changing)
//
// Updates a virtual sub-account's permission set and status.
type ModifyVirtualSubaccountService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewModifyVirtualSubaccountService(subAccountUid string, permList []VirtualSubaccountPerm, status VirtualSubaccountStatus) *ModifyVirtualSubaccountService {
	return &ModifyVirtualSubaccountService{c: c, body: map[string]any{
		"subAccountUid": subAccountUid,
		"permList":      permList,
		"status":        string(status),
	}}
}

func (s *ModifyVirtualSubaccountService) Do(ctx context.Context) (*ModifyVirtualSubaccountResult, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/modify-virtual-subaccount", s.body).WithSign()
	return request.Do[ModifyVirtualSubaccountResult](req)
}

// ModifyVirtualSubaccountResult reports whether the edit succeeded.
type ModifyVirtualSubaccountResult struct {
	Result string `json:"result"` // success, failure
}

// BatchCreateVirtualSubaccountAndApikeyService -- POST /api/v2/user/batch-create-virtual-subaccount-and-apikey (signed; state-changing)
//
// Creates a virtual sub-account together with its ApiKey in one call.
type BatchCreateVirtualSubaccountAndAPIKeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewBatchCreateVirtualSubaccountAndAPIKeyService(subAccountName, passphrase, label string, permList []VirtualSubaccountPerm) *BatchCreateVirtualSubaccountAndAPIKeyService {
	return &BatchCreateVirtualSubaccountAndAPIKeyService{c: c, body: map[string]any{
		"subAccountName": subAccountName,
		"passphrase":     passphrase,
		"label":          label,
		"permList":       permList,
	}}
}

func (s *BatchCreateVirtualSubaccountAndAPIKeyService) SetIPList(ipList []string) *BatchCreateVirtualSubaccountAndAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *BatchCreateVirtualSubaccountAndAPIKeyService) Do(ctx context.Context) ([]VirtualSubaccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/batch-create-virtual-subaccount-and-apikey", s.body).WithSign()
	resp, err := request.Do[[]VirtualSubaccountAPIKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// VirtualSubaccountApikey is a created sub-account ApiKey, including the
// (one-time-returned) secret key.
type VirtualSubaccountAPIKey struct {
	SubAccountUID    string                  `json:"subAccountUid"`
	SubAccountName   string                  `json:"subAccountName"`
	Label            string                  `json:"label"`
	SubAccountAPIKey string                  `json:"subAccountApiKey"`
	SecretKey        string                  `json:"secretKey"`
	PermList         []VirtualSubaccountPerm `json:"permList"`
	IPList           []string                `json:"ipList"`
}

// GetVirtualSubaccountListService -- GET /api/v2/user/virtual-subaccount-list (signed)
//
// Returns the main account's virtual sub-accounts, paginated.
type GetVirtualSubaccountListService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetVirtualSubaccountListService() *GetVirtualSubaccountListService {
	return &GetVirtualSubaccountListService{c: c, params: map[string]string{}}
}

func (s *GetVirtualSubaccountListService) SetLimit(limit string) *GetVirtualSubaccountListService {
	s.params["limit"] = limit
	return s
}

func (s *GetVirtualSubaccountListService) SetIDLessThan(idLessThan string) *GetVirtualSubaccountListService {
	s.params["idLessThan"] = idLessThan
	return s
}

func (s *GetVirtualSubaccountListService) SetStatus(status VirtualSubaccountStatus) *GetVirtualSubaccountListService {
	s.params["status"] = string(status)
	return s
}

func (s *GetVirtualSubaccountListService) Do(ctx context.Context) (*VirtualSubaccountList, error) {
	req := request.Get(ctx, s.c, "/api/v2/user/virtual-subaccount-list", s.params).WithSign()
	return request.Do[VirtualSubaccountList](req)
}

// VirtualSubaccountList is one page of virtual sub-accounts plus the paging
// cursor.
type VirtualSubaccountList struct {
	EndID          string              `json:"endId"`
	SubAccountList []VirtualSubaccount `json:"subAccountList"`
}

// VirtualSubaccount is a single virtual sub-account record.
type VirtualSubaccount struct {
	SubAccountUID  string                  `json:"subAccountUid"`
	SubAccountName string                  `json:"subAccountName"`
	Label          string                  `json:"label"`
	Status         VirtualSubaccountStatus `json:"status"`
	PermList       []VirtualSubaccountPerm `json:"permList"`
	CTime          time.Time               `json:"cTime"`
	UTime          time.Time               `json:"uTime"`
}

// CreateVirtualSubaccountApikeyService -- POST /api/v2/user/create-virtual-subaccount-apikey (signed; state-changing)
//
// Creates an ApiKey for an existing virtual sub-account.
type CreateVirtualSubaccountAPIKeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewCreateVirtualSubaccountAPIKeyService(subAccountUid, passphrase, label string, permList []VirtualSubaccountPerm) *CreateVirtualSubaccountAPIKeyService {
	return &CreateVirtualSubaccountAPIKeyService{c: c, body: map[string]any{
		"subAccountUid": subAccountUid,
		"passphrase":    passphrase,
		"label":         label,
		"permList":      permList,
	}}
}

func (s *CreateVirtualSubaccountAPIKeyService) SetIPList(ipList []string) *CreateVirtualSubaccountAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *CreateVirtualSubaccountAPIKeyService) Do(ctx context.Context) (*VirtualSubaccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/create-virtual-subaccount-apikey", s.body).WithSign()
	return request.Do[VirtualSubaccountAPIKey](req)
}

// ModifyVirtualSubaccountApikeyService -- POST /api/v2/user/modify-virtual-subaccount-apikey (signed; state-changing)
//
// Updates an existing virtual sub-account ApiKey's label, IP whitelist and
// permissions.
type ModifyVirtualSubaccountAPIKeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewModifyVirtualSubaccountAPIKeyService(subAccountUid, subAccountApiKey, passphrase, label string) *ModifyVirtualSubaccountAPIKeyService {
	return &ModifyVirtualSubaccountAPIKeyService{c: c, body: map[string]any{
		"subAccountUid":    subAccountUid,
		"subAccountApiKey": subAccountApiKey,
		"passphrase":       passphrase,
		"label":            label,
	}}
}

func (s *ModifyVirtualSubaccountAPIKeyService) SetIPList(ipList []string) *ModifyVirtualSubaccountAPIKeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *ModifyVirtualSubaccountAPIKeyService) SetPermList(permList []VirtualSubaccountPerm) *ModifyVirtualSubaccountAPIKeyService {
	s.body["permList"] = permList
	return s
}

func (s *ModifyVirtualSubaccountAPIKeyService) Do(ctx context.Context) (*VirtualSubaccountAPIKey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/modify-virtual-subaccount-apikey", s.body).WithSign()
	return request.Do[VirtualSubaccountAPIKey](req)
}

// GetVirtualSubaccountApikeyListService -- GET /api/v2/user/virtual-subaccount-apikey-list (signed)
//
// Returns the ApiKeys belonging to a virtual sub-account.
type GetVirtualSubaccountAPIKeyListService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetVirtualSubaccountAPIKeyListService(subAccountUid string) *GetVirtualSubaccountAPIKeyListService {
	return &GetVirtualSubaccountAPIKeyListService{c: c, params: map[string]string{
		"subAccountUid": subAccountUid,
	}}
}

func (s *GetVirtualSubaccountAPIKeyListService) Do(ctx context.Context) ([]VirtualSubaccountAPIKey, error) {
	req := request.Get(ctx, s.c, "/api/v2/user/virtual-subaccount-apikey-list", s.params).WithSign()
	resp, err := request.Do[[]VirtualSubaccountAPIKey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
