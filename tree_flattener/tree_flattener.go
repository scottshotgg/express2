package tree_flattener

import (
	"errors"
	"fmt"
	"sync"

	"github.com/scottshotgg/express2/builder"
)

type Flattener struct {
	IncludeChan chan string
	Wg          sync.WaitGroup
	IdentCounter int
}

func New() *Flattener {
	return &Flattener{
		IncludeChan: make(chan string, 10),
	}
}

func (f *Flattener) getIntType() *builder.Node {
	return &builder.Node{
		Type: "type",
		// Kind: "int",
		Value: "int",
	}
}

// expects an egroup
func (f *Flattener) getArrayType(node *builder.Node) (string, error) {
	if node.Type == "ident" {
		return node.Value.(string), nil
	}

	values, ok := node.Value.([]*builder.Node)
	if !ok {
		return "", errors.New("getArrayType: node value is not []*builder.Node")
	}

	if len(values) < 1 {
		return "", errors.New("getArrayType: empty array literal not supported")
	}

	typeOf := values[0].Kind
	for _, value := range values[1:] {
		if value.Kind != typeOf {
			return "", errors.New("getArrayType: mixed-type arrays not yet supported")
		}
	}

	return typeOf, nil
}

/*
	forin &{decl  0xc00000e280 map[] 0xc00000e2d0 0xc00000e4b0}
	forin &{ident  i map[] <nil> <nil>}
	forin &{array  [0xc00000e320 0xc00000e370 0xc00000e3c0 0xc00000e410 0xc00000e460] map[] <nil> <nil>}
	forin &{type  array map[dim:[0xc00000a1a0]] <nil> <nil>}
*/

func (f *Flattener) transformIdentToDecl(typeOf string, value interface{}, node *builder.Node) *builder.Node {
	switch typeOf {
	case "int":
		return &builder.Node{
			Type:     "decl",
			Value:    f.getIntType(),
			Left:     node,
			Metadata: map[string]interface{}{"mutable": true},
			Right: &builder.Node{
				Type:  "literal",
				Kind:  "int",
				Value: value,
			},
		}

	case "auto":
		return &builder.Node{
			Type: "decl",
			Value: &builder.Node{
				Type: "type",
				// Kind: "int",
				Value: "auto",
			},
			Left:     node,
			Metadata: map[string]interface{}{"mutable": true},
			Right: &builder.Node{
				Type:  "literal",
				Kind:  "auto",
				Value: value,
			},
		}
	}

	return nil
}

func (f *Flattener) transformArrayToDecl(typeOf string, node *builder.Node) *builder.Node {
	// Generate a unique identifier name for the array variable
	// If node.Value is a string, use it; otherwise generate a name
	var identValue string
	switch v := node.Value.(type) {
	case string:
		identValue = v
	default:
		// Generate a unique name like "arr_0", "arr_1", etc.
		identValue = fmt.Sprintf("arr_%d", f.IdentCounter)
		f.IdentCounter++
	}

	return &builder.Node{
		Type: "decl",
		// Type
		Value: &builder.Node{
			Type:     "type",
			Kind:     typeOf,
			Value:    "auto",
			Metadata: map[string]interface{}{
				// THIS NEEDS TO BE SET TO BE STATIC SIZE OF THE ARRAY
			},
		},
		// ident
		Left: &builder.Node{
			Type:  "ident",
			Value: identValue,
		},
		Right: node,
	}
}

func (f *Flattener) makeLengthCall(node *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "call",
		Value: &builder.Node{
			Type:  "ident",
			Value: "std::size",
		},
		Metadata: map[string]interface{}{
			"args": &builder.Node{
				Type:  "egroup",
				Value: []*builder.Node{node},
			},
		},
	}
}

func (f *Flattener) makeLTComp(lhs *builder.Node, rhs *builder.Node) *builder.Node {
	return &builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  lhs,
		Right: rhs,
	}
}

func (f *Flattener) makeIncrementOp(node *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "inc",
		Left: node,
	}
}

// Don't need any type information for this except for the array

func (f *Flattener) FlattenForIn(node *builder.Node) ([]*builder.Node, error) {
	start := node.Metadata["start"]
	if start == nil {
		return nil, nil
	}

	endNode := node.Metadata["end"].(*builder.Node)

	// int i = 0
	keyVar := f.transformIdentToDecl("int", 0, start.(*builder.Node))

	// Determine the array ident to use in the while condition.
	// If the end is already a named variable, use it directly to avoid
	// "auto numbers = numbers;" self-referential declarations.
	var arrayIdent *builder.Node
	var extraDecls []*builder.Node
	if endNode.Type == "ident" {
		arrayIdent = endNode
	} else {
		arrayType, err := f.getArrayType(endNode)
		if err != nil {
			return nil, err
		}
		arrayVar := f.transformArrayToDecl(arrayType, endNode)
		arrayIdent = arrayVar.Left
		extraDecls = []*builder.Node{arrayVar}
	}

	block := node.Value.(*builder.Node)

	// Flatten all statements in the block
	if err := f.FlattenNode(block); err != nil {
		return nil, err
	}

	stmts := append(block.Value.([]*builder.Node), f.makeIncrementOp(start.(*builder.Node)))
	while := &builder.Node{
		Type: "while",
		Value: &builder.Node{
			Type:  "block",
			Value: stmts,
		},
		Metadata: node.Metadata,
		Right:    endNode,
		Left:     f.makeLTComp(keyVar.Left, f.makeLengthCall(arrayIdent)),
	}

	result := []*builder.Node{keyVar}
	result = append(result, extraDecls...)
	result = append(result, while)
	return result, nil
}

