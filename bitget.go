// Package bitget is the entry point of the Bitget exchange Go SDK.
//
// Install: go get github.com/UnipayFI/go-bitget
// Import:  import bitget "github.com/UnipayFI/go-bitget"
//
// The SDK targets Bitget's Unified Trading Account (UTA, the /api/v3/* REST API
// and the v3 WebSocket streams). Authentication uses the HMAC-SHA256 scheme
// (ACCESS-KEY / ACCESS-SIGN / ACCESS-TIMESTAMP / ACCESS-PASSPHRASE). The core
// client/request/sign layers are product-agnostic, leaving room for a future
// classic-account module alongside uta.
//
// Quick start:
//
//	c := bitget.NewUTAClient(
//		client.WithAuth(apiKey, apiSecret, passphrase),
//	)
//	if err := c.SyncServerTime(ctx); err != nil { /* ... */ }
//	assets, err := c.NewGetAccountAssetsService().Do(ctx)
//
// WebSocket:
//
//	ws := bitget.NewUTAWebSocketClient(
//		client.WithWebSocketAuth(apiKey, apiSecret, passphrase),
//	)
//	done, _, err := ws.NewSubscribeTickerService(uta.WsInstTypeUSDTFutures, "BTCUSDT").
//		Do(ctx, func(p *request.WsPush[[]uta.WsTicker], err error) { /* ... */ })
package bitget

import (
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/uta"
)

// NewUTAClient constructs a REST client for the unified-account /api/v3/*
// endpoints.
func NewUTAClient(options ...client.Options) *uta.UTAClient {
	return uta.NewUTAClient(options...)
}

// NewUTAWebSocketClient constructs a WebSocket client for the unified-account
// public and private streams.
func NewUTAWebSocketClient(options ...client.WebSocketOptions) *uta.UTAWebSocketClient {
	return uta.NewUTAWebSocketClient(options...)
}
