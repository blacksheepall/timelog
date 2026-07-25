package router

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupDatasourceRoutes(group *gin.RouterGroup, deps Dependencies) {
	group.POST("/datasources/:name/sync", syncDatasourceHandler(deps))
}

func syncDatasourceHandler(deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "datasource name is required"))
			return
		}

		result, err := deps.Service.SyncMetrics(c.Request.Context(), name)
		if err != nil {
			// Distinguish "not found" from external API failures.
			if isDatasourceNotFound(err) {
				c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, err.Error()))
				return
			}
			c.JSON(http.StatusBadGateway, ErrorResponse(http.StatusBadGateway, err.Error()))
			return
		}

		msg := "Sync completed"
		if result.Failed > 0 {
			msg = "Sync completed with failures"
		}
		c.JSON(http.StatusOK, SuccessResponse(result, msg))
	}
}

func isDatasourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	// The registry returns errors containing "not found" / "not registered".
	return strings.Contains(err.Error(), "not found") || strings.Contains(err.Error(), "not registered")
}
