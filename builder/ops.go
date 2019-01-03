package builder

import (
	"github.com/scottshotgg/express-token"
)

func (b *Builder) ParseSet(n *Node) (*Node, error) {
	// This will be encountered when we have:
	// <expr> `:` <expr>

	// Step over the set token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the Expression
	b.Index++

	return &Node{
		Type:  "kv",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseBinOp(n *Node) (*Node, error) {
	var op = b.Tokens[b.Index].Value.String

	// Step over the operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "binop",
		Value: op,
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseLessThanExpression(n *Node) (*Node, error) {
	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: "<",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseGreaterThanExpression(n *Node) (*Node, error) {
	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: ">",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseIncrement(n *Node) (*Node, error) {
	return &Node{
		Type: "inc",
		Left: n,
	}, nil
}

func (b *Builder) ParseCall(n *Node) (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.LParen {
		return b.AppendTokenToError("Could not get left paren")
	}

	// We are not allowing for named arguments right now
	args, err := b.ParseGroupOfExpressions()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "call",
		Value: n,
		Metadata: map[string]interface{}{
			"args": args,
		},
	}, nil
}

func (b *Builder) ParseIndexExpression(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
	}

	if b.Tokens[b.Index].Type != token.LBracket {
		return b.AppendTokenToError("Could not get left bracket")
	}

	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// Step over the expression
	b.Index++

	return &Node{
		Type: "index",
		// Value: n,
		Left:  n,
		Right: expr,
	}, nil
}

func (b *Builder) ParseSelection(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
	}

	if b.Tokens[b.Index].Type != token.Accessor {
		return b.AppendTokenToError("Could not get selection operator")
	}

	// Step over the accessor
	b.Index++

	expr, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	// b.Index++

	return &Node{
		Type: "selection",
		// Value: n,
		Left:  n,
		Right: expr,
	}, nil
}
