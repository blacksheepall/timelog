package ports

import (
	"time"

	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
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
	CreateTimeLog(tl *gen.Timelog) error
	GetTimeLogByID(id int32) (*gen.Timelog, error)
	ListTimeLogs(conds ...interface{}) ([]gen.Timelog, error)
	ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error)
	ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]gen.Timelog, error)
	UpdateTimeLog(tl *gen.Timelog) error
	DeleteTimeLog(id int32) error
}

// CategoryRepository handles persistence for categories.
type CategoryRepository interface {
	CreateCategory(category *gen.Category) error
	GetCategoryByID(id int32) (*gen.Category, error)
	GetCategoryByName(name string, parentID *int32) (*gen.Category, error)
	ListCategories(conds ...interface{}) ([]gen.Category, error)
	ListCategoriesByLevel(level int32) ([]gen.Category, error)
	GetCategoriesByParentID(parentID *int32) ([]gen.Category, error)
	GetCategoryTree() ([]*model.CategoryNode, error)
	UpdateCategory(category *gen.Category) error
	MoveCategory(categoryID int32, newParentID *int32) error
}

// TaskRepository handles persistence for tasks.
type TaskRepository interface {
	CreateTask(task *gen.Task) error
	GetTaskByID(id int32) (*gen.Task, error)
	GetAllTasks(includeSuspended, includeCompleted bool) ([]gen.Task, error)
	GetTasksByDate(date time.Time, includeSuspended, includeCompleted bool) ([]gen.Task, error)
	GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error)
	UpdateTask(task *gen.Task) error
	DeleteTask(id int32) error
	MarkTaskAsCompleted(taskID int32) error
	MarkTaskAsIncomplete(taskID int32) error
	SuspendTask(taskID int32) error
	UnsuspendTask(taskID int32) error
	GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error)
	GetTaskStats(date time.Time) (map[string]interface{}, error)
}

// ConstraintRepository handles persistence for constraints.
type ConstraintRepository interface {
	CreateConstraint(constraint *gen.Constraint) error
	GetConstraintByID(id int32) (*gen.Constraint, error)
	GetAllConstraints() ([]gen.Constraint, error)
	GetActiveConstraints() ([]gen.Constraint, error)
	GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error)
	UpdateConstraint(constraint *gen.Constraint) error
	DeleteConstraint(id int32) error
	MarkConstraintAsCompleted(constraintID int32, endReason string) error
	MarkConstraintAsActive(constraintID int32) error
}

// PasskeyCredentialRepository handles persistence for WebAuthn credentials.
type PasskeyCredentialRepository interface {
	CreateWebAuthnCredential(credential *model.WebAuthnCredential) error
	GetWebAuthnCredentialByCredentialID(credentialID []byte) (*model.WebAuthnCredential, error)
	ListWebAuthnCredentials() ([]model.WebAuthnCredential, error)
	DeleteWebAuthnCredential(id uint) error
	UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error
}

// TempPasswordRepository handles persistence for temporary passwords.
type TempPasswordRepository interface {
	CreateTempPassword(tempPassword *model.TempPassword) error
	ListTempPasswords() ([]model.TempPassword, error)
	DeleteTempPassword(id uint) error
	DeleteExpiredTempPasswords(now time.Time) error
	GetTempPasswordByHash(hash string, now time.Time) (*model.TempPassword, error)
}
