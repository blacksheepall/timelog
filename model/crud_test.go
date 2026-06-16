package model

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	sqlite "github.com/ncruces/go-sqlite3/gormlite"
	"gorm.io/gorm"
)

func openMemoryDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	return db
}

func TestTimelogCRUD(t *testing.T) {
	db := openMemoryDB(t)
	if err := db.AutoMigrate(&gen.Timelog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	tl := &gen.Timelog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := CreateTimeLog(db, tl); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}
	if tl.ID == nil || *tl.ID == 0 {
		t.Fatal("expected ID")
	}

	got, err := GetTimeLogByID(db, *tl.ID)
	if err != nil || got.CategoryID != 1 {
		t.Fatalf("GetTimeLogByID: (%v, %v)", got, err)
	}

	logs, err := ListTimeLogs(db)
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogs: (%v, %v)", logs, err)
	}

	logs, err = ListTimeLogsWithOptions(db, 10, "start_time DESC")
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogsWithOptions: (%v, %v)", logs, err)
	}

	logs, err = ListTimeLogsByLocalDateRange(db, time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
	if err != nil || len(logs) != 1 {
		t.Fatalf("ListTimeLogsByLocalDateRange: (%v, %v)", logs, err)
	}

	endTime := time.Now().UTC()
	tl.EndTime = &endTime
	if err := UpdateTimeLog(db, tl); err != nil {
		t.Fatalf("UpdateTimeLog: %v", err)
	}

	if err := DeleteTimeLog(db, *tl.ID); err != nil {
		t.Fatalf("DeleteTimeLog: %v", err)
	}
}

func TestTaskCRUD(t *testing.T) {
	db := openMemoryDB(t)
	if err := db.AutoMigrate(&gen.Task{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	task := &gen.Task{Title: "Test", CategoryID: 1, DueDate: time.Now()}
	if err := CreateTask(db, task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}

	got, err := GetTaskByID(db, *task.ID)
	if err != nil || got.Title != "Test" {
		t.Fatalf("GetTaskByID: (%v, %v)", got, err)
	}

	tasks, err := GetAllTasks(db, true, true)
	if err != nil || len(tasks) != 1 {
		t.Fatalf("GetAllTasks: (%v, %v)", tasks, err)
	}

	tasks, err = GetTasksByDate(db, time.Now(), true, true)
	if err != nil || len(tasks) != 1 {
		t.Fatalf("GetTasksByDate: (%v, %v)", tasks, err)
	}

	tasks, err = GetTasksByDateRange(db, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	if err != nil || len(tasks) != 1 {
		t.Fatalf("GetTasksByDateRange: (%v, %v)", tasks, err)
	}

	if err := MarkTaskAsCompleted(db, *task.ID); err != nil {
		t.Fatalf("MarkTaskAsCompleted: %v", err)
	}
	if err := MarkTaskAsIncomplete(db, *task.ID); err != nil {
		t.Fatalf("MarkTaskAsIncomplete: %v", err)
	}
	if err := SuspendTask(db, *task.ID); err != nil {
		t.Fatalf("SuspendTask: %v", err)
	}
	if err := UnsuspendTask(db, *task.ID); err != nil {
		t.Fatalf("UnsuspendTask: %v", err)
	}

	if err := DeleteTask(db, *task.ID); err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
}

func TestConstraintCRUD(t *testing.T) {
	db := openMemoryDB(t)
	if err := db.AutoMigrate(&gen.Constraint{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	c := &gen.Constraint{Description: "Test", PunishmentQuote: "Q", StartDate: time.Now()}
	if err := CreateConstraint(db, c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	got, err := GetConstraintByID(db, *c.ID)
	if err != nil || got.Description != "Test" {
		t.Fatalf("GetConstraintByID: (%v, %v)", got, err)
	}

	constraints, err := GetAllConstraints(db)
	if err != nil || len(constraints) != 1 {
		t.Fatalf("GetAllConstraints: (%v, %v)", constraints, err)
	}

	constraints, err = GetConstraintsByDateRange(db, time.Now().Add(-time.Hour), time.Now().Add(time.Hour))
	if err != nil || len(constraints) != 1 {
		t.Fatalf("GetConstraintsByDateRange: (%v, %v)", constraints, err)
	}

	if err := MarkConstraintAsCompleted(db, *c.ID, "done"); err != nil {
		t.Fatalf("MarkConstraintAsCompleted: %v", err)
	}
	if err := MarkConstraintAsActive(db, *c.ID); err != nil {
		t.Fatalf("MarkConstraintAsActive: %v", err)
	}

	if err := DeleteConstraint(db, *c.ID); err != nil {
		t.Fatalf("DeleteConstraint: %v", err)
	}
}

func TestMetricCRUD(t *testing.T) {
	db := openMemoryDB(t)
	if err := db.AutoMigrate(&gen.Metric{}, &gen.MetricRecord{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	m := &gen.Metric{Name: "Test", MetricType: "counter", Unit: "count"}
	if err := CreateMetric(db, m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	got, err := GetMetricByID(db, *m.ID)
	if err != nil || got.Name != "Test" {
		t.Fatalf("GetMetricByID: (%v, %v)", got, err)
	}

	got2, err := GetMetricByName(db, "Test")
	if err != nil || got2.Name != "Test" {
		t.Fatalf("GetMetricByName: (%v, %v)", got2, err)
	}

	metrics, err := ListMetrics(db)
	if err != nil || len(metrics) != 1 {
		t.Fatalf("ListMetrics: (%v, %v)", metrics, err)
	}

	record := &gen.MetricRecord{MetricID: *m.ID, Value: 10}
	if err := CreateMetricRecord(db, record); err != nil {
		t.Fatalf("CreateMetricRecord: %v", err)
	}

	records, err := ListMetricRecordsByMetricID(db, *m.ID)
	if err != nil || len(records) != 1 {
		t.Fatalf("ListMetricRecordsByMetricID: (%v, %v)", records, err)
	}

	if err := UpdateMetricCurrentValue(db, *m.ID, 10, time.Now()); err != nil {
		t.Fatalf("UpdateMetricCurrentValue: %v", err)
	}

	m.Unit = "updated"
	if err := UpdateMetric(db, m); err != nil {
		t.Fatalf("UpdateMetric: %v", err)
	}

	if err := DeleteMetric(db, *m.ID); err != nil {
		t.Fatalf("DeleteMetric: %v", err)
	}
}
