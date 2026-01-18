package auth

import (
	"context"
	"errors"
)

type CtxKey string

const CtxAuthKey CtxKey = "auth"

func ContextWithAuthPayload(ctx context.Context, payload *JWTPayload) context.Context {
	return context.WithValue(ctx, CtxAuthKey, payload)
}

func AuthPayloadFromContext(ctx context.Context) (*JWTPayload, error) {
	payload, ok := ctx.Value(CtxAuthKey).(*JWTPayload)
	if !ok {
		return nil, errors.New("User is unauthorized")
	}
	return payload, nil
}
