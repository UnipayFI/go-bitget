// Package tax implements the Bitget classic-account Tax transaction-record REST endpoints under /api/v2/tax/.
package tax

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*TaxClient)(nil)

// TaxClient is the REST client for the classic-account Tax transaction-record REST endpoints under /api/v2/tax/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type TaxClient struct {
	*core.Client
}

// NewTaxClient constructs a classic-account TaxClient.
func NewTaxClient(options ...client.Options) *TaxClient {
	return &TaxClient{core.New(options...)}
}
