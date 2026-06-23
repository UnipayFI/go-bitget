// Package mix implements the Bitget classic-account Futures (Mix) REST (and, in
// the websocket files, the v2 Mix streams): market data, account, position,
// order and plan/trigger-order endpoints under the /api/v2/mix/ path namespace.
// It covers all three product types — USDT-FUTURES, USDC-FUTURES and
// COIN-FUTURES — selected per request via the productType parameter.
package mix

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*MixClient)(nil)

// MixClient is the REST client for the classic-account Futures (Mix) endpoints.
// It embeds the shared classic core, reusing the same HMAC signing/transport
// layer and SyncServerTime helper as every other classic product client.
type MixClient struct {
	*core.Client
}

// NewMixClient constructs a classic-account Futures (Mix) REST client.
func NewMixClient(options ...client.Options) *MixClient {
	return &MixClient{core.New(options...)}
}
