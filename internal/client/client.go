package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/telegram"
)

func Initialize(ctx context.Context, client *telegram.Client, runCh chan<- error) <-chan error {
	initCh := make(chan error, 1)

	go func() {
		runCh <- client.Run(ctx, func(ctx context.Context) error {
			status, err := client.Auth().Status(ctx)
			if err != nil {
				initCh <- fmt.Errorf("failed to get auth status: %w", err)
				return nil
			}

			if !status.Authorized {
				initCh <- fmt.Errorf("auth session is invalid. Please, update session")
				return nil
			}

			initCh <- nil

			return telegram.RunUntilCanceled(ctx, client)
		})
	}()

	return initCh

}

func WaitForInitialization(ctx context.Context, ch <-chan error) error {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case err := <-ch:
		if err != nil {
			return err
		}
	case <-timer.C:
		return fmt.Errorf("telegram client initialization timeout")
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
