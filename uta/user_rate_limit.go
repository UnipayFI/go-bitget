package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/request"
)

// SetRateLimitQuotaService -- POST /api/v3/user/set-rate-limit-quota (UTA mgt. read & write)
//
// Sets the rate-limit quota for one or more sub-accounts on a product line
// (category is "spot" or "futures"), single or batch with up to 50 UIDs per
// request. Only available to MM (Market Maker) and PRO institutional users: the
// per-account quota must not exceed the single-account cap, and the summed
// sub-account quota must not exceed the master-sub cap. The reply data is an
// empty object.
type SetRateLimitQuotaService struct {
	c    *UTAClient
	body map[string]any
}

func (c *UTAClient) NewSetRateLimitQuotaService(category string, uids []string, quota string) *SetRateLimitQuotaService {
	return &SetRateLimitQuotaService{c: c, body: map[string]any{
		"category": category,
		"uids":     uids,
		"quota":    quota,
	}}
}

func (s *SetRateLimitQuotaService) Do(ctx context.Context) (*any, error) {
	req := request.Post(ctx, s.c, "/api/v3/user/set-rate-limit-quota", s.body).WithSign()
	return request.Do[any](req)
}

// GetRateLimitQuotaService -- GET /api/v3/user/rate-limit-quota (UTA mgt. read)
//
// Returns the rate-limit quota configuration for a product line (category is
// "spot" or "futures"): a single account when uid is set, otherwise every
// sub-account under the master, paginated by cursor. Only available to MM
// (Market Maker) and PRO institutional users.
type GetRateLimitQuotaService struct {
	c      *UTAClient
	params map[string]string
}

func (c *UTAClient) NewGetRateLimitQuotaService(category string) *GetRateLimitQuotaService {
	return &GetRateLimitQuotaService{c: c, params: map[string]string{"category": category}}
}

// SetUID filters to a single account (all sub-accounts when omitted).
func (s *GetRateLimitQuotaService) SetUID(uid string) *GetRateLimitQuotaService {
	s.params["uid"] = uid
	return s
}

func (s *GetRateLimitQuotaService) SetCursor(cursor string) *GetRateLimitQuotaService {
	s.params["cursor"] = cursor
	return s
}

// SetLimit sets the page size (default 100, max 100).
func (s *GetRateLimitQuotaService) SetLimit(limit string) *GetRateLimitQuotaService {
	s.params["limit"] = limit
	return s
}

func (s *GetRateLimitQuotaService) Do(ctx context.Context) (*RateLimitQuota, error) {
	req := request.Get(ctx, s.c, "/api/v3/user/rate-limit-quota", s.params).WithSign()
	return request.Do[RateLimitQuota](req)
}

type RateLimitQuota struct {
	Category       string               `json:"category"` // spot, futures
	QuotaList      []RateLimitQuotaItem `json:"quotaList"`
	MasterSubQuota string               `json:"masterSubQuota"`
	MasterSubCap   string               `json:"masterSubCap"`
	Cursor         string               `json:"cursor"`
}

type RateLimitQuotaItem struct {
	UID   string `json:"uid"`
	Quota string `json:"quota"`
}
