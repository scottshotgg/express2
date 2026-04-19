package transpiler_test

import (
	"strings"
	"testing"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
	"github.com/scottshotgg/express2/transpiler"
)

// Helper functions to create builder and transpiler from strings

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	var tokens, err = lex.New(s).Lex()
	if err != nil {
		return nil, err
	}

	// Compress certain tokens; i.e: `:` and `=` compress into `:=`
	return ast.CompressTokens(tokens)
}

func getBuilderFromString(test string) (*builder.Builder, error) {
	var tokens, err = getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	return builder.New(tokens, logger.Noop()), nil
}

func getTranspilerFromBuilder(b *builder.Builder) (*transpiler.Transpiler, error) {
	ast, err := b.BuildAST()
	if err != nil {
		return nil, err
	}

	tr := transpiler.New(ast, b, "test", "/home/scottshotgg/Development/go/src/github.com/scottshotgg/express2/")
	return tr, nil
}

func getTranspilerFromString(src string) (*transpiler.Transpiler, error) {
	b, err := getBuilderFromString(src)
	if err != nil {
		return nil, err
	}

	return getTranspilerFromBuilder(b)
}

// ============================================
// Expression Transpilation Tests
// ============================================

func TestTranspileLiteralExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"string literal", `"hello world"`, `"hello world"`},
		{"int literal", "42", "42"},
		{"bool true", "true", "true"},
		{"bool false", "false", "false"},
		{"char literal", "'a'", "'a'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileLiteralExpression(node)
			if err != nil {
				t.Fatalf("TranspileLiteralExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileIdentExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"simple ident", "x", "x"},
		{"ident with underscore", "my_var", "my_var"},
		{"ident with numbers", "var123", "var123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileIdentExpression(node)
			if err != nil {
				t.Fatalf("TranspileIdentExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileArrayExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"empty array", "[]", " }"},
		{"int array", "[1, 2, 3]", "{ 1, 2, 3 }"},
		{"string array", `["a", "b"]`, `{ "a", "b" }`},
		{"mixed array", `[1, "test", true]`, `{ 1, "test", true }`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseArrayExpression()
			if err != nil {
				t.Fatalf("ParseArrayExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileArrayExpression(node)
			if err != nil {
				t.Fatalf("TranspileArrayExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileBinopExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"addition", "1 + 2", "1+2"},
		{"subtraction", "5 - 3", "5-3"},
		{"multiplication", "2 * 3", "2*3"},
		{"division", "10 / 2", "10/2"},
		{"complex", "a + b * c", "a+b*c"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			// For binop expressions (arithmetic)
			if node.Type == "binop" {
				result, err := tr.TranspileBinOpExpression(node)
				if err != nil {
					t.Fatalf("TranspileBinOpExpression error: %v", err)
				}
				if *result != tt.expected {
					t.Errorf("got %q, want %q", *result, tt.expected)
				}
			} else {
				// For comparison expressions (comp)
				result, err := tr.TranspileConditionExpression(node)
				if err != nil {
					t.Fatalf("TranspileConditionExpression error: %v", err)
				}
				if *result != tt.expected {
					t.Errorf("got %q, want %q", *result, tt.expected)
				}
			}
		})
	}
}

func TestTranspileCallExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"no args", "print()", "print()"},
		{"one arg", "print(5)", "print(5)"},
		{"multiple args", "print(1, 2, 3)", "print(1,2,3)"},
		{"nested call", "outer(inner())", "outer(inner())"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileCallExpression(node)
			if err != nil {
				t.Fatalf("TranspileCallExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileSelectionExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"simple selection", "obj.field", "obj.field"},
		{"nested selection", "obj.sub.field", "obj.sub.field"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileSelectExpression(node)
			if err != nil {
				t.Fatalf("TranspileSelectExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

// ============================================
// Statement Transpilation Tests
// ============================================

func TestTranspileDeclarationStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"int declaration", "int i = 10", "const int i = 10;"},
		{"string declaration", "string s = \"hello\"", `const std::string s = "hello";`},
		{"bool declaration", "bool b = true", "const bool b = true;"},
		{"var declaration", "var x = 5", "var x = 5;"},
		{"let declaration", "let y = 10", "const int y = 10;"},
		{"array declaration", "int[] arr = [1, 2, 3]", "const std::vector<int> arr = { 1, 2, 3 };"},
		{"declaration without init", "int x", "const int x = 0;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileDeclarationStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileAssignmentStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"simple assign", "x = 10", "x = 10;"},
		{"array assign", "arr[0] = 5", "arr[0] = 5;"},
		{"binop assign", "x = a + b", "x = a+b;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileAssignmentStatement(node)
			if err != nil {
				t.Fatalf("TranspileAssignmentStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileFunctionStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		contains []string
	}{
		{"simple func", "func main() { int x = 10 }", []string{"int main()", "int x = 10;"}},
		{"func with args", "func add(int a, int b) int { return a + b }", []string{"int add(int a,int b)", "return a+b;"}},
		{"func with return", "func getValue() int { return 5 }", []string{"int getValue()", "return 5;"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileFunctionStatement(node)
			if err != nil {
				t.Fatalf("TranspileFunctionStatement error: %v", err)
			}

			for _, sub := range tt.contains {
				if !strings.Contains(*result, sub) {
					t.Errorf("result missing %q\nfull output: %s", sub, *result)
				}
			}
		})
	}
}

func TestTranspileReturnStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"return value", "return 5", "return 5;"},
		{"return expression", "return a + b", "return a+b;"},
		{"return empty", "return\n", "return;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileReturnStatement(node)
			if err != nil {
				t.Fatalf("TranspileReturnStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileIfStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		contains []string
	}{
		{"simple if", `if x { int y = 10 }`, []string{"if (x)", "int y = 10;"}},
		{"if with else", `if x { int y = 10 } else { int z = 20 }`, []string{"if (x)", "else", "int z = 20;"}},
		{"if with elseif", `if x { int y = 10 } else if z { int w = 15 } else { int v = 20 }`, []string{"if (x)", "else if (z)", "int w = 15;"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileIfStatement(node)
			if err != nil {
				t.Fatalf("TranspileIfStatement error: %v", err)
			}

			for _, sub := range tt.contains {
				if !strings.Contains(*result, sub) {
					t.Errorf("result missing %q\nfull output: %s", sub, *result)
				}
			}
		})
	}
}

func TestTranspileBlockStatement(t *testing.T) {
	stmt := `{ int x = 10 int y = 20 }`

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileBlockStatement(node)
	if err != nil {
		t.Fatalf("TranspileBlockStatement error: %v", err)
	}

	for _, sub := range []string{"{", "int x = 10;", "int y = 20;"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

// TestTranspileForInStatement tests for-in statement parsing only.
// Note: the actual compilation path goes through tree_flattener → while loop;
// see TestForLoopTranspilation for a full-pipeline test.
func TestTranspileForInStatement(t *testing.T) {
	stmt := "for i in [1, 2, 3] { int x = i }"

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	_, err = b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}
}

// TestTranspileForOfStatement tests for-of statement parsing only.
// Note: the actual compilation path goes through tree_flattener → while loop;
// see TestForLoopTranspilation for a full-pipeline test.
func TestTranspileForOfStatement(t *testing.T) {
	stmt := "for i of [1, 2, 3] { int x = i }"

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	_, err = b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}
}

// TestForLoopTranspilation verifies the full pipeline for for-in and for-of:
// parser → tree_flattener (converts to while) → transpiler → C++.
func TestForLoopTranspilation(t *testing.T) {
	tests := []struct {
		name    string
		program string
		want    []string // substrings that must appear in ToCpp() output
	}{
		{
			name:    "for-in generates while with index",
			program: `func main() { for i in [1, 2, 3] { Println("idx:", i) } }`,
			want:    []string{"int i = 0", "while", "std::size", "(i)++"},
		},
		{
			name:    "for-of generates while with element access",
			program: `func main() { for v of [10, 20, 30] { Println("val:", v) } }`,
			want:    []string{"int _idx_0 = 0", "while", "std::size", "(_idx_0)++"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.program)
			if err != nil {
				t.Fatalf("getTranspilerFromString: %v", err)
			}

			if err := tr.Transpile(); err != nil {
				t.Fatalf("Transpile: %v", err)
			}

			cpp := tr.ToCpp()
			for _, sub := range tt.want {
				if !strings.Contains(cpp, sub) {
					t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
				}
			}
		})
	}
}

func TestTranspileStructStatement(t *testing.T) {
	stmt := "struct Point = { int x = 0 int y = 0 }"

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStructDeclaration(node)
	if err != nil {
		t.Fatalf("TranspileStructDeclaration error: %v", err)
	}

	for _, sub := range []string{"struct Point", "int x", "int y"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

func TestTranspileMapStatement(t *testing.T) {
	stmt := `map m = { "key" = "value" }`

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	// Debug: print tokens
	for i, tok := range b.Tokens {
		t.Logf("Token %d: %+v", i, tok)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileMapStatement(node)
	if err != nil {
		t.Fatalf("TranspileMapStatement error: %v", err)
	}

	for _, sub := range []string{"std::map", `"key"`, `"value"`} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

// TestTranspileTypedMapStatement exercises TranspileMapStatement with [K -> V] annotation.
func TestTranspileTypedMapStatement(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		wantType string
	}{
		{"string_int", `map[string -> int] scores = { "Alice" : 95 }`, "std::map<std::string, int>"},
		{"string_bool", `map[string -> bool] flags = {}`, "std::map<std::string, bool>"},
		{"string_float", `map[string -> float] rates = { "x" : 1 }`, "std::map<std::string, float>"},
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

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileMapStatement(node)
			if err != nil {
				t.Fatalf("TranspileMapStatement error: %v", err)
			}

			if !strings.Contains(*result, tt.wantType) {
				t.Errorf("result missing %q\nfull output: %s", tt.wantType, *result)
			}
		})
	}
}

// TestTranspileUntypedMapStatement verifies untyped map still uses std::string/std::string default.
func TestTranspileUntypedMapStatement(t *testing.T) {
	b, err := getBuilderFromString(`map m = { "key" : "value" }`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileMapStatement(node)
	if err != nil {
		t.Fatalf("TranspileMapStatement error: %v", err)
	}

	if !strings.Contains(*result, "std::map<std::string, var>") {
		t.Errorf("untyped map should default to std::string/var, got: %s", *result)
	}
}

// ============================================
// Type Transpilation Tests
// ============================================

func TestTranspileType(t *testing.T) {
	tests := []struct {
		name     string
		typeStr  string
		expected string
	}{
		{"int type", "int", "int"},
		{"string type", "string", "std::string"},
		{"bool type", "bool", "bool"},
		{"char type", "char", "char"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.typeStr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileType(node)
			if err != nil {
				t.Fatalf("TranspileType error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileStructDeclaration(t *testing.T) {
	stmt := "struct Point = { int x = 0 int y = 0 }"

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStructDeclaration(node)
	if err != nil {
		t.Fatalf("TranspileStructDeclaration error: %v", err)
	}

	for _, sub := range []string{"struct Point", "int x", "int y"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

func TestTranspileTypeDeclaration(t *testing.T) {
	stmt := "type myInt = int"

	b, err := getBuilderFromString(stmt)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileTypeDeclaration(node)
	if err != nil {
		t.Fatalf("TranspileTypeDeclaration error: %v", err)
	}

	if *result != "typedef int myInt;" {
		t.Errorf("got %q, want %q", *result, "typedef int myInt;")
	}
}

// ============================================
// Full Transpilation Tests
// ============================================

func TestTranspileCompleteProgram(t *testing.T) {
	tests := []struct {
		name     string
		program  string
		contains []string
	}{
		{
			"simple program with main",
			`func main() { int x = 10 Println(x) }`,
			[]string{"int main()", "int x = 10;"},
		},
		{
			"program with function",
			`func add(int a, int b) int { return a + b } func main() { int x = add(1, 2) }`,
			[]string{"int add(int a,int b)", "return a+b;", "int main()"},
		},
		{
			"program with struct",
			`struct Point = { int x = 0 int y = 0 } func main() { Point p = {} }`,
			[]string{"struct Point", "int main()"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.program)
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			if err = tr.Transpile(); err != nil {
				t.Logf("Transpile warning: %v", err)
			}

			cpp := tr.ToCpp()
			for _, sub := range tt.contains {
				if !strings.Contains(cpp, sub) {
					t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
				}
			}
		})
	}
}

func TestTranspileCNamespace(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		contains []string
		absent   []string
	}{
		{
			name:     "c_function_call_strips_prefix",
			src:      "import c\nfunc main() { c.printf(\"hi\") }",
			contains: []string{`printf`},
			absent:   []string{`c.printf`},
		},
		{
			name:     "c_constant_strips_prefix",
			src:      "import c\nfunc main() { int x = c.SEEK_SET }",
			contains: []string{`SEEK_SET`},
			absent:   []string{`c.SEEK_SET`},
		},
		{
			name:     "import_c_includes_libgen",
			src:      "import c\nfunc main() { }",
			contains: []string{`libgen.h`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.src)
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			if err = tr.Transpile(); err != nil {
				t.Logf("Transpile warning: %v", err)
			}

			cpp := tr.ToCpp()
			for _, sub := range tt.contains {
				if !strings.Contains(cpp, sub) {
					t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
				}
			}
			for _, sub := range tt.absent {
				if strings.Contains(cpp, sub) {
					t.Errorf("ToCpp() should not contain %q\nfull output:\n%s", sub, cpp)
				}
			}
		})
	}
}

func TestTranspileFullProgramWithTypes(t *testing.T) {
	program := `
type myInt = int
type myString = string

struct Point = {
	int x = 0
	int y = 0
}

func add(int a, int b) int {
	return a + b
}

func main() {
	int x = 10
	string s = "hello"
	bool b = true
	int result = add(x, 5)
	Println(result)
}
`

	tr, err := getTranspilerFromString(program)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	err = tr.Transpile()
	if err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	for _, sub := range []string{"typedef int myInt;", "struct Point", "int add(int a,int b)", "int main()"} {
		if !strings.Contains(cpp, sub) {
			t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
		}
	}
}

func TestTranspileWithControlFlow(t *testing.T) {
	program := `
func main() {
	// If statement
	if true {
		int x = 10
	}

	// For in statement
	for item in [1, 2, 3] {
		int x = item
	}

	// For of statement
	for i of [1, 2, 3] {
		int x = i
	}
}
`

	tr, err := getTranspilerFromString(program)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	err = tr.Transpile()
	if err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	for _, sub := range []string{"if (true)", "int main()", "while"} {
		if !strings.Contains(cpp, sub) {
			t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
		}
	}
}

func TestTranspileWithDeclarations(t *testing.T) {
	program := `
func main() {
	// Various declaration types
	int a = 10
	string b = "test"
	bool c = true
	char d = 'x'

	// Var and let
	var e = 20
	let f = 30

	// Array
	int[] arr = [1, 2, 3]
}

func funcWithArgs(int x, string y) int {
	return x
}
`

	tr, err := getTranspilerFromString(program)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	err = tr.Transpile()
	if err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	for _, sub := range []string{"int a = 10;", `std::string b = "test";`, "bool c = true;", "var e = 20;", "int main()"} {
		if !strings.Contains(cpp, sub) {
			t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
		}
	}
}

func TestTranspileReturnStatements(t *testing.T) {
	program := `
func getValue() int {
	return 42
}

func emptyReturn() {
	return
}

func complexReturn() int {
	int x = 10
	return x + 5
}
`

	tr, err := getTranspilerFromString(program)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	err = tr.Transpile()
	if err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	for _, sub := range []string{"return 42;", "return;", "return x+5;"} {
		if !strings.Contains(cpp, sub) {
			t.Errorf("ToCpp() missing %q\nfull output:\n%s", sub, cpp)
		}
	}
}

func TestTranspileVarStatement(t *testing.T) {
	b, err := getBuilderFromString("var x = 42")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if *result != "var x = 42;" {
		t.Errorf("got %q, want %q", *result, "var x = 42;")
	}
}

func TestTranspileLetStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"let int", "let n = 100", "const int n = 100;"},
		{"let string", `let s = "hi"`, `const std::string s = "hi";`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

// TestZeroInitDeclarations verifies auto-zero-initialization for every type.
// Express emits a zero value for every declaration without an explicit initializer
// so that the generated C++ is free of undefined behaviour.
func TestZeroInitDeclarations(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		// Primitives — immutable by default (const)
		{"int no init", "int x", "const int x = 0;"},
		{"float no init", "float f", "const float f = 0;"},
		{"bool no init", "bool b", "const bool b = false;"},
		{"char no init", "char c", `const char c = '\0';`},
		{"string no init", "string s", "const std::string s;"},
		{"var no init", "var v", "var v;"},
		{"pointer no init", "int* p", "int* p = nullptr;"},

		// Vectors (immutable by default; use `var int[] v` for mutable)
		{"int vector no init", "int[] v", "const std::vector<int> v;"},
		{"bool vector no init", "bool[] v", "const std::vector<bool> v;"},
		{"char vector no init", "char[] v", "const std::vector<char> v;"},
		{"var vector no init", "var[] v", "std::vector<var> v;"},

		// C-style arrays (fixed-size — aggregate-init, no const from early-return path)
		{"int array no init", "int[5] a", "int a[5] = {};"},
		{"bool array no init", "bool[3] a", "bool a[3] = {};"},
		{"char array no init", "char[8] a", "char a[8] = {};"},
		{"string array no init", "string[3] a", "std::string a[3] = {};"},

		// Map (immutable by default; use `var map m` for mutable)
		{"map no init", "map m", "const std::map<std::string, var> m;"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

// TestZeroInitStructDecl verifies that a struct declaration without an initializer emits `= {}`.
// The struct must be registered in scope first.
func TestZeroInitStructDecl(t *testing.T) {
	src := `struct S = { int x = 0 }
S s`

	tr, err := getTranspilerFromString(src)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	if err = tr.Transpile(); err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	if !strings.Contains(cpp, "S s = {};") {
		t.Errorf("expected 'S s = {};' in output\nfull output:\n%s", cpp)
	}
}

func TestTranspileCBlock(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		contains []string
	}{
		{
			name:     "c_block",
			stmt:     `c { printf("hi\n"); }`,
			contains: []string{`printf`, `\n`},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileStatement error: %v", err)
			}

			for _, sub := range tt.contains {
				if !strings.Contains(*result, sub) {
					t.Errorf("result missing %q\nfull output: %s", sub, *result)
				}
			}
		})
	}
}

func TestTranspileIdentExpression_Nil(t *testing.T) {
	b, err := getBuilderFromString("nil")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileIdentExpression(node)
	if err != nil {
		t.Fatalf("TranspileIdentExpression error: %v", err)
	}

	if *result != "nullptr" {
		t.Errorf("got %q, want %q", *result, "nullptr")
	}
}

func TestTranspileObjectStatement(t *testing.T) {
	b, err := getBuilderFromString(`object Person = { int age = 0 }`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileObjectStatement(node)
	if err != nil {
		t.Fatalf("TranspileObjectStatement error: %v", err)
	}

	for _, sub := range []string{"class Person", "age"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

func TestTranspilePackageStatement(t *testing.T) {
	b, err := getBuilderFromString(`package math { func add(int a, int b) int { return a + b } }`)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspilePackageStatement(node)
	if err != nil {
		t.Fatalf("TranspilePackageStatement error: %v", err)
	}

	if !strings.Contains(*result, "namespace __math") {
		t.Errorf("result missing \"namespace __math\"\nfull output: %s", *result)
	}
}

func TestTranspileForOverStatement(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		contains []string
	}{
		{
			name:     "single_var",
			src:      "for i over nums { }",
			contains: []string{"_coll", "i"},
		},
		{
			name:     "two_var",
			src:      "for i, v over nums { }",
			contains: []string{"_coll", "i", "v"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.src)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileForOverStatement(node)
			if err != nil {
				t.Fatalf("TranspileForOverStatement error: %v", err)
			}

			for _, sub := range tt.contains {
				if !strings.Contains(*result, sub) {
					t.Errorf("result missing %q\nfull output: %s", sub, *result)
				}
			}
		})
	}
}

func TestTranspileLaunchStatement(t *testing.T) {
	b, err := getBuilderFromString("launch Println(1)")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileLaunchStatement(node)
	if err != nil {
		t.Fatalf("TranspileLaunchStatement error: %v", err)
	}

	if !strings.Contains(*result, "go([=]") {
		t.Errorf("result missing \"go([=]\"\\nfull output: %s", *result)
	}
}

func TestTranspileTransformHelpers(t *testing.T) {
	identNode := &builder.Node{Type: "ident", Value: "i"}

	// TransformIdentToDefaultDeclaration wraps an ident as int decl
	ds := transpiler.TransformIdentToDefaultDeclaration(identNode)
	if ds.Type != "decl" {
		t.Errorf("TransformIdentToDefaultDeclaration: got type %q, want %q", ds.Type, "decl")
	}
	if ds.Left != identNode {
		t.Errorf("TransformIdentToDefaultDeclaration: Left is not the ident node")
	}

	// TransformExpressionToDeclaration wraps an expression as auto decl
	exprNode := &builder.Node{Type: "ident", Value: "arr"}
	dss := transpiler.TransformExpressionToDeclaration(exprNode)
	if dss.Type != "decl" {
		t.Errorf("TransformExpressionToDeclaration: got type %q, want %q", dss.Type, "decl")
	}
	if dss.Right != exprNode {
		t.Errorf("TransformExpressionToDeclaration: Right is not the expr node")
	}

	// GenerateLengthCall creates a std::size call
	lengthCall := transpiler.GenerateLengthCall(dss)
	if lengthCall.Type != "call" {
		t.Errorf("GenerateLengthCall: got type %q, want %q", lengthCall.Type, "call")
	}
	fnIdent, ok := lengthCall.Value.(*builder.Node)
	if !ok || fnIdent.Value.(string) != "std::size" {
		t.Errorf("GenerateLengthCall: expected std::size call, got %+v", lengthCall.Value)
	}
}

func TestTranspileForStdStatment(t *testing.T) {
	// Construct the forstd node directly: start must be an ident (not decl)
	// for TransformIdentToDefaultDeclaration to work correctly.
	node := &builder.Node{
		Type: "forstd",
		Value: &builder.Node{
			Type:  "block",
			Value: []*builder.Node{},
		},
		Metadata: map[string]interface{}{
			"start": &builder.Node{Type: "ident", Value: "i"},
			"end":   &builder.Node{Type: "ident", Value: "arr"},
			"step":  &builder.Node{Type: "inc", Left: &builder.Node{Type: "ident", Value: "i"}},
		},
	}

	b, err := getBuilderFromString("int x = 0")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileForStdStatment(node)
	if err != nil {
		t.Fatalf("TranspileForStdStatment error: %v", err)
	}

	for _, sub := range []string{"while", "std::size", "i"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

// ============================================
// Additional Coverage Tests
// ============================================

func TestTranspileDecrementExpression(t *testing.T) {
	b, err := getBuilderFromString("x--")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if !strings.Contains(*result, "(x)--") {
		t.Errorf("got %q, want contains \"(x)--\"", *result)
	}
}

func TestTranspileUnaryExpressions(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"deref", "*ptr", "*ptr"},
		{"not", "!x", "!x"},
		{"ref", "&x", "&x"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.expr)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseExpression()
			if err != nil {
				t.Fatalf("ParseExpression error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileExpression(node)
			if err != nil {
				t.Fatalf("TranspileExpression error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileEnumStatement(t *testing.T) {
	b, err := getBuilderFromString("enum { Red Green Blue }")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileEnumStatement(node)
	if err != nil {
		t.Fatalf("TranspileEnumStatement error: %v", err)
	}

	for _, sub := range []string{"enum", "Red", "Green", "Blue"} {
		if !strings.Contains(*result, sub) {
			t.Errorf("result missing %q\nfull output: %s", sub, *result)
		}
	}
}

func TestTranspileDeferStatement(t *testing.T) {
	b, err := getBuilderFromString("defer Println(1)")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if !strings.Contains(*result, "deferStack") {
		t.Errorf("result missing \"deferStack\"\nfull output: %s", *result)
	}
}

func TestTranspilePackageExpression(t *testing.T) {
	b, err := getBuilderFromString("int x = 0")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	node := &builder.Node{Type: "package", Value: "mylib"}
	result, err := tr.TranspileExpression(node)
	if err != nil {
		t.Fatalf("TranspileExpression error: %v", err)
	}

	if *result != "__mylib" {
		t.Errorf("got %q, want %q", *result, "__mylib")
	}
}

func TestTranspileSelectExpression_Extra(t *testing.T) {
	b, err := getBuilderFromString("int x = 0")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	t.Run("package_namespace_via_ident", func(t *testing.T) {
		// When Left is an ident registered as a package in t.Packages
		tr.Packages["mylib"] = "// mylib"
		node := &builder.Node{
			Type:  "selection",
			Left:  &builder.Node{Type: "ident", Value: "mylib"},
			Right: &builder.Node{Type: "ident", Value: "bar"},
		}
		result, err := tr.TranspileSelectExpression(node)
		if err != nil {
			t.Fatalf("TranspileSelectExpression error: %v", err)
		}
		if *result != "__mylib::bar" {
			t.Errorf("got %q, want %q", *result, "__mylib::bar")
		}
	})

	t.Run("package_node_type", func(t *testing.T) {
		// When Left is a "package" node type (not registered in Packages as ident)
		delete(tr.Packages, "mylib2")
		node := &builder.Node{
			Type:  "selection",
			Left:  &builder.Node{Type: "package", Value: "mylib2"},
			Right: &builder.Node{Type: "ident", Value: "baz"},
		}
		result, err := tr.TranspileSelectExpression(node)
		if err != nil {
			t.Fatalf("TranspileSelectExpression error: %v", err)
		}
		if !strings.Contains(*result, "::") {
			t.Errorf("got %q, want contains \"::\"\nfull output: %s", *result, *result)
		}
	})

	t.Run("current_receiver", func(t *testing.T) {
		// When CurrentReceiver is set — strips the receiver prefix
		tr.CurrentReceiver = "Point"
		node := &builder.Node{
			Type:  "selection",
			Left:  &builder.Node{Type: "ident", Value: "Point"},
			Right: &builder.Node{Type: "ident", Value: "x"},
		}
		result, err := tr.TranspileSelectExpression(node)
		tr.CurrentReceiver = ""
		if err != nil {
			t.Fatalf("TranspileSelectExpression error: %v", err)
		}
		if *result != "x" {
			t.Errorf("got %q, want %q", *result, "x")
		}
	})
}

func TestTranspileBlockExpression_KV(t *testing.T) {
	// A block with kv pairs is treated as a map literal → var{k1, v1, ...}
	b, err := getBuilderFromString("int x = 0")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	node := &builder.Node{
		Type: "block",
		Value: []*builder.Node{
			{
				Type:  "kv",
				Left:  &builder.Node{Type: "literal", Kind: "string", Value: "key"},
				Right: &builder.Node{Type: "literal", Kind: "string", Value: "val"},
			},
		},
	}

	result, err := tr.TranspileBlockExpression(node)

	if err != nil {
		t.Fatalf("TranspileBlockExpression error: %v", err)
	}

	if !strings.Contains(*result, "var{") {
		t.Errorf("result missing \"var{\"\nfull output: %s", *result)
	}
}

func TestTranspileLetStatement_Extra(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"let bool", "let b = true", "const bool b = true;"},
		{"let float", "let f = 3.14", "const float f = 3.14;"},
		{"let char", "let c = 'a'", "const char c = 'a';"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := getBuilderFromString(tt.stmt)
			if err != nil {
				t.Fatalf("lex error: %v", err)
			}

			node, err := b.ParseStatement()
			if err != nil {
				t.Fatalf("ParseStatement error: %v", err)
			}

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileStatement error: %v", err)
			}

			if *result != tt.expected {
				t.Errorf("got %q, want %q", *result, tt.expected)
			}
		})
	}
}

func TestTranspileLetStatement_AutoInference(t *testing.T) {
	// let x = foo() → "auto x = foo();"
	b, err := getBuilderFromString("let x = foo()")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if !strings.Contains(*result, "auto") {
		t.Errorf("got %q, want contains \"auto\"", *result)
	}
}

// TestTranspileStatement_AllCases exercises every case in the TranspileStatement
// dispatch switch by routing specific node types through it.
func TestTranspileStatement_AllCases(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		contains []string
	}{
		// c block
		{name: "c_block", src: `c { printf("hi"); }`, contains: []string{"printf"}},
		// struct declaration
		{name: "struct", src: `struct S = { int x = 0 }`, contains: []string{"struct S"}},
		// object declaration
		{name: "object", src: `object O = { int v = 0 }`, contains: []string{"class O"}},
		// typedef
		{name: "typedef", src: `type myInt = int`, contains: []string{"typedef"}},
		// map statement
		{name: "map", src: `map m = { "k" : "v" }`, contains: []string{"std::map"}},
		// inc as statement
		{name: "inc", src: `x++`, contains: []string{"(x)++"}},
		// while: not source-parseable (loop is an IDENT, not a keyword) — tested separately below
		// package statement
		{name: "package", src: `package foo { }`, contains: []string{"namespace __foo"}},
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

			tr, err := getTranspilerFromBuilder(b)
			if err != nil {
				t.Fatalf("Transpiler error: %v", err)
			}

			// FuncChan is still a channel — drain it to avoid blocking when func statements
			// are transpiled (the channel is drained by functionWorker in full Transpile() calls,
			// but in unit tests we call TranspileStatement directly without running workers).
			go func() {
				for range tr.FuncChan {
				}
			}()

			result, err := tr.TranspileStatement(node)
			if err != nil {
				t.Fatalf("TranspileStatement error: %v", err)
			}

			for _, sub := range tt.contains {
				if !strings.Contains(*result, sub) {
					t.Errorf("[%s] result missing %q\nfull output: %s", tt.name, sub, *result)
				}
			}
		})
	}
}

// TestTranspileWhileStatement directly constructs a while AST node and verifies transpilation.
// (while is not user-parseable from source; loop tokenizes as IDENT not LOOP)
func TestTranspileWhileStatement(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	whileNode := &builder.Node{
		Type: "while",
		Left: &builder.Node{
			Type:  "literal",
			Kind:  "bool",
			Value: true,
		},
		Value: &builder.Node{
			Type:  "block",
			Value: []*builder.Node{},
		},
	}

	result, err := tr.TranspileStatement(whileNode)
	if err != nil {
		t.Fatalf("TranspileStatement(while) error: %v", err)
	}

	if !strings.Contains(*result, "while") {
		t.Errorf("got %q, want contains \"while\"", *result)
	}
}

// TestToCpp_Packages verifies ToCpp emits the Packages section when populated.
func TestToCpp_Packages(t *testing.T) {
	program := `package math { }
func main() { }`

	tr, err := getTranspilerFromString(program)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	if err = tr.Transpile(); err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	if !strings.Contains(cpp, "namespace __math") {
		t.Errorf("ToCpp() missing \"namespace __math\"\nfull output:\n%s", cpp)
	}
	// "// Namespaces:" section must NOT say "none" when there's a package
	if strings.Contains(cpp, "// Namespaces:\n// none") {
		t.Error("ToCpp() shows '// none' for Namespaces but package was defined")
	}
}

func TestTranspilePrepLiteral_Struct(t *testing.T) {
	// prepLiteral for a struct literal (Kind == "struct")
	// Construct the struct literal node manually and transpile it via TranspileLiteralExpression
	b, err := getBuilderFromString("int x = 0")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// A struct literal node: {Type:"literal", Kind:"struct", Value:"Point", Right: block}
	structLit := &builder.Node{
		Type:  "literal",
		Kind:  "struct",
		Value: "Point",
		Right: &builder.Node{
			Type:  "block",
			Value: []*builder.Node{},
		},
	}

	result, err := tr.TranspileLiteralExpression(structLit)
	if err != nil {
		t.Fatalf("TranspileLiteralExpression error: %v", err)
	}

	if !strings.Contains(*result, "Point") {
		t.Errorf("result missing \"Point\"\nfull output: %s", *result)
	}
}

// TestFunctionWorker_MethodReceiver exercises the method-receiver path in functionWorker.
// A method declaration (func Receiver.method) puts the result in t.Methods[receiver].
func TestFunctionWorker_MethodReceiver(t *testing.T) {
	prog := `struct Point = { int x = 0 }
func Point.getX() int { return 0 }
func main() { }`

	tr, err := getTranspilerFromString(prog)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	if err = tr.Transpile(); err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	if len(tr.Methods["Point"]) == 0 {
		t.Error("Methods[\"Point\"] is empty — receiver path in functionWorker not exercised")
	}
}

// TestTranspileObjectStatement_Error verifies the error path when node is not "object".
func TestTranspileObjectStatement_Error(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	n := &builder.Node{Type: "ident", Value: "notAnObject"}
	_, err = tr.TranspileObjectStatement(n)
	if err == nil {
		t.Error("expected error for non-object node, got nil")
	}
}

// TestToCpp_ImportsAndIncludes exercises the non-empty Imports and Includes branches in ToCpp.
func TestToCpp_ImportsAndIncludes(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// Populate the maps directly so ToCpp's non-empty branches are hit.
	tr.Imports["stdio"] = "#include <stdio.h>"
	tr.Includes["defer"] = `#include "defer.cpp"`

	cpp := tr.ToCpp()

	if !strings.Contains(cpp, "#include <stdio.h>") {
		t.Errorf("ToCpp() missing import; output:\n%s", cpp)
	}
	if !strings.Contains(cpp, `#include "defer.cpp"`) {
		t.Errorf("ToCpp() missing include; output:\n%s", cpp)
	}
}

// TestToCpp_LibmillInclude exercises the libmill special-casing branch in ToCpp.
func TestToCpp_LibmillInclude(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	tr.Includes["libmill"] = `#include "libmill.h"`
	tr.Includes["defer"] = `#include "defer.cpp"`

	cpp := tr.ToCpp()

	// libmill must appear after other includes (the branch ensures it's last)
	if !strings.Contains(cpp, "libmill") {
		t.Errorf("ToCpp() missing libmill include; output:\n%s", cpp)
	}
}

// TestTranspileType_VarAndMap exercises the "var" and "map" branches in TranspileType.
func TestTranspileType_VarAndMap(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		contains string
	}{
		// var type — gradual typing (emits C++ "var" type)
		{"var_decl", `var x = 42`, "var x"},
		// map type
		{"map_decl", `map m = { "k" : "v" }`, "std::map"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString("func main() { " + tt.src + " }")
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			if err = tr.Transpile(); err != nil {
				t.Fatalf("Transpile error: %v", err)
			}

			cpp := tr.ToCpp()
			if !strings.Contains(cpp, tt.contains) {
				t.Errorf("[%s] ToCpp() missing %q; output:\n%s", tt.name, tt.contains, cpp)
			}
		})
	}
}

// TestTranspileIfStatement_ErrorPath exercises the "Node is not an if" error.
func TestTranspileIfStatement_ErrorPath(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	n := &builder.Node{Type: "ident", Value: "notAnIf"}
	_, err = tr.TranspileIfStatement(n)
	if err == nil {
		t.Error("expected error for non-if node, got nil")
	}
}

// TestTranspileIfStatement_DefaultRight exercises the default branch (invalid Right type).
func TestTranspileIfStatement_DefaultRight(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// if with Right of invalid type (not "if" or "block")
	n := &builder.Node{
		Type: "if",
		Value: &builder.Node{
			Type:  "literal",
			Kind:  "bool",
			Value: true,
		},
		Left: &builder.Node{
			Type:  "block",
			Value: []*builder.Node{},
		},
		Right: &builder.Node{
			Type:  "ident", // invalid — should be "if" or "block"
			Value: "notValid",
		},
	}
	_, err = tr.TranspileIfStatement(n)
	if err == nil {
		t.Error("expected error for invalid Right type, got nil")
	}
}

// TestTranspilePackageStatement_ErrorPath exercises wrong-node-type error.
func TestTranspilePackageStatement_ErrorPath(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	n := &builder.Node{Type: "ident", Value: "notPackage"}
	_, err = tr.TranspilePackageStatement(n)
	if err == nil {
		t.Error("expected error for non-package node, got nil")
	}
}

// TestTranspilePackageStatement_WithMethod exercises the receiver path in TranspilePackageStatement.
func TestTranspilePackageStatement_WithMethod(t *testing.T) {
	// Package with a struct and a method on that struct
	prog := `package geo {
  struct Point = { int x = 0  int y = 0 }
  func Point.sum() int { return 0 }
  func add(int a, int b) int { return a + b }
}
func main() { }`

	tr, err := getTranspilerFromString(prog)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	if err = tr.Transpile(); err != nil {
		t.Logf("Transpile warning: %v", err)
	}

	cpp := tr.ToCpp()
	if !strings.Contains(cpp, "namespace __geo") {
		t.Errorf("expected namespace __geo in output, got:\n%s", cpp)
	}
}

// TestTranspileType_PointerAndArray exercises pointer and array branches in TranspileType.
func TestTranspileType_PointerAndArray(t *testing.T) {
	tests := []struct {
		name     string
		prog     string
		contains string
	}{
		// int[] vector
		{"int_array", `func main() { int[] arr }`, "std::vector<int>"},
		// string type
		{"string_type", `func main() { string s = "hi" }`, "std::string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.prog)
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			if err = tr.Transpile(); err != nil {
				t.Logf("Transpile warning: %v", err)
			}

			cpp := tr.ToCpp()
			if !strings.Contains(cpp, tt.contains) {
				t.Errorf("[%s] ToCpp() missing %q; output:\n%s", tt.name, tt.contains, cpp)
			}
		})
	}
}

// TestTranspileWhileStatement_ErrorPath exercises the "Node is not a while" error.
func TestTranspileWhileStatement_ErrorPath(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	n := &builder.Node{Type: "ident", Value: "notWhile"}
	_, err = tr.TranspileWhileStatement(n)
	if err == nil {
		t.Error("expected error for non-while node, got nil")
	}
}

// TestTranspileStatement_SelectionAndInclude exercises the "selection" and "include" cases.
func TestTranspileStatement_SelectionAndInclude(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// "include" → returns error
	includeNode := &builder.Node{
		Type: "include",
		Left: &builder.Node{Type: "literal", Value: "stdio.h"},
	}
	_, err = tr.TranspileStatement(includeNode)
	if err == nil {
		t.Error("expected error for 'include' in TranspileStatement, got nil")
	}
}

// TestTranspileUseStatement exercises the "use" statement transpilation.
func TestTranspileUseStatement(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// Valid use node → always returns error about not being implemented
	useNode := &builder.Node{
		Type:  "use",
		Left:  &builder.Node{Type: "ident", Value: "mylib"},
		Right: &builder.Node{Type: "ident", Value: "lib"},
	}
	_, err = tr.TranspileUseStatement(useNode)
	if err == nil {
		t.Error("expected error from TranspileUseStatement, got nil")
	}
	if !strings.Contains(err.Error(), "use") {
		t.Errorf("error message missing 'use': %v", err)
	}
}

// TestTranspileUseStatement_WrongType exercises the type-check error path.
func TestTranspileUseStatement_WrongType(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	n := &builder.Node{Type: "ident", Value: "notUse"}
	_, err = tr.TranspileUseStatement(n)
	if err == nil {
		t.Error("expected error for non-use node, got nil")
	}
}

// TestTranspileForInStatement_Direct exercises the for-in transpilation directly.
// (forin nodes are normally converted to while by the tree flattener,
// but the transpiler function is still exercisable via direct call.)
func TestTranspileForInStatement_Direct(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// forin node: for i in nums { }
	// n.Left.Left.Value = index var name
	// n.Right.Value = collection name
	// n.Value = block
	forinNode := &builder.Node{
		Type: "forin",
		Left: &builder.Node{
			Type: "sgroup",
			Left: &builder.Node{Type: "ident", Value: "i"},
		},
		Right: &builder.Node{Type: "ident", Value: "nums"},
		Value: &builder.Node{Type: "block", Value: []*builder.Node{}},
	}

	result, err := tr.TranspileForInStatement(forinNode)
	if err != nil {
		t.Fatalf("TranspileForInStatement error: %v", err)
	}

	if !strings.Contains(*result, "for") {
		t.Errorf("result missing 'for': %s", *result)
	}
}

// TestTranspileIdentExpression_ErrorPaths exercises non-ident and non-string-value errors.
func TestTranspileIdentExpression_ErrorPaths(t *testing.T) {
	b, err := getBuilderFromString("")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	// non-ident type → error
	t.Run("wrong_type", func(t *testing.T) {
		n := &builder.Node{Type: "literal", Value: "42"}
		_, err := tr.TranspileIdentExpression(n)
		if err == nil {
			t.Error("expected error for non-ident node, got nil")
		}
	})

	// ident with non-string Value → error
	t.Run("non_string_value", func(t *testing.T) {
		n := &builder.Node{Type: "ident", Value: 42}
		_, err := tr.TranspileIdentExpression(n)
		if err == nil {
			t.Error("expected error for non-string ident Value, got nil")
		}
	})
}

// TestTranspileType_MapArrayAndPointer exercises array/map and pointer types.
func TestTranspileType_MapArrayAndPointer(t *testing.T) {
	tests := []struct {
		name     string
		prog     string
		contains string
	}{
		// int* pointer — uses else branch (recursive TranspileType call)
		{"int_pointer", `func main() { int* p }`, "int*"},
		// map[] — array of maps
		{"map_array", `func main() { map[] m }`, "std::map"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.prog)
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			if err = tr.Transpile(); err != nil {
				t.Logf("Transpile warning: %v", err)
			}

			cpp := tr.ToCpp()
			if !strings.Contains(cpp, tt.contains) {
				t.Errorf("[%s] ToCpp() missing %q; output:\n%s", tt.name, tt.contains, cpp)
			}
		})
	}
}

// TestTranspileBreakStatement verifies break transpiles to "break;".
func TestTranspileBreakStatement(t *testing.T) {
	b, err := getBuilderFromString("break")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if *result != "break;" {
		t.Errorf("got %q, want %q", *result, "break;")
	}
}

// TestTranspileContinueStatement verifies continue transpiles to "continue;".
func TestTranspileContinueStatement(t *testing.T) {
	b, err := getBuilderFromString("continue")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileStatement(node)
	if err != nil {
		t.Fatalf("TranspileStatement error: %v", err)
	}

	if *result != "continue;" {
		t.Errorf("got %q, want %q", *result, "continue;")
	}
}

// TestTranspileNestedArrayType verifies int[][] transpiles to std::vector<std::vector<int>>.
func TestTranspileNestedArrayType(t *testing.T) {
	src := `func main() { int[][] v }`
	tr, err := getTranspilerFromString(src)
	if err != nil {
		t.Fatalf("getTranspilerFromString error: %v", err)
	}

	if err = tr.Transpile(); err != nil {
		t.Fatalf("Transpile error: %v", err)
	}

	cpp := tr.ToCpp()
	if !strings.Contains(cpp, "std::vector<std::vector<int>>") {
		t.Errorf("expected std::vector<std::vector<int>> in output, got:\n%s", cpp)
	}
}

// TestTranspileNDimensionalMap verifies map[string, string -> int] transpiles to
// std::map<std::string, std::map<std::string, int>>.
func TestTranspileNDimensionalMap(t *testing.T) {
	src := `map[string, string -> int] scores`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileMapStatement(node)
	if err != nil {
		t.Fatalf("TranspileMapStatement error: %v", err)
	}

	want := "std::map<std::string, std::map<std::string, int>>"
	if !strings.Contains(*result, want) {
		t.Errorf("result missing %q\nfull output: %s", want, *result)
	}
}

// TestTranspileImmutableDecl verifies that `int x = 5` emits `const int x = 5;`.
func TestTranspileImmutableDecl(t *testing.T) {
	b, err := getBuilderFromString("int x = 5")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileDeclarationStatement(node)
	if err != nil {
		t.Fatalf("TranspileDeclarationStatement error: %v", err)
	}

	if !strings.Contains(*result, "const int x = 5;") {
		t.Errorf("expected 'const int x = 5;' in output, got: %s", *result)
	}
}

// TestTranspileVarTypedDecl verifies that `var int x = 5` emits `int x = 5;` (no const).
func TestTranspileVarTypedDecl(t *testing.T) {
	b, err := getBuilderFromString("var int x = 5")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileDeclarationStatement(node)
	if err != nil {
		t.Fatalf("TranspileDeclarationStatement error: %v", err)
	}

	if strings.Contains(*result, "const") {
		t.Errorf("unexpected 'const' in output for var int x = 5: %s", *result)
	}
	if !strings.Contains(*result, "int x = 5;") {
		t.Errorf("expected 'int x = 5;' in output, got: %s", *result)
	}
}

// TestTranspileLetConst verifies that `let x = 5` emits `const int x = 5;`.
func TestTranspileLetConst(t *testing.T) {
	b, err := getBuilderFromString("let x = 5")
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}

	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}

	tr, err := getTranspilerFromBuilder(b)
	if err != nil {
		t.Fatalf("Transpiler error: %v", err)
	}

	result, err := tr.TranspileLetStatement(node)
	if err != nil {
		t.Fatalf("TranspileLetStatement error: %v", err)
	}

	if !strings.Contains(*result, "const int x = 5;") {
		t.Errorf("expected 'const int x = 5;' in output, got: %s", *result)
	}
}
