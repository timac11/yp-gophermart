package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/handler/middleware"
	"github.com/timac11/yp-gophermart/internal/repository"
	"github.com/timac11/yp-gophermart/internal/service"
)

func InitRouter(conf *config.Config) (*chi.Mux, error) {
	router := chi.NewRouter()
	_, err := repository.NewPgClient(conf.DatabaseUri)

	if err != nil {
		return nil, err
	}
	// init application
	repo, err := repository.NewPgClient(conf.DatabaseUri)
	if err != nil {
		return nil, err
	}

	jwtControl := auth.JWTControl{TokenExp: (time.Duration(conf.JWTExpMinutes * uint(time.Minute))), Secret: conf.JWTSecret}

	application := NewApplication(repo, &jwtControl, &service.ServiceConfig{Attempts: conf.RetryAttempts, AttemptsInterval: conf.RetryInterval})

	// setup middlewares
	mParams := middleware.Params{JwtControl: &jwtControl}
	m := middleware.NewMiddleware(mParams)

	middlewares := []func(http.Handler) http.Handler{
		m.LoggingMiddleware,
		m.GzipMiddleware,
	}
	router.Use(middlewares...)

	router.Route("/api/user", func(router chi.Router) {
		router.Group(func(router chi.Router) {
			middlewares := []func(http.Handler) http.Handler{
				m.AuthCheckMiddleware,
			}
			router.Use(middlewares...)

			router.Post("/orders", application.CreateOrder)
			router.Get("/orders", application.GetOrders)
			router.Get("/balance", application.GetBalance)
			router.Post("/balance/withdraw", application.CreateWithdraw)
			router.Get("/withdrawals", application.GetWithdrawals)
		})

		router.Group(func(router chi.Router) {
			middlewares := []func(http.Handler) http.Handler{
				m.NotAuthCheckMiddleware,
			}
			router.Use(middlewares...)

			router.Post("/register", application.Register)
			router.Post("/login", application.Login)
		})
	})

	return router, nil
}
