package tinymcp

import (
	"context"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestSentinelErrors(t *testing.T) {
	if ErrNilServer == nil || ErrNilTool == nil || ErrNilHandler == nil || ErrRegistrationFailed == nil {
		t.Fatal("sentinel errors must be non-nil")
	}
}

func TestRegisterTool_nilServer(t *testing.T) {
	var s *TinyServer
	err := RegisterTool(s, "x", "y", func(_ context.Context, _ *mcp.CallToolRequest, _ echoArgs) (*mcp.CallToolResult, any, error) {
		return nil, nil, nil
	})
	if !errors.Is(err, ErrNilServer) {
		t.Fatalf("RegisterTool(nil): got %v, want ErrNilServer", err)
	}
}

func TestRegisterToolDef_nilTool(t *testing.T) {
	s := NewServer("test", "0.0.1")
	err := RegisterToolDef(s, nil, func(_ context.Context, _ *mcp.CallToolRequest, _ echoArgs) (*mcp.CallToolResult, any, error) {
		return nil, nil, nil
	})
	if !errors.Is(err, ErrNilTool) {
		t.Fatalf("RegisterToolDef(nil tool): got %v, want ErrNilTool", err)
	}
}

func TestRegisterToolDef_invalidInput(t *testing.T) {
	s := NewServer("test", "0.0.1")
	err := RegisterToolDef(s, &mcp.Tool{Name: "bad", Description: "x"},
		func(_ context.Context, _ *mcp.CallToolRequest, args string) (*mcp.CallToolResult, any, error) {
			return nil, nil, nil
		})
	if !errors.Is(err, ErrRegistrationFailed) {
		t.Fatalf("RegisterToolDef(invalid): got %v, want ErrRegistrationFailed", err)
	}
}

func TestMustRegisterTool_panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("MustRegisterTool(nil server): expected panic")
		}
	}()
	MustRegisterTool(nil, "x", "y", func(_ context.Context, _ *mcp.CallToolRequest, _ echoArgs) (*mcp.CallToolResult, any, error) {
		return nil, nil, nil
	})
}
