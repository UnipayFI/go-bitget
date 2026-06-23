// Package common implements the Bitget classic-account "common" REST endpoints:
// the cross-product utilities that are not specific to a single trading product
// — server time, announcements, trade fee rates, the all-account balance
// overview, coin conversion (Convert / BGB-Convert), virtual sub-account
// management, and the public big-data trading-insight feeds.
package common

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*CommonClient)(nil)

// CommonClient is the REST client for the classic-account common endpoints. It
// embeds the shared classic core, reusing the same HMAC signing/transport layer
// and SyncServerTime helper as every other classic product client.
type CommonClient struct {
	*core.Client
}

// NewCommonClient constructs a classic-account common REST client.
func NewCommonClient(options ...client.Options) *CommonClient {
	return &CommonClient{core.New(options...)}
}
