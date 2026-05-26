# Smithery publishing

Two ways to list **Tiny Go MCP Server** on [Smithery](https://smithery.ai):

| Mode | Who runs Docker? | Best for |
|------|------------------|----------|
| **URL (streamable HTTP)** | Nobody (you host HTTPS; users connect via Smithery Gateway) | Remote MCP, zero client setup |
| **MCPB (stdio bundle)** | End user (bundle runs `ghcr.io/kioie/tiny-go-mcp`) | Local stdio clients via Smithery install flow |

## URL listing (recommended for remote access)

Smithery proxies to your **public HTTPS** streamable HTTP endpoint. Users do not install Docker or a local binary.

### 1. Deploy the HTTP example

Runnable template: [`examples/http-deploy`](../examples/http-deploy).

```bash
go run ./examples/http-deploy
```

Production deploy (no Docker required for operators on Render/Railway — native Go build):

- **Render:** root dir `examples/http-deploy`, build `go build -o server .`, start `./server` — see [`render.yaml`](../examples/http-deploy/render.yaml)
- **Railway:** `railway up` from `examples/http-deploy`
- **Fly.io:** `fly deploy --config fly.http.toml` from repo root — see [`fly.http.toml`](../fly.http.toml)

The example exposes:

| Path | Purpose |
|------|---------|
| `/` | Streamable HTTP MCP |
| `/health` | Health check |
| `/.well-known/mcp/server-card.json` | Static metadata when live scan is blocked |

### 2. Publish on Smithery

```bash
smithery auth login
smithery mcp publish "https://YOUR_HOST" -n your-namespace/your-server
```

Requirements ([Smithery docs](https://smithery.ai/docs/build/publish)):

- **Streamable HTTP** transport (this repo’s `StartHTTP` / `StreamableHTTPHandler`)
- Public **HTTPS** URL (no trailing slash)
- **OAuth** only if your server requires auth (return **401**, not 403, for unauthenticated requests)

Optional: tunnel locally with `ngrok http 8080` before publishing a test listing.

### 3. Server card (optional but included in example)

If WAF/bot protection blocks Smithery’s scan (`User-Agent: SmitheryBot/1.0`), serve metadata at:

`/.well-known/mcp/server-card.json`

See [`examples/http-deploy/server-card.json`](../examples/http-deploy/server-card.json) and [SEP-1649](https://github.com/modelcontextprotocol/modelcontextprotocol/issues/1649).

### 4. Mount MCP on a subpath

If MCP lives at `/mcp` instead of `/`:

```go
mux.Handle("/mcp", mcpHandler)
```

Publish `https://YOUR_HOST/mcp` to Smithery.

---

## MCPB listing (stdio / local)

For stdio distribution through Smithery’s MCPB flow (clients download and run locally):

1. Build bundle from [`smithery.yaml`](../smithery.yaml):

```bash
npx mcp-bundler . .smithery/mcpb --entry=smithery/entry.js --inspect
npx @anthropic-ai/mcpb pack .smithery/mcpb server.mcpb
```

2. Publish:

```bash
smithery mcp publish server.mcpb -n kioie/tiny-go-mcp
```

This path **requires Docker on the end-user machine** (bundle runs `ghcr.io/kioie/tiny-go-mcp:1.1.0`).

**Live stdio listing:** https://smithery.ai/servers/kioie/tiny-go-mcp

**Live URL listing (HTTP):** https://smithery.ai/servers/kioie/tiny-go-mcp-http — upstream `https://tiny-go-mcp-http.fly.dev`

---

## Troubleshooting URL scans

| Symptom | Fix |
|---------|-----|
| 403 during scan | Whitelist `SmitheryBot/1.0` or disable bot fight mode; prefer returning **401** for auth |
| Scan timeout | Add `/.well-known/mcp/server-card.json` |
| Wrong tools on listing | Update server card or redeploy after tool changes |

See also [`docs/SUBMISSIONS.md`](./SUBMISSIONS.md) and [`docs/HTTP.md`](./HTTP.md).
