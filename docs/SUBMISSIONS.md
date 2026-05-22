# Community listings and launch copy

Steps and text for second-tier directories and social posts (P1).

## awesome-mcp-servers

- **PR:** https://github.com/punkpeye/awesome-mcp-servers/pull/6665 (open, CI green)
- **Action:** Wait for maintainer merge; comment if stale

## mcp.so

1. Open https://mcp.so and use **Submit Server** (or equivalent form).
2. Use this blurb:

**Name:** Tiny Go MCP Server  
**URL:** https://github.com/kioie/tiny-go-mcp-server  
**Category:** Framework / Developer Tools  
**Description:**

> Minimal Go library on the official MCP go-sdk — stdio, streamable HTTP, legacy SSE; tools, resources, prompts; struct-derived JSON Schema; ~5MB static binaries. Reference server on Glama and MCP Registry (`io.github.kioie/tiny-go-mcp` v1.1.0).

**Install:**

```bash
go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest
go install github.com/kioie/tiny-go-mcp-server/cmd/tiny-go-mcp@latest
```

## Smithery

1. Install CLI: `npm install -g @smithery/cli` (or `npx @smithery/cli`)
2. Authenticate: `smithery auth login`
3. Repo includes [`smithery.yaml`](../smithery.yaml) for stdio via Docker.
4. Publish:

```bash
smithery mcp publish --name @kioie/tiny-go-mcp --transport stdio
```

Or connect GitHub at https://smithery.ai/new if prompted.

## Project scaffold (`gonew`)

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
cd my-mcp
go run .
```

Requires Go 1.26+ and tagged submodule `template` (`template/v1.1.1+`).

## Launch post (Glama Discord / r/mcp)

**Title:** Tiny Go MCP Server v1.1.0 — minimal tinymcp on official go-sdk

**Body:**

We shipped **Tiny Go MCP Server** — a thin Go layer on the official `modelcontextprotocol/go-sdk`:

- **tinymcp** library: typed tools, resources, prompts, stdio + HTTP/SSE helpers
- Reference server: demo tools + resource + prompt
- **MCP Registry:** `io.github.kioie/tiny-go-mcp` v1.1.0
- **Glama:** https://glama.ai/mcp/servers/kioie/tiny-go-mcp-server
- **Docker:** `ghcr.io/kioie/tiny-go-mcp:1.1.0`

```bash
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest
```

Feedback and PRs welcome: https://github.com/kioie/tiny-go-mcp-server

---

After posting, link back from README **Discovery** section if desired.
