package handlers

import (
	"context"
	"log/slog"

	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/tg"
)

func DispatchNewMessages(bus *pubsub.Pubsub, d *tg.UpdateDispatcher) {

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

		user, ok := m.FromID.(*tg.PeerUser)
		if !ok {
			slog.Error("failed to get user from message")
			return nil
		}

		err := bus.Dispatch(ctx, pubsub.Event{
			Type: "new_message",
			Payload: map[string]any{
				"user_id":    user.UserID,
				"chat_id":    user.UserID,
				"message_id": m.ID,
				"message":    m.Message,
			},
		})
		if err != nil {
			slog.Error(
				"failed to dispatch new message",
				"error", err,
				"user_id", user.UserID,
				"chat_id", user.UserID,
				"message_id", m.ID,
				"message", m.Message,
			)
			return nil
		}

		slog.Info(
			"dispatched new message",
			"user_id", user.UserID,
			"chat_id", user.UserID,
			"message_id", m.ID,
			"message", m.Message,
		)
		return nil
	})
}
