package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/AmadeusAI-dev/telegram-service/internal/app"
	"github.com/AmadeusAI-dev/telegram-service/internal/config"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	cfg := config.Load()

	application, err := app.New(ctx, *cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := run(application, ctx); err != nil {
		log.Fatal(err)
	}
}

func run(application *app.App, ctx context.Context) error {
	defer shutdown(application)

	return application.Run(ctx)
}

func shutdown(application *app.App) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second*30,
	)
	defer cancel()

	err := application.Close(ctx)
	if err != nil {
		slog.Error("failed to close application", "error", err)
	}
}
