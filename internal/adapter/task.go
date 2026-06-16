package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// taskRepo implements ports.TaskRepository using the model layer.
type taskRepo struct {
	db model.DBProvider
}

var _ ports.TaskRepository = (*taskRepo)(nil)

func newTaskRepo(db model.DBProvider) *taskRepo {
	return &taskRepo{db: db}
}

func (r *taskRepo) CreateTask(task *domain.Task) error {
	g := toGenTask(task)
	if err := model.CreateTask(r.db.Db(), g); err != nil {
		return err
	}
	*task = *toDomainTask(g)
	return nil
}

func (r *taskRepo) GetTaskByID(id int32) (*domain.Task, error) {
	g, err := model.GetTaskByID(r.db.Db(), id)
	if err != nil {
		return nil, err
	}
	return toDomainTask(g), nil
}

func (r *taskRepo) GetAllTasks(includeSuspended, includeCompleted bool) ([]domain.Task, error) {
	list, err := model.GetAllTasks(r.db.Db(), includeSuspended, includeCompleted)
	if err != nil {
		return nil, err
	}
	return toDomainTasks(list), nil
}

func (r *taskRepo) GetTasksByDate(date time.Time, includeSuspended, includeCompleted bool) ([]domain.Task, error) {
	list, err := model.GetTasksByDate(r.db.Db(), date, includeSuspended, includeCompleted)
	if err != nil {
		return nil, err
	}
	return toDomainTasks(list), nil
}

func (r *taskRepo) GetTasksByDateRange(startDate, endDate time.Time) ([]domain.Task, error) {
	list, err := model.GetTasksByDateRange(r.db.Db(), startDate, endDate)
	if err != nil {
		return nil, err
	}
	return toDomainTasks(list), nil
}

func (r *taskRepo) UpdateTask(task *domain.Task) error {
	return model.UpdateTask(r.db.Db(), toGenTask(task))
}

func (r *taskRepo) DeleteTask(id int32) error {
	return model.DeleteTask(r.db.Db(), id)
}

func (r *taskRepo) MarkTaskAsCompleted(taskID int32) error {
	return model.MarkTaskAsCompleted(r.db.Db(), taskID)
}

func (r *taskRepo) MarkTaskAsIncomplete(taskID int32) error {
	return model.MarkTaskAsIncomplete(r.db.Db(), taskID)
}

func (r *taskRepo) SuspendTask(taskID int32) error {
	return model.SuspendTask(r.db.Db(), taskID)
}

func (r *taskRepo) UnsuspendTask(taskID int32) error {
	return model.UnsuspendTask(r.db.Db(), taskID)
}

func (r *taskRepo) GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]domain.Task, error) {
	list, err := model.GetCompletedTasksInDateRange(r.db.Db(), startDate, endDate)
	if err != nil {
		return nil, err
	}
	return toDomainTasks(list), nil
}

func (r *taskRepo) GetTaskStats(date time.Time) (map[string]interface{}, error) {
	return model.GetTaskStats(r.db.Db(), date)
}
