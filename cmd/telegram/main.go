package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
