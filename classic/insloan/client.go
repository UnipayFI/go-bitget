// Package insloan implements the Bitget classic-account Institutional Loan REST endpoints under /api/v2/spot/ins-loan/.
package insloan

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*InsLoanClient)(nil)

// InsLoanClient is the REST client for the classic-account Institutional Loan REST endpoints under /api/v2/spot/ins-loan/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type InsLoanClient struct {
	*core.Client
}

// NewInsLoanClient constructs a classic-account InsLoanClient.
func NewInsLoanClient(options ...client.Options) *InsLoanClient {
	return &InsLoanClient{core.New(options...)}
}
