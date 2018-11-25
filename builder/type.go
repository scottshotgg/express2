package builder

import (
	"github.com/pkg/errors"
	"github.com/scottshotgg/express-token"
)

func (b *Builder) ParseType() (*Node, error) {
	// Check ourselves ...
	if b.Tokens[b.Index].Type != token.Type {
		return b.AppendTokenToError("Could not get type")
	}

	// TODO: we would need to implement something like this
	// TODO: this is where we would also do pointers, need to do function types, etc
	// if typeOf == "map" {

	// }

	typeOf := b.Tokens[b.Index].Value.String

	if b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index+1].Type == token.LBracket {
		return b.ParseArrayType(typeOf)
	}

	// Increment over the type
	b.Index++

	return &Node{
		Type:  "type",
		Value: typeOf,
	}, nil
}

func (b *Builder) ParseArrayType(typeOf string) (*Node, error) {
	var dim []*Index

	// Look ahead at the next token here
	for b.Index < len(b.Tokens)-1 &&
		b.Tokens[b.Index+1].Type == token.LBracket {
		// Increment over the type token
		b.Index++

		expr, err := b.ParseArrayExpression()
		if err != nil {
			return nil, err
		}

		if expr.Value == nil {
			return nil, errors.Errorf("Array parse value was nil; %+v", expr)
		}

		nodesAssert, ok := expr.Value.([]*Node)
		if !ok {
			return nil, errors.Errorf("Invalid assertion; %+v", expr)
		}

		var dimValue Index

		switch len(nodesAssert) {

		case 1:
			dimValue.Type = nodesAssert[0].Kind
			dimValue.Value, ok = nodesAssert[0].Value.(int)
			if !ok {
				return nil, errors.Errorf("Could not assert array value to int; %+v", nodesAssert[0].Value)
			}

		case 0:
			dimValue.Type = "none"
			dimValue.Value = -1

		default:
			return nil, ErrMultDimArrInit
		}

		dim = append(dim, &dimValue)
	}

	b.Index++

	return &Node{
		Type:  "type",
		Value: "array",
		Metadata: map[string]interface{}{
			// "type": typeOf,
			"dim": dim,
		},
	}, nil
}