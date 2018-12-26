package tree_flattener

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/scottshotgg/express2/builder"
)

func getIntType() *builder.Node {
	return &builder.Node{
		Type: "type",
		// Kind: "int",
		Value: "int",
	}
}

// expects an egroup
func getArrayType(node *builder.Node) string {
	values := node.Value.([]*builder.Node)

	if len(values) < 1 {
		fmt.Println("not supporting empty array shit rn")
		os.Exit(9)
	}

	typeOf := values[0].Kind
	for _, value := range values[1:len(values)] {
		if value.Kind != typeOf {
			fmt.Println("not supporting dynamically typed arrays rn")
			os.Exit(9)
		}
	}

	return typeOf
}

/*
	forin &{decl  0xc00000e280 map[] 0xc00000e2d0 0xc00000e4b0}
	forin &{ident  i map[] <nil> <nil>}
	forin &{array  [0xc00000e320 0xc00000e370 0xc00000e3c0 0xc00000e410 0xc00000e460] map[] <nil> <nil>}
	forin &{type  array map[dim:[0xc00000a1a0]] <nil> <nil>}
*/

func transformIdentToDecl(typeOf string, node *builder.Node) *builder.Node {
	switch typeOf {
	case "int":
		return &builder.Node{
			Type:  "decl",
			Value: getIntType(),
			Left:  node,
			Right: &builder.Node{
				Type:  "literal",
				Kind:  "int",
				Value: 0,
			},
		}
	}

	return nil
}

func transformIdentAndValueToDecl(typeOf string, node, value *builder.Node) *builder.Node {
	switch typeOf {
	case "int":
		return &builder.Node{
			Type:  "decl",
			Value: getIntType(),
			Left:  node,
			Right: value,
		}
	}

	return nil
}

func transformArrayToDecl(typeOf string, node *builder.Node) *builder.Node {
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
			Value: "RANDOM_NAME_LATER",
		},
		Right: node,
	}
}

func transformIdentToAssignment(node *builder.Node, value *builder.Node) *builder.Node {
	return &builder.Node{
		Type:  "assignment",
		Left:  node,
		Right: value,
	}
}

func makeFunctionCall(node *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "call",
		Value: &builder.Node{
			Type:  "ident",
			Value: "std::size",
			Metadata: map[string]interface{}{
				"args": &builder.Node{
					Type:  "egroup",
					Value: []*builder.Node{node},
				},
			},
		},
	}
}

func makeLengthCall(node *builder.Node) *builder.Node {
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

func makeLTComp(lhs *builder.Node, rhs *builder.Node) *builder.Node {
	return &builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  lhs,
		Right: rhs,
	}
}

func makeIncrementOp(node *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "inc",
		Left: node,
	}
}

func makeRandomIdent() *builder.Node {
	return &builder.Node{
		Type:  "ident",
		Value: "RANDOM",
	}
}

// This needs to work with a scopeMap and then change the reference
// so that everyone referencing this var will feel the change
func anonymizeIdentName(n *builder.Node) error {
	if n == nil {
		return errors.New("Nil node ... anonymizeIdentName")
	}

	n.Value = n.Value.(string) + "_something_else"

	return nil
}

// Don't need any type information for this except for the array

func FlattenForIn(node *builder.Node) []*builder.Node {
	arrayType := getArrayType(node.Metadata["end"].(*builder.Node))
	// randomIdent := makeRandomIdent()
	start := node.Metadata["start"]
	if start == nil {
		// return nil, errors.New("No start amount ...")
		return nil
	}

	// err := anonymizeIdentName(start.(*builder.Node))
	// if err != nil {
	// 	return nil
	// }

	fmt.Printf("ident %+v\n", start)

	incVar := transformIdentToDecl("int", start.(*builder.Node))
	arrayVar := transformArrayToDecl(arrayType, node.Metadata["end"].(*builder.Node))
	block := node.Value.(*builder.Node)

	// Flatten all statements in the block
	var err = FlattenNode(block)
	if err != nil {
		log.Printf("err: %+v\n", err)
		return nil
	}

	stmts := append(block.Value.([]*builder.Node), makeIncrementOp(node.Metadata["start"].(*builder.Node)))
	while := &builder.Node{
		Type: "while",
		Value: &builder.Node{
			Type:  "block",
			Value: stmts,
		},
		Left: makeLTComp(incVar.Left, makeLengthCall(arrayVar.Left)),
	}

	// recurse, assign result to while.Value, return while.Value

	// make induction variable
	// make array if needed
	// make while loop with condition
	// recurse into block

	fmt.Println("something", []*builder.Node{
		incVar,
		arrayVar,
		while,
	})

	return []*builder.Node{
		incVar,
		arrayVar,
		while,
	}
}

