package main

import (
	"context"
	"fmt"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const infoResourceURI = "file:///info"

// RegisterResourcesAndPrompts adds demo resource and prompt handlers for MCP client tests.
func RegisterResourcesAndPrompts(server *tinymcp.TinyServer) error {
	if server == nil {
		return fmt.Errorf("cannot register resources on a nil server")
	}

	if err := tinymcp.RegisterTextResource(server,
		infoResourceURI,
		"info",
		"Demo: static server metadata. Use to test resources/read; do not use for live config—read env or APIs instead.",
		"text/plain",
		"Tiny Go MCP Server v1.1.3 — tinymcp reference (tools, resources, prompts)",
	); err != nil {
		return fmt.Errorf("failed to register info resource: %w", err)
	}

	if err := tinymcp.RegisterPrompt(server,
		"code_review",
		"Demo: code review prompt template. Use to test prompts/get; do not use for real reviews—use your agent workflow instead.",
		[]*mcp.PromptArgument{{Name: "code", Required: true, Description: "Source code to review"}},
		handleCodeReviewPrompt,
	); err != nil {
		return fmt.Errorf("failed to register code_review prompt: %w", err)
	}

	return nil
}

func handleCodeReviewPrompt(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	code := req.Params.Arguments["code"]
	if code == "" {
		return nil, tinymcp.RequiredPromptArgument("code")
	}
	return tinymcp.PromptResult("Code review",
		tinymcp.UserPromptMessage("Review this code for bugs and clarity:\n\n"+code),
	), nil
}
