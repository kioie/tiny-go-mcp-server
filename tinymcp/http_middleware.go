package tinymcp

import (
	"net/http"
	"strings"
)

// WithCrossOriginProtection wraps h with Go's cross-origin CSRF protection for browser MCP clients.
// Server-to-server clients without Origin/Sec-Fetch-Site headers are unaffected.
func WithCrossOriginProtection(h http.Handler) http.Handler {
	return http.NewCrossOriginProtection().Handler(h)
}

// BearerTokenAuth requires Authorization: Bearer <token> when token is non-empty.
// When token is empty, h is served without auth. Unauthenticated requests receive 401
// with WWW-Authenticate (preferred for Smithery and MCP gateways).
func BearerTokenAuth(token string, h http.Handler) http.Handler {
	if token == "" {
		return h
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, ok := bearerTokenFromHeader(r.Header.Get("Authorization"))
		if !ok || got != token {
			w.Header().Set("WWW-Authenticate", `Bearer realm="MCP"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func bearerTokenFromHeader(h string) (string, bool) {
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return "", false
	}
	return strings.TrimSpace(h[len(prefix):]), true
}
