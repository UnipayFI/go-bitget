package insloan

import (
	"testing"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

// Tolerable Bitget codes for an account that is not enrolled in the
// institutional-loan product (or whose key lacks the parent/child scope), plus
// the empty-data code. The endpoint + signing are still proven in these cases.
var insLoanOK = []string{"40099", "40034", "40014", "40068", "40054", "40037", "40029", "47001", "22001"}

func TestInsLoan(t *testing.T) {
	c := NewInsLoanClient(apitest.AuthOptions(t)...)
	if err := c.SyncServerTime(apitest.Ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := apitest.Ctx(t)

	// Risk units (no params; parent account only).
	if ru, err := c.NewGetRiskUnitService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "ins-loan/risk-unit", err, insLoanOK...) {
			t.Fatalf("risk unit: %v", err)
		}
	} else {
		t.Logf("risk units: %d", len(ru.RiskUnitId))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/ins-loan/risk-unit", nil, true)
		apitest.AssertCovers(t, "ins-loan/risk-unit", raw, ru)
	}

	// Loan orders (no required params).
	if orders, err := c.NewGetLoanOrderService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "ins-loan/loan-order", err, insLoanOK...) {
			t.Fatalf("loan order: %v", err)
		}
	} else {
		t.Logf("loan orders: %d", len(orders))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/ins-loan/loan-order", nil, true)
		apitest.AssertCovers(t, "ins-loan/loan-order", raw, orders)
	}

	// Repayment orders (no required params).
	if repaid, err := c.NewGetRepaidHistoryService().Do(cx); err != nil {
		if !apitest.Tolerable(t, "ins-loan/repaid-history", err, insLoanOK...) {
			t.Fatalf("repaid history: %v", err)
		}
	} else {
		t.Logf("repaid orders: %d", len(repaid))
		raw := apitest.FetchRawGet(t, c, cx, "/api/v2/spot/ins-loan/repaid-history", nil, true)
		apitest.AssertCovers(t, "ins-loan/repaid-history", raw, repaid)
	}
}
