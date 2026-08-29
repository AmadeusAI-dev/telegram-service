package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/AmadeusAI-dev/telegram-service/internal/client"
	"github.com/AmadeusAI-dev/telegram-service/internal/client/handlers"
	"github.com/AmadeusAI-dev/telegram-service/internal/config"
	"github.com/AmadeusAI-dev/telegram-service/internal/mcp/server"
	"github.com/AmadeusAI-dev/telegram-service/internal/mcp/tools"
	"github.com/TheKiryuKha/pubsub"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type App struct {
	pubsub *pubsub.Pubsub

	client       *telegram.Client
	peersManager *peers.Manager
	gaps         *updates.Manager

	mcpHandler *mcp.StreamableHTTPHandler
	url        string
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// @todo: recovery engine
	d := tg.NewUpdateDispatcher()
	var h telegram.UpdateHandler

	tgClient, err := telegram.ClientFromEnvironment(telegram.Options{
		UpdateHandler: telegram.UpdateHandlerFunc(func(ctx context.Context, u tg.UpdatesClass) error {
			return h.Handle(ctx, u)
		}),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	peerManager := peers.Options{}.Build(tgClient.API())
	gaps := updates.New(updates.Config{
		Handler:      peerManager.UpdateHook(d),
		AccessHasher: peerManager,
	})
	h = gaps

	server := mcp.NewServer(&mcp.Implementation{Name: "telegram-mcp", Version: "1.0.0"}, nil)

	sender := &client.Sender{Client: tgClient}
	mcp.AddTool(server, &mcp.Tool{
		Name:        "send_message",
		Description: "sends message to specific telegram chat, based on the chat_id",
	}, tools.SendMessageTool(sender))

	handler := mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return server
	}, nil)

	url := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	bus, err := pubsub.New(ctx, cfg.RabbitMq)
	if err != nil {
		return nil, fmt.Errorf("failed to init pubsub: %w", err)
	}

	userRepo := client.NewUserRepo(peerManager)
	handlers.DispatchNewMessages(bus, userRepo, &d)

	return &App{
		pubsub:       bus,
		client:       tgClient,
		mcpHandler:   handler,
		url:          url,
		peersManager: peerManager,
		gaps:         gaps,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	runTgCh := make(chan error, 1)
	runMcpCh := make(chan error, 1)

	initTgCh := client.Initialize(ctx, a.client, a.peersManager, a.gaps, runTgCh)
	if err := client.WaitForInitialization(ctx, initTgCh); err != nil {
		return err
	}

	server.Run(runMcpCh, a.url, a.mcpHandler)

	slog.Info("application started successfully")

	select {
	case err := <-runTgCh:
		return fmt.Errorf("telegram client stopped: %w", err)

	case err := <-runMcpCh:
		return fmt.Errorf("mcp http server stopped: %w", err)

	case <-ctx.Done():
		return nil
	}
}

func (a *App) Close(ctx context.Context) error {
	return a.pubsub.Close(ctx)
}
