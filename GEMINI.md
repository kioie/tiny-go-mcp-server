# Tiny Go MCP Server — Gemini CLI context

This repository is **Tiny Go MCP Server**: a thin Go helper library (`tinymcp`) on the official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).

## Layout

- `tinymcp/` — importable library
- `cmd/tiny-go-mcp/` — reference MCP server (demo tools)
- `examples/` — minimal, HTTP, http-deploy, resources
- `docs/` — HTTP, discovery, Smithery, Glama

## Review priorities for this repo

- Registration must not panic; return errors for nil server/handler at registration time
- MCP tool descriptions must state when to use, when not to, and sibling tools
- Handlers: `context.Context` first; success text via `tinymcp.TextResult`
- HTTP/SSE deploy paths: auth, localhost binding, cross-origin protection for public endpoints
- Race-safe tests (`go test -race`); table-driven tests where helpful
- Godoc on exported symbols; standard Go formatting
- Dual imports (`tinymcp` + `mcp`) are intentional — not a leaky-facade bug

See `AGENTS.md` for full agent/CI conventions.
