package tools

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var ctrCmpOpts = []cmp.Option{
	cmpopts.IgnoreUnexported(mcp.CallToolResult{}, mcp.GetPromptResult{}, mcp.ReadResourceResult{}),
	cmpopts.IgnoreFields(mcp.CallToolResult{}, "Meta"),
	cmpopts.IgnoreFields(mcp.GetPromptResult{}, "Meta"),
	cmpopts.IgnoreFields(mcp.ReadResourceResult{}, "Meta"),
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

func callTool(t testing.TB, cs *mcp.ClientSession, toolName string, args map[string]any) *mcp.CallToolResult {
	t.Helper()

	ctx := context.Background()
	params := &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	}

	got, err := cs.CallTool(ctx, params)
	if err != nil {
		t.Fatalf("error while calling tool: %v", err)
	}

	return got
}

func assertCallToolResultsMatch(t testing.TB, got, want *mcp.CallToolResult) {
	t.Helper()

	if diff := cmp.Diff(want, got, ctrCmpOpts...); diff != "" {
		t.Errorf("tools/call mismatch (-want +got):\n%s", diff)
	}
}
