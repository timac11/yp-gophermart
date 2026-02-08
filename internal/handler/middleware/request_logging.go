package middleware

import (
	"bytes"
	"net/http"
	"time"

	"github.com/timac11/yp-gophermart/internal/logger"
	"go.uber.org/zap"
)

type ResponseRecorder struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)                  // Record the body
	return r.ResponseWriter.Write(b) // Write to the actual response
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode                    // Record the status code
	r.ResponseWriter.WriteHeader(statusCode) // Write the actual header
}

func (m *Middleware) RequestLoggerMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			return
		}

		recorder := &ResponseRecorder{
			ResponseWriter: w,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}

		start := time.Now()
		h.ServeHTTP(recorder, r)
		end := time.Now()

		log := logger.LoggerFromContext(r.Context())

		log.Info(
			"Request info:",
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.String("duration", end.Sub(start).String()),
		)

		log.Info(
			"Response info",
			zap.Int("status", recorder.status),
			zap.Int("bytes", recorder.body.Len()),
		)
	})
}
