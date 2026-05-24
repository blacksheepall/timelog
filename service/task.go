package service

import (
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// CreateTask 创建任务
func CreateTask(task *gen.Task) error {
	dao := getDao()
	return model.CreateTask(dao.Db(), task)
}

// GetTaskByID 根据ID获取任务
func GetTaskByID(id int32) (*gen.Task, error) {
	dao := getDao()
	return model.GetTaskByID(dao.Db(), id)
}

// GetAllTasks 获取所有任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func GetAllTasks(includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	dao := getDao()
	return model.GetAllTasks(dao.Db(), includeSuspended, includeCompleted)
}

// GetTasksByDate 根据日期获取任务
// includeSuspended: 是否包含暂停的任务
// includeCompleted: 是否包含已完成的任务
func GetTasksByDate(date time.Time, includeSuspended bool, includeCompleted bool) ([]gen.Task, error) {
	dao := getDao()
	return model.GetTasksByDate(dao.Db(), date, includeSuspended, includeCompleted)
}

// GetTasksByDateRange 根据日期范围获取任务
func GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	dao := getDao()
	return model.GetTasksByDateRange(dao.Db(), startDate, endDate)
}

// UpdateTask 更新任务
func UpdateTask(task *gen.Task) error {
	dao := getDao()
	return model.UpdateTask(dao.Db(), task)
}

// DeleteTask 删除任务
func DeleteTask(id int32) error {
	dao := getDao()
	return model.DeleteTask(dao.Db(), id)
}

// MarkTaskAsCompleted 标记任务为完成
func MarkTaskAsCompleted(taskID int32) error {
	dao := getDao()
	return model.MarkTaskAsCompleted(dao.Db(), taskID)
}

// MarkTaskAsIncomplete 标记任务为未完成
func MarkTaskAsIncomplete(taskID int32) error {
	dao := getDao()
	return model.MarkTaskAsIncomplete(dao.Db(), taskID)
}

// SuspendTask 暂停任务
func SuspendTask(taskID int32) error {
	dao := getDao()
	return model.SuspendTask(dao.Db(), taskID)
}

// UnsuspendTask 取消暂停任务
func UnsuspendTask(taskID int32) error {
	dao := getDao()
	return model.UnsuspendTask(dao.Db(), taskID)
}

// GetCompletedTasksInDateRange 获取指定日期范围内的已完成任务
func GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	dao := getDao()
	return model.GetCompletedTasksInDateRange(dao.Db(), startDate, endDate)
}

// GetTaskStats 获取任务统计信息
func GetTaskStats(date time.Time) (map[string]interface{}, error) {
	dao := getDao()
	return model.GetTaskStats(dao.Db(), date)
}

// CompleteTaskWithTimelog 完成任务并创建时间记录
// 这是一个组合操作，将任务标记为完成，并可选地创建关联的时间记录
func CompleteTaskWithTimelog(taskID int32, createTimelog bool, timelogData *gen.Timelog) error {
	dao := getDao()

	// 开始事务
	tx := dao.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// 标记任务为完成
	if err := model.MarkTaskAsCompleted(tx.Db(), taskID); err != nil {
		tx.Rollback()
		return err
	}

	// 如果需要创建时间记录
	if createTimelog && timelogData != nil {
		timelogData.TaskID = &taskID
		if err := model.CreateTimeLog(tx.Db(), timelogData); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
