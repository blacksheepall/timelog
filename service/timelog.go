package service

import (
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// --- TimeLog Service ---

// CreateTimeLog 新增一条时间日志
func CreateTimeLog(tl *gen.Timelog) error {
	db := getDao().Db()
	return model.CreateTimeLog(db, tl)
}

// GetTimeLogByID 根据ID获取时间日志
func GetTimeLogByID(id int32) (*gen.Timelog, error) {
	db := getDao().Db()
	return model.GetTimeLogByID(db, id)
}

// ListTimeLogs 查询时间日志（可扩展条件）
func ListTimeLogs(conds ...interface{}) ([]gen.Timelog, error) {
	db := getDao().Db()
	return model.ListTimeLogs(db, conds...)
}

// ListTimeLogsWithOptions 查询时间日志（支持排序和限制）
func ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error) {
	db := getDao().Db()
	return model.ListTimeLogsWithOptions(db, limit, orderBy, conds...)
}

// ListTimeLogsByLocalDateRange queries timelogs whose start_time falls within
// the inclusive local-date range in Asia/Singapore (YYYY-MM-DD strings).
func ListTimeLogsByLocalDateRange(startDate, endDate string) ([]gen.Timelog, error) {
	db := getDao().Db()
	return model.ListTimeLogsByLocalDateRange(db, startDate, endDate)
}

// UpdateTimeLog 更新一条时间日志
func UpdateTimeLog(tl *gen.Timelog) error {
	db := getDao().Db()
	return model.UpdateTimeLog(db, tl)
}

// DeleteTimeLog 删除一条时间日志
func DeleteTimeLog(id int32) error {
	db := getDao().Db()
	return model.DeleteTimeLog(db, id)
}

// --- Category Service ---

// CreateCategory 创建分类
func CreateCategory(category *gen.Category) error {
	db := getDao().Db()
	return model.CreateCategory(db, category)
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(id int32) (*gen.Category, error) {
	db := getDao().Db()
	return model.GetCategoryByID(db, id)
}

// GetCategoryByName 根据名称获取分类
func GetCategoryByName(name string, parentID *int32) (*gen.Category, error) {
	db := getDao().Db()
	return model.GetCategoryByName(db, name, parentID)
}

// ListCategories 查询所有分类
func ListCategories(conds ...interface{}) ([]gen.Category, error) {
	db := getDao().Db()
	return model.ListCategories(db, conds...)
}

// ListCategoriesByLevel 按层级查询分类
func ListCategoriesByLevel(level int32) ([]gen.Category, error) {
	db := getDao().Db()
	return model.ListCategoriesByLevel(db, level)
}

// GetCategoriesByParentID 获取指定父分类下的子分类
func GetCategoriesByParentID(parentID *int32) ([]gen.Category, error) {
	db := getDao().Db()
	return model.GetCategoriesByParentID(db, parentID)
}

// GetCategoryTree 获取分类树
func GetCategoryTree() ([]*model.CategoryNode, error) {
	db := getDao().Db()
	return model.GetCategoryTree(db)
}

// UpdateCategory 更新分类
func UpdateCategory(category *gen.Category) error {
	db := getDao().Db()
	return model.UpdateCategory(db, category)
}

// MoveCategory 移动分类
func MoveCategory(categoryID int32, newParentID *int32) error {
	db := getDao().Db()
	return model.MoveCategory(db, categoryID, newParentID)
}
