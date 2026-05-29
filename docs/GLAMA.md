# Glama deployment (Docker)

How to list and host **Tiny Go MCP Server** on [Glama](https://glama.ai/mcp/servers).

## Files in this repo

| File | Purpose |
|------|---------|
| [`Dockerfile`](../Dockerfile) | Multi-stage Go build → `tiny-go-mcp` binary, stdio `ENTRYPOINT` |
| [`glama.json`](../glama.json) | Glama directory metadata, `mcpServers`, OCI package hint |
| [`.dockerignore`](../.dockerignore) | Smaller build context |

## Local Docker test

```bash
docker build -t tiny-go-mcp:local .
docker run -i --rm tiny-go-mcp:local
```

The container speaks MCP on **stdin/stdout**. Do not publish ports for stdio mode.

Optional stderr logs:

```bash
docker run -i --rm -e TINY_GO_MCP_VERBOSE=1 tiny-go-mcp:local
```

## Deploy on Glama (from GitHub)

1. Open [Glama MCP Hosting](https://glama.ai/mcp/hosting) or [Deploy a server](https://glama.ai/mcp/servers).
2. Choose **From GitHub**.
3. Install the **Glama GitHub App** on `kioie/tiny-go-mcp-server`.
4. Select the repository and branch (`main`).
5. Build settings (typical):
   - **Dockerfile path:** `Dockerfile`
   - **Context:** `.` (repository root)
   - **Entrypoint / CMD:** inferred from `ENTRYPOINT ["/usr/local/bin/tiny-go-mcp"]`
6. **Environment variables** (optional):
   - `TINY_GO_MCP_VERBOSE` = `1` if you want startup logs on stderr
7. Deploy. Glama wraps stdio and exposes a **Streamable HTTP** gateway URL for remote clients.

## Deploy on Glama (from Dockerfile / image)

If you publish the image to GHCR first (release tags push **linux/amd64** and **linux/arm64** via GitHub Actions):

```bash
docker build -t ghcr.io/kioie/tiny-go-mcp:1.1.3 .
docker push ghcr.io/kioie/tiny-go-mcp:1.1.3
```

Then in Glama choose **From a package or image** and use:

- **Image:** `ghcr.io/kioie/tiny-go-mcp:1.1.3`
- **Transport:** stdio (Glama wraps automatically)

3. Glama indexes from GitHub — after merging a release, open [server admin](https://glama.ai/mcp/servers/kioie/tiny-go-mcp-server/admin) and click **Sync Server** so tool/resource/prompt scores refresh.

## Directory listing

After deployment works:

1. Set deployment visibility to **public** in Glama (optional).
2. Ensure [`glama.json`](../glama.json) is on `main` with correct `maintainers` and `version`.
3. Glama indexes tools from the running server for search and badges.

## MCP Registry (`server.json`)

For the official MCP Registry (separate from Glama), see [`server.json`](../server.json) and [DISCOVERY.md](./DISCOVERY.md).

## Image details (for Glama form fields)

| Field | Value |
|-------|--------|
| **Base build** | `golang:1.26-alpine` → `alpine:3.21` |
| **Binary** | `/usr/local/bin/tiny-go-mcp` |
| **User** | `nobody` |
| **Exposed ports** | None (stdio) |
| **Health** | Glama gateway `/ping` (platform-side; not in container) |
| **Tools** | `add`, `subtract`, `greet` |
