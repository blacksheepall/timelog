package service

import (
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
)

// CreateConstraint 创建约束
func (s *Service) CreateConstraint(constraint *gen.Constraint) error {
	return s.constraintRepo.CreateConstraint(constraint)
}

// GetConstraintByID 根据ID获取约束
func (s *Service) GetConstraintByID(id int32) (*gen.Constraint, error) {
	return s.constraintRepo.GetConstraintByID(id)
}

// GetAllConstraints 获取所有约束
func (s *Service) GetAllConstraints() ([]gen.Constraint, error) {
	return s.constraintRepo.GetAllConstraints()
}

// GetActiveConstraints 获取活跃的约束
func (s *Service) GetActiveConstraints() ([]gen.Constraint, error) {
	return s.constraintRepo.GetActiveConstraints()
}

// GetConstraintsByDateRange 根据日期范围获取约束
func (s *Service) GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error) {
	return s.constraintRepo.GetConstraintsByDateRange(startDate, endDate)
}

// UpdateConstraint 更新约束
func (s *Service) UpdateConstraint(constraint *gen.Constraint) error {
	return s.constraintRepo.UpdateConstraint(constraint)
}

// DeleteConstraint 删除约束
func (s *Service) DeleteConstraint(id int32) error {
	return s.constraintRepo.DeleteConstraint(id)
}

// MarkConstraintAsCompleted 标记约束为完成
func (s *Service) MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	return s.constraintRepo.MarkConstraintAsCompleted(constraintID, endReason)
}

// MarkConstraintAsActive 重新激活约束
func (s *Service) MarkConstraintAsActive(constraintID int32) error {
	return s.constraintRepo.MarkConstraintAsActive(constraintID)
}

type ConstraintEvaluation struct {
	ConstraintID int32
	Passed       bool
	Actual       float64
	Target       float64
	Operator     string
}

func (s *Service) EvaluateConstraint(constraintID int32) (*ConstraintEvaluation, error) {
	c, err := s.constraintRepo.GetConstraintByID(constraintID)
	if err != nil {
		return nil, err
	}
	if c.MetricID == nil || c.MetricOperator == nil || c.MetricTargetValue == nil {
		return nil, fmt.Errorf("constraint has no metric rule")
	}

	metric, err := s.metricRepo.GetMetricByID(*c.MetricID)
	if err != nil {
		return nil, err
	}
	if metric.CurrentValue == nil {
		return nil, fmt.Errorf("metric has no current value")
	}

	passed, err := evaluateMetric(*metric.CurrentValue, *c.MetricTargetValue, *c.MetricOperator)
	if err != nil {
		return nil, err
	}

	return &ConstraintEvaluation{
		ConstraintID: constraintID,
		Passed:       passed,
		Actual:       *metric.CurrentValue,
		Target:       *c.MetricTargetValue,
		Operator:     *c.MetricOperator,
	}, nil
}

func evaluateMetric(actual, target float64, op string) (bool, error) {
	switch op {
	case "eq":
		return actual == target, nil
	case "ne":
		return actual != target, nil
	case "gt":
		return actual > target, nil
	case "gte":
		return actual >= target, nil
	case "lt":
		return actual < target, nil
	case "lte":
		return actual <= target, nil
	default:
		return false, fmt.Errorf("unsupported metric operator %q", op)
	}
}
