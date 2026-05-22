# Tiny Go MCP Server

[![CI](https://github.com/kioie/tiny-go-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/kioie/tiny-go-mcp-server/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kioie/tiny-go-mcp-server/tinymcp.svg)](https://pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp)
[![Go Report Card](https://goreportcard.com/badge/github.com/kioie/tiny-go-mcp-server/tinymcp)](https://goreportcard.com/report/github.com/kioie/tiny-go-mcp-server/tinymcp)
[![tiny-go-mcp-server MCP server](https://glama.ai/mcp/servers/kioie/tiny-go-mcp-server/badges/score.svg)](https://glama.ai/mcp/servers/kioie/tiny-go-mcp-server)

A lightweight **Model Context Protocol (MCP)** toolkit for Go. Build spec-compliant MCP servers (stdio, streamable HTTP, legacy SSE) with tools, resources, and prompts — minimal boilerplate and automatic JSON Schema from Go structs.

Built on the official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk).

**Requirements:** Go 1.26+ ([download](https://go.dev/dl/)).

---

## Why Tiny?

| | Tiny Go MCP Server | Full frameworks |
|---|---|---|
| **Goal** | Thin wrapper + tiny static binary | Full MCP feature surface |
| **Deps** | Official go-sdk only | Varies |
| **Binary** | ~5MB stripped, no runtime on host | Often larger stacks |
| **Schemas** | Inferred from struct tags | Manual or builder APIs |

Use this project as a **library** (`tinymcp` package) or as a **starting template** (`cmd/tiny-go-mcp`).

### When to use what

| | **tinymcp** (this repo) | **[mcp-go](https://github.com/mark3labs/mcp-go)** | **[go-sdk](https://github.com/modelcontextprotocol/go-sdk)** alone |
|---|---|---|---|
| **Best for** | Thin wrapper, tiny binary, official SDK | Rich helpers, large ecosystem | Full control, no extra layer |
| **Schema** | Struct tags → auto JSON Schema | Builder APIs / helpers | `AddTool` + generics yourself |
| **Transport** | stdio (`Start()`), streamable HTTP (`StartHTTP`), legacy SSE (`StartSSE`) | stdio, SSE, HTTP, … | All transports |
| **Deps** | go-sdk only | Standalone module | go-sdk only |

Choose **tinymcp** when you want the official protocol implementation with minimal boilerplate and a small static server binary.

### Transport

| Method | API | Typical clients |
|--------|-----|-----------------|
| **stdio** (default) | `Start()` | Cursor, Claude Desktop, Windsurf (local subprocess) |
| **Streamable HTTP** | `StartHTTP(addr, opts)` or `StreamableHTTPHandler` | Remote MCP clients, gateways, browser tools |
| **Legacy SSE** | `StartSSE(addr, opts)` or `SSEHandler` | Older clients on MCP 2024-11-05 SSE transport |

`Start()` runs **stdio** (stdin/stdout) — what most local AI clients expect.

For HTTP/SSE, tinymcp wraps the official go-sdk handlers with minimal options:

```go
// Streamable HTTP on :8080 (stateless demo mode)
log.Fatal(server.StartHTTP(":8080", &tinymcp.HTTPOptions{Stateless: true}))

// Or mount on your own mux (auth, TLS, path prefix)
handler, _ := tinymcp.StreamableHTTPHandler(server, nil)
http.Handle("/mcp", handler)
```

See [`docs/HTTP.md`](docs/HTTP.md) and [`examples/http`](examples/http). For advanced session routing or event stores, use `server.RawServer()` with the [go-sdk](https://github.com/modelcontextprotocol/go-sdk) directly.

---

## Quick start (library)

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
	_ = tinymcp.RegisterTool(s, "greet", "Greet someone by name", greet)
	log.Fatal(s.Start())
}

func greet(_ context.Context, _ *mcp.CallToolRequest, args greetArgs) (*mcp.CallToolResult, any, error) {
	return tinymcp.TextResult(fmt.Sprintf("Hello, %s!", args.Name)), nil, nil
}
```

See [`examples/minimal`](examples/minimal) for a runnable copy-paste example.

### Scaffold a new server

Requires tagged module [`template/`](template/) (v1.1.1+):

```bash
go install golang.org/x/tools/cmd/gonew@latest
gonew github.com/kioie/tiny-go-mcp-server/template@latest example.com/my-mcp my-mcp
cd my-mcp
go run .
```

Or copy [`examples/minimal`](examples/minimal) or [`template/`](template/) directly.

---

## Install the example server

Both install paths produce a binary named **`tiny-go-mcp`**:

```bash
# Installs to $(go env GOPATH)/bin/tiny-go-mcp
go install github.com/kioie/tiny-go-mcp-server/cmd/tiny-go-mcp@latest
```

Or build from source (binary in the repo root):

```bash
git clone https://github.com/kioie/tiny-go-mcp-server.git
cd tiny-go-mcp-server
make release   # → ./tiny-go-mcp
```

| Method | Binary name | Typical path |
|--------|-------------|--------------|
| `go install …/cmd/tiny-go-mcp` | `tiny-go-mcp` | `$(go env GOPATH)/bin/tiny-go-mcp` |
| `make build` / `make release` | `tiny-go-mcp` | `./tiny-go-mcp` in the repo |
| `make install` | `tiny-go-mcp` | `$(go env GOPATH)/bin/tiny-go-mcp` |

### Example tools (reference server)

These tools exist for **MCP integration demos**, not production logic. Agents should compute math and write greetings in-chat unless they are explicitly testing tool calls.

| Tool | When to use | When not to / alternative | Arguments |
|------|-------------|---------------------------|-----------|
| `add` | Test that the client can call an addition tool | Real arithmetic → compute locally or use a calculator MCP | `a`, `b` |
| `subtract` | Test subtraction wiring (use instead of `add` for subtraction tests) | Real arithmetic → compute locally | `a`, `b` |
| `greet` | Test a text-returning tool (use instead of `add`/`subtract` for messaging demos) | User-facing hello → reply in the conversation | `name` (required), `greeting` (optional) |

---

## Connect AI clients

MCP servers communicate over **stdio**. Point your client at the compiled binary path.

Template config: [`examples/mcp-client-config.json`](examples/mcp-client-config.json) (copy and set the absolute path to `tiny-go-mcp`).

**Logging:** The protocol uses stdin/stdout. Server logs (if any) go to **stderr** only. Set `TINY_GO_MCP_VERBOSE=1` on the server process to enable startup log lines.

### Cursor

Settings → **Features → MCP** → Add server:

- **Name**: `tiny-go-mcp`
- **Type**: `stdio`
- **Command**: `/absolute/path/to/tiny-go-mcp`

Or add to `.cursor/mcp.json` in your project:

```json
{
  "mcpServers": {
    "tiny-go-mcp": {
      "command": "/absolute/path/to/tiny-go-mcp"
    }
  }
}
```

### Claude Desktop

`~/Library/Application Support/Claude/claude_desktop_config.json` (macOS) or `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "tiny-go-mcp": {
      "command": "/absolute/path/to/tiny-go-mcp"
    }
  }
}
```

### Windsurf / Zed / other stdio clients

Use the same shape: **command** = absolute path to `tiny-go-mcp`, transport = stdio. Refer to your client’s MCP docs for the config file location.

### Tips for LLM-friendly tools

- Use clear **tool names** (`snake_case`) and descriptions that say **when to use**, **when not to**, and **which sibling tool** applies — models pick tools from these and often have overlapping options.
- Add `jsonschema` tags on struct fields so argument docs appear in the schema.
- Return human-readable text via `tinymcp.TextResult` for predictable client display.
- See [AGENTS.md](AGENTS.md) for conventions when extending this repo with AI assistants.

---

## Resources and prompts

Register read-only context and reusable prompt templates alongside tools:

```go
_ = tinymcp.RegisterTextResource(server, "file:///info", "info", "Server metadata", "text/plain", "…")
_ = tinymcp.RegisterPrompt(server, "code_review", "Review code", []*mcp.PromptArgument{
    {Name: "code", Required: true},
}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    return tinymcp.PromptResult("Review", tinymcp.UserPromptMessage("Review:\n"+req.Params.Arguments["code"])), nil
})
```

Runnable example: [`examples/resources`](examples/resources). For dynamic URI templates use `RegisterResourceTemplate`.

---

## Package API

```go
server := tinymcp.NewServer("name", "version")
tinymcp.RegisterTool(server, name, description, handler)     // typed handler, auto schema
tinymcp.RegisterTextResource(server, uri, name, desc, mime, text)
tinymcp.RegisterPrompt(server, name, desc, args, handler)
server.Start()                                              // stdio transport
server.StartHTTP(":8080", &tinymcp.HTTPOptions{})         // streamable HTTP
server.StartSSE(":8080", nil)                             // legacy SSE
tinymcp.StreamableHTTPHandler(server, nil)                // mount on custom http.Server
tinymcp.TextResult("message")                               // tool text helper
tinymcp.TextResource(uri, mime, text)                       // resource read helper
tinymcp.PromptResult(desc, tinymcp.UserPromptMessage("…")) // prompt helper
server.RawServer()                                          // escape hatch to go-sdk
```

Documentation: [pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp](https://pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp)

---

## Development

| Command | Description |
|---------|-------------|
| `make test` | Run tests with race detector |
| `make lint` | golangci-lint |
| `make coverage` | Coverage report |
| `make build` | Dev binary `./tiny-go-mcp` |
| `make release` | Stripped static binary |
| `make install` | `go install` → `$(go env GOPATH)/bin/tiny-go-mcp` |

### Smaller binaries (~1.8MB)

After `make release`, optionally pack with [UPX](https://upx.github.io/):

```bash
upx --best --lzma tiny-go-mcp
```

---

## Project structure

```
tinymcp/           # Library package
cmd/tiny-go-mcp/   # Reference MCP server
examples/minimal/  # Minimal stdio example
examples/http/     # Streamable HTTP example
examples/resources/ # Resources + prompts example
examples/mcp-client-config.json  # Cursor/Claude-style template
docs/HTTP.md       # stdio vs HTTP vs SSE
docs/DISCOVERY.md  # Registries and visibility
server.json        # MCP Registry metadata (publish with mcp-publisher)
.github/workflows/ # CI, lint, CodeQL, releases
```

---

## Releases

Tag a semver version (e.g. `v1.0.0`) to publish stable `go get` versions and trigger [GitHub Releases](https://github.com/kioie/tiny-go-mcp-server/releases) with cross-platform binaries.

```bash
git tag v1.0.0
git push origin v1.0.0
```

## Discovery and registries

See [docs/DISCOVERY.md](docs/DISCOVERY.md) for MCP Registry (`server.json`), awesome lists, and community directories. Listing copy and launch posts: [docs/SUBMISSIONS.md](docs/SUBMISSIONS.md). For Glama hosting with Docker, see [docs/GLAMA.md](docs/GLAMA.md).

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). CI runs tests, lint, and CodeQL; Dependabot keeps Go and Actions dependencies updated.

## License

MIT — see [LICENSE](LICENSE).
