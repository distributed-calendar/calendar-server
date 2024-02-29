package main

import (
	"os"

	"github.com/distributed-calendar/calendar-server/internal/app"
)

func main() {
	cfgPath := os.Getenv("CONFIG_PATH")

	app, err := app.NewApp(cfgPath)
	if err != nil {
		panic(err)
	}

	app.Run()
}
