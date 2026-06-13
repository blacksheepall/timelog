package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// tempPasswordRepo implements ports.TempPasswordRepository using the model layer.
type tempPasswordRepo struct {
	db model.DBProvider
}

var _ ports.TempPasswordRepository = (*tempPasswordRepo)(nil)

func newTempPasswordRepo(db model.DBProvider) *tempPasswordRepo {
	return &tempPasswordRepo{db: db}
}

func (r *tempPasswordRepo) CreateTempPassword(tempPassword *model.TempPassword) error {
	return model.CreateTempPassword(r.db.Db(), tempPassword)
}

func (r *tempPasswordRepo) ListTempPasswords() ([]model.TempPassword, error) {
	return model.ListTempPasswords(r.db.Db())
}

func (r *tempPasswordRepo) DeleteTempPassword(id uint) error {
	return model.DeleteTempPassword(r.db.Db(), id)
}

func (r *tempPasswordRepo) DeleteExpiredTempPasswords(now time.Time) error {
	return model.DeleteExpiredTempPasswords(r.db.Db(), now)
}

func (r *tempPasswordRepo) GetTempPasswordByHash(hash string, now time.Time) (*model.TempPassword, error) {
	return model.GetTempPasswordByHash(r.db.Db(), hash, now)
}
