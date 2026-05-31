package router

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/pkg/errs"
	"github.com/gin-gonic/gin"
)

// RegisterTimeLogRoutes 注册 TimeLog 相关路由
func RegisterTimeLogRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.POST("/timelogs", createTimeLogHandler(deps))
	group.GET("/timelogs", listTimeLogsHandler(deps))
	group.GET("/timelogs/:id", getTimeLogHandler(deps))
	group.PUT("/timelogs/:id", updateTimeLogHandler(deps))
	group.DELETE("/timelogs/:id", deleteTimeLogHandler(deps))

	// Category 相关路由
	group.GET("/categories", listCategoriesHandler(deps))
	group.GET("/categories/tree", getCategoryTreeHandler(deps))
	group.POST("/categories", createCategoryHandler(deps))
	group.GET("/categories/:id", getCategoryHandler(deps))
	group.PUT("/categories/:id", updateCategoryHandler(deps))

	group.POST("/categories/:id/move", moveCategoryHandler(deps))
}

func createTimeLogHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req timelogv1.CreateTimelogRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		tl, err := mapper.TimelogFromCreateRequest(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if err := deps.Service.CreateTimeLog(tl); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		createdLog, err := deps.Service.GetTimeLogByID(*tl.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(createdLog), "Time log created successfully"))
	}
}

func listTimeLogsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		limitStr := c.Query("limit")
		orderBy := c.Query("order")

		var tls []gen.Timelog
		var err error

		if limitStr != "" || orderBy != "" {
			limit := 0
			if limitStr != "" {
				if l, parseErr := strconv.Atoi(limitStr); parseErr == nil && l > 0 {
					limit = l
				}
			}

			if orderBy == "" {
				orderBy = "created_at DESC"
			}

			tls, err = deps.Service.ListTimeLogsWithOptions(limit, orderBy)
		} else {
			tls, err = deps.Service.ListTimeLogs()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogsToProto(tls), "Time logs retrieved successfully"))
	}
}

func getTimeLogHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		tl, err := deps.Service.GetTimeLogByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(tl), "Time log retrieved successfully"))
	}
}

func updateTimeLogHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		existing, err := deps.Service.GetTimeLogByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
			return
		}
		var req timelogv1.UpdateTimelogRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if err := mapper.ApplyTimelogUpdate(existing, &req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		existing.ID = &id
		if err := deps.Service.UpdateTimeLog(existing); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		updatedLog, err := deps.Service.GetTimeLogByID(*existing.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(updatedLog), "Time log updated successfully"))
	}
}

func deleteTimeLogHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var id int32
		if err := parseInt32Param(c, "id", &id); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}
		if err := deps.Service.DeleteTimeLog(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil, "Time log deleted successfully"))
	}
}

// parseUintParam 辅助函数
func parseUintParam(c *gin.Context, key string, out *uint) error {
	idStr := c.Param(key)
	var id64 uint64
	var err error
	if id64, err = parseUint(idStr); err != nil {
		return err
	}
	*out = uint(id64)
	return nil
}

func parseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// parseInt32Param 辅助函数
func parseInt32Param(c *gin.Context, key string, out *int32) error {
	idStr := c.Param(key)
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		return err
	}
	*out = int32(id64)
	return nil
}

// --- Category Handlers ---

func listCategoriesHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		levelStr := c.Query("level")
		parentIDStr := c.Query("parent_id")

		var categories []gen.Category
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
		category.ID = &id
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
