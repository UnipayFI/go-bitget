package spot

import (
	"strconv"
	"testing"
	"time"

	"github.com/UnipayFI/go-bitget/classic/internal/apitest"
)

func TestSpotWallet(t *testing.T) {
	c := NewSpotClient(apitest.AuthOptions(t)...)
	_ = c.SyncServerTime(apitest.Ctx(t))
	ctx := apitest.Ctx(t)

	now := time.Now()
	startMs := strconv.FormatInt(now.Add(-30*24*time.Hour).UnixMilli(), 10)
	endMs := strconv.FormatInt(now.UnixMilli(), 10)

	// Get Deposit Address -- private GET.
	{
		const label = "GetDepositAddress"
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/wallet/deposit-address",
			map[string]string{"coin": "USDT", "chain": "trc20"}, true)
		apitest.AssertCovers(t, label, raw, DepositAddress{})
	}

	// Get Deposit Records -- private GET.
	{
		const label = "GetDepositRecords"
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/wallet/deposit-records",
			map[string]string{"startTime": startMs, "endTime": endMs, "limit": "20"}, true)
		apitest.AssertCovers(t, label, raw, []DepositRecord{{}})
	}

	// Get Withdrawal Records -- private GET.
	{
		const label = "GetWithdrawalRecords"
		raw := apitest.FetchRawGet(t, c, ctx, "/api/v2/spot/wallet/withdrawal-records",
			map[string]string{"startTime": startMs, "endTime": endMs, "limit": "20"}, true)
		apitest.AssertCovers(t, label, raw, []WithdrawalRecord{{}})
	}
}
