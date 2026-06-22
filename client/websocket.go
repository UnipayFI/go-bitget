package client

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	bitgetCommon "github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/gorilla/websocket"
)

// WebSocketClient holds the configuration shared by every Bitget UTA stream.
// Bitget multiplexes all channels of a kind over one of two gateways (public
// vs. private); credentials are only needed for the private gateway's login.
type WebSocketClient struct {
	publicURL  string
	privateURL string
	apiKey     string
	apiSecret  string
	passphrase string
	signFn     SignFn
	logger     log.Logger
	dialer     *websocket.Dialer
}

type WebSocketOption struct {
	network    bitgetCommon.Network
	publicURL  string
	privateURL string
	apiKey     string
	apiSecret  string
	passphrase string
	signFn     SignFn
	logger     log.Logger
	dialer     *websocket.Dialer
}

type WebSocketOptions func(*WebSocketOption)

func defaultWebSocketOption() *WebSocketOption {
	return &WebSocketOption{
		network: bitgetCommon.Mainnet,
		logger:  log.GetDefaultLogger(),
		dialer:  defaultDialer(),
	}
}

func defaultDialer() *websocket.Dialer {
	return &websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: true,
	}
}

func NewWebSocketClient(options ...WebSocketOptions) *WebSocketClient {
	opt := defaultWebSocketOption()
	for _, option := range options {
		option(opt)
	}
	publicURL := opt.network.WsPublicURL()
	if opt.publicURL != "" {
		publicURL = opt.publicURL
	}
	privateURL := opt.network.WsPrivateURL()
	if opt.privateURL != "" {
		privateURL = opt.privateURL
	}
	return &WebSocketClient{
		publicURL:  publicURL,
		privateURL: privateURL,
		apiKey:     opt.apiKey,
		apiSecret:  opt.apiSecret,
		passphrase: opt.passphrase,
		signFn:     opt.signFn,
		logger:     opt.logger,
		dialer:     opt.dialer,
	}
}

func (c *WebSocketClient) GetPublicURL() string         { return c.publicURL }
func (c *WebSocketClient) GetPrivateURL() string        { return c.privateURL }
func (c *WebSocketClient) GetAPIKey() string            { return c.apiKey }
func (c *WebSocketClient) GetAPISecret() string         { return c.apiSecret }
func (c *WebSocketClient) GetPassphrase() string        { return c.passphrase }
func (c *WebSocketClient) GetSignFn() SignFn            { return c.signFn }
func (c *WebSocketClient) GetLogger() log.Logger        { return c.logger }
func (c *WebSocketClient) GetDialer() *websocket.Dialer { return c.dialer }

// WithWebSocketAuth sets the credentials used to log in to the private stream.
func WithWebSocketAuth(apiKey, apiSecret, passphrase string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		opt.apiKey = apiKey
		opt.apiSecret = apiSecret
		opt.passphrase = passphrase
	}
}

func WithWebSocketNetwork(network bitgetCommon.Network) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.network = network }
}

// WithWebSocketPublicURL overrides the public stream URL. Empty is ignored.
func WithWebSocketPublicURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.publicURL = u }
}

// WithWebSocketPrivateURL overrides the private stream URL. Empty is ignored.
func WithWebSocketPrivateURL(u string) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.privateURL = u }
}

func WithWebSocketLogger(logger log.Logger) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.logger = logger }
}

// WithWebSocketSignFn overrides the default HMAC login signer.
func WithWebSocketSignFn(signFn SignFn) WebSocketOptions {
	return func(opt *WebSocketOption) { opt.signFn = signFn }
}

// WithWebSocketProxy routes the stream dial through the given proxy (http,
// https, socks5, socks5h). Invalid URLs are logged and skipped.
func WithWebSocketProxy(proxyURL string) WebSocketOptions {
	return func(opt *WebSocketOption) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithWebSocketProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.dialer.Proxy = http.ProxyURL(u)
			opt.dialer.NetDialContext = nil
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithWebSocketProxy: socks5 setup failed: %v", err)
				return
			}
			opt.dialer.Proxy = nil
			opt.dialer.NetDialContext = dialCtx
		default:
			opt.logger.Errorf("WithWebSocketProxy: unsupported scheme %q", u.Scheme)
		}
	}
}
