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
2. Make your changes (library in `tinymcp/`, example server in `cmd/tiny-go-mcp/`).
3. Run checks locally:

   ```bash
   make test
   make lint    # if golangci-lint is installed
   ```

4. Open a pull request against `main`.

## Automated PR review

Non-draft pull requests trigger an [**OpenRouter**](https://openrouter.ai) review via [`.github/workflows/openrouter-pr-review.yml`](.github/workflows/openrouter-pr-review.yml).

**Repository admin setup (one time):**

1. Create an API key at [openrouter.ai/settings/keys](https://openrouter.ai/settings/keys).
2. Add Actions secret `OPENROUTER_API_KEY` (`Settings → Secrets and variables → Actions`).
3. Optional repository variable `OPENROUTER_MODEL` (default: `anthropic/claude-3.5-haiku`).

Reviews run on PR open/update. Re-run by commenting `/review` on the PR.

Review criteria follow [AGENTS.md](AGENTS.md). The workflow uses `continue-on-error: true` so transient model/API outages do not block merges.

## Adding tools to the example server

1. Define an arguments struct with `json` and `jsonschema` tags.
2. Implement a handler matching `mcp.ToolHandlerFor[YourArgs, any]`.
3. Register with `server.RegisterTool` in `RegisterTools`.
4. Add unit tests in `cmd/tiny-go-mcp/main_test.go`.

## Releases

Tag with semver (`v1.0.0`). The [Release workflow](.github/workflows/release.yml) builds cross-platform `tiny-go-mcp` binaries and attaches them to the GitHub release.

Before tagging:

1. Move `[Unreleased]` entries in [CHANGELOG.md](CHANGELOG.md) into a new version section with date.
2. Confirm [docs/STABILITY.md](docs/STABILITY.md) still matches any API changes.
3. Bump version strings in `server.json`, `glama.json`, and deploy examples if applicable.

## Questions

Open a [GitHub issue](https://github.com/kioie/tiny-go-mcp-server/issues) for bugs or feature requests.
