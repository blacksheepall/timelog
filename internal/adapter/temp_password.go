package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// tempPasswordRepo implements ports.TempPasswordRepository using the model layer.
type tempPasswordRepo struct {
	dao *model.Dao
}

var _ ports.TempPasswordRepository = (*tempPasswordRepo)(nil)

func newTempPasswordRepo(dao *model.Dao) *tempPasswordRepo {
	return &tempPasswordRepo{dao: dao}
}

func (r *tempPasswordRepo) CreateTempPassword(tempPassword *model.TempPassword) error {
	return model.CreateTempPassword(r.dao.Db(), tempPassword)
}

func (r *tempPasswordRepo) ListTempPasswords() ([]model.TempPassword, error) {
	return model.ListTempPasswords(r.dao.Db())
}

func (r *tempPasswordRepo) DeleteTempPassword(id uint) error {
	return model.DeleteTempPassword(r.dao.Db(), id)
}

func (r *tempPasswordRepo) DeleteExpiredTempPasswords(now time.Time) error {
	return model.DeleteExpiredTempPasswords(r.dao.Db(), now)
}

func (r *tempPasswordRepo) GetTempPasswordByHash(hash string, now time.Time) (*model.TempPassword, error) {
	return model.GetTempPasswordByHash(r.dao.Db(), hash, now)
}
