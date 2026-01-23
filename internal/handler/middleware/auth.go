package middleware

import (
	"net/http"

	"github.com/timac11/yp-gophermart/internal/auth"
)

func (m *Middleware) AuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		payload, error := m.jwtControl.GetSignedPayload(token)

		if error != nil {
			http.Error(w, "User unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := auth.ContextWithAuthPayload(r.Context(), payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
