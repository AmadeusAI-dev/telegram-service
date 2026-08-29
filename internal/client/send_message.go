package client

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
)

type Sender struct {
	Client *telegram.Client
}

func (s *Sender) Send(ctx context.Context, Username string, msg string) error {
	sender := message.NewSender(s.Client.API())

	_, err := sender.Resolve(Username).Text(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}
