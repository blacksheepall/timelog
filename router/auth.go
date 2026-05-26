package router

import (
	"net/http"

	"github.com/blacksheepaul/timelog/service"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(api *gin.RouterGroup) {
	api.POST("/auth/dev-login", devLoginHandler)
}

func devLoginHandler(c *gin.Context) {
	token, err := service.GenerateSessionToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	// 7 days TTL for dev convenience
	const devTokenTTL int64 = 7 * 24 * 3600
	if err := service.StoreSessionToken(token, devTokenTTL); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(
		gin.H{"token": token, "token_type": "Bearer", "expires_in": devTokenTTL},
		"dev login success",
	))
}
