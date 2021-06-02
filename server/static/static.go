package static

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed *
var static embed.FS

func InitFS(router *gin.Engine) {
	t, _ := template.ParseFS(static, "templates/*.html")
	router.SetHTMLTemplate(t)
	router.StaticFS("/static", http.FS(static))
}
