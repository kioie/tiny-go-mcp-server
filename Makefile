.PHONY: all build release test lint coverage clean install

BINARY_NAME=tiny-go-mcp
CMD=./cmd/tiny-go-mcp
LDFLAGS=-s -w

all: test build

build:
	go build -o $(BINARY_NAME) $(CMD)

release:
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(CMD)

install:
	go install $(CMD)

test:
	go test -race -v ./...

lint:
	golangci-lint run ./...

coverage:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out

clean:
	rm -f $(BINARY_NAME) tinymcp.exe coverage.out
	go clean
