package ast

import (
	"errors"
	"fmt"

	"github.com/scottshotgg/express-token"
)

type IterableType int

const (
	_ IterableType = iota
	Key
	Value
	Both
)

// Iterable is an abstract type in Express that represents the
// ability for an expression (specifically a literal) to be iterated over
type Iter interface {
	Node

	// This is just something to force the interface
	expressionNode()

	Next() *Literal
	Prev() *Literal

	Fields() []*Literal

	// TODO: this should have a Type function
	// TODO: this should have a Length function
}

type Iterable struct {
	Token    token.Token
	Variable *Ident
	Over     Expression
	Type     IterableType
}

func (i *Iterable) expressionNode() {}

// TokenLiteral returns the literal value of the token
func (i *Iterable) TokenLiteral() token.Token { return i.Token }

func (i *Iterable) Kind() NodeType { return IterableNode }

func (i *Iterable) String() string {
	// FIXME: just doing this to get it to compile
	return fmt.Sprintf("%+v", *i)
}

func NewIterable(t, i token.Token, v *Ident, o Expression) (*Iterable, error) {
	switch i.Value.String {
	case "in":
		return NewKeyIterable(t, v, o)

	case "of":
		return NewValueIterable(t, v, o)

	case "over":
		return NewBothIterable(t, v, o)

	default:
		return nil, errors.New("Could not deduce iterable type")
	}
}

func NewKeyIterable(t token.Token, v *Ident, o Expression) (*Iterable, error) {
	return &Iterable{
		Token:    t,
		Variable: v,
		Over:     o,
		Type:     Key,
	}, nil
}

func NewValueIterable(t token.Token, v *Ident, o Expression) (*Iterable, error) {
	return &Iterable{
		Token:    t,
		Variable: v,
		Over:     o,
		Type:     Value,
	}, nil
}

func NewBothIterable(t token.Token, v *Ident, o Expression) (*Iterable, error) {
	return &Iterable{
		Token:    t,
		Variable: v,
		Over:     o,
		Type:     Both,
	}, nil
}
