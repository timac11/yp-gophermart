package middleware

import (
	"github.com/timac11/yp-gophermart/internal/auth"
)

type Params struct {
	JwtControl *auth.JWTControl
}

type Middleware struct {
	jwtControl *auth.JWTControl
}

func NewMiddleware(params Params) *Middleware {
	mw := Middleware{
		jwtControl: params.JwtControl,
	}
	return &mw
}
