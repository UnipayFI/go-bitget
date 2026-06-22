package uta

import (
	"context"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.WsClient = (*UTAWebSocketClient)(nil)

// UTAWebSocketClient streams Bitget Unified Trading Account public and private
// channels. Public channels need no credentials; private channels (account,
// position, order, fill) require WithWebSocketAuth and log in automatically.
type UTAWebSocketClient struct {
	*client.WebSocketClient
}

// NewUTAWebSocketClient constructs a unified-account WebSocket client.
func NewUTAWebSocketClient(options ...client.WebSocketOptions) *UTAWebSocketClient {
	return &UTAWebSocketClient{client.NewWebSocketClient(options...)}
}

// WsInstType routes a public subscription to a product line.
type WsInstType string

const (
	WsInstTypeSpot        WsInstType = "spot"
	WsInstTypeUSDTFutures WsInstType = "usdt-futures"
	WsInstTypeCoinFutures WsInstType = "coin-futures"
	WsInstTypeUSDCFutures WsInstType = "usdc-futures"
)

// wsInstTypeUTA is the instType used by every private (account-scoped) channel.
const wsInstTypeUTA = "UTA"

// Subscribe is the low-level escape hatch: it subscribes to an arbitrary channel
// and delivers each data push's raw bytes. Prefer the typed NewSubscribe*
// services; use this for channels the SDK does not yet wrap.
func (c *UTAWebSocketClient) Subscribe(ctx context.Context, private bool, arg request.WsArg, cb func([]byte, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRaw(ctx, c, private, arg, cb)
}
