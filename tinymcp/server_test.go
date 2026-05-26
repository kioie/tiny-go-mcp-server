package tinymcp

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type echoArgs struct {
	Message string `json:"message" jsonschema:"Text to echo back"`
}

func TestNewServer_RegisterTool(t *testing.T) {
	s := NewServer("test", "0.0.1")
	err := RegisterTool(s, "echo", "Echoes input", func(_ context.Context, _ *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, any, error) {
		return TextResult(args.Message), nil, nil
	})
	if err != nil {
		t.Fatalf("RegisterTool: %v", err)
	}
	if s.RawServer() == nil {
		t.Fatal("expected non-nil RawServer")
	}
}

func TestRegisterTool_nilServer(t *testing.T) {
	var s *TinyServer
	err := RegisterTool(s, "x", "y", func(_ context.Context, _ *mcp.CallToolRequest, _ echoArgs) (*mcp.CallToolResult, any, error) {
		return nil, nil, nil
	})
	if err == nil {
		t.Fatal("expected error for nil server")
	}
}

func TestRegisterToolDef_invalidInput(t *testing.T) {
	s := NewServer("test", "0.0.1")
	err := RegisterToolDef(s, &mcp.Tool{Name: "bad", Description: "x"},
		func(_ context.Context, _ *mcp.CallToolRequest, args string) (*mcp.CallToolResult, any, error) {
			return nil, nil, nil
		})
	if err == nil {
		t.Fatal("expected error for invalid tool input type")
	}
}

func TestTextResult(t *testing.T) {
	r := TextResult("hello")
	if len(r.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(r.Content))
	}
	tc, ok := r.Content[0].(*mcp.TextContent)
	if !ok || tc.Text != "hello" {
		t.Fatalf("unexpected content: %#v", r.Content[0])
	}
}
