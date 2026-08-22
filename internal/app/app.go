package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AmadeusAI-dev/telegram-service/internal/client/handlers"
	"github.com/AmadeusAI-dev/telegram-service/internal/config"
	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

type App struct {
	pubsub *pubsub.Pubsub
	client *telegram.Client
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// @todo: recovery engine
	d := tg.NewUpdateDispatcher()

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler: d,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	bus, err := pubsub.New(ctx, cfg.RabbitMq)
	if err != nil {
		return nil, fmt.Errorf("failed to init pubsub: %w", err)
	}

	handlers.DispatchNewMessages(bus, &d)

	return &App{
		pubsub: bus,
		client: client,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	err := a.client.Run(ctx, func(ctx context.Context) error {
		status, err := a.client.Auth().Status(ctx)
		if err != nil {
			return fmt.Errorf("failed to get auth status: %w", err)
		}

		if !status.Authorized {
			return fmt.Errorf("auth session is invalid. Please, update session")
		}

		slog.Info("application started successfully")

		return telegram.RunUntilCanceled(ctx, a.client)
	})
	if err != nil {
		return fmt.Errorf("failed to run telegram client: %w", err)
	}

	return nil
}

func (a *App) Close(ctx context.Context) error {
	return a.pubsub.Close(ctx)
}
