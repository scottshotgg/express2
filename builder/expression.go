package builder

import (
	"encoding/json"
	"fmt"

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

	fmt.Println("type:", b.Tokens[b.Index].Type)

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

	// // Check and skip over the separator
	// if b.Tokens[b.Index].Type == token.RParen {
	// 	b.Index++
	// }

	return &Node{
		Type:  "egroup",
		Value: exprs,
	}, nil
}

func (b *Builder) ParseDerefExpression() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.PriOp &&
		b.Tokens[b.Index].Value.String == "*" {
		return nil, b.AppendTokenToError("Could not get deref")
	}

	// Step over the deref
	b.Index++

	ident, err := b.ParseExpression()
	if err != nil {
		return ident, err
	}

	fmt.Println("IDENT: WTFF", ident)

	var kind = "ident"

	switch ident.Type {
	case "type":
		kind = "type"
	}

	return &Node{
		Type: "deref",
		Left: ident,
		Kind: kind,
	}, nil
}

func (b *Builder) ParseRefExpression() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Ampersand &&
		b.Tokens[b.Index].Value.String == "&" {
		return nil, b.AppendTokenToError("Could not get ref")
	}

	// Look ahead and make sure it is an ident; you can't ref anything...
	if b.Tokens[b.Index+1].Type != token.Ident {
		return nil, b.AppendTokenToError("Could not get ident to ref")
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

func (b *Builder) nextOpFunc(tier int) opCallbackFn {
	var opFunc, ok = b.OpFuncMap[tier][b.Tokens[b.Index+1].Type]
	if !ok {
		return opFunc
	}

	// fmt.Printf("b.Tokens[b.Index-1] %+v\n", b.Tokens[b.Index-1].Value)
	// fmt.Printf("b.Tokens[b.Index+1]: %+v\n", b.Tokens[b.Index+1])

	// switch b.Tokens[b.Index-1].Type {
	// case token.RParen:
	// 	if b.Tokens[b.Index+1].Type == token.LParen {
	// 		b.Index++
	// 		return nil
	// 	}
	// }

	return opFunc
}

func (b *Builder) ParseExpression() (*Node, error) {
	term, err := b.ParseTerm2()
	if err != nil {
		return term, err
	}

	fmt.Println("term:", term)

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {
		var opFunc = b.nextOpFunc(2)
		if opFunc == nil {
			break
		}

		fmt.Println("OPFUNC2", b.Tokens[b.Index+1])

		// Step over the factor
		b.Index++

		term, err = opFunc(term)
		if err != nil {
			return term, err
		}

		fmt.Println("term1:", term)
	}

	return term, nil
}

func (b *Builder) ParseTerm2() (*Node, error) {
	factor, err := b.ParseTerm1()
	if err != nil {
		return factor, err
	}

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {
		var opFunc = b.nextOpFunc(1)
		if opFunc == nil {
			break
		}

		fmt.Println("OPFUNC1", b.Tokens[b.Index+1])

		// Step over the factor
		b.Index++

		factor, err = opFunc(factor)
		if err != nil {
			// if err == ErrOutOfTokens {
			// 	return factor,
			// }

			return factor, err
		}

		fmt.Println("term2:", factor)
	}

	return factor, nil
}

func (b *Builder) ParseTerm1() (*Node, error) {
	factor, err := b.ParseFactor()
	if err != nil {
		return factor, err
	}

	// LOOKAHEAD performed to figure out whether the expression is done
	for b.Index < len(b.Tokens)-1 {
		var opFunc = b.nextOpFunc(0)
		if opFunc == nil {
			break
		}

		fmt.Println("OPFUNC0", b.Tokens[b.Index+1])

		// Step over the factor
		b.Index++

		factor, err = opFunc(factor)
		if err != nil {
			// if err == ErrOutOfTokens {
			// 	return factor,
			// }

			return factor, err
		}

		fmt.Println("factor1:", factor)
	}

	return factor, nil
}

// This should go in the function map on the scopetree
var cFuncs = map[string]bool{
	"Println": true,
	"printf":  true,
	"sleep":   true,
	"msleep":  true,
	"now":     true,
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

	case token.Type:
		var node, err = b.ParseType(nil)
		if err != nil {
			return nil, err
		}

		var blob, _ = json.Marshal(node)
		fmt.Println("blob:", string(blob))

		return node, nil
		// return &Node{
		// 	Type:  "type",
		// 	Value: b.Tokens[b.Index].Value.String,
		// }, nil

	// Variable identifier
	case token.Ident:
		var typeOf = "ident"

		var value = b.Tokens[b.Index].Value.String

		// Check the scope map for the variable, if we already have a variable declared then use that
		var n = b.ScopeTree.GetType(value)
		if n != nil {
			typeOf = "type"

			var next = b.Tokens[b.Index+1]
			if next.Type == token.LBrace {
				b.Index++
				nn, err := b.ParseExpression()
				if err != nil {
					return nil, err
				}

				var blob, _ = json.Marshal(n)
				fmt.Println("nblob:", string(blob))

				return &Node{
					Type:  "literal",
					Kind:  n.Kind,
					Value: value,
					Right: nn,
				}, nil
			}
		}

		// Check the scope map for the variable, if we already have a variable declared then use that
		var nv = b.ScopeTree.Get(value)
		if nv != nil {
			if nv.Type == "program" {
				return &Node{
					Type:  "package",
					Value: value,
				}, nil
			}
		}

		return &Node{
			Type:  typeOf,
			Value: value,
		}, nil

	// Deref operator
	case token.PriOp:
		return b.ParseDerefExpression()

	// Ref operator
	case token.Ampersand:
		return b.ParseRefExpression()

	// Nested expression
	case token.LParen:
		var n, err = b.ParseGroupOfExpressions()
		if err != nil {
			return nil, err
		}

		if b.Tokens[b.Index].Type == token.RParen {
			b.Index++
		}

		return n, nil

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

	return nil, b.AppendTokenToError("Could not parse expression from token")
}

func (b *Builder) ParseNestedExpression() (*Node, error) {
	// Check ourselves
	if b.Tokens[b.Index].Type != token.LParen {
		return nil, b.AppendTokenToError("Could not get nested expression")
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
		return nil, b.AppendTokenToError("No right paren found at end of nested expression")
	}

	// Skip over the right paren
	b.Index++

	return expr, nil
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

	// // Step over the right bracket token
	// b.Index++

	return &Node{
		Type:  "array",
		Value: exprs,
	}, nil
}
