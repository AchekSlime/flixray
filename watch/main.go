package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/watch/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const port = "8093"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	hub := service.NewHub()

	err := hub.Configure("private_key")
	if err != nil {
		logrus.Errorf("error while NewHub(): %s", err)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.GET("/api/hello", BindHelloMessage)
	engine.GET("/api/start", hub.Start)
	engine.GET("/api/users_count", hub.GetUsersCount)
	engine.GET("/api/current_info", hub.GetCurrentVideoInfo)

	engine.GET("/ws/:roomName", hub.ServeWs)
}

func BindHelloMessage(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"msg": "hello",
	})
}
