package get

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RenderIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func RenderLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func RenderCreate(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func RenderUpload(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", nil)
}
