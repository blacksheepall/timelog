package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// taskRepo implements ports.TaskRepository using the model layer.
type taskRepo struct {
	dao *model.Dao
}

var _ ports.TaskRepository = (*taskRepo)(nil)

func newTaskRepo(dao *model.Dao) *taskRepo {
	return &taskRepo{dao: dao}
}

func (r *taskRepo) CreateTask(task *gen.Task) error {
	return model.CreateTask(r.dao.Db(), task)
}

func (r *taskRepo) GetTaskByID(id int32) (*gen.Task, error) {
	return model.GetTaskByID(r.dao.Db(), id)
}

func (r *taskRepo) GetAllTasks(includeSuspended, includeCompleted bool) ([]gen.Task, error) {
	return model.GetAllTasks(r.dao.Db(), includeSuspended, includeCompleted)
}

func (r *taskRepo) GetTasksByDate(date time.Time, includeSuspended, includeCompleted bool) ([]gen.Task, error) {
	return model.GetTasksByDate(r.dao.Db(), date, includeSuspended, includeCompleted)
}

func (r *taskRepo) GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return model.GetTasksByDateRange(r.dao.Db(), startDate, endDate)
}

func (r *taskRepo) UpdateTask(task *gen.Task) error {
	return model.UpdateTask(r.dao.Db(), task)
}

func (r *taskRepo) DeleteTask(id int32) error {
	return model.DeleteTask(r.dao.Db(), id)
}

func (r *taskRepo) MarkTaskAsCompleted(taskID int32) error {
	return model.MarkTaskAsCompleted(r.dao.Db(), taskID)
}

func (r *taskRepo) MarkTaskAsIncomplete(taskID int32) error {
	return model.MarkTaskAsIncomplete(r.dao.Db(), taskID)
}

func (r *taskRepo) SuspendTask(taskID int32) error {
	return model.SuspendTask(r.dao.Db(), taskID)
}

func (r *taskRepo) UnsuspendTask(taskID int32) error {
	return model.UnsuspendTask(r.dao.Db(), taskID)
}

func (r *taskRepo) GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return model.GetCompletedTasksInDateRange(r.dao.Db(), startDate, endDate)
}

func (r *taskRepo) GetTaskStats(date time.Time) (map[string]interface{}, error) {
	return model.GetTaskStats(r.dao.Db(), date)
}
