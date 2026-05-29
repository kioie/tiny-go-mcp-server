// lint-tools checks MCP tool descriptions follow AGENTS.md conventions:
// when to use, when not to, and sibling tools.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type toolDesc struct {
	file string
	name string
	desc string
}

var (
	registerToolRe = regexp.MustCompile(`(?s)RegisterTool\(\s*[^,]+,\s*"([^"]+)"\s*,\s*"((?:\\.|[^"\\])*)"`)
	toolBlockRe    = regexp.MustCompile(`(?s)&mcp\.Tool\{.*?Name:\s*"([^"]+)".*?Description:\s*"((?:\\.|[^"\\])*)"`)
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fatal(err)
	}

	targets := []string{
		filepath.Join(root, "cmd/tiny-go-mcp/main.go"),
		filepath.Join(root, "examples/http-deploy/main.go"),
	}

	var failed int
	for _, path := range targets {
		data, err := os.ReadFile(path)
		if err != nil {
			fatal(err)
		}
		content := string(data)
		for _, td := range extractDescriptions(path, content) {
			if msg := validate(td); msg != "" {
				fmt.Fprintf(os.Stderr, "%s: tool %q: %s\n", td.file, td.name, msg)
				failed++
			}
		}
	}

	if failed > 0 {
		fmt.Fprintf(os.Stderr, "lint-tools: %d tool description(s) failed — see AGENTS.md MCP tool conventions\n", failed)
		os.Exit(1)
	}
	fmt.Println("lint-tools: OK")
}

func extractDescriptions(path, content string) []toolDesc {
	var out []toolDesc
	for _, m := range registerToolRe.FindAllStringSubmatch(content, -1) {
		out = append(out, toolDesc{file: path, name: m[1], desc: unescape(m[2])})
	}
	for _, m := range toolBlockRe.FindAllStringSubmatch(content, -1) {
		out = append(out, toolDesc{file: path, name: m[1], desc: unescape(m[2])})
	}
	return out
}

func validate(td toolDesc) string {
	d := strings.ToLower(td.desc)
	switch {
	case !strings.Contains(d, "use "):
		return "missing when-to-use hint (include \"Use ...\")"
	case !strings.Contains(d, "not"):
		return "missing when-not-to hint (include \"do not use\" or \"not X\")"
	case !strings.Contains(d, "(not ") && !strings.Contains(d, " not "):
		return "missing sibling hint (e.g. \"Use add (not subtract)\")"
	}
	return ""
}

func unescape(s string) string {
	s = strings.ReplaceAll(s, `\"`, `"`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	return s
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "lint-tools: %v\n", err)
	os.Exit(2)
}
