package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/auth"
	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
)

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	user, err := parseUserFromBody(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	loginUser, err := app.userService.Login(ctx, user)

	if err != nil {
		handleUserError(w, r, err)
		return
	}

	signedString, err := app.jwtControl.BuildJWTString(auth.JWTPayload{UserID: loginUser.ID})

	if err != nil {
		handleUserError(w, r, err)
		return
	}

	w.Header().Set("Authorization", signedString)
	w.WriteHeader(http.StatusOK)
}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	user, err := parseUserFromBody(r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	registeredUser, err := app.userService.Register(ctx, user)

	if err != nil {
		handleUserError(w, r, err)
		return
	}

	signedString, err := app.jwtControl.BuildJWTString(auth.JWTPayload{UserID: registeredUser.ID})

	if err != nil {
		handleUserError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Authorization", signedString)
}

func parseUserFromBody(req *http.Request) (*model.UserLoginDto, error) {
	var user model.UserLoginDto

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func getUserErrorStatusCode(err error) int {
	if errors.IsEntityNotFoundErr(err) {
		return http.StatusNotFound
	} else if errors.IsEntityAlreadyExistsErr(err) {
		return http.StatusConflict
	} else if errors.IsInvalidPasswordError(err) {
		return http.StatusUnauthorized
	}

	return http.StatusInternalServerError
}

func handleUserError(w http.ResponseWriter, r *http.Request, err error) {
	code := getUserErrorStatusCode(err)
	if code == http.StatusInternalServerError {
		log := logger.LoggerFromContext(r.Context())
		log.Error(err.Error())
	}

	http.Error(w, http.StatusText(code), code)
}
