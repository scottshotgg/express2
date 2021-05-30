package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	token "github.com/scottshotgg/express-token"
)

// Add symbols to the map when parsing

// type Variable struct {
// 	Type    string
// 	Changed bool
// 	MaxType string
// 	Value   interface{}
// 	Props   map[string]Variable
// }

// Start off just a map[string]*Node doing ident nodes for now

// scopeTree is the entire tree
var (
	scopeTree   *ScopeTree
	currentTree *ScopeTree
)

type ScopeTree struct {
	Lock *sync.RWMutex

	// Node is the node that the spawned the scope
	// node *Node

	Imports map[string]*ScopeTree

	// Table is the map of vars
	Vars map[string]*Node

	// Types is the map of types
	Types map[string]*TypeValue

	// Global is a pointer to the global scope
	Global *ScopeTree

	// Parent is a pointer to the parent scope
	Parent *ScopeTree

	// Children ...
	Children map[string]*ScopeTree
}

func (st ScopeTree) MarshalJSON() ([]byte, error) {
	// Marhshal up an anonymous struct with only the data we want
	return json.Marshal(struct {
		Imports  map[string]*ScopeTree
		Vars     map[string]*Node
		Types    map[string]*TypeValue
		Children map[string]*ScopeTree
	}{
		st.Imports,
		st.Vars,
		st.Types,
		st.Children,
	})
}

// NewScopeTree will create a new global scope in the scopeTree variable
// func NewScopeTree(node *Node) *ScopeTree {
func NewScopeTree() *ScopeTree {
	// Since this is the global scope, it has no `parent` and its `global` pointer is recursive
	var scopeTree = &ScopeTree{
		Lock: &sync.RWMutex{},
		// node:  node,
		Imports:  map[string]*ScopeTree{},
		Vars:     map[string]*Node{},
		Types:    map[string]*TypeValue{},
		Children: map[string]*ScopeTree{},
	}

	// Grab all the types from the typemap and insert them into the scope trees typemap
	for value := range token.TypeMap {
		scopeTree.Types[value] = &TypeValue{
			Type: PrimitiveValue,
			Kind: value,
		}
	}

	// This might have some problems ...
	scopeTree.Global = scopeTree

	return scopeTree
}

// NewChild enumerates a new child scope
// func (st *ScopeTree) NewChild(node *Node) *ScopeTree {
func (st *ScopeTree) NewChildScope(name string) (*ScopeTree, error) {
	// On a new child, it might be needed, we could either COPY everything from the other scope ...
	// 	OR
	// (easier) Just defer to recursing up in the Get

	// Check for a child with the same name already
	if st.Children[name] != nil {
		return nil, errors.Errorf("There is already a scope with that name; %s", name)
	}

	var scope = &ScopeTree{
		Lock: &sync.RWMutex{},
		// node:   node,
		// TODO: fix this
		Vars:     map[string]*Node{},
		Types:    map[string]*TypeValue{},
		Parent:   st,
		Global:   st.Global,
		Children: map[string]*ScopeTree{},
	}

	st.Children[name] = scope

	return scope, nil
}

func (st *ScopeTree) GetImports() map[string]*ScopeTree {
	return st.Global.Imports
}

// NewChild enumerates a new child scope
// func (st *ScopeTree) NewChild(node *Node) *ScopeTree {
func (st *ScopeTree) NewPackageScope(name string) (*ScopeTree, error) {
	// On a new child, it might be needed, we could either COPY everything from the other scope ...
	// 	OR
	// (easier) Just defer to recursing up in the Get

	// Check for a child with the same name already
	if st.Imports[name] != nil {
		return nil, errors.Errorf("There is already a scope with that name; %s", name)
	}

	var scope = &ScopeTree{
		Lock: &sync.RWMutex{},
		// node:   node,
		// TODO: fix this
		Vars:     map[string]*Node{},
		Types:    map[string]*TypeValue{},
		Parent:   st,
		Global:   st.Global,
		Children: map[string]*ScopeTree{},
	}

	st.Imports[name] = scope

	return scope, nil
}

// Leave exits the current scope and crawl up to the parent scope
func (st *ScopeTree) Leave() (*ScopeTree, error) {
	if st.Parent == nil {
		return nil, errors.New("Already in top level scope")
	}

	return st.Parent, nil
}

func (st *ScopeTree) Declare(ref *Node) error {
	var (
		refName string
		ok      bool
	)

	switch ref.Type {
	case "function":
		refName = ref.Kind

	case "decl":
		// ref.Left.Value should be the name of the ident
		refName, ok = ref.Left.Value.(string)
		if !ok {
			blob, _ := json.Marshal(ref)
			fmt.Println("blobberino:", string(blob))
			return errors.Errorf("Node value was not a string %+v", ref)
		}

	case "package":
		log.Fatalln("fuck you, package")

	default:
		return errors.Errorf("Node type is not supported for declaration: %+v", ref)
	}

	// If we have designated this as a new declaration, we only need to search the current scope
	// to make sure it is not already defined
	// if ref.Type == "decl" {
	// Lock the map
	st.Lock.Lock()
	defer st.Lock.Unlock()

	// Search for the reference name in the current scope's symbol table
	var scopeRef = st.Vars[refName]
	// If it is not equal to nil then we already have something under that name in the CURRENT scope
	if scopeRef != nil {
		return errors.Errorf("Variable already exists: %s\nScopeRef:%+v\nRef:%+v\n", refName, scopeRef, ref)
	}

	// Put the ref into the table
	st.Vars[refName] = ref

	return nil
	// }
}

