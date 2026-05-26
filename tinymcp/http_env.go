package tinymcp

import "os"

// DisableLocalhostProtectionFromEnv reports whether TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION=1.
// Use only for local ngrok/tunnel testing on loopback listeners.
func DisableLocalhostProtectionFromEnv() bool {
	return os.Getenv("TINY_GO_MCP_DISABLE_LOCALHOST_PROTECTION") == "1"
}
