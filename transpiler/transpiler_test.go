package transpiler_test

import (
	"fmt"
	"testing"

	ast "github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
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

// func TestTranspiler(t *testing.T) {
// 	b, err := getBuilderFromString("interface.expr")
// 	if err != nil {
// 		return nil, err
// 	}

// 	node, err := b.BuildAST()
// 	fmt.Println("scopeTree", b.ScopeTree)

// 	trans := transpiler.New(node, b, "blah", "blah")

// 	err = trans.TranspileInterfaceDecl()
// }

func TestGenHelper(t *testing.T) {

}
