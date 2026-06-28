package service

import (
	"time"

	"github.com/blacksheepaul/timelog/core/timeutil"
)

func sgt() *time.Location {
	return timeutil.GetSingaporeLocation()
}

func FormatSGDateTime(t time.Time) string {
	return timeutil.FormatSGDateTime(t)
}

func FormatSGDate(t time.Time) string {
	return timeutil.FormatSGDate(t)
}

func FormatSGDateTimePtr(t *time.Time) string {
	return timeutil.FormatSGDateTimePtr(t)
}

func FormatDuration(d time.Duration) string {
	return timeutil.FormatDuration(d)
}

func ParseSGDateTime(value string) (time.Time, error) {
	return timeutil.ParseSGDateTime(value)
}

func TodaySGDateString() string {
	return timeutil.TodaySGDateString()
}
