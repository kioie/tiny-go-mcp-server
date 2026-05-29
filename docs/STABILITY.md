# API stability policy

**Tiny Go MCP Server** publishes the Go module `github.com/kioie/tiny-go-mcp-server/tinymcp`.

## Semver

- **Module tags** use [Semantic Versioning](https://semver.org/): `vMAJOR.MINOR.PATCH`.
- **`go get …/tinymcp@latest`** resolves to the newest tagged release on this repository.

## v1.x guarantee

For all **v1.x.y** releases:

- **No breaking changes** to exported symbols in package `tinymcp` (function signatures, exported types, exported constants).
- **Additive changes** are allowed: new functions, new optional fields on option structs, new examples.
- **Behavior fixes** that correct bugs (including security hardening) may ship in patch releases even if they change runtime behavior of incorrect usage.

Breaking changes to the public `tinymcp` API require a **v2** module path (`/tinymcp/v2`) per [Go module versioning](https://go.dev/ref/mod#module-path-version-suffixes).

## What is stable

| Stable in v1.x | Notes |
|----------------|-------|
| `NewServer`, `RegisterTool`, `RegisterToolDef` | Error-returning registration API |
| `TextResult`, resource/prompt helpers | |
| `Start`, `StartHTTP`, `StartSSE`, HTTP handlers | Options may gain fields |
| `RawServer()` escape hatch | Underlying go-sdk may evolve with dependency bumps |

## What is not guaranteed stable

| Area | Notes |
|------|-------|
| `cmd/tiny-go-mcp` reference server | Demo tools; not a library contract |
| `examples/*` | Illustrative; may change freely |
| Registry metadata (`server.json`, `glama.json`, …) | Operational config, not Go API |
| Transitive `go-sdk` behavior | Pin or test when upgrading `go-sdk` in your module |

## Dependency on go-sdk

`tinymcp` wraps the official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk). Handler code **must** import both `tinymcp` and `mcp` — this is intentional, not a leak to be abstracted away in v1.x.

## Release process

See [CHANGELOG.md](../CHANGELOG.md) and the version bump checklist in [CONTRIBUTING.md](../CONTRIBUTING.md#releases).
