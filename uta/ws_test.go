package uta

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/request"
)

func testPublicWsClient() *UTAWebSocketClient {
	var opts []client.WebSocketOptions
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewUTAWebSocketClient(opts...)
}

func testWsClient(t *testing.T) *UTAWebSocketClient {
	t.Helper()
	apiKey := os.Getenv("BITGET_API_KEY")
	apiSecret := os.Getenv("BITGET_API_SECRET")
	passphrase := os.Getenv("BITGET_PASSPHRASE")
	if apiKey == "" || apiSecret == "" || passphrase == "" {
		t.Skip("BITGET_API_KEY/SECRET/PASSPHRASE not set; skipping private ws test")
	}
	opts := []client.WebSocketOptions{client.WithWebSocketAuth(apiKey, apiSecret, passphrase)}
	if proxy := os.Getenv("BITGET_PROXY"); proxy != "" {
		opts = append(opts, client.WithWebSocketProxy(proxy))
	}
	return NewUTAWebSocketClient(opts...)
}

func TestWsPublicTicker(t *testing.T) {
	ws := testPublicWsClient()
	msgC := make(chan *request.WsPush[[]WsTicker], 4)
	done, _, err := ws.NewSubscribeTickerService(WsInstTypeUSDTFutures, "BTCUSDT").
		Do(context.Background(), func(p *request.WsPush[[]WsTicker], err error) {
			if err != nil {
				t.Logf("ws ticker err: %v", err)
				return
			}
			select {
			case msgC <- p:
			default:
			}
		})
	if err != nil {
		t.Fatalf("subscribe ticker: %v", err)
	}
	defer close(done)

	select {
	case p := <-msgC:
		if len(p.Data) == 0 {
			t.Fatal("ticker push had empty data")
		}
		t.Logf("ticker action=%s arg=%+v first=%+v", p.Action, p.Arg, p.Data[0])
		if p.Data[0].LastPrice.IsZero() {
			t.Error("ticker lastPrice is zero")
		}
		// USDT-FUTURES perp ticker carries nextFundingTime — guards against
		// WsTicker.NextFundingTime regressing (the live frame includes it).
		if p.Data[0].NextFundingTime.IsZero() {
			t.Error("perp ticker nextFundingTime is zero (WsTicker.NextFundingTime missing?)")
		}
	case <-time.After(15 * time.Second):
		t.Fatal("no ticker message within 15s")
	}
}

func TestWsPublicOrderBook(t *testing.T) {
	ws := testPublicWsClient()
	msgC := make(chan *request.WsPush[[]WsOrderBook], 4)
	done, _, err := ws.NewSubscribeOrderBookService(WsInstTypeUSDTFutures, "BTCUSDT").SetDepth("books5").
		Do(context.Background(), func(p *request.WsPush[[]WsOrderBook], err error) {
			if err != nil {
				t.Logf("ws book err: %v", err)
				return
			}
			select {
			case msgC <- p:
			default:
			}
		})
	if err != nil {
		t.Fatalf("subscribe books: %v", err)
	}
	defer close(done)

	select {
	case p := <-msgC:
		if len(p.Data) == 0 {
			t.Fatal("book push had empty data")
		}
		t.Logf("book action=%s asks=%d bids=%d seq=%d ts=%s",
			p.Action, len(p.Data[0].Asks), len(p.Data[0].Bids), p.Data[0].Seq, p.Data[0].Ts)
	case <-time.After(15 * time.Second):
		t.Fatal("no book message within 15s")
	}
}

func TestWsPublicRPIOrderBook(t *testing.T) {
	ws := testPublicWsClient()
	msgC := make(chan *request.WsPush[[]WsRPIOrderBook], 4)
	errC := make(chan error, 4)
	// RPI depth is only published for RPI-enabled spot symbols; use one from the
	// RPI symbol list.
	done, _, err := ws.NewSubscribeRPIOrderBookService(WsInstTypeSpot, "BTCUSDT").SetDepth("rpi-books5").
		Do(context.Background(), func(p *request.WsPush[[]WsRPIOrderBook], err error) {
			if err != nil {
				select {
				case errC <- err:
				default:
				}
				return
			}
			select {
			case msgC <- p:
			default:
			}
		})
	if err != nil {
		t.Fatalf("subscribe rpi-books: %v", err)
	}
	defer close(done)

	select {
	case e := <-errC:
		t.Skipf("RPI depth not available for this symbol: %v", e)
	case p := <-msgC:
		if len(p.Data) == 0 {
			t.Fatal("rpi-book push had empty data")
		}
		t.Logf("rpi-book action=%s asks=%d bids=%d seq=%d", p.Action, len(p.Data[0].Asks), len(p.Data[0].Bids), p.Data[0].Seq)
	case <-time.After(12 * time.Second):
		t.Skip("no rpi-book message (symbol may not be RPI-enabled)")
	}
}

