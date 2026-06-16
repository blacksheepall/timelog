package mapper

import (
	"fmt"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TaskToProto(t *domain.Task) *timelogv1.Task {
	if t == nil {
		return nil
	}
	return &timelogv1.Task{
		Id:               t.ID,
		Title:            t.Title,
		Description:      t.Description,
		CategoryId:       t.CategoryID,
		DueDate:          FormatDate(t.DueDate),
		EstimatedMinutes: t.EstimatedMinutes,
		IsCompleted:      t.IsCompleted,
		CompletedAt:      optionalString(FormatTimeUTCPtr(t.CompletedAt)),
		IsSuspended:      t.IsSuspended,
		CreatedAt:        FormatTimeUTC(t.CreatedAt),
		UpdatedAt:        FormatTimeUTC(t.UpdatedAt),
	}
}

func TasksToProto(tasks []domain.Task) []*timelogv1.Task {
	out := make([]*timelogv1.Task, 0, len(tasks))
	for i := range tasks {
		out = append(out, TaskToProto(&tasks[i]))
	}
	return out
}

func TaskFromCreateRequest(req *timelogv1.CreateTaskRequest) (*domain.Task, error) {
	if req == nil {
		return nil, nil
	}
	dueDate, err := ParseDate(req.DueDate)
	if err != nil {
		return nil, err
	}
	return &domain.Task{
		Title:            req.Title,
		Description:      req.Description,
		CategoryID:       req.CategoryId,
		DueDate:          dueDate,
		EstimatedMinutes: req.EstimatedMinutes,
	}, nil
}

func ApplyTaskUpdate(task *domain.Task, req *timelogv1.UpdateTaskRequest) error {
	if task == nil || req == nil {
		return nil
	}
	if req.Title != nil {
		task.Title = req.GetTitle()
	}
	if req.Description != nil {
		task.Description = req.GetDescription()
	}
	if req.CategoryId != nil {
		task.CategoryID = req.GetCategoryId()
	}
	if req.DueDate != nil {
		dueDate, err := ParseDate(req.GetDueDate())
		if err != nil {
			return err
		}
		task.DueDate = dueDate
	}
	if req.EstimatedMinutes != nil {
		task.EstimatedMinutes = req.GetEstimatedMinutes()
	}
	if req.IsCompleted != nil {
		task.IsCompleted = req.GetIsCompleted()
	}
	return nil
}

func TaskStatsToProto(stats map[string]interface{}) (*timelogv1.TaskStats, error) {
	total, err := numberToInt32(stats["total_tasks"])
	if err != nil {
		return nil, fmt.Errorf("total_tasks: %w", err)
	}
	completed, err := numberToInt32(stats["completed_tasks"])
	if err != nil {
		return nil, fmt.Errorf("completed_tasks: %w", err)
	}
	rate, err := numberToFloat64(stats["completion_rate"])
	if err != nil {
		return nil, fmt.Errorf("completion_rate: %w", err)
	}
	return &timelogv1.TaskStats{
		TotalTasks:     total,
		CompletedTasks: completed,
		CompletionRate: rate,
	}, nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func numberToInt32(value interface{}) (int32, error) {
	switch v := value.(type) {
	case int:
		return int32(v), nil
	case int32:
		return v, nil
	case int64:
		return int32(v), nil
	case float64:
		return int32(v), nil
	default:
		return 0, fmt.Errorf("unsupported number type %T", value)
	}
}

func numberToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, fmt.Errorf("unsupported number type %T", value)
	}
}
