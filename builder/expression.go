package builder

import (
	token "github.com/scottshotgg/express-token"
)

func (b *Builder) ParseGroupOfExpressions() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, b.AppendTokenToError("Could not get group of expressions")
	}

	// Skip over the left paren token
	b.Index++

	var (
		expr  *Node
		exprs []*Node
		err   error
	)

	for b.Index < len(b.Tokens) && b.Tokens[b.Index].Type != token.RParen {
		expr, err = b.ParseExpression()
		if err != nil {
			return expr, err
		}

		// Step over the expression token
		b.Index++

		exprs = append(exprs, expr)

		// Check and skip over the separator
		if b.Index < len(b.Tokens) && b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	return &Node{
		Type:  "egroup",
		Value: exprs,
	}, nil
}

// This should go in the function map on the scopetree
var cFuncs = map[string]bool{
	"Println": true,
	"printf":  true,
	"sleep":   true,
	"msleep":  true,
	"now":     true,
}

// ParseExpression is the public entry point for expression parsing.
// After it returns, b.Index is on the last consumed token.
func (b *Builder) ParseExpression() (*Node, error) {
	return b.prattParse(PrecNone)
}

func (b *Builder) ParseArrayExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LBracket {
		return nil, b.AppendTokenToError("Could not get array expression")
	}

	// Skip over the left bracket token
	b.Index++

	var (
		expr  *Node
		exprs []*Node
		err   error
	)

	for b.Index < len(b.Tokens) && b.Tokens[b.Index].Type != token.RBracket {
		expr, err = b.ParseExpression()
		if err != nil {
			return expr, err
		}

		b.Index++

		exprs = append(exprs, expr)

		// Check and skip over the separator
		if b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	return &Node{
		Type:  "array",
		Value: exprs,
	}, nil
}
