# Migrating to tinymcp v1.2.x

Guide for users upgrading from **v1.1.x** to **v1.2.x**. Full history: [CHANGELOG.md](../CHANGELOG.md).

## What changed

### Sentinel errors (registration)

Registration functions now return stable sentinel errors. Use `errors.Is` instead of matching error strings:

| Error | When |
|-------|------|
| `tinymcp.ErrNilServer` | Nil `*TinyServer` or `Start()` on nil |
| `tinymcp.ErrNilTool` | `RegisterToolDef` with nil `*mcp.Tool` |
| `tinymcp.ErrNilHandler` | Nil resource/prompt handler |
| `tinymcp.ErrRegistrationFailed` | go-sdk panic during registration (wrapped) |

```go
if err := tinymcp.RegisterTool(s, "x", "y", handler); err != nil {
    switch {
    case errors.Is(err, tinymcp.ErrNilServer):
        // ...
    case errors.Is(err, tinymcp.ErrRegistrationFailed):
        // invalid schema, bad template, etc.
    default:
        log.Fatal(err)
    }
}
```

### MustRegister helpers

Optional panic-at-startup variants: `MustRegisterTool`, `MustRegisterToolDef`, `MustRegisterPrompt`, `MustRegisterResource`, `MustRegisterResourceTemplate`, `MustRegisterTextResource`.

```go
tinymcp.MustRegisterTool(s, "greet", "Greet someone", greet)
```

### HTTP graceful shutdown

`ListenAndServeHTTP` now drains on **SIGINT** / **SIGTERM**. For custom lifecycle:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()
_ = tinymcp.ListenAndServeHTTPContext(ctx, addr, handler)
```

Or `StartHTTPContext` / `StartSSEContext`.

### HTTP middleware

```go
opts := (&tinymcp.HTTPOptions{Stateless: true}).WithMiddleware(requestLogger)
handler, _ := tinymcp.StreamableHTTPHandler(server, opts)
```

Apply auth/CORS **outside** `HTTPOptions.Middleware` (see [docs/HTTP.md](HTTP.md)).

## No action required

- Handler signatures unchanged (`mcp.CallToolRequest`, dual imports).
- `RegisterTool` / `Start()` / `TextResult` signatures unchanged.
- v1.x stability policy unchanged: [docs/STABILITY.md](STABILITY.md).

## Tagging

After upgrading, pin explicitly:

```bash
go get github.com/kioie/tiny-go-mcp-server/tinymcp@v1.2.0
```
