package uta

import (
	"os"
	"testing"

	"github.com/shopspring/decimal"
)

// TestGridValidate exercises the two grid validate endpoints. They are POSTs
// but create nothing — they only rule on a configuration — so they are safe to
// run against a live key without the BITGET_TEST_WRITE gate.
func TestGridValidate(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Spot/futures grid.
	{
		investment := []GridInvestment{{Coin: "USDT", Amount: decimal.NewFromInt(1000)}}
		verdict, err := c.NewValidateGridParamsService(CategoryUSDTFutures, "BTCUSDT",
			decimal.NewFromInt(50000), decimal.NewFromInt(90000), "20", GridOrderModeArithmetic, investment).
			SetGridType(GridTypeLong).
			SetLeverage("2").
			Do(cx)
		if err != nil {
			t.Fatalf("grid/validate: %v", err)
		}
		t.Logf("validate: %+v", verdict)
		raw := fetchRawPost(t, c, cx, "/api/v3/trade/grid/validate", map[string]any{
			"category":         string(CategoryUSDTFutures),
			"symbol":           "BTCUSDT",
			"minPrice":         "50000",
			"maxPrice":         "90000",
			"gridNum":          "20",
			"gridOrderMode":    string(GridOrderModeArithmetic),
			"investmentAmount": investment,
			"gridType":         string(GridTypeLong),
			"leverage":         "2",
		}, true)
		assertCovers(t, "trade/grid/validate", raw, verdict)
	}

	// Neutral futures grid.
	{
		verdict, err := c.NewValidateNeutralGridParamsService(CategoryUSDTFutures, "BTCUSDT",
			decimal.NewFromInt(50000), decimal.NewFromInt(90000), "20", GridOrderModeArithmetic).
			Do(cx)
		if err != nil {
			t.Fatalf("grid/validate-neutral: %v", err)
		}
		t.Logf("validate-neutral: %+v", verdict)
		raw := fetchRawPost(t, c, cx, "/api/v3/trade/grid/validate-neutral", map[string]any{
			"category":      string(CategoryUSDTFutures),
			"symbol":        "BTCUSDT",
			"minPrice":      "50000",
			"maxPrice":      "90000",
			"gridNum":       "20",
			"gridOrderMode": string(GridOrderModeArithmetic),
		}, true)
		assertCovers(t, "trade/grid/validate-neutral", raw, verdict)
	}
}

// TestGridBotDetail reads back a running grid bot. Bitget has no "list bots"
// endpoint, so the bot ID has to come from the environment: set
// BITGET_GRID_BOT_ID for a spot/futures grid, BITGET_NEUTRAL_GRID_BOT_ID for a
// neutral one.
func TestGridBotDetail(t *testing.T) {
	botID := os.Getenv("BITGET_GRID_BOT_ID")
	neutralBotID := os.Getenv("BITGET_NEUTRAL_GRID_BOT_ID")
	if botID == "" && neutralBotID == "" {
		t.Skip("BITGET_GRID_BOT_ID/BITGET_NEUTRAL_GRID_BOT_ID not set; skipping grid detail test")
	}
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	if botID != "" {
		detail, err := c.NewGetGridBotDetailService(botID).Do(cx)
		if err != nil {
			t.Fatalf("grid/bot-detail: %v", err)
		}
		t.Logf("bot detail: %+v", detail)
		raw := fetchRawGet(t, c, cx, "/api/v3/trade/grid/bot-detail", map[string]string{"botId": botID}, true)
		assertCovers(t, "trade/grid/bot-detail", raw, detail)

		orders, err := c.NewGetGridBotOrderDetailsService(botID).Do(cx)
		if err != nil {
			t.Fatalf("grid/list-details: %v", err)
		}
		t.Logf("orders: %d buy, %d sell", len(orders.BuyOrderList), len(orders.SellOrderList))
		raw = fetchRawGet(t, c, cx, "/api/v3/trade/grid/list-details", map[string]string{"botId": botID}, true)
		assertCovers(t, "trade/grid/list-details", raw, orders)
	}

	if neutralBotID != "" {
		detail, err := c.NewGetNeutralGridBotDetailService(neutralBotID).Do(cx)
		if err != nil {
			t.Fatalf("grid/neutral-bot-detail: %v", err)
		}
		t.Logf("neutral bot detail: %+v", detail)
		raw := fetchRawGet(t, c, cx, "/api/v3/trade/grid/neutral-bot-detail", map[string]string{"botId": neutralBotID}, true)
		assertCovers(t, "trade/grid/neutral-bot-detail", raw, detail)

		orders, err := c.NewGetNeutralGridBotOrderDetailsService(neutralBotID).Do(cx)
		if err != nil {
			t.Fatalf("grid/neutral-list-details: %v", err)
		}
		t.Logf("neutral orders: %d buy, %d sell", len(orders.BuyOrderList), len(orders.SellOrderList))
		raw = fetchRawGet(t, c, cx, "/api/v3/trade/grid/neutral-list-details", map[string]string{"botId": neutralBotID}, true)
		assertCovers(t, "trade/grid/neutral-list-details", raw, orders)
	}
}
