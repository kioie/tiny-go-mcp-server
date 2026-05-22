# Build from repository root (tinymcp module lives at repo root):
#   docker build -f examples/http-deploy/Dockerfile .
FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY tinymcp/ tinymcp/
COPY examples/http-deploy/ examples/http-deploy/

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /out/server \
    ./examples/http-deploy

FROM alpine:3.21

RUN apk add --no-cache ca-certificates

COPY --from=builder /out/server /usr/local/bin/server

ENV TINY_GO_MCP_VERBOSE=
USER nobody
EXPOSE 8080
CMD ["/usr/local/bin/server"]
