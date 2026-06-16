package adapter

import (
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/model"
)

// categoryRepo implements ports.CategoryRepository using the model layer.
type categoryRepo struct {
	db model.DBProvider
}

var _ ports.CategoryRepository = (*categoryRepo)(nil)

func newCategoryRepo(db model.DBProvider) *categoryRepo {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) CreateCategory(category *domain.Category) error {
	g := toGenCategory(category)
	if err := model.CreateCategory(r.db.Db(), g); err != nil {
		return err
	}
	*category = *toDomainCategory(g)
	return nil
}

func (r *categoryRepo) GetCategoryByID(id int32) (*domain.Category, error) {
	g, err := model.GetCategoryByID(r.db.Db(), id)
	if err != nil {
		return nil, err
	}
	return toDomainCategory(g), nil
}

func (r *categoryRepo) GetCategoryByName(name string, parentID *int32) (*domain.Category, error) {
	g, err := model.GetCategoryByName(r.db.Db(), name, parentID)
	if err != nil {
		return nil, err
	}
	return toDomainCategory(g), nil
}

func (r *categoryRepo) ListCategories(conds ...interface{}) ([]domain.Category, error) {
	list, err := model.ListCategories(r.db.Db(), conds...)
	if err != nil {
		return nil, err
	}
	return toDomainCategories(list), nil
}

func (r *categoryRepo) ListCategoriesByLevel(level int32) ([]domain.Category, error) {
	list, err := model.ListCategoriesByLevel(r.db.Db(), level)
	if err != nil {
		return nil, err
	}
	return toDomainCategories(list), nil
}

func (r *categoryRepo) GetCategoriesByParentID(parentID *int32) ([]domain.Category, error) {
	list, err := model.GetCategoriesByParentID(r.db.Db(), parentID)
	if err != nil {
		return nil, err
	}
	return toDomainCategories(list), nil
}

func (r *categoryRepo) GetCategoryTree() ([]*domain.CategoryNode, error) {
	tree, err := model.GetCategoryTree(r.db.Db())
	if err != nil {
		return nil, err
	}
	return toDomainCategoryNodes(tree), nil
}

func (r *categoryRepo) UpdateCategory(category *domain.Category) error {
	return model.UpdateCategory(r.db.Db(), toGenCategory(category))
}

func (r *categoryRepo) MoveCategory(categoryID int32, newParentID *int32) error {
	return model.MoveCategory(r.db.Db(), categoryID, newParentID)
}
