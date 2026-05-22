// Package tinymcp provides a lightweight wrapper around the official Model Context
// Protocol Go SDK. It reduces boilerplate for building MCP servers (stdio, streamable
// HTTP, or legacy SSE) that AI clients can discover and call.
package tinymcp

import (
	"context"
	"errors"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TinyServer wraps the official MCP Server with a simplified API for registering
// tools and starting stdio transport.
type TinyServer struct {
	server *mcp.Server
}

// NewServer creates a TinyServer identified by name and version for connected clients.
func NewServer(name, version string) *TinyServer {
	s := mcp.NewServer(
		&mcp.Implementation{Name: name, Version: version},
		nil,
	)
	return &TinyServer{server: s}
}

// RegisterTool registers a tool on s with automatic JSON Schema inference from the
// handler's input struct (use json and jsonschema struct tags). The handler must
// match mcp.ToolHandlerFor[In, Out] — typically return TextResult for text tools.
func RegisterTool[In, Out any](s *TinyServer, name, description string, handler mcp.ToolHandlerFor[In, Out]) error {
	if s == nil || s.server == nil {
		return errors.New("cannot register tool on a nil server")
	}
	mcp.AddTool(s.server, &mcp.Tool{
		Name:        name,
		Description: description,
	}, handler)
	return nil
}

// Start runs the server over stdin/stdout (stdio), the standard transport for local AI clients.
func (s *TinyServer) Start() error {
	if s == nil || s.server == nil {
		return errors.New("cannot start a nil server")
	}
	return s.server.Run(context.Background(), &mcp.StdioTransport{})
}

// RawServer returns the underlying *mcp.Server for advanced customization.
func (s *TinyServer) RawServer() *mcp.Server {
	if s == nil {
		return nil
	}
	return s.server
}

// TextResult builds a successful CallToolResult with a single text content block.
// Use this in tool handlers so LLM clients receive consistent, readable output.
func TextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: text}},
	}
}