func TestWsPublicLiquidation(t *testing.T) {
	ws := testPublicWsClient()
	errC := make(chan error, 4)
	dataC := make(chan *request.WsPush[[]WsLiquidation], 4)
	done, _, err := ws.NewSubscribeLiquidationService(WsInstTypeUSDTFutures, "BTCUSDT").
		Do(context.Background(), func(p *request.WsPush[[]WsLiquidation], err error) {
			if err != nil {
				select {
				case errC <- err:
				default:
				}
				return
			}
			select {
			case dataC <- p:
			default:
			}
		})
	if err != nil {
		t.Fatalf("subscribe liquidation: %v", err)
	}
	defer close(done)

	// Liquidations are event-driven (pushed only when they occur), so absence of
	// data is fine; only an error frame indicates a real problem.
	select {
	case e := <-errC:
		t.Fatalf("liquidation subscribe error: %v", e)
	case p := <-dataC:
		t.Logf("liquidation event: %d record(s)", len(p.Data))
	case <-time.After(8 * time.Second):
		t.Log("no liquidation event in 8s (event-driven); subscription accepted OK")
	}
}

// TestWsPrivateChannels confirms each private data channel logs in and accepts
// the subscription (no error frame). These channels are event-driven, so they
// push nothing on subscribe; the account channel (tested separately) is the one
// that snapshots immediately.
func TestWsPrivateChannels(t *testing.T) {
	ws := testWsClient(t)
	type sub struct {
		name string
		do   func(cb func([]byte, error)) (chan<- struct{}, <-chan struct{}, error)
	}
	raw := func(private bool, arg request.WsArg) func(func([]byte, error)) (chan<- struct{}, <-chan struct{}, error) {
		return func(cb func([]byte, error)) (chan<- struct{}, <-chan struct{}, error) {
			return ws.Subscribe(context.Background(), private, arg, cb)
		}
	}
	subs := []sub{
		{"order", raw(true, request.WsArg{InstType: "UTA", Topic: "order"})},
		{"fill", raw(true, request.WsArg{InstType: "UTA", Topic: "fill"})},
		{"position", raw(true, request.WsArg{InstType: "UTA", Topic: "position"})},
		{"strategy-order", raw(true, request.WsArg{InstType: "UTA", Topic: "strategy-order"})},
		{"fast-fill", raw(true, request.WsArg{InstType: "UTA", Topic: "fast-fill"})},
		{"adl-notification", raw(true, request.WsArg{InstType: "UTA", Topic: "adl-notification"})},
	}
	for _, s := range subs {
		errC := make(chan error, 1)
		done, _, err := s.do(func(_ []byte, e error) {
			if e != nil {
				select {
				case errC <- e:
				default:
				}
			}
		})
		if err != nil {
			t.Fatalf("%s: dial/login: %v", s.name, err)
		}
		select {
		case e := <-errC:
			close(done)
			t.Fatalf("%s: subscribe rejected: %v", s.name, e)
		case <-time.After(2 * time.Second):
			t.Logf("%s: subscription accepted (login OK, no error)", s.name)
		}
		close(done)
	}
}

func TestWsPrivateAccount(t *testing.T) {
	ws := testWsClient(t)
	msgC := make(chan *request.WsPush[[]WsAccount], 4)
	done, _, err := ws.NewSubscribeAccountService().
		Do(context.Background(), func(p *request.WsPush[[]WsAccount], err error) {
			if err != nil {
				t.Logf("ws account err: %v", err)
				return
			}
			select {
			case msgC <- p:
			default:
			}
		})
	if err != nil {
		t.Fatalf("subscribe account: %v", err)
	}
	defer close(done)

	select {
	case p := <-msgC:
		if len(p.Data) == 0 {
			t.Fatal("account push had empty data")
		}
		t.Logf("account action=%s totalEquity=%s coins=%d",
			p.Action, p.Data[0].TotalEquity, len(p.Data[0].Coin))
	case <-time.After(15 * time.Second):
		t.Fatal("no account message within 15s")
	}
}
