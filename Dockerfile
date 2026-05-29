# Tiny Go MCP Server — stdio MCP for Glama, Docker, and OCI registries.
# Glama: connect GitHub repo; build context . ; uses ENTRYPOINT below.
# Local: docker build -t tiny-go-mcp . && docker run -i --rm tiny-go-mcp

FROM golang:1.26-alpine AS builder

ARG TARGETARCH

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build \
    -ldflags="-s -w" \
    -o /out/tiny-go-mcp \
    ./cmd/tiny-go-mcp

FROM alpine:3.21

# Required by MCP Registry for OCI package verification (server.json packages[].identifier).
LABEL io.modelcontextprotocol.server.name="io.github.kioie/tiny-go-mcp"

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/tiny-go-mcp /usr/local/bin/tiny-go-mcp

# MCP stdio: protocol on stdin/stdout; keep stdout clean (logs → stderr only if TINY_GO_MCP_VERBOSE=1).
ENV TINY_GO_MCP_VERBOSE=

USER nobody

ENTRYPOINT ["/usr/local/bin/tiny-go-mcp"]
