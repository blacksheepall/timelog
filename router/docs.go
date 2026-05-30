//go:build !prod

package router

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed docs/openapi.yaml
//go:embed docs/redoc.html
var docsFS embed.FS

func setupDocs(r *gin.Engine) {
	docs, err := fs.Sub(docsFS, "docs")
	if err != nil {
		panic("docs filesystem: " + err.Error())
	}

	r.StaticFS("/docs", http.FS(docs))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/docs/redoc.html")
	})
}
