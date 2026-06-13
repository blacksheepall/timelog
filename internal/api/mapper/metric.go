package mapper

import (
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/service"
)

func MetricToProto(m *gen.Metric) *timelogv1.Metric {
	if m == nil {
		return nil
	}
	return &timelogv1.Metric{
		Id:             Int32Value(m.ID),
		Name:           m.Name,
		Description:    StringValue(m.Description),
		MetricType:     m.MetricType,
		Unit:           m.Unit,
		CurrentValue:   m.CurrentValue,
		LastRecordedAt: optionalString(FormatTimeUTCPtr(m.LastRecordedAt)),
		Extra:          m.Extra,
		CreatedAt:      FormatTimeUTCPtr(m.CreatedAt),
		UpdatedAt:      FormatTimeUTCPtr(m.UpdatedAt),
	}
}

func MetricsToProto(metrics []gen.Metric) []*timelogv1.Metric {
	out := make([]*timelogv1.Metric, 0, len(metrics))
	for i := range metrics {
		out = append(out, MetricToProto(&metrics[i]))
	}
	return out
}

func MetricRecordToProto(r *gen.MetricRecord) *timelogv1.MetricRecord {
	if r == nil {
		return nil
	}
	return &timelogv1.MetricRecord{
		Id:         Int32Value(r.ID),
		MetricId:   r.MetricID,
		Value:      r.Value,
		Source:     r.Source,
		RecordedAt: optionalString(FormatTimeUTCPtr(r.RecordedAt)),
		CreatedAt:  FormatTimeUTCPtr(r.CreatedAt),
	}
}

func MetricRecordsToProto(records []gen.MetricRecord) []*timelogv1.MetricRecord {
	out := make([]*timelogv1.MetricRecord, 0, len(records))
	for i := range records {
		out = append(out, MetricRecordToProto(&records[i]))
	}
	return out
}

func MetricFromCreateRequest(req *timelogv1.CreateMetricRequest) *gen.Metric {
	if req == nil {
		return nil
	}
	return &gen.Metric{
		Name:        req.Name,
		Description: &req.Description,
		MetricType:  req.MetricType,
		Unit:        req.Unit,
		Extra:       req.Extra,
	}
}

func ApplyMetricUpdate(m *gen.Metric, req *timelogv1.UpdateMetricRequest) {
	if m == nil || req == nil {
		return
	}
	if req.Name != nil {
		m.Name = req.GetName()
	}
	if req.Description != nil {
		m.Description = req.Description
	}
	if req.MetricType != nil {
		m.MetricType = req.GetMetricType()
	}
	if req.Unit != nil {
		m.Unit = req.GetUnit()
	}
	if req.Extra != nil {
		m.Extra = req.Extra
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
