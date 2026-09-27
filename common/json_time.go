package common

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

// timeFormat is the parsed `format` tag option of a time.Time field.
type timeFormat struct {
	pow10  uint64 // unix formats: 1e0, 1e3, 1e6 or 1e9 units per second
	layout string // layout formats: the time.Parse layout
}

// parseTimeFormat mirrors encoding/json/v2's interpretation of a time.Time
// `format` value. ok is false for formats this codec leaves to the standard
// library: no format (RFC 3339 default), RFC3339/RFC3339Nano (which get extra
// validation there) and invalid values (reported there).
func parseTimeFormat(format string) (f timeFormat, ok bool) {
	if format == "" {
		return f, false
	}
	// We assume that an exported constant in the time package will
	// always start with an uppercase ASCII letter.
	if c := format[0]; !('a' <= c && c <= 'z') && !('A' <= c && c <= 'Z') {
		return timeFormat{layout: format}, true
	}
	switch format {
	case "unix":
		return timeFormat{pow10: 1e0}, true
	case "unixmilli":
		return timeFormat{pow10: 1e3}, true
	case "unixmicro":
		return timeFormat{pow10: 1e6}, true
	case "unixnano":
		return timeFormat{pow10: 1e9}, true
	case "RFC3339", "RFC3339Nano":
		return f, false
	case "ANSIC":
		return timeFormat{layout: time.ANSIC}, true
	case "UnixDate":
		return timeFormat{layout: time.UnixDate}, true
	case "RubyDate":
		return timeFormat{layout: time.RubyDate}, true
	case "RFC822":
		return timeFormat{layout: time.RFC822}, true
	case "RFC822Z":
		return timeFormat{layout: time.RFC822Z}, true
	case "RFC850":
		return timeFormat{layout: time.RFC850}, true
	case "RFC1123":
		return timeFormat{layout: time.RFC1123}, true
	case "RFC1123Z":
		return timeFormat{layout: time.RFC1123Z}, true
	case "Kitchen":
		return timeFormat{layout: time.Kitchen}, true
	case "Stamp":
		return timeFormat{layout: time.Stamp}, true
	case "StampMilli":
		return timeFormat{layout: time.StampMilli}, true
	case "StampMicro":
		return timeFormat{layout: time.StampMicro}, true
	case "StampNano":
		return timeFormat{layout: time.StampNano}, true
	case "DateTime":
		return timeFormat{layout: time.DateTime}, true
	case "DateOnly":
		return timeFormat{layout: time.DateOnly}, true
	case "TimeOnly":
		return timeFormat{layout: time.TimeOnly}, true
	}
	// Reject any Go identifier in case new constants are supported.
	if strings.TrimFunc(format, isLetterOrDigit) == "" {
		return f, false
	}
	return timeFormat{layout: format}, true
}

func isLetterOrDigit(r rune) bool {
	return r == '_' || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || ('0' <= r && r <= '9')
}

// decodeTime decodes a time.Time field whose `format` is a unix unit or a
// layout, accepting the value quoted or bare and mapping Bitget's "not set"
// sentinels to the zero time. Anything else is decoded exactly as the
// standard format would decode it.
func decodeTime(dec *jsontext.Decoder, t *time.Time) error {
	f, ok := parseTimeFormat(formatOf(dec.Options()))
	if !ok {
		return errors.ErrUnsupported
	}
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}
	var b []byte
	switch val.Kind() {
	case 'n': // null
		*t = time.Time{}
		return nil
	case '"': // quoted string
		if b, err = unquote(val); err != nil {
			return err
		}
	case '0': // bare number
		if f.pow10 == 0 {
			return fmt.Errorf("bitget: cannot decode JSON number into %s time", f.layout)
		}
		b = val
	default:
		return fmt.Errorf("bitget: cannot decode %v value into time", val.Kind())
	}
	switch string(b) {
	case "", "0", "-1", "null": // "not set" sentinels
		*t = time.Time{}
		return nil
	}
	if f.pow10 != 0 {
		*t, err = parseTimeUnix(b, f.pow10)
	} else {
		*t, err = time.Parse(f.layout, string(b))
	}
	return err
}

