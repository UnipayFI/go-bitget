package request

import (
	"context"
	"maps"
	"net/http"
	"sort"
	"strings"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/pkg/log"
	"github.com/go-resty/resty/v2"
)

// Client is what every endpoint Service needs from a UTA/classic REST client.
// All getters are read-only; the concrete *client.Client satisfies it.
type Client interface {
	GetHttpClient() *resty.Client
	GetAPIKey() string
	GetAPISecret() string
	GetPassphrase() string
	GetLocale() string
	IsDemoTrading() bool
	GetLogger() log.Logger
	GetSignFn() SignFn
	TimestampMs() int64
}

type kv struct {
	Key   string
	Value string
}

type Request struct {
	client   Client
	r        *resty.Request
	method   string
	path     string
	query    []kv
	bodyJSON string
	needSign bool
	err      error
}

func newRequest(ctx context.Context, c Client, method, path string) *Request {
	r := c.GetHttpClient().R().
		SetHeader("User-Agent", common.GO_BITGET_USER_AGENT).
		SetHeader("locale", c.GetLocale()).
		SetContext(ctx)
	r.Method = method
	return &Request{
		client: c,
		r:      r,
		method: method,
		path:   path,
	}
}

// Get builds a GET request. Any params maps are merged and become the (sorted)
// query string, which is also part of the signed prehash.
func Get(ctx context.Context, c Client, path string, params ...map[string]string) *Request {
	r := newRequest(ctx, c, http.MethodGet, path)
	r.setQuery(params...)
	return r
}

// Post builds a POST request. Any body maps are merged and JSON-encoded once;
// the exact bytes sent are the bytes signed.
func Post(ctx context.Context, c Client, path string, body ...map[string]any) *Request {
	r := newRequest(ctx, c, http.MethodPost, path)
	r.setBody(mergeBody(body...))
	return r
}

func (r *Request) setQuery(params ...map[string]string) {
	merged := make(map[string]string)
	for _, p := range params {
		for k, v := range p {
			if v == "" {
				continue
			}
			merged[k] = v
		}
	}
	if len(merged) == 0 {
		return
	}
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	r.query = make([]kv, 0, len(keys))
	for _, k := range keys {
		r.query = append(r.query, kv{Key: k, Value: merged[k]})
	}
}

// SetBody overrides the request body with an arbitrary JSON-serializable value
// (used for batch endpoints whose body is an array or a nested struct rather
// than a flat map). The value is marshaled once and reused for signing.
func (r *Request) SetBody(body any) *Request {
	if r.err != nil {
		return r
	}
	data, err := common.JSONMarshal(body)
	if err != nil {
		r.err = err
		return r
	}
	r.bodyJSON = common.BytesToString(data)
	return r
}

func (r *Request) setBody(body map[string]any) {
	if len(body) == 0 {
		return
	}
	data, err := common.JSONMarshal(body)
	if err != nil {
		r.err = err
		return
	}
	r.bodyJSON = common.BytesToString(data)
}

func mergeBody(body ...map[string]any) map[string]any {
	merged := make(map[string]any)
	for _, b := range body {
		maps.Copy(merged, b)
	}
	return merged
}

// WithSign marks the request as private: ACCESS-* signing headers are attached
// at send time. Public market endpoints omit this.
func (r *Request) WithSign() *Request {
	r.needSign = true
	return r
}

// requestPath returns the path plus, when present, the "?"-prefixed query
// string. It is shared by the URL builder and the prehash so the bytes the
// client sends match the bytes the server signs.
func (r *Request) requestPath() string {
	if len(r.query) == 0 {
		return r.path
	}
	return r.path + "?" + encodeQuery(r.query)
}

// encodeQuery joins the already-sorted params as key=value pairs. Bitget signs
// the literal query string, and its endpoints only use plain ASCII values, so
// no percent-encoding is applied (matching the official SDKs).
func encodeQuery(params []kv) string {
	var b strings.Builder
	for i, p := range params {
		if i > 0 {
			b.WriteByte('&')
		}
		b.WriteString(p.Key)
		b.WriteByte('=')
		b.WriteString(p.Value)
	}
	return b.String()
}

func (r *Request) fullURL() string {
	base := strings.TrimSuffix(r.client.GetHttpClient().BaseURL, "/")
	path := r.path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	urlStr := base + path
	if len(r.query) > 0 {
		urlStr += "?" + encodeQuery(r.query)
	}
	return urlStr
}
