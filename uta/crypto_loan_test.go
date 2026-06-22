package uta

import (
	"strconv"
	"testing"
	"time"
)

func TestCryptoLoan(t *testing.T) {
	c := testClient(t)
	if err := c.SyncServerTime(ctx(t)); err != nil {
		t.Fatalf("sync time: %v", err)
	}
	cx := ctx(t)

	endTime := time.Now()
	startTime := endTime.Add(-30 * 24 * time.Hour)
	window := map[string]string{
		"startTime": strconv.FormatInt(startTime.UnixMilli(), 10),
		"endTime":   strconv.FormatInt(endTime.UnixMilli(), 10),
	}

	// Loan coins (object with loanInfos/pledgeInfos).
	coins, err := c.NewGetLoanCoinsService().Do(cx)
	if err != nil {
		t.Fatalf("loan coins: %v", err)
	}
	t.Logf("loan coins: loanInfos=%d pledgeInfos=%d", len(coins.LoanInfos), len(coins.PledgeInfos))
	raw := fetchRawGet(t, c, cx, "/api/v3/loan/coins", nil, true)
	assertCovers(t, "loan/coins", raw, coins)

	// Loan debts (object with pledgeInfos/loanInfos).
	if debts, err := c.NewGetLoanDebtsService().Do(cx); err != nil {
		if !tolerable(t, "loan/debts", err) {
			t.Fatalf("loan debts: %v", err)
		}
	} else {
		t.Logf("loan debts: pledgeInfos=%d loanInfos=%d", len(debts.PledgeInfos), len(debts.LoanInfos))
		raw = fetchRawGet(t, c, cx, "/api/v3/loan/debts", nil, true)
		assertCovers(t, "loan/debts", raw, debts)
	}

	// Ongoing borrows (array, no required params).
	ongoing, err := c.NewGetBorrowOngoingService().Do(cx)
	if err != nil {
		t.Fatalf("borrow ongoing: %v", err)
	}
	t.Logf("borrow ongoing: %d", len(ongoing))
	raw = fetchRawGet(t, c, cx, "/api/v3/loan/borrow-ongoing", nil, true)
	assertCovers(t, "loan/borrow-ongoing", raw, ongoing)

	// Borrow history (array within a time window).
	borrowHist, err := c.NewGetBorrowHistoryService(startTime, endTime).Do(cx)
	if err != nil {
		t.Fatalf("borrow history: %v", err)
	}
	t.Logf("borrow history: %d", len(borrowHist))
	raw = fetchRawGet(t, c, cx, "/api/v3/loan/borrow-history", window, true)
	assertCovers(t, "loan/borrow-history", raw, borrowHist)

	// Repay history (array within a time window).
	repayHist, err := c.NewGetRepayHistoryService(startTime, endTime).Do(cx)
	if err != nil {
		t.Fatalf("repay history: %v", err)
	}
	t.Logf("repay history: %d", len(repayHist))
	raw = fetchRawGet(t, c, cx, "/api/v3/loan/repay-history", window, true)
	assertCovers(t, "loan/repay-history", raw, repayHist)

	// Pledge rate history (array within a time window).
	pledgeHist, err := c.NewGetPledgeRateHistoryService(startTime, endTime).Do(cx)
	if err != nil {
		t.Fatalf("pledge rate history: %v", err)
	}
	t.Logf("pledge rate history: %d", len(pledgeHist))
	raw = fetchRawGet(t, c, cx, "/api/v3/loan/pledge-rate-history", window, true)
	assertCovers(t, "loan/pledge-rate-history", raw, pledgeHist)

	// Loan reduces / liquidations (array within a time window).
	reduces, err := c.NewGetLoanReducesService(startTime, endTime).Do(cx)
	if err != nil {
		t.Fatalf("loan reduces: %v", err)
	}
	t.Logf("loan reduces: %d", len(reduces))
	raw = fetchRawGet(t, c, cx, "/api/v3/loan/reduces", window, true)
	assertCovers(t, "loan/reduces", raw, reduces)
}
