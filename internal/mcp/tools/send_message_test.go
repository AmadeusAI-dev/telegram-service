package tools

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/gotd/td/tgerr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SpySender struct {
	calls    int
	username string
	message  string
}

func (s *SpySender) Send(ctx context.Context, username string, message string) error {
	s.calls++
	s.message = message
	s.username = username
	return nil
}

var (
	UserNameIsNotOcupied = tgerr.New(400, "USERNAME_NOT_OCCUPIED")
	UserNameIsInvalid    = tgerr.New(400, "USERNAME_INVALID")
	UnkownError          = errors.New("sender: unkown error")
)

type ErrorSender struct {
	errorToSend error
}

func (e *ErrorSender) Send(ctx context.Context, username string, message string) error {
	return e.errorToSend
}

func TestSendsMessage(t *testing.T) {
	server, cs := NewTestMcp(t)
	sender := &SpySender{}

	mcp.AddTool(server, SendMessageToolInfo(), SendMessageTool(sender))

	got := callTool(t, cs, "send_message", map[string]any{
		"username": "TheKiryuKha",
		"message":  "Hello, World!",
	})

	want := &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: `{"result":"message sent successfully"}`},
		},
		StructuredContent: map[string]any{"result": "message sent successfully"},
	}

	assertCallToolResultsMatch(t, got, want)

	if sender.calls != 1 {
		t.Fatalf("expected sender to be called %d, got %d", 1, sender.calls)
	}
	if sender.username != "TheKiryuKha" {
		t.Fatalf("expected chat_id to be %s, got %s", "TheKiryuKha", sender.username)
	}
	if sender.message != "Hello, World!" {
		t.Fatalf("expected sent message to be '%s', got '%s'", "Hello, World!", sender.message)
	}
}

func TestReturnsErrorsFromSender(t *testing.T) {
	tests := map[string]struct {
		senderError error
		mcpError    error
	}{
		"UserNameIsNotOcupied": {
			UserNameIsNotOcupied,
			UserNameNotFound,
		},
		"UserNameInvalid": {
			UserNameIsInvalid,
			UserNameNotFound,
		},
		"UnkownError": {
			UnkownError,
			FailedToSendMessage,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			server, cs := NewTestMcp(t)
			sender := &ErrorSender{test.senderError}

			mcp.AddTool(server, SendMessageToolInfo(), SendMessageTool(sender))

			got := callTool(t, cs, "send_message", map[string]any{
				"username": "TheKiryuKha",
				"message":  "Hi!",
			})

			want := &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: test.mcpError.Error(),
					},
				},
				IsError: true,
			}

			assertCallToolResultsMatch(t, got, want)
		})
	}

}

func TestInfo(t *testing.T) {
	want := &mcp.Tool{
		Name:        "send_message",
		Description: "sends message to specific telegram user, based on the username",
	}

	got := SendMessageToolInfo()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("mcp tool info mismatch. got: %v, want: %v", got, want)
	}
}
