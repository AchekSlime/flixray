package service

import (
	"fmt"
	"github.com/achekslime/flixray/auth/service/dto"
	"net/http"
	"strings"

	"github.com/achekslime/core/encoding"
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/gin-gonic/gin"
)

func (auth *AuthService) RegNewUser(context *gin.Context) {
	var user models.User

	// get body from request.
	if err := context.ShouldBindJSON(&user); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// достаем почту из мапы.
	_, ok := auth.confirmedMail[user.Email]
	if !ok {
		rest_api_utils.BindBadRequest(context, fmt.Errorf("mail unconfirmed"))
		return
	}

	// encrypt password.
	password, err := encoding.HashPassword(user.Password)
	if err != nil {
		rest_api_utils.BindInternalError(context, err)
		return
	}
	user.Password = *password

	// проверить что пользователя с таким мылом не существует.
	_, err = auth.storage.UserStorage.GetUserByEmail(user.Email)
	if err == nil {
		msg := fmt.Sprintf("user with email:%s already registered", user.Email)
		rest_api_utils.BindUnprocessableEntity(context, fmt.Errorf(msg))
		return
	} else if !strings.Contains(err.Error(), "sql: no rows in result set") {
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	// проверить что пользователя с таким name не существует.
	_, err = auth.storage.UserStorage.GetUserByName(user.Name)
	if err == nil {
		msg := fmt.Sprintf("user with name:%s already registered", user.Name)
		rest_api_utils.BindUnprocessableEntity(context, fmt.Errorf(msg))
		return
	} else if !strings.Contains(err.Error(), "sql: no rows in result set") {
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	// db saving.
	err = auth.storage.UserStorage.SaveUser(&user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			rest_api_utils.BindUnprocessableEntity(context, fmt.Errorf("user with this email already exists"))
			// удаляем из подтвержденных почт.
			delete(auth.confirmedMail, user.Email)
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	// удаляем из подтвержденных почт.
	delete(auth.confirmedMail, user.Email)
	response := dto.RegUserResponse{
		Name:  user.Name,
		Email: user.Email,
	}
	context.JSON(http.StatusCreated, response)
}
