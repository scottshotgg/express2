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
