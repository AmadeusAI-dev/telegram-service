package handlers

import (
	"context"
	"reflect"
	"testing"

	"github.com/AmadeusAI-dev/telegram-service/internal/client"
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

type SpyRepo struct {
	calls int
}

func (s *SpyRepo) Get(ctx context.Context, user_id int) (client.User, error) {
	s.calls++
	return client.User{ID: user_id, Username: "mr_TheKiryuKha"}, nil
}

var expectedEvent = pubsub.Event{
	Type: "new_message",
	Payload: map[string]any{
		"user_id":    456,
		"username":   "mr_TheKiryuKha",
		"message_id": 123,
		"message":    "Hello, World!",
	},
}

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		ownMessage    bool
		wantBusCalls  int
		wantEvents    []pubsub.Event
		wantRepoCalls int
	}{
		"new message": {
			ownMessage:    false,
			wantBusCalls:  1,
			wantEvents:    []pubsub.Event{expectedEvent},
			wantRepoCalls: 1,
		},
		"new own message": {
			ownMessage:    true,
			wantBusCalls:  0,
			wantEvents:    nil,
			wantRepoCalls: 0,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			d := tg.NewUpdateDispatcher()
			bus := &SpyBus{}
			repo := &SpyRepo{}

			DispatchNewMessages(bus, repo, &d)
			err := sentMessage(t, &d, test.ownMessage)
			if err != nil {
				t.Fatalf("expected no error while sending test message, got %v", err)
			}

			assertBusWork(t, bus, test.wantBusCalls, test.wantEvents)
			assertRepoWork(t, repo, test.wantRepoCalls)
		})
	}
}

func sentMessage(t testing.TB, d *tg.UpdateDispatcher, ownMessage bool) error {
	t.Helper()

	return d.Handle(context.Background(), &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateNewMessage{
				Message: &tg.Message{
					ID:      123,
					Message: "Hello, World!",
					Out:     ownMessage,
					FromID:  &tg.PeerUser{UserID: 456},
					PeerID:  &tg.PeerUser{UserID: 456},
				},
			},
		},
		Users: []tg.UserClass{
			&tg.User{ID: 456, Username: "mr_TheKiryuKha", FirstName: "Femboy"},
		},
	})
}

func assertBusWork(t testing.TB, bus *SpyBus, wantCalls int, wantEvents []pubsub.Event) {
	t.Helper()

	if bus.calls != wantCalls {
		t.Fatalf("expected %d bus calls, got %d", wantCalls, bus.calls)
	}

	if !reflect.DeepEqual(bus.events, wantEvents) {
		t.Fatalf("failed to assert that events are equal, got %#v, want %#v", bus.events, wantEvents)
	}
}

func assertRepoWork(t testing.TB, repo *SpyRepo, wantCalls int) {
	t.Helper()

	if repo.calls != wantCalls {
		t.Fatalf("expected %d repo calls, got %d", wantCalls, repo.calls)
	}
}
