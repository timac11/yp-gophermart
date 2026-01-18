package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/repository"
)

func InitRouter(conf *config.Config) (*chi.Mux, error) {
	router := chi.NewRouter()
	_, err := repository.NewPgClient(conf.DatabaseUri)

	if err != nil {
		return nil, err
	}

	router.Route("/api/user", func(router chi.Router) {
		router.Group(func(router chi.Router) {
			// must be auth routes (check it in middleware)
			router.Post("/orders", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/orders", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/balance", func(w http.ResponseWriter, r *http.Request) {})
			router.Post("/balance/withdraw", func(w http.ResponseWriter, r *http.Request) {})
			router.Get("/withdraws", func(w http.ResponseWriter, r *http.Request) {})
		})

		router.Group(func(router chi.Router) {
			// must not be auth routes (check it in middleware)
			router.Post("/register", func(w http.ResponseWriter, r *http.Request) {})
			router.Post("/login", func(w http.ResponseWriter, r *http.Request) {})
		})
	})

	return router, nil
}
