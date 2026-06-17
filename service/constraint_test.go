package service

import (
	"strings"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestGetConstraintByIDNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	_, err := svc.GetConstraintByID(9999)
	if err == nil {
		t.Fatal("expected error for non-existent constraint, got nil")
	}
}

func TestEvaluateConstraintMissingMetricRule(t *testing.T) {
	svc, _ := newTestService(t)

	c := &domain.Constraint{
		Description:     "No metric rule",
		PunishmentQuote: "Ouch",
		StartDate:       time.Now(),
		IsActive:        true,
	}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	_, err := svc.EvaluateConstraint(c.ID)
	if err == nil || !strings.Contains(err.Error(), "constraint has no metric rule") {
		t.Fatalf("expected 'constraint has no metric rule' error, got %v", err)
	}
}

func TestEvaluateConstraintMetricNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	c := &domain.Constraint{
		Description:       "Missing metric",
		PunishmentQuote:   "Ouch",
		StartDate:         time.Now(),
		IsActive:          true,
		MetricID:          int32Ptr(9999),
		MetricOperator:    strPtr("gte"),
		MetricTargetValue: float64Ptr(100),
	}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	_, err := svc.EvaluateConstraint(c.ID)
	if err == nil {
		t.Fatal("expected error when metric is not found, got nil")
	}
}

func TestEvaluateConstraintUnsupportedOperator(t *testing.T) {
	svc, _ := newTestService(t)

	m := &domain.Metric{Name: "pushups", MetricType: "counter", Unit: "count"}
	if err := svc.CreateMetric(m); err != nil {
		t.Fatalf("CreateMetric: %v", err)
	}
	m.CurrentValue = float64Ptr(10)
	if err := svc.UpdateMetric(m); err != nil {
		t.Fatalf("UpdateMetric: %v", err)
	}

	c := &domain.Constraint{
		Description:       "Unsupported operator",
		PunishmentQuote:   "Ouch",
		StartDate:         time.Now(),
		IsActive:          true,
		MetricID:          int32Ptr(m.ID),
		MetricOperator:    strPtr("unsupported"),
		MetricTargetValue: float64Ptr(5),
	}
	if err := svc.CreateConstraint(c); err != nil {
		t.Fatalf("CreateConstraint: %v", err)
	}

	_, err := svc.EvaluateConstraint(c.ID)
	if err == nil || !strings.Contains(err.Error(), "unsupported metric operator") {
		t.Fatalf("expected 'unsupported metric operator' error, got %v", err)
	}
}

func TestDeleteConstraintNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	if err := svc.DeleteConstraint(9999); err != nil {
		t.Fatalf("expected no error deleting non-existent constraint, got %v", err)
	}
}

func TestMarkConstraintAsCompletedNotFound(t *testing.T) {
	svc, _ := newTestService(t)

	if err := svc.MarkConstraintAsCompleted(9999, "done"); err != nil {
		t.Fatalf("expected no error marking non-existent constraint completed, got %v", err)
	}
}
