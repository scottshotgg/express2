package builder

import (
	"encoding/json"
	"fmt"

	token "github.com/scottshotgg/express-token"
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

	var blob, _ = json.Marshal(right)
	fmt.Println("rightblob:", string(blob))

	blob, _ = json.Marshal(n)
	fmt.Println("nnnnnnnnnnnblob:", string(blob))

	// Step over the expression
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

	right, err := b.ParseTerm1()
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

func (b *Builder) ParseLessOrEqualThanExpression(n *Node) (*Node, error) {
	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: "<=",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseGreaterOrEqualThanExpression(n *Node) (*Node, error) {
	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	return &Node{
		Type:  "comp",
		Value: ">=",
		Left:  n,
		Right: right,
	}, nil
}

func (b *Builder) ParseEqualityExpression(n *Node) (*Node, error) {
	fmt.Println("node:", n)

	// Step over the conditional operator token
	b.Index++

	right, err := b.ParseExpression()
	if err != nil {
		return nil, err
	}

	fmt.Println("right:", right)

	var nn = &Node{
		Type:  "comp",
		Value: "==",
		Left:  n,
		Right: right,
	}

	fmt.Println("nn:", *nn)

	return nn, nil
}

func (b *Builder) ParseIncrement(n *Node) (*Node, error) {
	return &Node{
		Type: "inc",
		Left: n,
	}, nil
}

func (b *Builder) ParseCall(n *Node) (*Node, error) {
	// We are not allowing for named arguments right now
	args, err := b.ParseGroupOfExpressions()
	if err != nil {
		return nil, err
	}

	// if b.Index < len(b.Tokens) && b.Tokens[b.Index].Type == token.RParen {
	// 	b.Index++
	// }

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
		return nil, b.AppendTokenToError("Could not get left bracket")
	}

	// b.Index++

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

// func (b *Builder) ParseBlockExpression(n *Node) (*Node, error) {
// 	if b.Index > len(b.Tokens)-1 {
// 		return nil, ErrOutOfTokens
// 	}

// 	if b.Tokens[b.Index].Type != token.LBrace {
// 		return nil, b.AppendTokenToError("Could not get left bracket")
// 	}

// 	b.Index++

// 	expr, err := b.ParseExpression()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Step over the expression
// 	b.Index++

// 	return &Node{
// 		Type: "index",
// 		// Value: n,
// 		Left:  n,
// 		Right: expr,
// 	}, nil
// }

func (b *Builder) ParseSelection(n *Node) (*Node, error) {
	if b.Index > len(b.Tokens)-1 {
		return nil, ErrOutOfTokens
	}

	if b.Tokens[b.Index].Type != token.Accessor {
		return nil, b.AppendTokenToError("Could not get selection operator")
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
