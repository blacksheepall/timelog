package middleware

import (
	"strings"

	"github.com/blacksheepaul/timelog/core/errs"
	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/gin-gonic/gin"
)

func Auth(store ports.SessionTokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := GetSessionFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(401, authErrorResponse(err.Error()))
			return
		}

		if !isValidUserToken(store, session) {
			c.AbortWithStatusJSON(401, authErrorResponse("Invalid or expired token"))
			return
		}

		c.Next()
	}
}

// authErrorResponse mirrors router.ApiResponse for middleware-level errors.
// The middleware package cannot import router (import cycle); the shape is
// pinned by TestAuthMiddleware401Envelope and router.TestEnvelopeContract.
func authErrorResponse(message string) gin.H {
	return gin.H{"data": nil, "message": message, "status": 401}
}

func isValidUserToken(store ports.SessionTokenStore, token string) bool {
	if store == nil {
		return false
	}
	// Only accept keys with auth_token: prefix to prevent passkey session misuse
	if _, ok := store.GetCache("auth_token:" + token); ok {
		return true
	}
	return false
}

func GetSessionFromHeader(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", errs.ErrAuthNoSession
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", errs.ErrAuthInvalidSession
	}

	return parts[1], nil
}
