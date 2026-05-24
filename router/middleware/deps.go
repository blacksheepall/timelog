package middleware

import "github.com/blacksheepaul/timelog/model"

var middlewareDaoProvider = model.GetDao

func getMiddlewareDAO() *model.Dao {
	return middlewareDaoProvider()
}

func setMiddlewareDAOProviderForTest(provider func() *model.Dao) {
	if provider == nil {
		middlewareDaoProvider = model.GetDao
		return
	}
	middlewareDaoProvider = provider
}
