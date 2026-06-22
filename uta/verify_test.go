package uta

import (
	"context"
	"errors"
	"maps"
	"sort"
	"testing"

	"github.com/UnipayFI/go-bitget/client"
	"github.com/UnipayFI/go-bitget/common"
	"github.com/UnipayFI/go-bitget/request"
)

// tolerable reports whether err is an expected "this account lacks the
// capability or simply has no data" Bitget response rather than a code bug, so
// capability-gated read tests (broker / copy-trading / P2P / loans) can treat it
// as a pass: the request path and signing were correct, the account just isn't
// a broker, copy trader, P2P merchant, or borrower.
func tolerable(t *testing.T, label string, err error) bool {
	t.Helper()
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.Code {
		case "40054", // data is empty
			"40014", // missing (copy-trading) permission
			"40029", // account is not a Broker
			"60039": // P2P user does not exist
			t.Logf("%s: account lacks this capability/data (code=%s) — endpoint+signing OK", label, apiErr.Code)
			return true
		}
	}
	return false
}

// fetchRawGet returns the raw JSON of the "data" field for a GET endpoint, used
// to diff the real response shape against the typed structs.
func fetchRawGet(t *testing.T, c *UTAClient, ctx context.Context, path string, params map[string]string, sign bool) []byte {
	t.Helper()
	req := request.Get(ctx, c, path, params)
	if sign {
		req = req.WithSign()
	}
	raw, err := request.DoRawData(req)
	if err != nil {
		t.Fatalf("raw GET %s: %v", path, err)
	}
	return raw
}

// fetchRawPost mirrors fetchRawGet for POST endpoints.
func fetchRawPost(t *testing.T, c *UTAClient, ctx context.Context, path string, body map[string]any, sign bool) []byte {
	t.Helper()
	req := request.Post(ctx, c, path, body)
	if sign {
		req = req.WithSign()
	}
	raw, err := request.DoRawData(req)
	if err != nil {
		t.Fatalf("raw POST %s: %v", path, err)
	}
	return raw
}

// assertCovers checks that every JSON key present in the real response (raw) is
// also produced by marshaling the typed value. It compares key *sets* (not
// values), recursing into nested objects and the first element of arrays, so a
// missing struct field surfaces as an uncovered key path. This is the backbone
// of the "fields match the real API" guarantee.
func assertCovers(t *testing.T, label string, raw []byte, typed any) {
	t.Helper()
	var rawAny any
	if err := common.JSONUnmarshal(raw, &rawAny); err != nil {
		t.Fatalf("%s: unmarshal raw: %v", label, err)
	}
	typedBytes, err := common.JSONMarshal(typed)
	if err != nil {
		t.Fatalf("%s: marshal typed: %v", label, err)
	}
	var haveAny any
	if err := common.JSONUnmarshal(typedBytes, &haveAny); err != nil {
		t.Fatalf("%s: unmarshal typed: %v", label, err)
	}

	var missing []string
	diffKeys(rawAny, haveAny, "", &missing)
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Errorf("%s: %d field(s) in real response NOT captured by struct:\n  %v", label, len(missing), missing)
		return
	}
	t.Logf("%s: OK, all response keys covered by struct", label)
}

// diffKeys walks raw and records the paths of keys absent from have.
func diffKeys(raw, have any, path string, out *[]string) {
	switch r := raw.(type) {
	case map[string]any:
		h, ok := have.(map[string]any)
		if !ok {
			*out = append(*out, path+" (expected object)")
			return
		}
		for k, rv := range r {
			child := path + "/" + k
			hv, present := h[k]
			if !present {
				*out = append(*out, child)
				continue
			}
			diffKeys(rv, hv, child, out)
		}
	case []any:
		h, ok := have.([]any)
		if !ok || len(r) == 0 || len(h) == 0 {
			return
		}
		// Merge keys across all raw elements so optional fields present only on
		// some rows are still checked against the (single-shape) struct.
		merged := map[string]any{}
		for _, e := range r {
			if em, ok := e.(map[string]any); ok {
				maps.Copy(merged, em)
			}
		}
		if len(merged) > 0 {
			diffKeys(merged, h[0], path+"[]", out)
		} else {
			diffKeys(r[0], h[0], path+"[]", out)
		}
	}
}
