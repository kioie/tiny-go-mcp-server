package tinymcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterResource adds a resource with a fixed URI. uri must be absolute (e.g. file:///docs/readme).
func RegisterResource(s *TinyServer, uri, name, description, mimeType string, handler mcp.ResourceHandler) error {
	srv, err := rawServer(s)
	if err != nil {
		return err
	}
	if handler == nil {
		return fmt.Errorf("%w: resource %q", ErrNilHandler, name)
	}
	return registerRecover(fmt.Sprintf("resource %q", name), func() {
		srv.AddResource(&mcp.Resource{
			URI:         uri,
			Name:        name,
			Description: description,
			MIMEType:    mimeType,
		}, handler)
	})
}

// MustRegisterResource is like RegisterResource but panics on error. Safe at process startup.
func MustRegisterResource(s *TinyServer, uri, name, description, mimeType string, handler mcp.ResourceHandler) {
	if err := RegisterResource(s, uri, name, description, mimeType, handler); err != nil {
		panic(err)
	}
}

// RegisterResourceTemplate adds a resource matched by URI template (e.g. file:///docs/{path}).
func RegisterResourceTemplate(s *TinyServer, uriTemplate, name, description, mimeType string, handler mcp.ResourceHandler) error {
	srv, err := rawServer(s)
	if err != nil {
		return err
	}
	if handler == nil {
		return fmt.Errorf("%w: resource template %q", ErrNilHandler, name)
	}
	return registerRecover(fmt.Sprintf("resource template %q", name), func() {
		srv.AddResourceTemplate(&mcp.ResourceTemplate{
			URITemplate: uriTemplate,
			Name:        name,
			Description: description,
			MIMEType:    mimeType,
		}, handler)
	})
}

// MustRegisterResourceTemplate is like RegisterResourceTemplate but panics on error.
func MustRegisterResourceTemplate(s *TinyServer, uriTemplate, name, description, mimeType string, handler mcp.ResourceHandler) {
	if err := RegisterResourceTemplate(s, uriTemplate, name, description, mimeType, handler); err != nil {
		panic(err)
	}
}

// RegisterTextResource registers a static text resource at uri.
func RegisterTextResource(s *TinyServer, uri, name, description, mimeType, text string) error {
	return RegisterResource(s, uri, name, description, mimeType,
		func(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
			if req.Params.URI != uri {
				return nil, mcp.ResourceNotFoundError(req.Params.URI)
			}
			return TextResource(uri, mimeType, text), nil
		})
}

// MustRegisterTextResource is like RegisterTextResource but panics on error.
func MustRegisterTextResource(s *TinyServer, uri, name, description, mimeType, text string) {
	if err := RegisterTextResource(s, uri, name, description, mimeType, text); err != nil {
		panic(err)
	}
}

// TextResource builds a ReadResourceResult with a single text body.
func TextResource(uri, mimeType, text string) *mcp.ReadResourceResult {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{TextResourceContents(uri, mimeType, text)},
	}
}

// TextResourceContents builds a text ResourceContents block for ReadResourceResult.
func TextResourceContents(uri, mimeType, text string) *mcp.ResourceContents {
	return &mcp.ResourceContents{
		URI:      uri,
		MIMEType: mimeType,
		Text:     text,
	}
}
