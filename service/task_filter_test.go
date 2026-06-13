package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

func TestListTasksByCompletionStatusFilters(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)
	db := dao.Db()

	completed := true
	pending := &gen.Task{
		Title:       "pending",
		CategoryID:  1,
		DueDate:     time.Now(),
		IsCompleted: nil,
	}
	done := &gen.Task{
		Title:       "done",
		CategoryID:  1,
		DueDate:     time.Now(),
		IsCompleted: &completed,
	}
	if err := model.CreateTask(db, pending); err != nil {
		t.Fatalf("seed pending: %v", err)
	}
	if err := model.CreateTask(db, done); err != nil {
		t.Fatalf("seed done: %v", err)
	}

	pendingOnly, err := svc.ListTasksByCompletionStatus("pending")
	if err != nil {
		t.Fatalf("pending: %v", err)
	}
	if len(pendingOnly) != 1 || pendingOnly[0].Title != "pending" {
		t.Fatalf("unexpected pending result: %+v", pendingOnly)
	}

	completedOnly, err := svc.ListTasksByCompletionStatus("completed")
	if err != nil {
		t.Fatalf("completed: %v", err)
	}
	if len(completedOnly) != 1 || completedOnly[0].Title != "done" {
		t.Fatalf("unexpected completed result: %+v", completedOnly)
	}
}
