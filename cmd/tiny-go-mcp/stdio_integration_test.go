package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestStdioRoundTrip_addTool builds the reference server and calls the add tool
// over a real stdio subprocess (newline-delimited JSON framing).
func TestStdioRoundTrip_addTool(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stdio subprocess integration test in short mode")
	}

	ctx := t.Context()
	bin := filepath.Join(t.TempDir(), "tiny-go-mcp")
	build := exec.Command("go", "build", "-o", bin, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "TINY_GO_MCP_VERBOSE=") // quiet stderr

	client := mcp.NewClient(&mcp.Implementation{Name: "integration-test", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer session.Close()

	result, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "add",
		Arguments: map[string]any{"a": 2, "b": 3},
	})
	if err != nil {
		t.Fatalf("CallTool add: %v", err)
	}
	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}
	tc, ok := result.Content[0].(*mcp.TextContent)
	expected := "Result of addition: 2 + 3 = 5"
	if !ok || tc.Text != expected {
		t.Fatalf("add(2,3) = %q, want %q", tc.Text, expected)
	}
}
