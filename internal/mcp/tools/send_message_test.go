package tools

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SpySender struct {
	calls   int
	chat_id int
	message string
}

func (s *SpySender) Send(ctx context.Context, chat_id int, message string) error {
	s.calls++
	s.message = message
	s.chat_id = chat_id
	return nil
}

func TestSendsMessage(t *testing.T) {
	ctx := context.Background()
	server, cs := NewTestMcp(t)
	sender := &SpySender{}

	mcp.AddTool(server, &mcp.Tool{Name: "send_message", Description: "sents message"}, SendMessageTool(sender))

	params := &mcp.CallToolParams{
		Name: "send_message",
		Arguments: map[string]any{
			"chat_id": 123,
			"message": "Hello, World!",
		},
	}
	res, err := cs.CallTool(ctx, params)
	if err != nil {
		t.Fatalf("error while calling tool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool call returned error: %v", res)
	}

	if sender.calls != 1 {
		t.Fatalf("expected sender to be called %d, got %d", 1, sender.calls)
	}
	if sender.chat_id != 123 {
		t.Fatalf("expected chat_id to be %d, got %d", 123, sender.chat_id)
	}
	if sender.message != "Hello, World!" {
		t.Fatalf("expected sent message to be '%s', got '%s'", "Hello, World!", sender.message)
	}
}

func NewTestMcp(t testing.TB) (*mcp.Server, *mcp.ClientSession) {
	t.Helper()

	ctx := context.Background()
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v1.0.0"}, nil)
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "v1.0.0"}, nil)
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatal(err)
	}
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		serverSession.Close()
		clientSession.Close()
	})

	return server, clientSession
}
