// Deployable streamable HTTP MCP for Smithery URL listing.
//
// End users connect via Smithery — no Docker on their machine. You host this
// binary on any HTTPS platform (Render, Railway, Fly, etc.).
//
// Run locally:
//
//	go run ./examples/http-deploy
//	ngrok http 8080   # optional: public URL for Smithery scan
//
// Publish to Smithery (after HTTPS is live):
//
//	smithery auth login
//	smithery mcp publish "https://YOUR_HOST" -n your-namespace/your-server
//
// See examples/http-deploy/README.md and docs/SMITHERY.md.
package main

import (
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

//go:embed server-card.json
var serverCardJSON []byte

type pingArgs struct {
	Message string `json:"message" jsonschema:"Message to echo back"`
}

func main() {
	addr := listenAddr()
	server := tinymcp.NewServer("tiny-go-mcp-http", "1.0.0")

	if err := tinymcp.RegisterTool(server, "ping", "Echo a message over HTTP MCP", ping); err != nil {
		log.Fatal(err)
	}

	mcpHandler, err := tinymcp.StreamableHTTPHandler(server, &tinymcp.HTTPOptions{
		Stateless:                  true, // simple default for hosted demos; use nil for full sessions
		DisableLocalhostProtection: true, // required behind reverse proxies (Fly, Render, tunnels)
	})
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	mux.Handle("/.well-known/mcp/server-card.json", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(serverCardJSON)
	}))
	mux.Handle("/", mcpHandler)

	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Printf("Tiny Go MCP HTTP deploy example listening on %s (MCP at /)", addr)
	}

	if err := http.ListenAndServe(addr, mux); err != nil {
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

func listenAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}
	if v := os.Getenv("TINY_GO_MCP_ADDR"); v != "" {
		return v
	}
	return ":8080"
}
