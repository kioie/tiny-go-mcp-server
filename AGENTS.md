# Agent instructions (Tiny Go MCP Server)

Guidance for AI coding agents working in this repository.

## Commands

- **Test all**: `make test` or `go test -race ./...`
- **Test package**: `go test -race -v ./tinymcp`
- **Test single**: `go test -run TestName ./path/to/pkg -v`
- **Lint**: `make lint` (requires [golangci-lint](https://golangci-lint.run/))
- **Coverage**: `make coverage`
- **Build example server**: `make build` → `./tiny-go-mcp`
- **Minimal example**: `go run ./examples/minimal`

## Layout

- `tinymcp/` — importable library (`github.com/kioie/tiny-go-mcp-server/tinymcp`)
- `cmd/tiny-go-mcp/` — reference MCP server (add, subtract, greet tools)
- `examples/minimal/` — smallest possible server for copy-paste

## API notes

- Tool registration: `tinymcp.RegisterTool(server, name, desc, handler)` — package-level generic, not a method on `TinyServer`.
- Handlers need both `tinymcp` and `github.com/modelcontextprotocol/go-sdk/mcp` imports.
- Optional server logs: set env `TINY_GO_MCP_VERBOSE=1` (stderr only; stdio is reserved for MCP).

## MCP tool conventions

- Tool **names**: lowercase snake_case (`my_tool`)
- Tool **descriptions**: one clear sentence; LLM clients use this to choose tools
- Input args: Go struct with `json` and `jsonschema` tags; schemas are inferred automatically
- Handlers: `func(ctx, *mcp.CallToolRequest, args In) (*mcp.CallToolResult, any, error)`
- Success text responses: use `tinymcp.TextResult("...")`
- Validation errors: return `error` from the handler (SDK maps to `isError` for clients)

## Code style

- Standard Go formatting (`gofmt` / `goimports`)
- Godoc on all exported symbols
- Table-driven tests where helpful
- `context.Context` as first parameter in handlers

## CI expectations

PRs should pass: tests (`-race`), golangci-lint, and build. Security scanning runs via CodeQL on `main`.

## Publishing

Registry metadata: `server.json`. Community listing steps: `docs/DISCOVERY.md`.
