package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/logger"
	"github.com/timac11/yp-gophermart/internal/model"
)

func (app *Application) GetBalance(w http.ResponseWriter, r *http.Request) {
	balance, err := app.balanceService.GetBalance(r.Context())
	log := logger.LoggerFromContext(r.Context())

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	returnBody, err := json.Marshal(balance)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnBody)
}

func (app *Application) CreateWithdraw(w http.ResponseWriter, r *http.Request) {
	createWithdraw, err := parseWithdrawFromBody(r)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = app.balanceService.CreateWithdraw(r.Context(), createWithdraw)
	log := logger.LoggerFromContext(r.Context())

	if err != nil {
		errorNum := http.StatusInternalServerError

		if errors.IsInvalidOrderNumErr(err) {
			errorNum = http.StatusBadRequest
		}

		log.Error(err.Error())

		http.Error(w, http.StatusText(errorNum), errorNum)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *Application) GetWithdrawals(w http.ResponseWriter, r *http.Request) {
	withdrawals, err := app.balanceService.GetWithdrawals(r.Context())
	log := logger.LoggerFromContext(r.Context())

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	returnBody, err := json.Marshal(withdrawals)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnBody)
}

func parseWithdrawFromBody(req *http.Request) (*model.CreateWithdraw, error) {
	var withdraw model.CreateWithdraw

	err := json.NewDecoder(req.Body).Decode(&withdraw)

	if err != nil {
		return nil, err
	}

	return &withdraw, nil
}
