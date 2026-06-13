package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// timelogRepo implements ports.TimelogRepository using the model layer.
type timelogRepo struct {
	dao *model.Dao
}

var _ ports.TimelogRepository = (*timelogRepo)(nil)

func newTimelogRepo(dao *model.Dao) *timelogRepo {
	return &timelogRepo{dao: dao}
}

func (r *timelogRepo) CreateTimeLog(tl *gen.Timelog) error {
	return model.CreateTimeLog(r.dao.Db(), tl)
}

func (r *timelogRepo) GetTimeLogByID(id int32) (*gen.Timelog, error) {
	return model.GetTimeLogByID(r.dao.Db(), id)
}

func (r *timelogRepo) ListTimeLogs(conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogs(r.dao.Db(), conds...)
}

func (r *timelogRepo) ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogsWithOptions(r.dao.Db(), limit, orderBy, conds...)
}

func (r *timelogRepo) ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]gen.Timelog, error) {
	return model.ListTimeLogsByLocalDateRange(r.dao.Db(), startDateStr, endDateStr)
}

func (r *timelogRepo) UpdateTimeLog(tl *gen.Timelog) error {
	return model.UpdateTimeLog(r.dao.Db(), tl)
}

func (r *timelogRepo) DeleteTimeLog(id int32) error {
	return model.DeleteTimeLog(r.dao.Db(), id)
}
