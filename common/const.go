package common

import "time"

const (
	GO_BITGET_USER_AGENT = "go-bitget/1.0"

	// REST endpoint. Bitget serves every business line (spot, futures, margin,
	// unified account) from a single domain; the product is encoded in the
	// request path, not the host.
	DEFAULT_REST_BASE_URL = "https://api.bitget.com"

	// WebSocket endpoints. The unified-account streams use the dedicated v3
	// gateway (topic/symbol args, lowercase instType), distinct from the v2
	// socket. Demo trading uses the wspap.bitget.com host instead.
	DEFAULT_WS_PUBLIC_URL  = "wss://ws.bitget.com/v3/ws/public"
	DEFAULT_WS_PRIVATE_URL = "wss://ws.bitget.com/v3/ws/private"

	// Classic-account (non-UTA) WebSocket endpoints. The classic streams use the
	// v2 gateway (instType/channel/instId args, uppercase instType like SPOT /
	// USDT-FUTURES), distinct from the unified-account v3 gateway above.
	DEFAULT_WS_V2_PUBLIC_URL  = "wss://ws.bitget.com/v2/ws/public"
	DEFAULT_WS_V2_PRIVATE_URL = "wss://ws.bitget.com/v2/ws/private"

	// Default ACCESS-* locale header. Bitget accepts "en-US", "zh-CN", etc.
	DEFAULT_LOCALE = "en-US"

	DEFAULT_KEEP_ALIVE_INTERVAL = 25 * time.Second
	DEFAULT_KEEP_ALIVE_TIMEOUT  = 60 * time.Second
)

// Network identifies which Bitget environment a client targets. Bitget has no
// separate testnet host; demo (paper) trading runs on the same domain and is
// toggled with the "paptrading: 1" header (see client.WithDemoTrading). The
// type is kept for forward symmetry with sibling SDKs and to leave room for a
// future dedicated environment.
type Network int

const (
	Mainnet Network = iota
)

// RestBaseURL returns the REST base URL for this network.
func (n Network) RestBaseURL() string {
	return DEFAULT_REST_BASE_URL
}

// WsPublicURL returns the public WebSocket URL for this network.
func (n Network) WsPublicURL() string {
	return DEFAULT_WS_PUBLIC_URL
}

// WsPrivateURL returns the private WebSocket URL for this network.
func (n Network) WsPrivateURL() string {
	return DEFAULT_WS_PRIVATE_URL
}
