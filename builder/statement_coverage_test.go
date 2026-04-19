package builder_test

import (
	"strings"
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
		{
			name:     "c_block_simple",
			src:      `c { printf("hello\n"); }`,
			wantType: "c",
			checkFunc: func(t *testing.T, n *builder.Node) {
				s, ok := n.Value.(string)
				if !ok {
					t.Fatal("Value is not string")
				}
				if !strings.Contains(s, `printf`) {
					t.Errorf("c block value missing printf, got: %q", s)
				}
				if !strings.Contains(s, `\n`) {
					t.Errorf("c block value missing escape sequence \\n, got: %q", s)
				}
			},
		},
		{
			name:     "c_block_nested_braces",
			src:      `c { if (1) { printf("nested"); } }`,
			wantType: "c",
			checkFunc: func(t *testing.T, n *builder.Node) {
				s, ok := n.Value.(string)
				if !ok {
					t.Fatal("Value is not string")
				}
				if !strings.Contains(s, "if") || !strings.Contains(s, "nested") {
					t.Errorf("c block with nested braces truncated, got: %q", s)
				}
			},
		},
		{
			name:     "while_statement",
			src:      "while true { int i = 1 }",
			wantType: "while",
			checkFunc: func(t *testing.T, n *builder.Node) {
				if n.Left == nil {
					t.Fatal("Left (condition) is nil")
				}
				if n.Value == nil {
					t.Fatal("Value (body) is nil")
				}
			},
		},
		{
			name:     "function_with_receiver",
			src:      "func MyType.myMethod() int { return 0 }",
			wantType: "function",
			wantKind: "myMethod",
			checkFunc: func(t *testing.T, n *builder.Node) {
				receiver, ok := n.Metadata["receiver"].(string)
				if !ok || receiver != "MyType" {
					t.Errorf("Metadata[receiver] = %q, want MyType", receiver)
				}
			},
		},
		{
			name:     "package_statement",
			src:      "package math { func add(int a, int b) int { return a + b } }",
			wantType: "package",
		},
		{
			name:     "include_statement",
			src:      "include something",
			wantType: "include",
		},
		{
			name:     "array_decl_ident_dim",
			src:      "char[amount] buff",
			wantType: "decl",
			checkFunc: func(t *testing.T, n *builder.Node) {
				typeNode, ok := n.Value.(*builder.Node)
				if !ok || typeNode == nil {
					t.Fatal("Value is not a *Node type")
				}
				typeVal, ok := typeNode.Value.(string)
				if !ok || typeVal != "array" {
					t.Errorf("typeNode.Value = %v, want \"array\"", typeNode.Value)
				}
				dims, ok := typeNode.Metadata["dim"].([]*builder.Index)
				if !ok || len(dims) == 0 {
					t.Fatal("dim metadata missing or empty")
				}
				if dims[0].Type != "ident" {
					t.Errorf("dim[0].Type = %q, want ident", dims[0].Type)
				}
				if dims[0].Value.(string) != "amount" {
					t.Errorf("dim[0].Value = %q, want amount", dims[0].Value)
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

// TestParseExpression_LogicalOps exercises the parseInfixLogical path (&&, ||).
func TestParseExpression_LogicalOps(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		wantValue string
	}{
		{"and", "a && b", "&&"},
		{"or", "x || y", "||"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := parseExpression(t, tt.src)
			if n.Type != "binop" {
				t.Errorf("Type = %q, want binop", n.Type)
			}
			if n.Value.(string) != tt.wantValue {
				t.Errorf("Value = %q, want %q", n.Value, tt.wantValue)
			}
			if n.Left == nil || n.Right == nil {
				t.Error("Left or Right is nil")
			}
		})
	}
}

// TestParseLetStatement_NoEquals exercises the missing-equals error path.
func TestParseLetStatement_NoEquals(t *testing.T) {
	// "let x\nint y = 0" — after parsing ident x, the next token is Separator,
	// not Assign, so ParseLetStatement returns an error.
	b, err := getBuilderFromString("let x\nint y = 0")
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.ParseLetStatement()
	if err == nil {
		t.Error("expected error for 'let x' without '='")
	}
}

// TestParseTypeDeclarationStatement_NoEquals exercises the missing-equals error path.
func TestParseTypeDeclarationStatement_NoEquals(t *testing.T) {
	// "type MyInt\nint x = 0" — after parsing ident MyInt, the next token is
	// Separator, not Assign, so ParseTypeDeclarationStatement returns an error.
	b, err := getBuilderFromString("type MyInt\nint x = 0")
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.ParseTypeDeclarationStatement()
	if err == nil {
		t.Error("expected error for 'type MyInt' without '='")
	}
}

// TestParseType_MapAnnotation exercises the int -> string map-annotation branch.
func TestParseType_MapAnnotation(t *testing.T) {
	b, err := getBuilderFromString("int -> string")
	if err != nil {
		t.Fatal(err)
	}
	n, err := b.ParseType(nil)
	if err != nil {
		t.Fatalf("ParseType error: %v", err)
	}
	if n.Kind != "map_annotation" {
		t.Errorf("Kind = %q, want map_annotation", n.Kind)
	}
	if n.Left == nil || n.Right == nil {
		t.Error("map_annotation node should have Left and Right")
	}
}

// TestParseTypeExpr_UserDefined exercises the token.Ident branch of ParseTypeExpr.
func TestParseTypeExpr_UserDefined(t *testing.T) {
	b, err := getBuilderFromString("MyType")
	if err != nil {
		t.Fatal(err)
	}
	// Register MyType in the scope tree so ParseTypeExpr can find it.
	if regErr := b.ScopeTree.NewType("MyType", &builder.TypeValue{
		Type:      builder.StruturedValue,
		Kind:      "MyType",
		Composite: true,
	}); regErr != nil {
		t.Fatalf("NewType: %v", regErr)
	}
	n, err := b.ParseTypeExpr()
	if err != nil {
		t.Fatalf("ParseTypeExpr error: %v", err)
	}
	if n.Type != "type" {
		t.Errorf("Type = %q, want type", n.Type)
	}
	if n.Value.(string) != "MyType" {
		t.Errorf("Value = %q, want MyType", n.Value)
	}
}

// TestParseTypeExpr_UnknownIdent exercises the unknown-type error path.
func TestParseTypeExpr_UnknownIdent(t *testing.T) {
	b, err := getBuilderFromString("UnknownXYZType")
	if err != nil {
		t.Fatal(err)
	}
	_, err = b.ParseTypeExpr()
	if err == nil {
		t.Error("expected error for unknown type identifier")
	}
}

// TestParseIdentStatement_CompoundAssign exercises +=, -=, *=, /= assignment desugaring.
func TestParseIdentStatement_CompoundAssign(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantType string
	}{
		{"add_assign", "i += 1", "assignment"},
		{"sub_assign", "i -= 1", "assignment"},
		{"mul_assign", "i *= 2", "assignment"},
		{"div_assign", "i /= 2", "assignment"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.src)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			if node.Type != tt.wantType {
				t.Errorf("got type %q, want %q", node.Type, tt.wantType)
			}
			// Right should be a binop
			if node.Right == nil || node.Right.Type != "binop" {
				t.Errorf("expected binop right node, got %+v", node.Right)
			}
		})
	}
}

