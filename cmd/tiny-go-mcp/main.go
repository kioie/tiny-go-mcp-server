package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddArguments defines input parameters for the 'add' tool.
type AddArguments struct {
	A int `json:"a" jsonschema:"The first number to add"`
	B int `json:"b" jsonschema:"The second number to add"`
}

// SubtractArguments defines input parameters for the 'subtract' tool.
type SubtractArguments struct {
	A int `json:"a" jsonschema:"The number to subtract from"`
	B int `json:"b" jsonschema:"The amount to subtract"`
}

// GreetArguments defines input parameters for the 'greet' tool.
type GreetArguments struct {
	Name     string `json:"name" jsonschema:"The name of the person to greet"`
	Greeting string `json:"greeting,omitempty" jsonschema:"An optional custom greeting phrase like 'Welcome'"`
}

func main() {
	// Initialize the Tiny Go MCP Server
	server := tinymcp.NewServer("tiny-go-mcp", "1.1.0")

	// Register tools to the server using the simplified tinymcp library
	if err := RegisterTools(server); err != nil {
		log.Fatalf("Failed to register tools: %v", err)
	}
	if err := RegisterResourcesAndPrompts(server); err != nil {
		log.Fatalf("Failed to register resources/prompts: %v", err)
	}

	// Serve the MCP server over standard input/output (stdio).
	// MCP uses stdin/stdout for the protocol; optional logs go to stderr (set TINY_GO_MCP_VERBOSE=1).
	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Println("Starting Tiny Go MCP server (stdio)")
	}
	if err := server.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal server error: %v\n", err)
		os.Exit(1)
	}
}

// RegisterTools configures and maps all tools to the given TinyServer.
func RegisterTools(server *tinymcp.TinyServer) error {
	if server == nil {
		return errors.New("cannot register tools on a nil server")
	}

	// 1. Add Tool
	if err := tinymcp.RegisterTool(server,
		"add",
		"Demo: adds two integers. Use only to test MCP addition calls; do not use for real math—compute in the agent instead. Use add (not subtract or greet) when verifying addition wiring.",
		handleAdd,
	); err != nil {
		return fmt.Errorf("failed to register 'add' tool: %w", err)
	}

	// 2. Subtract Tool
	if err := tinymcp.RegisterTool(server,
		"subtract",
		"Demo: subtracts B from A. Use only to test MCP subtraction calls; do not use for real math—compute in the agent instead. Use subtract (not add or greet) when verifying subtraction wiring.",
		handleSubtract,
	); err != nil {
		return fmt.Errorf("failed to register 'subtract' tool: %w", err)
	}

	// 3. Greet Tool
	if err := tinymcp.RegisterTool(server,
		"greet",
		"Demo: returns a greeting string. Use only to test MCP text tools; do not use for user-facing messages—reply in chat instead. Use greet (not add or subtract) for hello/welcome integration tests.",
		handleGreet,
	); err != nil {
		return fmt.Errorf("failed to register 'greet' tool: %w", err)
	}

	return nil
}

// handleAdd processes addition tool requests.
func handleAdd(ctx context.Context, req *mcp.CallToolRequest, args AddArguments) (*mcp.CallToolResult, any, error) {
	sum := Add(args.A, args.B)
	resultText := fmt.Sprintf("Result of addition: %d + %d = %d", args.A, args.B, sum)
	return tinymcp.TextResult(resultText), nil, nil
}

// handleSubtract processes subtraction tool requests.
func handleSubtract(ctx context.Context, req *mcp.CallToolRequest, args SubtractArguments) (*mcp.CallToolResult, any, error) {
	diff := Subtract(args.A, args.B)
	resultText := fmt.Sprintf("Result of subtraction: %d - %d = %d", args.A, args.B, diff)
	return tinymcp.TextResult(resultText), nil, nil
}

// handleGreet processes greet tool requests.
func handleGreet(ctx context.Context, req *mcp.CallToolRequest, args GreetArguments) (*mcp.CallToolResult, any, error) {
	if args.Name == "" {
		return nil, nil, errors.New("name argument is required and cannot be empty")
	}
	greeting := Greet(args.Name, args.Greeting)
	return tinymcp.TextResult(greeting), nil, nil
}

// Core Business Logic (isolated for straightforward unit testing)

// Add performs simple integer addition.
func Add(a, b int) int {
	return a + b
}

// Subtract performs integer subtraction.
func Subtract(a, b int) int {
	return a - b
}

// Greet generates a friendly greeting string.
func Greet(name, customGreeting string) string {
	phrase := "Hello"
	if customGreeting != "" {
		phrase = customGreeting
	}
	return fmt.Sprintf("%s, %s!", phrase, name)
}
