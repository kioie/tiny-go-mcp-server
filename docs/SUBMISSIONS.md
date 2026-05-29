# Community listings and launch copy

Steps and text for second-tier directories and social posts (P1).

## awesome-mcp-servers

- **Listed:** https://github.com/punkpeye/awesome-mcp-servers#frameworks (Frameworks section, merged [PR #6665](https://github.com/punkpeye/awesome-mcp-servers/pull/6665))
- **Description refresh:** [PR #6965](https://github.com/punkpeye/awesome-mcp-servers/pull/6965) (v1.1.3, HTTP/SSE, Smithery link)

## awesome-go

**Not yet listed.** Proposed PR to [avelino/awesome-go](https://github.com/avelino/awesome-go) under **Artificial Intelligence** (or a new **Model Context Protocol** subsection if maintainers prefer):

```markdown
- [tiny-go-mcp-server](https://github.com/kioie/tiny-go-mcp-server) - Thin helper on the official MCP go-sdk — stdio, streamable HTTP, legacy SSE; struct-tag JSON Schema; resources and prompts. Quick start: [docs/QUICKSTART.md](https://github.com/kioie/tiny-go-mcp-server/blob/main/docs/QUICKSTART.md).
```

Also list the official SDK if missing:

```markdown
- [go-sdk](https://github.com/modelcontextprotocol/go-sdk) - Official Go SDK for Model Context Protocol servers and clients (Google + MCP collaboration).
```

Submit via awesome-go PR; link back from README once merged.

## mcp.so

**Submitted:** https://github.com/chatmcp/mcpso/issues/2470 (awaiting triage)

1. Open https://mcp.so and use **Submit Server** (or equivalent form).
2. Use this blurb:

**Name:** Tiny Go MCP Server  
**URL:** https://github.com/kioie/tiny-go-mcp-server  
**Category:** Framework / Developer Tools  
**Description:**

> Minimal Go library on the official MCP go-sdk — stdio, streamable HTTP, legacy SSE; tools, resources, prompts; struct-derived JSON Schema; ~5MB static binaries. Reference server on Glama and MCP Registry (`io.github.kioie/tiny-go-mcp` v1.1.3).

**Install:**

```bash
go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest
go install github.com/kioie/tiny-go-mcp-server/cmd/tiny-go-mcp@latest
```

## Smithery

**Live listing (stdio MCPB):** https://smithery.ai/servers/kioie/tiny-go-mcp (icon uploaded)

### URL listing (no Docker for users)

**Live listing:** https://smithery.ai/servers/kioie/tiny-go-mcp-http  
**Upstream:** https://tiny-go-mcp-http.fly.dev  
**Deploy:** [`examples/http-deploy`](../examples/http-deploy) + [`fly.http.toml`](../fly.http.toml)

Host streamable HTTP and publish your HTTPS URL — see [`docs/SMITHERY.md`](./SMITHERY.md):

```bash
go run ./examples/http-deploy
# deploy to Render/Railway/Fly, then:
smithery auth login
smithery mcp publish "https://YOUR_HOST" -n your-namespace/your-server
```

### MCPB listing (stdio; requires Docker on client)

1. Install CLI: `npm install -g @smithery/cli` (or `npx @smithery/cli`)
2. Authenticate: `smithery auth login`
3. Build MCPB bundle from [`smithery.yaml`](../smithery.yaml):

```bash
npx mcp-bundler . .smithery/mcpb --entry=smithery/entry.js --inspect
npx @anthropic-ai/mcpb pack .smithery/mcpb server.mcpb
```

4. Publish (namespace must match your Smithery account, e.g. `kioie/tiny-go-mcp`):

```bash
smithery mcp publish server.mcpb -n kioie/tiny-go-mcp
```

Requires Docker on the client machine (bundle runs `ghcr.io/kioie/tiny-go-mcp:1.1.3`). Alternative: connect GitHub at https://smithery.ai/new.

## Project scaffold (`gonew`)

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
cd my-mcp
go run .
```

Requires Go 1.26+ and tagged submodule `template` (`template/v1.1.1+`).

### HTTP deploy (`template-http`)

```bash
gonew github.com/kioie/tiny-go-mcp-server/template-http@latest example.com/my-mcp-http my-mcp-http
cd my-mcp-http
go run .
```

Requires tagged `template-http` submodule. See [`template-http/README.md`](../template-http/README.md).

## Launch post (r/mcp)

**Posted:** https://www.reddit.com/r/mcp/comments/1tkru4a/tiny_go_mcp_server_v11_minimal_tinymcp_on/

**Title:** Tiny Go MCP Server v1.1 — minimal tinymcp on official go-sdk

**Body:**

We shipped **Tiny Go MCP Server** — a thin Go layer on the official `modelcontextprotocol/go-sdk`:

- **tinymcp** library: typed tools, resources, prompts, stdio + HTTP/SSE helpers
- Reference server: demo tools + resource + prompt
- **MCP Registry:** `io.github.kioie/tiny-go-mcp` v1.1.3
- **Glama:** https://glama.ai/mcp/servers/kioie/tiny-go-mcp-server
- **Smithery:** https://smithery.ai/servers/kioie/tiny-go-mcp
- **Docker:** `ghcr.io/kioie/tiny-go-mcp:1.1.3`

```bash
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest
```

Feedback and PRs welcome: https://github.com/kioie/tiny-go-mcp-server
