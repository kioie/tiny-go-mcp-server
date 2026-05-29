# TLS for HTTP MCP servers

Production MCP over HTTP must use **HTTPS**. tinymcp serves plain HTTP; terminate TLS in front of the server.

## Recommended: reverse proxy (simplest)

Most PaaS and gateways handle TLS for you:

| Platform | TLS |
|----------|-----|
| Fly.io, Render, Railway | Automatic HTTPS on `*.fly.dev`, `*.onrender.com`, etc. |
| Cloudflare, nginx, Caddy | Terminate TLS; proxy to loopback HTTP |

**Caddy** (automatic Let's Encrypt):

```text
your-host.example.com {
    reverse_proxy 127.0.0.1:8080
}
```

**nginx** snippet:

```nginx
server {
    listen 443 ssl;
    server_name your-host.example.com;
    ssl_certificate     /etc/letsencrypt/live/your-host/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-host/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Run tinymcp on loopback only (`127.0.0.1:8080`) behind the proxy. See [`examples/http-deploy`](../examples/http-deploy).

## Go-native TLS (advanced)

When you must terminate TLS in-process (no reverse proxy):

```go
srv := &http.Server{
    Addr:      ":8443",
    Handler:   mux,
    TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
}
log.Fatal(srv.ListenAndServeTLS("cert.pem", "key.pem"))
```

Use [`tinymcp.ListenAndServeHTTPContext`](HTTP.md) only for **plain HTTP**. For TLS, construct your own `http.Server` with the handler returned from `StreamableHTTPHandler`.

### Automatic certificates (autocert)

For a public domain pointing at your VM:

```go
import "golang.org/x/crypto/acme/autocert"

mux := http.NewServeMux()
mux.Handle("/", mcpHandler)

srv := &http.Server{
    Addr:    ":443",
    Handler: mux,
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        GetCertificate: autocert.Manager{
            Prompt: autocert.AcceptTOS,
            HostPolicy: autocert.HostWhitelist("mcp.example.com"),
            Cache:      autocert.DirCache("/var/lib/autocert"),
        }.GetCertificate,
    },
}
log.Fatal(srv.ListenAndServeTLS("", ""))
```

Requires ports **80** (HTTP-01 challenge) and **443** open. Prefer a managed proxy unless you operate the host yourself.

## Smithery / public URL listing

Smithery requires a **public HTTPS origin**. Use your PaaS URL or a reverse proxy with a valid certificate — self-signed certs will fail client scans.

See also: [`docs/HTTP.md`](HTTP.md), [`docs/SMITHERY.md`](SMITHERY.md).
