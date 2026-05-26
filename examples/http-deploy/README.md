# HTTP deploy example (Smithery URL listing)

Host a **streamable HTTP** MCP server so users connect through [Smithery](https://smithery.ai) — **no Docker required on their machine** (unlike the stdio MCPB bundle in [`smithery.yaml`](../../smithery.yaml)).

| Path | Purpose |
|------|---------|
| `/` | Streamable HTTP MCP (POST JSON-RPC) |
| `/health` | Platform health check |
| `/.well-known/mcp/server-card.json` | Static metadata for Smithery scan / listing |

Full Smithery guide: [`docs/SMITHERY.md`](../../docs/SMITHERY.md).

## Security

This example is **open by default** so Smithery can scan and proxy it (`TINY_GO_MCP_API_KEY` unset). Anyone who can reach your public URL can call MCP methods.

For production:

1. Set **`TINY_GO_MCP_API_KEY`** to require `Authorization: Bearer <key>` on MCP (`/`). Return **401** (not 403) when auth fails — Smithery-compatible.
2. Cross-origin CSRF protection is enabled on the MCP handler via `tinymcp.WithCrossOriginProtection` (server-to-server clients without `Origin` are unaffected).
3. **Do not deploy real tools or sensitive data** without reviewing auth and tool privileges.

See [`docs/HTTP.md`](../../docs/HTTP.md) and [`docs/SMITHERY.md`](../../docs/SMITHERY.md).

## Run locally

```bash
go run ./examples/http-deploy
# MCP: http://127.0.0.1:8080/
```

Optional public tunnel for Smithery testing (requires disabling localhost DNS-rebinding protection):

```bash
TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1 go run ./examples/http-deploy
ngrok http 8080
# Publish https://YOUR_SUBDOMAIN.ngrok-free.app (see below)
```

## Deploy (no Docker for end users)

Build a single static binary (`CGO_ENABLED=0 go build -o server .` from this directory) or let a PaaS compile Go for you.

### Render (native Go)

1. New **Web Service** → connect this repo.
2. **Root directory:** `examples/http-deploy`
3. **Build:** `go build -o server .`
4. **Start:** `./server`
5. Copy the `https://….onrender.com` URL.

Or use [`render.yaml`](./render.yaml) if deploying via Render Blueprint.

### Railway

```bash
cd examples/http-deploy
railway init
railway up
```

Railway sets `PORT` automatically; this example reads it.

### Fly.io

From the **repository root** (Dockerfile needs `tinymcp/` and `go.mod`):

```bash
fly launch --config fly.http.toml
fly deploy --config fly.http.toml
```

Live example: https://tiny-go-mcp-http.fly.dev → Smithery https://smithery.ai/servers/kioie/tiny-go-mcp-http

## Publish on Smithery (URL)

1. `smithery auth login`
2. Verify MCP with a **POST** JSON-RPC `initialize` to `/` (returns **200**). `GET /` returns **405** — that is expected in stateless mode.
3. Publish:

```bash
smithery mcp publish "https://YOUR_HOST" -n your-namespace/your-server
```

Use the **exact public HTTPS origin** Smithery should proxy (no trailing slash). MCP is served at `/`.

If scanning fails (WAF, auth wall), Smithery can use `/.well-known/mcp/server-card.json` — already included in this example.

When you add or change tools, resources, or prompts, update [`server-card.json`](./server-card.json) to match. CI runs `go test ./examples/http-deploy` to catch drift.

## Environment

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | — | Set by most PaaS (e.g. Render, Railway) |
| `TINY_GO_MCP_ADDR` | `127.0.0.1:8080` | Listen address when `PORT` is unset (PaaS sets `PORT` to bind on all interfaces) |
| `TINY_GO_MCP_API_KEY` | — | When set, MCP (`/`) requires `Authorization: Bearer <key>` |
| `TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION` | — | Set to `1` for local ngrok/tunnel testing only (disables DNS rebinding guard on loopback) |
| `TINY_GO_MCP_VERBOSE` | — | Set to `1` for stderr startup logs |
