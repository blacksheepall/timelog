package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

// 添加约束相关路由
func setupConstraintRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/constraints", listConstraintsHandler)
	group.POST("/constraints", createConstraintHandler)
	group.GET("/constraints/:id", getConstraintHandler)
	group.PUT("/constraints/:id", updateConstraintHandler)
	group.DELETE("/constraints/:id", deleteConstraintHandler)
	group.POST("/constraints/:id/complete", completeConstraintHandler)
	group.POST("/constraints/:id/reactivate", reactivateConstraintHandler)
}

// CreateConstraintHandler godoc
// @Summary 创建约束
// @Description 新增一项约束
// @Tags constraint
// @Accept json
// @Produce json
// @Param data body gen.Constraint true "约束数据"
// @Success 200 {object} gen.Constraint
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints [post]
func createConstraintHandler(c *gin.Context) {
	// 定义一个临时结构体来处理日期字符串
	var request struct {
		Description     string `json:"description" binding:"required"`
		EndReason       string `json:"end_reason"`
		PunishmentQuote string `json:"punishment_quote" binding:"required"`
		StartDate       string `json:"start_date" binding:"required"`
		EndDate         string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// 解析日期
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid start_date format, expected YYYY-MM-DD"))
		return
	}

	var endDate *time.Time
	if request.EndDate != "" {
		parsedEndDate, err := time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid end_date format, expected YYYY-MM-DD"))
			return
		}
		endDate = &parsedEndDate
	}

	// 创建约束
	trueValue := true
	constraint := &gen.Constraint{
		Description:     request.Description,
		EndReason:       &request.EndReason,
		PunishmentQuote: request.PunishmentQuote,
		StartDate:       startDate,
		EndDate:         endDate,
		IsActive:        &trueValue,
	}

	if err := service.CreateConstraint(constraint); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整信息
	if createdConstraint, err := service.GetConstraintByID(*constraint.ID); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(createdConstraint, "Constraint created successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(constraint, "Constraint created successfully"))
	}
}

// ListConstraintsHandler godoc
// @Summary 获取约束列表
// @Description 获取所有约束，支持按活跃状态过滤
// @Tags constraint
// @Produce json
// @Param active query bool false "是否只显示活跃约束"
// @Success 200 {array} gen.Constraint
// @Failure 500 {object} map[string]string
// @Router /api/constraints [get]
func listConstraintsHandler(c *gin.Context) {
	activeStr := c.Query("active")

	var constraints []gen.Constraint
	var err error

	if activeStr == "true" {
		constraints, err = service.GetActiveConstraints()
	} else {
		constraints, err = service.GetAllConstraints()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(constraints, "Constraints retrieved successfully"))
}

// GetConstraintHandler godoc
// @Summary 获取单个约束
// @Description 根据ID获取约束详情
// @Tags constraint
// @Produce json
// @Param id path int true "约束ID"
// @Success 200 {object} gen.Constraint
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints/{id} [get]
func getConstraintHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
		return
	}
	id := int32(id64)

	constraint, err := service.GetConstraintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(constraint, "Constraint retrieved successfully"))
}

// UpdateConstraintHandler godoc
// @Summary 更新约束
// @Description 更新约束信息
// @Tags constraint
// @Accept json
// @Produce json
// @Param id path int true "约束ID"
// @Param data body gen.Constraint true "约束数据"
// @Success 200 {object} gen.Constraint
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints/{id} [put]
func updateConstraintHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
		return
	}
	id := int32(id64)

	// 先检查约束是否存在
	existingConstraint, err := service.GetConstraintByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
		return
	}

	// 定义临时结构体处理日期字符串
	var request struct {
		Description     string `json:"description"`
		EndReason       string `json:"end_reason"`
		PunishmentQuote string `json:"punishment_quote"`
		StartDate       string `json:"start_date"`
		EndDate         string `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	// 解析日期
	if request.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", request.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid start_date format, expected YYYY-MM-DD"))
			return
		}
		existingConstraint.StartDate = startDate
	}

	if request.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", request.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid end_date format, expected YYYY-MM-DD"))
			return
		}
		existingConstraint.EndDate = &endDate
	} else if request.EndDate == "" {
		// 如果发送了空字符串，清空结束日期
		existingConstraint.EndDate = nil
	}

	// 更新其他字段
	if request.Description != "" {
		existingConstraint.Description = request.Description
	}
	if request.PunishmentQuote != "" {
		existingConstraint.PunishmentQuote = request.PunishmentQuote
	}
	if request.EndReason != "" {
		existingConstraint.EndReason = &request.EndReason
	}

	if err := service.UpdateConstraint(existingConstraint); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 重新查询以获取完整信息
	if updatedConstraint, err := service.GetConstraintByID(id); err == nil {
		c.JSON(http.StatusOK, SuccessResponse(updatedConstraint, "Constraint updated successfully"))
	} else {
		c.JSON(http.StatusOK, SuccessResponse(existingConstraint, "Constraint updated successfully"))
	}
}

// DeleteConstraintHandler godoc
// @Summary 删除约束
// @Description 删除指定约束
// @Tags constraint
// @Produce json
// @Param id path int true "约束ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints/{id} [delete]
func deleteConstraintHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
		return
	}
	id := int32(id64)

	// 先检查约束是否存在
	if _, err := service.GetConstraintByID(id); err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
		return
	}

	if err := service.DeleteConstraint(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint deleted successfully"))
}

// CompleteConstraintHandler godoc
// @Summary 标记约束为完成
// @Description 将约束标记为完成状态，记录结束理由
// @Tags constraint
// @Accept json
// @Produce json
// @Param id path int true "约束ID"
// @Param data body map[string]string true "结束理由"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints/{id}/complete [post]
func completeConstraintHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
		return
	}
	id := int32(id64)

	var requestData struct {
		EndReason string `json:"end_reason"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := service.MarkConstraintAsCompleted(id, requestData.EndReason); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint marked as completed"))
}

// ReactivateConstraintHandler godoc
// @Summary 重新激活约束
// @Description 将约束重新激活
// @Tags constraint
// @Produce json
// @Param id path int true "约束ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/constraints/{id}/reactivate [post]
func reactivateConstraintHandler(c *gin.Context) {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
		return
	}
	id := int32(id64)

	if err := service.MarkConstraintAsActive(id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint reactivated successfully"))
}
