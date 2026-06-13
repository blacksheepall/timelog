package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

type RecordMetricInput struct {
	MetricName string
	Value      float64
	Source     string
	RecordedAt string // RFC3339 or empty
}

type IncrementMetricInput struct {
	MetricName string
	Delta      float64
	Source     string
	RecordedAt string
}

func (s *Service) CreateMetric(metric *gen.Metric) error {
	return s.metricRepo.CreateMetric(metric)
}

func (s *Service) GetMetricByID(id int32) (*gen.Metric, error) {
	return s.metricRepo.GetMetricByID(id)
}

func (s *Service) GetMetricByName(name string) (*gen.Metric, error) {
	return s.metricRepo.GetMetricByName(name)
}

func (s *Service) ListMetrics() ([]gen.Metric, error) {
	return s.metricRepo.ListMetrics()
}

func (s *Service) UpdateMetric(metric *gen.Metric) error {
	return s.metricRepo.UpdateMetric(metric)
}

func (s *Service) DeleteMetric(id int32) error {
	return s.metricRepo.DeleteMetric(id)
}

func (s *Service) ListMetricRecords(metricID int32) ([]gen.MetricRecord, error) {
	return s.metricRepo.ListMetricRecordsByMetricID(metricID)
}

func (s *Service) RecordMetric(input RecordMetricInput) (*gen.Metric, error) {
	metric, err := s.metricRepo.GetMetricByName(input.MetricName)
	if err != nil {
		return nil, fmt.Errorf("metric not found: %w", err)
	}

	recordedAt, err := parseMetricRecordedAt(input.RecordedAt)
	if err != nil {
		return nil, err
	}

	record := &gen.MetricRecord{
		MetricID:   *metric.ID,
		Value:      input.Value,
		Source:     &input.Source,
		RecordedAt: &recordedAt,
	}
	if err := s.metricRepo.CreateMetricRecord(record); err != nil {
		return nil, err
	}

	if err := s.metricRepo.UpdateMetricCurrentValue(*metric.ID, input.Value, recordedAt); err != nil {
		return nil, err
	}

	return s.metricRepo.GetMetricByID(*metric.ID)
}

func (s *Service) IncrementMetric(input IncrementMetricInput) (*gen.Metric, error) {
	metric, err := s.metricRepo.GetMetricByName(input.MetricName)
	if err != nil {
		return nil, fmt.Errorf("metric not found: %w", err)
	}

	recordedAt, err := parseMetricRecordedAt(input.RecordedAt)
	if err != nil {
		return nil, err
	}

	current := 0.0
	if metric.CurrentValue != nil {
		current = *metric.CurrentValue
	}
	newValue := current + input.Delta

	record := &gen.MetricRecord{
		MetricID:   *metric.ID,
		Value:      input.Delta,
		Source:     &input.Source,
		RecordedAt: &recordedAt,
	}
	if err := s.metricRepo.CreateMetricRecord(record); err != nil {
		return nil, err
	}

	if err := s.metricRepo.UpdateMetricCurrentValue(*metric.ID, newValue, recordedAt); err != nil {
		return nil, err
	}

	return s.metricRepo.GetMetricByID(*metric.ID)
}

func parseMetricRecordedAt(value string) (time.Time, error) {
	if value == "" {
		return time.Now().UTC(), nil
	}
	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t.UTC(), nil
	}
	if t, err := time.Parse("2006-01-02 15:04:05", value); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("invalid recorded_at %q", value)
}
