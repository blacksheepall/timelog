package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// constraintRepo implements ports.ConstraintRepository using the model layer.
type constraintRepo struct {
	db model.DBProvider
}

var _ ports.ConstraintRepository = (*constraintRepo)(nil)

func newConstraintRepo(db model.DBProvider) *constraintRepo {
	return &constraintRepo{db: db}
}

func (r *constraintRepo) CreateConstraint(constraint *domain.Constraint) error {
	g := toGenConstraint(constraint)
	if err := model.CreateConstraint(r.db.Db(), g); err != nil {
		return err
	}
	*constraint = *toDomainConstraint(g)
	return nil
}

func (r *constraintRepo) GetConstraintByID(id int32) (*domain.Constraint, error) {
	g, err := model.GetConstraintByID(r.db.Db(), id)
	if err != nil {
		return nil, err
	}
	return toDomainConstraint(g), nil
}

func (r *constraintRepo) GetAllConstraints() ([]domain.Constraint, error) {
	list, err := model.GetAllConstraints(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainConstraints(list), nil
}

func (r *constraintRepo) GetActiveConstraints() ([]domain.Constraint, error) {
	list, err := model.GetActiveConstraints(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainConstraints(list), nil
}

func (r *constraintRepo) GetConstraintsByDateRange(startDate, endDate time.Time) ([]domain.Constraint, error) {
	list, err := model.GetConstraintsByDateRange(r.db.Db(), startDate, endDate)
	if err != nil {
		return nil, err
	}
	return toDomainConstraints(list), nil
}

func (r *constraintRepo) UpdateConstraint(constraint *domain.Constraint) error {
	return model.UpdateConstraint(r.db.Db(), toGenConstraint(constraint))
}

func (r *constraintRepo) DeleteConstraint(id int32) error {
	return model.DeleteConstraint(r.db.Db(), id)
}

func (r *constraintRepo) MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	return model.MarkConstraintAsCompleted(r.db.Db(), constraintID, endReason)
}

func (r *constraintRepo) MarkConstraintAsActive(constraintID int32) error {
	return model.MarkConstraintAsActive(r.db.Db(), constraintID)
}
