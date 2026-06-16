package ports

import (
	"context"
	"time"

	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/go-webauthn/webauthn/webauthn"
)

// SessionTokenStore validates bearer tokens without exposing the full cache.
type SessionTokenStore interface {
	GetCache(key string) (any, bool)
}

// CacheStore provides read/write access to the ephemeral cache.
type CacheStore interface {
	WriteCache(key string, value any, seconds int64)
	GetCache(key string) (any, bool)
}

// TimelogRepository handles persistence for time logs.
type TimelogRepository interface {
	CreateTimeLog(tl *domain.TimeLog) error
	GetTimeLogByID(id int32) (*domain.TimeLog, error)
	ListTimeLogs(conds ...interface{}) ([]domain.TimeLog, error)
	ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]domain.TimeLog, error)
	ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]domain.TimeLog, error)
	UpdateTimeLog(tl *domain.TimeLog) error
	DeleteTimeLog(id int32) error
}

// CategoryRepository handles persistence for categories.
type CategoryRepository interface {
	CreateCategory(category *domain.Category) error
	GetCategoryByID(id int32) (*domain.Category, error)
	GetCategoryByName(name string, parentID *int32) (*domain.Category, error)
	ListCategories(conds ...interface{}) ([]domain.Category, error)
	ListCategoriesByLevel(level int32) ([]domain.Category, error)
	GetCategoriesByParentID(parentID *int32) ([]domain.Category, error)
	GetCategoryTree() ([]*domain.CategoryNode, error)
	UpdateCategory(category *domain.Category) error
	MoveCategory(categoryID int32, newParentID *int32) error
}

// TaskRepository handles persistence for tasks.
type TaskRepository interface {
	CreateTask(task *domain.Task) error
	GetTaskByID(id int32) (*domain.Task, error)
	GetAllTasks(includeSuspended, includeCompleted bool) ([]domain.Task, error)
	GetTasksByDate(date time.Time, includeSuspended, includeCompleted bool) ([]domain.Task, error)
	GetTasksByDateRange(startDate, endDate time.Time) ([]domain.Task, error)
	UpdateTask(task *domain.Task) error
	DeleteTask(id int32) error
	MarkTaskAsCompleted(taskID int32) error
	MarkTaskAsIncomplete(taskID int32) error
	SuspendTask(taskID int32) error
	UnsuspendTask(taskID int32) error
	GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]domain.Task, error)
	GetTaskStats(date time.Time) (map[string]interface{}, error)
}

// ConstraintRepository handles persistence for constraints.
type ConstraintRepository interface {
	CreateConstraint(constraint *domain.Constraint) error
	GetConstraintByID(id int32) (*domain.Constraint, error)
	GetAllConstraints() ([]domain.Constraint, error)
	GetActiveConstraints() ([]domain.Constraint, error)
	GetConstraintsByDateRange(startDate, endDate time.Time) ([]domain.Constraint, error)
	UpdateConstraint(constraint *domain.Constraint) error
	DeleteConstraint(id int32) error
	MarkConstraintAsCompleted(constraintID int32, endReason string) error
	MarkConstraintAsActive(constraintID int32) error
}

// MetricRepository handles persistence for behavior metrics and their records.
type MetricRepository interface {
	CreateMetric(metric *domain.Metric) error
	GetMetricByID(id int32) (*domain.Metric, error)
	GetMetricByName(name string) (*domain.Metric, error)
	ListMetrics() ([]domain.Metric, error)
	UpdateMetric(metric *domain.Metric) error
	DeleteMetric(id int32) error
	CreateMetricRecord(record *domain.MetricRecord) error
	ListMetricRecordsByMetricID(metricID int32) ([]domain.MetricRecord, error)
	UpdateMetricCurrentValue(metricID int32, value float64, recordedAt time.Time) error
}

// PasskeyCredentialRepository handles persistence for WebAuthn credentials.
type PasskeyCredentialRepository interface {
	CreateWebAuthnCredential(credential *domain.PasskeyCredential) error
	GetWebAuthnCredentialByCredentialID(credentialID []byte) (*domain.PasskeyCredential, error)
	ListWebAuthnCredentials() ([]domain.PasskeyCredential, error)
	DeleteWebAuthnCredential(id int32) error
	UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error
}

// TempPasswordRepository handles persistence for temporary passwords.
type TempPasswordRepository interface {
	CreateTempPassword(tempPassword *domain.TempPassword) error
	ListTempPasswords() ([]domain.TempPassword, error)
	DeleteTempPassword(id int32) error
	DeleteExpiredTempPasswords(now time.Time) error
	GetTempPasswordByHash(hash string, now time.Time) (*domain.TempPassword, error)
}

// UnitOfWorkRepositories exposes repository adapters bound to a single
// transaction. It is passed to the callback executed inside UnitOfWork.Run.
type UnitOfWorkRepositories struct {
	TimelogRepo      TimelogRepository
	CategoryRepo     CategoryRepository
	TaskRepo         TaskRepository
	ConstraintRepo   ConstraintRepository
	MetricRepo       MetricRepository
	PasskeyRepo      PasskeyCredentialRepository
	TempPasswordRepo TempPasswordRepository
}

// UnitOfWork abstracts the persistence transaction boundary. Implementations
// should execute the provided callback inside a transaction and roll back on
// any non-nil error returned by the callback.
type UnitOfWork interface {
	Run(ctx context.Context, fn func(repos UnitOfWorkRepositories) error) error
}
