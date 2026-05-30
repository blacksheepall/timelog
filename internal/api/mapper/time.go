package mapper

import (
	"fmt"
	"time"
)

func FormatTimeUTC(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func FormatTimeUTCPtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTimeUTC(*t)
}

func FormatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func FormatDatePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatDate(*t)
}

func ParseTimeUTC(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid RFC3339 time %q: %w", value, err)
	}
	return t.UTC(), nil
}

func ParseOptionalTimeUTC(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	t, err := ParseTimeUTC(*value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func ParseDate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date %q: %w", value, err)
	}
	return t, nil
}

func ParseOptionalDate(value *string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	t, err := ParseDate(*value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func StringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func Int32Value(value *int32) int32 {
	if value == nil {
		return 0
	}
	return *value
}

func BoolValue(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}
