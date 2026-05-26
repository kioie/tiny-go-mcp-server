package tinymcp

import (
	"log/slog"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ServerOption configures optional [mcp.ServerOptions] passed to the underlying SDK.
type ServerOption func(*mcp.ServerOptions)

// NewServerWithOptions creates a TinyServer with explicit go-sdk server options.
// Pass nil for SDK defaults (same as [NewServer] without options).
func NewServerWithOptions(name, version string, options *mcp.ServerOptions) *TinyServer {
	if options == nil {
		return NewServer(name, version)
	}
	return newServer(name, version, WithSDKOptions(options))
}

// WithSDKOptions applies a full *mcp.ServerOptions value. Later options override
// overlapping fields set by earlier ones.
func WithSDKOptions(options *mcp.ServerOptions) ServerOption {
	return func(o *mcp.ServerOptions) {
		if options == nil {
			return
		}
		*o = *options
	}
}

// WithInstructions sets client-facing server instructions in the initialize result.
func WithInstructions(instructions string) ServerOption {
	return func(o *mcp.ServerOptions) {
		o.Instructions = instructions
	}
}

// WithLogger enables structured logging of server activity via the go-sdk.
func WithLogger(logger *slog.Logger) ServerOption {
	return func(o *mcp.ServerOptions) {
		o.Logger = logger
	}
}

// WithSchemaCache enables JSON schema caching (useful for stateless HTTP handlers).
func WithSchemaCache(cache *mcp.SchemaCache) ServerOption {
	return func(o *mcp.ServerOptions) {
		o.SchemaCache = cache
	}
}

func newServer(name, version string, opts ...ServerOption) *TinyServer {
	impl := &mcp.Implementation{Name: name, Version: version}
	if len(opts) == 0 {
		return &TinyServer{server: mcp.NewServer(impl, nil)}
	}

	sdkOpts := &mcp.ServerOptions{}
	for _, opt := range opts {
		opt(sdkOpts)
	}
	return &TinyServer{server: mcp.NewServer(impl, sdkOpts)}
}
