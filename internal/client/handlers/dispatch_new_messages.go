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

		err := bus.Dispatch(ctx, pubsub.Event{
			Type: "new_message",
			Payload: map[string]any{
				"message": m.Message,
				"chat_id": m.PeerID.TypeID(),
			},
		})
		if err != nil {
			slog.Error("failed to dispatch new message", "error", err, "chat_id", m.PeerID.TypeID(), "message_id", m.ID)
			return nil
		}

		slog.Info("dispatched new message", "chat_id", m.PeerID.TypeID(), "message_id", m.ID, "message", m.Message)
		return nil
	})
}
