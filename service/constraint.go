package service

import (
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// CreateConstraint 创建约束
func CreateConstraint(constraint *gen.Constraint) error {
	dao := getDao()
	return model.CreateConstraint(dao.Db(), constraint)
}

// GetConstraintByID 根据ID获取约束
func GetConstraintByID(id int32) (*gen.Constraint, error) {
	dao := getDao()
	return model.GetConstraintByID(dao.Db(), id)
}

// GetAllConstraints 获取所有约束
func GetAllConstraints() ([]gen.Constraint, error) {
	dao := getDao()
	return model.GetAllConstraints(dao.Db())
}

// GetActiveConstraints 获取活跃的约束
func GetActiveConstraints() ([]gen.Constraint, error) {
	dao := getDao()
	return model.GetActiveConstraints(dao.Db())
}

// GetConstraintsByDateRange 根据日期范围获取约束
func GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error) {
	dao := getDao()
	return model.GetConstraintsByDateRange(dao.Db(), startDate, endDate)
}

// UpdateConstraint 更新约束
func UpdateConstraint(constraint *gen.Constraint) error {
	dao := getDao()
	return model.UpdateConstraint(dao.Db(), constraint)
}

// DeleteConstraint 删除约束
func DeleteConstraint(id int32) error {
	dao := getDao()
	return model.DeleteConstraint(dao.Db(), id)
}

// MarkConstraintAsCompleted 标记约束为完成
func MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	dao := getDao()
	return model.MarkConstraintAsCompleted(dao.Db(), constraintID, endReason)
}

// MarkConstraintAsActive 重新激活约束
func MarkConstraintAsActive(constraintID int32) error {
	dao := getDao()
	return model.MarkConstraintAsActive(dao.Db(), constraintID)
}
