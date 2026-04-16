package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

func newDeclNode(name string, value int) *builder.Node {
	return &builder.Node{
		Type:  "decl",
		Left:  &builder.Node{Type: "ident", Value: name},
		Right: &builder.Node{Type: "literal", Kind: "int", Value: value},
	}
}

func TestScopeTree_DeclareAndGet(t *testing.T) {
	st := builder.NewScopeTree()
	n := newDeclNode("x", 42)
	if err := st.Declare(n); err != nil {
		t.Fatalf("Declare: %v", err)
	}
	got := st.Get("x")
	if got == nil {
		t.Fatal("Get returned nil")
	}
	if got != n {
		t.Error("Get returned different pointer")
	}
}

func TestScopeTree_DuplicateDeclareFails(t *testing.T) {
	st := builder.NewScopeTree()
	n := newDeclNode("x", 10)
	if err := st.Declare(n); err != nil {
		t.Fatalf("first Declare: %v", err)
	}
	if err := st.Declare(n); err == nil {
		t.Fatal("second Declare should have returned an error")
	}
}

func TestScopeTree_ChildScopeInheritsParent(t *testing.T) {
	parent := builder.NewScopeTree()
	n := newDeclNode("x", 1)
	if err := parent.Declare(n); err != nil {
		t.Fatalf("Declare: %v", err)
	}
	child, err := parent.NewChildScope("child")
	if err != nil {
		t.Fatalf("NewChildScope: %v", err)
	}
	got := child.Get("x")
	if got == nil {
		t.Fatal("child.Get did not find parent variable")
	}
}

func TestScopeTree_ChildScopeIsolation(t *testing.T) {
	parent := builder.NewScopeTree()
	child, err := parent.NewChildScope("child")
	if err != nil {
		t.Fatalf("NewChildScope: %v", err)
	}
	n := newDeclNode("y", 2)
	if err := child.Declare(n); err != nil {
		t.Fatalf("child.Declare: %v", err)
	}
	got := parent.Get("y")
	if got != nil {
		t.Fatal("parent.Get should not find child-only variable")
	}
}

func TestScopeTree_DuplicateChildScopeFails(t *testing.T) {
	st := builder.NewScopeTree()
	if _, err := st.NewChildScope("same"); err != nil {
		t.Fatalf("first NewChildScope: %v", err)
	}
	if _, err := st.NewChildScope("same"); err == nil {
		t.Fatal("second NewChildScope with same name should fail")
	}
}

func TestScopeTree_Leave(t *testing.T) {
	parent := builder.NewScopeTree()
	child, err := parent.NewChildScope("child")
	if err != nil {
		t.Fatalf("NewChildScope: %v", err)
	}
	got, err := child.Leave()
	if err != nil {
		t.Fatalf("Leave: %v", err)
	}
	if got != parent {
		t.Error("Leave did not return parent scope")
	}
}

func TestScopeTree_LeaveFromGlobalFails(t *testing.T) {
	st := builder.NewScopeTree()
	_, err := st.Leave()
	if err == nil {
		t.Fatal("Leave from global scope should return an error")
	}
}

func TestScopeTree_PrimitiveTypesPreloaded(t *testing.T) {
	st := builder.NewScopeTree()
	for _, typeName := range []string{"int", "float", "bool", "string", "char"} {
		tv := st.GetType(typeName)
		if tv == nil {
			t.Errorf("GetType(%q) = nil, want preloaded TypeValue", typeName)
		}
	}
}

func TestScopeTree_CustomType(t *testing.T) {
	st := builder.NewScopeTree()
	tv := &builder.TypeValue{Type: builder.PrimitiveValue, Kind: "myType"}
	if err := st.NewType("myType", tv); err != nil {
		t.Fatalf("NewType: %v", err)
	}
	got := st.GetType("myType")
	if got == nil {
		t.Fatal("GetType returned nil for custom type")
	}
	if got.Kind != "myType" {
		t.Errorf("Kind = %q, want myType", got.Kind)
	}
}

func TestScopeTree_AssignExisting(t *testing.T) {
	st := builder.NewScopeTree()
	declNode := newDeclNode("x", 10)
	if err := st.Declare(declNode); err != nil {
		t.Fatalf("Declare: %v", err)
	}
	assignNode := &builder.Node{
		Type:  "assignment",
		Left:  &builder.Node{Type: "ident", Value: "x"},
		Right: &builder.Node{Type: "literal", Kind: "int", Value: 20},
	}
	if err := st.Assign(assignNode); err != nil {
		t.Fatalf("Assign: %v", err)
	}
}

func TestScopeTree_AssignUndeclaredFails(t *testing.T) {
	st := builder.NewScopeTree()
	assignNode := &builder.Node{
		Type:  "assignment",
		Left:  &builder.Node{Type: "ident", Value: "undeclared"},
		Right: &builder.Node{Type: "literal", Kind: "int", Value: 0},
	}
	if err := st.Assign(assignNode); err == nil {
		t.Fatal("Assign undeclared variable should return an error")
	}
}

func TestScopeTree_DeclareFunctionAndGet(t *testing.T) {
	st := builder.NewScopeTree()
	funcNode := &builder.Node{
		Type: "function",
		Kind: "myFunc",
	}
	if err := st.Declare(funcNode); err != nil {
		t.Fatalf("Declare function: %v", err)
	}
	got := st.Get("myFunc")
	if got == nil {
		t.Fatal("Get returned nil for declared function")
	}
	if got != funcNode {
		t.Error("Get returned different pointer")
	}
}

func TestScopeTree_ImportScope(t *testing.T) {
	st := builder.NewScopeTree()
	pkgScope, err := st.NewPackageScope("mypkg")
	if err != nil {
		t.Fatalf("NewPackageScope: %v", err)
	}
	tv := &builder.TypeValue{Type: builder.PrimitiveValue, Kind: "int"}
	if err := pkgScope.NewType("MyType", tv); err != nil {
		t.Fatalf("pkgScope.NewType: %v", err)
	}
	got := st.GetImportedType("mypkg", "MyType")
	if got == nil {
		t.Fatal("GetImportedType returned nil")
	}
}

func TestScopeTree_Local(t *testing.T) {
	st := builder.NewScopeTree()
	n := newDeclNode("z", 99)
	if err := st.Declare(n); err != nil {
		t.Fatalf("Declare: %v", err)
	}
	got := st.Local("z")
	if got == nil {
		t.Fatal("Local returned nil for declared variable")
	}
	notGot := st.Local("nonexistent")
	if notGot != nil {
		t.Fatal("Local returned non-nil for nonexistent variable")
	}
}
