package service

import (
	"github.com/blacksheepaul/timelog/internal/domain"
)

// --- TimeLog Service ---

// CreateTimeLog 新增一条时间日志
func (s *Service) CreateTimeLog(tl *domain.TimeLog) error {
	return s.timelogRepo.CreateTimeLog(tl)
}

// GetTimeLogByID 根据ID获取时间日志
func (s *Service) GetTimeLogByID(id int32) (*domain.TimeLog, error) {
	return s.timelogRepo.GetTimeLogByID(id)
}

// ListTimeLogs 查询时间日志（可扩展条件）
func (s *Service) ListTimeLogs(conds ...interface{}) ([]domain.TimeLog, error) {
	return s.timelogRepo.ListTimeLogs(conds...)
}

// ListTimeLogsWithOptions 查询时间日志（支持排序和限制）
func (s *Service) ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]domain.TimeLog, error) {
	return s.timelogRepo.ListTimeLogsWithOptions(limit, orderBy, conds...)
}

// ListTimeLogsByLocalDateRange queries timelogs whose start_time falls within
// the inclusive local-date range in Asia/Singapore (YYYY-MM-DD strings).
func (s *Service) ListTimeLogsByLocalDateRange(startDate, endDate string) ([]domain.TimeLog, error) {
	return s.timelogRepo.ListTimeLogsByLocalDateRange(startDate, endDate)
}

// UpdateTimeLog 更新一条时间日志
func (s *Service) UpdateTimeLog(tl *domain.TimeLog) error {
	return s.timelogRepo.UpdateTimeLog(tl)
}

// DeleteTimeLog 删除一条时间日志
func (s *Service) DeleteTimeLog(id int32) error {
	return s.timelogRepo.DeleteTimeLog(id)
}

// --- Category Service ---

// CreateCategory 创建分类
func (s *Service) CreateCategory(category *domain.Category) error {
	return s.categoryRepo.CreateCategory(category)
}

// GetCategoryByID 根据ID获取分类
func (s *Service) GetCategoryByID(id int32) (*domain.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

// GetCategoryByName 根据名称获取分类
func (s *Service) GetCategoryByName(name string, parentID *int32) (*domain.Category, error) {
	return s.categoryRepo.GetCategoryByName(name, parentID)
}

// ListCategories 查询所有分类
func (s *Service) ListCategories(conds ...interface{}) ([]domain.Category, error) {
	return s.categoryRepo.ListCategories(conds...)
}

// ListCategoriesByLevel 按层级查询分类
func (s *Service) ListCategoriesByLevel(level int32) ([]domain.Category, error) {
	return s.categoryRepo.ListCategoriesByLevel(level)
}

// GetCategoriesByParentID 获取指定父分类下的子分类
func (s *Service) GetCategoriesByParentID(parentID *int32) ([]domain.Category, error) {
	return s.categoryRepo.GetCategoriesByParentID(parentID)
}

// GetCategoryTree 获取分类树
func (s *Service) GetCategoryTree() ([]*domain.CategoryNode, error) {
	return s.categoryRepo.GetCategoryTree()
}

// UpdateCategory 更新分类
func (s *Service) UpdateCategory(category *domain.Category) error {
	return s.categoryRepo.UpdateCategory(category)
}

// MoveCategory 移动分类
func (s *Service) MoveCategory(categoryID int32, newParentID *int32) error {
	return s.categoryRepo.MoveCategory(categoryID, newParentID)
}
