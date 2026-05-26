package tinymcp

import (
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// RegisterPrompt adds a prompt template. arguments may be nil for prompts with no parameters.
func RegisterPrompt(s *TinyServer, name, description string, arguments []*mcp.PromptArgument, handler mcp.PromptHandler) error {
	if handler == nil {
		return fmt.Errorf("cannot register prompt %q with a nil handler", name)
	}
	srv, err := rawServer(s)
	if err != nil {
		return err
	}
	return registerRecover(fmt.Sprintf("prompt %q", name), func() {
		srv.AddPrompt(&mcp.Prompt{
			Name:        name,
			Description: description,
			Arguments:   arguments,
		}, handler)
	})
}

// PromptResult builds a GetPromptResult from prompt messages.
// With no messages, returns an empty array (not JSON null).
func PromptResult(description string, messages ...*mcp.PromptMessage) *mcp.GetPromptResult {
	if messages == nil {
		messages = []*mcp.PromptMessage{}
	}
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
