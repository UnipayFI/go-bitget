// Package ws implements the Bitget classic-account (non-UTA) v2 WebSocket
// streams: public market channels (ticker, candlestick, depth/order-book,
// trades, auction) and private channels (account, orders, fill, positions, plan
// orders, ...) for spot, futures (mix) and margin, plus WebSocket order entry
// (op:"trade"). Subscriptions are routed by instType (SPOT / USDT-FUTURES /
// COIN-FUTURES / USDC-FUTURES) over a single shared v2 gateway, distinct from
// the unified-account v3 streams in package uta.
package ws

import (
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
)

var _ request.WsClient = (*WebSocketClient)(nil)

// WebSocketClient streams the classic-account v2 public and private channels.
// Public channels need no credentials; private channels require
// client.WithWebSocketAuth and log in automatically on first subscribe.
type WebSocketClient struct {
	*client.WebSocketClient
}

// NewWebSocketClient constructs a classic-account v2 WebSocket client. It
// defaults to the v2 gateway URLs (wss://ws.bitget.com/v2/ws/{public,private});
// pass client.WithWebSocketPublicURL / WithWebSocketPrivateURL to override.
func NewWebSocketClient(options ...client.WebSocketOptions) *WebSocketClient {
	opts := append([]client.WebSocketOptions{
		client.WithWebSocketPublicURL(common.DEFAULT_WS_V2_PUBLIC_URL),
		client.WithWebSocketPrivateURL(common.DEFAULT_WS_V2_PRIVATE_URL),
	}, options...)
	return &WebSocketClient{client.NewWebSocketClient(opts...)}
}
