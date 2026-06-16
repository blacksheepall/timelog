package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
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

func (s *Service) CreateTimeLogFromMCPInput(input CreateTimeLogMCPInput) (*domain.TimeLog, error) {
	if _, err := s.GetCategoryByID(input.CategoryID); err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	tl := &domain.TimeLog{
		UserID:     1,
		CategoryID: input.CategoryID,
	}

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
		tl.Remark = input.Remark
	}

	if err := s.CreateTimeLog(tl); err != nil {
		return nil, fmt.Errorf("failed to create time log: %w", err)
	}
	return s.GetTimeLogByID(tl.ID)
}

func (s *Service) UpdateTimeLogFromMCPInput(input UpdateTimeLogMCPInput) (*domain.TimeLog, error) {
	tl, err := s.GetTimeLogByID(input.ID)
	if err != nil {
		return nil, fmt.Errorf("time log not found: %w", err)
	}

	if input.CategoryID > 0 {
		if _, err := s.GetCategoryByID(input.CategoryID); err != nil {
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
		tl.Remark = input.Remark
	}

	if err := s.UpdateTimeLog(tl); err != nil {
		return nil, fmt.Errorf("failed to update time log: %w", err)
	}
	return s.GetTimeLogByID(input.ID)
}

func (s *Service) ListActiveTimeLogs() ([]domain.TimeLog, error) {
	return s.ListTimeLogsWithOptions(0, "start_time DESC", "end_time IS NULL")
}
