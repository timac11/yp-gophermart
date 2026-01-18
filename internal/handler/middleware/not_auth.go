package middleware

import (
	"net/http"
)

func (m *Middleware) NotAuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		payload, _ := m.authControl.GetSignedPayload(token)

		if payload != nil {
			// todo: return invlid status
			return
		}

		next.ServeHTTP(w, r)
	})
}
