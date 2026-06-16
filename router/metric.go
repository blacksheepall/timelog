package router

import (
	"net/http"
	"strconv"

	timelogv1 "github.com/blacksheepaul/timelog/gen/go/timelog/v1"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/gin-gonic/gin"
)

func setupMetricRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.GET("/metrics", listMetricsHandler(deps))
	group.POST("/metrics", createMetricHandler(deps))
	group.GET("/metrics/:id", getMetricHandler(deps))
	group.PUT("/metrics/:id", updateMetricHandler(deps))
	group.DELETE("/metrics/:id", deleteMetricHandler(deps))
	group.GET("/metrics/:id/records", listMetricRecordsHandler(deps))
}

func createMetricHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req timelogv1.CreateMetricRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		metric := mapper.MetricFromCreateRequest(&req)
		if err := deps.Service.CreateMetric(metric); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if created, err := deps.Service.GetMetricByID(metric.ID); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.MetricToProto(created), "Metric created successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.MetricToProto(metric), "Metric created successfully"))
		}
	}
}

func listMetricsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		metrics, err := deps.Service.ListMetrics()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.MetricsToProto(metrics), "Metrics retrieved successfully"))
	}
}

func getMetricHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := parseMetricIDParam(c)
		if id == 0 {
			return
		}
		metric, err := deps.Service.GetMetricByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Metric not found"))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.MetricToProto(metric), "Metric retrieved successfully"))
	}
}

func updateMetricHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := parseMetricIDParam(c)
		if id == 0 {
			return
		}

		existing, err := deps.Service.GetMetricByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Metric not found"))
			return
		}

		var req timelogv1.UpdateMetricRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
			return
		}

		mapper.ApplyMetricUpdate(existing, &req)
		if err := deps.Service.UpdateMetric(existing); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		if updated, err := deps.Service.GetMetricByID(id); err == nil {
			c.JSON(http.StatusOK, SuccessResponse(mapper.MetricToProto(updated), "Metric updated successfully"))
		} else {
			c.JSON(http.StatusOK, SuccessResponse(mapper.MetricToProto(existing), "Metric updated successfully"))
		}
	}
}

func deleteMetricHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := parseMetricIDParam(c)
		if id == 0 {
			return
		}
		if _, err := deps.Service.GetMetricByID(id); err != nil {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "Metric not found"))
			return
		}
		if err := deps.Service.DeleteMetric(id); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(nil, "Metric deleted successfully"))
	}
}

func listMetricRecordsHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := parseMetricIDParam(c)
		if id == 0 {
			return
		}
		records, err := deps.Service.ListMetricRecords(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(mapper.MetricRecordsToProto(records), "Metric records retrieved successfully"))
	}
}

func parseMetricIDParam(c *gin.Context) int32 {
	idStr := c.Param("id")
	id64, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "Invalid metric ID"))
		return 0
	}
	return int32(id64)
}
