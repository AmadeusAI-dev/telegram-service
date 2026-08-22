package main

import (
	"context"
	"log"

	"github.com/AmadeusAI-dev/telegram-service/internal/app"
	"github.com/AmadeusAI-dev/telegram-service/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	application, err := app.New(ctx, *cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
