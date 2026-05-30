package mapper

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

func TestTaskToProtoFormatsDateAndOptionalTime(t *testing.T) {
	id := int32(5)
	description := "Write tests"
	completed := true
	suspended := false
	completedAt := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	createdAt := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 5, 30, 13, 0, 0, 0, time.UTC)
	task := &gen.Task{
		ID:               &id,
		Title:            "IDL",
		Description:      &description,
		CategoryID:       3,
		DueDate:          time.Date(2026, 5, 30, 0, 0, 0, 0, time.UTC),
		EstimatedMinutes: 90,
		IsCompleted:      &completed,
		CompletedAt:      &completedAt,
		IsSuspended:      &suspended,
		CreatedAt:        &createdAt,
		UpdatedAt:        &updatedAt,
	}
	got := TaskToProto(task)
	if got.GetDueDate() != "2026-05-30" || got.GetCompletedAt() != "2026-05-30T12:00:00Z" || got.GetIsSuspended() {
		t.Fatalf("unexpected mapped task: %#v", got)
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
