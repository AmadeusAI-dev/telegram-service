package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/tg"
)

type SpyBus struct {
	calls  int
	events []pubsub.Event
}

func (s *SpyBus) Dispatch(ctx context.Context, event pubsub.Event) error {
	s.calls++
	s.events = append(s.events, event)
	return nil
}

var expectedEvent = pubsub.Event{
	Type: "new_message",
	Payload: map[string]any{
		"user_id":    int64(456),
		"chat_id":    int64(456),
		"message_id": 123,
		"message":    "Hello, World!",
	},
}

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		ownMessage bool
		wantCalls  int
		wantEvents []pubsub.Event
	}{
		"new message": {
			ownMessage: false,
			wantCalls:  1,
			wantEvents: []pubsub.Event{expectedEvent},
		},
		"new own message": {
			ownMessage: true,
			wantCalls:  0,
			wantEvents: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := tg.NewUpdateDispatcher()
			bus := &SpyBus{}

			DispatchNewMessages(bus, &d)
			err := sentMessage(t, &d, test.ownMessage)
			if err != nil {
				t.Fatalf("expected no error while sending test message, got %v", err)
			}

			assertBusWork(t, bus, test.wantCalls, test.wantEvents)
		})
	}
}

func sentMessage(t testing.TB, d *tg.UpdateDispatcher, ownMessage bool) error {
	t.Helper()

	return d.Handle(context.Background(), &tg.UpdateShort{
		Update: &tg.UpdateNewMessage{
			Message: &tg.Message{
				ID:      123,
				Message: "Hello, World!",
				Out:     ownMessage,
				FromID:  &tg.PeerUser{UserID: 456},
			},
		},
	})
}

func assertBusWork(t testing.TB, bus *SpyBus, wantCalls int, wantEvents []pubsub.Event) {
	t.Helper()

	if bus.calls != wantCalls {
		t.Fatalf("expected %d calls, got %d", wantCalls, bus.calls)
	}

	if !reflect.DeepEqual(bus.events, wantEvents) {
		t.Fatalf("failed to assert that events are equal, got %#v, want %#v", bus.events, wantEvents)
	}
}
