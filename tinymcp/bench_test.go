package tinymcp

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func BenchmarkRegisterTool(b *testing.B) {
	handler := func(_ context.Context, _ *mcp.CallToolRequest, _ echoArgs) (*mcp.CallToolResult, any, error) {
		return TextResult("ok"), nil, nil
	}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		s := NewServer("bench", "1.0.0")
		if err := RegisterTool(s, "echo", "Echo", handler); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkTextResult(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = TextResult("hello")
	}
}

func BenchmarkStreamableHTTPHandler_initialize(b *testing.B) {
	s := NewServer("bench-http", "1.0.0")
	if err := RegisterTool(s, "ping", "Ping", func(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
		return TextResult("pong"), nil, nil
	}); err != nil {
		b.Fatal(err)
	}
	h, err := StreamableHTTPHandler(s, &HTTPOptions{JSONResponse: true})
	if err != nil {
		b.Fatal(err)
	}
	ts := httptest.NewServer(h)
	b.Cleanup(ts.Close)

	body := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequestWithContext(b.Context(), http.MethodPost, ts.URL, bytes.NewReader(body))
		if err != nil {
			b.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json, text/event-stream")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			b.Fatalf("status = %d", resp.StatusCode)
		}
	}
}
