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
_ = tinymcp.RegisterTool(server, "ping", "Echo a message", pingHandler)

// Blocks; bind to loopback for local-only access. Add TLS or a reverse proxy in production.
log.Fatal(server.StartHTTP("127.0.0.1:8080", &tinymcp.HTTPOptions{
    Stateless: true, // simple demos; omit for full sessions + server-initiated messages
}))
```

Or mount the handler on an existing `http.Server` (auth, CORS, paths):

```go
handler, err := tinymcp.StreamableHTTPHandler(server, nil)
if err != nil {
    log.Fatal(err)
}
mux := http.NewServeMux()
mux.Handle("/mcp", handler)
log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))
```

Runnable example: [`examples/http`](../examples/http).

### HTTPOptions

| Field | Effect |
|-------|--------|
| `Stateless` | No session ID; one-shot requests (no server→client RPC) |
| `JSONResponse` | POST replies as `application/json` instead of SSE |
| `SessionTimeout` | Close idle sessions after duration |
| `DisableLocalhostProtection` | Disable DNS rebinding guard (use with care) |

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
- For local **ngrok/tunnel** testing, set `DisableLocalhostProtection: true` in `HTTPOptions` (included in [`examples/http-deploy`](../examples/http-deploy)).
- Put **TLS**, authentication, and rate limiting in front of public endpoints (reverse proxy or middleware wrapping the handler).
- Do not set `DisableLocalhostProtection` unless you understand [MCP security guidance](https://modelcontextprotocol.io/specification/2025-11-25/basic/security_best_practices).
- `StartHTTP` and `ListenAndServeHTTP` set `ReadHeaderTimeout` and `IdleTimeout` on the underlying `http.Server`.

## Escape hatch

For custom session routing, event stores, or advanced streamable options, use `server.RawServer()` with `mcp.NewStreamableHTTPHandler` directly — same underlying SDK, full control.
