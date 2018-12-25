package tree_flattener

import (
	"fmt"
	"os"

	"github.com/scottshotgg/express2/builder"
)

func getIntType() *builder.Node {
	return &builder.Node{
		Type: "type",
		Kind: "int",
	}
}

func transformIdentToDecl(typeOf string, node *builder.Node) *builder.Node {
	switch typeOf {
	case "int":
		return &builder.Node{
			Type:  "decl",
			Value: getIntType(),
			Left:  node,
			Right: &builder.Node{
				Type:  "int",
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

func transformArrayToDecl(typeOf string, node *builder.Node) *builder.Node {
	return &builder.Node{
		Type: "decl",
		// Type
		Value: &builder.Node{
			Type:     "type",
			Kind:     typeOf,
			Value:    "array",
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
			Metadata: map[string]interface{}{
				"args": &builder.Node{
					Type:  "egroup",
					Value: []*builder.Node{node},
				},
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

// Don't need any type information for this except for the array

func FlattenForIn(node *builder.Node) []*builder.Node {

	arrayType := getArrayType(node.Metadata["end"].(*builder.Node))

	incVar := transformIdentToDecl("int", node.Metadata["start"].(*builder.Node))
	arrayVar := transformArrayToDecl(arrayType, node.Metadata["end"].(*builder.Node))
	while := &builder.Node{
		Type:  "while",
		Value: node.Value,
		Left:  makeLTComp(incVar.Left, makeLengthCall(arrayVar.Left)),
	}

	// recurse, assign result to while.Value, return while.Value

	// make induction variable
	// make array if needed
	// make while loop with condition
	// recurse into block

	return []*builder.Node{
		incVar,
		arrayVar,
		while,
	}
}

func Flatten(node *builder.Node) {
	// fmt.Println("node", node)
	var (
		stmts    = node.Value.([]*builder.Node)
		newStmts = []*builder.Node{}
	)

	for _, stmt := range stmts {
		// fmt.Println("stmt", stmt)

		switch stmt.Type {
		case "forin":
			newStmts = append(newStmts, FlattenForIn(stmt)...)

		case "forof":
			newStmts = append(newStmts, FlattenForOf(stmt)...)

		default:
			newStmts = append(newStmts, stmt)
		}
	}

	node.Value = newStmts
}

func makeRandomIdent() *builder.Node {
	return &builder.Node{
		Type:  "ident",
		Value: "RANDOM",
	}
}

func transformIdentToAssignment(node *builder.Node, value *builder.Node) *builder.Node {
	return &builder.Node{
		Type:  "assignment",
		Left:  node,
		Right: value,
	}
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
