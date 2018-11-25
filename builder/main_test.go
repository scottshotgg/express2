package builder_test

import (
	"testing"

	"github.com/scottshotgg/express2/builder"
)

const (
	astFormatString  = "ast: %+v\n"
	errFormatString  = "err: %+v\n"
	jsonFormatString = "JSON: %s\n"
)

var (
	b        *builder.Builder
	node     *builder.Node
	err      error
	nodeJSON []byte
)

func TestNew(t *testing.T) {
	if builder.New(nil) == nil {
		t.Errorf(errFormatString, "Builder was nil for some reason")
	}
}

// func TestBuildAST(t *testing.T) {
// 	var totalString string

// 	var i int
// 	// Test each one individually
// 	for _, stmt := range statementTestMap {
// 		// if i > len(statementTestMap)-18 {
// 		// 	break
// 		// }

// 		i++

// 		// Accumulate a string containing all statements
// 		totalString += stmt

// 		b, err = getBuilderFromString(stmt)
// 		if err != nil {
// 			t.Errorf(errFormatString, err)
// 		}

// 		node, err = b.BuildAST()
// 		if err != nil {
// 			fmt.Println("before", b.Tokens[b.Index-3])
// 			fmt.Println("before", b.Tokens[b.Index-2])
// 			fmt.Println("before", b.Tokens[b.Index-1])
// 			t.Errorf(errFormatString, err)
// 			fmt.Println("after", b.Tokens[b.Index+1])
// 		}

// 		nodeJSON, _ = json.Marshal(node)
// 		fmt.Printf(jsonFormatString, nodeJSON)
// 	}

// 	b, err = getBuilderFromString(totalString)
// 	if err != nil {
// 		t.Errorf(errFormatString, err)
// 	}

// 	node, err = b.BuildAST()
// 	if err != nil {
// 		fmt.Println("before", b.Tokens[b.Index-3])
// 		fmt.Println("before", b.Tokens[b.Index-2])
// 		fmt.Println("before", b.Tokens[b.Index-1])
// 		t.Errorf(errFormatString, err)
// 		fmt.Println("after", b.Tokens[b.Index+1])
// 	}

// 	nodeJSON, _ = json.Marshal(node)
// 	fmt.Printf(jsonFormatString, nodeJSON)
// }

// func (b *Builder) TestBuildAST(t *testing.T) {

// }
