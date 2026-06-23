// Package bitget is the entry point of the Bitget exchange Go SDK.
//
// Install: go get github.com/UnipayFI/go-bitget
// Import:  import bitget "github.com/UnipayFI/go-bitget"
//
// The SDK covers BOTH of Bitget's account systems:
//
//   - Unified Trading Account (UTA): the /api/v3/* REST API and the v3
//     WebSocket streams, exposed through package uta.
//   - Classic account: the /api/v2/* REST API and the v2 WebSocket streams,
//     split by product line into packages under classic/ (spot, mix (futures),
//     margin, copy, earn, common, broker, affiliate, insloan, tax, p2p) plus the
//     shared classic/ws stream client.
//
// Authentication uses the HMAC-SHA256 scheme (ACCESS-KEY / ACCESS-SIGN /
// ACCESS-TIMESTAMP / ACCESS-PASSPHRASE); the core client/request/sign layers are
// shared by every product.
//
// Quick start (UTA):
//
//	c := bitget.NewUTAClient(client.WithAuth(apiKey, apiSecret, passphrase))
//	if err := c.SyncServerTime(ctx); err != nil { /* ... */ }
//	assets, err := c.NewGetAccountAssetsService().Do(ctx)
//
// Quick start (classic spot):
//
//	sp := bitget.NewSpotClient(client.WithAuth(apiKey, apiSecret, passphrase))
//	if err := sp.SyncServerTime(ctx); err != nil { /* ... */ }
//	tickers, err := sp.NewGetTickersService().Do(ctx)
package bitget

import (
	"github.com/UnipayFI/go-bitget/classic/affiliate"
	"github.com/UnipayFI/go-bitget/classic/broker"
	"github.com/UnipayFI/go-bitget/classic/common"
	"github.com/UnipayFI/go-bitget/classic/copy"
	"github.com/UnipayFI/go-bitget/classic/earn"
	"github.com/UnipayFI/go-bitget/classic/insloan"
	"github.com/UnipayFI/go-bitget/classic/margin"
	"github.com/UnipayFI/go-bitget/classic/mix"
	"github.com/UnipayFI/go-bitget/classic/p2p"
	"github.com/UnipayFI/go-bitget/classic/spot"
	"github.com/UnipayFI/go-bitget/classic/tax"
	"github.com/UnipayFI/go-bitget/classic/ws"
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

// --- Classic account (/api/v2/*) REST clients ---

// NewSpotClient constructs a classic-account Spot REST client.
func NewSpotClient(options ...client.Options) *spot.SpotClient {
	return spot.NewSpotClient(options...)
}

// NewMixClient constructs a classic-account Futures (Mix) REST client.
func NewMixClient(options ...client.Options) *mix.MixClient {
	return mix.NewMixClient(options...)
}

// NewMarginClient constructs a classic-account Margin (cross + isolated) REST client.
func NewMarginClient(options ...client.Options) *margin.MarginClient {
	return margin.NewMarginClient(options...)
}

// NewCommonClient constructs a classic-account common REST client (server time,
// announcements, convert, trade-rate, virtual sub-accounts, big-data insights).
func NewCommonClient(options ...client.Options) *common.CommonClient {
	return common.NewCommonClient(options...)
}

// NewCopyClient constructs a classic-account Copy-Trading REST client.
func NewCopyClient(options ...client.Options) *copy.CopyClient {
	return copy.NewCopyClient(options...)
}

// NewEarnClient constructs a classic-account Earn REST client.
func NewEarnClient(options ...client.Options) *earn.EarnClient {
	return earn.NewEarnClient(options...)
}

// NewBrokerClient constructs a classic-account Broker REST client.
func NewBrokerClient(options ...client.Options) *broker.BrokerClient {
	return broker.NewBrokerClient(options...)
}

// NewAffiliateClient constructs a classic-account Affiliate REST client.
func NewAffiliateClient(options ...client.Options) *affiliate.AffiliateClient {
	return affiliate.NewAffiliateClient(options...)
}

// NewInsLoanClient constructs a classic-account Institutional-Loan REST client.
func NewInsLoanClient(options ...client.Options) *insloan.InsLoanClient {
	return insloan.NewInsLoanClient(options...)
}

// NewTaxClient constructs a classic-account Tax REST client.
func NewTaxClient(options ...client.Options) *tax.TaxClient {
	return tax.NewTaxClient(options...)
}

// NewP2PClient constructs a classic-account P2P REST client.
func NewP2PClient(options ...client.Options) *p2p.P2PClient {
	return p2p.NewP2PClient(options...)
}

// NewClassicWebSocketClient constructs a WebSocket client for the classic-account
// v2 public and private streams (spot, futures, margin), including WebSocket
// order entry.
func NewClassicWebSocketClient(options ...client.WebSocketOptions) *ws.WebSocketClient {
	return ws.NewWebSocketClient(options...)
}
