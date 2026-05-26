# Security

## Reporting a vulnerability

Please report security issues privately via [GitHub Security Advisories](https://github.com/kioie/tiny-go-mcp-server/security/advisories/new) rather than in public issues.

We will acknowledge reports promptly and work on fixes as needed.

## Scope

This project is an MCP server **library** and reference binaries. The default transport is **stdio** (no network listener). Optional HTTP/SSE helpers and deploy examples (`examples/http`, `examples/http-deploy`) **do** expose network services when you run them — treat public deployments like any other web API: add authentication, TLS, and rate limiting. Review tool handlers you add — they run with the privileges of the host process that launches the MCP server.
