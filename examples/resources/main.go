// Resources and prompts example over stdio.
//
// Run: go run ./examples/resources
package main

import (
	"context"
	"log"
	"os"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	server := tinymcp.NewServer("tiny-go-mcp-resources", "1.0.0")

	if err := tinymcp.RegisterTextResource(server,
		"file:///info",
		"info",
		"Static server info for MCP clients",
		"text/plain",
		"Tiny Go MCP Server — resources/prompts example",
	); err != nil {
		log.Fatal(err)
	}

	if err := tinymcp.RegisterPrompt(server,
		"code_review",
		"Review code with a structured user message",
		[]*mcp.PromptArgument{{Name: "code", Required: true, Description: "Source code to review"}},
		codeReview,
	); err != nil {
		log.Fatal(err)
	}

	if os.Getenv("TINY_GO_MCP_VERBOSE") != "" {
		log.Println("Tiny Go MCP resources/prompts example (stdio)")
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}

func codeReview(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	code := req.Params.Arguments["code"]
	return tinymcp.PromptResult("Code review",
		tinymcp.UserPromptMessage("Review this code for bugs and clarity:\n\n"+code),
	), nil
}
