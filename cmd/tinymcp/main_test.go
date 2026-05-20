package main

import (
	"testing"

	"github.com/kioie/mcp-server/tinymcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// =============================================================================
// CORE BUSINESS LOGIC TESTS
// =============================================================================

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive numbers", 5, 3, 8},
		{"negative numbers", -5, -3, -8},
		{"mixed numbers", -5, 10, 5},
		{"zeros", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Add(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Add(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"positive difference", 10, 3, 7},
		{"negative difference", 3, 10, -7},
		{"mixed numbers", -5, -3, -2},
		{"zeros", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Subtract(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("Subtract(%d, %d) = %d; expected %d", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestGreet(t *testing.T) {
	tests := []struct {
		name           string
		user           string
		customGreeting string
		expected       string
	}{
		{"default greeting", "Eddy", "", "Hello, Eddy!"},
		{"custom greeting", "Eddy", "Welcome", "Welcome, Eddy!"},
		{"empty custom greeting", "World", "", "Hello, World!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Greet(tt.user, tt.customGreeting)
			if result != tt.expected {
				t.Errorf("Greet(%q, %q) = %q; expected %q", tt.user, tt.customGreeting, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// SERVER HANDLER INTEGRATION TESTS
// =============================================================================

func TestHandleAdd(t *testing.T) {
	ctx := t.Context()
	args := AddArguments{A: 12, B: 30}

	result, extra, err := handleAdd(ctx, nil, args)
	if err != nil {
		t.Fatalf("handleAdd returned unexpected error: %v", err)
	}
	if extra != nil {
		t.Errorf("handleAdd returned unexpected extra: %v", extra)
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	expectedText := "Result of addition: 12 + 30 = 42"
	if textContent.Text != expectedText {
		t.Errorf("expected text %q, got %q", expectedText, textContent.Text)
	}
}

func TestHandleSubtract(t *testing.T) {
	ctx := t.Context()
	args := SubtractArguments{A: 50, B: 15}

	result, extra, err := handleSubtract(ctx, nil, args)
	if err != nil {
		t.Fatalf("handleSubtract returned unexpected error: %v", err)
	}
	if extra != nil {
		t.Errorf("handleSubtract returned unexpected extra: %v", extra)
	}

	if len(result.Content) != 1 {
		t.Fatalf("expected 1 content block, got %d", len(result.Content))
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", result.Content[0])
	}

	expectedText := "Result of subtraction: 50 - 15 = 35"
	if textContent.Text != expectedText {
		t.Errorf("expected text %q, got %q", expectedText, textContent.Text)
	}
}

func TestHandleGreet(t *testing.T) {
	t.Run("valid greeting request", func(t *testing.T) {
		ctx := t.Context()
		args := GreetArguments{Name: "Alice", Greeting: "Bonjour"}
		result, extra, err := handleGreet(ctx, nil, args)
		if err != nil {
			t.Fatalf("handleGreet failed: %v", err)
		}
		if extra != nil {
			t.Errorf("unexpected extra: %v", extra)
		}

		if len(result.Content) != 1 {
			t.Fatalf("expected 1 content block, got %d", len(result.Content))
		}

		textContent, ok := result.Content[0].(*mcp.TextContent)
		if !ok {
			t.Fatalf("expected TextContent, got %T", result.Content[0])
		}

		expectedText := "Bonjour, Alice!"
		if textContent.Text != expectedText {
			t.Errorf("expected text %q, got %q", expectedText, textContent.Text)
		}
	})

	t.Run("missing name error", func(t *testing.T) {
		ctx := t.Context()
		args := GreetArguments{Name: ""}
		_, _, err := handleGreet(ctx, nil, args)
		if err == nil {
			t.Error("expected error when name argument is empty, got nil")
		}
	})
}

// =============================================================================
// SERVER SETUP & REGISTRATION TESTS
// =============================================================================

func TestRegisterTools(t *testing.T) {
	t.Run("nil server validation", func(t *testing.T) {
		err := RegisterTools(nil)
		if err == nil {
			t.Error("expected error when registering tools on nil server, got nil")
		}
	})

	t.Run("successful tool registration", func(t *testing.T) {
		server := tinymcp.NewServer("test-server", "1.0.0")

		err := RegisterTools(server)
		if err != nil {
			t.Fatalf("unexpected error registering tools: %v", err)
		}
	})
}
