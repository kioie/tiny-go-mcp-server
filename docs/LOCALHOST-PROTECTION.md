# Localhost protection (DNS rebinding)

Streamable HTTP and legacy SSE handlers in the official go-sdk (and tinymcp wrappers) include **localhost DNS-rebinding protection** by default. Requests that look like they target loopback but resolve to a remote address are rejected.

## When protection applies

- MCP HTTP/SSE bound to `127.0.0.1` or `localhost`
- Browser or tunnel clients connecting through loopback URLs

This mitigates [DNS rebinding](https://modelcontextprotocol.io/specification/2025-11-25/basic/security_best_practices) attacks against local MCP servers.

## Disabling protection

tinymcp exposes two ways to disable — **both are dangerous on untrusted networks**:

| API | Use |
|-----|-----|
| `HTTPOptions.DisableLocalhostProtection = true` | Code-level (avoid in production) |
| `TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1` | Env var; used via `tinymcp.DisableLocalhostProtectionFromEnv()` |

**Approved uses:**

- Local **ngrok** / **cloudflared** tunnel testing before Smithery publish
- CI/integration tests with loopback proxies

**Never** disable on a public bind (`0.0.0.0`) without TLS, auth, and a threat model review.

## Production guidance

| Deployment | Recommendation |
|--------------|----------------|
| PaaS (Fly, Render, Railway) | Bind to platform `PORT`; protection rarely applies — use **HTTPS + Bearer auth** |
| Loopback + reverse proxy | Proxy terminates TLS; tinymcp on `127.0.0.1:8080` — keep protection **enabled** |
| Public HTTP without auth | **Do not ship** — see [`docs/HTTP.md`](HTTP.md) and [`examples/http-deploy`](../examples/http-deploy) |

## Related

- [`docs/HTTP.md`](HTTP.md) — transport security
- [`docs/TLS.md`](TLS.md) — HTTPS termination
- MCP spec: [Security best practices](https://modelcontextprotocol.io/specification/2025-11-25/basic/security_best_practices)
