package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// timelogRepo implements ports.TimelogRepository using the model layer.
type timelogRepo struct {
	db model.DBProvider
}

var _ ports.TimelogRepository = (*timelogRepo)(nil)

func newTimelogRepo(db model.DBProvider) *timelogRepo {
	return &timelogRepo{db: db}
}

func (r *timelogRepo) CreateTimeLog(tl *gen.Timelog) error {
	return model.CreateTimeLog(r.db.Db(), tl)
}

func (r *timelogRepo) GetTimeLogByID(id int32) (*gen.Timelog, error) {
	return model.GetTimeLogByID(r.db.Db(), id)
}

func (r *timelogRepo) ListTimeLogs(conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogs(r.db.Db(), conds...)
}

func (r *timelogRepo) ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogsWithOptions(r.db.Db(), limit, orderBy, conds...)
}

func (r *timelogRepo) ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]gen.Timelog, error) {
	return model.ListTimeLogsByLocalDateRange(r.db.Db(), startDateStr, endDateStr)
}

func (r *timelogRepo) UpdateTimeLog(tl *gen.Timelog) error {
	return model.UpdateTimeLog(r.db.Db(), tl)
}

func (r *timelogRepo) DeleteTimeLog(id int32) error {
	return model.DeleteTimeLog(r.db.Db(), id)
}
