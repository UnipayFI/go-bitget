// Package earn implements the Bitget classic-account Earn REST endpoints (savings, sharkfin, loan) under /api/v2/earn/.
package earn

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*EarnClient)(nil)

// EarnClient is the REST client for the classic-account Earn REST endpoints (savings, sharkfin, loan) under /api/v2/earn/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type EarnClient struct {
	*core.Client
}

// NewEarnClient constructs a classic-account EarnClient.
func NewEarnClient(options ...client.Options) *EarnClient {
	return &EarnClient{core.New(options...)}
}
