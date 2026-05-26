# Discovery and publishing

Ways to make **Tiny Go MCP Server** visible to AI clients and the community.

## Go package (library)

1. Tag a semver release on GitHub (e.g. `v1.0.0`) so `go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest` resolves to a stable version.
2. pkg.go.dev indexes automatically: https://pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp
3. Optional: add a [Go Report Card](https://goreportcard.com/report/github.com/kioie/tiny-go-mcp-server/tinymcp) badge after the first tag.

## MCP server (reference binary)

The example server is **`tiny-go-mcp`** (stdio). Client config template: [`examples/mcp-client-config.json`](../examples/mcp-client-config.json).

### Official MCP Registry

1. Install the publisher CLI: https://modelcontextprotocol.io/registry/quickstart
2. Edit [`server.json`](../server.json) (namespace `io.github.kioie/tiny-go-mcp`, version aligned with release tag).
3. Authenticate: `mcp-publisher login github`
4. Validate: `mcp-publisher validate` (from repo root)
5. Publish: `mcp-publisher publish`
6. After tagging a release (e.g. `v1.1.3`), the Release workflow also pushes `ghcr.io/kioie/tiny-go-mcp:<version>`.

### Community directories (manual)

| Directory | How to submit |
|-----------|----------------|
| [awesome-mcp-servers](https://github.com/punkpeye/awesome-mcp-servers) | PR to README (library + server) — [PR #6665](https://github.com/punkpeye/awesome-mcp-servers/pull/6665) |
| [mcp.so](https://mcp.so) | Site form — blurb in [SUBMISSIONS.md](./SUBMISSIONS.md) |
| [Glama](https://glama.ai/mcp/servers) | GitHub + [`Dockerfile`](../Dockerfile) + [`glama.json`](../glama.json) — [GLAMA.md](./GLAMA.md) |
| [Smithery](https://smithery.ai) | URL: [`examples/http-deploy`](../examples/http-deploy) + [`docs/SMITHERY.md`](./SMITHERY.md); stdio MCPB: [`smithery.yaml`](../smithery.yaml) — [live](https://smithery.ai/servers/kioie/tiny-go-mcp) |

### Community posts

| Channel | Link |
|---------|------|
| [r/mcp](https://www.reddit.com/r/mcp/) | [Launch post (v1.1)](https://www.reddit.com/r/mcp/comments/1tkru4a/tiny_go_mcp_server_v11_minimal_tinymcp_on/) |

Official logo and banner: [`docs/BRANDING.md`](./BRANDING.md) · [`docs/assets/`](../docs/assets/)

## GitHub repository

- **About** description and **topics** (`mcp`, `go`, `model-context-protocol`, …) improve search.
- Pin the README **Quick start** and link to `examples/minimal`.
- GitHub **Releases** attach cross-platform binaries (see `.github/workflows/release.yml`).

## Positioning one-liner

> Minimal Go library on the official MCP SDK — stdio, HTTP, SSE; tools, resources, prompts; struct-derived JSON Schema; ~5MB static binaries.

Use this in directory listings and social posts.

## Release backlog

When cutting the next semver release, include postponed items:

- [ ] **Graceful HTTP/SSE shutdown** — `ListenAndServeHTTPContext`, SIGTERM/`Shutdown` drain (postponed from v1.1.3 planning)

## Descriptions for agents

Registry and directory text (`server.json`, `glama.json`) and each MCP **tool** description should tell models **when to use**, **when not to**, and **alternatives** (including sibling tools: “use `add` instead of `subtract` when …”). See [`README.md`](../README.md) “When to use what” and the reference tool table.
