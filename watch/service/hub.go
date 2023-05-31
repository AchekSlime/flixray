package service

import (
	"errors"
	"fmt"
	"github.com/achekslime/core/jwt"
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage"
	"github.com/achekslime/core/storage/postgres"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Hub struct {
	// room sessions.
	rooms map[string]*room

	storage    *storage.Storage
	jwtService *jwt.JwtService
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]*room),
	}
}

func (hub *Hub) Configure(jwtKey string) error {
	// storage configuration.
	newStorage, err := storage.NewStorage()
	if err != nil {
		return err
	}
	hub.storage = newStorage

	// jwt configuration.
	hub.jwtService = jwt.ConfigureJWT(jwtKey)
	return nil
}

func (hub *Hub) Authorize(context *gin.Context) *models.User {
	userClaims, err := jwt.GetTokenClaims(context, hub.jwtService)
	if err != nil {
		rest_api_utils.BindUnauthorized(context, err)
		return nil
	}

	// get user from db
	dbUser, err := hub.storage.UserStorage.GetUserByEmail(userClaims.Email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			rest_api_utils.BindUnauthorized(context, errors.New("user not found"))
			return nil
		}
		rest_api_utils.BindInternalError(context, err)
		return nil
	}

	return dbUser
}

func contains(s []models.Room, roomName string) bool {
	for _, v := range s {
		if v.Name == roomName {
			return true
		}
	}

	return false
}

// ServeWs handles websocket requests from the peer.
func (hub *Hub) ServeWs(ctx *gin.Context) {
	// Достаем room name.
	roomName := ctx.Param("roomName")

	// авторизуем.
	user := hub.Authorize(ctx)
	if user == nil {
		return
	}

	// достаем доступные комнаты.
	roomList, err := hub.storage.RoomStorage.GetAvailableRooms(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindUnprocessableEntity(ctx, errors.New("available rooms not found"))
			return
		}
		rest_api_utils.BindInternalError(ctx, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	// если комнаты нет в списке комнат, доступных пользователю.
	if !contains(roomList, roomName) {
		rest_api_utils.BindUnauthorized(ctx, fmt.Errorf("user doesnt have access for specified room:%s", roomName))
		return
	}

	// достаем комнату.
	roomDb, err := hub.storage.RoomStorage.GetRoomByName(roomName)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindUnprocessableEntity(ctx, errors.New("room not found"))
			return
		}
		rest_api_utils.BindInternalError(ctx, err)
		return
	}

	// проверяем создана ли сессия.
	if _, ok := hub.rooms[roomName]; !ok {
		rest_api_utils.BindUnprocessableEntity(ctx, errors.New("room session doesnt start"))
		return
	}

	// переводим на WS.
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// достаем комнату.
	currentRoom := hub.rooms[roomName]

	// создаем новое соединение с пользователем.
	clientConnection := newConnection(ws, user, roomDb)

	// добавляем в подписки комнаты новое соединение.
	currentRoom.addSub <- clientConnection

	// запускаем чтение/запись на соединении.
	go clientConnection.writeHandler()
	go clientConnection.readHandler(currentRoom.broadcast, currentRoom.delSub)
}
