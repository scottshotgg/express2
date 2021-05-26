package transpiler_test

import (
	"fmt"
	"io/ioutil"
	"testing"

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

	return builder.New(tokens), nil
}

func getASTFromString(test string) (*builder.Node, error) {
	b, err := getBuilderFromString(test)
	if err != nil {
		return nil, err
	}

	node, err := b.BuildAST()
	fmt.Println("scopeTree", b.ScopeTree)

	return node, err
}

func TestTranspiler(t *testing.T) {
	var contents, err = ioutil.ReadFile("simple_import.expr")
	if err != nil {
		t.Errorf("%+v", err)
	}

	ast, err := getASTFromString(string(contents))
	if err != nil {
		t.Errorf("%+v", err)
	}

	ast, err = builder.NewChecker(ast, builder.NewTypeResolver()).Execute()
	if err != nil {
		t.Errorf("%+v", err)
	}

	output, err := transpiler.NewC99(ast).Transpile()
	if err != nil {
		t.Errorf("%+v", err)
	}

	fmt.Println("output", *output)
}
