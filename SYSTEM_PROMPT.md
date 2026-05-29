# SYSTEM_PROMPT — Tiny Go MCP Server (`tinymcp`)

Token-efficient reference for AI agents **generating or reviewing code** against this library. For repo contribution workflows, see [AGENTS.md](AGENTS.md).

## Identity

- **Module:** `github.com/kioie/tiny-go-mcp-server/tinymcp`
- **Role:** Thin helper on official [`modelcontextprotocol/go-sdk`](https://github.com/modelcontextprotocol/go-sdk) — not a full MCP framework replacement.
- **Stability:** v1.x exported API is stable; see [docs/STABILITY.md](docs/STABILITY.md).

## Required imports

```go
import (
    "github.com/kioie/tiny-go-mcp-server/tinymcp"
    "github.com/modelcontextprotocol/go-sdk/mcp" // handler types, prompts, resources, advanced APIs
)
```

Dual imports are **intentional**. Do not hide `mcp` behind a facade.

## Minimal stdio server

```go
s := tinymcp.NewServer("my-mcp", "1.0.0")
if err := tinymcp.RegisterTool(s, "my_tool", "<description>", handler); err != nil {
    log.Fatal(err)
}
log.Fatal(s.Start())
```

Handler shape:

```go
func handler(ctx context.Context, _ *mcp.CallToolRequest, args myArgs) (*mcp.CallToolResult, any, error) {
    return tinymcp.TextResult("ok"), nil, nil
}
```

Args struct: `json` + `jsonschema` tags (schema inferred automatically).

## Registration API

| API | Use |
|-----|-----|
| `RegisterTool(s, name, desc, handler)` | Typed tool + auto schema |
| `MustRegisterTool(s, name, desc, handler)` | Same; panics at startup on error |
| `RegisterToolDef(s, &mcp.Tool{...}, handler)` | Full tool metadata (annotations) |
| `RegisterTextResource`, `RegisterResource`, `RegisterResourceTemplate` | Read-only context |
| `RegisterPrompt(s, name, desc, args, handler)` | Prompt templates |
| `NewServer(name, ver, tinymcp.WithInstructions(...))` | Optional server instructions |
| `RawServer()` | Escape hatch to underlying `*mcp.Server` |
| `errors.Is(err, tinymcp.ErrNilServer)` etc. | Programmatic registration error handling |
| `HTTPOptions.WithMiddleware(mw...)` | Logging, rate limits on MCP handler |
| `ListenAndServeHTTPContext(ctx, addr, h)` | Graceful shutdown when ctx canceled |

All registration functions **return errors** (nil server, nil handler, invalid tool). Check errors at startup, or use `MustRegister*` helpers.

## Transports

| Client need | API |
|-------------|-----|
| Cursor, Claude Desktop, local subprocess | `s.Start()` (stdio) |
| Remote / gateway | `s.StartHTTP(addr, &tinymcp.HTTPOptions{Stateless: true})` or `StreamableHTTPHandler` |
| Legacy SSE clients | `s.StartSSE(addr, opts)` or `SSEHandler` |

HTTP hardening helpers: `BearerTokenAuth`, `WithCrossOriginProtection`, `DisableLocalhostProtectionFromEnv`. See `examples/http-deploy`.

Logs: `TINY_GO_MCP_VERBOSE=1` → stderr only (stdio is protocol).

## Tool naming and descriptions

- **Names:** `snake_case` (`my_tool`).
- **Descriptions:** LLM clients pick tools from text. Each description must include:
  1. **When to use**
  2. **When not to** (and what to do instead)
  3. **Siblings** — e.g. “use `add` (not `subtract`) for addition tests”

Example:

```go
tinymcp.RegisterTool(s, "add",
    "Demo: adds two integers. Use only to test MCP addition wiring; compute math in the agent instead. Use add (not subtract or greet) when verifying addition.",
    handleAdd)
```

## Do / Don't

| Do | Don't |
|----|-------|
| Return `tinymcp.TextResult("...")` for success text | Panic on registration — check errors or use `MustRegister*` |
| Use `errors.Is(err, tinymcp.ErrNilServer)` for registration failures | String-match registration error messages |
| Use `context.Context` first in handlers | Put secrets in tool responses |
| Bind HTTP demos to loopback unless auth/TLS is configured | Assume tinymcp hides go-sdk types |
| Add `-race` tests for handlers | Ignore SIGTERM — use `ListenAndServeHTTP` or `ListenAndServeHTTPContext` |

## Layout pointers

| Path | Purpose |
|------|---------|
| `examples/minimal/` | Smallest copy-paste server |
| `examples/http-deploy/` | Smithery URL / Fly deploy pattern |
| `examples/resources/` | Resources + prompts |
| `template/` | `gonew` scaffold |
| `docs/HTTP.md` | stdio vs HTTP vs SSE |
| `CHANGELOG.md` | Release history |

## Codegen checklist

1. Both imports (`tinymcp`, `mcp`).
2. Register tools/resources/prompts with error checks before `Start()`.
3. Tool descriptions follow when / when-not / siblings pattern.
4. HTTP production: auth, TLS, timeouts — see `docs/HTTP.md` and `examples/http-deploy`.
5. Match existing style: `gofmt`, table-driven tests, `t.Context()` in tests.
