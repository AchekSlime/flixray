package service

import (
	"fmt"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/flixray/user_api/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (srv *UserService) FindByName(context *gin.Context) {
	// authorization.
	user := srv.Authorize(context)
	if user == nil {
		return
	}

	userFindName := context.Query("name")
	if userFindName == "" {
		rest_api_utils.BindBadRequest(context, fmt.Errorf("invalid name"))
		return
	}

	user, err := srv.storage.UserStorage.GetUserByName(userFindName)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			rest_api_utils.BindNoContent(context)
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	response := dto.FindByName{
		Email: user.Email,
		Name:  user.Name,
	}
	context.JSON(http.StatusOK, response)
}
