// Package timeutil holds timezone and formatting helpers so that the model
// package can remain focused on persistence.
package timeutil

import (
	"fmt"
	"time"
)

// location holds the configured application timezone. Defaults to
// Asia/Singapore for backward compatibility.
var location *time.Location

func init() {
	var err error
	location, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		// Fallback to UTC+8 if timezone data is not available
		location = time.FixedZone("SGT", 8*60*60)
	}
}

// SetTimezone configures the application timezone by IANA name (e.g.
// "Asia/Singapore"). It returns an error if the timezone cannot be loaded.
func SetTimezone(name string) error {
	loc, err := time.LoadLocation(name)
	if err != nil {
		return err
	}
	location = loc
	return nil
}

// GetLocation returns the currently configured application timezone.
func GetLocation() *time.Location {
	return location
}

// GetSingaporeLocation returns the currently configured application timezone.
// It is kept for backward compatibility; use GetLocation for new code.
func GetSingaporeLocation() *time.Location {
	return location
}

// ParseSGDateTime parses a "2006-01-02 15:04:05" string in the configured
// application timezone.
func ParseSGDateTime(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", value, location)
}

// FormatSGDateTime formats t as "2006-01-02 15:04:05" in the configured
// application timezone.
func FormatSGDateTime(t time.Time) string {
	return t.In(location).Format("2006-01-02 15:04:05")
}

// FormatSGDate formats t as "2006-01-02" in the configured application
// timezone.
func FormatSGDate(t time.Time) string {
	return t.In(location).Format("2006-01-02")
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

// TodaySGDateString returns today's date as "2006-01-02" in the configured
// application timezone.
func TodaySGDateString() string {
	return time.Now().In(location).Format("2006-01-02")
}
