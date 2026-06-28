package router

import (
	"net/http"
	"strconv"

	"github.com/blacksheepaul/timelog/core/audit"
	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/internal/domain"
	"github.com/gin-gonic/gin"
)

// RegisterTimeLogRoutes 注册 TimeLog 相关路由
func RegisterTimeLogRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.POST("/timelogs", createTimeLogHandler(deps))
	group.GET("/timelogs", listTimeLogsHandler(deps))
	group.GET("/timelogs/:id", getTimeLogHandler(deps))
	group.PUT("/timelogs/:id", updateTimeLogHandler(deps))
	group.DELETE("/timelogs/:id", deleteTimeLogHandler(deps))
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
		ctx := audit.WithSource(c.Request.Context(), audit.SourceHuman)
		if err := deps.Service.CreateTimeLog(ctx, tl); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		createdLog, err := deps.Service.GetTimeLogByID(tl.ID)
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

		var tls []domain.TimeLog
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
		existing.ID = id
		ctx := audit.WithSource(c.Request.Context(), audit.SourceHuman)
		if err := deps.Service.UpdateTimeLog(ctx, existing); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		updatedLog, err := deps.Service.GetTimeLogByID(existing.ID)
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
		ctx := audit.WithSource(c.Request.Context(), audit.SourceHuman)
		if err := deps.Service.DeleteTimeLog(ctx, id); err != nil {
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
