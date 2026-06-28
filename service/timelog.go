package service

import (
	"context"
	"fmt"

	"github.com/blacksheepaul/timelog/core/audit"
	"github.com/blacksheepaul/timelog/internal/domain"
)

// --- TimeLog Service ---

// CreateTimeLog 新增一条时间日志
func (s *Service) CreateTimeLog(ctx context.Context, tl *domain.TimeLog) error {
	if err := s.timelogRepo.CreateTimeLog(tl); err != nil {
		return err
	}
	s.logAudit(ctx, audit.ActionCreate, audit.EntityTimeLog, tl.ID, timelogPayload(tl))
	return nil
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
func (s *Service) UpdateTimeLog(ctx context.Context, tl *domain.TimeLog) error {
	if _, err := s.timelogRepo.GetTimeLogByID(tl.ID); err != nil {
		return fmt.Errorf("time log not found: %w", err)
	}
	if err := s.timelogRepo.UpdateTimeLog(tl); err != nil {
		return err
	}
	s.logAudit(ctx, audit.ActionUpdate, audit.EntityTimeLog, tl.ID, timelogPayload(tl))
	return nil
}

// DeleteTimeLog 删除一条时间日志
func (s *Service) DeleteTimeLog(ctx context.Context, id int32) error {
	if err := s.timelogRepo.DeleteTimeLog(id); err != nil {
		return err
	}
	s.logAudit(ctx, audit.ActionDelete, audit.EntityTimeLog, id, nil)
	return nil
}

// timelogPayload builds a serializable snapshot of a TimeLog for audit records.
func timelogPayload(tl *domain.TimeLog) map[string]any {
	if tl == nil {
		return nil
	}
	p := map[string]any{
		"category_id": tl.CategoryID,
		"start_time":  tl.StartTime,
		"remark":      tl.Remark,
	}
	if tl.EndTime != nil {
		p["end_time"] = *tl.EndTime
	}
	if tl.TaskID != nil {
		p["task_id"] = *tl.TaskID
	}
	return p
}

// logAudit records an auditable event when an audit logger is configured.
// Failures are swallowed so that audit problems do not break business logic.
func (s *Service) logAudit(ctx context.Context, action, entityType string, entityID int32, payload map[string]any) {
	if s.auditLogger == nil {
		return
	}
	_ = s.auditLogger.Log(ctx, action, entityType, entityID, payload)
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
