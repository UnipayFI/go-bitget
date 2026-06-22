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
		json.UnmarshalFromFunc(decodeDecimal),
	))
	marshalers = json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(encodeMillisTime),
		json.MarshalToFunc(encodeDecimal),
	))
)

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

func decodeMillisTime(dec *jsontext.Decoder, t *time.Time) error {
	tok, err := dec.ReadToken()
	if err != nil {
		return err
	}
	var s string
	switch tok.Kind() {
	case 'n': // null
		*t = time.Time{}
		return nil
	case '"': // quoted string
		s = tok.String()
	case '0': // bare number
		s = tok.String()
	default:
		return fmt.Errorf("bitget: cannot decode %v token into time.Time", tok.Kind())
	}
	switch s {
	case "", "0", "-1": // "not set" sentinels
		*t = time.Time{}
		return nil
	}
	ms, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("bitget: invalid millisecond timestamp %q: %w", s, err)
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
