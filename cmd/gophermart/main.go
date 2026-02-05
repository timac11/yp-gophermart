package main

import (
	"log"

	"github.com/timac11/yp-gophermart/internal/app"
)

func main() {
	err := app.RunApplication()
	if err != nil {
		log.Fatal(err)
	}
}
