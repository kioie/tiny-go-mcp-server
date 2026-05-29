package tinymcp

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func FuzzEchoHandler_args(f *testing.F) {
	f.Add([]byte(`{"message":"hello"}`))
	f.Add([]byte(`{"message":""}`))
	f.Add([]byte(`{}`))

	handler := func(_ context.Context, _ *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, any, error) {
		return TextResult(args.Message), nil, nil
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		var args echoArgs
		if err := json.Unmarshal(data, &args); err != nil {
			return
		}
		_, _, _ = handler(t.Context(), nil, args)
	})
}

func FuzzRegisterTool_invalidJSON(f *testing.F) {
	f.Add([]byte(`{"message":"x"}`))
	f.Add([]byte(`not-json`))

	s := NewServer("fuzz", "1.0.0")
	if err := RegisterTool(s, "echo", "Echo", func(_ context.Context, _ *mcp.CallToolRequest, args echoArgs) (*mcp.CallToolResult, any, error) {
		return TextResult(args.Message), nil, nil
	}); err != nil {
		f.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	ctx := context.Background()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		f.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "fuzz-client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		f.Fatal(err)
	}
	defer cs.Close()

	f.Fuzz(func(t *testing.T, data []byte) {
		if !json.Valid(data) {
			return
		}
		var args map[string]any
		if err := json.Unmarshal(data, &args); err != nil {
			return
		}
		_, _ = cs.CallTool(t.Context(), &mcp.CallToolParams{
			Name:      "echo",
			Arguments: args,
		})
	})
}
