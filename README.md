# Tiny Go MCP Server

[![CI](https://github.com/kioie/tiny-go-mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/kioie/tiny-go-mcp-server/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/kioie/tiny-go-mcp-server/tinymcp.svg)](https://pkg.go.dev/github.com/kioie/tiny-go-mcp-server/tinymcp)

A lightweight **Model Context Protocol (MCP)** toolkit for Go. Build spec-compliant, stdio-based MCP servers that AI clients can discover and call — with minimal boilerplate and automatic JSON Schema generation from Go structs.

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

| Tool | Description | Arguments |
|------|-------------|-----------|
| `add` | Add two integers | `a`, `b` |
| `subtract` | Subtract integers | `a`, `b` |
| `greet` | Personalized greeting | `name` (required), `greeting` (optional) |

---

## Connect AI clients

MCP servers communicate over **stdio**. Point your client at the compiled binary path.

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

`~/Library/Application Support/Claude/claude_desktop_config.json` (macOS):

```json
{
  "mcpServers": {
    "tiny-go-mcp": {
      "command": "/absolute/path/to/tiny-go-mcp"
    }
  }
}
```

### Tips for LLM-friendly tools

- Use clear **tool names** (`snake_case`) and **one-sentence descriptions** — models pick tools from these.
- Add `jsonschema` tags on struct fields so argument docs appear in the schema.
- Return human-readable text via `tinymcp.TextResult` for predictable client display.
- See [AGENTS.md](AGENTS.md) for conventions when extending this repo with AI assistants.

---

## Package API

```go
server := tinymcp.NewServer("name", "version")
tinymcp.RegisterTool(server, name, description, handler) // typed handler, auto schema
server.Start()                                              // stdio transport
tinymcp.TextResult("message")                               // helper for text tools
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
examples/minimal/  # Minimal example
.github/workflows/ # CI, lint, CodeQL, releases
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). CI runs tests, lint, and CodeQL; Dependabot keeps Go and Actions dependencies updated.

## License

MIT — see [LICENSE](LICENSE).
