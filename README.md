# go-bitget

[![Go Reference](https://pkg.go.dev/badge/github.com/UnipayFI/go-bitget.svg)](https://pkg.go.dev/github.com/UnipayFI/go-bitget)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-00ADD8?logo=go)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A Go SDK for the [Bitget](https://www.bitget.com/api-doc/uta/intro) exchange, covering both account systems.

| Account system | API | Aligned to |
|---|---|---|
| **UTA** — Unified Trading Account | `/api/v3` REST + v3 WebSocket | [2026-07-23](https://www.bitget.com/api-doc/uta/changelog) |
| **Classic** — per-product account | `/api/v2` REST + v2 WebSocket | [2026-08-04](https://www.bitget.com/api-doc/classic/changelog) |

Response structs are reconciled against the live API (not just the docs), so endpoints stay in sync with the dates above.

## Install

```bash
go get github.com/UnipayFI/go-bitget@latest
```

## Highlights

- One signing/transport core shared by UTA (`uta`) and Classic (`classic/*`).
- Fluent per-endpoint API: `NewXxxService(...).SetFoo(...).Do(ctx)`.
- Amounts as `decimal.Decimal`, ms timestamps as `time.Time` — Bitget's string-encoded numbers and `""`/`"0"`/`"-1"` "not set" sentinels are decoded for you.
- Every endpoint is tested against the live API, diffing real JSON keys against the struct.

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
	_ = c.SyncServerTime(ctx) // align clock to avoid signature drift

	// Public market data (no auth).
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
	fmt.Println("orderId:", ref.OrderID)
}
```

## Authentication

Pass credentials from the Bitget API-management page:

```go
c := bitget.NewUTAClient(client.WithAuth(apiKey, apiSecret, passphrase))
```

Requests are signed with HMAC-SHA256 over `timestamp + method + requestPath(+ "?" + query) + body`, base64-encoded into the `ACCESS-SIGN` header. For an RSA key or external signer, pass `client.WithSignFn(fn)`.

Other options: `WithProxy` (http/https/socks5), `WithBaseURL`, `WithLocale`, `WithDemoTrading`, `WithTimeOffset`, `WithLogger`, `WithHTTPClient`.

## WebSocket

```go
ws := bitget.NewUTAWebSocketClient(
	client.WithWebSocketAuth(apiKey, apiSecret, passphrase), // private channels only
)

// Public ticker.
done, _, _ := ws.NewSubscribeTickerService(uta.WsInstTypeUSDTFutures, "BTCUSDT").
	Do(ctx, func(p *request.WsPush[[]uta.WsTicker], err error) {
		if err != nil {
			return
		}
		fmt.Println(p.Action, p.Data[0].LastPrice)
	})
close(done) // unsubscribe + close

// Private account (auto login).
ws.NewSubscribeAccountService().Do(ctx, func(p *request.WsPush[[]uta.WsAccount], err error) {
	// p.Data[0].TotalEquity, p.Data[0].Coin, ...
})
```

Each `Do` returns `(done chan<- struct{}, stop <-chan struct{}, err error)`: close `done` to unsubscribe; `stop` closes when the reader exits. Ping/pong keepalive is automatic.

Orders can also be placed over a persistent, logged-in connection — a low-latency alternative to the REST trade endpoints:

```go
tc, _ := ws.DialTrade(ctx) // connect + login
defer tc.Close()
price := decimal.RequireFromString("30000")
ack, _ := tc.PlaceOrder(ctx, uta.CategorySpot, uta.WsNewOrder{
	Symbol: "BTCUSDT", Side: uta.SideBuy, OrderType: uta.OrderTypeLimit,
	Qty: decimal.RequireFromString("0.0001"), Price: &price,
})
// tc.ModifyOrder / tc.CancelOrder / tc.BatchPlaceOrders / ...
```

## Classic account

Each product line is a separate package (so `PlaceOrder`, `GetTickers`, `Account`, … don't collide), with the same `NewXxxService(...).SetFoo(...).Do(ctx)` shape:

```go
sp := bitget.NewSpotClient(client.WithAuth(apiKey, apiSecret, passphrase))
tickers, _ := sp.NewGetTickersService().SetSymbol("BTCUSDT").Do(ctx)

mx := bitget.NewMixClient(client.WithAuth(apiKey, apiSecret, passphrase))
pos, _ := mx.NewGetAllPositionService(mix.ProductTypeUSDTFutures).Do(ctx)
```

## Packages

**UTA** (`uta/`)

| Area | Files |
|------|-------|
| Market data | `market*.go` — instruments, tickers, orderbook, candles, funding rate, open interest, … |
| Account | `account*.go` — assets, settings, leverage, fee rate, records, transfer, deposit, withdrawal, … |
| Trade | `trade_*.go` — place/modify/cancel, batch, cancel-symbol, countdown-cancel, queries |
| Position / Strategy | `position.go` `strategy.go` — positions, ADL rank, trigger & TPSL plans |
| Copy / Earn / Loans / Tax | `copy.go` `earn.go` `crypto_loan.go` `ins_loan.go` `tax.go` |
| Broker / P2P / Sub-account | `broker.go` `p2p.go` `sub_account.go` |
| WebSocket | `ws_public.go` `ws_private.go` `ws_trade.go` |

**Classic** (`classic/`)

| Package | Scope |
|---------|-------|
| `common` | server time, announcements, trade-rate, all-account balance, convert / BGB-convert, virtual sub-accounts, insights |
| `spot` | market, trade, plan (trigger) orders, account, wallet/transfer/deposit/withdrawal |
| `mix` | futures market/account/position/trade/plan (USDT-/USDC-/COIN-M) |
| `margin` | cross + isolated: assets, borrow/repay, orders, fills, records |
| `copy` | futures & spot copy-trading: trader / follower / broker |
| `earn` | savings, shark-fin, loan: subscribe/redeem/borrow/repay + records |
| `broker` `affiliate` `insloan` `tax` `p2p` | broker sub-accounts & api-keys, affiliate, institutional loans, tax, p2p merchant |
| `ws` | v2 WebSocket — public + private channels (spot/mix/margin) + order entry |

**Core**

| Package | Scope |
|---------|-------|
| `bitget.go` | entry point: `NewUTAClient` + `NewSpotClient`/`NewMixClient`/… + WS clients |
| `client/` `request/` | REST client, options, HMAC signer, envelope decode, WS subscribe |
| `common/` | constants, global `time.Time` + `decimal.Decimal` JSON codec |
| `cmd/bgraw/` | dev tool: sign + dump any endpoint's raw response |

## Testing

Tests hit the live API and read credentials from the environment, skipping when unset:

```bash
export BITGET_API_KEY=...  BITGET_API_SECRET=...  BITGET_PASSPHRASE=...
export BITGET_PROXY=socks5://127.0.0.1:7890   # optional

go test ./uta/ -run TestAccountConfig -v            # one module at a time
BITGET_TEST_WRITE=1 go test ./uta/ -run TestOrder   # live order tests (tiny, reversible)
```

- Run **per module** (`-run TestXxx`) — Bitget rate-limits to ~1–20 req/s, so the full suite can trip HTTP 429.
- An **IP whitelist** on the key returns `40018 Invalid IP` from non-whitelisted egress IPs.
- Capability-gated reads (broker, copy-trading, P2P, loans) are skipped when the account lacks the capability — signing is still exercised.
- State-changing tests are gated behind `BITGET_TEST_WRITE=1` (minimal amounts, large-cap symbols). Destructive endpoints (downgrade, withdrawals, sub-account/broker creation) are implemented but never executed.

The `cmd/bgraw` helper dumps any endpoint's raw signed response:

```bash
go run ./cmd/bgraw GET /api/v3/account/info
go run ./cmd/bgraw GET /api/v3/account/financial-records "coin=USDT&limit=5"
```

## License

[MIT](LICENSE)
