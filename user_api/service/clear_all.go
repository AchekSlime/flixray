package service

import (
	"fmt"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (srv *UserService) ClearAll(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	err := srv.storage.RoomStorage.ClearAll()
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			rest_api_utils.BindNoContent(context)
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	context.JSON(http.StatusOK, "db cleared")
}
