package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

type metricRepo struct {
	db model.DBProvider
}

var _ ports.MetricRepository = (*metricRepo)(nil)

func newMetricRepo(db model.DBProvider) *metricRepo {
	return &metricRepo{db: db}
}

func (r *metricRepo) CreateMetric(metric *domain.Metric) error {
	g := toGenMetric(metric)
	if err := model.CreateMetric(r.db.Db(), g); err != nil {
		return err
	}
	*metric = *toDomainMetric(g)
	return nil
}

func (r *metricRepo) GetMetricByID(id int32) (*domain.Metric, error) {
	g, err := model.GetMetricByID(r.db.Db(), id)
	if err != nil {
		return nil, err
	}
	return toDomainMetric(g), nil
}

func (r *metricRepo) GetMetricByName(name string) (*domain.Metric, error) {
	g, err := model.GetMetricByName(r.db.Db(), name)
	if err != nil {
		return nil, err
	}
	return toDomainMetric(g), nil
}

func (r *metricRepo) ListMetrics() ([]domain.Metric, error) {
	list, err := model.ListMetrics(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainMetrics(list), nil
}

func (r *metricRepo) UpdateMetric(metric *domain.Metric) error {
	return model.UpdateMetric(r.db.Db(), toGenMetric(metric))
}

func (r *metricRepo) DeleteMetric(id int32) error {
	return model.DeleteMetric(r.db.Db(), id)
}

func (r *metricRepo) CreateMetricRecord(record *domain.MetricRecord) error {
	g := toGenMetricRecord(record)
	if err := model.CreateMetricRecord(r.db.Db(), g); err != nil {
		return err
	}
	*record = *toDomainMetricRecord(g)
	return nil
}

func (r *metricRepo) ListMetricRecordsByMetricID(metricID int32) ([]domain.MetricRecord, error) {
	list, err := model.ListMetricRecordsByMetricID(r.db.Db(), metricID)
	if err != nil {
		return nil, err
	}
	return toDomainMetricRecords(list), nil
}

func (r *metricRepo) UpdateMetricCurrentValue(metricID int32, value float64, recordedAt time.Time) error {
	return model.UpdateMetricCurrentValue(r.db.Db(), metricID, value, recordedAt)
}
