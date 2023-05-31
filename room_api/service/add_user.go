package service

import (
	"errors"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage/postgres"
	"github.com/achekslime/flixray/room_api/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (srv *RoomService) AddUser(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	// get room from request body.
	var roomRequest dto.AddUserRequest
	if err := context.ShouldBindJSON(&roomRequest); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// проверить пользователя
	userDb, err := srv.storage.UserStorage.GetUserByName(roomRequest.UserName)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindUnprocessableEntity(context, errors.New("user not found"))
			return
		}
		rest_api_utils.BindInternalError(context, err)
		return
	}

	// проверить комнату
	roomDb, err := srv.storage.RoomStorage.GetRoomByName(roomRequest.RoomName)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindUnprocessableEntity(context, errors.New("room not found"))
			return
		}
		rest_api_utils.BindInternalError(context, err)
		return
	}

	// db.
	err = srv.storage.RoomStorage.AddUserToRoom(roomDb.ID, userDb.ID)
	if err != nil {
		rest_api_utils.BindInternalError(context, err)
		return
	}

	context.JSON(http.StatusOK, nil)
}
