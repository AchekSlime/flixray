package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/minio/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const port = "8071"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	minioService := service.NewMinioService()
	err := minioService.Configure("private_key")
	if err != nil {
		logrus.Errorf("error while NewMinioService(): %s", err)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.GET("/api/:bucket/:file", minioService.GetFile)
}
