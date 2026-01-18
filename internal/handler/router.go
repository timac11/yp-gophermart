package handler

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/handler/middleware"
	"github.com/timac11/yp-gophermart/internal/repository"
)

func InitRouter(conf *config.Config) (*chi.Mux, error) {
	router := chi.NewRouter()
	_, err := repository.NewPgClient(conf.DatabaseUri)

	if err != nil {
		return nil, err
	}

	mParams := middleware.Params{Secret: conf.JWTSecret, TokenExp: (time.Duration(conf.JWTExpMinutes * int64(time.Minute)))}
	m := middleware.NewMiddleware(mParams)

	// setup middlewares
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

			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/balance", func(w http.ResponseWriter, r *http.Request) {})
			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/withdraws", func(w http.ResponseWriter, r *http.Request) {})
		})

		router.Group(func(router chi.Router) {
			middlewares := []func(http.Handler) http.Handler{
				m.NotAuthCheckMiddleware,
			}
			router.Use(middlewares...)

			router.Post("/register", func(w http.ResponseWriter, r *http.Request) {})
			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
		})
	})

	return router, nil
}
