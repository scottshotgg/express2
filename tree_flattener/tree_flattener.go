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

func identLTLength(lhs *builder.Node, rhs *builder.Node) *builder.Node {
	return &builder.Node{
		Type:  "comp",
		Value: "<",
		Left:  lhs,
		Right: &builder.Node{
			Type: "call",
			Value: &builder.Node{
				Type:  "ident",
				Value: "std::size",
				Metadata: map[string]interface{}{
					"args": &builder.Node{
						Type:  "egroup",
						Value: []*builder.Node{rhs},
					},
				},
			},
		},
	}
}

// Don't need any type information for this except for the array

func FlattenForIn(node *builder.Node) {

	arrayType := getArrayType(node.Metadata["end"].(*builder.Node))

	incVar := transformIdentToDecl("int", node.Metadata["start"].(*builder.Node))
	arrayVar := transformArrayToDecl(arrayType, node.Metadata["end"].(*builder.Node))
	while := &builder.Node{
		Type: "while",
		// Value: , // THIS NEEDS TO BE THE BLOCK AFTER IT IS CHECKED
		// i < len(arrayVar)
		Left: identLTLength(incVar.Left, arrayVar.Left),
	}
	fmt.Println("stuff", arrayType, incVar, arrayVar, while)

	// recurse, assign result to while.Value, return while.Value

	// make induction variable
	// make array if needed
	// make while loop with condition
	// recurse into block
}

func Flatten(node *builder.Node) {
	fmt.Println("node", node)
	stmts := node.Value.([]*builder.Node)

	for _, stmt := range stmts {
		fmt.Println("stmt", stmt)

		switch stmt.Type {
		case "forin":
			FlattenForIn(stmt)

		default:
			fmt.Println("not implemented", node)
			os.Exit(9)
		}
	}
}
