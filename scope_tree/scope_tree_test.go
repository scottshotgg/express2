package scope_tree_test

import (
	"fmt"
	"testing"

	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/scope_tree"
	"github.com/scottshotgg/express2/test"
)

const (
	astFormatString  = "ast: %+v\n"
	errFormatString  = "err: %+v\n"
	jsonFormatString = "JSON: %s\n"

	testProgram = `
		func main() {
			int i = 10
		}
	`
)

var (
	scopeTree *scope_tree.ScopeTree
	testNode  *builder.Node
)

func TestNew(t *testing.T) {
	// Parse something for a new scope
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("err %+v", err)
	}

	scopeTree = scope_tree.New(n)

	fmt.Printf("scopeTree: %+v\n", scopeTree)
}

func TestSetDeclaration(t *testing.T) {
	TestNew(t)

	var err error

	// Parse something for a new scope
	testNode, err = getStatementASTFromString(test.Tests[test.StatementTest]["decl"])
	if err != nil {
		t.Fatalf("err %+v", err)
	}

	err = scopeTree.Set(testNode)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestSetAssignment(t *testing.T) {
	// Set a variable first
	TestSetDeclaration(t)

	var (
		value = testNode.Right.Value.(int) + 10
		name  = testNode.Left.Value.(string)
		stmt  = fmt.Sprintf("%s = %d", name, value)
	)

	// Parse an assignment statement for the assignment variable
	node, err := getStatementASTFromString(stmt)
	if err != nil {
		t.Fatalf("err %+v", err)
	}

	err = scopeTree.Set(node)
	if err != nil {
		t.Fatalf("err %+v", err)
	}
}

func TestGetDeclaration(t *testing.T) {
	TestSetDeclaration(t)

	var ref = scopeTree.Get(testNode.Left.Value.(string))
	if ref == nil {
		t.Fatalf("Could not find node: %+v", ref)
	}

	if ref != testNode {
		t.Fatalf("Did not find expected node: %+v : %+v", ref, testNode)
	}

	if ref.Right.Value.(int) != testNode.Right.Value.(int) {
		t.Fatalf("Did not find expected value: %+v : %+v", ref.Value, testNode.Value)
	}

	fmt.Printf("ref: %+v\n", ref)
}

func TestGetAssignment(t *testing.T) {
	TestSetDeclaration(t)

	var ref = scopeTree.Get(testNode.Left.Value.(string))
	if ref == nil {
		t.Fatalf("Could not find node: %+v", ref)
	}

	fmt.Printf("ref: %+v\n", ref.Right)

	TestSetAssignment(t)

	ref = scopeTree.Get(testNode.Left.Value.(string))
	if ref == nil {
		t.Fatalf("Could not find node: %+v", ref)
	}

	if ref != testNode {
		t.Fatalf("Did not find expected node: %+v : %+v", ref, testNode)
	}

	if ref.Right.Value.(int) != testNode.Right.Value.(int) {
		t.Fatalf("Did not find expected value: %+v : %+v", ref.Value, testNode.Value)
	}

	fmt.Printf("ref: %+v\n", ref.Right)
}

func TestNewChild(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	parent := scope_tree.New(n)
	child := parent.NewChild(nil)
	if child == nil {
		t.Fatal("NewChild returned nil")
	}
}

func TestLeave(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	parent := scope_tree.New(n)
	child := parent.NewChild(nil)
	result, err := child.Leave()
	if err != nil {
		t.Fatalf("Leave: %v", err)
	}
	if result != parent {
		t.Error("Leave did not return parent scope")
	}
}

func TestLeaveFromGlobal(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	global := scope_tree.New(n)
	_, err = global.Leave()
	if err == nil {
		t.Fatal("Leave from global scope should return an error")
	}
}

func TestGetFromParent(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	parent := scope_tree.New(n)
	declNode, err := getStatementASTFromString("int x = 5")
	if err != nil {
		t.Fatalf("getStatementASTFromString: %v", err)
	}

	if err := parent.Set(declNode); err != nil {
		t.Fatalf("parent.Set: %v", err)
	}

	child := parent.NewChild(nil)
	got := child.Get("x")
	if got == nil {
		t.Fatal("child.Get did not find parent variable")
	}
}

func TestGetIsolation(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	parent := scope_tree.New(n)
	child := parent.NewChild(nil)

	declNode, err := getStatementASTFromString("int y = 7")
	if err != nil {
		t.Fatalf("getStatementASTFromString: %v", err)
	}

	if err := child.Set(declNode); err != nil {
		t.Fatalf("child.Set: %v", err)
	}

	got := parent.Get("y")
	if got != nil {
		t.Fatal("parent.Get should not find child-only variable")
	}
}

func TestDuplicateDecl(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	st := scope_tree.New(n)
	declNode, err := getStatementASTFromString("int z = 3")
	if err != nil {
		t.Fatalf("getStatementASTFromString: %v", err)
	}

	if err := st.Set(declNode); err != nil {
		t.Fatalf("first Set: %v", err)
	}

	// Second declaration of same name should fail
	if err := st.Set(declNode); err == nil {
		t.Fatal("second Set of same name should return an error")
	}
}

func TestSetAssignUndeclared(t *testing.T) {
	n, err := getASTFromString(testProgram)
	if err != nil {
		t.Fatalf("getASTFromString: %v", err)
	}

	st := scope_tree.New(n)
	assignNode, err := getStatementASTFromString("undeclared = 10")
	if err != nil {
		t.Fatalf("getStatementASTFromString: %v", err)
	}

	if err := st.Set(assignNode); err == nil {
		t.Fatal("Set assignment of undeclared variable should fail")
	}
}
