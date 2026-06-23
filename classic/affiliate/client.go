// Package affiliate implements the Bitget classic-account Affiliate/agent customer-info REST endpoints under /api/v2/broker/.
package affiliate

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*AffiliateClient)(nil)

// AffiliateClient is the REST client for the classic-account Affiliate/agent customer-info REST endpoints under /api/v2/broker/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type AffiliateClient struct {
	*core.Client
}

// NewAffiliateClient constructs a classic-account AffiliateClient.
func NewAffiliateClient(options ...client.Options) *AffiliateClient {
	return &AffiliateClient{core.New(options...)}
}
