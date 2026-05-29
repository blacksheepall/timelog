package service

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

func TestListTimeLogsByLocalDateRange(t *testing.T) {
	dao := setupTestModel()
	applyTestMigrations(t, dao)
	db := dao.Db()

	loc := model.GetSingaporeLocation()
	start := time.Date(2026, 5, 29, 10, 0, 0, 0, loc).UTC()
	end := time.Date(2026, 5, 29, 12, 0, 0, 0, loc).UTC()

	inRange := &gen.Timelog{CategoryID: 1, StartTime: start, Remark: strPtr("in")}
	outRange := &gen.Timelog{CategoryID: 1, StartTime: time.Date(2026, 5, 28, 10, 0, 0, 0, loc).UTC(), Remark: strPtr("out")}
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

	logs, err := ListTimeLogsByLocalDateRange("2026-05-29", "2026-05-29")
	if err != nil {
		t.Fatalf("ListTimeLogsByLocalDateRange: %v", err)
	}
	if len(logs) != 1 {
		t.Fatalf("expected 1 log, got %d", len(logs))
	}
}

func strPtr(s string) *string { return &s }

func applyTestMigrations(t *testing.T, dao *model.Dao) {
	t.Helper()
	rawDB := dao.RawDB
	if rawDB == nil {
		t.Fatal("raw database is nil")
	}

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
}
