// Package broker implements the Bitget classic-account Broker REST endpoints (commission, sub-accounts, api-keys) under /api/v2/broker/.
package broker

import (
	"github.com/UnipayFI/go-bitget/classic/internal/core"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.Client = (*BrokerClient)(nil)

// BrokerClient is the REST client for the classic-account Broker REST endpoints (commission, sub-accounts, api-keys) under /api/v2/broker/. It embeds the shared
// classic core, reusing the same HMAC signing/transport layer and
// SyncServerTime helper as every other classic product client.
type BrokerClient struct {
	*core.Client
}

// NewBrokerClient constructs a classic-account BrokerClient.
func NewBrokerClient(options ...client.Options) *BrokerClient {
	return &BrokerClient{core.New(options...)}
}
