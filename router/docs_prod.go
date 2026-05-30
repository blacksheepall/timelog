//go:build prod

package router

import "github.com/gin-gonic/gin"

func setupDocs(r *gin.Engine) {
	// API docs disabled in production build
}
