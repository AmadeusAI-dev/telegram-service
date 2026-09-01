package server

import (
	"net/http"

	"github.com/AmadeusAI-dev/telegram-service/internal/mcp/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func New(sender tools.Sender) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "telegram-mcp", Version: "1.0.0"}, nil)

	registerTools(server, sender)

	return server
}

func registerTools(server *mcp.Server, sender tools.Sender) {
	mcp.AddTool(server, tools.SendMessageToolInfo(), tools.SendMessageTool(sender))
}

func Run(ch chan<- error, url string, server *mcp.Server) {
	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	go func() {
		err := http.ListenAndServe(url, handler)
		if err != nil {
			ch <- err
		}
	}()
}
