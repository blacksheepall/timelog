package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/timeutil"
)

func TestCreateTimeLogFromMCPInputSetsUserAndParsesSGT(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)

	input := CreateTimeLogMCPInput{
		CategoryID: 1,
		StartTime:  "2026-05-29 10:00:00",
		EndTime:    "2026-05-29 11:00:00",
		Remark:     "mcp create",
	}

	created, err := svc.CreateTimeLogFromMCPInput(input)
	if err != nil {
		t.Fatalf("CreateTimeLogFromMCPInput: %v", err)
	}
	if created.UserID == nil || *created.UserID != 1 {
		t.Fatalf("expected default user_id 1, got %+v", created.UserID)
	}
	loc := timeutil.GetSingaporeLocation()
	wantStart := time.Date(2026, 5, 29, 10, 0, 0, 0, loc).UTC()
	if !created.StartTime.Equal(wantStart) {
		t.Fatalf("start_time mismatch: got %v want %v", created.StartTime, wantStart)
	}
}

func TestUpdateTimeLogFromMCPInputPartialRemark(t *testing.T) {
	svc, dao := setupTestModel()
	applyTestMigrations(t, dao)
	seedTestCategory(t, dao)
	db := dao.Db()

	remark := "before"
	end := time.Now().UTC()
	seed := &gen.Timelog{CategoryID: 1, StartTime: time.Now().UTC(), EndTime: &end, Remark: &remark}
	if err := model.CreateTimeLog(db, seed); err != nil {
		t.Fatalf("seed: %v", err)
	}

	updated, err := svc.UpdateTimeLogFromMCPInput(UpdateTimeLogMCPInput{
		ID:     *seed.ID,
		Remark: "after",
	})
	if err != nil {
		t.Fatalf("UpdateTimeLogFromMCPInput: %v", err)
	}
	if updated.Remark == nil || *updated.Remark != "after" {
		t.Fatalf("remark not updated: %+v", updated.Remark)
	}
}

func seedTestCategory(t *testing.T, dao *model.Dao) {
	t.Helper()
	level := int32(1)
	cat := &gen.Category{Name: "test", Level: &level}
	if err := model.CreateCategory(dao.Db(), cat); err != nil {
		t.Fatalf("seed category: %v", err)
	}
}
