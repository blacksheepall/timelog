package middleware

import (
	"strings"

	"github.com/blacksheepaul/timelog/internal/ports"
	"github.com/blacksheepaul/timelog/pkg/errs"
	"github.com/gin-gonic/gin"
)

func Auth(store ports.SessionTokenStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := GetSessionFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"msg": err.Error()})
			return
		}

		if !isValidUserToken(store, session) {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Invalid or expired token",
			})
			return
		}

		c.Next()
	}
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
