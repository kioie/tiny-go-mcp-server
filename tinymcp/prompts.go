package tinymcp

import (
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPrompt adds a prompt template. arguments may be nil for prompts with no parameters.
func RegisterPrompt(s *TinyServer, name, description string, arguments []*mcp.PromptArgument, handler mcp.PromptHandler) error {
	srv, err := rawServer(s)
	if err != nil {
		return err
	}
	srv.AddPrompt(&mcp.Prompt{
		Name:        name,
		Description: description,
		Arguments:   arguments,
	}, handler)
	return nil
}

// PromptResult builds a GetPromptResult from prompt messages.
func PromptResult(description string, messages ...*mcp.PromptMessage) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: description,
		Messages:    messages,
	}
}

// UserPromptMessage returns a user-role prompt message with text content.
func UserPromptMessage(text string) *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Role:    "user",
		Content: &mcp.TextContent{Text: text},
	}
}

// AssistantPromptMessage returns an assistant-role prompt message with text content.
func AssistantPromptMessage(text string) *mcp.PromptMessage {
	return &mcp.PromptMessage{
		Role:    "assistant",
		Content: &mcp.TextContent{Text: text},
	}
}
