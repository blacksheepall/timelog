package adapter

import (
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
	"github.com/blacksheepaul/timelog/model/gen"
)

// categoryRepo implements ports.CategoryRepository using the model layer.
type categoryRepo struct {
	dao *model.Dao
}

var _ ports.CategoryRepository = (*categoryRepo)(nil)

func newCategoryRepo(dao *model.Dao) *categoryRepo {
	return &categoryRepo{dao: dao}
}

func (r *categoryRepo) CreateCategory(category *gen.Category) error {
	return model.CreateCategory(r.dao.Db(), category)
}

func (r *categoryRepo) GetCategoryByID(id int32) (*gen.Category, error) {
	return model.GetCategoryByID(r.dao.Db(), id)
}

func (r *categoryRepo) GetCategoryByName(name string, parentID *int32) (*gen.Category, error) {
	return model.GetCategoryByName(r.dao.Db(), name, parentID)
}

func (r *categoryRepo) ListCategories(conds ...interface{}) ([]gen.Category, error) {
	return model.ListCategories(r.dao.Db(), conds...)
}

func (r *categoryRepo) ListCategoriesByLevel(level int32) ([]gen.Category, error) {
	return model.ListCategoriesByLevel(r.dao.Db(), level)
}

func (r *categoryRepo) GetCategoriesByParentID(parentID *int32) ([]gen.Category, error) {
	return model.GetCategoriesByParentID(r.dao.Db(), parentID)
}

func (r *categoryRepo) GetCategoryTree() ([]*model.CategoryNode, error) {
	return model.GetCategoryTree(r.dao.Db())
}

func (r *categoryRepo) UpdateCategory(category *gen.Category) error {
	return model.UpdateCategory(r.dao.Db(), category)
}

func (r *categoryRepo) MoveCategory(categoryID int32, newParentID *int32) error {
	return model.MoveCategory(r.dao.Db(), categoryID, newParentID)
}
