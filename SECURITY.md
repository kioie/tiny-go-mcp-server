# Security

## Reporting a vulnerability

Please report security issues privately via [GitHub Security Advisories](https://github.com/kioie/tiny-go-mcp-server/security/advisories/new) rather than in public issues.

We will acknowledge reports promptly and work on fixes as needed.

## Scope

This project is a local stdio MCP server library and reference binary. It does not expose network services by default. Review tool handlers you add — they run with the privileges of the host process that launches the MCP server.
