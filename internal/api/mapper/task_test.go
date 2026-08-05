package mapper

import (
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestTaskToProtoFormatsDateAndOptionalTime(t *testing.T) {
	completedAt := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	task := &domain.Task{
		ID:               5,
		Title:            "IDL",
		Description:      "Write tests",
		CategoryID:       3,
		DueDate:          time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		EstimatedMinutes: 90,
		IsCompleted:      true,
		CompletedAt:      &completedAt,
		IsSuspended:      false,
		CreatedAt:        time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC),
		UpdatedAt:        time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC),
	}
	got := TaskToProto(task)
	if got.GetDueDate() != "2026-05-30" || got.GetCompletedAt() != "2026-05-30T12:00:00Z" || got.GetIsSuspended() {
		t.Fatalf("unexpected mapped task: %#v", got)
	}
}

func TestTaskToProtoNil(t *testing.T) {
	if TaskToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestTasksToProto(t *testing.T) {
	tasks := []domain.Task{
		{ID: 1, Title: "A", DueDate: time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Title: "B", DueDate: time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)},
	}
	got := TasksToProto(tasks)
	if len(got) != 2 || got[0].GetTitle() != "A" || got[1].GetTitle() != "B" {
		t.Fatalf("unexpected mapped tasks: %#v", got)
	}
}

func TestTaskFromCreateRequest(t *testing.T) {
	description := "Desc"
	req := &timelogv1.CreateTaskRequest{
		Title:            "Test",
		Description:      &description,
		CategoryId:       3,
		DueDate:          "2026-05-30",
		EstimatedMinutes: 60,
	}
	got, err := TaskFromCreateRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "Test" || got.Description != "Desc" || got.CategoryID != 3 || got.EstimatedMinutes != 60 {
		t.Fatalf("unexpected domain task: %#v", got)
	}
}

func TestTaskFromCreateRequestNil(t *testing.T) {
	got, err := TaskFromCreateRequest(nil)
	if err != nil || got != nil {
		t.Fatalf("expected nil, got (%v, %v)", got, err)
	}
}

func TestApplyTaskUpdate(t *testing.T) {
	task := &domain.Task{Title: "Old", CategoryID: 1}
	newTitle := "New"
	newDesc := "NewDesc"
	newCategory := int32(2)
	newDue := "2026-06-01"
	newEst := int32(30)
	completed := true
	req := &timelogv1.UpdateTaskRequest{
		Title:            &newTitle,
		Description:      &newDesc,
		CategoryId:       &newCategory,
		DueDate:          &newDue,
		EstimatedMinutes: &newEst,
		IsCompleted:      &completed,
	}
	if err := ApplyTaskUpdate(task, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if task.Title != "New" || task.Description != "NewDesc" || task.CategoryID != 2 || task.EstimatedMinutes != 30 || !task.IsCompleted {
		t.Fatalf("unexpected updated task: %#v", task)
	}
}

func TestApplyTaskUpdateNil(t *testing.T) {
	if err := ApplyTaskUpdate(nil, nil); err != nil {
		t.Fatalf("expected no error for nil inputs, got %v", err)
	}
}

func TestTaskStatsToProto(t *testing.T) {
	got, err := TaskStatsToProto(map[string]interface{}{
		"total_tasks":     int64(4),
		"completed_tasks": int64(3),
		"completion_rate": 75.0,
	})
	if err != nil {
		t.Fatalf("TaskStatsToProto() error = %v", err)
	}
	if got.GetTotalTasks() != 4 || got.GetCompletedTasks() != 3 || got.GetCompletionRate() != 75 {
		t.Fatalf("unexpected stats: %#v", got)
	}
}
