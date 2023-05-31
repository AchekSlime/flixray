package service

import (
	"fmt"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/flixray/auth/service/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (auth *AuthService) Confirm(context *gin.Context) {
	// get body from request.
	var confirmRequest dto.ConfirmRequest
	if err := context.ShouldBindJSON(&confirmRequest); err != nil {
		rest_api_utils.BindBadRequest(context, err)
		return
	}

	// достаем почту из мапы.
	emailFromMap, ok := auth.unconfirmedMail[confirmRequest.Code]
	if !ok {
		rest_api_utils.BindNoContent(context)
		return
	}

	// кладем в подтвержденные тары.
	auth.confirmedMail[emailFromMap] = true

	// если мыло не совпадает.
	if emailFromMap != confirmRequest.Email {
		rest_api_utils.BindUnauthorized(context, fmt.Errorf("mail does not exist, try /api/reg_mail"))
		return
	}

	// убираем из мапы неподтвержденных.
	delete(auth.unconfirmedMail, confirmRequest.Code)

	// correct result.
	context.JSON(http.StatusOK, nil)
}
