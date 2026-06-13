package model

import (
	"fmt"

	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/errs"
	"gorm.io/gorm"
)

const MaxCategoryLevel = 3 // 最大层级：1=根, 2=子, 3=孙（0 为虚拟树根，无 DB 行）

// ValidateLevel 验证分类层级是否合法
func ValidateLevel(level int32) error {
	if level < 1 || level > MaxCategoryLevel {
		return fmt.Errorf("category level must be between 1 and %d", MaxCategoryLevel)
	}
	return nil
}

// GetFullPath 获取完整路径（包含自身）
func GetFullPath(category *gen.Category) string {
	if category.Path == nil || *category.Path == "/" {
		return fmt.Sprintf("/%s", category.Name)
	}
	return fmt.Sprintf("%s/%s", *category.Path, category.Name)
}

// --- CRUD ---

// CreateCategory 创建分类（自动计算level和path）
func CreateCategory(db *gorm.DB, category *gen.Category) error {
	if category.ParentID != nil && *category.ParentID <= 0 {
		return errs.ErrInvalidParentID
	}

	// 如果有父分类，计算level和path
	if category.ParentID != nil {
		parent, err := GetCategoryByID(db, *category.ParentID)
		if err != nil {
			return fmt.Errorf("parent category not found: %w", err)
		}
		newLevel := *parent.Level + 1
		if newLevel > MaxCategoryLevel {
			return fmt.Errorf("cannot create category: exceeds max level %d", MaxCategoryLevel)
		}
		category.Level = &newLevel
		parentPath := GetFullPath(parent)
		category.Path = &parentPath
	} else {
		rootLevel := int32(1)
		rootPath := "/"
		category.Level = &rootLevel
		category.Path = &rootPath
		category.ParentID = nil
	}

	return db.Create(category).Error
}

// GetCategoryByID 根据ID获取分类
func GetCategoryByID(db *gorm.DB, id int32) (*gen.Category, error) {
	var category gen.Category
	err := db.First(&category, id).Error
	return &category, err
}

