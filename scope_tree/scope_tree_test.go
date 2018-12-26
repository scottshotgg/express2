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

func TestNewChild(t *testing.T) {}
