package adapter

import (
	"time"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/go-webauthn/webauthn/webauthn"
)

// Repositories implements all repository ports using a single *model.Dao.
// It is a transitional adapter: in the future these may split into per-domain
// adapters or move to internal/adapter/sqlite.
type Repositories struct {
	dao *model.Dao
}

var (
	_ ports.TimelogRepository         = (*Repositories)(nil)
	_ ports.CategoryRepository        = (*Repositories)(nil)
	_ ports.TaskRepository            = (*Repositories)(nil)
	_ ports.ConstraintRepository      = (*Repositories)(nil)
	_ ports.PasskeyCredentialRepository = (*Repositories)(nil)
	_ ports.TempPasswordRepository    = (*Repositories)(nil)
	_ ports.CacheStore                = (*Repositories)(nil)
	_ ports.SessionTokenStore         = (*Repositories)(nil)
)

// NewRepositories creates a repository adapter backed by dao.
func NewRepositories(dao *model.Dao) *Repositories {
	return &Repositories{dao: dao}
}

// --- CacheStore / SessionTokenStore ---

func (r *Repositories) WriteCache(key string, value any, seconds int64) {
	r.dao.WriteCache(key, value, seconds)
}

func (r *Repositories) GetCache(key string) (any, bool) {
	return r.dao.GetCache(key)
}

// --- TimelogRepository ---

func (r *Repositories) CreateTimeLog(tl *gen.Timelog) error {
	return model.CreateTimeLog(r.dao.Db(), tl)
}

func (r *Repositories) GetTimeLogByID(id int32) (*gen.Timelog, error) {
	return model.GetTimeLogByID(r.dao.Db(), id)
}

func (r *Repositories) ListTimeLogs(conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogs(r.dao.Db(), conds...)
}

func (r *Repositories) ListTimeLogsWithOptions(limit int, orderBy string, conds ...interface{}) ([]gen.Timelog, error) {
	return model.ListTimeLogsWithOptions(r.dao.Db(), limit, orderBy, conds...)
}

func (r *Repositories) ListTimeLogsByLocalDateRange(startDateStr, endDateStr string) ([]gen.Timelog, error) {
	return model.ListTimeLogsByLocalDateRange(r.dao.Db(), startDateStr, endDateStr)
}

func (r *Repositories) UpdateTimeLog(tl *gen.Timelog) error {
	return model.UpdateTimeLog(r.dao.Db(), tl)
}

func (r *Repositories) DeleteTimeLog(id int32) error {
	return model.DeleteTimeLog(r.dao.Db(), id)
}

// --- CategoryRepository ---

func (r *Repositories) CreateCategory(category *gen.Category) error {
	return model.CreateCategory(r.dao.Db(), category)
}

func (r *Repositories) GetCategoryByID(id int32) (*gen.Category, error) {
	return model.GetCategoryByID(r.dao.Db(), id)
}

func (r *Repositories) GetCategoryByName(name string, parentID *int32) (*gen.Category, error) {
	return model.GetCategoryByName(r.dao.Db(), name, parentID)
}

func (r *Repositories) ListCategories(conds ...interface{}) ([]gen.Category, error) {
	return model.ListCategories(r.dao.Db(), conds...)
}

func (r *Repositories) ListCategoriesByLevel(level int32) ([]gen.Category, error) {
	return model.ListCategoriesByLevel(r.dao.Db(), level)
}

func (r *Repositories) GetCategoriesByParentID(parentID *int32) ([]gen.Category, error) {
	return model.GetCategoriesByParentID(r.dao.Db(), parentID)
}

func (r *Repositories) GetCategoryTree() ([]*model.CategoryNode, error) {
	return model.GetCategoryTree(r.dao.Db())
}

func (r *Repositories) UpdateCategory(category *gen.Category) error {
	return model.UpdateCategory(r.dao.Db(), category)
}

func (r *Repositories) MoveCategory(categoryID int32, newParentID *int32) error {
	return model.MoveCategory(r.dao.Db(), categoryID, newParentID)
}

// --- TaskRepository ---

func (r *Repositories) CreateTask(task *gen.Task) error {
	return model.CreateTask(r.dao.Db(), task)
}

func (r *Repositories) GetTaskByID(id int32) (*gen.Task, error) {
	return model.GetTaskByID(r.dao.Db(), id)
}

func (r *Repositories) GetAllTasks(includeSuspended, includeCompleted bool) ([]gen.Task, error) {
	return model.GetAllTasks(r.dao.Db(), includeSuspended, includeCompleted)
}

func (r *Repositories) GetTasksByDate(date time.Time, includeSuspended, includeCompleted bool) ([]gen.Task, error) {
	return model.GetTasksByDate(r.dao.Db(), date, includeSuspended, includeCompleted)
}

