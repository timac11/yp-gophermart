package handler

import (
	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/service"
)

type Application struct {
	userService    *service.UserService
	orderService   *service.OrderService
	balanceService *service.BalanceService
	jwtControl     *auth.JWTControl
}

func NewApplication(repository service.Repository, jwtControl *auth.JWTControl, config *service.ServiceConfig) *Application {
	userService := service.NewUserService(repository, config)
	orderService := service.NewOrderService(repository, config)
	balanceService := service.NewBalanceService(repository, config)

	return &Application{
		userService:    userService,
		orderService:   orderService,
		balanceService: balanceService,
		jwtControl:     jwtControl,
	}
}
