# Quick start

Get a working MCP server in minutes. For AI code generation, also read [SYSTEM_PROMPT.md](../SYSTEM_PROMPT.md) (token-efficient API cheat sheet).

## Prerequisites

- Go 1.26+ ([download](https://go.dev/dl/))

## 1. Stdio server (Cursor, Claude Desktop)

```bash
go get github.com/kioie/tiny-go-mcp-server/tinymcp@latest
```

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type greetArgs struct {
	Name string `json:"name" jsonschema:"Person to greet"`
}

func main() {
	s := tinymcp.NewServer("my-mcp", "1.0.0")
	if err := tinymcp.RegisterTool(s, "greet", "Greet someone by name", greet); err != nil {
		log.Fatal(err)
	}
	log.Fatal(s.Start())
}

func greet(_ context.Context, _ *mcp.CallToolRequest, args greetArgs) (*mcp.CallToolResult, any, error) {
	return tinymcp.TextResult(fmt.Sprintf("Hello, %s!", args.Name)), nil, nil
}
```

Runnable copy: [`examples/minimal`](../examples/minimal).

**Scaffold with gonew:**

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
cd my-mcp && go run .
```

## 2. Streamable HTTP (remote / Smithery)

```go
handler, err := tinymcp.StreamableHTTPHandler(server, &tinymcp.HTTPOptions{
	Stateless: true,
})
if err != nil {
	log.Fatal(err)
}
handler = tinymcp.WithCrossOriginProtection(handler)
handler = tinymcp.BearerTokenAuth(os.Getenv("TINY_GO_MCP_API_KEY"), handler)

mux := http.NewServeMux()
mux.Handle("/", handler)
log.Fatal(tinymcp.ListenAndServeHTTP("127.0.0.1:8080", mux))
```

Full deploy example with health check and server card: [`examples/http-deploy`](../examples/http-deploy). Smithery guide: [SMITHERY.md](./SMITHERY.md).

## 3. Connect a local client

Install the reference server:

```bash
go install github.com/kioie/tiny-go-mcp-server/cmd/tiny-go-mcp@latest
```

**Cursor** — add to MCP settings:

```json
{
  "mcpServers": {
    "tiny-go-mcp": {
      "command": "tiny-go-mcp"
    }
  }
}
```

## Next steps

| Goal | Doc |
|------|-----|
| stdio vs HTTP vs SSE | [HTTP.md](./HTTP.md) |
| Tool description conventions | [SYSTEM_PROMPT.md](../SYSTEM_PROMPT.md), [AGENTS.md](../AGENTS.md) |
| API stability | [STABILITY.md](./STABILITY.md) |
| Release history | [CHANGELOG.md](../CHANGELOG.md) |
| pkg.go.dev examples | [tinymcp package](https://pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp#pkg-examples) |
