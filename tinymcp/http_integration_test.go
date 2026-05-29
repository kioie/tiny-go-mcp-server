package tinymcp

import (
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"
)

// TestHTTPIntegration_middlewareAndContext verifies middleware wrapping and
// context-based graceful shutdown over a real TCP listener.
func TestHTTPIntegration_middlewareAndContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	t.Cleanup(cancel)

	var sawMiddleware bool
	s := NewServer("integration-http", "1.0.0")
	opts := (&HTTPOptions{JSONResponse: true}).WithMiddleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sawMiddleware = true
			next.ServeHTTP(w, r)
		})
	})
	h, err := StreamableHTTPHandler(s, opts)
	if err != nil {
		t.Fatal(err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	ln.Close()

	done := make(chan error, 1)
	go func() {
		done <- ListenAndServeHTTPContext(ctx, addr, h)
	}()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	reqBody := strings.NewReader(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`)
	req, err := http.NewRequestWithContext(t.Context(), http.MethodPost, "http://"+addr, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST initialize: %v", err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, body = %s", resp.StatusCode, body)
	}
	if !strings.Contains(string(body), `"protocolVersion"`) {
		t.Fatalf("expected initialize result, got: %s", body)
	}
	if !sawMiddleware {
		t.Fatal("middleware was not invoked")
	}

	cancel()
	select {
	case err := <-done:
		if err != nil && err != context.Canceled {
			t.Fatalf("ListenAndServeHTTPContext: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for graceful shutdown")
	}
}
