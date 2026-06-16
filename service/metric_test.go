package service

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/adapter"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/model"
)

func applyMigrations(t *testing.T, dao *model.Dao) {
	t.Helper()
	files, err := filepath.Glob("../model/migrations/*.up.sql")
	if err != nil {
		t.Fatalf("glob migrations: %v", err)
	}
	sort.Strings(files)
	rawDB := dao.RawDB
	if rawDB == nil {
		t.Fatal("raw database is nil")
	}
	if _, err := rawDB.Exec("PRAGMA foreign_keys = OFF"); err != nil {
		t.Fatalf("disable foreign keys: %v", err)
	}
	for _, f := range files {
		content, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read migration %s: %v", f, err)
		}
		if _, err := rawDB.Exec(string(content)); err != nil {
			t.Fatalf("apply migration %s: %v", f, err)
		}
	}
	if _, err := rawDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}
}

func newTestMetricService(t *testing.T) *Service {
	t.Helper()
	cfg := &config.Config{}
	cfg.Database.Host = ":memory:"
	cfg.Log.ORMLogLevel = 1

	dao, err := model.NewDao(cfg, FakeLogger{})
	if err != nil {
		t.Fatalf("NewDao: %v", err)
	}
	applyMigrations(t, dao)

	repos := adapter.NewRepositories(dao)
	return NewService(repos, repos, repos, repos, repos, repos, repos, repos, repos, FakeLogger{}, cfg, nil)
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrInt32(v int32) *int32 {
	return &v
}

func ptrString(v string) *string {
	return &v
}

func TestCreateAndGetMetric(t *testing.T) {
	svc := newTestMetricService(t)

	metric := &domain.Metric{
		Name:       "sleep_time",
		MetricType: "time",
		Unit:       "minutes",
	}
	if err := svc.CreateMetric(metric); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	if metric.ID == 0 {
		t.Fatal("expected metric ID after create")
	}

	got, err := svc.GetMetricByName("sleep_time")
	if err != nil {
		t.Fatalf("GetMetricByName: %v", err)
	}
	if got.Name != "sleep_time" {
		t.Fatalf("got name %q", got.Name)
	}
}

func TestRecordMetricUpdatesCurrentValue(t *testing.T) {
	svc := newTestMetricService(t)

	metric := &domain.Metric{
		Name:       "pushups",
		MetricType: "count",
		Unit:       "times",
	}
	if err := svc.CreateMetric(metric); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	updated, err := svc.RecordMetric(RecordMetricInput{
		MetricName: "pushups",
		Value:      20,
		Source:     "test",
		RecordedAt: "2026-06-13T08:00:00Z",
	})
	if err != nil {
		t.Fatalf("RecordMetric: %v", err)
	}
	if updated.CurrentValue == nil || *updated.CurrentValue != 20 {
		t.Fatalf("expected current value 20, got %v", updated.CurrentValue)
	}

	records, err := svc.ListMetricRecords(metric.ID)
	if err != nil {
		t.Fatalf("ListMetricRecords: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Value != 20 {
		t.Fatalf("expected record value 20, got %v", records[0].Value)
	}
}

func TestIncrementMetricAccumulates(t *testing.T) {
	svc := newTestMetricService(t)

	metric := &domain.Metric{
		Name:       "steps",
		MetricType: "count",
		Unit:       "steps",
	}
	if err := svc.CreateMetric(metric); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}

	if _, err := svc.RecordMetric(RecordMetricInput{MetricName: "steps", Value: 100, Source: "test"}); err != nil {
		t.Fatalf("RecordMetric: %v", err)
	}
	updated, err := svc.IncrementMetric(IncrementMetricInput{MetricName: "steps", Delta: 50, Source: "test"})
	if err != nil {
		t.Fatalf("IncrementMetric: %v", err)
	}
	if updated.CurrentValue == nil || *updated.CurrentValue != 150 {
		t.Fatalf("expected current value 150, got %v", updated.CurrentValue)
	}
}

func TestEvaluateConstraint(t *testing.T) {
	svc := newTestMetricService(t)

	metric := &domain.Metric{
		Name:       "sleep_time",
		MetricType: "time",
		Unit:       "minutes",
	}
	if err := svc.CreateMetric(metric); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	if _, err := svc.RecordMetric(RecordMetricInput{MetricName: "sleep_time", Value: 1380, Source: "test"}); err != nil {
		t.Fatalf("RecordMetric: %v", err)
	}

	constraint := &domain.Constraint{
		Description:       "sleep before 23:00",
		PunishmentQuote:   "stay focused",
		StartDate:         time.Now(),
		IsActive:          true,
		MetricID:          &metric.ID,
		MetricOperator:    ptrString("lte"),
		MetricTargetValue: ptrFloat64(1380),
	}
	if err := svc.CreateConstraint(constraint); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	eval, err := svc.EvaluateConstraint(constraint.ID)
	if err != nil {
		t.Fatalf("EvaluateConstraint: %v", err)
	}
	if !eval.Passed {
		t.Fatalf("expected constraint to pass")
	}
	if eval.Actual != 1380 || eval.Target != 1380 {
		t.Fatalf("unexpected actual/target values")
	}

	if _, err := svc.RecordMetric(RecordMetricInput{MetricName: "sleep_time", Value: 1430, Source: "test"}); err != nil {
		t.Fatalf("RecordMetric: %v", err)
	}
	eval, err = svc.EvaluateConstraint(constraint.ID)
	if err != nil {
		t.Fatalf("EvaluateConstraint: %v", err)
	}
	if eval.Passed {
		t.Fatalf("expected constraint to fail after exceeding target")
	}
}

func TestEvaluateConstraintWithoutMetricRule(t *testing.T) {
	svc := newTestMetricService(t)

	constraint := &domain.Constraint{
		Description:     "no metric",
		PunishmentQuote: "just do it",
		StartDate:       time.Now(),
		IsActive:        true,
	}
	if err := svc.CreateConstraint(constraint); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	_, err := svc.EvaluateConstraint(constraint.ID)
	if err == nil {
		t.Fatal("expected error for constraint without metric rule")
	}
}