func (r *Repositories) GetTasksByDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return model.GetTasksByDateRange(r.dao.Db(), startDate, endDate)
}

func (r *Repositories) UpdateTask(task *gen.Task) error {
	return model.UpdateTask(r.dao.Db(), task)
}

func (r *Repositories) DeleteTask(id int32) error {
	return model.DeleteTask(r.dao.Db(), id)
}

func (r *Repositories) MarkTaskAsCompleted(taskID int32) error {
	return model.MarkTaskAsCompleted(r.dao.Db(), taskID)
}

func (r *Repositories) MarkTaskAsIncomplete(taskID int32) error {
	return model.MarkTaskAsIncomplete(r.dao.Db(), taskID)
}

func (r *Repositories) SuspendTask(taskID int32) error {
	return model.SuspendTask(r.dao.Db(), taskID)
}

func (r *Repositories) UnsuspendTask(taskID int32) error {
	return model.UnsuspendTask(r.dao.Db(), taskID)
}

func (r *Repositories) GetCompletedTasksInDateRange(startDate, endDate time.Time) ([]gen.Task, error) {
	return model.GetCompletedTasksInDateRange(r.dao.Db(), startDate, endDate)
}

func (r *Repositories) GetTaskStats(date time.Time) (map[string]interface{}, error) {
	return model.GetTaskStats(r.dao.Db(), date)
}

// --- ConstraintRepository ---

func (r *Repositories) CreateConstraint(constraint *gen.Constraint) error {
	return model.CreateConstraint(r.dao.Db(), constraint)
}

func (r *Repositories) GetConstraintByID(id int32) (*gen.Constraint, error) {
	return model.GetConstraintByID(r.dao.Db(), id)
}

func (r *Repositories) GetAllConstraints() ([]gen.Constraint, error) {
	return model.GetAllConstraints(r.dao.Db())
}

func (r *Repositories) GetActiveConstraints() ([]gen.Constraint, error) {
	return model.GetActiveConstraints(r.dao.Db())
}

func (r *Repositories) GetConstraintsByDateRange(startDate, endDate time.Time) ([]gen.Constraint, error) {
	return model.GetConstraintsByDateRange(r.dao.Db(), startDate, endDate)
}

func (r *Repositories) UpdateConstraint(constraint *gen.Constraint) error {
	return model.UpdateConstraint(r.dao.Db(), constraint)
}

func (r *Repositories) DeleteConstraint(id int32) error {
	return model.DeleteConstraint(r.dao.Db(), id)
}

func (r *Repositories) MarkConstraintAsCompleted(constraintID int32, endReason string) error {
	return model.MarkConstraintAsCompleted(r.dao.Db(), constraintID, endReason)
}

func (r *Repositories) MarkConstraintAsActive(constraintID int32) error {
	return model.MarkConstraintAsActive(r.dao.Db(), constraintID)
}

// --- PasskeyCredentialRepository ---

func (r *Repositories) CreateWebAuthnCredential(credential *model.WebAuthnCredential) error {
	return model.CreateWebAuthnCredential(r.dao.Db(), credential)
}

func (r *Repositories) GetWebAuthnCredentialByCredentialID(credentialID []byte) (*model.WebAuthnCredential, error) {
	return model.GetWebAuthnCredentialByCredentialID(r.dao.Db(), credentialID)
}

func (r *Repositories) ListWebAuthnCredentials() ([]model.WebAuthnCredential, error) {
	return model.ListWebAuthnCredentials(r.dao.Db())
}

func (r *Repositories) DeleteWebAuthnCredential(id uint) error {
	return model.DeleteWebAuthnCredential(r.dao.Db(), id)
}

func (r *Repositories) UpdateWebAuthnCredentialAuth(credentialID []byte, credential *webauthn.Credential) error {
	return model.UpdateWebAuthnCredentialAuth(r.dao.Db(), credentialID, credential)
}

// --- TempPasswordRepository ---

func (r *Repositories) CreateTempPassword(tempPassword *model.TempPassword) error {
	return model.CreateTempPassword(r.dao.Db(), tempPassword)
}

func (r *Repositories) ListTempPasswords() ([]model.TempPassword, error) {
	return model.ListTempPasswords(r.dao.Db())
}

func (r *Repositories) DeleteTempPassword(id uint) error {
	return model.DeleteTempPassword(r.dao.Db(), id)
}

func (r *Repositories) DeleteExpiredTempPasswords(now time.Time) error {
	return model.DeleteExpiredTempPasswords(r.dao.Db(), now)
}

func (r *Repositories) GetTempPasswordByHash(hash string, now time.Time) (*model.TempPassword, error) {
	return model.GetTempPasswordByHash(r.dao.Db(), hash, now)
}
