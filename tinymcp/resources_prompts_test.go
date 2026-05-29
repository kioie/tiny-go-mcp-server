package tinymcp

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/jsonrpc"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterTextResource(t *testing.T) {
	ctx := t.Context()
	s := NewServer("test", "1.0.0")
	const uri = "file:///info"
	if err := RegisterTextResource(s, uri, "info", "Server info", "text/plain", "tinymcp demo"); err != nil {
		t.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	res, err := cs.ReadResource(ctx, &mcp.ReadResourceParams{URI: uri})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Contents) != 1 || res.Contents[0].Text != "tinymcp demo" {
		t.Fatalf("unexpected contents: %#v", res.Contents)
	}
}

func TestRegisterPrompt(t *testing.T) {
	ctx := t.Context()
	s := NewServer("test", "1.0.0")
	if err := RegisterPrompt(s, "greet", "Greeting prompt", []*mcp.PromptArgument{
		{Name: "name", Required: true},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		name := req.Params.Arguments["name"]
		return PromptResult("Say hello",
			UserPromptMessage("Hello, "+name+"!"),
		), nil
	}); err != nil {
		t.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	got, err := cs.GetPrompt(ctx, &mcp.GetPromptParams{
		Name:      "greet",
		Arguments: map[string]string{"name": "Ada"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(got.Messages))
	}
	tc, ok := got.Messages[0].Content.(*mcp.TextContent)
	if !ok || tc.Text != "Hello, Ada!" {
		t.Fatalf("unexpected message: %#v", got.Messages[0].Content)
	}
}

func TestRegisterResource_nilServer(t *testing.T) {
	var s *TinyServer
	if err := RegisterResource(s, "file:///x", "x", "", "", nil); !errors.Is(err, ErrNilServer) {
		t.Errorf("RegisterResource(nil): got %v, want ErrNilServer", err)
	}
	if err := RegisterTextResource(s, "file:///x", "x", "", "", ""); !errors.Is(err, ErrNilServer) {
		t.Errorf("RegisterTextResource(nil): got %v, want ErrNilServer", err)
	}
	if err := RegisterPrompt(s, "p", "", nil, nil); !errors.Is(err, ErrNilServer) {
		t.Errorf("RegisterPrompt(nil): got %v, want ErrNilServer", err)
	}
}

func TestRegister_nilHandlerOnServer(t *testing.T) {
	s := NewServer("test", "1.0.0")
	if err := RegisterPrompt(s, "p", "demo", nil, nil); !errors.Is(err, ErrNilHandler) {
		t.Fatalf("RegisterPrompt(nil handler): got %v, want ErrNilHandler", err)
	}
	if err := RegisterResource(s, "file:///x", "x", "", "text/plain", nil); !errors.Is(err, ErrNilHandler) {
		t.Fatalf("RegisterResource(nil handler): got %v, want ErrNilHandler", err)
	}
}

func TestRegisterResourceTemplate_invalidTemplate(t *testing.T) {
	s := NewServer("test", "1.0.0")
	err := RegisterResourceTemplate(s, "{{{", "bad", "", "text/plain", func(context.Context, *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return nil, nil
	})
	if !errors.Is(err, ErrRegistrationFailed) {
		t.Fatalf("expected ErrRegistrationFailed, got %v", err)
	}
}

func TestPromptResult_emptyMessagesNotNull(t *testing.T) {
	b, err := json.Marshal(PromptResult("desc"))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) == `{"description":"desc","messages":null}` {
		t.Fatalf("messages encoded as null: %s", b)
	}
	if string(b) != `{"description":"desc","messages":[]}` {
		t.Fatalf("unexpected JSON: %s", b)
	}
}

func TestRequiredPromptArgument(t *testing.T) {
	ctx := t.Context()
	s := NewServer("test", "1.0.0")
	if err := RegisterPrompt(s, "review", "Review code", []*mcp.PromptArgument{
		{Name: "code", Required: true},
	}, func(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		code := req.Params.Arguments["code"]
		if code == "" {
			return nil, RequiredPromptArgument("code")
		}
		return PromptResult("Review", UserPromptMessage(code)), nil
	}); err != nil {
		t.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := s.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	_, err = cs.GetPrompt(ctx, &mcp.GetPromptParams{Name: "review"})
	if err == nil {
		t.Fatal("expected error for missing required argument")
	}
	var rpcErr *jsonrpc.Error
	if !errors.As(err, &rpcErr) || rpcErr.Code != jsonrpc.CodeInvalidParams {
		t.Fatalf("expected InvalidParams, got %v", err)
	}
}
