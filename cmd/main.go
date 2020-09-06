package main

import (
	"context"
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

	if err = a.Run(ctx); err != nil {
		log.Panic(err)
	}
}
