package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/user_api/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const port = "8083"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	roomService := service.NewUserService()
	err := roomService.Configure("private_key")
	if err != nil {
		logrus.Errorf("error while NewRoomService(): %s", err)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.GET("/api/find_by_name", roomService.FindByName)
	engine.DELETE("/api/clear_all", roomService.ClearAll)
}
