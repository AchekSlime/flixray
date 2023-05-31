package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/gin-gonic/gin"
)

const port = "8080"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.GET("/hello", BindHelloMessage)
}

func BindHelloMessage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg": "hello",
	})
}
