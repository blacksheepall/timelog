package router

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/blacksheepaul/timelog/core/errs"
	"github.com/gin-gonic/gin"
)

// 添加分类相关路由
func setupCategoryRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/categories", listCategoriesHandler(deps))
	group.GET("/categories/tree", getCategoryTreeHandler(deps))
	group.POST("/categories", createCategoryHandler(deps))
	group.GET("/categories/:id", getCategoryHandler(deps))
	group.PUT("/categories/:id", updateCategoryHandler(deps))
	group.POST("/categories/:id/move", moveCategoryHandler(deps))
}

func listCategoriesHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		levelStr := c.Query("level")
		parentIDStr := c.Query("parent_id")

		var categories []domain.Category
		var err error

		if levelStr != "" {
			level, parseErr := strconv.ParseInt(levelStr, 10, 32)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "invalid level parameter"))
				return
			}
			categories, err = deps.Service.ListCategoriesByLevel(int32(level))
		} else if parentIDStr != "" {
			parentID, parseErr := strconv.ParseInt(parentIDStr, 10, 32)
			if parseErr != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "invalid parent_id parameter"))
				return
			}
			pid := int32(parentID)
			categories, err = deps.Service.GetCategoriesByParentID(&pid)
		} else {
			categories, err = deps.Service.ListCategories()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.CategoriesToProto(categories), "Categories retrieved successfully"))
	}
}

func getCategoryTreeHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		tree, err := deps.Service.GetCategoryTree()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryTreeToProto(tree), "Category tree retrieved successfully"))
	}
}

func createCategoryHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req timelogv1.CreateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		category, err := mapper.CategoryFromCreateRequest(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if err := deps.Service.CreateCategory(category); err != nil {
			c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category created successfully"))
	}
}

func getCategoryHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		category, err := deps.Service.GetCategoryByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category retrieved successfully"))
	}
}

func updateCategoryHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		category, err := deps.Service.GetCategoryByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		var req timelogv1.UpdateCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if err := mapper.ApplyCategoryUpdate(category, &req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		category.ID = id
		if err := deps.Service.UpdateCategory(category); err != nil {
			c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category updated successfully"))
	}
}

func moveCategoryHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		var req timelogv1.MoveCategoryRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := mapper.ValidateMoveCategoryRequest(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := deps.Service.MoveCategory(id, req.ParentId); err != nil {
			c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil, "Category moved successfully"))
	}
}

func categoryErrorStatus(err error) int {
	if errors.Is(err, errs.ErrInvalidParentID) {
		return http.StatusBadRequest
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "must not be provided"),
		strings.Contains(msg, "request is required"),
		strings.Contains(msg, "exceeds max level"),
		strings.Contains(msg, "cannot move category"),
		strings.Contains(msg, "parent category not found"):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
