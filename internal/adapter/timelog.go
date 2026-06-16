package adapter

import (
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// timelogRepo implements ports.TimelogRepository using the model layer.
type timelogRepo struct {
	db model.DBProvider
}

var _ ports.TimelogRepository = (*timelogRepo)(nil)

func newTimelogRepo(db model.DBProvider) *timelogRepo {
	return &timelogRepo{db: db}
}

func (r *timelogRepo) CreateTimeLog(tl *domain.TimeLog) error {
	g := toGenTimelog(tl)
	if err := model.CreateTimeLog(r.db.Db(), g); err != nil {
		return err
	}
	*tl = *toDomainTimelog(g)
	return nil
}

func (r *timelogRepo) GetTimeLogByID(id int32) (*domain.TimeLog, error) {
	g, err := model.GetTimeLogByID(r.db.Db(), id)
	if err != nil {
		return nil, err
	}
	return toDomainTimelog(g), nil
}

func (r *timelogRepo) ListTimeLogs(conds ...interface{}) ([]domain.TimeLog, error) {
	list, err := model.ListTimeLogs(r.db.Db(), conds...)
	if err != nil {
		return nil, err
	}
	return toDomainTimelogs(list), nil
}

func (r *timelogRepo) ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]domain.TimeLog, error) {
	list, err := model.ListTimeLogsWithOptions(r.db.Db(), limit, orderBy, conds...)
	if err != nil {
		return nil, err
	}
	return toDomainTimelogs(list), nil
}

func (r *timelogRepo) ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]domain.TimeLog, error) {
	list, err := model.ListTimeLogsByLocalDateRange(r.db.Db(), startDateStr, endDateStr)
	if err != nil {
		return nil, err
	}
	return toDomainTimelogs(list), nil
}

func (r *timelogRepo) UpdateTimeLog(tl *domain.TimeLog) error {
	return model.UpdateTimeLog(r.db.Db(), toGenTimelog(tl))
}

func (r *timelogRepo) DeleteTimeLog(id int32) error {
	return model.DeleteTimeLog(r.db.Db(), id)
}
