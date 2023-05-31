package main

import (
	"github.com/achekslime/core/rest_api"
	"github.com/achekslime/flixray/room_api/service"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const port = "8082"

func main() {
	engine := gin.New()
	InitRoutes(engine)

	apiRunner := rest_api.NewService()
	apiRunner.ConfigureServer(engine, port)
	apiRunner.Run()
}

func InitRoutes(engine *gin.Engine) {
	roomService := service.NewRoomService()
	err := roomService.Configure("private_key")
	if err != nil {
		logrus.Errorf("error while NewRoomService(): %s", err)
	}

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	engine.POST("/api/create", roomService.CreateNewRoom)
	engine.GET("/api/admin_rooms", roomService.GetAdminRooms)
	engine.GET("/api/available_rooms", roomService.GetAvailableRooms)
	engine.POST("/api/add_user", roomService.AddUser)
	engine.DELETE("/api/del_room", roomService.DelRoom)
}