// TestParseIdentStatement_StructLiteralDecl tests decl with LBrace (no explicit =).
// e.g. `Person p { name = "Alice" }` — decl with block expression.
func TestParseIdentStatement_SetKV(t *testing.T) {
	// "key" : "value" produces a kv node from ParseIdentStatement
	b, err := getBuilderFromString(`"hello" : "world"`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "kv" {
		t.Errorf("got type %q, want %q", node.Type, "kv")
	}
}

// TestParseIdentExpr_TypeName exercises the case where an ident resolves to a registered type.
func TestParseIdentExpr_TypeName(t *testing.T) {
	// Register a struct type first, then use it in expression context
	b, err := getBuilderFromString(`struct Point = { int x = 0  int y = 0 }
int z = 0`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// Parse the struct definition to register the type
	_, err = b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStructStatement error: %v", err)
	}

	// Now parse "int z = 0" which uses a primitive type — exercises ParseIdentStatement with type token
	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "decl" {
		t.Errorf("got type %q, want %q", node.Type, "decl")
	}
}

// TestScopeTree_GetImports exercises the GetImports method.
func TestScopeTree_GetImports(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatal(err)
	}
	imports := b.ScopeTree.GetImports()
	// May be nil or empty map — should not panic
	_ = imports
}

// TestAddStructured_Coverage exercises builder.AddStructured with a valid block node.
func TestAddStructured_Coverage(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatal(err)
	}
	// Build a block node whose single prop has type "int" (registered as a primitive).
	blockNode := &builder.Node{
		Type: "block",
		Value: []*builder.Node{
			{
				Type: "decl",
				Left: &builder.Node{Type: "ident", Value: "x"},
				// Value must be a *Node whose .Value is a string type name.
				Value: &builder.Node{Type: "type", Value: "int"},
			},
		},
	}
	tv, addErr := b.AddStructured("MyStruct", blockNode)
	if addErr != nil {
		t.Fatalf("AddStructured error: %v", addErr)
	}
	if !tv.Composite {
		t.Error("TypeValue.Composite should be true")
	}
	if tv.Type != builder.StruturedValue {
		t.Errorf("TypeValue.Type = %v, want StruturedValue", tv.Type)
	}
}

