package main

import (
	"context"
	"log"

	"github.com/AmadeusAI-dev/telegram-service/internal/client/handlers"
	"github.com/AmadeusAI-dev/telegram-service/internal/config"
	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := config.Load()

	bus, err := pubsub.New(ctx, config.RabbitMq)
	if err != nil {
		log.Fatalf("failed to create pubsub: %v", err)
	}

	// @todo: recovery engine
	d := tg.NewUpdateDispatcher()

	client, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler: d,
	})
	if err != nil {
		log.Fatalf("failed to create telegram client: %v", err)
	}

	handlers.DispatchNewMessages(bus, &d)

	err = client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			log.Fatalf("failed to get auth status: %v", err)
		}

		if !status.Authorized {
			log.Fatalf("auth session is invalid. Please, update session")
		}

		return telegram.RunUntilCanceled(ctx, client)
	})
}
