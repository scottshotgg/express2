package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestParseStatement_AllTypes(t *testing.T) {
	tests := []struct {
		name        string
		src         string
		wantType    string
		wantKind    string
		checkFunc   func(t *testing.T, n *builder.Node)
	}{
		{
			name:     "declaration",
			src:      "int i = 10",
			wantType: "decl",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (ident) is nil")
				}
				if n.Left.Value.(string) != "i" {
					t.Errorf("Left.Value = %q, want i", n.Left.Value)
				}
				if n.Right == nil {
					t.Fatal("Right (value) is nil")
				}
				if n.Value == nil {
					t.Fatal("Value (type node) is nil")
				}
			},
		},
		{
			name:     "simple_assignment",
			src:      "i = 10",
			wantType: "assignment",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Right == nil {
					t.Fatal("Right is nil")
				}
			},
		},
		{
			name:     "if_else",
			src:      `if true { int x = 1 } else { int y = 2 }`,
			wantType: "if",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (then-block) is nil")
				}
				if n.Right == nil {
					t.Fatal("Right (else-block) is nil")
				}
			},
		},
		{
			name:     "if_else_if",
			src:      `if true { int x = 1 } else if false { int y = 2 }`,
			wantType: "if",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Right == nil {
					t.Fatal("Right (else-if) is nil")
				}
				if n.Right.Type != "if" {
					t.Errorf("Right.Type = %q, want if", n.Right.Type)
				}
			},
		},
		{
			name:     "for_std",
			src:      "for int i = 0; i < 10; i++ { int k = 1 }",
			wantType: "forstd",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["start"] == nil {
					t.Error("Metadata[start] is nil")
				}
				if n.Metadata["end"] == nil {
					t.Error("Metadata[end] is nil")
				}
				if n.Metadata["step"] == nil {
					t.Error("Metadata[step] is nil")
				}
			},
		},
		{
			name:     "for_in",
			src:      "for i in [ 7, 8, 9 ] { j = 10 }",
			wantType: "forin",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["start"] == nil {
					t.Error("Metadata[start] is nil")
				}
				if n.Metadata["end"] == nil {
					t.Error("Metadata[end] is nil")
				}
			},
		},
		{
			name:     "for_of",
			src:      "for i of [ 7, 8, 9 ] { int i = 10 }",
			wantType: "forof",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["start"] == nil {
					t.Error("Metadata[start] is nil")
				}
				if n.Metadata["end"] == nil {
					t.Error("Metadata[end] is nil")
				}
			},
		},
		{
			name:     "for_over_single",
			src:      "for i over [ 1, 2, 3 ] { int k = 10 }",
			wantType: "forover",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["start"] == nil {
					t.Error("Metadata[start] is nil")
				}
				if n.Metadata["end"] == nil {
					t.Error("Metadata[end] is nil")
				}
				if n.Metadata["start2"] != nil {
					t.Error("Metadata[start2] should be nil for single-var form")
				}
			},
		},
		{
			name:     "for_over_two_vars",
			src:      "for i, j over [ 1, 2, 3 ] { int k = 10 }",
			wantType: "forover",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["start"] == nil {
					t.Error("Metadata[start] is nil")
				}
				if n.Metadata["start2"] == nil {
					t.Error("Metadata[start2] should be non-nil for two-var form")
				}
				if n.Metadata["end"] == nil {
					t.Error("Metadata[end] is nil")
				}
			},
		},
		{
			name:     "function_def",
			src:      "func something(int i, string s) int { return 10 }",
			wantType: "function",
			wantKind: "something",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Metadata["args"] == nil {
					t.Error("Metadata[args] is nil")
				}
				if n.Metadata["returns"] == nil {
					t.Error("Metadata[returns] is nil")
				}
			},
		},
		{
			name:     "function_no_return",
			src:      "func main() { int i = 1 }",
			wantType: "function",
			wantKind: "main",
		},
		{
			name:     "return_value",
			src:      "return 10",
			wantType: "return",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (return value) is nil")
				}
				if n.Left.Type != "literal" {
					t.Errorf("Left.Type = %q, want literal", n.Left.Type)
				}
			},
		},
		{
			name:     "block",
			src:      "{ int i = 10 int j = 20 }",
			wantType: "block",
			checkFunc: func(t *testing.T, n *builder.Node) {
				stmts, _ := n.Value.([]*builder.Node)
				if len(stmts) != 2 {
					t.Errorf("got %d stmts in block, want 2", len(stmts))
				}
			},
		},
		{
			name:     "empty_block",
			src:      "{ }",
			wantType: "block",
			checkFunc: func(t *testing.T, n *builder.Node) {
				stmts, _ := n.Value.([]*builder.Node)
				if len(stmts) != 0 {
					t.Errorf("got %d stmts in empty block, want 0", len(stmts))
				}
			},
		},
		{
			name:     "import_c",
			src:      "import c",
			wantType: "import",
			wantKind: "c",
		},
		{
			name:     "typedef",
			src:      "type myInt = int",
			wantType: "typedef",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "myInt" {
					t.Errorf("Left.Value = %q, want myInt", n.Left.Value)
				}
			},
		},
		{
			name:     "struct",
			src:      `struct something = { int i = 10 }`,
			wantType: "struct",
		},
		{
			name:     "object",
			src:      "object o = { int a = 6 }",
			wantType: "object",
		},
		{
			name:     "let",
			src:      "let something = 99",
			wantType: "let",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Value.(string) != "something" {
					t.Errorf("Left.Value = %q, want something", n.Left.Value)
				}
				if n.Right == nil {
					t.Fatal("Right is nil")
				}
			},
		},
		{
			name:     "deref_assign",
			src:      "*something = 10",
			wantType: "assignment",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Type != "deref" {
					t.Errorf("Left.Type = %q, want deref", n.Left.Type)
				}
			},
		},
		{
			name:     "index_assign",
			src:      `something[7] = "hey its me"`,
			wantType: "assignment",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Type != "index" {
					t.Errorf("Left.Type = %q, want index", n.Left.Type)
				}
			},
		},
		{
			name:     "selection_assign",
			src:      "some.thing = 10",
			wantType: "assignment",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left is nil")
				}
				if n.Left.Type != "selection" {
					t.Errorf("Left.Type = %q, want selection", n.Left.Type)
				}
			},
		},
		{
			name:     "launch",
			src:      "launch something()",
			wantType: "launch",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (call) is nil")
				}
			},
		},
		{
			name:     "defer",
			src:      "defer something()",
			wantType: "defer",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (call) is nil")
				}
			},
		},
		{
			name:     "enum_anon",
			src:      "enum { Red Green Blue }",
			wantType: "enum",
		},
		{
			name:     "binop_assign",
			src:      "i = 9 + 8 * 7",
			wantType: "assignment",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Right == nil {
					t.Fatal("Right is nil")
				}
				if n.Right.Type != "binop" {
					t.Errorf("Right.Type = %q, want binop", n.Right.Type)
				}
			},
		},
		{
			name:     "array_decl",
			src:      "int[] i = [ 8, 9, 0 ]",
			wantType: "decl",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Right == nil {
					t.Fatal("Right (array) is nil")
				}
				if n.Right.Type != "array" {
					t.Errorf("Right.Type = %q, want array", n.Right.Type)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			n := parseStatement(t, tt.src)
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
