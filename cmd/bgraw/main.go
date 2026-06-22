// Command bgraw signs and executes a single Bitget UTA REST call and pretty
// prints the raw response. It is a development aid for capturing the exact
// shape of private endpoints (which cannot be curled without HMAC signing) so
// the typed response structs can be reconciled against reality.
//
// Usage:
//
//	BITGET_API_KEY=... BITGET_API_SECRET=... BITGET_PASSPHRASE=... \
//	  go run ./cmd/bgraw GET  /api/v3/account/info
//	  go run ./cmd/bgraw GET  /api/v3/account/financial-records "coin=USDT&limit=5"
//	  go run ./cmd/bgraw POST /api/v3/trade/place-order '{"category":"SPOT", ...}'
//
// The third argument is the query string (GET) or JSON body (POST). Set
// BITGET_PROXY to route through a proxy.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/UnipayFI/go-bitget/client"
	bitgetCommon "github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "usage: bgraw <GET|POST> <path> [query-or-jsonbody]")
		os.Exit(2)
	}
	method := strings.ToUpper(os.Args[1])
	path := os.Args[2]
	arg := ""
	if len(os.Args) > 3 {
		arg = os.Args[3]
	}

	opts := []client.Options{
		client.WithAuth(
			os.Getenv("BITGET_API_KEY"),
			os.Getenv("BITGET_API_SECRET"),
			os.Getenv("BITGET_PASSPHRASE"),
		),
	}
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithProxy(proxy))
	}
	c := client.NewClient(opts...)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var req *request.Request
	switch method {
	case "GET":
		req = request.Get(ctx, c, path, parseQuery(arg)).WithSign()
	case "POST":
		body := map[string]any{}
		if arg != "" {
			if err := bitgetCommon.JSONUnmarshal([]byte(arg), &body); err != nil {
				fail("invalid json body: %v", err)
			}
		}
		req = request.Post(ctx, c, path, body).WithSign()
	default:
		fail("unsupported method %q", method)
	}

	body, err := request.DoRaw(req)
	if err != nil {
		fail("request error: %v", err)
	}
	fmt.Println(pretty(body))
}

func parseQuery(q string) map[string]string {
	out := map[string]string{}
	q = strings.TrimPrefix(q, "?")
	for pair := range strings.SplitSeq(q, "&") {
		if pair == "" {
			continue
		}
		k, v, _ := strings.Cut(pair, "=")
		out[k] = v
	}
	return out
}

func pretty(b []byte) string {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return string(b)
	}
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(b)
	}
	return string(out)
}

func fail(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
