package common

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/shopspring/decimal"
)

// Bitget encodes numbers and timestamps as JSON strings, inconsistently:
// timestamps are usually a quoted millisecond string ("1782057654329") but
// occasionally a bare number, and amounts/prices are quoted decimal strings.
// Both are emitted as "" (and timestamps as "0"/"-1") when "not set". The stock
// time.Time / shopspring decimal codecs reject the empty-string form, so we
// teach the JSON codec how to read/write both types once, globally. Every
// time.Time / decimal.Decimal field in this SDK is therefore a plain field with
// a plain json tag, and the conversions below apply.
var (
	unmarshalers = json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(decodeMillisTime),
		json.UnmarshalFromFunc(decodeMicrosTime),
		json.UnmarshalFromFunc(decodeDecimal),
	))
	marshalers = json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(encodeMillisTime),
		json.MarshalToFunc(encodeMicrosTime),
		json.MarshalToFunc(encodeDecimal),
	))
)

// MicrosTime is a timestamp Bitget encodes in microseconds rather than the
// usual milliseconds — currently the WebSocket trade gateway's receive/push
// times. Reading one with the millisecond codec above would misdate it by a
// factor of 1000, so it gets its own type; the embedded time.Time makes it
// usable as an ordinary time value.
type MicrosTime struct{ time.Time }

// JSONMarshal marshals v with Bitget's millisecond-time and decimal-string
// conventions applied.
func JSONMarshal(v any) ([]byte, error) {
	return json.Marshal(v, marshalers)
}

// JSONUnmarshal unmarshals data into v with Bitget's millisecond-time and
// decimal-string conventions applied.
func JSONUnmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v, unmarshalers)
}

// decodeTimestamp reads Bitget's quoted-or-bare timestamp token in the given
// unit, reporting set=false for null and the "not set" sentinels.
func decodeTimestamp(dec *jsontext.Decoder, unit string) (n int64, set bool, err error) {
	tok, err := dec.ReadToken()
	if err != nil {
		return 0, false, err
	}
	var s string
	switch tok.Kind() {
	case 'n': // null
		return 0, false, nil
	case '"': // quoted string
		s = tok.String()
	case '0': // bare number
		s = tok.String()
	default:
		return 0, false, fmt.Errorf("bitget: cannot decode %v token into %s timestamp", tok.Kind(), unit)
	}
	switch s {
	case "", "0", "-1": // "not set" sentinels
		return 0, false, nil
	}
	n, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, false, fmt.Errorf("bitget: invalid %s timestamp %q: %w", unit, s, err)
	}
	return n, true, nil
}

func decodeMillisTime(dec *jsontext.Decoder, t *time.Time) error {
	ms, set, err := decodeTimestamp(dec, "millisecond")
	if err != nil {
		return err
	}
	if !set {
		*t = time.Time{}
		return nil
	}
	*t = time.UnixMilli(ms)
	return nil
}

func encodeMillisTime(enc *jsontext.Encoder, t time.Time) error {
	if t.IsZero() {
		return enc.WriteToken(jsontext.String(""))
	}
	return enc.WriteToken(jsontext.String(strconv.FormatInt(t.UnixMilli(), 10)))
}

func decodeMicrosTime(dec *jsontext.Decoder, t *MicrosTime) error {
	us, set, err := decodeTimestamp(dec, "microsecond")
	if err != nil {
		return err
	}
	if !set {
		*t = MicrosTime{}
		return nil
	}
	*t = MicrosTime{time.UnixMicro(us)}
	return nil
}

func encodeMicrosTime(enc *jsontext.Encoder, t MicrosTime) error {
	if t.IsZero() {
		return enc.WriteToken(jsontext.String(""))
	}
	return enc.WriteToken(jsontext.String(strconv.FormatInt(t.UnixMicro(), 10)))
}

func decodeDecimal(dec *jsontext.Decoder, d *decimal.Decimal) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	var s string
	switch tok.Kind() {
	case 'n': // null
		*d = decimal.Zero
		return nil
	case '"': // quoted string
		s = tok.String()
	case '0': // bare number
		s = tok.String()
	default:
		return fmt.Errorf("bitget: cannot decode %v token into decimal", tok.Kind())
	}
	if s == "" {
		*d = decimal.Zero
		return nil
	}
	v, err := decimal.NewFromString(s)
	if err != nil {
		return fmt.Errorf("bitget: invalid decimal %q: %w", s, err)
	}
	*d = v
	return nil
}

func encodeDecimal(enc *jsontext.Encoder, d decimal.Decimal) error {
	return enc.WriteToken(jsontext.String(d.String()))
}
