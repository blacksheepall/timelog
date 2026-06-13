package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// categoryRepo implements ports.CategoryRepository using the model layer.
type categoryRepo struct {
	db model.DBProvider
}

var _ ports.CategoryRepository = (*categoryRepo)(nil)

func newCategoryRepo(db model.DBProvider) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(category *gen.Category) error {
	return model.CreateCategory(r.db.Db(), category)
}

func (r *categoryRepo) GetCategoryByID(id int32) (*gen.Category, error) {
	return model.GetCategoryByID(r.db.Db(), id)
}

func (r *categoryRepo) GetCategoryByName(name string, parentID *int32) (*gen.Category, error) {
	return model.GetCategoryByName(r.db.Db(), name, parentID)
}

func (r *categoryRepo) ListCategories(conds ...interface{}) ([]gen.Category, error) {
	return model.ListCategories(r.db.Db(), conds...)
}

func (r *categoryRepo) ListCategoriesByLevel(level int32) ([]gen.Category, error) {
	return model.ListCategoriesByLevel(r.db.Db(), level)
}

func (r *categoryRepo) GetCategoriesByParentID(parentID *int32) ([]gen.Category, error) {
	return model.GetCategoriesByParentID(r.db.Db(), parentID)
}

func (r *categoryRepo) GetCategoryTree() ([]*model.CategoryNode, error) {
	return model.GetCategoryTree(r.db.Db())
}

func (r *categoryRepo) UpdateCategory(category *gen.Category) error {
	return model.UpdateCategory(r.db.Db(), category)
}

func (r *categoryRepo) MoveCategory(categoryID int32, newParentID *int32) error {
	return model.MoveCategory(r.db.Db(), categoryID, newParentID)
}
