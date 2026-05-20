# Nimble Go MCP Server

A lightweight, spec-faithful, and high-performance **Model Context Protocol (MCP) server** written in Go using the official `github.com/modelcontextprotocol/go-sdk`. 

This server is designed to compile into an ultra-small statically linked binary (~5MB) with **zero runtime dependencies** on the host machine. It automatically generates JSON Schemas for its tools using standard Go struct reflection—eliminating manual schema definitions.

[![CI Build & Test](https://github.com/<your-username>/mcp-server/actions/workflows/ci.yml/badge.svg)](https://github.com/<your-username>/mcp-server/actions)

---

## 🛠️ Features

* **Zero Fluff**: 100% self-contained statically compiled binary requiring no Go runtime, node/python runtime, or virtual environment.
* **Automatic Schema Reflection**: Uses standard Go structs with `jsonschema` tags to dynamically define tool parameter configurations.
* **Spec-Faithful**: Leverages the official, Google-collaborative `go-sdk` for fully compliant MCP communications.
* **Highly Optimized**: Complete make task suite for producing stripped, size-optimized binaries.
* **Comprehensive Test Suite**: Exposes full coverage including unit testing for core helpers and mocked handler tests for complete JSON-RPC validation.

---

## 🚀 Exposed Tools

1. **`add`**: Adds two integers together and returns the sum (args: `a`, `b`).
2. **`subtract`**: Subtracts integer B from integer A (args: `a`, `b`).
3. **`greet`**: Personalized greeting generator (args: `name` (required), `greeting` (optional)).

---

## 📦 Getting Started

### 📋 Prerequisites
* [Go 1.25+](https://go.dev/dl/) installed locally (only required for building/testing; running the compiled binary has no prerequisites).

### 📥 Installation
Clone the repository:
```bash
git clone https://github.com/<your-username>/mcp-server.git
cd mcp-server
```

---

## 🛠️ Makefile Commands

The included `Makefile` covers all essential workflows:

### 🧪 Run the Test Suite
Executes all unit and integration mock tests with full verbose logs:
```bash
make test
```

### 🔨 Standard Dev Build
Compiles a standard binary with debug tables:
```bash
make build
```

### ⚡ Highly Optimized Statically Compiled Release Build
Strips DWARF debug symbols and symbol tables to yield an ultra-small, lightning-fast binary (down to ~5MB):
```bash
make release
```

### 🧹 Clean Build Artifacts
```bash
make clean
```

---

## 📉 Ultimate Compression: Down to ~1.8MB
If you want the absolute ultimate in size minimization, you can compress the stripped binary with **[UPX](https://upx.github.io/)** (Ultimate Packer for eXecutables):

```bash
# First create the stripped release binary
make release

# Pack with UPX
upx --best --lzma mcp-server
```
This compresses the executable down to **~1.8MB** while retaining instantaneous startup times.

---

## 🤖 Client Integration

To register this server with your favorite AI environments, compile the binary and add it to your configuration file.

### 1. Claude Desktop
Add this to your `claude_desktop_config.json` (usually at `~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "nimble-go-mcp": {
      "command": "/absolute/path/to/mcp-server"
    }
  }
}
```

### 2. Cursor IDE
1. Open Cursor settings and navigate to **Features > MCP**.
2. Click **+ Add New MCP Server**.
3. Fill in details:
   * **Name**: `nimble-go-mcp`
   * **Type**: `stdio`
   * **Command**: `/absolute/path/to/mcp-server`
