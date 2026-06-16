package mapper

import (
	"testing"
	"time"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func TestMetricToProto(t *testing.T) {
	current := 42.0
	lastRecorded := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	m := &domain.Metric{
		ID:             1,
		Name:           "Pages",
		Description:    "Pages read",
		MetricType:     "counter",
		Unit:           "page",
		CurrentValue:   &current,
		LastRecordedAt: &lastRecorded,
		Extra:          "extra",
		CreatedAt:      time.Date(2026, 5, 30, 10, 0, 0, 0, time.UTC),
		UpdatedAt:      time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
	}
	got := MetricToProto(m)
	if got.GetName() != "Pages" || got.GetDescription() != "Pages read" || got.GetCurrentValue() != 42 || got.GetExtra() != "extra" {
		t.Fatalf("unexpected mapped metric: %#v", got)
	}
}

func TestMetricToProtoNil(t *testing.T) {
	if MetricToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestMetricsToProto(t *testing.T) {
	metrics := []domain.Metric{
		{ID: 1, Name: "A"},
		{ID: 2, Name: "B"},
	}
	got := MetricsToProto(metrics)
	if len(got) != 2 || got[0].GetName() != "A" || got[1].GetName() != "B" {
		t.Fatalf("unexpected mapped metrics: %#v", got)
	}
}

func TestMetricRecordToProto(t *testing.T) {
	recordedAt := time.Date(2026, 5, 30, 12, 0, 0, 0, time.UTC)
	r := &domain.MetricRecord{
		ID:         1,
		MetricID:   2,
		Value:      10,
		Source:     "manual",
		RecordedAt: &recordedAt,
		CreatedAt:  time.Date(2026, 5, 30, 11, 0, 0, 0, time.UTC),
	}
	got := MetricRecordToProto(r)
	if got.GetMetricId() != 2 || got.GetValue() != 10 || got.GetSource() != "manual" {
		t.Fatalf("unexpected mapped record: %#v", got)
	}
}

func TestMetricRecordToProtoNil(t *testing.T) {
	if MetricRecordToProto(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestMetricRecordsToProto(t *testing.T) {
	records := []domain.MetricRecord{
		{ID: 1, Value: 10},
		{ID: 2, Value: 20},
	}
	got := MetricRecordsToProto(records)
	if len(got) != 2 || got[0].GetValue() != 10 || got[1].GetValue() != 20 {
		t.Fatalf("unexpected mapped records: %#v", got)
	}
}

func TestMetricFromCreateRequest(t *testing.T) {
	extra := "extra"
	req := &timelogv1.CreateMetricRequest{
		Name:        "Pages",
		Description: "Pages read",
		MetricType:  "counter",
		Unit:        "page",
		Extra:       &extra,
	}
	got := MetricFromCreateRequest(req)
	if got.Name != "Pages" || got.Extra != "extra" {
		t.Fatalf("unexpected domain metric: %#v", got)
	}
}

func TestMetricFromCreateRequestNil(t *testing.T) {
	if MetricFromCreateRequest(nil) != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestApplyMetricUpdate(t *testing.T) {
	m := &domain.Metric{Name: "Old", Description: "OldDesc", MetricType: "old", Unit: "old", Extra: "old"}
	newName := "New"
	newDesc := "NewDesc"
	newType := "counter"
	newUnit := "page"
	newExtra := "new"
	req := &timelogv1.UpdateMetricRequest{
		Name:        &newName,
		Description: &newDesc,
		MetricType:  &newType,
		Unit:        &newUnit,
		Extra:       &newExtra,
	}
	ApplyMetricUpdate(m, req)
	if m.Name != "New" || m.Description != "NewDesc" || m.MetricType != "counter" || m.Unit != "page" || m.Extra != "new" {
		t.Fatalf("unexpected updated metric: %#v", m)
	}
}

func TestApplyMetricUpdateNil(t *testing.T) {
	ApplyMetricUpdate(nil, nil)
}
