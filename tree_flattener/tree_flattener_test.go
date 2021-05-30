package tree_flattener_test

import (
	"fmt"
	"os"

	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/transpiler"
)

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	tokens, err := lex.New(s).Lex()
	if err != nil {
		return nil, err
	}

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	return ast.CompressTokens(tokens)
}

func getBuilderFromString(test string) (*builder.Builder, error) {
	tokens, err := getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return builder.New(tokens), nil
}

func getASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.BuildAST()
}

func getTranspilerFromString(test, name string) (*transpiler.Transpiler, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	ast, err := b.BuildAST()
	if err != nil {
		return nil, err
	}

	return transpiler.New(ast, b, name, os.Getenv("EXPRPATH")), nil
}

func getStatementASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	return b.ParseStatement()
}

// func TestFlattenForIn(t *testing.T) {
// 	testBytes, err := ioutil.ReadFile("test.expr")
// 	if err != nil {
// 		t.Fatalf("Could not read file: %+v", err)
// 	}

// 	test := string(testBytes)

// 	node, err := getASTFromString(test)
// 	if err != nil {
// 		t.Fatalf("Could not create transpiler: %+v", err)
// 	}

// 	tree_flattener.Flatten(node)

// 	fmt.Printf("\nNode: %+v\n", node)
// }

// func TestFlattenForOf(t *testing.T) {
// 	testBytes, err := ioutil.ReadFile("test.expr")
// 	if err != nil {
// 		t.Fatalf("Could not read file: %+v", err)
// 	}

// 	test := string(testBytes)

// 	node, err := getASTFromString(test)
// 	if err != nil {
// 		t.Fatalf("Could not create transpiler: %+v", err)
// 	}

// 	tree_flattener.Flatten(node)

// 	fmt.Printf("\nNode: %+v\n", node)
// }

// func TestTranspileFlattenedForIn(t *testing.T) {
// 	node, err := getStatementASTFromString(test.Tests[test.StatementTest]["forin"])
// 	if err != nil {
// 		t.Fatalf("Could not create transpiler: %+v", err)
// 	}

// 	tree_flattener.Flatten(node)

// 	fmt.Printf("\nNode: %+v\n", node)

// 	// tr, err := getTranspilerFromString(test, "main")
// 	// if err != nil {
// 	// 	t.Fatalf("Could not create transpiler: %+v", err)
// 	// }

// 	cpp, err := transpiler.Statement(node)
// 	if err != nil {
// 		t.Fatalf("Could not transpile to C++: %+v", err)
// 	}

// 	fmt.Printf("\nC++: %s\n\n", *cpp)
// }

// func TestFlatten(t *testing.T) {
// 	testBytes, err := ioutil.ReadFile("test.expr")
// 	if err != nil {
// 		t.Fatalf("Could not read file: %+v", err)
// 	}

// 	var test = string(testBytes)

// 	node, err := getASTFromString(test)
// 	if err != nil {
// 		t.Fatalf("Could not create transpiler: %+v", err)
// 	}

// 	// fmt.Printf("Before: %+v\n", node)
// 	stringy, _ := json.Marshal(node)

// 	fmt.Println(string(stringy))

// 	tree_flattener.Flatten(node)

// 	// fmt.Printf("After: %+v\n", node)
// 	stringy, _ = json.Marshal(node)

// 	fmt.Println(string(stringy))
// }
