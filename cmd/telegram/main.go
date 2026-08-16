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

		api := client.API()

		ds, err := api.HelpGetNearestDC(ctx)
		if err != nil {
			// will be removed; example code
			panic(err)
		}

		log.Println("ds: %w", ds)
		return nil
	})
}
