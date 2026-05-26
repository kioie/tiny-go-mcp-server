package tinymcp

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStreamableHTTPHandler_initialize(t *testing.T) {
	s := NewServer("test-http", "1.0.0")
	h, err := StreamableHTTPHandler(s, &HTTPOptions{JSONResponse: true})
	if err != nil {
		t.Fatalf("StreamableHTTPHandler: %v", err)
	}

	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	req, err := http.NewRequest(http.MethodPost, ts.URL, strings.NewReader(
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`,
	))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), `"protocolVersion"`) {
		t.Fatalf("expected initialize result, got: %s", body)
	}
}

func TestSSEHandler_connect(t *testing.T) {
	s := NewServer("test-sse", "1.0.0")
	h, err := SSEHandler(s, nil)
	if err != nil {
		t.Fatalf("SSEHandler: %v", err)
	}

	ts := httptest.NewServer(h)
	t.Cleanup(ts.Close)

	req, err := http.NewRequest(http.MethodGet, ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type = %q, want event-stream", ct)
	}
}

func TestHTTPHandlers_nilServer(t *testing.T) {
	var s *TinyServer
	if _, err := StreamableHTTPHandler(s, nil); err == nil {
		t.Error("StreamableHTTPHandler(nil): expected error")
	}
	if _, err := SSEHandler(s, nil); err == nil {
		t.Error("SSEHandler(nil): expected error")
	}
	if err := s.StartHTTP(":0", nil); err == nil {
		t.Error("StartHTTP(nil): expected error")
	}
	if err := s.StartSSE(":0", nil); err == nil {
		t.Error("StartSSE(nil): expected error")
	}
}

func TestBearerTokenAuth(t *testing.T) {
	ok := httptest.NewRecorder()
	BearerTokenAuth("secret", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(ok, httptest.NewRequest(http.MethodPost, "/", nil))
	if ok.Code != http.StatusUnauthorized {
		t.Fatalf("missing auth: status = %d, want 401", ok.Code)
	}

	withToken := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("Authorization", "Bearer secret")
	BearerTokenAuth("secret", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(withToken, req)
	if withToken.Code != http.StatusOK {
		t.Fatalf("valid auth: status = %d, want 200", withToken.Code)
	}

	open := httptest.NewRecorder()
	http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).ServeHTTP(open, httptest.NewRequest(http.MethodPost, "/", nil))
	BearerTokenAuth("", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})).ServeHTTP(open, httptest.NewRequest(http.MethodPost, "/", nil))
	if open.Code != http.StatusOK {
		t.Fatalf("empty token disables auth: status = %d, want 200", open.Code)
	}
}

func TestWithCrossOriginProtection(t *testing.T) {
	h := WithCrossOriginProtection(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodPost, "http://example.com/", nil)
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Origin", "https://evil.example")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("cross-origin POST: status = %d, want 403", rec.Code)
	}

	rec2 := httptest.NewRecorder()
	h.ServeHTTP(rec2, httptest.NewRequest(http.MethodPost, "http://example.com/", nil))
	if rec2.Code != http.StatusOK {
		t.Fatalf("server-to-server POST: status = %d, want 200", rec2.Code)
	}
}
