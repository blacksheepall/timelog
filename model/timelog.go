package model

import (
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/errs"
	"gorm.io/gorm"
)

// --- CRUD ---

// CreateTimeLog 新增一条时间日志
// 如果当前已存在 ongoing（end_time 为 NULL）的 timelog，则禁止写入
func CreateTimeLog(db *gorm.DB, tl *gen.Timelog) error {
	var count int64
	if err := db.Model(&gen.Timelog{}).Where("end_time IS NULL").Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errs.ErrOngoingTimeLogExists
	}
	return db.Create(tl).Error
}

// GetTimeLogByID 根据ID获取时间日志
func GetTimeLogByID(db *gorm.DB, id int32) (*gen.Timelog, error) {
	var tl gen.Timelog
	err := db.First(&tl, id).Error
	return &tl, err
}

// ListTimeLogs 查询时间日志（可扩展条件）
func ListTimeLogs(db *gorm.DB, conds ...interface{}) ([]gen.Timelog, error) {
	var tls []gen.Timelog
	err := db.Find(&tls, conds...).Error
	return tls, err
}

// ListTimeLogsWithOptions 查询时间日志（支持排序和限制）
func ListTimeLogsWithOptions(db *gorm.DB, limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error) {
	var tls []gen.Timelog
	query := db

	if len(conds) > 0 {
		query = query.Where(conds[0], conds[1:]...)
	}

	if orderBy != "" {
		query = query.Order(orderBy)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&tls).Error
	return tls, err
}

// UpdateTimeLog 更新时间日志
func UpdateTimeLog(db *gorm.DB, tl *gen.Timelog) error {
	return db.Save(tl).Error
}

// DeleteTimeLog 删除时间日志
func DeleteTimeLog(db *gorm.DB, id int32) error {
	return db.Delete(&gen.Timelog{}, id).Error
}

// 定义新加坡时区
var singaporeLocation *time.Location

func init() {
	var err error
	singaporeLocation, err = time.LoadLocation("Asia/Singapore")
	if err != nil {
		// Fallback to UTC+8 if timezone data is not available
		singaporeLocation = time.FixedZone("SGT", 8*60*60)
	}
}

// GetSingaporeLocation 返回新加坡时区，供其他包使用
func GetSingaporeLocation() *time.Location {
	return singaporeLocation
}

// ListTimeLogsByLocalDateRange 根据本地日期范围查询时间日志
// startDateStr 和 endDateStr 格式为 "YYYY-MM-DD"，会被解析为新加坡时区
// 数据库存储的是 UTC 时间，该函数会自动转换
func ListTimeLogsByLocalDateRange(db *gorm.DB, startDateStr, endDateStr string) ([]gen.Timelog, error) {
	startDate, err := time.ParseInLocation("2006-01-02", startDateStr, singaporeLocation)
	if err != nil {
		return nil, err
	}

	endDate, err := time.ParseInLocation("2006-01-02", endDateStr, singaporeLocation)
	if err != nil {
		return nil, err
	}

	// endDate 设置为当天 23:59:59 SGT
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	// 转换为 UTC 进行数据库查询
	startUTC := startDate.UTC()
	endUTC := endDate.UTC()

	return ListTimeLogsWithOptions(db, 0, "start_time ASC", "start_time >= ? AND start_time <= ?", startUTC, endUTC)
}
