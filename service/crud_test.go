package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestTaskCRUD(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)

	task := &domain.Task{Title: "Test", CategoryID: 1, DueDate: time.Now()}
	if err := svc.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	got, err := svc.GetTaskByID(task.ID)
	if err != nil || got.Title != "Test" {
		t.Fatalf("GetTaskByID: (%v, %v)", got, err)
	}

	all, err := svc.GetAllTasks(true, true)
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAllTasks: (%v, %v)", all, err)
	}

	task.Title = "Updated"
	if err := svc.UpdateTask(task); err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}

	if err := svc.MarkTaskAsCompleted(task.ID); err != nil {
		t.Fatalf("MarkTaskAsCompleted: %v", err)
	}
	if err := svc.MarkTaskAsIncomplete(task.ID); err != nil {
		t.Fatalf("MarkTaskAsIncomplete: %v", err)
	}
	if err := svc.SuspendTask(task.ID); err != nil {
		t.Fatalf("SuspendTask: %v", err)
	}
	if err := svc.UnsuspendTask(task.ID); err != nil {
		t.Fatalf("UnsuspendTask: %v", err)
	}

	if err := svc.DeleteTask(task.ID); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
}

func TestTimelogCRUD(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)

	tl := &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := svc.CreateTimeLog(tl); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}

	got, err := svc.GetTimeLogByID(tl.ID)
	if err != nil || got.CategoryID != 1 {
		t.Fatalf("GetTimeLogByID: (%v, %v)", got, err)
	}

	logs, err := svc.ListTimeLogs()
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogs: (%v, %v)", logs, err)
	}

	logs, err = svc.ListTimeLogsWithOptions(10, "start_time DESC")
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogsWithOptions: (%v, %v)", logs, err)
	}

	logs, err = svc.ListTimeLogsByLocalDateRange(time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogsByLocalDateRange: (%v, %v)", logs, err)
	}

	endTime := time.Now().UTC()
	tl.EndTime = &endTime
	if err := svc.UpdateTimeLog(tl); err != nil {
		t.Fatalf("UpdateTimeLog: %v", err)
	}

	if err := svc.DeleteTimeLog(tl.ID); err != nil {
		t.Fatalf("DeleteTimeLog: %v", err)
	}
}

func TestCategoryCRUD(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)

	cat := &domain.Category{Name: "Root"}
	if err := svc.CreateCategory(cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}

	got, err := svc.GetCategoryByID(cat.ID)
	if err != nil || got.Name != "Root" {
		t.Fatalf("GetCategoryByID: (%v, %v)", got, err)
	}

	got2, err := svc.GetCategoryByName("Root", nil)
	if err != nil || got2.Name != "Root" {
		t.Fatalf("GetCategoryByName: (%v, %v)", got2, err)
	}

	cats, err := svc.ListCategories()
	if err != nil || len(cats) < 1 {
		t.Fatalf("ListCategories: (%v, %v)", cats, err)
	}

	cats, err = svc.ListCategoriesByLevel(1)
	if err != nil || len(cats) < 1 {
		t.Fatalf("ListCategoriesByLevel: (%v, %v)", cats, err)
	}

	cats, err = svc.GetCategoriesByParentID(nil)
	if err != nil || len(cats) < 1 {
		t.Fatalf("GetCategoriesByParentID: (%v, %v)", cats, err)
	}

	tree, err := svc.GetCategoryTree()
	if err != nil || len(tree) < 1 {
		t.Fatalf("GetCategoryTree: (%v, %v)", tree, err)
	}

	cat.Color = "#FFFFFF"
	if err := svc.UpdateCategory(cat); err != nil {
		t.Fatalf("UpdateCategory: %v", err)
	}
}

func TestConstraintCRUD(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)

	c := &domain.Constraint{Description: "Test", PunishmentQuote: "Q", StartDate: time.Now(), IsActive: true}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	got, err := svc.GetConstraintByID(c.ID)
	if err != nil || got.Description != "Test" {
		t.Fatalf("GetConstraintByID: (%v, %v)", got, err)
	}

	all, err := svc.GetAllConstraints()
	if err != nil || len(all) != 1 {
		t.Fatalf("GetAllConstraints: (%v, %v)", all, err)
	}

	active, err := svc.GetActiveConstraints()
	if err != nil || len(active) < 1 {
		t.Fatalf("GetActiveConstraints: (%v, %v)", active, err)
	}

	c.Description = "Updated"
	if err := svc.UpdateConstraint(c); err != nil {
		t.Fatalf("UpdateConstraint: %v", err)
	}

	if err := svc.MarkConstraintAsCompleted(c.ID, "done"); err != nil {
		t.Fatalf("MarkConstraintAsCompleted: %v", err)
	}
	if err := svc.MarkConstraintAsActive(c.ID); err != nil {
		t.Fatalf("MarkConstraintAsActive: %v", err)
	}

	if err := svc.DeleteConstraint(c.ID); err != nil {
		t.Fatalf("DeleteConstraint: %v", err)
	}
}

func TestMetricCRUD(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)

	m := &domain.Metric{Name: "Test", MetricType: "counter", Unit: "count"}
	if err := svc.CreateMetric(m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	got, err := svc.GetMetricByID(m.ID)
	if err != nil || got.Name != "Test" {
		t.Fatalf("GetMetricByID: (%v, %v)", got, err)
	}

	metrics, err := svc.ListMetrics()
	if err != nil || len(metrics) != 1 {
		t.Fatalf("ListMetrics: (%v, %v)", metrics, err)
	}

	m.Unit = "updated"
	if err := svc.UpdateMetric(m); err != nil {
		t.Fatalf("UpdateMetric: %v", err)
	}

	if err := svc.DeleteMetric(m.ID); err != nil {
		t.Fatalf("DeleteMetric: %v", err)
	}
}

func TestCompleteTaskWithTimelog(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)

	task := &domain.Task{Title: "Complete", CategoryID: 1, DueDate: time.Now()}
	if err := svc.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	if err := svc.CompleteTaskWithTimelog(task.ID, true, &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("CompleteTaskWithTimelog: %v", err)
	}

	if err := svc.CompleteTaskWithTimelog(task.ID, false, nil); err != nil {
		t.Fatalf("CompleteTaskWithTimelog without timelog: %v", err)
	}
}

func TestListActiveTimeLogs(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)

	if err := svc.CreateTimeLog(&domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}

	logs, err := svc.ListActiveTimeLogs()
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListActiveTimeLogs: (%v, %v)", logs, err)
	}
}
