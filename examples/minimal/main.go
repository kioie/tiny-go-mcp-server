// Minimal example: a Tiny Go MCP server with one tool.
//
// Run: go run ./examples/minimal
// Or install: go install github.com/kioie/mcp-server/examples/minimal@latest
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kioie/mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type helloArgs struct {
	Name string `json:"name" jsonschema:"Person to greet (required)"`
}

func main() {
	server := tinymcp.NewServer("tiny-go-mcp-example", "1.0.0")

	if err := tinymcp.RegisterTool(server, "hello", "Say hello to someone", hello); err != nil {
		log.Fatal(err)
	}

	log.Println("Tiny Go MCP example server (stdio)")
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
