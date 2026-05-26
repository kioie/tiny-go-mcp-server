package main

import (
	"encoding/json"
	"testing"

	"github.com/kioie/tiny-go-mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type serverCard struct {
	ServerInfo struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"serverInfo"`
	Authentication struct {
		Required bool `json:"required"`
	} `json:"authentication"`
	Tools []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"tools"`
	Resources []struct {
		Name        string `json:"name"`
		URI         string `json:"uri"`
		Description string `json:"description"`
		MIMEType    string `json:"mimeType"`
	} `json:"resources"`
	Prompts []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"prompts"`
}

func TestServerCardMatchesLiveCapabilities(t *testing.T) {
	var card serverCard
	if err := json.Unmarshal(serverCardJSON, &card); err != nil {
		t.Fatalf("parse server-card.json: %v", err)
	}
	if card.ServerInfo.Name != "tiny-go-mcp-http" {
		t.Fatalf("serverInfo.name = %q, want tiny-go-mcp-http", card.ServerInfo.Name)
	}
	if card.ServerInfo.Version != serverVersion {
		t.Fatalf("serverInfo.version = %q, want %q (update server-card.json when bumping serverVersion)", card.ServerInfo.Version, serverVersion)
	}
	if card.Authentication.Required {
		t.Fatal("server-card authentication.required must be false for public Smithery demo")
	}

	ctx := t.Context()
	server := tinymcp.NewServer("tiny-go-mcp-http", serverVersion)
	if err := registerCapabilities(server); err != nil {
		t.Fatal(err)
	}

	t1, t2 := mcp.NewInMemoryTransports()
	if _, err := server.RawServer().Connect(ctx, t1, nil); err != nil {
		t.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	cs, err := client.Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != len(card.Tools) {
		t.Fatalf("tools count = %d, card has %d", len(tools.Tools), len(card.Tools))
	}
	for _, want := range card.Tools {
		got := findTool(tools.Tools, want.Name)
		if got == nil {
			t.Fatalf("live server missing tool %q from server-card.json", want.Name)
		}
		if got.Description != want.Description {
			t.Fatalf("tool %q description drift: update server-card.json or registerCapabilities", want.Name)
		}
	}

	resources, err := cs.ListResources(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(resources.Resources) != len(card.Resources) {
		t.Fatalf("resources count = %d, card has %d", len(resources.Resources), len(card.Resources))
	}
	for _, want := range card.Resources {
		got := findResource(resources.Resources, want.URI)
		if got == nil {
			t.Fatalf("live server missing resource %q from server-card.json", want.URI)
		}
		if got.Name != want.Name || got.Description != want.Description || got.MIMEType != want.MIMEType {
			t.Fatalf("resource %q metadata drift: update server-card.json or registerCapabilities", want.URI)
		}
	}

	prompts, err := cs.ListPrompts(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(prompts.Prompts) != len(card.Prompts) {
		t.Fatalf("prompts count = %d, card has %d", len(prompts.Prompts), len(card.Prompts))
	}
	for _, want := range card.Prompts {
		got := findPrompt(prompts.Prompts, want.Name)
		if got == nil {
			t.Fatalf("live server missing prompt %q from server-card.json", want.Name)
		}
		if got.Description != want.Description {
			t.Fatalf("prompt %q description drift: update server-card.json or registerCapabilities", want.Name)
		}
	}
}

func findTool(tools []*mcp.Tool, name string) *mcp.Tool {
	for _, t := range tools {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func findResource(resources []*mcp.Resource, uri string) *mcp.Resource {
	for _, r := range resources {
		if r.URI == uri {
			return r
		}
	}
	return nil
}

func findPrompt(prompts []*mcp.Prompt, name string) *mcp.Prompt {
	for _, p := range prompts {
		if p.Name == name {
			return p
		}
	}
	return nil
}
