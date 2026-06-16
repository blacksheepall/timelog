package mapper

import (
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/service"
)

func TestConstraintToProtoFormatsDates(t *testing.T) {
	endDate := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	metricID := int32(5)
	op := "gt"
	target := 100.0
	c := &domain.Constraint{
		ID:                9,
		Description:       "No social media",
		EndReason:         "done",
		PunishmentQuote:   "Pay the price",
		StartDate:         time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		EndDate:           &endDate,
		IsActive:          false,
		MetricID:          &metricID,
		MetricOperator:    &op,
		MetricTargetValue: &target,
		CreatedAt:         time.Date(2026, 5, 1, 8, 0, 0, 0, time.UTC),
		UpdatedAt:         time.Date(2026, 5, 2, 8, 0, 0, 0, time.UTC),
	}
	got := ConstraintToProto(c)
	if got.GetStartDate() != "2026-05-01" || got.GetEndDate() != "2026-05-10" || got.GetEndReason() != "done" || got.GetIsActive() {
		t.Fatalf("unexpected mapped constraint: %#v", got)
	}
	if got.GetMetricId() != 5 || got.GetMetricOperator() != "gt" || got.GetMetricTargetValue() != 100 {
		t.Fatalf("unexpected metric fields: %#v", got)
	}
}

func TestConstraintToProtoNil(t *testing.T) {
	if ConstraintToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestConstraintsToProto(t *testing.T) {
	constraints := []domain.Constraint{
		{ID: 1, Description: "A"},
		{ID: 2, Description: "B"},
	}
	got := ConstraintsToProto(constraints)
	if len(got) != 2 || got[0].GetDescription() != "A" || got[1].GetDescription() != "B" {
		t.Fatalf("unexpected mapped constraints: %#v", got)
	}
}

func TestConstraintFromCreateRequest(t *testing.T) {
	endDate := "2026-05-10"
	metricID := int32(5)
	op := "gt"
	target := 100.0
	req := &timelogv1.CreateConstraintRequest{
		Description:       "Test",
		PunishmentQuote:   "Quote",
		StartDate:         "2026-05-01",
		EndDate:           &endDate,
		MetricId:          &metricID,
		MetricOperator:    &op,
		MetricTargetValue: &target,
	}
	got, err := ConstraintFromCreateRequest(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Description != "Test" || got.PunishmentQuote != "Quote" || got.MetricID == nil || *got.MetricID != 5 {
		t.Fatalf("unexpected domain constraint: %#v", got)
	}
}

func TestConstraintFromCreateRequestNil(t *testing.T) {
	got, err := ConstraintFromCreateRequest(nil)
	if err != nil || got != nil {
		t.Fatalf("expected nil, got (%v, %v)", got, err)
	}
}

func TestApplyConstraintUpdate(t *testing.T) {
	c := &domain.Constraint{Description: "Old", PunishmentQuote: "OldQuote"}
	newDesc := "New"
	newQuote := "NewQuote"
	newStart := "2026-06-01"
	newEnd := "2026-06-10"
	metricID := int32(3)
	op := "lt"
	target := 50.0
	req := &timelogv1.UpdateConstraintRequest{
		Description:       &newDesc,
		PunishmentQuote:   &newQuote,
		StartDate:         &newStart,
		EndDate:           &newEnd,
		MetricId:          &metricID,
		MetricOperator:    &op,
		MetricTargetValue: &target,
	}
	if err := ApplyConstraintUpdate(c, req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Description != "New" || c.PunishmentQuote != "NewQuote" || c.MetricID == nil || *c.MetricID != 3 {
		t.Fatalf("unexpected updated constraint: %#v", c)
	}
}

func TestApplyConstraintUpdateNil(t *testing.T) {
	if err := ApplyConstraintUpdate(nil, nil); err != nil {
		t.Fatalf("expected no error for nil inputs, got %v", err)
	}
}

func TestConstraintEvaluationToProto(t *testing.T) {
	eval := &service.ConstraintEvaluation{
		ConstraintID: 1,
		Passed:       true,
		Actual:       120,
		Target:       100,
		Operator:     "gt",
	}
	got := ConstraintEvaluationToProto(eval)
	if got.GetConstraintId() != 1 || !got.GetPassed() || got.GetActual() != 120 || got.GetTarget() != 100 || got.GetOperator() != "gt" {
		t.Fatalf("unexpected mapped evaluation: %#v", got)
	}
}

func TestConstraintEvaluationToProtoNil(t *testing.T) {
	if ConstraintEvaluationToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}