// encodeTime is the inverse of decodeTime: a quoted value in the field's
// format, or "" for the zero time.
func encodeTime(enc *jsontext.Encoder, t time.Time) error {
	f, ok := parseTimeFormat(formatOf(enc.Options()))
	if !ok {
		return errors.ErrUnsupported
	}
	if t.IsZero() {
		return enc.WriteToken(jsontext.String(""))
	}
	var arr [64]byte
	b := append(arr[:0], '"')
	if f.pow10 != 0 {
		// Bitget timestamps are whole units, so truncate sub-unit precision
		// (as UnixMilli/UnixMicro do) rather than emit a fraction.
		t = t.Add(-time.Duration(t.Nanosecond() % int(1e9/f.pow10)))
		b = append(appendTimeUnix(b, t, f.pow10), '"')
	} else {
		var err error
		if b, err = jsontext.AppendQuote(b[:0], t.Format(f.layout)); err != nil {
			return err
		}
	}
	return enc.WriteValue(b)
}

// unquote returns the contents of a JSON string, avoiding a copy when it
// contains no escape sequences (timestamps never do).
func unquote(val jsontext.Value) ([]byte, error) {
	if b := val[1 : len(val)-1]; !containsByte(b, '\\') {
		return b, nil
	}
	return jsontext.AppendUnquote(nil, val)
}

func containsByte(b []byte, c byte) bool {
	for _, x := range b {
		if x == c {
			return true
		}
	}
	return false
}

// formatOf returns the `format` tag option of the struct field currently
// being marshaled or unmarshaled, or "" if it has none. Pass it the Options
// of the Encoder/Decoder given to a MarshalToFunc/UnmarshalFromFunc; they are
// only valid for the duration of that call.
//
// encoding/json/v2 has no public accessor for the format, so this reads the
// Format field of the options struct those methods return (a
// *jsonopts.Struct in Go 1.27). errFormatOf reports when that layout is not
// what this code expects, in which case formatOf always returns "".
//
// Format is only reset between struct fields, so it must only be consulted for
// the value a format-tagged field holds directly (time.Time or *time.Time), not
// for values nested inside a format-tagged type with its own JSON methods.
func formatOf(opts json.Options) string {
	if errFormatOf != nil || reflect.TypeOf(opts) != optionsType {
		return ""
	}
	// An interface value is a (type, data) word pair; data points at the struct.
	p := (*[2]unsafe.Pointer)(unsafe.Pointer(&opts))[1]
	return *(*string)(unsafe.Add(p, formatOffset))
}

var (
	optionsType  = reflect.TypeOf(new(jsontext.Decoder).Options())
	formatOffset uintptr
	errFormatOf  = func() error {
		if enc := reflect.TypeOf(new(jsontext.Encoder).Options()); enc != optionsType {
			return fmt.Errorf("bitget: encoder options type %v differs from decoder options type %v", enc, optionsType)
		}
		if optionsType == nil || optionsType.Kind() != reflect.Pointer || optionsType.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("bitget: unexpected encoding/json/v2 options type %v", optionsType)
		}
		sf, ok := optionsType.Elem().FieldByName("Format")
		if !ok || sf.Type.Kind() != reflect.String {
			return fmt.Errorf("bitget: encoding/json/v2 options type %v has no Format string field", optionsType)
		}
		// FieldByName reports the offset within the innermost embedded
		// struct; sum the offsets along the embedding path.
		st := optionsType.Elem()
		for _, i := range sf.Index {
			if st.Kind() != reflect.Struct {
				return fmt.Errorf("bitget: encoding/json/v2 options field Format is not embedded by value in %v", optionsType)
			}
			f := st.Field(i)
			formatOffset += f.Offset
			st = f.Type
		}
		return nil
	}()
)
