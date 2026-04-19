package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestNewChecker(t *testing.T) {
	ast := &builder.Node{Type: "program", Value: []*builder.Node{}}
	ch := builder.NewChecker(ast)
	if ch == nil {
		t.Fatal("NewChecker returned nil")
	}
}

func TestChecker_DummyPass(t *testing.T) {
	ast := &builder.Node{
		Type:  "program",
		Value: []*builder.Node{},
	}
	dummy := builder.NewDummyPass("test_pass")
	changed, err := dummy.Check(ast)
	if err != nil {
		t.Fatalf("DummyPass.Check error: %v", err)
	}
	// DummyPass returns true
	if !changed {
		t.Error("DummyPass.Check returned false, want true")
	}
}

func TestChecker_DummyPassName(t *testing.T) {
	dummy := builder.NewDummyPass("my_pass")
	if dummy.Name() != "my_pass" {
		t.Errorf("Name() = %q, want my_pass", dummy.Name())
	}
}

func TestChecker_AddPass(t *testing.T) {
	ast := &builder.Node{Type: "program", Value: []*builder.Node{}}
	ch := builder.NewChecker(ast)
	ch.AddPass(builder.NewDummyPass("added"))
	// No direct way to inspect ps — just verify it doesn't panic
}

func TestChecker_Execute_EmptyPasses(t *testing.T) {
	ast := &builder.Node{Type: "program", Value: []*builder.Node{}}
	ch := builder.NewChecker(ast)
	result, err := ch.Execute()
	if err != nil {
		t.Fatalf("Execute with no passes error: %v", err)
	}
	if result == nil {
		t.Fatal("Execute returned nil node")
	}
	if result.Type != "program" {
		t.Errorf("result.Type = %q, want program", result.Type)
	}
}

func TestNewTypeResolver(t *testing.T) {
	tr := builder.NewTypeResolver()
	if tr == nil {
		t.Fatal("NewTypeResolver returned nil")
	}
	if tr.Name() != "type_resolver" {
		t.Errorf("Name() = %q, want type_resolver", tr.Name())
	}
}

// TestTypeResolver_Check exercises the TypeResolver.Check branches directly.
func TestTypeResolver_Check(t *testing.T) {
	// Helper to build a program AST and run the type resolver through compiler path
	checkProgram := func(t *testing.T, src string) {
		t.Helper()
		b, err := getBuilderFromString(src)
		if err != nil {
			t.Fatalf("lex error: %v", err)
		}
		ast, err := b.BuildAST()
		if err != nil {
			t.Fatalf("BuildAST error: %v", err)
		}
		tr := builder.NewTypeResolverWithScope(b.ScopeTree)
		_, err = tr.Check(ast)
		if err != nil {
			t.Fatalf("Check error: %v", err)
		}
	}

	// literal: covered
	t.Run("literal", func(t *testing.T) {
		checkProgram(t, "func main() { int x = 42 }")
	})

	// decl with "ident" type — user-defined struct type
	t.Run("decl_ident_type", func(t *testing.T) {
		checkProgram(t, `
struct Point = { int x = 0  int y = 0 }
func main() {
  Point p = { x = 1  y = 2 }
}`)
	})

	// decl with "pointer" subtype — c pointer
	t.Run("decl_pointer_type", func(t *testing.T) {
		checkProgram(t, `
import c
func main() {
  c.FILE* f = c.fopen("x", "r")
}`)
	})

	// let statement — type inference
	t.Run("let_literal", func(t *testing.T) {
		checkProgram(t, `func main() { let x = 42 }`)
	})

	t.Run("let_bool", func(t *testing.T) {
		checkProgram(t, `func main() { let b = true }`)
	})

	t.Run("let_string", func(t *testing.T) {
		checkProgram(t, `func main() { let s = "hello" }`)
	})

	// function with explicit return type and arguments
	t.Run("function_with_args_and_return", func(t *testing.T) {
		checkProgram(t, `func add(int a, int b) int { return a + b }`)
	})

	// function without return (void)
	t.Run("function_void", func(t *testing.T) {
		checkProgram(t, `func greet() { int x = 1 }`)
	})

	// function named main (gets int return type injected)
	t.Run("function_main", func(t *testing.T) {
		checkProgram(t, `func main() { int x = 0 }`)
	})

	// typedef — type alias
	t.Run("typedef", func(t *testing.T) {
		checkProgram(t, `type MyInt = int
func main() { MyInt x = 42 }`)
	})

	// block with multiple statements
	t.Run("block_multiple_stmts", func(t *testing.T) {
		checkProgram(t, `func main() {
  int a = 1
  int b = 2
  int c = 3
}`)
	})

	// egroup with type ident resolution
	t.Run("egroup_resolution", func(t *testing.T) {
		checkProgram(t, `func foo() int { return 0 }
func main() { int x = 0 }`)
	})

	// sgroup (function arguments)
	t.Run("sgroup_args", func(t *testing.T) {
		checkProgram(t, `func greet(string name, int age) { int x = 0 }
func main() { }`)
	})

	// package statement
	t.Run("package_statement", func(t *testing.T) {
		checkProgram(t, `package math { func add(int a, int b) int { return a + b } }
func main() { }`)
	})

	// No-op node types — just verify they don't error
	t.Run("noop_nodes", func(t *testing.T) {
		tr := builder.NewTypeResolver()
		noopTypes := []string{
			"call", "selection", "assignment", "inc", "dec", "if", "for",
			"forin", "forof", "forover", "forstd", "binop", "comp", "deref",
			"ref", "index", "ident", "type", "return", "struct", "while",
			"kv", "map", "not", "enum", "defer",
		}
		for _, typ := range noopTypes {
			n := &builder.Node{Type: typ, Value: []*builder.Node{}}
			_, err := tr.Check(n)
			if err != nil {
				t.Errorf("Check(%q) unexpected error: %v", typ, err)
			}
		}
	})

	// default error path — unknown node type
	t.Run("default_error", func(t *testing.T) {
		tr := builder.NewTypeResolver()
		n := &builder.Node{Type: "unknown_xyz_type_987"}
		_, err := tr.Check(n)
		if err == nil {
			t.Error("expected error for unknown node type, got nil")
		}
	})
}