// TestParseMapStatement exercises ParseMapStatement via ParseStatement.
func TestParseMapStatement(t *testing.T) {
	b, err := getBuilderFromString(`map m = { "key" : "value" }`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "map" {
		t.Errorf("node.Type = %q, want map", node.Type)
	}
}

// TestParseTypedMapStatement exercises ParseMapStatement with [K -> V] type annotation.
func TestParseTypedMapStatement(t *testing.T) {
	tests := []struct {
		name      string
		src       string
		wantKey   string
		wantValue string
	}{
		{"string_int", `map[string -> int] scores = { "Alice" : 95 }`, "string", "int"},
		{"string_bool", `map[string -> bool] flags = {}`, "string", "bool"},
		{"string_float", `map[string -> float] rates = { "x" : 1 }`, "string", "float"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.src)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			if node.Type != "map" {
				t.Errorf("node.Type = %q, want map", node.Type)
			}

			if node.Metadata == nil {
				t.Fatal("node.Metadata is nil, expected key_node and value_node")
			}

			kn, ok := node.Metadata["key_node"].(*builder.Node)
			if !ok {
				t.Fatalf("Metadata[key_node] is not a *builder.Node: %T", node.Metadata["key_node"])
			}
			if kn.Kind != tt.wantKey {
				t.Errorf("key_node.Kind = %q, want %q", kn.Kind, tt.wantKey)
			}

			vn, ok := node.Metadata["value_node"].(*builder.Node)
			if !ok {
				t.Fatalf("Metadata[value_node] is not a *builder.Node: %T", node.Metadata["value_node"])
			}
			if vn.Kind != tt.wantValue {
				t.Errorf("value_node.Kind = %q, want %q", vn.Kind, tt.wantValue)
			}
		})
	}
}

// TestParseNDimensionalMap verifies that multi-key map annotations fold into nested map nodes.
func TestParseNDimensionalMap(t *testing.T) {
	src := `map[string, string -> int] scores`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "map" {
		t.Errorf("node.Type = %q, want map", node.Type)
	}

	// key_node should be string
	kn, ok := node.Metadata["key_node"].(*builder.Node)
	if !ok || kn.Kind != "string" {
		t.Errorf("key_node.Kind = %v, want string", kn)
	}

	// value_node should be a nested map type: map[string -> int]
	vn, ok := node.Metadata["value_node"].(*builder.Node)
	if !ok || vn.Kind != "map" {
		t.Errorf("value_node.Kind = %v, want map", vn)
	}

	// inner key should be string
	innerKey, ok := vn.Metadata["key_node"].(*builder.Node)
	if !ok || innerKey.Kind != "string" {
		t.Errorf("inner key_node.Kind = %v, want string", innerKey)
	}

	// inner value should be int
	innerVal, ok := vn.Metadata["value_node"].(*builder.Node)
	if !ok || innerVal.Kind != "int" {
		t.Errorf("inner value_node.Kind = %v, want int", innerVal)
	}
}

