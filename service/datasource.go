package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"gorm.io/gorm"
)

var (
	// ErrDatasourceNotFound is returned when the named datasource is not registered.
	ErrDatasourceNotFound = errors.New("datasource not found")
	// ErrSyncAllWritesFailed is returned when every metric write failed during sync.
	ErrSyncAllWritesFailed = errors.New("all metric writes failed")
)

// SyncMetricsResult reports how many data points were written for a datasource.
type SyncMetricsResult struct {
	Source string `json:"source"`
	Synced int    `json:"synced"`
	Failed int    `json:"failed"`
}

// SyncMetrics pulls data from the named external datasource and records each
// point through RecordMetric. Metrics are created on demand if they do not exist.
func (s *Service) SyncMetrics(ctx context.Context, sourceName string) (SyncMetricsResult, error) {
	result := SyncMetricsResult{Source: sourceName}

	if s.dataSourceRegistry == nil {
		return result, fmt.Errorf("%w: %q (registry not configured)", ErrDatasourceNotFound, sourceName)
	}

	source, err := s.dataSourceRegistry.Get(sourceName)
	if err != nil {
		return result, fmt.Errorf("%w: %s", ErrDatasourceNotFound, err)
	}

	points, err := source.Fetch(ctx)
	if err != nil {
		return result, fmt.Errorf("fetch from %q: %w", sourceName, err)
	}

	for _, point := range points {
		// Ensure the metric exists before recording; create on demand.
		if _, err := s.metricRepo.GetMetricByName(point.MetricName); err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				s.log.Error("failed to lookup metric", "name", point.MetricName, "error", err)
				result.Failed++
				continue
			}
			unit := point.Unit
			if unit == "" {
				unit = "count"
			}
			newMetric := &domain.Metric{
				Name:       point.MetricName,
				MetricType: "numeric",
				Unit:       unit,
			}
			if err := s.metricRepo.CreateMetric(newMetric); err != nil {
				s.log.Error("failed to create metric", "name", point.MetricName, "error", err)
				result.Failed++
				continue
			}
		}

		point.Source = sourceName
		if _, err := s.RecordMetric(RecordMetricInput{
			MetricName: point.MetricName,
			Value:      point.Value,
			Source:     point.Source,
			RecordedAt: point.RecordedAt.Format(time.RFC3339),
		}); err != nil {
			s.log.Error("failed to record metric", "name", point.MetricName, "error", err)
			result.Failed++
			continue
		}
		result.Synced++
	}

	if len(points) > 0 && result.Synced == 0 {
		return result, fmt.Errorf("%w: %d failed", ErrSyncAllWritesFailed, result.Failed)
	}

	return result, nil
}
