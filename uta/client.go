package uta

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*UTAClient)(nil)

// UTAClient is the REST client for Bitget's Unified Trading Account
// (/api/v3/*) endpoints. It embeds the shared core client, so a future
// classic-account client can reuse the same signing/transport layer with a
// different set of paths.
type UTAClient struct {
	*client.Client
}

// NewUTAClient constructs a unified-account REST client.
func NewUTAClient(options ...client.Options) *UTAClient {
	return &UTAClient{client.NewClient(options...)}
}

// SyncServerTime measures the client/server clock offset and stores it so that
// signed requests carry a timestamp the server accepts. Bitget rejects
// requests whose ACCESS-TIMESTAMP drifts too far from its own clock, so call
// this once at startup (and periodically for long-lived processes).
func (c *UTAClient) SyncServerTime(ctx context.Context) error {
	localBefore := time.Now().UnixMilli()
	resp, err := c.NewGetServerTimeService().Do(ctx)
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
