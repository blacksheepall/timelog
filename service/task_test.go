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

func TestGetTasksByDateAndRange(t *testing.T) {
	svc, dao := newTestService(t)
	testutil.SeedTestCategory(t, dao)

	today := time.Now().UTC()
	todayStart := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	tomorrowStart := todayStart.Add(24 * time.Hour)

	// Task due today
	todayTask := seedServiceTask(t, svc)
	todayTask.DueDate = todayStart.Add(1 * time.Hour)
	if err := svc.UpdateTask(todayTask); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	// Task due tomorrow
	tomorrowTask := seedServiceTask(t, svc)
	tomorrowTask.Title = "tomorrow"
	tomorrowTask.DueDate = tomorrowStart.Add(1 * time.Hour)
	if err := svc.UpdateTask(tomorrowTask); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	byDate, err := svc.GetTasksByDate(todayStart, true, true)
	if err != nil {
		t.Fatalf("GetTasksByDate: %v", err)
	}
	if len(byDate) != 1 || byDate[0].ID != todayTask.ID {
		t.Fatalf("expected 1 task today, got %d", len(byDate))
	}

	byRange, err := svc.GetTasksByDateRange(todayStart, tomorrowStart.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("GetTasksByDateRange: %v", err)
	}
	if len(byRange) != 2 {
		t.Fatalf("expected 2 tasks in range, got %d", len(byRange))
	}

	completed, err := svc.GetCompletedTasksInDateRange(todayStart, tomorrowStart.Add(24*time.Hour))
	if err != nil {
		t.Fatalf("GetCompletedTasksInDateRange: %v", err)
	}
	if len(completed) != 0 {
		t.Fatalf("expected 0 completed tasks, got %d", len(completed))
	}

	stats, err := svc.GetTaskStats(todayStart)
	if err != nil {
		t.Fatalf("GetTaskStats: %v", err)
	}
	if stats["total_tasks"].(int64) != 1 {
		t.Fatalf("expected total 1, got %v", stats["total_tasks"])
	}
}
