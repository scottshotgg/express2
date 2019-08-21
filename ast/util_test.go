package ast_test

import (
	"encoding/json"
	"fmt"

	"github.com/scottshotgg/express-ast"
	lex "github.com/scottshotgg/express-lex"
	token "github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/ast"
)

func getTokensFromString(s string) ([]token.Token, error) {
	// Lex and tokenize the source code
	var tokens, err = lex.New(s).Lex()
	if err != nil {
		return nil, err
	}

	fmt.Println("\nCompressing tokens ...")

	// Compress certain tokens;
	// i.e: `:` and `=` compress into `:=`
	return ast.CompressTokens(tokens)
}

func getASTFromString(test string) (*ast.AST, error) {
	var tokens, err = getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return ast.New(tokens), nil
}

func printTokensFromAST(b *ast.AST) {
	for _, token := range b.Tokens {
		fmt.Println(token)
	}
}

func printNode(node ast.Node) {
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}
