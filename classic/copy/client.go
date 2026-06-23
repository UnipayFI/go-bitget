// Package copy implements the Bitget classic-account Copy-Trading REST endpoints (futures + spot, trader/follower/broker) under /api/v2/copy/.
package copy

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*CopyClient)(nil)

// CopyClient is the REST client for the classic-account Copy-Trading REST endpoints (futures + spot, trader/follower/broker) under /api/v2/copy/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type CopyClient struct {
	*core.Client
}

// NewCopyClient constructs a classic-account CopyClient.
func NewCopyClient(options ...client.Options) *CopyClient {
	return &CopyClient{core.New(options...)}
}
