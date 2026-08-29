package client

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

type Sender struct {
	Client *telegram.Client
}

func (s *Sender) Send(ctx context.Context, chatID int, msg string) error {
	sender := message.NewSender(s.Client.API())

	_, err := sender.To(&tg.InputPeerUser{UserID: int64(chatID)}).Text(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}
