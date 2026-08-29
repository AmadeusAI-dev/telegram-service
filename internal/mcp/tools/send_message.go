package tools

import (
	"context"
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type Sender interface {
	Send(context.Context, int, string) error
}

type Input struct {
	ChatID  int    `json:"chat_id" jsonschema:"the id of the chat where message should be sent"`
	Message string `json:"message" jsonschema:"the message to be sent"`
}

type Output struct {
	Result string `json:"result" jsonschema:"result message"`
}

func SendMessageTool(sender Sender) mcp.ToolHandlerFor[Input, Output] {
	return func(ctx context.Context, req *mcp.CallToolRequest, input Input) (
		*mcp.CallToolResult,
		Output,
		error,
	) {
		err := sender.Send(ctx, input.ChatID, input.Message)
		if err != nil {
			slog.Error("failed to sent message", "err", err)
		}
		return nil, Output{Result: "message sent successfully"}, nil
	}
}
