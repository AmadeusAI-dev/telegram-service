package main

import (
	"context"
	"log"

	"github.com/gotd/td/telegram"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := telegram.ClientFromEnvironment(telegram.Options{})
	if err != nil {
		log.Fatalf("failed to create telegram client: %v", err)
	}

	err = client.Run(ctx, func(ctx context.Context) error {
		status, err := client.Auth().Status(ctx)
		if err != nil {
			log.Fatalf("failed to get auth status: %v", err)
		}

		if !status.Authorized {
			log.Fatalf("auth session is invalid. Please, update session")
		}

		// api := client.API()

		// it's temp example code; it will be removed soon
		self, _ := client.Self(ctx)
		log.Printf("Logged in as %s", self.Username)

		return nil
	})
}
