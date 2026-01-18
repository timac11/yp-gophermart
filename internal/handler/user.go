package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/model"
)

func (app *Application) Login(w http.ResponseWriter, r *http.Request) {
	_, err := app.parseUserFromBody(r)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func (app *Application) Register(w http.ResponseWriter, r *http.Request) {
	_, err := app.parseUserFromBody(r)

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

}

func (app *Application) parseUserFromBody(req *http.Request) (*model.UserDto, error) {
	var user model.UserDto

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
