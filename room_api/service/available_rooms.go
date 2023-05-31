package service

import (
	"fmt"
	"github.com/achekslime/flixray/room_api/service/dto"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage/postgres"
)

func (srv *RoomService) GetAvailableRooms(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	roomList, err := srv.storage.RoomStorage.GetAvailableRooms(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindNoContent(context)
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	roomListResponse := make([]dto.AvailableRoom, 0)
	for i, room := range roomList {
		roomResponse := dto.AvailableRoom{
			Name:      room.Name,
			IsPrivate: room.IsPrivate,
		}

		// get user from db
		dbUser, _ := srv.storage.UserStorage.GetUserByID(roomList[i].AdminID)
		responseUser := dto.AvailableUser{
			Name:  dbUser.Name,
			Email: dbUser.Email,
		}

		roomResponse.Admin = &responseUser
		roomListResponse = append(roomListResponse, roomResponse)
	}

	context.JSON(http.StatusOK, roomListResponse)
}
