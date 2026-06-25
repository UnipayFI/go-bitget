package ws

import (
	"context"
	"time"

	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
)

// InstType routes a classic v2 subscription to a product line. Bitget uses
// UPPERCASE instType on the v2 streams (unlike the v3 lowercase form).
type InstType string

const (
	InstTypeSpot        InstType = "SPOT"
	InstTypeUSDTFutures InstType = "USDT-FUTURES"
	InstTypeCoinFutures InstType = "COIN-FUTURES"
	InstTypeUSDCFutures InstType = "USDC-FUTURES"
	// InstTypeMargin is used by the cross/isolated margin channels.
	InstTypeMargin InstType = "MARGIN"
)

// WsArg identifies a classic v2 channel subscription: instType selects the
// product line, channel is the channel name, and instId narrows it to a symbol
// (coin is used by the few coin-scoped channels).
type WsArg struct {
	InstType string `json:"instType"`
	Channel  string `json:"channel"`
	InstID   string `json:"instId,omitempty"`
	Coin     string `json:"coin,omitempty"`
}

// InstAll is the instId sentinel that subscribes a private channel to every
// symbol rather than narrowing to one.
const InstAll = "default"

// WsPush is the envelope Bitget pushes for a v2 data event.
type WsPush[T any] struct {
	Action request.WsAction `json:"action"`
	Arg    WsArg            `json:"arg"`
	Data   T                `json:"data"`
	Ts     time.Time        `json:"ts"`
}

// WsHandler is invoked for every push (or error) on a subscription. The push's
// Data field is already decoded into the channel's typed slice.
type WsHandler[T any] func(*WsPush[[]T], error)

// Subscribe opens a dedicated connection to the public or private v2 gateway,
// logs in when private, subscribes to arg, and invokes cb for every data push
// (decoded into *WsPush[T]). It returns a done channel (close to stop) and a
// stop channel (closed when the reader exits).
func Subscribe[T any](ctx context.Context, c *WebSocketClient, private bool, arg WsArg, cb func(*WsPush[T], error)) (done chan<- struct{}, stop <-chan struct{}, err error) {
	return request.SubscribeRawArg(ctx, c, private, arg, func(message []byte, e error) {
		if e != nil {
			cb(nil, e)
			return
		}
		var push WsPush[T]
		if err := common.JSONUnmarshal(message, &push); err != nil {
			cb(nil, err)
			return
		}
		cb(&push, nil)
	})
}

// SubscribeRaw is the low-level escape hatch: it subscribes to an arbitrary
// channel and delivers each data push's raw bytes. Prefer the typed
// NewSubscribe* services; use this for channels the SDK does not yet wrap.
func (c *WebSocketClient) SubscribeRaw(ctx context.Context, private bool, arg WsArg, cb func([]byte, error)) (chan<- struct{}, <-chan struct{}, error) {
	return request.SubscribeRawArg(ctx, c, private, arg, cb)
}
