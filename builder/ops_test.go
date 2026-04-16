package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestParseGroupOfExpressions_Empty(t *testing.T) {
	n := parseExpression(t, "()")
	if n.Type != "egroup" {
		t.Fatalf("Type = %q, want egroup", n.Type)
	}
	items, _ := n.Value.([]*builder.Node)
	if len(items) != 0 {
		t.Errorf("got %d items, want 0", len(items))
	}
}

func TestParseGroupOfExpressions_Single(t *testing.T) {
	n := parseExpression(t, "(1)")
	if n.Type != "egroup" {
		t.Fatalf("Type = %q, want egroup", n.Type)
	}
	items, ok := n.Value.([]*builder.Node)
	if !ok || len(items) != 1 {
		t.Errorf("got %d items, want 1", len(items))
	}
}

func TestParseGroupOfExpressions_Multiple(t *testing.T) {
	n := parseExpression(t, `(1, i, "s")`)
	if n.Type != "egroup" {
		t.Fatalf("Type = %q, want egroup", n.Type)
	}
	items, ok := n.Value.([]*builder.Node)
	if !ok || len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestParseArrayExpression_Empty(t *testing.T) {
	n := parseExpression(t, "[ ]")
	if n.Type != "array" {
		t.Fatalf("Type = %q, want array", n.Type)
	}
	items, _ := n.Value.([]*builder.Node)
	if len(items) != 0 {
		t.Errorf("got %d items, want 0", len(items))
	}
}

func TestParseArrayExpression_Items(t *testing.T) {
	n := parseExpression(t, "[ 1, 2, 3 ]")
	if n.Type != "array" {
		t.Fatalf("Type = %q, want array", n.Type)
	}
	items, ok := n.Value.([]*builder.Node)
	if !ok || len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}

func TestParseArrayExpression_Nested(t *testing.T) {
	n := parseExpression(t, "[ [1], [2, 3] ]")
	if n.Type != "array" {
		t.Fatalf("Type = %q, want array", n.Type)
	}
	items, ok := n.Value.([]*builder.Node)
	if !ok || len(items) != 2 {
		t.Errorf("got %d items, want 2", len(items))
	}
	for i, item := range items {
		if item.Type != "array" {
			t.Errorf("items[%d].Type = %q, want array", i, item.Type)
		}
	}
}

func TestParseCall_NoArgs(t *testing.T) {
	n := parseExpression(t, "foo()")
	if n.Type != "call" {
		t.Fatalf("Type = %q, want call", n.Type)
	}
	args, ok := n.Metadata["args"].(*builder.Node)
	if !ok || args == nil {
		t.Fatal("Metadata[args] missing or wrong type")
	}
	items, _ := args.Value.([]*builder.Node)
	if len(items) != 0 {
		t.Errorf("got %d args, want 0", len(items))
	}
}

func TestParseCall_WithArgs(t *testing.T) {
	n := parseExpression(t, "foo(1, 2)")
	if n.Type != "call" {
		t.Fatalf("Type = %q, want call", n.Type)
	}
	args, ok := n.Metadata["args"].(*builder.Node)
	if !ok || args == nil {
		t.Fatal("Metadata[args] missing or wrong type")
	}
	items, ok := args.Value.([]*builder.Node)
	if !ok || len(items) != 2 {
		t.Errorf("got %d args, want 2", len(items))
	}
}
