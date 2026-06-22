package uta

import "testing"

func TestInsLoan(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	// Risk units (no params; parent account only).
	ru, err := c.NewGetInsLoanRiskUnitService().Do(cx)
	if err != nil {
		t.Fatalf("risk unit: %v", err)
	}
	t.Logf("risk units: %d", len(ru.RiskUnitId))
	raw := fetchRawGet(t, c, cx, "/api/v3/ins-loan/risk-unit", nil, true)
	assertCovers(t, "ins-loan/risk-unit", raw, ru)

	// Repayment orders (no required params).
	if repaid, err := c.NewGetInsLoanRepaidHistoryService().Do(cx); err != nil {
		if !tolerable(t, "ins-loan/repaid-history", err) {
			t.Fatalf("repaid history: %v", err)
		}
	} else {
		t.Logf("repaid orders: %d", len(repaid))
		raw = fetchRawGet(t, c, cx, "/api/v3/ins-loan/repaid-history", nil, true)
		assertCovers(t, "ins-loan/repaid-history", raw, repaid)
	}

	// Loan orders (no required params).
	if orders, err := c.NewGetInsLoanOrderService().Do(cx); err != nil {
		if !tolerable(t, "ins-loan/loan-order", err) {
			t.Fatalf("loan order: %v", err)
		}
	} else {
		t.Logf("loan orders: %d", len(orders))
		raw = fetchRawGet(t, c, cx, "/api/v3/ins-loan/loan-order", nil, true)
		assertCovers(t, "ins-loan/loan-order", raw, orders)
	}
}
