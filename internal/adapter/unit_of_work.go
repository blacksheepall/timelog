package adapter

import (
	"context"
	"fmt"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// sqliteUnitOfWork implements ports.UnitOfWork using GORM transactions.
type sqliteUnitOfWork struct {
	dao *model.Dao
}

var _ ports.UnitOfWork = (*sqliteUnitOfWork)(nil)

// newUnitOfWork creates a UnitOfWork backed by the given DAO.
func newUnitOfWork(dao *model.Dao) *sqliteUnitOfWork {
	return &sqliteUnitOfWork{dao: dao}
}

// Run executes fn inside a database transaction. If fn returns an error, the
// transaction is rolled back; otherwise it is committed.
func (u *sqliteUnitOfWork) Run(ctx context.Context, fn func(repos ports.UnitOfWorkRepositories) error) error {
	tx := u.dao.Begin()
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		}
	}()

	repos := ports.UnitOfWorkRepositories{
		TimelogRepo:      newTimelogRepo(tx),
		CategoryRepo:     newCategoryRepo(tx),
		TaskRepo:         newTaskRepo(tx),
		ConstraintRepo:   newConstraintRepo(tx),
		PasskeyRepo:      newPasskeyCredentialRepo(tx),
		TempPasswordRepo: newTempPasswordRepo(tx),
	}

	if err := fn(repos); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction rollback failed: %w (original: %w)", rbErr, err)
		}
		return err
	}

	return tx.Commit()
}
