package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

// assertBinop checks that n is a binop node with the given operator.
func assertBinop(t *testing.T, n *builder.Node, op string) {
	t.Helper()
	if n == nil {
		t.Fatalf("node is nil, want binop %q", op)
	}
	if n.Type != "binop" {
		t.Fatalf("Type = %q, want binop", n.Type)
	}
	if n.Value.(string) != op {
		t.Fatalf("Value = %q, want %q", n.Value, op)
	}
}

// assertComp checks that n is a comp node with the given operator.
func assertComp(t *testing.T, n *builder.Node, op string) {
	t.Helper()
	if n == nil {
		t.Fatalf("node is nil, want comp %q", op)
	}
	if n.Type != "comp" {
		t.Fatalf("Type = %q, want comp", n.Type)
	}
	if n.Value.(string) != op {
		t.Fatalf("Value = %q, want %q", n.Value, op)
	}
}

// assertLitInt checks that n is an int literal with the given value.
func assertLitInt(t *testing.T, n *builder.Node, v int) {
	t.Helper()
	if n == nil {
		t.Fatalf("node is nil, want literal int %d", v)
	}
	if n.Type != "literal" {
		t.Fatalf("Type = %q, want literal", n.Type)
	}
	if n.Value.(int) != v {
		t.Fatalf("Value = %v, want %d", n.Value, v)
	}
}

// assertIdent checks that n is an ident with the given name.
func assertIdent(t *testing.T, n *builder.Node, name string) {
	t.Helper()
	if n == nil {
		t.Fatalf("node is nil, want ident %q", name)
	}
	if n.Type != "ident" {
		t.Fatalf("Type = %q, want ident", n.Type)
	}
	if n.Value.(string) != name {
		t.Fatalf("Value = %q, want %q", n.Value, name)
	}
}

func TestPrattPrecedence_MulBeforeAdd(t *testing.T) {
	// 1 + 2 * 3 → +(1, *(2, 3))
	n := parseExpression(t, "1 + 2 * 3")
	assertBinop(t, n, "+")
	assertLitInt(t, n.Left, 1)
	assertBinop(t, n.Right, "*")
	assertLitInt(t, n.Right.Left, 2)
	assertLitInt(t, n.Right.Right, 3)
}

func TestPrattPrecedence_AddLeftAssoc(t *testing.T) {
	// 1 + 2 + 3 → +(+(1, 2), 3)
	n := parseExpression(t, "1 + 2 + 3")
	assertBinop(t, n, "+")
	assertBinop(t, n.Left, "+")
	assertLitInt(t, n.Left.Left, 1)
	assertLitInt(t, n.Left.Right, 2)
	assertLitInt(t, n.Right, 3)
}

func TestPrattPrecedence_MulLeftAssoc(t *testing.T) {
	// 2 * 3 * 4 → *(*(2, 3), 4)
	n := parseExpression(t, "2 * 3 * 4")
	assertBinop(t, n, "*")
	assertBinop(t, n.Left, "*")
	assertLitInt(t, n.Left.Left, 2)
	assertLitInt(t, n.Left.Right, 3)
	assertLitInt(t, n.Right, 4)
}

func TestPrattPrecedence_CompLowerThanAdd(t *testing.T) {
	// a + 1 < b + 2 → <(+(a,1), +(b,2))
	n := parseExpression(t, "a + 1 < b + 2")
	assertComp(t, n, "<")
	assertBinop(t, n.Left, "+")
	assertIdent(t, n.Left.Left, "a")
	assertLitInt(t, n.Left.Right, 1)
	assertBinop(t, n.Right, "+")
	assertIdent(t, n.Right.Left, "b")
	assertLitInt(t, n.Right.Right, 2)
}

func TestPrattPrecedence_EqLowestBinop(t *testing.T) {
	// a + 1 == b * 2 → ==(+(a,1), *(b,2))
	n := parseExpression(t, "a + 1 == b * 2")
	assertComp(t, n, "==")
	assertBinop(t, n.Left, "+")
	assertBinop(t, n.Right, "*")
}

func TestPrattPrecedence_CallHigherThanAdd(t *testing.T) {
	// foo() + 1 → +(call(foo), 1)
	n := parseExpression(t, "foo() + 1")
	assertBinop(t, n, "+")
	if n.Left == nil || n.Left.Type != "call" {
		t.Fatalf("Left.Type = %q, want call", n.Left.Type)
	}
	assertLitInt(t, n.Right, 1)
}

func TestPrattPrecedence_IndexHigherThanAdd(t *testing.T) {
	// arr[0] + 1 → +(idx(arr,0), 1)
	n := parseExpression(t, "arr[0] + 1")
	assertBinop(t, n, "+")
	if n.Left == nil || n.Left.Type != "index" {
		t.Fatalf("Left.Type = %q, want index", n.Left.Type)
	}
	assertLitInt(t, n.Right, 1)
}

func TestPrattPrecedence_ComplexBinop(t *testing.T) {
	// 9 + 8 * 7 → +(9, *(8, 7))
	n := parseExpression(t, "9 + 8 * 7")
	assertBinop(t, n, "+")
	assertLitInt(t, n.Left, 9)
	assertBinop(t, n.Right, "*")
	assertLitInt(t, n.Right.Left, 8)
	assertLitInt(t, n.Right.Right, 7)
}

func TestPrattPrecedence_SelectionHigherThanMul(t *testing.T) {
	// a.b * 2 → *(sel(a,b), 2)
	n := parseExpression(t, "a.b * 2")
	assertBinop(t, n, "*")
	if n.Left == nil || n.Left.Type != "selection" {
		t.Fatalf("Left.Type = %q, want selection", n.Left.Type)
	}
	assertLitInt(t, n.Right, 2)
}
