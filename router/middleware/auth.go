package middleware

import (
	"strings"

	"github.com/blacksheepaul/timelog/pkg/errs"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := GetSessionFromHeader(c)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"msg": err.Error()})
			return
		}

		if !isValidUserToken(session) {
			c.AbortWithStatusJSON(401, gin.H{
				"msg": "Invalid or expired token",
			})
			return
		}

		c.Next()
	}
}

var tokenValidator = isValidUserTokenWithPrefix

func isValidUserToken(token string) bool {
	return tokenValidator(token)
}

func isValidUserTokenWithPrefix(token string) bool {
	dao := getMiddlewareDAO()
	// Only accept keys with auth_token: prefix to prevent passkey session misuse
	if _, ok := dao.GetCache("auth_token:" + token); ok {
		return true
	}
	return false
}

type cacheReader interface {
	GetCache(key string) (any, bool)
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
