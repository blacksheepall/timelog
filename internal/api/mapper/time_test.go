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

func TestFormatDate(t *testing.T) {
	input := time.Date(2026, 5, 30, 20, 15, 0, 0, time.UTC)
	got := FormatDate(input)
	want := "2026-05-30"
	if got != want {
		t.Fatalf("FormatDate() = %q, want %q", got, want)
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
