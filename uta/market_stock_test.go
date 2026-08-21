package uta

import "testing"

func TestMarketStock(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// Split records.
	{
		records, err := c.NewGetSplitRecordsService().Do(cx)
		if err != nil {
			t.Fatalf("split-records: %v", err)
		}
		if len(records) == 0 {
			t.Fatal("no split records returned")
		}
		t.Logf("split record: %+v", records[0])
		raw := fetchRawGet(t, c, cx, "/api/v3/market/split-records", nil, false)
		assertCovers(t, "market/split-records", raw, records)
	}

	// Stock info.
	{
		params := map[string]string{"symbol": "RAAPLUSDT"}
		infos, err := c.NewGetStockInfoService().SetSymbol("RAAPLUSDT").Do(cx)
		if err != nil {
			t.Fatalf("stock-info: %v", err)
		}
		if len(infos) == 0 {
			t.Fatal("no stock info returned")
		}
		t.Logf("stock info: %+v", infos[0])
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/stock-info", params, false)
		assertCovers(t, "reality/market/stock-info", raw, infos)
	}

	// Market states.
	{
		states, err := c.NewGetMarketStatesService().Do(cx)
		if err != nil {
			t.Fatalf("states: %v", err)
		}
		if len(states.StateList) == 0 {
			t.Fatal("no trading sessions returned")
		}
		t.Logf("states: %+v", states)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/states", nil, false)
		assertCovers(t, "reality/market/states", raw, states)
	}

	// Market calendar.
	{
		calendar, err := c.NewGetMarketCalendarService().Do(cx)
		if err != nil {
			t.Fatalf("calendar: %v", err)
		}
		if len(calendar.RegularConfig) == 0 {
			t.Fatal("no regular closures returned")
		}
		t.Logf("calendar: %+v", calendar)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/calendar", nil, false)
		assertCovers(t, "reality/market/calendar", raw, calendar)
	}
}
