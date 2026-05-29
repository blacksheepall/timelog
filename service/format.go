package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model"
)

func sgt() *time.Location {
	return model.GetSingaporeLocation()
}

func FormatSGDateTime(t time.Time) string {
	return t.In(sgt()).Format("2006-01-02 15:04:05")
}

func FormatSGDate(t time.Time) string {
	return t.In(sgt()).Format("2006-01-02")
}

func FormatSGDateTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatSGDateTime(*t)
}

func FormatDuration(d time.Duration) string {
	return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes())%60)
}

func ParseSGDateTime(value string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", value, sgt())
}

func TodaySGDateString() string {
	return time.Now().In(sgt()).Format("2006-01-02")
}
