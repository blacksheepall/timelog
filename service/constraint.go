package service

import (
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

// --- Package-level wrappers (transitional) ---

func CreateConstraint(constraint *gen.Constraint) error { return defaultService.CreateConstraint(constraint) }
func GetConstraintByID(id int32) (*gen.Constraint, error) {
	return defaultService.GetConstraintByID(id)
}
func GetAllConstraints() ([]gen.Constraint, error)          { return defaultService.GetAllConstraints() }
func GetActiveConstraints() ([]gen.Constraint, error)       { return defaultService.GetActiveConstraints() }
func GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error) {
	return defaultService.GetConstraintsByDateRange(startDate, endDate)
}
func UpdateConstraint(constraint *gen.Constraint) error {
	return defaultService.UpdateConstraint(constraint)
}
func DeleteConstraint(id int32) error { return defaultService.DeleteConstraint(id) }
func MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	return defaultService.MarkConstraintAsCompleted(constraintID, endReason)
}
func MarkConstraintAsActive(constraintID int32) error {
	return defaultService.MarkConstraintAsActive(constraintID)
}
