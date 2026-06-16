package mapper

import (
	"testing"
	"time"
)

func TestFormatTimeUTC(t *testing.T) {
	loc := time.FixedZone("SGT", 8*60*60)
	input := time.Date(2026, 5, 30, 20, 15, 0, 0, loc)
	got := FormatTimeUTC(input)
	want := "2026-05-30T12:15:00Z"
	if got != want {
		t.Fatalf("FormatTimeUTC() = %q, want %q", got, want)
	}
}

func TestFormatTimeUTCZero(t *testing.T) {
	if got := FormatTimeUTC(time.Time{}); got != "" {
		t.Fatalf("expected empty string for zero time, got %q", got)
	}
}

func TestFormatTimeUTCPtr(t *testing.T) {
	input := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	if got := FormatTimeUTCPtr(&input); got != "2026-05-30T12:00:00Z" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := FormatTimeUTCPtr(nil); got != "" {
		t.Fatalf("expected empty string for nil, got %q", got)
	}
}

func TestFormatDate(t *testing.T) {
	input := time.Date(2026, 5, 30, 20, 15, 0, 0, time.UTC)
	got := FormatDate(input)
	want := "2026-05-30"
	if got != want {
		t.Fatalf("FormatDate() = %q, want %q", got, want)
	}
}

func TestFormatDateZero(t *testing.T) {
	if got := FormatDate(time.Time{}); got != "" {
		t.Fatalf("expected empty string for zero date, got %q", got)
	}
}

func TestFormatDatePtr(t *testing.T) {
	input := time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)
	if got := FormatDatePtr(&input); got != "2026-05-30" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := FormatDatePtr(nil); got != "" {
		t.Fatalf("expected empty string for nil, got %q", got)
	}
}

func TestParseTimeUTC(t *testing.T) {
	got, err := ParseTimeUTC("2026-05-30T20:15:00+08:00")
	if err != nil {
		t.Fatalf("ParseTimeUTC() error = %v", err)
	}
	want := "2026-05-30T12:15:00Z"
	if FormatTimeUTC(got) != want {
		t.Fatalf("ParseTimeUTC() = %q, want %q", FormatTimeUTC(got), want)
	}
}

func TestParseTimeUTCEmpty(t *testing.T) {
	got, err := ParseTimeUTC("")
	if err != nil || !got.IsZero() {
		t.Fatalf("expected zero time, got (%v, %v)", got, err)
	}
}

func TestParseOptionalTimeUTC(t *testing.T) {
	input := "2026-05-30T12:00:00Z"
	got, err := ParseOptionalTimeUTC(&input)
	if err != nil || got == nil {
		t.Fatalf("unexpected: (%v, %v)", got, err)
	}
	if got, err := ParseOptionalTimeUTC(nil); err != nil || got != nil {
		t.Fatalf("expected nil, got (%v, %v)", got, err)
	}
	empty := ""
	if got, err := ParseOptionalTimeUTC(&empty); err != nil || got != nil {
		t.Fatalf("expected nil for empty string, got (%v, %v)", got, err)
	}
}

func TestParseDate(t *testing.T) {
	got, err := ParseDate("2026-05-30")
	if err != nil {
		t.Fatalf("ParseDate() error = %v", err)
	}
	if FormatDate(got) != "2026-05-30" {
		t.Fatalf("unexpected: %q", FormatDate(got))
	}
}

func TestParseDateEmpty(t *testing.T) {
	got, err := ParseDate("")
	if err != nil || !got.IsZero() {
		t.Fatalf("expected zero date, got (%v, %v)", got, err)
	}
}

func TestParseOptionalDate(t *testing.T) {
	input := "2026-05-30"
	got, err := ParseOptionalDate(&input)
	if err != nil || got == nil {
		t.Fatalf("unexpected: (%v, %v)", got, err)
	}
	if got, err := ParseOptionalDate(nil); err != nil || got != nil {
		t.Fatalf("expected nil, got (%v, %v)", got, err)
	}
	empty := ""
	if got, err := ParseOptionalDate(&empty); err != nil || got != nil {
		t.Fatalf("expected nil for empty string, got (%v, %v)", got, err)
	}
}

func TestStringValue(t *testing.T) {
	s := "hello"
	if got := StringValue(&s); got != "hello" {
		t.Fatalf("unexpected: %q", got)
	}
	if got := StringValue(nil); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestInt32Value(t *testing.T) {
	v := int32(42)
	if got := Int32Value(&v); got != 42 {
		t.Fatalf("unexpected: %d", got)
	}
	if got := Int32Value(nil); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestBoolValue(t *testing.T) {
	b := true
	if got := BoolValue(&b); !got {
		t.Fatalf("expected true")
	}
	if got := BoolValue(nil); got {
		t.Fatalf("expected false")
	}
}
