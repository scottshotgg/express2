package symbol_table

import (
	"fmt"

	"github.com/scottshotgg/express2/builder"
)

// crawl the tree and assign values in the symbol map

type Variable struct {
	Type    string
	Changed bool
	MaxType string
	Value   interface{}
	Props   map[string]Variable
}

var (
	symbols = map[string]Variable{}
)

func Stuff(tree *builder.Node) {
	fmt.Println("tree", tree)

	stmts := tree.Value.([]*builder.Node)
	for _, stmt := range stmts {
		fmt.Println("stmt", stmt)

		switch stmt.Type {
		case "decl":
			// value is the type
			// left is the ident
			// right is the value
			symbols[stmt.Left.Value.(string)] = Variable{
				Type:  stmt.Value.(*builder.Node).Value.(string),
				Value: stmt.Right.Value,
			}
		}
	}

	fmt.Println("symbols", symbols)
}
