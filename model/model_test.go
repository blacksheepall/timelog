package model_test

import (
	"testing"

	"github.com/blacksheepaul/timelog/internal/testutil"
)

func TestNewDaoOpensInMemoryDatabase(t *testing.T) {
	dao := testutil.NewTestDAO(t)
	if dao.Db() == nil {
		t.Fatal("expected non-nil gorm DB")
	}
	if dao.RawDB == nil {
		t.Fatal("expected non-nil raw DB")
	}
}
