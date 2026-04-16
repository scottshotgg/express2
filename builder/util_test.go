package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/scottshotgg/express-ast"
	"github.com/scottshotgg/express-lex"
	"github.com/scottshotgg/express-token"
	"github.com/scottshotgg/express2/builder"
	"github.com/scottshotgg/express2/pkg/logger"
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

func getBuilderFromString(test string) (*builder.Builder, error) {
	var tokens, err = getTokensFromString(test)
	if err != nil {
		return nil, err
	}

	for _, token := range tokens {
		fmt.Println(token)
	}

	return builder.New(tokens, logger.Noop()), nil
}

func printTokensFromBuilder(b *builder.Builder) {
	for _, token := range b.Tokens {
		fmt.Println(token)
	}
}

func printNode(node builder.Node) {
	nodeJSON, _ = json.Marshal(node)
	fmt.Printf(jsonFormatString, nodeJSON)
}

func parseStatement(t *testing.T, src string) *builder.Node {
	t.Helper()
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}
	node, err := b.ParseStatement()
	if err != nil {
		t.Fatalf("ParseStatement error: %v", err)
	}
	if node == nil {
		t.Fatal("node was nil")
	}
	return node
}

func parseExpression(t *testing.T, src string) *builder.Node {
	t.Helper()
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}
	node, err := b.ParseExpression()
	if err != nil {
		t.Fatalf("ParseExpression error: %v", err)
	}
	if node == nil {
		t.Fatal("node was nil")
	}
	return node
}

func buildAST(t *testing.T, src string) *builder.Node {
	t.Helper()
	b, err := getBuilderFromString(src)
	if err != nil {
		t.Fatalf("lex error: %v", err)
	}
	node, err := b.BuildAST()
	if err != nil {
		t.Fatalf("BuildAST error: %v", err)
	}
	return node
}
