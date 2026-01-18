package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTPayload struct {
	UserID string
}

type JWTClaim struct {
	jwt.RegisteredClaims
	UserID string
}

type JWTControl struct {
	TokenExp time.Duration
	Secret   string
}

func (control *JWTControl) BuildJWTString(payload JWTPayload) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(control.TokenExp)),
		},
		UserID: payload.UserID,
	})

	tokenString, err := token.SignedString([]byte(control.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (control *JWTControl) GetSignedPayload(tokenString string) (*JWTPayload, error) {
	claim := &JWTClaim{}
	token, err := jwt.ParseWithClaims(tokenString, claim,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(control.Secret), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("Invalid token")
	}

	payload := JWTPayload{UserID: claim.UserID}

	return &payload, nil
}
