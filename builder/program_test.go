package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func TestProgram(t *testing.T) {
	src := `func main() { int i = 10 }`
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex: %v", err)
	}
	node, err := b.BuildAST()
	if err != nil {
		t.Fatalf("BuildAST: %v", err)
	}
	if node.Type != "program" {
		t.Errorf("Type = %q, want program", node.Type)
	}
	stmts, ok := node.Value.([]*builder.Node)
	if !ok {
		t.Fatalf("Value is not []*builder.Node")
	}
	if len(stmts) != 1 {
		t.Errorf("got %d stmts, want 1", len(stmts))
	}
}
