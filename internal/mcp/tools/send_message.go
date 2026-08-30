package tools

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
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
			slog.Error("failed to sent message", "err", err)
		}
		return nil, Output{Result: "message sent successfully"}, nil
	}
}
