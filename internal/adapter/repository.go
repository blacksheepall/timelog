package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// Repositories groups all per-domain SQLite adapters into a single composable
// structure. It is returned by NewRepositories so callers that only need one
// adapter can eventually migrate to the smaller constructors, while existing
// code continues to compile.
type Repositories struct {
	*cacheStore
	*timelogRepo
	*categoryRepo
	*taskRepo
	*constraintRepo
	*passkeyCredentialRepo
	*tempPasswordRepo
}

var (
	_ ports.TimelogRepository           = (*Repositories)(nil)
	_ ports.CategoryRepository          = (*Repositories)(nil)
	_ ports.TaskRepository              = (*Repositories)(nil)
	_ ports.ConstraintRepository        = (*Repositories)(nil)
	_ ports.PasskeyCredentialRepository = (*Repositories)(nil)
	_ ports.TempPasswordRepository      = (*Repositories)(nil)
	_ ports.CacheStore                  = (*Repositories)(nil)
	_ ports.SessionTokenStore           = (*Repositories)(nil)
)

// NewRepositories creates all repository adapters backed by dao.
func NewRepositories(dao *model.Dao) *Repositories {
	return &Repositories{
		cacheStore:            newCacheStore(dao),
		timelogRepo:           newTimelogRepo(dao),
		categoryRepo:          newCategoryRepo(dao),
		taskRepo:              newTaskRepo(dao),
		constraintRepo:        newConstraintRepo(dao),
		passkeyCredentialRepo: newPasskeyCredentialRepo(dao),
		tempPasswordRepo:      newTempPasswordRepo(dao),
	}
}
