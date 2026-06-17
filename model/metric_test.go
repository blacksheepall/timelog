package model_test

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/testutil"
	. "github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

func TestGetMetricByIDNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)

	_, err := GetMetricByID(dao.Db(), 9999)
	if err == nil {
		t.Fatal("expected error for non-existent metric ID, got nil")
	}
}

func TestGetMetricByNameNotFound(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)

	_, err := GetMetricByName(dao.Db(), "missing")
	if err == nil {
		t.Fatal("expected error for non-existent metric name, got nil")
	}
}

func TestUpdateMetricPersistsChanges(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	m := &gen.Metric{Name: "Original", MetricType: "counter", Unit: "count"}
	if err := CreateMetric(db, m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	m.Name = "Updated"
	if err := UpdateMetric(db, m); err != nil {
		t.Fatalf("UpdateMetric: %v", err)
	}

	got, err := GetMetricByID(db, *m.ID)
	if err != nil {
		t.Fatalf("GetMetricByID: %v", err)
	}
	if got.Name != "Updated" {
		t.Errorf("name: want Updated, got %s", got.Name)
	}
}

func TestDeleteMetricSoftDeletes(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	m := &gen.Metric{Name: "ToDelete", MetricType: "counter", Unit: "count"}
	if err := CreateMetric(db, m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	if err := DeleteMetric(db, *m.ID); err != nil {
		t.Fatalf("DeleteMetric: %v", err)
	}

	var softDeletedCount int64
	if err := db.Unscoped().Model(&gen.Metric{}).Where("id = ? AND deleted_at IS NOT NULL", *m.ID).Count(&softDeletedCount).Error; err != nil {
		t.Fatalf("count soft-deleted metric: %v", err)
	}
	if softDeletedCount != 1 {
		t.Errorf("expected 1 soft-deleted metric, got %d", softDeletedCount)
	}

	var visibleCount int64
	if err := db.Model(&gen.Metric{}).Where("id = ?", *m.ID).Count(&visibleCount).Error; err != nil {
		t.Fatalf("count visible metric: %v", err)
	}
	if visibleCount != 0 {
		t.Errorf("expected 0 visible metric, got %d", visibleCount)
	}
}

func TestListMetricRecordsByMetricIDEmpty(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)

	records, err := ListMetricRecordsByMetricID(dao.Db(), 9999)
	if err != nil {
		t.Fatalf("ListMetricRecordsByMetricID: %v", err)
	}
	if len(records) != 0 {
		t.Errorf("expected empty records, got %d", len(records))
	}
}

func TestUpdateMetricCurrentValuePersists(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	testutil.ApplyMigrations(t, dao)
	db := dao.Db()

	m := &gen.Metric{Name: "ValueMetric", MetricType: "numeric", Unit: "kg"}
	if err := CreateMetric(db, m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	recordedAt := time.Now().UTC().Truncate(time.Second)
	if err := UpdateMetricCurrentValue(db, *m.ID, 42.5, recordedAt); err != nil {
		t.Fatalf("UpdateMetricCurrentValue: %v", err)
	}

	got, err := GetMetricByID(db, *m.ID)
	if err != nil {
		t.Fatalf("GetMetricByID: %v", err)
	}
	if got.CurrentValue == nil || *got.CurrentValue != 42.5 {
		t.Errorf("current_value: want 42.5, got %v", got.CurrentValue)
	}
	if got.LastRecordedAt == nil || !got.LastRecordedAt.Equal(recordedAt) {
		t.Errorf("last_recorded_at: want %v, got %v", recordedAt, got.LastRecordedAt)
	}
}
