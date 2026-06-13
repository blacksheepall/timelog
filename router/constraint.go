package router

import (
	"net/http"
	"strconv"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/blacksheepaul/timelog/model/gen"
	"github.com/gin-gonic/gin"
)

// 添加约束相关路由
func setupConstraintRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/constraints", listConstraintsHandler(deps))
	group.POST("/constraints", createConstraintHandler(deps))
	group.GET("/constraints/:id", getConstraintHandler(deps))
	group.PUT("/constraints/:id", updateConstraintHandler(deps))
	group.DELETE("/constraints/:id", deleteConstraintHandler(deps))
	group.POST("/constraints/:id/complete", completeConstraintHandler(deps))
	group.POST("/constraints/:id/reactivate", reactivateConstraintHandler(deps))
	group.GET("/constraints/:id/evaluation", evaluateConstraintHandler(deps))
}

func createConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request timelogv1.CreateConstraintRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		constraint, err := mapper.ConstraintFromCreateRequest(&request)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := deps.Service.CreateConstraint(constraint); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		if createdConstraint, err := deps.Service.GetConstraintByID(*constraint.ID); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintToProto(createdConstraint), "Constraint created successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintToProto(constraint), "Constraint created successfully"))
		}
	}
}

func listConstraintsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		activeStr := c.Query("active")

		var constraints []gen.Constraint
		var err error

		if activeStr == "true" {
			constraints, err = deps.Service.GetActiveConstraints()
		} else {
			constraints, err = deps.Service.GetAllConstraints()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintsToProto(constraints), "Constraints retrieved successfully"))
	}
}

func getConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		constraint, err := deps.Service.GetConstraintByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintToProto(constraint), "Constraint retrieved successfully"))
	}
}

func updateConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		// 先检查约束是否存在
		existingConstraint, err := deps.Service.GetConstraintByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
			return
		}

		var request timelogv1.UpdateConstraintRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := mapper.ApplyConstraintUpdate(existingConstraint, &request); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := deps.Service.UpdateConstraint(existingConstraint); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 重新查询以获取完整信息
		if updatedConstraint, err := deps.Service.GetConstraintByID(id); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintToProto(updatedConstraint), "Constraint updated successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintToProto(existingConstraint), "Constraint updated successfully"))
		}
	}
}

func deleteConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		// 先检查约束是否存在
		if _, err := deps.Service.GetConstraintByID(id); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Constraint not found"))
			return
		}

		if err := deps.Service.DeleteConstraint(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint deleted successfully"))
	}
}

func completeConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		var requestData timelogv1.CompleteConstraintRequest
		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		if err := deps.Service.MarkConstraintAsCompleted(id, requestData.EndReason); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint marked as completed"))
	}
}

func reactivateConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		if err := deps.Service.MarkConstraintAsActive(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(nil, "Constraint reactivated successfully"))
	}
}

func evaluateConstraintHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id64, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid constraint ID"))
			return
		}
		id := int32(id64)

		eval, err := deps.Service.EvaluateConstraint(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(mapper.ConstraintEvaluationToProto(eval), "Constraint evaluated successfully"))
	}
}
