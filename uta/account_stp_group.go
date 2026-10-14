package uta

import (
	"context"
	"strings"
	"time"

	"github.com/UnipayFI/go-bitget/request"
)

// CreateSTPGroupService -- POST /api/v3/account/create-stp-groups (UTA mgt. read & write)
//
// Creates a custom Self-Trade Prevention (STP) group. Sub-accounts in the same
// group never match against each other; sub-accounts in different groups trade
// with each other normally. groupName (max 32 chars) must be unique. Only
// available to MM (Market Maker) whitelist users.
type CreateSTPGroupService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewCreateSTPGroupService(groupName string) *CreateSTPGroupService {
	return &CreateSTPGroupService{c: c, body: map[string]any{
		"groupName": groupName,
	}}
}

func (s *CreateSTPGroupService) Do(ctx context.Context) (*STPGroup, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/create-stp-groups", s.body).WithSign()
	return request.Do[STPGroup](req)
}

// GetSTPGroupsService -- GET /api/v3/account/stp-groups (UTA mgt. read)
//
// Returns the master account's custom STP groups with their member
// sub-accounts. Only available to MM (Market Maker) whitelist users.
type GetSTPGroupsService struct {
	c *UTAClient
}

func (c *UTAClient) NewGetSTPGroupsService() *GetSTPGroupsService {
	return &GetSTPGroupsService{c: c}
}

func (s *GetSTPGroupsService) Do(ctx context.Context) (*STPGroups, error) {
	req := request.Get(ctx, s.c, "/api/v3/account/stp-groups").WithSign()
	return request.Do[STPGroups](req)
}

type STPGroups struct {
	Groups []STPGroup `json:"groups"`
}

// STPGroup is a custom STP group. Members is only set by GetSTPGroupsService.
type STPGroup struct {
	GroupID     string           `json:"groupId"`
	GroupName   string           `json:"groupName"`
	CreatedTime time.Time        `json:"createdTime,format:unixmilli"`
	Members     []STPGroupMember `json:"members"`
}

type STPGroupMember struct {
	UID        string    `json:"uid"`
	JoinedTime time.Time `json:"joinedTime,format:unixmilli"`
}

// AddSTPGroupMembersService -- POST /api/v3/account/add-groups-members (UTA mgt. read & write)
//
// Adds up to 50 of the master account's sub-accounts to an STP group. A
// sub-account belongs to at most one group, so adding it to a new group
// removes it from its previous one. Only available to MM (Market Maker)
// whitelist users.
type AddSTPGroupMembersService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewAddSTPGroupMembersService(groupID string, uids []string) *AddSTPGroupMembersService {
	return &AddSTPGroupMembersService{c: c, body: map[string]any{
		"groupId": groupID,
		"uids":    strings.Join(uids, ","),
	}}
}

func (s *AddSTPGroupMembersService) Do(ctx context.Context) (*STPGroupMembersResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/add-groups-members", s.body).WithSign()
	return request.Do[STPGroupMembersResult](req)
}

// RemoveSTPGroupMembersService -- POST /api/v3/account/remove-groups-members (UTA mgt. read & write)
//
// Removes up to 50 sub-accounts from an STP group; a removed sub-account falls
// back to the default behavior of sharing its master account's STP scope.
// Only available to MM (Market Maker) whitelist users.
type RemoveSTPGroupMembersService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewRemoveSTPGroupMembersService(groupID string, uids []string) *RemoveSTPGroupMembersService {
	return &RemoveSTPGroupMembersService{c: c, body: map[string]any{
		"groupId": groupID,
		"uids":    strings.Join(uids, ","),
	}}
}

func (s *RemoveSTPGroupMembersService) Do(ctx context.Context) (*STPGroupMembersResult, error) {
	req := request.Post(ctx, s.c, "/api/v3/account/remove-groups-members", s.body).WithSign()
	return request.Do[STPGroupMembersResult](req)
}

// STPGroupMembersResult is the per-UID outcome of adding or removing STP group
// members.
type STPGroupMembersResult struct {
	CreatedTime time.Time     `json:"createdTime,format:unixmilli"`
	SuccessUIDs []STPGroupUID `json:"successUids"`
	FailedUIDs  []STPGroupUID `json:"failedUids"`
}

type STPGroupUID struct {
	UID string `json:"uid"`
}
