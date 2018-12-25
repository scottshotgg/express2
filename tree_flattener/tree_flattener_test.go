package tree_flattener_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/pkg/errors"
	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/transpiler"
	"github.com/scottshotgg/express2/tree_flattener"
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
	ast, err := getASTFromString(test)
	if err != nil {
		return nil, errors.Errorf("Could not create AST: %+v", err)
	}

	return transpiler.New(ast, name), nil
}

func TestFlattenForIn(t *testing.T) {
	testBytes, err := ioutil.ReadFile("test.expr")
	if err != nil {
		t.Fatalf("Could not read file: %+v", err)
	}

	test := string(testBytes)

	node, err := getASTFromString(test)
	if err != nil {
		t.Fatalf("Could not create transpiler: %+v", err)
	}

	tree_flattener.Flatten(node)

	fmt.Printf("\nNode: %+v\n", node)
}

func TestFlattenForOf(t *testing.T) {
	testBytes, err := ioutil.ReadFile("test.expr")
	if err != nil {
		t.Fatalf("Could not read file: %+v", err)
	}

	test := string(testBytes)

	node, err := getASTFromString(test)
	if err != nil {
		t.Fatalf("Could not create transpiler: %+v", err)
	}

	tree_flattener.Flatten(node)

	fmt.Printf("\nNode: %+v\n", node)
}
