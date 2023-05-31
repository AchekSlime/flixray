package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage/postgres"
)

func (srv *RoomService) GetAdminRooms(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	rooms, err := srv.storage.RoomStorage.GetRoomsByAdminID(user.ID)
	if err != nil {
		if strings.Contains(err.Error(), postgres.ErrSqlNoRows) {
			rest_api_utils.BindNoContent(context)
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	context.JSON(http.StatusOK, rooms)
}
