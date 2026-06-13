package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// cacheStore implements ephemeral cache ports backed by the DAO cache.
type cacheStore struct {
	dao *model.Dao
}

var (
	_ ports.CacheStore        = (*cacheStore)(nil)
	_ ports.SessionTokenStore = (*cacheStore)(nil)
)

// newCacheStore creates a cache adapter backed by the given DAO.
func newCacheStore(dao *model.Dao) *cacheStore {
	return &cacheStore{dao: dao}
}

func (c *cacheStore) WriteCache(key string, value any, seconds int64) {
	c.dao.WriteCache(key, value, seconds)
}

func (c *cacheStore) GetCache(key string) (any, bool) {
	return c.dao.GetCache(key)
}
