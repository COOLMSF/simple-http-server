package main

import (
	"./handler/get"
	"./handler/post"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", get.Index)
	r.GET("/login", get.Login)
	r.GET("/create", get.Create)
	r.POST("/login", post.Login)
	r.POST("/create", post.Create)
	r.Static("file", "saved")
	r.GET("/upload", get.SaveFile)
	r.POST("/upload", post.SaveFile)

	err := r.Run(":3333")
	if err != nil {
		log.Fatal("r.Run")
	}
}
