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
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// RegisterTimeLogRoutes 注册 TimeLog 相关路由
func RegisterTimeLogRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.POST("/timelogs", createTimeLogHandler)
	group.GET("/timelogs", listTimeLogsHandler)
	group.GET("/timelogs/:id", getTimeLogHandler)
	group.PUT("/timelogs/:id", updateTimeLogHandler)
	group.DELETE("/timelogs/:id", deleteTimeLogHandler)

	// Category 相关路由
	group.GET("/categories", listCategoriesHandler)
	group.GET("/categories/tree", getCategoryTreeHandler)
	group.POST("/categories", createCategoryHandler)
	group.GET("/categories/:id", getCategoryHandler)
	group.PUT("/categories/:id", updateCategoryHandler)

	group.POST("/categories/:id/move", moveCategoryHandler)
}

// CreateTimeLogHandler godoc
// @Summary 创建时间日志
// @Description 新增一条时间日志
// @Tags timelog
// @Accept json
// @Produce json
// @Param data body gen.Timelog true "时间日志数据"
// @Success 200 {object} gen.Timelog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs [post]
func createTimeLogHandler(c *gin.Context) {
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
	if err := service.CreateTimeLog(tl); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整信息
	createdLog, err := service.GetTimeLogByID(*tl.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(createdLog), "Time log created successfully"))
}

// ListTimeLogsHandler godoc
// @Summary 查询时间日志列表
// @Description 获取所有时间日志
// @Tags timelog
// @Produce json
// @Param limit query int false "Limit number of results"
// @Param order query string false "Order by field (default: created_at DESC)"
// @Success 200 {array} gen.Timelog
// @Failure 500 {object} map[string]string
// @Router /api/timelogs [get]
func listTimeLogsHandler(c *gin.Context) {
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

		tls, err = service.ListTimeLogsWithOptions(limit, orderBy)
	} else {
		tls, err = service.ListTimeLogs()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogsToProto(tls), "Time logs retrieved successfully"))
}

// GetTimeLogHandler godoc
// @Summary 查询单条时间日志
// @Description 根据ID获取时间日志
// @Tags timelog
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} gen.Timelog
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/timelogs/{id} [get]
func getTimeLogHandler(c *gin.Context) {
	var id int32
	if err := parseInt32Param(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	tl, err := service.GetTimeLogByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(tl), "Time log retrieved successfully"))
}

// UpdateTimeLogHandler godoc
// @Summary 更新时间日志
// @Description 根据ID更新时间日志
// @Tags timelog
// @Accept json
// @Produce json
// @Param id path int true "日志ID"
// @Param data body gen.Timelog true "时间日志数据"
// @Success 200 {object} gen.Timelog
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs/{id} [put]
func updateTimeLogHandler(c *gin.Context) {
	var id int32
	if err := parseInt32Param(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	existing, err := service.GetTimeLogByID(id)
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
	if err := service.UpdateTimeLog(existing); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整信息
	updatedLog, err := service.GetTimeLogByID(*existing.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(mapper.TimelogToProto(updatedLog), "Time log updated successfully"))
}

// DeleteTimeLogHandler godoc
// @Summary 删除时间日志
// @Description 根据ID删除时间日志
// @Tags timelog
// @Produce json
// @Param id path int true "日志ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/timelogs/{id} [delete]
func deleteTimeLogHandler(c *gin.Context) {
	var id int32
	if err := parseInt32Param(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	if err := service.DeleteTimeLog(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil, "Time log deleted successfully"))
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

// listCategoriesHandler godoc
// @Summary 查询分类列表
// @Description 获取所有分类（扁平列表）
// @Tags category
// @Produce json
// @Param level query int false "Filter by level (0, 1, 2)"
// @Param parent_id query int false "Filter by parent_id"
// @Success 200 {array} gen.Category
// @Failure 500 {object} map[string]string
// @Router /api/categories [get]
func listCategoriesHandler(c *gin.Context) {
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
		categories, err = service.ListCategoriesByLevel(int32(level))
	} else if parentIDStr != "" {
		parentID, parseErr := strconv.ParseInt(parentIDStr, 10, 32)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "invalid parent_id parameter"))
			return
		}
		pid := int32(parentID)
		categories, err = service.GetCategoriesByParentID(&pid)
	} else {
		categories, err = service.ListCategories()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.CategoriesToProto(categories), "Categories retrieved successfully"))
}

// getCategoryTreeHandler godoc
// @Summary 获取分类树
// @Description 获取树形结构的分类列表
// @Tags category
// @Produce json
// @Success 200 {array} model.CategoryNode
// @Failure 500 {object} map[string]string
// @Router /api/categories/tree [get]
func getCategoryTreeHandler(c *gin.Context) {
	tree, err := service.GetCategoryTree()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryTreeToProto(tree), "Category tree retrieved successfully"))
}

// createCategoryHandler godoc
// @Summary 创建分类
// @Description 新增一个分类（支持层级，最大深度3层）
// @Tags category
// @Accept json
// @Produce json
// @Param data body gen.Category true "分类数据"
// @Success 200 {object} gen.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories [post]
func createCategoryHandler(c *gin.Context) {
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
	if err := service.CreateCategory(category); err != nil {
		c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category created successfully"))
}

// getCategoryHandler godoc
// @Summary 查询单个分类
// @Description 根据ID获取分类
// @Tags category
// @Produce json
// @Param id path int true "分类ID"
// @Success 200 {object} gen.Category
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
func getCategoryHandler(c *gin.Context) {
	var id int32
	if err := parseInt32Param(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	category, err := service.GetCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category retrieved successfully"))
}

// updateCategoryHandler godoc
// @Summary 更新分类
// @Description 根据ID更新分类（不允许修改层级结构）
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param data body gen.Category true "分类数据"
// @Success 200 {object} gen.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories/{id} [put]
func updateCategoryHandler(c *gin.Context) {
	var id int32
	if err := parseInt32Param(c, "id", &id); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}
	category, err := service.GetCategoryByID(id)
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
	if err := service.UpdateCategory(category); err != nil {
		c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(mapper.CategoryToProto(category), "Category updated successfully"))
}

// moveCategoryHandler godoc
// @Summary 移动分类
// @Description 将分类移动到新的父分类下
// @Tags category
// @Accept json
// @Produce json
// @Param id path int true "分类ID"
// @Param data body object true "移动参数" {"parent_id": 0}
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/categories/{id}/move [post]
func moveCategoryHandler(c *gin.Context) {
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

	if err := service.MoveCategory(id, req.ParentId); err != nil {
		c.JSON(categoryErrorStatus(err), ErrorResponse(categoryErrorStatus(err), err.Error()))
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil, "Category moved successfully"))
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
