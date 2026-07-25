package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
)

type fakeDataSource struct {
	name   string
	points []domain.MetricDataPoint
	err    error
}

func (f *fakeDataSource) Name() string { return f.name }
func (f *fakeDataSource) Fetch(ctx context.Context) ([]domain.MetricDataPoint, error) {
	return f.points, f.err
}

type fakeRegistry struct {
	sources map[string]ports.DataSource
}

func (r *fakeRegistry) Get(name string) (ports.DataSource, error) {
	s, ok := r.sources[name]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return s, nil
}

func TestSyncMetrics(t *testing.T) {
	svc, _ := newTestService(t)
	reg := &fakeRegistry{
		sources: map[string]ports.DataSource{
			"maimemo": &fakeDataSource{
				name: "maimemo",
				points: []domain.MetricDataPoint{
					{MetricName: "今日已背单词", Value: 120, RecordedAt: time.Now().UTC(), Source: "maimemo"},
				},
			},
		},
	}
	svc.dataSourceRegistry = reg

	result, err := svc.SyncMetrics(context.Background(), "maimemo")
	if err != nil {
		t.Fatalf("SyncMetrics: %v", err)
	}
	if result.Synced != 1 {
		t.Fatalf("expected synced 1, got %d", result.Synced)
	}
	if result.Failed != 0 {
		t.Fatalf("expected failed 0, got %d", result.Failed)
	}

	metric, err := svc.GetMetricByName("今日已背单词")
	if err != nil {
		t.Fatalf("GetMetricByName: %v", err)
	}
	if metric.CurrentValue == nil || *metric.CurrentValue != 120 {
		t.Fatalf("expected current value 120, got %v", metric.CurrentValue)
	}
}

func TestSyncMetrics_SourceNotFound(t *testing.T) {
	svc, _ := newTestService(t)
	svc.dataSourceRegistry = &fakeRegistry{sources: map[string]ports.DataSource{}}

	_, err := svc.SyncMetrics(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error for missing source")
	}
}
