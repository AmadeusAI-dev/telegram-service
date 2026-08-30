package tools

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gotd/td/tg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var (
	UserNameNotFound    = errors.New("failed to send message: user not found")
	FailedToSendMessage = errors.New("failed to send message: please try again later")
)

type Sender interface {
	Send(context.Context, string, string) error
}

type Input struct {
	Username string `json:"username" jsonschema:"the username of the message recipient"`
	Message  string `json:"message" jsonschema:"the message to be sent"`
}

type Output struct {
	Result string `json:"result" jsonschema:"result message"`
}

func SendMessageToolInfo() *mcp.Tool {
	return &mcp.Tool{
		Name:        "send_message",
		Description: "sends message to specific telegram user, based on the username",
	}
}

func SendMessageTool(sender Sender) mcp.ToolHandlerFor[Input, Output] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input Input) (
		*mcp.CallToolResult,
		Output,
		error,
	) {
		err := sender.Send(ctx, input.Username, input.Message)
		if err != nil {
			switch {
			case tg.IsUsernameNotOccupied(err), tg.IsUsernameInvalid(err):
				return nil, Output{Result: "failed to sent message"}, UserNameNotFound
			}

			slog.Error(
				"failed to send message",
				"error", err,
				"username", input.Username,
				"message", input.Message,
			)
			return nil, Output{Result: "failed to sent message"}, FailedToSendMessage
		}

		return nil, Output{Result: "message sent successfully"}, nil
	}
}
