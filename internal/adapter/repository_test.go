package adapter

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

type fakeLogger struct{}

func (fakeLogger) Debug(...interface{})          {}
func (fakeLogger) Debugw(string, ...interface{}) {}
func (fakeLogger) Info(...interface{})           {}
func (fakeLogger) Infow(string, ...interface{})  {}
func (fakeLogger) Warn(...interface{})           {}
func (fakeLogger) Warnw(string, ...interface{})  {}
func (fakeLogger) Error(...interface{})          {}
func (fakeLogger) Errorw(string, ...interface{}) {}
func (fakeLogger) Fatal(...interface{})          {}
func (fakeLogger) Fatalw(string, ...interface{}) {}

func setupTestDao(t *testing.T) *model.Dao {
	t.Helper()
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	dao, err := model.NewDao(cfg, fakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	return dao
}

var (
	adapterMigrationsMu sync.Mutex
	adapterMigratedDBs  = map[*sql.DB]struct{}{}
)

func applyMigrations(t *testing.T, dao *model.Dao) {
	t.Helper()
	rawDB := dao.RawDB
	if rawDB == nil {
		t.Fatal("raw database is nil")
	}

	adapterMigrationsMu.Lock()
	if _, done := adapterMigratedDBs[rawDB]; done {
		adapterMigrationsMu.Unlock()
		return
	}
	adapterMigrationsMu.Unlock()

	migrationFiles, err := filepath.Glob("../../model/migrations/*.up.sql")
	if err != nil {
		t.Fatalf("list migrations: %v", err)
	}
	sort.Strings(migrationFiles)

	for _, migrationFile := range migrationFiles {
		content, err := os.ReadFile(migrationFile)
		if err != nil {
			t.Fatalf("read migration %s: %v", migrationFile, err)
		}
		if _, err := rawDB.Exec(string(content)); err != nil {
			t.Fatalf("apply migration %s: %v", migrationFile, err)
		}
	}

	adapterMigrationsMu.Lock()
	adapterMigratedDBs[rawDB] = struct{}{}
	adapterMigrationsMu.Unlock()
}

func seedCategory(t *testing.T, dao *model.Dao) {
	t.Helper()
	repo := newCategoryRepo(dao)
	cat := &domain.Category{Name: "test"}
	if err := repo.CreateCategory(cat); err != nil {
		t.Fatalf("seed category: %v", err)
	}
}

func TestNewRepositories(t *testing.T) {
	repos := NewRepositories(setupTestDao(t))
	if repos == nil {
		t.Fatal("expected non-nil repositories")
	}
}

func TestCategoryRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newCategoryRepo(dao)
	cat := &domain.Category{Name: "Root"}
	if err := repo.CreateCategory(cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	if cat.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestTaskRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	seedCategory(t, dao)
	repo := newTaskRepo(dao)
	task := &domain.Task{Title: "Test", CategoryID: 1, DueDate: time.Now()}
	if err := repo.CreateTask(task); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if task.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestTimelogRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	seedCategory(t, dao)
	repo := newTimelogRepo(dao)
	tl := &domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}
	if err := repo.CreateTimeLog(tl); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}
	if tl.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestConstraintRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newConstraintRepo(dao)
	c := &domain.Constraint{Description: "Test", PunishmentQuote: "Q", StartDate: time.Now()}
	if err := repo.CreateConstraint(c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}
	if c.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestMetricRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newMetricRepo(dao)
	m := &domain.Metric{Name: "Pages", MetricType: "counter", Unit: "page"}
	if err := repo.CreateMetric(m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	if m.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestMetricRecordRepoCreate(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newMetricRepo(dao)
	m := &domain.Metric{Name: "Pages", MetricType: "counter", Unit: "page"}
	if err := repo.CreateMetric(m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	record := &domain.MetricRecord{MetricID: m.ID, Value: 10}
	if err := repo.CreateMetricRecord(record); err != nil {
		t.Fatalf("CreateMetricRecord: %v", err)
	}
	if record.ID == 0 {
		t.Fatal("expected generated ID to be written back")
	}
}

func TestCategoryRepoGetByID(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newCategoryRepo(dao)
	cat := &domain.Category{Name: "Root"}
	if err := repo.CreateCategory(cat); err != nil {
		t.Fatalf("CreateCategory: %v", err)
	}
	got, err := repo.GetCategoryByID(cat.ID)
	if err != nil || got.Name != "Root" {
		t.Fatalf("GetCategoryByID: (%v, %v)", got, err)
	}
}

func TestTaskRepoGetAll(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	seedCategory(t, dao)
	repo := newTaskRepo(dao)
	if err := repo.CreateTask(&domain.Task{Title: "A", CategoryID: 1, DueDate: time.Now()}); err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	got, err := repo.GetAllTasks(true, true)
	if err != nil || len(got) != 1 {
		t.Fatalf("GetAllTasks: (%v, %v)", got, err)
	}
}

func TestTimelogRepoList(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	seedCategory(t, dao)
	repo := newTimelogRepo(dao)
	if err := repo.CreateTimeLog(&domain.TimeLog{StartTime: time.Now().UTC(), CategoryID: 1}); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}
	got, err := repo.ListTimeLogs()
	if err != nil || len(got) != 1 {
		t.Fatalf("ListTimeLogs: (%v, %v)", got, err)
	}
}

func TestConstraintRepoGetAll(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newConstraintRepo(dao)
	if err := repo.CreateConstraint(&domain.Constraint{Description: "A", PunishmentQuote: "Q", StartDate: time.Now()}); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}
	got, err := repo.GetAllConstraints()
	if err != nil || len(got) != 1 {
		t.Fatalf("GetAllConstraints: (%v, %v)", got, err)
	}
}

func TestMetricRepoList(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	repo := newMetricRepo(dao)
	if err := repo.CreateMetric(&domain.Metric{Name: "A", MetricType: "counter", Unit: "count"}); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	got, err := repo.ListMetrics()
	if err != nil || len(got) != 1 {
		t.Fatalf("ListMetrics: (%v, %v)", got, err)
	}
}

func TestCacheStore(t *testing.T) {
	dao := setupTestDao(t)
	store := newCacheStore(dao)
	store.WriteCache("key", "value", 60)
	got, ok := store.GetCache("key")
	if !ok || got != "value" {
		t.Fatalf("expected cached value, got (%v, %v)", got, ok)
	}
}

func TestUnitOfWorkRun(t *testing.T) {
	dao := setupTestDao(t)
	applyMigrations(t, dao)
	seedCategory(t, dao)
	uow := newUnitOfWork(dao)
	err := uow.Run(nil, func(repos ports.UnitOfWorkRepositories) error {
		task := &domain.Task{Title: "UOW", CategoryID: 1, DueDate: time.Now()}
		return repos.TaskRepo.CreateTask(task)
	})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
}