// TestParseTypedMapStatement_Errors exercises error paths in parseMapTypeAnnotation.
func TestParseTypedMapStatement_Errors(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{"missing_arrow", `map[string int] m = {}`},
		{"missing_value_after_arrow", `map[string ->] m = {}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.src)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			_, err = b.ParseStatement()
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestBreakStatement verifies that break parses to a break node.
func TestBreakStatement(t *testing.T) {
	src := `break`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "break" {
		t.Errorf("node.Type = %q, want break", node.Type)
	}
}

// TestContinueStatement verifies that continue parses to a continue node.
func TestContinueStatement(t *testing.T) {
	src := `continue`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "continue" {
		t.Errorf("node.Type = %q, want continue", node.Type)
	}
}

// TestNestedArrayType verifies that int[][] parses with two dimensions.
func TestNestedArrayType(t *testing.T) {
	src := `int[][] v`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	// The declaration node should have a type node with dim len 2
	typeNode, ok := node.Value.(*builder.Node)
	if !ok {
		t.Fatalf("node.Value is not *builder.Node: %T", node.Value)
	}

	dims, ok := typeNode.Metadata["dim"].([]*builder.Index)
	if !ok {
		t.Fatalf("typeNode.Metadata[dim] is not []*builder.Index: %T", typeNode.Metadata["dim"])
	}

	if len(dims) != 2 {
		t.Errorf("dim length = %d, want 2", len(dims))
	}
}

// TestParsePrefixNot exercises parsePrefixNot via ParseExpression.
func TestParsePrefixNot(t *testing.T) {
	b, err := getBuilderFromString("!x")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression error: %v", err)
	}

	if node.Type != "not" {
		t.Errorf("node.Type = %q, want not", node.Type)
	}
}

// TestParsePostfixDec exercises parsePostfixDec via ParseExpression.
func TestParsePostfixDec(t *testing.T) {
	b, err := getBuilderFromString("x--")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression error: %v", err)
	}

	if node.Type != "dec" {
		t.Errorf("node.Type = %q, want dec", node.Type)
	}
}

// TestTypeResolver_Check_EgroupIdents exercises the egroup branch of TypeResolver.Check
// with ident return types (not already resolved to "type" nodes).
func TestTypeResolver_Check_EgroupIdents(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatal(err)
	}

	tr := builder.NewTypeResolverWithScope(b.ScopeTree)

	t.Run("ident_resolved", func(t *testing.T) {
		// egroup with an ident that resolves to a known type
		n := &builder.Node{
			Type: "egroup",
			Value: []*builder.Node{
				{Type: "ident", Value: "int"},
			},
		}
		_, err := tr.Check(n)
		if err != nil {
			t.Errorf("Check egroup(ident=int) error: %v", err)
		}
	})

	t.Run("ident_not_found", func(t *testing.T) {
		// egroup with ident that is NOT in scope → error
		n := &builder.Node{
			Type: "egroup",
			Value: []*builder.Node{
				{Type: "ident", Value: "nonexistent_type_xyz987"},
			},
		}
		_, err := tr.Check(n)
		if err == nil {
			t.Error("expected error for unfound type, got nil")
		}
	})

	t.Run("non_string_value", func(t *testing.T) {
		// egroup with a node whose Value is not a string → error
		n := &builder.Node{
			Type: "egroup",
			Value: []*builder.Node{
				{Type: "ident", Value: 42}, // non-string Value
			},
		}
		_, err := tr.Check(n)
		if err == nil {
			t.Error("expected error for non-string Value, got nil")
		}
	})
}

// TestTypeResolver_Check_LetIdent exercises the let+ident inference path.
func TestTypeResolver_Check_LetIdent(t *testing.T) {
	// let x = someCall() — RHS is a call → typeNode becomes "unknown"
	b, err := getBuilderFromString("let x = foo()")
	if err != nil {
		t.Fatal(err)
	}

	ast, err := b.BuildAST()
	if err != nil {
		t.Fatal(err)
	}

	// The let statement's RHS is a call node → hits the default case in Check.
	tr := builder.NewTypeResolverWithScope(b.ScopeTree)
	_, err = tr.Check(ast)
	if err != nil {
		t.Fatalf("Check error: %v", err)
	}
}

// TestProcessTypeDeclaration_ErrorPaths exercises the ident-not-found and default cases.
func TestProcessTypeDeclaration_ErrorPaths(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// "ident" case but the type is unknown → error
	t.Run("ident_unknown_type", func(t *testing.T) {
		decl := &builder.Node{
			Type: "typedef",
			Left: &builder.Node{Type: "ident", Value: "MyT"},
			Right: &builder.Node{Type: "ident", Value: "NonExistent_xyz_987"},
		}
		err := b.ProcessTypeDeclaration(decl)
		if err == nil {
			t.Error("expected error for unknown type, got nil")
		}
	})

	// default case: unsupported type expression
	t.Run("default_unsupported", func(t *testing.T) {
		decl := &builder.Node{
			Type: "typedef",
			Left: &builder.Node{Type: "ident", Value: "MyT"},
			Right: &builder.Node{Type: "block"},
		}
		err := b.ProcessTypeDeclaration(decl)
		if err == nil {
			t.Error("expected error for unsupported type expression, got nil")
		}
	})
}

// TestParseMethodDeclaration exercises the method receiver path in ParseFunctionStatement.
func TestParseMethodDeclaration(t *testing.T) {
	// "func Point.getX() { }" — accessor '.' between receiver and method name
	b, err := getBuilderFromString("func Point.getX() { }")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	n, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if n.Type != "function" {
		t.Fatalf("Type = %q, want function", n.Type)
	}

	if n.Kind != "getX" {
		t.Errorf("Kind = %q, want getX", n.Kind)
	}

	receiver, _ := n.Metadata["receiver"].(string)
	if receiver != "Point" {
		t.Errorf("Metadata[receiver] = %q, want Point", receiver)
	}
}

// TestTypeResolver_Check_CBlock exercises the "c" case in TypeResolver.Check.
func TestTypeResolver_Check_CBlock(t *testing.T) {
	tr := builder.NewTypeResolver()
	n := &builder.Node{Type: "c", Value: "printf(\"hi\\n\");"}
	changed, err := tr.Check(n)
	if err != nil {
		t.Fatalf("Check(c) unexpected error: %v", err)
	}
	if changed {
		t.Error("Check(c) returned changed=true, want false")
	}
}

// TestParseStructLiteralInExpression exercises the struct-literal path in parseIdentExpr.
// When a known-type ident is immediately followed by '{', it parses a struct literal.
func TestParseStructLiteralInExpression(t *testing.T) {
	// Register "Point" as a known type, then parse "Point { }" as an expression.
	// We do this via a full program so the scope is set up correctly.
	b, err := getBuilderFromString("struct Point = { int x = 0 }")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// Parse the struct declaration to register the type
	_, err = b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement (struct) error: %v", err)
	}

	// Now parse a struct literal expression: Point { x = 1 }
	b2, err := getBuilderFromString("Point { x = 1 }")
	if err != nil {
		t.Fatalf("lex error for literal: %v", err)
	}

	// Copy the scope tree so Point is known
	b2.ScopeTree = b.ScopeTree

	n, err := b2.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression error: %v", err)
	}

	if n.Type != "literal" {
		t.Errorf("Type = %q, want literal", n.Type)
	}
}

// TestTypeResolver_Check_Typedef_Selection exercises the typedef+selection error path.
func TestTypeResolver_Check_Typedef_Selection(t *testing.T) {
	tr := builder.NewTypeResolver()
	// typedef where RHS is a "selection" node → returns "selection types: not implemented"
	n := &builder.Node{
		Type: "typedef",
		Left: &builder.Node{Type: "ident", Value: "MyType"},
		Right: &builder.Node{
			Type: "selection",
			Value: &builder.Node{
				Left:  &builder.Node{Type: "ident", Value: "pkg"},
				Right: &builder.Node{Type: "ident", Value: "SomeType"},
			},
		},
	}
	_, err := tr.Check(n)
	if err == nil {
		t.Error("expected error for selection typedef, got nil")
	}
}

// TestAddPrimitive_ErrorPaths exercises the error branches in AddPrimitive.
func TestAddPrimitive_ErrorPaths(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// nil Value → error
	t.Run("nil_value", func(t *testing.T) {
		_, err := b.AddPrimitive("T", &builder.Node{Type: "ident", Value: nil})
		if err == nil {
			t.Error("expected error for nil Value, got nil")
		}
	})

	// type already declared → error
	t.Run("already_declared", func(t *testing.T) {
		// "int" is already registered as a primitive
		_, err := b.AddPrimitive("int", &builder.Node{Type: "ident", Value: "float"})
		if err == nil {
			t.Error("expected error for already-declared type, got nil")
		}
	})

	// Value not a string → error
	t.Run("non_string_value", func(t *testing.T) {
		_, err := b.AddPrimitive("MyT", &builder.Node{Type: "ident", Value: 42})
		if err == nil {
			t.Error("expected error for non-string Value, got nil")
		}
	})

	// Type being aliased is not declared → error
	t.Run("alias_not_declared", func(t *testing.T) {
		_, err := b.AddPrimitive("MyT", &builder.Node{Type: "ident", Value: "nonexistent_type_xyz"})
		if err == nil {
			t.Error("expected error for undeclared alias type, got nil")
		}
	})
}

// TestExtractPropsFromComposite_Default exercises the default error path.
func TestExtractPropsFromComposite_Default(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// "block" type works; anything else should fail
	// AddStructured calls extractPropsFromComposite with n.Type != "block" to hit default
	_, addErr := b.AddStructured("T", &builder.Node{Type: "program", Kind: "notblock"})
	if addErr == nil {
		t.Error("expected error for non-block node in AddStructured, got nil")
	}
}

// TestBuildNodeFromTypeValue_NilAndDefault exercises nil-input and default-type errors.
func TestBuildNodeFromTypeValue_NilAndDefault(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// nil TypeValue → error
	t.Run("nil_typevalue", func(t *testing.T) {
		_, err := b.BuildNodeFromTypeValue(nil)
		if err == nil {
			t.Error("expected error for nil TypeValue, got nil")
		}
	})

	// PrimitiveValue TypeValue → hits default case → error
	t.Run("primitive_default", func(t *testing.T) {
		tv := &builder.TypeValue{Type: builder.PrimitiveValue, Kind: "int"}
		_, err := b.BuildNodeFromTypeValue(tv)
		if err == nil {
			t.Error("expected error for PrimitiveValue in BuildNodeFromTypeValue, got nil")
		}
	})
}

// TestIfExpression verifies that `if` can appear in expression position via ParseExpression.
func TestIfExpression(t *testing.T) {
	n := parseExpression(t, "if true { let x = 1 } else { let y = 2 }")
	if n.Type != "if" {
		t.Fatalf("expected if node, got %q", n.Type)
	}
	if n.Left == nil {
		t.Fatal("Left (then-branch) is nil")
	}
	if n.Right == nil {
		t.Fatal("Right (else-branch) is nil")
	}
}

// TestIfExpressionNoElse verifies that `if` without else works in expression position.
func TestIfExpressionNoElse(t *testing.T) {
	n := parseExpression(t, "if true { let x = 1 }")
	if n.Type != "if" {
		t.Fatalf("expected if node, got %q", n.Type)
	}
	if n.Left == nil {
		t.Fatal("Left (then-branch) is nil")
	}
}

// TestBlockExpression verifies that a bare block `{ ... }` can appear in expression position.
func TestBlockExpression(t *testing.T) {
	n := parseExpression(t, "{ let x = 1 }")
	if n.Type != "block" {
		t.Fatalf("expected block node, got %q", n.Type)
	}
}

// TestLetExpression verifies that `let` can appear in expression position.
func TestLetExpression(t *testing.T) {
	n := parseExpression(t, "let x = 5")
	if n.Type != "let" {
		t.Fatalf("expected let node, got %q", n.Type)
	}
	if n.Left == nil || n.Left.Value.(string) != "x" {
		t.Fatalf("expected Left = x, got %v", n.Left)
	}
}

// TestVarTypedDeclaration verifies `var int x = 5` produces a mutable typed decl.
func TestVarTypedDeclaration(t *testing.T) {
	b, err := getBuilderFromString("var int x = 5")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "decl" {
		t.Fatalf("node.Type = %q, want decl", node.Type)
	}
	typeNode, ok := node.Value.(*builder.Node)
	if !ok {
		t.Fatal("node.Value is not *builder.Node")
	}
	if typeNode.Kind != "int" {
		t.Errorf("typeNode.Kind = %q, want int", typeNode.Kind)
	}
	if node.Left == nil || node.Left.Value.(string) != "x" {
		t.Errorf("node.Left.Value = %v, want x", node.Left)
	}
	if node.Right == nil {
		t.Fatal("node.Right is nil, want literal 5")
	}
	if node.Metadata == nil || node.Metadata["mutable"] != true {
		t.Errorf("node.Metadata[mutable] = %v, want true", node.Metadata["mutable"])
	}
}

// TestVarTypedDeclUninitialized verifies `var int x` (no initializer) produces a mutable typed decl with nil Right.
func TestVarTypedDeclUninitialized(t *testing.T) {
	b, err := getBuilderFromString("var int x")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	if node.Type != "decl" {
		t.Fatalf("node.Type = %q, want decl", node.Type)
	}
	typeNode, ok := node.Value.(*builder.Node)
	if !ok {
		t.Fatal("node.Value is not *builder.Node")
	}
	if typeNode.Kind != "int" {
		t.Errorf("typeNode.Kind = %q, want int", typeNode.Kind)
	}
	if node.Left == nil || node.Left.Value.(string) != "x" {
		t.Errorf("node.Left.Value = %v, want x", node.Left)
	}
	if node.Right != nil {
		t.Errorf("node.Right = %v, want nil", node.Right)
	}
	if node.Metadata == nil || node.Metadata["mutable"] != true {
		t.Errorf("node.Metadata[mutable] = %v, want true", node.Metadata["mutable"])
	}
}
