// Deployable streamable HTTP MCP for Smithery URL listing.
//
// End users connect via Smithery — no Docker on their machine. You host this
// binary on any HTTPS platform (Render, Railway, Fly, etc.).
//
// Run locally:
//
//	go run ./examples/http-deploy
//	TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1 go run ./examples/http-deploy
//	ngrok http 8080   # optional: set env above for tunnel + Smithery scan
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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	serverVersion   = "1.1.2"
	infoResourceURI = "file:///info"
)

//go:embed server-card.json
var serverCardJSON []byte

type pingArgs struct {
	Message string `json:"message" jsonschema:"Message to echo back (defaults to pong when empty)"`
}

func main() {
	addr := listenAddr()
	server := tinymcp.NewServer("tiny-go-mcp-http", serverVersion)

	if err := registerCapabilities(server); err != nil {
		log.Fatal(err)
	}

	mcpHandler, err := tinymcp.StreamableHTTPHandler(server, &tinymcp.HTTPOptions{
		Stateless:                  true, // simple default for hosted demos; use nil for full sessions
		DisableLocalhostProtection: disableLocalhostProtection(),
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

	if err := tinymcp.ListenAndServeHTTP(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func registerCapabilities(server *tinymcp.TinyServer) error {
	openWorld := false
	destructive := false
	if err := tinymcp.RegisterToolDef(server, &mcp.Tool{
		Name: "ping",
		Description: "Demo: echo a message over streamable HTTP MCP. Use only to verify HTTP transport, " +
			"Smithery connectivity, or client wiring; do not use for production messaging—reply in chat instead. " +
			"Returns the input message or \"pong\" when empty.",
		Annotations: &mcp.ToolAnnotations{
			ReadOnlyHint:    true,
			IdempotentHint:  true,
			OpenWorldHint:   &openWorld,
			DestructiveHint: &destructive,
			Title:           "Ping",
		},
	}, ping); err != nil {
		return err
	}

	if err := tinymcp.RegisterTextResource(server,
		infoResourceURI,
		"info",
		"Demo: static server metadata for HTTP MCP. Use to test resources/read over streamable HTTP; do not use for live config—read env or APIs instead.",
		"text/plain",
		fmt.Sprintf("Tiny Go MCP Server (HTTP) v%s — tinymcp Smithery URL listing example", serverVersion),
	); err != nil {
		return err
	}

	return tinymcp.RegisterPrompt(server,
		"connectivity_check",
		"Demo: workflow prompt to verify MCP prompt support over HTTP. Use for Smithery score and integration tests; do not use for production workflows.",
		[]*mcp.PromptArgument{{Name: "target", Required: false, Description: "Optional label for the check (defaults to server name)"}},
		handleConnectivityPrompt,
	)
}

func ping(_ context.Context, _ *mcp.CallToolRequest, args pingArgs) (*mcp.CallToolResult, any, error) {
	msg := args.Message
	if msg == "" {
		msg = "pong"
	}
	return tinymcp.TextResult(msg), nil, nil
}

func handleConnectivityPrompt(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	target := req.Params.Arguments["target"]
	if target == "" {
		target = "tiny-go-mcp-http"
	}
	return tinymcp.PromptResult("Connectivity check",
		tinymcp.UserPromptMessage(fmt.Sprintf("Confirm MCP connectivity to %s is working.", target)),
	), nil
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

// disableLocalhostProtection opts out of the go-sdk DNS rebinding guard for loopback servers.
// Set TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1 only for local ngrok/tunnel testing.
func disableLocalhostProtection() bool {
	return os.Getenv("TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION") == "1"
}