// GetCategoryByName 根据名称和父ID获取分类（用于检查重复）
func GetCategoryByName(db *gorm.DB, name string, parentID *int32) (*gen.Category, error) {
	var category gen.Category
	query := db.Where("name = ?", name)
	if parentID != nil {
		query = query.Where("parent_id = ?", *parentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}
	err := query.First(&category).Error
	return &category, err
}

// ListCategories 查询分类列表（支持筛选条件）
func ListCategories(db *gorm.DB, conds ...interface{}) ([]gen.Category, error) {
	var categories []gen.Category
	err := db.Order("level ASC, sort_order DESC, name ASC").Find(&categories, conds...).Error
	return categories, err
}

// ListCategoriesByLevel 按层级查询分类
func ListCategoriesByLevel(db *gorm.DB, level int32) ([]gen.Category, error) {
	var categories []gen.Category
	err := db.Where("level = ?", level).Order("sort_order DESC, name ASC").Find(&categories).Error
	return categories, err
}

// GetCategoriesByParentID 获取指定父分类下的子分类
func GetCategoriesByParentID(db *gorm.DB, parentID *int32) ([]gen.Category, error) {
	var categories []gen.Category
	query := db
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Order("sort_order DESC, name ASC").Find(&categories).Error
	return categories, err
}

// UpdateCategory 更新分类（不允许更改层级结构）
func UpdateCategory(db *gorm.DB, category *gen.Category) error {
	// 获取原分类信息，防止更改层级
	existing, err := GetCategoryByID(db, *category.ID)
	if err != nil {
		return err
	}

	// 保留层级相关信息，只允许更新基本属性
	category.Level = existing.Level
	category.ParentID = existing.ParentID
	category.Path = existing.Path

	return db.Save(category).Error
}

// getAllDescendantIDs 获取分类及其所有后代的ID（使用ID-based递归查询）
func getAllDescendantIDs(db *gorm.DB, categoryID int32) ([]int32, error) {
	ids := []int32{categoryID}

	var children []gen.Category
	if err := db.Where("parent_id = ?", categoryID).Find(&children).Error; err != nil {
		return nil, err
	}

	for _, child := range children {
		childIDs, err := getAllDescendantIDs(db, *child.ID)
		if err != nil {
			return nil, err
		}
		ids = append(ids, childIDs...)
	}

	return ids, nil
}

// isDescendantOf 检查targetID是否是ancestorID的后代
func isDescendantOf(db *gorm.DB, targetID, ancestorID int32) (bool, error) {
	if targetID == ancestorID {
		return true, nil
	}

	var category gen.Category
	if err := db.First(&category, targetID).Error; err != nil {
		return false, err
	}

	if category.ParentID == nil || *category.ParentID == 0 {
		return false, nil
	}

	return isDescendantOf(db, *category.ParentID, ancestorID)
}

// MoveCategory 移动分类到新的父分类下
func MoveCategory(db *gorm.DB, categoryID int32, newParentID *int32) error {
	if newParentID != nil && *newParentID <= 0 {
		return errs.ErrInvalidParentID
	}

	return db.Transaction(func(tx *gorm.DB) error {
		category, err := GetCategoryByID(tx, categoryID)
		if err != nil {
			return err
		}

		newLevel := int32(1)
		newPath := "/"

		if newParentID != nil {
			if *newParentID == categoryID {
				return fmt.Errorf("cannot move category to itself")
			}

			isDescendant, err := isDescendantOf(tx, *newParentID, categoryID)
			if err != nil {
				return err
			}
			if isDescendant {
				return fmt.Errorf("cannot move category to its own child")
			}

			parent, err := GetCategoryByID(tx, *newParentID)
			if err != nil {
				return fmt.Errorf("parent category not found: %w", err)
			}

			newLevel = *parent.Level + 1
			if newLevel > MaxCategoryLevel {
				return fmt.Errorf("cannot move category: would exceed max level %d", MaxCategoryLevel)
			}
			parentPath := GetFullPath(parent)
			newPath = parentPath
		}

		oldPath := GetFullPath(category)
		oldLevel := *category.Level
		levelDelta := newLevel - oldLevel

		category.ParentID = newParentID
		category.Level = &newLevel
		category.Path = &newPath

		if err := tx.Save(category).Error; err != nil {
			return err
		}

		newFullPath := GetFullPath(category)
		return tx.Model(&gen.Category{}).Where("path LIKE ?", oldPath+"%").Updates(map[string]interface{}{
			"path":  gorm.Expr("REPLACE(path, ?, ?)", oldPath, newFullPath),
			"level": gorm.Expr("level + ?", levelDelta),
		}).Error
	})
}

// GetCategoryTree 获取分类树
func GetCategoryTree(db *gorm.DB) ([]*CategoryNode, error) {
	categories, err := ListCategories(db)
	if err != nil {
		return nil, err
	}

	return buildCategoryTree(categories), nil
}

// CategoryNode 分类树节点
type CategoryNode struct {
	Category gen.Category    `json:"category"`
	Children []*CategoryNode `json:"children,omitempty"`
}

// buildCategoryTree 构建分类树
func buildCategoryTree(categories []gen.Category) []*CategoryNode {
	nodeMap := make(map[int32]*CategoryNode)
	var roots []*CategoryNode

	// 第一遍：创建所有节点
	for i := range categories {
		node := &CategoryNode{
			Category: categories[i],
			Children: []*CategoryNode{},
		}
		nodeMap[*categories[i].ID] = node
	}

	// 第二遍：构建父子关系
	for i := range categories {
		node := nodeMap[*categories[i].ID]
		if categories[i].ParentID == nil || *categories[i].ParentID == 0 {
			roots = append(roots, node)
		} else {
			if parent, ok := nodeMap[*categories[i].ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return roots
}
