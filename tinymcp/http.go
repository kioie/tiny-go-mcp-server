package tinymcp

import (
	"errors"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTPOptions configures streamable HTTP transport (current MCP spec: POST + SSE).
type HTTPOptions struct {
	// Stateless uses a fresh session per request (no Mcp-Session-Id). Simpler for
	// stateless gateways; server-initiated requests are not supported.
	Stateless bool
	// JSONResponse returns application/json instead of text/event-stream for POST responses.
	JSONResponse bool
	// SessionTimeout closes idle sessions after this duration (zero = never).
	SessionTimeout time.Duration
	// DisableLocalhostProtection turns off DNS rebinding checks for localhost clients.
	// Only set when you understand the security tradeoff.
	DisableLocalhostProtection bool
}

// SSEOptions configures legacy SSE transport (MCP spec 2024-11-05).
type SSEOptions struct {
	// DisableLocalhostProtection turns off DNS rebinding checks for localhost clients.
	DisableLocalhostProtection bool
}

func (o *HTTPOptions) streamableOpts() *mcp.StreamableHTTPOptions {
	if o == nil {
		return nil
	}
	return &mcp.StreamableHTTPOptions{
		Stateless:                  o.Stateless,
		JSONResponse:               o.JSONResponse,
		SessionTimeout:             o.SessionTimeout,
		DisableLocalhostProtection: o.DisableLocalhostProtection,
	}
}

func (o *SSEOptions) sseOpts() *mcp.SSEOptions {
	if o == nil {
		return nil
	}
	return &mcp.SSEOptions{
		DisableLocalhostProtection: o.DisableLocalhostProtection,
	}
}

func rawServer(s *TinyServer) (*mcp.Server, error) {
	if s == nil || s.server == nil {
		return nil, errors.New("cannot use a nil server")
	}
	return s.server, nil
}

// StreamableHTTPHandler returns an http.Handler for streamable HTTP MCP clients
// (POST JSON-RPC; responses as JSON or SSE). Wrap with auth/CORS middleware as needed.
func StreamableHTTPHandler(s *TinyServer, opts *HTTPOptions) (http.Handler, error) {
	srv, err := rawServer(s)
	if err != nil {
		return nil, err
	}
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return srv
	}, opts.streamableOpts()), nil
}

// SSEHandler returns an http.Handler for legacy SSE MCP clients (2024-11-05 transport).
func SSEHandler(s *TinyServer, opts *SSEOptions) (http.Handler, error) {
	srv, err := rawServer(s)
	if err != nil {
		return nil, err
	}
	return mcp.NewSSEHandler(func(*http.Request) *mcp.Server {
		return srv
	}, opts.sseOpts()), nil
}

// StartHTTP listens on addr and serves streamable HTTP MCP until the server stops or errors.
func (s *TinyServer) StartHTTP(addr string, opts *HTTPOptions) error {
	h, err := StreamableHTTPHandler(s, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, h)
}

// StartSSE listens on addr and serves legacy SSE MCP until the server stops or errors.
func (s *TinyServer) StartSSE(addr string, opts *SSEOptions) error {
	h, err := SSEHandler(s, opts)
	if err != nil {
		return err
	}
	return http.ListenAndServe(addr, h)
}
