package model_test

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/testutil"
	. "github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

func seedConstraint(t *testing.T, db *gorm.DB) *gen.Constraint {
	t.Helper()
	c := &gen.Constraint{
		Description:     "daily focus",
		PunishmentQuote: "no focus penalty",
		StartDate:       time.Now().UTC(),
		IsActive:        boolPtr(true),
	}
	if err := CreateConstraint(db, c); err != nil {
		t.Fatalf("seed constraint: %v", err)
	}
	return c
}

func boolPtr(v bool) *bool { return &v }

func TestGetConstraintByIDNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	_, err := GetConstraintByID(db, 9999)
	if err == nil {
		t.Fatal("expected error for non-existent constraint")
	}
}

func TestGetActiveConstraintsEmpty(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	constraints, err := GetActiveConstraints(db)
	if err != nil {
		t.Fatalf("GetActiveConstraints: %v", err)
	}
	if len(constraints) != 0 {
		t.Fatalf("expected 0 active constraints, got %d", len(constraints))
	}
}

func TestUpdateConstraintPersistsChanges(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	c := seedConstraint(t, db)
	c.Description = "updated"
	if err := UpdateConstraint(db, c); err != nil {
		t.Fatalf("UpdateConstraint: %v", err)
	}

	got, err := GetConstraintByID(db, *c.ID)
	if err != nil {
		t.Fatalf("GetConstraintByID: %v", err)
	}
	if got.Description != "updated" {
		t.Fatalf("description not persisted: %s", got.Description)
	}
}

func TestDeleteConstraintSoftDeletes(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	c := seedConstraint(t, db)
	if err := DeleteConstraint(db, *c.ID); err != nil {
		t.Fatalf("DeleteConstraint: %v", err)
	}
	if _, err := GetConstraintByID(db, *c.ID); err == nil {
		t.Fatal("expected constraint to be soft deleted")
	}
}

func TestMarkConstraintAsCompletedNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	if err := MarkConstraintAsCompleted(db, 9999, "done"); err != nil {
		t.Fatalf("expected no error marking non-existent constraint, got %v", err)
	}
}

func TestMarkConstraintAsActiveNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	if err := MarkConstraintAsActive(db, 9999); err != nil {
		t.Fatalf("expected no error reactivating non-existent constraint, got %v", err)
	}
}
