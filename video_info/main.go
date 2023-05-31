package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/video_info/service"
	"github.com/gin-gonic/gin"
)

const port = "8084"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
	//info.GenerateVideoInfoConfig("minecraft.mp4")
}

func InitRoutes(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	playlist := service.NewVideoPlaylist()

	engine.GET("/video/list", service.VideoListController)
	engine.GET("/video/detailed_info", service.DetailedInfoController)
	engine.GET("/video/start", playlist.StartVideo)
	engine.GET("/video/current_duration", playlist.GetCurrentDuration)
}
