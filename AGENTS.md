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

- `tinymcp/` — importable library (`github.com/kioie/mcp-server/tinymcp`)
- `cmd/tinymcp/` — reference MCP server (add, subtract, greet tools)
- `examples/minimal/` — smallest possible server for copy-paste

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
