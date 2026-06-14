package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

type metricRepo struct {
	db model.DBProvider
}

var _ ports.MetricRepository = (*metricRepo)(nil)

func newMetricRepo(db model.DBProvider) *metricRepo {
	return &metricRepo{db: db}
}

func (r *metricRepo) CreateMetric(metric *gen.Metric) error {
	return model.CreateMetric(r.db.Db(), metric)
}

func (r *metricRepo) GetMetricByID(id int32) (*gen.Metric, error) {
	return model.GetMetricByID(r.db.Db(), id)
}

func (r *metricRepo) GetMetricByName(name string) (*gen.Metric, error) {
	return model.GetMetricByName(r.db.Db(), name)
}

func (r *metricRepo) ListMetrics() ([]gen.Metric, error) {
	return model.ListMetrics(r.db.Db())
}

func (r *metricRepo) UpdateMetric(metric *gen.Metric) error {
	return model.UpdateMetric(r.db.Db(), metric)
}

func (r *metricRepo) DeleteMetric(id int32) error {
	return model.DeleteMetric(r.db.Db(), id)
}

func (r *metricRepo) CreateMetricRecord(record *gen.MetricRecord) error {
	return model.CreateMetricRecord(r.db.Db(), record)
}

func (r *metricRepo) ListMetricRecordsByMetricID(metricID int32) ([]gen.MetricRecord, error) {
	return model.ListMetricRecordsByMetricID(r.db.Db(), metricID)
}

func (r *metricRepo) UpdateMetricCurrentValue(metricID int32, value float64, recordedAt time.Time) error {
	return model.UpdateMetricCurrentValue(r.db.Db(), metricID, value, recordedAt)
}
