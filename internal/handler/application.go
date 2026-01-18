package handler

import (
	"github.com/timac11/yp-gophermart/internal/service"
)

type Repository interface {
	service.UserRepository
}

type Application struct {
	userService *service.UserService
}

func NewApplication(repository Repository, config service.ServiceConfig) *Application {
	userService := service.NewUserService(repository, config)

	return &Application{
		userService: userService,
	}
}
