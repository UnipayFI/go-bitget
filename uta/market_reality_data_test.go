package uta

import "testing"

// realityCode is a listing that is present in every Reality basic-data
// endpoint, so one code exercises all nine.
const realityCode = "AAPL"

func TestMarketRealityData(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	params := map[string]string{"code": realityCode}

	// Company overview.
	{
		overview, err := c.NewGetCompanyOverviewService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("company-overview: %v", err)
		}
		t.Logf("overview: %+v", overview)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/company-overview", params, false)
		assertCovers(t, "reality/market/company-overview", raw, overview)
	}

	// Valuation indicators.
	{
		valuation, err := c.NewGetValuationIndicatorsService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("valuation-indicators: %v", err)
		}
		t.Logf("valuation: %+v", valuation)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/valuation-indicators", params, false)
		assertCovers(t, "reality/market/valuation-indicators", raw, valuation)
	}

	// Earnings forecast.
	{
		forecast, err := c.NewGetEarningsForecastService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("earnings-forecast: %v", err)
		}
		t.Logf("forecast: %+v", forecast)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/earnings-forecast", params, false)
		assertCovers(t, "reality/market/earnings-forecast", raw, forecast)
	}

	// Suspension and resumption info.
	{
		info, err := c.NewGetSuspensionResumptionInfoService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("suspension-resumption-info: %v", err)
		}
		t.Logf("suspension: %+v", info)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/suspension-resumption-info", params, false)
		assertCovers(t, "reality/market/suspension-resumption-info", raw, info)
	}

	// Dividends.
	{
		dividends, err := c.NewGetDividendsService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("dividends: %v", err)
		}
		if len(dividends.List) == 0 {
			t.Fatal("no dividend records returned")
		}
		t.Logf("dividend: %+v cursor=%s", dividends.List[0], dividends.Cursor)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/dividends", params, false)
		assertCovers(t, "reality/market/dividends", raw, dividends)
	}

	// Share capital change.
	{
		capital, err := c.NewGetShareCapitalChangeService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("share-capital-change: %v", err)
		}
		t.Logf("capital: %+v", capital)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/share-capital-change", params, false)
		assertCovers(t, "reality/market/share-capital-change", raw, capital)
	}

	// Insider trades.
	{
		trades, err := c.NewGetInnerTradesService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("inner-trades: %v", err)
		}
		if len(trades.List) == 0 {
			t.Fatal("no insider trades returned")
		}
		t.Logf("inner trade: %+v cursor=%s", trades.List[0], trades.Cursor)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/inner-trades", params, false)
		assertCovers(t, "reality/market/inner-trades", raw, trades)
	}

	// Executive shareholdings.
	{
		holdings, err := c.NewGetExecutiveShareholdingsService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("executive-shareholdings: %v", err)
		}
		if len(holdings.List) == 0 {
			t.Fatal("no executive shareholdings returned")
		}
		t.Logf("executive holding: %+v cursor=%s", holdings.List[0], holdings.Cursor)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/executive-shareholdings", params, false)
		assertCovers(t, "reality/market/executive-shareholdings", raw, holdings)
	}

	// Shareholder detail.
	{
		detail, err := c.NewGetShareholdDetailService(realityCode).Do(cx)
		if err != nil {
			t.Fatalf("sharehold-detail: %v", err)
		}
		if len(detail.List) == 0 {
			t.Fatal("no shareholders returned")
		}
		t.Logf("shareholder: %+v cursor=%s", detail.List[0], detail.Cursor)
		raw := fetchRawGet(t, c, cx, "/api/v3/reality/market/sharehold-detail", params, false)
		assertCovers(t, "reality/market/sharehold-detail", raw, detail)
	}
}
