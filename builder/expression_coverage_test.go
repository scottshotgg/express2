package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestParseExpression_AllTypes(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		wantType  string
		wantKind  string
		checkFunc func(t *testing.T, n *builder.Node)
	}{
		{
			name:     "int_literal",
			src:      "42",
			wantType: "literal",
			wantKind: "int",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(int) != 42 {
					t.Errorf("Value = %v, want 42", n.Value)
				}
			},
		},
		{
			name:     "string_literal",
			src:      `"hello"`,
			wantType: "literal",
			wantKind: "string",
		},
		{
			name:     "bool_true",
			src:      "true",
			wantType: "literal",
			wantKind: "bool",
		},
		{
			name:     "ident",
			src:      "myVar",
			wantType: "ident",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != "myVar" {
					t.Errorf("Value = %q, want myVar", n.Value)
				}
			},
		},
		{
			name:     "deref",
			src:      "*something",
			wantType: "deref",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "something" {
					t.Errorf("Left.Value = %q, want something", n.Left.Value)
				}
			},
		},
		{
			name:     "ref",
			src:      "&something",
			wantType: "ref",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "something" {
					t.Errorf("Left.Value = %q, want something", n.Left.Value)
				}
			},
		},
		{
			name:     "increment",
			src:      "i++",
			wantType: "inc",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "i" {
					t.Errorf("Left.Value = %q, want i", n.Left.Value)
				}
			},
		},
		{
			name:     "addition",
			src:      "1 + 2",
			wantType: "binop",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != "+" {
					t.Errorf("Value = %q, want +", n.Value)
				}
				if n.Left == nil || n.Right == nil {
					t.Fatal("Left or Right is nil")
				}
			},
		},
		{
			name:     "multiplication",
			src:      "3 * 4",
			wantType: "binop",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != "*" {
					t.Errorf("Value = %q, want *", n.Value)
				}
			},
		},
		{
			name:     "less_than",
			src:      "a < 10",
			wantType: "comp",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != "<" {
					t.Errorf("Value = %q, want <", n.Value)
				}
			},
		},
		{
			name:     "equality",
			src:      "i == 0",
			wantType: "comp",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != "==" {
					t.Errorf("Value = %q, want ==", n.Value)
				}
			},
		},
		{
			name:     "gte",
			src:      "a >= 5",
			wantType: "comp",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value.(string) != ">=" {
					t.Errorf("Value = %q, want >=", n.Value)
				}
			},
		},
		{
			name:     "selection",
			src:      "a.b.c",
			wantType: "selection",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Right == nil {
					t.Fatal("Right is nil")
				}
				if n.Right.Value.(string) != "c" {
					t.Errorf("Right.Value = %q, want c", n.Right.Value)
				}
			},
		},
		{
			name:     "index",
			src:      "arr[0]",
			wantType: "index",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "arr" {
					t.Errorf("Left.Value = %q, want arr", n.Left.Value)
				}
			},
		},
		{
			name:     "call_no_args",
			src:      "foo()",
			wantType: "call",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Value == nil {
					t.Fatal("Value (callee) is nil")
				}
			},
		},
		{
			name:     "call_with_args",
			src:      "foo(1, 2)",
			wantType: "call",
			checkFunc: func(t *testing.T, n *builder.Node) {
				args, ok := n.Metadata["args"].(*builder.Node)
				if !ok || args == nil {
					t.Fatal("args metadata missing")
				}
				items, _ := args.Value.([]*builder.Node)
				if len(items) != 2 {
					t.Errorf("got %d args, want 2", len(items))
				}
			},
		},
		{
			name:     "array_literal",
			src:      "[ 1, 2, 3 ]",
			wantType: "array",
			checkFunc: func(t *testing.T, n *builder.Node) {
				items, _ := n.Value.([]*builder.Node)
				if len(items) != 3 {
					t.Errorf("got %d items, want 3", len(items))
				}
			},
		},
		{
			name:     "empty_array",
			src:      "[ ]",
			wantType: "array",
			checkFunc: func(t *testing.T, n *builder.Node) {
				items, _ := n.Value.([]*builder.Node)
				if len(items) != 0 {
					t.Errorf("got %d items, want 0", len(items))
				}
			},
		},
		{
			name:     "grouped_expr",
			src:      "(1 + 2)",
			wantType: "egroup",
		},
		{
			name:     "block_expr",
			src:      "{ int i = 7 }",
			wantType: "block",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := parseExpression(t, tt.src)
			if n.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", n.Type, tt.wantType)
			}
			if tt.wantKind != "" && n.Kind != tt.wantKind {
				t.Errorf("Kind = %q, want %q", n.Kind, tt.wantKind)
			}
			if tt.checkFunc != nil {
				tt.checkFunc(t, n)
			}
		})
	}
}
