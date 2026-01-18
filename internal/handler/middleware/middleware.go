package middleware

import (
	"time"

	"github.com/timac11/yp-gophermart/internal/auth"
)

type Params struct {
	TokenExp time.Duration
	Secret   string
}

type Middleware struct {
	authControl auth.JWTControl
}

func NewMiddleware(params Params) *Middleware {
	control := auth.JWTControl{TokenExp: params.TokenExp, Secret: params.Secret}
	mw := Middleware{
		authControl: control,
	}
	return &mw
}
