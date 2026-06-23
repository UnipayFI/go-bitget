// Package margin implements the Bitget classic-account Margin (cross + isolated) REST endpoints under /api/v2/margin/.
package margin

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*MarginClient)(nil)

// MarginClient is the REST client for the classic-account Margin (cross + isolated) REST endpoints under /api/v2/margin/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type MarginClient struct {
	*core.Client
}

// NewMarginClient constructs a classic-account MarginClient.
func NewMarginClient(options ...client.Options) *MarginClient {
	return &MarginClient{core.New(options...)}
}
