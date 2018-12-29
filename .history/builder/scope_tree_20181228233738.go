package builder

import (
	"log"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express2/builder"
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
	lock *sync.RWMutex

	// Node is the node that the spawned the scope
	node *Node

	// Table is the map of symbols
	vars map[string]*builder.Node

	// Types is the map of symbols
	types map[string]*builder.TypeValue

	// Parent is a pointer to the parent scope
	parent *ScopeTree

	// Global is a pointer to the global scope
	global *ScopeTree
}

// New will create a new global scope in the scopeTree variable
func NewScopeTree(node *Node) *ScopeTree {
	// Since this is the global scope, it has no `parent` and its `global` pointer is recursive
	var scopeTree = &ScopeTree{
		lock:  &sync.RWMutex{},
		node:  node,
		table: map[string]*Node{},
	}

	scopeTree.global = scopeTree

	return scopeTree
}

// NewChild enumerates a new child scope
// func (st *ScopeTree) NewChild(node *Node) *ScopeTree {
func (st *ScopeTree) NewChild() *ScopeTree {
	// On a new child, it might be needed, we could either COPY everything from the other scope ...
	// 	OR
	// (easier) Just defer to recursing up in the Get
	return &ScopeTree{
		lock: &sync.RWMutex{},
		// node:   node,
		table:  map[string]*Node{},
		parent: st,
		global: st.global,
	}
}

// Leave exits the current scope and crawl up to the parent scope
func (st *ScopeTree) Leave() (*ScopeTree, error) {
	if st.parent == nil {
		return nil, errors.New("Already in top level scope")
	}

	return st.parent, nil
}

func (st *ScopeTree) Declare(ref *Node) error {
	// ref.Left.Value should be the name of the ident
	var refName, ok = ref.Value.(string)
	if !ok {
		return errors.Errorf("Node value was not a string %+v", ref)
	}

	// If we have designated this as a new declaration, we only need to search the current scope
	// to make sure it is not already defined
	// if ref.Type == "decl" {
	// Lock the map
	st.lock.Lock()
	defer st.lock.Unlock()

	// Search for the reference name in the current scope's symbol table
	var scopeRef = st.table[refName]
	// If it is not equal to nil then we already have something under that name in the CURRENT scope
	if scopeRef != nil {
		return errors.Errorf("Variable already exists: \nScopeRef:%+v\nRef:%+v\n", scopeRef, ref)
	}

	// Put the ref into the table
	st.table[refName] = ref

	return nil
	// }
}

func (st *ScopeTree) Assign(ref *Node) error {
	// ref.Left.Value should be the name of the ident
	var refName, ok = ref.Value.(string)
	if !ok {
		return errors.Errorf("Node value was not a string %+v", ref)
	}

	// If we have designated this as a new declaration, we only need to search the current scope
	// to make sure it is not already defined
	// if ref.Type == "decl" {
	// Lock the map
	st.lock.Lock()
	defer st.lock.Unlock()

	// If the ref is an assignment then we are expecting that variable to already be there
	// if ref.Type == "assignment" {
	// Not sure if this is going to work
	// Get the value
	var scopeRef = st.get(refName)
	// If it is equal to nil then we dont have something under that name in the ANY scope
	if scopeRef == nil {
		return errors.Errorf("Could not find variable: %+v", ref)
	}

	// assign to where ever this came from in the scope
	*scopeRef = *ref

	return nil
	// }
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
	st.lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.lock.Unlock()

	var ref = st.table[name]
	if ref != nil {
		// If we get something from the current scope then return
		return ref
	}

	// If we have a parent then check that
	if st.parent != nil {
		// Fetch from the parent if our scope doesn't have it
		return st.parent.Get(name)
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

	var ref = st.table[name]
	if ref != nil {
		// If we get something from the current scope then return
		return ref
	}

	// If we have a parent then check that
	if st.parent != nil {
		// Fetch from the parent if our scope doesn't have it
		return st.parent.Get(name)
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
	st.lock.Lock()
	// Don't know if we need to recursively lock ... it seems likely
	defer st.lock.Unlock()

	return st.table[name]
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
