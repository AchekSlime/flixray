package service

import (
	"errors"
	"fmt"
	"github.com/achekslime/core/encoding"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/flixray/auth/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (auth *AuthService) GenerateToken(context *gin.Context) {
	var request dto.TokenRequest

	// get request body
	if err := context.ShouldBindJSON(&request); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// get user from db
	dbUser, err := auth.storage.UserStorage.GetUserByEmail(request.Email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			rest_api_utils.BindUnprocessableEntity(context, errors.New("user not found"))
			return
		}
		rest_api_utils.BindInternalError(context, fmt.Errorf("db err: %s", err.Error()))
		return
	}

	// validate password
	credentialError := encoding.CheckPassword(request.Password, dbUser.Password)
	if credentialError != nil {
		rest_api_utils.BindUnauthorized(context, errors.New("invalid password"))
		return
	}

	tokenString, err := auth.jwtService.GenerateJWT(dbUser.Email)
	if err != nil {
		rest_api_utils.BindInternalError(context, fmt.Errorf("token generation err: %s", err.Error()))
		return
	}
	context.JSON(http.StatusOK, gin.H{"token": tokenString, "name": dbUser.Name, "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImltcHNmYWNlQGdtYWlsLnJ1IiwiZXhwIjoxNjg0MjgzNzU1fQ.BmAcD9-NG6L7UUFP7jOKyDT_PIUgjeFjyw6jCj4vsoo"})
}
