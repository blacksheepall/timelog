package testutil

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"testing"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/core/logger"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// FakeLogger is a no-op logger for tests.
type FakeLogger struct{}

var _ logger.Logger = (*FakeLogger)(nil)

func (FakeLogger) Debug(...interface{})                      {}
func (FakeLogger) Debugw(string, ...interface{})             {}
func (FakeLogger) Info(...interface{})                       {}
func (FakeLogger) Infow(string, ...interface{})              {}
func (FakeLogger) Warn(...interface{})                       {}
func (FakeLogger) Warnw(string, ...interface{})              {}
func (FakeLogger) Error(...interface{})                      {}
func (FakeLogger) Errorw(string, ...interface{})             {}
func (FakeLogger) Fatal(...interface{})                      {}
func (FakeLogger) Fatalw(string, ...interface{})             {}
func (FakeLogger) WithContext(context.Context) logger.Logger { return FakeLogger{} }

// NewTestConfig returns a test config using an in-memory SQLite database.
func NewTestConfig() *config.Config {
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1
	return cfg
}

// NewTestDAO creates and returns an in-memory SQLite DAO, failing the test on error.
func NewTestDAO(t *testing.T) *model.Dao {
	t.Helper()
	cfg := NewTestConfig()
	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	return dao
}

// ApplyMigrations runs model/migrations/*.up.sql against the provided DAO.
func ApplyMigrations(t *testing.T, dao *model.Dao) {
	t.Helper()
	rawDB := dao.RawDB
	if rawDB == nil {
		t.Fatal("raw database is nil")
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to determine testutil source path")
	}
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "..", "model", "migrations")
	migrationFiles, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		t.Fatalf("list migrations: %v", err)
	}
	if len(migrationFiles) == 0 {
		t.Fatalf("no migration files found in %s", migrationsDir)
	}
	sort.Strings(migrationFiles)

	for _, f := range migrationFiles {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read migration %s: %v", f, err)
		}
		if _, err := rawDB.Exec(string(content)); err != nil {
			t.Fatalf("apply migration %s: %v", f, err)
		}
	}
}

// SeedTestCategory inserts a root category named "test", failing the test on error.
func SeedTestCategory(t *testing.T, dao *model.Dao) *gen.Category {
	t.Helper()
	name := "test"
	path := "/"
	level := int32(1)
	cat := &gen.Category{
		Name:  name,
		Path:  &path,
		Level: &level,
	}
	if err := model.CreateCategory(dao.Db(), cat); err != nil {
		t.Fatalf("seed category: %v", err)
	}
	return cat
}
