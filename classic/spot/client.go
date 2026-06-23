// Package spot implements the Bitget classic-account Spot REST (and, in the
// websocket files, the v2 Spot streams): market data, spot trading, plan
// (trigger) orders, account information and the wallet/transfer endpoints under
// the /api/v2/spot/ path namespace.
package spot

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*SpotClient)(nil)

// SpotClient is the REST client for the classic-account Spot endpoints. It
// embeds the shared classic core, reusing the same HMAC signing/transport layer
// and SyncServerTime helper as every other classic product client.
type SpotClient struct {
	*core.Client
}

// NewSpotClient constructs a classic-account Spot REST client.
func NewSpotClient(options ...client.Options) *SpotClient {
	return &SpotClient{core.New(options...)}
}
