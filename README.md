# go-bitget

A Go SDK for the [Bitget](https://www.bitget.com/api-doc/uta/intro) exchange,
targeting the **Unified Trading Account (UTA)** — the `/api/v3/*` REST API and
the v3 WebSocket streams.

- **Faithful to the API.** Every response struct is reconciled against the
  *live* API, not just the docs (the docs are frequently incomplete). Each
  public endpoint and every testable private endpoint is covered by a test that
  diffs the real response's JSON keys against the struct.
- **Ergonomic.** Functional options (`WithAuth`, `WithProxy`, …), a fluent
  `NewXxxService(...).SetFoo(...).Do(ctx)` pattern per endpoint, and generics
  for response decoding.
- **Correct money & time.** Prices/quantities are `shopspring/decimal.Decimal`;
  millisecond timestamps are `time.Time`. Bitget's quirk of sending numbers and
  timestamps as JSON strings (and `""`/`"0"`/`"-1"` for "not set") is handled
  transparently.
- **Extensible.** The `client` / `request` / signing core is product-agnostic,
  leaving room for a classic-account module alongside `uta`.

## Install

```bash
go get github.com/UnipayFI/go-bitget
```

Requires Go 1.26+.

## Quick start

```go
package main

import (
	"context"
	"fmt"

	bitget "github.com/UnipayFI/go-bitget"
	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/uta"
	"github.com/shopspring/decimal"
)

func main() {
	ctx := context.Background()

	c := bitget.NewUTAClient(
		client.WithAuth("apiKey", "apiSecret", "passphrase"),
		// client.WithProxy("socks5://127.0.0.1:7890"),
		// client.WithDemoTrading(true),
	)

	// Align the request clock with the server (avoids signature drift).
	_ = c.SyncServerTime(ctx)

	// Public market data (no auth required).
	instruments, _ := c.NewGetInstrumentsService(uta.CategorySpot).
		SetSymbol("BTCUSDT").Do(ctx)
	fmt.Println(instruments[0].Symbol, instruments[0].PricePrecision)

	// Private account data.
	assets, _ := c.NewGetAccountAssetsService().Do(ctx)
	fmt.Println("equity:", assets.AccountEquity)

	// Place a limit order.
	ref, err := c.NewPlaceOrderService(uta.CategorySpot, "BTCUSDT",
		decimal.RequireFromString("0.0001"), uta.SideBuy, uta.OrderTypeLimit).
		SetPrice(decimal.RequireFromString("30000")).
		SetTimeInForce(uta.TimeInForceGTC).
		Do(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("orderId:", ref.OrderId)
}
```

## Authentication

Credentials come from the Bitget API-management page (key, secret, passphrase):

```go
c := bitget.NewUTAClient(client.WithAuth(apiKey, apiSecret, passphrase))
```

Requests are signed with HMAC-SHA256 over
`timestamp + method + requestPath(+ "?" + query) + body`, base64-encoded into the
`ACCESS-SIGN` header. To sign with an RSA key or an external signer instead,
pass `client.WithSignFn(fn)`.

Other options: `WithProxy` (http/https/socks5), `WithBaseURL`, `WithLocale`,
`WithDemoTrading`, `WithTimeOffset`, `WithLogger`, `WithHTTPClient`.

## Layout

```
bitget.go        entry point: NewUTAClient, NewUTAWebSocketClient
client/          core REST client, functional options, HMAC signer config, errors, WS client
request/         request builder, signing, Do[T] envelope decode, WebSocket subscribe
common/          constants, the global time.Time + decimal.Decimal JSON codec
pkg/log/         Logger interface
uta/             unified-account endpoints (one Service per endpoint) + WebSocket channels
cmd/bgraw/       dev tool: sign + dump any endpoint's raw response
```

### Endpoint modules (in `uta/`)

| Area | Files | Examples |
|------|-------|----------|
| Market data | `market*.go` | instruments, tickers, orderbook, candles, funding rate, open interest, risk reserve, discount rate, position tier, … |
| Account | `account*.go` | assets, info, settings, leverage, hold mode, fee rate, financial records, transfer, deposit, withdrawal, … |
| Trade | `trade_*.go` | place/modify/cancel order, batch, cancel-symbol, countdown-cancel, order/fill queries |
| Position | `position.go` | current/history positions, ADL rank, max open available |
| Strategy | `strategy.go` | trigger / TPSL plan orders |
| Copy / Earn / Tax | `copy.go` `earn.go` `tax.go` | elite copy trading, on-chain elite earn, tax records |
| Loans | `crypto_loan.go` `ins_loan.go` | crypto-backed loans, institutional loans |
| Broker / P2P / Sub-account | `broker.go` `p2p.go` `sub_account.go` | broker sub-accounts, P2P ads/orders, virtual sub-accounts |
| WebSocket | `ws_public.go` `ws_private.go` `ws_trade.go` | public: ticker, kline, orderbook, rpi-orderbook, trade, liquidation; private: account, position, order, fill, fast-fill, strategy-order, adl; order entry over WS: place/modify/cancel + batch |

## WebSocket

```go
ws := bitget.NewUTAWebSocketClient(
	client.WithWebSocketAuth(apiKey, apiSecret, passphrase), // private channels only
)

// Public ticker.
done, _, err := ws.NewSubscribeTickerService(uta.WsInstTypeUSDTFutures, "BTCUSDT").
	Do(ctx, func(p *request.WsPush[[]uta.WsTicker], err error) {
		if err != nil { return }
		fmt.Println(p.Action, p.Data[0].LastPrice)
	})
// ...
close(done) // unsubscribe + close

// Private account (auto login).
ws.NewSubscribeAccountService().Do(ctx, func(p *request.WsPush[[]uta.WsAccount], err error) {
	// p.Data[0].TotalEquity, p.Data[0].Coin, ...
})
```

Each `Do` returns `(done chan<- struct{}, stop <-chan struct{}, err error)`:
close `done` to unsubscribe; `stop` is closed when the reader exits. The client
sends Bitget's `ping`/`pong` keepalive automatically.

Orders can also be placed over a persistent, logged-in WebSocket connection (a
low-latency alternative to the REST trade endpoints):

```go
tc, _ := ws.DialTrade(ctx)       // connect + login
defer tc.Close()
price := decimal.RequireFromString("30000")
ack, err := tc.PlaceOrder(ctx, uta.CategorySpot, uta.WsNewOrder{
	Symbol: "BTCUSDT", Side: uta.SideBuy, OrderType: uta.OrderTypeLimit,
	Qty: decimal.RequireFromString("0.0001"), Price: &price,
})
// tc.ModifyOrder / tc.CancelOrder / tc.BatchPlaceOrders / ...
```

## Decimals and timestamps

All amounts/prices/rates are `decimal.Decimal`; all millisecond timestamps are
`time.Time`. A custom JSON codec (registered globally in `common/json.go`)
decodes Bitget's string-encoded numbers and timestamps — including the empty /
`"0"` / `"-1"` "not set" sentinels, which decode to the zero value — so you never
deal with raw strings.

## Testing

Tests run against the live API and read credentials from the environment; they
skip when unset:

```bash
export BITGET_API_KEY=...   BITGET_API_SECRET=...   BITGET_PASSPHRASE=...
export BITGET_PROXY=socks5://127.0.0.1:7890   # optional

go test ./uta/ -run TestAccountConfig -v            # one module at a time
BITGET_TEST_WRITE=1 go test ./uta/ -run TestOrder   # live order tests (tiny, reversible)
```

Every test hits the live API and verifies that the real response's JSON keys are
fully covered by the typed struct. A few notes:

- Run tests **per module** (`-run TestXxx`). Bitget rate-limits to ~1–20 req/s;
  running the whole suite at once can trip HTTP 429.
- If your API key has an **IP whitelist**, requests from a non-whitelisted IP
  return `40018 Invalid IP` — make sure your egress IP is whitelisted.
- Capability-gated reads (broker, copy-trading, P2P, loans) are skipped when the
  account lacks that capability (the endpoint + signing are still exercised).

State-changing tests (placing orders, changing settings) are gated behind
`BITGET_TEST_WRITE=1` and use minimal amounts on large-cap symbols. Destructive
endpoints (account downgrade, withdrawals, sub-account/broker creation) are
implemented but never executed by the suite.

The `cmd/bgraw` helper signs and dumps any endpoint's raw response, handy for
inspecting private payloads:

```bash
go run ./cmd/bgraw GET /api/v3/account/info
go run ./cmd/bgraw GET /api/v3/account/financial-records "coin=USDT&limit=5"
```

## License

See [LICENSE](LICENSE).
