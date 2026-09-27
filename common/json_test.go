package common

import (
	"encoding/json/v2"
	"fmt"
	"strings"
	"testing"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

func TestJSONSupport(t *testing.T) {
	if errJSONSupport != nil {
		t.Fatal(errJSONSupport)
	}
}

// TestUnixFormats pins each unix unit, and checks the units decode side by
// side without one being scaled into another.
func TestUnixFormats(t *testing.T) {
	var v struct {
		S  time.Time `json:"s,format:unix"`
		MS time.Time `json:"ms,format:unixmilli"`
		US time.Time `json:"us,format:unixmicro"`
		NS time.Time `json:"ns,format:unixnano"`
	}
	const payload = `{"s":"1750034397","ms":"1750034397008","us":"1750034396998123","ns":"1750034396998123456"}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := v.S.Unix(); got != 1750034397 {
		t.Errorf("s = %d, want 1750034397", got)
	}
	if got := v.MS.UnixMilli(); got != 1750034397008 {
		t.Errorf("ms = %d, want 1750034397008", got)
	}
	if got := v.US.UnixMicro(); got != 1750034396998123 {
		t.Errorf("us = %d, want 1750034396998123", got)
	}
	if got := v.NS.UnixNano(); got != 1750034396998123456 {
		t.Errorf("ns = %d, want 1750034396998123456", got)
	}
	// The WS trade gateway latency its microsecond receive/push times measure.
	if got := v.NS.Sub(v.US); got != 456*time.Nanosecond {
		t.Errorf("ns-us = %v, want 456ns", got)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal round-trip = %s, want %s", out, payload)
	}
}

// TestUnixMatchesStandard checks that normal values decode exactly like the
// standard `format` option (same instant, same UTC location), whether Bitget
// sends them quoted or bare.
func TestUnixMatchesStandard(t *testing.T) {
	std := jsonexp.ExperimentalSupportFormatTag(true)
	compared := 0
	defer func() {
		if compared < 20 {
			t.Errorf("only %d successful comparisons", compared)
		}
	}()
	for _, format := range []string{"unix", "unixmilli", "unixmicro", "unixnano"} {
		for _, n := range []string{"1750034397008", "1", "-1500", "1750034397.5", "12.000001", "999999999999", "9223372036854775807"} {
			var want, got1, got2 struct {
				T time.Time
			}
			wantErr := unmarshalWithTag(std, format, n, &want.T)
			err1 := unmarshalWithTag(nil, format, n, &got1.T)
			err2 := unmarshalWithTag(nil, format, `"`+n+`"`, &got2.T)
			if (wantErr != nil) != (err1 != nil) || (wantErr != nil) != (err2 != nil) {
				t.Errorf("%s %s: errors differ: std=%v bare=%v quoted=%v", format, n, wantErr, err1, err2)
				continue
			}
			if wantErr != nil {
				continue
			}
			compared++
			if want.T != got1.T || want.T != got2.T {
				t.Errorf("%s %s: std=%v bare=%v quoted=%v", format, n, want.T, got1.T, got2.T)
			}
		}
	}
}

// unmarshalWithTag decodes {"t":raw} into a field tagged with format, using
// the standard library alone when std is set and the Bitget codec otherwise.
func unmarshalWithTag(std json.Options, format, raw string, dst *time.Time) error {
	in := []byte(`{"t":` + raw + `}`)
	var err error
	switch format {
	case "unix":
		var v struct {
			T time.Time `json:"t,format:unix"`
		}
		err = decodeWith(std, in, &v)
		*dst = v.T
	case "unixmilli":
		var v struct {
			T time.Time `json:"t,format:unixmilli"`
		}
		err = decodeWith(std, in, &v)
		*dst = v.T
	case "unixmicro":
		var v struct {
			T time.Time `json:"t,format:unixmicro"`
		}
		err = decodeWith(std, in, &v)
		*dst = v.T
	case "unixnano":
		var v struct {
			T time.Time `json:"t,format:unixnano"`
		}
		err = decodeWith(std, in, &v)
		*dst = v.T
	}
	return err
}

func decodeWith(std json.Options, in []byte, v any) error {
	if std != nil {
		return json.Unmarshal(in, v, std)
	}
	return JSONUnmarshal(in, v)
}

// TestTimeNotSet covers the "not set" forms Bitget emits for timestamps.
func TestTimeNotSet(t *testing.T) {
	for _, raw := range []string{`""`, `"0"`, `"-1"`, `"null"`, `null`, `0`, `-1`} {
		var v struct {
			MS time.Time  `json:"ms,format:unixmilli"`
			US time.Time  `json:"us,format:unixmicro"`
			P  *time.Time `json:"p,format:unixmilli"`
		}
		in := `{"ms":` + raw + `,"us":` + raw + `,"p":` + raw + `}`
		if err := JSONUnmarshal([]byte(in), &v); err != nil {
			t.Errorf("unmarshal %s: %v", raw, err)
			continue
		}
		if !v.MS.IsZero() || !v.US.IsZero() || (v.P != nil && !v.P.IsZero()) {
			t.Errorf("unmarshal %s = %v %v %v, want zero times", raw, v.MS, v.US, v.P)
		}
	}
}

func TestLayoutFormats(t *testing.T) {
	var v struct {
		D  time.Time `json:"d,format:DateOnly"`
		DT time.Time `json:"dt,format:DateTime"`
		C  time.Time `json:"c,format:'2006/01/02'"`
		E  time.Time `json:"e,format:DateOnly"`
	}
	const payload = `{"d":"2026-08-16","dt":"2026-08-16 09:30:00","c":"2026/08/16","e":""}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	want := time.Date(2026, 8, 16, 0, 0, 0, 0, time.UTC)
	if !v.D.Equal(want) || !v.C.Equal(want) || !v.DT.Equal(want.Add(9*time.Hour+30*time.Minute)) || !v.E.IsZero() {
		t.Errorf("got %v %v %v %v", v.D, v.DT, v.C, v.E)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal = %s, want %s", out, payload)
	}
	if err := JSONUnmarshal([]byte(`{"d":20260816}`), &v); err == nil {
		t.Errorf("bare number into DateOnly: want error")
	}
}

// TestStandardFallback checks that fields without a unix/layout format keep
// the standard library's behaviour and errors.
func TestStandardFallback(t *testing.T) {
	var v struct {
		Plain time.Time `json:"plain"`
		RFC   time.Time `json:"rfc,format:RFC3339"`
	}
	if err := JSONUnmarshal([]byte(`{"plain":"2026-08-16T09:30:00Z","rfc":"2026-08-16T09:30:00+08:00"}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if err := JSONUnmarshal([]byte(`{"plain":"1750034397008"}`), &v); err == nil {
		t.Errorf("untagged time.Time accepted a millisecond string; it must require an explicit format")
	}
	var bad struct {
		T time.Time `json:"t,format:bogus"`
	}
	if err := JSONUnmarshal([]byte(`{"t":"1"}`), &bad); err == nil || !strings.Contains(err.Error(), "format") {
		t.Errorf("invalid format error = %v", err)
	}
	if _, err := JSONMarshal(bad); err == nil || !strings.Contains(err.Error(), "format") {
		t.Errorf("invalid format marshal error = %v", err)
	}
}

func TestOtherStandardFormats(t *testing.T) {
	var v struct {
		D time.Duration     `json:"d,format:units"`
		B []byte            `json:"b,format:hex"`
		F float64           `json:"f,format:nonfinite"`
		M map[string]string `json:"m,format:emitnull"`
		S []int             `json:"s,format:emitempty"`
	}
	const payload = `{"d":"1h30m0s","b":"0102","f":"NaN","m":null,"s":[]}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	v.S = nil
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal = %s, want %s", out, payload)
	}
}

func TestFormatOfScope(t *testing.T) {
	type inner struct {
		T time.Time `json:"t"`
	}
	var v struct {
		P  *time.Time  `json:"p,format:unixmilli"`
		In inner       `json:"in"`
		Ts []time.Time `json:"ts"`
	}
	// The nested/element times carry no format, so they must fall back to
	// RFC 3339 rather than inherit the sibling field's unixmilli.
	in := `{"p":"1750034397008","in":{"t":"2026-08-16T00:00:00Z"},"ts":["2026-08-16T00:00:00Z"]}`
	if err := JSONUnmarshal([]byte(in), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.P == nil || v.P.UnixMilli() != 1750034397008 {
		t.Errorf("p = %v", v.P)
	}
}

func TestDecimalCodec(t *testing.T) {
	var v struct {
		Quoted decimal.Decimal `json:"quoted"`
		Bare   decimal.Decimal `json:"bare"`
		Empty  decimal.Decimal `json:"empty"`
		Null   decimal.Decimal `json:"null"`
	}
	if err := JSONUnmarshal([]byte(`{"quoted":"65000.5","bare":0.001,"empty":"","null":null}`), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if v.Quoted.String() != "65000.5" || v.Bare.String() != "0.001" || !v.Empty.IsZero() || !v.Null.IsZero() {
		t.Errorf("got %v %v %v %v", v.Quoted, v.Bare, v.Empty, v.Null)
	}
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if want := `{"quoted":"65000.5","bare":"0.001","empty":"0","null":"0"}`; string(out) != want {
		t.Errorf("marshal = %s, want %s", out, want)
	}
}

// TestEncodeTruncatesToUnit keeps encoded timestamps whole units, like
// UnixMilli/UnixMicro, since Bitget never sends fractional ones.
func TestEncodeTruncatesToUnit(t *testing.T) {
	var v struct {
		MS  time.Time `json:"ms,format:unixmilli"`
		US  time.Time `json:"us,format:unixmicro"`
		Neg time.Time `json:"neg,format:unixmilli"`
	}
	v.MS = time.Unix(1750034397, 8_123_456)
	v.US = v.MS
	v.Neg = time.Unix(-1, 999_999_999)
	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := fmt.Sprintf(`{"ms":"%d","us":"%d","neg":"%d"}`, v.MS.UnixMilli(), v.US.UnixMicro(), v.Neg.UnixMilli())
	if string(out) != want {
		t.Errorf("marshal = %s, want %s", out, want)
	}
}
