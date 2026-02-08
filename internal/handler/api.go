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

func NewApplication(repository service.Repository, jwtControl *auth.JWTControl) *Application {
	userService := service.NewUserService(repository)
	orderService := service.NewOrderService(repository)
	balanceService := service.NewBalanceService(repository)

	return &Application{
		userService:    userService,
		orderService:   orderService,
		balanceService: balanceService,
		jwtControl:     jwtControl,
	}
}
