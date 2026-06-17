package model_test

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/testutil"
	. "github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

func seedTask(t *testing.T, db *gorm.DB) *gen.Task {
	t.Helper()
	task := &gen.Task{
		Title:            "test task",
		CategoryID:       1,
		DueDate:          time.Now().UTC(),
		EstimatedMinutes: 30,
	}
	if err := CreateTask(db, task); err != nil {
		t.Fatalf("seed task: %v", err)
	}
	return task
}

func TestCreateTaskRejectsMissingCategory(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	task := &gen.Task{Title: "orphan", CategoryID: 999, DueDate: time.Now().UTC()}
	if err := CreateTask(db, task); err == nil {
		t.Fatal("expected error for missing category foreign key")
	}
}

func TestGetTaskByIDNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	_, err := GetTaskByID(db, 9999)
	if err == nil {
		t.Fatal("expected error for non-existent task")
	}
}

func TestGetAllTasksFilters(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	active := seedTask(t, db)
	suspended := seedTask(t, db)
	suspended.Title = "suspended"
	if err := SuspendTask(db, *suspended.ID); err != nil {
		t.Fatalf("suspend: %v", err)
	}
	completed := seedTask(t, db)
	completed.Title = "completed"
	if err := MarkTaskAsCompleted(db, *completed.ID); err != nil {
		t.Fatalf("complete: %v", err)
	}

	tasks, err := GetAllTasks(db, false, false)
	if err != nil {
		t.Fatalf("GetAllTasks: %v", err)
	}
	if len(tasks) != 1 || *tasks[0].ID != *active.ID {
		t.Fatalf("expected only active task, got %d", len(tasks))
	}
}

func TestUpdateTaskPersistsChanges(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	task := seedTask(t, db)
	newTitle := "updated"
	task.Title = newTitle
	if err := UpdateTask(db, task); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	got, err := GetTaskByID(db, *task.ID)
	if err != nil {
		t.Fatalf("GetTaskByID: %v", err)
	}
	if got.Title != "updated" {
		t.Fatalf("title not persisted: %s", got.Title)
	}
}

func TestDeleteTaskSoftDeletes(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	task := seedTask(t, db)
	if err := DeleteTask(db, *task.ID); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
	if _, err := GetTaskByID(db, *task.ID); err == nil {
		t.Fatal("expected task to be soft deleted")
	}
	var count int64
	db.Unscoped().Model(&gen.Task{}).Where("id = ?", *task.ID).Count(&count)
	if count != 1 {
		t.Fatalf("expected soft-deleted record to remain, got count %d", count)
	}
}

func TestMarkTaskAsCompletedNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	if err := MarkTaskAsCompleted(db, 9999); err != nil {
		t.Fatalf("expected no error marking non-existent task, got %v", err)
	}
}

func TestSuspendAndUnsuspendTask(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	task := seedTask(t, db)
	if err := SuspendTask(db, *task.ID); err != nil {
		t.Fatalf("SuspendTask: %v", err)
	}
	got, _ := GetTaskByID(db, *task.ID)
	if got.IsSuspended == nil || !*got.IsSuspended {
		t.Fatal("expected task suspended")
	}
	if err := UnsuspendTask(db, *task.ID); err != nil {
		t.Fatalf("UnsuspendTask: %v", err)
	}
	got, _ = GetTaskByID(db, *task.ID)
	if got.IsSuspended == nil || *got.IsSuspended {
		t.Fatal("expected task not suspended")
	}
}

func TestGetTaskStatsEmpty(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	stats, err := GetTaskStats(db, time.Now().UTC())
	if err != nil {
		t.Fatalf("GetTaskStats: %v", err)
	}
	if stats["total_tasks"].(int64) != 0 || stats["completed_tasks"].(int64) != 0 {
		t.Fatalf("expected empty stats, got %v", stats)
	}
	if stats["completion_rate"].(float64) != 0 {
		t.Fatalf("expected 0 completion rate, got %v", stats["completion_rate"])
	}
}

func TestGetCompletedTasksInDateRangeEmpty(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	now := time.Now().UTC()
	tasks, err := GetCompletedTasksInDateRange(db, now.Add(-24*time.Hour), now)
	if err != nil {
		t.Fatalf("GetCompletedTasksInDateRange: %v", err)
	}
	if len(tasks) != 0 {
		t.Fatalf("expected 0 tasks, got %d", len(tasks))
	}
}
