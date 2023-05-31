package service

import (
	"errors"
	"fmt"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage/postgres"
	"github.com/achekslime/flixray/watch/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

func (hub *Hub) Start(ctx *gin.Context) {
	// Достаем room name.
	roomName := ctx.Query("room_name")
	if roomName == "" {
		rest_api_utils.BindBadRequest(ctx, fmt.Errorf("invalid room_name querry"))
		return
	}

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

	// получаем комнату.
	var currentRoom *room
	// если нет еще такой комнаты.
	if _, ok := hub.rooms[roomName]; !ok {
		currentRoom = newRoom()
		hub.rooms[roomName] = currentRoom
		go currentRoom.Run()
	} else {
		rest_api_utils.BindUnprocessableEntity(ctx, errors.New("session already started"))
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"video_url": "https://app.flixray.ru/srv/streaming/minecraft.m3u8"})
}

func (hub *Hub) GetUsersCount(ctx *gin.Context) {
	// Достаем room name.
	roomName := ctx.Query("room_name")
	if roomName == "" {
		rest_api_utils.BindBadRequest(ctx, fmt.Errorf("invalid room_name querry"))
		return
	}

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

	// получаем комнату.
	var currentRoom *room
	// если нет еще такой комнаты.
	if _, ok := hub.rooms[roomName]; !ok {
		rest_api_utils.BindUnprocessableEntity(ctx, errors.New("session doesnt start"))
		return
	}
	// достаем комнату.
	currentRoom = hub.rooms[roomName]

	// собираем список юзеров.
	users := make([]dto.User, 0)
	for k := range currentRoom.connections {
		users = append(users, dto.User{
			Name:  k.user.Name,
			Email: k.user.Email,
		})
	}

	// отдаем ответ.
	ctx.JSON(http.StatusOK, dto.CurrentUsersResponse{
		UsersCount: len(users),
		Users:      users,
	})
}

func (hub *Hub) GetCurrentVideoInfo(ctx *gin.Context) {
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

	// собираем ответ.
	response := dto.CurrentVideoInfoResponse{
		Rooms: make([]dto.RoomInfo, 0),
	}
	for i := range roomList {
		if _, ok := hub.rooms[roomList[i].Name]; !ok {
			continue
		}

		// get admin from db
		admin, err := hub.storage.UserStorage.GetUserByID(roomList[i].AdminID)
		if err != nil {
			if strings.Contains(err.Error(), "sql: no rows in result set") {
				rest_api_utils.BindUnauthorized(ctx, errors.New("admin not found"))
				return
			}
			rest_api_utils.BindInternalError(ctx, err)
			return
		}

		currentRoom := hub.rooms[roomList[i].Name]
		var timing int64
		if currentRoom.status == "PAUSE" {
			timing = currentRoom.timePassed
		} else {
			//timing = (int64)(time.Now().Second()*1000 - currentRoom.timeStart.Second()*1000)
			timing = time.Since(currentRoom.timeStart).Milliseconds()
		}
		roomInfo := dto.RoomInfo{
			Name:      roomList[i].Name,
			Timing:    timing,
			Url:       currentRoom.currentVideoUrl,
			UserCount: len(currentRoom.connections),
			Duration:  currentRoom.duration,
			IsPrivate: roomList[i].IsPrivate,
			AdminName: admin.Name,
		}

		response.Rooms = append(response.Rooms, roomInfo)
	}

	// отдаем ответ.
	ctx.JSON(http.StatusOK, response)
}
