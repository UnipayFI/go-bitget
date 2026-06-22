package uta

import (
	"context"
	"testing"
	"time"
)

func TestTradingData(t *testing.T) {
	c := testPublicClient()
	// These endpoints are paced at ~1s each (rate limit 1/sec/IP); the ~22 calls
	// outlast the shared 20s ctx helper, so use a longer deadline here.
	cx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	spotSymbol := "BTCUSDT"
	futuresSymbol := "BTCUSDT"

	// These endpoints are rate-limited to 1 request/second/IP and each section
	// issues two requests (Service.Do + raw fetch), so pace the calls.
	pause := func() { time.Sleep(1100 * time.Millisecond) }

	// Spot whale net flow.
	pause()
	whale, err := c.NewGetSpotWhaleNetFlowService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("spot whale net flow: %v", err)
	}
	if len(whale) == 0 {
		t.Fatal("no spot whale net flow returned")
	}
	t.Logf("whale first: %+v", whale[0])
	pause()
	raw := fetchRawGet(t, c, cx, "/api/v3/market/spot-whale-flow", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/spot-whale-flow", raw, whale)

	// Spot fund flow (single object).
	pause()
	fund, err := c.NewGetSpotFundFlowService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("spot fund flow: %v", err)
	}
	t.Logf("fund: %+v", fund)
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/spot-fund-flow", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/spot-fund-flow", raw, fund)

	// Spot net flow.
	pause()
	netFlow, err := c.NewGetSpotNetFlowService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("spot net flow: %v", err)
	}
	if len(netFlow) == 0 {
		t.Fatal("no spot net flow returned")
	}
	t.Logf("net flow first: %+v", netFlow[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/spot-net-flow", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/spot-net-flow", raw, netFlow)

	// Margin long/short ratio.
	pause()
	marginLS, err := c.NewGetMarginLongShortService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("margin long short: %v", err)
	}
	if len(marginLS) == 0 {
		t.Fatal("no margin long short returned")
	}
	t.Logf("margin long short first: %+v", marginLS[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/margin-long-short", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/margin-long-short", raw, marginLS)

	// Margin loan growth.
	pause()
	loanGrowth, err := c.NewGetMarginLoanGrowthService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("margin loan growth: %v", err)
	}
	if len(loanGrowth) == 0 {
		t.Fatal("no margin loan growth returned")
	}
	t.Logf("loan growth first: %+v", loanGrowth[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/margin-loan-growth", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/margin-loan-growth", raw, loanGrowth)

	// Margin isolated borrow.
	pause()
	isoBorrow, err := c.NewGetMarginIsolatedBorrowService(spotSymbol).Do(cx)
	if err != nil {
		t.Fatalf("margin isolated borrow: %v", err)
	}
	if len(isoBorrow) == 0 {
		t.Fatal("no margin isolated borrow returned")
	}
	t.Logf("isolated borrow first: %+v", isoBorrow[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/margin-isolated-borrow", map[string]string{"symbol": spotSymbol}, false)
	assertCovers(t, "market/margin-isolated-borrow", raw, isoBorrow)

	// Futures active buy/sell volume.
	pause()
	activeBS, err := c.NewGetFuturesActiveBuySellService(futuresSymbol).Do(cx)
	if err != nil {
		t.Fatalf("futures active buy sell: %v", err)
	}
	if len(activeBS) == 0 {
		t.Fatal("no futures active buy sell returned")
	}
	t.Logf("active buy sell first: %+v", activeBS[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/futures-active-buy-sell", map[string]string{"symbol": futuresSymbol}, false)
	assertCovers(t, "market/futures-active-buy-sell", raw, activeBS)

	// Futures long/short ratio.
	pause()
	futLS, err := c.NewGetFuturesLongShortService(futuresSymbol).Do(cx)
	if err != nil {
		t.Fatalf("futures long short: %v", err)
	}
	if len(futLS) == 0 {
		t.Fatal("no futures long short returned")
	}
	t.Logf("futures long short first: %+v", futLS[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/futures-long-short", map[string]string{"symbol": futuresSymbol}, false)
	assertCovers(t, "market/futures-long-short", raw, futLS)

	// Futures position long/short ratio.
	pause()
	posLS, err := c.NewGetFuturesPositionLongShortService(futuresSymbol).Do(cx)
	if err != nil {
		t.Fatalf("futures position long short: %v", err)
	}
	if len(posLS) == 0 {
		t.Fatal("no futures position long short returned")
	}
	t.Logf("position long short first: %+v", posLS[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/futures-position-long-short", map[string]string{"symbol": futuresSymbol}, false)
	assertCovers(t, "market/futures-position-long-short", raw, posLS)

	// Futures account long/short ratio.
	pause()
	acctLS, err := c.NewGetFuturesAccountLongShortService(futuresSymbol).Do(cx)
	if err != nil {
		t.Fatalf("futures account long short: %v", err)
	}
	if len(acctLS) == 0 {
		t.Fatal("no futures account long short returned")
	}
	t.Logf("account long short first: %+v", acctLS[0])
	pause()
	raw = fetchRawGet(t, c, cx, "/api/v3/market/futures-account-long-short", map[string]string{"symbol": futuresSymbol}, false)
	assertCovers(t, "market/futures-account-long-short", raw, acctLS)
}
