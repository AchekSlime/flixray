package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/achekslime/core/models"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/flixray/room_api/service/dto"
	"github.com/gin-gonic/gin"

	"github.com/achekslime/core/storage/postgres"
)

func (srv *RoomService) CreateNewRoom(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	// get room from request body.
	var roomRequest dto.Room
	if err := context.ShouldBindJSON(&roomRequest); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// create room.
	room := models.Room{
		Name:      roomRequest.Name,
		AdminID:   user.ID,
		IsPrivate: roomRequest.IsPrivate,
	}

	id, err := srv.storage.RoomStorage.SaveRoom(&room)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrUniqueConstraintDuplicate) {
			rest_api_utils.BindUnprocessableEntity(context, fmt.Errorf("room with this name already exists"))
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	srv.addUsersToRoom(id, roomRequest)
	context.JSON(http.StatusCreated, room)
}

func (srv *RoomService) addUsersToRoom(id int, room dto.Room) {
	for i := range room.EmailList {
		user, err := srv.storage.UserStorage.GetUserByEmail(room.EmailList[i])
		if err != nil {
			continue
		}
		srv.storage.RoomStorage.AddUserToRoom(id, user.ID)
	}
}
