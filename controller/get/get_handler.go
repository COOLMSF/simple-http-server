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

func RenderGif(c *gin.Context) {
	c.HTML(http.StatusOK, "gif.html", nil)
}

func RenderWebshell(c *gin.Context) {
	c.HTML(http.StatusOK, "webshell.html", nil)
}

func RenderPlayground(c *gin.Context) {
	c.HTML(http.StatusOK, "playground.html", nil)
}

