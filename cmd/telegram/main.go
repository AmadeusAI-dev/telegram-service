package main

import (
	"context"
	"log"
	"log/slog"
	"os/signal"
	"syscall"

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
	defer func() {
		err := application.Close(ctx)
		if err != nil {
			slog.Error("failed to close application", "error", err)
		}
	}()

	return application.Run(ctx)
}
