package tinymcp

import (
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestNewServer_defaultOptions(t *testing.T) {
	s := NewServer("test", "1.0.0")
	if s == nil || s.RawServer() == nil {
		t.Fatal("expected server")
	}
}

func TestNewServer_withInstructions(t *testing.T) {
	ctx := t.Context()
	const instructions = "Use ping only for connectivity tests."

	s := NewServer("test", "1.0.0", WithInstructions(instructions))

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	if got := cs.InitializeResult().Instructions; got != instructions {
		t.Fatalf("Instructions = %q, want %q", got, instructions)
	}
}

func TestNewServerWithOptions(t *testing.T) {
	ctx := t.Context()
	const instructions = "Configured via NewServerWithOptions."

	s := NewServerWithOptions("test", "1.0.0", &mcp.ServerOptions{
		Instructions: instructions,
	})

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	if got := cs.InitializeResult().Instructions; got != instructions {
		t.Fatalf("Instructions = %q, want %q", got, instructions)
	}
}
