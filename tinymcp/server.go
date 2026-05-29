// Package tinymcp provides a lightweight wrapper around the official Model Context
// Protocol Go SDK. It reduces boilerplate for building MCP servers (stdio, streamable
// HTTP, or legacy SSE) that AI clients can discover and call.
package tinymcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TinyServer wraps the official MCP Server with a simplified API for registering
// tools, resources, prompts, and starting transports.
type TinyServer struct {
	server *mcp.Server
}

// NewServer creates a TinyServer identified by name and version for connected clients.
// Optional server configuration uses functional options or [NewServerWithOptions].
func NewServer(name, version string, opts ...ServerOption) *TinyServer {
	return newServer(name, version, opts...)
}

// RegisterTool registers a tool on s with automatic JSON Schema inference from the
// handler's input struct (use json and jsonschema struct tags). The handler must
// match mcp.ToolHandlerFor[In, Out] — typically return TextResult for text tools.
// Invalid tool definitions surface as returned errors instead of panicking.
func RegisterTool[In, Out any](s *TinyServer, name, description string, handler mcp.ToolHandlerFor[In, Out]) error {
	if s == nil || s.server == nil {
		return ErrNilServer
	}
	return RegisterToolDef(s, &mcp.Tool{
		Name:        name,
		Description: description,
	}, handler)
}

// MustRegisterTool is like RegisterTool but panics on error. Safe at process startup.
func MustRegisterTool[In, Out any](s *TinyServer, name, description string, handler mcp.ToolHandlerFor[In, Out]) {
	if err := RegisterTool(s, name, description, handler); err != nil {
		panic(err)
	}
}

// RegisterToolDef registers a tool with full go-sdk Tool metadata (annotations, etc.).
// Invalid tool definitions surface as returned errors instead of panicking.
func RegisterToolDef[In, Out any](s *TinyServer, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, Out]) error {
	if s == nil || s.server == nil {
		return ErrNilServer
	}
	if tool == nil {
		return ErrNilTool
	}
	name := tool.Name
	if name == "" {
		name = "<unnamed>"
	}
	return registerRecover(fmt.Sprintf("tool %q", name), func() {
		mcp.AddTool(s.server, tool, handler)
	})
}

// MustRegisterToolDef is like RegisterToolDef but panics on error. Safe at process startup.
func MustRegisterToolDef[In, Out any](s *TinyServer, tool *mcp.Tool, handler mcp.ToolHandlerFor[In, Out]) {
	if err := RegisterToolDef(s, tool, handler); err != nil {
		panic(err)
	}
}

// Start runs the server over stdin/stdout (stdio), the standard transport for local AI clients.
func (s *TinyServer) Start() error {
	if s == nil || s.server == nil {
		return ErrNilServer
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
