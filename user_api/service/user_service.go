package service

import (
	"errors"
	"github.com/achekslime/core/jwt"
	"github.com/achekslime/core/models"
	"github.com/achekslime/core/rest_api_utils"
	"github.com/achekslime/core/storage"
	"github.com/gin-gonic/gin"
	"strings"
)

type UserService struct {
	storage    *storage.Storage
	jwtService *jwt.JwtService
}

func NewUserService() *UserService {
	return &UserService{}
}

func (srv *UserService) Configure(jwtKey string) error {
	// storage configuration.
	newStorage, err := storage.NewStorage()
	if err != nil {
		return err
	}
	srv.storage = newStorage

	// jwt configuration.
	srv.jwtService = jwt.ConfigureJWT(jwtKey)
	return nil
}

func (srv *UserService) Authorize(context *gin.Context) *models.User {
	userClaims, err := srv.GetTokenClaims(context)
	if err != nil {
		rest_api_utils.BindUnauthorized(context, err)
		return nil
	}

	// get user from db
	dbUser, err := srv.storage.UserStorage.GetUserByEmail(userClaims.Email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			rest_api_utils.BindUnauthorized(context, errors.New("user not found"))
			return nil
		}
		rest_api_utils.BindInternalError(context, err)
		return nil
	}

	return dbUser
}

func (srv *UserService) GetTokenClaims(context *gin.Context) (*jwt.JWTClaim, error) {
	token, err := rest_api_utils.ExtractToken(context)
	if err != nil {
		return nil, err
	}
	// validate token.
	userClaims, err := srv.jwtService.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return userClaims, nil
}
