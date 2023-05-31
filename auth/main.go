package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/auth/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const port = "8081"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	authService := service.NewAuthService()
	err := authService.Configure("private_key")
	if err != nil {
		logrus.Errorf("error while NewAuthService(): %s", err)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.POST("/api/reg_mail", authService.RegMail)
	engine.POST("/api/confirm", authService.Confirm)
	engine.POST("/api/reg_user", authService.RegNewUser)
	engine.POST("/api/token", authService.GenerateToken)
}
