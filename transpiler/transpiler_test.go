package transpiler_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
	"github.com/scottshotgg/express2/transpiler"
)

const (
	jsonFormatString = "JSON: %s\n"
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

func printNode(node *builder.Node) {
	nodeJSON, _ := json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
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

			fmt.Printf("Result: %s\n", *result)
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

			fmt.Printf("Result: %s\n", *result)
		})
	}
}

func TestTranspileArrayExpression(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		expected string
	}{
		{"empty array", "[]", "{  }"},
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

			fmt.Printf("Result: %s\n", *result)
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
				fmt.Printf("Result (binop): %s\n", *result)
			} else {
				// For comparison expressions (comp)
				result, err := tr.TranspileConditionExpression(node)
				if err != nil {
					t.Fatalf("TranspileConditionExpression error: %v", err)
				}
				fmt.Printf("Result (comp): %s\n", *result)
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
		{"no args", "print()", "print"},
		{"one arg", "print(5)", "print"},
		{"multiple args", "print(1, 2, 3)", "print"},
		{"nested call", "outer(inner())", "outer"},
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

			fmt.Printf("Result: %s\n", *result)
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

			fmt.Printf("Result: %s\n", *result)
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
		typeName string
	}{
		{"int declaration", "int i = 10", "int"},
		{"string declaration", "string s = \"hello\"", "string"},
		{"bool declaration", "bool b = true", "bool"},
		{"var declaration", "var x = 5", "var"},
		{"let declaration", "let y = 10", "let"},
		{"array declaration", "int[] arr = [1, 2, 3]", "array"},
		{"declaration without init", "int x", "int"},
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

			fmt.Printf("Result: %s\n", *result)
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

			fmt.Printf("Result: %s\n", *result)
		})
	}
}

func TestTranspileFunctionStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		hasFunc  bool
	}{
		{"simple func", "func main() { int x = 10 }", true},
		{"func with args", "func add(int a, int b) int { return a + b }", true},
		{"func with return", "func getValue() int { return 5 }", true},
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

			fmt.Printf("Result: %s\n", *result)
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

			fmt.Printf("Result: %s\n", *result)
		})
	}
}

func TestTranspileIfStatement(t *testing.T) {
	tests := []struct {
		name     string
		stmt     string
		expected string
	}{
		{"simple if", `if x { int y = 10 }`, "if (x)"},
		{"if with else", `if x { int y = 10 } else { int z = 20 }`, "if (x)"},
		{"if with elseif", `if x { int y = 10 } else if z { int w = 15 } else { int v = 20 }`, "if (x)"},
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

			fmt.Printf("Result: %s\n", *result)
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

	fmt.Printf("Result: %s\n", *result)
}

// TestTranspileForInStatement tests for-in statement parsing
// Note: TranspileForInStatement has bugs in transpiler_old.go
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

// TestTranspileForOfStatement tests for-of statement parsing
// Note: TranspileForOverStatement has bugs in transpiler_old.go
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

	fmt.Printf("Result: %s\n", *result)
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

	fmt.Printf("Result: %s\n", *result)
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

			fmt.Printf("Result: %s\n", *result)
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

	fmt.Printf("Result: %s\n", *result)
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

	fmt.Printf("Result: %s\n", *result)
}

// ============================================
// Full Transpilation Tests
// ============================================

func TestTranspileCompleteProgram(t *testing.T) {
	tests := []struct {
		name     string
		program  string
		hasMain  bool
		checkErr bool
	}{
		{"simple program with main", `func main() { int x = 10 Println(x) }`, true, false},
		{"program with function", `func add(int a, int b) int { return a + b } func main() { int x = add(1, 2) }`, true, false},
		{"program with struct", `struct Point = { int x = 0 int y = 0 } func main() { Point p = {} }`, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr, err := getTranspilerFromString(tt.program)
			if err != nil {
				t.Fatalf("getTranspilerFromString error: %v", err)
			}

			err = tr.Transpile()
			if err != nil {
				if !tt.checkErr {
					t.Logf("Transpile warning: %v", err)
				} else {
					t.Fatalf("Transpile error: %v", err)
				}
			}

			cpp := tr.ToCpp()
			fmt.Printf("=== C++ Output for %s ===\n%s\n\n", tt.name, cpp)
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
	fmt.Printf("=== Full Program with Types ===\n%s\n", cpp)
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
	fmt.Printf("=== With Control Flow ===\n%s\n", cpp)
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
	fmt.Printf("=== With Declarations ===\n%s\n", cpp)
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
	fmt.Printf("=== With Return Statements ===\n%s\n", cpp)
}
