package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// constraintRepo implements ports.ConstraintRepository using the model layer.
type constraintRepo struct {
	dao *model.Dao
}

var _ ports.ConstraintRepository = (*constraintRepo)(nil)

func newConstraintRepo(dao *model.Dao) *constraintRepo {
	return &constraintRepo{dao: dao}
}

func (r *constraintRepo) CreateConstraint(constraint *gen.Constraint) error {
	return model.CreateConstraint(r.dao.Db(), constraint)
}

func (r *constraintRepo) GetConstraintByID(id int32) (*gen.Constraint, error) {
	return model.GetConstraintByID(r.dao.Db(), id)
}

func (r *constraintRepo) GetAllConstraints() ([]gen.Constraint, error) {
	return model.GetAllConstraints(r.dao.Db())
}

func (r *constraintRepo) GetActiveConstraints() ([]gen.Constraint, error) {
	return model.GetActiveConstraints(r.dao.Db())
}

func (r *constraintRepo) GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error) {
	return model.GetConstraintsByDateRange(r.dao.Db(), startDate, endDate)
}

func (r *constraintRepo) UpdateConstraint(constraint *gen.Constraint) error {
	return model.UpdateConstraint(r.dao.Db(), constraint)
}

func (r *constraintRepo) DeleteConstraint(id int32) error {
	return model.DeleteConstraint(r.dao.Db(), id)
}

func (r *constraintRepo) MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	return model.MarkConstraintAsCompleted(r.dao.Db(), constraintID, endReason)
}

func (r *constraintRepo) MarkConstraintAsActive(constraintID int32) error {
	return model.MarkConstraintAsActive(r.dao.Db(), constraintID)
}
