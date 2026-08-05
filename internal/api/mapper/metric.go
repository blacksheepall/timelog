package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/service"
)

func MetricToProto(m *domain.Metric) *timelogv1.Metric {
	if m == nil {
		return nil
	}
	return &timelogv1.Metric{
		Id:             m.ID,
		Name:           m.Name,
		Description:    optionalString(m.Description),
		MetricType:     m.MetricType,
		Unit:           m.Unit,
		CurrentValue:   m.CurrentValue,
		LastRecordedAt: optionalString(FormatTimeUTCPtr(m.LastRecordedAt)),
		Extra:          optionalString(m.Extra),
		CreatedAt:      FormatTimeUTC(m.CreatedAt),
		UpdatedAt:      FormatTimeUTC(m.UpdatedAt),
	}
}

func MetricsToProto(metrics []domain.Metric) []*timelogv1.Metric {
	out := make([]*timelogv1.Metric, 0, len(metrics))
	for i := range metrics {
		out = append(out, MetricToProto(&metrics[i]))
	}
	return out
}

func MetricRecordToProto(r *domain.MetricRecord) *timelogv1.MetricRecord {
	if r == nil {
		return nil
	}
	return &timelogv1.MetricRecord{
		Id:         r.ID,
		MetricId:   r.MetricID,
		Value:      r.Value,
		Source:     &r.Source,
		RecordedAt: optionalString(FormatTimeUTCPtr(r.RecordedAt)),
		CreatedAt:  FormatTimeUTC(r.CreatedAt),
	}
}

func MetricRecordsToProto(records []domain.MetricRecord) []*timelogv1.MetricRecord {
	out := make([]*timelogv1.MetricRecord, 0, len(records))
	for i := range records {
		out = append(out, MetricRecordToProto(&records[i]))
	}
	return out
}

func MetricFromCreateRequest(req *timelogv1.CreateMetricRequest) *domain.Metric {
	if req == nil {
		return nil
	}
	return &domain.Metric{
		Name:        req.Name,
		Description: req.GetDescription(),
		MetricType:  req.MetricType,
		Unit:        req.Unit,
		Extra:       req.GetExtra(),
	}
}

func ApplyMetricUpdate(m *domain.Metric, req *timelogv1.UpdateMetricRequest) {
	if m == nil || req == nil {
		return
	}
	if req.Name != nil {
		m.Name = req.GetName()
	}
	if req.Description != nil {
		m.Description = req.GetDescription()
	}
	if req.MetricType != nil {
		m.MetricType = req.GetMetricType()
	}
	if req.Unit != nil {
		m.Unit = req.GetUnit()
	}
	if req.Extra != nil {
		m.Extra = req.GetExtra()
	}
}

func ConstraintEvaluationToProto(eval *service.ConstraintEvaluation) *timelogv1.ConstraintEvaluation {
	if eval == nil {
		return nil
	}
	return &timelogv1.ConstraintEvaluation{
		ConstraintId: eval.ConstraintID,
		Passed:       eval.Passed,
		Actual:       eval.Actual,
		Target:       eval.Target,
		Operator:     eval.Operator,
	}
}
