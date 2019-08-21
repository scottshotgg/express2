package builder

import (
	"fmt"

	"github.com/scottshotgg/express-token"
)

func (b *Builder) ParseGroupOfExpressions() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get group of expressions")
	}

	// Skip over the left paren token
	b.Index++

	var (
		expr  *Node
		exprs []*Node
		err   error
	)

	for b.Tokens[b.Index].Type != token.RParen {
		expr, err = b.ParseExpression()
		if err != nil {
			return expr, err
		}

		// Step over the expression token
		b.Index++

		exprs = append(exprs, expr)

		// Check and skip over the separator
		if b.Tokens[b.Index].Type == token.Separator {
			b.Index++
		}
	}

	// // Step over the right paren token
	// b.Index++

	return &Node{
		Type:  "egroup",
		Value: exprs,
	}, nil
}

func (b *Builder) ParseDerefExpression() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.PriOp &&
		b.Tokens[b.Index].Value.String == "*" {
		return b.AppendTokenToError("Could not get deref")
	}

	// Look ahead and make sure it is an ident;you can't deref just anything...
	if b.Tokens[b.Index+1].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident to deref")
	}

	// Step over the deref
	b.Index++

	ident, err := b.ParseExpression()
	if err != nil {
		return ident, err
	}

	return &Node{
		Type: "deref",
		Left: ident,
	}, nil
}

func (b *Builder) ParseRefExpression() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Ampersand &&
		b.Tokens[b.Index].Value.String == "&" {
		return b.AppendTokenToError("Could not get ref")
	}

	// Look ahead and make sure it is an ident; you can't ref anything...
	if b.Tokens[b.Index+1].Type != token.Ident {
		return b.AppendTokenToError("Could not get ident to ref")
	}

	// Step over the deref
	b.Index++

	// Will probably have to change this to just parse the ident instead
	// so we don't have problems with operator precedence
	ident, err := b.ParseExpression()
	if err != nil {
		return ident, err
	}

	return &Node{
		Type: "ref",
		Left: ident,
	}, nil
}

func (b *Builder) ParseExpression() (*Node, error) {
	term, err := b.ParseTerm()
	if err != nil {
		return term, err
	}

	var (
		ok     bool
		opFunc opCallbackFn
	)

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {
		// Look for a tier2 operator in the func map
		opFunc, ok = b.OpFuncMap[1][b.Tokens[b.Index+1].Type]
		if !ok {
			break
		}

		// Step over the factor
		b.Index++

		term, err = opFunc(term)
		if err != nil {
			return term, err
		}
	}

	return term, nil
}

func (b *Builder) ParseTerm() (*Node, error) {
	factor, err := b.ParseFactor()
	if err != nil {
		return factor, err
	}

	var (
		ok     bool
		opFunc opCallbackFn
	)

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {

		// Look for a tier1 operator in the func map
		opFunc, ok = b.OpFuncMap[0][b.Tokens[b.Index+1].Type]
		if !ok {
			break
		}
		fmt.Println("OPFUNC", b.Tokens[b.Index+1])

		// Step over the factor
		b.Index++

		factor, err = opFunc(factor)
		if err != nil {
			// if err == ErrOutOfTokens {
			// 	return factor,
			// }

			return factor, err
		}
	}

	return factor, nil
}

func (b *Builder) ParseFactor() (*Node, error) {
	// Here we will switch on the type and determine whether we have:
	// - literal
	// - ident
	// - call
	// - index operation
	// - selection operation
	// - block
	// - array
	// - nil

	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
		// return nil, nil
	}

	switch b.Tokens[b.Index].Type {
	// Any literal value
	case token.Literal:
		return &Node{
			Type:  "literal",
			Kind:  b.Tokens[b.Index].Value.Type,
			Value: b.Tokens[b.Index].Value.True,
		}, nil

	// Variable identifier
	case token.Ident:
		// // Check the scope map for the variable, if we already have a variable declared then use that
		// if node != nil {
		// 	// TODO: might need to fix this
		// 	return node, nil
		// }

		return &Node{
			Type:  "ident",
			Value: b.Tokens[b.Index].Value.String,
		}, nil

	// Deref operator
	case token.PriOp:
		return b.ParseDerefExpression()

	// Ref operator
	case token.Ampersand:
		return b.ParseRefExpression()

	// Nested expression
	case token.LParen:
		return b.ParseNestedExpression()

	// Array expression
	case token.LBracket:
		return b.ParseArrayExpression()

	// Named block
	case token.LBrace:
		var a, c = b.ParseBlockStatement()
		// If this is an expression, then whatever called ParseExpression
		// is going to increment the index again ...
		b.Index--
		return a, c
	}

	return b.AppendTokenToError("Could not parse expression from token")
}

func (b *Builder) ParseNestedExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get nested expression")
	}

	// Skip over the left paren
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return expr, err
	}

	// Skip over the expression
	b.Index++

	if b.Tokens[b.Index].Type != token.RParen {
		return b.AppendTokenToError("No right paren found at end of nested expression")
	}

	// Skip over the right paren
	b.Index++

	return expr, nil
}

func (b *Builder) ParseArrayExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LBracket {
		return b.AppendTokenToError("Could not get array expression")
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

	// // Step over the right bracket token
	// b.Index++

	return &Node{
		Type:  "array",
		Value: exprs,
	}, nil
}
