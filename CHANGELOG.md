# Changelog

All notable changes to **Tiny Go MCP Server** are documented here.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) for the Go module `github.com/kioie/tiny-go-mcp-server/tinymcp`.

Stability guarantees for the public API are in [docs/STABILITY.md](docs/STABILITY.md).

## [Unreleased]

### Added

- OpenRouter automated PR review workflow (`.github/workflows/openrouter-pr-review.yml`).
- README Philosophy section (thin helper on official go-sdk, intentional dual imports).
- Phase 2 discoverability: `docs/QUICKSTART.md`, README tinymcp vs go-sdk comparison, pkg.go.dev `Example*` functions.
- GitHub issue templates and CHANGELOG validation workflow.
- awesome-go submission draft in `docs/SUBMISSIONS.md`.

### Changed

- Dependabot bumps: go-sdk v1.6.1, GitHub Actions updates.
- Expanded release version-bump checklist in `CONTRIBUTING.md`.
- HTTP middleware stack documentation in `examples/http-deploy/README.md`.

### Documentation

- awesome-mcp-servers listing status in `docs/DISCOVERY.md` and `docs/SUBMISSIONS.md`.

### Added (API — v1.2.0)

- `MustRegisterTool`, `MustRegisterToolDef`, `MustRegisterPrompt`, `MustRegisterResource`, `MustRegisterResourceTemplate`, `MustRegisterTextResource` — panic-at-startup registration helpers.
- Sentinel errors: `ErrNilServer`, `ErrNilTool`, `ErrNilHandler`, `ErrRegistrationFailed` (use `errors.Is`).
- `HTTPOptions.Middleware` and `HTTPOptions.WithMiddleware` for handler wrapping (logging, rate limits).
- Graceful HTTP shutdown: `ListenAndServeHTTPContext`, `StartHTTPContext`, `StartSSEContext`; `ListenAndServeHTTP` drains on SIGINT/SIGTERM.

### Changed (API)

- Registration and HTTP handler errors now return sentinel errors instead of ad-hoc strings (check with `errors.Is`, not string matching).

### Added (testing — v1.2.0)

- Stdio subprocess integration test for reference server (`cmd/tiny-go-mcp`).
- Fuzz tests for tool argument JSON decoding (`tinymcp/fuzz_test.go`).
- HTTP integration test for middleware + context shutdown.
- Benchmarks: `RegisterTool`, `TextResult`, streamable HTTP initialize.
- CI coverage artifact upload and 70% gate on `tinymcp` + `cmd/tiny-go-mcp`.
- v1.2 migration guide: `docs/MIGRATION-v1.2.md`.

### Added (deployment — v1.2.0)

- GHCR Docker images built for **linux/amd64** and **linux/arm64** on release.
- `docs/TLS.md` — HTTPS via reverse proxy, nginx/Caddy, and optional Go autocert.
- HTTP `ReadTimeout` (30s) on `ListenAndServeHTTP` / `ListenAndServeHTTPContext` for slow-loris protection.

### Added (agent tooling — v1.2.0)

- `make lint-tools` — validates MCP tool descriptions (when / when-not / siblings) in reference servers.
- `template-http/` gonew scaffold for streamable HTTP deploy (Smithery/Fly path).
- `docs/LOCALHOST-PROTECTION.md` — DNS rebinding security advisory.
- CI: multi-arch Docker build verification (no push); Fly.io/Render examples in `docs/TLS.md`.

## [1.1.3] - 2026-05-26

### Added

- `NewServer` / `NewServerWithOptions` optional `ServerOptions` passthrough.
- HTTP middleware: `BearerTokenAuth`, `WithCrossOriginProtection`, `DisableLocalhostProtectionFromEnv`.
- Env-gated localhost protection disable for `examples/http-deploy`.
- Server-card drift test in `examples/http-deploy`.

### Fixed

- Registration hardening: `registerRecover`, nil server/handler checks, ordered error reporting.
- HTTP deploy docs and Fly/Smithery security guidance.
- Audit polish (low-priority findings L1, L2, L6, L5, L7, M2).

### Changed

- Version strings aligned across `server.json`, `glama.json`, `smithery.yaml`, and deploy examples.

## [1.1.2] - 2026-05-22

### Added

- `examples/http-deploy` for Smithery URL listing and public HTTP demo.
- Smithery MCPB publish assets and live stdio listing docs.
- Official logo, banner, and branding assets (`docs/BRANDING.md`).

### Documentation

- Adoption listing status in `docs/SUBMISSIONS.md`.
- r/mcp launch post recorded in `docs/DISCOVERY.md`.

## [1.1.1] - 2026-05-22

### Added

- `template/` module for `gonew github.com/kioie/tiny-go-mcp-server/template@latest`.
- Smithery and second-tier directory submission copy in `docs/SUBMISSIONS.md`.

## [1.1.0] - 2026-05-22

### Added

- Streamable HTTP and legacy SSE transports (`StartHTTP`, `StartSSE`, `StreamableHTTPHandler`, `SSEHandler`).
- Resources and prompts helpers (`RegisterTextResource`, `RegisterResource`, `RegisterPrompt`, …).
- `Dockerfile`, `glama.json`, and Glama hosting docs.
- MCP Registry `server.json` and GHCR image push on release.
- Reference server resources/prompts (`file:///info`, `code_review` prompt).

### Documentation

- Agent-oriented tool description guidance in README and `AGENTS.md`.
- Glama score badge on README.

## [1.0.0] - 2026-05-20

### Added

- Initial stable release of `tinymcp` on official `modelcontextprotocol/go-sdk`.
- Stdio reference server `cmd/tiny-go-mcp` (add, subtract, greet tools).
- Struct-tag JSON Schema inference via `RegisterTool`.
- `examples/minimal`, CI (test, lint, CodeQL), and cross-platform release binaries.

[Unreleased]: https://github.com/kioie/tiny-go-mcp-server/compare/v1.1.3...HEAD
[1.1.3]: https://github.com/kioie/tiny-go-mcp-server/compare/v1.1.2...v1.1.3
[1.1.2]: https://github.com/kioie/tiny-go-mcp-server/compare/v1.1.1...v1.1.2
[1.1.1]: https://github.com/kioie/tiny-go-mcp-server/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/kioie/tiny-go-mcp-server/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/kioie/tiny-go-mcp-server/releases/tag/v1.0.0