func (st *ScopeTree) Assign(ref *Node) error {
	// ref.Left.Value should be the name of the ident
	var refName, ok = ref.Left.Value.(string)
	if !ok {
		return errors.Errorf("Node value was not a string %+v", ref.Left)
	}

	st.Lock.Lock()
	defer st.Lock.Unlock()

	var scopeRef = st.get(refName)
	// If it is equal to nil then we dont have something under that name in the ANY scope
	if scopeRef == nil {
		return errors.Errorf("Could not find variable: %+v", refName)
	}

	// assign to where ever this came from in the scope
	// TODO: I don't think this is going to work
	// *scopeRef = *ref

	return nil
	// }
}

func (st *ScopeTree) NewType(key string, ref *TypeValue) error {
	// If we have designated this as a new declaration, we only need to search the current scope
	// to make sure it is not already defined
	// if ref.Type == "decl" {
	// Lock the map
	st.Lock.Lock()
	defer st.Lock.Unlock()

	st.Types[key] = ref

	return nil
}

func (st *ScopeTree) GetImportedType(packageName, name string) *TypeValue {
	// If st is nil then we have a problem
	if st == nil {
		log.Printf("Current scope was nil ...")
		os.Exit(9)
	}

	// The Node in the current scope is not allowed to act as a ref as of right now
	// Search for the reference name in the current scope's symbol table
	st.Lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.Lock.Unlock()

	var imports = st.Global.Imports[packageName]
	if imports == nil {
		return nil
	}

	// Imports are always found in the global scope
	return imports.Types[name]
}

func (st *ScopeTree) GetType(name string) *TypeValue {
	// If st is nil then we have a problem
	if st == nil {
		log.Printf("Current scope was nil ...")
		os.Exit(9)
	}

	// The Node in the current scope is not allowed to act as a ref as of right now
	// Search for the reference name in the current scope's symbol table
	st.Lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.Lock.Unlock()

	var ref = st.Types[name]
	if ref != nil {
		fmt.Println("found it in the types", name)
		// If we get something from the current scope then return
		fmt.Println("ref", *ref)
		return ref
	}

	// If we have a parent then check that
	if st.Parent != nil {
		fmt.Println("going to the parents", name)
		// Fetch from the parent if our scope doesn't have it
		return st.Parent.GetType(name)
	}

	return nil
}

// Get will recursively search up the scope tree to verify whether that reference can be found
func (st *ScopeTree) Get(name string) *Node {
	// If st is nil then we have a problem
	if st == nil {
		log.Printf("Current scope was nil ...")
		os.Exit(9)
	}

	// The Node in the current scope is not allowed to act as a ref as of right now
	// Search for the reference name in the current scope's symbol table
	st.Lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.Lock.Unlock()

	var ref = st.Vars[name]
	if ref != nil {
		// If we get something from the current scope then return
		return ref
	}

	// If we have a parent then check that
	if st.Parent != nil {
		// Fetch from the parent if our scope doesn't have it
		return st.Parent.Get(name)
	}

	return nil
}

// Get will recursively search up the scope tree to verify whether that reference can be found
func (st *ScopeTree) get(name string) *Node {
	// If st is nil then we have a problem
	if st == nil {
		log.Printf("Current scope was nil ...")
		os.Exit(9)
	}

	var ref = st.Vars[name]
	if ref != nil {
		// If we get something from the current scope then return
		return ref
	}

	// If we have a parent then check that
	if st.Parent != nil {
		// Fetch from the parent if our scope doesn't have it
		return st.Parent.Get(name)
	}

	return nil
}

// Get will recursively search up the scope tree to verify whether that reference can be found
func (st *ScopeTree) Local(name string) *Node {
	// If st is nil then we have a problem
	if st == nil {
		log.Printf("Current scope was nil ...")
		os.Exit(9)
	}

	// The Node in the current scope is not allowed to act as a ref as of right now
	// Search for the reference name in the current scope's symbol table
	st.Lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.Lock.Unlock()

	return st.Vars[name]
}

// // SetGlobal will set the reference in the global scope
// func (st *ScopeTree) SetGlobal() {}

// // Global will search the global scope for the reference
// func (st *ScopeTree) GetGlobal() {}

// // Local will search the current scope for the reference
// func (st *ScopeTree) GetLocal() {}

// // FromTop will search for the reference in the scope X amount from the top
// func (st *ScopeTree) GetFromTop(name string, x int, allowParentSearch bool) {}

// // FromTop will search for the reference in the scope X amount from the top
// func (st *ScopeTree) GetFromBottom(name string, x int, allowParentSearch bool) {}
