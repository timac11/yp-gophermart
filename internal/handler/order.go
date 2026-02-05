package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/errors"
	"github.com/timac11/yp-gophermart/internal/logger"
)

func (app *Application) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := app.orderService.GetOrders(r.Context())
	log := logger.LoggerFromContext(r.Context())

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	returnBody, err := json.Marshal(orders)

	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(returnBody)
}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orderNum := string(bodyBytes)

	_, err = app.orderService.CreateOrder(r.Context(), orderNum)
	if err != nil {
		if errors.IsEntityAlreadyExistsErr(err) {
			w.WriteHeader(http.StatusOK)
			return
		}

		handleOrderError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getOrderErrorStatusCode(err error) int {
	if errors.IsInvalidOrderNumErr(err) {
		return http.StatusUnprocessableEntity
	}

	return http.StatusInternalServerError
}

func handleOrderError(w http.ResponseWriter, r *http.Request, err error) {
	code := getOrderErrorStatusCode(err)
	if code == http.StatusInternalServerError {
		log := logger.LoggerFromContext(r.Context())
		log.Error(err.Error())
	}

	http.Error(w, http.StatusText(code), code)
}
