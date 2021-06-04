package main

import (
	"./controller/get"
	"./controller/post"
	"./pkg"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)


func main() {
	r := engine()
	r.Use(gin.Logger())
	if err := engine().RunTLS(":3333", "server.crt", "server.key"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()

	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	r.LoadHTMLGlob("view/*")

	r.GET("/", get.RenderIndex)
	r.GET("/view/login.html", get.RenderLogin)
	r.GET("/view/create.html", get.RenderCreate)
	r.GET("/view/upload.html", get.RenderUpload)
	r.GET("/view/gif.html", get.RenderGif)
	r.GET("/view/webshell.html", get.RenderWebshell)
	r.GET("/view/playground.html", get.RenderPlayground)
	r.GET("/logout", pkg.Logout)

	r.POST("/login", post.Login)
	r.POST("/playground", post.Playground)
	r.POST("/gif", post.Gif)
	r.POST("/create", post.Create)

	private := r.Group("/private")
	// private.Use(AuthRequired)
	private.Use(pkg.AuthSessionMiddle)
	{
		private.GET("/me", pkg.Me)
		private.GET("/status", pkg.Status)
		private.POST("/upload", post.Upload)
		private.POST("/webshell", post.WebShell)
		private.StaticFS("/download", http.Dir("download"))
	}
	return r
}
