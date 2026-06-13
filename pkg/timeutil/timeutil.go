// Package timeutil holds timezone and formatting helpers so that the model
// package can remain focused on persistence.
package timeutil

import (
	"fmt"
	"time"
)

// singaporeLocation is loaded once and reused.
var singaporeLocation *time.Location

func init() {
	var err error
	singaporeLocation, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		// Fallback to UTC+8 if timezone data is not available
		singaporeLocation = time.FixedZone("SGT", 8*60*60)
	}
}

// GetSingaporeLocation returns the Asia/Singapore timezone.
func GetSingaporeLocation() *time.Location {
	return singaporeLocation
}

// ParseSGDateTime parses a "2006-01-02 15:04:05" string in Asia/Singapore.
func ParseSGDateTime(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", value, singaporeLocation)
}

// FormatSGDateTime formats t as "2006-01-02 15:04:05" in Asia/Singapore.
func FormatSGDateTime(t time.Time) string {
	return t.In(singaporeLocation).Format("2006-01-02 15:04:05")
}

// FormatSGDate formats t as "2006-01-02" in Asia/Singapore.
func FormatSGDate(t time.Time) string {
	return t.In(singaporeLocation).Format("2006-01-02")
}

// FormatSGDateTimePtr formats a pointer to time as FormatSGDateTime does.
func FormatSGDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatSGDateTime(*t)
}

// FormatDuration formats d as "Xh Ym".
func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

// TodaySGDateString returns today's date as "2006-01-02" in Asia/Singapore.
func TodaySGDateString() string {
	return time.Now().In(singaporeLocation).Format("2006-01-02")
}
