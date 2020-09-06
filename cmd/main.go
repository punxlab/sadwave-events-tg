package main

import (
	"context"
	"github.com/punxlab/sadwave-events-tg/internal/app/ping"
	"github.com/punxlab/sadwave-events-tg/internal/config"
	"log"

	"github.com/punxlab/sadwave-events-tg/internal/app"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Panic(err)
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Panic(err)
	}

	ping.Serve()

	log.Print("starting the app")
	defer log.Print("the app has been finished")
	if err = a.Run(ctx); err != nil {
		log.Panic(err)
	}
}
