package handler

import (
	"net/http"
	"encoding/json"

	"github.com/timac11/yp-gophermart/internal/model"
)

func Login(w http.ResponseWriter, r *http.Request) error {
	_, err := parseUserFromBody(r)

	if err != nil {
		return err
	}

	return nil
}

func Register(w http.ResponseWriter, r *http.Request) error {
	_, err := parseUserFromBody(r)

	if err != nil {
		return err
	}

	return nil
}

func parseUserFromBody(req *http.Request) (*model.User, error) {
	var user model.User

	err := json.NewDecoder(req.Body).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}