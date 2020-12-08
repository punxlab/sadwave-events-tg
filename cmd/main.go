package main

import (
	"context"
	"log"

	"github.com/punxlab/sadwave-events-tg/internal/app"
	"github.com/punxlab/sadwave-events-tg/internal/app/ping"
	"github.com/punxlab/sadwave-events-tg/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err.Error())
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Panic(err.Error())
	}

	ping.Serve()

	log.Print("starting the app")
	defer log.Print("the app has been finished")
	if err = a.Run(ctx); err != nil {
		log.Panic(err.Error())
	}
}
