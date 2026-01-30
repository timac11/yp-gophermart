package handler

import (
	"encoding/json"
	_errors "errors"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/model"
)

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	user, err := app.parseUserFromBody(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	loginUser, err := app.userService.Login(ctx, user)

	if err != nil {
		code := handleError(err)
		http.Error(w, err.Error(), code)
		return
	}

	signedString, err := app.jwtControl.BuildJWTString(auth.JWTPayload{UserID: loginUser.Id})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", signedString)
	w.WriteHeader(http.StatusOK)
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	user, err := app.parseUserFromBody(r)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	registeredUser, err := app.userService.Register(ctx, user)

	if err != nil {
		code := handleError(err)
		http.Error(w, err.Error(), code)
		return
	}

	signedString, err := app.jwtControl.BuildJWTString(auth.JWTPayload{UserID: registeredUser.Id})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", signedString)
}

func (app *Application) parseUserFromBody(req *http.Request) (*model.UserLoginDto, error) {
	var user model.UserLoginDto

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func handleError(err error) int {
	var entityError *errors.EntityError

	if _errors.As(err, &entityError) {
		switch entityError.Type {
		case errors.EntityNotFound:
			return http.StatusNotFound
		case errors.EntityAlreadyExists:
			return http.StatusConflict
		default:
			return http.StatusInternalServerError
		}
	}

	var passErr *errors.InvalidPasswordError
	if _errors.As(err, &passErr) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}
