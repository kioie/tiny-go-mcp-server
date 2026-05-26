// Streamable HTTP MCP example (POST + SSE per current MCP spec).
//
// Run: go run ./examples/http
//      TINY_GO_MCP_ADDR=0.0.0.0:8080 go run ./examples/http   # listen on all interfaces
//      TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1 go run ./examples/http   # ngrok/tunnel testing
package main

import (
	"context"
	"log"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type pingArgs struct {
	Message string `json:"message" jsonschema:"Message to echo back"`
}

func main() {
	addr := envOr("TINY_GO_MCP_ADDR", "127.0.0.1:8080")
	server := tinymcp.NewServer("tiny-go-mcp-http", "1.0.0")

	if err := tinymcp.RegisterTool(server, "ping", "Echo a message over HTTP MCP", ping); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Printf("Tiny Go MCP HTTP example listening on %s", addr)
	}

	// Stateless=true is a simple default for demos; use nil or Stateful opts for full sessions.
	if err := server.StartHTTP(addr, &tinymcp.HTTPOptions{
		Stateless:                  true,
		DisableLocalhostProtection: tinymcp.DisableLocalhostProtectionFromEnv(),
	}); err != nil {
		log.Fatal(err)
	}
}

func ping(_ context.Context, _ *mcp.CallToolRequest, args pingArgs) (*mcp.CallToolResult, any, error) {
	msg := args.Message
	if msg == "" {
		msg = "pong"
	}
	return tinymcp.TextResult(msg), nil, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
