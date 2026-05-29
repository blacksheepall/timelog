package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model"
)

func TestFormatSGDateTime(t *testing.T) {
	loc := model.GetSingaporeLocation()
	utc := time.Date(2026, 5, 29, 2, 30, 0, 0, time.UTC) // 10:30 SGT
	got := FormatSGDateTime(utc)
	want := utc.In(loc).Format("2006-01-02 15:04:05")
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestFormatDuration(t *testing.T) {
	got := FormatDuration(2*time.Hour + 15*time.Minute)
	if got != "2h 15m" {
		t.Fatalf("got %q", got)
	}
}
