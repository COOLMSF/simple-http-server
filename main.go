package main

import (
	"./controller/get"
	"./controller/post"
	_ "./middleware"
	"./pkg"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {

	r := gin.Default()

	/*
	 * HTML render
	 */
	r.LoadHTMLGlob("view/*")
	r.GET("/", get.RenderIndex)
	r.GET("/view/login.html", get.RenderLogin)
	r.GET("/view/create.html", get.RenderCreate)
	r.GET("/view/upload.html", get.RenderUpload)

	r.POST("/create", post.Create)

	sess := r.Group("/", pkg.EnableCookieSession())
	{
		sess.POST("/login", post.Login)

		/*
		 * Got problem here, session not working...
		 */
		// authorized := sess.Group("/auth", pkg.AuthSessionMiddle())
		authorized := sess.Group("/auth")
		{
			authorized.POST("/upload", post.Upload)
			authorized.StaticFS("/download", http.Dir("download"))
		}
	}

	// err := r.Run(":3333")
	err := r.RunTLS(":3333", "server.crt", "server.key")
	if err != nil {
		log.Fatal("r.Run")
	}
}