// func FlattenBlock(node *builder.Node) []*builder.Node {
// 	var newStmts []*builder.Node

// 	for _, stmt := range node.Value.([]*builder.Node) {
// 		var err = FlattenNode(stmt)
// 		if err != nil {
// 			return err
// 		}
// 	}

// }

func makeIncludeNode(importName string) *builder.Node {
	return &builder.Node{
		Type: "include",
		Left: &builder.Node{
			Type:  "literal",
			Value: importName,
		},
	}
}

var includeChan = make(chan string, 10)

func Flatten(node *builder.Node) ([]*builder.Node, error) {
	if node.Type != "program" {
		return nil, errors.New("Flatten must be called with a tree; `program` node")
	}

	var (
		includes []*builder.Node
		wg       sync.WaitGroup
	)

	// Spin off a worker to process to includes that are found
	wg.Add(1)
	go func() {
		defer wg.Done()

		// Keep a map to track which includes we already have
		var (
			includesMap = map[string]struct{}{}
			ok          bool
		)

		for includeName := range includeChan {
			fmt.Println("includeName", includeName)
			// If it's already in the map then just skip it
			_, ok = includesMap[includeName]
			if ok {
				continue
			}

			includesMap[includeName] = struct{}{}
			includes = append(includes, makeIncludeNode(includeName))
		}
	}()

	// Flatten all nodes in the program
	for _, n := range node.Value.([]*builder.Node) {
		var err = FlattenNode(n)
		if err != nil {
			return nil, err
		}
	}

	// Close the channel and alert the import worker that we are done
	close(includeChan)

	// Wait for all extraneous imports to be transpiled
	wg.Wait()

	// Turn the node into a block, this will allow for all of the anonymous idents to
	// avoid confliction with current idents, but will also preserve the scope

	return includes, nil
}

func FlattenNode(node *builder.Node) error {
	var newStmts []*builder.Node

	switch node.Type {
	case "forin":
		includeChan <- "array"
		newStmts = append(newStmts, FlattenForIn(node)...)

	case "forof":
		includeChan <- "array"
		newStmts = append(newStmts, FlattenForOf(node)...)

	case "function":
		var err = FlattenNode(node.Value.(*builder.Node))
		if err != nil {
			return err
		}

	case "block":
		for _, stmt := range node.Value.([]*builder.Node) {
			var err = FlattenNode(stmt)
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

func FlattenForOf(node *builder.Node) []*builder.Node {
	randomIdent := makeRandomIdent()
	arrayType := getArrayType(node.Metadata["end"].(*builder.Node))
	incVar := transformIdentToDecl("int", node.Metadata["start"].(*builder.Node))
	indVar := transformIdentToDecl(arrayType, randomIdent)
	arrayVar := transformArrayToDecl(arrayType, node.Metadata["end"].(*builder.Node))
	block := node.Value.(*builder.Node)
	stmts := append(block.Value.([]*builder.Node), makeIncrementOp(node.Metadata["start"].(*builder.Node)))
	while := &builder.Node{
		Type: "while",
		// Value: , // THIS NEEDS TO BE THE BLOCK AFTER IT IS CHECKED
		Value: &builder.Node{
			Type: "block",
			Value: append(
				[]*builder.Node{
					transformIdentToAssignment(randomIdent, &builder.Node{
						Type:  "selection",
						Left:  node.Metadata["end"].(*builder.Node),
						Right: node.Metadata["start"].(*builder.Node),
					})},
				stmts...),
		},
		Left: makeLTComp(incVar.Left, arrayVar.Left),
	}

	return []*builder.Node{
		incVar,
		indVar,
		arrayVar,
		while,
	}
}
