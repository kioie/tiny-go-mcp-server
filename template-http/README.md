# template-http — gonew scaffold for HTTP MCP

Streamable HTTP server with auth/CORS middleware pre-wired. For stdio, use [`template/`](../template/).

## Create a project

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kioie/tiny-go-mcp-server/template-http@latest example.com/my-mcp-http my-mcp-http
cd my-mcp-http
go run .
```

Requires Go 1.26+ and a tagged `template-http` release on this repo.

## Environment

| Variable | Description |
|----------|-------------|
| `PORT` | PaaS listen port |
| `TINY_GO_MCP_ADDR` | Default `127.0.0.1:8080` when `PORT` unset |
| `TINY_GO_MCP_API_KEY` | Optional Bearer auth on `/` |
| `TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION` | Set to `1` for ngrok/tunnel testing only |
| `MCP_SERVER_NAME` / `MCP_SERVER_VERSION` | Server identity in MCP initialize |

See [`examples/http-deploy`](../examples/http-deploy) and [`docs/SMITHERY.md`](../docs/SMITHERY.md).
