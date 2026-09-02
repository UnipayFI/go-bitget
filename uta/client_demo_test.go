package uta

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/UnipayFI/go-bitget/client"
)

func TestSyncServerTimeIncludesDemoHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("paptrading"); got != "1" {
			t.Errorf("paptrading header = %q, want 1", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code":"00000","msg":"success","data":{"serverTime":"1788275950592"}}`))
	}))
	defer server.Close()
	c := NewUTAClient(client.WithBaseURL(server.URL), client.WithDemoTrading(true))
	if err := c.SyncServerTime(context.Background()); err != nil {
		t.Fatal(err)
	}
}
