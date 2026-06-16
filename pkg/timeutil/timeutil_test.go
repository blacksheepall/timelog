package timeutil

import (
	"testing"
	"time"
)

func TestGetSingaporeLocation(t *testing.T) {
	loc := GetSingaporeLocation()
	if loc == nil {
		t.Fatal("expected non-nil location")
	}
}

func TestParseSGDateTime(t *testing.T) {
	got, err := ParseSGDateTime("2026-05-30 10:00:00")
	if err != nil {
		t.Fatalf("ParseSGDateTime: %v", err)
	}
	if FormatSGDateTime(got) != "2026-05-30 10:00:00" {
		t.Fatalf("unexpected: %s", FormatSGDateTime(got))
	}
}

func TestFormatSGDate(t *testing.T) {
	input := time.Date(2026, 5, 30, 10, 0, 0, 0, time.UTC)
	got := FormatSGDate(input)
	if got != "2026-05-30" && got != "2026-05-31" {
		t.Fatalf("unexpected date: %s", got)
	}
}

func TestFormatSGDateTimePtr(t *testing.T) {
	input := time.Date(2026, 5, 30, 10, 0, 0, 0, time.UTC)
	if got := FormatSGDateTimePtr(&input); got == "" {
		t.Fatal("expected non-empty")
	}
	if got := FormatSGDateTimePtr(nil); got != "" {
		t.Fatalf("expected empty for nil, got %q", got)
	}
}

func TestFormatDuration(t *testing.T) {
	if got := FormatDuration(90 * time.Minute); got != "1h 30m" {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestTodaySGDateString(t *testing.T) {
	if got := TodaySGDateString(); got == "" {
		t.Fatal("expected non-empty date string")
	}
}
