package handlers

import (
	"context"
	"log/slog"

	"github.com/AmadeusAI-dev/telegram-service/internal/client"
	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/tg"
)

type Bus interface {
	Dispatch(context.Context, pubsub.Event) error
}

type UserRepo interface {
	Get(context.Context, int) (client.User, error)
}

func DispatchNewMessages(bus Bus, repo UserRepo, d *tg.UpdateDispatcher) {
	d.OnNewMessage(func(ctx context.Context, e tg.Entities, u *tg.UpdateNewMessage) error {
		m, ok := u.Message.(*tg.Message)
		if !ok {
			slog.Error("failed to get new message from update")
			return nil
		}

		if m.Out {
			// ignore its own messages
			return nil
		}

		PeerUser, ok := m.PeerID.(*tg.PeerUser)
		if !ok {
			slog.Error("failed to get PeerUser from message")
			return nil
		}

		user, err := repo.Get(ctx, int(PeerUser.UserID))
		if err != nil {
			slog.Error("failed to get user", "errors", err)
		}

		err = bus.Dispatch(ctx, pubsub.Event{
			Type: "new_message",
			Payload: map[string]any{
				"user_id":    user.ID,
				"username":   user.Username,
				"message_id": m.ID,
				"message":    m.Message,
			},
		})
		if err != nil {
			slog.Error(
				"failed to dispatch new message",
				"error", err,
				"user_id", user.ID,
				"username", user.Username,
				"message_id", m.ID,
				"message", m.Message,
			)
			return nil
		}

		slog.Info(
			"dispatched new message",
			"user_id", user.ID,
			"username", user.Username,
			"message_id", m.ID,
			"message", m.Message,
		)
		return nil
	})
}
