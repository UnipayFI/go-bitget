package ws

import (
	"context"
	"testing"
	"time"
)

// awaitOK waits up to 12s for a subscription callback to report either a decode
// error or a successful data push (signalled via res).
func awaitOK(t *testing.T, label string, res <-chan error) {
	t.Helper()
	select {
	case err := <-res:
		if err != nil {
			t.Errorf("%s: %v", label, err)
			return
		}
		t.Logf("%s: OK (push received, typed decode clean)", label)
	case <-time.After(12 * time.Second):
		t.Errorf("%s: timed out waiting for a push", label)
	}
}

// TestWsPublicChannels live-subscribes to every public classic v2 channel and
// verifies a push arrives AND decodes cleanly into the typed struct (a wrong
// field type surfaces as a decode error in the callback).
func TestWsPublicChannels(t *testing.T) {
	c := NewWebSocketClient(publicWSOptions()...)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	const (
		spot = "BTCUSDT"
		perp = "BTCUSDT"
	)

	t.Run("spot-ticker", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeSpotTickerService(spot).Do(ctx, func(p *WsPush[[]SpotWsTicker], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "spot-ticker", res)
	})

	t.Run("spot-candle", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeSpotCandleService(spot, "1m").Do(ctx, func(p *WsPush[[]SpotWsCandle], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "spot-candle", res)
	})

	t.Run("spot-trade", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeSpotTradeService(spot).Do(ctx, func(p *WsPush[[]SpotWsTrade], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "spot-trade", res)
	})

	t.Run("spot-orderbook", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeSpotOrderBookService(spot).Do(ctx, func(p *WsPush[[]SpotWsOrderBook], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "spot-orderbook", res)
	})

	t.Run("mix-ticker", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeMixTickerService(InstTypeUSDTFutures, perp).Do(ctx, func(p *WsPush[[]MixWsTicker], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "mix-ticker", res)
	})

	t.Run("mix-candle", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeMixCandleService(InstTypeUSDTFutures, perp, "1m").Do(ctx, func(p *WsPush[[]MixWsCandle], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "mix-candle", res)
	})

	t.Run("mix-orderbook", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeMixOrderBookService(InstTypeUSDTFutures, perp).Do(ctx, func(p *WsPush[[]MixWsOrderBook], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "mix-orderbook", res)
	})

	t.Run("mix-trade", func(t *testing.T) {
		res := make(chan error, 1)
		done, _, err := c.NewSubscribeMixTradeService(InstTypeUSDTFutures, perp).Do(ctx, func(p *WsPush[[]MixWsTrade], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "mix-trade", res)
	})

	t.Run("margin-cross-index-price", func(t *testing.T) {
		res := make(chan error, 1)
		// Cross-margin index-price uses instId "default" (all symbols), not a pair.
		done, _, err := c.NewSubscribeMarginCrossIndexPriceService("default").Do(ctx, func(p *WsPush[[]MarginCrossIndexPrice], e error) {
			if e != nil {
				trySend(res, e)
			} else if len(p.Data) > 0 {
				trySend(res, nil)
			}
		})
		if err != nil {
			t.Fatalf("subscribe: %v", err)
		}
		defer close(done)
		awaitOK(t, "margin-cross-index-price", res)
	})
}

func trySend(ch chan<- error, v error) {
	select {
	case ch <- v:
	default:
	}
}
