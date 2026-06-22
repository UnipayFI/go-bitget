package client

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"

	bitgetCommon "github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
)

// SignFn produces the ACCESS-SIGN value from the prehash string. The default
// implementation is HMAC-SHA256(secret, prehash) base64-encoded; supply a
// custom one via WithSignFn to sign with an RSA private key instead (in which
// case secret carries the PEM-encoded key).
type SignFn = func(secret, prehash string) (signature string, err error)

type Option struct {
	apiKey       string
	apiSecret    string
	passphrase   string
	network      bitgetCommon.Network
	baseURL      string
	locale       string
	demoTrading  bool
	logger       log.Logger
	signFn       SignFn
	client       *resty.Client
	timeOffsetMs int64
}

type Options func(*Option)

func defaultOption() *Option {
	return &Option{
		network: bitgetCommon.Mainnet,
		locale:  bitgetCommon.DEFAULT_LOCALE,
		logger:  log.GetDefaultLogger(),
		client:  defaultHttpClient(),
	}
}

func defaultHttpClient() *resty.Client {
	return resty.New().
		SetJSONMarshaler(bitgetCommon.JSONMarshal).
		SetJSONUnmarshaler(bitgetCommon.JSONUnmarshal)
}

// WithAuth sets the API credentials used to sign private requests. All three
// values come from the Bitget API-management page; passphrase is the one set
// when the key was created.
func WithAuth(apiKey, apiSecret, passphrase string) Options {
	return func(opt *Option) {
		opt.apiKey = apiKey
		opt.apiSecret = apiSecret
		opt.passphrase = passphrase
	}
}

// WithNetwork selects the Bitget environment. Bitget exposes a single
// production domain, so this currently only accepts common.Mainnet; it exists
// for forward symmetry with sibling SDKs.
func WithNetwork(network bitgetCommon.Network) Options {
	return func(opt *Option) {
		opt.network = network
	}
}

// WithBaseURL overrides the REST base URL derived from WithNetwork. Use it to
// point the client at a custom or proxied endpoint. An empty value is ignored.
func WithBaseURL(baseURL string) Options {
	return func(opt *Option) {
		opt.baseURL = baseURL
	}
}

// WithLocale sets the ACCESS locale header (e.g. "en-US", "zh-CN"). It only
// affects the language of server-side error messages. Defaults to "en-US".
func WithLocale(locale string) Options {
	return func(opt *Option) {
		if locale != "" {
			opt.locale = locale
		}
	}
}

// WithDemoTrading routes requests to Bitget's demo (paper) trading environment
// by attaching the "paptrading: 1" header. The REST domain is unchanged.
func WithDemoTrading(enabled bool) Options {
	return func(opt *Option) {
		opt.demoTrading = enabled
	}
}

func WithLogger(logger log.Logger) Options {
	return func(opt *Option) {
		opt.logger = logger
	}
}

// WithSignFn replaces the default HMAC-SHA256 signer. Use it for RSA-signed
// keys or to delegate signing to an HSM / remote signer.
func WithSignFn(signFn SignFn) Options {
	return func(opt *Option) {
		opt.signFn = signFn
	}
}

// WithTimeOffset sets a fixed client/server clock offset in milliseconds. The
// request timestamp is computed as localMillis - timeOffsetMs. Usually set
// automatically via the client's SyncServerTime helper.
func WithTimeOffset(timeOffsetMs int64) Options {
	return func(opt *Option) {
		opt.timeOffsetMs = timeOffsetMs
	}
}

// WithHTTPClient supplies a pre-configured resty client (custom transport,
// timeouts, TLS, etc.). The JSON (un)marshalers and base URL are still set by
// the SDK afterwards.
func WithHTTPClient(client *resty.Client) Options {
	return func(opt *Option) {
		if client != nil {
			opt.client = client
		}
	}
}

// WithProxy routes all REST traffic through the given proxy. Supported
// schemes: http, https, socks5, socks5h. Pass userinfo in the URL for
// authenticated proxies. Invalid URLs are logged and skipped.
func WithProxy(proxyURL string) Options {
	return func(opt *Option) {
		if proxyURL == "" {
			return
		}
		u, err := url.Parse(proxyURL)
		if err != nil {
			opt.logger.Errorf("WithProxy: invalid proxy URL %q: %v", proxyURL, err)
			return
		}
		switch strings.ToLower(u.Scheme) {
		case "http", "https":
			opt.client.SetProxy(proxyURL)
		case "socks5", "socks5h":
			dialCtx, err := socks5DialContext(u)
			if err != nil {
				opt.logger.Errorf("WithProxy: socks5 setup failed: %v", err)
				return
			}
			transport := cloneDefaultTransport()
			transport.Proxy = nil
			transport.DialContext = dialCtx
			opt.client.SetTransport(transport)
		default:
			opt.logger.Errorf("WithProxy: unsupported scheme %q (want http, https, socks5, socks5h)", u.Scheme)
		}
	}
}

// socks5DialContext builds a DialContext that tunnels TCP through the SOCKS5
// proxy described by u. socks5h is accepted as an alias of socks5: the SOCKS5
// dialer in golang.org/x/net/proxy already resolves hostnames remotely.
func socks5DialContext(u *url.URL) (func(ctx context.Context, network, addr string) (net.Conn, error), error) {
	su := *u
	if strings.EqualFold(su.Scheme, "socks5h") {
		su.Scheme = "socks5"
	}
	pd, err := proxy.FromURL(&su, proxy.Direct)
	if err != nil {
		return nil, err
	}
	if cd, ok := pd.(proxy.ContextDialer); ok {
		return cd.DialContext, nil
	}
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return pd.Dial(network, addr)
	}, nil
}

func cloneDefaultTransport() *http.Transport {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		return t.Clone()
	}
	panic(fmt.Sprintf("bitget: http.DefaultTransport is not *http.Transport (got %T)", http.DefaultTransport))
}
