.PHONY: all build release test clean

# The binary name
BINARY_NAME=mcp-server

all: test build

# Standard build for development
build:
	go build -o $(BINARY_NAME) main.go

# Stripped, size-optimized release build
release:
	go build -ldflags="-s -w" -o $(BINARY_NAME) main.go

# Run all unit and integration tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean
