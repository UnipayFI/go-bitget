// Package core holds the REST client base shared by every classic-account
// product package (classic/spot, classic/mix, classic/margin, ...). Each
// product client embeds *core.Client to inherit the product-agnostic plumbing:
// credentials, transport and signing (via the embedded *client.Client) plus
// server-time clock synchronization. The product path namespaces differ
// (/api/v2/spot, /api/v2/mix, ...), but the HMAC auth scheme and the v2 public
// server-time endpoint are shared across the whole classic API, so this base
// is defined once rather than copied into each package.
package core

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*Client)(nil)

// Client is the shared base embedded by each classic product client.
type Client struct {
	*client.Client
}

// New constructs the shared base from the standard client options.
func New(options ...client.Options) *Client {
	return &Client{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp the server accepts. Bitget rejects
// requests whose ACCESS-TIMESTAMP drifts too far from its own clock, so call
// this once at startup (and periodically for long-lived processes). All classic
// products share the v2 public server-time endpoint.
func (c *Client) SyncServerTime(ctx context.Context) error {
	localBefore := time.Now().UnixMilli()
	req := request.Get(ctx, c, "/api/v2/public/time")
	resp, err := request.Do[serverTime](req)
	if err != nil {
		return err
	}
	localAfter := time.Now().UnixMilli()
	local := (localBefore + localAfter) / 2
	c.SetTimeOffset(local - resp.ServerTime.UnixMilli())
	c.GetLogger().Infof("Time sync: local=%d, server=%d, offset=%dms",
		local, resp.ServerTime.UnixMilli(), c.GetTimeOffsetMs())
	return nil
}

// serverTime is the minimal shape of GET /api/v2/public/time used by
// SyncServerTime. The public Service that exposes server time to users lives in
// classic/common.
type serverTime struct {
	ServerTime time.Time `json:"serverTime"`
}
