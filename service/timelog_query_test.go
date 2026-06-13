package service

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/timeutil"
)

func TestListTimeLogsByLocalDateRange(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)
	db := dao.Db()

	loc := timeutil.GetSingaporeLocation()
	start := time.Date(2026, 7, 15, 10, 0, 0, 0, loc).UTC()
	end := time.Date(2026, 7, 15, 12, 0, 0, 0, loc).UTC()

	inRange := &gen.Timelog{CategoryID: 1, StartTime: start, Remark: strPtr("in")}
	outRange := &gen.Timelog{CategoryID: 1, StartTime: time.Date(2026, 7, 14, 10, 0, 0, 0, loc).UTC(), Remark: strPtr("out")}
	endUTC := end
	inRange.EndTime = &endUTC
	outEnd := time.Date(2026, 5, 28, 12, 0, 0, 0, loc).UTC()
	outRange.EndTime = &outEnd
	if err := model.CreateTimeLog(db, inRange); err != nil {
		t.Fatalf("seed in-range log: %v", err)
	}
	if err := model.CreateTimeLog(db, outRange); err != nil {
		t.Fatalf("seed out-of-range log: %v", err)
	}

	logs, err := svc.ListTimeLogsByLocalDateRange("2026-07-15", "2026-07-15")
	if err != nil {
		t.Fatalf("ListTimeLogsByLocalDateRange: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
}

func strPtr(s string) *string { return &s }

var (
	testMigrationsMu sync.Mutex
	migratedTestDBs  = map[*sql.DB]struct{}{}
)

func applyTestMigrations(t *testing.T, dao *model.Dao) {
	t.Helper()
	rawDB := dao.RawDB
	if rawDB == nil {
		t.Fatal("raw database is nil")
	}

	testMigrationsMu.Lock()
	if _, done := migratedTestDBs[rawDB]; done {
		testMigrationsMu.Unlock()
		return
	}
	testMigrationsMu.Unlock()

	migrationFiles, err := filepath.Glob("../model/migrations/*.up.sql")
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

	testMigrationsMu.Lock()
	migratedTestDBs[rawDB] = struct{}{}
	testMigrationsMu.Unlock()
}
