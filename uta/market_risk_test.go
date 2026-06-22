package uta

import "testing"

func TestMarketRisk(t *testing.T) {
	c := testPublicClient()
	cx := ctx(t)

	// Risk Reserve (Daily).
	rr, err := c.NewGetRiskReserveService(CategoryUSDTFutures, "BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("risk-reserve: %v", err)
	}
	rrParams := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}
	rrRaw := fetchRawGet(t, c, cx, "/api/v3/market/risk-reserve", rrParams, false)
	assertCovers(t, "market/risk-reserve", rrRaw, rr)

	// Risk Reserve (Hourly).
	rrh, err := c.NewGetRiskReserveHourService(CategoryUSDTFutures, "BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("risk-reserve-hour: %v", err)
	}
	rrhRaw := fetchRawGet(t, c, cx, "/api/v3/market/risk-reserve-hour", rrParams, false)
	assertCovers(t, "market/risk-reserve-hour", rrhRaw, rrh)

	// Risk Reserve All.
	rra, err := c.NewGetRiskReserveAllService(CategoryUSDTFutures).Do(cx)
	if err != nil {
		t.Fatalf("risk-reserve-all: %v", err)
	}
	rraParams := map[string]string{"category": string(CategoryUSDTFutures)}
	rraRaw := fetchRawGet(t, c, cx, "/api/v3/market/risk-reserve-all", rraParams, false)
	assertCovers(t, "market/risk-reserve-all", rraRaw, rra)

	// Discount Rate.
	dr, err := c.NewGetDiscountRateService().Do(cx)
	if err != nil {
		t.Fatalf("discount-rate: %v", err)
	}
	drRaw := fetchRawGet(t, c, cx, "/api/v3/market/discount-rate", nil, false)
	assertCovers(t, "market/discount-rate", drRaw, dr)

	// Margin Loans.
	ml, err := c.NewGetMarginLoansService("USDT").Do(cx)
	if err != nil {
		t.Fatalf("margin-loans: %v", err)
	}
	mlParams := map[string]string{"coin": "USDT"}
	mlRaw := fetchRawGet(t, c, cx, "/api/v3/market/margin-loans", mlParams, false)
	assertCovers(t, "market/margin-loans", mlRaw, ml)

	// Position Tier.
	pt, err := c.NewGetPositionTierService(CategoryUSDTFutures).SetSymbol("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("position-tier: %v", err)
	}
	ptParams := map[string]string{"category": string(CategoryUSDTFutures), "symbol": "BTCUSDT"}
	ptRaw := fetchRawGet(t, c, cx, "/api/v3/market/position-tier", ptParams, false)
	assertCovers(t, "market/position-tier", ptRaw, pt)

	// Index Components.
	ic, err := c.NewGetIndexComponentsService("BTCUSDT").Do(cx)
	if err != nil {
		t.Fatalf("index-components: %v", err)
	}
	icParams := map[string]string{"symbol": "BTCUSDT"}
	icRaw := fetchRawGet(t, c, cx, "/api/v3/market/index-components", icParams, false)
	assertCovers(t, "market/index-components", icRaw, ic)

	// Proof Of Reserves.
	por, err := c.NewGetProofOfReservesService().Do(cx)
	if err != nil {
		t.Fatalf("proof-of-reserves: %v", err)
	}
	porRaw := fetchRawGet(t, c, cx, "/api/v3/market/proof-of-reserves", nil, false)
	assertCovers(t, "market/proof-of-reserves", porRaw, por)
}
