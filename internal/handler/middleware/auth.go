package middleware

import (
	"net/http"

	"github.com/timac11/yp-gophermart/internal/auth"
)

func (m *Middleware) AuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		payload, error := m.authControl.GetSignedPayload(token)

		if error != nil {
			// todo: return 401 status
			return
		}

		ctx := auth.ContextWithAuthPayload(r.Context(), payload)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
