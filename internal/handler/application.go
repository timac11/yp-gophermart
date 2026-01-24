package handler

import (
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/service"
)

type Repository interface {
	service.UserRepository
	service.OrderRepository
}

type Application struct {
	userService  *service.UserService
	orderService *service.OrderService
	jwtControl   *auth.JWTControl
}

func NewApplication(repository Repository, jwtControl *auth.JWTControl, config *service.ServiceConfig) *Application {
	userService := service.NewUserService(repository, config)
	orderService := service.NewOrderService(repository, config)

	return &Application{
		userService:  userService,
		orderService: orderService,
		jwtControl:   jwtControl,
	}
}
