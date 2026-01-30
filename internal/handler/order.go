package handler

import (
	"io"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/common/util"
)

func (app *Application) GetOrders(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) CreateOrder(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	orderNum := string(bodyBytes)

	if !util.CheckOrderNum(orderNum) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err = app.orderService.CreateOrder(r.Context(), orderNum)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
