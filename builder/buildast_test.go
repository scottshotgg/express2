package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestBuildAST_SingleFunc(t *testing.T) {
	n := buildAST(t, "func main() { int i = 10 }")
	if n.Type != "program" {
		t.Fatalf("Type = %q, want program", n.Type)
	}
	stmts, ok := n.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 1 {
		t.Fatalf("got %d stmts, want 1", len(stmts))
	}
	if stmts[0].Type != "function" {
		t.Errorf("stmts[0].Type = %q, want function", stmts[0].Type)
	}
}

func TestBuildAST_MultiDecl(t *testing.T) {
	n := buildAST(t, "int x = 1 int y = 2 int z = 3")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 3 {
		t.Errorf("got %d stmts, want 3", len(stmts))
	}
	for i, s := range stmts {
		if s.Type != "decl" {
			t.Errorf("stmts[%d].Type = %q, want decl", i, s.Type)
		}
	}
}

func TestBuildAST_ImportCAndFunc(t *testing.T) {
	n := buildAST(t, "import c func main() { int i = 1 }")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 2 {
		t.Errorf("got %d stmts, want 2", len(stmts))
	}
	if stmts[0].Type != "import" {
		t.Errorf("stmts[0].Type = %q, want import", stmts[0].Type)
	}
	if stmts[1].Type != "function" {
		t.Errorf("stmts[1].Type = %q, want function", stmts[1].Type)
	}
}

func TestBuildAST_EmptyFunc(t *testing.T) {
	n := buildAST(t, "func main() { }")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok || len(stmts) != 1 {
		t.Fatalf("expected 1 stmt")
	}
	fn := stmts[0]
	if fn.Type != "function" {
		t.Fatalf("Type = %q, want function", fn.Type)
	}
	body, ok := fn.Value.(*builder.Node)
	if !ok || body == nil {
		t.Fatal("function body is nil or wrong type")
	}
	bodyStmts, _ := body.Value.([]*builder.Node)
	if len(bodyStmts) != 0 {
		t.Errorf("body has %d stmts, want 0", len(bodyStmts))
	}
}

func TestBuildAST_MultipleFuncs(t *testing.T) {
	n := buildAST(t, "func a() { } func b() { }")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 2 {
		t.Errorf("got %d stmts, want 2", len(stmts))
	}
	if stmts[0].Type != "function" || stmts[0].Kind != "a" {
		t.Errorf("stmts[0] = {Type:%q Kind:%q}, want function a", stmts[0].Type, stmts[0].Kind)
	}
	if stmts[1].Type != "function" || stmts[1].Kind != "b" {
		t.Errorf("stmts[1] = {Type:%q Kind:%q}, want function b", stmts[1].Type, stmts[1].Kind)
	}
}

func TestBuildAST_LetAndTypedef(t *testing.T) {
	n := buildAST(t, "type myInt = int let x = 42")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 2 {
		t.Errorf("got %d stmts, want 2", len(stmts))
	}
	if stmts[0].Type != "typedef" {
		t.Errorf("stmts[0].Type = %q, want typedef", stmts[0].Type)
	}
	if stmts[1].Type != "let" {
		t.Errorf("stmts[1].Type = %q, want let", stmts[1].Type)
	}
}

func TestBuildAST_FuncWithReturn(t *testing.T) {
	n := buildAST(t, "func add(int a, int b) int { return a + b }")
	stmts, ok := n.Value.([]*builder.Node)
	if !ok || len(stmts) != 1 {
		t.Fatalf("expected 1 stmt")
	}
	fn := stmts[0]
	if fn.Type != "function" || fn.Kind != "add" {
		t.Errorf("got {Type:%q Kind:%q}, want function add", fn.Type, fn.Kind)
	}
	if fn.Metadata["args"] == nil {
		t.Error("Metadata[args] is nil")
	}
	if fn.Metadata["returns"] == nil {
		t.Error("Metadata[returns] is nil")
	}
}
