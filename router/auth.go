package router

import (
	"net/http"

	"github.com/blacksheepaul/timelog/core/config"
	"github.com/blacksheepaul/timelog/internal/api/mapper"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(api *gin.RouterGroup, cfg *config.Config, deps Dependencies) {
	api.POST("/auth/dev-login", devLoginHandler(cfg, deps))
}

func devLoginHandler(cfg *config.Config, deps Dependencies) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !cfg.DevMode {
			c.JSON(http.StatusNotFound, ErrorResponse(http.StatusNotFound, "not found"))
			return
		}

		token, err := deps.Service.GenerateSessionToken()
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		// 7 days TTL for dev convenience
		const devTokenTTL int64 = 7 * 24 * 3600
		if err := deps.Service.StoreSessionToken(token, devTokenTTL); err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
			return
		}

		c.JSON(http.StatusOK, SuccessResponse(
			mapper.LoginResponse(token, "Bearer", devTokenTTL),
			"dev login success",
		))
	}
}
