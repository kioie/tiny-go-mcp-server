// Scaffold a new tinymcp MCP server.
//
// Create a project: gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type helloArgs struct {
	Name string `json:"name" jsonschema:"Person to greet (required)"`
}

func main() {
	name := envOr("MCP_SERVER_NAME", "my-mcp")
	version := envOr("MCP_SERVER_VERSION", "0.1.0")
	server := tinymcp.NewServer(name, version)

	if err := tinymcp.RegisterTool(server, "hello", "Say hello to someone by name", hello); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Printf("Starting %s %s (stdio)", name, version)
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func hello(_ context.Context, _ *mcp.CallToolRequest, args helloArgs) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return nil, nil, fmt.Errorf("name is required")
	}
	return tinymcp.TextResult(fmt.Sprintf("Hello, %s!", args.Name)), nil, nil
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
