package model

import (
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	"gorm.io/gorm"
)

func CreateMetric(db *gorm.DB, metric *gen.Metric) error {
	return db.Create(metric).Error
}

func GetMetricByID(db *gorm.DB, id int32) (*gen.Metric, error) {
	var metric gen.Metric
	err := db.First(&metric, id).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func GetMetricByName(db *gorm.DB, name string) (*gen.Metric, error) {
	var metric gen.Metric
	err := db.Where("name = ?", name).First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

func ListMetrics(db *gorm.DB) ([]gen.Metric, error) {
	var metrics []gen.Metric
	err := db.Find(&metrics).Error
	return metrics, err
}

func UpdateMetric(db *gorm.DB, metric *gen.Metric) error {
	return db.Save(metric).Error
}

func DeleteMetric(db *gorm.DB, id int32) error {
	return db.Delete(&gen.Metric{}, id).Error
}

func CreateMetricRecord(db *gorm.DB, record *gen.MetricRecord) error {
	return db.Create(record).Error
}

func ListMetricRecordsByMetricID(db *gorm.DB, metricID int32) ([]gen.MetricRecord, error) {
	var records []gen.MetricRecord
	err := db.Where("metric_id = ?", metricID).Order("created_at DESC").Find(&records).Error
	return records, err
}

func UpdateMetricCurrentValue(db *gorm.DB, metricID int32, value float64, recordedAt time.Time) error {
	return db.Model(&gen.Metric{}).Where("id = ?", metricID).Updates(map[string]interface{}{
		"current_value":    value,
		"last_recorded_at": recordedAt,
	}).Error
}
