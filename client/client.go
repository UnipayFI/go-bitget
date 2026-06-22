package client

import (
	"time"

	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Client is the shared, product-agnostic REST core. Every Bitget business line
// (unified account, classic spot/futures/margin, broker, ...) is just a set of
// request paths layered on top of this same signing + transport machinery, so
// the core carries no product-specific state.
type Client struct {
	client *resty.Client

	apiKey       string
	apiSecret    string
	passphrase   string
	locale       string
	demoTrading  bool
	logger       log.Logger
	signFn       SignFn
	timeOffsetMs int64
}

func NewClient(options ...Options) *Client {
	opt := defaultOption()
	for _, option := range options {
		option(opt)
	}

	baseURL := opt.network.RestBaseURL()
	if opt.baseURL != "" {
		baseURL = opt.baseURL
	}
	opt.client.SetBaseURL(baseURL)

	return &Client{
		client:       opt.client,
		apiKey:       opt.apiKey,
		apiSecret:    opt.apiSecret,
		passphrase:   opt.passphrase,
		locale:       opt.locale,
		demoTrading:  opt.demoTrading,
		logger:       opt.logger,
		signFn:       opt.signFn,
		timeOffsetMs: opt.timeOffsetMs,
	}
}

func (c *Client) GetHttpClient() *resty.Client { return c.client }

func (c *Client) GetAPIKey() string { return c.apiKey }

func (c *Client) GetAPISecret() string { return c.apiSecret }

func (c *Client) GetPassphrase() string { return c.passphrase }

func (c *Client) GetLocale() string { return c.locale }

func (c *Client) IsDemoTrading() bool { return c.demoTrading }

func (c *Client) GetLogger() log.Logger { return c.logger }

func (c *Client) GetSignFn() SignFn { return c.signFn }

func (c *Client) GetTimeOffsetMs() int64 { return c.timeOffsetMs }

func (c *Client) SetTimeOffset(offsetMs int64) { c.timeOffsetMs = offsetMs }

// TimestampMs returns the current request timestamp in milliseconds, adjusted
// by the configured client/server clock offset. This is the value placed in
// the ACCESS-TIMESTAMP header and signed as part of the prehash string.
func (c *Client) TimestampMs() int64 {
	return time.Now().UnixMilli() - c.timeOffsetMs
}
