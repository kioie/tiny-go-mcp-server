# Contributing

Thank you for contributing to **Tiny Go MCP Server**.

## Prerequisites

- Go 1.26+ ([download](https://go.dev/dl/))
- Optional: [golangci-lint](https://golangci-lint.run/) for local linting

## Setup

```bash
git clone https://github.com/kioie/tiny-go-mcp-server.git
cd tiny-go-mcp-server
go mod download
```

## Workflow

1. Create a branch from `main`.
2. Make your changes (library in `tinymcp/`, example server in `cmd/tinymcp/`).
3. Run checks locally:

   ```bash
   make test
   make lint    # if golangci-lint is installed
   ```

4. Open a pull request against `main`.

## Adding tools to the example server

1. Define an arguments struct with `json` and `jsonschema` tags.
2. Implement a handler matching `mcp.ToolHandlerFor[YourArgs, any]`.
3. Register with `server.RegisterTool` in `RegisterTools`.
4. Add unit tests in `cmd/tinymcp/main_test.go`.

## Releases

Tag with semver (`v1.0.0`). The [Release workflow](.github/workflows/release.yml) builds cross-platform `tiny-go-mcp` binaries and attaches them to the GitHub release.

## Questions

Open a [GitHub issue](https://github.com/kioie/tiny-go-mcp-server/issues) for bugs or feature requests.
