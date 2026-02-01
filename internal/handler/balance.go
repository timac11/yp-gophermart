package handler

import (
	"encoding/json"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/logger"
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
