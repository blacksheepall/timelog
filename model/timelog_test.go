package model_test

import (
	"errors"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/testutil"
	. "github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/errs"
)

func TestCreateTimeLogRejectsOngoingConflict(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	first := &gen.Timelog{CategoryID: 1, StartTime: time.Now().UTC()}
	if err := CreateTimeLog(db, first); err != nil {
		t.Fatalf("CreateTimeLog first: %v", err)
	}

	second := &gen.Timelog{CategoryID: 1, StartTime: time.Now().UTC()}
	if err := CreateTimeLog(db, second); !errors.Is(err, errs.ErrOngoingTimeLogExists) {
		t.Fatalf("expected ErrOngoingTimeLogExists, got %v", err)
	}
}

func TestCreateTimeLogRejectsMissingCategory(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	tl := &gen.Timelog{CategoryID: 999, StartTime: time.Now().UTC()}
	if err := CreateTimeLog(db, tl); err == nil {
		t.Fatal("expected error for missing category foreign key")
	}
}

func TestGetTimeLogByIDNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	_, err := GetTimeLogByID(db, 9999)
	if err == nil {
		t.Fatal("expected error for non-existent timelog")
	}
}

func TestListTimeLogsByLocalDateRangeRejectsInvalidDate(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	if _, err := ListTimeLogsByLocalDateRange(db, "not-a-date", "2026-07-15"); err == nil {
		t.Fatal("expected error for invalid start date")
	}
	if _, err := ListTimeLogsByLocalDateRange(db, "2026-07-15", "not-a-date"); err == nil {
		t.Fatal("expected error for invalid end date")
	}
}

func TestUpdateTimeLogPersistsChanges(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	tl := &gen.Timelog{CategoryID: 1, StartTime: time.Now().UTC()}
	if err := CreateTimeLog(db, tl); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}

	newRemark := "updated"
	tl.Remark = &newRemark
	if err := UpdateTimeLog(db, tl); err != nil {
		t.Fatalf("UpdateTimeLog: %v", err)
	}

	got, err := GetTimeLogByID(db, *tl.ID)
	if err != nil {
		t.Fatalf("GetTimeLogByID: %v", err)
	}
	if got.Remark == nil || *got.Remark != "updated" {
		t.Fatalf("remark not persisted: %v", got.Remark)
	}
}

func TestDeleteTimeLogNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	if err := DeleteTimeLog(db, 9999); err != nil {
		t.Fatalf("expected no error deleting non-existent timelog, got %v", err)
	}
}

func TestDeleteTimeLogRemovesRecord(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	testutil.SeedTestCategory(t, dao)
	db := dao.Db()

	tl := &gen.Timelog{CategoryID: 1, StartTime: time.Now().UTC()}
	if err := CreateTimeLog(db, tl); err != nil {
		t.Fatalf("CreateTimeLog: %v", err)
	}
	if err := DeleteTimeLog(db, *tl.ID); err != nil {
		t.Fatalf("DeleteTimeLog: %v", err)
	}
	if _, err := GetTimeLogByID(db, *tl.ID); err == nil {
		t.Fatal("expected timelog to be deleted")
	}
}
