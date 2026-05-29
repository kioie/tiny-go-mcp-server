// HTTP deploy scaffold for tinymcp (Smithery URL listing, Fly/Render/Railway).
//
// Create a project:
//
//	gonew github.com/kioie/tiny-go-mcp-server/template-http@latest example.com/my-mcp-http my-mcp-http
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type pingArgs struct {
	Message string `json:"message" jsonschema:"Message to echo back (defaults to pong when empty)"`
}

func main() {
	name := envOr("MCP_SERVER_NAME", "my-mcp-http")
	version := envOr("MCP_SERVER_VERSION", "0.1.0")
	addr := listenAddr()

	server := tinymcp.NewServer(name, version)
	if err := registerTools(server); err != nil {
		log.Fatal(err)
	}

	mcpHandler, err := tinymcp.StreamableHTTPHandler(server, &tinymcp.HTTPOptions{
		Stateless:                  true,
		DisableLocalhostProtection: tinymcp.DisableLocalhostProtectionFromEnv(),
	})
	if err != nil {
		log.Fatal(err)
	}
	mcpHandler = tinymcp.WithCrossOriginProtection(mcpHandler)
	mcpHandler = tinymcp.BearerTokenAuth(os.Getenv("TINY_GO_MCP_API_KEY"), mcpHandler)

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	mux.Handle("/", mcpHandler)

	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Printf("Starting %s %s on %s (streamable HTTP)", name, version, addr)
	}
	if err := tinymcp.ListenAndServeHTTP(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func registerTools(server *tinymcp.TinyServer) error {
	return tinymcp.RegisterTool(server, "ping",
		"Demo: echo a message over HTTP MCP. Use only to verify HTTP transport and deployment wiring; do not use for production messaging—reply in chat instead. Use ping (not other tools) when testing connectivity.",
		ping,
	)
}

func ping(_ context.Context, _ *mcp.CallToolRequest, args pingArgs) (*mcp.CallToolResult, any, error) {
	msg := args.Message
	if msg == "" {
		msg = "pong"
	}
	return tinymcp.TextResult(msg), nil, nil
}

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	if v := os.Getenv("TINY_GO_MCP_ADDR"); v != "" {
		return v
	}
	return "127.0.0.1:8080"
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
