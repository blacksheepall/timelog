package service

import (
	"context"
	"fmt"
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model/gen"
)

// CreateTask 创建任务
func (s *Service) CreateTask(task *gen.Task) error {
	return s.taskRepo.CreateTask(task)
}

// GetTaskByID 根据ID获取任务
func (s *Service) GetTaskByID(id int32) (*gen.Task, error) {
	return s.taskRepo.GetTaskByID(id)
}

// GetAllTasks 获取所有任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func (s *Service) GetAllTasks(includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	return s.taskRepo.GetAllTasks(includeSuspended, includeCompleted)
}

// ListTasksByCompletionStatus returns tasks filtered by completion state.
// status must be "completed", "pending", or "all".
func (s *Service) ListTasksByCompletionStatus(status string) ([]gen.Task, error) {
	tasks, err := s.taskRepo.GetAllTasks(true, true)
	if err != nil {
		return nil, err
	}
	if status == "all" {
		return tasks, nil
	}

	filtered := make([]gen.Task, 0, len(tasks))
	for _, task := range tasks {
		isCompleted := task.IsCompleted != nil && *task.IsCompleted
		switch status {
		case "completed":
			if isCompleted {
				filtered = append(filtered, task)
			}
		case "pending":
			if !isCompleted {
				filtered = append(filtered, task)
			}
		default:
			return nil, fmt.Errorf("invalid status %q: want completed, pending, or all", status)
		}
	}
	return filtered, nil
}

// GetTasksByDate 根据日期获取任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func (s *Service) GetTasksByDate(date time.Time, includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	return s.taskRepo.GetTasksByDate(date, includeSuspended, includeCompleted)
}

// GetTasksByDateRange 根据日期范围获取任务
func (s *Service) GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return s.taskRepo.GetTasksByDateRange(startDate, endDate)
}

// UpdateTask 更新任务
func (s *Service) UpdateTask(task *gen.Task) error {
	return s.taskRepo.UpdateTask(task)
}

// DeleteTask 删除任务
func (s *Service) DeleteTask(id int32) error {
	return s.taskRepo.DeleteTask(id)
}

// MarkTaskAsCompleted 标记任务为完成
func (s *Service) MarkTaskAsCompleted(taskID int32) error {
	return s.taskRepo.MarkTaskAsCompleted(taskID)
}

// MarkTaskAsIncomplete 标记任务为未完成
func (s *Service) MarkTaskAsIncomplete(taskID int32) error {
	return s.taskRepo.MarkTaskAsIncomplete(taskID)
}

// SuspendTask 暂停任务
func (s *Service) SuspendTask(taskID int32) error {
	return s.taskRepo.SuspendTask(taskID)
}

// UnsuspendTask 取消暂停任务
func (s *Service) UnsuspendTask(taskID int32) error {
	return s.taskRepo.UnsuspendTask(taskID)
}

// GetCompletedTasksInDateRange 获取指定日期范围内的已完成任务
func (s *Service) GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return s.taskRepo.GetCompletedTasksInDateRange(startDate, endDate)
}

// GetTaskStats 获取任务统计信息
func (s *Service) GetTaskStats(date time.Time) (map[string]interface{}, error) {
	return s.taskRepo.GetTaskStats(date)
}

// CompleteTaskWithTimelog 完成任务并创建时间记录
// 这是一个组合操作，将任务标记为完成，并可选地创建关联的时间记录
func (s *Service) CompleteTaskWithTimelog(taskID int32, createTimelog bool, timelogData *gen.Timelog) error {
	return s.unitOfWork.Run(context.Background(), func(repos ports.UnitOfWorkRepositories) error {
		// 标记任务为完成
		if err := repos.TaskRepo.MarkTaskAsCompleted(taskID); err != nil {
			return err
		}

		// 如果需要创建时间记录
		if createTimelog && timelogData != nil {
			timelogData.TaskID = &taskID
			if err := repos.TimelogRepo.CreateTimeLog(timelogData); err != nil {
				return err
			}
		}

		return nil
	})
}

// --- Package-level wrappers (transitional) ---

func CreateTask(task *gen.Task) error { return defaultService.CreateTask(task) }
func GetTaskByID(id int32) (*gen.Task, error) { return defaultService.GetTaskByID(id) }
func GetAllTasks(includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	return defaultService.GetAllTasks(includeSuspended, includeCompleted)
}
func ListTasksByCompletionStatus(status string) ([]gen.Task, error) {
	return defaultService.ListTasksByCompletionStatus(status)
}
func GetTasksByDate(date time.Time, includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	return defaultService.GetTasksByDate(date, includeSuspended, includeCompleted)
}
func GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return defaultService.GetTasksByDateRange(startDate, endDate)
}
func UpdateTask(task *gen.Task) error        { return defaultService.UpdateTask(task) }
func DeleteTask(id int32) error               { return defaultService.DeleteTask(id) }
func MarkTaskAsCompleted(taskID int32) error  { return defaultService.MarkTaskAsCompleted(taskID) }
func MarkTaskAsIncomplete(taskID int32) error { return defaultService.MarkTaskAsIncomplete(taskID) }
func SuspendTask(taskID int32) error          { return defaultService.SuspendTask(taskID) }
func UnsuspendTask(taskID int32) error        { return defaultService.UnsuspendTask(taskID) }
func GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return defaultService.GetCompletedTasksInDateRange(startDate, endDate)
}
func GetTaskStats(date time.Time) (map[string]interface{}, error) {
	return defaultService.GetTaskStats(date)
}
func CompleteTaskWithTimelog(taskID int32, createTimelog bool, timelogData *gen.Timelog) error {
	return defaultService.CompleteTaskWithTimelog(taskID, createTimelog, timelogData)
}
