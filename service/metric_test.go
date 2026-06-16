package service

import (
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestCreateAndGetMetric(t *testing.T) {
	svc, _ := newTestService(t)

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
	svc, _ := newTestService(t)

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
	svc, _ := newTestService(t)

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
	svc, _ := newTestService(t)

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
		MetricOperator:    strPtr("lte"),
		MetricTargetValue: float64Ptr(1380),
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

func TestEvaluateMetricOperators(t *testing.T) {
	cases := []struct {
		op     string
		actual float64
		target float64
		want   bool
	}{
		{"eq", 5, 5, true},
		{"eq", 5, 4, false},
		{"ne", 5, 4, true},
		{"ne", 5, 5, false},
		{"gt", 5, 4, true},
		{"gt", 4, 5, false},
		{"gte", 5, 5, true},
		{"gte", 4, 5, false},
		{"lt", 4, 5, true},
		{"lt", 5, 4, false},
		{"lte", 5, 5, true},
		{"lte", 5, 4, false},
	}
	for _, tc := range cases {
		got, err := evaluateMetric(tc.actual, tc.target, tc.op)
		if err != nil {
			t.Fatalf("evaluateMetric(%s): %v", tc.op, err)
		}
		if got != tc.want {
			t.Fatalf("evaluateMetric(%s, %v, %v) = %v, want %v", tc.op, tc.actual, tc.target, got, tc.want)
		}
	}

	if _, err := evaluateMetric(0, 0, "unknown"); err == nil {
		t.Fatal("expected error for unknown operator")
	}
}

func TestParseMetricRecordedAt(t *testing.T) {
	got, err := parseMetricRecordedAt("2026-06-13T08:00:00Z")
	if err != nil || got.UTC().Format(time.RFC3339) != "2026-06-13T08:00:00Z" {
		t.Fatalf("unexpected: (%v, %v)", got, err)
	}

	got, err = parseMetricRecordedAt("2026-06-13 08:00:00")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	emptyStart := time.Now().UTC()
	got, err = parseMetricRecordedAt("")
	if err != nil || got.Before(emptyStart) {
		t.Fatalf("expected current time, got (%v, %v)", got, err)
	}

	if _, err := parseMetricRecordedAt("not-a-time"); err == nil {
		t.Fatal("expected error for invalid time")
	}
}

func TestEvaluateConstraintWithoutMetricRule(t *testing.T) {
	svc, _ := newTestService(t)

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
