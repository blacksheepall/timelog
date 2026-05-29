package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

type CreateTimeLogMCPInput struct {
	CategoryID int32
	StartTime  string
	EndTime    string
	TaskID     int32
	Remark     string
}

type UpdateTimeLogMCPInput struct {
	ID         int32
	CategoryID int32
	StartTime  string
	EndTime    string
	TaskID     int32
	Remark     string
}

func CreateTimeLogFromMCPInput(input CreateTimeLogMCPInput) (*gen.Timelog, error) {
	if _, err := GetCategoryByID(input.CategoryID); err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	tl := &gen.Timelog{CategoryID: input.CategoryID}
	defaultUserID := int32(1)
	tl.UserID = &defaultUserID

	if input.StartTime != "" {
		st, err := ParseSGDateTime(input.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		tl.StartTime = st.UTC()
	} else {
		tl.StartTime = time.Now().UTC()
	}

	if input.EndTime != "" {
		et, err := ParseSGDateTime(input.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		etUTC := et.UTC()
		tl.EndTime = &etUTC
	}

	if input.TaskID > 0 {
		tl.TaskID = &input.TaskID
	}
	if input.Remark != "" {
		tl.Remark = &input.Remark
	}

	if err := CreateTimeLog(tl); err != nil {
		return nil, fmt.Errorf("failed to create time log: %w", err)
	}
	return GetTimeLogByID(*tl.ID)
}

func UpdateTimeLogFromMCPInput(input UpdateTimeLogMCPInput) (*gen.Timelog, error) {
	tl, err := GetTimeLogByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("time log not found: %w", err)
	}

	if input.CategoryID > 0 {
		if _, err := GetCategoryByID(input.CategoryID); err != nil {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		tl.CategoryID = input.CategoryID
	}

	if input.StartTime != "" {
		st, err := ParseSGDateTime(input.StartTime)
		if err != nil {
			return nil, fmt.Errorf("invalid start_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		tl.StartTime = st.UTC()
	}

	if input.EndTime != "" {
		et, err := ParseSGDateTime(input.EndTime)
		if err != nil {
			return nil, fmt.Errorf("invalid end_time format, expected YYYY-MM-DD HH:MM:SS: %w", err)
		}
		etUTC := et.UTC()
		tl.EndTime = &etUTC
	}

	if input.TaskID > 0 {
		tl.TaskID = &input.TaskID
	}
	if input.Remark != "" {
		tl.Remark = &input.Remark
	}

	if err := UpdateTimeLog(tl); err != nil {
		return nil, fmt.Errorf("failed to update time log: %w", err)
	}
	return GetTimeLogByID(input.ID)
}

func ListActiveTimeLogs() ([]gen.Timelog, error) {
	return ListTimeLogsWithOptions(0, "start_time DESC", "end_time IS NULL")
}
