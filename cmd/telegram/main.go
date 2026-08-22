package main

import (
	"context"
	"fmt"
	"log"

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

	d.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		m, ok := u.Message.(*tg.Message)
		if !ok {
			return nil
		}

		err = bus.Dispatch(ctx, pubsub.Event{
			Type: "new_message",
			Payload: map[string]any{
				"message": m.Message,
				"chat_id": m.PeerID.TypeID(),
			},
		})
		if err != nil {
			// just loggging instead of failing
			log.Fatalf("failed to dispatch message: %v", err)
		}

		fmt.Printf("message: %s\n", m.Message)
		return nil
	})

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