// func (f * Flattener) FlattenBlock(node *builder.Node) []*builder.Node {
// 	var newStmts []*builder.Node

// 	for _, stmt := range node.Value.([]*builder.Node) {
// 		var err = FlattenNode(stmt)
// 		if err != nil {
// 			return err
// 		}
// 	}

// }

func (f *Flattener) makeIncludeNode(importName string) *builder.Node {
	return &builder.Node{
		Type: "include",
		Left: &builder.Node{
			Type:  "literal",
			Value: importName,
		},
	}
}

func (f *Flattener) Flatten(node *builder.Node) ([]*builder.Node, error) {
	if node.Type != "program" {
		return nil, errors.New("Flatten must be called with a tree; `program` node")
	}

	var (
		includes []*builder.Node
	)

	// Spin off a worker to process to includes that are found
	f.Wg.Add(1)
	go func() {
		defer f.Wg.Done()

		// Keep a map to track which includes we already have
		var (
			includesMap = map[string]struct{}{}
			ok          bool
		)

		for includeName := range f.IncludeChan {
			// If it's already in the map then just skip it
			_, ok = includesMap[includeName]
			if ok {
				continue
			}

			includesMap[includeName] = struct{}{}
			includes = append(includes, f.makeIncludeNode(includeName))
		}
	}()

	f.IncludeChan <- "iostream"

	// Flatten all nodes in the program
	for _, n := range node.Value.([]*builder.Node) {
		var err = f.FlattenNode(n)
		if err != nil {
			return nil, err
		}
	}

	// Close the channel and alert the import worker that we are done
	close(f.IncludeChan)

	// Wait for all extraneous imports to be transpiled
	f.Wg.Wait()

	// Turn the node into a block, this will allow for all of the anonymous idents to
	// avoid confliction with current idents, but will also preserve the scope

	return includes, nil
}

func (f *Flattener) FlattenNode(node *builder.Node) error {
	var newStmts []*builder.Node

	switch node.Type {
	case "forin":
		f.IncludeChan <- "array"
		nodes, err := f.FlattenForIn(node)
		if err != nil {
			return err
		}
		newStmts = append(newStmts, nodes...)

	case "forof":
		f.IncludeChan <- "array"
		nodes, err := f.FlattenForOf(node)
		if err != nil {
			return err
		}
		newStmts = append(newStmts, nodes...)

	case "function":
		var err = f.FlattenNode(node.Value.(*builder.Node))
		if err != nil {
			return err
		}

	case "block":
		for _, stmt := range node.Value.([]*builder.Node) {
			var err = f.FlattenNode(stmt)
			if err != nil {
				return err
			}
		}

	/*
		We were gonna call `flatten` on the entire tree and recurse through it
		implement that step later when we need it.
	*/
	default:
		// The node can stay the same
		return nil
	}

	// If we acquired new statements then the block is now that
	if len(newStmts) > 0 {
		*node = builder.Node{
			Type:  "block",
			Value: newStmts,
		}
	}

	return nil
}

func (f *Flattener) FlattenForOf(node *builder.Node) ([]*builder.Node, error) {
	// Generate a unique internal index counter name
	idxIdentName := fmt.Sprintf("_idx_%d", f.IdentCounter)
	f.IdentCounter++

	idxIdent := &builder.Node{Type: "ident", Value: idxIdentName}
	// int _idx_N = 0
	idxVar := f.transformIdentToDecl("int", 0, idxIdent)

	endNode := node.Metadata["end"].(*builder.Node)

	// Determine the array ident to use in the while condition.
	// If the end is already a named variable, use it directly to avoid
	// "auto numbers = numbers;" self-referential declarations.
	var arrayIdent *builder.Node
	var extraDecls []*builder.Node
	if endNode.Type == "ident" {
		arrayIdent = endNode
	} else {
		arrayType, err := f.getArrayType(endNode)
		if err != nil {
			return nil, err
		}
		arrayVar := f.transformArrayToDecl(arrayType, endNode)
		arrayIdent = arrayVar.Left
		extraDecls = []*builder.Node{arrayVar}
	}

	block := node.Value.(*builder.Node)

	// Flatten all statements in the block
	if err := f.FlattenNode(block); err != nil {
		return nil, err
	}

	// auto value = arr_N[_idx_N]  -- declared at the top of the while body
	elemDecl := &builder.Node{
		Type: "decl",
		Value: &builder.Node{
			Type:  "type",
			Value: "auto",
		},
		Left: node.Metadata["start"].(*builder.Node),
		Right: &builder.Node{
			Type:  "index",
			Left:  arrayIdent,
			Right: idxIdent,
		},
	}

	// Increment _idx BEFORE the user body so that `continue` cannot skip it.
	// The element value is captured into `elemDecl` first, so the pre-increment
	// index is still used for the element access.
	stmts := append(
		[]*builder.Node{elemDecl, f.makeIncrementOp(idxIdent)},
		block.Value.([]*builder.Node)...,
	)

	while := &builder.Node{
		Type: "while",
		Value: &builder.Node{
			Type:  "block",
			Value: stmts,
		},
		Left: f.makeLTComp(idxIdent, f.makeLengthCall(arrayIdent)),
	}

	result := []*builder.Node{idxVar}
	result = append(result, extraDecls...)
	result = append(result, while)
	return result, nil
}
