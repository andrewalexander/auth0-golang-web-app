package main

import (
	"app"

	log "github.com/sirupsen/logrus"
)

func main() {
	if err := app.Init(); err != nil {
		log.Fatal(err)
	}
	StartServer()
}
