package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/testutil"
)

func seedServiceTask(t *testing.T, svc *Service) *domain.Task {
	t.Helper()
	task := &domain.Task{
		Title:            "test",
		CategoryID:       1,
		DueDate:          time.Now().UTC(),
		EstimatedMinutes: 30,
	}
	if err := svc.CreateTask(task); err != nil {
		t.Fatalf("seed task: %v", err)
	}
	return task
}

func TestCreateTaskRejectsMissingCategory(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.CreateTask(&domain.Task{Title: "orphan", CategoryID: 999, DueDate: time.Now().UTC()}); err == nil {
		t.Fatal("expected error for missing category")
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	_, err := svc.GetTaskByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestListTasksByCompletionStatusInvalid(t *testing.T) {
	svc, _ := newTestService(t)
	if _, err := svc.ListTasksByCompletionStatus("invalid"); err == nil {
		t.Fatal("expected error for invalid status")
	}
}

func TestDeleteTaskNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.DeleteTask(9999); err != nil {
		t.Fatalf("expected no error deleting non-existent task, got %v", err)
	}
}

func TestMarkTaskAsCompletedNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	if err := svc.MarkTaskAsCompleted(9999); err != nil {
		t.Fatalf("expected no error marking non-existent task, got %v", err)
	}
}

func TestCompleteTaskWithTimelogInvalidState(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	task := seedServiceTask(t, svc)
	if err := svc.CompleteTaskWithTimelog(task.ID, true, nil); err == nil {
		t.Fatal("expected error when createTimelog true but timelogData nil")
	}
}
