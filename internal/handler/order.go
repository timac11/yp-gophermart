package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/common/util"
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
	log := logger.LoggerFromContext(r.Context())

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	orderNum := string(bodyBytes)

	if !util.CheckOrderNum(orderNum) {
		log.Error("Order num is not valid")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = app.orderService.CreateOrder(r.Context(), orderNum)

	if err != nil {
		if errors.IsEntityAlreadyExists(err) {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Error(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
