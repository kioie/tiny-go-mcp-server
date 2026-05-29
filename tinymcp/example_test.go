package tinymcp_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type greetArgs struct {
	Name string `json:"name" jsonschema:"Person to greet"`
}

func ExampleRegisterTool() {
	s := tinymcp.NewServer("demo", "1.0.0")
	err := tinymcp.RegisterTool(s, "greet", "Greet someone by name", func(_ context.Context, _ *mcp.CallToolRequest, args greetArgs) (*mcp.CallToolResult, any, error) {
		return tinymcp.TextResult("Hello, " + args.Name), nil, nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleStreamableHTTPHandler() {
	s := tinymcp.NewServer("demo", "1.0.0")
	if err := tinymcp.RegisterTool(s, "ping", "Echo ping", func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		return tinymcp.TextResult("pong"), nil, nil
	}); err != nil {
		log.Fatal(err)
	}

	handler, err := tinymcp.StreamableHTTPHandler(s, &tinymcp.HTTPOptions{Stateless: true})
	if err != nil {
		log.Fatal(err)
	}
	handler = tinymcp.WithCrossOriginProtection(handler)
	handler = tinymcp.BearerTokenAuth("secret", handler)

	http.Handle("/mcp", handler)
	// Listen with tinymcp.ListenAndServeHTTP or your own http.Server + TLS.
}

func ExampleTextResult() {
	r := tinymcp.TextResult("hello")
	tc := r.Content[0].(*mcp.TextContent)
	fmt.Println(tc.Text)
	// Output: hello
}

func ExampleBearerTokenAuth() {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	h := tinymcp.BearerTokenAuth("secret", inner)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	fmt.Println(rec.Body.String())
	// Output: ok
}
