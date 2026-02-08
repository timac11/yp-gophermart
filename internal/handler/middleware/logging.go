package middleware

import (
	"net/http"

	"github.com/timac11/yp-gophermart/internal/logger"
	"go.uber.org/zap"
)

func (*Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lg := logger.LoggerFromContext(r.Context())
		lg = lg.With(
			zap.String("url", r.URL.String()),
			zap.String("method", r.Method),
		)
		ctx := logger.ContextWithLogger(r.Context(), lg)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
