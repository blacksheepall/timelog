package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

type tempPasswordRepo struct {
	db model.DBProvider
}

var _ ports.TempPasswordRepository = (*tempPasswordRepo)(nil)

func newTempPasswordRepo(db model.DBProvider) *tempPasswordRepo {
	return &tempPasswordRepo{db: db}
}

func (r *tempPasswordRepo) CreateTempPassword(tempPassword *domain.TempPassword) error {
	return model.CreateTempPassword(r.db.Db(), toModelTempPassword(tempPassword))
}

func (r *tempPasswordRepo) ListTempPasswords() ([]domain.TempPassword, error) {
	list, err := model.ListTempPasswords(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainTempPasswords(list), nil
}

func (r *tempPasswordRepo) DeleteTempPassword(id int32) error {
	return model.DeleteTempPassword(r.db.Db(), uint(id))
}

func (r *tempPasswordRepo) DeleteExpiredTempPasswords(now time.Time) error {
	return model.DeleteExpiredTempPasswords(r.db.Db(), now)
}

func (r *tempPasswordRepo) GetTempPasswordByHash(hash string, now time.Time) (*domain.TempPassword, error) {
	m, err := model.GetTempPasswordByHash(r.db.Db(), hash, now)
	if err != nil {
		return nil, err
	}
	return toDomainTempPassword(m), nil
}
