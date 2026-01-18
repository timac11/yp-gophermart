package main

import (
	"log"
	"net/http"

	"github.com/timac11/yp-gophermart/internal/config"
	"github.com/timac11/yp-gophermart/internal/handler"
)

func main() {
	runServer()
}

func runServer() {
	conf := config.InitConfig()
	mux, err := handler.InitRouter(conf)

	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(conf.Address, mux)
	if err != nil {
		log.Fatal(err)
	}
}
