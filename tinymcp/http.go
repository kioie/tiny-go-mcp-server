package tinymcp

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const defaultShutdownTimeout = 30 * time.Second

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
	// Middleware wraps the MCP handler (first entry is innermost). Use for logging,
	// rate limits, etc. Apply auth/CORS outside via StreamableHTTPHandler return value
	// or additional wrapping on your mux.
	Middleware []func(http.Handler) http.Handler
}

// WithMiddleware appends middleware wrappers to opts and returns opts for chaining.
func (o *HTTPOptions) WithMiddleware(mw ...func(http.Handler) http.Handler) *HTTPOptions {
	if o == nil {
		o = &HTTPOptions{}
	}
	o.Middleware = append(o.Middleware, mw...)
	return o
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
		return nil, ErrNilServer
	}
	return s.server, nil
}

func applyMiddleware(h http.Handler, middleware []func(http.Handler) http.Handler) http.Handler {
	for _, mw := range middleware {
		h = mw(h)
	}
	return h
}

// StreamableHTTPHandler returns an http.Handler for streamable HTTP MCP clients
// (POST JSON-RPC; responses as JSON or SSE). Wrap with auth/CORS middleware as needed.
func StreamableHTTPHandler(s *TinyServer, opts *HTTPOptions) (http.Handler, error) {
	srv, err := rawServer(s)
	if err != nil {
		return nil, err
	}
	var h http.Handler = mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return srv
	}, opts.streamableOpts())
	if opts != nil {
		h = applyMiddleware(h, opts.Middleware)
	}
	return h, nil
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

func newHTTPServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}
}

// ListenAndServeHTTP starts addr with handler using timeouts suited to MCP HTTP (SSE-friendly).
// It shuts down gracefully on SIGINT or SIGTERM (for container deploys).
func ListenAndServeHTTP(addr string, handler http.Handler) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := ListenAndServeHTTPContext(ctx, addr, handler)
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

// ListenAndServeHTTPContext serves until ctx is canceled, then drains active connections.
func ListenAndServeHTTPContext(ctx context.Context, addr string, handler http.Handler) error {
	if ctx == nil {
		ctx = context.Background()
	}
	srv := newHTTPServer(addr, handler)

	errCh := make(chan error, 1)
	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			errCh <- nil
			return
		}
		errCh <- err
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		if err := <-errCh; err != nil {
			return err
		}
		return ctx.Err()
	}
}

// StartHTTP listens on addr and serves streamable HTTP MCP until the server stops or errors.
func (s *TinyServer) StartHTTP(addr string, opts *HTTPOptions) error {
	h, err := StreamableHTTPHandler(s, opts)
	if err != nil {
		return err
	}
	return ListenAndServeHTTP(addr, h)
}

// StartHTTPContext serves streamable HTTP until ctx is canceled.
func (s *TinyServer) StartHTTPContext(ctx context.Context, addr string, opts *HTTPOptions) error {
	h, err := StreamableHTTPHandler(s, opts)
	if err != nil {
		return err
	}
	return ListenAndServeHTTPContext(ctx, addr, h)
}

// StartSSE listens on addr and serves legacy SSE MCP until the server stops or errors.
func (s *TinyServer) StartSSE(addr string, opts *SSEOptions) error {
	h, err := SSEHandler(s, opts)
	if err != nil {
		return err
	}
	return ListenAndServeHTTP(addr, h)
}

// StartSSEContext serves legacy SSE until ctx is canceled.
func (s *TinyServer) StartSSEContext(ctx context.Context, addr string, opts *SSEOptions) error {
	h, err := SSEHandler(s, opts)
	if err != nil {
		return err
	}
	return ListenAndServeHTTPContext(ctx, addr, h)
}
