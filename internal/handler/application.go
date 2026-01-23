package handler

import (
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/service"
)

type Repository interface {
	service.UserRepository
}

type Application struct {
	userService *service.UserService
	jwtControl  *auth.JWTControl
}

func NewApplication(repository Repository, jwtControl *auth.JWTControl, config *service.ServiceConfig) *Application {
	userService := service.NewUserService(repository, config)

	return &Application{
		userService: userService,
		jwtControl:  jwtControl,
	}
}
