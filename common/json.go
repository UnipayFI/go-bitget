package common

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"time"

	jsonexp "github.com/go-json-experiment/json"
	"github.com/shopspring/decimal"
)

// Bitget encodes numbers and timestamps as JSON strings, inconsistently:
// timestamps are usually a quoted digit string ("1782057654329") but
// occasionally a bare number, and amounts/prices are quoted decimal strings.
// Both are emitted as "" (and timestamps as "0"/"-1"/"null") when "not set".
//
// Every time.Time field declares its wire format with the standard `format`
// tag option (e.g. `json:"cTime,format:unixmilli"`), which Go 1.27's
// encoding/json/v2 only honours when ExperimentalSupportFormatTag is set. The
// time codec below keeps the standard semantics of that format and only adds
// Bitget's quirks on top: quoted-or-bare numbers and the "not set" sentinels.
// decimal.Decimal fields stay plain fields with a plain json tag.
var (
	unmarshalOptions json.Options
	marshalOptions   json.Options

	// errJSONSupport is non-nil when the running Go release no longer
	// provides the `format` tag hooks this codec relies on.
	errJSONSupport = initJSON()
)

// JSONMarshal marshals v with Bitget's time and decimal conventions applied.
func JSONMarshal(v any) ([]byte, error) {
	if errJSONSupport != nil {
		return nil, errJSONSupport
	}
	return json.Marshal(v, marshalOptions)
}

// JSONUnmarshal unmarshals data into v with Bitget's time and decimal
// conventions applied.
func JSONUnmarshal(data []byte, v any) error {
	if errJSONSupport != nil {
		return errJSONSupport
	}
	return json.Unmarshal(data, v, unmarshalOptions)
}

// initJSON builds the codec options and round-trips a probe through them, so
// a Go release that drops the experimental `format` tag support (by panicking
// on the unknown option or by ignoring it) or changes the options layout
// formatOf reads fails loudly instead of silently misdating fields.
func initJSON() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("bitget: encoding/json/v2 format tag support unavailable: %v", r)
		}
	}()
	if errFormatOf != nil {
		return errFormatOf
	}
	formatTag := jsonexp.ExperimentalSupportFormatTag(true)
	unmarshalOptions = json.JoinOptions(formatTag, json.WithUnmarshalers(json.JoinUnmarshalers(
		json.UnmarshalFromFunc(decodeTime),
		json.UnmarshalFromFunc(decodeDecimal),
	)))
	marshalOptions = json.JoinOptions(formatTag, json.WithMarshalers(json.JoinMarshalers(
		json.MarshalToFunc(encodeTime),
		json.MarshalToFunc(encodeDecimal),
	)))

	const payload = `{"t":"1750034396998123"}`
	var probe struct {
		T time.Time `json:"t,format:unixmicro"`
	}
	if err := json.Unmarshal([]byte(payload), &probe, unmarshalOptions); err != nil {
		return fmt.Errorf("bitget: encoding/json/v2 format tag support unavailable: %w", err)
	}
	if got := probe.T.UnixMicro(); got != 1750034396998123 {
		return fmt.Errorf("bitget: encoding/json/v2 format tag probe decoded %d, want 1750034396998123", got)
	}
	if out, err := json.Marshal(probe, marshalOptions); err != nil || string(out) != payload {
		return fmt.Errorf("bitget: encoding/json/v2 format tag probe encoded %s (%v), want %s", out, err, payload)
	}
	return nil
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
