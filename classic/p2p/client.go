// Package p2p implements the Bitget classic-account P2P merchant REST endpoints under /api/v2/p2p/.
package p2p

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*P2PClient)(nil)

// P2PClient is the REST client for the classic-account P2P merchant REST endpoints under /api/v2/p2p/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type P2PClient struct {
	*core.Client
}

// NewP2PClient constructs a classic-account P2PClient.
func NewP2PClient(options ...client.Options) *P2PClient {
	return &P2PClient{core.New(options...)}
}
