package get

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func Create(c *gin.Context) {
	c.HTML(http.StatusOK, "create.html", nil)
}

func SaveFile(c *gin.Context) {
}