# HTTP and SSE transports

`tinymcp` defaults to **stdio** (`Start()`), which local clients (Cursor, Claude Desktop) expect. For remote clients, gateways, or browser-based tools, use **streamable HTTP** or legacy **SSE** via thin helpers on the official [go-sdk](https://github.com/modelcontextprotocol/go-sdk).

## Which transport?

| Transport | tinymcp API | Best for |
|-----------|-------------|----------|
| **stdio** | `Start()` | Local AI clients, subprocess MCP |
| **Streamable HTTP** | `StartHTTP` / `StreamableHTTPHandler` | Current MCP spec; remote clients, Glama-style gateways, POST + SSE |
| **Legacy SSE** | `StartSSE` / `SSEHandler` | Older clients on MCP spec 2024-11-05 |

Prefer **streamable HTTP** for new remote deployments unless a client explicitly requires legacy SSE.

## Streamable HTTP (recommended)

```go
server := tinymcp.NewServer("my-mcp", "1.0.0")
if err := tinymcp.RegisterTool(server, "ping", "Echo a message", pingHandler); err != nil {
	log.Fatal(err)
}

// Blocks; bind to loopback for local-only access. Add TLS or a reverse proxy in production.
log.Fatal(server.StartHTTP("127.0.0.1:8080", &tinymcp.HTTPOptions{
    Stateless: true, // simple demos; omit for full sessions + server-initiated messages
}))
```

Or mount the handler on an existing `http.Server` (auth, CORS, paths):

```go
opts := (&tinymcp.HTTPOptions{Stateless: true}).WithMiddleware(
    requestLogger, // your func(http.Handler) http.Handler
)
handler, err := tinymcp.StreamableHTTPHandler(server, opts)
if err != nil {
	log.Fatal(err)
}
handler = tinymcp.WithCrossOriginProtection(handler)
handler = tinymcp.BearerTokenAuth(os.Getenv("MCP_API_KEY"), handler)
mux := http.NewServeMux()
mux.Handle("/mcp", handler)
log.Fatal(tinymcp.ListenAndServeHTTP("127.0.0.1:8080", mux))
```

`ListenAndServeHTTP` drains active connections on **SIGINT** or **SIGTERM** (PaaS deploys). For custom lifecycle control:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
if err := tinymcp.ListenAndServeHTTPContext(ctx, addr, mux); err != nil && !errors.Is(err, context.Canceled) {
	log.Fatal(err)
}
```

Runnable example: [`examples/http`](../examples/http).

### HTTPOptions

| Field | Effect |
|-------|--------|
| `Stateless` | No session ID; one-shot requests (no server→client RPC) |
| `JSONResponse` | POST replies as `application/json` instead of SSE |
| `SessionTimeout` | Close idle sessions after duration |
| `DisableLocalhostProtection` | Disable DNS rebinding guard (use with care) |
| `Middleware` | Wrap MCP handler (first = innermost); chain with `WithMiddleware` |

### Middleware order

Recommended wrapping (inner → outer):

1. `StreamableHTTPHandler` (MCP core)
2. `HTTPOptions.Middleware` — logging, rate limits, metrics
3. `WithCrossOriginProtection` — browser clients
4. `BearerTokenAuth` — optional API key

See [`examples/http-deploy`](../examples/http-deploy) for a full stack.

### Client connection

Point an MCP client that supports streamable HTTP at your server URL (often the root path or `/mcp`). See your client’s docs for `url` / `transport: http` config.

### Smithery URL listing (no Docker for users)

Host streamable HTTP on a public HTTPS URL and publish to [Smithery](https://smithery.ai) so clients connect via the Smithery Gateway — end users do not run Docker or a local binary.

Deploy template: [`examples/http-deploy`](../examples/http-deploy) (includes `/health` and `/.well-known/mcp/server-card.json`). Full steps: [`docs/SMITHERY.md`](./SMITHERY.md).

## Legacy SSE (2024-11-05)

```go
log.Fatal(server.StartSSE(":8080", nil))
```

Clients use `SSEClientTransport` with your server’s SSE endpoint URL. Prefer streamable HTTP for new work.

## Security

- Bind to `127.0.0.1:8080` when only local access is needed (`TINY_GO_MCP_ADDR=0.0.0.0:8080` to listen on all interfaces).
- For local **ngrok/tunnel** testing on loopback, set `TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1` (`tinymcp.DisableLocalhostProtectionFromEnv()` in [`examples/http`](../examples/http) and [`examples/http-deploy`](../examples/http-deploy)).
- Put **TLS** in front of public endpoints — see [`docs/TLS.md`](TLS.md). Authentication and rate limiting: reverse proxy or middleware. Helpers: `tinymcp.BearerTokenAuth`, `tinymcp.WithCrossOriginProtection` (used in [`examples/http-deploy`](../examples/http-deploy)).
- Do not set `DisableLocalhostProtection` unless you understand [MCP security guidance](https://modelcontextprotocol.io/specification/2025-11-25/basic/security_best_practices).
- `StartHTTP` and `ListenAndServeHTTP` set `ReadHeaderTimeout` (10s), `ReadTimeout` (30s), and `IdleTimeout` (120s) on the underlying `http.Server`, and drain on SIGINT/SIGTERM. `WriteTimeout` is intentionally unset so streamable HTTP/SSE responses can be long-lived.

## Escape hatch

For custom session routing, event stores, or advanced streamable options, use `server.RawServer()` with `mcp.NewStreamableHTTPHandler` directly — same underlying SDK, full control.
