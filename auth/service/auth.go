package service

import (
	"github.com/achekslime/core/jwt"
	"github.com/achekslime/core/storage"
)

type AuthService struct {
	storage         *storage.Storage
	jwtService      *jwt.JwtService
	unconfirmedMail map[int64]string
	confirmedMail   map[string]bool
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (auth *AuthService) Configure(jwtKey string) error {
	// storage configuration.
	newStorage, err := storage.NewStorage()
	if err != nil {
		return err
	}
	auth.storage = newStorage

	// configure map.
	auth.unconfirmedMail = make(map[int64]string)
	auth.confirmedMail = make(map[string]bool)

	// jwt configuration.
	auth.jwtService = jwt.ConfigureJWT(jwtKey)
	return nil
}
