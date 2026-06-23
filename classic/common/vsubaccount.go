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
	SubAccountUid  string                  `json:"subaAccountUid"`
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
type BatchCreateVirtualSubaccountAndApikeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewBatchCreateVirtualSubaccountAndApikeyService(subAccountName, passphrase, label string, permList []VirtualSubaccountPerm) *BatchCreateVirtualSubaccountAndApikeyService {
	return &BatchCreateVirtualSubaccountAndApikeyService{c: c, body: map[string]any{
		"subAccountName": subAccountName,
		"passphrase":     passphrase,
		"label":          label,
		"permList":       permList,
	}}
}

func (s *BatchCreateVirtualSubaccountAndApikeyService) SetIpList(ipList []string) *BatchCreateVirtualSubaccountAndApikeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *BatchCreateVirtualSubaccountAndApikeyService) Do(ctx context.Context) ([]VirtualSubaccountApikey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/batch-create-virtual-subaccount-and-apikey", s.body).WithSign()
	resp, err := request.Do[[]VirtualSubaccountApikey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

// VirtualSubaccountApikey is a created sub-account ApiKey, including the
// (one-time-returned) secret key.
type VirtualSubaccountApikey struct {
	SubAccountUid    string                  `json:"subAccountUid"`
	SubAccountName   string                  `json:"subAccountName"`
	Label            string                  `json:"label"`
	SubAccountApiKey string                  `json:"subAccountApiKey"`
	SecretKey        string                  `json:"secretKey"`
	PermList         []VirtualSubaccountPerm `json:"permList"`
	IpList           []string                `json:"ipList"`
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

func (s *GetVirtualSubaccountListService) SetIdLessThan(idLessThan string) *GetVirtualSubaccountListService {
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
	EndId          string              `json:"endId"`
	SubAccountList []VirtualSubaccount `json:"subAccountList"`
}

// VirtualSubaccount is a single virtual sub-account record.
type VirtualSubaccount struct {
	SubAccountUid  string                  `json:"subAccountUid"`
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
type CreateVirtualSubaccountApikeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewCreateVirtualSubaccountApikeyService(subAccountUid, passphrase, label string, permList []VirtualSubaccountPerm) *CreateVirtualSubaccountApikeyService {
	return &CreateVirtualSubaccountApikeyService{c: c, body: map[string]any{
		"subAccountUid": subAccountUid,
		"passphrase":    passphrase,
		"label":         label,
		"permList":      permList,
	}}
}

func (s *CreateVirtualSubaccountApikeyService) SetIpList(ipList []string) *CreateVirtualSubaccountApikeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *CreateVirtualSubaccountApikeyService) Do(ctx context.Context) (*VirtualSubaccountApikey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/create-virtual-subaccount-apikey", s.body).WithSign()
	return request.Do[VirtualSubaccountApikey](req)
}

// ModifyVirtualSubaccountApikeyService -- POST /api/v2/user/modify-virtual-subaccount-apikey (signed; state-changing)
//
// Updates an existing virtual sub-account ApiKey's label, IP whitelist and
// permissions.
type ModifyVirtualSubaccountApikeyService struct {
	c    *CommonClient
	body map[string]any
}

func (c *CommonClient) NewModifyVirtualSubaccountApikeyService(subAccountUid, subAccountApiKey, passphrase, label string) *ModifyVirtualSubaccountApikeyService {
	return &ModifyVirtualSubaccountApikeyService{c: c, body: map[string]any{
		"subAccountUid":    subAccountUid,
		"subAccountApiKey": subAccountApiKey,
		"passphrase":       passphrase,
		"label":            label,
	}}
}

func (s *ModifyVirtualSubaccountApikeyService) SetIpList(ipList []string) *ModifyVirtualSubaccountApikeyService {
	s.body["ipList"] = ipList
	return s
}

func (s *ModifyVirtualSubaccountApikeyService) SetPermList(permList []VirtualSubaccountPerm) *ModifyVirtualSubaccountApikeyService {
	s.body["permList"] = permList
	return s
}

func (s *ModifyVirtualSubaccountApikeyService) Do(ctx context.Context) (*VirtualSubaccountApikey, error) {
	req := request.Post(ctx, s.c, "/api/v2/user/modify-virtual-subaccount-apikey", s.body).WithSign()
	return request.Do[VirtualSubaccountApikey](req)
}

// GetVirtualSubaccountApikeyListService -- GET /api/v2/user/virtual-subaccount-apikey-list (signed)
//
// Returns the ApiKeys belonging to a virtual sub-account.
type GetVirtualSubaccountApikeyListService struct {
	c      *CommonClient
	params map[string]string
}

func (c *CommonClient) NewGetVirtualSubaccountApikeyListService(subAccountUid string) *GetVirtualSubaccountApikeyListService {
	return &GetVirtualSubaccountApikeyListService{c: c, params: map[string]string{
		"subAccountUid": subAccountUid,
	}}
}

func (s *GetVirtualSubaccountApikeyListService) Do(ctx context.Context) ([]VirtualSubaccountApikey, error) {
	req := request.Get(ctx, s.c, "/api/v2/user/virtual-subaccount-apikey-list", s.params).WithSign()
	resp, err := request.Do[[]VirtualSubaccountApikey](req)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}
