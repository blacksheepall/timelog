package mapper

import (
	"fmt"

	"github.com/blacksheepaul/timelog/core/errs"
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/domain"
)

func CategoryToProto(c *domain.Category) *timelogv1.Category {
	if c == nil {
		return nil
	}
	return &timelogv1.Category{
		Id:          c.ID,
		Name:        c.Name,
		Color:       c.Color,
		Description: c.Description,
		ParentId:    c.ParentID,
		Level:       c.Level,
		SortOrder:   c.SortOrder,
		Path:        c.Path,
		CreatedAt:   FormatTimeUTC(c.CreatedAt),
		UpdatedAt:   FormatTimeUTC(c.UpdatedAt),
	}
}

func CategoriesToProto(categories []domain.Category) []*timelogv1.Category {
	out := make([]*timelogv1.Category, 0, len(categories))
	for i := range categories {
		out = append(out, CategoryToProto(&categories[i]))
	}
	return out
}

func CategoryTreeNodeToProto(node *domain.CategoryNode) *timelogv1.CategoryTreeNode {
	if node == nil {
		return nil
	}
	out := &timelogv1.CategoryTreeNode{
		Category: CategoryToProto(&node.Category),
		Children: make([]*timelogv1.CategoryTreeNode, 0, len(node.Children)),
	}
	for _, child := range node.Children {
		out.Children = append(out.Children, CategoryTreeNodeToProto(child))
	}
	return out
}

func CategoryTreeToProto(nodes []*domain.CategoryNode) []*timelogv1.CategoryTreeNode {
	out := make([]*timelogv1.CategoryTreeNode, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, CategoryTreeNodeToProto(node))
	}
	return out
}

func CategoryFromCreateRequest(req *timelogv1.CreateCategoryRequest) (*domain.Category, error) {
	if req == nil {
		return nil, fmt.Errorf("request is required")
	}
	if req.Level != nil {
		return nil, fmt.Errorf("level is computed by the server and must not be provided")
	}
	if req.SortOrder != nil {
		return nil, fmt.Errorf("sort_order is computed by the server and must not be provided")
	}
	if req.ParentId != nil && *req.ParentId <= 0 {
		return nil, errs.ErrInvalidParentID
	}
	return &domain.Category{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
		ParentID:    req.ParentId,
	}, nil
}

func ValidateMoveCategoryRequest(req *timelogv1.MoveCategoryRequest) error {
	if req == nil {
		return fmt.Errorf("request is required")
	}
	if req.ParentId != nil && *req.ParentId <= 0 {
		return errs.ErrInvalidParentID
	}
	return nil
}

func ApplyCategoryUpdate(category *domain.Category, req *timelogv1.UpdateCategoryRequest) error {
	if category == nil || req == nil {
		return nil
	}
	if req.ParentId != nil {
		return errs.ErrInvalidParentID
	}
	if req.Level != nil {
		return fmt.Errorf("level is computed by the server and must not be provided")
	}
	if req.SortOrder != nil {
		return fmt.Errorf("sort_order is computed by the server and must not be provided")
	}
	if req.Name != nil {
		category.Name = req.GetName()
	}
	if req.Color != nil {
		category.Color = req.GetColor()
	}
	if req.Description != nil {
		category.Description = req.GetDescription()
	}
	return nil
}
