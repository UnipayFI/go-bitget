package common

import (
	"testing"
	"time"
)

// TestMicrosTimeCodec pins MicrosTime to microsecond precision, and checks it
// does not disturb the millisecond codec sharing the same JSON value space: the
// two units must decode side by side without either being scaled by 1000.
func TestMicrosTimeCodec(t *testing.T) {
	var v struct {
		CTime       time.Time  `json:"cTime"`
		ReceiveTime MicrosTime `json:"receiveTime"`
		PushTime    MicrosTime `json:"pushTime"`
	}
	const payload = `{"cTime":"1750034397008","receiveTime":"1750034396998123","pushTime":"1750034397076456"}`
	if err := JSONUnmarshal([]byte(payload), &v); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got := v.CTime.UnixMilli(); got != 1750034397008 {
		t.Errorf("cTime = %d ms, want 1750034397008", got)
	}
	if got := v.ReceiveTime.UnixMicro(); got != 1750034396998123 {
		t.Errorf("receiveTime = %d us, want 1750034396998123", got)
	}
	if got := v.PushTime.UnixMicro(); got != 1750034397076456 {
		t.Errorf("pushTime = %d us, want 1750034397076456", got)
	}
	// The gateway latency these two fields exist to measure.
	if got := v.PushTime.Sub(v.ReceiveTime.Time); got != 78333*time.Microsecond {
		t.Errorf("pushTime-receiveTime = %v, want 78.333ms", got)
	}

	out, err := JSONMarshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(out) != payload {
		t.Errorf("marshal round-trip = %s, want %s", out, payload)
	}
}

// TestMicrosTimeNotSet covers the "not set" forms Bitget emits for timestamps.
func TestMicrosTimeNotSet(t *testing.T) {
	for _, raw := range []string{`""`, `"0"`, `"-1"`, `null`, `0`} {
		var v struct {
			T MicrosTime `json:"t"`
		}
		if err := JSONUnmarshal([]byte(`{"t":`+raw+`}`), &v); err != nil {
			t.Errorf("unmarshal %s: %v", raw, err)
			continue
		}
		if !v.T.IsZero() {
			t.Errorf("unmarshal %s = %v, want zero time", raw, v.T)
		}
	}
}
